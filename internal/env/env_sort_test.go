package env

import (
	"testing"
)

func TestSortedKeys_AscendingOrder(t *testing.T) {
	m := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	keys := SortedKeys(m, DefaultSortOptions())
	want := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("index %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSortedKeys_DescendingOrder(t *testing.T) {
	m := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	opts := SortOptions{Order: SortDesc, CaseInsensitive: true}
	keys := SortedKeys(m, opts)
	want := []string{"ZEBRA", "MANGO", "APPLE"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("index %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSortedKeys_ByValue(t *testing.T) {
	m := map[string]string{"A": "zebra", "B": "apple", "C": "mango"}
	opts := SortOptions{Order: SortByValue, CaseInsensitive: true}
	keys := SortedKeys(m, opts)
	want := []string{"B", "C", "A"} // apple, mango, zebra
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("index %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSortedKeys_CaseSensitive(t *testing.T) {
	m := map[string]string{"b": "1", "A": "2"}
	opts := SortOptions{Order: SortAsc, CaseInsensitive: false}
	keys := SortedKeys(m, opts)
	// uppercase 'A' (65) < lowercase 'b' (98) in ASCII
	if keys[0] != "A" || keys[1] != "b" {
		t.Errorf("expected [A b], got %v", keys)
	}
}

func TestSortedKeys_EmptyMap(t *testing.T) {
	keys := SortedKeys(map[string]string{}, DefaultSortOptions())
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %v", keys)
	}
}

func TestSortedMap_ReturnsPairs(t *testing.T) {
	m := map[string]string{"B": "beta", "A": "alpha"}
	pairs := SortedMap(m, DefaultSortOptions())
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	if pairs[0][0] != "A" || pairs[0][1] != "alpha" {
		t.Errorf("first pair: got %v, want [A alpha]", pairs[0])
	}
	if pairs[1][0] != "B" || pairs[1][1] != "beta" {
		t.Errorf("second pair: got %v, want [B beta]", pairs[1])
	}
}

func TestSortedMap_PreservesValues(t *testing.T) {
	m := map[string]string{"KEY": "value", "OTHER": "data"}
	pairs := SortedMap(m, DefaultSortOptions())
	for _, p := range pairs {
		if m[p[0]] != p[1] {
			t.Errorf("value mismatch for key %q: got %q, want %q", p[0], p[1], m[p[0]])
		}
	}
}
