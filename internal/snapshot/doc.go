// Package snapshot provides utilities for capturing, saving, and loading
// named snapshots of environment variable sets.
//
// A Snapshot records the name, source file, timestamp, and full map of
// key-value pairs for a given environment state. Snapshots can be persisted
// to disk as JSON files and reloaded later, enabling point-in-time comparisons
// across deployments or environments.
//
// Typical usage:
//
//	vars, err := parser.ParseEnvFile("production.env")
//	if err != nil { ... }
//
//	snap := snapshot.New("prod-2024-01-15", "production.env", vars)
//	if err := snapshot.Save(snap, "snapshots/prod-2024-01-15.json"); err != nil { ... }
//
//	// Later, load and compare:
//	old, err := snapshot.Load("snapshots/prod-2024-01-14.json")
//	new, err := snapshot.Load("snapshots/prod-2024-01-15.json")
//	changes := diff.Compare(old.Vars, new.Vars)
package snapshot
