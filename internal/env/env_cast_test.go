package env

import (
	"testing"
)

func TestCastMap_NoOp_StringType(t *testing.T) {
	src := map[string]string{"HOST": "localhost", "PORT": "8080"}
	opts := DefaultCastOptions()
	out, results, err := CastMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", out["PORT"])
	}
}

func TestCastMap_IntType_Valid(t *testing.T) {
	src := map[string]string{"PORT": "8080", "TIMEOUT": "30"}
	opts := DefaultCastOptions()
	opts.TargetType = CastInt
	out, results, err := CastMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if HasCastErrors(results) {
		t.Error("expected no cast errors")
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", out["PORT"])
	}
}

func TestCastMap_IntType_Invalid_SkipOnError(t *testing.T) {
	src := map[string]string{"PORT": "not-a-number"}
	opts := DefaultCastOptions()
	opts.TargetType = CastInt
	opts.SkipOnError = true
	out, results, err := CastMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasCastErrors(results) {
		t.Error("expected cast errors in results")
	}
	// original value preserved on skip
	if out["PORT"] != "not-a-number" {
		t.Errorf("expected original value preserved, got %s", out["PORT"])
	}
}

func TestCastMap_IntType_Invalid_ErrorOnFailure(t *testing.T) {
	src := map[string]string{"PORT": "abc"}
	opts := DefaultCastOptions()
	opts.TargetType = CastInt
	opts.SkipOnError = false
	_, _, err := CastMap(src, opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCastMap_BoolType_Normalises(t *testing.T) {
	src := map[string]string{"ENABLED": "yes", "VERBOSE": "0", "DEBUG": "TRUE"}
	opts := DefaultCastOptions()
	opts.TargetType = CastBool
	out, results, err := CastMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if HasCastErrors(results) {
		t.Error("expected no cast errors")
	}
	if out["ENABLED"] != "true" {
		t.Errorf("expected true, got %s", out["ENABLED"])
	}
	if out["VERBOSE"] != "false" {
		t.Errorf("expected false, got %s", out["VERBOSE"])
	}
	if out["DEBUG"] != "true" {
		t.Errorf("expected true, got %s", out["DEBUG"])
	}
}

func TestCastMap_SpecificKeys_OnlyCastsThem(t *testing.T) {
	src := map[string]string{"PORT": "8080", "NAME": "hello"}
	opts := DefaultCastOptions()
	opts.TargetType = CastInt
	opts.Keys = []string{"PORT"}
	_, results, err := CastMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "PORT" {
		t.Errorf("expected result for PORT, got %s", results[0].Key)
	}
}

func TestCastMap_TrimSpace_StripsBeforeCast(t *testing.T) {
	src := map[string]string{"PORT": "  8080  "}
	opts := DefaultCastOptions()
	opts.TargetType = CastInt
	opts.TrimSpace = true
	_, results, err := CastMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if HasCastErrors(results) {
		t.Errorf("expected no errors after trim, got error: %v", results[0].Err)
	}
}

func TestCastMap_FloatType_Valid(t *testing.T) {
	src := map[string]string{"RATIO": "3.14", "SCALE": "1.0"}
	opts := DefaultCastOptions()
	opts.TargetType = CastFloat
	_, results, err := CastMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if HasCastErrors(results) {
		t.Error("expected no cast errors for valid floats")
	}
}
