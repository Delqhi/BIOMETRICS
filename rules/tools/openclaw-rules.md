# OpenClaw Rules - Comprehensive Guide

**Version:** 1.0  
**Status:** ACTIVE  
**Last Updated:** 2026-02-20  
**Project:** BIOMETRICS  
**Location:** `/Users/jeremy/dev/BIOMETRICS/rules/tools/openclaw-rules.md`

---

## Table of Contents

1. [OpenClaw Installation](#1-openclaw-installation)
2. [Agent Configuration](#2-agent-configuration)
3. [MCP Server Setup](#3-mcp-server-setup)
4. [Workflow Automation](#4-workflow-automation)
5. [Configuration](#5-configuration)
6. [Gateway Management](#6-gateway-management)
7. [Best Practices](#7-best-practices)
8. [Troubleshooting](#8-troubleshooting)

---

## 1. OpenClaw Installation

### 1.1 Prerequisites

Before installing OpenClaw, ensure your system meets the following requirements:

- **Operating System:** macOS 12+, Linux (Ubuntu 20.04+), or Windows 11 with WSL2
- **Docker:** Docker Engine 20.10+ and Docker Compose 2.0+
- **Node.js:** Version 18.0 or higher
- **Python:** Version 3.9 or higher (for MCP server integrations)
- **Network:** Access to NVIDIA NIM API (for GPU models) or alternative model providers

### 1.2 Installation Steps

OpenClaw can be installed via multiple methods. Choose the one that best fits your environment:

#### Method A: Homebrew Installation (macOS)

```bash
# Install OpenClaw via Homebrew
brew install openclaw/tap/openclaw

# Verify installation
openclaw --version

# Check available commands
openclaw --help
```

#### Method B: Direct Binary Download

```bash
# Download the latest binary for your platform
# For macOS ARM64:
curl -L -o /usr/local/bin/openclaw https://github.com/anomalyco/opencode/releases/latest/download/openclaw-darwin-arm64

# For macOS x86_64:
curl -L -o /usr/local/bin/openclaw https://github.com/anomalyco/opencode/releases/latest/download/openclaw-darwin-x86_64

# For Linux:
curl -L -o /usr/local/bin/openclaw https://github.com/anomalyco/opencode/releases/latest/download/openclaw-linux-x86_64

# Make executable
chmod +x /usr/local/bin/openclaw

# Verify installation
openclaw --version
```

#### Method C: Python pip Installation

```bash
# Create virtual environment (recommended)
python3 -m venv openclaw-env
source openclaw-env/bin/activate

# Install OpenClaw
pip install openclaw

# Verify installation
openclaw --version
```

#### Method D: Docker Installation

```bash
# Pull the OpenClaw Docker image
docker pull anomalyco/openclaw:latest

# Create alias for convenient CLI usage
alias openclaw='docker run --rm -it -v ~/.openclaw:/root/.openclaw anomalyco/openclaw'

# Verify installation
openclaw --version
```

### 1.3 Configuration Directory Setup

OpenClaw requires a configuration directory at `~/.openclaw/`. This directory contains all configuration files, credentials, and runtime data.

```bash
# Create OpenClaw configuration directory
mkdir -p ~/.openclaw

# Set appropriate permissions (restrict access to protect credentials)
chmod 700 ~/.openclaw

# Verify directory exists
ls -la ~/.openclaw/
```

### 1.4 Initial Configuration

After installation, you need to configure OpenClaw with your model providers and agents. The main configuration file is `~/.openclaw/openclaw.json`.

#### Creating the Initial Configuration

```bash
# Create default configuration file
openclaw init

# This creates ~/.openclaw/openclaw.json with default settings
```

#### Verification Commands

After installation, verify that OpenClaw is properly installed and configured:

```bash
# Check OpenClaw version
openclaw --version

# Verify configuration is valid
openclaw config validate

# Check available models
openclaw models

# List configured agents
openclaw agents list

# Test connectivity to default provider
openclaw doctor
```

### 1.5 Environment Variables

OpenClaw supports several environment variables for configuration override:

```bash
# Provider API Keys
export NVIDIA_API_KEY="nvapi-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
export OPENROUTER_API_KEY="sk-or-v1-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
export ANTHROPIC_API_KEY="sk-ant-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

# OpenClaw Configuration
export OPENCLAW_CONFIG_PATH="/path/to/custom/config.json"
export OPENCLAW_DATA_DIR="/path/to/data/directory"
export OPENCLAW_LOG_LEVEL="debug"

# Gateway Configuration
export OPENCLAW_GATEWAY_PORT="18789"
export OPENCLAW_GATEWAY_HOST="0.0.0.0"

# MCP Server Configuration
export MCP_SERVER_TIMEOUT="30000"
export MCP_MAX_CONNECTIONS="10"
```

---

## 2. Agent Configuration

### 2.1 Agent Presets

OpenClaw comes with predefined agent presets that are optimized for different task types. Each preset has specific configurations for model selection, temperature, max tokens, and specialized capabilities.

#### Available Agent Presets

| Preset Name | Primary Use Case | Best For |
|-------------|------------------|----------|
| **sisyphus** | Code implementation | Writing, refactoring, fixing bugs |
| **prometheus** | Planning and architecture | System design, roadmaps, technical decisions |
| **atlas** | Heavy lifting | Data processing, bulk operations |
| **oracle** | Architecture review | Code review, best practices validation |
| **librarian** | Research and documentation | Finding information, writing docs |
| **explore** | Code discovery | Finding patterns, understanding codebases |
| **metis** | Deep analysis | Complex problem solving, research |
| **momus** | Markdown generation | Writing documentation, reports |

#### Agent Preset Configuration Example

```json
{
  "agents": {
    "presets": {
      "sisyphus": {
        "description": "Primary code implementation agent",
        "model": {
          "provider": "nvidia",
          "model_id": "qwen/qwen3.5-397b-a17b",
          "temperature": 0.7,
          "max_tokens": 32768
        },
        "capabilities": {
          "code_generation": true,
          "refactoring": true,
          "bug_fixing": true,
          "testing": true
        },
        "fallback_chain": [
          "nvidia/moonshotai/kimi-k2.5",
          "nvidia/meta/llama-3.3-70b-instruct"
        ]
      },
      "prometheus": {
        "description": "Planning and architecture agent",
        "model": {
          "provider": "nvidia",
          "model_id": "qwen/qwen3.5-397b-a17b",
          "temperature": 0.5,
          "max_tokens": 16384
        },
        "capabilities": {
          "planning": true,
          "architecture": true,
          "roadmapping": true,
          "estimation": true
        }
      },
      "atlas": {
        "description": "Heavy lifting and data processing",
        "model": {
          "provider": "opencode-zen",
          "model_id": "kimi-k2.5-free",
          "temperature": 0.3,
          "max_tokens": 65536
        },
        "capabilities": {
          "bulk_processing": true,
          "data_transformation": true,
          "batch_operations": true
        }
      }
    }
  }
}
```

### 2.2 Model Assignment

Model assignment defines which AI models are used for different task types. OpenClaw supports multiple providers and allows fine-grained control over model selection.

#### Model Assignment Rules

```json
{
  "model_assignment": {
    "rules": [
      {
        "task_type": "code_generation",
        "primary_model": "nvidia/qwen/qwen3.5-397b-a17b",
        "fallbacks": [
          "nvidia/moonshotai/kimi-k2.5",
          "opencode-zen/grok-code"
        ],
        "timeout_ms": 120000
      },
      {
        "task_type": "code_review",
        "primary_model": "nvidia/meta/llama-3.3-70b-instruct",
        "fallbacks": [
          "nvidia/mistralai/mistral-large-3-675b-instruct-2512"
        ],
        "timeout_ms": 90000
      },
      {
        "task_type": "research",
        "primary_model": "opencode-zen/kimi-k2.5-free",
        "fallbacks": [],
        "timeout_ms": 60000
      },
      {
        "task_type": "documentation",
        "primary_model": "opencode-zen/minimax-m2.5-free",
        "fallbacks": [],
        "timeout_ms": 45000
      }
    ]
  }
}
```

### 2.3 Fallback Chains

Fallback chains ensure task continuity when the primary model fails or is rate-limited. OpenClaw automatically attempts the next model in the chain if the current one fails.

#### Fallback Chain Configuration

```json
{
  "fallback_chains": {
    "default": {
      "chain": [
        {
          "provider": "nvidia",
          "model_id": "qwen/qwen3.5-397b-a17b",
          "priority": 1,
          "timeout_ms": 120000,
          "retry_on_error": true
        },
        {
          "provider": "nvidia",
          "model_id": "moonshotai/kimi-k2.5",
          "priority": 2,
          "timeout_ms": 90000,
          "retry_on_error": true
        },
        {
          "provider": "nvidia",
          "model_id": "meta/llama-3.3-70b-instruct",
          "priority": 3,
          "timeout_ms": 90000,
          "retry_on_error": true
        },
        {
          "provider": "opencode-zen",
          "model_id": "grok-code",
          "priority": 4,
          "timeout_ms": 60000,
          "retry_on_error": false
        }
      ],
      "continue_on_429": true,
      "wait_time_on_429_ms": 60000
    },
    "fast": {
      "chain": [
        {
          "provider": "opencode-zen",
          "model_id": "minimax-m2.5-free",
          "priority": 1,
          "timeout_ms": 30000,
          "retry_on_error": false
        }
      ]
    }
  }
}
```

### 2.4 Agent Swarms

Agent swarms enable parallel execution of multiple agents for complex tasks. This is a powerful pattern for achieving results faster by distributing work across multiple specialized agents.

#### Swarm Configuration

```json
{
  "swarms": {
    "default_minimum": 5,
    "execution_mode": "parallel",
    "result_aggregation": "merge",
    "presets": {
      "code_review_swarm": {
        "agents": [
          {
            "preset": "oracle",
            "role": "primary_reviewer"
          },
          {
            "preset": "sisyphus",
            "role": "implementation_checker"
          },
          {
            "preset": "librarian",
            "role": "documentation_verifier"
          },
          {
            "preset": "explore",
            "role": "pattern_finder"
          },
          {
            "preset": "metis",
            "role": "security_analyst"
          }
        ],
        "coordination": {
          "strategy": "round_robin",
          "max_iterations": 3,
          "consensus_threshold": 0.7
        }
      },
      "research_swarm": {
        "agents": [
          {
            "preset": "metis",
            "role": "primary_researcher"
          },
          {
            "preset": "librarian",
            "role": "fact_checker"
          },
          {
            "preset": "explore",
            "role": "source_finder"
          },
          {
            "preset": "atlas",
            "role": "data_aggregator"
          },
          {
            "preset": "momus",
            "role": "report_writer"
          }
        ],
        "coordination": {
          "strategy": "hierarchical",
          "max_iterations": 5,
          "consensus_threshold": 0.8
        }
      }
    }
  }
}
```

---

## 3. MCP Server Setup

### 3.1 Local MCPs

Local MCPs run on the same machine as OpenClaw and communicate via stdio. These are typically lighter-weight integrations that don't require containerization.

#### Installing Local MCPs

```bash
# Install Serena (Orchestration MCP)
uvx serena start-mcp-server

# Install Tavily (Web Search MCP)
npx @tavily/claude-mcp

# Install Canva (Design MCP)
npx @canva/claude-mcp

# Install Context7 (Documentation MCP)
npx @anthropics/context7-mcp

# Install Chrome DevTools MCP
npx @anthropics/chrome-devtools-mcp
```

#### Local MCP Configuration

```json
{
  "mcp": {
    "servers": {
      "serena": {
        "type": "local",
        "command": "uvx",
        "args": ["serena", "start-mcp-server"],
        "enabled": true,
        "environment": {}
      },
      "tavily": {
        "type": "local",
        "command": "npx",
        "args": ["@tavily/claude-mcp"],
        "enabled": true,
        "environment": {
          "TAVILY_API_KEY": "${TAVILY_API_KEY}"
        }
      },
      "canva": {
        "type": "local",
        "command": "npx",
        "args": ["@canva/claude-mcp"],
        "enabled": true,
        "environment": {
          "CANVA_API_KEY": "${CANVA_API_KEY}"
        }
      },
      "context7": {
        "type": "local",
        "command": "npx",
        "args": ["@anthropics/context7-mcp"],
        "enabled": true,
        "environment": {}
      }
    }
  }
}
```

### 3.2 Remote MCPs (Docker-based)

Remote MCPs run in Docker containers and communicate via HTTP. This is the recommended approach for more complex integrations that require isolated environments.

#### Docker MCP Setup

```bash
# Pull MCP Docker images
docker pull anomalyco/serena-mcp:latest
docker pull anomalyco/skyvern-mcp:latest

# Run MCP container
docker run -d \
  --name serena-mcp \
  -p 8000:8000 \
  -e MCP_API_KEY="your-api-key" \
  anomalyco/serena-mcp:latest

# Verify container is running
docker ps | grep serena-mcp
```

#### Remote MCP Configuration

```json
{
  "mcp": {
    "servers": {
      "serena-remote": {
        "type": "remote",
        "url": "http://localhost:8000",
        "enabled": true,
        "timeout_ms": 30000,
        "retry_attempts": 3
      },
      "skyvern": {
        "type": "remote",
        "url": "http://agent-06-skyvern-solver:8030",
        "enabled": true,
        "timeout_ms": 60000,
        "retry_attempts": 2
      },
      "linear": {
        "type": "remote",
        "url": "https://mcp.linear.app/sse",
        "enabled": true,
        "timeout_ms": 15000,
        "auth": {
          "type": "bearer",
          "token": "${LINEAR_API_KEY}"
        }
      }
    }
  }
}
```

### 3.3 MCP Wrapper Pattern

The MCP Wrapper Pattern is used when Docker containers expose HTTP APIs but OpenClaw expects stdio communication. This is essential for integrating with containerized services.

#### Wrapper Implementation

```javascript
#!/usr/bin/env node
/**
 * MCP Wrapper for Docker HTTP Services
 * Converts stdio MCP protocol to HTTP API calls
 */

const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
const { StdioServerTransport } = require('@modelcontextprotocol/sdk/server/stdio.js');
const axios = require('axios');

const API_URL = process.env.API_URL || 'http://localhost:8000';
const API_KEY = process.env.API_KEY;

const server = new Server(
  { name: 'http-service-mcp', version: '1.0.0' },
  { capabilities: { tools: {} } }
);

// Example tool: Execute HTTP request
async function httpRequestTool(params) {
  const { method, endpoint, data } = params;
  const url = `${API_URL}${endpoint}`;
  
  const response = await axios({
    method,
    url,
    data,
    headers: {
      'Authorization': `Bearer ${API_KEY}`,
      'Content-Type': 'application/json'
    }
  });
  
  return response.data;
}

// Example tool: Query service status
async function healthCheck() {
  try {
    const response = await axios.get(`${API_URL}/health`);
    return { status: 'healthy', ...response.data };
  } catch (error) {
    return { status: 'unhealthy', error: error.message };
  }
}

server.setRequestHandler(ListToolsRequestSchema, async () => ({
  tools: [
    {
      name: 'http_request',
      description: 'Make HTTP request to the service',
      inputSchema: {
        type: 'object',
        properties: {
          method: { type: 'string', enum: ['GET', 'POST', 'PUT', 'DELETE'] },
          endpoint: { type: 'string' },
          data: { type: 'object' }
        },
        required: ['method', 'endpoint']
      }
    },
    {
      name: 'health_check',
      description: 'Check service health status',
      inputSchema: {
        type: 'object',
        properties: {}
      }
    }
  ]
}));

server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args } = request.params;
  
  try {
    switch (name) {
      case 'http_request':
        return { content: [{ type: 'text', text: JSON.stringify(await httpRequestTool(args)) }] };
      case 'health_check':
        return { content: [{ type: 'text', text: JSON.stringify(await healthCheck()) }] };
      default:
        throw new Error(`Unknown tool: ${name}`);
    }
  } catch (error) {
    return { content: [{ type: 'text', text: `Error: ${error.message}` }], isError: true };
  }
});

const transport = new StdioServerTransport();
server.connect(transport).catch(console.error);
```

#### Wrapper Configuration Example

```json
{
  "mcp": {
    "servers": {
      "custom-service-mcp": {
        "type": "local",
        "command": ["node", "/path/to/wrapper.js"],
        "enabled": true,
        "environment": {
          "API_URL": "http://custom-service:8080",
          "API_KEY": "${CUSTOM_SERVICE_API_KEY}"
        }
      }
    }
  }
}
```

### 3.4 Configuration Examples

#### Complete MCP Configuration

```json
{
  "mcp": {
    "global_settings": {
      "timeout_ms": 30000,
      "max_retries": 3,
      "retry_delay_ms": 1000,
      "connection_pool_size": 10
    },
    "servers": {
      "serena": {
        "type": "local",
        "command": "uvx",
        "args": ["serena", "start-mcp-server"],
        "enabled": true,
        "environment": {},
        "auto_restart": true,
        "health_check_interval_ms": 60000
      },
      "tavily": {
        "type": "local",
        "command": "npx",
        "args": ["@tavily/claude-mcp"],
        "enabled": true,
        "environment": {
          "TAVILY_API_KEY": "${TAVILY_API_KEY}"
        }
      },
      "canva": {
        "type": "local",
        "command": "npx",
        "args": ["@canva/claude-mcp"],
        "enabled": true,
        "environment": {
          "CANVA_API_KEY": "${CANVA_API_KEY}",
          "CANVA_CLIENT_ID": "${CANVA_CLIENT_ID}",
          "CANVA_CLIENT_SECRET": "${CANVA_CLIENT_SECRET}"
        }
      },
      "context7": {
        "type": "local",
        "command": "npx",
        "args": ["@anthropics/context7-mcp"],
        "enabled": true,
        "environment": {}
      },
      "chrome-devtools": {
        "type": "local",
        "command": "npx",
        "args": ["@anthropics/chrome-devtools-mcp"],
        "enabled": true,
        "environment": {
          "CDP_ENDPOINT": "ws://localhost:9222"
        }
      },
      "linear": {
        "type": "remote",
        "url": "https://mcp.linear.app/sse",
        "enabled": true,
        "timeout_ms": 15000,
        "auth": {
          "type": "bearer",
          "token": "${LINEAR_API_KEY}"
        }
      },
      "gh_grep": {
        "type": "remote",
        "url": "https://mcp.grep.app",
        "enabled": true,
        "timeout_ms": 20000
      }
    }
  }
}
```

---

## 4. Workflow Automation

### 4.1 Task Delegation

Task delegation is the core mechanism for distributing work to agents. OpenClaw supports both synchronous and asynchronous delegation patterns.

#### Basic Task Delegation

```bash
# Delegate a simple task to sisyphus
openclaw delegate --agent sisyphus --task "Fix the login bug in auth.ts"

# Delegate with specific model
openclaw delegate --agent sisyphus --model nvidia/qwen/qwen3.5-397b-a17b --task "Write tests for user service"

# Delegate with priority
openclaw delegate --agent prometheus --priority high --task "Design the new API architecture"
```

#### Programmatic Task Delegation

```typescript
import { OpenClawClient } from '@openclaw/sdk';

const client = new OpenClawClient({
  gatewayUrl: 'http://localhost:18789'
});

// Simple delegation
const result = await client.delegate({
  agent: 'sisyphus',
  task: 'Implement user authentication',
  context: {
    files: ['/path/to/auth/module.ts'],
    requirements: ['JWT tokens', 'OAuth2 support']
  }
});

// Delegation with callbacks
await client.delegate({
  agent: 'prometheus',
  task: 'Create project roadmap',
  context: {
    project_name: 'BIOMETRICS',
    scope: 'Phase 2.6'
  },
  onProgress: (progress) => {
    console.log(`Progress: ${progress.percentage}%`);
  },
  onComplete: (result) => {
    console.log(`Result: ${result.output}`);
  }
});
```

### 4.2 Swarm Mode (5+ Agents)

Swarm mode enables parallel execution of multiple agents for complex tasks. This is essential for achieving results faster and more reliably.

#### Starting a Swarm

```bash
# Start a code review swarm
openclaw swarm start --preset code_review_swarm --task "Review the entire auth module"

# Start a research swarm
openclaw swarm start --preset research_swarm --task "Research best practices for biometric authentication"

# Custom swarm with specific agents
openclaw swarm start \
  --agents sisyphus,oracle,librarian,explore,metis \
  --task "Implement and document the new payment integration"
```

#### Swarm Configuration

```json
{
  "swarm": {
    "default_minimum": 5,
    "max_parallel": 10,
    "timeout_ms": 300000,
    "aggregation": {
      "strategy": "weighted_vote",
      "weights": {
        "oracle": 0.3,
        "sisyphus": 0.25,
        "librarian": 0.2,
        "explore": 0.15,
        "metis": 0.1
      }
    }
  }
}
```

### 4.3 Orchestration Patterns

OpenClaw supports several orchestration patterns for different workflow scenarios.

#### Sequential Orchestration

```typescript
// Sequential: Each agent waits for the previous to complete
async function sequentialWorkflow(task: string, agents: string[]) {
  let context = {};
  
  for (const agent of agents) {
    const result = await openclaw.delegate({
      agent,
      task,
      context
    });
    
    // Pass output as context to next agent
    context = { ...context, ...result };
  }
  
  return context;
}

// Example: Documentation generation
await sequentialWorkflow(
  'Generate API documentation',
  ['explore', 'librarian', 'momus']
);
```

#### Parallel Orchestration

```typescript
// Parallel: All agents work simultaneously
async function parallelWorkflow(task: string, agents: string[]) {
  const results = await Promise.all(
    agents.map(agent => 
      openclaw.delegate({ agent, task })
    )
  );
  
  // Aggregate results
  return aggregateResults(results);
}

// Example: Multiple perspectives on code
await parallelWorkflow(
  'Review the authentication module',
  ['oracle', 'sisyphus', 'metis', 'explore', 'librarian']
);
```

#### Hierarchical Orchestration

```typescript
// Hierarchical: Supervisor delegates to sub-agents
async function hierarchicalWorkflow(task: string) {
  // Phase 1: Planning
  const plan = await openclaw.delegate({
    agent: 'prometheus',
    task: `Create implementation plan for: ${task}`
  });
  
  // Phase 2: Parallel implementation
  const implementations = await Promise.all(
    plan.phases.map(phase =>
      openclaw.delegate({
        agent: 'sisyphus',
        task: phase.description,
        context: { phase }
      })
    )
  );
  
  // Phase 3: Aggregation
  return openclaw.delegate({
    agent: 'oracle',
    task: 'Synthesize implementation results',
    context: { implementations }
  });
}
```

### 4.4 Result Collection

Result collection handles aggregating outputs from multiple agents into coherent responses.

#### Result Aggregation Strategies

```json
{
  "aggregation": {
    "strategies": {
      "merge": {
        "description": "Combine all outputs into single response",
        "priority": "all",
        "deduplication": true
      },
      "weighted_vote": {
        "description": "Weight votes by agent expertise",
        "weights": {
          "oracle": 0.3,
          "sisyphus": 0.25,
          "metis": 0.2,
          "librarian": 0.15,
          "explore": 0.1
        },
        "threshold": 0.7
      },
      "first_valid": {
        "description": "Return first non-error result",
        "priority": "order"
      },
      "consensus": {
        "description": "Only return if agents agree",
        "threshold": 0.8,
        "require_all": false
      }
    },
    "default_strategy": "weighted_vote"
  }
}
```

---

## 5. Configuration

### 5.1 OpenClaw Configuration Structure

The main configuration file is located at `~/.openclaw/openclaw.json`. This file contains all settings for providers, agents, MCP servers, and runtime behavior.

#### Complete Configuration Schema

```json
{
  "version": "1.0",
  "env": {
    "NVIDIA_API_KEY": "nvapi-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    "OPENROUTER_API_KEY": "sk-or-v1-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    "ANTHROPIC_API_KEY": "sk-ant-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    "TAVILY_API_KEY": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    "LINEAR_API_KEY": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  },
  "models": {
    "providers": {
      "nvidia": {
        "baseUrl": "https://integrate.api.nvidia.com/v1",
        "api": "openai-completions",
        "models": [
          "qwen/qwen3.5-397b-a17b",
          "moonshotai/kimi-k2.5",
          "meta/llama-3.3-70b-instruct",
          "mistralai/mistral-large-3-675b-instruct-2512"
        ]
      },
      "opencode-zen": {
        "baseUrl": "https://api.opencode.ai/v1",
        "api": "openai-completions",
        "models": [
          "zen/big-pickle",
          "zen/uncensored",
          "grok-code",
          "kimi-k2.5-free",
          "minimax-m2.5-free"
        ]
      }
    }
  },
  "agents": {
    "defaults": {
      "model": {
        "primary": "nvidia/moonshotai/kimi-k2.5",
        "fallbacks": [
          "nvidia/meta/llama-3.3-70b-instruct",
          "nvidia/mistralai/mistral-large-3-675b-instruct-2512",
          "opencode-zen/grok-code"
        ]
      },
      "timeout_ms": 120000,
      "max_retries": 3
    },
    "presets": { }
  },
  "mcp": {
    "servers": { }
  },
  "gateway": {
    "host": "0.0.0.0",
    "port": 18789,
    "max_connections": 100,
    "timeout_ms": 300000
  },
  "logging": {
    "level": "info",
    "format": "json",
    "output": "stdout"
  }
}
```

### 5.2 Environment Variables

Environment variables provide a secure way to configure API keys and sensitive credentials without storing them in configuration files.

#### Setting Environment Variables

```bash
# For bash/zsh
export NVIDIA_API_KEY="nvapi-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
export OPENROUTER_API_KEY="sk-or-v1-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

# For fish
set -x NVIDIA_API_KEY "nvapi-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

# For Windows (PowerShell)
$env:NVIDIA_API_KEY = "nvapi-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

#### Environment Variable Precedence

1. **Runtime environment variables** (highest priority)
2. **Configuration file env section**
3. **Default values**

### 5.3 Provider Configuration

Providers define how OpenClaw connects to different AI model backends.

#### NVIDIA NIM Provider

```json
{
  "models": {
    "providers": {
      "nvidia": {
        "baseUrl": "https://integrate.api.nvidia.com/v1",
        "api": "openai-completions",
        "timeout_ms": 120000,
        "rate_limit": {
          "requests_per_minute": 40,
          "burst": 10
        },
        "models": [
          "qwen/qwen3.5-397b-a17b",
          "moonshotai/kimi-k2.5",
          "meta/llama-3.3-70b-instruct",
          "mistralai/mistral-large-3-675b-instruct-2512"
        ]
      }
    }
  }
}
```

#### OpenCode Zen Provider

```json
{
  "models": {
    "providers": {
      "opencode-zen": {
        "baseUrl": "https://api.opencode.ai/v1",
        "api": "openai-completions",
        "timeout_ms": 60000,
        "models": [
          "zen/big-pickle",
          "zen/uncensored",
          "grok-code",
          "kimi-k2.5-free",
          "minimax-m2.5-free"
        ]
      }
    }
  }
}
```

### 5.4 Agent Defaults

Agent defaults provide fallback settings when specific configurations are not provided.

```json
{
  "agents": {
    "defaults": {
      "model": {
        "primary": "nvidia/moonshotai/kimi-k2.5",
        "fallbacks": [
          "nvidia/meta/llama-3.3-70b-instruct",
          "nvidia/mistralai/mistral-large-3-675b-instruct-2512",
          "opencode-zen/grok-code"
        ]
      },
      "temperature": 0.7,
      "max_tokens": 32768,
      "timeout_ms": 120000,
      "max_retries": 3,
      "retry_delay_ms": 1000
    },
    "behavior": {
      "auto_fallback": true,
      "continue_on_429": true,
      "wait_time_on_429_ms": 60000,
      "parallel_execution": true,
      "max_parallel_agents": 5
    }
  }
}
```

---

## 6. Gateway Management

### 6.1 Gateway Start/Stop/Restart

The OpenClaw gateway is the central service that handles agent communication and orchestration.

#### Starting the Gateway

```bash
# Start gateway in foreground (for debugging)
openclaw gateway start

# Start gateway in background (production)
openclaw gateway start --background

# Start with custom port
openclaw gateway start --port 18789 --host 0.0.0.0

# Start with debug logging
openclaw gateway start --log-level debug
```

#### Stopping the Gateway

```bash
# Stop gateway gracefully
openclaw gateway stop

# Force stop (kills process)
openclaw gateway stop --force

# Stop and remove data
openclaw gateway stop --clean
```

#### Restarting the Gateway

```bash
# Restart gateway
openclaw gateway restart

# Restart with configuration reload
openclaw gateway restart --reload-config

# Restart and clear cache
openclaw gateway restart --clear-cache
```

### 6.2 Health Checks

Regular health checks ensure the gateway and its dependencies are functioning correctly.

#### Manual Health Check

```bash
# Run health check
openclaw doctor

# Run health check with auto-fix
openclaw doctor --fix

# Run specific health checks
openclaw doctor --check gateway
openclaw doctor --check models
openclaw doctor --check mcp
openclaw doctor --check network
```

#### Automated Health Monitoring

```json
{
  "health": {
    "enabled": true,
    "interval_ms": 60000,
    "checks": {
      "gateway": {
        "endpoint": "/health",
        "timeout_ms": 5000
      },
      "models": {
        "enabled": true,
        "test_model": "nvidia/moonshotai/kimi-k2.5"
      },
      "mcp": {
        "enabled": true,
        "servers": ["serena", "tavily"]
      },
      "network": {
        "enabled": true,
        "targets": [
          "https://integrate.api.nvidia.com/v1",
          "https://api.opencode.ai/v1"
        ]
      }
    },
    "notifications": {
      "on_failure": {
        "email": "admin@example.com",
        "slack_webhook": "https://hooks.slack.com/..."
      }
    }
  }
}
```

### 6.3 Doctor Command

The `doctor` command is a comprehensive diagnostic tool that checks and repairs OpenClaw installation.

#### Doctor Command Options

```bash
# Full diagnostic
openclaw doctor

# Diagnostic with auto-repair
openclaw doctor --fix

# Verbose output
openclaw doctor --verbose

# Specific category
openclaw doctor --category config
openclaw doctor --category network
openclaw doctor --category models
openclaw doctor --category mcp
```

#### Common Doctor Repairs

```bash
# Repair configuration permissions
openclaw doctor --fix-permissions

# Repair network configuration
openclaw doctor --fix-network

# Repair model cache
openclaw doctor --fix-cache

# Reset to defaults
openclaw doctor --reset
```

### 6.4 Logs and Debugging

Logs provide detailed information about OpenClaw's operation for troubleshooting and monitoring.

#### Viewing Logs

```bash
# View gateway logs
openclaw logs

# View with specific level
openclaw logs --level debug
openclaw logs --level info
openclaw logs --level warn
openclaw logs --level error

# Follow logs in real-time
openclaw logs --follow

# View specific number of lines
openclaw logs --lines 100

# Filter logs
openclaw logs --filter "error"
openclaw logs --filter "agent=sisyphus"
```

#### Debugging Options

```bash
# Enable debug mode
openclaw gateway start --debug

# Trace specific agent
openclaw gateway start --trace-agent sisyphus

# Trace specific MCP
openclaw gateway start --trace-mcp serena

# Profile performance
openclaw gateway start --profile
```

---

## 7. Best Practices

### 7.1 Agent Swarms (Minimum 5 Agents)

For complex tasks, always use at least 5 agents in parallel to leverage diverse perspectives and expertise.

#### Swarm Best Practices

```typescript
// Always use minimum 5 agents for complex tasks
const CODE_REVIEW_SWARM = [
  'oracle',      // Primary reviewer (architecture + best practices)
  'sisyphus',    // Implementation checker
  'librarian',   // Documentation verifier
  'explore',     // Pattern finder
  'metis',       // Security analyst
];

const RESEARCH_SWARM = [
  'metis',       // Primary researcher
  'librarian',   // Fact checker
  'explore',     // Source finder
  'atlas',       // Data aggregator
  'momus',       // Report writer
];

// Execute swarm with proper configuration
await openclaw.swarm.start({
  agents: CODE_REVIEW_SWARM,
  task: 'Comprehensive code review',
  coordination: {
    strategy: 'parallel',
    consensus_required: true,
    threshold: 0.7
  }
});
```

### 7.2 Parallel Execution

Always prefer parallel execution over sequential when agents don't depend on each other's output.

#### Parallel Execution Rules

```typescript
// Good: Independent tasks run in parallel
const results = await Promise.all([
  openclaw.delegate({ agent: 'explore', task: 'Find all auth files' }),
  openclaw.delegate({ agent: 'librarian', task: 'Check auth documentation' }),
  openclaw.delegate({ agent: 'atlas', task: 'Count auth-related code' }),
]);

// Bad: Sequential when parallel is possible
const result1 = await openclaw.delegate({ agent: 'explore', task: 'Find auth files' });
const result2 = await openclaw.delegate({ agent: 'librarian', task: 'Check docs' });
const result3 = await openclaw.delegate({ agent: 'atlas', task: 'Count code' });
```

### 7.3 Resource Management

Proper resource management ensures optimal performance and cost efficiency.

#### Resource Management Guidelines

```json
{
  "resource_management": {
    "max_concurrent_agents": 10,
    "max_parallel_mcp_calls": 5,
    "timeout_defaults": {
      "fast_task": 30000,
      "normal_task": 60000,
      "complex_task": 120000
    },
    "model_selection": {
      "prefer_free": true,
      "fallback_to_paid": true,
      "cost_optimization": {
        "enabled": true,
        "max_cost_per_request": 0.01
      }
    },
    "caching": {
      "enabled": true,
      "ttl_seconds": 3600,
      "max_size_mb": 500
    }
  }
}
```

### 7.4 Cost Optimization (FREE-first)

Always prefer free models and services to minimize costs while maintaining quality.

#### Cost Optimization Strategy

```typescript
// Priority: Free models first
const MODEL_PRIORITY = [
  // Tier 1: Completely Free
  'opencode-zen/minimax-m2.5-free',    // $0.00
  'opencode-zen/kimi-k2.5-free',        // $0.00
  
  // Tier 2: Free Tier (with limits)
  'nvidia/qwen/qwen3.5-397b-a17b',     // 40 RPM free
  'nvidia/moonshotai/kimi-k2.5',        // Free tier available
  
  // Tier 3: Paid (fallback only)
  'nvidia/meta/llama-3.3-70b-instruct', // Pay per use
];

// Implementation
async function getOptimalModel(taskType: string): Promise<string> {
  const cache = await getCachedModel(taskType);
  if (cache) return cache;
  
  for (const model of MODEL_PRIORITY) {
    if (await isModelAvailable(model)) {
      await cacheModel(taskType, model);
      return model;
    }
  }
  
  throw new Error('No models available');
}
```

---

## 8. Troubleshooting

### 8.1 Common Errors

This section covers the most common errors encountered when using OpenClaw and their solutions.

#### Error: "Connection refused to gateway"

**Cause:** Gateway is not running or not accessible.

**Solution:**
```bash
# Check if gateway is running
openclaw gateway status

# Start gateway if not running
openclaw gateway start

# Check port availability
lsof -i :18789

# Verify firewall rules
firewall-cmd --list-ports
```

#### Error: "Model not found"

**Cause:** Model ID is incorrect or model is not available.

**Solution:**
```bash
# List available models
openclaw models

# Verify model ID in configuration
cat ~/.openclaw/openclaw.json | grep model_id

# Check provider status
openclaw doctor --check models
```

#### Error: "MCP server connection timeout"

**Cause:** MCP server is not responding or network issues.

**Solution:**
```bash
# Check MCP server status
openclaw mcp list

# Restart MCP server
openclaw mcp restart serena

# Check MCP server logs
openclaw logs --filter "mcp"
```

#### Error: "Rate limit exceeded (HTTP 429)"

**Cause:** Too many requests to the model provider.

**Solution:**
```bash
# Wait 60 seconds for rate limit reset
sleep 60

# Use fallback model
openclaw delegate --model <fallback-model> --task "..."

# Check current rate limit status
openclaw doctor --check network
```

### 8.2 Gateway Issues

#### Gateway Won't Start

```bash
# Check configuration validity
openclaw config validate

# Check port availability
netstat -an | grep 18789

# Check logs for errors
openclaw logs --level error

# Try starting in foreground to see errors
openclaw gateway start --foreground
```

#### Gateway Performance Issues

```bash
# Check resource usage
openclaw gateway status

# Enable performance profiling
openclaw gateway start --profile

# Check concurrent connections
openclaw gateway status --connections

# Adjust max connections in config
openclaw config set gateway.max_connections 100
```

### 8.3 MCP Connection Problems

#### MCP Server Not Responding

```bash
# List all MCP servers
openclaw mcp list

# Test specific MCP server
openclaw mcp test serena

# Restart MCP server
openclaw mcp restart serena

# Check MCP server logs
openclaw logs --filter "mcp=serena"
```

#### MCP Authentication Errors

```bash
# Verify environment variables
env | grep API_KEY

# Update MCP configuration
openclaw config set mcp.servers.serena.environment.API_KEY "new-key"

# Test authentication
openclaw mcp test serena --auth
```

### 8.4 Agent Failures

#### Agent Returns Error

```bash
# Get detailed error information
openclaw delegate --agent sisyphus --task "test" --verbose

# Check agent logs
openclaw logs --filter "agent=sisyphus"

# Verify model availability
openclaw models | grep <model-id>

# Try with different agent preset
openclaw delegate --agent atlas --task "test"
```

#### Agent Timeout

```bash
# Increase timeout
openclaw delegate --timeout 180000 --task "complex task"

# Use faster model
openclaw delegate --model opencode-zen/minimax-m2.5-free --task "task"

# Check network latency
openclaw doctor --check network

# Break task into smaller pieces
```

---

## Appendix A: Command Reference

### Core Commands

| Command | Description |
|---------|-------------|
| `openclaw --version` | Show version |
| `openclaw --help` | Show help |
| `openclaw init` | Initialize configuration |
| `openclaw config validate` | Validate configuration |
| `openclaw models` | List available models |
| `openclaw agents list` | List configured agents |

### Gateway Commands

| Command | Description |
|---------|-------------|
| `openclaw gateway start` | Start gateway |
| `openclaw gateway stop` | Stop gateway |
| `openclaw gateway restart` | Restart gateway |
| `openclaw gateway status` | Show gateway status |
| `openclaw doctor` | Run diagnostics |

### Agent Commands

| Command | Description |
|---------|-------------|
| `openclaw delegate` | Delegate task to agent |
| `openclaw swarm start` | Start agent swarm |
| `openclaw swarm status` | Show swarm status |
| `openclaw agents list` | List agents |

### MCP Commands

| Command | Description |
|---------|-------------|
| `openclaw mcp list` | List MCP servers |
| `openclaw mcp test` | Test MCP server |
| `openclaw mcp restart` | Restart MCP server |

---

## Appendix B: Configuration Examples

### Minimal Configuration

```json
{
  "env": {
    "NVIDIA_API_KEY": "nvapi-xxx"
  },
  "models": {
    "providers": {
      "nvidia": {
        "baseUrl": "https://integrate.api.nvidia.com/v1",
        "api": "openai-completions"
      }
    }
  },
  "agents": {
    "defaults": {
      "model": {
        "primary": "nvidia/moonshotai/kimi-k2.5"
      }
    }
  }
}
```

### Production Configuration

See Section 5.1 for the complete production configuration.

---

## Appendix C: Quick Reference

### Key Commands

```bash
# Setup
openclaw init
openclaw doctor --fix

# Models
openclaw models | grep nvidia

# Gateway
openclaw gateway restart
openclaw doctor

# Agents
openclaw agents list
openclaw delegate --agent sisyphus --task "Hello"

# MCP
openclaw mcp list
openclaw mcp test serena
```

---

**Document Version:** 1.0  
**Created:** 2026-02-20  
**Last Updated:** 2026-02-20  
**Status:** Active

---

*This document is part of the BIOMETRICS project rules documentation.*
