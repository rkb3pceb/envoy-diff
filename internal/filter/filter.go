// Package filter provides utilities for filtering diff results
// based on key patterns, prefixes, or change types.
package filter

import (
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// Options holds the configuration for filtering diff results.
type Options struct {
	// Prefix filters entries whose keys start with the given prefix.
	Prefix string

	// KeyContains filters entries whose keys contain the given substring.
	KeyContains string

	// OnlyChanged, when true, excludes Unchanged entries from results.
	OnlyChanged bool

	// Types is an optional allowlist of ChangeTypes to include.
	// If empty, all types are included.
	Types []diff.ChangeType
}

// Apply filters a slice of diff.Entry according to the provided Options
// and returns a new slice containing only the matching entries.
func Apply(entries []diff.Entry, opts Options) []diff.Entry {
	var result []diff.Entry

	allowedTypes := buildTypeSet(opts.Types)

	for _, e := range entries {
		if opts.OnlyChanged && e.Type == diff.Unchanged {
			continue
		}

		if len(allowedTypes) > 0 {
			if _, ok := allowedTypes[e.Type]; !ok {
				continue
			}
		}

		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}

		if opts.KeyContains != "" && !strings.Contains(e.Key, opts.KeyContains) {
			continue
		}

		result = append(result, e)
	}

	return result
}

// buildTypeSet converts a slice of ChangeType into a set (map) for O(1) lookup.
func buildTypeSet(types []diff.ChangeType) map[diff.ChangeType]struct{} {
	if len(types) == 0 {
		return nil
	}
	s := make(map[diff.ChangeType]struct{}, len(types))
	for _, t := range types {
		s[t] = struct{}{}
	}
	return s
}
