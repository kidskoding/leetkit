package runner

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// testPython prefers the repository's uv-managed .venv, then PATH, and skips
// when no interpreter with pytest is available.
func testPython(t *testing.T) string {
	t.Helper()
	venv, err := filepath.Abs(filepath.Join("..", "..", ".venv", "bin", "python"))
	if err != nil {
		t.Fatal(err)
	}
	py, err := FindPython(Candidate{Path: venv, Source: "repo venv"})
	if err != nil {
		t.Skipf("no Python found: %v", err)
	}
	if _, err := PytestVersion(context.Background(), py.Path); err != nil {
		t.Skipf("%s: %v (run `uv sync` or `nix develop`)", py.Path, err)
	}
	return py.Path
}

const tests = `from starter import add

def test_small():
    assert add(1, 2) == 3

def test_negative():
    assert add(-1, -1) == -2
`

func run(t *testing.T, code string, timeout time.Duration) *Result {
	t.Helper()
	r := Runner{Python: testPython(t), Timeout: timeout}
	res, err := r.Run(context.Background(), Job{
		Tests:    []byte(tests),
		Code:     []byte(code),
		CodePath: "/work/add/starter.py",
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	return res
}

func TestRunPassing(t *testing.T) {
	res := run(t, "def add(a, b):\n    return a + b\n", 30*time.Second)
	if res.Status != StatusPassed || res.Count(OutcomePassed) != 2 {
		t.Fatalf("status %v, cases %+v\n%s", res.Status, res.Cases, res.Output)
	}
}

func TestRunFailing(t *testing.T) {
	res := run(t, "def add(a, b):\n    return a - b\n", 30*time.Second)
	if res.Status != StatusFailed || res.Count(OutcomeFailed) != 2 {
		t.Fatalf("status %v, cases %+v", res.Status, res.Cases)
	}
	if msg := res.Cases[0].Message; !strings.Contains(msg, "assert -1 == 3") {
		t.Errorf("message %q lacks the assertion", msg)
	}
}

func TestRunSyntaxError(t *testing.T) {
	res := run(t, "def add(a, b)\n    return a + b\n", 30*time.Second)
	if res.Status != StatusError {
		t.Fatalf("status %v, want error", res.Status)
	}
	msg := res.Cases[0].Message
	if !strings.Contains(msg, "SyntaxError") || !strings.Contains(msg, "/work/add/starter.py") {
		t.Errorf("message should name the error and the learner's file, got:\n%s", msg)
	}
	if strings.Contains(res.Output, "leetkit-run-") {
		t.Errorf("output leaks the scratch directory:\n%s", res.Output)
	}
}

func TestRunRuntimeException(t *testing.T) {
	res := run(t, "def add(a, b):\n    raise ValueError('boom')\n", 30*time.Second)
	if res.Status != StatusFailed || !strings.Contains(res.Cases[0].Message, "ValueError: boom") {
		t.Fatalf("status %v, cases %+v", res.Status, res.Cases)
	}
}

func TestRunTimeout(t *testing.T) {
	start := time.Now()
	res := run(t, "def add(a, b):\n    while True:\n        pass\n", time.Second)
	if res.Status != StatusTimedOut {
		t.Fatalf("status %v, want timed out", res.Status)
	}
	if elapsed := time.Since(start); elapsed > 10*time.Second {
		t.Errorf("timeout took %s to take effect", elapsed)
	}
}

func TestRunCanceled(t *testing.T) {
	r := Runner{Python: testPython(t), Timeout: time.Minute}
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(500*time.Millisecond, cancel)
	res, err := r.Run(ctx, Job{Tests: []byte(tests), Code: []byte("def add(a, b):\n    while True: pass\n")})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != StatusCanceled {
		t.Fatalf("status %v, want canceled", res.Status)
	}
}

func TestRunMissingPytest(t *testing.T) {
	py, err := FindPython()
	if err != nil {
		t.Skip(err)
	}
	if _, err := PytestVersion(context.Background(), py.Path); !errors.Is(err, ErrPytestMissing) {
		t.Skipf("%s has pytest; cannot exercise the missing-pytest path", py.Path)
	}
	_, err = Runner{Python: py.Path, Timeout: 30 * time.Second}.Run(context.Background(), Job{Tests: []byte(tests)})
	if !errors.Is(err, ErrPytestMissing) {
		t.Errorf("err = %v, want ErrPytestMissing", err)
	}
}

func TestRunMissingInterpreter(t *testing.T) {
	_, err := Runner{Python: filepath.Join(t.TempDir(), "nope")}.Run(context.Background(), Job{})
	if err == nil || !strings.Contains(err.Error(), "start") {
		t.Errorf("err = %v, want start failure", err)
	}
}

func TestFindPython(t *testing.T) {
	dir := t.TempDir()
	fake := filepath.Join(dir, "mypython")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)

	got, err := FindPython(Candidate{Path: "", Source: "config"}, Candidate{Path: fake, Source: "venv"})
	if err != nil || got.Path != fake || got.Source != "venv" {
		t.Errorf("got %+v, %v", got, err)
	}

	_, err = FindPython(Candidate{Path: filepath.Join(dir, "missing"), Source: "LEETKIT_PYTHON", Explicit: true})
	if err == nil || !strings.Contains(err.Error(), "LEETKIT_PYTHON") {
		t.Errorf("explicit missing interpreter: err = %v", err)
	}

	if _, err := FindPython(Candidate{Path: filepath.Join(dir, "missing"), Source: "venv"}); !errors.Is(err, ErrNoPython) {
		t.Errorf("no interpreter on PATH: err = %v, want ErrNoPython", err)
	}
}

func TestCappedBuffer(t *testing.T) {
	b := &cappedBuffer{limit: 4}
	n, err := b.Write([]byte("abcdef"))
	if n != 6 || err != nil {
		t.Fatalf("Write = %d, %v", n, err)
	}
	if got := b.String(); got != "abcd\n[output truncated]\n" {
		t.Errorf("got %q", got)
	}
}
