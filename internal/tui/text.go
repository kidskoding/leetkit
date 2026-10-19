package tui

import (
	"strings"
	"unicode"
)

// reflow joins hard-wrapped source lines into paragraphs so they can be
// wrapped to the window width. Blank lines separate paragraphs; list items
// and indented lines keep their own line breaks.
func reflow(s string) string {
	var out []string
	var para []string
	flush := func() {
		if len(para) > 0 {
			out = append(out, strings.Join(para, " "))
			para = nil
		}
	}
	for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "":
			flush()
			out = append(out, "")
		case isListItem(trimmed) || startsIndented(line):
			flush()
			out = append(out, strings.TrimRight(line, " "))
		default:
			para = append(para, trimmed)
		}
	}
	flush()
	return strings.Join(out, "\n")
}

func isListItem(s string) bool {
	if strings.HasPrefix(s, "- ") || strings.HasPrefix(s, "* ") {
		return true
	}
	digits := strings.TrimLeftFunc(s, unicode.IsDigit)
	return len(digits) < len(s) && strings.HasPrefix(digits, ". ")
}

func startsIndented(line string) bool {
	return strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t")
}
