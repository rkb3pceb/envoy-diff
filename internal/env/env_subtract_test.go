package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubtractMap_NoOptions_ReturnsCopy(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	res := SubtractMap(base, DefaultSubtractOptions())
	assert.Equal(t, base, res.Result)
	assert.Empty(t, res.Removed)
	assert.False(t, res.HasSubtracted())
}

func TestSubtractMap_RemovesExplicitKeys(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	opts := DefaultSubtractOptions()
	opts.Keys = []string{"A", "C"}
	res := SubtractMap(base, opts)
	assert.Equal(t, map[string]string{"B": "2"}, res.Result)
	assert.Equal(t, []string{"A", "C"}, res.Removed)
	assert.True(t, res.HasSubtracted())
}

func TestSubtractMap_RemovesByPrefix(t *testing.T) {
	base := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "APP_NAME": "myapp"}
	opts := DefaultSubtractOptions()
	opts.Prefixes = []string{"DB_"}
	res := SubtractMap(base, opts)
	assert.Equal(t, map[string]string{"APP_NAME": "myapp"}, res.Result)
	assert.ElementsMatch(t, []string{"DB_HOST", "DB_PORT"}, res.Removed)
}

func TestSubtractMap_CaseInsensitive_Key(t *testing.T) {
	base := map[string]string{"Secret": "abc", "KEEP": "yes"}
	opts := DefaultSubtractOptions()
	opts.Keys = []string{"secret"}
	opts.CaseInsensitive = true
	res := SubtractMap(base, opts)
	assert.Equal(t, map[string]string{"KEEP": "yes"}, res.Result)
	assert.Equal(t, []string{"Secret"}, res.Removed)
}

func TestSubtractMap_CaseInsensitive_Prefix(t *testing.T) {
	base := map[string]string{"Db_Host": "h", "DB_PORT": "p", "APP": "a"}
	opts := DefaultSubtractOptions()
	opts.Prefixes = []string{"db_"}
	opts.CaseInsensitive = true
	res := SubtractMap(base, opts)
	assert.Equal(t, map[string]string{"APP": "a"}, res.Result)
	assert.ElementsMatch(t, []string{"Db_Host", "DB_PORT"}, res.Removed)
}

func TestSubtractMap_DoesNotMutateInput(t *testing.T) {
	base := map[string]string{"X": "1", "Y": "2"}
	opts := DefaultSubtractOptions()
	opts.Keys = []string{"X"}
	_ = SubtractMap(base, opts)
	assert.Equal(t, map[string]string{"X": "1", "Y": "2"}, base)
}

func TestSubtractMap_MissingKey_IsIgnored(t *testing.T) {
	base := map[string]string{"A": "1"}
	opts := DefaultSubtractOptions()
	opts.Keys = []string{"NONEXISTENT"}
	res := SubtractMap(base, opts)
	assert.Equal(t, base, res.Result)
	assert.Empty(t, res.Removed)
}
