package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kidskoding/leetkit/exercises"
	"github.com/kidskoding/leetkit/internal/config"
	"github.com/kidskoding/leetkit/internal/runner"
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

func writeWork(t *testing.T, env Env, id string, code []byte) {
	t.Helper()
	path := filepath.Join(env.Paths.WorkspaceDir(), id, "starter.py")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, code, 0o644); err != nil {
		t.Fatal(err)
	}
}

func solution(t *testing.T) []byte {
	t.Helper()
	data, err := exercises.FS.ReadFile("iteration/enumerate-01/solution.py")
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestCheckSolutionPasses(t *testing.T) {
	usePytestPython(t)
	env := tempEnv(t)
	writeWork(t, env, "enumerate-01", solution(t))
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
		{[]string{"check", "a", "b"}, exitUsage, "expected at most one exercise id"},
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

func TestCheckAll(t *testing.T) {
	usePytestPython(t)
	tests := []struct {
		name  string
		work  []byte
		code  int
		wants []string
	}{
		{"not started", nil, exitOK, []string{"○ enumerate-01", "not started", "0/1 passing"}},
		{"solved", solution(t), exitOK, []string{"✓ enumerate-01", "6/6 passed", "1/1 passing"}},
		{"failing", []byte("def find_index(nums, target):\n    return 0\n"), exitError,
			[]string{"✗ enumerate-01", "1 failing", "leetkit check <id>"}},
		{"syntax error", []byte("def find_index(\n"), exitError, []string{"could not run"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := tempEnv(t)
			if tt.work != nil {
				writeWork(t, env, "enumerate-01", tt.work)
			}
			code, out, errOut := run(t, env, "check")
			if code != tt.code {
				t.Errorf("exit %d, want %d; stderr %s", code, tt.code, errOut)
			}
			for _, want := range tt.wants {
				if !strings.Contains(out, want) {
					t.Errorf("output missing %q:\n%s", want, out)
				}
			}
		})
	}
}

func TestCheckAllDoesNotCreateWorkingCopies(t *testing.T) {
	usePytestPython(t)
	env := tempEnv(t)
	run(t, env, "check")
	if _, err := os.Stat(env.Paths.WorkspaceDir()); !os.IsNotExist(err) {
		t.Errorf("check created %s", env.Paths.WorkspaceDir())
	}
}

func TestPackagedDefaultPython(t *testing.T) {
	t.Setenv("LEETKIT_PYTHON", "")
	fake := filepath.Join(t.TempDir(), "python3")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	a := newApp(Env{Paths: tempEnv(t).Paths, DefaultPython: fake})
	py, err := a.findPython(a.cfg)
	if err != nil || py.Path != fake || py.Source != "package default" {
		t.Errorf("got %+v, %v", py, err)
	}

	t.Setenv("LEETKIT_PYTHON", fake+"-override")
	if _, err := a.findPython(a.cfg); err == nil {
		t.Error("$LEETKIT_PYTHON should take precedence over the packaged default")
	}
}
