// Package report renders test results for both the CLI and the dashboard.
package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/kidskoding/leetkit/internal/runner"
	"github.com/kidskoding/leetkit/internal/theme"
)

// Result renders a summary line and every test in run order: passes as a
// checklist, failures and errors with their message underneath. Verbose
// appends pytest's full output.
func Result(th theme.Theme, res *runner.Result, verbose bool) string {
	var b strings.Builder
	passed, failed := res.Count(runner.OutcomePassed), res.Count(runner.OutcomeFailed)
	elapsed := th.Faint.Render(res.Duration.Round(10 * time.Millisecond).String())
	switch res.Status {
	case runner.StatusPassed:
		fmt.Fprintf(&b, "%s %s\n", th.Pass.Render(fmt.Sprintf("✓ All %d tests passed", passed)), elapsed)
	case runner.StatusFailed:
		fmt.Fprintf(&b, "%s %s, %s %s\n", th.Fail.Render("✗"),
			th.Good.Render(fmt.Sprintf("%d passed", passed)),
			th.Bad.Render(fmt.Sprintf("%d failed", failed)), elapsed)
	case runner.StatusError:
		fmt.Fprintln(&b, th.Fail.Render("✗ "+res.Problem))
	default:
		fmt.Fprintln(&b, th.Warn.Render("! "+res.Problem))
	}
	if len(res.Cases) > 0 {
		b.WriteString("\n")
	}
	for i, c := range res.Cases {
		switch {
		case c.Outcome == runner.OutcomePassed:
			fmt.Fprintf(&b, "  %s\n", th.Good.Render("✓ "+c.Name))
		case c.Outcome == runner.OutcomeSkipped:
			fmt.Fprintf(&b, "  %s\n", th.Faint.Render("○ "+c.Name+" (skipped)"))
		case c.Collection:
			// Not a real test: show the error itself, usually a SyntaxError.
			fmt.Fprintf(&b, "%s\n", Indent(c.Message, "  "))
		default:
			fmt.Fprintf(&b, "  %s\n%s\n", th.Bad.Render("✗ "+c.Name), Indent(c.Message, "      "))
			// A blank line keeps a failure's details apart from the next test.
			if i < len(res.Cases)-1 {
				b.WriteString("\n")
			}
		}
	}
	if verbose && res.Output != "" {
		fmt.Fprintf(&b, "\n%s\n%s", th.Faint.Render("─── pytest output ───"), res.Output)
	}
	return b.String()
}

func Indent(s, prefix string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i, l := range lines {
		if l != "" {
			lines[i] = prefix + l
		}
	}
	return strings.Join(lines, "\n")
}

// Meter renders a fixed-width completion bar such as "████░░░░".
func Meter(th theme.Theme, done, total, width int) string {
	filled := 0
	if total > 0 {
		filled = min(done*width/total, width)
	}
	return th.Good.Render(strings.Repeat("█", filled)) + th.Faint.Render(strings.Repeat("░", width-filled))
}
