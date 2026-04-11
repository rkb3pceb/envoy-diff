package merge_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/merge"
)

func TestMerge_LastWins(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "base"}
	b := map[string]string{"FOO": "2", "BAZ": "new"}

	res, err := merge.Merge(merge.StrategyLast, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["FOO"] != "2" {
		t.Errorf("expected FOO=2, got %s", res.Env["FOO"])
	}
	if res.Env["BAR"] != "base" {
		t.Errorf("expected BAR=base, got %s", res.Env["BAR"])
	}
	if res.Env["BAZ"] != "new" {
		t.Errorf("expected BAZ=new, got %s", res.Env["BAZ"])
	}
}

func TestMerge_FirstWins(t *testing.T) {
	a := map[string]string{"FOO": "original"}
	b := map[string]string{"FOO": "override"}

	res, err := merge.Merge(merge.StrategyFirst, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["FOO"] != "original" {
		t.Errorf("expected FOO=original, got %s", res.Env["FOO"])
	}
}

func TestMerge_ErrorOnConflict(t *testing.T) {
	a := map[string]string{"SECRET": "abc"}
	b := map[string]string{"SECRET": "xyz"}

	res, err := merge.Merge(merge.StrategyError, a, b)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
	if len(res.Conflicts) == 0 {
		t.Error("expected conflicts to be recorded")
	}
}

func TestMerge_NoConflict_NoError(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}

	res, err := merge.Merge(merge.StrategyError, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Env))
	}
}

func TestMerge_TracksOverrides(t *testing.T) {
	a := map[string]string{"KEY": "v1"}
	b := map[string]string{"KEY": "v2"}

	res, _ := merge.Merge(merge.StrategyLast, a, b)
	if len(res.Overrides["KEY"]) < 2 {
		t.Errorf("expected at least 2 override entries for KEY, got %v", res.Overrides["KEY"])
	}
}

func TestMerge_EmptySources(t *testing.T) {
	res, err := merge.Merge(merge.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 0 {
		t.Errorf("expected empty env, got %d keys", len(res.Env))
	}
}
