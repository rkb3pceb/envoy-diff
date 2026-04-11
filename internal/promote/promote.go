// Package promote provides utilities for promoting environment variable
// sets between deployment stages (e.g. staging → production), applying
// redaction, policy checks, and generating a diff summary.
package promote

import (
	"fmt"

	"github.com/user/envoy-diff/internal/audit"
	"github.com/user/envoy-diff/internal/diff"
	"github.com/user/envoy-diff/internal/policy"
	"github.com/user/envoy-diff/internal/redact"
)

// Stage represents a named deployment stage.
type Stage struct {
	Name string
	Env  map[string]string
}

// Result holds the outcome of a promotion evaluation.
type Result struct {
	From       string
	To         string
	Changes    []diff.Change
	Findings   []audit.Finding
	Violations []policy.Violation
	Blocked    bool
}

// Options controls promotion behaviour.
type Options struct {
	RedactCtx *redact.Context
	Policy    *policy.Policy
}

// DefaultOptions returns Options with safe defaults.
func DefaultOptions() Options {
	return Options{
		RedactCtx: redact.DefaultContext(),
		Policy:    policy.Empty(),
	}
}

// Evaluate compares two stages, audits the diff, and checks policy.
func Evaluate(from, to Stage, opts Options) (*Result, error) {
	if from.Name == "" || to.Name == "" {
		return nil, fmt.Errorf("promote: stage names must not be empty")
	}

	changes := diff.Compare(from.Env, to.Env)

	if opts.RedactCtx != nil {
		for i, c := range changes {
			changes[i].OldValue = opts.RedactCtx.Apply(c.Key, c.OldValue)
			changes[i].NewValue = opts.RedactCtx.Apply(c.Key, c.NewValue)
		}
	}

	findings := audit.Audit(changes)

	var violations []policy.Violation
	blocked := false
	if opts.Policy != nil {
		violations = opts.Policy.Evaluate(changes)
		blocked = policy.HasBlockers(violations)
	}

	return &Result{
		From:       from.Name,
		To:         to.Name,
		Changes:    changes,
		Findings:   findings,
		Violations: violations,
		Blocked:    blocked,
	}, nil
}
