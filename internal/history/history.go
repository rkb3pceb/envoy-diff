// Package history manages a local history of snapshot comparisons,
// allowing users to review past diffs and track environment drift over time.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single recorded diff event.
type Entry struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	OldFile   string    `json:"old_file"`
	NewFile   string    `json:"new_file"`
	Added     int       `json:"added"`
	Removed   int       `json:"removed"`
	Modified  int       `json:"modified"`
	Findings  int       `json:"findings"`
}

// Store holds the collection of history entries.
type Store struct {
	Entries []Entry `json:"entries"`
	path    string
}

// Open loads an existing history store from disk, or returns an empty one.
func Open(path string) (*Store, error) {
	s := &Store{path: path}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("history: read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, s); err != nil {
		return nil, fmt.Errorf("history: parse %s: %w", path, err)
	}
	return s, nil
}

// Append adds a new entry to the store and persists it to disk.
func (s *Store) Append(e Entry) error {
	e.Timestamp = time.Now().UTC()
	s.Entries = append(s.Entries, e)
	return s.save()
}

// Last returns the most recent entry, if any.
func (s *Store) Last() (Entry, bool) {
	if len(s.Entries) == 0 {
		return Entry{}, false
	}
	return s.Entries[len(s.Entries)-1], true
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("history: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	return os.WriteFile(s.path, data, 0o644)
}
