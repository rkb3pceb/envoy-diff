package env

import (
	"testing"

	"github.com/wvdschel-personal/envoy-diff/internal/diff"
)

func TestDiffMaps_NoChanges(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new_ := map[string]string{"A": "1", "B": "2"}

	r := DiffMaps(old, new_, DefaultDiffOptions())
	if len(r.Changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(r.Changes))
	}
	if r.Added != 0 || r.Removed != 0 || r.Modified != 0 {
		t.Errorf("unexpected counts: %+v", r)
	}
}

func TestDiffMaps_DetectsAdded(t *testing.T) {
	old := map[string]string{"A": "1"}
	new_ := map[string]string{"A": "1", "B": "2"}

	r := DiffMaps(old, new_, DefaultDiffOptions())
	if r.Added != 1 {
		t.Fatalf("expected 1 added, got %d", r.Added)
	}
	if r.Changes[0].Type != diff.Added || r.Changes[0].Key != "B" {
		t.Errorf("unexpected change: %+v", r.Changes[0])
	}
}

func TestDiffMaps_DetectsRemoved(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new_ := map[string]string{"A": "1"}

	r := DiffMaps(old, new_, DefaultDiffOptions())
	if r.Removed != 1 {
		t.Fatalf("expected 1 removed, got %d", r.Removed)
	}
	if r.Changes[0].Key != "B" || r.Changes[0].Type != diff.Removed {
		t.Errorf("unexpected change: %+v", r.Changes[0])
	}
}

func TestDiffMaps_DetectsModified(t *testing.T) {
	old := map[string]string{"KEY": "old"}
	new_ := map[string]string{"KEY": "new"}

	r := DiffMaps(old, new_, DefaultDiffOptions())
	if r.Modified != 1 {
		t.Fatalf("expected 1 modified, got %d", r.Modified)
	}
	c := r.Changes[0]
	if c.OldValue != "old" || c.NewValue != "new" {
		t.Errorf("unexpected values: %+v", c)
	}
}

func TestDiffMaps_IncludeUnchanged(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new_ := map[string]string{"A": "1", "B": "99"}

	opts := DefaultDiffOptions()
	opts.IncludeUnchanged = true

	r := DiffMaps(old, new_, opts)
	if r.Unchanged != 1 {
		t.Fatalf("expected 1 unchanged, got %d", r.Unchanged)
	}
	if r.Modified != 1 {
		t.Fatalf("expected 1 modified, got %d", r.Modified)
	}
	if len(r.Changes) != 2 {
		t.Errorf("expected 2 total changes, got %d", len(r.Changes))
	}
}

func TestDiffMaps_SortedByKey(t *testing.T) {
	old := map[string]string{"Z": "1", "A": "1", "M": "1"}
	new_ := map[string]string{"Z": "2", "A": "2", "M": "2"}

	r := DiffMaps(old, new_, DefaultDiffOptions())
	if len(r.Changes) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(r.Changes))
	}
	keys := []string{r.Changes[0].Key, r.Changes[1].Key, r.Changes[2].Key}
	expected := []string{"A", "M", "Z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("position %d: expected %s, got %s", i, expected[i], k)
		}
	}
}
