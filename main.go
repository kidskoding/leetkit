// Command leetkit is an interactive trainer for Python coding-interview syntax.
package main

import (
	"os"

	"github.com/anirudhkonidala/leetkit/exercises"
	"github.com/anirudhkonidala/leetkit/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], cli.Env{
		Exercises: exercises.FS,
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
	}))
}
