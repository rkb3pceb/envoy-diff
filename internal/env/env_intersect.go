package env

// IntersectOptions controls how IntersectMaps behaves.
type IntersectOptions struct {
	// KeepValues determines which map's values are used for matching keys.
	// "first" keeps values from the first map, "last" keeps from the last (default).
	KeepValues string

	// RequireEqual only includes keys where all maps share the same value.
	RequireEqual bool
}

// DefaultIntersectOptions returns sensible defaults for IntersectMaps.
func DefaultIntersectOptions() IntersectOptions {
	return IntersectOptions{
		KeepValues: "last",
	}
}

// IntersectResult holds the output of IntersectMaps.
type IntersectResult struct {
	// Map contains keys present in all input maps.
	Map map[string]string

	// Dropped is the count of keys excluded because they were not in all maps.
	Dropped int

	// Conflicts is the count of keys excluded due to RequireEqual mismatch.
	Conflicts int
}

// IntersectMaps returns only keys that appear in every provided map.
// At least two maps must be provided; if fewer are given the first map is
// returned as-is with no keys dropped.
func IntersectMaps(opts IntersectOptions, maps ...map[string]string) IntersectResult {
	result := IntersectResult{Map: make(map[string]string)}

	if len(maps) == 0 {
		return result
	}
	if len(maps) == 1 {
		for k, v := range maps[0] {
			result.Map[k] = v
		}
		return result
	}

	// Seed candidate keys from the first map.
	for key, val := range maps[0] {
		presentInAll := true
		conflict := false

		for _, m := range maps[1:] {
			other, ok := m[key]
			if !ok {
				presentInAll = false
				break
			}
			if opts.RequireEqual && other != val {
				conflict = true
				break
			}
			// Track the "last" value as we iterate.
			if opts.KeepValues == "last" {
				val = other
			}
		}

		switch {
		case !presentInAll:
			result.Dropped++
		case conflict:
			result.Conflicts++
		default:
			result.Map[key] = val
		}
	}

	// Count keys in subsequent maps that were not in the first map.
	seen := maps[0]
	for _, m := range maps[1:] {
		for k := range m {
			if _, inFirst := seen[k]; !inFirst {
				result.Dropped++
			}
		}
	}

	return result
}

// HasIntersectResult returns true when the result map is non-empty.
func HasIntersectResult(r IntersectResult) bool {
	return len(r.Map) > 0
}
