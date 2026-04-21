package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
	"github.com/yourorg/envoy-diff/internal/parser"
)

var uniqueCmd = &cobra.Command{
	Use:   "unique <file>",
	Short: "Remove keys with duplicate values from an env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runUnique,
}

func init() {
	uniqueCmd.Flags().Bool("keep-last", false, "Keep the last key for each duplicate value instead of the first")
	uniqueCmd.Flags().Bool("case-insensitive", false, "Treat values as equal regardless of case")
	uniqueCmd.Flags().Bool("show-removed", false, "Print removed keys to stderr")
	rootCmd.AddCommand(uniqueCmd)
}

func runUnique(cmd *cobra.Command, args []string) error {
	path := args[0]

	raw, err := parser.ParseEnvFile(path)
	if err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	keepLast, _ := cmd.Flags().GetBool("keep-last")
	caseInsensitive, _ := cmd.Flags().GetBool("case-insensitive")
	showRemoved, _ := cmd.Flags().GetBool("show-removed")

	opts := env.DefaultUniqueOptions()
	opts.KeepFirst = !keepLast
	opts.CaseSensitive = !caseInsensitive

	result := env.UniqueMap(raw, opts)

	if showRemoved && env.HasUniqueChanges(result) {
		fmt.Fprintf(os.Stderr, "# removed duplicate-value keys: %s\n",
			strings.Join(result.Removed, ", "))
	}

	keys := env.SortedKeys(result.Map, env.DefaultSortOptions())
	for _, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, result.Map[k])
	}

	return nil
}
