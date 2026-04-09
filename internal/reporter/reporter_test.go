package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envoy-diff/internal/audit"
	"github.com/yourorg/envoy-diff/internal/diff"
	"github.com/yourorg/envoy-diff/internal/reporter"
)

func TestBuild_CountsChanges(t *testing.T) {
	changes := []diff.Change{
		{Key: "A", Type: diff.Added},
		{Key: "B", Type: diff.Removed},
		{Key: "C", Type: diff.Modified},
		{Key: "D", Type: diff.Unchanged},
		{Key: "E", Type: diff.Unchanged},
	}
	s := reporter.Build(changes, nil)
	if s.TotalKeys != 5 {
		t.Errorf("TotalKeys: want 5, got %d", s.TotalKeys)
	}
	if s.Added != 1 || s.Removed != 1 || s.Modified != 1 || s.Unchanged != 2 {
		t.Errorf("unexpected change counts: %+v", s)
	}
}

func TestBuild_CountsFindings(t *testing.T) {
	findings := []audit.Finding{
		{Severity: audit.High},
		{Severity: audit.High},
		{Severity: audit.Medium},
		{Severity: audit.Low},
	}
	s := reporter.Build(nil, findings)
	if s.HighRisk != 2 {
		t.Errorf("HighRisk: want 2, got %d", s.HighRisk)
	}
	if s.MediumRisk != 1 {
		t.Errorf("MediumRisk: want 1, got %d", s.MediumRisk)
	}
	if s.LowRisk != 1 {
		t.Errorf("LowRisk: want 1, got %d", s.LowRisk)
	}
}

func TestWrite_ContainsSummaryHeader(t *testing.T) {
	var buf bytes.Buffer
	s := reporter.Summary{TotalKeys: 3, Added: 1, Removed: 1, Modified: 1}
	reporter.Write(&buf, s)
	out := buf.String()
	if !strings.Contains(out, "SUMMARY") {
		t.Error("expected SUMMARY header in output")
	}
	if !strings.Contains(out, "Total keys") {
		t.Error("expected 'Total keys' in output")
	}
}

func TestWrite_AuditSectionOnlyWhenFindings(t *testing.T) {
	var buf bytes.Buffer
	reporter.Write(&buf, reporter.Summary{TotalKeys: 2, Unchanged: 2})
	if strings.Contains(buf.String(), "AUDIT FINDINGS") {
		t.Error("should not print audit section when no findings")
	}

	buf.Reset()
	reporter.Write(&buf, reporter.Summary{HighRisk: 1})
	if !strings.Contains(buf.String(), "AUDIT FINDINGS") {
		t.Error("expected audit section when findings present")
	}
}
