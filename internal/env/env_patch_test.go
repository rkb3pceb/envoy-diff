package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func baseMap() map[string]string {
	return map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_URL":   "postgres://localhost/dev",
	}
}

func TestPatchMap_SetAddsNewKey(t *testing.T) {
	ops := []PatchOp{{Op: "set", Key: "NEW_KEY", Value: "hello"}}
	r, err := PatchMap(baseMap(), ops, DefaultPatchOptions())
	require.NoError(t, err)
	assert.Equal(t, "hello", r.Env["NEW_KEY"])
	assert.Len(t, r.Applied, 1)
}

func TestPatchMap_SetOverwritesExisting(t *testing.T) {
	ops := []PatchOp{{Op: "set", Key: "APP_PORT", Value: "9090"}}
	r, err := PatchMap(baseMap(), ops, DefaultPatchOptions())
	require.NoError(t, err)
	assert.Equal(t, "9090", r.Env["APP_PORT"])
	assert.Len(t, r.Applied, 1)
}

func TestPatchMap_SetSkipsWhenOverwriteDisabled(t *testing.T) {
	opts := DefaultPatchOptions()
	opts.AllowOverwrite = false
	ops := []PatchOp{{Op: "set", Key: "APP_PORT", Value: "9090"}}
	r, err := PatchMap(baseMap(), ops, opts)
	require.NoError(t, err)
	assert.Equal(t, "8080", r.Env["APP_PORT"])
	assert.Len(t, r.Skipped, 1)
	assert.Len(t, r.Applied, 0)
}

func TestPatchMap_DeleteRemovesKey(t *testing.T) {
	ops := []PatchOp{{Op: "delete", Key: "APP_PORT"}}
	r, err := PatchMap(baseMap(), ops, DefaultPatchOptions())
	require.NoError(t, err)
	_, exists := r.Env["APP_PORT"]
	assert.False(t, exists)
	assert.Len(t, r.Applied, 1)
}

func TestPatchMap_DeleteMissingKey_Skipped(t *testing.T) {
	ops := []PatchOp{{Op: "delete", Key: "DOES_NOT_EXIST"}}
	r, err := PatchMap(baseMap(), ops, DefaultPatchOptions())
	require.NoError(t, err)
	assert.Len(t, r.Skipped, 1)
}

func TestPatchMap_DeleteMissingKey_ErrorOnMissing(t *testing.T) {
	opts := DefaultPatchOptions()
	opts.ErrorOnMissing = true
	ops := []PatchOp{{Op: "delete", Key: "DOES_NOT_EXIST"}}
	_, err := PatchMap(baseMap(), ops, opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DOES_NOT_EXIST")
}

func TestPatchMap_RenameKey(t *testing.T) {
	ops := []PatchOp{{Op: "rename", Key: "APP_HOST", To: "SERVICE_HOST"}}
	r, err := PatchMap(baseMap(), ops, DefaultPatchOptions())
	require.NoError(t, err)
	assert.Equal(t, "localhost", r.Env["SERVICE_HOST"])
	_, old := r.Env["APP_HOST"]
	assert.False(t, old)
}

func TestPatchMap_UnknownOp_ReturnsError(t *testing.T) {
	ops := []PatchOp{{Op: "upsert", Key: "X", Value: "y"}}
	_, err := PatchMap(baseMap(), ops, DefaultPatchOptions())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown op")
}

func TestPatchMap_DoesNotMutateInput(t *testing.T) {
	input := baseMap()
	ops := []PatchOp{{Op: "set", Key: "APP_PORT", Value: "1234"}}
	_, err := PatchMap(input, ops, DefaultPatchOptions())
	require.NoError(t, err)
	assert.Equal(t, "8080", input["APP_PORT"])
}

func TestHasPatchApplied_TrueWhenApplied(t *testing.T) {
	ops := []PatchOp{{Op: "set", Key: "NEW", Value: "val"}}
	r, _ := PatchMap(baseMap(), ops, DefaultPatchOptions())
	assert.True(t, HasPatchApplied(r))
}

func TestHasPatchApplied_FalseWhenNoneApplied(t *testing.T) {
	opts := DefaultPatchOptions()
	opts.AllowOverwrite = false
	ops := []PatchOp{{Op: "set", Key: "APP_PORT", Value: "9090"}}
	r, _ := PatchMap(baseMap(), ops, opts)
	assert.False(t, HasPatchApplied(r))
}
