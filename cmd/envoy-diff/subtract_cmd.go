package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
	"github.com/yourorg/envoy-diff/internal/parser"
)

var subtractCmd = &cobra.Command{
	Use:   "subtract <file>",
	Short: "Remove keys from an env file by name or prefix",
	Args:  cobra.ExactArgs(1),
	RunE:  runSubtract,
}

func init() {
	subtractCmd.Flags().StringSliceP("key", "k", nil, "Exact key(s) to remove")
	subtractCmd.Flags().StringSliceP("prefix", "p", nil, "Remove keys matching prefix(es)")
	subtractCmd.Flags().BoolP("ignore-case", "i", false, "Case-insensitive key and prefix matching")
	subtractCmd.Flags().BoolP("quiet", "q", false, "Suppress removed-key summary")
	rootCmd.AddCommand(subtractCmd)
}

func runSubtract(cmd *cobra.Command, args []string) error {
	keys, _ := cmd.Flags().GetStringSlice("key")
	prefixes, _ := cmd.Flags().GetStringSlice("prefix")
	ci, _ := cmd.Flags().GetBool("ignore-case")
	quiet, _ := cmd.Flags().GetBool("quiet")

	if len(keys) == 0 && len(prefixes) == 0 {
		return fmt.Errorf("at least one --key or --prefix must be specified")
	}

	m, err := parser.ParseEnvFile(args[0])
	if err != nil {
		return fmt.Errorf("parse %s: %w", args[0], err)
	}

	opts := env.DefaultSubtractOptions()
	opts.Keys = keys
	opts.Prefixes = prefixes
	opts.CaseInsensitive = ci

	res := env.SubtractMap(m, opts)

	for _, k := range env.SortedKeys(res.Result) {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, res.Result[k])
	}

	if !quiet && res.HasSubtracted() {
		fmt.Fprintf(os.Stderr, "# removed %d key(s): %s\n",
			len(res.Removed), strings.Join(res.Removed, ", "))
	}

	return nil
}
