package tui

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/kidskoding/leetkit/internal/editor"
	"github.com/kidskoding/leetkit/internal/runner"
)

func assertContains(t *testing.T, m Model, wants ...string) {
	t.Helper()
	out := plain(m)
	for _, w := range wants {
		if !strings.Contains(out, w) {
			t.Errorf("view missing %q:\n%s", w, out)
		}
	}
}

func TestOpensOnFirstPendingAndRunsTests(t *testing.T) {
	f := newFixture(t, nil).withPython(t)
	if _, err := f.store.Record("enumerate-01", true); err != nil {
		t.Fatal(err)
	}
	m := start(t, f.opts)
	assertContains(t, m, "zip-01 · Pair Up", "Describe zip-01.", "double(2) → 4", "── Tests", "✗ 0 passed, 1 failed",
		"Progress", "1/3", "Current", "zip-01/starter.py")
	if e, _ := f.store.Entry("zip-01"); e.Attempts != 1 {
		t.Errorf("run not recorded: %+v", e)
	}
}

func TestResumesSavedExercise(t *testing.T) {
	f := newFixture(t, runner.ErrNoPython)
	if err := f.store.SetCurrent("sets-01"); err != nil {
		t.Fatal(err)
	}
	m := start(t, f.opts)
	assertContains(t, m, "sets-01 · Spot Duplicates")
}

func TestRunnerErrorIsShown(t *testing.T) {
	m := start(t, newFixture(t, runner.ErrNoPython).opts)
	assertContains(t, m, "No Python interpreter found", "leetkit doctor")
}

func TestPassingRunCompletesExercise(t *testing.T) {
	f := newFixture(t, nil).withPython(t)
	f.writeWork(t, "enumerate-01", solutionCode)
	m := start(t, f.opts)
	assertContains(t, m, "✓ All 1 tests passed", "✓ Exercise complete!", "Objective for enumerate-01.", "n next")
	if !f.store.Completed("enumerate-01") {
		t.Error("completion not saved")
	}

	m, cmd := update(t, m, press("n"))
	m = settle(t, m, cmd, idle)
	assertContains(t, m, "zip-01 · Pair Up")
	if f.store.Current() != "zip-01" {
		t.Errorf("current = %q, want zip-01", f.store.Current())
	}
}

func TestNextRequiresCompletion(t *testing.T) {
	m := start(t, newFixture(t, runner.ErrNoPython).opts)
	m, _ = update(t, m, press("n"))
	assertContains(t, m, "Solve this exercise first", "enumerate-01 · Find an Index")
}

func TestNextSkipsCompletedAndReportsAllDone(t *testing.T) {
	f := newFixture(t, runner.ErrNoPython)
	for _, id := range []string{"enumerate-01", "zip-01"} {
		if _, err := f.store.Record(id, true); err != nil {
			t.Fatal(err)
		}
	}
	if err := f.store.SetCurrent("enumerate-01"); err != nil {
		t.Fatal(err)
	}
	m := start(t, f.opts)
	m, cmd := update(t, m, press("n"))
	m = settle(t, m, cmd, idle)
	assertContains(t, m, "sets-01 · Spot Duplicates")

	if _, err := f.store.Record("sets-01", true); err != nil {
		t.Fatal(err)
	}
	m, _ = update(t, m, press("n"))
	assertContains(t, m, "Every exercise is complete")
}

func TestHints(t *testing.T) {
	m := start(t, newFixture(t, runner.ErrNoPython).opts)
	assertContains(t, m, "t hint 0/1")
	m, _ = update(t, m, press("t"))
	assertContains(t, m, "── Hint 1/1", "First hint for enumerate-01.")
	m, _ = update(t, m, press("t"))
	assertContains(t, m, "No more hints")
}

func TestListSelectsExercise(t *testing.T) {
	f := newFixture(t, runner.ErrNoPython)
	m := start(t, f.opts)
	m, _ = update(t, m, press("h"))
	assertContains(t, m, "Exercises", "○ enumerate-01", "Objective for enumerate-01.")
	m, _ = update(t, m, press("j"), press("down"))
	assertContains(t, m, "Objective for sets-01.")
	m, cmd := update(t, m, press("l"))
	m = settle(t, m, cmd, idle)
	if m.screen != screenWatch {
		t.Fatal("enter did not return to the watch screen")
	}
	assertContains(t, m, "sets-01 · Spot Duplicates")
	if f.store.Current() != "sets-01" {
		t.Errorf("current = %q", f.store.Current())
	}
}

func TestListBackKeysDoNotQuit(t *testing.T) {
	m := start(t, newFixture(t, runner.ErrNoPython).opts)
	for _, open := range []string{"h", "left", "esc"} {
		m, _ = update(t, m, press(open))
		if m.screen != screenList {
			t.Fatalf("%s did not open the list", open)
		}
		m, _ = update(t, m, press("esc"))
	}
	for _, k := range []string{"esc", "q"} {
		var cmd tea.Cmd
		m, _ = update(t, m, press("h"))
		m, cmd = update(t, m, press(k))
		if isQuit(cmd) || m.screen != screenWatch {
			t.Errorf("%s in list: quit=%v screen=%v", k, isQuit(cmd), m.screen)
		}
	}
}

func TestListSearch(t *testing.T) {
	m := start(t, newFixture(t, runner.ErrNoPython).opts)
	m, _ = update(t, m, press("h"), press("/"))
	m, cmd := update(t, m, press("q"))
	if isQuit(cmd) || m.screen != screenList {
		t.Fatal("typing q into search left the list")
	}
	m, _ = update(t, m, press("backspace"))
	for _, k := range []string{"s", "e", "t"} {
		m, cmd = update(t, m, press(k))
	}
	m = settle(t, m, cmd, func(m Model) bool { return m.selectedIndex() == 2 })
	m, cmd = update(t, m, press("enter"), press("enter"))
	m = settle(t, m, cmd, idle)
	assertContains(t, m, "sets-01 · Spot Duplicates")
}

func TestQuit(t *testing.T) {
	m := start(t, newFixture(t, runner.ErrNoPython).opts)
	for _, k := range []string{"q", "ctrl+c"} {
		if _, cmd := update(t, m, press(k)); !isQuit(cmd) {
			t.Errorf("%s did not quit", k)
		}
	}
}

func TestStaleRunResultIgnored(t *testing.T) {
	m := start(t, newFixture(t, runner.ErrNoPython).opts)
	stale := runDoneMsg{seq: m.run.seq - 1, id: "enumerate-01", res: &runner.Result{Status: runner.StatusPassed}}
	m, _ = update(t, m, stale)
	if m.opts.Store.Completed("enumerate-01") {
		t.Error("stale result was recorded")
	}
}

func TestTooSmall(t *testing.T) {
	m := start(t, newFixture(t, runner.ErrNoPython).opts)
	m, _ = update(t, m, tea.WindowSizeMsg{Width: 40, Height: 10})
	assertContains(t, m, "Terminal too small")
	m, _ = update(t, m, size)
	if strings.Contains(plain(m), "Terminal too small") {
		t.Error("message persisted after resize")
	}
}

func TestViewFitsWindow(t *testing.T) {
	m := start(t, newFixture(t, runner.ErrNoPython).opts)
	m, _ = update(t, m, press("t"), press("?"))
	for _, scr := range []screen{screenWatch, screenList, screenProgress} {
		m.screen = scr
		for _, sz := range []tea.WindowSizeMsg{{Width: 60, Height: 12}, size, {Width: 200, Height: 50}} {
			m, _ = update(t, m, sz)
			lines := strings.Split(m.View().Content, "\n")
			if len(lines) > sz.Height {
				t.Errorf("screen %d at %dx%d: %d lines", scr, sz.Width, sz.Height, len(lines))
			}
			for i, l := range lines {
				if w := lipgloss.Width(l); w > sz.Width {
					t.Errorf("screen %d at %dx%d: line %d is %d wide", scr, sz.Width, sz.Height, i, w)
					break
				}
			}
		}
	}
}

func TestEditWithoutEditorShowsFix(t *testing.T) {
	f := newFixture(t, runner.ErrNoPython)
	f.opts.EditorErr = editor.ErrNoEditor
	m := start(t, f.opts)
	m, cmd := update(t, m, press("e"))
	if cmd != nil {
		t.Error("edit without an editor returned a command")
	}
	assertContains(t, m, "No editor", "Set $EDITOR")
}

func TestEditOpensWorkingCopy(t *testing.T) {
	f := newFixture(t, runner.ErrNoPython)
	f.opts.Editor = editor.Editor{Argv: []string{"true"}, Source: "test"}
	m := start(t, f.opts)
	m, cmd := update(t, m, press("e"))
	if cmd == nil {
		t.Fatal("e returned no command")
	}
	if _, err := os.Stat(f.opts.Workspace.Path("enumerate-01")); err != nil {
		t.Errorf("working copy not created before opening the editor: %v", err)
	}
}

func TestEditorExitRerunsTests(t *testing.T) {
	f := newFixture(t, nil).withPython(t)
	m := start(t, f.opts)
	f.writeWork(t, "enumerate-01", solutionCode)
	m, cmd := update(t, m, editorDoneMsg{})
	m = settle(t, m, cmd, idle)
	assertContains(t, m, "✓ Exercise complete!")
}

func TestEditorStartFailure(t *testing.T) {
	f := newFixture(t, runner.ErrNoPython)
	f.opts.Editor = editor.Editor{Argv: []string{"missing-editor"}}
	m := start(t, f.opts)
	m, cmd := update(t, m, editorDoneMsg{err: exec.ErrNotFound})
	if cmd != nil {
		t.Error("a failed editor launch should not rerun tests")
	}
	assertContains(t, m, "Could not start missing-editor")
}

func TestSaveTriggersRerun(t *testing.T) {
	f := newFixture(t, nil).withPython(t)
	m := start(t, f.opts)
	path := f.opts.Workspace.Path("enumerate-01")
	base := time.Now().Add(-time.Minute)

	m, cmd := update(t, m, fileCheckMsg{path: path, mod: base, ok: true})
	if m.run.active {
		t.Fatal("first observation started a run")
	}
	m, _ = update(t, m, fileCheckMsg{path: path, mod: base, ok: true})
	if m.run.active {
		t.Fatal("unchanged file started a run")
	}

	f.writeWork(t, "enumerate-01", solutionCode)
	m, cmd = update(t, m, fileCheckMsg{path: path, mod: base.Add(time.Second), ok: true})
	if !m.run.active {
		t.Fatal("changed file did not start a run")
	}
	m = settle(t, m, cmd, idle)
	assertContains(t, m, "✓ Exercise complete!")
}

func TestWatchIgnoresOtherFilesAndEditing(t *testing.T) {
	f := newFixture(t, runner.ErrNoPython)
	m := start(t, f.opts)
	path := f.opts.Workspace.Path("enumerate-01")
	t0 := time.Now()
	m, _ = update(t, m, fileCheckMsg{path: path, mod: t0, ok: true})
	seq := m.run.seq

	m, _ = update(t, m, fileCheckMsg{path: "/elsewhere/starter.py", mod: t0.Add(time.Second), ok: true})
	m.editing = true
	m, _ = update(t, m, fileCheckMsg{path: path, mod: t0.Add(2 * time.Second), ok: true})
	m.editing = false
	m, _ = update(t, m, fileCheckMsg{path: path, mod: t0.Add(2 * time.Second), ok: true})
	if m.run.seq != seq {
		t.Errorf("runs started: %d", m.run.seq-seq)
	}
}

func TestCheckFileReportsModTime(t *testing.T) {
	f := newFixture(t, runner.ErrNoPython)
	f.writeWork(t, "zip-01", starterCode)
	path := f.opts.Workspace.Path("zip-01")
	msg := checkFile(path)().(fileCheckMsg)
	if !msg.ok || msg.path != path || msg.mod.IsZero() {
		t.Errorf("got %+v", msg)
	}
	if missing := checkFile(path + ".missing")().(fileCheckMsg); missing.ok {
		t.Error("missing file reported ok")
	}
}

func TestResetFromList(t *testing.T) {
	f := newFixture(t, runner.ErrNoPython)
	f.writeWork(t, "zip-01", solutionCode)
	if _, err := f.store.Record("zip-01", true); err != nil {
		t.Fatal(err)
	}
	m := start(t, f.opts)
	m, _ = update(t, m, press("h"), press("j"), press("r"))
	assertContains(t, m, "Reset zip-01 to the starter code?", "y/n")

	m, _ = update(t, m, press("n"))
	if !f.store.Completed("zip-01") {
		t.Fatal("declining the prompt reset progress")
	}

	m, _ = update(t, m, press("r"), press("y"))
	data, err := os.ReadFile(f.opts.Workspace.Path("zip-01"))
	if err != nil || string(data) != starterCode {
		t.Errorf("working copy = %q, %v; want starter", data, err)
	}
	if _, ok := f.store.Entry("zip-01"); ok {
		t.Error("progress not cleared")
	}
	assertContains(t, m, "Reset zip-01 to the starter code.", "○ zip-01")
}

func TestExplanationUnlocksWhenSolved(t *testing.T) {
	f := newFixture(t, runner.ErrNoPython)
	m := start(t, f.opts)
	m, _ = update(t, m, press("x"))
	assertContains(t, m, "Solve this exercise to unlock its explanation.")
	if m.explaining {
		t.Fatal("explanation opened before solving")
	}

	if _, err := f.store.Record("enumerate-01", true); err != nil {
		t.Fatal(err)
	}
	m, _ = update(t, m, press("x"))
	assertContains(t, m, "── Explanation · enumerate-01", "x back to exercise")
	m, _ = update(t, m, press("x"))
	assertContains(t, m, "Describe enumerate-01.", "x explanation")
}

func TestProgressScreen(t *testing.T) {
	f := newFixture(t, runner.ErrNoPython)
	for id, passed := range map[string]bool{"enumerate-01": true, "zip-01": false} {
		if _, err := f.store.Record(id, passed); err != nil {
			t.Fatal(err)
		}
	}
	m := start(t, f.opts)
	for _, from := range []screen{screenWatch, screenList} {
		m.screen = from
		m, _ = update(t, m, press("p"))
		if m.screen != screenProgress {
			t.Fatalf("p from screen %d did not open progress", from)
		}
		assertContains(t, m, "Progress", "Overall", "1/3", "Iteration and Unpacking", "1/2",
			"Dictionaries and Sets", "0/1", "1 solved · 2 test runs")
		m, _ = update(t, m, press("esc"))
		if m.screen != from {
			t.Errorf("esc returned to %d, want %d", m.screen, from)
		}
	}
	m, _ = update(t, m, press("p"))
	if _, cmd := update(t, m, press("q")); !isQuit(cmd) {
		t.Error("q did not quit from the progress screen")
	}
}

func TestFullHelpToggle(t *testing.T) {
	m := start(t, newFixture(t, runner.ErrNoPython).opts)
	if strings.Contains(plain(m), "progress") {
		t.Fatal("short help should not list every key")
	}
	m, _ = update(t, m, press("?"))
	assertContains(t, m, "progress", "j/k scroll", "? less")
	m, _ = update(t, m, press("?"))
	if strings.Contains(plain(m), "progress") {
		t.Error("? did not collapse the help")
	}
}

func TestGlobalKeysIgnoredWhileSearching(t *testing.T) {
	m := start(t, newFixture(t, runner.ErrNoPython).opts)
	m, _ = update(t, m, press("h"), press("/"), press("p"))
	if m.screen != screenList {
		t.Error("typing p into search opened the progress screen")
	}
}
