package cli

import (
	"context"
	"embed"
	"fmt"
	"strings"
)

//go:embed completions
var completions embed.FS

var shells = []string{"bash", "zsh", "fish"}

func (a *app) completion(_ context.Context, args []string) error {
	if len(args) != 1 {
		return usagef("expected a shell: %s", strings.Join(shells, ", "))
	}
	script, err := completions.ReadFile("completions/leetkit." + args[0])
	if err != nil {
		return usagef("unsupported shell %q; choose %s", args[0], strings.Join(shells, ", "))
	}
	_, err = fmt.Fprint(a.Stdout, string(script))
	return err
}
