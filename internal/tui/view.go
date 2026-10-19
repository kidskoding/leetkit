package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anirudhkonidala/leetkit/internal/exercise"
)

const (
	listMinWidth = 34
	barMinWidth  = 10
	barMaxWidth  = 24
)

// layout derives component sizes from the window size. Panels have a
// rounded border and one cell of horizontal padding; lipgloss widths include both.
func (m *Model) layout() {
	m.help.SetWidth(m.width)
	listWidth, bodyHeight := m.paneSizes()
	m.list.SetSize(listWidth-4, bodyHeight-2)
}

func (m Model) paneSizes() (listWidth, bodyHeight int) {
	listWidth = max(listMinWidth, m.width*2/5)
	bodyHeight = m.height - lipgloss.Height(m.header()) - lipgloss.Height(m.footer())
	return listWidth, bodyHeight
}

func (m Model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	v.WindowTitle = "LeetKit"
	return v
}

func (m Model) render() string {
	if m.width == 0 {
		return ""
	}
	if m.width < minWidth || m.height < minHeight {
		msg := fmt.Sprintf("Terminal too small: need %d×%d, have %d×%d.\nEnlarge the window or press q to quit.",
			minWidth, minHeight, m.width, m.height)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.th.Warn.Render(msg))
	}
	listWidth, bodyHeight := m.paneSizes()
	left := m.panel(m.list.View(), listWidth, bodyHeight, true)
	right := m.panel(m.preview(m.width-listWidth-4, bodyHeight-2), m.width-listWidth, bodyHeight, false)
	return lipgloss.JoinVertical(lipgloss.Left, m.header(), lipgloss.JoinHorizontal(lipgloss.Top, left, right), m.footer())
}

func (m Model) panel(content string, width, height int, focused bool) string {
	border := m.th.Palette.Muted
	if focused {
		border = m.th.Palette.Accent
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Width(width).
		Height(height).
		MaxHeight(height).
		Padding(0, 1).
		Render(content)
}

func (m Model) header() string {
	all := m.opts.Catalog.All()
	ids := make([]string, len(all))
	for i, ex := range all {
		ids[i] = ex.ID
	}
	done := m.opts.Store.CountCompleted(ids)
	ratio := 0.0
	if len(all) > 0 {
		ratio = float64(done) / float64(len(all))
	}
	count := m.th.Text.Render(fmt.Sprintf("%d/%d", done, len(all)))
	title := m.th.Heading.Render("LEETKIT")
	subtitle := m.th.Faint.Render("  Python interview toolkit")
	// Drop the subtitle before shrinking the progress bar on narrow windows.
	room := m.width - 2 - lipgloss.Width(title) - lipgloss.Width(count) - 2
	if room-lipgloss.Width(subtitle) >= barMinWidth {
		title += subtitle
		room -= lipgloss.Width(subtitle)
	}
	bar := m.bar
	bar.SetWidth(min(max(room, 0), barMaxWidth))
	right := bar.ViewAs(ratio) + " " + count
	gap := max(m.width-2-lipgloss.Width(title)-lipgloss.Width(right), 1)
	return " " + title + strings.Repeat(" ", gap) + right + " "
}

func (m Model) footer() string {
	if m.notice != "" {
		return " " + m.th.Warn.Render(m.notice)
	}
	return " " + m.help.View(m.keys)
}

func (m Model) preview(width, height int) string {
	ex, ok := m.selected()
	if !ok {
		return m.th.Faint.Render("No exercise selected.")
	}
	var b strings.Builder
	b.WriteString(m.th.ID.Bold(true).Render(ex.ID) + m.th.Faint.Render(" · ") + m.th.Text.Bold(true).Render(ex.Title) + "\n")
	b.WriteString(m.th.Faint.Render(m.moduleTitle(ex)+" · "+ex.Difficulty+" · "+ex.TypeLabel()) + "\n")
	section := func(title, body string) {
		b.WriteString("\n" + m.th.Heading.Render(title) + "\n" + m.th.Text.Render(strings.TrimSpace(body)) + "\n")
	}
	section("Objective", ex.Objective)
	section("Problem", ex.Description)
	if len(ex.Examples) > 0 {
		var lines []string
		for _, e := range ex.Examples {
			lines = append(lines, m.th.ID.Render(e.Input)+m.th.Faint.Render(" → ")+m.th.Text.Render(e.Output))
		}
		b.WriteString("\n" + m.th.Heading.Render("Examples") + "\n" + strings.Join(lines, "\n") + "\n")
	}
	return lipgloss.NewStyle().Width(width).MaxHeight(height).Render(b.String())
}

func (m Model) moduleTitle(ex *exercise.Exercise) string {
	for _, mod := range m.opts.Catalog.Modules {
		if mod.ID == ex.Module {
			return mod.Title
		}
	}
	return ex.Module
}
