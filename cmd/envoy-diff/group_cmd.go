package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
)

var (
	groupDelimiter      string
	groupDepth          int
	groupExcludeOther   bool
	groupUngroupedLabel string
)

func init() {
	groupCmd := &cobra.Command{
		Use:   "group <file>",
		Short: "Group environment variables by key prefix",
		Args:  cobra.ExactArgs(1),
		RunE:  runGroup,
	}
	groupCmd.Flags().StringVar(&groupDelimiter, "delimiter", "_", "delimiter used to split key prefix")
	groupCmd.Flags().IntVar(&groupDepth, "depth", 0, "number of prefix segments to use as group label")
	groupCmd.Flags().BoolVar(&groupExcludeOther, "no-other", false, "exclude keys that have no delimiter")
	groupCmd.Flags().StringVar(&groupUngroupedLabel, "other-label", "OTHER", "label for ungrouped keys")
	rootCmd.AddCommand(groupCmd)
}

func runGroup(cmd *cobra.Command, args []string) error {
	m, err := loadEnvFile(args[0])
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[0], err)
	}

	opts := env.DefaultGroupOptions()
	opts.Delimiter = groupDelimiter
	opts.MaxDepth = groupDepth
	opts.IncludeUngrouped = !groupExcludeOther
	opts.UngroupedLabel = groupUngroupedLabel

	groups := env.GroupMap(m, opts)
	if len(groups) == 0 {
		fmt.Fprintln(os.Stdout, "(no groups found)")
		return nil
	}

	w := os.Stdout
	for _, g := range groups {
		fmt.Fprintf(w, "[%s] (%d keys)\n", g.Label, len(g.Keys))
		for _, k := range g.Keys {
			fmt.Fprintf(w, "  %s=%s\n", k, g.Vars[k])
		}
		fmt.Fprintln(w, strings.Repeat("-", 40))
	}
	return nil
}
