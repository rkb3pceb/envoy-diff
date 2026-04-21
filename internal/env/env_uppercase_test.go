package env

import (
	"testing"
)

func TestUppercaseMap_NoOp_WhenBothDisabled(t *testing.T) {
	src := map[string]string{"foo": "bar", "baz": "qux"}
	opts := UppercaseOptions{Keys: false, Values: false}
	res := UppercaseMap(src, opts)
	if len(res.Changed) != 0 {
		t.Fatalf("expected no changes, got %v", res.Changed)
	}
	if res.Map["foo"] != "bar" {
		t.Errorf("unexpected value: %s", res.Map["foo"])
	}
}

func TestUppercaseMap_UppercasesKeys(t *testing.T) {
	src := map[string]string{"db_host": "localhost", "port": "5432"}
	opts := DefaultUppercaseOptions()
	res := UppercaseMap(src, opts)
	if _, ok := res.Map["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
	if _, ok := res.Map["PORT"]; !ok {
		t.Error("expected PORT in result")
	}
	if !HasUppercaseChanges(res) {
		t.Error("expected changes to be reported")
	}
}

func TestUppercaseMap_UppercasesValues(t *testing.T) {
	src := map[string]string{"MODE": "production", "LOG": "debug"}
	opts := UppercaseOptions{Keys: false, Values: true}
	res := UppercaseMap(src, opts)
	if res.Map["MODE"] != "PRODUCTION" {
		t.Errorf("expected PRODUCTION, got %s", res.Map["MODE"])
	}
	if res.Map["LOG"] != "DEBUG" {
		t.Errorf("expected DEBUG, got %s", res.Map["LOG"])
	}
}

func TestUppercaseMap_OnlyKeys_LimitsScope(t *testing.T) {
	src := map[string]string{"db_host": "localhost", "app_name": "myapp"}
	opts := UppercaseOptions{Keys: true, OnlyKeys: []string{"db_host"}}
	res := UppercaseMap(src, opts)
	if _, ok := res.Map["DB_HOST"]; !ok {
		t.Error("expected DB_HOST")
	}
	if _, ok := res.Map["app_name"]; !ok {
		t.Error("expected app_name to remain unchanged")
	}
	if len(res.Changed) != 1 {
		t.Errorf("expected 1 change, got %d", len(res.Changed))
	}
}

func TestUppercaseMap_DoesNotMutateInput(t *testing.T) {
	src := map[string]string{"key": "value"}
	opts := DefaultUppercaseOptions()
	_ = UppercaseMap(src, opts)
	if _, ok := src["key"]; !ok {
		t.Error("original map was mutated")
	}
}

func TestUppercaseMap_AlreadyUppercase_NoChange(t *testing.T) {
	src := map[string]string{"HOST": "localhost"}
	opts := DefaultUppercaseOptions()
	res := UppercaseMap(src, opts)
	if HasUppercaseChanges(res) {
		t.Error("expected no changes when keys already uppercase")
	}
}

func TestUppercaseMap_EmptyMap_ReturnsEmpty(t *testing.T) {
	res := UppercaseMap(map[string]string{}, DefaultUppercaseOptions())
	if len(res.Map) != 0 {
		t.Errorf("expected empty map, got %v", res.Map)
	}
}
