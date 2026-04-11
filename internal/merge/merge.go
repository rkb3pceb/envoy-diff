// Package merge provides utilities for merging multiple env maps
// with configurable override strategies.
package merge

import "fmt"

// Strategy defines how conflicts are resolved when merging env maps.
type Strategy int

const (
	// StrategyLast means the last source wins on conflict.
	StrategyLast Strategy = iota
	// StrategyFirst means the first source wins on conflict.
	StrategyFirst
	// StrategyError means a conflict returns an error.
	StrategyError
)

// Result holds the merged map and metadata about the merge.
type Result struct {
	Env      map[string]string
	Overrides map[string][]string // key -> list of source indices that defined it
	Conflicts []string            // keys with conflicts (StrategyError only)
}

// Merge combines multiple env maps according to the given strategy.
// Sources are applied in order: index 0 is the base.
func Merge(strategy Strategy, sources ...map[string]string) (*Result, error) {
	result := &Result{
		Env:      make(map[string]string),
		Overrides: make(map[string][]string),
	}

	for i, src := range sources {
		for k, v := range src {
			label := fmt.Sprintf("source[%d]", i)
			if _, exists := result.Env[k]; exists {
				result.Overrides[k] = append(result.Overrides[k], label)
				switch strategy {
				case StrategyFirst:
					continue
				case StrategyError:
					result.Conflicts = append(result.Conflicts, k)
					continue
				default: // StrategyLast
					result.Env[k] = v
				}
			} else {
				result.Env[k] = v
				result.Overrides[k] = []string{label}
			}
		}
	}

	if strategy == StrategyError && len(result.Conflicts) > 0 {
		return result, fmt.Errorf("merge conflict on keys: %v", result.Conflicts)
	}

	return result, nil
}
