package resolve_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/resolve"
)

func TestExpand_NoReferences(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	got := resolve.Expand("plain-value", env, resolve.DefaultOptions())
	if got != "plain-value" {
		t.Errorf("expected plain-value, got %q", got)
	}
}

func TestExpand_ResolvesReference(t *testing.T) {
	env := map[string]string{"HOST": "db.internal", "DSN": "postgres://${HOST}/mydb"}
	got := resolve.Expand(env["DSN"], env, resolve.DefaultOptions())
	if got != "postgres://db.internal/mydb" {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestExpand_MissingRef_LeftAsIs(t *testing.T) {
	env := map[string]string{}
	got := resolve.Expand("${MISSING}", env, resolve.DefaultOptions())
	if got != "${MISSING}" {
		t.Errorf("expected reference unchanged, got %q", got)
	}
}

func TestMap_ResolvesAllValues(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  "${BASE_URL}/api/v1",
		"PLAIN":    "no-refs",
	}
	out := resolve.Map(env, resolve.DefaultOptions())
	if out["API_URL"] != "https://example.com/api/v1" {
		t.Errorf("API_URL not resolved: %q", out["API_URL"])
	}
	if out["PLAIN"] != "no-refs" {
		t.Errorf("PLAIN changed unexpectedly: %q", out["PLAIN"])
	}
}

func TestMap_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"A": "${B}", "B": "hello"}
	resolve.Map(env, resolve.DefaultOptions())
	if env["A"] != "${B}" {
		t.Error("original map was mutated")
	}
}

func TestHasUnresolved_True(t *testing.T) {
	if !resolve.HasUnresolved("${FOO}") {
		t.Error("expected unresolved reference to be detected")
	}
}

func TestHasUnresolved_False(t *testing.T) {
	if resolve.HasUnresolved("fully-resolved-value") {
		t.Error("expected no unresolved references")
	}
}

func TestUnresolvedKeys_ReturnsMissingRefs(t *testing.T) {
	env := map[string]string{
		"GOOD": "static",
		"BAD":  "${MISSING_VAR}",
	}
	keys := resolve.UnresolvedKeys(env)
	if len(keys) != 1 || keys[0] != "BAD" {
		t.Errorf("expected [BAD], got %v", keys)
	}
}
