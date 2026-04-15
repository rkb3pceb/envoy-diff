package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
	"github.com/yourorg/envoy-diff/internal/parser"
)

var pinErrorOnViolation bool

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin <pinned-file> <env-file>",
		Short: "Check that pinned keys have not changed in an env file",
		Args:  cobra.ExactArgs(2),
		RunE:  runPin,
	}
	pinCmd.Flags().BoolVar(&pinErrorOnViolation, "error", false, "exit non-zero if any pin is violated")
	rootCmd.AddCommand(pinCmd)
}

func runPin(cmd *cobra.Command, args []string) error {
	pinnedPath := args[0]
	envPath := args[1]

	pinned, err := parser.ParseEnvFile(pinnedPath)
	if err != nil {
		return fmt.Errorf("reading pinned file: %w", err)
	}

	envMap, err := parser.ParseEnvFile(envPath)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	opts := env.DefaultPinOptions()
	opts.ErrorOnViolation = pinErrorOnViolation

	result, err := env.CheckPins(pinned, envMap, opts)
	if err != nil {
		// Print violations before returning the error
		printPinViolations(result)
		return err
	}

	if !result.HasViolations() {
		fmt.Fprintln(cmd.OutOrStdout(), "✔ all pinned keys match")
		return nil
	}

	printPinViolations(result)
	return nil
}

func printPinViolations(result env.PinResult) {
	for _, v := range result.Violations {
		actual := v.Actual
		if actual == "" {
			actual = "<missing>"
		}
		fmt.Fprintf(os.Stderr, "  VIOLATION  %s: pinned=%q actual=%q\n", v.Key, v.Pinned, actual)
	}
}
