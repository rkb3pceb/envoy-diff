package env

import (
	"testing"
)

func TestAddPrefix_AddsToAllKeys(t *testing.T) {
	m := map[string]string{"FOO": "1", "BAR": "2"}
	got := AddPrefix(m, "APP_", DefaultPrefixOptions())
	if _, ok := got["APP_FOO"]; !ok {
		t.Error("expected APP_FOO")
	}
	if _, ok := got["APP_BAR"]; !ok {
		t.Error("expected APP_BAR")
	}
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
}

func TestAddPrefix_StripExisting_NeverDoubles(t *testing.T) {
	m := map[string]string{"APP_FOO": "1", "BAR": "2"}
	opts := DefaultPrefixOptions()
	opts.StripExisting = true
	got := AddPrefix(m, "APP_", opts)
	if _, ok := got["APP_APP_FOO"]; ok {
		t.Error("should not double-prefix APP_FOO")
	}
	if _, ok := got["APP_FOO"]; !ok {
		t.Error("expected APP_FOO after strip+add")
	}
	if _, ok := got["APP_BAR"]; !ok {
		t.Error("expected APP_BAR")
	}
}

func TestAddPrefix_DoesNotMutateInput(t *testing.T) {
	m := map[string]string{"KEY": "val"}
	_ = AddPrefix(m, "P_", DefaultPrefixOptions())
	if _, ok := m["P_KEY"]; ok {
		t.Error("input map should not be mutated")
	}
}

func TestStripPrefix_RemovesFromMatchingKeys(t *testing.T) {
	m := map[string]string{"APP_FOO": "1", "APP_BAR": "2", "OTHER": "3"}
	got := StripPrefix(m, "APP_", DefaultPrefixOptions())
	if _, ok := got["FOO"]; !ok {
		t.Error("expected FOO")
	}
	if _, ok := got["BAR"]; !ok {
		t.Error("expected BAR")
	}
	if _, ok := got["OTHER"]; !ok {
		t.Error("expected OTHER to pass through")
	}
}

func TestStripPrefix_IgnoreCase(t *testing.T) {
	m := map[string]string{"app_foo": "1"}
	opts := DefaultPrefixOptions()
	opts.IgnoreCase = true
	got := StripPrefix(m, "APP_", opts)
	if _, ok := got["foo"]; !ok {
		t.Error("expected foo after case-insensitive strip")
	}
}

func TestHasPrefixedKeys_True(t *testing.T) {
	m := map[string]string{"APP_X": "1", "Y": "2"}
	if !HasPrefixedKeys(m, "APP_", false) {
		t.Error("expected true")
	}
}

func TestHasPrefixedKeys_False(t *testing.T) {
	m := map[string]string{"FOO": "1", "BAR": "2"}
	if HasPrefixedKeys(m, "APP_", false) {
		t.Error("expected false")
	}
}

func TestHasPrefixedKeys_IgnoreCase(t *testing.T) {
	m := map[string]string{"app_secret": "x"}
	if !HasPrefixedKeys(m, "APP_", true) {
		t.Error("expected true with ignore-case")
	}
}
