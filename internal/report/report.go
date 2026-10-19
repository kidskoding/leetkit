// Package report renders test results for both the CLI and the dashboard.
package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/anirudhkonidala/leetkit/internal/runner"
	"github.com/anirudhkonidala/leetkit/internal/theme"
)

// Result renders a summary line followed by each failing or erroring test.
// Verbose appends pytest's full output.
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
	for _, c := range res.Cases {
		if c.Outcome != runner.OutcomeFailed && c.Outcome != runner.OutcomeError {
			continue
		}
		fmt.Fprintf(&b, "\n  %s\n%s\n", th.Bad.Render("✗ "+c.Name), Indent(c.Message, "    "))
	}
	if verbose && res.Output != "" {
		fmt.Fprintf(&b, "\n%s\n%s", th.Faint.Render("─── pytest output ───"), res.Output)
	}
	return b.String()
}

func Indent(s, prefix string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i, l := range lines {
		lines[i] = prefix + l
	}
	return strings.Join(lines, "\n")
}
