package tui

import (
	"errors"
	"os/exec"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/kidskoding/leetkit/internal/exercise"
)

type editorDoneMsg struct{ err error }

// edit hands the terminal to the learner's editor. Tests rerun when it exits.
func (m *Model) edit() tea.Cmd {
	if m.opts.EditorErr != nil {
		m.notice = "No editor: " + m.opts.EditorErr.Error() + ". Set $EDITOR or `editor` in config.toml."
		return nil
	}
	ex := m.currentExercise()
	starter, err := m.opts.Catalog.ReadFile(ex, exercise.StarterFile)
	if err != nil {
		m.notice = "Could not read the starter code: " + err.Error()
		return nil
	}
	path, err := m.opts.Workspace.Ensure(ex.ID, starter)
	if err != nil {
		m.notice = err.Error()
		return nil
	}
	m.run.cancel()
	m.editing = true
	return tea.ExecProcess(m.opts.Editor.Command(path), func(err error) tea.Msg { return editorDoneMsg{err} })
}

func (m *Model) finishEdit(msg editorDoneMsg) tea.Cmd {
	m.editing = false
	// The run below covers the edit; a zero baseline keeps the file watcher
	// from starting a second run for the same change.
	m.lastMod = time.Time{}
	var exitErr *exec.ExitError
	switch {
	case msg.err == nil:
	case errors.As(msg.err, &exitErr):
		// The editor ran but exited non-zero (for example vim's :cq); the file may still have changed.
		m.notice = "Editor exited with " + exitErr.String() + "."
	default:
		m.notice = "Could not start " + m.opts.Editor.String() + ": " + msg.err.Error()
		return nil
	}
	return m.startRun()
}
