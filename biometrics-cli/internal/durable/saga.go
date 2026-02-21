package durable

import (
	"fmt"
	"sync"
	"time"
)

type Saga struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Steps       []*SagaStep `json:"steps"`
	Status      string      `json:"status"` // "running", "completed", "compensating", "rolled_back", "failed"
	StartedAt   time.Time   `json:"started_at"`
	CompletedAt *time.Time  `json:"completed_at,omitempty"`
	CurrentStep int         `json:"current_step"`
	mu          sync.RWMutex
}

type SagaStep struct {
	Name         string                 `json:"name"`
	Action       func() error           `json:"-"`
	Compensate   func() error           `json:"-"`
	ActionInput  map[string]interface{} `json:"action_input,omitempty"`
	ActionOutput map[string]interface{} `json:"action_output,omitempty"`
	Status       string                 `json:"status"` // "pending", "running", "completed", "failed", "compensated"
	Error        string                 `json:"error,omitempty"`
	StartedAt    time.Time              `json:"started_at"`
	CompletedAt  *time.Time             `json:"completed_at,omitempty"`
}

type SagaExecutor struct {
	sagas       map[string]*Saga
	sagaDir     string
	persistence *FilePersistence
	mu          sync.RWMutex
}

func NewSagaExecutor(sagaDir string) *SagaExecutor {
	return &SagaExecutor{
		sagas:       make(map[string]*Saga),
		sagaDir:     sagaDir,
		persistence: NewFilePersistence(sagaDir),
	}
}

func (se *SagaExecutor) CreateSaga(name string, steps []*SagaStep) *Saga {
	saga := &Saga{
		ID:          fmt.Sprintf("saga_%d_%s", time.Now().UnixNano(), name),
		Name:        name,
		Steps:       steps,
		Status:      "running",
		StartedAt:   time.Now(),
		CurrentStep: 0,
	}

	se.mu.Lock()
	defer se.mu.Unlock()
	se.sagas[saga.ID] = saga

	return saga
}

func (se *SagaExecutor) Execute(sagaID string) error {
	se.mu.RLock()
	saga, exists := se.sagas[sagaID]
	se.mu.RUnlock()

	if !exists {
		return fmt.Errorf("saga not found: %s", sagaID)
	}

	saga.mu.Lock()
	defer saga.mu.Unlock()

	for saga.CurrentStep < len(saga.Steps) {
		step := saga.Steps[saga.CurrentStep]
		step.Status = "running"
		step.StartedAt = time.Now()

		err := step.Action()
		if err != nil {
			step.Status = "failed"
			step.Error = err.Error()
			saga.Status = "compensating"

			se.compensate(saga)
			return fmt.Errorf("saga failed at step %d: %w", saga.CurrentStep, err)
		}

		step.Status = "completed"
		now := time.Now()
		step.CompletedAt = &now
		saga.CurrentStep++
	}

	saga.Status = "completed"
	now := time.Now()
	saga.CompletedAt = &now

	return nil
}

func (se *SagaExecutor) compensate(saga *Saga) {
	saga.mu.Lock()
	defer saga.mu.Unlock()

	for i := saga.CurrentStep - 1; i >= 0; i-- {
		step := saga.Steps[i]
		if step.Status == "completed" && step.Compensate != nil {
			step.Status = "compensating"

			if err := step.Compensate(); err != nil {
				step.Error = fmt.Sprintf("compensation failed: %v", err)
				continue
			}

			step.Status = "compensated"
			now := time.Now()
			step.CompletedAt = &now
		}
	}

	saga.Status = "rolled_back"
}

func (se *SagaExecutor) GetSaga(sagaID string) (*Saga, error) {
	se.mu.RLock()
	defer se.mu.RUnlock()

	if saga, exists := se.sagas[sagaID]; exists {
		return saga, nil
	}

	return nil, fmt.Errorf("saga not found: %s", sagaID)
}

func (se *SagaExecutor) GetStats() map[string]interface{} {
	se.mu.RLock()
	defer se.mu.RUnlock()

	running := 0
	completed := 0
	rolledBack := 0
	failed := 0

	for _, saga := range se.sagas {
		switch saga.Status {
		case "running":
			running++
		case "completed":
			completed++
		case "rolled_back":
			rolledBack++
		case "failed", "compensating":
			failed++
		}
	}

	return map[string]interface{}{
		"total_sagas": len(se.sagas),
		"running":     running,
		"completed":   completed,
		"rolled_back": rolledBack,
		"failed":      failed,
	}
}
