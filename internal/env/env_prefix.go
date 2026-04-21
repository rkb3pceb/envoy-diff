package env

import "strings"

// DefaultPrefixOptions returns safe defaults for prefix operations.
func DefaultPrefixOptions() PrefixOptions {
	return PrefixOptions{
		IgnoreCase: false,
		StripExisting: false,
	}
}

// PrefixOptions controls how AddPrefix and StripPrefix behave.
type PrefixOptions struct {
	// IgnoreCase makes prefix matching case-insensitive during strip.
	IgnoreCase bool
	// StripExisting removes the prefix before adding, avoiding double-prefix.
	StripExisting bool
}

// AddPrefix returns a new map with prefix prepended to every key.
// If StripExisting is true the prefix is stripped first so it is never doubled.
func AddPrefix(m map[string]string, prefix string, opts PrefixOptions) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		key := k
		if opts.StripExisting {
			key = stripPrefix(key, prefix, opts.IgnoreCase)
		}
		out[prefix+key] = v
	}
	return out
}

// StripPrefix returns a new map with prefix removed from every matching key.
// Keys that do not carry the prefix are passed through unchanged.
func StripPrefix(m map[string]string, prefix string, opts PrefixOptions) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[stripPrefix(k, prefix, opts.IgnoreCase)] = v
	}
	return out
}

// HasPrefixedKeys reports whether any key in m starts with prefix.
func HasPrefixedKeys(m map[string]string, prefix string, ignoreCase bool) bool {
	for k := range m {
		if hasPrefix(k, prefix, ignoreCase) {
			return true
		}
	}
	return false
}

func stripPrefix(key, prefix string, ignoreCase bool) string {
	if hasPrefix(key, prefix, ignoreCase) {
		return key[len(prefix):]
	}
	return key
}

func hasPrefix(key, prefix string, ignoreCase bool) bool {
	if ignoreCase {
		return strings.HasPrefix(strings.ToLower(key), strings.ToLower(prefix))
	}
	return strings.HasPrefix(key, prefix)
}
