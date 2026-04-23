# envoy-conf

> A CLI tool to diff and validate environment variable configs across deployment targets before rollout.

---

## Installation

```bash
go install github.com/yourorg/envoy-conf@latest
```

Or download a prebuilt binary from the [Releases](https://github.com/yourorg/envoy-conf/releases) page.

---

## Usage

**Diff configs between two targets:**

```bash
envoy-conf diff --source staging --target production
```

**Validate a config against a schema or baseline:**

```bash
envoy-conf validate --env production --schema ./schema.yaml
```

**Example output:**

```
[MISSING]  DATABASE_URL     found in staging, not in production
[MISMATCH] LOG_LEVEL        staging=debug | production=debug ✓
[EXTRA]    DEBUG_MODE       found in production, not in staging
```

### Common Flags

| Flag | Description |
|------|-------------|
| `--source` | Source deployment target |
| `--target` | Target deployment target |
| `--schema` | Path to validation schema file |
| `--output` | Output format: `text`, `json`, `yaml` |
| `--strict` | Exit with non-zero code on any diff |

---

## Configuration

`envoy-conf` reads targets from a local `envoy-conf.yaml` file or environment variables. See [docs/configuration.md](docs/configuration.md) for full details.

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

MIT © [yourorg](https://github.com/yourorg)