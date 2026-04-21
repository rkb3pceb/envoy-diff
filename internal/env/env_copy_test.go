package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyKeys_NoKeys_CopiesAll(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{}
	opts := DefaultCopyOptions()

	res, err := CopyKeys(src, dst, nil, opts)
	require.NoError(t, err)
	assert.Len(t, res.Copied, 2)
	assert.Equal(t, "1", dst["A"])
	assert.Equal(t, "2", dst["B"])
}

func TestCopyKeys_SpecificKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}
	opts := DefaultCopyOptions()

	res, err := CopyKeys(src, dst, []string{"A", "C"}, opts)
	require.NoError(t, err)
	assert.Len(t, res.Copied, 2)
	assert.Equal(t, "1", dst["A"])
	assert.Equal(t, "3", dst["C"])
	_, exists := dst["B"]
	assert.False(t, exists)
}

func TestCopyKeys_DestPrefix(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	dst := map[string]string{}
	opts := DefaultCopyOptions()
	opts.DestPrefix = "NEW_"

	res, err := CopyKeys(src, dst, []string{"FOO"}, opts)
	require.NoError(t, err)
	assert.Len(t, res.Copied, 1)
	assert.Equal(t, "bar", dst["NEW_FOO"])
}

func TestCopyKeys_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := map[string]string{"X": "new"}
	dst := map[string]string{"X": "old"}
	opts := DefaultCopyOptions()

	res, err := CopyKeys(src, dst, []string{"X"}, opts)
	require.NoError(t, err)
	assert.Len(t, res.Skipped, 1)
	assert.Equal(t, "old", dst["X"])
}

func TestCopyKeys_OverwriteExisting(t *testing.T) {
	src := map[string]string{"X": "new"}
	dst := map[string]string{"X": "old"}
	opts := DefaultCopyOptions()
	opts.Overwrite = true

	res, err := CopyKeys(src, dst, []string{"X"}, opts)
	require.NoError(t, err)
	assert.Len(t, res.Copied, 1)
	assert.Equal(t, "new", dst["X"])
}

func TestCopyKeys_MissingKey_SilentByDefault(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{}
	opts := DefaultCopyOptions()

	res, err := CopyKeys(src, dst, []string{"MISSING"}, opts)
	require.NoError(t, err)
	assert.Len(t, res.Missing, 1)
	assert.Len(t, res.Copied, 0)
}

func TestCopyKeys_MissingKey_ErrorOnMissing(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{}
	opts := DefaultCopyOptions()
	opts.ErrorOnMissing = true

	_, err := CopyKeys(src, dst, []string{"MISSING"}, opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MISSING")
}

func TestHasCopied_TrueWhenCopied(t *testing.T) {
	assert.True(t, HasCopied(CopyResult{Copied: []string{"K"}}))
	assert.False(t, HasCopied(CopyResult{}))
}

func TestCopyKeys_DestPrefix_DoesNotAffectOriginalKey(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	dst := map[string]string{}
	opts := DefaultCopyOptions()
	opts.DestPrefix = "NEW_"

	_, err := CopyKeys(src, dst, []string{"FOO"}, opts)
	require.NoError(t, err)
	_, exists := dst["FOO"]
	assert.False(t, exists, "original key should not be present in dst when DestPrefix is set")
	assert.Equal(t, "bar", dst["NEW_FOO"])
}
