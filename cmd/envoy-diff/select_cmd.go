package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
)

func init() {
	var keys []string
	var prefixes []string
	var caseInsensitive bool
	var invert bool

	cmd := &cobra.Command{
		Use:   "select <file>",
		Short: "Select a subset of keys from an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSelect(args[0], keys, prefixes, caseInsensitive, invert)
		},
	}

	cmd.Flags().StringSliceVarP(&keys, "key", "k", nil, "Exact key(s) to select (comma-separated)")
	cmd.Flags().StringSliceVarP(&prefixes, "prefix", "p", nil, "Key prefix(es) to select (comma-separated)")
	cmd.Flags().BoolVar(&caseInsensitive, "case-insensitive", false, "Match keys case-insensitively")
	cmd.Flags().BoolVar(&invert, "invert", false, "Return keys NOT matching the criteria")

	rootCmd.AddCommand(cmd)
}

func runSelect(file string, keys, prefixes []string, caseInsensitive, invert bool) error {
	sources, err := env.LoadSources([]string{file})
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}
	flat := env.Flatten(sources)

	opts := env.DefaultSelectOptions()
	opts.Keys = keys
	opts.Prefixes = prefixes
	opts.CaseSensitive = !caseInsensitive
	opts.Invert = invert

	out := env.SelectMap(flat, opts)

	for _, k := range env.Keys(out) {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, strings.ReplaceAll(out[k], "\n", "\\n"))
	}
	return nil
}
