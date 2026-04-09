// Package baseline provides functionality for establishing and comparing
// environment variable baselines across deployment environments.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a named reference point for environment variables.
type Baseline struct {
	Name      string            `json:"name"`
	Env       string            `json:"env"`
	Vars      map[string]string `json:"vars"`
	CreatedAt time.Time         `json:"created_at"`
}

// Store manages persisted baselines on disk.
type Store struct {
	path string
	items map[string]*Baseline
}

// Open loads or initializes a baseline store at the given path.
func Open(path string) (*Store, error) {
	s := &Store{path: path, items: make(map[string]*Baseline)}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("baseline: read store: %w", err)
	}
	if err := json.Unmarshal(data, &s.items); err != nil {
		return nil, fmt.Errorf("baseline: parse store: %w", err)
	}
	return s, nil
}

// Set stores a baseline under the given name.
func (s *Store) Set(name, env string, vars map[string]string) error {
	s.items[name] = &Baseline{
		Name:      name,
		Env:       env,
		Vars:      vars,
		CreatedAt: time.Now().UTC(),
	}
	return s.flush()
}

// Get retrieves a baseline by name. Returns nil if not found.
func (s *Store) Get(name string) *Baseline {
	return s.items[name]
}

// List returns all stored baseline names.
func (s *Store) List() []string {
	names := make([]string, 0, len(s.items))
	for k := range s.items {
		names = append(names, k)
	}
	return names
}

// Delete removes a baseline by name.
func (s *Store) Delete(name string) error {
	if _, ok := s.items[name]; !ok {
		return fmt.Errorf("baseline: %q not found", name)
	}
	delete(s.items, name)
	return s.flush()
}

func (s *Store) flush() error {
	data, err := json.MarshalIndent(s.items, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write store: %w", err)
	}
	return nil
}
