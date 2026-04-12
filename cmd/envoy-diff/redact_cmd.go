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
	redactLevel   string
	redactExtra   []string
	redactDisable bool
	redactList    bool
)

func init() {
	redactCmd := &cobra.Command{
		Use:   "redact <envfile>",
		Short: "Print an env file with sensitive values masked",
		Args:  cobra.ExactArgs(1),
		RunE:  runRedact,
	}

	redactCmd.Flags().StringVar(&redactLevel, "level", "partial", "Masking level: full or partial")
	redactCmd.Flags().StringArrayVar(&redactExtra, "pattern", nil, "Extra sensitive key patterns (repeatable)")
	redactCmd.Flags().BoolVar(&redactDisable, "no-redact", false, "Disable redaction (pass values through)")
	redactCmd.Flags().BoolVar(&redactList, "list", false, "Only list sensitive key names, do not print values")

	rootCmd.AddCommand(redactCmd)
}

func runRedact(cmd *cobra.Command, args []string) error {
	path := args[0]

	m, err := parser.ParseEnvFile(path)
	if err != nil {
		return fmt.Errorf("parse %q: %w", path, err)
	}

	if redactList {
		keys := env.SensitiveKeys(m, redactExtra)
		sort.Strings(keys)
		if len(keys) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "no sensitive keys detected")
			return nil
		}
		for _, k := range keys {
			fmt.Fprintln(cmd.OutOrStdout(), k)
		}
		return nil
	}

	opts := env.RedactOptions{
		Enabled:       !redactDisable,
		Level:         redactLevel,
		ExtraPatterns: redactExtra,
	}

	redacted := env.RedactMap(m, opts)

	keys := make([]string, 0, len(redacted))
	for k := range redacted {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := cmd.OutOrStdout()
	for _, k := range keys {
		fmt.Fprintf(w, "%s=%s\n", k, redacted[k])
	}

	if !redactDisable {
		sensitive := env.SensitiveKeys(m, redactExtra)
		if len(sensitive) > 0 {
			fmt.Fprintf(os.Stderr, "redacted %d sensitive key(s)\n", len(sensitive))
		}
	}

	return nil
}
