package audit_test

import (
	"testing"

	"github.com/user/envoy-diff/internal/audit"
	"github.com/user/envoy-diff/internal/diff"
)

func TestAudit_NoFindings(t *testing.T) {
	changes := []diff.Change{
		{Key: "APP_ENV", Type: diff.Modified, OldValue: "staging", NewValue: "production"},
	}
	findings := audit.Audit(changes)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestAudit_SensitiveModifiedIsHigh(t *testing.T) {
	changes := []diff.Change{
		{Key: "DB_PASSWORD", Type: diff.Modified, OldValue: "old", NewValue: "new"},
	}
	findings := audit.Audit(changes)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != audit.SeverityHigh {
		t.Errorf("expected HIGH severity, got %s", findings[0].Severity)
	}
}

func TestAudit_SensitiveAddedIsMedium(t *testing.T) {
	changes := []diff.Change{
		{Key: "STRIPE_API_KEY", Type: diff.Added, NewValue: "sk_live_abc"},
	}
	findings := audit.Audit(changes)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != audit.SeverityMedium {
		t.Errorf("expected MEDIUM severity, got %s", findings[0].Severity)
	}
}

func TestAudit_RemovedVariableIsMedium(t *testing.T) {
	changes := []diff.Change{
		{Key: "FEATURE_FLAG", Type: diff.Removed, OldValue: "true"},
	}
	findings := audit.Audit(changes)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != audit.SeverityMedium {
		t.Errorf("expected MEDIUM severity, got %s", findings[0].Severity)
	}
}

func TestAudit_MixedChanges(t *testing.T) {
	changes := []diff.Change{
		{Key: "APP_PORT", Type: diff.Modified, OldValue: "8080", NewValue: "9090"},
		{Key: "JWT_SECRET", Type: diff.Modified, OldValue: "old_secret", NewValue: "new_secret"},
		{Key: "DATABASE_URL", Type: diff.Added, NewValue: "postgres://localhost/db"},
		{Key: "OLD_VAR", Type: diff.Removed, OldValue: "value"},
	}
	findings := audit.Audit(changes)
	// JWT_SECRET (high), DATABASE_URL (medium), OLD_VAR (medium)
	if len(findings) != 3 {
		t.Errorf("expected 3 findings, got %d", len(findings))
	}
}
