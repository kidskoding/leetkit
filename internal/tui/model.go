// Package tui implements the interactive LeetKit dashboard. Like Rustlings,
// it opens on the current exercise and reruns its tests on demand; the
// exercise list is a secondary screen.
package tui

import (
	"context"
	"errors"
	"fmt"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/kidskoding/leetkit/internal/exercise"
	pstore "github.com/kidskoding/leetkit/internal/progress"
	"github.com/kidskoding/leetkit/internal/runner"
	"github.com/kidskoding/leetkit/internal/theme"
	"github.com/kidskoding/leetkit/internal/workspace"
)

const (
	minWidth  = 60
	minHeight = 12
)

type Options struct {
	Catalog   *exercise.Catalog
	Store     *pstore.Store
	Workspace workspace.Workspace
	Runner    runner.Runner
	// RunnerErr explains why tests cannot run, such as a missing Python. The
	// dashboard still opens so learners can browse exercises.
	RunnerErr error
	Theme     theme.Theme
	// Notice is shown in the footer until the first key press.
	Notice string
}

type screen int

const (
	screenWatch screen = iota
	screenList
)

type Model struct {
	opts      Options
	th        theme.Theme
	exercises []*exercise.Exercise

	screen  screen
	current int
	hints   int

	run runState

	watchKeys watchKeys
	listKeys  listKeys
	list      list.Model
	vp        viewport.Model
	spinner   spinner.Model
	help      help.Model
	bar       progress.Model
	notice    string

	width, height int
}

func Run(ctx context.Context, opts Options) error {
	final, err := tea.NewProgram(New(opts), tea.WithContext(ctx)).Run()
	if m, ok := final.(Model); ok {
		m.run.cancel()
	}
	if errors.Is(err, tea.ErrProgramKilled) && ctx.Err() != nil {
		return nil
	}
	return err
}

func New(opts Options) Model {
	th := opts.Theme
	m := Model{
		opts:      opts,
		th:        th,
		exercises: opts.Catalog.All(),
		watchKeys: newWatchKeys(),
		listKeys:  newListKeys(),
		notice:    opts.Notice,
	}
	m.list = newList(th, m.exercises)
	m.vp = viewport.New()
	m.vp.KeyMap = viewport.KeyMap{}
	m.spinner = spinner.New(spinner.WithSpinner(spinner.MiniDot), spinner.WithStyle(th.Heading))

	m.help = help.New()
	m.help.Styles.ShortKey, m.help.Styles.FullKey = th.Heading, th.Heading
	m.help.Styles.ShortDesc, m.help.Styles.FullDesc = th.Faint, th.Faint
	m.help.Styles.ShortSeparator, m.help.Styles.FullSeparator = th.Faint, th.Faint

	m.bar = progress.New(progress.WithoutPercentage(), progress.WithColors(th.Palette.Success))
	m.bar.EmptyColor = th.Palette.Surface

	m.current = m.startIndex()
	m.refreshItems()
	return m
}

// startIndex resumes at the saved exercise, else the first unfinished one.
func (m Model) startIndex() int {
	saved := m.opts.Store.Current()
	for i, ex := range m.exercises {
		if ex.ID == saved {
			return i
		}
	}
	if i := m.nextPending(-1); i >= 0 {
		return i
	}
	return 0
}

// nextPending returns the first unfinished exercise after from, wrapping
// around, or -1 when every exercise is complete.
func (m Model) nextPending(from int) int {
	n := len(m.exercises)
	for step := 1; step <= n; step++ {
		i := (from + step + n) % n
		if !m.opts.Store.Completed(m.exercises[i].ID) {
			return i
		}
	}
	return -1
}

func (m Model) currentExercise() *exercise.Exercise {
	if len(m.exercises) == 0 {
		return nil
	}
	return m.exercises[m.current]
}

func (m *Model) refreshItems() {
	items := make([]list.Item, len(m.exercises))
	for i, ex := range m.exercises {
		items[i] = item{ex: ex, module: m.moduleTitle(ex), status: statusOf(m.opts.Store, ex.ID)}
	}
	m.list.SetItems(items)
}

func (m Model) moduleTitle(ex *exercise.Exercise) string {
	for _, mod := range m.opts.Catalog.Modules {
		if mod.ID == ex.Module {
			return mod.Title
		}
	}
	return ex.Module
}

func (m Model) Init() tea.Cmd {
	return m.startRunCmd()
}

// startRunCmd defers the first run to Update, because Init cannot change the model.
func (m Model) startRunCmd() tea.Cmd {
	return func() tea.Msg { return requestRunMsg{} }
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.layout()
		return m, nil
	case requestRunMsg:
		return m, m.startRun()
	case runDoneMsg:
		return m, m.finishRun(msg)
	case spinner.TickMsg:
		if !m.run.active {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		m.refreshWatch()
		return m, cmd
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			m.run.cancel()
			return m, tea.Quit
		}
		m.notice = ""
	}
	if m.screen == screenList {
		return m.updateList(msg)
	}
	return m.updateWatch(msg)
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
	if len(m.exercises) == 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			m.th.Faint.Render("No exercises found. Press q to quit."))
	}
	if m.screen == screenList {
		return m.listView()
	}
	return m.watchView()
}

func (m *Model) layout() {
	m.help.SetWidth(m.width - 2)
	m.layoutWatch()
	m.layoutList()
}
