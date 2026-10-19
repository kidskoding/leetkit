// Package editor resolves the learner's external editor into a command.
package editor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var ErrNoEditor = errors.New("no editor found")

type Editor struct {
	Argv   []string
	Source string
}

// Resolve picks the editor from configured, $VISUAL, $EDITOR, then a
// platform fallback. Values may contain arguments, such as "code --wait".
func Resolve(configured string) (Editor, error) {
	settings := []struct{ value, source string }{
		{configured, "config"},
		{os.Getenv("VISUAL"), "$VISUAL"},
		{os.Getenv("EDITOR"), "$EDITOR"},
	}
	for _, s := range settings {
		if strings.TrimSpace(s.value) == "" {
			continue
		}
		argv, err := splitArgs(s.value)
		if err != nil {
			return Editor{}, fmt.Errorf("editor from %s: %w", s.source, err)
		}
		if _, err := exec.LookPath(argv[0]); err != nil {
			return Editor{}, fmt.Errorf("editor %q from %s: %w", argv[0], s.source, err)
		}
		return Editor{Argv: argv, Source: s.source}, nil
	}
	fallbacks := []string{"nvim", "vim", "vi", "nano"}
	if runtime.GOOS == "windows" {
		fallbacks = []string{"notepad"}
	}
	for _, name := range fallbacks {
		if _, err := exec.LookPath(name); err == nil {
			return Editor{Argv: []string{name}, Source: "fallback"}, nil
		}
	}
	return Editor{}, ErrNoEditor
}

// Command returns an unstarted command that opens file; stdio is left for the
// caller to attach.
func (e Editor) Command(file string) *exec.Cmd {
	args := append(append([]string{}, e.Argv[1:]...), file)
	return exec.Command(e.Argv[0], args...)
}

func (e Editor) String() string { return strings.Join(e.Argv, " ") }

// splitArgs splits a command line using POSIX-shell quoting rules for quotes
// and backslashes. Expansions and operators are not supported: no shell runs.
func splitArgs(s string) ([]string, error) {
	var (
		args    []string
		cur     strings.Builder
		inArg   bool
		quote   rune
		escaped bool
	)
	for _, r := range s {
		switch {
		case escaped:
			cur.WriteRune(r)
			escaped = false
		case r == '\\' && quote != '\'':
			escaped, inArg = true, true
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				cur.WriteRune(r)
			}
		case r == '\'' || r == '"':
			quote, inArg = r, true
		case r == ' ' || r == '\t' || r == '\n':
			if inArg {
				args = append(args, cur.String())
				cur.Reset()
				inArg = false
			}
		default:
			cur.WriteRune(r)
			inArg = true
		}
	}
	if quote != 0 || escaped {
		return nil, fmt.Errorf("unterminated quote or escape in %q", s)
	}
	if inArg {
		args = append(args, cur.String())
	}
	if len(args) == 0 {
		return nil, errors.New("empty command")
	}
	return args, nil
}
