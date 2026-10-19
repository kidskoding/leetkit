package tui

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func lipglossWidth(s string) int { return lipgloss.Width(s) }

func plain(m Model) string { return ansi.Strip(m.View().Content) }

// settle runs the commands an update produced and feeds their messages back,
// as the Bubble Tea runtime would. Commands that block, such as cursor blink
// timers, are abandoned after a short wait.
func settle(t *testing.T, m Model, cmd tea.Cmd) Model {
	t.Helper()
	queue := []tea.Cmd{cmd}
	for len(queue) > 0 {
		c := queue[0]
		queue = queue[1:]
		if c == nil {
			continue
		}
		done := make(chan tea.Msg, 1)
		go func() { done <- c() }()
		var msg tea.Msg
		select {
		case msg = <-done:
		case <-time.After(50 * time.Millisecond):
			continue
		}
		if batch, ok := msg.(tea.BatchMsg); ok {
			queue = append(queue, batch...)
			continue
		}
		if msg == nil {
			continue
		}
		next, c2 := m.Update(msg)
		m = next.(Model)
		queue = append(queue, c2)
	}
	return m
}
