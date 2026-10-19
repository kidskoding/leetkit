package tui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/kidskoding/leetkit/internal/exercise"
	"github.com/kidskoding/leetkit/internal/report"
	"github.com/kidskoding/leetkit/internal/runner"
)

// statusLines is the fixed block under the scrolling content: a blank line,
// progress, current file, and the key line.
const statusLines = 4

func (m Model) updateWatch(msg tea.Msg) (tea.Model, tea.Cmd) {
	k, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	keys := m.watchKeys
	switch {
	case key.Matches(k, keys.Quit):
		m.run.cancel()
		return m, tea.Quit
	case key.Matches(k, keys.Edit):
		return m, m.edit()
	case key.Matches(k, keys.Run):
		return m, m.startRun()
	case key.Matches(k, keys.Next):
		return m, m.next()
	case key.Matches(k, keys.Hint):
		m.revealHint()
	case key.Matches(k, keys.Explain):
		m.toggleExplanation()
	case key.Matches(k, keys.List):
		m.screen = screenList
		if m.list.FilterState() == list.Unfiltered {
			m.list.Select(m.current)
		}
	case key.Matches(k, keys.Scroll):
		switch k.String() {
		case "up", "k":
			m.vp.ScrollUp(1)
		case "down", "j":
			m.vp.ScrollDown(1)
		case "pgup":
			m.vp.PageUp()
		case "pgdown":
			m.vp.PageDown()
		}
	}
	return m, nil
}

func (m *Model) next() tea.Cmd {
	if !m.opts.Store.Completed(m.currentExercise().ID) {
		m.notice = "Solve this exercise first, or press h to pick another."
		return nil
	}
	i := m.nextPending(m.current)
	if i < 0 {
		m.notice = "Every exercise is complete. Nice work!"
		return nil
	}
	m.setCurrent(i)
	return m.startRun()
}

func (m *Model) toggleExplanation() {
	if !m.explaining && !m.opts.Store.Completed(m.currentExercise().ID) {
		m.notice = "Solve this exercise to unlock its explanation."
		return
	}
	m.explaining = !m.explaining
	m.refreshWatch()
	m.vp.GotoTop()
}

func (m *Model) setCurrent(i int) {
	m.current, m.hints, m.explaining = i, 0, false
	if err := m.opts.Store.SetCurrent(m.exercises[i].ID); err != nil {
		m.notice = "Could not save progress: " + err.Error()
	}
	m.refreshWatch()
	m.vp.GotoTop()
}

func (m *Model) revealHint() {
	total := len(m.currentExercise().Hints)
	if m.hints >= total {
		m.notice = "No more hints for this exercise."
		return
	}
	m.hints++
	m.refreshWatch()
	m.vp.GotoBottom()
}

func (m *Model) layoutWatch() {
	m.vp.SetWidth(max(m.width-2, 0))
	m.vp.SetHeight(max(m.height-statusLines, 0))
	m.refreshWatch()
}

// refreshWatch rebuilds the scrolling content; call it after any state change.
func (m *Model) refreshWatch() {
	if m.width == 0 || m.currentExercise() == nil {
		return
	}
	m.watchKeys.solved = m.opts.Store.Completed(m.currentExercise().ID)
	desc := "explanation"
	if m.explaining {
		desc = "back to exercise"
	}
	m.watchKeys.Explain.SetHelp("x", desc)
	m.watchKeys.Hint.SetHelp("t", fmt.Sprintf("hint %d/%d", m.hints, len(m.currentExercise().Hints)))
	m.vp.SetContent(m.watchContent(m.vp.Width()))
}

func (m Model) watchView() string {
	content := lipgloss.NewStyle().PaddingLeft(1).Render(m.vp.View())
	footer := " " + m.help.View(m.watchKeys)
	if m.notice != "" {
		footer = " " + m.th.Warn.Render(m.notice)
	}
	path := truncateLeft(shortenHome(m.opts.Workspace.Path(m.currentExercise().ID)), m.width-11)
	status := strings.Join([]string{
		"",
		" " + m.th.Faint.Render("Progress ") + m.progressBar(),
		" " + m.th.Faint.Render("Current  ") + m.th.ID.Render(path),
		lipgloss.NewStyle().MaxWidth(m.width).Render(footer),
	}, "\n")
	return lipgloss.JoinVertical(lipgloss.Left, content, status)
}

func (m Model) watchContent(width int) string {
	ex := m.currentExercise()
	th := m.th
	if m.explaining {
		doc, err := m.opts.Catalog.ReadFile(ex, exercise.ExplanationFile)
		if err != nil {
			return th.Fail.Render("✗ Could not read the explanation: " + err.Error())
		}
		return m.rule("Explanation · "+ex.ID, width) + "\n\n" + renderMarkdown(th, string(doc), width)
	}
	var b strings.Builder
	b.WriteString(th.ID.Bold(true).Render(ex.ID) + th.Faint.Render(" · ") + th.Text.Bold(true).Render(ex.Title) + "\n")
	b.WriteString(th.Faint.Render(m.moduleTitle(ex)+" · "+ex.Difficulty+" · "+ex.TypeLabel()) + "\n\n")
	b.WriteString(th.Text.Render(strings.TrimSpace(ex.Description)) + "\n")
	for i, e := range ex.Examples {
		label := "         "
		if i == 0 {
			label = "Example  "
		}
		b.WriteString("\n" + th.Faint.Render(label) + th.ID.Render(e.Input) + th.Faint.Render(" → ") + th.Text.Render(e.Output))
	}
	if len(ex.Examples) > 0 {
		b.WriteString("\n")
	}

	b.WriteString("\n" + m.rule("Tests", width) + "\n")
	switch {
	case m.run.active:
		b.WriteString(m.spinner.View() + th.Faint.Render(" Running tests…") + "\n")
	case m.run.err != nil:
		b.WriteString(th.Fail.Render("✗ "+m.runErrorText()) + "\n")
	case m.run.result != nil:
		b.WriteString(report.Result(th, m.run.result, false))
		if m.run.result.Status == runner.StatusPassed {
			b.WriteString("\n" + th.Pass.Render("✓ Exercise complete!") + " " + th.Text.Render(ex.Objective) + "\n")
			b.WriteString(th.Faint.Render("Press ") + th.Heading.Render("n") + th.Faint.Render(" for the next exercise or ") +
				th.Heading.Render("x") + th.Faint.Render(" to read the explanation.") + "\n")
		}
	}

	for i := range m.hints {
		b.WriteString("\n" + m.rule(fmt.Sprintf("Hint %d/%d", i+1, len(ex.Hints)), width) + "\n")
		b.WriteString(th.Text.Render(ex.Hints[i]) + "\n")
	}
	return lipgloss.NewStyle().Width(width).Render(b.String())
}

func (m Model) runErrorText() string {
	err := m.run.err
	switch {
	case errors.Is(err, runner.ErrPytestMissing):
		return fmt.Sprintf("pytest is not installed for %s. Quit and run `leetkit doctor` for setup steps.", m.opts.Runner.Python)
	case errors.Is(err, runner.ErrNoPython):
		return "No Python interpreter found. Quit and run `leetkit doctor` for setup steps."
	}
	return err.Error()
}

func (m Model) rule(title string, width int) string {
	label := "── " + title + " "
	return m.th.Heading.Render(label) + m.th.Faint.Render(strings.Repeat("─", max(width-lipgloss.Width(label), 0)))
}

func (m Model) progressBar() string {
	ids := make([]string, len(m.exercises))
	for i, ex := range m.exercises {
		ids[i] = ex.ID
	}
	done := m.opts.Store.CountCompleted(ids)
	bar := m.bar
	bar.SetWidth(min(max(m.width-24, 10), 40))
	ratio := 0.0
	if len(ids) > 0 {
		ratio = float64(done) / float64(len(ids))
	}
	return bar.ViewAs(ratio) + " " + m.th.Text.Render(fmt.Sprintf("%d/%d", done, len(ids)))
}

// truncateLeft keeps the end of s, which for paths is the informative part.
func truncateLeft(s string, width int) string {
	r := []rune(s)
	if len(r) <= width || width < 2 {
		return s
	}
	return "…" + string(r[len(r)-width+1:])
}

func shortenHome(path string) string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return path
	}
	if rel, err := filepath.Rel(home, path); err == nil && !strings.HasPrefix(rel, "..") {
		return filepath.Join("~", rel)
	}
	return path
}
