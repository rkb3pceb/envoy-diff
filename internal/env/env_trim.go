package env

import (
	"strings"
)

// TrimOptions controls how values are trimmed.
type TrimOptions struct {
	// TrimKeys removes leading/trailing whitespace from keys.
	TrimKeys bool
	// TrimValues removes leading/trailing whitespace from values.
	TrimValues bool
	// TrimPrefix removes a specific prefix from all values.
	TrimPrefix string
	// TrimSuffix removes a specific suffix from all values.
	TrimSuffix string
	// TrimChars removes the given characters from both ends of values.
	TrimChars string
}

// DefaultTrimOptions returns TrimOptions with sensible defaults.
func DefaultTrimOptions() TrimOptions {
	return TrimOptions{
		TrimKeys:   false,
		TrimValues: true,
	}
}

// TrimResult holds the outcome of a trim operation.
type TrimResult struct {
	Output   map[string]string
	Modified []string
}

// TrimMap applies trimming rules to the given env map.
// It returns a new map and a result describing what changed.
func TrimMap(input map[string]string, opts TrimOptions) TrimResult {
	out := make(map[string]string, len(input))
	modified := []string{}

	for k, v := range input {
		newKey := k
		if opts.TrimKeys {
			newKey = strings.TrimSpace(k)
		}

		newVal := v
		if opts.TrimValues {
			newVal = strings.TrimSpace(newVal)
		}
		if opts.TrimPrefix != "" {
			newVal = strings.TrimPrefix(newVal, opts.TrimPrefix)
		}
		if opts.TrimSuffix != "" {
			newVal = strings.TrimSuffix(newVal, opts.TrimSuffix)
		}
		if opts.TrimChars != "" {
			newVal = strings.Trim(newVal, opts.TrimChars)
		}

		if newKey != k || newVal != v {
			modified = append(modified, k)
		}
		out[newKey] = newVal
	}

	return TrimResult{Output: out, Modified: modified}
}

// HasTrimChanges returns true if any entries were modified.
func HasTrimChanges(r TrimResult) bool {
	return len(r.Modified) > 0
}
