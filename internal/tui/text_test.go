package tui

import "testing"

func TestReflow(t *testing.T) {
	in := "Implement f(x).\n\nReturn the index of the first\noccurrence, or -1 if\nabsent.\n\n- keep this\n- and this\n  indented stays\n1. numbered\n"
	want := "Implement f(x).\n\nReturn the index of the first occurrence, or -1 if absent.\n\n- keep this\n- and this\n  indented stays\n1. numbered"
	if got := reflow(in); got != want {
		t.Errorf("reflow =\n%q\nwant\n%q", got, want)
	}
}
