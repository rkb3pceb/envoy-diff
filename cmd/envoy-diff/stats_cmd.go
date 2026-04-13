package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/env"
	"github.com/yourorg/envoy-diff/internal/parser"
)

var statsTopN int

func init() {
	statsCmd := &cobra.Command{
		Use:   "stats <file>",
		Short: "Print aggregate statistics for an env file",
		Args:  cobra.ExactArgs(1),
		RunE:  runStats,
	}
	statsCmd.Flags().IntVar(&statsTopN, "top", 5, "number of top prefixes to display")
	rootCmd.AddCommand(statsCmd)
}

func runStats(cmd *cobra.Command, args []string) error {
	path := args[0]

	envMap, err := parser.ParseEnvFile(path)
	if err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	opts := env.StatsOptions{TopN: statsTopN}
	r := env.Stats(envMap, opts)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "File:\t%s\n", path)
	fmt.Fprintf(w, "Total keys:\t%d\n", r.Total)
	fmt.Fprintf(w, "Empty values:\t%d\n", r.Empty)
	fmt.Fprintf(w, "Sensitive keys:\t%d\n", r.Sensitive)
	fmt.Fprintf(w, "Unique values:\t%d\n", r.Unique)
	w.Flush()

	if len(r.TopPrefixes) > 0 {
		fmt.Println()
		fmt.Println("Top prefixes:")
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for _, p := range r.TopPrefixes {
			fmt.Fprintf(tw, "  %s\t%d keys\n", p.Prefix, p.Count)
		}
		tw.Flush()
	}

	return nil
}
