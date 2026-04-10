package schema_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-diff/internal/schema"
)

func writeSchemaFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "schema.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_ValidSchema(t *testing.T) {
	path := writeSchemaFile(t, `
rules:
  - key: PORT
    type: int
    required: true
  - key: API_URL
    type: url
  - key: ENABLED
    type: bool
`)
	rules, err := schema.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rules))
	}
	if rules[0].Key != "PORT" || rules[0].Type != schema.TypeInt || !rules[0].Required {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
	if rules[1].Type != schema.TypeURL {
		t.Errorf("expected url type for API_URL")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := schema.Load("/nonexistent/schema.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	path := writeSchemaFile(t, `{not: [valid yaml`)
	_, err := schema.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoad_UnknownType(t *testing.T) {
	path := writeSchemaFile(t, `
rules:
  - key: FOO
    type: uuid
`)
	_, err := schema.Load(path)
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestLoad_EmptyRules(t *testing.T) {
	path := writeSchemaFile(t, `rules: []`)
	rules, err := schema.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected 0 rules, got %d", len(rules))
	}
}

func TestEmpty_ReturnsNil(t *testing.T) {
	if schema.Empty() != nil {
		t.Error("expected Empty() to return nil")
	}
}
