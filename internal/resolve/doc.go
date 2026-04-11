// Package resolve implements environment variable reference expansion
// for envoy-diff.
//
// It supports the standard ${VAR} and $VAR syntax used in shell and
// Docker Compose files. Resolution is performed against a caller-
// supplied map, with an optional fallback to the host OS environment.
//
// Example usage:
//
//	env := map[string]string{
//		"HOST":    "db.internal",
//		"DB_URL":  "postgres://${HOST}/app",
//	}
//	resolved := resolve.Map(env, resolve.DefaultOptions())
//	// resolved["DB_URL"] == "postgres://db.internal/app"
//
// Unresolved references (where the referenced key is absent) are left
// in place so that callers can detect and report them via UnresolvedKeys.
//
// Interpolation is non-recursive: if a substituted value itself contains
// a ${VAR} reference, that reference is not expanded. This mirrors the
// behaviour of Docker Compose and avoids infinite-loop edge cases.
//
// # Options
//
// DefaultOptions returns an Options value with sensible defaults:
//   - FallbackToOS: false (do not read from os.Getenv)
//   - KeepUnresolved: true (leave unmatched references intact)
//
// Pass a customised Options to Map or Value to override these defaults.
package resolve
