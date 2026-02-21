# Audit Package - Comprehensive Documentation

**Package:** github.com/delqhi/biometrics/pkg/audit  
**Version:** 1.0.0  
**Date:** February 2026  
**Status:** Sprint 1 Feature

---

## Table of Contents

1. [Package Overview](#package-overview)
2. [Installation](#installation)
3. [Quick Start](#quick-start)
4. [API Reference](#api-reference)
5. [Event Types](#event-types)
6. [Storage Backends](#storage-backends)
7. [Query Examples](#query-examples)
8. [Configuration Options](#configuration-options)

---

## Package Overview

The audit package provides comprehensive security audit logging for the Biometrics CLI. It enables organizations to meet compliance requirements (SOC2, HIPAA, GDPR) by recording all security-relevant events in a tamper-proof manner.

### Key Features

- **Event Logging:** Capture authentication, authorization, data access, and system events
- **Multiple Storage Backends:** File-based, in-memory, or database storage
- **Query API:** Powerful querying with filters, sorting, and pagination
- **Export Capabilities:** Export to JSON, CSV, or XML formats
- **Log Rotation:** Automatic log rotation based on size and time
- **Integrity Verification:** SHA-256 hashing for tamper detection
- **Compression:** Optional gzip compression for storage efficiency
- **Encryption:** Optional encryption for sensitive audit data

### Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         AUDIT PACKAGE ARCHITECTURE                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────┐      ┌─────────────┐      ┌─────────────┐        │
│  │   Client    │      │   Client    │      │   Client    │        │
│  │  Code A     │      │  Code B     │      │  Code C     │        │
│  └──────┬──────┘      └──────┬──────┘      └──────┬──────┘        │
│         │                     │                     │                │
│         └─────────────────────┼─────────────────────┘                │
│                               │                                       │
│                               ▼                                       │
│                    ┌─────────────────────┐                           │
│                    │     Auditor         │                           │
│                    │  (Main Interface)   │                           │
│                    └──────────┬──────────┘                           │
│                               │                                       │
│         ┌─────────────────────┼─────────────────────┐                │
│         │                     │                     │                │
│         ▼                     ▼                     ▼                │
│  ┌─────────────┐      ┌─────────────┐      ┌─────────────┐        │
│  │ Event Queue  │      │  Hash Chain │      │   Stats     │        │
│  │  (Buffered)  │      │  (Integrity)│      │  Collector  │        │
│  └──────┬──────┘      └─────────────┘      └─────────────┘        │
│         │                                                         │
│         ▼                                                         │
│  ┌─────────────────────────────────────────┐                      │
│  │           Storage Backend               │                      │
│  │  ┌──────────┐  ┌──────────┐  ┌──────┐ │                      │
│  │  │  File    │  │  Memory  │  │Custom │ │                      │
│  │  │ Storage  │  │ Storage  │  │      │ │                      │
│  │  └──────────┘  └──────────┘  └──────┘ │                      │
│  └─────────────────────────────────────────┘                      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Installation

### Prerequisites

- Go 1.21 or later
- For file storage: Write access to the configured directory
- For Redis storage: Redis 6.0+ (optional)

### Installation Command

```bash
go get github.com/delqhi/biometrics/pkg/audit
```

### Dependencies

The audit package depends on the following:

```go
import (
    "context"           // Standard library
    "encoding/json"     // Standard library
    "fmt"              // Standard library
    "sync"             // Standard library
    "time"             // Standard library
    "crypto/sha256"   // Standard library
    "encoding/hex"     // Standard library
)
```

No external dependencies are required for the core functionality. Optional compression requires:

```bash
go get compress/gzips
```

---

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/delqhi/biometrics/pkg/audit"
)

func main() {
    // Create auditor with default configuration
    auditor, err := audit.NewAuditor(nil)
    if err != nil {
        log.Fatalf("Failed to create auditor: %v", err)
    }
    defer auditor.Stop()
    
    // Log an authentication event
    err = auditor.LogAuthentication(
        "user-123",  // User ID
        true,        // Success
        "oauth2",    // Method
        "192.168.1.100", // IP Address
    )
    if err != nil {
        log.Printf("Failed to log event: %v", err)
    }
    
    // Log authorization
    err = auditor.LogAuthorization(
        "user-123",
        "/api/admin",
        "DELETE",
        true,
    )
    if err != nil {
        log.Printf("Failed to log event: %v", err)
    }
    
    fmt.Println("Audit logging configured successfully")
}
```

### With Custom Configuration

```go
func customConfig() *audit.Auditor {
    config := &audit.AuditConfig{
        StoragePath:        "/var/log/biometrics/audit",
        StorageType:        audit.StorageTypeFile,
        MaxSize:            100 * 1024 * 1024, // 100MB
        RetentionDays:      90,
        EnableCompression:  true,
        EnableEncryption:   false,
        FlushInterval:      5 * 1e9, // 5 seconds (in nanoseconds)
        QueueSize:          1000,
    }
    
    auditor, err := audit.NewAuditor(config)
    if err != nil {
        panic(err)
    }
    
    return auditor
}
```

### Logging Different Event Types

```go
func logVariousEvents(auditor *audit.Auditor) {
    // Authentication events
    auditor.LogAuthentication("user-123", true, "password", "10.0.0.1")
    auditor.LogAuthentication("user-456", false, "oauth2", "10.0.0.2")
    
    // Authorization events
    auditor.LogAuthorization("user-123", "/api/users", "GET", true)
    auditor.LogAuthorization("user-456", "/api/admin", "DELETE", false)
    
    // Data access events
    auditor.LogDataAccess("user-123", "medical-records", "READ", "rec-001")
    auditor.LogDataAccess("user-123", "medical-records", "WRITE", "rec-002")
    auditor.LogDataAccess("user-456", "medical-records", "DELETE", "rec-003")
    
    // Security events
    auditor.LogSecurityEvent(
        audit.EventSecurityAlert,
        "system",
        "Multiple failed login attempts from IP 10.0.0.55",
        "high",
    )
    
    // System events
    auditor.LogSystemEvent(
        audit.EventSystemStart,
        "api-server",
        "started",
        "ok",
    )
}
```

---

## API Reference

### Types

#### AuditConfig

Configuration for the audit system:

```go
type AuditConfig struct {
    StoragePath        string        // Path for storage (file storage)
    StorageType        StorageType   // file, memory
    MaxSize            int64         // Max file size before rotation
    RetentionDays      int           // Days to retain logs
    EnableCompression  bool          // Enable gzip compression
    EnableEncryption   bool          // Enable encryption
    FlushInterval      time.Duration // How often to flush to storage
    QueueSize          int           // Event queue size
}
```

#### AuditEvent

Represents a single audit event:

```go
type AuditEvent struct {
    ID        uint64                 // Unique event ID
    Timestamp time.Time              // Event timestamp
    EventType EventType              // Type of event
    Actor     string                 // Who performed the action
    Action    string                 // What was done
    Resource  string                 // Resource affected
    Metadata  map[string]interface{} // Additional data
    Hash      string                 // Event hash (integrity)
    PrevHash  string                 // Previous event hash (chain)
}
```

#### AuditQuery

Query parameters for searching events:

```go
type AuditQuery struct {
    StartTime  time.Time   // Filter start time
    EndTime    time.Time   // Filter end time
    EventTypes []EventType // Filter by event types
    Actors     []string    // Filter by actors
    Resources  []string    // Filter by resources
    Actions    []string    // Filter by actions
    Limit      int         // Max results
    Offset     int         // Pagination offset
    SortBy     string      // Sort field (timestamp, event_type, actor)
    SortOrder  string      // Sort order (asc, desc)
}
```

#### AuditStats

Statistics about the audit log:

```go
type AuditStats struct {
    TotalEvents     uint64            // Total events logged
    EventsByType    map[string]uint64 // Count by event type
    EventsByActor   map[string]uint64 // Count by actor
    LastEventTime   time.Time         // Most recent event
    StorageSize     int64             // Total storage used
    AvgEventsPerDay float64           // Average events per day
}
```

### Functions

#### NewAuditor

Creates a new auditor instance:

```go
func NewAuditor(config *AuditConfig) (*Auditor, error)
```

#### DefaultAuditConfig

Returns default configuration:

```go
func DefaultAuditConfig() *AuditConfig
```

### Methods

#### Auditor Methods

##### Log

Log a generic event:

```go
func (a *Auditor) Log(
    eventType EventType,
    actor string,
    action string,
    resource string,
    metadata map[string]interface{},
) error
```

##### LogAuthentication

Log an authentication attempt:

```go
func (a *Auditor) LogAuthentication(
    userID string,
    success bool,
    method string,
    ip string,
) error
```

##### LogAuthorization

Log an authorization decision:

```go
func (a *Auditor) LogAuthorization(
    userID string,
    resource string,
    action string,
    allowed bool,
) error
```

##### LogDataAccess

Log data access:

```go
func (a *Auditor) LogDataAccess(
    userID string,
    dataType string,
    operation string,
    recordID string,
) error
```

##### LogSecurityEvent

Log a security event:

```go
func (a *Auditor) LogSecurityEvent(
    eventType EventType,
    actor string,
    details string,
    severity string,
) error
```

##### LogSystemEvent

Log a system event:

```go
func (a *Auditor) LogSystemEvent(
    eventType EventType,
    component string,
    action string,
    status string,
) error
```

##### Query

Query audit events:

```go
func (a *Auditor) Query(ctx context.Context, query *AuditQuery) (*AuditQueryResult, error)
```

##### Export

Export events to a format:

```go
func (a *Auditor) Export(
    startTime time.Time,
    endTime time.Time,
    format ExportFormat,
) ([]byte, error)
```

##### GetStats

Get audit statistics:

```go
func (a *Auditor) GetStats() (*AuditStats, error)
```

##### RotateLogs

Trigger log rotation:

```go
func (a *Auditor) RotateLogs() error
```

##### Cleanup

Remove old logs:

```go
func (a *Auditor) Cleanup() error
```

##### Stop

Stop the auditor:

```go
func (a *Auditor) Stop()
```

---

## Event Types

The package defines the following event types:

### Authentication Events

| Event Type | Description | Security Critical |
|------------|-------------|------------------|
| `authentication.success` | Successful authentication | No |
| `authentication.failure` | Failed authentication attempt | Yes |

### Authorization Events

| Event Type | Description | Security Critical |
|------------|-------------|------------------|
| `authorization.granted` | Access granted | No |
| `authorization.denied` | Access denied | Yes |

### Data Events

| Event Type | Description | Security Critical |
|------------|-------------|------------------|
| `data.access` | Data was accessed | Yes |
| `data.create` | Data was created | No |
| `data.update` | Data was modified | No |
| `data.delete` | Data was deleted | Yes |

### Configuration Events

| Event Type | Description | Security Critical |
|------------|-------------|------------------|
| `config.change` | Configuration was changed | Yes |

### System Events

| Event Type | Description | Security Critical |
|------------|-------------|------------------|
| `system.start` | System started | No |
| `system.stop` | System stopped | No |
| `system.error` | System error occurred | Yes |

### Security Events

| Event Type | Description | Security Critical |
|------------|-------------|------------------|
| `security.alert` | Security alert triggered | Yes |
| `security.violation` | Security policy violated | Yes |

### Helper Functions

```go
// Get category for an event type
category := audit.EventAuthSuccess.Category()  // Returns: "authentication"

// Check if security critical
isCritical := audit.EventAuthFailure.IsSecurityCritical()  // Returns: true

// Get severity for event type
severity := audit.GetSeverityForEvent(audit.EventAuthFailure)  // Returns: SeverityHigh
```

---

## Storage Backends

### File Storage

Default backend that writes events to JSON Lines files:

```go
config := &audit.AuditConfig{
    StoragePath:        "/var/log/biometrics/audit",
    StorageType:        audit.StorageTypeFile,
    MaxSize:            100 * 1024 * 1024,  // 100MB per file
    EnableCompression:  true,               // Gzip compression
}
```

Features:
- Automatic file rotation by size
- Optional gzip compression
- Daily file naming (`audit_20260220.log`)
- Automatic `.gz` suffix for compressed files

### Memory Storage

In-memory storage for testing or short-lived applications:

```go
config := &audit.AuditConfig{
    StorageType: audit.StorageTypeMemory,
    QueueSize:   10000,  // Max events in memory
}
```

Features:
- Fastest performance
- No persistence (events lost on restart)
- Automatic oldest event eviction when full

### Custom Storage

Implement the AuditStorage interface for custom backends:

```go
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
```

---

## Query Examples

### Basic Queries

```go
// Get all events from the last 24 hours
query := &audit.AuditQuery{
    StartTime: time.Now().Add(-24 * time.Hour),
    Limit:     100,
    SortBy:    "timestamp",
    SortOrder: "desc",
}
results, err := auditor.Query(context.Background(), query)
```

### Filter by Event Type

```go
// Find all authentication failures
query := &audit.AuditQuery{
    StartTime:  time.Now().Add(-7 * 24 * time.Hour),
    EventTypes: []audit.EventType{audit.EventAuthFailure},
    Limit:      1000,
}
```

### Filter by Actor

```go
// Find all events for a specific user
query := &audit.AuditQuery{
    Actors:    []string{"user-123"},
    Limit:     100,
    SortBy:    "timestamp",
    SortOrder: "desc",
}
```

### Complex Query

```go
// Find failed auth attempts in the last hour from a specific IP
query := &audit.AuditQuery{
    StartTime:  time.Now().Add(-1 * time.Hour),
    EndTime:    time.Now(),
    EventTypes: []audit.EventType{audit.EventAuthFailure},
    Limit:      100,
    SortBy:     "timestamp",
    SortOrder:  "desc",
}

results, err := auditor.Query(context.Background(), query)
for _, event := range results.Events {
    if ip, ok := event.GetDetailString("ip"); ok && ip == "10.0.0.55" {
        fmt.Printf("Suspicious: %s at %s\n", event.Actor, event.Timestamp)
    }
}
```

### Pagination

```go
// Paginate through results
page := 0
pageSize := 100

for {
    query := &audit.AuditQuery{
        StartTime: time.Now().Add(-30 * 24 * time.Hour),
        Limit:     pageSize,
        Offset:    page * pageSize,
    }
    
    results, err := auditor.Query(context.Background(), query)
    if err != nil {
        break
    }
    
    // Process results
    processEvents(results.Events)
    
    if !results.HasMore {
        break
    }
    page++
}
```

### Query Result Processing

```go
func processResults(results *audit.AuditQueryResult) {
    fmt.Printf("Total events: %d\n", results.TotalCount)
    fmt.Printf("Returned: %d\n", len(results.Events))
    fmt.Printf("Has more: %v\n", results.HasMore)
    
    for _, event := range results.Events {
        fmt.Printf("[%s] %s: %s %s %s\n",
            event.Timestamp.Format("2006-01-02 15:04:05"),
            event.EventType,
            event.Actor,
            event.Action,
            event.Resource,
        )
    }
}
```

---

## Configuration Options

### Complete Configuration Example

```go
config := &audit.AuditConfig{
    // Storage configuration
    StoragePath:        "/var/log/biometrics/audit",
    StorageType:        audit.StorageTypeFile,
    
    // Rotation settings
    MaxSize:            100 * 1024 * 1024,  // 100MB
    RetentionDays:      90,
    
    // Compression and encryption
    EnableCompression:  true,
    EnableEncryption:   false,  // Set true for sensitive data
    
    // Performance tuning
    FlushInterval:      5 * time.Second,
    QueueSize:          1000,
}
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `BIOMETRICS_AUDIT_PATH` | Storage path | `/tmp/biometrics/audit` |
| `BIOMETRICS_AUDIT_TYPE` | Storage type | `file` |
| `BIOMETRICS_AUDIT_MAX_SIZE` | Max file size | `104857600` (100MB) |
| `BIOMETRICS_AUDIT_RETENTION` | Retention days | `90` |
| `BIOMETRICS_AUDIT_COMPRESS` | Enable compression | `true` |

---

## Error Handling

```go
// Proper error handling
func handleAuditErrors() {
    auditor, err := audit.NewAuditor(nil)
    if err != nil {
        // Handle configuration errors
        fmt.Printf("Configuration error: %v\n", err)
        return
    }
    
    // Log with error handling
    if err := auditor.LogAuthentication("user", true, "oauth", "ip"); err != nil {
        switch {
        case err.Error() == "audit queue full, event dropped":
            // Queue is full - might need to increase QueueSize
            fmt.Printf("Warning: Audit queue full\n")
        default:
            fmt.Printf("Audit error: %v\n", err)
        }
    }
}
```

---

## Best Practices

### 1. Always Log Authentication Events

```go
// Good: Always log authentication
func loginHandler(w http.ResponseWriter, r *http.Request) {
    // ... authentication logic ...
    
    if success {
        auditor.LogAuthentication(userID, true, method, getClientIP(r))
    } else {
        auditor.LogAuthentication(userID, false, method, getClientIP(r))
    }
}
```

### 2. Use Appropriate Retention Periods

```go
// Configuration based on compliance requirements
config := &audit.AuditConfig{
    // SOC2: 90 days minimum
    // HIPAA: 6 years
    // GDPR: 2 years
    RetentionDays: 2555,  // ~7 years for HIPAA
}
```

### 3. Monitor Storage Usage

```go
// Regular monitoring
func monitorAuditStats(auditor *audit.Auditor) {
    stats, err := auditor.GetStats()
    if err != nil {
        return
    }
    
    // Alert if storage is growing too fast
    if stats.AvgEventsPerDay > 100000 {
        fmt.Println("Warning: High audit event volume")
    }
}
```

---

## Related Documentation

- [Security Configuration](../SECURITY.md)
- [Performance Configuration](../PERFORMANCE.md)
- [API Reference](../docs/api/)

---

*Document Version: 1.0.0*  
*Last Updated: February 2026*  
*Compliant with Enterprise Practices Feb 2026*
