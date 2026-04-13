package env

import (
	"strings"
	"testing"
)

func TestMaskMap_NoSensitiveKeys_PassesThrough(t *testing.T) {
	m := map[string]string{
		"APP_NAME": "envoy",
		"PORT":     "8080",
	}
	res := MaskMap(m, DefaultMaskOptions())
	if res.Map["APP_NAME"] != "envoy" {
		t.Errorf("expected APP_NAME unchanged, got %q", res.Map["APP_NAME"])
	}
	if len(res.MaskedKeys) != 0 {
		t.Errorf("expected no masked keys, got %v", res.MaskedKeys)
	}
}

func TestMaskMap_SensitiveKey_FullMask(t *testing.T) {
	m := map[string]string{
		"DB_PASSWORD": "supersecret",
	}
	opts := DefaultMaskOptions()
	opts.Level = MaskFull
	res := MaskMap(m, opts)
	if res.Map["DB_PASSWORD"] != opts.Placeholder {
		t.Errorf("expected placeholder, got %q", res.Map["DB_PASSWORD"])
	}
	if len(res.MaskedKeys) != 1 || res.MaskedKeys[0] != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD in masked keys, got %v", res.MaskedKeys)
	}
}

func TestMaskMap_SensitiveKey_PartialMask_LongValue(t *testing.T) {
	m := map[string]string{
		"API_TOKEN": "abcdefghij",
	}
	opts := DefaultMaskOptions()
	opts.Level = MaskPartial
	res := MaskMap(m, opts)
	got := res.Map["API_TOKEN"]
	if !strings.HasPrefix(got, "ab") || !strings.HasSuffix(got, "ij") {
		t.Errorf("expected partial mask with edges, got %q", got)
	}
	if strings.Contains(got, "cdefgh") {
		t.Errorf("expected middle to be masked, got %q", got)
	}
}

func TestMaskMap_SensitiveKey_PartialMask_ShortValue(t *testing.T) {
	m := map[string]string{
		"API_SECRET": "abc",
	}
	opts := DefaultMaskOptions()
	opts.Level = MaskPartial
	res := MaskMap(m, opts)
	if res.Map["API_SECRET"] != opts.Placeholder {
		t.Errorf("expected placeholder for short value, got %q", res.Map["API_SECRET"])
	}
}

func TestMaskMap_ExtraPatterns_MatchCustomKey(t *testing.T) {
	m := map[string]string{
		"DEPLOY_PASSPHRASE": "hunter2",
		"REGION":            "us-east-1",
	}
	opts := DefaultMaskOptions()
	opts.ExtraPatterns = []string{"passphrase"}
	res := MaskMap(m, opts)
	if res.Map["DEPLOY_PASSPHRASE"] != opts.Placeholder {
		t.Errorf("expected custom pattern to mask value")
	}
	if res.Map["REGION"] != "us-east-1" {
		t.Errorf("expected REGION unchanged")
	}
}

func TestMaskMap_DoesNotMutateInput(t *testing.T) {
	m := map[string]string{
		"DB_PASSWORD": "original",
	}
	_ = MaskMap(m, DefaultMaskOptions())
	if m["DB_PASSWORD"] != "original" {
		t.Errorf("input map was mutated")
	}
}

func TestMaskMap_CaseInsensitivePattern(t *testing.T) {
	m := map[string]string{
		"db_password": "lowercase",
	}
	res := MaskMap(m, DefaultMaskOptions())
	if res.Map["db_password"] != DefaultMaskOptions().Placeholder {
		t.Errorf("expected lowercase key to be masked")
	}
}
