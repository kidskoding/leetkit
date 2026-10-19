package cli

import (
	"context"
	"flag"
	"fmt"

	"charm.land/lipgloss/v2"
)

func (a *app) list(_ context.Context, args []string) error {
	fset := flag.NewFlagSet("list", flag.ContinueOnError)
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
	if len(cat.All()) == 0 {
		fmt.Fprintln(a.Stdout, "No exercises found.")
		return nil
	}
	idWidth, titleWidth, levelWidth := 0, 0, 0
	for _, ex := range cat.All() {
		idWidth = max(idWidth, lipgloss.Width(ex.ID))
		titleWidth = max(titleWidth, lipgloss.Width(ex.Title))
		levelWidth = max(levelWidth, lipgloss.Width(ex.Difficulty))
	}
	for i, m := range cat.Modules {
		if i > 0 {
			fmt.Fprintln(a.Stdout)
		}
		fmt.Fprintln(a.Stdout, a.th.Heading.Render(m.Title))
		for _, ex := range m.Exercises {
			fmt.Fprintf(a.Stdout, "  %s  %s  %s  %s\n",
				a.th.ID.Render(pad(ex.ID, idWidth)),
				a.th.Text.Render(pad(ex.Title, titleWidth)),
				a.th.Faint.Render(pad(ex.Difficulty, levelWidth)),
				a.th.Faint.Render(ex.TypeLabel()))
		}
	}
	return nil
}
