package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envoy-diff/internal/parser"
	"github.com/user/envoy-diff/internal/policy"
	"github.com/user/envoy-diff/internal/promote"
	"github.com/user/envoy-diff/internal/redact"
	"github.com/user/envoy-diff/internal/reporter"
)

var (
	promoteFromName  string
	promoteToName    string
	promotePolicy    string
	promoteNoRedact  bool
)

var promoteCmd = &cobra.Command{
	Use:   "promote <from-file> <to-file>",
	Short: "Evaluate an environment promotion between two stages",
	Args:  cobra.ExactArgs(2),
	RunE:  runPromote,
}

func init() {
	promoteCmd.Flags().StringVar(&promoteFromName, "from", "staging", "source stage name")
	promoteCmd.Flags().StringVar(&promoteToName, "to", "production", "target stage name")
	promoteCmd.Flags().StringVar(&promotePolicy, "policy", "", "path to policy YAML file")
	promoteCmd.Flags().BoolVar(&promoteNoRedact, "no-redact", false, "disable value redaction")
	rootCmd.AddCommand(promoteCmd)
}

func runPromote(cmd *cobra.Command, args []string) error {
	fromEnv, err := parser.ParseEnvFile(args[0])
	if err != nil {
		return fmt.Errorf("reading from-file: %w", err)
	}
	toEnv, err := parser.ParseEnvFile(args[1])
	if err != nil {
		return fmt.Errorf("reading to-file: %w", err)
	}

	opts := promote.DefaultOptions()
	if promoteNoRedact {
		opts.RedactCtx = redact.NewContext(false)
	}
	if promotePolicy != "" {
		p, err := policy.Load(promotePolicy)
		if err != nil {
			return fmt.Errorf("loading policy: %w", err)
		}
		opts.Policy = p
	}

	result, err := promote.Evaluate(
		promote.Stage{Name: promoteFromName, Env: fromEnv},
		promote.Stage{Name: promoteToName, Env: toEnv},
		opts,
	)
	if err != nil {
		return err
	}

	rep := reporter.Build(result.Changes, result.Findings)
	if err := reporter.Write(rep, os.Stdout); err != nil {
		return err
	}

	if result.Blocked {
		fmt.Fprintf(os.Stderr, "\npromotion BLOCKED by policy (%d violation(s))\n", len(result.Violations))
		return fmt.Errorf("policy violations block promotion")
	}
	return nil
}
