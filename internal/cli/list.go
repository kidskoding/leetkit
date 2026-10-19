package cli

import (
	"context"
	"flag"
	"fmt"
	"text/tabwriter"
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
	tw := tabwriter.NewWriter(env.Stdout, 0, 0, 2, ' ', 0)
	for i, m := range cat.Modules {
		if i > 0 {
			fmt.Fprintln(tw)
		}
		fmt.Fprintln(tw, m.Title)
		for _, ex := range m.Exercises {
			fmt.Fprintf(tw, "  %s\t%s\t%s\t%s\n", ex.ID, ex.Title, ex.Difficulty, ex.TypeLabel())
		}
	}
	return tw.Flush()
}
