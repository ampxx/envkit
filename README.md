# envkit

A CLI tool for managing and validating environment variable configs across multiple deployment targets.

---

## Installation

```bash
go install github.com/yourname/envkit@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/yourname/envkit/releases).

---

## Usage

Define your environment schema in a `.envkit.yaml` file:

```yaml
targets:
  production:
    required:
      - DATABASE_URL
      - API_KEY
      - PORT
    optional:
      - LOG_LEVEL
  staging:
    required:
      - DATABASE_URL
      - API_KEY
```

Then validate your environment against a target:

```bash
# Validate current environment against a deployment target
envkit validate --target production

# Check a specific .env file
envkit validate --target staging --file .env.staging

# List all defined targets
envkit targets list
```

Example output:

```
✔ DATABASE_URL   set
✔ API_KEY        set
✗ PORT           missing (required)

1 error found for target: production
```

---

## Why envkit?

- Catch missing environment variables before deployment
- Manage configs across `production`, `staging`, `dev`, and more
- Simple YAML-based schema with no runtime dependencies

---

## License

MIT © yourname