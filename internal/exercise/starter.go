package exercise

import (
	"fmt"
	"strings"
)

const headerWidth = 76

// Starter returns the code a learner starts from: the bundled starter.py
// preceded by a comment header with the problem statement and examples, so
// the problem is visible in the editor, as Rustlings does.
func (c *Catalog) Starter(ex *Exercise) ([]byte, error) {
	code, err := c.ReadFile(ex, StarterFile)
	if err != nil {
		return nil, err
	}
	return append([]byte(Header(ex)), code...), nil
}

func Header(ex *Exercise) string {
	var lines []string
	add := func(s string) { lines = append(lines, s) }
	add(ex.ID + " · " + ex.Title)
	add("")
	for _, l := range strings.Split(strings.TrimSpace(ex.Description), "\n") {
		add(l)
	}
	for i, e := range ex.Examples {
		add("")
		add(fmt.Sprintf("Example %d:", i+1))
		add("  Input:  " + e.Input)
		add("  Output: " + e.Output)
		for j, l := range wrap("Explanation: "+e.Explanation, headerWidth-4) {
			if j > 0 {
				l = "  " + l
			}
			add("  " + l)
		}
	}
	add("")
	add("Run the tests in LeetKit (r) or with `leetkit check " + ex.ID + "`.")

	var b strings.Builder
	for _, l := range lines {
		b.WriteString(strings.TrimRight("# "+l, " ") + "\n")
	}
	b.WriteString("\n")
	return b.String()
}

func wrap(s string, width int) []string {
	var lines []string
	line := ""
	for _, word := range strings.Fields(s) {
		if line != "" && len(line)+1+len(word) > width {
			lines = append(lines, line)
			line = word
			continue
		}
		if line != "" {
			line += " "
		}
		line += word
	}
	return append(lines, line)
}
