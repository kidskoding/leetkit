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
		m.revealResults = true
		return m, m.startRun()
	case key.Matches(k, keys.Next):
		return m, m.step(1)
	case key.Matches(k, keys.Prev):
		return m, m.step(-1)
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

// step moves through the curriculum in order, wrapping at either end.
func (m *Model) step(delta int) tea.Cmd {
	n := len(m.exercises)
	m.setCurrent((m.current + delta + n) % n)
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
	m.run.result, m.run.err, m.revealResults = nil, nil, false
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

// The exercise screen has a fixed header and status bar around a scrolling
// body, so the title and the latest result are always visible.
const watchHeaderLines = 3

func (m *Model) layoutWatch() {
	m.vp.SetWidth(max(m.width-2, 0))
	m.vp.SetHeight(max(m.height-watchHeaderLines-lipgloss.Height(m.watchStatus()), 0))
	m.refreshWatch()
}

// refreshWatch rebuilds the scrolling content; call it after any state change.
func (m *Model) refreshWatch() {
	if m.width == 0 || m.currentExercise() == nil {
		return
	}
	m.watchKeys.solved = m.opts.Store.Completed(m.currentExercise().ID)
	desc := "explain"
	if m.explaining {
		desc = "back"
	}
	m.watchKeys.Explain.SetHelp("x", desc)
	content, testsLine := m.watchContent(m.vp.Width())
	m.vp.SetContent(content)
	m.testsLine = testsLine
}

// showResults scrolls the body so the Tests section is at the top.
func (m *Model) showResults() {
	m.vp.SetYOffset(m.testsLine)
}

func (m Model) watchView() string {
	status := m.watchStatus()
	// The status bar can grow (full help), so size the viewport per render.
	vp := m.vp
	vp.SetHeight(max(m.height-watchHeaderLines-lipgloss.Height(status), 0))
	body := lipgloss.NewStyle().PaddingLeft(1).Render(vp.View())
	return lipgloss.JoinVertical(lipgloss.Left, m.watchHeader(), body, status)
}

func (m Model) watchHeader() string {
	ex := m.currentExercise()
	th := m.th
	pos := th.Faint.Render(fmt.Sprintf("%d / %d", m.current+1, len(m.exercises)))
	title := lipgloss.NewStyle().MaxWidth(max(m.width-lipgloss.Width(pos)-4, 1)).Render(th.Text.Bold(true).Render(ex.Title))
	gap := strings.Repeat(" ", max(m.width-lipgloss.Width(title)-lipgloss.Width(pos)-2, 1))
	meta := ex.ID + " · " + m.moduleTitle(ex) + " · " + ex.Difficulty + " · " + ex.TypeLabel()
	return strings.Join([]string{
		" " + title + gap + pos,
		" " + lipgloss.NewStyle().MaxWidth(m.width-2).Render(th.Faint.Render(meta)),
		m.separator(),
	}, "\n")
}

// watchStatus is the fixed bar under the body: the latest result with overall
// progress, the learner's file, and the key help or a notice.
func (m Model) watchStatus() string {
	th := m.th
	solved, total := m.solvedCount()
	meterWidth := min(max(m.width/5, 8), 24)
	right := report.Meter(th, solved, total, meterWidth) + " " + th.Faint.Render(fmt.Sprintf("%d/%d solved", solved, total))
	left := lipgloss.NewStyle().MaxWidth(max(m.width-lipgloss.Width(right)-4, 1)).Render(m.resultSummary())
	gap := strings.Repeat(" ", max(m.width-lipgloss.Width(left)-lipgloss.Width(right)-2, 1))

	footer := m.help.View(m.watchKeys)
	if m.notice != "" {
		footer = th.Warn.Render(m.notice)
	}
	path := truncateLeft(shortenHome(m.opts.Workspace.Path(m.currentExercise().ID)), m.width-2)
	return strings.Join([]string{
		m.separator(),
		" " + left + gap + right,
		" " + th.ID.Render(path),
		lipgloss.NewStyle().MaxWidth(m.width).PaddingLeft(1).Render(footer),
	}, "\n")
}

// resultSummary is the one-line state of the latest run.
func (m Model) resultSummary() string {
	th := m.th
	switch {
	case m.run.active:
		return m.spinner.View() + th.Faint.Render(" Running tests…")
	case m.run.err != nil:
		return th.Fail.Render("✗ Tests cannot run")
	case m.run.result == nil:
		return th.Faint.Render("Not run yet")
	}
	res := m.run.result
	passed, n := res.Count(runner.OutcomePassed), len(res.Cases)
	switch res.Status {
	case runner.StatusPassed:
		return th.Pass.Render("✓ Solved") + th.Faint.Render(fmt.Sprintf(" · %d of %d tests passing", passed, n))
	case runner.StatusFailed:
		return th.Fail.Render("✗ " + fmt.Sprintf("%d of %d tests passing", passed, n))
	case runner.StatusError:
		return th.Fail.Render("✗ Your code has an error")
	case runner.StatusTimedOut:
		return th.Warn.Render("! Timed out")
	}
	return th.Warn.Render("! " + res.Status.String())
}

func (m Model) solvedCount() (solved, total int) {
	ids := make([]string, len(m.exercises))
	for i, ex := range m.exercises {
		ids[i] = ex.ID
	}
	return m.opts.Store.CountCompleted(ids), len(ids)
}

func (m Model) separator() string {
	return m.th.Faint.Render(strings.Repeat("─", max(m.width, 0)))
}

// watchContent renders the scrolling body and reports the line where the
// Tests section starts, so runs can scroll straight to their results.
func (m Model) watchContent(width int) (content string, testsLine int) {
	ex := m.currentExercise()
	th := m.th
	wrap := lipgloss.NewStyle().Width(width)
	if m.explaining {
		doc, err := m.opts.Catalog.ReadFile(ex, exercise.ExplanationFile)
		if err != nil {
			return th.Fail.Render("✗ Could not read the explanation: " + err.Error()), 0
		}
		return m.rule("Explanation", width) + "\n\n" + renderMarkdown(th, string(doc), width), 0
	}

	var top strings.Builder
	top.WriteString(th.Text.Render(reflow(ex.Description)) + "\n")
	for i, e := range ex.Examples {
		top.WriteString("\n" + m.example(i+1, e, width) + "\n")
	}
	problem := strings.TrimRight(wrap.Render(top.String()), "\n ")

	var b strings.Builder
	b.WriteString(m.rule("Tests", width) + "\n")
	switch {
	case m.run.err != nil:
		b.WriteString(wrap.Render(th.Fail.Render("✗ "+m.runErrorText())) + "\n")
	case m.run.active && m.run.result == nil:
		b.WriteString(m.spinner.View() + th.Faint.Render(" Running tests…") + "\n")
	case m.run.result != nil:
		res := report.Result(th, m.run.result, false)
		if m.run.result.Status == runner.StatusPassed {
			// Put the next step right under the summary line.
			summary, rest, _ := strings.Cut(res, "\n")
			res = summary + "\n" + th.Faint.Render("Press ") + th.Heading.Render("n") +
				th.Faint.Render(" for the next exercise or ") + th.Heading.Render("x") +
				th.Faint.Render(" for the explanation.") + "\n" + rest
		}
		b.WriteString(res)
	default:
		b.WriteString(th.Faint.Render("Save your file or press r to run the tests.") + "\n")
	}

	if m.hints > 0 {
		b.WriteString("\n" + m.rule(fmt.Sprintf("Hints %d/%d", m.hints, len(ex.Hints)), width) + "\n")
		for i := range m.hints {
			b.WriteString(wrapHanging(th, fmt.Sprintf("%d. ", i+1), th.Text.Render(ex.Hints[i]), width) + "\n")
		}
	}
	return problem + "\n\n" + wrap.Render(b.String()), lipgloss.Height(problem) + 1
}

// example renders one example LeetCode-style: input, output, and why.
func (m Model) example(n int, e exercise.Example, width int) string {
	th := m.th
	label := func(s string) string { return th.Faint.Render(pad(s, 15)) }
	why := lipgloss.NewStyle().Width(max(width-15, 10)).Render(th.Text.Render(e.Explanation))
	whyLines := strings.Split(why, "\n")
	for i := 1; i < len(whyLines); i++ {
		whyLines[i] = strings.Repeat(" ", 15) + whyLines[i]
	}
	return strings.Join([]string{
		th.Text.Bold(true).Render(fmt.Sprintf("Example %d", n)),
		label("  Input:") + th.ID.Render(e.Input),
		label("  Output:") + th.ID.Render(e.Output),
		label("  Explanation:") + strings.Join(whyLines, "\n"),
	}, "\n")
}

func (m Model) runErrorText() string {
	err := m.run.err
	switch {
	case errors.Is(err, runner.ErrPytestMissing):
		return fmt.Sprintf("pytest is not installed for %s.\n\n%s", m.opts.Runner.Python, m.opts.SetupHelp)
	case errors.Is(err, runner.ErrNoPython):
		return "No Python interpreter found.\n\n" + m.opts.SetupHelp
	}
	return err.Error()
}

func (m Model) rule(title string, width int) string {
	label := "── " + title + " "
	return m.th.Heading.Render(label) + m.th.Faint.Render(strings.Repeat("─", max(width-lipgloss.Width(label), 0)))
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
