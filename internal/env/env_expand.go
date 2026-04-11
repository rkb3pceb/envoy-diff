package env

import (
	"fmt"
	"os"
	"strings"
)

// ExpandOptions controls how environment variable expansion behaves.
type ExpandOptions struct {
	// FallbackToOS allows falling back to OS environment variables
	// when a key is not found in the provided map.
	FallbackToOS bool
	// ErrorOnMissing causes Expand to return an error if a referenced
	// variable is not found in the map (or OS, if FallbackToOS is true).
	ErrorOnMissing bool
}

// DefaultExpandOptions returns sensible defaults for expansion.
func DefaultExpandOptions() ExpandOptions {
	return ExpandOptions{
		FallbackToOS:   false,
		ErrorOnMissing: false,
	}
}

// ExpandMap expands all values in the given map using other entries in the
// same map as a source for variable references (${VAR} or $VAR syntax).
// It does not mutate the input map.
func ExpandMap(env map[string]string, opts ExpandOptions) (map[string]string, error) {
	out := make(map[string]string, len(env))
	var missing []string

	lookup := func(key string) string {
		if v, ok := env[key]; ok {
			return v
		}
		if opts.FallbackToOS {
			if v, ok := os.LookupEnv(key); ok {
				return v
			}
		}
		missing = append(missing, key)
		return ""
	}

	for k, v := range env {
		out[k] = os.Expand(v, lookup)
	}

	if opts.ErrorOnMissing && len(missing) > 0 {
		return out, fmt.Errorf("unresolved variable references: %s", strings.Join(dedupe(missing), ", "))
	}

	return out, nil
}

// MissingRefs returns all variable references in the map values that cannot
// be resolved from within the map itself (or the OS, if fallbackToOS is true).
func MissingRefs(env map[string]string, fallbackToOS bool) []string {
	var missing []string
	for _, v := range env {
		os.Expand(v, func(key string) string {
			if _, ok := env[key]; ok {
				return env[key]
			}
			if fallbackToOS {
				if val, ok := os.LookupEnv(key); ok {
					return val
				}
			}
			missing = append(missing, key)
			return ""
		})
	}
	return dedupe(missing)
}

func dedupe(ss []string) []string {
	seen := make(map[string]struct{}, len(ss))
	out := ss[:0]
	for _, s := range ss {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}
