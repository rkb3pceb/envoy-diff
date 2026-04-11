package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/envoy-diff/internal/env"
	"github.com/your-org/envoy-diff/internal/redact"
)

var (
	flattenShowOverrides bool
	flattenNoRedact      bool
)

func init() {
	flattenCmd := &cobra.Command{
		Use:   "flatten <file1> <file2> [fileN...]",
		Short: "Merge multiple env files into a single resolved map",
		Args:  cobra.MinimumNArgs(2),
		RunE:  runFlatten,
	}
	flattenCmd.Flags().BoolVar(&flattenShowOverrides, "show-overrides", false, "print keys that were overridden")
	flattenCmd.Flags().BoolVar(&flattenNoRedact, "no-redact", false, "disable sensitive value redaction")
	rootCmd.AddCommand(flattenCmd)
}

func runFlatten(cmd *cobra.Command, args []string) error {
	sources, err := env.LoadSources(args)
	if err != nil {
		return fmt.Errorf("flatten: %w", err)
	}

	merged, overridden := env.Flatten(sources)

	if flattenShowOverrides && len(overridden) > 0 {
		fmt.Fprintln(os.Stderr, "# overridden keys:")
		for _, k := range overridden {
			fmt.Fprintf(os.Stderr, "#   %s\n", k)
		}
	}

	for _, k := range env.Keys(merged) {
		v := merged[k]
		if !flattenNoRedact && redact.IsSensitive(k) {
			v = redact.Value(k, v)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
	}
	return nil
}
