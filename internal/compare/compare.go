// Package compare provides utilities for comparing two env maps
// and producing a structured summary suitable for reporting or export.
package compare

import (
	"sort"

	"github.com/your-org/envoy-diff/internal/diff"
	"github.com/your-org/envoy-diff/internal/redact"
)

// Summary holds aggregate statistics and the full change list
// produced by comparing two environment variable maps.
type Summary struct {
	Added    int
	Removed  int
	Modified int
	Unchanged int
	Changes  []diff.Change
}

// Options controls how the comparison is performed.
type Options struct {
	// RedactSensitive masks sensitive values in the returned changes.
	RedactSensitive bool
}

// Run compares oldEnv and newEnv and returns a Summary.
// Keys are processed in deterministic (sorted) order.
func Run(oldEnv, newEnv map[string]string, opts Options) Summary {
	changes := diff.Compare(oldEnv, newEnv)

	if opts.RedactSensitive {
		for i, c := range changes {
			if redact.IsSensitive(c.Key) {
				changes[i].OldValue = redact.Value(c.Key, c.OldValue)
				changes[i].NewValue = redact.Value(c.Key, c.NewValue)
			}
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	s := Summary{Changes: changes}
	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			s.Added++
		case diff.Removed:
			s.Removed++
		case diff.Modified:
			s.Modified++
		case diff.Unchanged:
			s.Unchanged++
		}
	}
	return s
}

// HasChanges returns true if the summary contains any non-unchanged entries.
func (s Summary) HasChanges() bool {
	return s.Added+s.Removed+s.Modified > 0
}
