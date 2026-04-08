# envoy-diff

> CLI tool to diff and audit environment variable changes across deployment configs

---

## Installation

```bash
go install github.com/yourusername/envoy-diff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envoy-diff.git
cd envoy-diff
go build -o envoy-diff .
```

---

## Usage

Compare environment variables between two deployment config files:

```bash
envoy-diff compare staging.yaml production.yaml
```

Audit a single config for missing or sensitive variables:

```bash
envoy-diff audit deployment.yaml --strict
```

**Example output:**

```
~ DATABASE_URL   [changed]
+ NEW_FEATURE_FLAG  "true"
- DEPRECATED_KEY
```

Supported formats: `.yaml`, `.json`, `.env`

---

## Flags

| Flag | Description |
|------|-------------|
| `--strict` | Exit with non-zero code if differences are found |
| `--format` | Output format: `text`, `json`, `table` (default: `text`) |
| `--ignore` | Comma-separated list of keys to ignore |

---

## License

[MIT](LICENSE)