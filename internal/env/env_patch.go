package env

import "fmt"

// PatchOptions controls how patch operations are applied.
type PatchOptions struct {
	// ErrorOnMissing causes Patch to return an error if a key to delete or
	// update does not exist in the base map.
	ErrorOnMissing bool

	// AllowOverwrite controls whether set operations may overwrite existing
	// keys. Defaults to true.
	AllowOverwrite bool
}

// DefaultPatchOptions returns sensible defaults for PatchOptions.
func DefaultPatchOptions() PatchOptions {
	return PatchOptions{
		ErrorOnMissing: false,
		AllowOverwrite: true,
	}
}

// PatchOp represents a single patch operation.
type PatchOp struct {
	// Op is one of "set", "delete", or "rename".
	Op  string
	Key string
	// Value is used by the "set" op.
	Value string
	// To is used by the "rename" op.
	To string
}

// PatchResult holds the outcome of a PatchMap call.
type PatchResult struct {
	Env     map[string]string
	Applied []PatchOp
	Skipped []PatchOp
}

// PatchMap applies a sequence of PatchOps to a copy of base.
// Operations are applied in order; later ops see the result of earlier ones.
func PatchMap(base map[string]string, ops []PatchOp, opts PatchOptions) (PatchResult, error) {
	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}

	result := PatchResult{Env: out}

	for _, op := range ops {
		switch op.Op {
		case "set":
			if _, exists := out[op.Key]; exists && !opts.AllowOverwrite {
				result.Skipped = append(result.Skipped, op)
				continue
			}
			out[op.Key] = op.Value
			result.Applied = append(result.Applied, op)

		case "delete":
			if _, exists := out[op.Key]; !exists {
				if opts.ErrorOnMissing {
					return PatchResult{}, fmt.Errorf("patch delete: key %q not found", op.Key)
				}
				result.Skipped = append(result.Skipped, op)
				continue
			}
			delete(out, op.Key)
			result.Applied = append(result.Applied, op)

		case "rename":
			val, exists := out[op.Key]
			if !exists {
				if opts.ErrorOnMissing {
					return PatchResult{}, fmt.Errorf("patch rename: key %q not found", op.Key)
				}
				result.Skipped = append(result.Skipped, op)
				continue
			}
			delete(out, op.Key)
			out[op.To] = val
			result.Applied = append(result.Applied, op)

		default:
			return PatchResult{}, fmt.Errorf("patch: unknown op %q", op.Op)
		}
	}

	return result, nil
}

// HasPatchApplied returns true when at least one op was applied.
func HasPatchApplied(r PatchResult) bool {
	return len(r.Applied) > 0
}

// Summary returns a human-readable description of the patch result,
// listing how many operations were applied and skipped.
func (r PatchResult) Summary() string {
	return fmt.Sprintf("patch result: %d applied, %d skipped", len(r.Applied), len(r.Skipped))
}
