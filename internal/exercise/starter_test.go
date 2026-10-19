package exercise

import (
	"strings"
	"testing"
)

func TestStarterHasProblemHeader(t *testing.T) {
	cat, err := Load(validFS())
	if err != nil {
		t.Fatal(err)
	}
	ex, _ := cat.Get("zip-01")
	ex.Examples[0].Explanation = strings.Repeat("word ", 30)
	got, err := cat.Starter(ex)
	if err != nil {
		t.Fatal(err)
	}
	s := string(got)
	for _, want := range []string{
		"# zip-01 · Title zip-01\n#\n# Do the thing.\n",
		"# Example 1:\n#   Input:  x\n#   Output: y\n#   Explanation: word",
		"#     word",
		"leetkit check zip-01",
		"\n\n# starter.py",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in:\n%s", want, s)
		}
	}
	for _, line := range strings.Split(s, "\n") {
		if len(line) > headerWidth+2 {
			t.Errorf("line too long (%d): %q", len(line), line)
		}
	}
}

func TestExamplesRequired(t *testing.T) {
	fsys := validFS()
	f := fsys[zipDir+MetadataFile]
	f.Data = []byte(strings.Replace(string(f.Data), `explanation = "Because."`, "", 1))
	if _, err := Load(fsys); err == nil || !strings.Contains(err.Error(), "needs input, output, and explanation") {
		t.Errorf("err = %v", err)
	}
}
