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
package resolve
