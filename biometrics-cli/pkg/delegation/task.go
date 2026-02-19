package delegation

import (
	"sync"
	"time"
)

type TaskType string

const (
	TaskTypeCode          TaskType = "code"
	TaskTypeResearch      TaskType = "research"
	TaskTypeArchitecture  TaskType = "architecture"
	TaskTypeDocumentation TaskType = "documentation"
	TaskTypeTesting       TaskType = "testing"
	TaskTypeReview        TaskType = "review"
)

type Priority int

const (
	PriorityCritical Priority = 0
	PriorityHigh     Priority = 1
	PriorityNormal   Priority = 2
	PriorityLow      Priority = 3
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
)

type Task struct {
	ID         string                 `json:"id"`
	Type       TaskType               `json:"type"`
	Priority   Priority               `json:"priority"`
	Payload    interface{}            `json:"payload"`
	Context    map[string]interface{} `json:"context"`
	Status     TaskStatus             `json:"status"`
	CreatedAt  time.Time              `json:"created_at"`
	DeadlineAt time.Time              `json:"deadline_at"`
	RetryCount int                    `json:"retry_count"`
	MaxRetries int                    `json:"max_retries"`
	mu         sync.RWMutex
}

func NewTask(id string, taskType TaskType, priority Priority, payload interface{}) *Task {
	return &Task{
		ID:         id,
		Type:       taskType,
		Priority:   priority,
		Payload:    payload,
		Context:    make(map[string]interface{}),
		Status:     TaskStatusPending,
		CreatedAt:  time.Now(),
		DeadlineAt: time.Now().Add(30 * time.Minute),
		RetryCount: 0,
		MaxRetries: 3,
	}
}

func (t *Task) SetStatus(status TaskStatus) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Status = status
}

func (t *Task) GetStatus() TaskStatus {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Status
}

func (t *Task) SetContext(key string, value interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Context[key] = value
}

func (t *Task) GetContext(key string) interface{} {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Context[key]
}

func (t *Task) IncrementRetry() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.RetryCount++
	return t.RetryCount <= t.MaxRetries
}
