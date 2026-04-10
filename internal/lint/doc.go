// Package lint implements static analysis rules for environment variable maps.
//
// It checks for common misconfigurations and style violations:
//
//   - Keys that are not UPPER_SNAKE_CASE
//   - Keys containing whitespace
//   - Empty values
//   - Values that appear to be unfilled placeholders (e.g. "CHANGEME", "<TOKEN>")
//
// Usage:
//
//	env, _ := parser.ParseEnvFile("production.env")
//	findings := lint.Lint(env)
//	for _, f := range findings {
//		fmt.Println(f)
//	}
//
// Findings carry a Severity of either SeverityWarn or SeverityError,
// allowing callers to decide whether to block a deployment or merely
// surface an advisory notice.
package lint
