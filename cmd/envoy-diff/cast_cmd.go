package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
)

var castCmd = &cobra.Command{
	Use:   "cast <file>",
	Short: "Validate or normalise env values to a target type",
	Args:  cobra.ExactArgs(1),
	RunE:  runCast,
}

var (
	castType      string
	castKeys      []string
	castSkipError bool
)

func init() {
	rootCmd.AddCommand(castCmd)
	castCmd.Flags().StringVarP(&castType, "type", "t", "string", "target type: string|int|float|bool")
	castCmd.Flags().StringSliceVarP(&castKeys, "key", "k", nil, "keys to cast (default: all)")
	castCmd.Flags().BoolVar(&castSkipError, "skip-errors", true, "skip keys that fail casting")
}

func runCast(cmd *cobra.Command, args []string) error {
	sources, err := env.LoadSources(args)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}
	src := env.Flatten(sources)

	opts := env.DefaultCastOptions()
	opts.TargetType = env.CastType(castType)
	opts.Keys = castKeys
	opts.SkipOnError = castSkipError

	_, results, err := env.CastMap(src, opts)
	if err != nil {
		return fmt.Errorf("cast: %w", err)
	}

	w := cmd.OutOrStdout()
	if len(results) == 0 {
		fmt.Fprintln(w, "no keys matched")
		return nil
	}

	fmt.Fprintf(w, "%-30s %-10s %s\n", "KEY", "STATUS", "VALUE")
	fmt.Fprintln(w, strings.Repeat("-", 60))
	for _, r := range results {
		status := "ok"
		val := r.Casted
		if r.Err != nil {
			status = "FAIL"
			val = r.Original
			fmt.Fprintf(os.Stderr, "warn: %s: %v\n", r.Key, r.Err)
		}
		fmt.Fprintf(w, "%-30s %-10s %s\n", r.Key, status, val)
	}

	if env.HasCastErrors(results) && !castSkipError {
		return fmt.Errorf("one or more keys failed casting")
	}
	return nil
}
