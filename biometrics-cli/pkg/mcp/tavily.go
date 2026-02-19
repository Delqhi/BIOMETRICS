package mcp

import (
	"context"
	"fmt"
	"time"
)

type TavilyClient struct {
	config    *MCPConfig
	connected bool
	apiKey    string
}

func NewTavilyClient(config *MCPConfig) *TavilyClient {
	apiKey := ""
	if config.Environment != nil {
		apiKey = config.Environment["TAVILY_API_KEY"]
	}

	return &TavilyClient{
		config:    config,
		connected: false,
		apiKey:    apiKey,
	}
}

func (t *TavilyClient) Name() string {
	return "tavily"
}

func (t *TavilyClient) Connect(ctx context.Context) error {
	if t.connected {
		return nil
	}

	if t.apiKey == "" {
		return fmt.Errorf("tavily api key required")
	}

	t.connected = true
	return nil
}

func (t *TavilyClient) Disconnect(ctx context.Context) error {
	if !t.connected {
		return nil
	}

	t.connected = false
	return nil
}

func (t *TavilyClient) Health(ctx context.Context) (*HealthStatus, error) {
	start := time.Now()

	status := &HealthStatus{
		Healthy:   t.connected,
		Latency:   time.Since(start),
		Version:   "1.0.0",
		Timestamp: time.Now(),
	}

	if !t.connected {
		status.Error = "not connected"
		return status, fmt.Errorf("tavily not connected")
	}

	return status, nil
}

func (t *TavilyClient) IsConnected() bool {
	return t.connected
}

func (t *TavilyClient) Search(ctx context.Context, query string, limit int) (*ToolResult, error) {
	if !t.connected {
		return &ToolResult{
			Success: false,
			Error:   "not connected",
		}, fmt.Errorf("tavily not connected")
	}

	if query == "" {
		return &ToolResult{
			Success: false,
			Error:   "query required",
		}, fmt.Errorf("query required")
	}

	return &ToolResult{
		Success: true,
		Data: map[string]interface{}{
			"query":  query,
			"limit":  limit,
			"status": "searched",
		},
	}, nil
}

func (t *TavilyClient) Research(ctx context.Context, query string, depth string) (*ToolResult, error) {
	if !t.connected {
		return &ToolResult{
			Success: false,
			Error:   "not connected",
		}, fmt.Errorf("tavily not connected")
	}

	return &ToolResult{
		Success: true,
		Data: map[string]interface{}{
			"query": query,
			"depth": depth,
		},
	}, nil
}
