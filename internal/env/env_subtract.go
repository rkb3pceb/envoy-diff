package env

// SubtractOptions controls how SubtractMap removes keys.
type SubtractOptions struct {
	// Keys lists explicit keys to remove from base.
	Keys []string
	// Prefixes removes any key whose prefix matches one of these values.
	Prefixes []string
	// CaseInsensitive makes key and prefix matching case-insensitive.
	CaseInsensitive bool
}

// DefaultSubtractOptions returns a SubtractOptions with safe defaults.
func DefaultSubtractOptions() SubtractOptions {
	return SubtractOptions{}
}

// SubtractResult holds the output of SubtractMap.
type SubtractResult struct {
	// Result is the base map with matching keys removed.
	Result map[string]string
	// Removed lists the keys that were actually deleted.
	Removed []string
}

// HasSubtracted reports whether any keys were removed.
func (r SubtractResult) HasSubtracted() bool {
	return len(r.Removed) > 0
}

// SubtractMap removes keys from base according to opts.
// base is never mutated; a new map is returned inside SubtractResult.
func SubtractMap(base map[string]string, opts SubtractOptions) SubtractResult {
	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		if opts.CaseInsensitive {
			keySet[strings.ToLower(k)] = struct{}{}
		} else {
			keySet[k] = struct{}{}
		}
	}

	out := make(map[string]string, len(base))
	var removed []string

	for k, v := range base {
		compare := k
		if opts.CaseInsensitive {
			compare = strings.ToLower(k)
		}

		if _, hit := keySet[compare]; hit {
			removed = append(removed, k)
			continue
		}

		if matchesAnyPrefixCI(compare, opts.Prefixes, opts.CaseInsensitive) {
			removed = append(removed, k)
			continue
		}

		out[k] = v
	}

	sort.Strings(removed)
	return SubtractResult{Result: out, Removed: removed}
}

func matchesAnyPrefixCI(key string, prefixes []string, ci bool) bool {
	for _, p := range prefixes {
		pfx := p
		if ci {
			pfx = strings.ToLower(p)
		}
		if strings.HasPrefix(key, pfx) {
			return true
		}
	}
	return false
}
