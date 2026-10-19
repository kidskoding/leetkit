package tui

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/anirudhkonidala/leetkit/internal/exercise"
	"github.com/anirudhkonidala/leetkit/internal/progress"
	"github.com/anirudhkonidala/leetkit/internal/theme"
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
	ex     *exercise.Exercise
	module string
	status status
}

func (i item) FilterValue() string {
	return i.ex.ID + " " + i.ex.Title + " " + i.module
}

type delegate struct {
	th      theme.Theme
	idWidth int
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
	id := fmt.Sprintf("%-*s", d.idWidth, it.ex.ID)
	cursor, idStyle, titleStyle := "  ", d.th.ID, d.th.Faint
	if index == m.Index() {
		cursor = d.th.Heading.Render("▌ ")
		idStyle, titleStyle = idStyle.Bold(true), d.th.Text
	}
	line := cursor + mark + " " + idStyle.Render(id) + "  " + titleStyle.Render(it.ex.Title)
	fmt.Fprint(w, lipgloss.NewStyle().MaxWidth(m.Width()).Render(line))
}
