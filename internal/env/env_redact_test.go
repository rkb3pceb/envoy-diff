package env

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/redact"
)

func TestRedactMap_DisabledPassesThrough(t *testing.T) {
	m := map[string]string{"DB_PASSWORD": "secret123", "APP_ENV": "production"}
	opts := DefaultRedactOptions()
	opts.Enabled = false

	out := RedactMap(m, opts)
	if out["DB_PASSWORD"] != "secret123" {
		t.Errorf("expected plain value, got %q", out["DB_PASSWORD"])
	}
}

func TestRedactMap_MasksSensitiveKey(t *testing.T) {
	m := map[string]string{"DB_PASSWORD": "supersecret", "APP_ENV": "staging"}
	opts := DefaultRedactOptions()
	opts.Level = "full"

	out := RedactMap(m, opts)
	if out["DB_PASSWORD"] != redact.Placeholder {
		t.Errorf("expected placeholder, got %q", out["DB_PASSWORD"])
	}
	if out["APP_ENV"] != "staging" {
		t.Errorf("expected unchanged value, got %q", out["APP_ENV"])
	}
}

func TestRedactMap_PartialLevel_ShowsEdges(t *testing.T) {
	m := map[string]string{"API_SECRET_KEY": "abcdefghij"}
	opts := DefaultRedactOptions()
	opts.Level = "partial"

	out := RedactMap(m, opts)
	if out["API_SECRET_KEY"] == "abcdefghij" {
		t.Error("expected value to be masked")
	}
	if out["API_SECRET_KEY"] == "" {
		t.Error("expected non-empty masked value")
	}
}

func TestRedactMap_ExtraPatterns(t *testing.T) {
	m := map[string]string{"MY_CUSTOM_CRED": "topsecret", "NORMAL_VAR": "hello"}
	opts := DefaultRedactOptions()
	opts.Level = "full"
	opts.ExtraPatterns = []string{"custom_cred"}

	out := RedactMap(m, opts)
	if out["MY_CUSTOM_CRED"] != redact.Placeholder {
		t.Errorf("expected placeholder for extra pattern, got %q", out["MY_CUSTOM_CRED"])
	}
	if out["NORMAL_VAR"] != "hello" {
		t.Errorf("expected unchanged value, got %q", out["NORMAL_VAR"])
	}
}

func TestRedactMap_DoesNotMutateInput(t *testing.T) {
	m := map[string]string{"DB_PASSWORD": "original"}
	opts := DefaultRedactOptions()
	opts.Level = "full"

	_ = RedactMap(m, opts)
	if m["DB_PASSWORD"] != "original" {
		t.Error("input map was mutated")
	}
}

func TestSensitiveKeys_ReturnsMatchingKeys(t *testing.T) {
	m := map[string]string{
		"DB_PASSWORD": "x",
		"APP_NAME":    "envoy",
		"AUTH_TOKEN":  "tok",
	}

	keys := SensitiveKeys(m, nil)
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[k] = true
	}

	if !keySet["DB_PASSWORD"] {
		t.Error("expected DB_PASSWORD to be sensitive")
	}
	if !keySet["AUTH_TOKEN"] {
		t.Error("expected AUTH_TOKEN to be sensitive")
	}
	if keySet["APP_NAME"] {
		t.Error("expected APP_NAME to not be sensitive")
	}
}

func TestSensitiveKeys_EmptyMap(t *testing.T) {
	keys := SensitiveKeys(map[string]string{}, nil)
	if len(keys) != 0 {
		t.Errorf("expected no keys, got %d", len(keys))
	}
}
