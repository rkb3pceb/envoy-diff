package env

// UppercaseOptions controls how keys and values are uppercased.
type UppercaseOptions struct {
	// UppercaseKeys converts all map keys to uppercase.
	UppercaseKeys bool
	// UppercaseValues converts all map values to uppercase.
	UppercaseValues bool
	// OnlyKeys limits the operation to a specific set of keys (empty = all).
	OnlyKeys []string
}

// DefaultUppercaseOptions returns UppercaseOptions with no-op defaults.
func DefaultUppercaseOptions() UppercaseOptions {
	return UppercaseOptions{}
}

// UppercaseMap returns a new map with keys and/or values uppercased
// according to opts. The original map is never mutated.
func UppercaseMap(src map[string]string, opts UppercaseOptions) map[string]string {
	if !opts.UppercaseKeys && !opts.UppercaseValues {
		out := make(map[string]string, len(src))
		for k, v := range src {
			out[k] = v
		}
		return out
	}

	scope := make(map[string]struct{}, len(opts.OnlyKeys))
	for _, k := range opts.OnlyKeys {
		scope[k] = struct{}{}
	}
	inScope := func(k string) bool {
		if len(scope) == 0 {
			return true
		}
		_, ok := scope[k]
		return ok
	}

	out := make(map[string]string, len(src))
	for k, v := range src {
		newKey := k
		newVal := v
		if inScope(k) {
			if opts.UppercaseKeys {
				newKey = asciiUpper(k)
			}
			if opts.UppercaseValues {
				newVal = asciiUpper(v)
			}
		}
		out[newKey] = newVal
	}
	return out
}

// HasUppercaseChanges reports whether UppercaseMap would produce a different
// map than the input under the given options.
func HasUppercaseChanges(src map[string]string, opts UppercaseOptions) bool {
	result := UppercaseMap(src, opts)
	if len(result) != len(src) {
		return true
	}
	for k, v := range result {
		orig, ok := src[k]
		if !ok || orig != v {
			return true
		}
	}
	return false
}

// asciiUpper uppercases ASCII letters only, leaving non-ASCII bytes untouched.
func asciiUpper(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}
	return string(b)
}
