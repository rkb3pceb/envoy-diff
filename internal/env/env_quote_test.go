package env

import (
	"testing"
)

func TestQuoteMap_NoOp_WhenStyleNone(t *testing.T) {
	m := map[string]string{"KEY": "hello world", "B": "plain"}
	out, err := QuoteMap(m, QuoteOptions{Style: QuoteStyleNone})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "hello world" {
		t.Errorf("expected unquoted value, got %q", out["KEY"])
	}
}

func TestQuoteMap_DoubleQuotesAllValues(t *testing.T) {
	m := map[string]string{"A": "foo", "B": "bar baz"}
	out, err := QuoteMap(m, QuoteOptions{Style: QuoteStyleDouble})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range out {
		if len(v) < 2 || v[0] != '"' || v[len(v)-1] != '"' {
			t.Errorf("key %q: expected double-quoted value, got %q", k, v)
		}
	}
}

func TestQuoteMap_SingleQuotesAllValues(t *testing.T) {
	m := map[string]string{"X": "value"}
	out, err := QuoteMap(m, QuoteOptions{Style: QuoteStyleSingle})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X"] != "'value'" {
		t.Errorf("expected 'value', got %q", out["X"])
	}
}

func TestQuoteMap_Auto_QuotesOnlyWhenNeeded(t *testing.T) {
	m := map[string]string{"PLAIN": "nospecial", "SPACED": "has space"}
	out, err := QuoteMap(m, QuoteOptions{Style: QuoteStyleAuto})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PLAIN"] != "nospecial" {
		t.Errorf("plain value should not be quoted, got %q", out["PLAIN"])
	}
	if out["SPACED"] != `"has space"` {
		t.Errorf("spaced value should be double-quoted, got %q", out["SPACED"])
	}
}

func TestQuoteMap_SkipEmpty_LeavesEmptyUntouched(t *testing.T) {
	m := map[string]string{"EMPTY": ""}
	out, err := QuoteMap(m, QuoteOptions{Style: QuoteStyleDouble, SkipEmpty: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["EMPTY"] != "" {
		t.Errorf("expected empty string, got %q", out["EMPTY"])
	}
}

func TestQuoteMap_SpecificKeys_OnlyQuotesThem(t *testing.T) {
	m := map[string]string{"A": "alpha", "B": "beta"}
	out, err := QuoteMap(m, QuoteOptions{Style: QuoteStyleDouble, Keys: []string{"A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != `"alpha"` {
		t.Errorf("expected A to be quoted, got %q", out["A"])
	}
	if out["B"] != "beta" {
		t.Errorf("expected B to be unchanged, got %q", out["B"])
	}
}

func TestQuoteMap_InvalidStyle_ReturnsError(t *testing.T) {
	m := map[string]string{"K": "v"}
	_, err := QuoteMap(m, QuoteOptions{Style: "bogus"})
	if err == nil {
		t.Fatal("expected error for unknown style")
	}
}

func TestHasQuoteChanges_DetectsChange(t *testing.T) {
	orig := map[string]string{"K": "v"}
	quoted := map[string]string{"K": `"v"`}
	if !HasQuoteChanges(orig, quoted) {
		t.Error("expected changes to be detected")
	}
}

func TestHasQuoteChanges_NoChange(t *testing.T) {
	orig := map[string]string{"K": "v"}
	if HasQuoteChanges(orig, orig) {
		t.Error("expected no changes")
	}
}
