package workspace

import (
	"os"
	"testing"
)

func read(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestEnsureKeepsLearnerChanges(t *testing.T) {
	ws := Workspace{Dir: t.TempDir()}
	path, err := ws.Ensure("enumerate-01", []byte("starter"))
	if err != nil {
		t.Fatal(err)
	}
	if got := read(t, path); got != "starter" {
		t.Fatalf("first Ensure wrote %q", got)
	}
	if err := os.WriteFile(path, []byte("learner work"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ws.Ensure("enumerate-01", []byte("starter")); err != nil {
		t.Fatal(err)
	}
	if got := read(t, path); got != "learner work" {
		t.Errorf("Ensure overwrote learner work: %q", got)
	}
}

func TestResetRestoresStarter(t *testing.T) {
	ws := Workspace{Dir: t.TempDir()}
	path, err := ws.Ensure("zip-01", []byte("starter"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("broken"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ws.Reset("zip-01", []byte("starter")); err != nil {
		t.Fatal(err)
	}
	if got := read(t, path); got != "starter" {
		t.Errorf("after Reset got %q", got)
	}
}

func TestEmptyDirIsRejected(t *testing.T) {
	var ws Workspace
	if _, err := ws.Ensure("x", nil); err == nil {
		t.Error("Ensure with empty Dir wrote relative to the working directory")
	}
	if _, err := ws.Reset("x", nil); err == nil {
		t.Error("Reset with empty Dir wrote relative to the working directory")
	}
}
