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
	"strings"

	"github.com/kidskoding/leetkit/internal/config"
	"github.com/kidskoding/leetkit/internal/editor"
	"github.com/kidskoding/leetkit/internal/report"
	"github.com/kidskoding/leetkit/internal/runner"
	"github.com/kidskoding/leetkit/internal/theme"
)

type checklist struct {
	w      io.Writer
	th     theme.Theme
	failed bool
}

func (r *checklist) ok(name, format string, a ...any) {
	fmt.Fprintf(r.w, "%s %s %s\n", r.th.Pass.Render("✓"), r.th.Text.Render(pad(name, 10)),
		r.th.Faint.Render(fmt.Sprintf(format, a...)))
}

func (r *checklist) fail(name, problem, fix string) {
	r.failed = true
	fmt.Fprintf(r.w, "%s %s %s\n%s\n", r.th.Fail.Render("✗"), r.th.Text.Render(pad(name, 10)),
		r.th.Bad.Render(problem), report.Indent(fix, "    "))
}

func (a *app) doctor(ctx context.Context, args []string) error {
	fset := flag.NewFlagSet("doctor", flag.ContinueOnError)
	if err := parseFlags(fset, args); err != nil {
		return err
	}
	if fset.NArg() > 0 {
		return usagef("unexpected argument %q", fset.Arg(0))
	}
	r := &checklist{w: a.Stdout, th: a.th}

	cfg := a.cfg
	switch {
	case a.cfgErr != nil:
		r.fail("Config", a.cfgErr.Error(), "Fix or remove the config file.")
		cfg = config.Default()
	case fileExists(a.Paths.ConfigFile):
		r.ok("Config", "%s", a.Paths.ConfigFile)
	default:
		r.ok("Config", "%s (not found; using defaults)", a.Paths.ConfigFile)
	}

	if cat, err := a.catalog(); err != nil {
		r.fail("Exercises", err.Error(), "The bundled curriculum is broken. Please report this bug.")
	} else {
		r.ok("Exercises", "%d bundled", len(cat.All()))
	}

	py, err := a.findPython(cfg)
	if errors.Is(err, runner.ErrNoPython) {
		r.fail("Python", err.Error(), "Install Python 3.10 or newer, then:\n"+pytestSetup(a.Paths))
	} else if err != nil {
		r.fail("Python", err.Error(), "Correct or unset that setting.")
	} else if version, err := runner.PythonVersion(ctx, py.Path); err != nil {
		r.fail("Python", fmt.Sprintf("%s (from %s) does not run: %v", py.Path, py.Source, err),
			"Point LEETKIT_PYTHON or `python` in the config file at a working interpreter.")
	} else {
		r.ok("Python", "%s %s (from %s)", py.Path, version, py.Source)
		if version, err := runner.PytestVersion(ctx, py.Path); errors.Is(err, runner.ErrPytestMissing) {
			r.fail("pytest", "not installed for "+py.Path, pytestSetup(a.Paths))
		} else if err != nil {
			r.fail("pytest", err.Error(), pytestSetup(a.Paths))
		} else {
			r.ok("pytest", "%s", version)
		}
	}

	if ed, err := editor.Resolve(cfg.Editor); err != nil {
		r.fail("Editor", err.Error(), "Set $EDITOR (for example `export EDITOR=nvim`) or `editor` in "+a.Paths.ConfigFile+".")
	} else {
		r.ok("Editor", "%s (from %s)", ed, ed.Source)
	}

	r.ok("Theme", "%s (available: %s)", a.th.Name, strings.Join(theme.Names(), ", "))

	if err := checkWritable(a.Paths.DataDir); err != nil {
		r.fail("Data", err.Error(), "LeetKit stores your code and progress here; make it writable.")
	} else {
		r.ok("Data", "%s", a.Paths.DataDir)
	}

	if r.failed {
		return errSilent
	}
	fmt.Fprintln(a.Stdout, "\n"+a.th.Pass.Render("All checks passed."))
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
