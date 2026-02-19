# MCP (Model Context Protocol)

**Purpose:** Model Context Protocol implementation for AI model integration

## Overview

This package implements the Model Context Protocol (MCP) for standardizing communication between the biometrics CLI and various AI models and services.

## Features

- Standardized model interface
- Streaming support
- Tool calling
- Context management
- Connection pooling

## Usage

### Creating MCP Client

```go
import "github.com/delqhi/biometrics/pkg/mcp"

client, err := mcp.NewClient(ctx, mcp.Config{
    Endpoint: "https://api.nvidia.com/v1",
    Model: "qwen/qwen3.5-397b-a17b",
    Timeout: 120 * time.Second,
})
```

### Making Requests

```go
// Simple request
response, err := client.Complete(ctx, &mcp.Request{
    Prompt: "Analyze this biometric sample",
    Options: mcp.Options{
        Temperature: 0.7,
        MaxTokens: 4096,
    },
})

// Streaming
stream, err := client.CompleteStream(ctx, req)
for chunk := range stream {
    fmt.Print(chunk.Content)
}
```

### Tool Calling

```go
response, err := client.Complete(ctx, &mcp.Request{
    Prompt: "Process the biometric data",
    Tools: []mcp.Tool{
        {
            Name: "process_image",
            Description: "Process a biometric image",
            Parameters: schema,
        },
    },
})
```

## Configuration

```yaml
mcp:
  providers:
    nvidia:
      endpoint: https://integrate.api.nvidia.com/v1
      model: qwen/qwen3.5-397b-a17b
      timeout: 120s
    opencode:
      endpoint: https://api.opencode.ai/v1
      model: kimi-k2.5-free
      timeout: 60s
```

## Connection Pooling

```go
client, err := mcp.NewPooledClient(ctx, mcp.PoolConfig{
    Endpoint: endpoint,
    MaxConnections: 10,
    IdleTimeout: 5 * time.Minute,
})
```

## Error Handling

```go
_, err := client.Complete(ctx, req)
if err != nil {
    switch {
    case errors.Is(err, mcp.ErrRateLimited):
        // Handle rate limiting
    case errors.Is(err, mcp.ErrModelUnavailable):
        // Switch to fallback
    case errors.Is(err, mcp.ErrTimeout):
        // Retry with backoff
    }
}
```

## Related Documentation

- [Model Configuration](../config/models.md)
- [Provider Setup](../config/providers.md)
- [Agent Integration](../agents/)
