package cli

import (
	"context"
	"fmt"

	"github.com/kidskoding/leetkit/internal/config"
	"github.com/kidskoding/leetkit/internal/progress"
	"github.com/kidskoding/leetkit/internal/runner"
	"github.com/kidskoding/leetkit/internal/tui"
	"github.com/kidskoding/leetkit/internal/workspace"
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
	cfg := a.cfg
	var notice string
	switch {
	case store.Recovered != "":
		notice = "Progress file was corrupt; moved it to " + store.Recovered + " and started fresh."
	case a.cfgErr != nil:
		cfg = config.Default()
		notice = "Config error; using defaults. Run `leetkit doctor` for details."
	}
	py, pyErr := a.findPython(cfg)
	return tui.Run(ctx, tui.Options{
		Catalog:   cat,
		Store:     store,
		Workspace: workspace.Workspace{Dir: a.Paths.WorkspaceDir()},
		Runner:    runner.Runner{Python: py.Path, Timeout: cfg.Timeout},
		RunnerErr: pyErr,
		Theme:     a.th,
		Notice:    notice,
	})
}
