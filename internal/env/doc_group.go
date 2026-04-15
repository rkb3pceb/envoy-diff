// Package env provides utilities for working with environment variable maps.
//
// # Group
//
// GroupMap partitions a flat env map into labelled groups based on a shared
// key prefix. By default the first segment before "_" is used as the group
// label, but the delimiter and depth are configurable.
//
// Example:
//
//	groups := env.GroupMap(m, env.DefaultGroupOptions())
//	for _, g := range groups {
//	    fmt.Printf("[%s]\n", g.Label)
//	    for _, k := range g.Keys {
//	        fmt.Printf("  %s=%s\n", k, g.Vars[k])
//	    }
//	}
//
// Keys that contain no delimiter are placed into an "OTHER" catch-all group
// unless IncludeUngrouped is set to false.
package env
