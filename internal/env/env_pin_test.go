package env

import (
	"testing"
)

func TestCheckPins_NoViolations_AllMatch(t *testing.T) {
	pinned := map[string]string{"APP_ENV": "production", "LOG_LEVEL": "info"}
	env := map[string]string{"APP_ENV": "production", "LOG_LEVEL": "info", "EXTRA": "value"}

	result, err := CheckPins(pinned, env, DefaultPinOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasViolations() {
		t.Errorf("expected no violations, got %d", len(result.Violations))
	}
}

func TestCheckPins_DetectsChangedValue(t *testing.T) {
	pinned := map[string]string{"APP_ENV": "production"}
	env := map[string]string{"APP_ENV": "staging"}

	result, err := CheckPins(pinned, env, DefaultPinOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasViolations() {
		t.Fatal("expected a violation")
	}
	v := result.Violations[0]
	if v.Key != "APP_ENV" || v.Pinned != "production" || v.Actual != "staging" {
		t.Errorf("unexpected violation: %+v", v)
	}
}

func TestCheckPins_DetectsMissingKey(t *testing.T) {
	pinned := map[string]string{"REQUIRED_KEY": "value"}
	env := map[string]string{}

	result, _ := CheckPins(pinned, env, DefaultPinOptions())
	if !result.HasViolations() {
		t.Fatal("expected violation for missing key")
	}
	if result.Violations[0].Actual != "" {
		t.Errorf("expected empty actual for missing key")
	}
}

func TestCheckPins_ErrorOnViolation_ReturnsError(t *testing.T) {
	pinned := map[string]string{"APP_ENV": "production"}
	env := map[string]string{"APP_ENV": "dev"}

	opts := DefaultPinOptions()
	opts.ErrorOnViolation = true

	_, err := CheckPins(pinned, env, opts)
	if err == nil {
		t.Fatal("expected an error when ErrorOnViolation is set")
	}
}

func TestCheckPins_ErrorOnViolation_NoErrorWhenClean(t *testing.T) {
	pinned := map[string]string{"APP_ENV": "production"}
	env := map[string]string{"APP_ENV": "production"}

	opts := DefaultPinOptions()
	opts.ErrorOnViolation = true

	_, err := CheckPins(pinned, env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckPins_MultipleViolations_AllReported(t *testing.T) {
	pinned := map[string]string{"A": "1", "B": "2", "C": "3"}
	env := map[string]string{"A": "X", "B": "2", "C": "Y"}

	result, _ := CheckPins(pinned, env, DefaultPinOptions())
	if len(result.Violations) != 2 {
		t.Errorf("expected 2 violations, got %d", len(result.Violations))
	}
}
