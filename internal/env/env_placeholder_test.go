package env

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":   "myapp",
		"DB_PASS":    "CHANGEME",
		"API_KEY":    "TODO",
		"PORT":       "8080",
		"SECRET_KEY": "real-secret",
	}
}

func TestFindPlaceholders_DetectsKnownMarkers(t *testing.T) {
	results := FindPlaceholders(baseEnv(), DefaultPlaceholderOptions())
	if len(results) != 2 {
		t.Fatalf("expected 2 placeholders, got %d", len(results))
	}
	keys := map[string]bool{}
	for _, r := range results {
		keys[r.Key] = true
	}
	if !keys["DB_PASS"] || !keys["API_KEY"] {
		t.Errorf("expected DB_PASS and API_KEY to be flagged, got %v", keys)
	}
}

func TestFindPlaceholders_CleanEnv_ReturnsEmpty(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "5432"}
	results := FindPlaceholders(env, DefaultPlaceholderOptions())
	if len(results) != 0 {
		t.Errorf("expected no placeholders, got %d", len(results))
	}
}

func TestFindPlaceholders_CustomMarker(t *testing.T) {
	env := map[string]string{"TOKEN": "MYMARKER", "HOST": "localhost"}
	opts := DefaultPlaceholderOptions()
	opts.Markers = []string{"MYMARKER"}
	results := FindPlaceholders(env, opts)
	if len(results) != 1 || results[0].Key != "TOKEN" {
		t.Errorf("expected TOKEN to be flagged, got %+v", results)
	}
}

func TestSubstitutePlaceholders_ReplacesValue(t *testing.T) {
	opts := DefaultPlaceholderOptions()
	opts.Substitutions = map[string]string{"DB_PASS": "s3cr3t"}
	out, findings, err := SubstitutePlaceholders(baseEnv(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected substituted value, got %q", out["DB_PASS"])
	}
	resolved := 0
	for _, f := range findings {
		if f.Resolved {
			resolved++
		}
	}
	if resolved != 1 {
		t.Errorf("expected 1 resolved finding, got %d", resolved)
	}
}

func TestSubstitutePlaceholders_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"DB_PASS": "CHANGEME"}
	opts := DefaultPlaceholderOptions()
	opts.Substitutions = map[string]string{"DB_PASS": "new"}
	SubstitutePlaceholders(env, opts)
	if env["DB_PASS"] != "CHANGEME" {
		t.Error("input map was mutated")
	}
}

func TestSubstitutePlaceholders_ErrorOnUnresolved(t *testing.T) {
	opts := DefaultPlaceholderOptions()
	opts.ErrorOnUnresolved = true
	_, _, err := SubstitutePlaceholders	if err == nil {
		t.Error("expected error for unresolved placeholders")
	}
}

func TestSubstitutePlaceholders_NoError_WhenAllResolved(t *testing.T) {
	opts := DefaultPlaceholderOptions()
	opts.ErrorOnUnresolved = true
	opts.Substitutions = map[string]string{"DB_PASS": "x", "API_KEY": "y"}
	_, _, err := SubstitutePlaceholders(baseEnv(), opts)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHasPlaceholders_TrueWhenPresent(t *testing.T) {
	if !HasPlaceholders(baseEnv(), DefaultPlaceholderOptions()) {
		t.Error("expected HasPlaceholders to return true")
	}
}

func TestHasPlaceholders_FalseWhenClean(t *testing.T) {
	env := map[string]string{"HOST": "prod.example.com"}
	if HasPlaceholders(env, DefaultPlaceholderOptions()) {
		t.Error("expected HasPlaceholders to return false")
	}
}
