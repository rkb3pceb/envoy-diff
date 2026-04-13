package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
)

var (
	maskLevel        string
	maskPlaceholder  string
	maskExtraPattern []string
	maskShowKeys     bool
)

func init() {
	maskCmd := &cobra.Command{
		Use:   "mask <file>",
		Short: "Mask sensitive values in an env file",
		Args:  cobra.ExactArgs(1),
		RunE:  runMask,
	}
	maskCmd.Flags().StringVar(&maskLevel, "level", "full", "Mask level: full or partial")
	maskCmd.Flags().StringVar(&maskPlaceholder, "placeholder", "[REDACTED]", "Placeholder for fully masked values")
	maskCmd.Flags().StringArrayVar(&maskExtraPattern, "pattern", nil, "Extra sensitive key patterns (repeatable)")
	maskCmd.Flags().BoolVar(&maskShowKeys, "show-keys", false, "Print which keys were masked to stderr")
	rootCmd.AddCommand(maskCmd)
}

func runMask(cmd *cobra.Command, args []string) error {
	sources, err := env.LoadSources(args)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}
	m := env.Flatten(sources)

	opts := env.DefaultMaskOptions()
	opts.Level = env.MaskLevel(maskLevel)
	opts.Placeholder = maskPlaceholder
	opts.ExtraPatterns = maskExtraPattern

	result := env.MaskMap(m, opts)

	if maskShowKeys && len(result.MaskedKeys) > 0 {
		sort.Strings(result.MaskedKeys)
		fmt.Fprintln(os.Stderr, "masked keys:")
		for _, k := range result.MaskedKeys {
			fmt.Fprintf(os.Stderr, "  %s\n", k)
		}
	}

	keys := env.Keys(result.Map)
	for _, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, result.Map[k])
	}
	return nil
}
