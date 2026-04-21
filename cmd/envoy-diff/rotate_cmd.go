package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
	"github.com/yourorg/envoy-diff/internal/parser"
)

var rotateCmd = &cobra.Command{
	Use:   "rotate <file>",
	Short: "Rename environment variable keys according to rotation rules",
	Args:  cobra.ExactArgs(1),
	RunE:  runRotate,
}

func init() {
	rotateCmd.Flags().StringSliceP("rule", "r", nil, "Rotation rules as OLD=NEW pairs (repeatable)")
	rotateCmd.Flags().Bool("keep-old", false, "Retain the original key alongside the new key")
	rotateCmd.Flags().Bool("error-on-missing", false, "Return an error if a source key is absent")
	rotateCmd.Flags().Bool("error-on-conflict", false, "Return an error if the destination key already exists")
	rootCmd.AddCommand(rotateCmd)
}

func runRotate(cmd *cobra.Command, args []string) error {
	rawRules, _ := cmd.Flags().GetStringSlice("rule")
	keepOld, _ := cmd.Flags().GetBool("keep-old")
	errorOnMissing, _ := cmd.Flags().GetBool("error-on-missing")
	errorOnConflict, _ := cmd.Flags().GetBool("error-on-conflict")

	if len(rawRules) == 0 {
		return fmt.Errorf("at least one --rule OLD=NEW is required")
	}

	rules := make(map[string]string, len(rawRules))
	for _, r := range rawRules {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid rule %q: expected OLD=NEW", r)
		}
		rules[parts[0]] = parts[1]
	}

	src, err := parser.ParseEnvFile(args[0])
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[0], err)
	}

	opts := env.DefaultRotateOptions()
	opts.Rules = rules
	opts.KeepOld = keepOld
	opts.ErrorOnMissing = errorOnMissing
	opts.ErrorOnConflict = errorOnConflict

	out, res, err := env.RotateMap(src, opts)
	if err != nil {
		return err
	}

	for _, k := range env.SortedKeys(out, env.DefaultSortOptions()) {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, out[k])
	}

	if len(res.Skipped) > 0 {
		fmt.Fprintf(os.Stderr, "skipped keys: %s\n", strings.Join(res.Skipped, ", "))
	}
	return nil
}
