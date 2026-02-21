package heartbeat

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"biometrics-cli/internal/metrics"
)

type Status string

const (
	StatusAlive Status = "alive"
	StatusBusy  Status = "busy"
	StatusIdle  Status = "idle"
	StatusStuck Status = "stuck"
	StatusDead  Status = "dead"
)

type Heartbeat struct {
	AgentID     string                 `json:"agent_id"`
	SessionID   string                 `json:"session_id"`
	Status      Status                 `json:"status"`
	Model       string                 `json:"model,omitempty"`
	Load        float64                `json:"load"`
	TasksDone   int                    `json:"tasks_done"`
	CurrentTask string                 `json:"current_task,omitempty"`
	StartedAt   time.Time              `json:"started_at"`
	LastBeat    time.Time              `json:"last_beat"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type HeartbeatConfig struct {
	Interval       time.Duration
	Timeout        time.Duration
	StuckThreshold time.Duration
	MaxLoad        float64
	PID            int
}

type Monitor struct {
	mu          sync.RWMutex
	heartbeats  map[string]*Heartbeat
	config      *HeartbeatConfig
	alertChan   chan *Alert
	ctx         context.Context
	cancel      context.CancelFunc
	persistence *FilePersistence
}

type Alert struct {
	AgentID   string                 `json:"agent_id"`
	Type      string                 `json:"type"` // "timeout", "stuck", "dead", "high_load"
	Message   string                 `json:"message"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

func NewHeartbeatMonitor(config *HeartbeatConfig) *Monitor {
	ctx, cancel := context.WithCancel(context.Background())

	if config.Interval == 0 {
		config.Interval = 30 * time.Second
	}
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Minute
	}
	if config.StuckThreshold == 0 {
		config.StuckThreshold = 30 * time.Minute
	}

	return &Monitor{
		heartbeats:  make(map[string]*Heartbeat),
		config:      config,
		alertChan:   make(chan *Alert, 100),
		ctx:         ctx,
		cancel:      cancel,
		persistence: NewFilePersistence("./heartbeat-data"),
	}
}

func (m *Monitor) RegisterAgent(agentID, sessionID, model string) *Heartbeat {
	m.mu.Lock()
	defer m.mu.Unlock()

	hb := &Heartbeat{
		AgentID:   agentID,
		SessionID: sessionID,
		Status:    StatusAlive,
		Model:     model,
		Load:      0,
		TasksDone: 0,
		StartedAt: time.Now(),
		LastBeat:  time.Now(),
		Metadata:  make(map[string]interface{}),
	}

	m.heartbeats[agentID] = hb
	metrics.HeartbeatsRegistered.WithLabelValues(agentID).Inc()

	m.persistence.Save(hb)

	return hb
}

func (m *Monitor) Beat(agentID string, status Status, currentTask string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	hb, exists := m.heartbeats[agentID]
	if !exists {
		return fmt.Errorf("agent %s not registered", agentID)
	}

	hb.Status = status
	hb.CurrentTask = currentTask
	hb.LastBeat = time.Now()

	if status == StatusAlive {
		hb.Load = 0
	}

	m.persistence.Save(hb)
	metrics.HeartbeatsReceived.WithLabelValues(agentID).Inc()

	return nil
}

func (m *Monitor) UpdateLoad(agentID string, load float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	hb, exists := m.heartbeats[agentID]
	if !exists {
		return fmt.Errorf("agent %s not registered", agentID)
	}

	hb.Load = load
	hb.LastBeat = time.Now()

	if load > m.config.MaxLoad {
		m.alertChan <- &Alert{
			AgentID:   agentID,
			Type:      "high_load",
			Message:   fmt.Sprintf("Agent %s load exceeded max: %.2f", agentID, load),
			Timestamp: time.Now(),
		}
	}

	return nil
}

func (m *Monitor) IncrementTasksDone(agentID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if hb, exists := m.heartbeats[agentID]; exists {
		hb.TasksDone++
		metrics.HeartbeatTasksDone.WithLabelValues(agentID).Inc()
	}
}

func (m *Monitor) Start(ctx context.Context) {
	go m.runMonitorLoop(ctx)
}

func (m *Monitor) runMonitorLoop(ctx context.Context) {
	ticker := time.NewTicker(m.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.checkHeartbeats()
		}
	}
}

func (m *Monitor) checkHeartbeats() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()

	for agentID, hb := range m.heartbeats {
		sinceLastBeat := now.Sub(hb.LastBeat)

		if sinceLastBeat > m.config.Timeout {
			hb.Status = StatusDead
			m.alertChan <- &Alert{
				AgentID:   agentID,
				Type:      "timeout",
				Message:   fmt.Sprintf("Agent %s heartbeat timeout (>%v)", agentID, m.config.Timeout),
				Timestamp: now,
				Metadata: map[string]interface{}{
					"last_beat": hb.LastBeat,
					"timeout":   m.config.Timeout,
				},
			}
			metrics.HeartbeatTimeouts.WithLabelValues(agentID).Inc()
		}

		if sinceLastBeat > m.config.StuckThreshold && hb.Status == StatusBusy {
			hb.Status = StatusStuck
			m.alertChan <- &Alert{
				AgentID:   agentID,
				Type:      "stuck",
				Message:   fmt.Sprintf("Agent %s stuck on task: %s", agentID, hb.CurrentTask),
				Timestamp: now,
				Metadata: map[string]interface{}{
					"current_task": hb.CurrentTask,
					"duration":     sinceLastBeat,
				},
			}
			metrics.HeartbeatStuckAgents.WithLabelValues(agentID).Inc()
		}

		m.persistence.Save(hb)
	}
}

func (m *Monitor) GetHeartbeat(agentID string) (*Heartbeat, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if hb, exists := m.heartbeats[agentID]; exists {
		return hb, nil
	}

	return nil, fmt.Errorf("agent %s not found", agentID)
}

func (m *Monitor) GetAllHeartbeats() []*Heartbeat {
	m.mu.RLock()
	defer m.mu.RUnlock()

	heartbeats := make([]*Heartbeat, 0, len(m.heartbeats))
	for _, hb := range m.heartbeats {
		heartbeats = append(heartbeats, hb)
	}

	return heartbeats
}

func (m *Monitor) GetAlerts() <-chan *Alert {
	return m.alertChan
}

func (m *Monitor) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	alive := 0
	busy := 0
	idle := 0
	stuck := 0
	dead := 0

	for _, hb := range m.heartbeats {
		switch hb.Status {
		case StatusAlive:
			alive++
		case StatusBusy:
			busy++
		case StatusIdle:
			idle++
		case StatusStuck:
			stuck++
		case StatusDead:
			dead++
		}
	}

	return map[string]interface{}{
		"total_agents": len(m.heartbeats),
		"alive":        alive,
		"busy":         busy,
		"idle":         idle,
		"stuck":        stuck,
		"dead":         dead,
		"config":       m.config,
	}
}

func (m *Monitor) Stop() {
	m.cancel()
}

type FilePersistence struct {
	dir string
}

func NewFilePersistence(dir string) *FilePersistence {
	os.MkdirAll(dir, 0755)
	return &FilePersistence{dir: dir}
}

func (fp *FilePersistence) Save(hb *Heartbeat) error {
	filename := filepath.Join(fp.dir, fmt.Sprintf("%s.json", hb.AgentID))
	data, err := json.MarshalIndent(hb, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
