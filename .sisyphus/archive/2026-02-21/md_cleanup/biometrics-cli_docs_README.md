# Biometrics CLI Documentation

## Overview

This directory contains comprehensive documentation for the biometrics CLI application. The documentation covers installation, usage, configuration, and advanced features.

## Documentation Structure

```
docs/
├── getting-started.md      # Quick start guide
├── installation.md          # Installation instructions
├── configuration.md         # Configuration reference
├── commands/               # Command-specific docs
│   ├── auth.md            # Authentication commands
│   ├── users.md           # User management
│   └── system.md          # System commands
├── api/                   # API reference
│   ├── authentication.md  # API auth
│   ├── endpoints.md       # REST endpoints
│   └── errors.md         # Error codes
├── advanced/              # Advanced topics
│   ├── performance.md     # Performance tuning
│   ├── security.md       # Security best practices
│   └── integration.md    # Integration guides
└── troubleshooting.md     # Common issues
```

## Quick Start

### Installation
```bash
# Using Homebrew
brew install biometrics-cli

# Using Go
go install github.com/org/biometrics-cli@latest
```

### Basic Usage
```bash
# Initialize configuration
biometrics-cli init

# Enroll a new user
biometrics-cli auth enroll --user-id john@example.com

# Verify identity
biometrics-cli auth verify --user-id john@example.com
```

## Configuration

### Config File Location
- Linux/macOS: `~/.biometrics-cli/config.yaml`
- Windows: `%USERPROFILE%\.biometrics-cli\config.yaml`

### Example Configuration
```yaml
api:
  base_url: "https://api.biometrics.example.com"
  timeout: 30
  retry:
    max_attempts: 3
    backoff: "exponential"

auth:
  method: "oauth2"
  client_id: "your-client-id"
  client_secret: "your-secret"

biometrics:
  default_threshold: 0.95
  liveness_detection: true

logging:
  level: "info"
  file: "~/.biometrics-cli/logs/app.log"
```

## Commands

### Authentication Commands

| Command | Description |
|---------|-------------|
| `auth enroll` | Enroll new biometric template |
| `auth verify` | Verify identity |
| `auth update` | Update existing template |
| `auth delete` | Remove enrollment |

### User Management

| Command | Description |
|---------|-------------|
| `users list` | List all users |
| `users info` | Get user details |
| `users export` | Export user data |

### System Commands

| Command | Description |
|---------|-------------|
| `system status` | Check system health |
| `system config` | View configuration |
| `system version` | Show version info |

## API Reference

### REST Endpoints

```
POST /api/v1/auth/enroll
POST /api/v1/auth/verify
POST /api/v1/auth/update
DELETE /api/v1/auth/{user_id}

GET  /api/v1/users
GET  /api/v1/users/{user_id}
DELETE /api/v1/users/{user_id}

GET  /api/v1/health
GET  /api/v1/metrics
```

### Authentication
All API requests require Bearer token authentication:
```
Authorization: Bearer <your-token>
```

## Advanced Features

### Performance Optimization
- Connection pooling for API requests
- Local caching of templates
- Batch processing for bulk operations

### Security
- Encrypted credential storage
- TLS 1.3 for all communications
- Audit logging for all operations

### Integration
- Webhook support for events
- Custom plugin system
- Export to multiple formats

## Troubleshooting

### Common Issues

1. **Authentication Failed**
   - Check API credentials
   - Verify network connectivity
   - Ensure correct base URL

2. **Biometric Verification Failed**
   - Check image/audio quality
   - Verify liveness detection
   - Ensure proper lighting

3. **Performance Issues**
   - Increase timeout values
   - Enable local caching
   - Check network latency

## Getting Help

- Documentation: https://docs.biometrics-cli.example.com
- GitHub Issues: https://github.com/org/biometrics-cli/issues
- Support Email: support@biometrics-cli.example.com

## Version

Current Version: 1.5.0
Last Updated: 2026-02-20
