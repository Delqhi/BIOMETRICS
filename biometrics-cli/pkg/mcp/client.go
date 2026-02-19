package mcp

import (
	"context"
	"time"
)

// MCPClient defines the interface for all MCP server clients
type MCPClient interface {
	// Name returns the MCP server name
	Name() string

	// Connect establishes connection to the MCP server
	Connect(ctx context.Context) error

	// Disconnect closes the connection
	Disconnect(ctx context.Context) error

	// Health checks the server health
	Health(ctx context.Context) (*HealthStatus, error)

	// IsConnected returns connection status
	IsConnected() bool
}

// HealthStatus represents MCP server health
type HealthStatus struct {
	Healthy   bool          `json:"healthy"`
	Latency   time.Duration `json:"latency"`
	Version   string        `json:"version"`
	Timestamp time.Time     `json:"timestamp"`
	Error     string        `json:"error,omitempty"`
}

// ToolCall represents a tool invocation request
type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
	Timeout   time.Duration          `json:"timeout"`
}

// ToolResult represents a tool invocation result
type ToolResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
}

// MCPConfig holds MCP server configuration
type MCPConfig struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"` // local, remote, docker
	Command     []string          `json:"command,omitempty"`
	URL         string            `json:"url,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Timeout     time.Duration     `json:"timeout"`
	Enabled     bool              `json:"enabled"`
}

// DefaultTimeout is the default MCP request timeout
const DefaultTimeout = 120 * time.Second // Qwen 3.5 needs 120s
