package report

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"

	"github.com/kidskoding/leetkit/internal/runner"
	"github.com/kidskoding/leetkit/internal/theme"
)

func render(t *testing.T, res *runner.Result, verbose bool) string {
	t.Helper()
	th, err := theme.Load("")
	if err != nil {
		t.Fatal(err)
	}
	return ansi.Strip(Result(th, res, verbose))
}

func TestResult(t *testing.T) {
	failing := &runner.Result{
		Status:   runner.StatusFailed,
		Duration: 234 * time.Millisecond,
		Output:   "raw pytest output",
		Cases: []runner.Case{
			{Name: "test_ok", Outcome: runner.OutcomePassed},
			{Name: "test_bad", Outcome: runner.OutcomeFailed, Message: "assert 1 == 2\n +  where 1 = f()"},
		},
	}
	tests := []struct {
		name    string
		res     *runner.Result
		verbose bool
		want    []string
		absent  []string
	}{
		{"passed", &runner.Result{Status: runner.StatusPassed, Cases: []runner.Case{{Outcome: runner.OutcomePassed}}},
			false, []string{"✓ All 1 tests passed"}, []string{"✗"}},
		{"failed", failing, false,
			[]string{"✗ 1 passed, 1 failed 230ms", "  ✗ test_bad\n    assert 1 == 2\n     +  where 1 = f()"},
			[]string{"test_ok", "raw pytest output"}},
		{"verbose", failing, true, []string{"─── pytest output ───\nraw pytest output"}, nil},
		{"error", &runner.Result{Status: runner.StatusError, Problem: "Your code could not be tested.",
			Cases: []runner.Case{{Name: "test_exercise", Outcome: runner.OutcomeError, Message: "SyntaxError: expected ':'"}}},
			false, []string{"✗ Your code could not be tested.", "SyntaxError: expected ':'"}, nil},
		{"timeout", &runner.Result{Status: runner.StatusTimedOut, Problem: "Tests did not finish within 10s."},
			false, []string{"! Tests did not finish within 10s."}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := render(t, tt.res, tt.verbose)
			for _, w := range tt.want {
				if !strings.Contains(out, w) {
					t.Errorf("missing %q in:\n%s", w, out)
				}
			}
			for _, a := range tt.absent {
				if strings.Contains(out, a) {
					t.Errorf("unexpected %q in:\n%s", a, out)
				}
			}
		})
	}
}
