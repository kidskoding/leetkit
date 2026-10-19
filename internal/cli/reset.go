package cli

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/kidskoding/leetkit/internal/workspace"
)

func (a *app) reset(_ context.Context, args []string) error {
	fset := flag.NewFlagSet("reset", flag.ContinueOnError)
	yes := fset.Bool("yes", false, "skip the confirmation prompt")
	fset.BoolVar(yes, "y", false, "skip the confirmation prompt")
	if err := parseFlags(fset, args); err != nil {
		return err
	}
	if fset.NArg() != 1 {
		return usagef("expected one exercise id")
	}
	cat, err := a.catalog()
	if err != nil {
		return err
	}
	ex, err := findExercise(cat, fset.Arg(0))
	if err != nil {
		return err
	}
	if !*yes {
		fmt.Fprintf(a.Stdout, "Reset %s to the starter code and clear its progress? Your changes will be lost. [y/N] ", ex.ID)
		ok, err := confirm(a.Stdin)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Fprintln(a.Stdout, a.th.Faint.Render("Nothing changed."))
			return nil
		}
	}

	starter, err := cat.Starter(ex)
	if err != nil {
		return fmt.Errorf("read starter: %w", err)
	}
	path, err := workspace.Workspace{Dir: a.Paths.WorkspaceDir()}.Reset(ex.ID, starter)
	if err != nil {
		return err
	}
	store, err := a.openStore()
	if err != nil {
		return err
	}
	if err := store.Reset(ex.ID); err != nil {
		return err
	}
	fmt.Fprintf(a.Stdout, "%s %s\n%s %s\n", a.th.Pass.Render("✓ Reset"), a.th.ID.Render(ex.ID),
		a.th.Faint.Render("Your code:"), a.th.ID.Render(path))
	return nil
}

// confirm reads one answer line; only "y" or "yes" confirms. With no input,
// such as a closed pipe, it refuses and says how to skip the prompt.
func confirm(in io.Reader) (bool, error) {
	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil && line == "" {
		return false, errors.New("no answer on stdin; pass --yes to reset without confirmation")
	}
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		return true, nil
	}
	return false, nil
}
