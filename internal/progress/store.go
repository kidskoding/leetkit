// Package progress records exercise attempts and completions in a JSON file.
package progress

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sync"
	"time"

	"github.com/anirudhkonidala/leetkit/internal/atomicfile"
)

const schemaVersion = 1

type Entry struct {
	Attempts    int       `json:"attempts"`
	LastAttempt time.Time `json:"last_attempt"`
	// CompletedAt is the first successful run; later failures never clear it.
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	LastPassed  *time.Time `json:"last_passed,omitempty"`
}

func (e Entry) Completed() bool { return e.CompletedAt != nil }

type file struct {
	Version int `json:"version"`
	// Current is the exercise the dashboard resumes at.
	Current   string           `json:"current,omitempty"`
	Exercises map[string]Entry `json:"exercises"`
}

// Store is safe for concurrent use. Every write re-reads the file first, so
// two LeetKit processes do not discard each other's records.
type Store struct {
	// Recovered is the backup path when a corrupt file was found and moved aside.
	Recovered string

	path    string
	now     func() time.Time
	mu      sync.Mutex
	current string
	entries map[string]Entry
}

func Open(path string) (*Store, error) {
	s := &Store{path: path, now: func() time.Time { return time.Now().UTC() }}
	if err := s.reload(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) Entry(id string) (Entry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[id]
	return e, ok
}

func (s *Store) Current() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current
}

func (s *Store) SetCurrent(id string) error {
	return s.update(func(map[string]Entry) { s.current = id })
}

func (s *Store) Completed(id string) bool {
	e, _ := s.Entry(id)
	return e.Completed()
}

func (s *Store) CountCompleted(ids []string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for _, id := range ids {
		if s.entries[id].Completed() {
			n++
		}
	}
	return n
}

// Record counts an attempt and, if passed, a successful completion.
func (s *Store) Record(id string, passed bool) (Entry, error) {
	var updated Entry
	err := s.update(func(entries map[string]Entry) {
		e := entries[id]
		now := s.now()
		e.Attempts++
		e.LastAttempt = now
		if passed {
			e.LastPassed = &now
			if e.CompletedAt == nil {
				e.CompletedAt = &now
			}
		}
		entries[id] = e
		updated = e
	})
	return updated, err
}

func (s *Store) Reset(id string) error {
	return s.update(func(entries map[string]Entry) { delete(entries, id) })
}

func (s *Store) update(apply func(map[string]Entry)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.reload(); err != nil {
		return err
	}
	apply(s.entries)
	data, err := json.MarshalIndent(file{Version: schemaVersion, Current: s.current, Exercises: s.entries}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode progress: %w", err)
	}
	if err := atomicfile.Write(s.path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("save progress: %w", err)
	}
	return nil
}

// reload reads the file. A corrupt file is moved aside rather than
// overwritten so the learner can recover it by hand.
func (s *Store) reload() error {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, fs.ErrNotExist) {
		s.current, s.entries = "", map[string]Entry{}
		return nil
	}
	if err != nil {
		return fmt.Errorf("read progress: %w", err)
	}
	var f file
	if err := json.Unmarshal(data, &f); err != nil {
		backup := fmt.Sprintf("%s.corrupt-%s", s.path, s.now().Format("20060102T150405Z"))
		if err := os.Rename(s.path, backup); err != nil {
			return fmt.Errorf("progress file %s is corrupt and could not be moved aside: %w", s.path, err)
		}
		s.Recovered = backup
		s.current, s.entries = "", map[string]Entry{}
		return nil
	}
	if f.Version > schemaVersion {
		return fmt.Errorf("%s was written by a newer LeetKit (format %d); upgrade LeetKit", s.path, f.Version)
	}
	if f.Exercises == nil {
		f.Exercises = map[string]Entry{}
	}
	s.current, s.entries = f.Current, f.Exercises
	return nil
}
