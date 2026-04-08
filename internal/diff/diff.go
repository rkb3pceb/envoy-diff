// Package diff provides functionality for comparing two sets of environment
// variables and producing a structured report of additions, removals, and
// modifications between them.
package diff

// ChangeType represents the kind of change detected for an environment variable.
type ChangeType string

const (
	// Added indicates a key present in the new env but not the old.
	Added ChangeType = "added"
	// Removed indicates a key present in the old env but not the new.
	Removed ChangeType = "removed"
	// Modified indicates a key present in both envs but with a different value.
	Modified ChangeType = "modified"
)

// Change describes a single environment variable change.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Result holds the full diff between two environment variable maps.
type Result struct {
	Changes []Change
}

// HasChanges returns true if any differences were detected.
func (r *Result) HasChanges() bool {
	return len(r.Changes) > 0
}

// Added returns only the changes of type Added.
func (r *Result) Added() []Change {
	return r.filter(Added)
}

// Removed returns only the changes of type Removed.
func (r *Result) Removed() []Change {
	return r.filter(Removed)
}

// Modified returns only the changes of type Modified.
func (r *Result) Modified() []Change {
	return r.filter(Modified)
}

func (r *Result) filter(ct ChangeType) []Change {
	out := make([]Change, 0)
	for _, c := range r.Changes {
		if c.Type == ct {
			out = append(out, c)
		}
	}
	return out
}

// Compare takes two env maps (old and new) and returns a Result describing
// all detected changes. Keys are compared in a deterministic order.
func Compare(oldEnv, newEnv map[string]string) *Result {
	result := &Result{}

	// Detect removed and modified keys.
	for key, oldVal := range oldEnv {
		newVal, exists := newEnv[key]
		if !exists {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Removed,
				OldValue: oldVal,
			})
		} else if oldVal != newVal {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Modified,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}

	// Detect added keys.
	for key, newVal := range newEnv {
		if _, exists := oldEnv[key]; !exists {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Added,
				NewValue: newVal,
			})
		}
	}

	return result
}
