package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/merge"
	"github.com/yourorg/envoy-diff/internal/parser"
)

var mergeStrategy string

func init() {
	mergeCmd := &cobra.Command{
		Use:   "merge <file1> <file2> [fileN...]",
		Short: "Merge multiple env files into one",
		Long: `Merge two or more .env files using a configurable conflict strategy.

Strategies:
  last   - last file to define a key wins (default)
  first  - first definition is kept
  error  - any duplicate key causes a non-zero exit`,
		Args:    cobra.MinimumNArgs(2),
		RunE:    runMerge,
		SilenceUsage: true,
	}

	mergeCmd.Flags().StringVarP(&mergeStrategy, "strategy", "s", "last",
		"conflict resolution strategy: last, first, error")

	rootCmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, args []string) error {
	var strat merge.Strategy
	switch strings.ToLower(mergeStrategy) {
	case "last":
		strat = merge.StrategyLast
	case "first":
		strat = merge.StrategyFirst
	case "error":
		strat = merge.StrategyError
	default:
		return fmt.Errorf("unknown strategy %q: choose last, first, or error", mergeStrategy)
	}

	sources := make([]map[string]string, 0, len(args))
	for _, path := range args {
		env, err := parser.ParseEnvFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}
		sources = append(sources, env)
	}

	res, err := merge.Merge(strat, sources...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "merge error: %v\n", err)
		return err
	}

	keys := make([]string, 0, len(res.Env))
	for k := range res.Env {
		keys = append(keys, k)
	}
	// stable output
	sortStrings(keys)
	for _, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, res.Env[k])
	}

	if len(res.Conflicts) > 0 {
		fmt.Fprintf(os.Stderr, "conflicts: %s\n", strings.Join(res.Conflicts, ", "))
	}
	return nil
}

func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
