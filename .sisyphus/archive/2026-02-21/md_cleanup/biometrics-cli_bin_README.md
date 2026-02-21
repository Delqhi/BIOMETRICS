# CLI Binaries

**Purpose:** Compiled binary executables for the biometrics CLI application

## Overview

This directory contains the compiled binary executables for the biometrics CLI. These are ready-to-run binaries for various platforms.

## Binaries

| Binary | Platform | Architecture | Status |
|--------|----------|--------------|--------|
| `biometrics` | macOS/Linux | ARM64 (Apple Silicon) | Current |
| `biometrics-darwin-amd64` | macOS | x86_64 | Available |
| `biometrics-linux-amd64` | Linux | x86_64 | Available |
| `biometrics-windows-amd64.exe` | Windows | x86_64 | Available |

## Installation

### Quick Install (Recommended)
```bash
# Download latest release
curl -L -o biometrics https://github.com/delqhi/biometrics/releases/latest/download/biometrics

# Make executable
chmod +x biometrics

# Move to PATH
sudo mv biometrics /usr/local/bin/
```

### From Source
```bash
cd /Users/jeremy/dev/BIOMETRICS/biometrics-cli
go build -o biometrics ./cmd/biometrics
```

## Verification

### Check Version
```bash
./biometrics version
```

### Check Hash
```bash
sha256sum biometrics
```

## Usage

### Basic Command
```bash
biometrics [command] [flags]
```

### Common Commands

```bash
# Display help
biometrics --help

# Process biometric data
biometrics process --input sample.bmp

# Verify identity
biometrics verify --input sample.bmp --threshold 0.95
```

## Configuration

The CLI looks for configuration in:

1. Current directory: `./biometrics.yaml`
2. User home: `~/.biometrics/config.yaml`
3. System-wide: `/etc/biometrics/config.yaml`

## Permissions

The binary requires:

- **Read access** to input files
- **Write access** to output directories
- **Execute permission** on the binary itself

## Troubleshooting

### Permission Denied
```bash
chmod +x biometrics
```

### Binary Not Found
Ensure the binary is in your PATH:
```bash
export PATH=$PATH:/path/to/binaries
```

### Architecture Mismatch
Ensure you downloaded the correct binary for your platform:
```bash
uname -m  # Check architecture
```

## Version Compatibility

| Binary Version | Go Version | Minimum OS |
|---------------|------------|------------|
| 1.x | 1.21+ | macOS 11+, Ubuntu 20.04+ |

## Related Documentation

- [User Guide](../docs/user-guide.md)
- [Command Reference](../docs/commands.md)
- [Configuration](../docs/configuration.md)
- [Release Notes](../CHANGELOG.md)
