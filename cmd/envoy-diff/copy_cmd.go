package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
	"github.com/yourorg/envoy-diff/internal/parser"
)

var (
	copyDestPrefix    string
	copySourcePrefix  string
	copyOverwrite     bool
	copyErrorMissing  bool
	copyKeys          []string
)

func init() {
	copyCmd := &cobra.Command{
		Use:   "copy <src-file> <dst-file>",
		Short: "Copy keys from one env file into another",
		Args:  cobra.ExactArgs(2),
		RunE:  runCopy,
	}

	copyCmd.Flags().StringVar(&copySourcePrefix, "src-prefix", "", "only copy keys with this prefix")
	copyCmd.Flags().StringVar(&copyDestPrefix, "dest-prefix", "", "prepend prefix to copied keys")
	copyCmd.Flags().BoolVar(&copyOverwrite, "overwrite", false, "overwrite existing keys in dst")
	copyCmd.Flags().BoolVar(&copyErrorMissing, "error-missing", false, "error when a requested key is absent")
	copyCmd.Flags().StringSliceVar(&copyKeys, "keys", nil, "comma-separated list of keys to copy (default: all)")

	rootCmd.AddCommand(copyCmd)
}

func runCopy(cmd *cobra.Command, args []string) error {
	srcPath, dstPath := args[0], args[1]

	srcMap, err := parser.ParseEnvFile(srcPath)
	if err != nil {
		return fmt.Errorf("reading source %q: %w", srcPath, err)
	}

	dstMap, err := parser.ParseEnvFile(dstPath)
	if err != nil {
		// Treat a missing dst file as an empty map.
		if os.IsNotExist(err) {
			dstMap = map[string]string{}
		} else {
			return fmt.Errorf("reading destination %q: %w", dstPath, err)
		}
	}

	opts := env.DefaultCopyOptions()
	opts.SourcePrefix = copySourcePrefix
	opts.DestPrefix = copyDestPrefix
	opts.Overwrite = copyOverwrite
	opts.ErrorOnMissing = copyErrorMissing

	result, err := env.CopyKeys(srcMap, dstMap, copyKeys, opts)
	if err != nil {
		return err
	}

	// Write merged result back to dst file.
	f, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("writing destination %q: %w", dstPath, err)
	}
	defer f.Close()

	for k, v := range dstMap {
		if _, werr := fmt.Fprintf(f, "%s=%s\n", k, v); werr != nil {
			return werr
		}
	}

	fmt.Fprintf(cmd.OutOrStdout(), "copied %d key(s), skipped %d, missing %d\n",
		len(result.Copied), len(result.Skipped), len(result.Missing))

	if len(result.Copied) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "  copied: %s\n", strings.Join(result.Copied, ", "))
	}

	return nil
}
