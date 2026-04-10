package schema_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/schema"
)

func TestValidate_NoFindings_WhenRulesEmpty(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	findings := schema.Validate(env, nil)
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

func TestValidate_RequiredKey_Missing(t *testing.T) {
	env := map[string]string{}
	rules := []schema.Rule{{Key: "PORT", Type: schema.TypeInt, Required: true}}
	findings := schema.Validate(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "PORT" {
		t.Errorf("unexpected key %q", findings[0].Key)
	}
}

func TestValidate_IntType_Valid(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	rules := []schema.Rule{{Key: "PORT", Type: schema.TypeInt}}
	findings := schema.Validate(env, rules)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %v", findings)
	}
}

func TestValidate_IntType_Invalid(t *testing.T) {
	env := map[string]string{"PORT": "not-a-number"}
	rules := []schema.Rule{{Key: "PORT", Type: schema.TypeInt}}
	findings := schema.Validate(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestValidate_BoolType_Valid(t *testing.T) {
	for _, v := range []string{"true", "false", "1", "0", "TRUE", "FALSE"} {
		env := map[string]string{"ENABLE_FEATURE": v}
		rules := []schema.Rule{{Key: "ENABLE_FEATURE", Type: schema.TypeBool}}
		if findings := schema.Validate(env, rules); len(findings) != 0 {
			t.Errorf("value %q should be valid bool, got findings: %v", v, findings)
		}
	}
}

func TestValidate_URLType_Invalid(t *testing.T) {
	env := map[string]string{"API_URL": "not-a-url"}
	rules := []schema.Rule{{Key: "API_URL", Type: schema.TypeURL}}
	findings := schema.Validate(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestValidate_URLType_Valid(t *testing.T) {
	env := map[string]string{"API_URL": "https://example.com/api"}
	rules := []schema.Rule{{Key: "API_URL", Type: schema.TypeURL}}
	if findings := schema.Validate(env, rules); len(findings) != 0 {
		t.Errorf("expected no findings, got %v", findings)
	}
}

func TestValidate_EmailType_Valid(t *testing.T) {
	env := map[string]string{"ADMIN_EMAIL": "admin@example.com"}
	rules := []schema.Rule{{Key: "ADMIN_EMAIL", Type: schema.TypeEmail}}
	if findings := schema.Validate(env, rules); len(findings) != 0 {
		t.Errorf("expected no findings, got %v", findings)
	}
}

func TestValidate_EmailType_Invalid(t *testing.T) {
	env := map[string]string{"ADMIN_EMAIL": "not-an-email"}
	rules := []schema.Rule{{Key: "ADMIN_EMAIL", Type: schema.TypeEmail}}
	findings := schema.Validate(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestHasErrors_TrueWhenFindings(t *testing.T) {
	findings := []schema.Finding{{Key: "X", Message: "bad"}}
	if !schema.HasErrors(findings) {
		t.Error("expected HasErrors to return true")
	}
}

func TestHasErrors_FalseWhenEmpty(t *testing.T) {
	if schema.HasErrors(nil) {
		t.Error("expected HasErrors to return false")
	}
}
