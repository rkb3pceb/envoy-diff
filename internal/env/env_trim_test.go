package env

import (
	"testing"
)

func TestTrimMap_NoOp_DefaultOptions(t *testing.T) {
	input := map[string]string{"KEY": "value", "OTHER": "clean"}
	opts := DefaultTrimOptions()
	r := TrimMap(input, opts)
	if r.Output["KEY"] != "value" {
		t.Errorf("expected 'value', got %q", r.Output["KEY"])
	}
	if HasTrimChanges(r) {
		t.Error("expected no changes")
	}
}

func TestTrimMap_TrimValues_RemovesWhitespace(t *testing.T) {
	input := map[string]string{"KEY": "  hello  ", "B": "\tworld\n"}
	opts := DefaultTrimOptions()
	r := TrimMap(input, opts)
	if r.Output["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", r.Output["KEY"])
	}
	if r.Output["B"] != "world" {
		t.Errorf("expected 'world', got %q", r.Output["B"])
	}
	if !HasTrimChanges(r) {
		t.Error("expected changes")
	}
}

func TestTrimMap_TrimKeys_RemovesWhitespace(t *testing.T) {
	input := map[string]string{" KEY ": "value"}
	opts := DefaultTrimOptions()
	opts.TrimKeys = true
	r := TrimMap(input, opts)
	if _, ok := r.Output["KEY"]; !ok {
		t.Error("expected trimmed key 'KEY' to exist")
	}
}

func TestTrimMap_TrimPrefix_RemovesFromValues(t *testing.T) {
	input := map[string]string{"URL": "https://example.com", "OTHER": "plain"}
	opts := TrimOptions{TrimPrefix: "https://"}
	r := TrimMap(input, opts)
	if r.Output["URL"] != "example.com" {
		t.Errorf("expected 'example.com', got %q", r.Output["URL"])
	}
	if r.Output["OTHER"] != "plain" {
		t.Errorf("expected 'plain', got %q", r.Output["OTHER"])
	}
}

func TestTrimMap_TrimSuffix_RemovesFromValues(t *testing.T) {
	input := map[string]string{"KEY": "value/", "CLEAN": "ok"}
	opts := TrimOptions{TrimSuffix: "/"}
	r := TrimMap(input, opts)
	if r.Output["KEY"] != "value" {
		t.Errorf("expected 'value', got %q", r.Output["KEY"])
	}
}

func TestTrimMap_TrimChars_RemovesBothEnds(t *testing.T) {
	input := map[string]string{"KEY": "***secret***"}
	opts := TrimOptions{TrimChars: "*"}
	r := TrimMap(input, opts)
	if r.Output["KEY"] != "secret" {
		t.Errorf("expected 'secret', got %q", r.Output["KEY"])
	}
}

func TestTrimMap_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"KEY": "  spaced  "}
	opts := DefaultTrimOptions()
	TrimMap(input, opts)
	if input["KEY"] != "  spaced  " {
		t.Error("input map was mutated")
	}
}
