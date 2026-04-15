package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
	"github.com/yourorg/envoy-diff/internal/parser"
)

var trimCmd = &cobra.Command{
	Use:   "trim <file>",
	Short: "Trim whitespace or affixes from env variable values",
	Args:  cobra.ExactArgs(1),
	RunE:  runTrim,
}

func init() {
	trimCmd.Flags().Bool("keys", false, "also trim whitespace from keys")
	trimCmd.Flags().String("prefix", "", "remove this prefix from all values")
	trimCmd.Flags().String("suffix", "", "remove this suffix from all values")
	trimCmd.Flags().String("chars", "", "trim these characters from both ends of values")
	trimCmd.Flags().Bool("no-trim-values", false, "disable default whitespace trimming on values")
	rootCmd.AddCommand(trimCmd)
}

func runTrim(cmd *cobra.Command, args []string) error {
	path := args[0]

	parsed, err := parser.ParseEnvFile(path)
	if err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	opts := env.DefaultTrimOptions()

	if v, _ := cmd.Flags().GetBool("keys"); v {
		opts.TrimKeys = true
	}
	if v, _ := cmd.Flags().GetBool("no-trim-values"); v {
		opts.TrimValues = false
	}
	if v, _ := cmd.Flags().GetString("prefix"); v != "" {
		opts.TrimPrefix = v
	}
	if v, _ := cmd.Flags().GetString("suffix"); v != "" {
		opts.TrimSuffix = v
	}
	if v, _ := cmd.Flags().GetString("chars"); v != "" {
		opts.TrimChars = v
	}

	result := env.TrimMap(parsed, opts)

	keys := make([]string, 0, len(result.Output))
	for k := range result.Output {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, result.Output[k])
	}

	if env.HasTrimChanges(result) {
		fmt.Fprintf(os.Stderr, "trimmed %d key(s)\n", len(result.Modified))
	}

	return nil
}
