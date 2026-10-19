package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/anirudhkonidala/leetkit/internal/exercise"
	"github.com/anirudhkonidala/leetkit/internal/runner"
	"github.com/anirudhkonidala/leetkit/internal/theme"
	"github.com/anirudhkonidala/leetkit/internal/workspace"
)

// errSilent signals a non-zero exit after the command already reported why.
var errSilent = errors.New("silent failure")

type checker struct {
	cat    *exercise.Catalog
	ws     workspace.Workspace
	runner runner.Runner
	py     runner.Candidate
}

func runCheck(ctx context.Context, env Env, args []string) error {
	fset := flag.NewFlagSet("check", flag.ContinueOnError)
	verbose := fset.Bool("v", false, "show full pytest output")
	if err := parseFlags(fset, args); err != nil {
		return err
	}
	if fset.NArg() > 1 {
		return usagef("expected at most one exercise id")
	}
	cat, err := loadCatalog(env)
	if err != nil {
		return err
	}
	var ex *exercise.Exercise
	if fset.NArg() == 1 {
		if ex, err = findExercise(cat, fset.Arg(0)); err != nil {
			return err
		}
	}
	cfg, err := loadConfig(env)
	if err != nil {
		return err
	}
	py, err := findPython(env, cfg)
	if err != nil {
		return fmt.Errorf("%w\nRun `leetkit doctor` for setup steps.", err)
	}
	c := checker{
		cat:    cat,
		ws:     workspace.Workspace{Dir: env.Paths.WorkspaceDir()},
		runner: runner.Runner{Python: py.Path, Timeout: cfg.Timeout},
		py:     py,
	}
	if ex != nil {
		return c.one(ctx, env.Stdout, ex, *verbose)
	}
	return c.all(ctx, env.Stdout)
}

func (c checker) one(ctx context.Context, w io.Writer, ex *exercise.Exercise, verbose bool) error {
	starter, err := c.cat.ReadFile(ex, exercise.StarterFile)
	if err != nil {
		return fmt.Errorf("read starter: %w", err)
	}
	path, err := c.ws.Ensure(ex.ID, starter)
	if err != nil {
		return err
	}
	res, err := c.check(ctx, ex, path)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s %s %s\n", theme.ID.Bold(true).Render(ex.ID), theme.Faint.Render("·"), theme.Text.Render(ex.Title))
	printResult(w, res, verbose)
	fmt.Fprintf(w, "\n%s %s\n", theme.Faint.Render("Your code:"), theme.ID.Render(path))
	if res.Status != runner.StatusPassed {
		return errSilent
	}
	return nil
}

type checkOutcome struct {
	res     *runner.Result
	err     error
	started bool
}

// all checks every exercise the learner has opened. Untouched exercises are
// reported as not started instead of running their failing starter code.
func (c checker) all(ctx context.Context, w io.Writer) error {
	if _, err := runner.PytestVersion(ctx, c.py.Path); err != nil {
		return c.pytestError(err)
	}
	exs := c.cat.All()
	outcomes := make([]checkOutcome, len(exs))
	sem := make(chan struct{}, runtime.NumCPU())
	var wg sync.WaitGroup
	for i, ex := range exs {
		path := c.ws.Path(ex.ID)
		if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
			continue
		}
		outcomes[i].started = true
		wg.Go(func() {
			sem <- struct{}{}
			defer func() { <-sem }()
			outcomes[i].res, outcomes[i].err = c.check(ctx, ex, path)
		})
	}
	wg.Wait()

	idWidth, titleWidth := 0, 0
	for _, ex := range exs {
		idWidth = max(idWidth, lipgloss.Width(ex.ID))
		titleWidth = max(titleWidth, lipgloss.Width(ex.Title))
	}
	var passing, failing, notStarted int
	i := 0
	for mi, m := range c.cat.Modules {
		if mi > 0 {
			fmt.Fprintln(w)
		}
		fmt.Fprintln(w, theme.Heading.Render(m.Title))
		for _, ex := range m.Exercises {
			o := outcomes[i]
			i++
			mark, status := describe(o)
			switch {
			case !o.started:
				notStarted++
			case o.err == nil && o.res.Status == runner.StatusPassed:
				passing++
			default:
				failing++
			}
			fmt.Fprintf(w, "  %s %s  %s  %s\n", mark,
				theme.ID.Render(pad(ex.ID, idWidth)), theme.Text.Render(pad(ex.Title, titleWidth)), status)
		}
	}

	summary := []string{
		theme.Pass.Render(fmt.Sprintf("%d/%d passing", passing, len(exs))),
		theme.Bad.Render(fmt.Sprintf("%d failing", failing)),
		theme.Faint.Render(fmt.Sprintf("%d not started", notStarted)),
	}
	fmt.Fprintf(w, "\n%s\n", strings.Join(summary, theme.Faint.Render(" · ")))
	if failing > 0 {
		fmt.Fprintln(w, theme.Faint.Render("Run `leetkit check <id>` for details."))
		return errSilent
	}
	return nil
}

func (c checker) check(ctx context.Context, ex *exercise.Exercise, path string) (*runner.Result, error) {
	tests, err := c.cat.ReadFile(ex, exercise.TestFile)
	if err != nil {
		return nil, fmt.Errorf("read tests: %w", err)
	}
	code, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read your code: %w", err)
	}
	res, err := c.runner.Run(ctx, runner.Job{Tests: tests, Code: code, CodePath: path})
	if err != nil {
		return nil, c.pytestError(err)
	}
	return res, nil
}

func (c checker) pytestError(err error) error {
	if errors.Is(err, runner.ErrPytestMissing) {
		return fmt.Errorf("pytest is not installed for %s (from %s).\nRun `leetkit doctor` for setup steps.", c.py.Path, c.py.Source)
	}
	return err
}

func describe(o checkOutcome) (mark, status string) {
	switch {
	case !o.started:
		return theme.Faint.Render("○"), theme.Faint.Render("not started")
	case o.err != nil:
		return theme.Fail.Render("✗"), theme.Bad.Render(firstLine(o.err.Error()))
	}
	r := o.res
	passed := r.Count(runner.OutcomePassed)
	switch r.Status {
	case runner.StatusPassed:
		return theme.Pass.Render("✓"), theme.Good.Render(fmt.Sprintf("%d/%d passed", passed, len(r.Cases)))
	case runner.StatusFailed:
		return theme.Fail.Render("✗"), theme.Bad.Render(fmt.Sprintf("%d/%d passed", passed, len(r.Cases)))
	case runner.StatusError:
		return theme.Fail.Render("✗"), theme.Bad.Render("could not run (syntax or import error)")
	default:
		return theme.Warn.Render("!"), theme.Warn.Render(r.Status.String())
	}
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
	elapsed := theme.Faint.Render(res.Duration.Round(10 * time.Millisecond).String())
	switch res.Status {
	case runner.StatusPassed:
		fmt.Fprintf(w, "%s %s\n", theme.Pass.Render(fmt.Sprintf("✓ All %d tests passed", passed)), elapsed)
	case runner.StatusFailed:
		fmt.Fprintf(w, "%s %s, %s %s\n", theme.Fail.Render("✗"),
			theme.Good.Render(fmt.Sprintf("%d passed", passed)),
			theme.Bad.Render(fmt.Sprintf("%d failed", failed)), elapsed)
	case runner.StatusError:
		fmt.Fprintln(w, theme.Fail.Render("✗ "+res.Problem))
	default:
		fmt.Fprintln(w, theme.Warn.Render("! "+res.Problem))
	}
	for _, c := range res.Cases {
		if c.Outcome != runner.OutcomeFailed && c.Outcome != runner.OutcomeError {
			continue
		}
		fmt.Fprintf(w, "\n  %s\n%s\n", theme.Bad.Render("✗ "+c.Name), indent(c.Message, "    "))
	}
	if verbose && res.Output != "" {
		fmt.Fprintf(w, "\n%s\n%s", theme.Faint.Render("─── pytest output ───"), res.Output)
	}
}

func pad(s string, width int) string {
	return s + strings.Repeat(" ", max(width-lipgloss.Width(s), 0))
}

func firstLine(s string) string {
	line, _, _ := strings.Cut(s, "\n")
	return line
}

func indent(s, prefix string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i, l := range lines {
		lines[i] = prefix + l
	}
	return strings.Join(lines, "\n")
}
