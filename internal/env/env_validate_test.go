package env

import (
	"testing"
)

func TestValidateMap_NoFindings_WhenAllClean(t *testing.T) {
	m := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	}
	opts := DefaultValidateOptions()
	findings := ValidateMap(m, opts)
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

func TestValidateMap_LowercaseKey_ReturnsWarning(t *testing.T) {
	m := map[string]string{"app_host": "localhost"}
	opts := DefaultValidateOptions()
	findings := ValidateMap(m, opts)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Level != LevelWarning {
		t.Errorf("expected warning, got %s", findings[0].Level)
	}
}

func TestValidateMap_EmptyValue_ReturnsError_WhenForbidEmpty(t *testing.T) {
	m := map[string]string{"DB_PASS": ""}
	opts := DefaultValidateOptions()
	opts.ForbidEmpty = true
	findings := ValidateMap(m, opts)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Level != LevelError {
		t.Errorf("expected error level, got %s", findings[0].Level)
	}
}

func TestValidateMap_EmptyValue_NoFinding_WhenNotForbidden(t *testing.T) {
	m := map[string]string{"DB_PASS": ""}
	opts := DefaultValidateOptions()
	opts.ForbidEmpty = false
	findings := ValidateMap(m, opts)
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

func TestValidateMap_KeyTooLong_ReturnsError(t *testing.T) {
	longKey := "ABCDEFGHIJ_ABCDEFGHIJ_ABCDEFGHIJ" // 33 chars
	m := map[string]string{longKey: "val"}
	opts := DefaultValidateOptions()
	opts.MaxKeyLength = 20
	findings := ValidateMap(m, opts)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Level != LevelError {
		t.Errorf("expected error, got %s", findings[0].Level)
	}
}

func TestValidateMap_AllowedPrefixes_ViolationIsWarning(t *testing.T) {
	m := map[string]string{
		"APP_HOST": "localhost",
		"LEGACY_X": "old",
	}
	opts := DefaultValidateOptions()
	opts.AllowedPrefixes = []string{"APP_"}
	findings := ValidateMap(m, opts)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for LEGACY_X, got %d", len(findings))
	}
	if findings[0].Key != "LEGACY_X" {
		t.Errorf("expected finding for LEGACY_X, got %s", findings[0].Key)
	}
	if findings[0].Level != LevelWarning {
		t.Errorf("expected warning, got %s", findings[0].Level)
	}
}

func TestHasValidationErrors_TrueOnError(t *testing.T) {
	findings := []ValidationFinding{
		{Key: "X", Message: "bad", Level: LevelError},
	}
	if !HasValidationErrors(findings) {
		t.Error("expected HasValidationErrors to return true")
	}
}

func TestHasValidationErrors_FalseOnWarningsOnly(t *testing.T) {
	findings := []ValidationFinding{
		{Key: "X", Message: "warn", Level: LevelWarning},
	}
	if HasValidationErrors(findings) {
		t.Error("expected HasValidationErrors to return false for warnings only")
	}
}
