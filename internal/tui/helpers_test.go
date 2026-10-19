package tui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/kidskoding/leetkit/internal/exercise"
	"github.com/kidskoding/leetkit/internal/progress"
	"github.com/kidskoding/leetkit/internal/runner"
	"github.com/kidskoding/leetkit/internal/theme"
	"github.com/kidskoding/leetkit/internal/workspace"
)

const (
	starterCode  = "def double(x):\n    pass\n"
	solutionCode = "def double(x):\n    return x * 2\n"
	testCode     = "from starter import double\n\ndef test_double():\n    assert double(2) == 4\n"
)

func fixtureCatalog(t *testing.T) *exercise.Catalog {
	t.Helper()
	fsys := fstest.MapFS{"modules.toml": {Data: []byte(`
[[modules]]
id = "iteration"
title = "Iteration and Unpacking"

[[modules]]
id = "hashing"
title = "Dictionaries and Sets"
`)}}
	moduleDirs := map[string]string{"iteration": "01_iteration", "hashing": "02_hashing"}
	add := func(module, id, title string, order int) {
		dir := fmt.Sprintf("%s/%02d_%s/", moduleDirs[module], order, id)
		fsys[dir+exercise.MetadataFile] = &fstest.MapFile{Data: []byte(fmt.Sprintf(`
id = %q
title = %q
module = %q
difficulty = "beginner"
type = "syntax"
description = "Describe %s."
objective = "Objective for %s."
interview = "Relevant."
hints = ["First hint for %s."]

[[examples]]
input = "double(2)"
output = "4"
explanation = "Two doubled is four."
`, id, title, module, id, id, id))}
		files := map[string]string{
			exercise.StarterFile:     starterCode,
			exercise.TestFile:        testCode,
			exercise.SolutionFile:    solutionCode,
			exercise.ExplanationFile: "# " + id,
		}
		for name, data := range files {
			fsys[dir+name] = &fstest.MapFile{Data: []byte(data)}
		}
	}
	add("iteration", "enumerate-01", "Find an Index", 1)
	add("iteration", "zip-01", "Pair Up", 2)
	add("hashing", "sets-01", "Spot Duplicates", 1)
	cat, err := exercise.Load(fsys)
	if err != nil {
		t.Fatal(err)
	}
	return cat
}

type fixture struct {
	opts  Options
	store *progress.Store
}

// newFixture builds options whose runner reports runErr instead of running
// Python; pass nil and call withPython for real test runs.
func newFixture(t *testing.T, runErr error) fixture {
	t.Helper()
	dir := t.TempDir()
	store, err := progress.Open(filepath.Join(dir, "progress.json"))
	if err != nil {
		t.Fatal(err)
	}
	th, err := theme.Load("")
	if err != nil {
		t.Fatal(err)
	}
	return fixture{
		store: store,
		opts: Options{
			Catalog:   fixtureCatalog(t),
			Store:     store,
			Workspace: workspace.Workspace{Dir: filepath.Join(dir, "workspace")},
			RunnerErr: runErr,
			Theme:     th,
		},
	}
}

// withPython configures a real interpreter with pytest, or skips the test.
func (f fixture) withPython(t *testing.T) fixture {
	t.Helper()
	venv, err := filepath.Abs(filepath.Join("..", "..", ".venv", "bin", "python"))
	if err != nil {
		t.Fatal(err)
	}
	py, err := runner.FindPython(runner.Candidate{Path: venv})
	if err == nil {
		_, err = runner.PytestVersion(context.Background(), py.Path)
	}
	if err != nil {
		t.Skipf("no Python with pytest: %v", err)
	}
	f.opts.RunnerErr = nil
	f.opts.Runner = runner.Runner{Python: py.Path, Timeout: 30 * time.Second}
	return f
}

func (f fixture) writeWork(t *testing.T, id, code string) {
	t.Helper()
	path := f.opts.Workspace.Path(id)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}
}

var size = tea.WindowSizeMsg{Width: 100, Height: 30}

// start creates the model, sizes it, and lets the initial run finish.
func start(t *testing.T, opts Options) Model {
	t.Helper()
	m := New(opts)
	m, _ = update(t, m, size)
	return settle(t, m, m.Init(), idle)
}

func idle(m Model) bool { return !m.run.active && (m.run.result != nil || m.run.err != nil) }

func update(t *testing.T, m Model, msgs ...tea.Msg) (Model, tea.Cmd) {
	t.Helper()
	var cmd tea.Cmd
	for _, msg := range msgs {
		var next tea.Model
		next, cmd = m.Update(msg)
		m = next.(Model)
	}
	return m, cmd
}

// settle executes commands concurrently, as the Bubble Tea runtime does, and
// feeds their messages back until done reports true.
func settle(t *testing.T, m Model, cmd tea.Cmd, done func(Model) bool) Model {
	t.Helper()
	msgs := make(chan tea.Msg, 64)
	launch := func(c tea.Cmd) {
		if c != nil {
			go func() { msgs <- c() }()
		}
	}
	launch(cmd)
	deadline := time.After(20 * time.Second)
	for !done(m) {
		select {
		case msg := <-msgs:
			switch msg := msg.(type) {
			case nil, spinner.TickMsg:
				continue
			case tea.BatchMsg:
				for _, c := range msg {
					launch(c)
				}
				continue
			}
			var c tea.Cmd
			m, c = update(t, m, msg)
			launch(c)
		case <-deadline:
			t.Fatalf("model did not settle; view:\n%s", plain(m))
		}
	}
	return m
}

func press(s string) tea.KeyPressMsg {
	switch s {
	case "down":
		return tea.KeyPressMsg{Code: tea.KeyDown}
	case "left":
		return tea.KeyPressMsg{Code: tea.KeyLeft}
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	case "esc":
		return tea.KeyPressMsg{Code: tea.KeyEscape}
	case "backspace":
		return tea.KeyPressMsg{Code: tea.KeyBackspace}
	case "ctrl+c":
		return tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}
	}
	r := []rune(s)[0]
	return tea.KeyPressMsg{Code: r, Text: s}
}

func isQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

func plain(m Model) string { return ansi.Strip(m.View().Content) }
