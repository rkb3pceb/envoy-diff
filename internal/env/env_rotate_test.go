package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRotateMap_NoRules_ReturnsCopy(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	opts := DefaultRotateOptions()
	out, res, err := RotateMap(src, opts)
	require.NoError(t, err)
	assert.Equal(t, src, out)
	assert.False(t, res.HasRotated())
}

func TestRotateMap_BasicRename(t *testing.T) {
	src := map[string]string{"OLD_KEY": "value", "OTHER": "x"}
	opts := DefaultRotateOptions()
	opts.Rules = map[string]string{"OLD_KEY": "NEW_KEY"}
	out, res, err := RotateMap(src, opts)
	require.NoError(t, err)
	assert.Equal(t, "value", out["NEW_KEY"])
	_, oldPresent := out["OLD_KEY"]
	assert.False(t, oldPresent)
	assert.Contains(t, res.Rotated, "OLD_KEY")
}

func TestRotateMap_KeepOld_RetainsBothKeys(t *testing.T) {
	src := map[string]string{"OLD_KEY": "value"}
	opts := DefaultRotateOptions()
	opts.Rules = map[string]string{"OLD_KEY": "NEW_KEY"}
	opts.KeepOld = true
	out, res, err := RotateMap(src, opts)
	require.NoError(t, err)
	assert.Equal(t, "value", out["NEW_KEY"])
	assert.Equal(t, "value", out["OLD_KEY"])
	assert.Contains(t, res.Rotated, "OLD_KEY")
}

func TestRotateMap_MissingKey_SkippedByDefault(t *testing.T) {
	src := map[string]string{"A": "1"}
	opts := DefaultRotateOptions()
	opts.Rules = map[string]string{"MISSING": "NEW"}
	out, res, err := RotateMap(src, opts)
	require.NoError(t, err)
	_, newPresent := out["NEW"]
	assert.False(t, newPresent)
	assert.Contains(t, res.Skipped, "MISSING")
}

func TestRotateMap_MissingKey_ErrorOnMissing(t *testing.T) {
	src := map[string]string{"A": "1"}
	opts := DefaultRotateOptions()
	opts.Rules = map[string]string{"MISSING": "NEW"}
	opts.ErrorOnMissing = true
	_, _, err := RotateMap(src, opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MISSING")
}

func TestRotateMap_ConflictKey_SkippedByDefault(t *testing.T) {
	src := map[string]string{"OLD": "v1", "NEW": "v2"}
	opts := DefaultRotateOptions()
	opts.Rules = map[string]string{"OLD": "NEW"}
	out, res, err := RotateMap(src, opts)
	require.NoError(t, err)
	assert.Equal(t, "v2", out["NEW"]) // original preserved
	assert.Contains(t, res.Conflict, "NEW")
	assert.Contains(t, res.Skipped, "OLD")
}

func TestRotateMap_ConflictKey_ErrorOnConflict(t *testing.T) {
	src := map[string]string{"OLD": "v1", "NEW": "v2"}
	opts := DefaultRotateOptions()
	opts.Rules = map[string]string{"OLD": "NEW"}
	opts.ErrorOnConflict = true
	_, _, err := RotateMap(src, opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "NEW")
}

func TestRotateMap_DoesNotMutateInput(t *testing.T) {
	src := map[string]string{"OLD": "val"}
	orig := map[string]string{"OLD": "val"}
	opts := DefaultRotateOptions()
	opts.Rules = map[string]string{"OLD": "NEW"}
	_, _, err := RotateMap(src, opts)
	require.NoError(t, err)
	assert.Equal(t, orig, src)
}
