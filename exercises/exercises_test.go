package exercises_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kidskoding/leetkit/exercises"
	"github.com/kidskoding/leetkit/internal/exercise"
	"github.com/kidskoding/leetkit/internal/runner"
)

func TestBundledCurriculumLoads(t *testing.T) {
	cat, err := exercise.Load(exercises.FS)
	if err != nil {
		t.Fatalf("bundled curriculum is invalid:\n%v", err)
	}
	if len(cat.All()) == 0 {
		t.Fatal("bundled curriculum has no exercises")
	}
}

// TestCurriculum runs every exercise twice: the reference solution must pass,
// and the untouched starter must fail its tests without a syntax or import
// error. Run one exercise with -run 'TestCurriculum/<id>'.
func TestCurriculum(t *testing.T) {
	cat, err := exercise.Load(exercises.FS)
	if err != nil {
		t.Fatal(err)
	}
	py := pytestPython(t)
	r := runner.Runner{Python: py, Timeout: time.Minute}
	for _, ex := range cat.All() {
		t.Run(ex.ID, func(t *testing.T) {
			t.Parallel()
			tests := read(t, cat, ex, exercise.TestFile)

			res := run(t, r, tests, read(t, cat, ex, exercise.SolutionFile))
			if res.Status != runner.StatusPassed {
				t.Errorf("solution: %s\n%s", res.Status, res.Output)
			}
			if n := len(res.Cases); n < 3 {
				t.Errorf("only %d tests; cover normal and edge cases", n)
			}

			res = run(t, r, tests, read(t, cat, ex, exercise.StarterFile))
			if res.Status != runner.StatusFailed {
				t.Errorf("starter: want failing tests, got %s\n%s", res.Status, res.Output)
			}
		})
	}
}

func read(t *testing.T, cat *exercise.Catalog, ex *exercise.Exercise, name string) []byte {
	t.Helper()
	data, err := cat.ReadFile(ex, name)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func run(t *testing.T, r runner.Runner, tests, code []byte) *runner.Result {
	t.Helper()
	res, err := r.Run(context.Background(), runner.Job{Tests: tests, Code: code, CodePath: "starter.py"})
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func pytestPython(t *testing.T) string {
	t.Helper()
	venv, err := filepath.Abs(filepath.Join("..", ".venv", "bin", "python"))
	if err != nil {
		t.Fatal(err)
	}
	py, err := runner.FindPython(runner.Candidate{Path: venv})
	if err == nil {
		_, err = runner.PytestVersion(context.Background(), py.Path)
	}
	if err != nil {
		t.Skipf("no Python with pytest (run `uv sync` or `nix develop`): %v", err)
	}
	return py.Path
}

// The README advertises the exercise count in a badge and in the curriculum
// section; keep both in step with the catalog.
func TestReadmeExerciseCount(t *testing.T) {
	cat, err := exercise.Load(exercises.FS)
	if err != nil {
		t.Fatal(err)
	}
	readme, err := os.ReadFile(filepath.Join("..", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	n := len(cat.All())
	for _, want := range []string{
		fmt.Sprintf("badge/exercises-%d-", n),
		fmt.Sprintf("%d exercises across", n),
	} {
		if !strings.Contains(string(readme), want) {
			t.Errorf("README.md does not mention %q; update it to %d exercises", want, n)
		}
	}
}
