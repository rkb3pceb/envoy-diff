// Package history provides a persistent log of envoy-diff comparison runs.
//
// Each time a diff is performed, a summary Entry can be recorded to a local
// JSON store on disk. This enables users to audit environment drift over time,
// review past comparison results, and detect patterns in configuration changes.
//
// Typical usage:
//
//	store, err := history.Open("~/.envoy-diff/history.json")
//	if err != nil { ... }
//
//	err = store.Append(history.Entry{
//	    ID:       "run-001",
//	    OldFile:  "staging.env",
//	    NewFile:  "production.env",
//	    Added:    3,
//	    Removed:  1,
//	    Modified: 5,
//	    Findings: 2,
//	})
package history
