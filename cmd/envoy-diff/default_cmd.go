package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/envoy-diff/internal/env"
)

var (
	defaultOverwriteEmpty bool
	defaultOverwriteAll   bool
	defaultPairs          []string
)

func init() {
	defaultCmd := &cobra.Command{
		Use:   "default <env-file>",
		Short: "Apply default values to missing or empty keys in an env file",
		Args:  cobra.ExactArgs(1),
		RunE:  runDefault,
	}

	defaultCmd.Flags().StringArrayVarP(&defaultPairs, "set", "s", nil,
		"Default key=value pair (repeatable)")
	defaultCmd.Flags().BoolVar(&defaultOverwriteEmpty, "overwrite-empty", false,
		"Replace keys that exist but have an empty value")
	defaultCmd.Flags().BoolVar(&defaultOverwriteAll, "overwrite-all", false,
		"Replace all keys unconditionally")

	rootCmd.AddCommand(defaultCmd)
}

func runDefault(cmd *cobra.Command, args []string) error {
	sources, err := env.LoadSources(args)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}
	src := env.Flatten(sources)

	defaults := make(map[string]string, len(defaultPairs))
	for _, pair := range defaultPairs {
		k, v, ok := splitKeyValue(pair)
		if !ok {
			return fmt.Errorf("invalid --set value %q: expected key=value", pair)
		}
		defaults[k] = v
	}

	opts := env.DefaultDefaultOptions()
	opts.Defaults = defaults
	opts.OverwriteEmpty = defaultOverwriteEmpty
	opts.OverwriteAll = defaultOverwriteAll

	out, res := env.ApplyDefaults(src, opts)

	for _, k := range env.Keys([]map[string]string{out}) {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, out[k])
	}

	if len(res.Applied) > 0 {
		fmt.Fprintf(os.Stderr, "applied defaults: %v\n", res.Applied)
	}
	return nil
}

// splitKeyValue splits "KEY=VALUE" into (key, value, true).
// Returns ("", "", false) when no '=' is present.
func splitKeyValue(pair string) (string, string, bool) {
	for i, c := range pair {
		if c == '=' {
			return pair[:i], pair[i+1:], true
		}
	}
	return "", "", false
}
