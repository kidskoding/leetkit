package cli

import (
	"context"
	"fmt"

	"github.com/anirudhkonidala/leetkit/internal/progress"
	"github.com/anirudhkonidala/leetkit/internal/tui"
)

func (a *app) dashboard(ctx context.Context, args []string) error {
	if len(args) > 0 {
		return usagef("unexpected argument %q", args[0])
	}
	if !a.Interactive {
		fmt.Fprint(a.Stdout, usage)
		return nil
	}
	cat, err := a.catalog()
	if err != nil {
		return err
	}
	store, err := progress.Open(a.Paths.ProgressFile())
	if err != nil {
		return err
	}
	var notice string
	switch {
	case store.Recovered != "":
		notice = "Progress file was corrupt; moved it to " + store.Recovered + " and started fresh."
	case a.cfgErr != nil:
		notice = "Config error; using defaults. Run `leetkit doctor` for details."
	}
	return tui.Run(ctx, tui.Options{Catalog: cat, Store: store, Theme: a.th, Notice: notice})
}
