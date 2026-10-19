// Command leetkit is an interactive trainer for Python coding-interview syntax.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/term"

	"github.com/kidskoding/leetkit/exercises"
	"github.com/kidskoding/leetkit/internal/cli"
	"github.com/kidskoding/leetkit/internal/config"
)

// defaultPython lets packagers bundle an interpreter with pytest, for example
// -ldflags "-X main.defaultPython=/nix/store/...-python3/bin/python3".
var defaultPython string

func main() {
	paths, err := config.DefaultPaths()
	if err != nil {
		fmt.Fprintf(os.Stderr, "leetkit: %v\n", err)
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	code := cli.Run(ctx, os.Args[1:], cli.Env{
		Exercises: exercises.FS,
		Paths:     paths,
		// colorprofile downsamples colors for the terminal and strips them when
		// output is not a TTY or NO_COLOR is set.
		Stdout: colorprofile.NewWriter(os.Stdout, os.Environ()),
		Stderr: colorprofile.NewWriter(os.Stderr, os.Environ()),

		Interactive:   term.IsTerminal(os.Stdin.Fd()) && term.IsTerminal(os.Stdout.Fd()),
		DefaultPython: defaultPython,
	})
	stop()
	os.Exit(code)
}
