package promote_test

import (
	"testing"

	"github.com/user/envoy-diff/internal/diff"
	"github.com/user/envoy-diff/internal/promote"
	"github.com/user/envoy-diff/internal/redact"
)

func stageFrom() promote.Stage {
	return promote.Stage{
		Name: "staging",
		Env: map[string]string{
			"APP_PORT": "8080",
			"DB_PASSWORD": "secret",
			"FEATURE_FLAG": "true",
		},
	}
}

func stageTo() promote.Stage {
	return promote.Stage{
		Name: "production",
		Env: map[string]string{
			"APP_PORT":    "8080",
			"DB_PASSWORD": "newSecret",
			"LOG_LEVEL":   "info",
		},
	}
}

func TestEvaluate_DetectsChanges(t *testing.T) {
	result, err := promote.Evaluate(stageFrom(), stageTo(), promote.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Changes) == 0 {
		t.Error("expected changes, got none")
	}
}

func TestEvaluate_StageMeta(t *testing.T) {
	result, _ := promote.Evaluate(stageFrom(), stageTo(), promote.DefaultOptions())
	if result.From != "staging" || result.To != "production" {
		t.Errorf("unexpected stage names: from=%q to=%q", result.From, result.To)
	}
}

func TestEvaluate_EmptyStageNameErrors(t *testing.T) {
	_, err := promote.Evaluate(promote.Stage{Name: ""}, stageTo(), promote.DefaultOptions())
	if err == nil {
		t.Error("expected error for empty stage name")
	}
}

func TestEvaluate_RedactsSensitiveValues(t *testing.T) {
	result, err := promote.Evaluate(stageFrom(), stageTo(), promote.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, c := range result.Changes {
		if c.Key == "DB_PASSWORD" {
			if c.OldValue == "secret" || c.NewValue == "newSecret" {
				t.Error("sensitive value was not redacted")
			}
			return
		}
	}
	t.Error("DB_PASSWORD change not found")
}

func TestEvaluate_NoRedact_WhenDisabled(t *testing.T) {
	opts := promote.Options{
		RedactCtx: redact.NewContext(false),
	}
	result, _ := promote.Evaluate(stageFrom(), stageTo(), opts)
	for _, c := range result.Changes {
		if c.Key == "DB_PASSWORD" && c.Type == diff.Modified {
			if c.OldValue != "secret" {
				t.Errorf("expected raw value, got %q", c.OldValue)
			}
			return
		}
	}
}

func TestEvaluate_AuditFindingsPresent(t *testing.T) {
	result, _ := promote.Evaluate(stageFrom(), stageTo(), promote.DefaultOptions())
	if len(result.Findings) == 0 {
		t.Error("expected audit findings for sensitive key change")
	}
}
