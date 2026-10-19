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
)

type Env struct {
	Exercises fs.FS
	Paths     config.Paths
	Stdout    io.Writer
	Stderr    io.Writer
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
  leetkit check [-v] <id> run an exercise's tests
  leetkit doctor          check Python, pytest, and editor setup
  leetkit help            show this help
`

type usageError struct{ error }

func usagef(format string, a ...any) error {
	return usageError{fmt.Errorf(format, a...)}
}

var commands = map[string]func(ctx context.Context, env Env, args []string) error{
	"list":   runList,
	"check":  runCheck,
	"doctor": runDoctor,
}

// Run executes the command named in args and returns the process exit code.
func Run(ctx context.Context, args []string, env Env) int {
	if len(args) == 0 {
		fmt.Fprint(env.Stdout, usage)
		return exitOK
	}
	switch args[0] {
	case "help", "-h", "--help":
		fmt.Fprint(env.Stdout, usage)
		return exitOK
	}
	cmd, ok := commands[args[0]]
	if !ok {
		fmt.Fprintf(env.Stderr, "leetkit: unknown command %q\n\n%s", args[0], usage)
		return exitUsage
	}
	if err := cmd(ctx, env, args[1:]); err != nil {
		if errors.Is(err, errSilent) {
			return exitError
		}
		fmt.Fprintf(env.Stderr, "leetkit %s: %v\n", args[0], err)
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

func loadCatalog(env Env) (*exercise.Catalog, error) {
	cat, err := exercise.Load(env.Exercises)
	if err != nil {
		return nil, fmt.Errorf("load exercises: %w", err)
	}
	return cat, nil
}

func loadConfig(env Env) (config.Config, error) {
	cfg, err := config.Load(env.Paths.ConfigFile)
	if err != nil {
		return cfg, fmt.Errorf("%w\nFix or remove the config file.", err)
	}
	return cfg, nil
}

func findPython(env Env, cfg config.Config) (runner.Candidate, error) {
	return runner.FindPython(
		runner.Candidate{Path: os.Getenv("LEETKIT_PYTHON"), Source: "$LEETKIT_PYTHON", Explicit: true},
		runner.Candidate{Path: cfg.Python, Source: "config", Explicit: true},
		runner.Candidate{Path: env.Paths.VenvPython(), Source: "LeetKit venv"},
	)
}
