package redact

import (
	"testing"
)

func TestDefaultContext_EnabledByDefault(t *testing.T) {
	ctx := DefaultContext()
	if !ctx.Enabled() {
		t.Fatal("expected DefaultContext to be enabled")
	}
}

func TestNewContext_DisabledSkipsRedaction(t *testing.T) {
	ctx := NewContext(nil, MaskFull, false)
	result := ctx.Apply("SECRET_KEY", "super-secret")
	if result != "super-secret" {
		t.Errorf("expected original value when disabled, got %q", result)
	}
}

func TestContext_ShouldRedact_SensitiveKey(t *testing.T) {
	ctx := DefaultContext()
	if !ctx.ShouldRedact("DATABASE_PASSWORD") {
		t.Error("expected DATABASE_PASSWORD to be redacted")
	}
}

func TestContext_ShouldRedact_SafeKey(t *testing.T) {
	ctx := DefaultContext()
	if ctx.ShouldRedact("APP_ENV") {
		t.Error("expected APP_ENV to not be redacted")
	}
}

func TestContext_Apply_MasksValue(t *testing.T) {
	ctx := NewContext(nil, MaskFull, true)
	result := ctx.Apply("API_SECRET", "abc123")
	if result == "abc123" {
		t.Error("expected value to be masked")
	}
}

func TestContext_Apply_CustomRuleUsesPlaceholder(t *testing.T) {
	rs := NewRuleSet([]Rule{
		{Pattern: "INTERNAL_TOKEN", Placeholder: "[TOKEN]"},
	})
	ctx := NewContext(rs, MaskFull, true)
	result := ctx.Apply("INTERNAL_TOKEN", "xyz")
	if result != "[TOKEN]" {
		t.Errorf("expected '[TOKEN]', got %q", result)
	}
}

func TestContext_ApplyMap_RedactsOnlySensitive(t *testing.T) {
	ctx := DefaultContext()
	env := map[string]string{
		"APP_ENV":    "production",
		"DB_PASSWORD": "s3cr3t",
	}
	out := ctx.ApplyMap(env)
	if out["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV unchanged, got %q", out["APP_ENV"])
	}
	if out["DB_PASSWORD"] == "s3cr3t" {
		t.Error("expected DB_PASSWORD to be masked")
	}
}

func TestContext_ApplyMap_DisabledPassesThrough(t *testing.T) {
	ctx := NewContext(nil, MaskFull, false)
	env := map[string]string{
		"DB_PASSWORD": "s3cr3t",
	}
	out := ctx.ApplyMap(env)
	if out["DB_PASSWORD"] != "s3cr3t" {
		t.Errorf("expected original value when disabled, got %q", out["DB_PASSWORD"])
	}
}
