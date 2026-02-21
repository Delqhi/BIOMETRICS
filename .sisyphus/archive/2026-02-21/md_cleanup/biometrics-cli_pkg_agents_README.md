# AI Agent Package

**Purpose:** AI agent abstractions, integrations, and management

## Overview

This package provides interfaces and implementations for AI agent integrations used throughout the biometrics CLI.

## Components

### Agent Interface

```go
type Agent interface {
    // Process handles a prompt and returns a response
    Process(ctx context.Context, prompt string) (Response, error)
    
    // Stream processes with streaming responses
    Stream(ctx context.Context, prompt string) (<-chan Response, error)
    
    // Name returns the agent's identifier
    Name() string
}
```

### Supported Agents

| Agent | Provider | Capabilities |
|-------|----------|--------------|
| Sisyphus | NVIDIA Qwen 3.5 | Code generation, planning |
| Prometheus | NVIDIA Qwen 3.5 | Strategic planning |
| Oracle | NVIDIA Qwen 3.5 | Architecture review |
| Atlas | Kimi K2.5 | Heavy lifting |
| Librarian | OpenCode ZEN | Documentation |

## Usage

### Creating an Agent

```go
import "github.com/delqhi/biometrics/pkg/agents"

// Create agent from config
agent, err := agents.NewAgent(ctx, agents.Config{
    Name: "sisyphus",
    Model: "qwen/qwen3.5-397b-a17b",
    Provider: "nvidia",
})
```

### Processing Requests

```go
// Simple request
response, err := agent.Process(ctx, "Write a function to process biometrics")

// Streaming
stream, err := agent.Stream(ctx, "Generate code for...")
for chunk := range stream {
    fmt.Print(chunk.Content)
}
```

### Agent Pool

```go
// Create pool of agents
pool := agents.NewPool([]agents.Agent{sisyphus, atlas})

// Distribute work
results, err := pool.ProcessAll(ctx, tasks)
```

## Configuration

```yaml
agents:
  sisyphus:
    provider: nvidia
    model: qwen/qwen3.5-397b-a17b
    max_tokens: 32768
    temperature: 0.7
    
  atlas:
    provider: opencode
    model: kimi-k2.5-free
    max_tokens: 65536
```

## Rate Limiting

The package includes built-in rate limiting:

```go
config := agents.Config{
    Name: "agent-name",
    RateLimit: agents.RateLimit{
        RequestsPerMinute: 60,
        TokensPerMinute: 100000,
    },
}
```

## Error Handling

```go
response, err := agent.Process(ctx, prompt)
if err != nil {
    switch {
    case errors.Is(err, agents.ErrRateLimited):
        // Handle rate limit
    case errors.Is(err, agents.ErrModelUnavailable):
        // Fallback to another agent
    default:
        return err
    }
}
```

## Testing

```bash
go test ./pkg/agents/... -v -run TestAgent
go test ./pkg/agents/... -v -cover
```

## Related Documentation

- [Agent Configuration](../config/agents.md)
- [Provider Setup](../config/providers.md)
- [Delegation Patterns](../patterns/delegation/)
