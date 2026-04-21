package env

import "strings"

// DefaultSelectOptions returns a SelectOptions with safe defaults.
func DefaultSelectOptions() SelectOptions {
	return SelectOptions{
		CaseSensitive: true,
	}
}

// SelectOptions controls how keys are selected from an env map.
type SelectOptions struct {
	// Keys is the explicit list of keys to select.
	Keys []string

	// Prefixes selects all keys that start with any of these prefixes.
	Prefixes []string

	// CaseSensitive controls whether key matching is case-sensitive.
	CaseSensitive bool

	// Invert returns all keys NOT matching the selection criteria.
	Invert bool
}

// SelectMap returns a new map containing only the entries that match the
// provided SelectOptions. If no keys or prefixes are specified the original
// map is returned as-is (shallow copy).
func SelectMap(src map[string]string, opts SelectOptions) map[string]string {
	result := make(map[string]string, len(src))

	if len(opts.Keys) == 0 && len(opts.Prefixes) == 0 {
		for k, v := range src {
			result[k] = v
		}
		return result
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		norm := k
		if !opts.CaseSensitive {
			norm = strings.ToUpper(k)
		}
		keySet[norm] = struct{}{}
	}

	for k, v := range src {
		matched := selectMatches(k, keySet, opts)
		if opts.Invert {
			matched = !matched
		}
		if matched {
			result[k] = v
		}
	}
	return result
}

func selectMatches(key string, keySet map[string]struct{}, opts SelectOptions) bool {
	norm := key
	if !opts.CaseSensitive {
		norm = strings.ToUpper(key)
	}
	if _, ok := keySet[norm]; ok {
		return true
	}
	for _, p := range opts.Prefixes {
		pn := p
		if !opts.CaseSensitive {
			pn = strings.ToUpper(p)
		}
		if strings.HasPrefix(norm, pn) {
			return true
		}
	}
	return false
}

// HasSelectedKeys reports whether SelectMap produced any entries.
func HasSelectedKeys(m map[string]string) bool {
	return len(m) > 0
}
