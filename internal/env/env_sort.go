package env

import (
	"sort"
	"strings"
)

// SortOrder defines the ordering strategy for environment variable keys.
type SortOrder int

const (
	// SortAsc sorts keys in ascending (A-Z) order.
	SortAsc SortOrder = iota
	// SortDesc sorts keys in descending (Z-A) order.
	SortDesc
	// SortByValue sorts keys by their associated value, ascending.
	SortByValue
)

// DefaultSortOptions returns a SortOptions with ascending key order.
func DefaultSortOptions() SortOptions {
	return SortOptions{
		Order:           SortAsc,
		CaseInsensitive: true,
	}
}

// SortOptions controls how SortedKeys behaves.
type SortOptions struct {
	// Order specifies the sort direction or strategy.
	Order SortOrder
	// CaseInsensitive performs case-folded comparison when true.
	CaseInsensitive bool
}

// compareStrings returns true if a should come before b, given the
// case-insensitive flag and whether the ordering is ascending.
func compareStrings(a, b string, caseInsensitive, ascending bool) bool {
	if caseInsensitive {
		a, b = strings.ToLower(a), strings.ToLower(b)
	}
	if ascending {
		return a < b
	}
	return a > b
}

// SortedKeys returns the keys of m ordered according to opts.
func SortedKeys(m map[string]string, opts SortOptions) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	switch opts.Order {
	case SortDesc:
		sort.Slice(keys, func(i, j int) bool {
			return compareStrings(keys[i], keys[j], opts.CaseInsensitive, false)
		})
	case SortByValue:
		sort.Slice(keys, func(i, j int) bool {
			return compareStrings(m[keys[i]], m[keys[j]], opts.CaseInsensitive, true)
		})
	default: // SortAsc
		sort.Slice(keys, func(i, j int) bool {
			return compareStrings(keys[i], keys[j], opts.CaseInsensitive, true)
		})
	}

	return keys
}

// SortedMap returns a slice of [2]string pairs {key, value} ordered by opts.
func SortedMap(m map[string]string, opts SortOptions) [][2]string {
	keys := SortedKeys(m, opts)
	pairs := make([][2]string, len(keys))
	for i, k := range keys {
		pairs[i] = [2]string{k, m[k]}
	}
	return pairs
}
