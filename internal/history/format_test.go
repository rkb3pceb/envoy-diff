package history_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/envoy-diff/internal/history"
)

func sampleEntry(id string) history.Entry {
	return history.Entry{
		ID:        id,
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		OldFile:   "configs/staging.env",
		NewFile:   "configs/production.env",
		Added:     3,
		Removed:   1,
		Modified:  2,
		Findings:  1,
	}
}

func TestPrint_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	history.Print(&buf, nil)
	if !strings.Contains(buf.String(), "No history") {
		t.Errorf("expected no-history message, got: %s", buf.String())
	}
}

func TestPrint_ShowsHeader(t *testing.T) {
	var buf bytes.Buffer
	history.Print(&buf, []history.Entry{sampleEntry("abc123")})
	out := buf.String()
	for _, col := range []string{"ID", "TIMESTAMP", "OLD FILE", "NEW FILE", "FINDINGS"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected column %q in output", col)
		}
	}
}

func TestPrint_ShowsEntryData(t *testing.T) {
	var buf bytes.Buffer
	history.Print(&buf, []history.Entry{sampleEntry("deadbeef-long-id")})
	out := buf.String()
	if !strings.Contains(out, "deadbeef") {
		t.Errorf("expected truncated ID in output")
	}
	if !strings.Contains(out, "staging.env") {
		t.Errorf("expected base filename in output")
	}
	if !strings.Contains(out, "production.env") {
		t.Errorf("expected base filename in output")
	}
}

func TestPrint_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	entries := []history.Entry{sampleEntry("id-one"), sampleEntry("id-two")}
	history.Print(&buf, entries)
	out := buf.String()
	if !strings.Contains(out, "id-one") || !strings.Contains(out, "id-two") {
		t.Errorf("expected both entries in output")
	}
}
