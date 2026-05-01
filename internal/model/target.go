package model

// TargetFile is the top-level structure for a target definition file.
type TargetFile struct {
	Targets []Target `json:"targets" yaml:"targets"`
}

// Target represents a single Prometheus scrape target.
type Target struct {
	Name        string            `json:"name"         yaml:"name"`
	Job         string            `json:"job"          yaml:"job"`
	Address     string            `json:"address"      yaml:"address"`
	Port        int               `json:"port"         yaml:"port"`
	Scheme      string            `json:"scheme"       yaml:"scheme"`
	MetricsPath string            `json:"metrics_path" yaml:"metrics_path"`
	Environment string            `json:"environment"  yaml:"environment"`
	Team        string            `json:"team"         yaml:"team"`
	Labels      map[string]string `json:"labels"       yaml:"labels"`
}

// ApplyDefaults fills in optional fields with sensible defaults.
func (t *Target) ApplyDefaults() {
	if t.Scheme == "" {
		t.Scheme = "http"
	}
	if t.MetricsPath == "" {
		t.MetricsPath = "/metrics"
	}
	if t.Environment == "" {
		t.Environment = "unknown"
	}
	if t.Team == "" {
		t.Team = "unknown"
	}
	if t.Labels == nil {
		t.Labels = map[string]string{}
	}
}
