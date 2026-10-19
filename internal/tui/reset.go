package tui

import (
	tea "charm.land/bubbletea/v2"
)

// reset restores the starter code and clears the exercise's progress.
func (m *Model) reset(id string) tea.Cmd {
	ex, ok := m.opts.Catalog.Get(id)
	if !ok {
		return nil
	}
	starter, err := m.opts.Catalog.Starter(ex)
	if err == nil {
		_, err = m.opts.Workspace.Reset(id, starter)
	}
	if err == nil {
		err = m.opts.Store.Reset(id)
	}
	if err != nil {
		m.notice = "Reset failed: " + err.Error()
		return nil
	}
	m.notice = "Reset " + id + " to the starter code."
	m.refreshItems()
	if id != m.currentExercise().ID {
		return nil
	}
	m.hints = 0
	return m.startRun()
}
