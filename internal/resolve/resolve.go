// Package resolve provides environment variable resolution,
// expanding references like ${VAR} or $VAR within env values.
package resolve

import (
	"os"
	"strings"
)

// Options controls resolution behaviour.
type Options struct {
	// FallbackToOS allows falling back to the host OS environment
	// when a variable is not found in the provided map.
	FallbackToOS bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{FallbackToOS: false}
}

// Map resolves all values in the provided env map, expanding
// ${VAR} and $VAR references using other keys in the same map.
// Unresolved references are left as-is unless FallbackToOS is set.
func Map(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = Expand(v, env, opts)
	}
	return out
}

// Expand resolves variable references within a single value string.
func Expand(value string, env map[string]string, opts Options) string {
	return os.Expand(value, func(key string) string {
		if val, ok := env[key]; ok {
			return val
		}
		if opts.FallbackToOS {
			return os.Getenv(key)
		}
		// Return the original reference so callers can detect it.
		return "${" + key + "}"
	})
}

// HasUnresolved reports whether the value still contains unresolved
// variable references after expansion.
func HasUnresolved(value string) bool {
	return strings.Contains(value, "${") || strings.Contains(value, "$")
}

// UnresolvedKeys returns the keys whose values contain unresolved
// variable references.
func UnresolvedKeys(env map[string]string) []string {
	var keys []string
	for k, v := range env {
		if HasUnresolved(v) {
			keys = append(keys, k)
		}
	}
	return keys
}
