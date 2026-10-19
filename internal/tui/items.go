package tui

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/kidskoding/leetkit/internal/exercise"
	"github.com/kidskoding/leetkit/internal/progress"
	"github.com/kidskoding/leetkit/internal/theme"
)

type status int

const (
	statusNotStarted status = iota
	statusAttempted
	statusDone
)

func statusOf(store *progress.Store, id string) status {
	e, ok := store.Entry(id)
	switch {
	case !ok:
		return statusNotStarted
	case e.Completed():
		return statusDone
	default:
		return statusAttempted
	}
}

type item struct {
	num    int
	ex     *exercise.Exercise
	module string
	status status
}

func (i item) FilterValue() string {
	return i.ex.ID + " " + i.ex.Title + " " + i.module
}

type delegate struct {
	th       theme.Theme
	numWidth int
}

func (d delegate) Height() int                         { return 1 }
func (d delegate) Spacing() int                        { return 0 }
func (d delegate) Update(tea.Msg, *list.Model) tea.Cmd { return nil }

func (d delegate) Render(w io.Writer, m list.Model, index int, li list.Item) {
	it, ok := li.(item)
	if !ok {
		return
	}
	var mark string
	switch it.status {
	case statusDone:
		mark = d.th.Pass.Render("✓")
	case statusAttempted:
		mark = d.th.Fail.Render("✗")
	default:
		mark = d.th.Faint.Render("○")
	}
	cursor, title := "  ", d.th.Text.Render(it.ex.Title)
	if index == m.Index() {
		cursor, title = d.th.Heading.Render("▌ "), d.th.Heading.Render(it.ex.Title)
	}
	num := d.th.Faint.Render(fmt.Sprintf("%0*d", d.numWidth, it.num))
	fmt.Fprint(w, ansi.Truncate(cursor+num+" "+mark+" "+title, m.Width(), "…"))
}
