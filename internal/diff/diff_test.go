package diff_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func TestCompare_NoChanges(t *testing.T) {
	old := map[string]string{"FOO": "bar", "BAZ": "qux"}
	new := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := diff.Compare(old, new)
	if result.HasChanges() {
		t.Errorf("expected no changes, got %d", len(result.Changes))
	}
}

func TestCompare_DetectsAdded(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new := map[string]string{"FOO": "bar", "NEW_KEY": "value"}

	result := diff.Compare(old, new)
	added := result.Added()
	if len(added) != 1 {
		t.Fatalf("expected 1 added change, got %d", len(added))
	}
	if added[0].Key != "NEW_KEY" || added[0].NewValue != "value" {
		t.Errorf("unexpected added change: %+v", added[0])
	}
}

func TestCompare_DetectsRemoved(t *testing.T) {
	old := map[string]string{"FOO": "bar", "OLD_KEY": "gone"}
	new := map[string]string{"FOO": "bar"}

	result := diff.Compare(old, new)
	removed := result.Removed()
	if len(removed) != 1 {
		t.Fatalf("expected 1 removed change, got %d", len(removed))
	}
	if removed[0].Key != "OLD_KEY" || removed[0].OldValue != "gone" {
		t.Errorf("unexpected removed change: %+v", removed[0])
	}
}

func TestCompare_DetectsModified(t *testing.T) {
	old := map[string]string{"FOO": "original"}
	new := map[string]string{"FOO": "updated"}

	result := diff.Compare(old, new)
	modified := result.Modified()
	if len(modified) != 1 {
		t.Fatalf("expected 1 modified change, got %d", len(modified))
	}
	if modified[0].OldValue != "original" || modified[0].NewValue != "updated" {
		t.Errorf("unexpected modified change: %+v", modified[0])
	}
}

func TestCompare_MixedChanges(t *testing.T) {
	old := map[string]string{"KEEP": "same", "CHANGE": "old", "DROP": "bye"}
	new := map[string]string{"KEEP": "same", "CHANGE": "new", "GAIN": "hello"}

	result := diff.Compare(old, new)
	if len(result.Added()) != 1 {
		t.Errorf("expected 1 added, got %d", len(result.Added()))
	}
	if len(result.Removed()) != 1 {
		t.Errorf("expected 1 removed, got %d", len(result.Removed()))
	}
	if len(result.Modified()) != 1 {
		t.Errorf("expected 1 modified, got %d", len(result.Modified()))
	}
}

func TestCompare_EmptyEnvs(t *testing.T) {
	result := diff.Compare(map[string]string{}, map[string]string{})
	if result.HasChanges() {
		t.Error("expected no changes for two empty envs")
	}
}
