# CLI Commands Module

## Overview

The `cmd` module contains the command-line interface implementations for the biometrics CLI. This module provides all user-facing commands for interacting with the biometric authentication system.

## Available Commands

### Authentication Commands

```bash
biometrics-cli auth enroll    # Enroll new biometric template
biometrics-cli auth verify    # Verify identity
biometrics-cli auth update    # Update biometric template
biometrics-cli auth delete    # Remove biometric enrollment
```

### Management Commands

```bash
biometrics-cli users list      # List all enrolled users
biometrics-cli users info      # Show user details
biometrics-cli users export   # Export user data
```

### System Commands

```bash
biometrics-cli system status   # Check system health
biometrics-cli system config   # View configuration
biometrics-cli system logs     # View system logs
```

## Command Structure

Each command follows a consistent structure:

```
Command Name: <action>
Description: <what it does>
Options:
  - --verbose: Enable verbose output
  - --format: Output format (json, yaml, table)
  - --config: Custom config file path
```

## Usage Examples

### Enroll New User
```bash
biometrics-cli auth enroll \
  --user-id user123 \
  --biometric face,voice \
  --quality high
```

### Verify Identity
```bash
biometrics-cli auth verify \
  --user-id user123 \
  --biometric face \
  --threshold 0.95
```

### Check System Status
```bash
biometrics-cli system status \
  --verbose \
  --format json
```

## Output Formats

### JSON Output
```bash
biometrics-cli users list --format json
{
  "users": [
    {"id": "user123", "enrolled": "2026-01-15"},
    {"id": "user456", "enrolled": "2026-01-20"}
  ]
}
```

### Table Output
```bash
biometrics-cli users list --format table
+----------+------------+
| User ID  | Enrolled   |
+----------+------------+
| user123  | 2026-01-15 |
| user456  | 2026-01-20 |
+----------+------------+
```

## Configuration

Commands can be configured via:
1. CLI flags (highest priority)
2. Environment variables
3. Config file (`~/.biometrics-cli.yaml`)
4. Default values (lowest priority)

### Environment Variables
```bash
export BIOMETRICS_API_URL="https://api.example.com"
export BIOMETRICS_API_KEY="your-api-key"
export BIOMETRICS_TIMEOUT=30
```

## Error Handling

All commands return appropriate exit codes:
- `0`: Success
- `1`: General error
- `2`: Invalid arguments
- `3`: Authentication error
- `4`: Resource not found

## Logging

Enable debug logging:
```bash
biometrics-cli --debug <command>
```

Logs are written to:
- Console (stdout)
- File: `~/.biometrics-cli/logs/cli.log`

## Extending Commands

To add new commands:
1. Create command file in `cmd/`
2. Register in `cmd/registry.go`
3. Add help text in `cmd/help.go`

## Dependencies

- Go 1.21+
- biometrics-core library
- OAuth2 client credentials

## See Also

- [CLI Documentation](../docs/cli.md)
- [API Documentation](../docs/api.md)
- [Configuration Guide](../docs/configuration.md)
