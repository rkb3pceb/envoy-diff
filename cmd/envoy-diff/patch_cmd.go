package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
	"github.com/yourorg/envoy-diff/internal/parser"
)

var (
	patchSet    []string
	patchDelete []string
	patchRename []string
	patchStrict bool
)

func init() {
	patchCmd := &cobra.Command{
		Use:   "patch <env-file>",
		Short: "Apply set/delete/rename patches to an env file",
		Args:  cobra.ExactArgs(1),
		RunE:  runPatch,
	}

	patchCmd.Flags().StringArrayVar(&patchSet, "set", nil, "Set a key: KEY=VALUE")
	patchCmd.Flags().StringArrayVar(&patchDelete, "delete", nil, "Delete a key by name")
	patchCmd.Flags().StringArrayVar(&patchRename, "rename", nil, "Rename a key: OLD=NEW")
	patchCmd.Flags().BoolVar(&patchStrict, "strict", false, "Error if a key to delete/rename is missing")

	rootCmd.AddCommand(patchCmd)
}

func runPatch(cmd *cobra.Command, args []string) error {
	parsed, err := parser.ParseEnvFile(args[0])
	if err != nil {
		return fmt.Errorf("patch: reading %s: %w", args[0], err)
	}

	var ops []env.PatchOp

	for _, s := range patchSet {
		k, v, ok := strings.Cut(s, "=")
		if !ok {
			return fmt.Errorf("patch: --set %q must be KEY=VALUE", s)
		}
		ops = append(ops, env.PatchOp{Op: "set", Key: k, Value: v})
	}

	for _, k := range patchDelete {
		ops = append(ops, env.PatchOp{Op: "delete", Key: k})
	}

	for _, s := range patchRename {
		old, newKey, ok := strings.Cut(s, "=")
		if !ok {
			return fmt.Errorf("patch: --rename %q must be OLD=NEW", s)
		}
		ops = append(ops, env.PatchOp{Op: "rename", Key: old, To: newKey})
	}

	opts := env.DefaultPatchOptions()
	opts.ErrorOnMissing = patchStrict

	result, err := env.PatchMap(parsed, ops, opts)
	if err != nil {
		return fmt.Errorf("patch: %w", err)
	}

	keys := env.SortedKeys(result.Env, env.DefaultSortOptions())
	for _, k := range keys {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, result.Env[k])
	}

	if len(result.Skipped) > 0 {
		fmt.Fprintf(os.Stderr, "skipped %d op(s)\n", len(result.Skipped))
	}

	return nil
}
