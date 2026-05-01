package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sdfnj/ptdgen/internal/model"
	"gopkg.in/yaml.v3"
)

// Parse reads and parses a YAML or JSON target definition file, applying defaults.
func Parse(path string) (*model.TargetFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %q: %w", path, err)
	}

	var tf model.TargetFile
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &tf); err != nil {
			return nil, fmt.Errorf("invalid YAML in %q: %w", path, err)
		}
	case ".json":
		if err := json.Unmarshal(data, &tf); err != nil {
			return nil, fmt.Errorf("invalid JSON in %q: %w", path, err)
		}
	default:
		return nil, fmt.Errorf("unsupported file extension %q (use .yaml, .yml, or .json)", ext)
	}

	for i := range tf.Targets {
		tf.Targets[i].ApplyDefaults()
	}

	return &tf, nil
}
