package templates

// InitReadme is the README.md written into the output directory by ptdgen init.
const InitReadme = `# Prometheus Target Discovery — Generated Project

This directory was created by **ptdgen init**.

## What's inside

| Path | Description |
|------|-------------|
` + "`targets.yaml`" + ` | Primary target definition (edit this) |
` + "`targets.json`" + ` | Same targets in JSON format |
` + "`generated/prometheus-file-sd.json`" + ` | Prometheus file_sd_configs compatible output |
` + "`generated/alloy-discovery.alloy`" + ` | Grafana Alloy starter config |
` + "`generated/scrape-config-example.yaml`" + ` | Example Prometheus scrape_config |
` + "`examples/prometheus.yml`" + ` | Full example prometheus.yml |
` + "`examples/alloy.river`" + ` | Full example Alloy river config |

---

## How to edit targets

Open ` + "`targets.yaml`" + ` and add or modify target entries:

` + "```yaml" + `
targets:
  - name: my-server-01
    job: my-job
    address: 192.168.1.100
    port: 9100
    scheme: http
    metrics_path: /metrics
    environment: production
    team: platform
    labels:
      region: us-east-1
` + "```" + `

---

## How to validate targets

` + "```bash" + `
ptdgen validate --file ./targets.yaml
` + "```" + `

---

## How to regenerate output

` + "```bash" + `
# Regenerate all formats
ptdgen generate --file ./targets.yaml --output ./generated --format all --force

# Prometheus only
ptdgen generate --file ./targets.yaml --output ./generated --format prometheus --force

# Alloy only
ptdgen generate --file ./targets.yaml --output ./generated --format alloy --force
` + "```" + `

---

## Using generated output with Prometheus

1. Copy ` + "`generated/prometheus-file-sd.json`" + ` to your Prometheus host, e.g.:
   ` + "`/etc/prometheus/file_sd/prometheus-file-sd.json`" + `

2. Add to your ` + "`prometheus.yml`" + `:

` + "```yaml" + `
scrape_configs:
  - job_name: file-sd-generated-targets
    file_sd_configs:
      - files:
          - /etc/prometheus/file_sd/prometheus-file-sd.json
        refresh_interval: 30s
` + "```" + `

3. Reload Prometheus. It will pick up targets automatically.

---

## Using generated output with Grafana Alloy

1. Copy ` + "`generated/prometheus-file-sd.json`" + ` to your Alloy host, e.g.:
   ` + "`/etc/alloy/file_sd/prometheus-file-sd.json`" + `

2. Copy or reference ` + "`generated/alloy-discovery.alloy`" + ` in your Alloy config.

3. Update the ` + "`prometheus.remote_write`" + ` URL to point to your Mimir or Prometheus.

4. Reload Alloy.

---

## Limitations

- This tool generates static files. Targets do not update automatically.
- Re-run ` + "`ptdgen generate`" + ` after every change to ` + "`targets.yaml`" + `.
- No authentication or secret management — do not store credentials in label values.
- No Kubernetes or cloud provider integration in this version.

---

## Security notes

- Do **not** include passwords, tokens, or API keys in label values.
- The ` + "`remote_write`" + ` URL in the Alloy config is a placeholder — update before use.
- Keep generated files inside your trusted infrastructure; treat them as configuration, not secrets.
`

// SampleTargetsYAML is written to targets.yaml by ptdgen init.
const SampleTargetsYAML = `targets:
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
      region: kr
      role: backend

  - name: worker-01
    job: app-worker
    address: 10.10.10.12
    port: 9100
    scheme: http
    metrics_path: /metrics
    environment: production
    team: platform
    labels:
      service: worker
      region: kr
      role: worker

  - name: db-exporter-01
    job: app-database
    address: 10.10.10.20
    port: 9187
    scheme: http
    metrics_path: /metrics
    environment: production
    team: data
    labels:
      service: postgres
      region: kr
      role: primary
`

// SampleTargetsJSON is written to targets.json by ptdgen init.
const SampleTargetsJSON = `{
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
      "labels": {
        "service": "api",
        "region": "kr",
        "role": "backend"
      }
    },
    {
      "name": "worker-01",
      "job": "app-worker",
      "address": "10.10.10.12",
      "port": 9100,
      "scheme": "http",
      "metrics_path": "/metrics",
      "environment": "production",
      "team": "platform",
      "labels": {
        "service": "worker",
        "region": "kr",
        "role": "worker"
      }
    }
  ]
}
`
