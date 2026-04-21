// Package env provides utilities for loading, transforming, and inspecting
// environment variable maps.
//
// # Key Rotation
//
// RotateMap renames keys in an env map according to a set of OLD→NEW rules.
// It is useful when migrating services to a new variable naming convention
// without losing existing values.
//
// Basic usage:
//
//	opts := env.DefaultRotateOptions()
//	opts.Rules = map[string]string{
//		"DB_PASS": "DATABASE_PASSWORD",
//	}
//	out, result, err := env.RotateMap(src, opts)
//
// Options:
//   - KeepOld: when true the original key is not deleted, leaving both names
//     present in the output (useful for zero-downtime rollouts).
//   - ErrorOnMissing: return an error when a source key named in a rule does
//     not exist in the input map.
//   - ErrorOnConflict: return an error when the destination key already exists
//     in the input map.
//
// The returned RotateResult lists which keys were rotated, skipped, or
// caused a conflict, allowing callers to report outcomes to the user.
package env
