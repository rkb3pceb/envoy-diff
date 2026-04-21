package env

import (
	"testing"
)

func TestSplitMap_NoOp_SingleValues(t *testing.T) {
	src := map[string]string{"A": "hello", "B": "world"}
	r, err := SplitMap(src, DefaultSplitOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if HasSplitChanges(r) {
		t.Error("expected no split changes for single-value keys")
	}
	if r.Map["A"] != "hello" || r.Map["B"] != "world" {
		t.Error("expected original values preserved")
	}
}

func TestSplitMap_SplitsCommaDelimited(t *testing.T) {
	src := map[string]string{"HOSTS": "a,b,c"}
	opts := DefaultSplitOptions()
	r, err := SplitMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasSplitChanges(r) {
		t.Fatal("expected split changes")
	}
	if r.Map["HOSTS_1"] != "a" || r.Map["HOSTS_2"] != "b" || r.Map["HOSTS_3"] != "c" {
		t.Errorf("unexpected map: %v", r.Map)
	}
	if _, ok := r.Map["HOSTS"]; ok {
		t.Error("original key should be removed when KeepOriginal=false")
	}
}

func TestSplitMap_KeepOriginal(t *testing.T) {
	src := map[string]string{"TAGS": "x,y"}
	opts := DefaultSplitOptions()
	opts.KeepOriginal = true
	r, err := SplitMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Map["TAGS"] != "x,y" {
		t.Error("expected original key to be kept")
	}
	if r.Map["TAGS_1"] != "x" || r.Map["TAGS_2"] != "y" {
		t.Errorf("unexpected split values: %v", r.Map)
	}
}

func TestSplitMap_ZeroIndexBase(t *testing.T) {
	src := map[string]string{"V": "p,q"}
	opts := DefaultSplitOptions()
	opts.IndexBase = 0
	r, err := SplitMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Map["V_0"] != "p" || r.Map["V_1"] != "q" {
		t.Errorf("unexpected keys: %v", r.Map)
	}
}

func TestSplitMap_SpecificKeys_OnlySplitsThem(t *testing.T) {
	src := map[string]string{"A": "x,y", "B": "m,n"}
	opts := DefaultSplitOptions()
	opts.Keys = []string{"A"}
	r, err := SplitMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Map["B"] != "m,n" {
		t.Error("expected B to remain unsplit")
	}
	if r.Map["A_1"] != "x" {
		t.Error("expected A to be split")
	}
}

func TestSplitMap_SkipEmpty_DropsBlankParts(t *testing.T) {
	src := map[string]string{"LIST": "a,,b"}
	opts := DefaultSplitOptions()
	opts.SkipEmpty = true
	r, err := SplitMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.Map["LIST_3"]; ok {
		t.Error("expected no third key when empty parts are skipped")
	}
	if r.Map["LIST_1"] != "a" || r.Map["LIST_2"] != "b" {
		t.Errorf("unexpected values: %v", r.Map)
	}
}

func TestSplitMap_EmptyDelimiter_ReturnsError(t *testing.T) {
	opts := DefaultSplitOptions()
	opts.Delimiter = ""
	_, err := SplitMap(map[string]string{"K": "v"}, opts)
	if err == nil {
		t.Error("expected error for empty delimiter")
	}
}

func TestSplitMap_CustomDelimiter(t *testing.T) {
	src := map[string]string{"PORTS": "80|443|8080"}
	opts := DefaultSplitOptions()
	opts.Delimiter = "|"
	r, err := SplitMap(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Map["PORTS_1"] != "80" || r.Map["PORTS_2"] != "443" || r.Map["PORTS_3"] != "8080" {
		t.Errorf("unexpected map: %v", r.Map)
	}
}
