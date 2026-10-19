package tui

import "charm.land/bubbles/v2/key"

// Navigation follows ranger/lf conventions: j/k move, l enters an exercise
// from the list, and h leaves it for the list. Arrow keys mirror them.
type watchKeys struct {
	// showNext hides n from the help until the exercise is solved.
	showNext bool

	Next   key.Binding
	Edit   key.Binding
	Hint   key.Binding
	Run    key.Binding
	List   key.Binding
	Scroll key.Binding
	Quit   key.Binding
}

func newWatchKeys() watchKeys {
	return watchKeys{
		Next:   key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "next")),
		Hint:   key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "hint")),
		Edit:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		Run:    key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "run")),
		List:   key.NewBinding(key.WithKeys("h", "left", "esc"), key.WithHelp("h/←", "list")),
		Scroll: key.NewBinding(key.WithKeys("j", "k", "down", "up", "pgdown", "pgup"), key.WithHelp("j/k", "scroll")),
		Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	}
}

func (k watchKeys) ShortHelp() []key.Binding {
	keys := []key.Binding{k.Edit, k.Hint, k.Run, k.List, k.Scroll, k.Quit}
	if k.showNext {
		keys = append([]key.Binding{k.Next}, keys...)
	}
	return keys
}

func (k watchKeys) FullHelp() [][]key.Binding { return [][]key.Binding{k.ShortHelp()} }

type listKeys struct {
	Move   key.Binding
	Open   key.Binding
	Search key.Binding
	Reset  key.Binding
	Back   key.Binding
}

func newListKeys() listKeys {
	return listKeys{
		Move:   key.NewBinding(key.WithKeys("j", "k", "down", "up"), key.WithHelp("j/k", "move")),
		Open:   key.NewBinding(key.WithKeys("l", "right", "enter"), key.WithHelp("l/enter", "open")),
		Search: key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		Reset:  key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reset")),
		Back:   key.NewBinding(key.WithKeys("esc", "q"), key.WithHelp("esc", "back")),
	}
}

func (k listKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Move, k.Open, k.Search, k.Reset, k.Back}
}

func (k listKeys) FullHelp() [][]key.Binding { return [][]key.Binding{k.ShortHelp()} }
