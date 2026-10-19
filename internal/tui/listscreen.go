package tui

import (
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/kidskoding/leetkit/internal/exercise"
	"github.com/kidskoding/leetkit/internal/theme"
)

const listMinWidth = 34

func newList(th theme.Theme, exercises []*exercise.Exercise) list.Model {
	idWidth := 0
	for _, ex := range exercises {
		idWidth = max(idWidth, lipgloss.Width(ex.ID))
	}
	l := list.New(nil, delegate{th: th, idWidth: idWidth}, 0, 0)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetStatusBarItemName("exercise", "exercises")
	l.KeyMap.Quit.SetEnabled(false)
	l.KeyMap.ForceQuit.SetEnabled(false)
	l.KeyMap.ShowFullHelp.SetEnabled(false)
	l.KeyMap.CloseFullHelp.SetEnabled(false)
	// h and l belong to screen navigation, so paging keeps only page keys.
	l.KeyMap.PrevPage = key.NewBinding(key.WithKeys("pgup", "b", "u"))
	l.KeyMap.NextPage = key.NewBinding(key.WithKeys("pgdown", "f", "d"))
	styles := list.DefaultStyles(true)
	styles.StatusBar = th.Faint.PaddingLeft(2)
	styles.StatusBarActiveFilter = th.Text
	styles.StatusBarFilterCount = th.Faint
	styles.NoItems = th.Faint.PaddingLeft(2)
	styles.Filter.Focused.Prompt = th.Heading
	styles.Filter.Blurred.Prompt = th.Faint
	styles.DefaultFilterCharacterMatch = th.Heading.Underline(true)
	l.Styles = styles
	return l
}

func (m Model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok && m.confirmReset != "" {
		id := m.confirmReset
		m.confirmReset = ""
		if k.String() == "y" {
			return m, m.reset(id)
		}
		return m, nil
	}
	if k, ok := msg.(tea.KeyPressMsg); ok && !m.list.SettingFilter() {
		switch {
		case key.Matches(k, m.listKeys.Reset):
			if i := m.selectedIndex(); i >= 0 {
				m.confirmReset = m.exercises[i].ID
			}
			return m, nil
		case key.Matches(k, m.listKeys.Open):
			if i := m.selectedIndex(); i >= 0 {
				m.screen = screenWatch
				if i != m.current {
					m.setCurrent(i)
					return m, m.startRun()
				}
			}
			return m, nil
		// Esc clears an applied filter before leaving the list.
		case key.Matches(k, m.listKeys.Back) && !(k.String() == "esc" && m.list.FilterState() == list.FilterApplied):
			m.screen = screenWatch
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) selectedIndex() int {
	it, ok := m.list.SelectedItem().(item)
	if !ok {
		return -1
	}
	for i, ex := range m.exercises {
		if ex.ID == it.ex.ID {
			return i
		}
	}
	return -1
}

// Panels have a rounded border and one cell of horizontal padding; lipgloss
// widths include both.
func (m Model) listPaneSizes() (listWidth, bodyHeight int) {
	return max(listMinWidth, m.width*2/5), m.height - 2
}

func (m *Model) layoutList() {
	listWidth, bodyHeight := m.listPaneSizes()
	m.list.SetSize(listWidth-4, bodyHeight-2)
}

func (m Model) listView() string {
	listWidth, bodyHeight := m.listPaneSizes()
	header := " " + m.th.Heading.Render("Exercises") + "  " + m.progressBar()
	footer := " " + m.help.View(m.listKeys)
	switch {
	case m.confirmReset != "":
		footer = " " + m.th.Warn.Render("Reset "+m.confirmReset+" to the starter code? Your changes will be lost.") +
			" " + m.th.Heading.Render("y") + m.th.Faint.Render("/") + m.th.Heading.Render("n")
	case m.notice != "":
		footer = " " + m.th.Warn.Render(m.notice)
	}
	footer = lipgloss.NewStyle().MaxWidth(m.width).Render(footer)
	left := m.panel(m.list.View(), listWidth, bodyHeight, true)
	right := m.panel(m.preview(m.width-listWidth-4, bodyHeight-2), m.width-listWidth, bodyHeight, false)
	return lipgloss.JoinVertical(lipgloss.Left, header, lipgloss.JoinHorizontal(lipgloss.Top, left, right), footer)
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

func (m Model) preview(width, height int) string {
	i := m.selectedIndex()
	if i < 0 {
		return m.th.Faint.Render("No exercise selected.")
	}
	ex := m.exercises[i]
	th := m.th
	var b strings.Builder
	b.WriteString(th.ID.Bold(true).Render(ex.ID) + th.Faint.Render(" · ") + th.Text.Bold(true).Render(ex.Title) + "\n")
	b.WriteString(th.Faint.Render(m.moduleTitle(ex)+" · "+ex.Difficulty+" · "+ex.TypeLabel()) + "\n")
	b.WriteString("\n" + th.Heading.Render("Objective") + "\n" + th.Text.Render(strings.TrimSpace(ex.Objective)) + "\n")
	b.WriteString("\n" + th.Heading.Render("Problem") + "\n" + th.Text.Render(strings.TrimSpace(ex.Description)) + "\n")
	return lipgloss.NewStyle().Width(width).MaxHeight(height).Render(b.String())
}
