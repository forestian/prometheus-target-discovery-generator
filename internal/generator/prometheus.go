package generator

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/sdfnj/ptdgen/internal/model"
)

// prometheusFileSd is a single entry in Prometheus file_sd_configs JSON.
type prometheusFileSd struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

// buildPrometheusFileSd converts targets to Prometheus file_sd JSON bytes.
func buildPrometheusFileSd(targets []model.Target) ([]byte, error) {
	entries := make([]prometheusFileSd, 0, len(targets))

	for _, t := range targets {
		labels := map[string]string{
			"job":              t.Job,
			"instance":         t.Name,
			"environment":      t.Environment,
			"team":             t.Team,
			"__scheme__":       t.Scheme,
			"__metrics_path__": t.MetricsPath,
		}
		// Merge custom labels; custom labels do NOT override reserved keys.
		for k, v := range t.Labels {
			if _, reserved := labels[k]; !reserved {
				labels[k] = v
			}
		}

		entries = append(entries, prometheusFileSd{
			Targets: []string{fmt.Sprintf("%s:%d", t.Address, t.Port)},
			Labels:  labels,
		})
	}

	// Stable output: sort entries by instance name.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Labels["instance"] < entries[j].Labels["instance"]
	})

	return json.MarshalIndent(entries, "", "  ")
}
