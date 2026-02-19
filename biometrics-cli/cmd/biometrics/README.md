# Biometrics CLI Command

**Purpose:** Main CLI binary implementation for the biometrics command-line tool

## Overview

This directory contains the main entry point for the `biometrics` CLI command. It uses the Cobra framework for command-line interface construction.

## Files

| File | Description |
|------|-------------|
| `main.go` | Main entry point for the biometrics CLI |

## Architecture

```
biometrics CLI
├── Root Command (biometrics)
│   ├── auth Command
│   ├── config Command  
│   ├── process Command
│   ├── verify Command
│   └── version Command
├── Global Flags
├── Configuration Loading
└── Plugin System
```

## Commands

### Root Command
```bash
biometrics [command]
```

### Available Commands

| Command | Description |
|---------|-------------|
| `auth` | Authentication and authorization |
| `config` | Configuration management |
| `process` | Process biometric data |
| `verify` | Verify biometric samples |
| `version` | Display version information |

## Usage Examples

### Basic Usage
```bash
# Process a biometric sample
biometrics process --input sample.bmp --type face

# Verify identity
biometrics verify --input sample.bmp --database users.db

# Manage authentication
biometrics auth login --provider azure
```

### With Configuration
```bash
# Use custom config
biometrics --config /path/to/config.yaml process --input sample.bmp

# Verbose output
biometrics -v --debug process --input sample.bmp
```

## Configuration

The CLI can be configured via:

1. **Config file:** `~/.biometrics/config.yaml`
2. **Environment variables:** `BIOMETRICS_*`
3. **Command-line flags:** `--flag value`

### Example Configuration

```yaml
# ~/.biometrics/config.yaml
version: "1.0"
provider:
  type: azure
  tenant_id: "your-tenant-id"
processing:
  workers: 4
  batch_size: 100
logging:
  level: info
  format: json
```

## Building

### Build for Current Platform
```bash
cd /Users/jeremy/dev/BIOMETRICS/biometrics-cli
go build -o biometrics ./cmd/biometrics
```

### Cross-Platform Build
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o biometrics-linux ./cmd/biometrics

# Windows
GOOS=windows GOARCH=amd64 go build -o biometrics.exe ./cmd/biometrics
```

## Testing

```bash
go test ./cmd/biometrics/... -v
```

## Integration

This CLI can be integrated with:

- **Shell scripts** for automation
- **CI/CD pipelines** for automated workflows
- **Other Go programs** via import
- **External tools** via stdin/stdout

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `BIOMETRICS_CONFIG_PATH` | Config file path | `~/.biometrics/config.yaml` |
| `BIOMETRICS_LOG_LEVEL` | Logging level | `info` |
| `BIOMETRICS_API_KEY` | API authentication key | - |

## Related Documentation

- [CLI User Guide](../docs/user-guide.md)
- [Configuration Reference](../docs/configuration.md)
- [Command Reference](../docs/commands.md)
- [Examples](../examples/)
