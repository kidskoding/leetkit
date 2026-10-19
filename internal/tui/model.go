// Package tui implements the interactive LeetKit dashboard.
package tui

import (
	"context"
	"errors"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anirudhkonidala/leetkit/internal/exercise"
	pstore "github.com/anirudhkonidala/leetkit/internal/progress"
	"github.com/anirudhkonidala/leetkit/internal/theme"
)

const (
	minWidth  = 60
	minHeight = 12
)

type Options struct {
	Catalog *exercise.Catalog
	Store   *pstore.Store
	Theme   theme.Theme
	// Notice is shown in the footer until the first key press, for example a
	// warning that a corrupt progress file was moved aside.
	Notice string
}

type Model struct {
	opts   Options
	th     theme.Theme
	keys   keyMap
	list   list.Model
	help   help.Model
	bar    progress.Model
	notice string

	width, height int
}

func Run(ctx context.Context, opts Options) error {
	_, err := tea.NewProgram(New(opts), tea.WithContext(ctx)).Run()
	if errors.Is(err, tea.ErrProgramKilled) && ctx.Err() != nil {
		return nil
	}
	return err
}

func New(opts Options) Model {
	th := opts.Theme
	idWidth := 0
	for _, ex := range opts.Catalog.All() {
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
	styles := list.DefaultStyles(true)
	styles.StatusBar = th.Faint.PaddingLeft(2)
	styles.StatusBarActiveFilter = th.Text
	styles.StatusBarFilterCount = th.Faint
	styles.NoItems = th.Faint.PaddingLeft(2)
	styles.Filter.Focused.Prompt = th.Heading
	styles.Filter.Blurred.Prompt = th.Faint
	styles.DefaultFilterCharacterMatch = th.Heading.Underline(true)
	l.Styles = styles

	h := help.New()
	h.Styles.ShortKey, h.Styles.FullKey = th.ID, th.ID
	h.Styles.ShortDesc, h.Styles.FullDesc = th.Faint, th.Faint
	h.Styles.ShortSeparator, h.Styles.FullSeparator = th.Faint, th.Faint

	bar := progress.New(progress.WithoutPercentage(), progress.WithColors(th.Palette.Success))
	bar.EmptyColor = th.Palette.Surface

	m := Model{opts: opts, th: th, keys: newKeyMap(), list: l, help: h, bar: bar, notice: opts.Notice}
	m.refreshItems()
	return m
}

func (m *Model) refreshItems() {
	var items []list.Item
	for _, mod := range m.opts.Catalog.Modules {
		for _, ex := range mod.Exercises {
			items = append(items, item{ex: ex, module: mod.Title, status: statusOf(m.opts.Store, ex.ID)})
		}
	}
	m.list.SetItems(items)
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.layout()
		return m, nil
	case tea.KeyPressMsg:
		m.notice = ""
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// While the filter has focus, every key is text for the filter.
		if !m.list.SettingFilter() {
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keys.Help):
				m.help.ShowAll = !m.help.ShowAll
				m.layout()
				return m, nil
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) selected() (*exercise.Exercise, bool) {
	it, ok := m.list.SelectedItem().(item)
	if !ok {
		return nil, false
	}
	return it.ex, true
}
