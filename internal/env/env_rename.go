package env

import "fmt"

// RenameOptions controls how keys are renamed in an env map.
type RenameOptions struct {
	// Rules maps old key names to new key names.
	Rules map[string]string
	// ErrorOnMissing causes Rename to return an error if a rule references
	// a key that does not exist in the source map.
	ErrorOnMissing bool
	// ErrorOnConflict causes Rename to return an error if the destination
	// key already exists in the map before the rename.
	ErrorOnConflict bool
}

// DefaultRenameOptions returns a RenameOptions with safe defaults.
func DefaultRenameOptions() RenameOptions {
	return RenameOptions{
		Rules:           map[string]string{},
		ErrorOnMissing:  false,
		ErrorOnConflict: false,
	}
}

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	// Map is the resulting env map after renames.
	Map map[string]string
	// Applied lists the old→new pairs that were successfully renamed.
	Applied []RenameApplied
	// Skipped lists old keys that were not found in the source map.
	Skipped []string
}

// RenameApplied records a single successful rename.
type RenameApplied struct {
	OldKey string
	NewKey string
}

// RenameMap applies rename rules to src and returns a new map.
// src is never mutated.
func RenameMap(src map[string]string, opts RenameOptions) (RenameResult, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}

	result := RenameResult{}

	for oldKey, newKey := range opts.Rules {
		val, exists := out[oldKey]
		if !exists {
			if opts.ErrorOnMissing {
				return RenameResult{}, fmt.Errorf("rename: source key %q not found", oldKey)
			}
			result.Skipped = append(result.Skipped, oldKey)
			continue
		}
		if _, conflict := out[newKey]; conflict && oldKey != newKey {
			if opts.ErrorOnConflict {
				return RenameResult{}, fmt.Errorf("rename: destination key %q already exists", newKey)
			}
		}
		delete(out, oldKey)
		out[newKey] = val
		result.Applied = append(result.Applied, RenameApplied{OldKey: oldKey, NewKey: newKey})
	}

	result.Map = out
	return result, nil
}

// HasRenameApplied returns true if at least one rename rule was applied.
func HasRenameApplied(r RenameResult) bool {
	return len(r.Applied) > 0
}
