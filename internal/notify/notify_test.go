package notify_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/user/envoy-diff/internal/audit"
	"github.com/user/envoy-diff/internal/diff"
	"github.com/user/envoy-diff/internal/notify"
)

// fakeNotifier records whether Send was called.
type fakeNotifier struct {
	called bool
	err    error
}

func (f *fakeNotifier) Send(_ notify.Event) error {
	f.called = true
	return f.err
}

func highFinding() audit.Finding {
	return audit.Finding{Key: "SECRET", Severity: audit.SeverityHigh, Message: "modified"}
}

func medFinding() audit.Finding {
	return audit.Finding{Key: "DB_PASS", Severity: audit.SeverityMedium, Message: "added"}
}

func sampleEvent(findings ...audit.Finding) notify.Event {
	return notify.Event{
		Changes:  []diff.Change{{Key: "SECRET", Type: diff.Modified}},
		Findings: findings,
		Summary:  "1 change, 1 finding",
	}
}

func TestDispatch_FiresOnHighThreshold(t *testing.T) {
	n := &fakeNotifier{}
	d := notify.New(notify.LevelHigh, n)
	_ = d.Dispatch(sampleEvent(highFinding()))
	if !n.called {
		t.Fatal("expected notifier to be called for high severity")
	}
}

func TestDispatch_SkipsWhenBelowThreshold(t *testing.T) {
	n := &fakeNotifier{}
	d := notify.New(notify.LevelHigh, n)
	_ = d.Dispatch(sampleEvent(medFinding()))
	if n.called {
		t.Fatal("expected notifier NOT to be called when only medium findings")
	}
}

func TestDispatch_LevelAny_AlwaysFires(t *testing.T) {
	n := &fakeNotifier{}
	d := notify.New(notify.LevelAny, n)
	_ = d.Dispatch(sampleEvent()) // no findings
	if !n.called {
		t.Fatal("expected notifier to be called for LevelAny")
	}
}

func TestDispatch_ReturnsErrorFromNotifier(t *testing.T) {
	n := &fakeNotifier{err: errors.New("boom")}
	d := notify.New(notify.LevelAny, n)
	err := d.Dispatch(sampleEvent())
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}

func TestStdoutNotifier_WritesOutput(t *testing.T) {
	var buf bytes.Buffer
	sn := notify.NewStdoutNotifier(&buf)
	e := sampleEvent(highFinding())
	if err := sn.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[envoy-diff]") {
		t.Errorf("expected header in output, got: %s", out)
	}
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected finding key in output, got: %s", out)
	}
}
