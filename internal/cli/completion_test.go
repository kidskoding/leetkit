package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestListIDs(t *testing.T) {
	code, out, _ := run(t, bundled(), "list", "--ids")
	ids := strings.Fields(out)
	if code != exitOK || len(ids) < 2 || ids[0] != "enumerate-01" || strings.Contains(out, " ") {
		t.Errorf("exit %d, output %q", code, out)
	}
}

func TestCompletionScriptsCoverCommands(t *testing.T) {
	for _, sh := range shells {
		code, out, errOut := run(t, bundled(), "completion", sh)
		if code != exitOK {
			t.Fatalf("%s: exit %d, %s", sh, code, errOut)
		}
		for name := range commands {
			if !strings.Contains(out, name) {
				t.Errorf("%s script does not complete %q", sh, name)
			}
		}
		if !strings.Contains(out, "list --ids") {
			t.Errorf("%s script does not complete exercise ids", sh)
		}
	}
	for _, args := range [][]string{{"completion"}, {"completion", "powershell"}} {
		if code, _, _ := run(t, bundled(), args...); code != exitUsage {
			t.Errorf("%v: exit %d, want usage error", args, code)
		}
	}
}

// fakeLeetkit puts a stand-in leetkit on PATH that prints two exercise ids.
func fakeLeetkit(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	script := "#!/bin/sh\nprintf 'enumerate-01\\nzip-01\\n'\n"
	if err := os.WriteFile(filepath.Join(dir, "leetkit"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir + string(os.PathListSeparator) + os.Getenv("PATH")
}

func completionScript(t *testing.T, sh string) string {
	t.Helper()
	data, err := completions.ReadFile("completions/leetkit." + sh)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "leetkit."+sh)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestBashCompletion(t *testing.T) {
	// Minimal bash builds (such as the Nix build sandbox's) lack programmable completion.
	if err := exec.Command("bash", "-c", "type compgen complete").Run(); err != nil {
		t.Skip("bash with programmable completion not available")
	}
	script := completionScript(t, "bash")
	cases := map[string][]string{
		"leetkit ch":       {"check"},
		"leetkit check en": {"enumerate-01"},
		"leetkit reset -":  {"--yes", "-y"},
		"leetkit run ":     {"enumerate-01", "zip-01"},
	}
	for line, want := range cases {
		words := strings.Split(line, " ")
		bashArgs := append([]string{"-c", `source "$1"; shift; COMP_WORDS=("$@"); COMP_CWORD=$(($# - 1)); ` +
			`_leetkit; printf '%s\n' "${COMPREPLY[@]}"`, "bash", script}, words...)
		cmd := exec.Command("bash", bashArgs...)
		cmd.Env = append(os.Environ(), "PATH="+fakeLeetkit(t))
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%q: %v\n%s", line, err, out)
		}
		got := strings.Fields(string(out))
		slices.Sort(got)
		if !slices.Equal(got, want) {
			t.Errorf("%q completes to %q, want %q", line, got, want)
		}
	}
}

func TestFishCompletion(t *testing.T) {
	if _, err := exec.LookPath("fish"); err != nil {
		t.Skip("fish not installed")
	}
	script := completionScript(t, "fish")
	cmd := exec.Command("fish", "--no-config", "-c", "source $argv[1]; complete -C 'leetkit check en'", script)
	cmd.Env = append(os.Environ(), "PATH="+fakeLeetkit(t))
	out, err := cmd.CombinedOutput()
	if err != nil || !strings.HasPrefix(string(out), "enumerate-01") {
		t.Errorf("fish completion = %q, %v", out, err)
	}
}

func TestZshCompletionParses(t *testing.T) {
	if _, err := exec.LookPath("zsh"); err != nil {
		t.Skip("zsh not installed")
	}
	if out, err := exec.Command("zsh", "-n", completionScript(t, "zsh")).CombinedOutput(); err != nil {
		t.Errorf("zsh -n: %v\n%s", err, out)
	}
}
