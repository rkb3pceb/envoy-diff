package env

import "strings"

// DedupeOptions controls how duplicate detection works.
type DedupeOptions struct {
	// CaseInsensitive treats keys as equal regardless of case.
	CaseInsensitive bool
	// KeepFirst retains the first occurrence instead of the last.
	KeepFirst bool
	// ReportOnly returns the duplicates found without modifying the map.
	ReportOnly bool
}

// DefaultDedupeOptions returns sensible defaults for deduplication.
func DefaultDedupeOptions() DedupeOptions {
	return DedupeOptions{
		CaseInsensitive: false,
		KeepFirst:       false,
		ReportOnly:      false,
	}
}

// DuplicateEntry records a key that appeared more than once across sources.
type DuplicateEntry struct {
	Key    string
	Values []string
}

// DedupeResult holds the deduplicated map and any duplicate entries found.
type DedupeResult struct {
	Map        map[string]string
	Duplicates []DuplicateEntry
}

// HasDuplicates returns true when at least one duplicate key was detected.
func (r DedupeResult) HasDuplicates() bool {
	return len(r.Duplicates) > 0
}

// DedupeMap removes or reports duplicate keys across the provided ordered
// list of env maps. Later maps win by default (KeepFirst=false).
func DedupeMap(sources []map[string]string, opts DedupeOptions) DedupeResult {
	// Track all values seen per normalised key.
	type entry struct {
		originalKey string
		values      []string
	}
	seen := make(map[string]*entry)
	order := []string{} // insertion order of normalised keys

	for _, src := range sources {
		for k, v := range src {
			norm := k
			if opts.CaseInsensitive {
				norm = strings.ToLower(k)
			}
			if e, exists := seen[norm]; exists {
				e.values = append(e.values, v)
				if !opts.KeepFirst {
					e.originalKey = k
				}
			} else {
				seen[norm] = &entry{originalKey: k, values: []string{v}}
				order = append(order, norm)
			}
		}
	}

	out := make(map[string]string, len(seen))
	var dupes []DuplicateEntry

	for _, norm := range order {
		e := seen[norm]
		if len(e.values) > 1 {
			dupes = append(dupes, DuplicateEntry{Key: e.originalKey, Values: e.values})
		}
		if !opts.ReportOnly {
			out[e.originalKey] = e.values[len(e.values)-1]
			if opts.KeepFirst {
				out[e.originalKey] = e.values[0]
			}
		}
	}

	return DedupeResult{Map: out, Duplicates: dupes}
}
