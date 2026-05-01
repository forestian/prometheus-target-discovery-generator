package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/sdfnj/ptdgen/internal/parser"
	"github.com/sdfnj/ptdgen/internal/validate"
	"github.com/spf13/cobra"
)

var validateFile string

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a target definition file",
	Long:  `Parse and validate a YAML or JSON target definition file, reporting any errors.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if validateFile == "" {
			return fmt.Errorf("--file is required")
		}

		tf, err := parser.Parse(validateFile)
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

		fmt.Printf("OK: %s is valid (%d target(s))\n", validateFile, len(tf.Targets))
		names := make([]string, 0, len(tf.Targets))
		for _, t := range tf.Targets {
			names = append(names, t.Name)
		}
		fmt.Printf("    targets: %s\n", strings.Join(names, ", "))
		return nil
	},
}

func init() {
	validateCmd.Flags().StringVar(&validateFile, "file", "", "Path to target definition file (.yaml or .json)")
}
