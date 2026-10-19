package tui

import (
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/kidskoding/leetkit/internal/theme"
)

// renderMarkdown styles the small Markdown subset used by explanation.md:
// headings, fenced code, bullets, and inline code. It avoids pulling in a
// full Markdown renderer for a handful of short documents.
func renderMarkdown(th theme.Theme, src string, width int) string {
	var out, para []string
	inFence := false
	flush := func() {
		if len(para) > 0 {
			out = append(out, lipgloss.NewStyle().Width(width).Render(inlineCode(th, strings.Join(para, " "))))
			para = nil
		}
	}
	for _, line := range strings.Split(strings.TrimSpace(src), "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "```"):
			flush()
			inFence = !inFence
		case inFence:
			out = append(out, "  "+th.ID.Render(line))
		case trimmed == "":
			flush()
			out = append(out, "")
		case strings.HasPrefix(trimmed, "#"):
			flush()
			out = append(out, th.Heading.Render(strings.TrimSpace(strings.TrimLeft(trimmed, "#"))))
		case strings.HasPrefix(trimmed, "- "), strings.HasPrefix(trimmed, "* "):
			flush()
			out = append(out, wrapHanging(th, "• ", inlineCode(th, trimmed[2:]), width))
		default:
			// Source paragraphs are hard-wrapped; join them so they wrap to the window.
			para = append(para, trimmed)
		}
	}
	flush()
	return strings.Join(out, "\n")
}

func inlineCode(th theme.Theme, s string) string {
	parts := strings.Split(s, "`")
	if len(parts)%2 == 0 {
		return th.Text.Render(s)
	}
	var b strings.Builder
	for i, p := range parts {
		if i%2 == 1 {
			b.WriteString(th.ID.Render(p))
		} else {
			b.WriteString(th.Text.Render(p))
		}
	}
	return b.String()
}

func wrapHanging(th theme.Theme, bullet, text string, width int) string {
	body := lipgloss.NewStyle().Width(max(width-lipgloss.Width(bullet), 1)).Render(text)
	indent := strings.Repeat(" ", lipgloss.Width(bullet))
	lines := strings.Split(body, "\n")
	for i := range lines {
		if i == 0 {
			lines[i] = th.Heading.Render(bullet) + lines[i]
		} else {
			lines[i] = indent + lines[i]
		}
	}
	return strings.Join(lines, "\n")
}
