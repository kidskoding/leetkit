// Package theme defines LeetKit's color palettes and the styles built from them.
package theme

import (
	"fmt"
	"image/color"
	"maps"
	"slices"
	"strings"

	"charm.land/lipgloss/v2"
)

// Default follows the terminal's own color scheme.
const Default = "terminal"

// Palette assigns colors to roles rather than hues, so every theme maps onto
// the same UI semantics.
type Palette struct {
	Background color.Color
	Surface    color.Color
	Foreground color.Color
	Muted      color.Color
	// Accent marks headings and focus.
	Accent color.Color
	// Link marks identifiers and paths.
	Link    color.Color
	Success color.Color
	Failure color.Color
	Warning color.Color
}

var palettes = map[string]Palette{
	// terminal uses ANSI palette indices, which the terminal renders with its
	// own color scheme, and leaves text and background at terminal defaults.
	"terminal": {
		Background: lipgloss.NoColor{},
		Surface:    lipgloss.BrightBlack,
		Foreground: lipgloss.NoColor{},
		Muted:      lipgloss.BrightBlack,
		Accent:     lipgloss.Magenta,
		Link:       lipgloss.Blue,
		Success:    lipgloss.Green,
		Failure:    lipgloss.Red,
		Warning:    lipgloss.Yellow,
	},
	"gruvbox": {
		Background: lipgloss.Color("#282828"),
		Surface:    lipgloss.Color("#3c3836"),
		Foreground: lipgloss.Color("#ebdbb2"),
		Muted:      lipgloss.Color("#928374"),
		Accent:     lipgloss.Color("#fabd2f"),
		Link:       lipgloss.Color("#83a598"),
		Success:    lipgloss.Color("#b8bb26"),
		Failure:    lipgloss.Color("#fb4934"),
		Warning:    lipgloss.Color("#fe8019"),
	},
	"gruvbox-light": {
		Background: lipgloss.Color("#fbf1c7"),
		Surface:    lipgloss.Color("#ebdbb2"),
		Foreground: lipgloss.Color("#3c3836"),
		Muted:      lipgloss.Color("#928374"),
		Accent:     lipgloss.Color("#b57614"),
		Link:       lipgloss.Color("#076678"),
		Success:    lipgloss.Color("#79740e"),
		Failure:    lipgloss.Color("#9d0006"),
		Warning:    lipgloss.Color("#af3a03"),
	},
	"tokyonight": {
		Background: lipgloss.Color("#1a1b26"),
		Surface:    lipgloss.Color("#24283b"),
		Foreground: lipgloss.Color("#c0caf5"),
		Muted:      lipgloss.Color("#565f89"),
		Accent:     lipgloss.Color("#bb9af7"),
		Link:       lipgloss.Color("#7aa2f7"),
		Success:    lipgloss.Color("#9ece6a"),
		Failure:    lipgloss.Color("#f7768e"),
		Warning:    lipgloss.Color("#e0af68"),
	},
}

func Names() []string {
	return slices.Sorted(maps.Keys(palettes))
}

type Theme struct {
	Name    string
	Palette Palette

	Heading lipgloss.Style
	ID      lipgloss.Style
	Text    lipgloss.Style
	Faint   lipgloss.Style
	Pass    lipgloss.Style
	Fail    lipgloss.Style
	Warn    lipgloss.Style
	Good    lipgloss.Style
	Bad     lipgloss.Style
}

// Load returns the named built-in theme; an empty name selects Default.
func Load(name string) (Theme, error) {
	if name == "" {
		name = Default
	}
	p, ok := palettes[name]
	if !ok {
		return Theme{}, fmt.Errorf("unknown theme %q (available: %s)", name, strings.Join(Names(), ", "))
	}
	return New(name, p), nil
}

func New(name string, p Palette) Theme {
	fg := lipgloss.NewStyle().Foreground
	return Theme{
		Name:    name,
		Palette: p,
		Heading: fg(p.Accent).Bold(true),
		ID:      fg(p.Link),
		Text:    fg(p.Foreground),
		Faint:   fg(p.Muted),
		Pass:    fg(p.Success).Bold(true),
		Fail:    fg(p.Failure).Bold(true),
		Warn:    fg(p.Warning).Bold(true),
		Good:    fg(p.Success),
		Bad:     fg(p.Failure),
	}
}
