package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
	"github.com/yourorg/envoy-diff/internal/parser"
)

var (
	prefixAdd        string
	prefixStrip      string
	prefixIgnoreCase bool
	prefixNoDouble   bool
)

func init() {
	prefixCmd := &cobra.Command{
		Use:   "prefix <file>",
		Short: "Add or strip a prefix from env variable keys",
		Args:  cobra.ExactArgs(1),
		RunE:  runPrefix,
	}
	prefixCmd.Flags().StringVar(&prefixAdd, "add", "", "prefix to add to all keys")
	prefixCmd.Flags().StringVar(&prefixStrip, "strip", "", "prefix to strip from matching keys")
	prefixCmd.Flags().BoolVar(&prefixIgnoreCase, "ignore-case", false, "case-insensitive prefix matching")
	prefixCmd.Flags().BoolVar(&prefixNoDouble, "no-double", false, "strip existing prefix before adding (prevents doubling)")
	rootCmd.AddCommand(prefixCmd)
}

func runPrefix(cmd *cobra.Command, args []string) error {
	if prefixAdd == "" && prefixStrip == "" {
		return fmt.Errorf("at least one of --add or --strip is required")
	}

	parsed, err := parser.ParseEnvFile(args[0])
	if err != nil {
		return fmt.Errorf("parse %s: %w", args[0], err)
	}

	opts := env.DefaultPrefixOptions()
	opts.IgnoreCase = prefixIgnoreCase
	opts.StripExisting = prefixNoDouble

	result := parsed
	if prefixStrip != "" {
		result = env.StripPrefix(result, prefixStrip, opts)
	}
	if prefixAdd != "" {
		result = env.AddPrefix(result, prefixAdd, opts)
	}

	for _, k := range env.Keys(result) {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, result[k])
	}
	return nil
}
