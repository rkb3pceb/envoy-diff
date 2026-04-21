package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectMap_NoOptions_ReturnsCopy(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	out := SelectMap(src, DefaultSelectOptions())
	assert.Equal(t, src, out)
	// must be a copy
	out["C"] = "3"
	assert.NotContains(t, src, "C")
}

func TestSelectMap_SpecificKeys(t *testing.T) {
	src := map[string]string{"FOO": "a", "BAR": "b", "BAZ": "c"}
	opts := DefaultSelectOptions()
	opts.Keys = []string{"FOO", "BAZ"}
	out := SelectMap(src, opts)
	assert.Equal(t, map[string]string{"FOO": "a", "BAZ": "c"}, out)
}

func TestSelectMap_PrefixFilter(t *testing.T) {
	src := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_URL": "pg://"}
	opts := DefaultSelectOptions()
	opts.Prefixes = []string{"APP_"}
	out := SelectMap(src, opts)
	assert.Equal(t, map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}, out)
}

func TestSelectMap_CaseInsensitive_Keys(t *testing.T) {
	src := map[string]string{"FOO": "val", "BAR": "other"}
	opts := DefaultSelectOptions()
	opts.CaseSensitive = false
	opts.Keys = []string{"foo"}
	out := SelectMap(src, opts)
	assert.Equal(t, map[string]string{"FOO": "val"}, out)
}

func TestSelectMap_CaseInsensitive_Prefix(t *testing.T) {
	src := map[string]string{"APP_HOST": "h", "DB_URL": "u"}
	opts := DefaultSelectOptions()
	opts.CaseSensitive = false
	opts.Prefixes = []string{"app_"}
	out := SelectMap(src, opts)
	assert.Equal(t, map[string]string{"APP_HOST": "h"}, out)
}

func TestSelectMap_Invert_ReturnsNonMatching(t *testing.T) {
	src := map[string]string{"FOO": "1", "BAR": "2", "BAZ": "3"}
	opts := DefaultSelectOptions()
	opts.Keys = []string{"FOO"}
	opts.Invert = true
	out := SelectMap(src, opts)
	assert.NotContains(t, out, "FOO")
	assert.Contains(t, out, "BAR")
	assert.Contains(t, out, "BAZ")
}

func TestSelectMap_Invert_WithPrefix(t *testing.T) {
	src := map[string]string{"APP_HOST": "h", "APP_PORT": "p", "DB_URL": "u"}
	opts := DefaultSelectOptions()
	opts.Prefixes = []string{"APP_"}
	opts.Invert = true
	out := SelectMap(src, opts)
	assert.Equal(t, map[string]string{"DB_URL": "u"}, out)
}

func TestHasSelectedKeys_True(t *testing.T) {
	assert.True(t, HasSelectedKeys(map[string]string{"X": "1"}))
}

func TestHasSelectedKeys_False(t *testing.T) {
	assert.False(t, HasSelectedKeys(map[string]string{}))
}
