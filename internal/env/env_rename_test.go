package env

import (
	"testing"
)

func TestRenameMap_NoRules_ReturnsCopy(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	opts := DefaultRenameOptions()
	res, err := RenameMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Map) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(res.Map))
	}
	if HasRenameApplied(res) {
		t.Error("expected no renames applied")
	}
}

func TestRenameMap_BasicRename(t *testing.T) {
	src := map[string]string{"OLD_KEY": "value"}
	opts := DefaultRenameOptions()
	opts.Rules = map[string]string{"OLD_KEY": "NEW_KEY"}
	res, err := RenameMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Map["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY to exist")
	}
	if _, ok := res.Map["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if !HasRenameApplied(res) {
		t.Error("expected rename to be recorded")
	}
}

func TestRenameMap_MissingKey_SkippedByDefault(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	opts := DefaultRenameOptions()
	opts.Rules = map[string]string{"MISSING": "NEW"}
	res, err := RenameMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING" {
		t.Errorf("expected MISSING in skipped, got %v", res.Skipped)
	}
}

func TestRenameMap_MissingKey_ErrorOnMissing(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	opts := DefaultRenameOptions()
	opts.Rules = map[string]string{"MISSING": "NEW"}
	opts.ErrorOnMissing = true
	_, err := RenameMap(src, opts)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRenameMap_ConflictError(t *testing.T) {
	src := map[string]string{"OLD": "v1", "NEW": "v2"}
	opts := DefaultRenameOptions()
	opts.Rules = map[string]string{"OLD": "NEW"}
	opts.ErrorOnConflict = true
	_, err := RenameMap(src, opts)
	if err == nil {
		t.Fatal("expected error for conflicting destination key")
	}
}

func TestRenameMap_DoesNotMutateInput(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	opts := DefaultRenameOptions()
	opts.Rules = map[string]string{"A": "Z"}
	_, err := RenameMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := src["A"]; !ok {
		t.Error("input map was mutated: A was removed")
	}
}

func TestRenameMap_PreservesValue(t *testing.T) {
	src := map[string]string{"DB_HOST": "localhost"}
	opts := DefaultRenameOptions()
	opts.Rules = map[string]string{"DB_HOST": "DATABASE_HOST"}
	res, err := RenameMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Map["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected value 'localhost', got %q", res.Map["DATABASE_HOST"])
	}
}
