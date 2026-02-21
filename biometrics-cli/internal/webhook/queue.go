package webhook

import (
	"biometrics-cli/internal/circuit"
	"biometrics-cli/internal/metrics"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type EventQueue struct {
	mu             sync.RWMutex
	queues         map[string][]*QueuedEvent
	processing     map[string]bool
	workers        int
	bufferSize     int
	circuitBreaker *circuit.CircuitBreaker
	metrics        *QueueMetrics
	ctx            context.Context
	cancel         context.CancelFunc
}

type QueuedEvent struct {
	ID          string          `json:"id"`
	Event       string          `json:"event"`
	Payload     json.RawMessage `json:"payload"`
	Retries     int             `json:"retries"`
	Priority    int             `json:"priority"`
	CreatedAt   time.Time       `json:"created_at"`
	ProcessedAt *time.Time      `json:"processed_at,omitempty"`
}

type QueueMetrics struct {
	Enqueued  int64 `json:"enqueued"`
	Dequeued  int64 `json:"dequeued"`
	Failed    int64 `json:"failed"`
	Processed int64 `json:"processed"`
	Retried   int64 `json:"retried"`
}

var globalQueue *EventQueue

func NewEventQueue(workers int, bufferSize int) *EventQueue {
	ctx, cancel := context.WithCancel(context.Background())

	breaker := circuit.GetOrCreate("webhook-queue", &circuit.CircuitBreakerConfig{
		Name:         "webhook-queue",
		MaxFailures:  5,
		Timeout:      30 * time.Second,
		HalfOpenMax:  3,
		ResetTimeout: 60 * time.Second,
	})

	return &EventQueue{
		queues:         make(map[string][]*QueuedEvent),
		processing:     make(map[string]bool),
		workers:        workers,
		bufferSize:     bufferSize,
		circuitBreaker: breaker,
		metrics:        &QueueMetrics{},
		ctx:            ctx,
		cancel:         cancel,
	}
}

func GetGlobalQueue() *EventQueue {
	if globalQueue == nil {
		globalQueue = NewEventQueue(10, 1000)
	}
	return globalQueue
}

func (eq *EventQueue) Enqueue(event string, payload json.RawMessage, priority int) error {
	eq.mu.Lock()
	defer eq.mu.Unlock()

	queuedEvent := &QueuedEvent{
		ID:        fmt.Sprintf("%d-%s", time.Now().UnixNano(), event),
		Event:     event,
		Payload:   payload,
		Retries:   0,
		Priority:  priority,
		CreatedAt: time.Now(),
	}

	if _, exists := eq.queues[event]; !exists {
		eq.queues[event] = make([]*QueuedEvent, 0)
	}

	eq.queues[event] = append(eq.queues[event], queuedEvent)
	eq.metrics.Enqueued++

	metrics.WebhookQueueDepth.Set(float64(eq.TotalDepth()))

	return nil
}

func (eq *EventQueue) Dequeue(event string) (*QueuedEvent, error) {
	eq.mu.Lock()
	defer eq.mu.Unlock()

	queue, exists := eq.queues[event]
	if !exists || len(queue) == 0 {
		return nil, fmt.Errorf("no events in queue for: %s", event)
	}

	eventObj := queue[0]
	eq.queues[event] = queue[1:]
	eq.metrics.Dequeued++

	return eventObj, nil
}

func (eq *EventQueue) Start(ctx context.Context, handler func(*QueuedEvent) error) {
	for i := 0; i < eq.workers; i++ {
		workerID := i
		go eq.runWorker(ctx, workerID, handler)
	}
}

func (eq *EventQueue) runWorker(ctx context.Context, workerID int, handler func(*QueuedEvent) error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			eq.mu.RLock()
			events := make([]string, 0, len(eq.queues))
			for event := range eq.queues {
				if !eq.processing[event] && len(eq.queues[event]) > 0 {
					events = append(events, event)
				}
			}
			eq.mu.RUnlock()

			for _, event := range events {
				if err := eq.processEvent(event, handler); err != nil {
					eq.metrics.Failed++
				}
			}
		}
	}
}

func (eq *EventQueue) processEvent(event string, handler func(*QueuedEvent) error) error {
	eventObj, err := eq.Dequeue(event)
	if err != nil {
		return err
	}

	eq.mu.Lock()
	eq.processing[event] = true
	eq.mu.Unlock()

	defer func() {
		eq.mu.Lock()
		eq.processing[event] = false
		eq.mu.Unlock()
	}()

	err = eq.circuitBreaker.Execute(func() error {
		return handler(eventObj)
	})

	if err != nil {
		eq.metrics.Failed++
		if eventObj.Retries < 3 {
			eventObj.Retries++
			eq.metrics.Retried++
			return eq.Enqueue(event, eventObj.Payload, eventObj.Priority)
		}
		return fmt.Errorf("event %s failed after 3 retries: %w", eventObj.ID, err)
	}

	now := time.Now()
	eventObj.ProcessedAt = &now
	eq.metrics.Processed++

	return nil
}

func (eq *EventQueue) Depth(event string) int {
	eq.mu.RLock()
	defer eq.mu.RUnlock()
	return len(eq.queues[event])
}

func (eq *EventQueue) TotalDepth() int {
	eq.mu.RLock()
	defer eq.mu.RUnlock()

	total := 0
	for _, queue := range eq.queues {
		total += len(queue)
	}
	return total
}

func (eq *EventQueue) GetMetrics() *QueueMetrics {
	eq.mu.RLock()
	defer eq.mu.RUnlock()
	return &QueueMetrics{
		Enqueued:  eq.metrics.Enqueued,
		Dequeued:  eq.metrics.Dequeued,
		Failed:    eq.metrics.Failed,
		Processed: eq.metrics.Processed,
		Retried:   eq.metrics.Retried,
	}
}

func (eq *EventQueue) Stop() {
	eq.cancel()
}

func (eq *EventQueue) IsHealthy() bool {
	state := eq.circuitBreaker.GetState()
	return state == circuit.StateClosed
}
