package history_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-diff/internal/history"
)

func tempStorePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "history.json")
}

func TestOpen_NewStore_IsEmpty(t *testing.T) {
	s, err := history.Open(tempStorePath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Entries) != 0 {
		t.Errorf("expected empty store, got %d entries", len(s.Entries))
	}
}

func TestAppend_PersistsEntry(t *testing.T) {
	path := tempStorePath(t)
	s, _ := history.Open(path)

	e := history.Entry{ID: "abc123", OldFile: "a.env", NewFile: "b.env", Added: 2, Removed: 1}
	if err := s.Append(e); err != nil {
		t.Fatalf("append failed: %v", err)
	}

	s2, err := history.Open(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if len(s2.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(s2.Entries))
	}
	if s2.Entries[0].ID != "abc123" {
		t.Errorf("expected id abc123, got %s", s2.Entries[0].ID)
	}
}

func TestAppend_MultipleEntries_PreservesOrder(t *testing.T) {
	path := tempStorePath(t)
	s, _ := history.Open(path)

	ids := []string{"first", "second", "third"}
	for _, id := range ids {
		if err := s.Append(history.Entry{ID: id}); err != nil {
			t.Fatalf("append failed for id %s: %v", id, err)
		}
	}

	s2, err := history.Open(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if len(s2.Entries) != len(ids) {
		t.Fatalf("expected %d entries, got %d", len(ids), len(s2.Entries))
	}
	for i, id := range ids {
		if s2.Entries[i].ID != id {
			t.Errorf("entry[%d]: expected %s, got %s", i, id, s2.Entries[i].ID)
		}
	}
}

func TestLast_ReturnsNewest(t *testing.T) {
	path := tempStorePath(t)
	s, _ := history.Open(path)

	_ = s.Append(history.Entry{ID: "first"})
	_ = s.Append(history.Entry{ID: "second"})

	last, ok := s.Last()
	if !ok {
		t.Fatal("expected a last entry")
	}
	if last.ID != "second" {
		t.Errorf("expected 'second', got %s", last.ID)
	}
}

func TestLast_EmptyStore(t *testing.T) {
	s, _ := history.Open(tempStorePath(t))
	_, ok := s.Last()
	if ok {
		t.Error("expected no last entry for empty store")
	}
}

func TestOpen_InvalidJSON(t *testing.T) {
	path := tempStorePath(t)
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := history.Open(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
