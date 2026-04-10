package main

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-diff/internal/parser"
	"github.com/yourorg/envoy-diff/internal/resolve"
)

var resolveOSFallback bool

var resolveCmd = &cobra.Command{
	Use:   "resolve <env-file>",
	Short: "Expand variable references within an env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runResolve,
}

func init() {
	resolveCmd.Flags().BoolVar(&resolveOSFallback, "os-fallback", false,
		"fall back to host OS environment for missing references")
	rootCmd.AddCommand(resolveCmd)
}

func runResolve(cmd *cobra.Command, args []string) error {
	envMap, err := parser.ParseEnvFile(args[0])
	if err != nil {
		return fmt.Errorf("parse %s: %w", args[0], err)
	}

	opts := resolve.Options{FallbackToOS: resolveOSFallback}
	resolved := resolve.Map(envMap, opts)

	unresolved := resolve.UnresolvedKeys(resolved)
	if len(unresolved) > 0 {
		sort.Strings(unresolved)
		fmt.Fprintln(cmd.ErrOrStderr(), "warning: unresolved references in keys:")
		for _, k := range unresolved {
			fmt.Fprintf(cmd.ErrOrStderr(), "  %s=%s\n", k, resolved[k])
		}
	}

	keys := make([]string, 0, len(resolved))
	for k := range resolved {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, resolved[k])
	}
	return nil
}
