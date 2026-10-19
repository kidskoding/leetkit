package editor

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestSplitArgs(t *testing.T) {
	tests := map[string][]string{
		"nvim":                         {"nvim"},
		"  code   --wait ":             {"code", "--wait"},
		`"/opt/My Editor/bin/ed" -n`:   {"/opt/My Editor/bin/ed", "-n"},
		`emacs -nw --eval '(foo "x")'`: {"emacs", "-nw", "--eval", `(foo "x")`},
		`ed a\ b ""`:                   {"ed", "a b", ""},
	}
	for in, want := range tests {
		got, err := splitArgs(in)
		if err != nil || !slices.Equal(got, want) {
			t.Errorf("splitArgs(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, bad := range []string{`vim "unclosed`, `vim \`, "   "} {
		if _, err := splitArgs(bad); err == nil {
			t.Errorf("splitArgs(%q) succeeded", bad)
		}
	}
}

// fakePath puts executables with the given names on an isolated PATH.
func fakePath(t *testing.T, names ...string) {
	t.Helper()
	dir := t.TempDir()
	for _, name := range names {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("#!/bin/sh\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", dir)
}

func TestResolveOrder(t *testing.T) {
	fakePath(t, "cfg-ed", "visual-ed", "env-ed", "vim")
	t.Setenv("VISUAL", "visual-ed -f")
	t.Setenv("EDITOR", "env-ed")

	tests := []struct {
		configured, visual, editor string
		want                       Editor
	}{
		{"cfg-ed --wait", "visual-ed -f", "env-ed", Editor{[]string{"cfg-ed", "--wait"}, "config"}},
		{"", "visual-ed -f", "env-ed", Editor{[]string{"visual-ed", "-f"}, "$VISUAL"}},
		{"", "", "env-ed", Editor{[]string{"env-ed"}, "$EDITOR"}},
		{"", "", "", Editor{[]string{"vim"}, "fallback"}},
	}
	for _, tt := range tests {
		t.Setenv("VISUAL", tt.visual)
		t.Setenv("EDITOR", tt.editor)
		got, err := Resolve(tt.configured)
		if err != nil || !slices.Equal(got.Argv, tt.want.Argv) || got.Source != tt.want.Source {
			t.Errorf("Resolve(%q) with VISUAL=%q EDITOR=%q = %+v, %v; want %+v",
				tt.configured, tt.visual, tt.editor, got, err, tt.want)
		}
	}
}

func TestResolveMissingConfiguredEditor(t *testing.T) {
	fakePath(t, "vim")
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "subl -w")
	_, err := Resolve("")
	if err == nil || !strings.Contains(err.Error(), `"subl" from $EDITOR`) {
		t.Errorf("err = %v; want a missing $EDITOR error instead of a silent fallback", err)
	}
}

func TestResolveNoEditor(t *testing.T) {
	fakePath(t)
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "")
	if _, err := Resolve(""); !errors.Is(err, ErrNoEditor) {
		t.Errorf("err = %v, want ErrNoEditor", err)
	}
}

func TestCommand(t *testing.T) {
	e := Editor{Argv: []string{"code", "--wait"}}
	cmd := e.Command("/tmp/a b.py")
	if !slices.Equal(cmd.Args, []string{"code", "--wait", "/tmp/a b.py"}) {
		t.Errorf("args = %q", cmd.Args)
	}
	if len(e.Argv) != 2 {
		t.Error("Command mutated Argv")
	}
}
