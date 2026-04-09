// Package redact provides helpers for identifying and masking sensitive
// environment variable values before they are surfaced in diffs, reports,
// or any other human-readable output.
//
// A key is considered sensitive when its name contains one of a curated
// set of substrings such as PASSWORD, TOKEN, SECRET, or API_KEY. The
// comparison is case-insensitive so that both DB_PASSWORD and db_password
// are treated equally.
//
// Usage:
//
//	if redact.IsSensitive(key) {
//		// handle accordingly
//	}
//
//	// Redact a single value
//	displayValue := redact.Value(key, rawValue)
//
//	// Redact an entire map before printing
//	safeEnv := redact.Map(envMap)
package redact
