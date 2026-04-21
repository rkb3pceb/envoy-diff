package env

import (
	"fmt"
	"sort"
	"strings"
)

// DefaultJoinOptions returns a JoinOptions with sensible defaults.
func DefaultJoinOptions() JoinOptions {
	return JoinOptions{
		Separator:    ",",
		OutputKey:    "",
		SortValues:   false,
		SkipEmpty:    true,
		Distinct:     false,
	}
}

// JoinOptions controls how JoinMap combines multiple key values into one.
type JoinOptions struct {
	// Keys is the list of source keys whose values will be joined.
	Keys []string

	// OutputKey is the destination key that receives the joined value.
	// If empty, the first key in Keys is used.
	OutputKey string

	// Separator is placed between each value segment. Defaults to ",".
	Separator string

	// SortValues sorts the collected values before joining.
	SortValues bool

	// SkipEmpty omits values that are empty or whitespace-only.
	SkipEmpty bool

	// Distinct removes duplicate values before joining.
	Distinct bool

	// RemoveSources deletes the source keys after joining.
	RemoveSources bool
}

// JoinResult holds the output of a JoinMap call.
type JoinResult struct {
	// Out is the resulting env map.
	Out map[string]string

	// OutputKey is the key that was written.
	OutputKey string

	// Parts are the individual values that were joined.
	Parts []string

	// Skipped lists source keys that were omitted (empty value).
	Skipped []string
}

// JoinMap collects values from the specified keys and joins them into a single
// key using the configured separator. The source map is never mutated.
func JoinMap(src map[string]string, opts JoinOptions) (JoinResult, error) {
	if len(opts.Keys) == 0 {
		return JoinResult{Out: copyMap(src)}, nil
	}

	outKey := opts.OutputKey
	if outKey == "" {
		outKey = opts.Keys[0]
	}

	var parts []string
	var skipped []string
	seen := map[string]bool{}

	for _, k := range opts.Keys {
		v, ok := src[k]
		if !ok {
			continue
		}
		if opts.SkipEmpty && strings.TrimSpace(v) == "" {
			skipped = append(skipped, k)
			continue
		}
		if opts.Distinct {
			if seen[v] {
				continue
			}
			seen[v] = true
		}
		parts = append(parts, v)
	}

	if opts.SortValues {
		sort.Strings(parts)
	}

	out := copyMap(src)
	out[outKey] = strings.Join(parts, opts.Separator)

	if opts.RemoveSources {
		for _, k := range opts.Keys {
			if k != outKey {
				delete(out, k)
			}
		}
	}

	return JoinResult{
		Out:       out,
		OutputKey: outKey,
		Parts:     parts,
		Skipped:   skipped,
	}, nil
}

// HasJoinResult returns true when the result contains a non-empty joined value.
func HasJoinResult(r JoinResult) bool {
	return r.OutputKey != "" && r.Out[r.OutputKey] != ""
}

// copyMap returns a shallow copy of m. It is defined locally to avoid import
// cycles; other env_*.go files in this package do the same.
func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// joinValidateKeys returns an error if any requested key is not present in src.
func joinValidateKeys(src map[string]string, keys []string) error {
	for _, k := range keys {
		if _, ok := src[k]; !ok {
			return fmt.Errorf("join: key %q not found in source map", k)
		}
	}
	return nil
}
