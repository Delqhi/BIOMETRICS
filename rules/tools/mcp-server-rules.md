# MCP Server Rules - Model Context Protocol Integration Guide

**Version:** 1.0  
**Last Updated:** 2026-02-20  
**Status:** ACTIVE - MANDATORY COMPLIANCE

---

## Table of Contents

1. [MCP Architecture](#1-mcp-architecture)
2. [MCP Wrapper Pattern](#2-mcp-wrapper-pattern)
3. [Local MCP Servers](#3-local-mcp-servers)
4. [Remote MCP Servers](#4-remote-mcp-servers)
5. [MCP Configuration](#5-mcp-configuration)
6. [Creating New MCPs](#6-creating-new-mcps)
7. [Best Practices](#7-best-practices)
8. [Troubleshooting](#8-troubleshooting)

---

## 1. MCP Architecture

### 1.1 What is MCP (Model Context Protocol)

The Model Context Protocol (MCP) is an open-standard protocol that enables AI models to interact with external tools, services, and data sources in a standardized manner. MCP serves as a bridge between AI assistants and the real world, allowing models to execute actions, retrieve information, and maintain state across interactions.

MCP was designed to solve a fundamental problem in AI development: the fragmentation of tool integrations. Before MCP, each AI model required custom implementations for interacting with different services, databases, and APIs. MCP standardizes these interactions, making it possible to:

- Connect AI models to any data source without custom code
- Share tool definitions across different AI platforms and frameworks
- Maintain consistent authentication and security patterns
- Enable seamless tool discovery and composition
- Support both synchronous and asynchronous tool execution

The protocol defines three core concepts that form the foundation of all MCP integrations:

**Resources** represent data that can be read or written through MCP. Resources are identified by URIs and can be files, database records, API responses, or any other form of structured data. Resources enable AI models to access external information in a controlled, auditable manner.

**Tools** are executable functions that AI models can invoke. Each tool has a defined schema specifying its name, description, and input parameters. Tools can perform actions like querying databases, calling APIs, executing code, or manipulating files. The tool definition includes comprehensive parameter descriptions that help AI models understand when and how to use each tool.

**Prompts** are reusable prompt templates that can be configured within MCP. Prompts allow teams to standardize common AI interactions and ensure consistent behavior across different use cases. Prompts can include variable placeholders that get filled at runtime.

### 1.2 Local vs Remote MCPs

MCP servers can operate in two primary modes: local and remote. Understanding the distinction between these modes is crucial for designing effective AI integrations.

**Local MCP Servers** run on the same machine as the AI application or are accessed through local network connections. They communicate through standard input and output (stdio) streams, making them ideal for tools that require direct system access or need to operate without network latency. Local MCPs are typically used for:

- File system operations
- Local development tools
- System command execution
- Database connections on localhost
- Integration with locally installed applications

Local MCPs offer several advantages including minimal latency, no network dependency, simpler security configuration, and direct access to local resources. However, they require the MCP server process to be running on the same machine as the AI application, which can limit scalability and distribution.

**Remote MCP Servers** operate over network connections, typically using HTTP or WebSocket protocols. They enable AI models to access services and tools hosted on remote servers, in containers, or in cloud environments. Remote MCPs are essential for:

- Cloud-based services and APIs
- Distributed microservices architectures
- Shared infrastructure resources
- Services requiring centralized management
- High-availability deployments

Remote MCPs provide greater flexibility in terms of deployment and scaling but introduce network latency and require proper network security configuration. They also need authentication mechanisms to ensure secure access to remote resources.

The choice between local and remote MCPs depends on specific requirements:

| Aspect | Local MCP | Remote MCP |
|--------|-----------|------------|
| Latency | Minimal (<1ms) | Network dependent (10-100ms+) |
| Deployment | Process-based | Container/Service-based |
| Scaling | Limited | Horizontal scaling possible |
| Security | Local access controls | Network + auth required |
| Maintenance | Per-machine setup | Centralized management |

### 1.3 Docker-based MCPs

Docker containers provide an excellent foundation for deploying MCP servers in production environments. Containerization offers consistent runtime environments, easy scaling, and simplified deployment workflows. When deploying MCPs as Docker containers, several architectural patterns apply.

**Containerized MCP Architecture** follows a microservices pattern where each MCP server runs in its own container. This isolation ensures that dependencies don't conflict, enables independent scaling, and simplifies maintenance. Common configurations include:

```yaml
# docker-compose.yml example for MCP server
version: '3.8'
services:
  mcp-server:
    build: ./mcp-server
    ports:
      - "50001:50001"
    environment:
      - API_KEY=${MCP_API_KEY}
      - DATABASE_URL=${DATABASE_URL}
      - LOG_LEVEL=info
    volumes:
      - ./config:/app/config:ro
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:50001/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

**Multi-Container MCP Setups** are common when MCP servers require supporting services like databases, caches, or message queues. For example, an MCP server that provides vector search capabilities might require both a PostgreSQL database and a Redis cache:

```yaml
version: '3.8'
services:
  vector-mcp:
    build: ./vector-mcp
    depends_on:
      - postgres
      - redis
    environment:
      - DATABASE_URL=postgresql://user:pass@postgres:5432/vectors
      - REDIS_URL=redis://redis:6379
    networks:
      - mcp-network

  postgres:
    image: pgvector/pgvector:pg16
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
      - POSTGRES_DB=vectors
    volumes:
      - postgres-data:/var/lib/postgresql/data
    networks:
      - mcp-network

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis-data:/data
    networks:
      - mcp-network

networks:
  mcp-network:
    driver: bridge

volumes:
  postgres-data:
  redis-data:
```

**Container Networking** is critical for Docker-based MCPs. Containers communicate through a shared Docker network, which provides DNS-based service discovery. Each container can reach other containers by their service name. For external access, ports are mapped to the host machine using the `ports` configuration.

### 1.4 stdio vs HTTP Communication

MCP supports two primary communication patterns: stdio and HTTP. Each pattern has specific use cases, advantages, and implementation requirements.

**stdio Communication** is the traditional method for local MCP servers. In this pattern, the MCP client spawns the server as a subprocess, and all communication occurs through standard input and output streams. Messages are JSON-RPC formatted and sent line-by-line:

```
Client -> stdin -> MCP Server
Client <- stdout <- MCP Server
```

The stdio pattern offers several advantages for local operations:

- Simple implementation without network setup
- No port conflicts or network security concerns
- Automatic process lifecycle management
- Minimal resource overhead
- Direct access to local resources

However, stdio has limitations for remote or distributed scenarios:

- No direct network accessibility
- Limited concurrent connection support
- Difficult to scale horizontally
- No built-in load balancing

**HTTP Communication** uses REST or WebSocket APIs for client-server interaction. This pattern is essential for remote MCPs and enables sophisticated deployment scenarios:

```typescript
// HTTP MCP Server Example (Express.js)
import express from 'express';
import { Server } from '@modelcontextprotocol/sdk/server/index.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';

const app = express();
app.use(express.json());

// MCP server instance
const mcpServer = new Server(
  { name: 'http-mcp-server', version: '1.0.0' },
  { capabilities: { tools: {} } }
);

// HTTP endpoint for tool invocation
app.post('/mcp/tools/call', async (req, res) => {
  const { name, arguments: args } = req.body;
  
  try {
    const result = await mcpServer.handleToolCall(name, args);
    res.json({ success: true, result });
  } catch (error) {
    res.status(500).json({ success: false, error: error.message });
  }
});

// HTTP endpoint for tool listing
app.get('/mcp/tools', async (req, res) => {
  const tools = await mcpServer.getTools();
  res.json({ tools });
});

// Health check endpoint
app.get('/health', (req, res) => {
  res.json({ status: 'healthy' });
});

app.listen(50001, () => {
  console.log('MCP Server running on port 50001');
});
```

**WebSocket Communication** provides real-time, bidirectional communication for MCP servers that require persistent connections or streaming data. This pattern is particularly useful for:

- Real-time data streaming
- Long-running operations with progress updates
- Interactive tool execution
- Collaborative sessions

```typescript
// WebSocket MCP Server Example
import { WebSocketServer } from 'ws';
import { Server } from '@modelcontextprotocol/sdk/server/index.js';

const wss = new WebSocketServer({ port: 50002 });
const mcpServer = new Server(
  { name: 'websocket-mcp-server', version: '1.0.0' },
  { capabilities: { tools: {}, resources: {} } }
);

wss.on('connection', (ws) => {
  ws.on('message', async (message) => {
    const data = JSON.parse(message.toString());
    
    switch (data.type) {
      case 'tool_call':
        const result = await mcpServer.handleToolCall(data.name, data.args);
        ws.send(JSON.stringify({ type: 'tool_result', result }));
        break;
        
      case 'resource_read':
        const resource = await mcpServer.readResource(data.uri);
        ws.send(JSON.stringify({ type: 'resource_data', resource }));
        break;
    }
  });
});

console.log('WebSocket MCP Server running on port 50002');
```

---

## 2. MCP Wrapper Pattern

### 2.1 Why Wrappers Are Needed

The MCP protocol is designed to work with stdio-based communication, which presents challenges when integrating with HTTP-based services like Docker containers, cloud APIs, and existing web services. MCP wrappers solve this fundamental compatibility problem by translating between the stdio protocol used by MCP clients and the HTTP protocols used by modern services.

Wrappers provide several critical functions that enable seamless integration:

**Protocol Translation** is the primary function of MCP wrappers. The wrapper receives JSON-RPC messages from stdin, translates them into HTTP requests, sends them to the target service, and translates the HTTP responses back into JSON-RPC format for stdout. This translation happens transparently, making the HTTP service appear as a native stdio MCP server.

**Authentication Management** is handled by wrappers, which can maintain API keys, tokens, and credentials without exposing them to the MCP client. Wrappers can implement OAuth flows, JWT token refresh, and other authentication patterns without requiring changes to the MCP client or the underlying service.

**Connection Pooling** and resource management are handled by wrappers, which can maintain persistent connections to backend services, implement retry logic, and manage rate limiting. This improves performance and reliability compared to creating new connections for each request.

**Error Transformation** allows wrappers to convert service-specific error formats into standard JSON-RPC error responses. This ensures consistent error handling across different backend services and enables proper error propagation to MCP clients.

The following diagram illustrates the wrapper architecture:

```
┌─────────────┐     stdio      ┌─────────────┐     HTTP      ┌─────────────┐
│ MCP Client  │ ◄─────────────►│ MCP Wrapper │ ◄─────────────►│ HTTP Service│
│ (OpenCode)  │   JSON-RPC     │ (stdio/HTTP)│   REST/WS     │ (Docker/API)│
└─────────────┘                └─────────────┘               └─────────────┘
```

### 2.2 Node.js Wrapper Template

The following template provides a complete Node.js implementation for an MCP wrapper that connects to HTTP-based services:

```javascript
#!/usr/bin/env node

/**
 * MCP HTTP Wrapper Template
 * Converts stdio-based MCP communication to HTTP API calls
 */

const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
const { StdioServerTransport } = require('@modelcontextprotocol/sdk/server/stdio.js');
const axios = require('axios');

// Configuration from environment variables
const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';
const API_KEY = process.env.API_KEY || '';
const REQUEST_TIMEOUT = parseInt(process.env.REQUEST_TIMEOUT || '30000', 10);

// Create axios instance with default configuration
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: REQUEST_TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
    ...(API_KEY && { 'Authorization': `Bearer ${API_KEY}` })
  }
});

// MCP Server instance
const server = new Server(
  { name: 'http-mcp-wrapper', version: '1.0.0' },
  { capabilities: { tools: {} } }
);

// Tool definitions - these define what the MCP client can invoke
const toolDefinitions = [
  {
    name: 'http_get',
    description: 'Perform a GET request to the API',
    inputSchema: {
      type: 'object',
      properties: {
        endpoint: {
          type: 'string',
          description: 'API endpoint path (e.g., /users/1)'
        },
        params: {
          type: 'object',
          description: 'Query parameters'
        }
      },
      required: ['endpoint']
    }
  },
  {
    name: 'http_post',
    description: 'Perform a POST request to the API',
    inputSchema: {
      type: 'object',
      properties: {
        endpoint: {
          type: 'string',
          description: 'API endpoint path'
        },
        data: {
          type: 'object',
          description: 'Request body data'
        }
      },
      required: ['endpoint', 'data']
    }
  },
  {
    name: 'http_put',
    description: 'Perform a PUT request to update resources',
    inputSchema: {
      type: 'object',
      properties: {
        endpoint: {
          type: 'string',
          description: 'API endpoint path'
        },
        data: {
          type: 'object',
          description: 'Update data'
        }
      },
      required: ['endpoint', 'data']
    }
  },
  {
    name: 'http_delete',
    description: 'Perform a DELETE request',
    inputSchema: {
      type: 'object',
      properties: {
        endpoint: {
          type: 'string',
          description: 'API endpoint path'
        }
      },
      required: ['endpoint']
    }
  }
];

// Register tool handlers
server.setRequestHandler('tools/list', async () => {
  return {
    tools: toolDefinitions
  };
});

server.setRequestHandler('tools/call', async (request) => {
  const { name, arguments: args } = request.params;
  
  try {
    let response;
    
    switch (name) {
      case 'http_get':
        response = await apiClient.get(args.endpoint, { params: args.params });
        return {
          content: [{
            type: 'text',
            text: JSON.stringify(response.data, null, 2)
          }]
        };
        
      case 'http_post':
        response = await apiClient.post(args.endpoint, args.data);
        return {
          content: [{
            type: 'text',
            text: JSON.stringify(response.data, null, 2)
          }]
        };
        
      case 'http_put':
        response = await apiClient.put(args.endpoint, args.data);
        return {
          content: [{
            type: 'text',
            text: JSON.stringify(response.data, null, 2)
          }]
        };
        
      case 'http_delete':
        response = await apiClient.delete(args.endpoint);
        return {
          content: [{
            type: 'text',
            text: JSON.stringify(response.data, null, 2)
          }]
        };
        
      default:
        throw new Error(`Unknown tool: ${name}`);
    }
  } catch (error) {
    // Transform errors to JSON-RPC error format
    const errorMessage = error.response?.data?.message || error.message;
    return {
      content: [{
        type: 'text',
        text: `Error: ${errorMessage}`
      }],
      isError: true
    };
  }
});

// Connect to stdio transport and start server
const transport = new StdioServerTransport();
server.connect(transport).catch((error) => {
  console.error('Failed to connect MCP server:', error);
  process.exit(1);
});

// Handle graceful shutdown
process.on('SIGINT', () => {
  console.log('Received SIGINT, shutting down gracefully');
  server.close();
  process.exit(0);
});

process.on('SIGTERM', () => {
  console.log('Received SIGTERM, shutting down gracefully');
  server.close();
  process.exit(0);
});
```

**Advanced Wrapper with Connection Pooling:**

```javascript
#!/usr/bin/env node

/**
 * Advanced MCP HTTP Wrapper with Connection Pooling
 * Includes retry logic, caching, and error handling
 */

const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
const { StdioServerTransport } = require('@modelcontextprotocol/sdk/server/stdio.js');
const axios = require('axios');

// Configuration
const config = {
  apiBaseUrl: process.env.API_BASE_URL || 'http://localhost:8080',
  apiKey: process.env.API_KEY || '',
  timeout: parseInt(process.env.REQUEST_TIMEOUT || '30000', 10),
  retries: parseInt(process.env.MAX_RETRIES || '3', 10),
  retryDelay: parseInt(process.env.RETRY_DELAY || '1000', 10),
  cacheEnabled: process.env.CACHE_ENABLED === 'true',
  cacheTTL: parseInt(process.env.CACHE_TTL || '60000', 10)
};

// Create axios instance with connection pooling
const apiClient = axios.create({
  baseURL: config.apiBaseUrl,
  timeout: config.timeout,
  headers: {
    'Content-Type': 'application/json',
    ...(config.apiKey && { 'Authorization': `Bearer ${config.apiKey}` })
  },
  // Connection pool configuration
  httpAgent: new (require('http').Agent)({
    maxSockets: 25,
    maxFreeSockets: 10,
    timeout: config.timeout,
    keepAlive: true
  }),
  httpsAgent: new (require('https').Agent)({
    maxSockets: 25,
    maxFreeSockets: 10,
    timeout: config.timeout,
    keepAlive: true
  })
});

// Simple in-memory cache
const cache = new Map();

function getCached(key) {
  if (!config.cacheEnabled) return null;
  const cached = cache.get(key);
  if (cached && Date.now() < cached.expiresAt) {
    return cached.value;
  }
  cache.delete(key);
  return null;
}

function setCache(key, value) {
  if (!config.cacheEnabled) return;
  cache.set(key, {
    value,
    expiresAt: Date.now() + config.cacheTTL
  });
}

// Retry logic with exponential backoff
async function withRetry(fn, operationName) {
  let lastError;
  
  for (let attempt = 0; attempt <= config.retries; attempt++) {
    try {
      return await fn();
    } catch (error) {
      lastError = error;
      
      // Don't retry on client errors (4xx)
      if (error.response?.status >= 400 && error.response?.status < 500) {
        throw error;
      }
      
      if (attempt < config.retries) {
        const delay = config.retryDelay * Math.pow(2, attempt);
        console.error(`${operationName} failed (attempt ${attempt + 1}), retrying in ${delay}ms:`, error.message);
        await new Promise(resolve => setTimeout(resolve, delay));
      }
    }
  }
  
  throw lastError;
}

// MCP Server setup
const server = new Server(
  { name: 'advanced-http-mcp-wrapper', version: '1.0.0' },
  { capabilities: { tools: {}, resources: {} } }
);

// Tool definitions
const tools = [
  {
    name: 'api_query',
    description: 'Query the API with GET request and optional caching',
    inputSchema: {
      type: 'object',
      properties: {
        endpoint: { type: 'string', description: 'API endpoint path' },
        params: { type: 'object', description: 'Query parameters' },
        useCache: { type: 'boolean', description: 'Enable caching for this request' }
      },
      required: ['endpoint']
    }
  },
  {
    name: 'api_create',
    description: 'Create a new resource via POST',
    inputSchema: {
      type: 'object',
      properties: {
        endpoint: { type: 'string', description: 'API endpoint path' },
        data: { type: 'object', description: 'Resource data to create' }
      },
      required: ['endpoint', 'data']
    }
  }
];

server.setRequestHandler('tools/list', async () => ({ tools }));

server.setRequestHandler('tools/call', async (request) => {
  const { name, arguments: args } = request.params;
  
  try {
    switch (name) {
      case 'api_query': {
        const cacheKey = `${args.endpoint}:${JSON.stringify(args.params || {})}`;
        
        if (args.useCache !== false) {
          const cached = getCached(cacheKey);
          if (cached) {
            return { content: [{ type: 'text', text: JSON.stringify(cached) }] };
          }
        }
        
        const response = await withRetry(
          () => apiClient.get(args.endpoint, { params: args.params }),
          `GET ${args.endpoint}`
        );
        
        if (args.useCache !== false) {
          setCache(cacheKey, response.data);
        }
        
        return { content: [{ type: 'text', text: JSON.stringify(response.data) }] };
      }
      
      case 'api_create': {
        const response = await withRetry(
          () => apiClient.post(args.endpoint, args.data),
          `POST ${args.endpoint}`
        );
        return { content: [{ type: 'text', text: JSON.stringify(response.data) }] };
      }
      
      default:
        throw new Error(`Unknown tool: ${name}`);
    }
  } catch (error) {
    const errorMessage = error.response?.data?.message || error.message;
    return {
      content: [{ type: 'text', text: `Error: ${errorMessage}` }],
      isError: true
    };
  }
});

// Start server
const transport = new StdioServerTransport();
server.connect(transport).catch((error) => {
  console.error('Failed to connect MCP server:', error);
  process.exit(1);
});
```

### 2.3 Python Wrapper Template

Python wrappers use similar patterns but leverage Python's extensive library ecosystem for HTTP communication:

```python
#!/usr/bin/env python3
"""
MCP HTTP Wrapper - Python Implementation
Translates stdio-based MCP communication to HTTP API calls
"""

import os
import sys
import json
import asyncio
from typing import Any, Dict, Optional
from dataclasses import dataclass

import httpx
from mcp.server import Server
from mcp.server.stdio import StdioServerTransport
from mcp.types import Tool, TextContent
from mcp.server.handlers import ListToolsHandler, CallToolHandler


@dataclass
class Config:
    """Configuration loaded from environment variables"""
    api_base_url: str = os.getenv("API_BASE_URL", "http://localhost:8080")
    api_key: str = os.getenv("API_KEY", "")
    request_timeout: int = int(os.getenv("REQUEST_TIMEOUT", "30"))
    max_retries: int = int(os.getenv("MAX_RETRIES", "3"))
    retry_delay: float = float(os.getenv("RETRY_DELAY", "1.0"))


class HTTPWrapper:
    """HTTP client wrapper with retry logic"""
    
    def __init__(self, config: Config):
        self.config = config
        self.client = httpx.AsyncClient(
            base_url=config.api_base_url,
            timeout=config.request_timeout,
            headers={
                "Content-Type": "application/json",
                **({"Authorization": f"Bearer {config.api_key}"} if config.api_key else {})
            }
        )
    
    async def get(self, endpoint: str, params: Optional[Dict] = None) -> Dict[str, Any]:
        """Perform GET request with retry logic"""
        for attempt in range(self.config.max_retries):
            try:
                response = await self.client.get(endpoint, params=params)
                response.raise_for_status()
                return response.json()
            except httpx.HTTPStatusError:
                raise
            except Exception as e:
                if attempt < self.config.max_retries - 1:
                    await asyncio.sleep(self.config.retry_delay * (2 ** attempt))
                else:
                    raise
    
    async def post(self, endpoint: str, data: Dict[str, Any]) -> Dict[str, Any]:
        """Perform POST request with retry logic"""
        for attempt in range(self.config.max_retries):
            try:
                response = await self.client.post(endpoint, json=data)
                response.raise_for_status()
                return response.json()
            except httpx.HTTPStatusError:
                raise
            except Exception as e:
                if attempt < self.config.max_retries - 1:
                    await asyncio.sleep(self.config.retry_delay * (2 ** attempt))
                else:
                    raise
    
    async def put(self, endpoint: str, data: Dict[str, Any]) -> Dict[str, Any]:
        """Perform PUT request with retry logic"""
        for attempt in range(self.config.max_retries):
            try:
                response = await self.client.put(endpoint, json=data)
                response.raise_for_status()
                return response.json()
            except httpx.HTTPStatusError:
                raise
            except Exception as e:
                if attempt < self.config.max_retries - 1:
                    await asyncio.sleep(self.config.retry_delay * (2 ** attempt))
                else:
                    raise
    
    async def delete(self, endpoint: str) -> Dict[str, Any]:
        """Perform DELETE request"""
        response = await self.client.delete(endpoint)
        response.raise_for_status()
        return response.json() if response.content else {}
    
    async def close(self):
        """Close the HTTP client"""
        await self.client.aclose()


# Define MCP tools
TOOLS = [
    Tool(
        name="http_get",
        description="Perform a GET request to the API",
        inputSchema={
            "type": "object",
            "properties": {
                "endpoint": {"type": "string", "description": "API endpoint path"},
                "params": {"type": "object", "description": "Query parameters"}
            },
            "required": ["endpoint"]
        }
    ),
    Tool(
        name="http_post",
        description="Perform a POST request to create a resource",
        inputSchema={
            "type": "object",
            "properties": {
                "endpoint": {"type": "string", "description": "API endpoint path"},
                "data": {"type": "object", "description": "Request body data"}
            },
            "required": ["endpoint", "data"]
        }
    ),
    Tool(
        name="http_put",
        description="Perform a PUT request to update a resource",
        inputSchema={
            "type": "object",
            "properties": {
                "endpoint": {"type": "string", "description": "API endpoint path"},
                "data": {"type": "object", "description": "Update data"}
            },
            "required": ["endpoint", "data"]
        }
    ),
    Tool(
        name="http_delete",
        description="Perform a DELETE request",
        inputSchema={
            "type": "object",
            "properties": {
                "endpoint": {"type": "string", "description": "API endpoint path"}
            },
            "required": ["endpoint"]
        }
    )
]


class MCPHTTPWrapper:
    """Main MCP server implementation"""
    
    def __init__(self):
        self.config = Config()
        self.http = HTTPWrapper(self.config)
        self.server = Server(
            "http-mcp-wrapper",
            "1.0.0",
            tools=TOOLS
        )
        
        @self.server.list_tools()
        async def list_tools() -> list[Tool]:
            return TOOLS
        
        @self.server.call_tool()
        async def call_tool(name: str, arguments: dict) -> list[TextContent]:
            try:
                result = await self._handle_tool(name, arguments)
                return [TextContent(type="text", text=json.dumps(result, indent=2))]
            except Exception as e:
                return [TextContent(type="text", text=f"Error: {str(e)}")]
    
    async def _handle_tool(self, name: str, arguments: dict) -> Any:
        """Route tool calls to appropriate HTTP methods"""
        endpoint = arguments.get("endpoint")
        
        if name == "http_get":
            return await self.http.get(endpoint, arguments.get("params"))
        elif name == "http_post":
            return await self.http.post(endpoint, arguments.get("data", {}))
        elif name == "http_put":
            return await self.http.put(endpoint, arguments.get("data", {}))
        elif name == "http_delete":
            return await self.http.delete(endpoint)
        else:
            raise ValueError(f"Unknown tool: {name}")
    
    async def run(self):
        """Start the MCP server"""
        transport = StdioServerTransport()
        await self.server.run(transport, self.server.create_initialization_options())


async def main():
    """Entry point"""
    wrapper = MCPHTTPWrapper()
    try:
        await wrapper.run()
    finally:
        await wrapper.http.close()


if __name__ == "__main__":
    asyncio.run(main())
```

### 2.4 stdio to HTTP Conversion

The core challenge in MCP wrapper development is converting between the stdio-based JSON-RPC protocol used by MCP clients and the HTTP protocol used by backend services. This section details the conversion process.

**Message Flow:**

1. **Client to Server (stdio → HTTP):**
   - MCP client sends JSON-RPC message to stdin
   - Wrapper reads and parses the JSON-RPC message
   - Wrapper extracts method name, parameters, and request ID
   - Wrapper translates to appropriate HTTP request
   - HTTP request is sent to backend service

2. **Server to Client (HTTP → stdio):**
   - Backend service returns HTTP response
   - Wrapper extracts response data or error
   - Wrapper creates JSON-RPC response message
   - Wrapper writes response to stdout
   - MCP client reads response from stdout

**Detailed Implementation:**

```typescript
// stdio-to-http-bridge.ts - Core conversion logic

interface JSONRPCRequest {
  jsonrpc: '2.0';
  id: number | string;
  method: string;
  params?: Record<string, unknown>;
}

interface JSONRPCResponse {
  jsonrpc: '2.0';
  id: number | string;
  result?: unknown;
  error?: {
    code: number;
    message: string;
    data?: unknown;
  };
}

class StdioToHTTPBridge {
  private requestId = 0;
  private pendingRequests = new Map<number, {
    resolve: (value: unknown) => void;
    reject: (error: Error) => void;
  }>();
  
  constructor(
    private httpClient: AxiosInstance,
    private logger: Logger
  ) {}
  
  /**
   * Process incoming JSON-RPC request from stdin
   */
  async processStdinMessage(message: string): Promise<void> {
    try {
      const request: JSONRPCRequest = JSON.parse(message);
      
      if (request.method === 'tools/list') {
        // Handle tool listing
        const tools = await this.listTools();
        this.sendResponse(request.id, tools);
      } else if (request.method === 'tools/call') {
        // Handle tool invocation
        const { name, arguments: args } = request.params;
        const result = await this.handleToolCall(name, args);
        this.sendResponse(request.id, result);
      } else if (request.method.startsWith('resources/')) {
        // Handle resource operations
        const result = await this.handleResourceOperation(request.method, request.params);
        this.sendResponse(request.id, result);
      }
    } catch (error) {
      this.logger.error('Failed to process message:', error);
      // Send error response
      this.sendErrorResponse(-32700, 'Parse error', null);
    }
  }
  
  /**
   * Handle tool call - translate to HTTP request
   */
  private async handleToolCall(toolName: string, args: Record<string, unknown>): Promise<unknown> {
    this.logger.info(`Calling tool: ${toolName}`, args);
    
    // Map tool names to HTTP endpoints and methods
    const toolMapping = this.getToolMapping(toolName);
    
    if (!toolMapping) {
      throw new Error(`Unknown tool: ${toolName}`);
    }
    
    const { method, endpoint, transformParams } = toolMapping;
    const httpParams = transformParams ? transformParams(args) : args;
    
    // Make HTTP request
    const response = await this.httpClient.request({
      method,
      url: endpoint,
      ...(method === 'GET' ? { params: httpParams } : { data: httpParams })
    });
    
    return response.data;
  }
  
  /**
   * Send JSON-RPC response to stdout
   */
  private sendResponse(id: number | string, result: unknown): void {
    const response: JSONRPCResponse = {
      jsonrpc: '2.0',
      id,
      result
    };
    console.log(JSON.stringify(response));
  }
  
  /**
   * Send JSON-RPC error response to stdout
   */
  private sendErrorResponse(code: number, message: string, id: number | string | null): void {
    const response: JSONRPCResponse = {
      jsonrpc: '2.0',
      id: id ?? -1,
      error: { code, message }
    };
    console.error(JSON.stringify(response));
  }
}
```

---

## 3. Local MCP Servers

### 3.1 Serena MCP (Orchestration)

Serena MCP serves as the central orchestration layer for AI agent coordination. It provides tools for managing complex multi-agent workflows, coordinating task delegation, and maintaining state across agent interactions.

**Installation:**

```bash
pip install serena-mcp
# or
npm install -g serena-mcp
```

**Configuration:**

```json
{
  "mcp": {
    "serena": {
      "type": "local",
      "command": ["serena", "start-mcp-server"],
      "enabled": true,
      "environment": {
        "SERENA_API_KEY": "${SERENA_API_KEY}",
        "SERENA_LOG_LEVEL": "info"
      }
    }
  }
}
```

**Available Tools:**

| Tool | Description | Parameters |
|------|-------------|-------------|
| `orchestrate_workflow` | Execute a multi-step workflow | workflow: object, context: object |
| `delegate_task` | Delegate a task to another agent | task: object, agent: string |
| `coordinate_agents` | Coordinate multiple agents | agents: string[], task: object |
| `gather_results` | Collect results from multiple agents | taskIds: string[] |

### 3.2 Tavily MCP (Web Search)

Tavily MCP provides web search capabilities with optimized results for AI applications. It offers semantic search, content extraction, and news aggregation.

**Installation:**

```bash
npm install -g @tavily/claude-mcp
```

**Configuration:**

```json
{
  "mcp": {
    "tavily": {
      "type": "local",
      "command": ["npx", "@tavily/claude-mcp"],
      "enabled": true,
      "environment": {
        "TAVILY_API_KEY": "${TAVILY_API_KEY}"
      }
    }
  }
}
```

**Available Tools:**

| Tool | Description | Parameters |
|------|-------------|-------------|
| `search` | Perform a web search | query: string, max_results?: number |
| `search_news` | Search for recent news | query: string, days?: number |
| `extract_content` | Extract content from URL | url: string |

### 3.3 Canva MCP (Design)

Canva MCP integrates with Canva's design API, enabling AI assistants to create and manage designs programmatically.

**Installation:**

```bash
npm install -g @canva/claude-mcp
```

**Configuration:**

```json
{
  "mcp": {
    "canva": {
      "type": "local",
      "command": ["npx", "@canva/claude-mcp"],
      "enabled": true,
      "environment": {
        "CANVA_API_KEY": "${CANVA_API_KEY}"
      }
    }
  }
}
```

**Available Tools:**

| Tool | Description | Parameters |
|------|-------------|-------------|
| `create_design` | Create a new design | template_id: string, data: object |
| `get_design` | Retrieve design details | design_id: string |
| `export_design` | Export design to file | design_id: string, format: string |

### 3.4 Context7 MCP (Documentation)

Context7 MCP provides access to comprehensive documentation and code examples for various libraries and frameworks.

**Installation:**

```bash
npm install -g @anthropics/context7-mcp
```

**Configuration:**

```json
{
  "mcp": {
    "context7": {
      "type": "local",
      "command": ["npx", "@anthropics/context7-mcp"],
      "enabled": true
    }
  }
}
```

**Available Tools:**

| Tool | Description | Parameters |
|------|-------------|-------------|
| `resolve_library` | Resolve library ID for documentation | libraryName: string, query: string |
| `query_docs` | Query documentation and examples | libraryId: string, query: string |

### 3.5 Skyvern MCP (Browser)

Skyvern MCP provides browser automation capabilities with visual AI for element detection and autonomous navigation.

**Installation:**

```bash
pip install skyvern-mcp
python -m skyvern.mcp.server
```

**Configuration:**

```json
{
  "mcp": {
    "skyvern": {
      "type": "local",
      "command": ["python", "-m", "skyvern.mcp.server"],
      "enabled": true,
      "environment": {
        "SKYVERN_API_KEY": "${SKYVERN_API_KEY}",
        "LLM_PROVIDER": "mistral",
        "DATABASE_URL": "postgresql://user:pass@host:5432/db"
      },
      "port": 50006
    }
  }
}
```

**Available Tools:**

| Tool | Description | Parameters |
|------|-------------|-------------|
| `navigate_and_solve` | Navigate to URL and solve task | url: string, task: string |
| `analyze_screenshot` | Analyze screenshot with visual AI | screenshot: string, task: string |
| `extract_elements` | Extract clickable elements | url: string |

### 3.6 Chrome DevTools MCP

Chrome DevTools MCP provides direct access to Chrome's debugging protocol, enabling sophisticated browser automation and inspection.

**Installation:**

```bash
npm install -g @anthropics/chrome-devtools-mcp
```

**Configuration:**

```json
{
  "mcp": {
    "chrome-devtools": {
      "type": "local",
      "command": ["npx", "@anthropics/chrome-devtools-mcp"],
      "enabled": true,
      "environment": {
        "CDP_ENDPOINT": "ws://localhost:9222"
      }
    }
  }
}
```

**Available Tools:**

| Tool | Description | Parameters |
|------|-------------|-------------|
| `create_session` | Create new browser session | options?: object |
| `navigate` | Navigate to URL | sessionId: string, url: string |
| `screenshot` | Capture screenshot | sessionId: string |
| `click` | Click element | sessionId: string, selector: string |
| `type` | Type text | sessionId: string, selector: string, text: string |

---

## 4. Remote MCP Servers

### 4.1 Docker Container as MCP

Deploying MCP servers as Docker containers provides isolation, consistency, and easy scaling. This section covers the architecture and implementation patterns.

**Architecture Overview:**

```
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│  OpenCode       │      │  MCP Wrapper    │      │  Docker         │
│  Client         │─────►│  (stdio/HTTP)   │─────►│  Container      │
│                 │      │                 │      │  (HTTP Service) │
└─────────────────┘      └─────────────────┘      └─────────────────┘
```

**Dockerfile for MCP Server:**

```dockerfile
FROM node:20-alpine

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm ci --only=production

# Copy application code
COPY dist/ ./dist/

# Create non-root user
RUN addgroup -g 1001 -S nodejs && \
    adduser -S mcpuser -u 1001
USER mcpuser

# Expose port
EXPOSE 50001

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD node -e "require('http').get('http://localhost:50001/health', (r) => process.exit(r.statusCode === 200 ? 0 : 1))"

# Start command
CMD ["node", "dist/index.js"]
```

**docker-compose.yml:**

```yaml
version: '3.8'
services:
  my-mcp:
    build: 
      context: .
      dockerfile: Dockerfile
    ports:
      - "50001:50001"
    environment:
      - NODE_ENV=production
      - API_KEY=${MCP_API_KEY}
      - LOG_LEVEL=info
    volumes:
      - ./config:/app/config:ro
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "node", "-e", "require('http').get('http://localhost:50001/health', (r) => process.exit(r.statusCode === 200 ? 0 : 1))"]
      interval: 30s
      timeout: 10s
      retries: 3

networks:
  default:
    name: mcp-network
```

### 4.2 HTTP API Endpoints

MCP wrappers expose HTTP endpoints that conform to RESTful patterns. This section documents the standard endpoint structure.

**Standard Endpoint Structure:**

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/mcp/tools` | List available tools |
| POST | `/mcp/tools/call` | Call a tool |
| GET | `/mcp/resources` | List resources |
| GET | `/mcp/resources/{uri}` | Read a resource |

**Example Implementation:**

```typescript
// http-endpoints.ts

import express, { Request, Response, NextFunction } from 'express';
import { Server } from '@modelcontextprotocol/sdk/server/index.js';

const app = express();
app.use(express.json());

// Create MCP server instance
const mcpServer = new Server({ name: 'mcp-http-server', version: '1.0.0' }, {
  capabilities: { tools: {}, resources: {} }
});

// Health check endpoint
app.get('/health', (req: Request, res: Response) => {
  res.json({
    status: 'healthy',
    version: '1.0.0',
    timestamp: new Date().toISOString()
  });
});

// List available tools
app.get('/mcp/tools', async (req: Request, res: Response) => {
  try {
    const tools = await mcpServer.getTools();
    res.json({ tools });
  } catch (error) {
    res.status(500).json({ error: 'Failed to list tools' });
  }
});

// Call a tool
app.post('/mcp/tools/call', async (req: Request, res: Response) => {
  const { name, arguments: args } = req.body;
  
  if (!name) {
    return res.status(400).json({ error: 'Missing tool name' });
  }
  
  try {
    const result = await mcpServer.handleToolCall(name, args);
    res.json({ success: true, result });
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unknown error';
    res.status(500).json({ success: false, error: message });
  }
});

// List resources
app.get('/mcp/resources', async (req: Request, res: Response) => {
  try {
    const resources = await mcpServer.getResources();
    res.json({ resources });
  } catch (error) {
    res.status(500).json({ error: 'Failed to list resources' });
  }
});

// Read a specific resource
app.get('/mcp/resources/:uri', async (req: Request, res: Response) => {
  const { uri } = req.params;
  
  try {
    const resource = await mcpServer.readResource(uri);
    res.json({ resource });
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unknown error';
    res.status(404).json({ error: message });
  }
});

// Error handling middleware
app.use((err: Error, req: Request, res: Response, next: NextFunction) => {
  console.error('Unhandled error:', err);
  res.status(500).json({ error: 'Internal server error' });
});

const PORT = process.env.PORT || 50001;
app.listen(PORT, () => {
  console.log(`MCP HTTP Server running on port ${PORT}`);
});
```

### 4.3 Authentication

Authentication is critical for protecting MCP servers and ensuring only authorized clients can access tools and resources.

**API Key Authentication:**

```typescript
// api-key-auth.ts

import { Request, Response, NextFunction } from 'express';

const API_KEY = process.env.API_KEY;

export function apiKeyAuth(req: Request, res: Response, next: NextFunction): void {
  const providedKey = req.headers['x-api-key'] as string;
  
  if (!providedKey) {
    res.status(401).json({ error: 'API key required' });
    return;
  }
  
  if (providedKey !== API_KEY) {
    res.status(403).json({ error: 'Invalid API key' });
    return;
  }
  
  next();
}

// Apply to routes
app.use('/mcp', apiKeyAuth);
app.use('/health', (req, res, next) => next()); // Health check without auth
```

**JWT Token Authentication:**

```typescript
// jwt-auth.ts

import jwt from 'jsonwebtoken';
import { Request, Response, NextFunction } from 'express';

const JWT_SECRET = process.env.JWT_SECRET || 'secret';

export interface AuthenticatedRequest extends Request {
  user?: {
    id: string;
    roles: string[];
  };
}

export function jwtAuth(req: AuthenticatedRequest, res: Response, next: NextFunction): void {
  const authHeader = req.headers.authorization;
  
  if (!authHeader?.startsWith('Bearer ')) {
    res.status(401).json({ error: 'Bearer token required' });
    return;
  }
  
  const token = authHeader.substring(7);
  
  try {
    const decoded = jwt.verify(token, JWT_SECRET) as { id: string; roles: string[] };
    req.user = decoded;
    next();
  } catch (error) {
    res.status(403).json({ error: 'Invalid token' });
  }
}

// Role-based access control
export function requireRole(...roles: string[]) {
  return (req: AuthenticatedRequest, res: Response, next: NextFunction): void => {
    if (!req.user) {
      res.status(401).json({ error: 'Authentication required' });
      return;
    }
    
    const hasRole = roles.some(role => req.user!.roles.includes(role));
    if (!hasRole) {
      res.status(403).json({ error: 'Insufficient permissions' });
      return;
    }
    
    next();
  };
}
```

### 4.4 Configuration Examples

**Complete Remote MCP Configuration:**

```json
{
  "mcp": {
    "my-remote-mcp": {
      "type": "remote",
      "url": "https://mcp.example.com",
      "enabled": true,
      "timeout": 60000,
      "headers": {
        "X-API-Key": "${MCP_API_KEY}"
      }
    }
  }
}
```

**Docker-based MCP with Authentication:**

```yaml
version: '3.8'
services:
  authenticated-mcp:
    build: ./mcp-server
    ports:
      - "50003:50003"
    environment:
      - NODE_ENV=production
      - API_KEY=${MCP_API_KEY}
      - JWT_SECRET=${JWT_SECRET}
      - DATABASE_URL=postgresql://user:pass@postgres:5432/db
      - REDIS_URL=redis://redis:6379
      - LOG_LEVEL=info
      - RATE_LIMIT=100
    secrets:
      - api_key
      - jwt_secret
    restart: unless-stopped

secrets:
  api_key:
    file: ./secrets/api_key
  jwt_secret:
    file: ./secrets/jwt_secret
```

---

## 5. MCP Configuration

### 5.1 opencode.json MCP Config

OpenCode uses a JSON configuration file to define MCP server connections. This section provides comprehensive configuration examples.

**Basic Configuration:**

```json
{
  "mcp": {
    "server-name": {
      "type": "local",
      "command": ["npx", "package-name"],
      "enabled": true
    }
  }
}
```

**Advanced Configuration:**

```json
{
  "mcp": {
    "serena": {
      "type": "local",
      "command": ["uvx", "serena", "start-mcp-server"],
      "enabled": true,
      "environment": {
        "SERENA_API_KEY": "${SERENA_API_KEY}",
        "LOG_LEVEL": "info"
      }
    },
    "tavily": {
      "type": "local",
      "command": ["npx", "@tavily/claude-mcp"],
      "enabled": true,
      "environment": {
        "TAVILY_API_KEY": "${TAVILY_API_KEY}"
      }
    },
    "context7": {
      "type": "local",
      "command": ["npx", "@anthropics/context7-mcp"],
      "enabled": true
    },
    "skyvern": {
      "type": "local",
      "command": ["python", "-m", "skyvern.mcp.server"],
      "enabled": true,
      "environment": {
        "SKYVERN_API_KEY": "${SKYVERN_API_KEY}",
        "LLM_PROVIDER": "mistral",
        "DATABASE_URL": "postgresql://skyvern:pass@room-03-postgres:5432/skyvern"
      }
    },
    "custom-http-mcp": {
      "type": "remote",
      "url": "https://custom-mcp.example.com",
      "enabled": true,
      "timeout": 30000,
      "headers": {
        "X-API-Key": "${CUSTOM_MCP_API_KEY}"
      }
    }
  }
}
```

**Environment Variable Best Practices:**

```bash
# .env file - never commit to version control
SERENA_API_KEY=sk_serena_xxxxxxxxxxxxx
TAVILY_API_KEY=tvly_xxxxxxxxxxxxx
SKYVERN_API_KEY=sk_skyvern_xxxxxxxxxxxxx
MCP_API_KEY=mcp_xxxxxxxxxxxxx
CUSTOM_MCP_API_KEY=custom_xxxxxxxxxxxxx
```

### 5.2 Environment Variables

MCP servers use environment variables for configuration. This section documents common environment variables and their usage.

**Common Environment Variables:**

| Variable | Description | Example |
|----------|-------------|---------|
| `API_KEY` | API authentication key | `sk_xxxxxxxxxxxxx` |
| `API_BASE_URL` | Base URL for HTTP MCPs | `https://api.example.com` |
| `DATABASE_URL` | Database connection string | `postgresql://user:pass@host:5432/db` |
| `REDIS_URL` | Redis connection string | `redis://host:6379` |
| `LOG_LEVEL` | Logging verbosity | `debug`, `info`, `warn`, `error` |
| `REQUEST_TIMEOUT` | HTTP request timeout (ms) | `30000` |
| `MAX_RETRIES` | Number of retry attempts | `3` |
| `CACHE_ENABLED` | Enable response caching | `true` / `false` |

**Loading Environment Variables:**

```typescript
// config.ts
import dotenv from 'dotenv';

// Load .env file
dotenv.config();

export const config = {
  apiKey: process.env.API_KEY || '',
  apiBaseUrl: process.env.API_BASE_URL || 'http://localhost:8080',
  timeout: parseInt(process.env.REQUEST_TIMEOUT || '30000', 10),
  retries: parseInt(process.env.MAX_RETRIES || '3', 10),
  logLevel: process.env.LOG_LEVEL || 'info'
};
```

### 5.3 Port Assignment (50000-59999)

Following the AGENTS.md mandate, MCP servers must use unique ports in the 50000-59999 range to avoid conflicts.

**Port Assignment Table:**

| Service | Port | Protocol | Status |
|---------|------|----------|--------|
| Custom MCP 1 | 50001 | HTTP | Available |
| Custom MCP 2 | 50002 | WebSocket | Available |
| Skyvern MCP | 50006 | HTTP | Reserved |
| Local MCPs | 50010-50099 | stdio/HTTP | Available |
| Container MCPs | 50100-50999 | HTTP | Reserved |

**Port Configuration in docker-compose:**

```yaml
services:
  mcp-server-1:
    ports:
      - "50001:50001"  # Custom mapping
  
  mcp-server-2:
    ports:
      - "50002:50002"  # WebSocket MCP
```

### 5.4 Service Discovery

Service discovery enables MCP servers to find and connect to other services dynamically without hardcoded addresses.

**DNS-Based Service Discovery:**

```yaml
version: '3.8'
services:
  mcp-server:
    build:: my-mcp-service
    networks .
    hostname:
      - mcp-network

networks:
  mcp-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.25.0.0/16
```

**Environment-Based Discovery:**

```typescript
// service-discovery.ts

interface ServiceEndpoint {
  host: string;
  port: number;
  protocol: 'http' | 'https' | 'ws' | 'wss';
}

class ServiceDiscovery {
  private services = new Map<string, ServiceEndpoint>();
  
  constructor() {
    this.discoverServices();
  }
  
  private discoverServices(): void {
    // Discover from environment variables
    const services = [
      { name: 'postgres', env: 'POSTGRES_HOST', port: 5432, defaultHost: 'room-03-postgres' },
      { name: 'redis', env: 'REDIS_HOST', port: 6379, defaultHost: 'room-04-redis' },
      { name: 'skyvern', env: 'SKYVERN_HOST', port: 50006, defaultHost: 'agent-06-skyvern' }
    ];
    
    for (const service of services) {
      const host = process.env[service.env] || service.defaultHost;
      this.services.set(service.name, {
        host,
        port: parseInt(process.env[`${service.name.toUpperCase()}_PORT`] || String(service.port), 10),
        protocol: process.env[`${service.name.toUpperCase()}_PROTOCOL`] as any || 'http'
      });
    }
  }
  
  getService(name: string): ServiceEndpoint | undefined {
    return this.services.get(name);
  }
  
  getDatabaseURL(serviceName: string = 'postgres'): string {
    const endpoint = this.getService(serviceName);
    if (!endpoint) throw new Error(`Service ${serviceName} not found`);
    
    const user = process.env[`${serviceName.toUpperCase()}_USER`] || 'user';
    const password = process.env[`${serviceName.toUpperCase()}_PASSWORD`] || 'pass';
    const db = process.env[`${serviceName.toUpperCase()}_DB`] || 'db';
    
    return `${endpoint.protocol}://${user}:${password}@${endpoint.host}:${endpoint.port}/${db}`;
  }
}
```

---

## 6. Creating New MCPs

### 6.1 Step-by-Step Guide

This section provides a comprehensive guide for creating new MCP servers from scratch.

**Step 1: Define the MCP Purpose**

Before implementation, clearly define:
- What functionality the MCP will provide
- What tools and resources it will expose
- What backend services it will connect to
- What authentication methods are required

**Step 2: Choose Implementation Language:**

| Language | Use Case | Pros | Cons |
|----------|----------|------|------|
| Node.js | Web APIs, TypeScript | Large ecosystem, async | Memory usage |
| Python | ML, Data, AI | Rich libraries | Slower startup |
| Go | High performance | Fast, small binaries | Complex error handling |

**Step 3: Set Up Project Structure:**

```bash
mkdir my-mcp-server
cd my-mcp-server
npm init -y
npm install @modelcontextprotocol/sdk axios express
```

**Step 4: Implement the MCP Server:**

```typescript
// src/index.ts
import { Server } from '@modelcontextprotocol/sdk/server/index.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import { Tool } from '@modelcontextprotocol/sdk/types.js';

const server = new Server(
  { name: 'my-mcp-server', version: '1.0.0' },
  { capabilities: { tools: {} } }
);

// Define tools
const tools: Tool[] = [
  {
    name: 'my_tool',
    description: 'Description of what the tool does',
    inputSchema: {
      type: 'object',
      properties: {
        param1: { type: 'string', description: 'First parameter' },
        param2: { type: 'number', description: 'Second parameter' }
      },
      required: ['param1']
    }
  }
];

server.setRequestHandler('tools/list', async () => ({ tools }));

server.setRequestHandler('tools/call', async (request) => {
  const { name, arguments: args } = request.params;
  
  if (name === 'my_tool') {
    // Implement tool logic
    const result = { success: true, data: args.param1 };
    return { content: [{ type: 'text', text: JSON.stringify(result) }] };
  }
  
  throw new Error(`Unknown tool: ${name}`);
});

const transport = new StdioServerTransport();
server.connect(transport);
```

**Step 5: Add HTTP Support (Optional):**

```typescript
// src/http-server.ts
import express from 'express';
import { Server } from '@modelcontextprotocol/sdk/server/index.js';

export function createHTTPServer(mcpServer: Server, port: number = 50001) {
  const app = express();
  app.use(express.json());
  
  app.get('/health', (req, res) => {
    res.json({ status: 'healthy' });
  });
  
  app.get('/mcp/tools', async (req, res) => {
    const tools = await mcpServer.getTools();
    res.json({ tools });
  });
  
  app.post('/mcp/tools/call', async (req, res) => {
    const { name, arguments: args } = req.body;
    try {
      const result = await mcpServer.handleToolCall(name, args);
      res.json({ success: true, result });
    } catch (error) {
      res.status(500).json({ success: false, error: (error as Error).message });
    }
  });
  
  return app.listen(port, () => {
    console.log(`MCP HTTP Server running on port ${port}`);
  });
}
```

**Step 6: Create Dockerfile:**

```dockerfile
FROM node:20-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY dist/ ./dist/
EXPOSE 50001
CMD ["node", "dist/index.js"]
```

### 6.2 Wrapper Implementation

When wrapping existing HTTP services as MCP servers, follow these implementation patterns.

**HTTP Service to MCP Wrapper:**

```typescript
// wrapper.ts - Complete wrapper implementation

import { Server } from '@modelcontextprotocol/sdk/server/index.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import axios, { AxiosInstance } from 'axios';

// Configuration
const config = {
  baseURL: process.env.API_BASE_URL || 'http://localhost:8080',
  apiKey: process.env.API_KEY || '',
  timeout: parseInt(process.env.TIMEOUT || '30000', 10)
};

// HTTP client
const httpClient: AxiosInstance = axios.create({
  baseURL: config.baseURL,
  timeout: config.timeout,
  headers: {
    'Content-Type': 'application/json',
    ...(config.apiKey && { 'Authorization': `Bearer ${config.apiKey}` })
  }
});

// MCP Server
const server = new Server(
  { name: 'api-wrapper', version: '1.0.0' },
  { capabilities: { tools: {} } }
);

// Tool definitions
const tools = [
  {
    name: 'get_data',
    description: 'Fetch data from the API',
    inputSchema: {
      type: 'object',
      properties: {
        endpoint: { type: 'string', description: 'API endpoint' },
        params: { type: 'object', description: 'Query parameters' }
      },
      required: ['endpoint']
    }
  },
  {
    name: 'create_data',
    description: 'Create new data via API',
    inputSchema: {
      type: 'object',
      properties: {
        endpoint: { type: 'string', description: 'API endpoint' },
        data: { type: 'object', description: 'Data to create' }
      },
      required: ['endpoint', 'data']
    }
  }
];

server.setRequestHandler('tools/list', async () => ({ tools }));

server.setRequestHandler('tools/call', async (request) => {
  const { name, arguments: args } = request.params;
  
  try {
    let response;
    
    switch (name) {
      case 'get_data':
        response = await httpClient.get(args.endpoint, { params: args.params });
        break;
      case 'create_data':
        response = await httpClient.post(args.endpoint, args.data);
        break;
      default:
        throw new Error(`Unknown tool: ${name}`);
    }
    
    return {
      content: [{ type: 'text', text: JSON.stringify(response.data) }]
    };
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unknown error';
    return {
      content: [{ type: 'text', text: `Error: ${message}` }],
      isError: true
    };
  }
});

// Start server
const transport = new StdioServerTransport();
server.connect(transport).catch(console.error);
```

### 6.3 Testing MCPs

Testing is critical for ensuring MCP reliability. This section covers testing strategies.

**Unit Testing Tools:**

```typescript
// __tests__/mcp-tools.test.ts

import { Server } from '@modelcontextprotocol/sdk/server/index.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';

describe('MCP Server Tools', () => {
  let server: Server;
  let transport: StdioServerTransport;
  
  beforeEach(() => {
    server = new Server(
      { name: 'test-server', version: '1.0.0' },
      { capabilities: { tools: {} } }
    );
    transport = new StdioServerTransport();
  });
  
  test('should list available tools', async () => {
    const tools = await server.getTools();
    expect(tools).toBeDefined();
    expect(Array.isArray(tools)).toBe(true);
  });
  
  test('should call tool successfully', async () => {
    const result = await server.handleToolCall('test_tool', { param: 'value' });
    expect(result).toBeDefined();
  });
  
  test('should handle tool errors', async () => {
    await expect(
      server.handleToolCall('nonexistent_tool', {})
    ).rejects.toThrow();
  });
});
```

**Integration Testing:**

```typescript
// __tests__/integration.test.ts

import axios from 'axios';

describe('MCP HTTP Integration', () => {
  const baseURL = process.env.MCP_BASE_URL || 'http://localhost:50001';
  const httpClient = axios.create({ baseURL });
  
  test('health endpoint returns healthy status', async () => {
    const response = await httpClient.get('/health');
    expect(response.status).toBe(200);
    expect(response.data.status).toBe('healthy');
  });
  
  test('can list tools', async () => {
    const response = await httpClient.get('/mcp/tools');
    expect(response.status).toBe(200);
    expect(response.data.tools).toBeDefined();
  });
  
  test('can call tool via HTTP', async () => {
    const response = await httpClient.post('/mcp/tools/call', {
      name: 'test_tool',
      arguments: { param: 'value' }
    });
    expect(response.status).toBe(200);
    expect(response.data.success).toBe(true);
  });
});
```

### 6.4 Documentation Requirements

Every MCP server must include comprehensive documentation covering:

**Required Documentation:**

| Document | Location | Content |
|----------|----------|---------|
| README.md | Project root | Installation, quick start, overview |
| API.md | docs/ | HTTP API endpoints, request/response formats |
| TOOLS.md | docs/ | Available tools, parameters, examples |
| CONFIG.md | docs/ | Configuration options, environment variables |
| TROUBLESHOOTING.md | docs/ | Common issues and solutions |

**README Template:**

```markdown
# MCP Server Name

Brief description of what this MCP server provides.

## Features

- Feature 1
- Feature 2
- Feature 3

## Installation

```bash
npm install
npm run build
```

## Configuration

| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| API_KEY | API authentication key | - |
| API_BASE_URL | Base URL for API | http://localhost:8080 |
| TIMEOUT | Request timeout (ms) | 30000 |

## Usage

### CLI Mode

```bash
node dist/index.js
```

### HTTP Mode

```bash
HTTP_PORT=50001 node dist/http-server.js
```

## Available Tools

| Tool | Description | Parameters |
|------|-------------|-------------|
| tool_name | Tool description | param1, param2 |

## Docker

```bash
docker build -t my-mcp .
docker run -p 50001:50001 -e API_KEY=xxx my-mcp
```

## License

MIT
```

---

## 7. Best Practices

### 7.1 Connection Pooling

Connection pooling is essential for high-performance MCP servers that handle multiple requests.

**HTTP Connection Pooling:**

```typescript
// connection-pool.ts

import axios, { AxiosInstance, Pool } from 'axios';

class ConnectionPoolManager {
  private pools: Map<string, AxiosInstance> = new Map();
  
  createPool(name: string, config: {
    baseURL: string;
    maxSockets?: number;
    timeout?: number;
  }): AxiosInstance {
    if (this.pools.has(name)) {
      return this.pools.get(name)!;
    }
    
    const pool = axios.create({
      baseURL: config.baseURL,
      timeout: config.timeout || 30000,
      httpAgent: new (require('http').Agent)({
        maxSockets: config.maxSockets || 25,
        keepAlive: true
      }),
      httpsAgent: new (require('https').Agent)({
        maxSockets: config.maxSockets || 25,
        keepAlive: true
      })
    });
    
    this.pools.set(name, pool);
    return pool;
  }
  
  getPool(name: string): AxiosInstance | undefined {
    return this.pools.get(name);
  }
  
  closeAll(): void {
    for (const pool of this.pools.values()) {
      pool.close();
    }
    this.pools.clear();
  }
}

export const poolManager = new ConnectionPoolManager();
```

**WebSocket Connection Pooling:**

```typescript
// websocket-pool.ts

import WebSocket from 'ws';

interface PooledConnection {
  ws: WebSocket;
  inUse: boolean;
  createdAt: number;
}

class WebSocketPool {
  private pool: PooledConnection[] = [];
  private maxPoolSize: number;
  private connectionUrl: string;
  
  constructor(connectionUrl: string, maxPoolSize: number = 10) {
    this.connectionUrl = connectionUrl;
    this.maxPoolSize = maxPoolSize;
  }
  
  async acquire(): Promise<WebSocket> {
    // Find available connection
    const available = this.pool.find(conn => !conn.inUse);
    
    if (available) {
      available.inUse = true;
      if (available.ws.readyState !== WebSocket.OPEN) {
        // Recreate if closed
        available.ws = await this.createConnection();
      }
      return available.ws;
    }
    
    // Create new if under limit
    if (this.pool.length < this.maxPoolSize) {
      const ws = await this.createConnection();
      this.pool.push({ ws, inUse: true, createdAt: Date.now() });
      return ws;
    }
    
    // Wait for available connection
    return new Promise((resolve) => {
      const checkInterval = setInterval(() => {
        const conn = this.pool.find(c => !c.inUse);
        if (conn) {
          clearInterval(checkInterval);
          conn.inUse = true;
          resolve(conn.ws);
        }
      }, 100);
    });
  }
  
  release(ws: WebSocket): void {
    const conn = this.pool.find(c => c.ws === ws);
    if (conn) {
      conn.inUse = false;
    }
  }
  
  private createConnection(): Promise<WebSocket> {
    return new Promise((resolve, reject) => {
      const ws = new WebSocket(this.connectionUrl);
      
      ws.on('open', () => resolve(ws));
      ws.on('error', reject);
      
      // Timeout
      setTimeout(() => {
        if (ws.readyState === WebSocket.CONNECTING) {
          ws.close();
          reject(new Error('Connection timeout'));
        }
      }, 10000);
    });
  }
}
```

### 7.2 Error Handling

Robust error handling ensures MCP servers provide useful feedback and recover gracefully from failures.

**Error Handler Implementation:**

```typescript
// error-handler.ts

interface MCPError {
  code: number;
  message: string;
  data?: unknown;
}

class MCPErrorHandler {
  // Standard JSON-RPC error codes
  static readonly PARSE_ERROR = -32700;
  static readonly INVALID_REQUEST = -32600;
  static readonly METHOD_NOT_FOUND = -32601;
  static readonly INVALID_PARAMS = -32602;
  static readonly INTERNAL_ERROR = -32603;
  
  // Custom error codes
  static readonly AUTHENTICATION_ERROR = -32001;
  static readonly RATE_LIMIT_ERROR = -32002;
  static readonly SERVICE_UNAVAILABLE = -32003;
  static readonly TIMEOUT_ERROR = -32004;
  
  static handleError(error: unknown): MCPError {
    if (error instanceof AuthenticationError) {
      return {
        code: this.AUTHENTICATION_ERROR,
        message: 'Authentication failed',
        data: error.message
      };
    }
    
    if (error instanceof RateLimitError) {
      return {
        code: this.RATE_LIMIT_ERROR,
        message: 'Rate limit exceeded',
        data: { retryAfter: error.retryAfter }
      };
    }
    
    if (error instanceof TimeoutError) {
      return {
        code: this.TIMEOUT_ERROR,
        message: 'Request timeout',
        data: error.message
      };
    }
    
    // Default to internal error
    return {
      code: this.INTERNAL_ERROR,
      message: 'Internal server error',
      data: error instanceof Error ? error.message : 'Unknown error'
    };
  }
  
  static toJSONRPCError(id: number | string | null, error: MCPError): string {
    return JSON.stringify({
      jsonrpc: '2.0',
      id,
      error: {
        code: error.code,
        message: error.message,
        data: error.data
      }
    });
  }
}

// Custom error classes
class AuthenticationError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'AuthenticationError';
  }
}

class RateLimitError extends Error {
  retryAfter: number;
  
  constructor(retryAfter: number) {
    super(`Rate limited, retry after ${retryAfter} seconds`);
    this.name = 'RateLimitError';
    this.retryAfter = retryAfter;
  }
}

class TimeoutError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'TimeoutError';
  }
}
```

### 7.3 Logging

Comprehensive logging enables debugging and monitoring of MCP server behavior.

**Logging Implementation:**

```typescript
// logger.ts

enum LogLevel {
  DEBUG = 0,
  INFO = 1,
  WARN = 2,
  ERROR = 3
}

interface LogEntry {
  timestamp: string;
  level: string;
  message: string;
  context?: Record<string, unknown>;
}

class Logger {
  private level: LogLevel;
  private logs: LogEntry[] = [];
  private maxLogs: number;
  
  constructor(level: LogLevel = LogLevel.INFO, maxLogs: number = 1000) {
    this.level = level;
    this.maxLogs = maxLogs;
  }
  
  debug(message: string, context?: Record<string, unknown>): void {
    this.log(LogLevel.DEBUG, message, context);
  }
  
  info(message: string, context?: Record<string, unknown>): void {
    this.log(LogLevel.INFO, message, context);
  }
  
  warn(message: string, context?: Record<string, unknown>): void {
    this.log(LogLevel.WARN, message, context);
  }
  
  error(message: string, context?: Record<string, unknown>): void {
    this.log(LogLevel.ERROR, message, context);
  }
  
  private log(level: LogLevel, message: string, context?: Record<string, unknown>): void {
    if (level < this.level) return;
    
    const entry: LogEntry = {
      timestamp: new Date().toISOString(),
      level: LogLevel[level],
      message,
      context
    };
    
    this.logs.push(entry);
    
    // Trim logs if needed
    if (this.logs.length > this.maxLogs) {
      this.logs.shift();
    }
    
    // Console output
    const color = this.getLevelColor(level);
    console.log(`${color}[${entry.timestamp}] ${entry.level}:${'\x1b[0m'} ${message}`, context || '');
  }
  
  private getLevelColor(level: LogLevel): string {
    switch (level) {
      case LogLevel.DEBUG: return '\x1b[36m'; // Cyan
      case LogLevel.INFO: return '\x1b[32m'; // Green
      case LogLevel.WARN: return '\x1b[33m'; // Yellow
      case LogLevel.ERROR: return '\x1b[31m'; // Red
    }
  }
  
  getLogs(): LogEntry[] {
    return [...this.logs];
  }
  
  getLogsByLevel(level: LogLevel): LogEntry[] {
    return this.logs.filter(log => log.level === LogLevel[level]);
  }
}

export const logger = new Logger(
  (process.env.LOG_LEVEL?.toUpperCase() as any) || LogLevel.INFO
);
```

### 7.4 Performance Optimization

Optimizing MCP server performance involves multiple strategies.

**Caching Strategy:**

```typescript
// cache.ts

interface CacheEntry<T> {
  value: T;
  expiresAt: number;
}

class CacheManager {
  private cache = new Map<string, CacheEntry<unknown>>();
  private defaultTTL: number;
  
  constructor(defaultTTL: number = 60000) {
    this.defaultTTL = defaultTTL;
    
    // Periodic cleanup
    setInterval(() => this.cleanup(), 60000);
  }
  
  set<T>(key: string, value: T, ttl?: number): void {
    const expiresAt = Date.now() + (ttl || this.defaultTTL);
    this.cache.set(key, { value, expiresAt });
  }
  
  get<T>(key: string): T | undefined {
    const entry = this.cache.get(key) as CacheEntry<T> | undefined;
    
    if (!entry) return undefined;
    
    if (Date.now() > entry.expiresAt) {
      this.cache.delete(key);
      return undefined;
    }
    
    return entry.value;
  }
  
  delete(key: string): void {
    this.cache.delete(key);
  }
  
  clear(): void {
    this.cache.clear();
  }
  
  private cleanup(): void {
    const now = Date.now();
    for (const [key, entry] of this.cache.entries()) {
      if (now > entry.expiresAt) {
        this.cache.delete(key);
      }
    }
  }
  
  stats(): { size: number; hitRate: number } {
    return {
      size: this.cache.size,
      hitRate: 0 // Would need hit/miss tracking
    };
  }
}
```

---

## 8. Troubleshooting

### 8.1 Connection Issues

Common connection problems and their solutions.

**stdio Connection Failed:**

```
Error: Failed to connect to MCP server
```

**Solutions:**

1. Verify the command is correct:
```json
{
  "mcp": {
    "server": {
      "command": ["correct", "command"],
      "enabled": true
    }
  }
}
```

2. Check if the package is installed:
```bash
npm list -g package-name
```

3. Verify working directory:
```json
{
  "mcp": {
    "server": {
      "command": ["npx", "package-name"],
      "cwd": "/path/to/working/directory"
    }
  }
}
```

**HTTP Connection Failed:**

```
Error: connect ECONNREFUSED 127.0.0.1:50001
```

**Solutions:**

1. Verify the service is running:
```bash
docker ps | grep mcp
curl http://localhost:50001/health
```

2. Check port configuration:
```bash
netstat -tlnp | grep 50001
```

3. Verify network connectivity:
```bash
ping container-name
curl http://container-name:50001/health
```

### 8.2 Timeout Problems

Timeout errors indicate the server is taking too long to respond.

**Timeout Error:**

```
Error: Request timeout after 30000ms
```

**Solutions:**

1. Increase timeout in configuration:
```json
{
  "mcp": {
    "server": {
      "timeout": 60000
    }
  }
}
```

2. Optimize server response time:
```typescript
// Add caching
app.get('/mcp/tools', cacheMiddleware('tools'), (req, res) => {
  res.json({ tools: getCachedTools() });
});
```

3. Implement async responses:
```typescript
app.post('/mcp/tools/call', async (req, res) => {
  const { name, arguments: args } = req.body;
  
  // Return immediately with job ID
  const jobId = await startBackgroundJob(name, args);
  res.json({ jobId });
});

// Poll for result
app.get('/mcp/jobs/:jobId', async (req, res) => {
  const result = await getJobResult(req.params.jobId);
  res.json(result);
});
```

### 8.3 Authentication Failures

Authentication errors prevent access to protected MCP servers.

**Error:**

```
Error: 401 Unauthorized
Error: 403 Forbidden
```

**Solutions:**

1. Verify API key is set:
```bash
echo $API_KEY
```

2. Check environment variable loading:
```typescript
// Verify in MCP server
console.log('API_KEY:', process.env.API_KEY ? 'set' : 'NOT SET');
```

3. Validate token format:
```bash
# JWT tokens should be three parts separated by dots
echo $JWT_TOKEN | cut -d. -f1,2,3
```

4. Check authentication endpoint:
```typescript
app.use('/mcp', (req, res, next) => {
  console.log('Auth header:', req.headers.authorization);
  next();
});
```

### 8.4 Common Errors

**JSON Parse Error:**

```
Error: Parse error - Invalid JSON
```

Cause: Malformed JSON in request or response
Solution: Validate JSON format, check for trailing commas

**Tool Not Found:**

```
Error: Unknown tool: tool_name
```

Cause: Tool name doesn't match registered tools
Solution: Check tool list with `/mcp/tools` endpoint

**Resource Not Found:**

```
Error: Resource not found: file:///resource/uri
```

Cause: Invalid or inaccessible resource URI
Solution: Verify resource exists and URI is correct

**Memory Error:**

```
FATAL ERROR: JavaScript heap out of memory
```

Cause: Memory leak or insufficient memory
Solution:
- Increase Node.js memory limit: `node --max-old-space-size=4096 server.js`
- Profile and fix memory leaks
- Add connection pooling limits

---

## Appendix A: Quick Reference

### MCP Configuration Checklist

- [ ] Define tools and resources
- [ ] Implement MCP server interface
- [ ] Add error handling
- [ ] Implement logging
- [ ] Add health check endpoint
- [ ] Configure authentication
- [ ] Write tests
- [ ] Create Dockerfile
- [ ] Document API and tools
- [ ] Test in production-like environment

### Common Commands

```bash
# Test MCP server locally
node dist/index.js

# Test HTTP endpoint
curl http://localhost:50001/health

# List tools via HTTP
curl http://localhost:50001/mcp/tools

# Call tool via HTTP
curl -X POST http://localhost:50001/mcp/tools/call \
  -H "Content-Type: application/json" \
  -d '{"name": "tool_name", "arguments": {}}'

# Docker build and run
docker build -t mcp-server .
docker run -p 50001:50001 -e API_KEY=xxx mcp-server
```

---

## Appendix B: Reference

### MCP Protocol Specification

The MCP protocol follows JSON-RPC 2.0 specification with additional conventions:

**Request Format:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "tool_name",
    "arguments": {}
  }
}
```

**Response Format:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {}
}
```

**Error Format:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32601,
    "message": "Method not found"
  }
}
```

---

**Document Version:** 1.0  
**Last Updated:** 2026-02-20  
**Maintained By:** BIOMETRICS Team  
**Status:** ACTIVE - MANDATORY COMPLIANCE
