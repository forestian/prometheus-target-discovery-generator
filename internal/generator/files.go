package generator

const scrapeConfigExampleTemplate = `# Example Prometheus scrape_config that consumes the file_sd output from ptdgen.
# Mount generated/prometheus-file-sd.json into your Prometheus container and
# update the path below.

scrape_configs:
  - job_name: file-sd-generated-targets
    file_sd_configs:
      - files:
          - /etc/prometheus/file_sd/prometheus-file-sd.json
        refresh_interval: 30s
`

const prometheusYmlTemplate = `# Example prometheus.yml showing how to wire up file_sd targets.
# Copy prometheus-file-sd.json to /etc/prometheus/file_sd/ (or adjust path).

global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: file-sd-generated-targets
    file_sd_configs:
      - files:
          - /etc/prometheus/file_sd/prometheus-file-sd.json
        refresh_interval: 30s
`

const alloyRiverExampleTemplate = `// Example alloy.river / alloy config for file_sd-based discovery.
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
`

func buildScrapeConfigExample() []byte {
	return []byte(scrapeConfigExampleTemplate)
}

func buildPrometheusYml() []byte {
	return []byte(prometheusYmlTemplate)
}

func buildAlloyRiverExample() []byte {
	return []byte(alloyRiverExampleTemplate)
}
