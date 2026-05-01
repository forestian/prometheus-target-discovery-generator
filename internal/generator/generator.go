package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdfnj/ptdgen/internal/model"
)

// Format represents the output format selection.
type Format string

const (
	FormatPrometheus Format = "prometheus"
	FormatAlloy      Format = "alloy"
	FormatAll        Format = "all"
)

// Options controls generator behaviour.
type Options struct {
	Format Format
	Force  bool
}

// Generate writes the selected output files into outDir.
func Generate(tf *model.TargetFile, outDir string, opts Options) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("cannot create output directory %q: %w", outDir, err)
	}

	type fileSpec struct {
		name    string
		content []byte
	}

	var files []fileSpec

	switch opts.Format {
	case FormatPrometheus:
		data, err := buildPrometheusFileSd(tf.Targets)
		if err != nil {
			return fmt.Errorf("building prometheus file_sd: %w", err)
		}
		files = append(files, fileSpec{"prometheus-file-sd.json", data})
		files = append(files, fileSpec{"scrape-config-example.yaml", buildScrapeConfigExample()})

	case FormatAlloy:
		files = append(files, fileSpec{"alloy-discovery.alloy", buildAlloyDiscovery()})

	case FormatAll:
		data, err := buildPrometheusFileSd(tf.Targets)
		if err != nil {
			return fmt.Errorf("building prometheus file_sd: %w", err)
		}
		files = append(files,
			fileSpec{"prometheus-file-sd.json", data},
			fileSpec{"alloy-discovery.alloy", buildAlloyDiscovery()},
			fileSpec{"scrape-config-example.yaml", buildScrapeConfigExample()},
		)

	default:
		return fmt.Errorf("unknown format %q (use prometheus, alloy, or all)", opts.Format)
	}

	for _, f := range files {
		dest := filepath.Join(outDir, f.name)
		if err := WriteFilePublic(dest, f.content, opts.Force); err != nil {
			return err
		}
		fmt.Printf("  wrote %s\n", dest)
	}

	return nil
}

// WriteFilePublic writes content to path. Fails if file exists and force is false.
func WriteFilePublic(path string, content []byte, force bool) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("file already exists: %q — use --force to overwrite", path)
		}
	}
	return os.WriteFile(path, content, 0o644)
}
