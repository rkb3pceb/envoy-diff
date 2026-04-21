package env

import (
	"testing"
)

func TestIntersectMaps_NoMaps_ReturnsEmpty(t *testing.T) {
	r := IntersectMaps(DefaultIntersectOptions())
	if len(r.Map) != 0 {
		t.Fatalf("expected empty map, got %v", r.Map)
	}
}

func TestIntersectMaps_SingleMap_ReturnsCopy(t *testing.T) {
	m := map[string]string{"A": "1", "B": "2"}
	r := IntersectMaps(DefaultIntersectOptions(), m)
	if len(r.Map) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(r.Map))
	}
}

func TestIntersectMaps_CommonKeys_Returned(t *testing.T) {
	a := map[string]string{"X": "1", "Y": "2", "Z": "3"}
	b := map[string]string{"X": "9", "Y": "2"}
	r := IntersectMaps(DefaultIntersectOptions(), a, b)
	if _, ok := r.Map["X"]; !ok {
		t.Error("expected X in result")
	}
	if _, ok := r.Map["Y"]; !ok {
		t.Error("expected Y in result")
	}
	if _, ok := r.Map["Z"]; ok {
		t.Error("did not expect Z in result")
	}
}

func TestIntersectMaps_KeepLast_UsesLastValue(t *testing.T) {
	a := map[string]string{"K": "first"}
	b := map[string]string{"K": "last"}
	opts := DefaultIntersectOptions()
	opts.KeepValues = "last"
	r := IntersectMaps(opts, a, b)
	if r.Map["K"] != "last" {
		t.Errorf("expected 'last', got %q", r.Map["K"])
	}
}

func TestIntersectMaps_KeepFirst_UsesFirstValue(t *testing.T) {
	a := map[string]string{"K": "first"}
	b := map[string]string{"K": "last"}
	opts := DefaultIntersectOptions()
	opts.KeepValues = "first"
	r := IntersectMaps(opts, a, b)
	if r.Map["K"] != "first" {
		t.Errorf("expected 'first', got %q", r.Map["K"])
	}
}

func TestIntersectMaps_RequireEqual_ExcludesConflicts(t *testing.T) {
	a := map[string]string{"SAME": "v", "DIFF": "a"}
	b := map[string]string{"SAME": "v", "DIFF": "b"}
	opts := DefaultIntersectOptions()
	opts.RequireEqual = true
	r := IntersectMaps(opts, a, b)
	if _, ok := r.Map["SAME"]; !ok {
		t.Error("expected SAME in result")
	}
	if _, ok := r.Map["DIFF"]; ok {
		t.Error("did not expect DIFF in result due to conflict")
	}
	if r.Conflicts != 1 {
		t.Errorf("expected 1 conflict, got %d", r.Conflicts)
	}
}

func TestIntersectMaps_ThreeMaps_AllMustContainKey(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"A": "1", "B": "2", "C": "3"}
	c := map[string]string{"A": "1"}
	r := IntersectMaps(DefaultIntersectOptions(), a, b, c)
	if _, ok := r.Map["A"]; !ok {
		t.Error("expected A")
	}
	if _, ok := r.Map["B"]; ok {
		t.Error("B should be absent; not in third map")
	}
}

func TestHasIntersectResult_TrueWhenNonEmpty(t *testing.T) {
	r := IntersectResult{Map: map[string]string{"K": "v"}}
	if !HasIntersectResult(r) {
		t.Error("expected true")
	}
}

func TestHasIntersectResult_FalseWhenEmpty(t *testing.T) {
	r := IntersectResult{Map: map[string]string{}}
	if HasIntersectResult(r) {
		t.Error("expected false")
	}
}
