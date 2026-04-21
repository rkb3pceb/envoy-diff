package env

import (
	"strings"
)

// UppercaseOptions controls the behaviour of UppercaseMap.
type UppercaseOptions struct {
	// Keys uppercases all map keys when true.
	Keys bool

	// Values uppercases all map values when true.
	Values bool

	// OnlyKeys restricts key uppercasing to the listed keys.
	// When empty every key is considered.
	OnlyKeys []string
}

// DefaultUppercaseOptions returns a sensible default: uppercase keys only.
func DefaultUppercaseOptions() UppercaseOptions {
	return UppercaseOptions{Keys: true}
}

// UppercaseResult holds the output of UppercaseMap.
type UppercaseResult struct {
	Map      map[string]string
	Changed  []string // keys whose name or value was altered
}

// UppercaseMap applies uppercasing to keys and/or values of src according to
// opts. The original map is never mutated.
func UppercaseMap(src map[string]string, opts UppercaseOptions) UppercaseResult {
	out := make(map[string]string, len(src))
	var changed []string

	allowSet := make(map[string]struct{}, len(opts.OnlyKeys))
	for _, k := range opts.OnlyKeys {
		allowSet[k] = struct{}{}
	}

	for k, v := range src {
		newKey := k
		newVal := v
		modified := false

		if opts.Keys {
			if len(allowSet) == 0 {
				newKey = strings.ToUpper(k)
			} else if _, ok := allowSet[k]; ok {
				newKey = strings.ToUpper(k)
			}
			if newKey != k {
				modified = true
			}
		}

		if opts.Values {
			if len(allowSet) == 0 {
				newVal = strings.ToUpper(v)
			} else if _, ok := allowSet[k]; ok {
				newVal = strings.ToUpper(v)
			}
			if newVal != v {
				modified = true
			}
		}

		out[newKey] = newVal
		if modified {
			changed = append(changed, k)
		}
	}

	return UppercaseResult{Map: out, Changed: changed}
}

// HasUppercaseChanges returns true when at least one key or value was altered.
func HasUppercaseChanges(r UppercaseResult) bool {
	return len(r.Changed) > 0
}
