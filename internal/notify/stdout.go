package notify

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// StdoutNotifier writes a human-readable notification to an io.Writer.
type StdoutNotifier struct {
	w io.Writer
}

// NewStdoutNotifier returns a StdoutNotifier writing to w.
// If w is nil, os.Stdout is used.
func NewStdoutNotifier(w io.Writer) *StdoutNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &StdoutNotifier{w: w}
}

// Send prints the event summary and findings to the writer.
func (s *StdoutNotifier) Send(e Event) error {
	lines := []string{
		"[envoy-diff] notification",
		fmt.Sprintf("  summary : %s", e.Summary),
		fmt.Sprintf("  changes : %d", len(e.Changes)),
		fmt.Sprintf("  findings: %d", len(e.Findings)),
	}
	for _, f := range e.Findings {
		lines = append(lines, fmt.Sprintf("    [%s] %s — %s", f.Severity, f.Key, f.Message))
	}
	_, err := fmt.Fprintln(s.w, strings.Join(lines, "\n"))
	return err
}

// SendAll sends multiple events in sequence, returning the first error
// encountered. Remaining events are still attempted after an error.
func (s *StdoutNotifier) SendAll(events []Event) error {
	var firstErr error
	for _, e := range events {
		if err := s.Send(e); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
