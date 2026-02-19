package mcp

import (
	"context"
	"time"
)

type MCPClient interface {
	Name() string
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Health(ctx context.Context) (*HealthStatus, error)
	IsConnected() bool
}

type HealthStatus struct {
	Healthy   bool          `json:"healthy"`
	Latency   time.Duration `json:"latency"`
	Version   string        `json:"version"`
	Timestamp time.Time     `json:"timestamp"`
	Error     string        `json:"error,omitempty"`
}

type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
	Timeout   time.Duration          `json:"timeout"`
}

type ToolResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
}

type MCPConfig struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Command     []string          `json:"command,omitempty"`
	URL         string            `json:"url,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Timeout     time.Duration     `json:"timeout"`
	Enabled     bool              `json:"enabled"`
}

const DefaultTimeout = 120 * time.Second
