package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterMap_NoOptions_ReturnsAll(t *testing.T) {
	src := map[string]string{"APP_HOST": "localhost", "DB_URL": "postgres://", "PORT": "8080"}
	out := FilterMap(src, DefaultFilterOptions())
	assert.Equal(t, src, out)
}

func TestFilterMap_PrefixFilter_KeepsMatching(t *testing.T) {
	src := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_URL": "postgres://"}
	opts := DefaultFilterOptions()
	opts.Prefixes = []string{"APP_"}
	out := FilterMap(src, opts)
	assert.Equal(t, map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}, out)
}

func TestFilterMap_ExcludePrefix_RemovesMatching(t *testing.T) {
	src := map[string]string{"APP_HOST": "localhost", "SECRET_KEY": "abc", "PORT": "8080"}
	opts := DefaultFilterOptions()
	opts.ExcludePrefixes = []string{"SECRET_"}
	out := FilterMap(src, opts)
	assert.NotContains(t, out, "SECRET_KEY")
	assert.Contains(t, out, "APP_HOST")
	assert.Contains(t, out, "PORT")
}

func TestFilterMap_Contains_CaseInsensitive(t *testing.T) {
	src := map[string]string{"DATABASE_URL": "pg://", "DB_PASS": "s3cr3t", "APP_HOST": "localhost"}
	opts := DefaultFilterOptions()
	opts.Contains = "db"
	out := FilterMap(src, opts)
	assert.Contains(t, out, "DATABASE_URL")
	assert.Contains(t, out, "DB_PASS")
	assert.NotContains(t, out, "APP_HOST")
}

func TestFilterMap_OnlyNonEmpty_RemovesBlanks(t *testing.T) {
	src := map[string]string{"HOST": "localhost", "TOKEN": "", "PORT": "8080"}
	opts := DefaultFilterOptions()
	opts.OnlyNonEmpty = true
	out := FilterMap(src, opts)
	assert.NotContains(t, out, "TOKEN")
	assert.Len(t, out, 2)
}

func TestFilterMap_MultiplePrefixes_KeepsAll(t *testing.T) {
	src := map[string]string{"APP_X": "1", "DB_Y": "2", "LOG_Z": "3", "OTHER": "4"}
	opts := DefaultFilterOptions()
	opts.Prefixes = []string{"APP_", "DB_"}
	out := FilterMap(src, opts)
	assert.Len(t, out, 2)
	assert.Contains(t, out, "APP_X")
	assert.Contains(t, out, "DB_Y")
}

func TestFilterMap_DoesNotMutateInput(t *testing.T) {
	src := map[string]string{"A": "1", "B": ""}
	opts := DefaultFilterOptions()
	opts.OnlyNonEmpty = true
	_ = FilterMap(src, opts)
	assert.Len(t, src, 2, "original map must not be modified")
}

func TestHasFilteredKeys_TrueWhenSomethingRemoved(t *testing.T) {
	src := map[string]string{"APP_A": "1", "DB_B": "2"}
	opts := DefaultFilterOptions()
	opts.Prefixes = []string{"APP_"}
	assert.True(t, HasFilteredKeys(src, opts))
}

func TestHasFilteredKeys_FalseWhenNothingRemoved(t *testing.T) {
	src := map[string]string{"APP_A": "1", "APP_B": "2"}
	opts := DefaultFilterOptions()
	opts.Prefixes = []string{"APP_"}
	assert.False(t, HasFilteredKeys(src, opts))
}
