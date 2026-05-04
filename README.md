# ptdgen — Prometheus Target Discovery Generator

Generate Prometheus `file_sd_configs` and Grafana Alloy discovery configs from a single YAML source of truth.

---

## Demo

GIF demo coming soon.

---

## Quick Start

**Download a prebuilt binary** from [GitHub Releases](https://github.com/forestian/prometheus-target-discovery-generator/releases):

```bash
# Linux / macOS
tar -xzf ptdgen_<version>_<os>_<arch>.tar.gz
chmod +x ptdgen
./ptdgen version

# Windows — extract the archive and run:
ptdgen.exe version
```

**Or build from source** (Go 1.21+):

```bash
git clone https://github.com/forestian/prometheus-target-discovery-generator
cd prometheus-target-discovery-generator
go build -o ptdgen .
./ptdgen version
```

---

## Quick Demo

```bash
$ ptdgen generate --file ./targets.yaml --output ./generated
Generating all output into ./generated ...
Done. 3 target(s) processed.

$ ls ./generated/
prometheus-file-sd.json   alloy-discovery.alloy   scrape-config-example.yaml
```

---

## Use Cases

- Replace hand-edited `scrape_configs` with a validated, version-controlled YAML file
- Bootstrap Grafana Alloy discovery configs for new environments in seconds
- Validate scrape target definitions in CI before deploying Prometheus changes
- Maintain a single source of truth across multiple Prometheus instances
- Scaffold new monitoring projects with `ptdgen init`

---

## Commands

### `ptdgen init`

Scaffold a new project directory with sample targets and pre-generated output:

```bash
ptdgen init --output ./my-targets
ptdgen init --output ./my-targets --force   # overwrite existing
```

Creates:

```
my-targets/
├── README.md
├── targets.yaml                         ← edit this
├── targets.json                         ← JSON alternative
├── generated/
│   ├── prometheus-file-sd.json
│   ├── alloy-discovery.alloy
│   └── scrape-config-example.yaml
└── examples/
    ├── prometheus.yml
    └── alloy.river
```

---

### `ptdgen validate`

Validate a target definition file and report errors:

```bash
ptdgen validate --file ./targets.yaml
ptdgen validate --file ./targets.json
```

Exit code 0 on success, 1 on any validation error.

---

### `ptdgen generate`

Parse, validate, and write output files:

```bash
# Generate all formats (default)
ptdgen generate --file ./targets.yaml --output ./generated

# Prometheus file_sd only
ptdgen generate --file ./targets.yaml --output ./generated --format prometheus

# Grafana Alloy config only
ptdgen generate --file ./targets.yaml --output ./generated --format alloy

# Overwrite existing files
ptdgen generate --file ./targets.yaml --output ./generated --force
```

| Flag | Default | Description |
|------|---------|-------------|
| `--file` | (required) | Path to target definition file |
| `--output` | `./generated` | Output directory |
| `--format` | `all` | `prometheus`, `alloy`, or `all` |
| `--force` | `false` | Overwrite existing files |

---

## Typical Workflow

```bash
# 1. Scaffold a new project
ptdgen init --output ./my-targets

# 2. Edit your targets
$EDITOR ./my-targets/targets.yaml

# 3. Validate
ptdgen validate --file ./my-targets/targets.yaml

# 4. Generate
ptdgen generate --file ./my-targets/targets.yaml \
  --output ./my-targets/generated \
  --force

# 5. Copy generated/prometheus-file-sd.json to your Prometheus host
# 6. Point prometheus.yml file_sd_configs at it
```

---

## Target Definition Format

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

JSON format is also supported — see [`examples/targets.yaml`](examples/targets.yaml).

---

## Example Output

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

## Validation Rules

- `name` must not be empty and must be unique across all targets
- `job` must not be empty
- `address` must not be empty
- `port` must be between 1 and 65535
- Duplicate `address:port` under the same `job` → error
- Existing output files → error unless `--force`

---

## Limitations

- Read-only and local-only — ptdgen never connects to or modifies Prometheus, Alloy, or any live system
- No secrets, tokens, or credentials are read or written
- Generated output should be reviewed before applying to production
- No automatic target reloading — re-run `ptdgen generate` after changes
- No Kubernetes, Vault, or cloud provider integration

---

## Roadmap

- GitHub Actions example workflow for GitOps-style target management
- Watch mode to regenerate on file change
- Additional output formats (e.g. plain `scrape_configs` YAML block)
- Multi-file target definitions
- Improved validation errors with line number references

---

Part of the [Forestian Cloud Native Toolkit](https://github.com/forestian) — small CLI tools for Kubernetes, observability, GitOps, and platform engineering.
