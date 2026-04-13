package env

import (
	"sort"
	"strings"
)

// StatsReport holds aggregate statistics about an env map.
type StatsReport struct {
	Total        int
	Empty        int
	Sensitive    int
	Unique       int
	PrefixCounts map[string]int
	TopPrefixes  []PrefixCount
}

// PrefixCount pairs a prefix with its key count.
type PrefixCount struct {
	Prefix string
	Count  int
}

// DefaultStatsOptions returns options with top-prefix limit of 5.
func DefaultStatsOptions() StatsOptions {
	return StatsOptions{TopN: 5}
}

// StatsOptions controls Stats behaviour.
type StatsOptions struct {
	// TopN limits how many prefixes appear in TopPrefixes.
	TopN int
}

// Stats computes aggregate statistics for the given env map.
func Stats(env map[string]string, opts StatsOptions) StatsReport {
	seen := make(map[string]bool)
	prefixCounts := make(map[string]int)

	r := StatsReport{
		Total:        len(env),
		PrefixCounts: prefixCounts,
	}

	for k, v := range env {
		if v == "" {
			r.Empty++
		}
		if isSensitiveKey(k) {
			r.Sensitive++
		}
		if !seen[v] {
			seen[v] = true
			r.Unique++
		}
		if idx := strings.Index(k, "_"); idx > 0 {
			pfx := k[:idx]
			prefixCounts[pfx]++
		}
	}

	// Build sorted top-N prefix list.
	pairs := make([]PrefixCount, 0, len(prefixCounts))
	for p, c := range prefixCounts {
		pairs = append(pairs, PrefixCount{Prefix: p, Count: c})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count != pairs[j].Count {
			return pairs[i].Count > pairs[j].Count
		}
		return pairs[i].Prefix < pairs[j].Prefix
	})
	if opts.TopN > 0 && len(pairs) > opts.TopN {
		pairs = pairs[:opts.TopN]
	}
	r.TopPrefixes = pairs
	return r
}

// isSensitiveKey returns true when the key name suggests a secret.
func isSensitiveKey(k string) bool {
	upper := strings.ToUpper(k)
	for _, pat := range []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "CREDENTIAL", "PRIVATE"} {
		if strings.Contains(upper, pat) {
			return true
		}
	}
	return false
}
