package env

import (
	"testing"
)

func TestUppercaseMap_NoOp_WhenBothDisabled(t *testing.T) {
	src := map[string]string{"db_host": "localhost", "port": "5432"}
	opts := DefaultUppercaseOptions()
	out := UppercaseMap(src, opts)
	for k, v := range src {
		if out[k] != v {
			t.Errorf("key %q: got %q, want %q", k, out[k], v)
		}
	}
}

func TestUppercaseMap_UppercasesKeys(t *testing.T) {
	src := map[string]string{"db_host": "localhost", "api_key": "secret"}
	opts := DefaultUppercaseOptions()
	opts.UppercaseKeys = true
	out := UppercaseMap(src, opts)
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected key DB_HOST")
	}
	if _, ok := out["API_KEY"]; !ok {
		t.Error("expected key API_KEY")
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestUppercaseMap_UppercasesValues(t *testing.T) {
	src := map[string]string{"ENV": "production", "MODE": "release"}
	opts := DefaultUppercaseOptions()
	opts.UppercaseValues = true
	out := UppercaseMap(src, opts)
	if out["ENV"] != "PRODUCTION" {
		t.Errorf("ENV: got %q, want PRODUCTION", out["ENV"])
	}
	if out["MODE"] != "RELEASE" {
		t.Errorf("MODE: got %q, want RELEASE", out["MODE"])
	}
}

func TestUppercaseMap_OnlyKeys_LimitsScope(t *testing.T) {
	src := map[string]string{"app_name": "myapp", "log_level": "debug"}
	opts := DefaultUppercaseOptions()
	opts.UppercaseKeys = true
	opts.OnlyKeys = []string{"app_name"}
	out := UppercaseMap(src, opts)
	if _, ok := out["APP_NAME"]; !ok {
		t.Error("expected APP_NAME to be uppercased")
	}
	if _, ok := out["log_level"]; !ok {
		t.Error("expected log_level to remain unchanged")
	}
}

func TestUppercaseMap_DoesNotMutateInput(t *testing.T) {
	src := map[string]string{"key": "value"}
	opts := DefaultUppercaseOptions()
	opts.UppercaseKeys = true
	opts.UppercaseValues = true
	_ = UppercaseMap(src, opts)
	if _, ok := src["key"]; !ok {
		t.Error("original map was mutated")
	}
	if src["key"] != "value" {
		t.Error("original value was mutated")
	}
}

func TestUppercaseMap_HasChanges_DetectsModification(t *testing.T) {
	src := map[string]string{"db_host": "localhost"}
	opts := DefaultUppercaseOptions()
	opts.UppercaseKeys = true
	if !HasUppercaseChanges(src, opts) {
		t.Error("expected HasUppercaseChanges to return true")
	}
}

func TestUppercaseMap_HasChanges_NoModification(t *testing.T) {
	src := map[string]string{"DB_HOST": "localhost"}
	opts := DefaultUppercaseOptions()
	opts.UppercaseKeys = true
	if HasUppercaseChanges(src, opts) {
		t.Error("expected HasUppercaseChanges to return false")
	}
}
