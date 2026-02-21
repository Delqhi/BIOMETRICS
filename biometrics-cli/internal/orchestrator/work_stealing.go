package orchestrator

import (
	"biometrics-cli/internal/lock"
	"biometrics-cli/internal/metrics"
	"fmt"
	"sync"
	"time"
)

type WorkStealer struct {
	mu             sync.RWMutex
	agentQueues    map[string][]*Task
	agentLoads     map[string]int
	minLoad        int
	maxLoad        int
	stealThreshold int
	stolenTasks    map[string][]*Task
	lock           *lock.DistributedLock
}

type Task struct {
	ID          string
	AgentID     string
	Priority    int
	Payload     interface{}
	CreatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
}

type StealResult struct {
	TaskID    string
	FromAgent string
	ToAgent   string
	Timestamp time.Time
}

func NewWorkStealer(minLoad, maxLoad, stealThreshold int) *WorkStealer {
	return &WorkStealer{
		agentQueues:    make(map[string][]*Task),
		agentLoads:     make(map[string]int),
		minLoad:        minLoad,
		maxLoad:        maxLoad,
		stealThreshold: stealThreshold,
		stolenTasks:    make(map[string][]*Task),
		lock:           lock.NewDistributedLock(""),
	}
}

func (ws *WorkStealer) RegisterAgent(agentID string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if _, exists := ws.agentQueues[agentID]; !exists {
		ws.agentQueues[agentID] = make([]*Task, 0)
		ws.agentLoads[agentID] = 0
	}
}

func (ws *WorkStealer) UnregisterAgent(agentID string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	delete(ws.agentQueues, agentID)
	delete(ws.agentLoads, agentID)
}

func (ws *WorkStealer) Enqueue(agentID string, task *Task) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if _, exists := ws.agentQueues[agentID]; !exists {
		return fmt.Errorf("agent %s not registered", agentID)
	}

	ws.agentQueues[agentID] = append(ws.agentQueues[agentID], task)
	ws.agentLoads[agentID]++

	metrics.WorkStealingTasksEnqueued.WithLabelValues(agentID).Inc()

	return nil
}

func (ws *WorkStealer) Dequeue(agentID string) (*Task, error) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	queue, exists := ws.agentQueues[agentID]
	if !exists {
		return nil, fmt.Errorf("agent %s not registered", agentID)
	}

	if len(queue) == 0 {
		return nil, fmt.Errorf("no tasks in queue for agent %s", agentID)
	}

	task := queue[0]
	ws.agentQueues[agentID] = queue[1:]
	ws.agentLoads[agentID]--

	now := time.Now()
	task.StartedAt = &now

	return task, nil
}

func (ws *WorkStealer) Steal(fromAgent, toAgent string) (*Task, error) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	fromQueue, exists := ws.agentQueues[fromAgent]
	if !exists || len(fromQueue) == 0 {
		return nil, fmt.Errorf("no tasks to steal from agent %s", fromAgent)
	}

	task := fromQueue[len(fromQueue)-1]
	ws.agentQueues[fromAgent] = fromQueue[:len(fromQueue)-1]
	ws.agentLoads[fromAgent]--

	task.AgentID = toAgent
	ws.agentQueues[toAgent] = append(ws.agentQueues[toAgent], task)
	ws.agentLoads[toAgent]++

	if ws.stolenTasks[toAgent] == nil {
		ws.stolenTasks[toAgent] = make([]*Task, 0)
	}
	ws.stolenTasks[toAgent] = append(ws.stolenTasks[toAgent], task)

	metrics.WorkStealingTasksStolen.Inc()

	return task, nil
}

func (ws *WorkStealer) FindVictims() []string {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	var victims []string
	for agentID, load := range ws.agentLoads {
		if load > ws.stealThreshold {
			victims = append(victims, agentID)
		}
	}

	return victims
}

func (ws *WorkStealer) GetLoad(agentID string) int {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.agentLoads[agentID]
}

func (ws *WorkStealer) GetTotalLoad() int {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	total := 0
	for _, load := range ws.agentLoads {
		total += load
	}
	return total
}

func (ws *WorkStealer) GetAgentLoads() map[string]int {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	loads := make(map[string]int)
	for agentID, load := range ws.agentLoads {
		loads[agentID] = load
	}
	return loads
}

func (ws *WorkStealer) Rebalance() []StealResult {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	var results []StealResult

	victims := ws.findVictimsLocked()
	receivers := ws.findReceiversLocked()

	for _, victim := range victims {
		for _, receiver := range receivers {
			victimLoad := ws.agentLoads[victim]
			receiverLoad := ws.agentLoads[receiver]

			if victimLoad-receiverLoad > ws.stealThreshold {
				task, err := ws.stealTaskLocked(victim, receiver)
				if err == nil {
					results = append(results, StealResult{
						TaskID:    task.ID,
						FromAgent: victim,
						ToAgent:   receiver,
						Timestamp: time.Now(),
					})
				}
			}
		}
	}

	return results
}

func (ws *WorkStealer) findVictimsLocked() []string {
	var victims []string
	for agentID, load := range ws.agentLoads {
		if load > ws.minLoad+ws.stealThreshold {
			victims = append(victims, agentID)
		}
	}
	return victims
}

func (ws *WorkStealer) findReceiversLocked() []string {
	var receivers []string
	for agentID, load := range ws.agentLoads {
		if load < ws.minLoad {
			receivers = append(receivers, agentID)
		}
	}
	return receivers
}

func (ws *WorkStealer) stealTaskLocked(fromAgent, toAgent string) (*Task, error) {
	queue := ws.agentQueues[fromAgent]
	if len(queue) == 0 {
		return nil, fmt.Errorf("no tasks to steal")
	}

	task := queue[len(queue)-1]
	ws.agentQueues[fromAgent] = queue[:len(queue)-1]
	ws.agentLoads[fromAgent]--

	task.AgentID = toAgent
	ws.agentQueues[toAgent] = append(ws.agentQueues[toAgent], task)
	ws.agentLoads[toAgent]++

	return task, nil
}

func (ws *WorkStealer) StartRebalancer(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				results := ws.Rebalance()
				if len(results) > 0 {
					metrics.WorkStealingRebalances.Inc()
				}
			}
		}
	}()
}

func (ws *WorkStealer) Stop() {
}

func (ws *WorkStealer) GetStats() map[string]interface{} {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	return map[string]interface{}{
		"total_agents":    len(ws.agentQueues),
		"total_tasks":     ws.GetTotalLoad(),
		"agent_loads":     ws.agentLoads,
		"stolen_tasks":    len(ws.stolenTasks),
		"steal_threshold": ws.stealThreshold,
	}
}
