package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sdfnj/ptdgen/internal/parser"
)

const sampleYAML = `
targets:
  - name: api-server-01
    job: app-api
    address: 10.10.10.11
    port: 9100
    scheme: http
    metrics_path: /metrics
    environment: production
    team: platform
    labels:
      service: api
`

const sampleJSON = `{
  "targets": [
    {
      "name": "api-server-01",
      "job": "app-api",
      "address": "10.10.10.11",
      "port": 9100,
      "scheme": "http",
      "metrics_path": "/metrics",
      "environment": "production",
      "team": "platform",
      "labels": {"service": "api"}
    }
  ]
}`

const minimalYAML = `
targets:
  - name: minimal-target
    job: minimal-job
    address: 192.168.1.1
    port: 8080
`

func writeTempFile(t *testing.T, content, ext string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "ptdgen-*"+ext)
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("WriteString: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseYAML(t *testing.T) {
	path := writeTempFile(t, sampleYAML, ".yaml")
	tf, err := parser.Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tf.Targets) != 1 {
		t.Fatalf("want 1 target, got %d", len(tf.Targets))
	}
	tgt := tf.Targets[0]
	if tgt.Name != "api-server-01" {
		t.Errorf("want name=api-server-01, got %q", tgt.Name)
	}
	if tgt.Port != 9100 {
		t.Errorf("want port=9100, got %d", tgt.Port)
	}
}

func TestParseJSON(t *testing.T) {
	path := writeTempFile(t, sampleJSON, ".json")
	tf, err := parser.Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tf.Targets) != 1 {
		t.Fatalf("want 1 target, got %d", len(tf.Targets))
	}
	if tf.Targets[0].Job != "app-api" {
		t.Errorf("want job=app-api, got %q", tf.Targets[0].Job)
	}
}

func TestParseDefaults(t *testing.T) {
	path := writeTempFile(t, minimalYAML, ".yaml")
	tf, err := parser.Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tgt := tf.Targets[0]
	if tgt.Scheme != "http" {
		t.Errorf("default scheme: want http, got %q", tgt.Scheme)
	}
	if tgt.MetricsPath != "/metrics" {
		t.Errorf("default metrics_path: want /metrics, got %q", tgt.MetricsPath)
	}
	if tgt.Environment != "unknown" {
		t.Errorf("default environment: want unknown, got %q", tgt.Environment)
	}
	if tgt.Team != "unknown" {
		t.Errorf("default team: want unknown, got %q", tgt.Team)
	}
	if tgt.Labels == nil {
		t.Error("default labels: want empty map, got nil")
	}
}

func TestParseUnsupportedExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "targets.toml")
	os.WriteFile(path, []byte("foo=bar"), 0o644)
	_, err := parser.Parse(path)
	if err == nil {
		t.Fatal("want error for .toml extension, got nil")
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, err := parser.Parse("/nonexistent/path/targets.yaml")
	if err == nil {
		t.Fatal("want error for missing file, got nil")
	}
}

func TestParseInvalidYAML(t *testing.T) {
	// yaml.v3 errors on unclosed flow sequences.
	path := writeTempFile(t, "targets:\n  - port: [unclosed\n", ".yaml")
	_, err := parser.Parse(path)
	if err == nil {
		t.Fatal("want error for invalid YAML, got nil")
	}
}

func TestParseInvalidJSON(t *testing.T) {
	path := writeTempFile(t, "{bad json", ".json")
	_, err := parser.Parse(path)
	if err == nil {
		t.Fatal("want error for invalid JSON, got nil")
	}
}
