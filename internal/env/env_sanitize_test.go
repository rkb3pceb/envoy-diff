package env

import (
	"testing"
)

func TestSanitizeMap_NoOp_DefaultOptions(t *testing.T) {
	input := map[string]string{"FOO": "bar", "BAZ": "qux"}
	opts := DefaultSanitizeOptions()
	r := SanitizeMap(input, opts)
	if r.Map["FOO"] != "bar" || r.Map["BAZ"] != "qux" {
		t.Error("expected values to be unchanged")
	}
	if len(r.Renamed) != 0 || len(r.Dropped) != 0 {
		t.Error("expected no renames or drops")
	}
}

func TestSanitizeMap_TrimSpace_TrimsValues(t *testing.T) {
	input := map[string]string{"KEY": "  hello  "}
	opts := DefaultSanitizeOptions()
	opts.TrimSpace = true
	r := SanitizeMap(input, opts)
	if r.Map["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", r.Map["KEY"])
	}
}

func TestSanitizeMap_RemoveEmpty_DropsBlankValues(t *testing.T) {
	input := map[string]string{"KEEP": "value", "DROP": ""}
	opts := DefaultSanitizeOptions()
	opts.RemoveEmpty = true
	r := SanitizeMap(input, opts)
	if _, ok := r.Map["DROP"]; ok {
		t.Error("expected DROP to be removed")
	}
	if r.Map["KEEP"] != "value" {
		t.Error("expected KEEP to remain")
	}
	if len(r.Dropped) != 1 || r.Dropped[0] != "DROP" {
		t.Errorf("expected Dropped=[DROP], got %v", r.Dropped)
	}
}

func TestSanitizeMap_NormalizeKeys_Uppercases(t *testing.T) {
	input := map[string]string{"my_key": "val"}
	opts := DefaultSanitizeOptions()
	opts.NormalizeKeys = true
	r := SanitizeMap(input, opts)
	if _, ok := r.Map["MY_KEY"]; !ok {
		t.Error("expected MY_KEY to exist after normalization")
	}
	if r.Renamed["my_key"] != "MY_KEY" {
		t.Errorf("expected rename my_key->MY_KEY, got %v", r.Renamed)
	}
}

func TestSanitizeMap_ReplaceInvalidChars(t *testing.T) {
	input := map[string]string{"my-key.name": "val"}
	opts := DefaultSanitizeOptions()
	opts.ReplaceInvalidChars = true
	r := SanitizeMap(input, opts)
	if _, ok := r.Map["my_key_name"]; !ok {
		t.Errorf("expected my_key_name, got keys: %v", r.Map)
	}
	if r.Renamed["my-key.name"] != "my_key_name" {
		t.Errorf("unexpected rename map: %v", r.Renamed)
	}
}

func TestSanitizeMap_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"key": "  spaced  "}
	opts := DefaultSanitizeOptions()
	opts.TrimSpace = true
	SanitizeMap(input, opts)
	if input["key"] != "  spaced  " {
		t.Error("input map was mutated")
	}
}

func TestHasSanitizeChanges_ReturnsFalse_WhenClean(t *testing.T) {
	r := SanitizeResult{Map: map[string]string{"A": "1"}, Renamed: map[string]string{}, Dropped: nil}
	if HasSanitizeChanges(r) {
		t.Error("expected no changes")
	}
}

func TestHasSanitizeChanges_ReturnsTrue_WhenRenamed(t *testing.T) {
	r := SanitizeResult{
		Map:     map[string]string{"A_B": "1"},
		Renamed: map[string]string{"a-b": "A_B"},
		Dropped: nil,
	}
	if !HasSanitizeChanges(r) {
		t.Error("expected changes detected")
	}
}
