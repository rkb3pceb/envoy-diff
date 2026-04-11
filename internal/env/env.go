// Package env provides utilities for loading, merging, and resolving
// environment variable maps from multiple sources in a defined priority order.
package env

import (
	"fmt"
	"sort"

	"github.com/your-org/envoy-diff/internal/parser"
)

// Source represents a named environment variable source.
type Source struct {
	Name string
	Vars map[string]string
}

// LoadSources parses each file path into a named Source.
// Files are returned in the same order as paths.
func LoadSources(paths []string) ([]Source, error) {
	sources := make([]Source, 0, len(paths))
	for _, p := range paths {
		vars, err := parser.ParseEnvFile(p)
		if err != nil {
			return nil, fmt.Errorf("env: loading %q: %w", p, err)
		}
		sources = append(sources, Source{Name: p, Vars: vars})
	}
	return sources, nil
}

// Flatten merges all sources into a single map. Later sources override
// earlier ones (last-write-wins). Returns the merged map and a list of
// keys that were overridden at least once.
func Flatten(sources []Source) (map[string]string, []string) {
	result := make(map[string]string)
	overrideSet := make(map[string]bool)

	for _, src := range sources {
		for k, v := range src.Vars {
			if _, exists := result[k]; exists {
				overrideSet[k] = true
			}
			result[k] = v
		}
	}

	overridden := make([]string, 0, len(overrideSet))
	for k := range overrideSet {
		overridden = append(overridden, k)
	}
	sort.Strings(overridden)
	return result, overridden
}

// Keys returns a sorted slice of all keys present in the map.
func Keys(vars map[string]string) []string {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
