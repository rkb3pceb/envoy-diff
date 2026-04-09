package redact_test

import (
	"testing"

	"github.com/yourusername/envoy-diff/internal/redact"
)

func TestIsSensitive_MatchesKnownPatterns(t *testing.T) {
	sensitiveKeys := []string{
		"DB_PASSWORD",
		"API_KEY",
		"AUTH_TOKEN",
		"AWS_SECRET",
		"PRIVATE_KEY",
		"ACCESS_KEY_ID",
		"SIGNING_KEY",
		"app_secret",    // lowercase
		"Github_Token",  // mixed case
	}
	for _, key := range sensitiveKeys {
		if !redact.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_AllowsNonSensitive(t *testing.T) {
	safeKeys := []string{
		"APP_ENV",
		"LOG_LEVEL",
		"PORT",
		"DATABASE_HOST",
		"FEATURE_FLAG",
	}
	for _, key := range safeKeys {
		if redact.IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestValue_RedactsSensitive(t *testing.T) {
	got := redact.Value("DB_PASSWORD", "supersecret")
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}

func TestValue_PassesThroughSafe(t *testing.T) {
	got := redact.Value("APP_ENV", "production")
	if got != "production" {
		t.Errorf("expected 'production', got %q", got)
	}
}

func TestMap_RedactsOnlySensitiveKeys(t *testing.T) {
	input := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "hunter2",
		"PORT":        "8080",
		"API_KEY":     "abc123",
	}
	out := redact.Map(input)

	if out["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be unchanged")
	}
	if out["PORT"] != "8080" {
		t.Errorf("PORT should be unchanged")
	}
	if out["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("DB_PASSWORD should be redacted")
	}
	if out["API_KEY"] != "[REDACTED]" {
		t.Errorf("API_KEY should be redacted")
	}
}

func TestMap_DoesNotMutateOriginal(t *testing.T) {
	input := map[string]string{"DB_PASSWORD": "secret"}
	_ = redact.Map(input)
	if input["DB_PASSWORD"] != "secret" {
		t.Error("original map should not be mutated")
	}
}
