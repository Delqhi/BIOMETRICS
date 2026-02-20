package audit

import (
	"compress/gzip"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type StorageType string

const (
	StorageTypeFile   StorageType = "file"
	StorageTypeMemory StorageType = "memory"
)

type AuditStorage interface {
	Store(event *AuditEvent) error
	Query(ctx context.Context, query *AuditQuery) (*AuditQueryResult, error)
	Export(startTime, endTime time.Time, format ExportFormat) ([]byte, error)
	GetStats() (*AuditStats, error)
	Flush() error
	Rotate() error
	Cleanup(cutoff time.Time) error
	Close() error
	IsClosed() bool
}

type FileStorage struct {
	config      *AuditConfig
	filePath    string
	currentFile *os.File
	writer      io.Writer
	mu          sync.RWMutex
	closed      bool
	eventCount  int64
	totalSize   int64
}

type MemoryStorage struct {
	config     *AuditConfig
	events     []*AuditEvent
	mu         sync.RWMutex
	closed     bool
	maxEvents  int
	eventIndex map[uint64]int
}

func NewAuditStorage(config *AuditConfig) (AuditStorage, error) {
	switch config.StorageType {
	case StorageTypeFile, "":
		return NewFileStorage(config)
	case StorageTypeMemory:
		return NewMemoryStorage(config)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", config.StorageType)
	}
}

func NewFileStorage(config *AuditConfig) (*FileStorage, error) {
	if err := os.MkdirAll(config.StoragePath, 0700); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	storage := &FileStorage{
		config:   config,
		filePath: getAuditFilePath(config.StoragePath),
	}

	if err := storage.openFile(); err != nil {
		return nil, err
	}

	return storage, nil
}

func NewMemoryStorage(config *AuditConfig) (*MemoryStorage, error) {
	maxEvents := config.QueueSize * 10
	if maxEvents < 10000 {
		maxEvents = 10000
	}

	return &MemoryStorage{
		config:     config,
		events:     make([]*AuditEvent, 0, maxEvents),
		maxEvents:  maxEvents,
		eventIndex: make(map[uint64]int),
	}, nil
}

func (s *FileStorage) Store(event *AuditEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return fmt.Errorf("storage is closed")
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if _, err := s.writer.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write event: %w", err)
	}

	s.eventCount++
	s.totalSize += int64(len(data)) + 1

	if s.totalSize >= s.config.MaxSize {
		if err := s.rotateLocked(); err != nil {
			return fmt.Errorf("failed to rotate log: %w", err)
		}
	}

	return nil
}

func (s *FileStorage) Query(ctx context.Context, query *AuditQuery) (*AuditQueryResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, fmt.Errorf("storage is closed")
	}

	events, err := s.loadEvents()
	if err != nil {
		return nil, err
	}

	filtered := filterEvents(events, query)
	sortEvents(filtered, query.SortBy, query.SortOrder)

	totalCount := len(filtered)
	hasMore := query.Offset+query.Limit < totalCount

	if query.Limit > 0 && query.Offset < len(filtered) {
		end := query.Offset + query.Limit
		if end > len(filtered) {
			end = len(filtered)
		}
		filtered = filtered[query.Offset:end]
	}

	return &AuditQueryResult{
		Events:     filtered,
		TotalCount: totalCount,
		HasMore:    hasMore,
	}, nil
}

func (s *FileStorage) Export(startTime, endTime time.Time, format ExportFormat) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, fmt.Errorf("storage is closed")
	}

	events, err := s.loadEvents()
	if err != nil {
		return nil, err
	}

	filtered := FilterEventsByTimeRange(events, startTime, endTime)

	switch format {
	case ExportFormatJSON:
		return json.MarshalIndent(filtered, "", "  ")
	case ExportFormatCSV:
		return exportToCSV(filtered)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

func (s *FileStorage) GetStats() (*AuditStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, fmt.Errorf("storage is closed")
	}

	events, err := s.loadEvents()
	if err != nil {
		return nil, err
	}

	stats := &AuditStats{
		TotalEvents:   uint64(len(events)),
		EventsByType:  make(map[string]uint64),
		EventsByActor: make(map[string]uint64),
	}

	for _, event := range events {
		stats.EventsByType[string(event.EventType)]++
		stats.EventsByActor[event.Actor]++

		if event.Timestamp.After(stats.LastEventTime) {
			stats.LastEventTime = event.Timestamp
		}
	}

	stats.StorageSize = s.totalSize

	if len(events) > 0 {
		firstEvent := events[0].Timestamp
		lastEvent := events[len(events)-1].Timestamp
		days := lastEvent.Sub(firstEvent).Hours() / 24
		if days > 0 {
			stats.AvgEventsPerDay = float64(len(events)) / days
		}
	}

	return stats, nil
}

func (s *FileStorage) Flush() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return fmt.Errorf("storage is closed")
	}

	if syncer, ok := s.writer.(interface{ Sync() error }); ok {
		return syncer.Sync()
	}

	return nil
}

func (s *FileStorage) Rotate() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rotateLocked()
}

func (s *FileStorage) rotateLocked() error {
	if s.closed {
		return fmt.Errorf("storage is closed")
	}

	if s.currentFile != nil {
		s.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	rotatedPath := fmt.Sprintf("%s.%s", s.filePath, timestamp)

	if err := os.Rename(s.filePath, rotatedPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to rotate log file: %w", err)
	}

	s.eventCount = 0
	s.totalSize = 0

	return s.openFile()
}

func (s *FileStorage) Cleanup(cutoff time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return fmt.Errorf("storage is closed")
	}

	files, err := filepath.Glob(s.filePath + ".*")
	if err != nil {
		return fmt.Errorf("failed to find log files: %w", err)
	}

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			if err := os.Remove(file); err != nil {
				fmt.Printf("Failed to remove old log file %s: %v\n", file, err)
			}
		}
	}

	return nil
}

func (s *FileStorage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true

	if s.currentFile != nil {
		return s.currentFile.Close()
	}

	return nil
}

func (s *FileStorage) IsClosed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.closed
}

func (s *FileStorage) openFile() error {
	file, err := os.OpenFile(s.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	s.currentFile = file

	if s.config.EnableCompression {
		s.writer = gzip.NewWriter(file)
	} else {
		s.writer = file
	}

	return nil
}

func (s *FileStorage) loadEvents() ([]*AuditEvent, error) {
	files, err := filepath.Glob(s.filePath + "*")
	if err != nil {
		return nil, err
	}

	sort.Strings(files)

	var events []*AuditEvent
	for _, file := range files {
		fileEvents, err := s.loadFile(file)
		if err != nil {
			continue
		}
		events = append(events, fileEvents...)
	}

	return events, nil
}

func (s *FileStorage) loadFile(path string) ([]*AuditEvent, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var reader io.Reader = file
	if strings.HasSuffix(path, ".gz") {
		reader, err = gzip.NewReader(file)
		if err != nil {
			return nil, err
		}
	}

	var events []*AuditEvent
	decoder := json.NewDecoder(reader)
	for {
		var event AuditEvent
		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		events = append(events, &event)
	}

	return events, nil
}

func (s *MemoryStorage) Store(event *AuditEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return fmt.Errorf("storage is closed")
	}

	if len(s.events) >= s.maxEvents {
		s.events = s.events[1:]
	}

	s.events = append(s.events, event)
	s.eventIndex[event.ID] = len(s.events) - 1

	return nil
}

func (s *MemoryStorage) Query(ctx context.Context, query *AuditQuery) (*AuditQueryResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, fmt.Errorf("storage is closed")
	}

	filtered := filterEvents(s.events, query)
	sortEvents(filtered, query.SortBy, query.SortOrder)

	totalCount := len(filtered)
	hasMore := query.Offset+query.Limit < totalCount

	if query.Limit > 0 && query.Offset < len(filtered) {
		end := query.Offset + query.Limit
		if end > len(filtered) {
			end = len(filtered)
		}
		filtered = filtered[query.Offset:end]
	}

	return &AuditQueryResult{
		Events:     filtered,
		TotalCount: totalCount,
		HasMore:    hasMore,
	}, nil
}

func (s *MemoryStorage) Export(startTime, endTime time.Time, format ExportFormat) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, fmt.Errorf("storage is closed")
	}

	filtered := FilterEventsByTimeRange(s.events, startTime, endTime)

	switch format {
	case ExportFormatJSON:
		return json.MarshalIndent(filtered, "", "  ")
	case ExportFormatCSV:
		return exportToCSV(filtered)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

func (s *MemoryStorage) GetStats() (*AuditStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, fmt.Errorf("storage is closed")
	}

	stats := &AuditStats{
		TotalEvents:   uint64(len(s.events)),
		EventsByType:  make(map[string]uint64),
		EventsByActor: make(map[string]uint64),
	}

	for _, event := range s.events {
		stats.EventsByType[string(event.EventType)]++
		stats.EventsByActor[event.Actor]++

		if event.Timestamp.After(stats.LastEventTime) {
			stats.LastEventTime = event.Timestamp
		}
	}

	return stats, nil
}

func (s *MemoryStorage) Flush() error {
	return nil
}

func (s *MemoryStorage) Rotate() error {
	return nil
}

func (s *MemoryStorage) Cleanup(cutoff time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return fmt.Errorf("storage is closed")
	}

	filtered := make([]*AuditEvent, 0)
	for _, event := range s.events {
		if !event.Timestamp.Before(cutoff) {
			filtered = append(filtered, event)
		}
	}

	s.events = filtered
	return nil
}

func (s *MemoryStorage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}

func (s *MemoryStorage) IsClosed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.closed
}

func getAuditFilePath(storagePath string) string {
	timestamp := time.Now().Format("20060102")
	return filepath.Join(storagePath, fmt.Sprintf("audit_%s.log", timestamp))
}

func filterEvents(events []*AuditEvent, query *AuditQuery) []*AuditEvent {
	filtered := make([]*AuditEvent, 0)

	for _, event := range events {
		if !query.StartTime.IsZero() && event.Timestamp.Before(query.StartTime) {
			continue
		}
		if !query.EndTime.IsZero() && event.Timestamp.After(query.EndTime) {
			continue
		}

		if len(query.EventTypes) > 0 {
			found := false
			for _, et := range query.EventTypes {
				if event.EventType == et {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if len(query.Actors) > 0 {
			found := false
			for _, actor := range query.Actors {
				if event.Actor == actor {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if len(query.Resources) > 0 {
			found := false
			for _, resource := range query.Resources {
				if event.Resource == resource {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		filtered = append(filtered, event)
	}

	return filtered
}

func sortEvents(events []*AuditEvent, sortBy, sortOrder string) {
	if sortBy == "" {
		sortBy = "timestamp"
	}

	sort.Slice(events, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "timestamp":
			less = events[i].Timestamp.Before(events[j].Timestamp)
		case "event_type":
			less = string(events[i].EventType) < string(events[j].EventType)
		case "actor":
			less = events[i].Actor < events[j].Actor
		default:
			less = events[i].Timestamp.Before(events[j].Timestamp)
		}

		if sortOrder == "desc" {
			return !less
		}
		return less
	})
}

func exportToCSV(events []*AuditEvent) ([]byte, error) {
	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	headers := []string{"id", "timestamp", "event_type", "actor", "action", "resource", "metadata"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, event := range events {
		metadataJSON, _ := json.Marshal(event.Metadata)
		record := []string{
			fmt.Sprintf("%d", event.ID),
			event.Timestamp.Format(time.RFC3339),
			string(event.EventType),
			event.Actor,
			event.Action,
			event.Resource,
			string(metadataJSON),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return []byte(buf.String()), nil
}
