package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/sdfnj/ptdgen/internal/generator"
	"github.com/sdfnj/ptdgen/internal/parser"
	"github.com/sdfnj/ptdgen/internal/validate"
	"github.com/spf13/cobra"
)

var (
	generateFile   string
	generateOutput string
	generateFormat string
	generateForce  bool
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Prometheus / Alloy configuration from a target definition file",
	Long: `Parse a YAML or JSON target definition file, validate it, and write
output files for Prometheus file_sd_configs and/or Grafana Alloy discovery.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if generateFile == "" {
			return fmt.Errorf("--file is required")
		}

		format := generator.Format(strings.ToLower(generateFormat))
		switch format {
		case generator.FormatPrometheus, generator.FormatAlloy, generator.FormatAll:
		default:
			return fmt.Errorf("invalid --format %q (use prometheus, alloy, or all)", generateFormat)
		}

		tf, err := parser.Parse(generateFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		errs := validate.Validate(tf)
		if len(errs) > 0 {
			fmt.Fprintf(os.Stderr, "Validation failed (%d error(s)):\n", len(errs))
			for _, e := range errs {
				fmt.Fprintf(os.Stderr, "  - %s\n", e)
			}
			os.Exit(1)
		}

		fmt.Printf("Generating %s output into %s ...\n", format, generateOutput)
		if err := generator.Generate(tf, generateOutput, generator.Options{
			Format: format,
			Force:  generateForce,
		}); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Done. %d target(s) processed.\n", len(tf.Targets))
		return nil
	},
}

func init() {
	generateCmd.Flags().StringVar(&generateFile, "file", "", "Path to target definition file (.yaml or .json)")
	generateCmd.Flags().StringVar(&generateOutput, "output", "./generated", "Output directory")
	generateCmd.Flags().StringVar(&generateFormat, "format", "all", "Output format: prometheus, alloy, or all")
	generateCmd.Flags().BoolVar(&generateForce, "force", false, "Overwrite existing output files")
}
