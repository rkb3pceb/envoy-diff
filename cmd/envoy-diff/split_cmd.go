package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
	"github.com/yourorg/envoy-diff/internal/parser"
)

func init() {
	var (
		delimiter    string
		keyTemplate  string
		indexBase    int
		keys         []string
		skipEmpty    bool
		trimParts    bool
		keepOriginal bool
	)

	cmd := &cobra.Command{
		Use:   "split <file>",
		Short: "Split multi-value environment variables into indexed keys",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSplit(args[0], env.SplitOptions{
				Delimiter:    delimiter,
				KeyTemplate:  keyTemplate,
				IndexBase:    indexBase,
				Keys:         keys,
				SkipEmpty:    skipEmpty,
				TrimParts:    trimParts,
				KeepOriginal: keepOriginal,
			})
		},
	}

	defaults := env.DefaultSplitOptions()
	cmd.Flags().StringVarP(&delimiter, "delimiter", "d", defaults.Delimiter, "Value delimiter")
	cmd.Flags().StringVar(&keyTemplate, "key-template", defaults.KeyTemplate, "Key template ({KEY} and {INDEX} placeholders)")
	cmd.Flags().IntVar(&indexBase, "index-base", defaults.IndexBase, "Starting index for generated keys")
	cmd.Flags().StringSliceVarP(&keys, "key", "k", nil, "Keys to split (default: all)")
	cmd.Flags().BoolVar(&skipEmpty, "skip-empty", defaults.SkipEmpty, "Skip empty parts after splitting")
	cmd.Flags().BoolVar(&trimParts, "trim", defaults.TrimParts, "Trim whitespace from each part")
	cmd.Flags().BoolVar(&keepOriginal, "keep-original", false, "Keep original key alongside split keys")

	rootCmd.AddCommand(cmd)
}

func runSplit(file string, opts env.SplitOptions) error {
	raw, err := parser.ParseEnvFile(file)
	if err != nil {
		return fmt.Errorf("split: parse %q: %w", file, err)
	}

	result, err := env.SplitMap(raw, opts)
	if err != nil {
		return fmt.Errorf("split: %w", err)
	}

	keys := make([]string, 0, len(result.Map))
	for k := range result.Map {
		keys = append(keys, k)
	}
	sortStrings(keys)

	w := os.Stdout
	for _, k := range keys {
		v := result.Map[k]
		if strings.ContainsAny(v, " \t\n") {
			fmt.Fprintf(w, "%s=%q\n", k, v)
		} else {
			fmt.Fprintf(w, "%s=%s\n", k, v)
		}
	}

	if env.HasSplitChanges(result) {
		fmt.Fprintf(os.Stderr, "split: expanded %d key(s): %s\n",
			len(result.SplitKeys), strings.Join(result.SplitKeys, ", "))
	}
	return nil
}
