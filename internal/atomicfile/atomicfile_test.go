package atomicfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteCreatesAndReplaces(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "file.json")
	for _, content := range []string{"first", "second"} {
		if err := Write(path, []byte(content), 0o644); err != nil {
			t.Fatalf("Write(%q): %v", content, err)
		}
		got, err := os.ReadFile(path)
		if err != nil || string(got) != content {
			t.Fatalf("read back %q, %v; want %q", got, err, content)
		}
	}
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("temp files left behind: %v", entries)
	}
}

func TestWriteFailsWhenParentIsFile(t *testing.T) {
	parent := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(parent, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Write(filepath.Join(parent, "child"), []byte("x"), 0o644); err == nil {
		t.Error("Write succeeded under a regular file")
	}
}
