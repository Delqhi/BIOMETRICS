package audit

import (
	"os"
	"testing"
	"time"
)

func BenchmarkAuditorLog(b *testing.B) {
	tmpDir := b.TempDir()
	config := &AuditConfig{
		StoragePath:   tmpDir,
		StorageType:   StorageTypeMemory,
		QueueSize:     10000,
		FlushInterval: 1 * time.Minute,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		b.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		auditor.Log(EventAuthSuccess, "benchmark-user", "login", "auth-system", map[string]interface{}{
			"ip":      "192.168.1.1",
			"method":  "password",
			"session": "sess-123",
		})
	}

	b.StopTimer()
}

func BenchmarkAuditorLogParallel(b *testing.B) {
	tmpDir := b.TempDir()
	config := &AuditConfig{
		StoragePath:   tmpDir,
		StorageType:   StorageTypeMemory,
		QueueSize:     100000,
		FlushInterval: 1 * time.Minute,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		b.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			auditor.Log(EventAuthSuccess, "benchmark-user", "login", "auth-system", map[string]interface{}{
				"ip":      "192.168.1.1",
				"method":  "password",
				"session": "sess-123",
			})
			i++
		}
	})

	b.StopTimer()
}

func BenchmarkFileStorageWrite(b *testing.B) {
	tmpDir := b.TempDir()
	config := &AuditConfig{
		StoragePath:       tmpDir,
		StorageType:       StorageTypeFile,
		EnableCompression: false,
		MaxSize:           100 * 1024 * 1024,
	}

	storage, err := NewAuditStorage(config)
	if err != nil {
		b.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "benchmark-user",
		Action:    "login",
		Resource:  "auth-system",
		Metadata: map[string]interface{}{
			"ip":      "192.168.1.1",
			"method":  "password",
			"session": "sess-123",
			"details": "Additional metadata for testing",
			"extra":   "More data to simulate real scenarios",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		event.ID = uint64(i)
		if err := storage.Store(event); err != nil {
			b.Fatalf("Failed to store event: %v", err)
		}
	}

	b.StopTimer()
	if err := storage.Flush(); err != nil {
		b.Fatalf("Failed to flush: %v", err)
	}
}

func BenchmarkMemoryStorageWrite(b *testing.B) {
	config := &AuditConfig{
		StorageType: StorageTypeMemory,
		QueueSize:   100000,
	}

	storage, err := NewAuditStorage(config)
	if err != nil {
		b.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "benchmark-user",
		Action:    "login",
		Resource:  "auth-system",
		Metadata: map[string]interface{}{
			"ip":      "192.168.1.1",
			"method":  "password",
			"session": "sess-123",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		event.ID = uint64(i)
		if err := storage.Store(event); err != nil {
			b.Fatalf("Failed to store event: %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkFileStorageQuery(b *testing.B) {
	tmpDir := b.TempDir()
	config := &AuditConfig{
		StoragePath:       tmpDir,
		StorageType:       StorageTypeFile,
		EnableCompression: false,
	}

	storage, err := NewAuditStorage(config)
	if err != nil {
		b.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	for i := 0; i < 1000; i++ {
		event := &AuditEvent{
			ID:        uint64(i),
			Timestamp: time.Now().Add(-time.Duration(i) * time.Minute),
			EventType: EventAuthSuccess,
			Actor:     "benchmark-user",
			Action:    "login",
			Resource:  "auth-system",
		}
		storage.Store(event)
	}
	storage.Flush()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := &AuditQuery{
			Limit:      100,
			EventTypes: []EventType{EventAuthSuccess},
		}
		if _, err := storage.Query(nil, query); err != nil {
			b.Fatalf("Query failed: %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkMemoryStorageQuery(b *testing.B) {
	config := &AuditConfig{
		StorageType: StorageTypeMemory,
		QueueSize:   10000,
	}

	storage, err := NewAuditStorage(config)
	if err != nil {
		b.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	for i := 0; i < 1000; i++ {
		event := &AuditEvent{
			ID:        uint64(i),
			Timestamp: time.Now().Add(-time.Duration(i) * time.Minute),
			EventType: EventAuthSuccess,
			Actor:     "benchmark-user",
			Action:    "login",
			Resource:  "auth-system",
		}
		storage.Store(event)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := &AuditQuery{
			Limit:      100,
			EventTypes: []EventType{EventAuthSuccess},
		}
		if _, err := storage.Query(nil, query); err != nil {
			b.Fatalf("Query failed: %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkFileStorageExport(b *testing.B) {
	tmpDir := b.TempDir()
	config := &AuditConfig{
		StoragePath:       tmpDir,
		StorageType:       StorageTypeFile,
		EnableCompression: false,
	}

	storage, err := NewAuditStorage(config)
	if err != nil {
		b.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	for i := 0; i < 500; i++ {
		event := &AuditEvent{
			ID:        uint64(i),
			Timestamp: time.Now(),
			EventType: EventAuthSuccess,
			Actor:     "benchmark-user",
			Action:    "login",
			Resource:  "auth-system",
		}
		storage.Store(event)
	}
	storage.Flush()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := storage.Export(time.Now().Add(-time.Hour), time.Now().Add(time.Hour), ExportFormatJSON); err != nil {
			b.Fatalf("Export failed: %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkMemoryStorageExport(b *testing.B) {
	config := &AuditConfig{
		StorageType: StorageTypeMemory,
		QueueSize:   10000,
	}

	storage, err := NewAuditStorage(config)
	if err != nil {
		b.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	for i := 0; i < 500; i++ {
		event := &AuditEvent{
			ID:        uint64(i),
			Timestamp: time.Now(),
			EventType: EventAuthSuccess,
			Actor:     "benchmark-user",
			Action:    "login",
			Resource:  "auth-system",
		}
		storage.Store(event)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := storage.Export(time.Now().Add(-time.Hour), time.Now().Add(time.Hour), ExportFormatJSON); err != nil {
			b.Fatalf("Export failed: %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkAuditEventSerialization(b *testing.B) {
	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "benchmark-user",
		Action:    "login",
		Resource:  "auth-system",
		Metadata: map[string]interface{}{
			"ip":        "192.168.1.100",
			"method":    "oauth2",
			"session":   "sess-abc123",
			"useragent": "Mozilla/5.0",
			"country":   "DE",
			"city":      "Berlin",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := event.ToJSON(); err != nil {
			b.Fatalf("Serialization failed: %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkAuditHash(b *testing.B) {
	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "benchmark-user",
		Action:    "login",
		Resource:  "auth-system",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		event.ID = uint64(i)
		if _, err := CreateAuditHash(event); err != nil {
			b.Fatalf("Hash creation failed: %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkFileStorageGetStats(b *testing.B) {
	tmpDir := b.TempDir()
	config := &AuditConfig{
		StoragePath:       tmpDir,
		StorageType:       StorageTypeFile,
		EnableCompression: false,
	}

	storage, err := NewAuditStorage(config)
	if err != nil {
		b.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	for i := 0; i < 1000; i++ {
		event := &AuditEvent{
			ID:        uint64(i),
			Timestamp: time.Now(),
			EventType: EventAuthSuccess,
			Actor:     "benchmark-user",
			Action:    "login",
			Resource:  "auth-system",
		}
		storage.Store(event)
	}
	storage.Flush()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := storage.GetStats(); err != nil {
			b.Fatalf("GetStats failed: %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkMemoryStorageGetStats(b *testing.B) {
	config := &AuditConfig{
		StorageType: StorageTypeMemory,
		QueueSize:   10000,
	}

	storage, err := NewAuditStorage(config)
	if err != nil {
		b.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	for i := 0; i < 1000; i++ {
		event := &AuditEvent{
			ID:        uint64(i),
			Timestamp: time.Now(),
			EventType: EventAuthSuccess,
			Actor:     "benchmark-user",
			Action:    "login",
			Resource:  "auth-system",
		}
		storage.Store(event)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := storage.GetStats(); err != nil {
			b.Fatalf("GetStats failed: %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkCompressedFileStorageWrite(b *testing.B) {
	tmpDir := b.TempDir()
	config := &AuditConfig{
		StoragePath:       tmpDir,
		StorageType:       StorageTypeFile,
		EnableCompression: true,
		MaxSize:           100 * 1024 * 1024,
	}

	storage, err := NewAuditStorage(config)
	if err != nil {
		b.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "benchmark-user",
		Action:    "login",
		Resource:  "auth-system",
		Metadata: map[string]interface{}{
			"ip":      "192.168.1.1",
			"method":  "password",
			"session": "sess-123",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		event.ID = uint64(i)
		if err := storage.Store(event); err != nil {
			b.Fatalf("Failed to store event: %v", err)
		}
	}

	b.StopTimer()
	if err := storage.Flush(); err != nil {
		b.Fatalf("Failed to flush: %v", err)
	}
}

func BenchmarkLargeMetadataWrite(b *testing.B) {
	config := &AuditConfig{
		StorageType: StorageTypeMemory,
		QueueSize:   10000,
	}

	storage, err := NewAuditStorage(config)
	if err != nil {
		b.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	largeMetadata := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
		"key5": "value5",
		"nested": map[string]interface{}{
			"inner1": "innerValue1",
			"inner2": "innerValue2",
			"list":   []string{"a", "b", "c", "d", "e"},
		},
	}

	event := &AuditEvent{
		ID:        1,
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		Actor:     "benchmark-user",
		Action:    "login",
		Resource:  "auth-system",
		Metadata:  largeMetadata,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		event.ID = uint64(i)
		if err := storage.Store(event); err != nil {
			b.Fatalf("Failed to store event: %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkFileRotation(b *testing.B) {
	tmpDir := b.TempDir()
	config := &AuditConfig{
		StoragePath:       tmpDir,
		StorageType:       StorageTypeFile,
		EnableCompression: false,
		MaxSize:           1024 * 1024,
	}

	storage, err := NewAuditStorage(config)
	if err != nil {
		b.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	for i := 0; i < 100; i++ {
		event := &AuditEvent{
			ID:        uint64(i),
			Timestamp: time.Now(),
			EventType: EventAuthSuccess,
			Actor:     "benchmark-user",
			Action:    "login",
			Resource:  "auth-system",
		}
		storage.Store(event)
	}
	storage.Flush()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := storage.Rotate(); err != nil {
			b.Fatalf("Rotate failed: %v", err)
		}
	}

	b.StopTimer()

	files, _ := os.ReadDir(tmpDir)
	b.Logf("Files after rotation: %d", len(files))
}

func BenchmarkFilterEvents(b *testing.B) {
	events := make([]*AuditEvent, 1000)
	for i := 0; i < 1000; i++ {
		eventType := EventAuthSuccess
		if i%3 == 0 {
			eventType = EventAuthFailure
		}
		events[i] = &AuditEvent{
			ID:        uint64(i),
			Timestamp: time.Now(),
			EventType: eventType,
			Actor:     "benchmark-user",
			Action:    "login",
			Resource:  "auth-system",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filtered := FilterEventsByType(events, EventAuthSuccess)
		_ = filtered
	}

	b.StopTimer()
}

func BenchmarkConcurrentWrites(b *testing.B) {
	tmpDir := b.TempDir()
	config := &AuditConfig{
		StoragePath:   tmpDir,
		StorageType:   StorageTypeMemory,
		QueueSize:     1000000,
		FlushInterval: time.Hour,
	}

	auditor, err := NewAuditor(config)
	if err != nil {
		b.Fatalf("Failed to create auditor: %v", err)
	}
	defer auditor.Stop()

	b.ResetTimer()
	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			auditor.Log(EventAuthSuccess, "user", "action", "resource", nil)
		}
	})

	b.StopTimer()
}
