package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
)

var intersectCmd = &cobra.Command{
	Use:   "intersect <file1> <file2> [fileN...]",
	Short: "Output keys common to all provided env files",
	Args:  cobra.MinimumNArgs(2),
	RunE:  runIntersect,
}

func init() {
	intersectCmd.Flags().Bool("require-equal", false, "only include keys whose values are identical across all files")
	intersectCmd.Flags().String("keep", "last", "which file's value to keep: first or last")
	rootCmd.AddCommand(intersectCmd)
}

func runIntersect(cmd *cobra.Command, args []string) error {
	requireEqual, _ := cmd.Flags().GetBool("require-equal")
	keep, _ := cmd.Flags().GetString("keep")

	if keep != "first" && keep != "last" {
		return fmt.Errorf("--keep must be 'first' or 'last', got %q", keep)
	}

	maps := make([]map[string]string, 0, len(args))
	for _, path := range args {
		m, err := env.LoadSources(path)
		if err != nil {
			return fmt.Errorf("loading %s: %w", path, err)
		}
		maps = append(maps, env.Flatten(m))
	}

	opts := env.DefaultIntersectOptions()
	opts.KeepValues = keep
	opts.RequireEqual = requireEqual

	result := env.IntersectMaps(opts, maps...)

	if !env.HasIntersectResult(result) {
		fmt.Fprintln(os.Stderr, "no common keys found")
		return nil
	}

	keys := make([]string, 0, len(result.Map))
	for k := range result.Map {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, result.Map[k])
	}

	if result.Dropped > 0 || result.Conflicts > 0 {
		fmt.Fprintf(os.Stderr, "# dropped=%d conflicts=%d\n", result.Dropped, result.Conflicts)
	}

	return nil
}
