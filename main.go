// Command leetkit is an interactive trainer for Python coding-interview syntax.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/anirudhkonidala/leetkit/exercises"
	"github.com/anirudhkonidala/leetkit/internal/cli"
	"github.com/anirudhkonidala/leetkit/internal/config"
)

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
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
	})
	stop()
	os.Exit(code)
}
