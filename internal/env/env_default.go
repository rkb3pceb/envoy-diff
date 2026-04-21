package env

// DefaultOptions controls how default values are applied to an env map.
type DefaultOptions struct {
	// Defaults is a map of key -> default value.
	// A default is applied only when the key is absent or (optionally) empty.
	Defaults map[string]string

	// OverwriteEmpty replaces existing keys whose value is the empty string.
	OverwriteEmpty bool

	// OverwriteAll replaces every key unconditionally (acts as a forced merge).
	OverwriteAll bool
}

// DefaultDefaultOptions returns a DefaultOptions with safe zero-value settings.
func DefaultDefaultOptions() DefaultOptions {
	return DefaultOptions{
		Defaults:       map[string]string{},
		OverwriteEmpty: false,
		OverwriteAll:   false,
	}
}

// DefaultResult carries the outcome of ApplyDefaults.
type DefaultResult struct {
	// Applied lists keys whose value was set from the defaults map.
	Applied []string
	// Skipped lists keys that already had a value and were not overwritten.
	Skipped []string
}

// HasApplied returns true when at least one default was applied.
func (r DefaultResult) HasApplied() bool { return len(r.Applied) > 0 }

// ApplyDefaults fills missing (or empty) keys in src using opts.Defaults.
// It never mutates src; the returned map is always a new copy.
func ApplyDefaults(src map[string]string, opts DefaultOptions) (map[string]string, DefaultResult) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}

	var result DefaultResult

	for k, def := range opts.Defaults {
		existing, exists := out[k]
		switch {
		case opts.OverwriteAll:
			out[k] = def
			result.Applied = append(result.Applied, k)
		case !exists:
			out[k] = def
			result.Applied = append(result.Applied, k)
		case opts.OverwriteEmpty && existing == "":
			out[k] = def
			result.Applied = append(result.Applied, k)
		default:
			result.Skipped = append(result.Skipped, k)
		}
	}

	sortStringsInPlace(result.Applied)
	sortStringsInPlace(result.Skipped)
	return out, result
}

func sortStringsInPlace(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
