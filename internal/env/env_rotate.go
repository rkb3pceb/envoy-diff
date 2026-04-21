package env

import "fmt"

// RotateOptions controls how key rotation is applied.
type RotateOptions struct {
	// Rules maps old key names to new key names.
	Rules map[string]string
	// KeepOld retains the original key alongside the new one.
	KeepOld bool
	// ErrorOnMissing returns an error if a rule references a key not in the map.
	ErrorOnMissing bool
	// ErrorOnConflict returns an error if the new key already exists.
	ErrorOnConflict bool
}

// DefaultRotateOptions returns safe defaults for RotateMap.
func DefaultRotateOptions() RotateOptions {
	return RotateOptions{
		Rules:           map[string]string{},
		KeepOld:         false,
		ErrorOnMissing:  false,
		ErrorOnConflict: false,
	}
}

// RotateResult holds the outcome of a rotation operation.
type RotateResult struct {
	Rotated  []string // old key names that were rotated
	Skipped  []string // old key names skipped (missing or conflict)
	Conflict []string // new key names that already existed
}

// HasRotated returns true when at least one key was rotated.
func (r RotateResult) HasRotated() bool { return len(r.Rotated) > 0 }

// RotateMap renames keys according to the provided rules, producing a new map.
func RotateMap(src map[string]string, opts RotateOptions) (map[string]string, RotateResult, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}

	var result RotateResult

	for oldKey, newKey := range opts.Rules {
		val, exists := out[oldKey]
		if !exists {
			if opts.ErrorOnMissing {
				return nil, result, fmt.Errorf("rotate: key %q not found", oldKey)
			}
			result.Skipped = append(result.Skipped, oldKey)
			continue
		}

		if _, conflict := out[newKey]; conflict {
			if opts.ErrorOnConflict {
				return nil, result, fmt.Errorf("rotate: destination key %q already exists", newKey)
			}
			result.Conflict = append(result.Conflict, newKey)
			result.Skipped = append(result.Skipped, oldKey)
			continue
		}

		out[newKey] = val
		if !opts.KeepOld {
			delete(out, oldKey)
		}
		result.Rotated = append(result.Rotated, oldKey)
	}

	return out, result, nil
}
