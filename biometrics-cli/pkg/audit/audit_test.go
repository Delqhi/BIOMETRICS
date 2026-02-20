package audit

import (
	"context"
	"testing"
	"time"
)

func TestNewAuditor(t *testing.T) {
	tmpDir := t.TempDir()
	config := &AuditConfig{
		StoragePath:   tmpDir,
		StorageType:   StorageTypeFile,
		FlushInterval: 100 * time.Millisecond,
		QueueSize:     100,
		RetentionDays: 30,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		t.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	if auditor.storage == nil {
		t.Error("Storage should not be nil")
	}

	if auditor.eventQueue == nil {
		t.Error("Event queue should not be nil")
	}
}

func TestDefaultAuditConfig(t *testing.T) {
	config := DefaultAuditConfig()

	if config.StoragePath == "" {
		t.Error("StoragePath should not be empty")
	}

	if config.MaxSize <= 0 {
		t.Error("MaxSize should be positive")
	}

	if config.RetentionDays <= 0 {
		t.Error("RetentionDays should be positive")
	}
}

func TestAuditor_Log(t *testing.T) {
	tmpDir := t.TempDir()
	config := &AuditConfig{
		StoragePath:   tmpDir,
		StorageType:   StorageTypeMemory,
		QueueSize:     100,
		FlushInterval: time.Hour,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		t.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	err = auditor.Log(EventAuthSuccess, "user123", "login", "auth-system", map[string]interface{}{
		"ip": "192.168.1.1",
	})
	if err != nil {
		t.Errorf("Failed to log event: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	stats, err := auditor.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.TotalEvents != 1 {
		t.Errorf("Expected 1 event, got %d", stats.TotalEvents)
	}
}

func TestAuditor_LogAuthentication(t *testing.T) {
	tmpDir := t.TempDir()
	config := &AuditConfig{
		StoragePath:   tmpDir,
		StorageType:   StorageTypeMemory,
		QueueSize:     100,
		FlushInterval: time.Hour,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		t.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	err = auditor.LogAuthentication("user123", true, "password", "192.168.1.1")
	if err != nil {
		t.Errorf("Failed to log authentication: %v", err)
	}

	err = auditor.LogAuthentication("user456", false, "password", "192.168.1.2")
	if err != nil {
		t.Errorf("Failed to log failed authentication: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	stats, err := auditor.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.TotalEvents != 2 {
		t.Errorf("Expected 2 events, got %d", stats.TotalEvents)
	}
}

func TestAuditor_LogAuthorization(t *testing.T) {
	tmpDir := t.TempDir()
	config := &AuditConfig{
		StoragePath: tmpDir,
		StorageType: StorageTypeMemory,
		QueueSize:   100,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		t.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	err = auditor.LogAuthorization("user123", "/api/data", "read", true)
	if err != nil {
		t.Errorf("Failed to log authorization: %v", err)
	}

	err = auditor.LogAuthorization("user456", "/api/admin", "write", false)
	if err != nil {
		t.Errorf("Failed to log denied authorization: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	stats, err := auditor.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.TotalEvents != 2 {
		t.Errorf("Expected 2 events, got %d", stats.TotalEvents)
	}
}

func TestAuditor_LogDataAccess(t *testing.T) {
	tmpDir := t.TempDir()
	config := &AuditConfig{
		StoragePath: tmpDir,
		StorageType: StorageTypeMemory,
		QueueSize:   100,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		t.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	err = auditor.LogDataAccess("user123", "users", "read", "record-456")
	if err != nil {
		t.Errorf("Failed to log data access: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	stats, err := auditor.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.TotalEvents != 1 {
		t.Errorf("Expected 1 event, got %d", stats.TotalEvents)
	}
}

func TestAuditor_Query(t *testing.T) {
	tmpDir := t.TempDir()
	config := &AuditConfig{
		StoragePath:   tmpDir,
		StorageType:   StorageTypeMemory,
		QueueSize:     100,
		FlushInterval: time.Hour,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		t.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	for i := 0; i < 10; i++ {
		eventType := EventAuthSuccess
		if i%2 == 0 {
			eventType = EventAuthFailure
		}

		auditor.Log(eventType, "user123", "login", "auth-system", nil)
	}

	time.Sleep(100 * time.Millisecond)

	query := &AuditQuery{
		EventTypes: []EventType{EventAuthSuccess},
		Limit:      5,
	}

	result, err := auditor.Query(context.Background(), query)
	if err != nil {
		t.Fatalf("Failed to query events: %v", err)
	}

	if result.TotalCount != 5 {
		t.Errorf("Expected 5 events, got %d", result.TotalCount)
	}

	if len(result.Events) != 5 {
		t.Errorf("Expected 5 events in result, got %d", len(result.Events))
	}
}

func TestAuditor_Export(t *testing.T) {
	tmpDir := t.TempDir()
	config := &AuditConfig{
		StoragePath: tmpDir,
		StorageType: StorageTypeMemory,
		QueueSize:   100,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		t.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	now := time.Now()
	auditor.Log(EventAuthSuccess, "user123", "login", "auth-system", nil)
	auditor.Log(EventDataAccess, "user123", "read", "data", nil)

	time.Sleep(50 * time.Millisecond)

	data, err := auditor.Export(now.Add(-time.Hour), now.Add(time.Hour), ExportFormatJSON)
	if err != nil {
		t.Fatalf("Failed to export: %v", err)
	}

	if len(data) == 0 {
		t.Error("Exported data should not be empty")
	}
}

func TestAuditor_GetStats(t *testing.T) {
	tmpDir := t.TempDir()
	config := &AuditConfig{
		StoragePath: tmpDir,
		StorageType: StorageTypeMemory,
		QueueSize:   100,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		t.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	for i := 0; i < 5; i++ {
		auditor.Log(EventAuthSuccess, "user123", "login", "auth-system", nil)
	}

	for i := 0; i < 3; i++ {
		auditor.Log(EventDataAccess, "user456", "read", "data", nil)
	}

	time.Sleep(50 * time.Millisecond)

	stats, err := auditor.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.TotalEvents != 8 {
		t.Errorf("Expected 8 total events, got %d", stats.TotalEvents)
	}

	if stats.EventsByType["authentication.success"] != 5 {
		t.Errorf("Expected 5 auth events, got %d", stats.EventsByType["authentication.success"])
	}

	if stats.EventsByActor["user123"] != 5 {
		t.Errorf("Expected 5 events from user123, got %d", stats.EventsByActor["user123"])
	}
}

func TestAuditor_QueueFull(t *testing.T) {
	tmpDir := t.TempDir()
	config := &AuditConfig{
		StoragePath:   tmpDir,
		StorageType:   StorageTypeMemory,
		QueueSize:     5,
		FlushInterval: time.Hour,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		t.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	for i := 0; i < 10; i++ {
		err := auditor.Log(EventAuthSuccess, "user123", "login", "auth-system", nil)
		if i < 5 && err != nil {
			t.Errorf("Should not get error for first 5 events: %v", err)
		}
	}

	queueSize := auditor.GetEventQueueSize()
	if queueSize > 5 {
		t.Errorf("Queue size should not exceed capacity, got %d", queueSize)
	}
}

func TestAuditor_Stop(t *testing.T) {
	tmpDir := t.TempDir()
	config := &AuditConfig{
		StoragePath: tmpDir,
		StorageType: StorageTypeMemory,
		QueueSize:   100,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		t.Fatalf("Failed to create auditor: %v", err)
	}

	auditor.Log(EventAuthSuccess, "user123", "login", "auth-system", nil)
	auditor.Stop()

	if !auditor.storage.IsClosed() {
		t.Error("Storage should be closed after stop")
	}
}

func TestEventTypeMethods(t *testing.T) {
	eventType := EventAuthSuccess

	if eventType.String() != "authentication.success" {
		t.Errorf("Wrong string representation: %s", eventType.String())
	}

	if eventType.Category() != "authentication" {
		t.Errorf("Wrong category: %s", eventType.Category())
	}

	if eventType.IsSecurityCritical() {
		t.Error("Auth success should NOT be security critical")
	}

	failureEvent := EventAuthFailure
	if !failureEvent.IsSecurityCritical() {
		t.Error("Auth failure should be security critical")
	}
}

func TestGetSeverityForEvent(t *testing.T) {
	tests := []struct {
		eventType EventType
		expected  SeverityLevel
	}{
		{EventSecurityViolation, SeverityCritical},
		{EventAuthFailure, SeverityHigh},
		{EventDataDelete, SeverityMedium},
		{EventAuthSuccess, SeverityLow},
	}

	for _, test := range tests {
		severity := GetSeverityForEvent(test.eventType)
		if severity != test.expected {
			t.Errorf("Wrong severity for %s: expected %s, got %s", test.eventType, test.expected, severity)
		}
	}
}

func TestCreateAuditHash(t *testing.T) {
	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "user123",
	}

	hash1, err := CreateAuditHash(event)
	if err != nil {
		t.Fatalf("Failed to create hash: %v", err)
	}

	if len(hash1) != 64 {
		t.Errorf("Hash should be 64 characters, got %d", len(hash1))
	}

	hash2, err := CreateAuditHash(event)
	if err != nil {
		t.Fatalf("Failed to create second hash: %v", err)
	}

	if hash1 != hash2 {
		t.Error("Same event should produce same hash")
	}
}

func TestVerifyAuditIntegrity(t *testing.T) {
	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "user123",
	}

	hash, err := CreateAuditHash(event)
	if err != nil {
		t.Fatalf("Failed to create hash: %v", err)
	}

	valid, err := VerifyAuditIntegrity(event, hash)
	if err != nil {
		t.Fatalf("Failed to verify integrity: %v", err)
	}

	if !valid {
		t.Error("Event should be valid")
	}

	event.Actor = "modified"
	valid, err = VerifyAuditIntegrity(event, hash)
	if err != nil {
		t.Fatalf("Failed to verify modified event: %v", err)
	}

	if valid {
		t.Error("Modified event should not be valid")
	}
}

func TestMemoryStorage(t *testing.T) {
	config := DefaultAuditConfig()
	config.StorageType = StorageTypeMemory

	storage, err := NewAuditStorage(config)
	if err != nil {
		t.Fatalf("Failed to create memory storage: %v", err)
	}
	defer storage.Close()

	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "user123",
	}

	err = storage.Store(event)
	if err != nil {
		t.Errorf("Failed to store event: %v", err)
	}

	stats, err := storage.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.TotalEvents != 1 {
		t.Errorf("Expected 1 event, got %d", stats.TotalEvents)
	}
}

func TestFileStorage(t *testing.T) {
	tmpDir := t.TempDir()
	config := DefaultAuditConfig()
	config.StoragePath = tmpDir
	config.StorageType = StorageTypeFile
	config.EnableCompression = false

	storage, err := NewAuditStorage(config)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}
	defer storage.Close()

	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "user123",
	}

	err = storage.Store(event)
	if err != nil {
		t.Errorf("Failed to store event: %v", err)
	}

	err = storage.Flush()
	if err != nil {
		t.Errorf("Failed to flush: %v", err)
	}

	stats, err := storage.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.TotalEvents != 1 {
		t.Errorf("Expected 1 event, got %d", stats.TotalEvents)
	}
}

func TestAuditor_Health(t *testing.T) {
	tmpDir := t.TempDir()
	config := &AuditConfig{
		StoragePath: tmpDir,
		StorageType: StorageTypeMemory,
		QueueSize:   100,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		t.Fatalf("Failed to create auditor: %v", err)
	}

	if !auditor.IsHealthy() {
		t.Error("Auditor should be healthy")
	}

	auditor.Stop()

	if auditor.IsHealthy() {
		t.Error("Auditor should not be healthy after stop")
	}
}

func TestAuditEvent_Validate(t *testing.T) {
	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "user1",
		Action:    "login",
		Resource:  "auth",
	}

	if err := event.Validate(); err != nil {
		t.Errorf("Valid event should pass validation: %v", err)
	}

	invalidEvent := &AuditEvent{}
	if err := invalidEvent.Validate(); err == nil {
		t.Error("Invalid event should fail validation")
	}
}

func TestAuditEvent_ToJSON(t *testing.T) {
	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "user1",
		Action:    "login",
		Resource:  "auth",
		Metadata:  map[string]interface{}{"ip": "127.0.0.1"},
	}

	data, err := event.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("JSON output should not be empty")
	}

	str, err := event.ToJSONString()
	if err != nil {
		t.Fatalf("ToJSONString failed: %v", err)
	}

	if str == "" {
		t.Error("JSON string should not be empty")
	}
}

func TestAuditEvent_GetDetail(t *testing.T) {
	event := &AuditEvent{
		Metadata: map[string]interface{}{
			"ip":      "127.0.0.1",
			"count":   42,
			"enabled": true,
		},
	}

	val, ok := event.GetDetail("ip")
	if !ok || val != "127.0.0.1" {
		t.Error("GetDetail should return correct value")
	}

	str, ok := event.GetDetailString("ip")
	if !ok || str != "127.0.0.1" {
		t.Error("GetDetailString should return correct string")
	}

	num, ok := event.GetDetailInt("count")
	if !ok || num != 42 {
		t.Error("GetDetailInt should return correct int")
	}

	_, ok = event.GetDetail("nonexistent")
	if ok {
		t.Error("GetDetail should return false for nonexistent key")
	}
}

func TestAuditEvent_AddDetail(t *testing.T) {
	event := &AuditEvent{}
	event.AddDetail("key1", "value1")
	event.AddDetail("key2", 123)

	if val, ok := event.GetDetail("key1"); !ok || val != "value1" {
		t.Error("AddDetail should add key-value pair")
	}
}

func TestAuditEvent_WithHash(t *testing.T) {
	event := &AuditEvent{ID: 1}
	event.WithHash("abc123").WithPrevHash("def456")

	if event.Hash != "abc123" {
		t.Error("WithHash should set hash")
	}
	if event.PrevHash != "def456" {
		t.Error("WithPrevHash should set prev_hash")
	}
}

func TestNewAuditEvent(t *testing.T) {
	event := NewAuditEvent(EventAuthSuccess, "user1", "login", "auth")

	if event.EventType != EventAuthSuccess {
		t.Error("NewAuditEvent should set event type")
	}
	if event.Actor != "user1" {
		t.Error("NewAuditEvent should set actor")
	}
	if event.Metadata == nil {
		t.Error("NewAuditEvent should initialize metadata map")
	}
}

func TestNewLoginSuccessEvent(t *testing.T) {
	event := NewLoginSuccessEvent("user1", "127.0.0.1", "password")

	if event.EventType != EventAuthSuccess {
		t.Error("Should be auth success event")
	}
	if ip, _ := event.GetDetailString("ip"); ip != "127.0.0.1" {
		t.Error("Should include IP")
	}
}

func TestNewLoginFailureEvent(t *testing.T) {
	event := NewLoginFailureEvent("user1", "127.0.0.1", "password", "invalid credentials")

	if event.EventType != EventAuthFailure {
		t.Error("Should be auth failure event")
	}
	if reason, _ := event.GetDetailString("failure_reason"); reason != "invalid credentials" {
		t.Error("Should include failure reason")
	}
}

func TestNewResourceAccessEvent(t *testing.T) {
	event := NewResourceAccessEvent("user1", "resource1", "read", true)
	if event.EventType != EventAuthzGranted {
		t.Error("Should be authorization granted")
	}

	deniedEvent := NewResourceAccessEvent("user1", "resource1", "write", false)
	if deniedEvent.EventType != EventAuthzDenied {
		t.Error("Should be authorization denied")
	}
}

func TestNewConfigChangeEvent(t *testing.T) {
	event := NewConfigChangeEvent("admin", "max_connections", "100", "200")

	if event.EventType != EventConfigChange {
		t.Error("Should be config change event")
	}
	if key, _ := event.GetDetailString("config_key"); key != "max_connections" {
		t.Error("Should include config key")
	}
}

func TestNewSecurityAlertEvent(t *testing.T) {
	event := NewSecurityAlertEvent("system", "intrusion", "Detected unauthorized access", SeverityHigh)

	if event.EventType != EventSecurityAlert {
		t.Error("Should be security alert event")
	}
	if event.GetSeverity() != SeverityHigh {
		t.Error("Security alert should have high severity")
	}
}

func TestAuditEvent_Clone(t *testing.T) {
	original := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "user1",
		Metadata:  map[string]interface{}{"key": "value"},
		Hash:      "abc123",
	}

	clone := original.Clone()

	if clone.ID != original.ID {
		t.Error("Clone should have same ID")
	}
	if clone.Metadata["key"] != original.Metadata["key"] {
		t.Error("Clone should have same metadata")
	}

	clone.Metadata["key"] = "modified"
	if original.Metadata["key"] == "modified" {
		t.Error("Clone should be deep copy")
	}
}

func TestCountEventsByType(t *testing.T) {
	events := []*AuditEvent{
		{EventType: EventAuthSuccess},
		{EventType: EventAuthSuccess},
		{EventType: EventAuthFailure},
	}

	counts := CountEventsByType(events)

	if counts[EventAuthSuccess] != 2 {
		t.Error("Should count 2 auth success events")
	}
	if counts[EventAuthFailure] != 1 {
		t.Error("Should count 1 auth failure event")
	}
}

func TestGetUniqueActors(t *testing.T) {
	events := []*AuditEvent{
		{Actor: "user1"},
		{Actor: "user2"},
		{Actor: "user1"},
	}

	actors := GetUniqueActors(events)

	if len(actors) != 2 {
		t.Errorf("Should have 2 unique actors, got %d", len(actors))
	}
}

func TestGetEventsInTimeWindow(t *testing.T) {
	now := time.Now()
	events := []*AuditEvent{
		{Timestamp: now.Add(-2 * time.Hour)},
		{Timestamp: now.Add(-30 * time.Minute)},
		{Timestamp: now.Add(-10 * time.Minute)},
	}

	recent := GetEventsInTimeWindow(events, 1*time.Hour)

	if len(recent) != 2 {
		t.Errorf("Should have 2 events in last hour, got %d", len(recent))
	}
}

func TestAuditEvent_GetSeverity(t *testing.T) {
	event := &AuditEvent{EventType: EventAuthFailure}
	if event.GetSeverity() != SeverityHigh {
		t.Error("Auth failure should be high severity")
	}
}

func TestAuditEvent_IsSecurityCritical(t *testing.T) {
	event := &AuditEvent{EventType: EventAuthFailure}
	if !event.IsSecurityCritical() {
		t.Error("Auth failure should be security critical")
	}
}

func TestAuditEvent_GetCategory(t *testing.T) {
	event := &AuditEvent{EventType: EventAuthSuccess}
	if event.GetCategory() != "authentication" {
		t.Error("Should return authentication category")
	}
}

func TestStorageTypeConstants(t *testing.T) {
	if StorageTypeFile != "file" {
		t.Error("StorageTypeFile should be 'file'")
	}
	if StorageTypeMemory != "memory" {
		t.Error("StorageTypeMemory should be 'memory'")
	}
}

func TestExportFormatConstants(t *testing.T) {
	if ExportFormatJSON != "json" {
		t.Error("ExportFormatJSON should be 'json'")
	}
	if ExportFormatCSV != "csv" {
		t.Error("ExportFormatCSV should be 'csv'")
	}
}

func TestSeverityLevelConstants(t *testing.T) {
	if SeverityLow != "low" {
		t.Error("SeverityLow should be 'low'")
	}
	if SeverityCritical != "critical" {
		t.Error("SeverityCritical should be 'critical'")
	}
}

func TestMemoryStorage_Store(t *testing.T) {
	config := DefaultAuditConfig()
	config.StorageType = StorageTypeMemory

	storage, err := NewAuditStorage(config)
	if err != nil {
		t.Fatalf("Failed to create memory storage: %v", err)
	}
	defer storage.Close()

	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "user1",
		Action:    "login",
		Resource:  "auth",
	}

	if err := storage.Store(event); err != nil {
		t.Errorf("Failed to store event: %v", err)
	}
}

func TestMemoryStorage_Query(t *testing.T) {
	config := DefaultAuditConfig()
	config.StorageType = StorageTypeMemory

	storage, err := NewAuditStorage(config)
	if err != nil {
		t.Fatalf("Failed to create memory storage: %v", err)
	}
	defer storage.Close()

	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "user1",
		Action:    "login",
		Resource:  "auth",
	}
	storage.Store(event)

	ctx := context.Background()
	query := &AuditQuery{
		Limit: 10,
	}

	result, err := storage.Query(ctx, query)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(result.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(result.Events))
	}
}

func TestMemoryStorage_Export(t *testing.T) {
	config := DefaultAuditConfig()
	config.StorageType = StorageTypeMemory

	storage, err := NewAuditStorage(config)
	if err != nil {
		t.Fatalf("Failed to create memory storage: %v", err)
	}
	defer storage.Close()

	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "user1",
	}
	storage.Store(event)

	data, err := storage.Export(time.Now().Add(-time.Hour), time.Now(), ExportFormatJSON)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("Export should return data")
	}
}

func TestMemoryStorage_Rotate(t *testing.T) {
	config := DefaultAuditConfig()
	config.StorageType = StorageTypeMemory

	storage, err := NewAuditStorage(config)
	if err != nil {
		t.Fatalf("Failed to create memory storage: %v", err)
	}
	defer storage.Close()

	if err := storage.Rotate(); err != nil {
		t.Errorf("Rotate should not fail: %v", err)
	}
}

func TestMemoryStorage_Cleanup(t *testing.T) {
	config := DefaultAuditConfig()
	config.StorageType = StorageTypeMemory

	storage, err := NewAuditStorage(config)
	if err != nil {
		t.Fatalf("Failed to create memory storage: %v", err)
	}
	defer storage.Close()

	cutoff := time.Now().Add(-30 * 24 * time.Hour)
	if err := storage.Cleanup(cutoff); err != nil {
		t.Errorf("Cleanup should not fail: %v", err)
	}
}

func TestMemoryStorage_IsClosed(t *testing.T) {
	config := DefaultAuditConfig()
	config.StorageType = StorageTypeMemory

	storage, err := NewAuditStorage(config)
	if err != nil {
		t.Fatalf("Failed to create memory storage: %v", err)
	}

	if storage.IsClosed() {
		t.Error("Storage should not be closed initially")
	}

	storage.Close()

	if !storage.IsClosed() {
		t.Error("Storage should be closed after Close()")
	}
}

func TestAuditor_LogSecurityEvent(t *testing.T) {
	tmpDir := t.TempDir()
	config := &AuditConfig{
		StoragePath: tmpDir,
		StorageType: StorageTypeMemory,
		QueueSize:   100,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		t.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	err = auditor.LogSecurityEvent(EventSecurityAlert, "user1", "Test alert", "high")
	if err != nil {
		t.Errorf("LogSecurityEvent failed: %v", err)
	}
}

func TestAuditor_LogSystemEvent(t *testing.T) {
	tmpDir := t.TempDir()
	config := &AuditConfig{
		StoragePath: tmpDir,
		StorageType: StorageTypeMemory,
		QueueSize:   100,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		t.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	err = auditor.LogSystemEvent(EventSystemStart, "api-server", "start", "success")
	if err != nil {
		t.Errorf("LogSystemEvent failed: %v", err)
	}
}

func TestAuditor_RotateLogs(t *testing.T) {
	tmpDir := t.TempDir()
	config := &AuditConfig{
		StoragePath: tmpDir,
		StorageType: StorageTypeMemory,
		QueueSize:   100,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		t.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	if err := auditor.RotateLogs(); err != nil {
		t.Errorf("RotateLogs should not fail: %v", err)
	}
}

func TestAuditor_Cleanup(t *testing.T) {
	tmpDir := t.TempDir()
	config := &AuditConfig{
		StoragePath:   tmpDir,
		StorageType:   StorageTypeMemory,
		RetentionDays: 30,
		QueueSize:     100,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		t.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	if err := auditor.Cleanup(); err != nil {
		t.Errorf("Cleanup should not fail: %v", err)
	}
}

func TestAuditor_GetEventQueueSize(t *testing.T) {
	tmpDir := t.TempDir()
	config := &AuditConfig{
		StoragePath: tmpDir,
		StorageType: StorageTypeMemory,
		QueueSize:   100,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		t.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	size := auditor.GetEventQueueSize()
	if size < 0 {
		t.Error("Queue size should not be negative")
	}
}

func TestFileStorage_Query(t *testing.T) {
	tmpDir := t.TempDir()
	config := DefaultAuditConfig()
	config.StoragePath = tmpDir
	config.StorageType = StorageTypeFile
	config.EnableCompression = false

	storage, err := NewAuditStorage(config)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}
	defer storage.Close()

	event1 := &AuditEvent{
		ID:        1,
		Timestamp: time.Now().Add(-2 * time.Hour),
		EventType: EventAuthSuccess,
		Actor:     "user1",
		Action:    "login",
		Resource:  "auth",
	}
	event2 := &AuditEvent{
		ID:        2,
		Timestamp: time.Now(),
		EventType: EventAuthFailure,
		Actor:     "user2",
		Action:    "login",
		Resource:  "auth",
	}

	storage.Store(event1)
	storage.Store(event2)
	storage.Flush()

	ctx := context.Background()
	query := &AuditQuery{
		Limit:     10,
		StartTime: time.Now().Add(-1 * time.Hour),
		SortBy:    "timestamp",
		SortOrder: "asc",
	}

	result, err := storage.Query(ctx, query)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(result.Events) != 1 {
		t.Errorf("Expected 1 event in last hour, got %d", len(result.Events))
	}
}

func TestFileStorage_Export(t *testing.T) {
	tmpDir := t.TempDir()
	config := DefaultAuditConfig()
	config.StoragePath = tmpDir
	config.StorageType = StorageTypeFile
	config.EnableCompression = false

	storage, err := NewAuditStorage(config)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}
	defer storage.Close()

	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "user1",
	}
	storage.Store(event)
	storage.Flush()

	data, err := storage.Export(time.Now().Add(-time.Hour), time.Now(), ExportFormatJSON)
	if err != nil {
		t.Fatalf("Export JSON failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("Export should return data")
	}

	csvData, err := storage.Export(time.Now().Add(-time.Hour), time.Now(), ExportFormatCSV)
	if err != nil {
		t.Fatalf("Export CSV failed: %v", err)
	}

	if len(csvData) == 0 {
		t.Error("CSV export should return data")
	}
}

func TestFileStorage_Rotate(t *testing.T) {
	tmpDir := t.TempDir()
	config := DefaultAuditConfig()
	config.StoragePath = tmpDir
	config.StorageType = StorageTypeFile
	config.EnableCompression = false
	config.MaxSize = 10 * 1024 * 1024

	storage, err := NewAuditStorage(config)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}
	defer storage.Close()

	for i := 0; i < 5; i++ {
		event := &AuditEvent{
			ID:        uint64(i),
			Timestamp: time.Now(),
			EventType: EventAuthSuccess,
			Actor:     "user1",
		}
		if err := storage.Store(event); err != nil {
			t.Logf("Store event %d: %v", i, err)
		}
	}

	if err := storage.Flush(); err != nil {
		t.Errorf("Flush failed: %v", err)
	}

	if err := storage.Rotate(); err != nil {
		t.Errorf("Rotate failed: %v", err)
	}
}

func TestFileStorage_Cleanup(t *testing.T) {
	tmpDir := t.TempDir()
	config := DefaultAuditConfig()
	config.StoragePath = tmpDir
	config.StorageType = StorageTypeFile
	config.EnableCompression = false
	config.RetentionDays = 1

	storage, err := NewAuditStorage(config)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}
	defer storage.Close()

	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "user1",
	}
	storage.Store(event)
	storage.Flush()

	cutoff := time.Now().Add(-2 * 24 * time.Hour)
	if err := storage.Cleanup(cutoff); err != nil {
		t.Errorf("Cleanup failed: %v", err)
	}
}

func TestFileStorage_IsClosed(t *testing.T) {
	tmpDir := t.TempDir()
	config := DefaultAuditConfig()
	config.StoragePath = tmpDir
	config.StorageType = StorageTypeFile
	config.EnableCompression = false

	storage, err := NewAuditStorage(config)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}

	if storage.IsClosed() {
		t.Error("Storage should not be closed initially")
	}

	storage.Close()

	if !storage.IsClosed() {
		t.Error("Storage should be closed after Close()")
	}
}

func TestFileStorage_Flush(t *testing.T) {
	tmpDir := t.TempDir()
	config := DefaultAuditConfig()
	config.StoragePath = tmpDir
	config.StorageType = StorageTypeFile
	config.EnableCompression = false

	storage, err := NewAuditStorage(config)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}
	defer storage.Close()

	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "user1",
	}
	storage.Store(event)

	if err := storage.Flush(); err != nil {
		t.Errorf("Flush failed: %v", err)
	}
}

func TestFilterEvents(t *testing.T) {
	events := []*AuditEvent{
		{EventType: EventAuthSuccess, Actor: "user1", Timestamp: time.Now()},
		{EventType: EventAuthFailure, Actor: "user2", Timestamp: time.Now().Add(-time.Hour)},
		{EventType: EventAuthSuccess, Actor: "user1", Timestamp: time.Now().Add(-2 * time.Hour)},
	}

	filtered := FilterEventsByType(events, EventAuthSuccess)
	if len(filtered) != 2 {
		t.Errorf("Should filter 2 success events, got %d", len(filtered))
	}

	filtered = FilterEventsByActor(events, "user1")
	if len(filtered) != 2 {
		t.Errorf("Should filter 2 user1 events, got %d", len(filtered))
	}

	start := time.Now().Add(-90 * time.Minute)
	end := time.Now()
	filtered = FilterEventsByTimeRange(events, start, end)
	if len(filtered) != 2 {
		t.Errorf("Should filter 2 events in range, got %d", len(filtered))
	}
}

func TestValidateEvent(t *testing.T) {
	validEvent := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "user1",
		Action:    "login",
		Resource:  "auth",
	}

	if err := validEvent.Validate(); err != nil {
		t.Errorf("Valid event should pass: %v", err)
	}

	invalidEvents := []*AuditEvent{
		{},
		{ID: 1},
		{ID: 1, Timestamp: time.Now()},
		{ID: 1, Timestamp: time.Now(), EventType: EventAuthSuccess},
		{ID: 1, Timestamp: time.Now(), EventType: EventAuthSuccess, Actor: "user1"},
		{ID: 1, Timestamp: time.Now(), EventType: EventAuthSuccess, Actor: "user1", Action: "login"},
	}

	for i, event := range invalidEvents {
		if err := event.Validate(); err == nil {
			t.Errorf("Invalid event %d should fail validation", i)
		}
	}
}

func TestGetDetailInt(t *testing.T) {
	event := &AuditEvent{
		Metadata: map[string]interface{}{
			"count":  42,
			"float":  3.14,
			"string": "not a number",
		},
	}

	if val, ok := event.GetDetailInt("count"); !ok || val != 42 {
		t.Error("GetDetailInt should work with int")
	}

	if val, ok := event.GetDetailInt("float"); !ok || val != 3 {
		t.Error("GetDetailInt should work with float")
	}

	if _, ok := event.GetDetailInt("string"); ok {
		t.Error("GetDetailInt should fail for string")
	}
}

func TestEventTypeCategory(t *testing.T) {
	tests := []struct {
		eventType EventType
		category  string
	}{
		{EventAuthSuccess, "authentication"},
		{EventAuthFailure, "authentication"},
		{EventAuthzGranted, "authorization"},
		{EventAuthzDenied, "authorization"},
		{EventDataAccess, "data"},
		{EventDataCreate, "data"},
		{EventDataUpdate, "data"},
		{EventDataDelete, "data"},
		{EventConfigChange, "configuration"},
		{EventSystemStart, "system"},
		{EventSystemStop, "system"},
		{EventSystemError, "system"},
		{EventSecurityAlert, "security"},
		{EventSecurityViolation, "security"},
		{EventCompliance, "compliance"},
	}

	for _, test := range tests {
		if test.eventType.Category() != test.category {
			t.Errorf("%s should have category %s, got %s", test.eventType, test.category, test.eventType.Category())
		}
	}
}

func TestGetSeverityForEventAllTypes(t *testing.T) {
	tests := []struct {
		eventType EventType
		severity  SeverityLevel
	}{
		{EventSecurityViolation, SeverityCritical},
		{EventSystemError, SeverityCritical},
		{EventSecurityAlert, SeverityHigh},
		{EventAuthFailure, SeverityHigh},
		{EventAuthzDenied, SeverityHigh},
		{EventDataDelete, SeverityMedium},
		{EventConfigChange, SeverityMedium},
		{EventAuthSuccess, SeverityLow},
		{EventAuthzGranted, SeverityLow},
		{EventDataAccess, SeverityLow},
	}

	for _, test := range tests {
		if GetSeverityForEvent(test.eventType) != test.severity {
			t.Errorf("%s should have severity %s, got %s", test.eventType, test.severity, GetSeverityForEvent(test.eventType))
		}
	}
}
