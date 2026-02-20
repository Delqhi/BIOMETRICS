package audit

import (
	"encoding/json"
	"fmt"
	"time"
)

type EventType string

const (
	EventAuthSuccess       EventType = "authentication.success"
	EventAuthFailure       EventType = "authentication.failure"
	EventAuthzGranted      EventType = "authorization.granted"
	EventAuthzDenied       EventType = "authorization.denied"
	EventDataAccess        EventType = "data.access"
	EventDataCreate        EventType = "data.create"
	EventDataUpdate        EventType = "data.update"
	EventDataDelete        EventType = "data.delete"
	EventConfigChange      EventType = "config.change"
	EventSystemStart       EventType = "system.start"
	EventSystemStop        EventType = "system.stop"
	EventSystemError       EventType = "system.error"
	EventSecurityAlert     EventType = "security.alert"
	EventSecurityViolation EventType = "security.violation"
	EventCompliance        EventType = "compliance.check"
)

type AuditEvent struct {
	ID        uint64                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	EventType EventType              `json:"event_type"`
	Actor     string                 `json:"actor"`
	Action    string                 `json:"action"`
	Resource  string                 `json:"resource"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Hash      string                 `json:"hash,omitempty"`
	PrevHash  string                 `json:"prev_hash,omitempty"`
}

type AuditQuery struct {
	StartTime  time.Time
	EndTime    time.Time
	EventTypes []EventType
	Actors     []string
	Resources  []string
	Actions    []string
	Limit      int
	Offset     int
	SortBy     string
	SortOrder  string
}

type AuditQueryResult struct {
	Events     []*AuditEvent
	TotalCount int
	HasMore    bool
}

type AuditStats struct {
	TotalEvents     uint64            `json:"total_events"`
	EventsByType    map[string]uint64 `json:"events_by_type"`
	EventsByActor   map[string]uint64 `json:"events_by_actor"`
	LastEventTime   time.Time         `json:"last_event_time"`
	StorageSize     int64             `json:"storage_size_bytes"`
	AvgEventsPerDay float64           `json:"avg_events_per_day"`
}

type ExportFormat string

const (
	ExportFormatJSON ExportFormat = "json"
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatXML  ExportFormat = "xml"
)

type SeverityLevel string

const (
	SeverityLow      SeverityLevel = "low"
	SeverityMedium   SeverityLevel = "medium"
	SeverityHigh     SeverityLevel = "high"
	SeverityCritical SeverityLevel = "critical"
)

func (e EventType) String() string {
	return string(e)
}

func (e EventType) Category() string {
	switch e {
	case EventAuthSuccess, EventAuthFailure:
		return "authentication"
	case EventAuthzGranted, EventAuthzDenied:
		return "authorization"
	case EventDataAccess, EventDataCreate, EventDataUpdate, EventDataDelete:
		return "data"
	case EventConfigChange:
		return "configuration"
	case EventSystemStart, EventSystemStop, EventSystemError:
		return "system"
	case EventSecurityAlert, EventSecurityViolation:
		return "security"
	case EventCompliance:
		return "compliance"
	default:
		return "unknown"
	}
}

func (e EventType) IsSecurityCritical() bool {
	switch e {
	case EventAuthFailure, EventAuthzDenied, EventSecurityAlert, EventSecurityViolation:
		return true
	default:
		return false
	}
}

func GetSeverityForEvent(eventType EventType) SeverityLevel {
	switch eventType {
	case EventSecurityViolation, EventSystemError:
		return SeverityCritical
	case EventSecurityAlert, EventAuthFailure, EventAuthzDenied:
		return SeverityHigh
	case EventDataDelete, EventConfigChange:
		return SeverityMedium
	default:
		return SeverityLow
	}
}

func FilterEventsByType(events []*AuditEvent, eventType EventType) []*AuditEvent {
	filtered := make([]*AuditEvent, 0)
	for _, event := range events {
		if event.EventType == eventType {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

func FilterEventsByActor(events []*AuditEvent, actor string) []*AuditEvent {
	filtered := make([]*AuditEvent, 0)
	for _, event := range events {
		if event.Actor == actor {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

func FilterEventsByTimeRange(events []*AuditEvent, start, end time.Time) []*AuditEvent {
	filtered := make([]*AuditEvent, 0)
	for _, event := range events {
		if !event.Timestamp.Before(start) && !event.Timestamp.After(end) {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

func CountEventsByType(events []*AuditEvent) map[EventType]int {
	counts := make(map[EventType]int)
	for _, event := range events {
		counts[event.EventType]++
	}
	return counts
}

func GetUniqueActors(events []*AuditEvent) []string {
	actorSet := make(map[string]bool)
	for _, event := range events {
		actorSet[event.Actor] = true
	}

	actors := make([]string, 0, len(actorSet))
	for actor := range actorSet {
		actors = append(actors, actor)
	}
	return actors
}

func GetEventsInTimeWindow(events []*AuditEvent, window time.Duration) []*AuditEvent {
	now := time.Now()
	cutoff := now.Add(-window)
	return FilterEventsByTimeRange(events, cutoff, now)
}

func (e *AuditEvent) Validate() error {
	if e.ID == 0 {
		return fmt.Errorf("event ID is required")
	}
	if e.Timestamp.IsZero() {
		return fmt.Errorf("event timestamp is required")
	}
	if e.EventType == "" {
		return fmt.Errorf("event type is required")
	}
	if e.Actor == "" {
		return fmt.Errorf("actor is required")
	}
	if e.Action == "" {
		return fmt.Errorf("action is required")
	}
	if e.Resource == "" {
		return fmt.Errorf("resource is required")
	}
	return nil
}

func (e *AuditEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e *AuditEvent) ToJSONString() (string, error) {
	data, err := e.ToJSON()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (e *AuditEvent) GetDetail(key string) (interface{}, bool) {
	if e.Metadata == nil {
		return nil, false
	}
	val, ok := e.Metadata[key]
	return val, ok
}

func (e *AuditEvent) GetDetailString(key string) (string, bool) {
	val, ok := e.GetDetail(key)
	if !ok {
		return "", false
	}
	str, ok := val.(string)
	return str, ok
}

func (e *AuditEvent) GetDetailInt(key string) (int, bool) {
	val, ok := e.GetDetail(key)
	if !ok {
		return 0, false
	}
	switch v := val.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

func (e *AuditEvent) AddDetail(key string, value interface{}) {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
}

func (e *AuditEvent) GetSeverity() SeverityLevel {
	return GetSeverityForEvent(e.EventType)
}

func (e *AuditEvent) IsSecurityCritical() bool {
	return e.EventType.IsSecurityCritical()
}

func (e *AuditEvent) GetCategory() string {
	return e.EventType.Category()
}

func (e *AuditEvent) WithHash(hash string) *AuditEvent {
	e.Hash = hash
	return e
}

func (e *AuditEvent) WithPrevHash(hash string) *AuditEvent {
	e.PrevHash = hash
	return e
}

func NewAuditEvent(eventType EventType, actor, action, resource string) *AuditEvent {
	return &AuditEvent{
		Timestamp: time.Now().UTC(),
		EventType: eventType,
		Actor:     actor,
		Action:    action,
		Resource:  resource,
		Metadata:  make(map[string]interface{}),
	}
}

func NewLoginSuccessEvent(userID, ip, method string) *AuditEvent {
	event := NewAuditEvent(EventAuthSuccess, userID, "login", "auth-system")
	event.AddDetail("ip", ip)
	event.AddDetail("method", method)
	event.AddDetail("success", true)
	return event
}

func NewLoginFailureEvent(userID, ip, method, reason string) *AuditEvent {
	event := NewAuditEvent(EventAuthFailure, userID, "login", "auth-system")
	event.AddDetail("ip", ip)
	event.AddDetail("method", method)
	event.AddDetail("success", false)
	event.AddDetail("failure_reason", reason)
	return event
}

func NewResourceAccessEvent(userID, resource, action string, allowed bool) *AuditEvent {
	eventType := EventAuthzGranted
	if !allowed {
		eventType = EventAuthzDenied
	}
	event := NewAuditEvent(eventType, userID, action, resource)
	event.AddDetail("allowed", allowed)
	return event
}

func NewConfigChangeEvent(userID, configKey, oldValue, newValue string) *AuditEvent {
	event := NewAuditEvent(EventConfigChange, userID, "config-change", "config-system")
	event.AddDetail("config_key", configKey)
	event.AddDetail("old_value", oldValue)
	event.AddDetail("new_value", newValue)
	return event
}

func NewSecurityAlertEvent(actor, alertType, description string, severity SeverityLevel) *AuditEvent {
	event := NewAuditEvent(EventSecurityAlert, actor, "security-alert", "security-system")
	event.AddDetail("alert_type", alertType)
	event.AddDetail("description", description)
	event.AddDetail("severity", severity)
	return event
}

func (e *AuditEvent) Clone() *AuditEvent {
	clone := &AuditEvent{
		ID:        e.ID,
		Timestamp: e.Timestamp,
		EventType: e.EventType,
		Actor:     e.Actor,
		Action:    e.Action,
		Resource:  e.Resource,
		Hash:      e.Hash,
		PrevHash:  e.PrevHash,
	}

	if e.Metadata != nil {
		clone.Metadata = make(map[string]interface{}, len(e.Metadata))
		for k, v := range e.Metadata {
			clone.Metadata[k] = v
		}
	}

	return clone
}
