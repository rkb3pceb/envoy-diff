package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
	"github.com/yourorg/envoy-diff/internal/parser"
)

var (
	placeholderErrorOnUnresolved bool
	placeholderSubstitutions     []string
)

func init() {
	placeholderCmd := &cobra.Command{
		Use:   "placeholder <file>",
		Short: "Find and substitute placeholder values in an env file",
		Args:  cobra.ExactArgs(1),
		RunE:  runPlaceholder,
	}
	placeholderCmd.Flags().BoolVar(&placeholderErrorOnUnresolved, "error-on-unresolved", false,
		"Exit with error if any placeholder remains unsubstituted")
	placeholderCmd.Flags().StringArrayVar(&placeholderSubstitutions, "set", nil,
		"Substitute a placeholder: KEY=VALUE (repeatable)")
	rootCmd.AddCommand(placeholderCmd)
}

func runPlaceholder(cmd *cobra.Command, args []string) error {
	envMap, err := parser.ParseEnvFile(args[0])
	if err != nil {
		return fmt.Errorf("parse %s: %w", args[0], err)
	}

	opts := env.DefaultPlaceholderOptions()
	opts.ErrorOnUnresolved = placeholderErrorOnUnresolved

	for _, s := range placeholderSubstitutions {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --set value %q: expected KEY=VALUE", s)
		}
		opts.Substitutions[parts[0]] = parts[1]
	}

	out, findings, err := env.SubstitutePlaceholders(envMap, opts)
	if err != nil {
		return err
	}

	if len(findings) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No placeholders found.")
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%-30s %-15s %-10s %s\n", "KEY", "MARKER", "RESOLVED", "VALUE")
	fmt.Fprintln(cmd.OutOrStdout(), strings.Repeat("-", 70))
	for _, f := range findings {
		resolvedStr := "no"
		displayVal := f.Value
		if f.Resolved {
			resolvedStr = "yes"
			displayVal = out[f.Key]
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%-30s %-15s %-10s %s\n",
			f.Key, f.Marker, resolvedStr, displayVal)
	}

	unresolved := 0
	for _, f := range findings {
		if !f.Resolved {
			unresolved++
		}
	}
	if unresolved > 0 {
		fmt.Fprintf(os.Stderr, "\n%d unresolved placeholder(s) remain.\n", unresolved)
	}
	return nil
}
