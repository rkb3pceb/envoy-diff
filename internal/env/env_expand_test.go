package env

import (
	"os"
	"testing"
)

func TestExpandMap_NoReferences(t *testing.T) {
	input := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	out, err := ExpandMap(input, DefaultExpandOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "myapp" || out["PORT"] != "8080" {
		t.Errorf("values should pass through unchanged, got %v", out)
	}
}

func TestExpandMap_ResolvesInternalReference(t *testing.T) {
	input := map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  "${BASE_URL}/api",
	}
	out, err := ExpandMap(input, DefaultExpandOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_URL"] != "https://example.com/api" {
		t.Errorf("expected expanded URL, got %q", out["API_URL"])
	}
}

func TestExpandMap_MissingRef_SilentByDefault(t *testing.T) {
	input := map[string]string{
		"GREETING": "Hello ${UNKNOWN_VAR}",
	}
	out, err := ExpandMap(input, DefaultExpandOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// os.Expand replaces missing vars with empty string
	if out["GREETING"] != "Hello " {
		t.Errorf("expected 'Hello ', got %q", out["GREETING"])
	}
}

func TestExpandMap_ErrorOnMissing(t *testing.T) {
	input := map[string]string{
		"VAL": "${DOES_NOT_EXIST}",
	}
	opts := ExpandOptions{ErrorOnMissing: true}
	_, err := ExpandMap(input, opts)
	if err == nil {
		t.Fatal("expected error for missing reference, got nil")
	}
}

func TestExpandMap_FallbackToOS(t *testing.T) {
	t.Setenv("OS_VAR", "from-os")
	input := map[string]string{
		"COMBINED": "prefix-${OS_VAR}",
	}
	opts := ExpandOptions{FallbackToOS: true}
	out, err := ExpandMap(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["COMBINED"] != "prefix-from-os" {
		t.Errorf("expected 'prefix-from-os', got %q", out["COMBINED"])
	}
}

func TestExpandMap_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{
		"HOST": "localhost",
		"ADDR": "${HOST}:9000",
	}
	original := input["ADDR"]
	_, _ = ExpandMap(input, DefaultExpandOptions())
	if input["ADDR"] != original {
		t.Errorf("input map was mutated")
	}
}

func TestMissingRefs_ReturnsUnresolved(t *testing.T) {
	input := map[string]string{
		"A": "${B}",
		"C": "${D} and ${E}",
	}
	missing := MissingRefs(input, false)
	if len(missing) != 3 {
		t.Errorf("expected 3 missing refs, got %d: %v", len(missing), missing)
	}
}

func TestMissingRefs_FallbackToOS_ReducesMissing(t *testing.T) {
	os.Setenv("KNOWN_OS_VAR", "yes")
	defer os.Unsetenv("KNOWN_OS_VAR")
	input := map[string]string{
		"X": "${KNOWN_OS_VAR} and ${TRULY_MISSING}",
	}
	missing := MissingRefs(input, true)
	if len(missing) != 1 || missing[0] != "TRULY_MISSING" {
		t.Errorf("expected only TRULY_MISSING, got %v", missing)
	}
}
