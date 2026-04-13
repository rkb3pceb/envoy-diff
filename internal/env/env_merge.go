package env

import "fmt"

// MergeStrategy controls how conflicts between sources are resolved.
type MergeStrategy int

const (
	// MergeLastWins overwrites existing keys with values from later sources.
	MergeLastWins MergeStrategy = iota
	// MergeFirstWins keeps the first value seen for a key.
	MergeFirstWins
	// MergeErrorOnConflict returns an error if the same key appears in multiple sources.
	MergeErrorOnConflict
)

// MergeOptions configures the behaviour of MergeMaps.
type MergeOptions struct {
	Strategy  MergeStrategy
	// Overrides is an optional set of key=value pairs applied last, always winning.
	Overrides map[string]string
}

// DefaultMergeOptions returns sensible defaults (last-wins, no overrides).
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{Strategy: MergeLastWins}
}

// MergeResult holds the merged map and metadata about the operation.
type MergeResult struct {
	Merged    map[string]string
	// Conflicts lists keys that appeared in more than one source (all strategies).
	Conflicts []string
}

// MergeMaps merges multiple env maps according to opts.
// Sources are processed in order; opts.Overrides are applied last.
func MergeMaps(sources []map[string]string, opts MergeOptions) (MergeResult, error) {
	merged := make(map[string]string)
	conflicts := []string{}
	seen := make(map[string]bool)

	for _, src := range sources {
		for k, v := range src {
			if seen[k] {
				conflicts = append(conflicts, k)
				switch opts.Strategy {
				case MergeErrorOnConflict:
					return MergeResult{}, fmt.Errorf("merge conflict: key %q appears in multiple sources", k)
				case MergeFirstWins:
					continue
				default: // MergeLastWins
					merged[k] = v
				}
			} else {
				merged[k] = v
				seen[k] = true
			}
		}
	}

	for k, v := range opts.Overrides {
		merged[k] = v
	}

	return MergeResult{Merged: merged, Conflicts: dedupe(conflicts)}, nil
}

// HasMergeConflicts returns true when the result contains at least one conflict.
func HasMergeConflicts(r MergeResult) bool {
	return len(r.Conflicts) > 0
}
