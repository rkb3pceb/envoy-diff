package env

import (
	"testing"
)

func TestInspect_ExistingKey(t *testing.T) {
	env := map[string]string{"APP_PORT": "8080"}
	r := Inspect(env, "APP_PORT", DefaultInspectOptions())
	if !r.Exists {
		t.Fatal("expected key to exist")
	}
	if r.Value != "8080" {
		t.Errorf("unexpected value: %s", r.Value)
	}
	if r.Empty {
		t.Error("expected non-empty")
	}
}

func TestInspect_MissingKey(t *testing.T) {
	env := map[string]string{}
	r := Inspect(env, "MISSING", DefaultInspectOptions())
	if r.Exists {
		t.Error("expected key to be missing")
	}
	if r.Length != 0 {
		t.Errorf("expected length 0, got %d", r.Length)
	}
}

func TestInspect_SensitiveKeyRedacted(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "s3cr3t"}
	opts := DefaultInspectOptions()
	opts.Redact = true
	r := Inspect(env, "DB_PASSWORD", opts)
	if !r.Sensitive {
		t.Error("expected key to be marked sensitive")
	}
	if r.Value != "[REDACTED]" {
		t.Errorf("expected redacted value, got %q", r.Value)
	}
}

func TestInspect_SensitiveKeyNotRedacted(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "s3cr3t"}
	opts := DefaultInspectOptions()
	opts.Redact = false
	r := Inspect(env, "DB_PASSWORD", opts)
	if r.Value != "s3cr3t" {
		t.Errorf("expected plain value, got %q", r.Value)
	}
}

func TestInspect_EmptyValue(t *testing.T) {
	env := map[string]string{"EMPTY_KEY": ""}
	r := Inspect(env, "EMPTY_KEY", DefaultInspectOptions())
	if !r.Exists {
		t.Error("expected key to exist")
	}
	if !r.Empty {
		t.Error("expected Empty to be true")
	}
}

func TestInspect_MetadataNumeric(t *testing.T) {
	env := map[string]string{"PORT": "3000"}
	r := Inspect(env, "PORT", DefaultInspectOptions())
	if r.Metadata["numeric"] != "true" {
		t.Errorf("expected numeric=true, got %q", r.Metadata["numeric"])
	}
	if r.Metadata["boolean_like"] != "false" {
		t.Errorf("expected boolean_like=false, got %q", r.Metadata["boolean_like"])
	}
}

func TestInspect_MetadataBoolLike(t *testing.T) {
	env := map[string]string{"FEATURE_FLAG": "true"}
	r := Inspect(env, "FEATURE_FLAG", DefaultInspectOptions())
	if r.Metadata["boolean_like"] != "true" {
		t.Errorf("expected boolean_like=true, got %q", r.Metadata["boolean_like"])
	}
}

func TestInspectAll_SortedByKey(t *testing.T) {
	env := map[string]string{
		"ZEBRA": "z",
		"ALPHA": "a",
		"MANGO": "m",
	}
	results := InspectAll(env, DefaultInspectOptions())
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Key != "ALPHA" || results[1].Key != "MANGO" || results[2].Key != "ZEBRA" {
		t.Errorf("unexpected order: %v", []string{results[0].Key, results[1].Key, results[2].Key})
	}
}
