package tui

import (
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

// Polling keeps LeetKit dependency-free and works across editors that save by
// writing a temp file and renaming it, which confuses some inotify watchers.
const watchInterval = 500 * time.Millisecond

type fileCheckMsg struct {
	path string
	mod  time.Time
	ok   bool
}

func checkFile(path string) tea.Cmd {
	return tea.Tick(watchInterval, func(time.Time) tea.Msg {
		info, err := os.Stat(path)
		if err != nil {
			return fileCheckMsg{path: path}
		}
		return fileCheckMsg{path: path, mod: info.ModTime(), ok: true}
	})
}

// onFileCheck reruns tests when the current exercise's file changes on disk,
// such as a save from an editor in another pane. The first observation of a
// file only records a baseline.
func (m *Model) onFileCheck(msg fileCheckMsg) tea.Cmd {
	path := m.opts.Workspace.Path(m.currentExercise().ID)
	var run tea.Cmd
	switch {
	case msg.path != path || !msg.ok:
		m.lastMod = time.Time{}
	case m.editing:
		m.lastMod = msg.mod
	default:
		changed := !m.lastMod.IsZero() && !msg.mod.Equal(m.lastMod)
		m.lastMod = msg.mod
		if changed {
			m.revealResults = true
			run = m.startRun()
		}
	}
	return tea.Batch(run, checkFile(path))
}
