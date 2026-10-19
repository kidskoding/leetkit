package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/anirudhkonidala/leetkit/internal/exercise"
	"github.com/anirudhkonidala/leetkit/internal/runner"
	"github.com/anirudhkonidala/leetkit/internal/workspace"
)

// errSilent signals a non-zero exit after the command already reported why.
var errSilent = errors.New("silent failure")

func runCheck(ctx context.Context, env Env, args []string) error {
	fset := flag.NewFlagSet("check", flag.ContinueOnError)
	verbose := fset.Bool("v", false, "show full pytest output")
	if err := parseFlags(fset, args); err != nil {
		return err
	}
	if fset.NArg() != 1 {
		return usagef("expected one exercise id")
	}
	cat, err := loadCatalog(env)
	if err != nil {
		return err
	}
	ex, err := findExercise(cat, fset.Arg(0))
	if err != nil {
		return err
	}
	cfg, err := loadConfig(env)
	if err != nil {
		return err
	}
	py, err := findPython(env, cfg)
	if err != nil {
		return fmt.Errorf("%w\nRun `leetkit doctor` for setup steps.", err)
	}

	tests, err := cat.ReadFile(ex, exercise.TestFile)
	if err != nil {
		return fmt.Errorf("read tests: %w", err)
	}
	starter, err := cat.ReadFile(ex, exercise.StarterFile)
	if err != nil {
		return fmt.Errorf("read starter: %w", err)
	}
	ws := workspace.Workspace{Dir: env.Paths.WorkspaceDir()}
	codePath, err := ws.Ensure(ex.ID, starter)
	if err != nil {
		return err
	}
	code, err := os.ReadFile(codePath)
	if err != nil {
		return fmt.Errorf("read your code: %w", err)
	}

	r := runner.Runner{Python: py.Path, Timeout: cfg.Timeout}
	res, err := r.Run(ctx, runner.Job{Tests: tests, Code: code, CodePath: codePath})
	if errors.Is(err, runner.ErrPytestMissing) {
		return fmt.Errorf("pytest is not installed for %s (from %s).\nRun `leetkit doctor` for setup steps.", py.Path, py.Source)
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(env.Stdout, "%s · %s\n", ex.ID, ex.Title)
	printResult(env.Stdout, res, *verbose)
	fmt.Fprintf(env.Stdout, "\nYour code: %s\n", codePath)
	if res.Status != runner.StatusPassed {
		return errSilent
	}
	return nil
}

func findExercise(cat *exercise.Catalog, id string) (*exercise.Exercise, error) {
	ex, ok := cat.Get(id)
	if !ok {
		return nil, fmt.Errorf("no exercise %q. Run `leetkit list` to see available ids.", id)
	}
	return ex, nil
}

func printResult(w io.Writer, res *runner.Result, verbose bool) {
	passed, failed := res.Count(runner.OutcomePassed), res.Count(runner.OutcomeFailed)
	elapsed := res.Duration.Round(10 * time.Millisecond)
	switch res.Status {
	case runner.StatusPassed:
		fmt.Fprintf(w, "✓ All %d tests passed (%s)\n", passed, elapsed)
	case runner.StatusFailed:
		fmt.Fprintf(w, "✗ %d passed, %d failed (%s)\n", passed, failed, elapsed)
	default:
		fmt.Fprintf(w, "✗ %s\n", res.Problem)
	}
	for _, c := range res.Cases {
		if c.Outcome != runner.OutcomeFailed && c.Outcome != runner.OutcomeError {
			continue
		}
		fmt.Fprintf(w, "\n✗ %s\n%s\n", c.Name, indent(c.Message, "    "))
	}
	if verbose && res.Output != "" {
		fmt.Fprintf(w, "\n--- pytest output ---\n%s", res.Output)
	}
}

func indent(s, prefix string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i, l := range lines {
		lines[i] = prefix + l
	}
	return strings.Join(lines, "\n")
}
