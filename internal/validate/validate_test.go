package validate_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/validate"
)

func TestValidate_NoFindings_WhenRulesEmpty(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	findings := validate.Validate(env, validate.Rule{})
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

func TestValidate_RequiredKey_Missing(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost"}
	rule := validate.Rule{RequiredKeys: []string{"APP_PORT"}}
	findings := validate.Validate(env, rule)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "APP_PORT" {
		t.Errorf("unexpected key: %s", findings[0].Key)
	}
	if findings[0].Severity != validate.SeverityError {
		t.Errorf("expected error severity, got %s", findings[0].Severity)
	}
}

func TestValidate_RequiredKey_EmptyValue(t *testing.T) {
	env := map[string]string{"APP_PORT": ""}
	rule := validate.Rule{RequiredKeys: []string{"APP_PORT"}}
	findings := validate.Validate(env, rule)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for empty value, got %d", len(findings))
	}
}

func TestValidate_ForbiddenKey_Present(t *testing.T) {
	env := map[string]string{"DEBUG": "true", "APP_HOST": "localhost"}
	rule := validate.Rule{ForbiddenKeys: []string{"DEBUG"}}
	findings := validate.Validate(env, rule)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "DEBUG" {
		t.Errorf("unexpected key: %s", findings[0].Key)
	}
	if findings[0].Severity != validate.SeverityError {
		t.Errorf("expected error severity, got %s", findings[0].Severity)
	}
}

func TestValidate_AllowedPrefixes_ViolationIsWarning(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "INTERNAL_SECRET": "x"}
	rule := validate.Rule{AllowedPrefixes: []string{"APP_"}}
	findings := validate.Validate(env, rule)
	if len(findings) != 1 {
		t.Fatalf("expected 1 warning finding, got %d", len(findings))
	}
	if findings[0].Severity != validate.SeverityWarning {
		t.Errorf("expected warning severity, got %s", findings[0].Severity)
	}
}

func TestHasErrors_ReturnsTrueOnError(t *testing.T) {
	findings := []validate.Finding{
		{Key: "X", Message: "missing", Severity: validate.SeverityError},
	}
	if !validate.HasErrors(findings) {
		t.Error("expected HasErrors to return true")
	}
}

func TestHasErrors_ReturnsFalseOnWarningsOnly(t *testing.T) {
	findings := []validate.Finding{
		{Key: "X", Message: "warn", Severity: validate.SeverityWarning},
	}
	if validate.HasErrors(findings) {
		t.Error("expected HasErrors to return false for warnings-only")
	}
}
