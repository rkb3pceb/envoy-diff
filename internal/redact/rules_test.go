package redact

import (
	"testing"
)

func TestNewRuleSet_Len(t *testing.T) {
	rs := NewRuleSet([]Rule{
		{Pattern: "token"},
		{Pattern: "cert"},
	})
	if rs.Len() != 2 {
		t.Fatalf("expected 2 rules, got %d", rs.Len())
	}
}

func TestEmpty_HasNoRules(t *testing.T) {
	rs := Empty()
	if rs.Len() != 0 {
		t.Fatalf("expected 0 rules, got %d", rs.Len())
	}
}

func TestMatches_FindsSubstringCaseInsensitive(t *testing.T) {
	rs := NewRuleSet([]Rule{
		{Pattern: "token", Placeholder: "[TOKEN]"},
	})

	tests := []struct {
		key   string
		want  bool
	}{
		{"AUTH_TOKEN", true},
		{"auth_token", true},
		{"Token", true},
		{"API_KEY", false},
	}

	for _, tc := range tests {
		_, got := rs.Matches(tc.key)
		if got != tc.want {
			t.Errorf("Matches(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}

func TestMatches_ReturnsCorrectRule(t *testing.T) {
	rs := NewRuleSet([]Rule{
		{Pattern: "cert", Placeholder: "[CERT]"},
	})
	r, ok := rs.Matches("TLS_CERT_PATH")
	if !ok {
		t.Fatal("expected match, got none")
	}
	if r.Placeholder != "[CERT]" {
		t.Errorf("expected placeholder [CERT], got %q", r.Placeholder)
	}
}

func TestRedactValue_UsesRulePlaceholder(t *testing.T) {
	rs := NewRuleSet([]Rule{
		{Pattern: "internal", Placeholder: "[INTERNAL]"},
	})
	got := rs.RedactValue("INTERNAL_SECRET", "mysecret")
	if got != "[INTERNAL]" {
		t.Errorf("expected [INTERNAL], got %q", got)
	}
}

func TestRedactValue_DefaultPlaceholderWhenEmpty(t *testing.T) {
	rs := NewRuleSet([]Rule{
		{Pattern: "magic", Placeholder: ""},
	})
	got := rs.RedactValue("MAGIC_KEY", "value")
	if got != Placeholder {
		t.Errorf("expected default placeholder %q, got %q", Placeholder, got)
	}
}

func TestRedactValue_FallsBackToBuiltIn(t *testing.T) {
	rs := Empty()
	// PASSWORD is caught by built-in IsSensitive patterns.
	got := rs.RedactValue("DB_PASSWORD", "hunter2")
	if got != Placeholder {
		t.Errorf("expected redacted value, got %q", got)
	}
}

func TestRedactValue_PassesThroughSafe(t *testing.T) {
	rs := Empty()
	got := rs.RedactValue("APP_ENV", "production")
	if got != "production" {
		t.Errorf("expected 'production', got %q", got)
	}
}
