package env

import "strings"

// FilterOptions controls which keys are included in the output map.
type FilterOptions struct {
	// Prefixes limits output to keys that start with any of the given prefixes.
	Prefixes []string

	// Contains limits output to keys whose name contains the given substring
	// (case-insensitive). Empty string disables this filter.
	Contains string

	// ExcludePrefixes removes keys that start with any of the given prefixes.
	ExcludePrefixes []string

	// OnlyNonEmpty removes keys whose value is the empty string.
	OnlyNonEmpty bool
}

// DefaultFilterOptions returns a FilterOptions with no restrictions applied.
func DefaultFilterOptions() FilterOptions {
	return FilterOptions{}
}

// FilterMap returns a new map containing only the entries from src that satisfy
// all active filter criteria in opts.
func FilterMap(src map[string]string, opts FilterOptions) map[string]string {
	out := make(map[string]string, len(src))

	for k, v := range src {
		if !matchesPrefixes(k, opts.Prefixes) {
			continue
		}
		if matchesAnyPrefix(k, opts.ExcludePrefixes) {
			continue
		}
		if opts.Contains != "" && !strings.Contains(strings.ToLower(k), strings.ToLower(opts.Contains)) {
			continue
		}
		if opts.OnlyNonEmpty && v == "" {
			continue
		}
		out[k] = v
	}

	return out
}

// HasFilteredKeys returns true when FilterMap would remove at least one key.
func HasFilteredKeys(src map[string]string, opts FilterOptions) bool {
	return len(FilterMap(src, opts)) < len(src)
}

func matchesPrefixes(key string, prefixes []string) bool {
	if len(prefixes) == 0 {
		return true
	}
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

func matchesAnyPrefix(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
