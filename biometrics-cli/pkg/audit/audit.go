package audit

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Auditor struct {
	config     *AuditConfig
	storage    AuditStorage
	eventQueue chan *AuditEvent
	wg         sync.WaitGroup
	stopChan   chan struct{}
	mu         sync.RWMutex
	eventID    uint64
}

type AuditConfig struct {
	StoragePath       string
	StorageType       StorageType
	MaxSize           int64
	RetentionDays     int
	EnableCompression bool
	EnableEncryption  bool
	FlushInterval     time.Duration
	QueueSize         int
}

func DefaultAuditConfig() *AuditConfig {
	return &AuditConfig{
		StoragePath:       "/tmp/biometrics/audit",
		StorageType:       StorageTypeFile,
		MaxSize:           100 * 1024 * 1024,
		RetentionDays:     90,
		EnableCompression: true,
		EnableEncryption:  false,
		FlushInterval:     5 * time.Second,
		QueueSize:         1000,
	}
}

func NewAuditor(config *AuditConfig) (*Auditor, error) {
	if config == nil {
		config = DefaultAuditConfig()
	}

	// Ensure FlushInterval is not zero to prevent ticker panic
	if config.FlushInterval <= 0 {
		config.FlushInterval = 5 * time.Second
	}

	// Ensure QueueSize has a reasonable minimum
	if config.QueueSize <= 0 {
		config.QueueSize = 1000
	}

	storage, err := NewAuditStorage(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create audit storage: %w", err)
	}

	auditor := &Auditor{
		config:     config,
		storage:    storage,
		eventQueue: make(chan *AuditEvent, config.QueueSize),
		stopChan:   make(chan struct{}),
	}

	auditor.wg.Add(1)
	go auditor.processEvents()

	auditor.wg.Add(1)
	go auditor.flushPeriodically()

	return auditor, nil
}

func (a *Auditor) Log(eventType EventType, actor string, action string, resource string, metadata map[string]interface{}) error {
	event := &AuditEvent{
		ID:        a.generateEventID(),
		Timestamp: time.Now().UTC(),
		EventType: eventType,
		Actor:     actor,
		Action:    action,
		Resource:  resource,
		Metadata:  metadata,
	}

	select {
	case a.eventQueue <- event:
		return nil
	default:
		return fmt.Errorf("audit queue full, event dropped")
	}
}

func (a *Auditor) LogAuthentication(userID string, success bool, method string, ip string) error {
	eventType := EventAuthSuccess
	if !success {
		eventType = EventAuthFailure
	}

	return a.Log(eventType, userID, "authentication", "auth-system", map[string]interface{}{
		"method": method,
		"ip":     ip,
	})
}

func (a *Auditor) LogAuthorization(userID string, resource string, action string, allowed bool) error {
	eventType := EventAuthzGranted
	if !allowed {
		eventType = EventAuthzDenied
	}

	return a.Log(eventType, userID, action, resource, map[string]interface{}{
		"allowed": allowed,
	})
}

func (a *Auditor) LogDataAccess(userID string, dataType string, operation string, recordID string) error {
	return a.Log(EventDataAccess, userID, operation, dataType, map[string]interface{}{
		"record_id": recordID,
	})
}

func (a *Auditor) LogSecurityEvent(eventType EventType, actor string, details string, severity string) error {
	return a.Log(eventType, actor, "security-event", "security-system", map[string]interface{}{
		"details":  details,
		"severity": severity,
	})
}

func (a *Auditor) LogSystemEvent(eventType EventType, component string, action string, status string) error {
	return a.Log(eventType, "system", action, component, map[string]interface{}{
		"status": status,
	})
}

func (a *Auditor) processEvents() {
	defer a.wg.Done()

	for {
		select {
		case event := <-a.eventQueue:
			if err := a.storage.Store(event); err != nil {
				fmt.Printf("Failed to store audit event: %v\n", err)
			}
		case <-a.stopChan:
			drainQueue(a.eventQueue, a.storage)
			return
		}
	}
}

func (a *Auditor) flushPeriodically() {
	defer a.wg.Done()

	ticker := time.NewTicker(a.config.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := a.storage.Flush(); err != nil {
				fmt.Printf("Failed to flush audit storage: %v\n", err)
			}
		case <-a.stopChan:
			return
		}
	}
}

func (a *Auditor) Query(ctx context.Context, query *AuditQuery) (*AuditQueryResult, error) {
	return a.storage.Query(ctx, query)
}

func (a *Auditor) Export(startTime, endTime time.Time, format ExportFormat) ([]byte, error) {
	return a.storage.Export(startTime, endTime, format)
}

func (a *Auditor) GetStats() (*AuditStats, error) {
	return a.storage.GetStats()
}

func (a *Auditor) RotateLogs() error {
	return a.storage.Rotate()
}

func (a *Auditor) Cleanup() error {
	cutoff := time.Now().AddDate(0, 0, -a.config.RetentionDays)
	return a.storage.Cleanup(cutoff)
}

func (a *Auditor) Stop() {
	close(a.stopChan)
	a.wg.Wait()
	if err := a.storage.Close(); err != nil {
		fmt.Printf("Failed to close audit storage: %v\n", err)
	}
}

func (a *Auditor) generateEventID() uint64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.eventID++
	return a.eventID
}

func CreateAuditHash(event *AuditEvent) (string, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

func VerifyAuditIntegrity(event *AuditEvent, expectedHash string) (bool, error) {
	actualHash, err := CreateAuditHash(event)
	if err != nil {
		return false, err
	}

	return actualHash == expectedHash, nil
}

func drainQueue(queue chan *AuditEvent, storage AuditStorage) {
	for {
		select {
		case event := <-queue:
			if err := storage.Store(event); err != nil {
				fmt.Printf("Failed to drain audit event: %v\n", err)
			}
		default:
			return
		}
	}
}

func (a *Auditor) GetEventQueueSize() int {
	return len(a.eventQueue)
}

func (a *Auditor) IsHealthy() bool {
	return a.storage != nil && !a.storage.IsClosed()
}
