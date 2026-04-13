package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeMaps_NoConflict_MergesAll(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"C": "3"}
	res, err := MergeMaps([]map[string]string{a, b}, DefaultMergeOptions())
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"A": "1", "B": "2", "C": "3"}, res.Merged)
	assert.Empty(t, res.Conflicts)
}

func TestMergeMaps_LastWins_OverwritesKey(t *testing.T) {
	a := map[string]string{"KEY": "old"}
	b := map[string]string{"KEY": "new"}
	opts := DefaultMergeOptions()
	res, err := MergeMaps([]map[string]string{a, b}, opts)
	require.NoError(t, err)
	assert.Equal(t, "new", res.Merged["KEY"])
	assert.Contains(t, res.Conflicts, "KEY")
}

func TestMergeMaps_FirstWins_KeepsOriginal(t *testing.T) {
	a := map[string]string{"KEY": "original"}
	b := map[string]string{"KEY": "ignored"}
	opts := MergeOptions{Strategy: MergeFirstWins}
	res, err := MergeMaps([]map[string]string{a, b}, opts)
	require.NoError(t, err)
	assert.Equal(t, "original", res.Merged["KEY"])
	assert.Contains(t, res.Conflicts, "KEY")
}

func TestMergeMaps_ErrorOnConflict_ReturnsError(t *testing.T) {
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "2"}
	opts := MergeOptions{Strategy: MergeErrorOnConflict}
	_, err := MergeMaps([]map[string]string{a, b}, opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "X")
}

func TestMergeMaps_OverridesAppliedLast(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	opts := MergeOptions{
		Strategy:  MergeLastWins,
		Overrides: map[string]string{"B": "override", "C": "new"},
	}
	res, err := MergeMaps([]map[string]string{a}, opts)
	require.NoError(t, err)
	assert.Equal(t, "override", res.Merged["B"])
	assert.Equal(t, "new", res.Merged["C"])
}

func TestMergeMaps_EmptySources_ReturnsEmpty(t *testing.T) {
	res, err := MergeMaps([]map[string]string{}, DefaultMergeOptions())
	require.NoError(t, err)
	assert.Empty(t, res.Merged)
	assert.Empty(t, res.Conflicts)
}

func TestHasMergeConflicts_TrueWhenConflicts(t *testing.T) {
	r := MergeResult{Conflicts: []string{"KEY"}}
	assert.True(t, HasMergeConflicts(r))
}

func TestHasMergeConflicts_FalseWhenNone(t *testing.T) {
	r := MergeResult{}
	assert.False(t, HasMergeConflicts(r))
}
