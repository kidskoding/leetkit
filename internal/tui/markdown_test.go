package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	"github.com/kidskoding/leetkit/internal/theme"
)

func TestRenderMarkdown(t *testing.T) {
	th, err := theme.Load("")
	if err != nil {
		t.Fatal(err)
	}
	src := "# enumerate()\n\nUse `enumerate(nums)` here.\n\n```python\nfor i, x in enumerate(nums):\n    pass\n```\n\n- first `a`\n- second\n"
	got := ansi.Strip(renderMarkdown(th, src, 40))
	for _, want := range []string{"enumerate()\n", "Use enumerate(nums) here.", "  for i, x in enumerate(nums):", "      pass", "• first a", "• second"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	if strings.Contains(got, "```") || strings.Contains(got, "# ") {
		t.Errorf("markup leaked:\n%s", got)
	}
}
