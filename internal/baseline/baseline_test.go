package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-diff/internal/baseline"
)

func tempStore(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baselines.json")
}

func TestOpen_NewStore_IsEmpty(t *testing.T) {
	s, err := baseline.Open(tempStore(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.List(); len(got) != 0 {
		t.Errorf("expected empty store, got %v", got)
	}
}

func TestSet_And_Get_RoundTrip(t *testing.T) {
	path := tempStore(t)
	s, _ := baseline.Open(path)

	vars := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	if err := s.Set("prod", "production", vars); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	bl := s.Get("prod")
	if bl == nil {
		t.Fatal("expected baseline, got nil")
	}
	if bl.Vars["APP_ENV"] != "production" {
		t.Errorf("expected production, got %s", bl.Vars["APP_ENV"])
	}
	if bl.Env != "production" {
		t.Errorf("expected env production, got %s", bl.Env)
	}
}

func TestSet_PersistsToDisk(t *testing.T) {
	path := tempStore(t)
	s1, _ := baseline.Open(path)
	_ = s1.Set("staging", "staging", map[string]string{"DEBUG": "true"})

	s2, err := baseline.Open(path)
	if err != nil {
		t.Fatalf("reopen failed: %v", err)
	}
	if bl := s2.Get("staging"); bl == nil {
		t.Error("expected baseline to persist across open")
	}
}

func TestGet_Missing_ReturnsNil(t *testing.T) {
	s, _ := baseline.Open(tempStore(t))
	if bl := s.Get("nonexistent"); bl != nil {
		t.Errorf("expected nil, got %v", bl)
	}
}

func TestDelete_RemovesBaseline(t *testing.T) {
	path := tempStore(t)
	s, _ := baseline.Open(path)
	_ = s.Set("tmp", "dev", map[string]string{"X": "1"})

	if err := s.Delete("tmp"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if bl := s.Get("tmp"); bl != nil {
		t.Error("expected nil after delete")
	}
}

func TestDelete_NotFound_ReturnsError(t *testing.T) {
	s, _ := baseline.Open(tempStore(t))
	if err := s.Delete("ghost"); err == nil {
		t.Error("expected error deleting nonexistent baseline")
	}
}

func TestOpen_InvalidJSON_ReturnsError(t *testing.T) {
	path := tempStore(t)
	_ = os.WriteFile(path, []byte("not-json{"), 0o644)
	if _, err := baseline.Open(path); err == nil {
		t.Error("expected parse error for invalid JSON")
	}
}
