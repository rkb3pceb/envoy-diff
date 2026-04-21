package env

import "sort"

// DefaultUniqueOptions returns a UniqueOptions with safe defaults.
func DefaultUniqueOptions() UniqueOptions {
	return UniqueOptions{
		ByValue:       true,
		CaseSensitive: true,
		KeepFirst:     true,
	}
}

// UniqueOptions controls how UniqueMap deduplicates entries.
type UniqueOptions struct {
	// ByValue removes keys whose values are duplicated across the map.
	ByValue bool

	// CaseSensitive controls whether value comparison is case-sensitive.
	CaseSensitive bool

	// KeepFirst retains the first key encountered for a duplicate value.
	// When false the last key (alphabetical order) is kept.
	KeepFirst bool
}

// UniqueResult holds the output of UniqueMap.
type UniqueResult struct {
	// Map contains the deduplicated key-value pairs.
	Map map[string]string

	// Removed lists keys that were dropped because their value was a duplicate.
	Removed []string
}

// UniqueMap returns a copy of m with duplicate values removed according to opts.
// Keys are evaluated in sorted order so results are deterministic.
func UniqueMap(m map[string]string, opts UniqueOptions) UniqueResult {
	result := UniqueResult{
		Map: make(map[string]string, len(m)),
	}

	if !opts.ByValue {
		for k, v := range m {
			result.Map[k] = v
		}
		return result
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// seen maps normalised value → winning key
	seen := make(map[string]string, len(m))

	for _, k := range keys {
		v := m[k]
		norm := v
		if !opts.CaseSensitive {
			norm = asciiLower(v)
		}

		if winner, dup := seen[norm]; dup {
			if opts.KeepFirst {
				// current key loses
				result.Removed = append(result.Removed, k)
			} else {
				// current key wins; evict previous winner
				delete(result.Map, winner)
				result.Removed = append(result.Removed, winner)
				result.Map[k] = v
				seen[norm] = k
			}
		} else {
			seen[norm] = k
			result.Map[k] = v
		}
	}

	sort.Strings(result.Removed)
	return result
}

// HasUniqueChanges reports whether any keys were removed.
func HasUniqueChanges(r UniqueResult) bool {
	return len(r.Removed) > 0
}

// asciiLower returns s with ASCII uppercase letters lowercased.
func asciiLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}
