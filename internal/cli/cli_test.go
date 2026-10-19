package cli

import (
	"bytes"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/anirudhkonidala/leetkit/exercises"
)

func run(t *testing.T, env Env, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errOut bytes.Buffer
	env.Stdout, env.Stderr = &out, &errOut
	code = Run(args, env)
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
