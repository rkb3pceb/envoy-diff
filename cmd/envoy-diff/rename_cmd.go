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
	renamePairs         []string
	renameErrorMissing  bool
	renameErrorConflict bool
	renameOutput        string
)

func init() {
	renameCmd := &cobra.Command{
		Use:   "rename <file>",
		Short: "Rename keys in an env file",
		Long: `Rename one or more keys in an env file using OLD=NEW pairs.

Example:
  envoy-diff rename .env --rule DB_HOST=DATABASE_HOST --rule APP_PORT=PORT
`,
		Args:    cobra.ExactArgs(1),
		RunE:    runRename,
		SilenceUsage: true,
	}

	renameCmd.Flags().StringArrayVar(&renamePairs, "rule", nil, "rename rule in OLD=NEW format (repeatable)")
	renameCmd.Flags().BoolVar(&renameErrorMissing, "error-missing", false, "error if a source key does not exist")
	renameCmd.Flags().BoolVar(&renameErrorConflict, "error-conflict", false, "error if the destination key already exists")
	renameCmd.Flags().StringVarP(&renameOutput, "output", "o", "text", "output format: text or env")

	rootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, args []string) error {
	src, err := parser.ParseEnvFile(args[0])
	if err != nil {
		return fmt.Errorf("rename: failed to parse %s: %w", args[0], err)
	}

	rules := make(map[string]string, len(renamePairs))
	for _, pair := range renamePairs {
		for i, ch := range pair {
			if ch == '=' {
				old, newKey := pair[:i], pair[i+1:]
				if old == "" || newKey == "" {
					return fmt.Errorf("rename: invalid rule %q — expected OLD=NEW", pair)
				}
				rules[old] = newKey
				break
			}
		}
	}

	opts := env.DefaultRenameOptions()
	opts.Rules = rules
	opts.ErrorOnMissing = renameErrorMissing
	opts.ErrorOnConflict = renameErrorConflict

	res, err := env.RenameMap(src, opts)
	if err != nil {
		return fmt.Errorf("rename: %w", err)
	}

	if renameOutput == "env" {
		keys := make([]string, 0, len(res.Map))
		for k := range res.Map {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, res.Map[k])
		}
		return nil
	}

	for _, a := range res.Applied {
		fmt.Fprintf(os.Stdout, "renamed: %s → %s\n", a.OldKey, a.NewKey)
	}
	for _, s := range res.Skipped {
		fmt.Fprintf(os.Stdout, "skipped: %s (not found)\n", s)
	}
	if len(res.Applied) == 0 && len(res.Skipped) == 0 {
		fmt.Fprintln(os.Stdout, "no rename rules matched")
	}
	return nil
}
