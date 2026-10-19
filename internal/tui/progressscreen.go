package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/kidskoding/leetkit/internal/report"
)

func (m Model) updateProgress(msg tea.Msg) (tea.Model, tea.Cmd) {
	k, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch {
	case key.Matches(k, m.progressKeys.Quit):
		m.run.cancel()
		return m, tea.Quit
	case key.Matches(k, m.progressKeys.Back):
		m.screen = m.prevScreen
	}
	return m, nil
}

func (m Model) progressView() string {
	th := m.th
	labelWidth := lipgloss.Width("Overall")
	for _, mod := range m.opts.Catalog.Modules {
		labelWidth = max(labelWidth, lipgloss.Width(mod.Title))
	}
	barWidth := min(max(m.width-labelWidth-16, 10), 40)
	row := func(label string, done, total int, style lipgloss.Style) string {
		pct := 0
		if total > 0 {
			pct = done * 100 / total
		}
		return " " + style.Render(pad(label, labelWidth)) + "  " + report.Meter(th, done, total, barWidth) + "  " +
			th.Text.Render(fmt.Sprintf("%d/%d", done, total)) + th.Faint.Render(fmt.Sprintf("  %d%%", pct))
	}

	var all []string
	lines := []string{" " + th.Heading.Render("Progress"), ""}
	var modules []string
	attempts := 0
	for _, mod := range m.opts.Catalog.Modules {
		ids := make([]string, len(mod.Exercises))
		for i, ex := range mod.Exercises {
			ids[i] = ex.ID
			if e, ok := m.opts.Store.Entry(ex.ID); ok {
				attempts += e.Attempts
			}
		}
		all = append(all, ids...)
		modules = append(modules, row(mod.Title, m.opts.Store.CountCompleted(ids), len(ids), th.Text))
	}
	done := m.opts.Store.CountCompleted(all)
	lines = append(lines, row("Overall", done, len(all), th.Text.Bold(true)), "")
	lines = append(lines, modules...)
	lines = append(lines, "", " "+th.Faint.Render(fmt.Sprintf("%d solved · %d test runs", done, attempts)))

	body := lipgloss.NewStyle().MaxWidth(m.width).Render(strings.Join(lines, "\n"))
	footer := " " + m.help.View(m.progressKeys)
	gap := max(m.height-lipgloss.Height(body)-lipgloss.Height(footer), 0)
	return body + strings.Repeat("\n", gap+1) + footer
}

func pad(s string, width int) string {
	return s + strings.Repeat(" ", max(width-lipgloss.Width(s), 0))
}
