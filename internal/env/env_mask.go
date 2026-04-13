package env

import (
	"strings"
)

// MaskLevel controls how much of a sensitive value is revealed.
type MaskLevel string

const (
	// MaskFull replaces the entire value with the placeholder.
	MaskFull MaskLevel = "full"
	// MaskPartial reveals the first and last two characters.
	MaskPartial MaskLevel = "partial"
)

// DefaultMaskOptions returns sensible defaults for MaskMap.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		Level:       MaskFull,
		Placeholder: "[REDACTED]",
		Patterns:    defaultSensitivePatterns(),
	}
}

// MaskOptions controls which keys are masked and how.
type MaskOptions struct {
	// Level determines how much of the value is hidden.
	Level MaskLevel
	// Placeholder is substituted for fully masked values.
	Placeholder string
	// Patterns is a list of substring patterns (case-insensitive) that
	// identify sensitive keys.
	Patterns []string
	// ExtraPatterns are appended to Patterns at call time.
	ExtraPatterns []string
}

// MaskResult holds the masked map and which keys were masked.
type MaskResult struct {
	Map        map[string]string
	MaskedKeys []string
}

// MaskMap returns a copy of m with sensitive values masked according to opts.
func MaskMap(m map[string]string, opts MaskOptions) MaskResult {
	patterns := append(opts.Patterns, opts.ExtraPatterns...)
	out := make(map[string]string, len(m))
	var masked []string

	for k, v := range m {
		if isSensitiveByPatterns(k, patterns) {
			masked = append(masked, k)
			out[k] = maskValue(v, opts)
		} else {
			out[k] = v
		}
	}
	return MaskResult{Map: out, MaskedKeys: masked}
}

func isSensitiveByPatterns(key string, patterns []string) bool {
	lower := strings.ToLower(key)
	for _, p := range patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

func maskValue(v string, opts MaskOptions) string {
	if opts.Level == MaskPartial && len(v) > 6 {
		return v[:2] + strings.Repeat("*", len(v)-4) + v[len(v)-2:]
	}
	return opts.Placeholder
}

func defaultSensitivePatterns() []string {
	return []string{
		"password", "passwd", "secret", "token",
		"api_key", "apikey", "private_key", "auth",
		"credential", "cert", "private",
	}
}
