package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdfnj/ptdgen/internal/generator"
	"github.com/sdfnj/ptdgen/internal/model"
	"github.com/sdfnj/ptdgen/internal/parser"
	"github.com/sdfnj/ptdgen/internal/templates"
	"github.com/sdfnj/ptdgen/internal/validate"
	"github.com/spf13/cobra"
)

var (
	initOutput string
	initForce  bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an example project directory with sample target definitions",
	Long: `Creates a demo directory with sample targets.yaml, targets.json,
pre-generated output files, and example Prometheus / Alloy configs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		outDir := initOutput

		// Refuse to overwrite unless --force.
		if info, err := os.Stat(outDir); err == nil && info.IsDir() {
			if !initForce {
				return fmt.Errorf("directory %q already exists — use --force to overwrite", outDir)
			}
		}

		dirs := []string{
			outDir,
			filepath.Join(outDir, "generated"),
			filepath.Join(outDir, "examples"),
		}
		for _, d := range dirs {
			if err := os.MkdirAll(d, 0o755); err != nil {
				return fmt.Errorf("cannot create directory %q: %w", d, err)
			}
		}

		write := func(path string, content []byte) error {
			if err := generator.WriteFilePublic(path, content, initForce); err != nil {
				return err
			}
			fmt.Printf("  created %s\n", path)
			return nil
		}

		// README.md
		if err := write(filepath.Join(outDir, "README.md"), []byte(templates.InitReadme)); err != nil {
			return err
		}

		// targets.yaml
		if err := write(filepath.Join(outDir, "targets.yaml"), []byte(templates.SampleTargetsYAML)); err != nil {
			return err
		}

		// targets.json
		if err := write(filepath.Join(outDir, "targets.json"), []byte(templates.SampleTargetsJSON)); err != nil {
			return err
		}

		// examples/prometheus.yml
		if err := write(filepath.Join(outDir, "examples", "prometheus.yml"), buildPrometheusYml()); err != nil {
			return err
		}

		// examples/alloy.river
		if err := write(filepath.Join(outDir, "examples", "alloy.river"), buildAlloyRiver()); err != nil {
			return err
		}

		// Parse the just-written targets.yaml and generate files.
		targetsFile := filepath.Join(outDir, "targets.yaml")
		tf, err := parser.Parse(targetsFile)
		if err != nil {
			return fmt.Errorf("internal: cannot parse sample targets: %w", err)
		}

		errs := validate.Validate(tf)
		if len(errs) > 0 {
			return fmt.Errorf("internal: sample targets are invalid: %v", errs)
		}

		generatedDir := filepath.Join(outDir, "generated")
		if err := generator.Generate(tf, generatedDir, generator.Options{
			Format: generator.FormatAll,
			Force:  initForce,
		}); err != nil {
			return err
		}

		fmt.Printf("\nInitialized project at %s\n", outDir)
		fmt.Println("\nNext steps:")
		fmt.Printf("  ptdgen validate --file %s\n", filepath.Join(outDir, "targets.yaml"))
		fmt.Printf("  ptdgen generate --file %s --output %s --format all --force\n",
			filepath.Join(outDir, "targets.yaml"),
			generatedDir)
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&initOutput, "output", "./target-discovery-demo", "Output directory to create")
	initCmd.Flags().BoolVar(&initForce, "force", false, "Overwrite existing directory and files")
}

// buildPrometheusYml returns the example prometheus.yml content.
func buildPrometheusYml() []byte {
	return []byte(`# Example prometheus.yml showing how to wire up file_sd targets.
# Copy generated/prometheus-file-sd.json to /etc/prometheus/file_sd/ (or adjust path).

global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: file-sd-generated-targets
    file_sd_configs:
      - files:
          - /etc/prometheus/file_sd/prometheus-file-sd.json
        refresh_interval: 30s
`)
}

// buildAlloyRiver returns the example alloy.river content.
func buildAlloyRiver() []byte {
	return []byte(`// Example Grafana Alloy river config for file_sd-based discovery.
// PLACEHOLDER: update paths and remote_write URL before using in production.

discovery.file "ptdgen_targets" {
  files            = ["/etc/alloy/file_sd/prometheus-file-sd.json"]
  refresh_interval = "30s"
}

prometheus.scrape "ptdgen_targets" {
  targets    = discovery.file.ptdgen_targets.targets
  forward_to = [prometheus.remote_write.mimir.receiver]
}

prometheus.remote_write "mimir" {
  endpoint {
    // PLACEHOLDER: replace with your Mimir or Prometheus remote_write URL.
    url = "http://mimir-nginx.monitoring.svc:80/api/v1/push"
  }
}
`)
}

// buildTargetsJSON converts a TargetFile to pretty-printed JSON bytes.
func buildTargetsJSON(tf *model.TargetFile) ([]byte, error) {
	return json.MarshalIndent(tf, "", "  ")
}
