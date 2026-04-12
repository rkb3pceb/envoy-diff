package env

import (
	"testing"
)

func TestValidateSchema_NoFindings_WhenRulesEmpty(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	findings := ValidateSchema(env, nil)
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(findings))
	}
}

func TestValidateSchema_RequiredKey_Missing(t *testing.T) {
	env := map[string]string{}
	rules := []SchemaRule{{Key: "DB_HOST", Type: SchemaTypeString, Required: true}}
	findings := ValidateSchema(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if !findings[0].Error {
		t.Error("expected finding to be an error")
	}
}

func TestValidateSchema_IntType_Valid(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	rules := []SchemaRule{{Key: "PORT", Type: SchemaTypeInt}}
	findings := ValidateSchema(env, rules)
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(findings))
	}
}

func TestValidateSchema_IntType_Invalid(t *testing.T) {
	env := map[string]string{"PORT": "not-a-number"}
	rules := []SchemaRule{{Key: "PORT", Type: SchemaTypeInt}}
	findings := ValidateSchema(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestValidateSchema_BoolType_Valid(t *testing.T) {
	for _, val := range []string{"true", "false", "1", "0", "True", "FALSE"} {
		env := map[string]string{"DEBUG": val}
		rules := []SchemaRule{{Key: "DEBUG", Type: SchemaTypeBool}}
		findings := ValidateSchema(env, rules)
		if len(findings) != 0 {
			t.Errorf("value %q: expected 0 findings, got %d", val, len(findings))
		}
	}
}

func TestValidateSchema_BoolType_Invalid(t *testing.T) {
	env := map[string]string{"DEBUG": "yes"}
	rules := []SchemaRule{{Key: "DEBUG", Type: SchemaTypeBool}}
	findings := ValidateSchema(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestValidateSchema_URLType_Valid(t *testing.T) {
	env := map[string]string{"API_URL": "https://example.com"}
	rules := []SchemaRule{{Key: "API_URL", Type: SchemaTypeURL}}
	findings := ValidateSchema(env, rules)
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(findings))
	}
}

func TestValidateSchema_URLType_Invalid(t *testing.T) {
	env := map[string]string{"API_URL": "ftp://example.com"}
	rules := []SchemaRule{{Key: "API_URL", Type: SchemaTypeURL}}
	findings := ValidateSchema(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestHasSchemaErrors_ReturnsTrueOnError(t *testing.T) {
	findings := []SchemaFinding{{Key: "X", Message: "bad", Error: true}}
	if !HasSchemaErrors(findings) {
		t.Error("expected HasSchemaErrors to return true")
	}
}

func TestHasSchemaErrors_ReturnsFalseWhenNone(t *testing.T) {
	if HasSchemaErrors(nil) {
		t.Error("expected HasSchemaErrors to return false for nil")
	}
}
