package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

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
	m, _ = update(t, m, press("t"))
	for _, scr := range []screen{screenWatch, screenList} {
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
