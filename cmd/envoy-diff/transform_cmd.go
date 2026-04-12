package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
	"github.com/yourorg/envoy-diff/internal/parser"
)

var transformCmd = &cobra.Command{
	Use:   "transform <file>",
	Short: "Apply key/value transformations to an env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runTransform,
}

func init() {
	transformCmd.Flags().String("prefix-add", "", "prepend a prefix to every key")
	transformCmd.Flags().String("prefix-strip", "", "strip a prefix from matching keys")
	transformCmd.Flags().Bool("uppercase", false, "convert all keys to UPPER_CASE")
	transformCmd.Flags().Bool("lowercase", false, "convert all keys to lower_case")
	transformCmd.Flags().Bool("drop-empty", false, "remove entries with empty values")
	rootCmd.AddCommand(transformCmd)
}

func runTransform(cmd *cobra.Command, args []string) error {
	path := args[0]

	parsed, err := parser.ParseEnvFile(path)
	if err != nil {
		return fmt.Errorf("transform: parse %q: %w", path, err)
	}

	opts := env.DefaultTransformOptions()
	opts.PrefixAdd, _ = cmd.Flags().GetString("prefix-add")
	opts.PrefixStrip, _ = cmd.Flags().GetString("prefix-strip")
	opts.UppercaseKeys, _ = cmd.Flags().GetBool("uppercase")
	opts.LowercaseKeys, _ = cmd.Flags().GetBool("lowercase")
	opts.DropEmpty, _ = cmd.Flags().GetBool("drop-empty")

	out, err := env.TransformMap(parsed, opts)
	if err != nil {
		return fmt.Errorf("transform: %w", err)
	}

	keys := make([]string, 0, len(out))
	for k := range out {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := cmd.OutOrStdout()
	if w == nil {
		w = os.Stdout
	}
	for _, k := range keys {
		fmt.Fprintf(w, "%s=%s\n", k, out[k])
	}
	return nil
}
