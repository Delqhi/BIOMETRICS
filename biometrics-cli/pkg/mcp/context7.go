package mcp

import (
	"context"
	"fmt"
	"time"
)

type Context7Client struct {
	config    *MCPConfig
	connected bool
}

func NewContext7Client(config *MCPConfig) *Context7Client {
	return &Context7Client{
		config:    config,
		connected: false,
	}
}

func (c *Context7Client) Name() string {
	return "context7"
}

func (c *Context7Client) Connect(ctx context.Context) error {
	if c.connected {
		return nil
	}

	c.connected = true
	return nil
}

func (c *Context7Client) Disconnect(ctx context.Context) error {
	if !c.connected {
		return nil
	}

	c.connected = false
	return nil
}

func (c *Context7Client) Health(ctx context.Context) (*HealthStatus, error) {
	start := time.Now()

	status := &HealthStatus{
		Healthy:   c.connected,
		Latency:   time.Since(start),
		Version:   "1.0.0",
		Timestamp: time.Now(),
	}

	if !c.connected {
		status.Error = "not connected"
		return status, fmt.Errorf("context7 not connected")
	}

	return status, nil
}

func (c *Context7Client) IsConnected() bool {
	return c.connected
}

func (c *Context7Client) ResolveLibrary(ctx context.Context, libraryName string) (*ToolResult, error) {
	if !c.connected {
		return &ToolResult{
			Success: false,
			Error:   "not connected",
		}, fmt.Errorf("context7 not connected")
	}

	if libraryName == "" {
		return &ToolResult{
			Success: false,
			Error:   "library name required",
		}, fmt.Errorf("library name required")
	}

	return &ToolResult{
		Success: true,
		Data: map[string]interface{}{
			"library":   libraryName,
			"libraryId": fmt.Sprintf("/%s/docs", libraryName),
		},
	}, nil
}

func (c *Context7Client) QueryDocs(ctx context.Context, libraryId string, query string) (*ToolResult, error) {
	if !c.connected {
		return &ToolResult{
			Success: false,
			Error:   "not connected",
		}, fmt.Errorf("context7 not connected")
	}

	return &ToolResult{
		Success: true,
		Data: map[string]interface{}{
			"libraryId": libraryId,
			"query":     query,
		},
	}, nil
}
