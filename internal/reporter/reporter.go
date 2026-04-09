// Package reporter provides summary reporting for envoy-diff audit results.
package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/envoy-diff/internal/audit"
	"github.com/yourorg/envoy-diff/internal/diff"
)

// Summary holds aggregated counts and metadata for a diff+audit run.
type Summary struct {
	TotalKeys  int
	Added      int
	Removed    int
	Modified   int
	Unchanged  int
	HighRisk   int
	MediumRisk int
	LowRisk    int
}

// Build constructs a Summary from a slice of changes and audit findings.
func Build(changes []diff.Change, findings []audit.Finding) Summary {
	s := Summary{TotalKeys: len(changes)}
	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			s.Added++
		case diff.Removed:
			s.Removed++
		case diff.Modified:
			s.Modified++
		case diff.Unchanged:
			s.Unchanged++
		}
	}
	for _, f := range findings {
		switch f.Severity {
		case audit.High:
			s.HighRisk++
		case audit.Medium:
			s.MediumRisk++
		case audit.Low:
			s.LowRisk++
		}
	}
	return s
}

// Write renders the summary as a human-readable block to w.
func Write(w io.Writer, s Summary) {
	fmt.Fprintln(w, strings.Repeat("─", 40))
	fmt.Fprintln(w, "SUMMARY")
	fmt.Fprintln(w, strings.Repeat("─", 40))
	fmt.Fprintf(w, "  Total keys : %d\n", s.TotalKeys)
	fmt.Fprintf(w, "  Added      : %d\n", s.Added)
	fmt.Fprintf(w, "  Removed    : %d\n", s.Removed)
	fmt.Fprintf(w, "  Modified   : %d\n", s.Modified)
	fmt.Fprintf(w, "  Unchanged  : %d\n", s.Unchanged)
	if s.HighRisk+s.MediumRisk+s.LowRisk > 0 {
		fmt.Fprintln(w, strings.Repeat("─", 40))
		fmt.Fprintln(w, "AUDIT FINDINGS")
		fmt.Fprintf(w, "  High       : %d\n", s.HighRisk)
		fmt.Fprintf(w, "  Medium     : %d\n", s.MediumRisk)
		fmt.Fprintf(w, "  Low        : %d\n", s.LowRisk)
	}
	fmt.Fprintln(w, strings.Repeat("─", 40))
}
