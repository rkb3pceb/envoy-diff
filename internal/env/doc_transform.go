// Package env provides utilities for loading, transforming, and validating
// environment variable maps from one or more .env files.
//
// # Transform
//
// TransformMap applies a pipeline of key/value mutations to a map without
// mutating the original. Supported built-in operations:
//
//   - PrefixAdd    – prepend a string to every key
//   - PrefixStrip  – remove a prefix from keys that carry it
//   - UppercaseKeys / LowercaseKeys – normalise key casing (mutually exclusive)
//   - DropEmpty    – discard entries whose value is the empty string
//   - Extra        – ordered slice of user-supplied TransformFn callbacks
//
// Example:
//
//	opts := env.DefaultTransformOptions()
//	opts.PrefixStrip = "APP_"
//	opts.UppercaseKeys = true
//	out, err := env.TransformMap(src, opts)
//
// TransformFn signature:
//
//	type TransformFn func(key, value string) (newKey, newValue string, keep bool)
//
// Returning keep=false drops the entry from the output map.
package env
