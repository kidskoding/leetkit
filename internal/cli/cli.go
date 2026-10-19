// Package cli implements the leetkit command-line interface.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"text/tabwriter"

	"github.com/anirudhkonidala/leetkit/internal/exercise"
)

type Env struct {
	Exercises fs.FS
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
  leetkit help            show this help
`

type usageError struct{ error }

func usagef(format string, a ...any) error {
	return usageError{fmt.Errorf(format, a...)}
}

var commands = map[string]func(env Env, args []string) error{
	"list": runList,
}

// Run executes the command named in args and returns the process exit code.
func Run(args []string, env Env) int {
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
	if err := cmd(env, args[1:]); err != nil {
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

func runList(env Env, args []string) error {
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
