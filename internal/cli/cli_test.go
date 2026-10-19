package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/charmbracelet/colorprofile"

	"github.com/kidskoding/leetkit/exercises"
	"github.com/kidskoding/leetkit/internal/config"
)

func run(t *testing.T, env Env, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errOut bytes.Buffer
	// Same writer as main; buffers are not terminals, so styling is stripped.
	env.Stdout = colorprofile.NewWriter(&out, nil)
	env.Stderr = colorprofile.NewWriter(&errOut, nil)
	if env.Paths == (config.Paths{}) {
		dir := t.TempDir()
		env.Paths = config.Paths{ConfigFile: filepath.Join(dir, "config.toml"), DataDir: filepath.Join(dir, "data")}
	}
	code = Run(context.Background(), args, env)
	return code, out.String(), errOut.String()
}

func bundled() Env { return Env{Exercises: exercises.FS} }

func TestList(t *testing.T) {
	code, out, errOut := run(t, bundled(), "list")
	if code != exitOK {
		t.Fatalf("exit %d, stderr: %s", code, errOut)
	}
	for _, want := range []string{"Iteration and Unpacking", "enumerate-01", "Find an Index", "Syntax Drill"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
}

func TestListRejectsArguments(t *testing.T) {
	code, _, errOut := run(t, bundled(), "list", "extra")
	if code != exitUsage || !strings.Contains(errOut, `unexpected argument "extra"`) {
		t.Errorf("exit %d, stderr %q", code, errOut)
	}
}

func TestListEmptyCurriculum(t *testing.T) {
	env := Env{Exercises: fstest.MapFS{"modules.toml": {Data: []byte("[[modules]]\nid = \"a\"\ntitle = \"A\"\n")}}}
	code, out, _ := run(t, env, "list")
	if code != exitOK || !strings.Contains(out, "No exercises found") {
		t.Errorf("exit %d, output %q", code, out)
	}
}

func TestListBrokenCurriculum(t *testing.T) {
	code, _, errOut := run(t, Env{Exercises: fstest.MapFS{}}, "list")
	if code != exitError || !strings.Contains(errOut, "load exercises") {
		t.Errorf("exit %d, stderr %q", code, errOut)
	}
}

func TestUnknownCommand(t *testing.T) {
	code, _, errOut := run(t, bundled(), "frobnicate")
	if code != exitUsage || !strings.Contains(errOut, `unknown command "frobnicate"`) {
		t.Errorf("exit %d, stderr %q", code, errOut)
	}
}

func TestHelp(t *testing.T) {
	for _, arg := range []string{"help", "-h", "--help"} {
		code, out, _ := run(t, bundled(), arg)
		if code != exitOK || !strings.Contains(out, "Usage:") {
			t.Errorf("%s: exit %d, output %q", arg, code, out)
		}
	}
}

func TestUnknownFlag(t *testing.T) {
	code, _, errOut := run(t, bundled(), "list", "-x")
	if code != exitUsage || !strings.Contains(errOut, "flag provided but not defined: -x") {
		t.Errorf("exit %d, stderr %q", code, errOut)
	}
}

func TestThemeFromConfig(t *testing.T) {
	tests := []struct {
		config, want string
		wantErr      bool
	}{
		{"", "terminal", false},
		{`theme = "gruvbox"`, "gruvbox", false},
		{`theme = "nope"`, "terminal", true},
	}
	for _, tt := range tests {
		dir := t.TempDir()
		paths := config.Paths{ConfigFile: filepath.Join(dir, "config.toml"), DataDir: dir}
		if tt.config != "" {
			if err := os.WriteFile(paths.ConfigFile, []byte(tt.config), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		a := newApp(Env{Paths: paths})
		if a.th.Name != tt.want || (a.cfgErr != nil) != tt.wantErr {
			t.Errorf("config %q: theme %q, cfgErr %v; want %q, error %v", tt.config, a.th.Name, a.cfgErr, tt.want, tt.wantErr)
		}
	}
}

func TestForcedColorUsesThemePalette(t *testing.T) {
	var out bytes.Buffer
	env := Env{
		Exercises: exercises.FS,
		Paths:     config.Paths{ConfigFile: filepath.Join(t.TempDir(), "config.toml")},
		Stdout:    colorprofile.NewWriter(&out, []string{"CLICOLOR_FORCE=1", "COLORTERM=truecolor"}),
		Stderr:    &bytes.Buffer{},
	}
	if code := Run(context.Background(), []string{"list"}, env); code != exitOK {
		t.Fatalf("exit %d", code)
	}
	// The terminal theme emits ANSI palette indices, never fixed RGB values.
	if !strings.Contains(out.String(), "\x1b[36m") || strings.Contains(out.String(), "38;2;") {
		t.Errorf("want ANSI palette colors only, got %q", out.String())
	}
}

func TestNoArgumentsWithoutTerminalPrintsUsage(t *testing.T) {
	code, out, _ := run(t, bundled())
	if code != exitOK || !strings.Contains(out, "Usage:") {
		t.Errorf("exit %d, output %q", code, out)
	}
}

func TestRunNeedsTerminalAndValidID(t *testing.T) {
	for args, want := range map[string]string{
		"run":              "expected one exercise id",
		"run nope-01":      `no exercise "nope-01"`,
		"run enumerate-01": "leetkit check enumerate-01",
	} {
		code, _, errOut := run(t, bundled(), strings.Fields(args)...)
		if code == exitOK || !strings.Contains(errOut, want) {
			t.Errorf("%s: exit %d, stderr %q", args, code, errOut)
		}
	}
}

func TestVersion(t *testing.T) {
	env := bundled()
	env.Version = "1.2.3"
	for _, arg := range []string{"version", "--version"} {
		code, out, _ := run(t, env, arg)
		if code != exitOK || out != "leetkit 1.2.3\n" {
			t.Errorf("%s: exit %d, output %q", arg, code, out)
		}
	}
}
