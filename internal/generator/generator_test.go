package generator_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sdfnj/ptdgen/internal/generator"
	"github.com/sdfnj/ptdgen/internal/model"
)

func sampleTargetFile() *model.TargetFile {
	t1 := model.Target{
		Name:        "api-server-01",
		Job:         "app-api",
		Address:     "10.10.10.11",
		Port:        9100,
		Environment: "production",
		Team:        "platform",
		Labels:      map[string]string{"service": "api", "region": "kr"},
	}
	t1.ApplyDefaults()
	return &model.TargetFile{Targets: []model.Target{t1}}
}

func TestGeneratePrometheus(t *testing.T) {
	outDir := t.TempDir()
	tf := sampleTargetFile()

	err := generator.Generate(tf, outDir, generator.Options{
		Format: generator.FormatPrometheus,
		Force:  false,
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "prometheus-file-sd.json"))
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	var entries []struct {
		Targets []string          `json:"targets"`
		Labels  map[string]string `json:"labels"`
	}
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}

	entry := entries[0]
	if len(entry.Targets) != 1 || entry.Targets[0] != "10.10.10.11:9100" {
		t.Errorf("want target=10.10.10.11:9100, got %v", entry.Targets)
	}
	if entry.Labels["job"] != "app-api" {
		t.Errorf("want job=app-api, got %q", entry.Labels["job"])
	}
	if entry.Labels["instance"] != "api-server-01" {
		t.Errorf("want instance=api-server-01, got %q", entry.Labels["instance"])
	}
	if entry.Labels["environment"] != "production" {
		t.Errorf("want environment=production, got %q", entry.Labels["environment"])
	}
	if entry.Labels["team"] != "platform" {
		t.Errorf("want team=platform, got %q", entry.Labels["team"])
	}
	if entry.Labels["__scheme__"] != "http" {
		t.Errorf("want __scheme__=http, got %q", entry.Labels["__scheme__"])
	}
	if entry.Labels["__metrics_path__"] != "/metrics" {
		t.Errorf("want __metrics_path__=/metrics, got %q", entry.Labels["__metrics_path__"])
	}
	if entry.Labels["service"] != "api" {
		t.Errorf("want service=api, got %q", entry.Labels["service"])
	}
	if entry.Labels["region"] != "kr" {
		t.Errorf("want region=kr, got %q", entry.Labels["region"])
	}
}

func TestGenerateAlloy(t *testing.T) {
	outDir := t.TempDir()
	tf := sampleTargetFile()

	err := generator.Generate(tf, outDir, generator.Options{
		Format: generator.FormatAlloy,
		Force:  false,
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "alloy-discovery.alloy"))
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "discovery.file") {
		t.Error("alloy config should contain discovery.file")
	}
	if !strings.Contains(content, "prometheus.scrape") {
		t.Error("alloy config should contain prometheus.scrape")
	}
	if !strings.Contains(content, "prometheus.remote_write") {
		t.Error("alloy config should contain prometheus.remote_write")
	}
}

func TestGenerateAll(t *testing.T) {
	outDir := t.TempDir()
	tf := sampleTargetFile()

	err := generator.Generate(tf, outDir, generator.Options{
		Format: generator.FormatAll,
		Force:  false,
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	for _, name := range []string{
		"prometheus-file-sd.json",
		"alloy-discovery.alloy",
		"scrape-config-example.yaml",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Errorf("expected file %q to exist: %v", name, err)
		}
	}
}

func TestGenerateOverwriteProtection(t *testing.T) {
	outDir := t.TempDir()
	tf := sampleTargetFile()

	opts := generator.Options{Format: generator.FormatPrometheus, Force: false}

	// First run: should succeed.
	if err := generator.Generate(tf, outDir, opts); err != nil {
		t.Fatalf("first Generate: %v", err)
	}

	// Second run without --force: should fail.
	if err := generator.Generate(tf, outDir, opts); err == nil {
		t.Error("want error on overwrite without --force")
	}
}

func TestGenerateOverwriteWithForce(t *testing.T) {
	outDir := t.TempDir()
	tf := sampleTargetFile()

	opts := generator.Options{Format: generator.FormatPrometheus, Force: true}

	if err := generator.Generate(tf, outDir, opts); err != nil {
		t.Fatalf("first Generate: %v", err)
	}
	if err := generator.Generate(tf, outDir, opts); err != nil {
		t.Fatalf("second Generate with force: %v", err)
	}
}

func TestGenerateInvalidFormat(t *testing.T) {
	outDir := t.TempDir()
	tf := sampleTargetFile()

	err := generator.Generate(tf, outDir, generator.Options{Format: "invalid"})
	if err == nil {
		t.Error("want error for invalid format")
	}
}

func TestGenerateDeterministic(t *testing.T) {
	// Two runs with the same input should produce identical output.
	tf := &model.TargetFile{
		Targets: func() []model.Target {
			targets := []model.Target{
				{Name: "z-last", Job: "job1", Address: "10.0.0.3", Port: 9100},
				{Name: "a-first", Job: "job1", Address: "10.0.0.1", Port: 9100},
				{Name: "m-middle", Job: "job1", Address: "10.0.0.2", Port: 9100},
			}
			for i := range targets {
				targets[i].ApplyDefaults()
			}
			return targets
		}(),
	}

	out1 := t.TempDir()
	out2 := t.TempDir()
	opts := generator.Options{Format: generator.FormatPrometheus}

	if err := generator.Generate(tf, out1, opts); err != nil {
		t.Fatalf("run1: %v", err)
	}
	if err := generator.Generate(tf, out2, opts); err != nil {
		t.Fatalf("run2: %v", err)
	}

	d1, _ := os.ReadFile(filepath.Join(out1, "prometheus-file-sd.json"))
	d2, _ := os.ReadFile(filepath.Join(out2, "prometheus-file-sd.json"))
	if string(d1) != string(d2) {
		t.Error("output is not deterministic across two runs")
	}

	// Also verify sorted order: a-first < m-middle < z-last.
	var entries []struct {
		Labels map[string]string `json:"labels"`
	}
	if err := json.Unmarshal(d1, &entries); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	names := []string{entries[0].Labels["instance"], entries[1].Labels["instance"], entries[2].Labels["instance"]}
	if names[0] != "a-first" || names[1] != "m-middle" || names[2] != "z-last" {
		t.Errorf("want sorted order [a-first m-middle z-last], got %v", names)
	}
}
