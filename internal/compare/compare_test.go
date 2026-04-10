package compare_test

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/compare"
	"github.com/your-org/envoy-diff/internal/diff"
)

func TestRun_NoChanges(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := compare.Run(env, env, compare.Options{})
	if s.HasChanges() {
		t.Errorf("expected no changes, got added=%d removed=%d modified=%d", s.Added, s.Removed, s.Modified)
	}
	if s.Unchanged != 2 {
		t.Errorf("expected 2 unchanged, got %d", s.Unchanged)
	}
}

func TestRun_CountsCorrectly(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2", "C": "3"}
	new := map[string]string{"A": "1", "B": "changed", "D": "4"}
	s := compare.Run(old, new, compare.Options{})

	if s.Added != 1 {
		t.Errorf("expected 1 added, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", s.Removed)
	}
	if s.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", s.Modified)
	}
	if s.Unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", s.Unchanged)
	}
	if !s.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestRun_ChangesSortedByKey(t *testing.T) {
	old := map[string]string{"Z": "1", "A": "1", "M": "1"}
	new := map[string]string{"Z": "2", "A": "2", "M": "2"}
	s := compare.Run(old, new, compare.Options{})

	for i := 1; i < len(s.Changes); i++ {
		if s.Changes[i-1].Key > s.Changes[i].Key {
			t.Errorf("changes not sorted: %s > %s", s.Changes[i-1].Key, s.Changes[i].Key)
		}
	}
}

func TestRun_RedactSensitiveValues(t *testing.T) {
	old := map[string]string{"DB_PASSWORD": "secret123"}
	new := map[string]string{"DB_PASSWORD": "newsecret"}
	s := compare.Run(old, new, compare.Options{RedactSensitive: true})

	if len(s.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(s.Changes))
	}
	c := s.Changes[0]
	if c.Type != diff.Modified {
		t.Errorf("expected Modified, got %v", c.Type)
	}
	if c.OldValue == "secret123" || c.NewValue == "newsecret" {
		t.Errorf("sensitive values were not redacted: old=%q new=%q", c.OldValue, c.NewValue)
	}
}

func TestRun_NoRedact_WhenDisabled(t *testing.T) {
	old := map[string]string{"DB_PASSWORD": "secret123"}
	new := map[string]string{"DB_PASSWORD": "newsecret"}
	s := compare.Run(old, new, compare.Options{RedactSensitive: false})

	if s.Changes[0].OldValue != "secret123" {
		t.Errorf("expected plain old value, got %q", s.Changes[0].OldValue)
	}
}
