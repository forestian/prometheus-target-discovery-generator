# ptdgen — Prometheus Target Discovery Generator

`ptdgen` is a local CLI tool that generates Prometheus `file_sd_configs`-compatible
scrape target configuration and Grafana Alloy `discovery.file` starter configs
from a simple YAML or JSON target definition file.

It is the MVP foundation for a future SaaS that manages Prometheus scrape targets
for DevOps, SRE, and platform engineering teams.

---

## Why ptdgen?

Managing Prometheus `scrape_configs` by hand across multiple environments is
error-prone and repetitive. `ptdgen` gives you:

- A single source of truth for your scrape targets (one YAML / JSON file)
- Validated, deterministic output on every run
- Drop-in compatibility with Prometheus `file_sd_configs` and Grafana Alloy
- A foundation for GitOps: commit targets, generate on CI, deploy output

---

## Install / build

**Requirements:** Go 1.21+

```bash
git clone https://github.com/sdfnj/ptdgen
cd ptdgen
go build -o ptdgen .
```

Or run without installing:

```bash
go run . <command>
```

---

## Install from GitHub Releases

Download a prebuilt binary from the [GitHub Releases](https://github.com/forestian/prometheus-target-discovery-generator/releases) page.

**Linux / macOS:**

```bash
tar -xzf ptdgen_<version>_<os>_<arch>.tar.gz
chmod +x ptdgen
./ptdgen version
```

**Windows:**

Download the Windows archive, extract it, and run:

```
ptdgen.exe version
```

---

## Commands

### `ptdgen version`

```
ptdgen version
# ptdgen version 0.1.0
```

---

### `ptdgen init`

Creates an example project directory with sample targets and pre-generated output.

```bash
ptdgen init --output ./target-discovery-demo
ptdgen init --output ./target-discovery-demo --force   # overwrite existing
```

Creates:

```
target-discovery-demo/
├── README.md
├── targets.yaml                         ← edit this
├── targets.json                         ← edit this (JSON alternative)
├── generated/
│   ├── prometheus-file-sd.json          ← Prometheus file_sd_configs input
│   ├── alloy-discovery.alloy            ← Grafana Alloy starter config
│   └── scrape-config-example.yaml       ← example Prometheus scrape_config
└── examples/
    ├── prometheus.yml                   ← full example prometheus.yml
    └── alloy.river                      ← full example Alloy river config
```

---

### `ptdgen validate`

Validates a target definition file and reports errors.

```bash
ptdgen validate --file ./targets.yaml
ptdgen validate --file ./targets.json
```

Exit code 0 on success, 1 on any validation error.

---

### `ptdgen generate`

Reads and validates targets, then writes output files.

```bash
# Generate all formats (default)
ptdgen generate --file ./targets.yaml --output ./generated --format all

# Prometheus file_sd only
ptdgen generate --file ./targets.yaml --output ./generated --format prometheus

# Grafana Alloy config only
ptdgen generate --file ./targets.yaml --output ./generated --format alloy

# Overwrite existing files
ptdgen generate --file ./targets.yaml --output ./generated --force
```

Flags:

| Flag | Default | Description |
|------|---------|-------------|
| `--file` | (required) | Path to target definition file |
| `--output` | `./generated` | Output directory |
| `--format` | `all` | `prometheus`, `alloy`, or `all` |
| `--force` | `false` | Overwrite existing files |

---

## Example workflow

```bash
# 1. Scaffold a new project
ptdgen init --output ./my-targets

# 2. Edit targets
$EDITOR ./my-targets/targets.yaml

# 3. Validate
ptdgen validate --file ./my-targets/targets.yaml

# 4. Generate output
ptdgen generate --file ./my-targets/targets.yaml \
  --output ./my-targets/generated \
  --format all \
  --force

# 5. Copy generated/prometheus-file-sd.json to your Prometheus host
# 6. Point prometheus.yml file_sd_configs at it
```

---

## Target definition format

### YAML (`targets.yaml`)

```yaml
targets:
  - name: api-server-01          # unique name (required)
    job: app-api                 # Prometheus job label (required)
    address: 10.10.10.11         # host / IP (required)
    port: 9100                   # port 1–65535 (required)
    scheme: http                 # http or https (default: http)
    metrics_path: /metrics       # (default: /metrics)
    environment: production      # (default: unknown)
    team: platform               # (default: unknown)
    labels:                      # arbitrary extra labels
      service: api
      region: kr
```

### JSON (`targets.json`)

```json
{
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
      "labels": { "service": "api", "region": "kr" }
    }
  ]
}
```

---

## Generated output examples

### `prometheus-file-sd.json`

```json
[
  {
    "targets": ["10.10.10.11:9100"],
    "labels": {
      "job": "app-api",
      "instance": "api-server-01",
      "environment": "production",
      "team": "platform",
      "service": "api",
      "region": "kr",
      "__scheme__": "http",
      "__metrics_path__": "/metrics"
    }
  }
]
```

### `alloy-discovery.alloy`

```hcl
discovery.file "generated_targets" {
  files = ["/etc/alloy/file_sd/prometheus-file-sd.json"]
  refresh_interval = "30s"
}

prometheus.scrape "generated_targets" {
  targets    = discovery.file.generated_targets.targets
  forward_to = [prometheus.remote_write.default.receiver]
}

prometheus.remote_write "default" {
  endpoint {
    url = "http://mimir-nginx.monitoring.svc:80/api/v1/push"
  }
}
```

---

## Validation rules

- `name` must not be empty
- `job` must not be empty
- `address` must not be empty
- `port` must be between 1 and 65535
- Duplicate target names → error
- Duplicate `address:port` under the same `job` → error
- Existing output files → error unless `--force`

---

## MVP limitations

- No automatic target reloading; re-run `ptdgen generate` after changes
- No Kubernetes integration
- No Vault / secret management
- No cloud provider integration
- No authentication or multi-user support
- No web UI or API server

---

## Future roadmap (not yet implemented)

- Web UI to manage scrape targets
- REST API for target registration
- Vault integration for secret injection
- Kubernetes service discovery integration
- GitOps export (auto-commit generated files)
- Slack / PagerDuty alerts when targets change
- Mimir / Prometheus integration
- Team-based target ownership and RBAC
