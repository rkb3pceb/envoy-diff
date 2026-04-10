// Package policy provides rule-based enforcement for environment variable changes.
// Rules can block or warn when specific keys are added, removed, or modified.
package policy

import (
	"fmt"
	"strings"

	"github.com/your-org/envoy-diff/internal/diff"
)

// Severity represents the enforcement level of a policy rule.
type Severity string

const (
	SeverityBlock Severity = "block"
	SeverityWarn  Severity = "warn"
)

// Rule defines a single policy constraint.
type Rule struct {
	Name     string   `yaml:"name"`
	Keys     []string `yaml:"keys"`
	Types    []string `yaml:"types"` // added, removed, modified
	Severity Severity `yaml:"severity"`
	Message  string   `yaml:"message"`
}

// Violation is produced when a change breaches a rule.
type Violation struct {
	Rule     string
	Severity Severity
	Key      string
	Message  string
}

// Policy holds a set of rules and evaluates them against a diff.
type Policy struct {
	Rules []Rule
}

// Evaluate checks all changes against the policy rules and returns any violations.
func (p *Policy) Evaluate(changes []diff.Change) []Violation {
	var violations []Violation
	for _, change := range changes {
		for _, rule := range p.Rules {
			if matchesRule(rule, change) {
				msg := rule.Message
				if msg == "" {
					msg = fmt.Sprintf("key %q violated rule %q", change.Key, rule.Name)
				}
				violations = append(violations, Violation{
					Rule:     rule.Name,
					Severity: rule.Severity,
					Key:      change.Key,
					Message:  msg,
				})
			}
		}
	}
	return violations
}

// HasBlockers returns true if any violation has block severity.
func HasBlockers(violations []Violation) bool {
	for _, v := range violations {
		if v.Severity == SeverityBlock {
			return true
		}
	}
	return false
}

func matchesRule(rule Rule, change diff.Change) bool {
	keyMatch := false
	for _, k := range rule.Keys {
		if strings.EqualFold(k, change.Key) || strings.HasSuffix(k, "*") && strings.HasPrefix(change.Key, strings.TrimSuffix(k, "*")) {
			keyMatch = true
			break
		}
	}
	if !keyMatch {
		return false
	}
	for _, t := range rule.Types {
		if strings.EqualFold(t, string(change.Type)) {
			return true
		}
	}
	return false
}
