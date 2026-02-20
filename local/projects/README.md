# Local Projects Directory

## Overview

This directory contains local project configurations, workspace setups, and development environment definitions used for local development and testing.

## Contents

### Project Templates

| Directory | Description |
|-----------|-------------|
| templates/ | Project starter templates |
| examples/ | Example project configurations |

### Development Environments

| Directory | Description |
|-----------|-------------|
| docker/ | Docker Compose configurations |
| kubernetes/ | K8s development setups |
| vagrant/ | Vagrant VM configurations |

### Workspace Configs

| File | Description |
|------|-------------|
| workspace.code-workspace | VS Code workspace |
| workspace.code-workspace | JetBrains workspace |

## Project Structure

### Template Structure
```
templates/
├── nodejs/
│   ├── package.json
│   ├── src/
│   └── test/
├── python/
│   ├── requirements.txt
│   ├── src/
│   └── tests/
└── go/
    ├── go.mod
    ├── cmd/
    └── pkg/
```

## Local Development

### Quick Start
```bash
# Create new project
biometrics-cli project create my-project --template nodejs

# Start local environment
cd my-project
docker-compose up

# Run development server
biometrics-cli dev start
```

### Environment Variables
```bash
# .env.local
BIOMETRICS_API_URL=http://localhost:8080
BIOMETRICS_API_KEY=test-key
BIOMETRICS_DEBUG=true
```

## Docker Setup

### docker-compose.yml
```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://db:5432
    depends_on:
      - db
      - redis

  db:
    image: postgres:15
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

## Tools

### Local CLI
```bash
# Development commands
biometrics-cli dev start      # Start dev environment
biometrics-cli dev stop       # Stop environment
biometrics-cli dev logs       # View logs
biometrics-cli dev restart    # Restart services
```

### Debugging
```bash
# Attach debugger
biometrics-cli dev debug --port 9229

# View logs
biometrics-cli dev logs --follow
```

## Configuration

### Local Settings
```json
{
  "dev": {
    "autoReload": true,
    "port": 8080,
    "debugPort": 9229,
    "environment": "development"
  }
}
```

## Maintenance

- Clean up unused projects
- Update templates regularly
- Remove temporary files

## Best Practices

1. **Isolate**: Use separate directories per project
2. **Document**: Add README to each project
3. **Version**: Use version control
4. **Clean**: Remove when done

## See Also

- [Docker Configurations](./docker/)
- [Scripts Directory](../scripts/)
- [Configuration Docs](../docs/configuration.md)
