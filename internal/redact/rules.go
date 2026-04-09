package redact

// Rule defines a custom redaction rule with a key pattern and optional replacement.
type Rule struct {
	// Pattern is a substring or exact key name to match against env var keys.
	Pattern string
	// Placeholder overrides the default redaction placeholder for this rule.
	Placeholder string
}

// RuleSet holds a collection of user-defined redaction rules.
type RuleSet struct {
	rules []Rule
}

// NewRuleSet constructs a RuleSet from a slice of Rules.
func NewRuleSet(rules []Rule) *RuleSet {
	return &RuleSet{rules: rules}
}

// Empty returns a RuleSet with no rules.
func Empty() *RuleSet {
	return &RuleSet{}
}

// Matches reports whether the given key matches any rule in the set.
// It returns the matching Rule and true if found.
func (rs *RuleSet) Matches(key string) (Rule, bool) {
	for _, r := range rs.rules {
		if containsFold(key, r.Pattern) {
			return r, true
		}
	}
	return Rule{}, false
}

// RedactValue applies the first matching rule to value, returning the
// redacted string. If no rule matches and the key is not sensitive by
// default pattern, the original value is returned unchanged.
func (rs *RuleSet) RedactValue(key, value string) string {
	if r, ok := rs.Matches(key); ok {
		ph := r.Placeholder
		if ph == "" {
			ph = Placeholder
		}
		return ph
	}
	// Fall back to built-in sensitivity check.
	return Value(key, value)
}

// Len returns the number of rules in the set.
func (rs *RuleSet) Len() int {
	return len(rs.rules)
}

// containsFold is a case-insensitive substring check.
func containsFold(s, substr string) bool {
	if substr == "" {
		return false
	}
	s2 := toLower(s)
	sub2 := toLower(substr)
	return contains(s2, sub2)
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		b[i] = c
	}
	return string(b)
}

func contains(s, sub string) bool {
	if len(sub) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
