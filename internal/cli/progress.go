package cli

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/kidskoding/leetkit/internal/theme"
)

const progressBarWidth = 24

func (a *app) progress(_ context.Context, args []string) error {
	fset := flag.NewFlagSet("progress", flag.ContinueOnError)
	if err := parseFlags(fset, args); err != nil {
		return err
	}
	if fset.NArg() > 0 {
		return usagef("unexpected argument %q", fset.Arg(0))
	}
	cat, err := a.catalog()
	if err != nil {
		return err
	}
	store, err := a.openStore()
	if err != nil {
		return err
	}

	labelWidth := lipgloss.Width("Overall")
	for _, m := range cat.Modules {
		labelWidth = max(labelWidth, lipgloss.Width(m.Title))
	}
	row := func(label string, done, total int, style lipgloss.Style) string {
		pct := 0
		if total > 0 {
			pct = done * 100 / total
		}
		return fmt.Sprintf("%s  %s  %s %s", style.Render(pad(label, labelWidth)), bar(a.th, done, total),
			a.th.Text.Render(fmt.Sprintf("%d/%d", done, total)), a.th.Faint.Render(fmt.Sprintf("%d%%", pct)))
	}

	var all []string
	var lines []string
	attempts := 0
	var last time.Time
	for _, m := range cat.Modules {
		ids := make([]string, len(m.Exercises))
		for i, ex := range m.Exercises {
			ids[i] = ex.ID
			if e, ok := store.Entry(ex.ID); ok {
				attempts += e.Attempts
				if e.LastAttempt.After(last) {
					last = e.LastAttempt
				}
			}
		}
		all = append(all, ids...)
		lines = append(lines, row(m.Title, store.CountCompleted(ids), len(ids), a.th.Text))
	}
	done := store.CountCompleted(all)

	fmt.Fprintln(a.Stdout, row("Overall", done, len(all), a.th.Heading))
	fmt.Fprintln(a.Stdout)
	for _, l := range lines {
		fmt.Fprintln(a.Stdout, l)
	}
	summary := fmt.Sprintf("%d solved · %d test runs", done, attempts)
	if !last.IsZero() {
		summary += " · last activity " + last.Local().Format("2006-01-02 15:04")
	}
	fmt.Fprintln(a.Stdout, "\n"+a.th.Faint.Render(summary))
	return nil
}

func bar(th theme.Theme, done, total int) string {
	filled := 0
	if total > 0 {
		filled = done * progressBarWidth / total
	}
	return th.Good.Render(strings.Repeat("█", filled)) + th.Faint.Render(strings.Repeat("░", progressBarWidth-filled))
}
