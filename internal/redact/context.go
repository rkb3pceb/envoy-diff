package redact

// Context holds redaction configuration for a diff session,
// combining pattern-based detection with custom rule overrides.
type Context struct {
	rules   *RuleSet
	level   MaskLevel
	enabled bool
}

// NewContext creates a Context with the given RuleSet and mask level.
// If rules is nil, the default pattern list is used via an empty RuleSet.
func NewContext(rules *RuleSet, level MaskLevel, enabled bool) *Context {
	if rules == nil {
		rules = Empty()
	}
	return &Context{
		rules:   rules,
		level:   level,
		enabled: enabled,
	}
}

// DefaultContext returns a Context with sensible defaults:
// partial masking enabled using the built-in pattern list.
func DefaultContext() *Context {
	return NewContext(Empty(), MaskPartial, true)
}

// Enabled reports whether redaction is active.
func (c *Context) Enabled() bool {
	return c.enabled
}

// ShouldRedact reports whether the given key should be redacted,
// checking both custom rules and built-in sensitive patterns.
func (c *Context) ShouldRedact(key string) bool {
	if !c.enabled {
		return false
	}
	if c.rules.Matches(key) != nil {
		return true
	}
	return IsSensitive(key)
}

// Apply returns the masked value for key if redaction applies,
// otherwise returns the original value unchanged.
func (c *Context) Apply(key, value string) string {
	if !c.ShouldRedact(key) {
		return value
	}
	// Prefer rule-specific placeholder when a custom rule matches.
	if r := c.rules.Matches(key); r != nil {
		return r.RedactValue(value, c.level)
	}
	return Mask(value, c.level)
}

// ApplyMap redacts all sensitive keys in the provided map,
// returning a new map with safe values replaced.
func (c *Context) ApplyMap(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = c.Apply(k, v)
	}
	return out
}
