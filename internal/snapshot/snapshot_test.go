package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envoy-diff/internal/snapshot"
)

func TestNew_SetsFieldsCorrectly(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := snapshot.New("test-snap", "staging.env", vars)

	if s.Name != "test-snap" {
		t.Errorf("expected name %q, got %q", "test-snap", s.Name)
	}
	if s.Source != "staging.env" {
		t.Errorf("expected source %q, got %q", "staging.env", s.Source)
	}
	if len(s.Vars) != 2 {
		t.Errorf("expected 2 vars, got %d", len(s.Vars))
	}
	if s.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	if s.Timestamp.Location() != time.UTC {
		t.Error("expected UTC timestamp")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	vars := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
	orig := snapshot.New("prod-snap", "prod.env", vars)

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "snap.json")

	if err := snapshot.Save(orig, path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Name != orig.Name {
		t.Errorf("name mismatch: got %q, want %q", loaded.Name, orig.Name)
	}
	if loaded.Source != orig.Source {
		t.Errorf("source mismatch: got %q, want %q", loaded.Source, orig.Source)
	}
	if len(loaded.Vars) != len(orig.Vars) {
		t.Errorf("vars length mismatch: got %d, want %d", len(loaded.Vars), len(orig.Vars))
	}
	for k, v := range orig.Vars {
		if loaded.Vars[k] != v {
			t.Errorf("var %q: got %q, want %q", k, loaded.Vars[k], v)
		}
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error loading missing file, got nil")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	s := snapshot.New("x", "y", map[string]string{})
	err := snapshot.Save(s, "/nonexistent/dir/snap.json")
	if err == nil {
		t.Error("expected error saving to invalid path, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "bad.json")
	if err := os.WriteFile(path, []byte("not json{"), 0644); err != nil {
		t.Fatalf("failed to write bad file: %v", err)
	}
	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error on invalid JSON, got nil")
	}
}
