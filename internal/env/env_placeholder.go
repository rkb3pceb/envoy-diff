package env

import (
	"fmt"
	"strings"
)

// PlaceholderOptions controls how placeholder detection and substitution behaves.
type PlaceholderOptions struct {
	// Markers lists strings that indicate an unresolved placeholder (e.g. "CHANGEME", "TODO").
	Markers []string
	// Substitutions maps placeholder marker values to replacement values.
	Substitutions map[string]string
	// ErrorOnUnresolved causes SubstitutePlaceholders to return an error if any
	// placeholder remains after substitution.
	ErrorOnUnresolved bool
}

// DefaultPlaceholderOptions returns sensible defaults.
func DefaultPlaceholderOptions() PlaceholderOptions {
	return PlaceholderOptions{
		Markers: []string{"CHANGEME", "TODO", "FIXME", "PLACEHOLDER", "<REPLACE>"},
		Substitutions: map[string]string{},
	}
}

// PlaceholderResult holds information about a single placeholder finding.
type PlaceholderResult struct {
	Key      string
	Value    string
	Marker   string
	Resolved bool
}

// FindPlaceholders returns all entries whose values match a known placeholder marker.
func FindPlaceholders(env map[string]string, opts PlaceholderOptions) []PlaceholderResult {
	var results []PlaceholderResult
	for _, k := range SortedKeys(env, DefaultSortOptions()) {
		v := env[k]
		for _, marker := range opts.Markers {
			if strings.EqualFold(v, marker) || strings.Contains(strings.ToUpper(v), strings.ToUpper(marker)) {
				results = append(results, PlaceholderResult{
					Key:    k,
					Value:  v,
					Marker: marker,
				})
				break
			}
		}
	}
	return results
}

// SubstitutePlaceholders replaces placeholder values using opts.Substitutions.
// It returns a new map and a list of results describing each placeholder found.
// If ErrorOnUnresolved is true and any placeholder remains, an error is returned.
func SubstitutePlaceholders(env map[string]string, opts PlaceholderOptions) (map[string]string, []PlaceholderResult, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	findings := FindPlaceholders(out, opts)
	for i, f := range findings {
		if replacement, ok := opts.Substitutions[f.Key]; ok {
			out[f.Key] = replacement
			findings[i].Resolved = true
		}
	}

	if opts.ErrorOnUnresolved {
		var unresolved []string
		for _, f := range findings {
			if !f.Resolved {
				unresolved = append(unresolved, f.Key)
			}
		}
		if len(unresolved) > 0 {
			return out, findings, fmt.Errorf("unresolved placeholders: %s", strings.Join(unresolved, ", "))
		}
	}

	return out, findings, nil
}

// HasPlaceholders returns true if any placeholder markers are found in env.
func HasPlaceholders(env map[string]string, opts PlaceholderOptions) bool {
	return len(FindPlaceholders(env, opts)) > 0
}
