package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-diff/internal/audit"
	"github.com/yourorg/envoy-diff/internal/diff"
	"github.com/yourorg/envoy-diff/internal/formatter"
	"github.com/yourorg/envoy-diff/internal/parser"
)

var (
	format     string
	auditMode  bool
	noColor    bool
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "envoy-diff [old-env-file] [new-env-file]",
	Short: "Diff and audit environment variable changes across deployment configs",
	Long: `envoy-diff is a CLI tool that compares two environment configuration files,
identifies changes, and provides security audit recommendations.`,
	Args: cobra.ExactArgs(2),
	RunE: runDiff,
}

func init() {
	rootCmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
	rootCmd.Flags().BoolVarP(&auditMode, "audit", "a", false, "Enable security audit mode")
	rootCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}

func runDiff(cmd *cobra.Command, args []string) error {
	oldFile, newFile := args[0], args[1]

	// Parse environment files
	oldEnv, err := parser.ParseEnvFile(oldFile)
	if err != nil {
		return fmt.Errorf("failed to parse old env file: %w", err)
	}

	newEnv, err := parser.ParseEnvFile(newFile)
	if err != nil {
		return fmt.Errorf("failed to parse new env file: %w", err)
	}

	// Compare environments
	changes := diff.Compare(oldEnv, newEnv)

	// Create formatter
	fmt, err := formatter.New(format, noColor)
	if err != nil {
		return fmt.Errorf("invalid format: %w", err)
	}

	// Output results
	output, err := fmt.Format(changes)
	if err != nil {
		return fmt.Errorf("formatting failed: %w", err)
	}
	fmt.Println(output)

	// Run audit if enabled
	if auditMode {
		findings := audit.Audit(changes)
		if len(findings) > 0 {
			fmt.Println("\n=== Security Audit Findings ===")
			for _, finding := range findings {
				fmt.Printf("[%s] %s: %s\n", finding.Severity, finding.Variable, finding.Message)
			}
			return fmt.Errorf("audit found %d security concern(s)", len(findings))
		}
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
