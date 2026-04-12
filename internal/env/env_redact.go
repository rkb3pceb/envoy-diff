package env

import (
	"github.com/yourorg/envoy-diff/internal/redact"
)

// RedactOptions controls how sensitive values are masked in an env map.
type RedactOptions struct {
	// Enabled toggles redaction. When false, values are returned as-is.
	Enabled bool
	// Level controls masking depth: "full" replaces entirely, "partial" shows edges.
	Level string
	// ExtraPatterns are additional substring patterns treated as sensitive.
	ExtraPatterns []string
}

// DefaultRedactOptions returns sensible defaults for redaction.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		Enabled: true,
		Level:   "partial",
	}
}

// RedactMap returns a copy of m with sensitive values masked according to opts.
// Keys are never modified. Only values whose keys match sensitive patterns are
// replaced; all other values pass through unchanged.
func RedactMap(m map[string]string, opts RedactOptions) map[string]string {
	out := make(map[string]string, len(m))

	if !opts.Enabled {
		for k, v := range m {
			out[k] = v
		}
		return out
	}

	patterns := append(redact.PatternList, opts.ExtraPatterns...)
	rs := redact.NewRuleSet(patterns, redact.Placeholder)

	for k, v := range m {
		if rs.Matches(k) != nil {
			out[k] = redact.Mask(v, opts.Level)
		} else {
			out[k] = v
		}
	}
	return out
}

// SensitiveKeys returns the subset of keys in m that are considered sensitive.
func SensitiveKeys(m map[string]string, extraPatterns []string) []string {
	patterns := append(redact.PatternList, extraPatterns...)
	rs := redact.NewRuleSet(patterns, redact.Placeholder)

	var keys []string
	for k := range m {
		if rs.Matches(k) != nil {
			keys = append(keys, k)
		}
	}
	return keys
}
