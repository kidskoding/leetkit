package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/anirudhkonidala/leetkit/internal/config"
	"github.com/anirudhkonidala/leetkit/internal/editor"
	"github.com/anirudhkonidala/leetkit/internal/runner"
)

type report struct {
	w      io.Writer
	failed bool
}

func (r *report) ok(name, format string, a ...any) {
	fmt.Fprintf(r.w, "✓ %-10s %s\n", name, fmt.Sprintf(format, a...))
}

func (r *report) fail(name, problem, fix string) {
	r.failed = true
	fmt.Fprintf(r.w, "✗ %-10s %s\n%s\n", name, problem, indent(fix, "    "))
}

func runDoctor(ctx context.Context, env Env, args []string) error {
	fset := flag.NewFlagSet("doctor", flag.ContinueOnError)
	if err := parseFlags(fset, args); err != nil {
		return err
	}
	if fset.NArg() > 0 {
		return usagef("unexpected argument %q", fset.Arg(0))
	}
	r := &report{w: env.Stdout}

	cfg, err := config.Load(env.Paths.ConfigFile)
	switch {
	case err != nil:
		r.fail("Config", err.Error(), "Fix or remove the config file.")
		cfg = config.Default()
	case fileExists(env.Paths.ConfigFile):
		r.ok("Config", "%s", env.Paths.ConfigFile)
	default:
		r.ok("Config", "%s (not found; using defaults)", env.Paths.ConfigFile)
	}

	if cat, err := loadCatalog(env); err != nil {
		r.fail("Exercises", err.Error(), "The bundled curriculum is broken. Please report this bug.")
	} else {
		r.ok("Exercises", "%d bundled", len(cat.All()))
	}

	py, err := findPython(env, cfg)
	if errors.Is(err, runner.ErrNoPython) {
		r.fail("Python", err.Error(), "Install Python 3.10 or newer, then:\n"+pytestSetup(env.Paths))
	} else if err != nil {
		r.fail("Python", err.Error(), "Correct or unset that setting.")
	} else if version, err := runner.PythonVersion(ctx, py.Path); err != nil {
		r.fail("Python", fmt.Sprintf("%s (from %s) does not run: %v", py.Path, py.Source, err),
			"Point LEETKIT_PYTHON or `python` in the config file at a working interpreter.")
	} else {
		r.ok("Python", "%s %s (from %s)", py.Path, version, py.Source)
		if version, err := runner.PytestVersion(ctx, py.Path); errors.Is(err, runner.ErrPytestMissing) {
			r.fail("pytest", "not installed for "+py.Path, pytestSetup(env.Paths))
		} else if err != nil {
			r.fail("pytest", err.Error(), pytestSetup(env.Paths))
		} else {
			r.ok("pytest", "%s", version)
		}
	}

	if ed, err := editor.Resolve(cfg.Editor); err != nil {
		r.fail("Editor", err.Error(), "Set $EDITOR (for example `export EDITOR=nvim`) or `editor` in "+env.Paths.ConfigFile+".")
	} else {
		r.ok("Editor", "%s (from %s)", ed, ed.Source)
	}

	if err := checkWritable(env.Paths.DataDir); err != nil {
		r.fail("Data", err.Error(), "LeetKit stores your code and progress here; make it writable.")
	} else {
		r.ok("Data", "%s", env.Paths.DataDir)
	}

	if r.failed {
		return errSilent
	}
	fmt.Fprintln(env.Stdout, "\nAll checks passed.")
	return nil
}

// pytestSetup prints commands that create LeetKit's own venv, which
// FindPython picks up automatically. Nothing is installed for the user.
func pytestSetup(p config.Paths) string {
	venv := p.VenvDir()
	if _, err := exec.LookPath("uv"); err == nil {
		return fmt.Sprintf("Create a LeetKit environment with uv:\n  uv venv %s\n  uv pip install --python %s pytest", venv, venv)
	}
	python := "python3"
	if runtime.GOOS == "windows" {
		python = "py -3"
	}
	return fmt.Sprintf("Create a LeetKit environment:\n  %s -m venv %s\n  %s -m pip install pytest",
		python, venv, p.VenvPython())
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func checkWritable(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".doctor-*")
	if err != nil {
		return fmt.Errorf("%s is not writable: %w", dir, err)
	}
	name := f.Name()
	f.Close()
	return os.Remove(name)
}
