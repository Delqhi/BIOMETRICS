package mcp

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type MCPRegistry struct {
	clients map[string]MCPClient
	configs map[string]*MCPConfig
	mu      sync.RWMutex
}

func NewMCPRegistry() *MCPRegistry {
	return &MCPRegistry{
		clients: make(map[string]MCPClient),
		configs: make(map[string]*MCPConfig),
	}
}

func (r *MCPRegistry) Register(config *MCPConfig, client MCPClient) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !config.Enabled {
		return nil
	}

	r.configs[config.Name] = config
	r.clients[config.Name] = client

	return nil
}

func (r *MCPRegistry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.clients[name]; !exists {
		return fmt.Errorf("mcp %s not found", name)
	}

	delete(r.clients, name)
	delete(r.configs, name)

	return nil
}

func (r *MCPRegistry) GetClient(name string) (MCPClient, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	client, exists := r.clients[name]
	if !exists {
		return nil, fmt.Errorf("mcp %s not found", name)
	}

	return client, nil
}

func (r *MCPRegistry) ConnectAll(ctx context.Context) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for name, client := range r.clients {
		if err := client.Connect(ctx); err != nil {
			return fmt.Errorf("failed to connect %s: %w", name, err)
		}
	}

	return nil
}

func (r *MCPRegistry) DisconnectAll(ctx context.Context) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for name, client := range r.clients {
		if err := client.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to disconnect %s: %w", name, err)
		}
	}

	return nil
}

func (r *MCPRegistry) HealthCheck(ctx context.Context) map[string]*HealthStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make(map[string]*HealthStatus)

	for name, client := range r.clients {
		status, err := client.Health(ctx)
		if err != nil {
			results[name] = &HealthStatus{
				Healthy:   false,
				Latency:   0,
				Version:   "unknown",
				Timestamp: time.Now(),
				Error:     err.Error(),
			}
		} else {
			results[name] = status
		}
	}

	return results
}

func (r *MCPRegistry) ListClients() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.clients))
	for name := range r.clients {
		names = append(names, name)
	}

	return names
}

func (r *MCPRegistry) GetActiveClients() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, client := range r.clients {
		if client.IsConnected() {
			count++
		}
	}

	return count
}
