package tui

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	tea "charm.land/bubbletea/v2"

	"github.com/anirudhkonidala/leetkit/internal/exercise"
	"github.com/anirudhkonidala/leetkit/internal/progress"
	"github.com/anirudhkonidala/leetkit/internal/theme"
)

func fixtureCatalog(t *testing.T) *exercise.Catalog {
	t.Helper()
	fsys := fstest.MapFS{"modules.toml": {Data: []byte(`
[[modules]]
id = "iteration"
title = "Iteration and Unpacking"
order = 1

[[modules]]
id = "hashing"
title = "Dictionaries and Sets"
order = 2
`)}}
	add := func(module, id, title string, order int) {
		dir := module + "/" + id + "/"
		fsys[dir+exercise.MetadataFile] = &fstest.MapFile{Data: []byte(fmt.Sprintf(`
id = %q
title = %q
module = %q
difficulty = "beginner"
type = "syntax"
order = %d
description = "Describe %s."
objective = "Objective for %s."
interview = "Relevant."
hints = ["Hint."]

[[examples]]
input = "f(1)"
output = "2"
`, id, title, module, order, id, id))}
		for _, name := range []string{exercise.StarterFile, exercise.TestFile, exercise.SolutionFile, exercise.ExplanationFile} {
			fsys[dir+name] = &fstest.MapFile{Data: []byte("# " + name)}
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

func newModel(t *testing.T) (Model, *progress.Store) {
	t.Helper()
	store, err := progress.Open(filepath.Join(t.TempDir(), "progress.json"))
	if err != nil {
		t.Fatal(err)
	}
	th, err := theme.Load("")
	if err != nil {
		t.Fatal(err)
	}
	return New(Options{Catalog: fixtureCatalog(t), Store: store, Theme: th}), store
}

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

func press(s string) tea.KeyPressMsg {
	switch s {
	case "down":
		return tea.KeyPressMsg{Code: tea.KeyDown}
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	case "esc":
		return tea.KeyPressMsg{Code: tea.KeyEscape}
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

var size = tea.WindowSizeMsg{Width: 100, Height: 30}

func TestRendersListAndPreview(t *testing.T) {
	m, _ := newModel(t)
	m, _ = update(t, m, size)
	out := m.View().Content
	for _, want := range []string{"LEETKIT", "0/3", "enumerate-01", "zip-01", "sets-01", "Objective for enumerate-01", "f(1)"} {
		if !strings.Contains(out, want) {
			t.Errorf("view missing %q", want)
		}
	}
	if !m.View().AltScreen {
		t.Error("view should use the alternate screen")
	}
}

func TestNavigationUpdatesPreview(t *testing.T) {
	m, _ := newModel(t)
	m, _ = update(t, m, size, press("j"), press("down"))
	if ex, _ := m.selected(); ex.ID != "sets-01" {
		t.Fatalf("selected %s, want sets-01", ex.ID)
	}
	if !strings.Contains(m.View().Content, "Objective for sets-01") {
		t.Error("preview did not follow the cursor")
	}
}

func TestStatusMarksAndProgress(t *testing.T) {
	m, store := newModel(t)
	if _, err := store.Record("enumerate-01", true); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Record("zip-01", false); err != nil {
		t.Fatal(err)
	}
	m.refreshItems()
	m, _ = update(t, m, size)
	out := plain(m)
	for _, want := range []string{"✓ enumerate-01", "✗ zip-01", "○ sets-01", "1/3"} {
		if !strings.Contains(out, want) {
			t.Errorf("view missing %q:\n%s", want, out)
		}
	}
}

func TestQuit(t *testing.T) {
	m, _ := newModel(t)
	for _, k := range []string{"q", "ctrl+c"} {
		if _, cmd := update(t, m, size, press(k)); !isQuit(cmd) {
			t.Errorf("%s did not quit", k)
		}
	}
}

func TestSearchKeepsFocus(t *testing.T) {
	m, _ := newModel(t)
	m, cmd := update(t, m, size, press("/"), press("q"))
	if isQuit(cmd) {
		t.Fatal("typing q into the search field quit the app")
	}
	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyBackspace})
	for _, k := range []string{"s", "e", "t"} {
		m, cmd = update(t, m, press(k))
		m = settle(t, m, cmd)
	}
	m, _ = update(t, m, press("enter"))
	if ex, _ := m.selected(); ex.ID != "sets-01" {
		t.Errorf("search selected %v, want sets-01", ex.ID)
	}
}

func TestTooSmall(t *testing.T) {
	m, _ := newModel(t)
	m, _ = update(t, m, tea.WindowSizeMsg{Width: 40, Height: 10})
	if !strings.Contains(m.View().Content, "Terminal too small") {
		t.Error("missing too-small message")
	}
	m, _ = update(t, m, size)
	if strings.Contains(m.View().Content, "Terminal too small") {
		t.Error("message persisted after resize")
	}
}

func TestNoticeClearsOnKeyPress(t *testing.T) {
	m, _ := newModel(t)
	m.notice = "progress file was corrupt"
	m, _ = update(t, m, size)
	if !strings.Contains(m.View().Content, "progress file was corrupt") {
		t.Fatal("notice not shown")
	}
	m, _ = update(t, m, press("j"))
	if strings.Contains(m.View().Content, "progress file was corrupt") {
		t.Error("notice not cleared")
	}
}

func TestViewFitsWindow(t *testing.T) {
	m, _ := newModel(t)
	for _, sz := range []tea.WindowSizeMsg{{Width: 60, Height: 12}, size, {Width: 200, Height: 50}} {
		m, _ = update(t, m, sz)
		out := m.View().Content
		lines := strings.Split(out, "\n")
		if len(lines) > sz.Height {
			t.Errorf("%dx%d: %d lines", sz.Width, sz.Height, len(lines))
		}
		for i, l := range lines {
			if w := lipglossWidth(l); w > sz.Width {
				t.Errorf("%dx%d: line %d is %d wide", sz.Width, sz.Height, i, w)
				break
			}
		}
	}
}
