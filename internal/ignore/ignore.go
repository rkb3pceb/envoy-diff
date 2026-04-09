// Package ignore provides functionality to load and apply ignore rules
// for environment variable keys during diff and audit operations.
package ignore

import (
	"bufio"
	"os"
	"strings"
)

// Rules holds a set of key patterns to ignore during diffing.
type Rules struct {
	exact    map[string]struct{}
	prefixes []string
}

// Load reads an ignore file from the given path and returns a Rules instance.
// Each line in the file is either an exact key name or a prefix ending with '*'.
// Lines starting with '#' and blank lines are ignored.
func Load(path string) (*Rules, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := &Rules{
		exact: make(map[string]struct{}),
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasSuffix(line, "*") {
			r.prefixes = append(r.prefixes, strings.TrimSuffix(line, "*"))
		} else {
			r.exact[line] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return r, nil
}

// Empty returns a Rules instance with no ignore rules applied.
func Empty() *Rules {
	return &Rules{exact: make(map[string]struct{})}
}

// Match reports whether the given key matches any ignore rule.
func (r *Rules) Match(key string) bool {
	if _, ok := r.exact[key]; ok {
		return true
	}
	for _, prefix := range r.prefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

// FilterKeys returns only the keys from the provided slice that do NOT match
// any ignore rule.
func (r *Rules) FilterKeys(keys []string) []string {
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		if !r.Match(k) {
			out = append(out, k)
		}
	}
	return out
}
