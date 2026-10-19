package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anirudhkonidala/leetkit/exercises"
	"github.com/anirudhkonidala/leetkit/internal/config"
	"github.com/anirudhkonidala/leetkit/internal/runner"
)

// usePytestPython points LEETKIT_PYTHON at an interpreter with pytest, or skips.
func usePytestPython(t *testing.T) {
	t.Helper()
	venv, err := filepath.Abs(filepath.Join("..", "..", ".venv", "bin", "python"))
	if err != nil {
		t.Fatal(err)
	}
	py, err := runner.FindPython(runner.Candidate{Path: venv})
	if err == nil {
		_, err = runner.PytestVersion(context.Background(), py.Path)
	}
	if err != nil {
		t.Skipf("no Python with pytest: %v", err)
	}
	t.Setenv("LEETKIT_PYTHON", py.Path)
}

func tempEnv(t *testing.T) Env {
	dir := t.TempDir()
	return Env{
		Exercises: exercises.FS,
		Paths:     config.Paths{ConfigFile: filepath.Join(dir, "config.toml"), DataDir: filepath.Join(dir, "data")},
	}
}

func TestCheckStarterFails(t *testing.T) {
	usePytestPython(t)
	env := tempEnv(t)
	code, out, errOut := run(t, env, "check", "enumerate-01")
	if code != exitError {
		t.Fatalf("exit %d, stderr %s", code, errOut)
	}
	for _, want := range []string{"✗ 0 passed, 6 failed", "✗ test_target_present", "assert None == 1", "Your code: "} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
	if errOut != "" {
		t.Errorf("unexpected stderr: %s", errOut)
	}
}

func TestCheckSolutionPasses(t *testing.T) {
	usePytestPython(t)
	env := tempEnv(t)
	solution, err := exercises.FS.ReadFile("iteration/enumerate-01/solution.py")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(env.Paths.WorkspaceDir(), "enumerate-01", "starter.py")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, solution, 0o644); err != nil {
		t.Fatal(err)
	}
	code, out, errOut := run(t, env, "check", "enumerate-01")
	if code != exitOK || !strings.Contains(out, "✓ All 6 tests passed") {
		t.Errorf("exit %d\nstdout:\n%s\nstderr:\n%s", code, out, errOut)
	}
}

func TestCheckArguments(t *testing.T) {
	tests := []struct {
		args []string
		code int
		want string
	}{
		{[]string{"check"}, exitUsage, "expected one exercise id"},
		{[]string{"check", "a", "b"}, exitUsage, "expected one exercise id"},
		{[]string{"check", "nope-01"}, exitError, `no exercise "nope-01"`},
	}
	for _, tt := range tests {
		code, _, errOut := run(t, tempEnv(t), tt.args...)
		if code != tt.code || !strings.Contains(errOut, tt.want) {
			t.Errorf("%v: exit %d, stderr %q", tt.args, code, errOut)
		}
	}
}

func TestCheckBadPythonSetting(t *testing.T) {
	t.Setenv("LEETKIT_PYTHON", filepath.Join(t.TempDir(), "missing-python"))
	code, _, errOut := run(t, tempEnv(t), "check", "enumerate-01")
	if code != exitError || !strings.Contains(errOut, "$LEETKIT_PYTHON") || !strings.Contains(errOut, "leetkit doctor") {
		t.Errorf("exit %d, stderr %q", code, errOut)
	}
}

func TestDoctor(t *testing.T) {
	usePytestPython(t)
	t.Setenv("EDITOR", "sh")
	code, out, _ := run(t, tempEnv(t), "doctor")
	if code != exitOK || !strings.Contains(out, "All checks passed") || !strings.Contains(out, "✓ pytest") {
		t.Errorf("exit %d\n%s", code, out)
	}
}

func TestDoctorReportsProblems(t *testing.T) {
	t.Setenv("LEETKIT_PYTHON", filepath.Join(t.TempDir(), "missing-python"))
	env := tempEnv(t)
	if err := os.WriteFile(env.Paths.ConfigFile, []byte("timeout = 5\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, out, _ := run(t, env, "doctor")
	if code != exitError {
		t.Errorf("exit %d, want %d", code, exitError)
	}
	for _, want := range []string{"✗ Config", "duration string", "✗ Python", "$LEETKIT_PYTHON is set to", "Correct or unset"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
}
