package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
	"github.com/yourorg/envoy-diff/internal/parser"
)

var (
	expandFallbackOS   bool
	expandErrorMissing bool
	expandShowMissing  bool
)

func init() {
	expandCmd := &cobra.Command{
		Use:   "expand <file>",
		Short: "Expand variable references within an env file",
		Long: `Reads an env file and resolves internal variable references
(e.g. API_URL=${BASE_URL}/api). Optionally falls back to OS environment
variables and reports any unresolved references.`,
		Args: cobra.ExactArgs(1),
		RunE: runExpand,
	}

	expandCmd.Flags().BoolVar(&expandFallbackOS, "os-fallback", false,
		"fall back to OS environment variables for unresolved references")
	expandCmd.Flags().BoolVar(&expandErrorMissing, "error-missing", false,
		"return a non-zero exit code if any references remain unresolved")
	expandCmd.Flags().BoolVar(&expandShowMissing, "show-missing", false,
		"print unresolved variable references to stderr before output")

	rootCmd.AddCommand(expandCmd)
}

func runExpand(cmd *cobra.Command, args []string) error {
	path := args[0]

	raw, err := parser.ParseEnvFile(path)
	if err != nil {
		return fmt.Errorf("parse %q: %w", path, err)
	}

	if expandShowMissing {
		missing := env.MissingRefs(raw, expandFallbackOS)
		if len(missing) > 0 {
			fmt.Fprintf(os.Stderr, "unresolved references: %v\n", missing)
		}
	}

	opts := env.ExpandOptions{
		FallbackToOS:   expandFallbackOS,
		ErrorOnMissing: expandErrorMissing,
	}

	expanded, err := env.ExpandMap(raw, opts)
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(expanded))
	for k := range expanded {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, expanded[k])
	}

	return nil
}
