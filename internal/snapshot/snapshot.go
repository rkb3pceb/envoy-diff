// Package snapshot provides functionality to capture and persist
// environment variable states for later comparison.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a captured state of environment variables at a point in time.
type Snapshot struct {
	Name      string            `json:"name"`
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Vars      map[string]string `json:"vars"`
}

// New creates a new Snapshot from a map of environment variables.
func New(name, source string, vars map[string]string) *Snapshot {
	return &Snapshot{
		Name:      name,
		Timestamp: time.Now().UTC(),
		Source:    source,
		Vars:      vars,
	}
}

// Save writes the snapshot to a JSON file at the given path.
func Save(s *Snapshot, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: failed to create file %q: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: failed to encode snapshot: %w", err)
	}
	return nil
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: failed to open file %q: %w", path, err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: failed to decode snapshot: %w", err)
	}
	return &s, nil
}
