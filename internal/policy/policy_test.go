package policy_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/envoy-diff/internal/diff"
	"github.com/your-org/envoy-diff/internal/policy"
)

func writePolicy(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "policy.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_ValidPolicy(t *testing.T) {
	p := writePolicy(t, `version: "1"
rules:
  - name: no-remove-db-url
    keys: [DATABASE_URL]
    types: [removed]
    severity: block
    message: "DATABASE_URL must not be removed"
`)
	pol, err := policy.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pol.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(pol.Rules))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := policy.Load("/nonexistent/policy.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidSeverity(t *testing.T) {
	p := writePolicy(t, `rules:
  - name: bad
    keys: [FOO]
    types: [added]
    severity: critical
`)
	_, err := policy.Load(p)
	if err == nil {
		t.Fatal("expected validation error for invalid severity")
	}
}

func TestEvaluate_NoViolations(t *testing.T) {
	pol := &policy.Policy{
		Rules: []policy.Rule{
			{Name: "r1", Keys: []string{"SECRET"}, Types: []string{"removed"}, Severity: policy.SeverityBlock},
		},
	}
	changes := []diff.Change{
		{Key: "SECRET", Type: diff.Added, NewValue: "abc"},
	}
	violations := pol.Evaluate(changes)
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
}

func TestEvaluate_BlockViolation(t *testing.T) {
	pol := &policy.Policy{
		Rules: []policy.Rule{
			{Name: "no-remove-secret", Keys: []string{"SECRET"}, Types: []string{"removed"}, Severity: policy.SeverityBlock, Message: "do not remove SECRET"},
		},
	}
	changes := []diff.Change{
		{Key: "SECRET", Type: diff.Removed, OldValue: "old"},
	}
	violations := pol.Evaluate(changes)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !policy.HasBlockers(violations) {
		t.Error("expected HasBlockers to be true")
	}
}

func TestEvaluate_WildcardKey(t *testing.T) {
	pol := &policy.Policy{
		Rules: []policy.Rule{
			{Name: "db-warn", Keys: []string{"DB_*"}, Types: []string{"modified"}, Severity: policy.SeverityWarn},
		},
	}
	changes := []diff.Change{
		{Key: "DB_HOST", Type: diff.Modified, OldValue: "a", NewValue: "b"},
	}
	violations := pol.Evaluate(changes)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if policy.HasBlockers(violations) {
		t.Error("expected no blockers for warn severity")
	}
}

func TestEmpty_ReturnsNoViolations(t *testing.T) {
	pol := policy.Empty()
	changes := []diff.Change{
		{Key: "ANY", Type: diff.Added, NewValue: "v"},
	}
	if v := pol.Evaluate(changes); len(v) != 0 {
		t.Fatalf("expected 0 violations from empty policy, got %d", len(v))
	}
}
