package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/kidskoding/leetkit/internal/progress"
)

func resetFixture(t *testing.T) (Env, string) {
	t.Helper()
	env := tempEnv(t)
	writeWork(t, env, "enumerate-01", []byte("my work"))
	store, err := progress.Open(env.Paths.ProgressFile())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Record("enumerate-01", true); err != nil {
		t.Fatal(err)
	}
	return env, env.Paths.WorkspaceDir() + "/enumerate-01/starter.py"
}

func completed(t *testing.T, env Env) bool {
	t.Helper()
	store, err := progress.Open(env.Paths.ProgressFile())
	if err != nil {
		t.Fatal(err)
	}
	return store.Completed("enumerate-01")
}

func TestResetConfirmed(t *testing.T) {
	for _, tt := range []struct {
		name  string
		args  []string
		stdin string
	}{
		{"prompt yes", []string{"reset", "enumerate-01"}, "y\n"},
		{"flag", []string{"reset", "--yes", "enumerate-01"}, ""},
		{"short flag", []string{"reset", "-y", "enumerate-01"}, ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			env, path := resetFixture(t)
			env.Stdin = strings.NewReader(tt.stdin)
			code, out, errOut := run(t, env, tt.args...)
			if code != exitOK || !strings.Contains(out, "✓ Reset enumerate-01") {
				t.Fatalf("exit %d\nstdout: %s\nstderr: %s", code, out, errOut)
			}
			data, err := os.ReadFile(path)
			if err != nil || !strings.Contains(string(data), "def find_index") || strings.Contains(string(data), "my work") {
				t.Errorf("working copy not restored: %q, %v", data, err)
			}
			if completed(t, env) {
				t.Error("progress not cleared")
			}
		})
	}
}

func TestResetDeclined(t *testing.T) {
	for _, stdin := range []string{"n\n", "\n", "whatever\n"} {
		env, path := resetFixture(t)
		env.Stdin = strings.NewReader(stdin)
		code, out, _ := run(t, env, "reset", "enumerate-01")
		data, _ := os.ReadFile(path)
		if code != exitOK || !strings.Contains(out, "Nothing changed") || string(data) != "my work" || !completed(t, env) {
			t.Errorf("stdin %q: exit %d, out %q, file %q", stdin, code, out, data)
		}
	}
}

func TestResetWithoutAnswer(t *testing.T) {
	env, path := resetFixture(t)
	env.Stdin = strings.NewReader("")
	code, _, errOut := run(t, env, "reset", "enumerate-01")
	data, _ := os.ReadFile(path)
	if code != exitError || !strings.Contains(errOut, "pass --yes") || string(data) != "my work" {
		t.Errorf("exit %d, stderr %q, file %q", code, errOut, data)
	}
}

func TestResetArguments(t *testing.T) {
	for args, want := range map[string]string{
		"reset":           "expected one exercise id",
		"reset nope-01 y": "expected one exercise id",
		"reset nope-01":   `no exercise "nope-01"`,
	} {
		env := tempEnv(t)
		env.Stdin = strings.NewReader("y\n")
		_, _, errOut := run(t, env, strings.Fields(args)...)
		if !strings.Contains(errOut, want) {
			t.Errorf("%s: stderr %q, want %q", args, errOut, want)
		}
	}
}
