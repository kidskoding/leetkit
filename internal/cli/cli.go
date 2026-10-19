// Package cli implements the leetkit command-line interface.
package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/anirudhkonidala/leetkit/internal/config"
	"github.com/anirudhkonidala/leetkit/internal/exercise"
	"github.com/anirudhkonidala/leetkit/internal/runner"
	"github.com/anirudhkonidala/leetkit/internal/theme"
)

type Env struct {
	Exercises fs.FS
	Paths     config.Paths
	Stdout    io.Writer
	Stderr    io.Writer
	// Interactive reports whether stdin and stdout are a terminal, which the
	// dashboard requires.
	Interactive bool
	// DefaultPython is a packaged interpreter tried after user settings and
	// the LeetKit venv, before PATH.
	DefaultPython string
}

const (
	exitOK    = 0
	exitError = 1
	exitUsage = 2
)

const usage = `LeetKit - master Python's toolkit for coding interviews.

Usage:
  leetkit                 open the interactive dashboard
  leetkit list            list modules and exercises
  leetkit check           check every exercise you have started
  leetkit check [-v] <id> run one exercise's tests
  leetkit doctor          check Python, pytest, and editor setup
  leetkit help            show this help
`

type usageError struct{ error }

func usagef(format string, a ...any) error {
	return usageError{fmt.Errorf(format, a...)}
}

var commands = map[string]func(a *app, ctx context.Context, args []string) error{
	"list":   (*app).list,
	"check":  (*app).check,
	"doctor": (*app).doctor,
}

// app is the per-invocation state shared by commands.
type app struct {
	Env
	cfg config.Config
	// cfgErr is kept rather than returned so commands that do not need
	// settings still work, and doctor can report it.
	cfgErr error
	th     theme.Theme
}

func newApp(env Env) *app {
	a := &app{Env: env}
	a.cfg, a.cfgErr = config.Load(env.Paths.ConfigFile)
	th, err := theme.Load(a.cfg.Theme)
	if a.cfgErr != nil || err != nil {
		th, _ = theme.Load(theme.Default)
	}
	a.th = th
	return a
}

// Run executes the command named in args and returns the process exit code.
// With no arguments it opens the dashboard.
func Run(ctx context.Context, args []string, env Env) int {
	name, rest := "", args
	if len(args) > 0 {
		name, rest = args[0], args[1:]
	}
	switch name {
	case "help", "-h", "--help":
		fmt.Fprint(env.Stdout, usage)
		return exitOK
	}
	a := newApp(env)
	cmd, prefix := (*app).dashboard, "leetkit:"
	if name != "" {
		var ok bool
		if cmd, ok = commands[name]; !ok {
			fmt.Fprintf(env.Stderr, "%s unknown command %q\n\n%s", a.th.Fail.Render("leetkit:"), name, usage)
			return exitUsage
		}
		prefix = "leetkit " + name + ":"
	}
	if err := cmd(a, ctx, rest); err != nil {
		if errors.Is(err, errSilent) {
			return exitError
		}
		fmt.Fprintf(env.Stderr, "%s %v\n", a.th.Fail.Render(prefix), err)
		if errors.As(err, new(usageError)) {
			fmt.Fprintf(env.Stderr, "\n%s", usage)
			return exitUsage
		}
		return exitError
	}
	return exitOK
}

func parseFlags(fset *flag.FlagSet, args []string) error {
	fset.SetOutput(io.Discard)
	if err := fset.Parse(args); err != nil {
		return usageError{err}
	}
	return nil
}

func (a *app) catalog() (*exercise.Catalog, error) {
	cat, err := exercise.Load(a.Exercises)
	if err != nil {
		return nil, fmt.Errorf("load exercises: %w", err)
	}
	return cat, nil
}

func (a *app) config() (config.Config, error) {
	if a.cfgErr != nil {
		return a.cfg, fmt.Errorf("%w\nFix or remove the config file.", a.cfgErr)
	}
	return a.cfg, nil
}

func (a *app) findPython(cfg config.Config) (runner.Candidate, error) {
	return runner.FindPython(
		runner.Candidate{Path: os.Getenv("LEETKIT_PYTHON"), Source: "$LEETKIT_PYTHON", Explicit: true},
		runner.Candidate{Path: cfg.Python, Source: "config", Explicit: true},
		runner.Candidate{Path: a.Paths.VenvPython(), Source: "LeetKit venv"},
		runner.Candidate{Path: a.DefaultPython, Source: "package default"},
	)
}
