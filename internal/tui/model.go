// Package tui implements the interactive LeetKit dashboard. Like Rustlings,
// it opens on the current exercise and reruns its tests on demand; the
// exercise list is a secondary screen.
package tui

import (
	"context"
	"errors"
	"fmt"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/kidskoding/leetkit/internal/editor"
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
	Editor    editor.Editor
	// EditorErr explains why no editor could be resolved.
	EditorErr error
	Theme     theme.Theme
	// Notice is shown in the footer until the first key press.
	Notice string
}

type screen int

const (
	screenWatch screen = iota
	screenList
	screenProgress
)

type Model struct {
	opts      Options
	th        theme.Theme
	exercises []*exercise.Exercise

	screen     screen
	prevScreen screen
	current    int
	hints      int
	editing    bool
	// explaining shows explanation.md instead of the exercise details.
	explaining bool
	// confirmReset holds the exercise id awaiting a y/n reset confirmation.
	confirmReset string
	// lastMod is the working copy's last seen modification time.
	lastMod time.Time

	run runState

	watchKeys    watchKeys
	listKeys     listKeys
	progressKeys progressKeys
	list         list.Model
	vp           viewport.Model
	spinner      spinner.Model
	help         help.Model
	bar          progress.Model
	notice       string

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
		opts:         opts,
		th:           th,
		exercises:    opts.Catalog.All(),
		watchKeys:    newWatchKeys(),
		listKeys:     newListKeys(),
		progressKeys: newProgressKeys(),
		notice:       opts.Notice,
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
	m.bar.EmptyColor = th.Palette.Muted

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
	cmds := []tea.Cmd{m.startRunCmd()}
	if ex := m.currentExercise(); ex != nil {
		cmds = append(cmds, checkFile(m.opts.Workspace.Path(ex.ID)))
	}
	return tea.Batch(cmds...)
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
	case editorDoneMsg:
		return m, m.finishEdit(msg)
	case fileCheckMsg:
		return m, m.onFileCheck(msg)
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
		if m.acceptsGlobalKeys() {
			switch msg.String() {
			case "p":
				m.prevScreen, m.screen = m.screen, screenProgress
				return m, nil
			case "?":
				m.help.ShowAll = !m.help.ShowAll
				desc := "more"
				if m.help.ShowAll {
					desc = "less"
				}
				m.watchKeys.Help.SetHelp("?", desc)
				m.listKeys.Help.SetHelp("?", desc)
				m.layout()
				return m, nil
			}
		}
	}
	switch m.screen {
	case screenList:
		return m.updateList(msg)
	case screenProgress:
		return m.updateProgress(msg)
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
	switch m.screen {
	case screenList:
		return m.listView()
	case screenProgress:
		return m.progressView()
	}
	return m.watchView()
}

// acceptsGlobalKeys is false while a text field or prompt owns the keyboard.
func (m Model) acceptsGlobalKeys() bool {
	switch {
	case m.screen == screenProgress, m.confirmReset != "":
		return false
	case m.screen == screenList && m.list.SettingFilter():
		return false
	}
	return true
}

func (m *Model) layout() {
	m.help.SetWidth(m.width - 2)
	m.layoutWatch()
	m.layoutList()
}
