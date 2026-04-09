// Package baseline provides a persistent store for named environment variable
// baselines used by envoy-diff.
//
// A baseline captures the state of an environment's variables at a specific
// point in time, identified by a human-readable name (e.g. "prod-2024-01").
// Baselines can be used as the "old" side of a diff to compare against a
// current deployment config without needing to keep the original file around.
//
// Usage:
//
//	store, err := baseline.Open("/path/to/baselines.json")
//	if err != nil { ... }
//
//	// Save current vars as a named baseline
//	err = store.Set("prod-before-deploy", "production", currentVars)
//
//	// Retrieve later for diffing
//	bl := store.Get("prod-before-deploy")
//	if bl != nil {
//		changes := diff.Compare(bl.Vars, newVars)
//	}
package baseline
