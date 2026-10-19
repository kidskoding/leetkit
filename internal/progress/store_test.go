package progress

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func openAt(t *testing.T, path string, now time.Time) *Store {
	t.Helper()
	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	s.now = func() time.Time { return now }
	return s
}

func tempPath(t *testing.T) string {
	return filepath.Join(t.TempDir(), "progress.json")
}

var (
	t1 = time.Date(2026, 10, 19, 12, 0, 0, 0, time.UTC)
	t2 = t1.Add(time.Hour)
	t3 = t2.Add(time.Hour)
)

func TestOpenMissingFile(t *testing.T) {
	s := openAt(t, tempPath(t), t1)
	if _, ok := s.Entry("enumerate-01"); ok {
		t.Error("empty store has an entry")
	}
}

func TestRecordKeepsFirstCompletion(t *testing.T) {
	path := tempPath(t)
	s := openAt(t, path, t1)
	if _, err := s.Record("x", false); err != nil {
		t.Fatal(err)
	}
	s.now = func() time.Time { return t2 }
	if _, err := s.Record("x", true); err != nil {
		t.Fatal(err)
	}
	s.now = func() time.Time { return t3 }
	e, err := s.Record("x", false)
	if err != nil {
		t.Fatal(err)
	}
	if e.Attempts != 3 || !e.LastAttempt.Equal(t3) || !e.Completed() || !e.CompletedAt.Equal(t2) || !e.LastPassed.Equal(t2) {
		t.Errorf("entry = %+v", e)
	}

	reopened := openAt(t, path, t3)
	got, ok := reopened.Entry("x")
	if !ok || got.Attempts != 3 || !got.CompletedAt.Equal(t2) {
		t.Errorf("after reopen: %+v, %v", got, ok)
	}
}

func TestCountCompleted(t *testing.T) {
	s := openAt(t, tempPath(t), t1)
	for id, passed := range map[string]bool{"a": true, "b": false, "c": true} {
		if _, err := s.Record(id, passed); err != nil {
			t.Fatal(err)
		}
	}
	if got := s.CountCompleted([]string{"a", "b", "c", "d"}); got != 2 {
		t.Errorf("CountCompleted = %d, want 2", got)
	}
}

func TestReset(t *testing.T) {
	path := tempPath(t)
	s := openAt(t, path, t1)
	if _, err := s.Record("x", true); err != nil {
		t.Fatal(err)
	}
	if err := s.Reset("x"); err != nil {
		t.Fatal(err)
	}
	if s.Completed("x") || openAt(t, path, t1).Completed("x") {
		t.Error("Reset did not clear progress")
	}
}

func TestCorruptFileIsMovedAside(t *testing.T) {
	path := tempPath(t)
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := openAt(t, path, t1)
	if s.Recovered == "" {
		t.Fatal("Recovered not set")
	}
	backup, err := os.ReadFile(s.Recovered)
	if err != nil || string(backup) != "{not json" {
		t.Errorf("backup = %q, %v", backup, err)
	}
	if _, err := s.Record("x", true); err != nil {
		t.Fatalf("store unusable after recovery: %v", err)
	}
}

func TestNewerFormatIsRejected(t *testing.T) {
	path := tempPath(t)
	if err := os.WriteFile(path, []byte(`{"version": 99, "exercises": {}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Open(path); err == nil || !strings.Contains(err.Error(), "newer LeetKit") {
		t.Errorf("err = %v", err)
	}
}

func TestWritesMergeAcrossStores(t *testing.T) {
	path := tempPath(t)
	a := openAt(t, path, t1)
	b := openAt(t, path, t1)
	if _, err := a.Record("from-a", true); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Record("from-b", true); err != nil {
		t.Fatal(err)
	}
	if !openAt(t, path, t1).Completed("from-a") {
		t.Error("second store overwrote the first store's record")
	}
}

func TestConcurrentRecords(t *testing.T) {
	s := openAt(t, tempPath(t), t1)
	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			if _, err := s.Record("x", false); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if e, _ := s.Entry("x"); e.Attempts != 20 {
		t.Errorf("attempts = %d, want 20", e.Attempts)
	}
}

func TestCurrentPersists(t *testing.T) {
	path := tempPath(t)
	s := openAt(t, path, t1)
	if s.Current() != "" {
		t.Fatalf("new store current = %q", s.Current())
	}
	if err := s.SetCurrent("zip-01"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Record("enumerate-01", true); err != nil {
		t.Fatal(err)
	}
	if got := openAt(t, path, t1).Current(); got != "zip-01" {
		t.Errorf("after reopen current = %q, want zip-01", got)
	}
}
