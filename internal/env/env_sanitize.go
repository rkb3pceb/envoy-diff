package env

import (
	"strings"
)

// SanitizeOptions controls how env map sanitization is performed.
type SanitizeOptions struct {
	// TrimSpace removes leading/trailing whitespace from values.
	TrimSpace bool
	// NormalizeKeys uppercases all keys.
	NormalizeKeys bool
	// RemoveEmpty drops keys with empty values after trimming.
	RemoveEmpty bool
	// ReplaceInvalidChars replaces characters in keys that are not
	// alphanumeric or underscore with an underscore.
	ReplaceInvalidChars bool
}

// DefaultSanitizeOptions returns sensible defaults for sanitization.
func DefaultSanitizeOptions() SanitizeOptions {
	return SanitizeOptions{
		TrimSpace:           true,
		NormalizeKeys:       false,
		RemoveEmpty:         false,
		ReplaceInvalidChars: false,
	}
}

// SanitizeResult holds the output of a sanitize operation.
type SanitizeResult struct {
	Map      map[string]string
	// Renamed tracks keys that were renamed due to normalization or invalid chars.
	Renamed  map[string]string // old -> new
	// Dropped lists keys removed because their value was empty.
	Dropped  []string
}

// SanitizeMap applies sanitization rules to the provided env map.
// It does not mutate the input map.
func SanitizeMap(input map[string]string, opts SanitizeOptions) SanitizeResult {
	out := make(map[string]string, len(input))
	renamed := make(map[string]string)
	var dropped []string

	for k, v := range input {
		if opts.TrimSpace {
			v = strings.TrimSpace(v)
		}
		if opts.RemoveEmpty && v == "" {
			dropped = append(dropped, k)
			continue
		}

		newKey := k
		if opts.ReplaceInvalidChars {
			newKey = sanitizeKey(newKey)
		}
		if opts.NormalizeKeys {
			newKey = strings.ToUpper(newKey)
		}
		if newKey != k {
			renamed[k] = newKey
		}
		out[newKey] = v
	}

	return SanitizeResult{
		Map:     out,
		Renamed: renamed,
		Dropped: dropped,
	}
}

// HasSanitizeChanges returns true if any keys were renamed or dropped.
func HasSanitizeChanges(r SanitizeResult) bool {
	return len(r.Renamed) > 0 || len(r.Dropped) > 0
}

func sanitizeKey(key string) string {
	var b strings.Builder
	for _, ch := range key {
		if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') ||
			(ch >= '0' && ch <= '9') || ch == '_' {
			b.WriteRune(ch)
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}
