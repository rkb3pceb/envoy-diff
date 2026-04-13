package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
	"github.com/yourorg/envoy-diff/internal/parser"
)

var sanitizeCmd = &cobra.Command{
	Use:   "sanitize <file>",
	Short: "Sanitize env file keys and values",
	Args:  cobra.ExactArgs(1),
	RunE:  runSanitize,
}

func init() {
	sanitizeCmd.Flags().Bool("trim", true, "trim whitespace from values")
	sanitizeCmd.Flags().Bool("normalize-keys", false, "uppercase all keys")
	sanitizeCmd.Flags().Bool("remove-empty", false, "drop keys with empty values")
	sanitizeCmd.Flags().Bool("fix-keys", false, "replace invalid key characters with underscore")
	rootCmd.AddCommand(sanitizeCmd)
}

func runSanitize(cmd *cobra.Command, args []string) error {
	path := args[0]

	parsed, err := parser.ParseEnvFile(path)
	if err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	trim, _ := cmd.Flags().GetBool("trim")
	normalize, _ := cmd.Flags().GetBool("normalize-keys")
	removeEmpty, _ := cmd.Flags().GetBool("remove-empty")
	fixKeys, _ := cmd.Flags().GetBool("fix-keys")

	opts := env.SanitizeOptions{
		TrimSpace:           trim,
		NormalizeKeys:       normalize,
		RemoveEmpty:         removeEmpty,
		ReplaceInvalidChars: fixKeys,
	}

	result := env.SanitizeMap(parsed, opts)

	if env.HasSanitizeChanges(result) {
		if len(result.Renamed) > 0 {
			fmt.Fprintln(os.Stderr, "# renamed keys:")
			for old, newKey := range result.Renamed {
				fmt.Fprintf(os.Stderr, "#   %s -> %s\n", old, newKey)
			}
		}
		if len(result.Dropped) > 0 {
			fmt.Fprintln(os.Stderr, "# dropped keys:")
			for _, k := range result.Dropped {
				fmt.Fprintf(os.Stderr, "#   %s\n", k)
			}
		}
	}

	keys := make([]string, 0, len(result.Map))
	for k := range result.Map {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%s=%s\n", k, result.Map[k])
	}
	return nil
}
