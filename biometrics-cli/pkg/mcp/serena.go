package mcp

import (
	"context"
	"fmt"
	"time"
)

type SerenaClient struct {
	config     *MCPConfig
	connected  bool
	clientID   string
	lastHealth *HealthStatus
}

func NewSerenaClient(config *MCPConfig) *SerenaClient {
	return &SerenaClient{
		config:    config,
		connected: false,
		clientID:  "biometrics-cli",
	}
}

func (s *SerenaClient) Name() string {
	return "serena"
}

func (s *SerenaClient) Connect(ctx context.Context) error {
	if s.connected {
		return nil
	}

	s.connected = true
	return nil
}

func (s *SerenaClient) Disconnect(ctx context.Context) error {
	if !s.connected {
		return nil
	}

	s.connected = false
	return nil
}

func (s *SerenaClient) Health(ctx context.Context) (*HealthStatus, error) {
	start := time.Now()

	status := &HealthStatus{
		Healthy:   s.connected,
		Latency:   time.Since(start),
		Version:   "1.0.0",
		Timestamp: time.Now(),
	}

	if !s.connected {
		status.Error = "not connected"
		return status, fmt.Errorf("serena not connected")
	}

	s.lastHealth = status
	return status, nil
}

func (s *SerenaClient) IsConnected() bool {
	return s.connected
}

func (s *SerenaClient) Orchestrate(ctx context.Context, task string) (*ToolResult, error) {
	if !s.connected {
		return &ToolResult{
			Success: false,
			Error:   "not connected",
		}, fmt.Errorf("serena not connected")
	}

	return &ToolResult{
		Success: true,
		Data: map[string]interface{}{
			"task":     task,
			"status":   "orchestrated",
			"agent_id": s.clientID,
		},
	}, nil
}
