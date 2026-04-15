package env

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestEncodeMap_Base64_AllKeys(t *testing.T) {
	env := map[string]string{
		"TOKEN": "secret",
		"NAME":  "alice",
	}
	opts := DefaultEncodeOptions()
	res := EncodeMap(env, opts)

	if len(res.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	expect := base64.StdEncoding.EncodeToString([]byte("secret"))
	if res.Map["TOKEN"] != expect {
		t.Errorf("TOKEN: got %q, want %q", res.Map["TOKEN"], expect)
	}
	if len(res.Encoded) != 2 {
		t.Errorf("expected 2 encoded keys, got %d", len(res.Encoded))
	}
}

func TestEncodeMap_Base64URL(t *testing.T) {
	env := map[string]string{"DATA": "hello world"}
	opts := DefaultEncodeOptions()
	opts.Format = EncodeBase64URL
	res := EncodeMap(env, opts)

	expect := base64.URLEncoding.EncodeToString([]byte("hello world"))
	if res.Map["DATA"] != expect {
		t.Errorf("got %q, want %q", res.Map["DATA"], expect)
	}
}

func TestEncodeMap_Hex(t *testing.T) {
	env := map[string]string{"VAL": "abc"}
	opts := DefaultEncodeOptions()
	opts.Format = EncodeHex
	res := EncodeMap(env, opts)

	expect := fmt.Sprintf("%x", []byte("abc"))
	if res.Map["VAL"] != expect {
		t.Errorf("got %q, want %q", res.Map["VAL"], expect)
	}
}

func TestEncodeMap_SpecificKeys_OnlyEncodesThem(t *testing.T) {
	env := map[string]string{
		"SECRET": "mysecret",
		"PLAIN":  "plaintext",
	}
	opts := DefaultEncodeOptions()
	opts.Keys = []string{"SECRET"}
	res := EncodeMap(env, opts)

	if res.Map["PLAIN"] != "plaintext" {
		t.Errorf("PLAIN should be unchanged, got %q", res.Map["PLAIN"])
	}
	if len(res.Encoded) != 1 || res.Encoded[0] != "SECRET" {
		t.Errorf("expected only SECRET to be encoded, got %v", res.Encoded)
	}
}

func TestEncodeMap_SkipEmpty_LeavesEmptyUntouched(t *testing.T) {
	env := map[string]string{"EMPTY": "", "FULL": "value"}
	opts := DefaultEncodeOptions()
	opts.SkipEmpty = true
	res := EncodeMap(env, opts)

	if res.Map["EMPTY"] != "" {
		t.Errorf("expected EMPTY to be untouched, got %q", res.Map["EMPTY"])
	}
}

func TestEncodeMap_Decode_RoundTrip(t *testing.T) {
	original := map[string]string{"TOKEN": "roundtrip-value"}

	encOpts := DefaultEncodeOptions()
	encoded := EncodeMap(original, encOpts)

	decOpts := DefaultEncodeOptions()
	decOpts.Decode = true
	decoded := EncodeMap(encoded.Map, decOpts)

	if decoded.Map["TOKEN"] != original["TOKEN"] {
		t.Errorf("round-trip failed: got %q, want %q", decoded.Map["TOKEN"], original["TOKEN"])
	}
}

func TestEncodeMap_InvalidBase64_ReturnsError(t *testing.T) {
	env := map[string]string{"BAD": "!!!not-base64!!!"}
	opts := DefaultEncodeOptions()
	opts.Decode = true
	res := EncodeMap(env, opts)

	if len(res.Errors) == 0 {
		t.Error("expected an error for invalid base64 input")
	}
	if _, ok := res.Map["BAD"]; ok {
		if res.Map["BAD"] != "!!!not-base64!!!" {
			t.Error("original value should be preserved on error")
		}
	}
}

func TestEncodeMap_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	opts := DefaultEncodeOptions()
	_ = EncodeMap(env, opts)

	if env["KEY"] != "value" {
		t.Error("EncodeMap must not mutate the input map")
	}
}
