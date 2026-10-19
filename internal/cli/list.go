package cli

import (
	"context"
	"flag"
	"fmt"

	"charm.land/lipgloss/v2"

	"github.com/anirudhkonidala/leetkit/internal/theme"
)

func runList(_ context.Context, env Env, args []string) error {
	fset := flag.NewFlagSet("list", flag.ContinueOnError)
	if err := parseFlags(fset, args); err != nil {
		return err
	}
	if fset.NArg() > 0 {
		return usagef("unexpected argument %q", fset.Arg(0))
	}
	cat, err := loadCatalog(env)
	if err != nil {
		return err
	}
	if len(cat.All()) == 0 {
		fmt.Fprintln(env.Stdout, "No exercises found.")
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
			fmt.Fprintln(env.Stdout)
		}
		fmt.Fprintln(env.Stdout, theme.Heading.Render(m.Title))
		for _, ex := range m.Exercises {
			fmt.Fprintf(env.Stdout, "  %s  %s  %s  %s\n",
				theme.ID.Render(pad(ex.ID, idWidth)),
				theme.Text.Render(pad(ex.Title, titleWidth)),
				theme.Faint.Render(pad(ex.Difficulty, levelWidth)),
				theme.Faint.Render(ex.TypeLabel()))
		}
	}
	return nil
}
