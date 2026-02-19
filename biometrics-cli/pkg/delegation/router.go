package delegation

import (
	"errors"
	"sync"
	"time"
)

type AgentCapability struct {
	Name         string   `json:"name"`
	AgentID      string   `json:"agent_id"`
	Capabilities []string `json:"capabilities"`
	Load         int      `json:"load"`
	Healthy      bool     `json:"healthy"`
}

type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

type CircuitBreaker struct {
	state     CircuitState
	failures  int
	threshold int
	timeout   time.Duration
	lastFail  time.Time
	mu        sync.RWMutex
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:     CircuitClosed,
		threshold: threshold,
		timeout:   timeout,
	}
}

func (cb *CircuitBreaker) CanExecute() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		if time.Since(cb.lastFail) > cb.timeout {
			return true
		}
		return false
	case CircuitHalfOpen:
		return true
	}
	return false
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	cb.state = CircuitClosed
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	cb.lastFail = time.Now()

	if cb.failures >= cb.threshold {
		cb.state = CircuitOpen
	}
}

type DelegationRouter struct {
	agents          map[string]*AgentCapability
	circuitBreakers map[string]*CircuitBreaker
	affinityMap     map[string]string
	mu              sync.RWMutex
}

func NewDelegationRouter() *DelegationRouter {
	return &DelegationRouter{
		agents:          make(map[string]*AgentCapability),
		circuitBreakers: make(map[string]*CircuitBreaker),
		affinityMap:     make(map[string]string),
	}
}

func (r *DelegationRouter) RegisterAgent(agent *AgentCapability) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents[agent.AgentID] = agent
	r.circuitBreakers[agent.AgentID] = NewCircuitBreaker(3, 30*time.Second)
}

func (r *DelegationRouter) Route(task *Task) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var candidates []*AgentCapability

	for _, agent := range r.agents {
		if !agent.Healthy {
			continue
		}

		cb := r.circuitBreakers[agent.AgentID]
		if !cb.CanExecute() {
			continue
		}

		if r.hasCapability(agent, string(task.Type)) {
			candidates = append(candidates, agent)
		}
	}

	if len(candidates) == 0 {
		return "", errors.New("no available agents with required capability")
	}

	selected := r.selectBestAgent(candidates, task)
	r.affinityMap[task.ID] = selected.AgentID

	return selected.AgentID, nil
}

func (r *DelegationRouter) hasCapability(agent *AgentCapability, capability string) bool {
	for _, cap := range agent.Capabilities {
		if cap == capability {
			return true
		}
	}
	return false
}

func (r *DelegationRouter) selectBestAgent(candidates []*AgentCapability, task *Task) *AgentCapability {
	var best *AgentCapability
	minLoad := int(^uint(0) >> 1)

	for _, candidate := range candidates {
		if candidate.Load < minLoad {
			minLoad = candidate.Load
			best = candidate
		}
	}

	return best
}

func (r *DelegationRouter) GetAffinityAgent(taskID string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.affinityMap[taskID]
}

func (r *DelegationRouter) UpdateAgentLoad(agentID string, load int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if agent, exists := r.agents[agentID]; exists {
		agent.Load = load
	}
}

func (r *DelegationRouter) MarkAgentHealthy(agentID string, healthy bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if agent, exists := r.agents[agentID]; exists {
		agent.Healthy = healthy
	}
}

func (r *DelegationRouter) RecordSuccess(agentID string) {
	r.mu.RLock()
	cb := r.circuitBreakers[agentID]
	r.mu.RUnlock()

	if cb != nil {
		cb.RecordSuccess()
	}
}

func (r *DelegationRouter) RecordFailure(agentID string) {
	r.mu.RLock()
	cb := r.circuitBreakers[agentID]
	r.mu.RUnlock()

	if cb != nil {
		cb.RecordFailure()
	}
}
