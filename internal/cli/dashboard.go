package cli

import (
	"context"
	"fmt"

	"github.com/kidskoding/leetkit/internal/config"
	"github.com/kidskoding/leetkit/internal/editor"
	"github.com/kidskoding/leetkit/internal/exercise"
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
	return a.openDashboard(ctx, cat, "")
}

// run opens the dashboard on one exercise.
func (a *app) run(ctx context.Context, args []string) error {
	if len(args) != 1 {
		return usagef("expected one exercise id")
	}
	cat, err := a.catalog()
	if err != nil {
		return err
	}
	ex, err := findExercise(cat, args[0])
	if err != nil {
		return err
	}
	if !a.Interactive {
		return fmt.Errorf("the dashboard needs a terminal; use `leetkit check %s` to run its tests", ex.ID)
	}
	return a.openDashboard(ctx, cat, ex.ID)
}

func (a *app) openDashboard(ctx context.Context, cat *exercise.Catalog, start string) error {
	store, err := progress.Open(a.Paths.ProgressFile())
	if err != nil {
		return err
	}
	// The dashboard owns the terminal, so it shows the recovery notice itself
	// rather than using openStore's stderr warning.
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
	ed, edErr := editor.Resolve(cfg.Editor)
	return tui.Run(ctx, tui.Options{
		Catalog:   cat,
		Store:     store,
		Workspace: workspace.Workspace{Dir: a.Paths.WorkspaceDir()},
		Runner:    runner.Runner{Python: py.Path, Timeout: cfg.Timeout},
		RunnerErr: pyErr,
		Editor:    ed,
		EditorErr: edErr,
		Theme:     a.th,
		Notice:    notice,
		Start:     start,
	})
}
