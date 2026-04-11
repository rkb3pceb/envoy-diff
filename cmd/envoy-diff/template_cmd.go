package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-diff/internal/template"
)

var (
	tmplVars        []string
	tmplAllowMissing bool
	tmplOutput      string
)

func init() {
	templateCmd := &cobra.Command{
		Use:   "template <file>",
		Short: "Render an env template file",
		Long: `Render a .env template that uses Go template syntax.

The template context is populated from the current OS environment.
Extra variables can be injected with --var KEY=VALUE flags.`,
		Args:    cobra.ExactArgs(1),
		RunE:    runTemplate,
		Example: "  envoy-diff template deploy/prod.env.tmpl --var REGION=us-east-1",
	}

	templateCmd.Flags().StringArrayVar(&tmplVars, "var", nil, "Extra variables as KEY=VALUE (repeatable)")
	templateCmd.Flags().BoolVar(&tmplAllowMissing, "allow-missing", false, "Treat undefined template vars as empty string")
	templateCmd.Flags().StringVarP(&tmplOutput, "output", "o", "", "Write rendered output to file instead of stdout")

	rootCmd.AddCommand(templateCmd)
}

func runTemplate(cmd *cobra.Command, args []string) error {
	opts := template.DefaultOptions()
	opts.AllowMissing = tmplAllowMissing

	for _, kv := range tmplVars {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --var %q: expected KEY=VALUE", kv)
		}
		opts.Vars[parts[0]] = parts[1]
	}

	out, err := template.RenderFile(args[0], opts)
	if err != nil {
		return err
	}

	if tmplOutput != "" {
		if err := os.WriteFile(tmplOutput, out, 0o644); err != nil {
			return fmt.Errorf("write output file: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "rendered template written to %s\n", tmplOutput)
		return nil
	}

	_, err = cmd.OutOrStdout().Write(out)
	return err
}
