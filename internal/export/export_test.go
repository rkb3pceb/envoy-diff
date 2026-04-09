package export_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/envoy-diff/internal/diff"
	"github.com/yourorg/envoy-diff/internal/export"
)

var sampleChanges = []diff.Change{
	{Key: "DB_HOST", Type: diff.Modified, OldValue: "localhost", NewValue: "prod-db"},
	{Key: "API_KEY", Type: diff.Added, OldValue: "", NewValue: "secret123"},
	{Key: "LEGACY", Type: diff.Removed, OldValue: "old", NewValue: ""},
}

var sampleOpts = export.Options{
	Format:    export.FormatCSV,
	OldFile:   "old.env",
	NewFile:   "new.env",
	Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
}

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []export.Format{export.FormatCSV, export.FormatMarkdown} {
		e, err := export.New(f)
		if err != nil {
			t.Errorf("New(%q) unexpected error: %v", f, err)
		}
		if e == nil {
			t.Errorf("New(%q) returned nil exporter", f)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := export.New("xml")
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestCSVExporter_WritesHeader(t *testing.T) {
	e, _ := export.New(export.FormatCSV)
	var buf bytes.Buffer
	if err := e.Write(&buf, sampleChanges, sampleOpts); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if lines[0] != "key,type,old_value,new_value" {
		t.Errorf("unexpected header: %s", lines[0])
	}
	if len(lines) != 4 { // header + 3 rows
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
}

func TestCSVExporter_NoChanges(t *testing.T) {
	e, _ := export.New(export.FormatCSV)
	var buf bytes.Buffer
	if err := e.Write(&buf, nil, sampleOpts); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header line, got %d lines", len(lines))
	}
}

func TestMarkdownExporter_ContainsTableHeader(t *testing.T) {
	e, _ := export.New(export.FormatMarkdown)
	var buf bytes.Buffer
	opts := sampleOpts
	opts.Format = export.FormatMarkdown
	if err := e.Write(&buf, sampleChanges, opts); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "| Key | Type |") {
		t.Error("expected markdown table header not found")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected key DB_HOST in output")
	}
}

func TestMarkdownExporter_NoChanges(t *testing.T) {
	e, _ := export.New(export.FormatMarkdown)
	var buf bytes.Buffer
	if err := e.Write(&buf, nil, sampleOpts); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if !strings.Contains(buf.String(), "No changes detected") {
		t.Error("expected no-changes message in markdown output")
	}
}
