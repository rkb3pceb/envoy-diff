package env

import (
	"sort"

	"github.com/wvdschel-personal/envoy-diff/internal/diff"
)

// DiffOptions controls behaviour of DiffMaps.
type DiffOptions struct {
	// IncludeUnchanged includes entries with no change in the result.
	IncludeUnchanged bool
}

// DefaultDiffOptions returns sensible defaults for DiffMaps.
func DefaultDiffOptions() DiffOptions {
	return DiffOptions{
		IncludeUnchanged: false,
	}
}

// DiffResult holds the outcome of comparing two env maps.
type DiffResult struct {
	Changes []diff.Change
	Added   int
	Removed int
	Modified int
	Unchanged int
}

// DiffMaps compares oldEnv and newEnv and returns a DiffResult.
// Changes are sorted alphabetically by key.
func DiffMaps(oldEnv, newEnv map[string]string, opts DiffOptions) DiffResult {
	seen := make(map[string]struct{})
	var changes []diff.Change

	for key, oldVal := range oldEnv {
		seen[key] = struct{}{}
		newVal, exists := newEnv[key]
		switch {
		case !exists:
			changes = append(changes, diff.Change{
				Key:      key,
				OldValue: oldVal,
				NewValue: "",
				Type:     diff.Removed,
			})
		case newVal != oldVal:
			changes = append(changes, diff.Change{
				Key:      key,
				OldValue: oldVal,
				NewValue: newVal,
				Type:     diff.Modified,
			})
		default:
			if opts.IncludeUnchanged {
				changes = append(changes, diff.Change{
					Key:      key,
					OldValue: oldVal,
					NewValue: newVal,
					Type:     diff.Unchanged,
				})
			}
		}
	}

	for key, newVal := range newEnv {
		if _, ok := seen[key]; ok {
			continue
		}
		changes = append(changes, diff.Change{
			Key:      key,
			OldValue: "",
			NewValue: newVal,
			Type:     diff.Added,
		})
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	result := DiffResult{Changes: changes}
	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			result.Added++
		case diff.Removed:
			result.Removed++
		case diff.Modified:
			result.Modified++
		case diff.Unchanged:
			result.Unchanged++
		}
	}
	return result
}
