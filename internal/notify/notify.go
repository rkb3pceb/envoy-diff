// Package notify provides notification hooks that fire when environment
// variable diffs exceed configurable severity thresholds.
package notify

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envoy-diff/internal/audit"
	"github.com/user/envoy-diff/internal/diff"
)

// Level represents the minimum severity that triggers a notification.
type Level string

const (
	LevelHigh   Level = "high"
	LevelMedium Level = "medium"
	LevelAny    Level = "any"
)

// Event holds the data passed to every Notifier.
type Event struct {
	Changes  []diff.Change
	Findings []audit.Finding
	Summary  string
}

// Notifier is the interface implemented by all notification backends.
type Notifier interface {
	Send(e Event) error
}

// Dispatcher routes an Event to one or more Notifiers based on threshold.
type Dispatcher struct {
	notifiers []Notifier
	threshold Level
	out        io.Writer
}

// New creates a Dispatcher with the given threshold and notifiers.
func New(threshold Level, notifiers ...Notifier) *Dispatcher {
	return &Dispatcher{
		notifiers: notifiers,
		threshold: threshold,
		out:        os.Stderr,
	}
}

// Dispatch evaluates findings against the threshold and fires notifiers.
func (d *Dispatcher) Dispatch(e Event) error {
	if !d.shouldFire(e.Findings) {
		return nil
	}
	var errs []string
	for _, n := range d.notifiers {
		if err := n.Send(e); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("notify: %s", strings.Join(errs, "; "))
	}
	return nil
}

func (d *Dispatcher) shouldFire(findings []audit.Finding) bool {
	if d.threshold == LevelAny {
		return true
	}
	for _, f := range findings {
		switch d.threshold {
		case LevelHigh:
			if f.Severity == audit.SeverityHigh {
				return true
			}
		case LevelMedium:
			if f.Severity == audit.SeverityHigh || f.Severity == audit.SeverityMedium {
				return true
			}
		}
	}
	return false
}
