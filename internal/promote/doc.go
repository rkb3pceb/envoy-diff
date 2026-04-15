// Package promote implements stage-to-stage environment promotion logic
// for envoy-diff.
//
// A promotion compares the environment variable set of a source stage
// (e.g. "staging") against a target stage (e.g. "production"), producing
// a structured Result that includes:
//
//   - The list of diff.Change entries between the two stages.
//   - Audit findings (sensitivity, risk level) via the audit package.
//   - Policy violations and a Blocked flag when a policy file is loaded.
//
// Sensitive values are automatically redacted using a redact.Context so
// that secrets are never surfaced in promotion reports.
//
// # Policy files
//
// A policy file (YAML or JSON) can be supplied via Options.PolicyPath to
// enforce rules such as "no high-risk changes to production" or "require
// approval for keys matching SECRET_*". When any rule is violated,
// Result.Blocked is set to true and Result.Violations lists each breach.
//
// Basic usage:
//
//	from := promote.Stage{Name: "staging",    Env: stagingVars}
//	to   := promote.Stage{Name: "production", Env: prodVars}
//
//	result, err := promote.Evaluate(from, to, promote.DefaultOptions())
//	if err != nil { ... }
//	if result.Blocked { ... }
package promote
