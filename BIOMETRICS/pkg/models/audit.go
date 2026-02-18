package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID            uuid.UUID              `json:"id" db:"id"`
	TenantID      *uuid.UUID             `json:"tenant_id,omitempty" db:"tenant_id"`
	UserID        *uuid.UUID             `json:"user_id,omitempty" db:"user_id"`
	ActorID       *uuid.UUID             `json:"actor_id,omitempty" db:"actor_id"`
	Action        string                 `json:"action" db:"action"`
	Category      string                 `json:"category" db:"category"`
	ResourceType  string                 `json:"resource_type" db:"resource_type"`
	ResourceID    string                 `json:"resource_id" db:"resource_id"`
	Description   string                 `json:"description" db:"description"`
	OldValues     map[string]interface{} `json:"old_values" db:"old_values"`
	NewValues     map[string]interface{} `json:"new_values" db:"new_values"`
	Changes       map[string]interface{} `json:"changes" db:"changes"`
	IPAddress     string                 `json:"ip_address" db:"ip_address"`
	UserAgent     string                 `json:"user_agent" db:"user_agent"`
	Location      string                 `json:"location" db:"location"`
	SessionID     string                 `json:"session_id" db:"session_id"`
	RequestID     string                 `json:"request_id" db:"request_id"`
	CorrelationID string                 `json:"correlation_id" db:"correlation_id"`
	Level         string                 `json:"level" db:"level"`
	Severity      string                 `json:"severity" db:"severity"`
	Status        string                 `json:"status" db:"status"`
	ErrorMessage  string                 `json:"error_message" db:"error_message"`
	Duration      int                    `json:"duration" db:"duration"`
	Metadata      map[string]interface{} `json:"metadata" db:"metadata"`
	Tags          []string               `json:"tags" db:"tags"`
	Version       int                    `json:"version" db:"version"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
}

type AuditCategory string

const (
	AuditCategoryAuth          AuditCategory = "authentication"
	AuditCategoryAuthorization AuditCategory = "authorization"
	AuditCategoryData          AuditCategory = "data"
	AuditCategoryConfig        AuditCategory = "configuration"
	AuditCategorySystem        AuditCategory = "system"
	AuditCategorySecurity      AuditCategory = "security"
	AuditCategoryCompliance    AuditCategory = "compliance"
)

type AuditLevel string

const (
	AuditLevelInfo     AuditLevel = "info"
	AuditLevelWarning  AuditLevel = "warning"
	AuditLevelError    AuditLevel = "error"
	AuditLevelCritical AuditLevel = "critical"
)

type AuditSeverity string

const (
	AuditSeverityLow      AuditSeverity = "low"
	AuditSeverityMedium   AuditSeverity = "medium"
	AuditSeverityHigh     AuditSeverity = "high"
	AuditSeverityCritical AuditSeverity = "critical"
)

type AuditStatus string

const (
	AuditStatusSuccess AuditStatus = "success"
	AuditStatusFailure AuditStatus = "failure"
	AuditStatusPending AuditStatus = "pending"
)

type AuditFilter struct {
	UserID       string `query:"user_id"`
	ActorID      string `query:"actor_id"`
	Action       string `query:"action"`
	Category     string `query:"category"`
	ResourceType string `query:"resource_type"`
	ResourceID   string `query:"resource_id"`
	TenantID     string `query:"tenant_id"`
	Level        string `query:"level"`
	Severity     string `query:"severity"`
	Status       string `query:"status"`
	IPAddress    string `query:"ip_address"`
	FromDate     string `query:"from_date"`
	ToDate       string `query:"to_date"`
	Search       string `query:"search"`
	Page         int    `query:"page"`
	PageSize     int    `query:"page_size"`
}

type CreateAuditLogInput struct {
	TenantID      string                 `json:"tenant_id"`
	UserID        string                 `json:"user_id"`
	ActorID       string                 `json:"actor_id"`
	Action        string                 `json:"action" binding:"required"`
	Category      string                 `json:"category" binding:"required"`
	ResourceType  string                 `json:"resource_type"`
	ResourceID    string                 `json:"resource_id"`
	Description   string                 `json:"description"`
	OldValues     map[string]interface{} `json:"old_values"`
	NewValues     map[string]interface{} `json:"new_values"`
	IPAddress     string                 `json:"ip_address"`
	UserAgent     string                 `json:"user_agent"`
	Location      string                 `json:"location"`
	SessionID     string                 `json:"session_id"`
	RequestID     string                 `json:"request_id"`
	CorrelationID string                 `json:"correlation_id"`
	Level         string                 `json:"level"`
	Severity      string                 `json:"severity"`
	Duration      int                    `json:"duration"`
	Metadata      map[string]interface{} `json:"metadata"`
	Tags          []string               `json:"tags"`
}

type AuditLogResponse struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     *uuid.UUID `json:"tenant_id,omitempty"`
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	ActorID      *uuid.UUID `json:"actor_id,omitempty"`
	Action       string     `json:"action"`
	Category     string     `json:"category"`
	ResourceType string     `json:"resource_type"`
	ResourceID   string     `json:"resource_id"`
	Description  string     `json:"description"`
	IPAddress    string     `json:"ip_address"`
	UserAgent    string     `json:"user_agent"`
	Location     string     `json:"location"`
	Level        string     `json:"level"`
	Severity     string     `json:"severity"`
	Status       string     `json:"status"`
	ErrorMessage string     `json:"error_message"`
	Duration     int        `json:"duration"`
	Tags         []string   `json:"tags"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (a *AuditLog) ToResponse() *AuditLogResponse {
	return &AuditLogResponse{
		ID:           a.ID,
		TenantID:     a.TenantID,
		UserID:       a.UserID,
		ActorID:      a.ActorID,
		Action:       a.Action,
		Category:     a.Category,
		ResourceType: a.ResourceType,
		ResourceID:   a.ResourceID,
		Description:  a.Description,
		IPAddress:    a.IPAddress,
		UserAgent:    a.UserAgent,
		Location:     a.Location,
		Level:        a.Level,
		Severity:     a.Severity,
		Status:       a.Status,
		ErrorMessage: a.ErrorMessage,
		Duration:     a.Duration,
		Tags:         a.Tags,
		CreatedAt:    a.CreatedAt,
	}
}
