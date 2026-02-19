# Go Packages

**Purpose:** Reusable Go packages and libraries for the biometrics CLI

## Overview

This directory contains reusable Go packages that provide core functionality for the biometrics CLI application.

## Structure

```
pkg/
├── auth/         # Authentication
├── workflows/   # Workflow management
├── agents/      # AI agent integrations
├── mcp/         # Model Context Protocol
├── delegation/  # Task delegation patterns
└── ratelimit/  # Rate limiting
```

## Packages

### auth/

Authentication and authorization utilities.

```go
import "github.com/delqhi/biometrics/pkg/auth"

// Usage
provider := auth.NewProvider("azure")
token, err := provider.Authenticate(credentials)
```

### workflows/

Workflow management and orchestration.

```go
import "github.com/delqhi/biometrics/pkg/workflows"

// Usage
wf := workflows.NewWorkflow(ctx, config)
result, err := wf.Execute(ctx, input)
```

### agents/

AI agent integrations and abstractions.

```go
import "github.com/delqhi/biometrics/pkg/agents"

// Usage
agent := agents.New(agentConfig)
response, err := agent.Process(ctx, prompt)
```

### mcp/

Model Context Protocol implementation.

```go
import "github.com/delqhi/biometrics/pkg/mcp"

// Usage
client := mcp.NewClient(endpoint, opts)
result, err := client.Call(ctx, method, params)
```

### delegation/

Task delegation patterns and strategies.

```go
import "github.com/delqhi/biometrics/pkg/delegation"

// Usage
strategy := delegation.NewStrategy(config)
tasks, err := strategy.Distribute(workItems)
```

### ratelimit/

Rate limiting utilities.

```go
import "github.com/delqhi/biometrics/pkg/ratelimit"

// Usage
limiter := ratelimit.NewLimiter(rps)
err := limiter.Acquire(ctx)
```

## Usage in Commands

Import packages in commands:
```go
package main

import (
    "github.com/delqhi/biometrics/pkg/auth"
    "github.com/delqhi/biometrics/pkg/workflows"
)

func main() {
    // Use packages
}
```

## Testing

```bash
# Test all packages
go test ./pkg/... -v

# Test specific package
go test ./pkg/auth/... -v

# Run with coverage
go test ./pkg/... -coverprofile=coverage.out
```

## Documentation

Each package includes:
- README.md (this file)
- Usage examples
- Unit tests
- API documentation

## Versioning

Packages follow semantic versioning:
- Major: Breaking changes
- Minor: New features
- Patch: Bug fixes

## Related Documentation

- [CLI Commands](../commands/)
- [Configuration](../docs/configuration.md)
- [API Reference](../docs/api/)
