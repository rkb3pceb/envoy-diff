package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyDefaults_NoDefaults_ReturnsCopy(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	opts := DefaultDefaultOptions()
	out, res := ApplyDefaults(src, opts)
	assert.Equal(t, src, out)
	assert.False(t, res.HasApplied())
}

func TestApplyDefaults_MissingKey_AppliesDefault(t *testing.T) {
	src := map[string]string{"A": "1"}
	opts := DefaultDefaultOptions()
	opts.Defaults = map[string]string{"B": "default_b"}
	out, res := ApplyDefaults(src, opts)
	require.Equal(t, "default_b", out["B"])
	assert.Contains(t, res.Applied, "B")
	assert.Empty(t, res.Skipped)
}

func TestApplyDefaults_ExistingKey_IsSkipped(t *testing.T) {
	src := map[string]string{"A": "original"}
	opts := DefaultDefaultOptions()
	opts.Defaults = map[string]string{"A": "default_a"}
	out, res := ApplyDefaults(src, opts)
	assert.Equal(t, "original", out["A"])
	assert.Contains(t, res.Skipped, "A")
	assert.Empty(t, res.Applied)
}

func TestApplyDefaults_OverwriteEmpty_ReplacesBlankValue(t *testing.T) {
	src := map[string]string{"A": ""}
	opts := DefaultDefaultOptions()
	opts.OverwriteEmpty = true
	opts.Defaults = map[string]string{"A": "filled"}
	out, res := ApplyDefaults(src, opts)
	assert.Equal(t, "filled", out["A"])
	assert.Contains(t, res.Applied, "A")
}

func TestApplyDefaults_OverwriteEmpty_KeepsNonEmpty(t *testing.T) {
	src := map[string]string{"A": "keep"}
	opts := DefaultDefaultOptions()
	opts.OverwriteEmpty = true
	opts.Defaults = map[string]string{"A": "default_a"}
	out, res := ApplyDefaults(src, opts)
	assert.Equal(t, "keep", out["A"])
	assert.Contains(t, res.Skipped, "A")
}

func TestApplyDefaults_OverwriteAll_ReplacesEverything(t *testing.T) {
	src := map[string]string{"A": "original", "B": ""}
	opts := DefaultDefaultOptions()
	opts.OverwriteAll = true
	opts.Defaults = map[string]string{"A": "new_a", "B": "new_b", "C": "new_c"}
	out, res := ApplyDefaults(src, opts)
	assert.Equal(t, "new_a", out["A"])
	assert.Equal(t, "new_b", out["B"])
	assert.Equal(t, "new_c", out["C"])
	assert.Len(t, res.Applied, 3)
	assert.Empty(t, res.Skipped)
}

func TestApplyDefaults_DoesNotMutateInput(t *testing.T) {
	src := map[string]string{"A": "1"}
	opts := DefaultDefaultOptions()
	opts.Defaults = map[string]string{"B": "2"}
	_, _ = ApplyDefaults(src, opts)
	_, ok := src["B"]
	assert.False(t, ok, "original map must not be mutated")
}

func TestApplyDefaults_AppliedKeys_AreSorted(t *testing.T) {
	src := map[string]string{}
	opts := DefaultDefaultOptions()
	opts.Defaults = map[string]string{"Z": "z", "A": "a", "M": "m"}
	_, res := ApplyDefaults(src, opts)
	assert.Equal(t, []string{"A", "M", "Z"}, res.Applied)
}
