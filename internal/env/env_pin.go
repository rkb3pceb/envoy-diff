package env

import (
	"fmt"
	"sort"
)

// PinOptions controls how key pinning is applied.
type PinOptions struct {
	// PinnedKeys is the set of keys whose values must not change.
	PinnedKeys []string
	// ErrorOnViolation returns an error if any pinned key differs.
	ErrorOnViolation bool
}

// DefaultPinOptions returns sensible defaults.
func DefaultPinOptions() PinOptions {
	return PinOptions{
		ErrorOnViolation: false,
	}
}

// PinViolation describes a key whose value diverged from its pinned value.
type PinViolation struct {
	Key      string
	Pinned   string
	Actual   string
}

// PinResult holds the outcome of a pin check.
type PinResult struct {
	Violations []PinViolation
}

// HasViolations reports whether any pinned key was changed.
func (r PinResult) HasViolations() bool {
	return len(r.Violations) > 0
}

// CheckPins compares env against a set of pinned key=value pairs.
// pinned is a map of key -> expected value. env is the map being validated.
func CheckPins(pinned, env map[string]string, opts PinOptions) (PinResult, error) {
	var result PinResult

	keys := make([]string, 0, len(pinned))
	for k := range pinned {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		expected := pinned[k]
		actual, exists := env[k]
		if !exists || actual != expected {
			result.Violations = append(result.Violations, PinViolation{
				Key:    k,
				Pinned: expected,
				Actual: actual,
			})
		}
	}

	if opts.ErrorOnViolation && result.HasViolations() {
		return result, fmt.Errorf("pin check failed: %d key(s) violated", len(result.Violations))
	}
	return result, nil
}
