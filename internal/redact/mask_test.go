package redact

import (
	"strings"
	"testing"
)

func TestMask_EmptyValue_ReturnsPlaceholder(t *testing.T) {
	got := Mask("", MaskFull)
	if got != Placeholder {
		t.Errorf("expected %q, got %q", Placeholder, got)
	}
}

func TestMask_FullLevel_ReturnsPlaceholder(t *testing.T) {
	got := Mask("supersecret", MaskFull)
	if got != Placeholder {
		t.Errorf("expected %q, got %q", Placeholder, got)
	}
}

func TestMask_PartialLevel_LongValue_ShowsEdges(t *testing.T) {
	value := "abcdefgh"
	got := Mask(value, MaskPartial)
	if !strings.HasPrefix(got, "ab") {
		t.Errorf("expected prefix 'ab', got %q", got)
	}
	if !strings.HasSuffix(got, "gh") {
		t.Errorf("expected suffix 'gh', got %q", got)
	}
	if !strings.Contains(got, "****") {
		t.Errorf("expected masked middle in %q", got)
	}
}

func TestMask_PartialLevel_ShortValue_ReturnsPlaceholder(t *testing.T) {
	got := Mask("abc", MaskPartial)
	if got != Placeholder {
		t.Errorf("short value should use placeholder, got %q", got)
	}
}

func TestMaskMap_MasksSensitiveKeys(t *testing.T) {
	input := map[string]string{
		"DATABASE_PASSWORD": "hunter2secret",
		"APP_NAME":          "myapp",
	}
	out := MaskMap(input, MaskFull)
	if out["DATABASE_PASSWORD"] != Placeholder {
		t.Errorf("expected sensitive key to be masked, got %q", out["DATABASE_PASSWORD"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected non-sensitive key to pass through, got %q", out["APP_NAME"])
	}
}

func TestMaskMap_PreservesAllKeys(t *testing.T) {
	input := map[string]string{
		"SECRET_KEY": "topsecret",
		"PORT":       "8080",
		"HOST":       "localhost",
	}
	out := MaskMap(input, MaskFull)
	if len(out) != len(input) {
		t.Errorf("expected %d keys, got %d", len(input), len(out))
	}
}
