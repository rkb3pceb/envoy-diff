package env

import (
	"strings"
	"testing"
)

func TestTransformMap_NoOp(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := TransformMap(src, DefaultTransformOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(src) {
		t.Fatalf("expected %d entries, got %d", len(src), len(out))
	}
}

func TestTransformMap_PrefixAdd(t *testing.T) {
	src := map[string]string{"HOST": "localhost"}
	opts := DefaultTransformOptions()
	opts.PrefixAdd = "APP_"
	out, _ := TransformMap(src, opts)
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected key APP_HOST")
	}
}

func TestTransformMap_PrefixStrip(t *testing.T) {
	src := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	opts := DefaultTransformOptions()
	opts.PrefixStrip = "APP_"
	out, _ := TransformMap(src, opts)
	if _, ok := out["HOST"]; !ok {
		t.Error("expected key HOST after strip")
	}
	if _, ok := out["PORT"]; !ok {
		t.Error("expected key PORT after strip")
	}
}

func TestTransformMap_UppercaseKeys(t *testing.T) {
	src := map[string]string{"db_host": "localhost"}
	opts := DefaultTransformOptions()
	opts.UppercaseKeys = true
	out, _ := TransformMap(src, opts)
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected uppercased key DB_HOST")
	}
}

func TestTransformMap_LowercaseKeys(t *testing.T) {
	src := map[string]string{"DB_HOST": "localhost"}
	opts := DefaultTransformOptions()
	opts.LowercaseKeys = true
	out, _ := TransformMap(src, opts)
	if _, ok := out["db_host"]; !ok {
		t.Error("expected lowercased key db_host")
	}
}

func TestTransformMap_MutuallyExclusiveCaseFlags(t *testing.T) {
	opts := DefaultTransformOptions()
	opts.UppercaseKeys = true
	opts.LowercaseKeys = true
	_, err := TransformMap(map[string]string{}, opts)
	if err == nil {
		t.Error("expected error for conflicting case flags")
	}
}

func TestTransformMap_DropEmpty(t *testing.T) {
	src := map[string]string{"FOO": "bar", "EMPTY": ""}
	opts := DefaultTransformOptions()
	opts.DropEmpty = true
	out, _ := TransformMap(src, opts)
	if _, ok := out["EMPTY"]; ok {
		t.Error("expected EMPTY to be dropped")
	}
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
}

func TestTransformMap_ExtraFn(t *testing.T) {
	src := map[string]string{"SECRET": "topsecret", "HOST": "localhost"}
	opts := DefaultTransformOptions()
	opts.Extra = []TransformFn{
		func(k, v string) (string, string, bool) {
			if strings.Contains(k, "SECRET") {
				return k, "***", true
			}
			return k, v, true
		},
	}
	out, _ := TransformMap(src, opts)
	if out["SECRET"] != "***" {
		t.Errorf("expected SECRET to be redacted, got %q", out["SECRET"])
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST unchanged, got %q", out["HOST"])
	}
}

func TestTransformMap_DoesNotMutateInput(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	opts := DefaultTransformOptions()
	opts.PrefixAdd = "X_"
	TransformMap(src, opts)
	if _, ok := src["FOO"]; !ok {
		t.Error("original map was mutated")
	}
}
