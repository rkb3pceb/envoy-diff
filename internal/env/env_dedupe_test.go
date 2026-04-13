package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDedupeMap_NoDuplicates_ReturnsCopy(t *testing.T) {
	src := []map[string]string{{"A": "1", "B": "2"}}
	res := DedupeMap(src, DefaultDedupeOptions())
	assert.False(t, res.HasDuplicates())
	assert.Equal(t, map[string]string{"A": "1", "B": "2"}, res.Map)
}

func TestDedupeMap_LastWins_ByDefault(t *testing.T) {
	sources := []map[string]string{
		{"KEY": "first"},
		{"KEY": "second"},
	}
	res := DedupeMap(sources, DefaultDedupeOptions())
	require.True(t, res.HasDuplicates())
	assert.Equal(t, "second", res.Map["KEY"])
	assert.Equal(t, []string{"first", "second"}, res.Duplicates[0].Values)
}

func TestDedupeMap_KeepFirst_RetainsFirstValue(t *testing.T) {
	sources := []map[string]string{
		{"KEY": "first"},
		{"KEY": "second"},
	}
	opts := DefaultDedupeOptions()
	opts.KeepFirst = true
	res := DedupeMap(sources, opts)
	assert.Equal(t, "first", res.Map["KEY"])
}

func TestDedupeMap_CaseInsensitive_DetectsDuplicates(t *testing.T) {
	sources := []map[string]string{
		{"key": "lower"},
		{"KEY": "upper"},
	}
	opts := DefaultDedupeOptions()
	opts.CaseInsensitive = true
	res := DedupeMap(sources, opts)
	assert.True(t, res.HasDuplicates())
	assert.Len(t, res.Duplicates, 1)
}

func TestDedupeMap_ReportOnly_DoesNotPopulateMap(t *testing.T) {
	sources := []map[string]string{
		{"X": "1"},
		{"X": "2"},
	}
	opts := DefaultDedupeOptions()
	opts.ReportOnly = true
	res := DedupeMap(sources, opts)
	assert.True(t, res.HasDuplicates())
	assert.Empty(t, res.Map)
}

func TestDedupeMap_MultipleSources_AllDuplicatesReported(t *testing.T) {
	sources := []map[string]string{
		{"A": "1", "B": "x"},
		{"A": "2", "C": "y"},
		{"A": "3", "B": "z"},
	}
	res := DedupeMap(sources, DefaultDedupeOptions())
	assert.Len(t, res.Duplicates, 2)
	assert.Equal(t, "3", res.Map["A"])
	assert.Equal(t, "z", res.Map["B"])
	assert.Equal(t, "y", res.Map["C"])
}
