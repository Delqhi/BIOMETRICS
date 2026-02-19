# BIOMETRICS API Reference

Complete API documentation for the BIOMETRICS AI Agent Orchestration Platform.

## Base URLs

| Environment | URL |
|-------------|-----|
| Production | `https://api.biometrics.dev/v1` |
| Development | `http://localhost:8080/v1` |

## Authentication

All API endpoints require Bearer token authentication using NVIDIA API keys.

```bash
curl -H "Authorization: Bearer nvapi-YOUR_KEY" \
  https://api.biometrics.dev/v1/health
```

See [Authentication Guide](auth.md) for detailed setup instructions.

## Quick Start

### 1. Health Check

```bash
curl https://api.biometrics.dev/v1/health
```

**Response:**
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime": 86400,
  "timestamp": "2026-02-19T10:30:00Z"
}
```

### 2. List Available Agents

```bash
curl -H "Authorization: Bearer nvapi-YOUR_KEY" \
  https://api.biometrics.dev/v1/agents
```

**Response:**
```json
{
  "agents": [
    {
      "id": "agent-sisyphus-001",
      "name": "Sisyphus",
      "role": "coder",
      "status": "busy",
      "model": "qwen/qwen3.5-397b-a17b",
      "currentTask": "task-123",
      "tasksCompleted": 157,
      "createdAt": "2026-02-01T00:00:00Z",
      "lastActiveAt": "2026-02-19T10:29:00Z"
    }
  ],
  "total": 6
}
```

### 3. Create a Session

```bash
curl -X POST \
  -H "Authorization: Bearer nvapi-YOUR_KEY" \
  -H "Content-Type: application/json" \
  -d '{"title": "Build REST API"}' \
  https://api.biometrics.dev/v1/sessions
```

**Response:**
```json
{
  "id": "ses_abc123",
  "title": "Build REST API",
  "agentId": "agent-sisyphus-001",
  "status": "active",
  "messages": [],
  "createdAt": "2026-02-19T10:30:00Z",
  "updatedAt": "2026-02-19T10:30:00Z"
}
```

### 4. Send a Prompt

```bash
curl -X POST \
  -H "Authorization: Bearer nvapi-YOUR_KEY" \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Build a REST API with Express"}' \
  https://api.biometrics.dev/v1/sessions/ses_abc123/prompt
```

**Response:**
```json
{
  "id": "msg_xyz789",
  "content": "I'll help you build a REST API with Express...",
  "model": "qwen/qwen3.5-397b-a17b",
  "usage": {
    "promptTokens": 45,
    "completionTokens": 230,
    "totalTokens": 275
  },
  "createdAt": "2026-02-19T10:30:15Z"
}
```

## Core Resources

### Agents

AI agents that form the development swarm.

| Agent | Role | Model | Purpose |
|-------|------|-------|---------|
| **Sisyphus** | Coder | Qwen 3.5 397B | Main code implementation |
| **Prometheus** | Planner | Qwen 3.5 397B | Strategic planning |
| **Oracle** | Architect | Qwen 3.5 397B | Architecture review |
| **Atlas** | Heavy Lifting | Kimi K2.5 | Complex tasks |
| **Librarian** | Documenter | OpenCode ZEN (FREE) | Documentation |
| **Explore** | Explorer | OpenCode ZEN (FREE) | Code discovery |

### Sessions

Interactive sessions for agent conversations and task execution.

### Tasks

Individual units of work assigned to agents.

### Models

Available AI models for agent execution.

## Rate Limits

| Tier | RPM | Notes |
|------|-----|-------|
| Free | 40 RPM | NVIDIA Free Tier |
| Enterprise | Unlimited | Custom limits |

**HTTP 429 Response:**
```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests. Please wait 60 seconds."
  }
}
```

## Error Handling

All errors follow this format:

```json
{
  "error": {
    "code": "AGENT_NOT_FOUND",
    "message": "Agent with ID agent-001 not found",
    "details": {
      "agentId": "agent-001"
    }
  }
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Invalid or missing API key |
| `AGENT_NOT_FOUND` | 404 | Specified agent doesn't exist |
| `SESSION_NOT_FOUND` | 404 | Specified session doesn't exist |
| `TASK_NOT_FOUND` | 404 | Specified task doesn't exist |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Server error |

## SDKs & Libraries

### JavaScript/TypeScript

```typescript
import { BiometricsClient } from '@biometrics/sdk';

const client = new BiometricsClient('nvapi-YOUR_KEY');

// Create session
const session = await client.sessions.create({
  title: 'Build REST API'
});

// Send prompt
const response = await client.sessions.prompt(session.id, {
  prompt: 'Build a REST API with Express'
});

console.log(response.content);
```

### Python

```python
from biometrics import BiometricsClient

client = BiometricsClient('nvapi-YOUR_KEY')

# Create session
session = client.sessions.create(title='Build REST API')

# Send prompt
response = client.sessions.prompt(session.id, 
    prompt='Build a REST API with Express')

print(response.content)
```

### Go

```go
import "github.com/biometrics/go-sdk"

client := biometrics.NewClient("nvapi-YOUR_KEY")

// Create session
session, _ := client.Sessions.Create(context.Background(), 
    biometrics.CreateSessionRequest{
        Title: "Build REST API",
    })

// Send prompt
response, _ := client.Sessions.Prompt(context.Background(), 
    session.ID, 
    biometrics.PromptRequest{
        Prompt: "Build a REST API with Express",
    })

fmt.Println(response.Content)
```

## Webhooks

Subscribe to real-time events:

```bash
curl -X POST \
  -H "Authorization: Bearer nvapi-YOUR_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://your-server.com/webhook",
    "events": ["task.completed", "task.failed", "agent.status_changed"]
  }' \
  https://api.biometrics.dev/v1/webhooks
```

### Available Events

| Event | Description |
|-------|-------------|
| `task.created` | New task created |
| `task.started` | Task execution started |
| `task.completed` | Task finished successfully |
| `task.failed` | Task failed with error |
| `task.cancelled` | Task was cancelled |
| `agent.status_changed` | Agent status updated |
| `session.created` | New session created |
| `session.completed` | Session finished |

## Best Practices

### 1. Connection Pooling

Reuse HTTP connections for better performance:

```typescript
// ❌ BAD - New connection each request
const response1 = await fetch('/agents');
const response2 = await fetch('/tasks');

// ✅ GOOD - Reuse connection
const client = new BiometricsClient('nvapi-YOUR_KEY');
const agents = await client.agents.list();
const tasks = await client.tasks.list();
```

### 2. Error Handling

Always handle errors gracefully:

```typescript
try {
  const task = await client.tasks.create({
    title: 'Build feature'
  });
} catch (error) {
  if (error.code === 'RATE_LIMIT_EXCEEDED') {
    // Wait and retry
    await sleep(60000);
    return retry();
  }
  throw error;
}
```

### 3. Pagination

Use pagination for large lists:

```typescript
// Get all agents in batches
const allAgents = [];
let offset = 0;
const limit = 20;

while (true) {
  const response = await client.agents.list({ limit, offset });
  allAgents.push(...response.agents);
  
  if (response.agents.length < limit) break;
  offset += limit;
}
```

### 4. Streaming

Enable streaming for long-running tasks:

```typescript
const stream = await client.sessions.prompt(sessionId, {
  prompt: 'Build a complete application',
  stream: true
});

for await (const chunk of stream) {
  process.stdout.write(chunk.content);
}
```

## Changelog

### v1.0.0 (2026-02-19)

- Initial API release
- Agent management endpoints
- Session management endpoints
- Task management endpoints
- Model listing endpoints
- Health check endpoints

## Support

- **Documentation**: https://github.com/Delqhi/BIOMETRICS/docs
- **Issues**: https://github.com/Delqhi/BIOMETRICS/issues
- **Discord**: https://discord.gg/biometrics
- **Email**: support@biometrics.dev

---

**Version:** 1.0.0 | **Last Updated:** February 2026
