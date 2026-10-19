// Package theme defines LeetKit's Tokyo Night palette and shared styles.
package theme

import "charm.land/lipgloss/v2"

var (
	Background = lipgloss.Color("#1a1b26")
	Surface    = lipgloss.Color("#24283b")
	Foreground = lipgloss.Color("#c0caf5")
	Blue       = lipgloss.Color("#7aa2f7")
	Purple     = lipgloss.Color("#bb9af7")
	Green      = lipgloss.Color("#9ece6a")
	Yellow     = lipgloss.Color("#e0af68")
	Red        = lipgloss.Color("#f7768e")
	Muted      = lipgloss.Color("#565f89")
)

var (
	Heading = lipgloss.NewStyle().Foreground(Purple).Bold(true)
	ID      = lipgloss.NewStyle().Foreground(Blue)
	Text    = lipgloss.NewStyle().Foreground(Foreground)
	Faint   = lipgloss.NewStyle().Foreground(Muted)
	Pass    = lipgloss.NewStyle().Foreground(Green).Bold(true)
	Fail    = lipgloss.NewStyle().Foreground(Red).Bold(true)
	Warn    = lipgloss.NewStyle().Foreground(Yellow).Bold(true)
	Good    = lipgloss.NewStyle().Foreground(Green)
	Bad     = lipgloss.NewStyle().Foreground(Red)
)
