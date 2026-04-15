package env

import (
	"sort"
	"strings"
)

// GroupOptions controls how keys are grouped.
type GroupOptions struct {
	// Delimiter separates the prefix from the rest of the key (default "_").
	Delimiter string
	// MaxDepth is the number of delimiter segments used to form the group key.
	// 0 means use only the first segment.
	MaxDepth int
	// IncludeUngrouped places keys with no delimiter into a special group.
	IncludeUngrouped bool
	// UngroupedLabel is the name of the catch-all group (default "OTHER").
	UngroupedLabel string
}

// DefaultGroupOptions returns sensible defaults.
func DefaultGroupOptions() GroupOptions {
	return GroupOptions{
		Delimiter:        "_",
		MaxDepth:         0,
		IncludeUngrouped: true,
		UngroupedLabel:   "OTHER",
	}
}

// GroupResult holds a single group produced by GroupMap.
type GroupResult struct {
	Label string
	Keys  []string
	Vars  map[string]string
}

// GroupMap partitions env vars by a common key prefix.
// Groups are returned sorted by label.
func GroupMap(m map[string]string, opts GroupOptions) []GroupResult {
	if opts.Delimiter == "" {
		opts.Delimiter = "_"
	}
	if opts.UngroupedLabel == "" {
		opts.UngroupedLabel = "OTHER"
	}

	buckets := map[string]map[string]string{}

	for k, v := range m {
		label := groupLabel(k, opts.Delimiter, opts.MaxDepth)
		if label == "" {
			if !opts.IncludeUngrouped {
				continue
			}
			label = opts.UngroupedLabel
		}
		if buckets[label] == nil {
			buckets[label] = map[string]string{}
		}
		buckets[label][k] = v
	}

	results := make([]GroupResult, 0, len(buckets))
	for label, vars := range buckets {
		keys := make([]string, 0, len(vars))
		for k := range vars {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		results = append(results, GroupResult{Label: label, Keys: keys, Vars: vars})
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Label < results[j].Label })
	return results
}

// groupLabel derives the group label for a key.
func groupLabel(key, delim string, depth int) string {
	parts := strings.Split(key, delim)
	if len(parts) < 2 {
		return ""
	}
	end := depth + 1
	if end <= 0 || end > len(parts)-1 {
		end = 1
	}
	return strings.Join(parts[:end], delim)
}
