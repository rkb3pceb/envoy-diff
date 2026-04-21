package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUniqueMap_NoDuplicates_ReturnsCopy(t *testing.T) {
	m := map[string]string{"A": "1", "B": "2", "C": "3"}
	r := UniqueMap(m, DefaultUniqueOptions())
	assert.Equal(t, m, r.Map)
	assert.Empty(t, r.Removed)
	assert.False(t, HasUniqueChanges(r))
}

func TestUniqueMap_KeepFirst_RemovesLaterDuplicate(t *testing.T) {
	m := map[string]string{"A": "same", "B": "same", "C": "other"}
	opts := DefaultUniqueOptions()
	opts.KeepFirst = true
	r := UniqueMap(m, opts)
	// sorted keys: A, B, C — A wins, B removed
	require.Contains(t, r.Map, "A")
	assert.NotContains(t, r.Map, "B")
	assert.Contains(t, r.Map, "C")
	assert.Equal(t, []string{"B"}, r.Removed)
	assert.True(t, HasUniqueChanges(r))
}

func TestUniqueMap_KeepLast_RemovesEarlierDuplicate(t *testing.T) {
	m := map[string]string{"A": "same", "B": "same", "C": "other"}
	opts := DefaultUniqueOptions()
	opts.KeepFirst = false
	r := UniqueMap(m, opts)
	// sorted keys: A, B — B wins, A removed
	assert.NotContains(t, r.Map, "A")
	assert.Contains(t, r.Map, "B")
	assert.Equal(t, []string{"A"}, r.Removed)
}

func TestUniqueMap_CaseInsensitive_DetectsDuplicates(t *testing.T) {
	m := map[string]string{"A": "Hello", "B": "hello", "C": "world"}
	opts := DefaultUniqueOptions()
	opts.CaseSensitive = false
	r := UniqueMap(m, opts)
	assert.Len(t, r.Map, 2)
	assert.Equal(t, []string{"B"}, r.Removed)
}

func TestUniqueMap_CaseSensitive_KeepsBothVariants(t *testing.T) {
	m := map[string]string{"A": "Hello", "B": "hello"}
	opts := DefaultUniqueOptions()
	opts.CaseSensitive = true
	r := UniqueMap(m, opts)
	assert.Len(t, r.Map, 2)
	assert.Empty(t, r.Removed)
}

func TestUniqueMap_ByValueFalse_ReturnsCopyUnchanged(t *testing.T) {
	m := map[string]string{"X": "dup", "Y": "dup"}
	opts := DefaultUniqueOptions()
	opts.ByValue = false
	r := UniqueMap(m, opts)
	assert.Equal(t, m, r.Map)
	assert.Empty(t, r.Removed)
}

func TestUniqueMap_DoesNotMutateInput(t *testing.T) {
	m := map[string]string{"A": "v", "B": "v"}
	orig := map[string]string{"A": "v", "B": "v"}
	UniqueMap(m, DefaultUniqueOptions())
	assert.Equal(t, orig, m)
}

func TestUniqueMap_MultipleDuplicateGroups(t *testing.T) {
	m := map[string]string{
		"A": "alpha", "B": "alpha",
		"C": "beta", "D": "beta",
		"E": "gamma",
	}
	r := UniqueMap(m, DefaultUniqueOptions())
	assert.Len(t, r.Map, 3)
	assert.Len(t, r.Removed, 2)
}
