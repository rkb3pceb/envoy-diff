// Package env provides utilities for loading, transforming, and inspecting
// environment variable maps.
//
// # ApplyDefaults
//
// ApplyDefaults fills gaps in an env map by applying a caller-supplied set of
// default key/value pairs.  Three modes are supported:
//
//   - Normal (default): a default is written only when the key is entirely
//     absent from the source map.
//
//   - OverwriteEmpty: additionally replaces keys whose value is the empty
//     string ("").
//
//   - OverwriteAll: replaces every key in the defaults map unconditionally,
//     effectively acting as a forced merge from the defaults side.
//
// The source map is never mutated; ApplyDefaults always returns a fresh copy.
//
// Example:
//
//	opts := env.DefaultDefaultOptions()
//	opts.Defaults = map[string]string{"PORT": "8080", "LOG_LEVEL": "info"}
//	opts.OverwriteEmpty = true
//
//	out, result := env.ApplyDefaults(src, opts)
//	fmt.Println(result.Applied) // ["LOG_LEVEL" "PORT"] (sorted)
package env
