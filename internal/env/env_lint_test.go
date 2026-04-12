package env

import (
	"testing"
)

func TestLintMap_NoFindings_CleanEnv(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
	}
	findings := LintMap(env, DefaultLintOptions())
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d: %+v", len(findings), findings)
	}
}

func TestLintMap_LowercaseKey_ReturnsWarning(t *testing.T) {
	env := map[string]string{"database_url": "postgres://localhost/db"}
	findings := LintMap(env, DefaultLintOptions())
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != LintWarning {
		t.Errorf("expected warning, got %s", findings[0].Severity)
	}
}

func TestLintMap_KeyWithSpace_ReturnsError(t *testing.T) {
	env := map[string]string{"BAD KEY": "value"}
	findings := LintMap(env, DefaultLintOptions())
	var found bool
	for _, f := range findings {
		if f.Severity == LintError && f.Key == "BAD KEY" {
			found = true
		}
	}
	if !found {
		t.Error("expected an error finding for key with whitespace")
	}
}

func TestLintMap_EmptyValue_ReturnsWarning(t *testing.T) {
	env := map[string]string{"MY_VAR": ""}
	findings := LintMap(env, DefaultLintOptions())
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != LintWarning {
		t.Errorf("expected warning severity, got %s", findings[0].Severity)
	}
}

func TestLintMap_PlaceholderValue_ReturnsWarning(t *testing.T) {
	env := map[string]string{"API_KEY": "changeme"}
	findings := LintMap(env, DefaultLintOptions())
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}
	if findings[0].Severity != LintWarning {
		t.Errorf("expected warning, got %s", findings[0].Severity)
	}
}

func TestLintMap_KeyTooLong_ReturnsError(t *testing.T) {
	long := "ABCDEFGHIJ_ABCDEFGHIJ_ABCDEFGHIJ_ABCDEFGHIJ_ABCDEFGHIJ_ABCDEFGHIJ"
	env := map[string]string{long: "value"}
	opts := DefaultLintOptions()
	opts.MaxKeyLength = 10
	findings := LintMap(env, opts)
	var found bool
	for _, f := range findings {
		if f.Severity == LintError {
			found = true
		}
	}
	if !found {
		t.Error("expected error finding for key exceeding max length")
	}
}

func TestHasLintErrors_ReturnsFalse_WhenOnlyWarnings(t *testing.T) {
	findings := []LintFinding{
		{Key: "X", Message: "warn", Severity: LintWarning},
	}
	if HasLintErrors(findings) {
		t.Error("expected no errors")
	}
}

func TestHasLintErrors_ReturnsTrue_WhenErrorPresent(t *testing.T) {
	findings := []LintFinding{
		{Key: "X", Message: "err", Severity: LintError},
	}
	if !HasLintErrors(findings) {
		t.Error("expected errors to be detected")
	}
}

func TestLintMap_DisabledChecks_NoFindings(t *testing.T) {
	env := map[string]string{"lower_key": "changeme"}
	opts := LintOptions{
		ForbidLowercase:    false,
		ForbidPlaceholders: false,
		ForbidEmpty:        false,
		ForbidWhitespace:   false,
	}
	findings := LintMap(env, opts)
	if len(findings) != 0 {
		t.Errorf("expected no findings with all checks disabled, got %d", len(findings))
	}
}
