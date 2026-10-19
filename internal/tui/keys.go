package tui

import "charm.land/bubbles/v2/key"

// Navigation follows ranger/lf conventions: j/k move, l enters an exercise
// from the list, and h leaves it for the list. Arrow keys mirror them.
type watchKeys struct {
	// solved shows n and x in the help once the exercise passes.
	solved bool

	Next     key.Binding
	Explain  key.Binding
	Edit     key.Binding
	Hint     key.Binding
	Run      key.Binding
	List     key.Binding
	Scroll   key.Binding
	Progress key.Binding
	Help     key.Binding
	Quit     key.Binding
}

func newWatchKeys() watchKeys {
	return watchKeys{
		Next:     key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "next")),
		Explain:  key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "explanation")),
		Edit:     key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		Hint:     key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "hint")),
		Run:      key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "run")),
		List:     key.NewBinding(key.WithKeys("h", "left", "esc"), key.WithHelp("h/←", "list")),
		Scroll:   key.NewBinding(key.WithKeys("j", "k", "down", "up", "pgdown", "pgup"), key.WithHelp("j/k", "scroll")),
		Progress: key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "progress")),
		Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "more")),
		Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	}
}

func (k watchKeys) solvedKeys() []key.Binding {
	if k.solved {
		return []key.Binding{k.Next, k.Explain}
	}
	return nil
}

func (k watchKeys) ShortHelp() []key.Binding {
	return append(k.solvedKeys(), k.Edit, k.Hint, k.Run, k.List, k.Help, k.Quit)
}

func (k watchKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		append(k.solvedKeys(), k.Edit, k.Hint, k.Run),
		{k.List, k.Scroll, k.Progress},
		{k.Help, k.Quit},
	}
}

type listKeys struct {
	Move     key.Binding
	Open     key.Binding
	Search   key.Binding
	Reset    key.Binding
	Progress key.Binding
	Help     key.Binding
	Back     key.Binding
}

func newListKeys() listKeys {
	return listKeys{
		Move:     key.NewBinding(key.WithKeys("j", "k", "down", "up"), key.WithHelp("j/k", "move")),
		Open:     key.NewBinding(key.WithKeys("l", "right", "enter"), key.WithHelp("l/enter", "open")),
		Search:   key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		Reset:    key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reset")),
		Progress: key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "progress")),
		Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "more")),
		Back:     key.NewBinding(key.WithKeys("esc", "q"), key.WithHelp("esc", "back")),
	}
}

func (k listKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Move, k.Open, k.Search, k.Reset, k.Help, k.Back}
}

func (k listKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Move, k.Open, k.Search}, {k.Reset, k.Progress}, {k.Help, k.Back}}
}

type progressKeys struct {
	Back key.Binding
	Quit key.Binding
}

func newProgressKeys() progressKeys {
	return progressKeys{
		Back: key.NewBinding(key.WithKeys("p", "esc", "h", "left"), key.WithHelp("p/esc", "back")),
		Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	}
}

func (k progressKeys) ShortHelp() []key.Binding  { return []key.Binding{k.Back, k.Quit} }
func (k progressKeys) FullHelp() [][]key.Binding { return [][]key.Binding{k.ShortHelp()} }
