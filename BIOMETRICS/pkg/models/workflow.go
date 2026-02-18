package models

import (
	"time"

	"github.com/google/uuid"
)

type Workflow struct {
	ID            uuid.UUID              `json:"id" db:"id"`
	TenantID      *uuid.UUID             `json:"tenant_id,omitempty" db:"tenant_id"`
	Name          string                 `json:"name" db:"name"`
	Slug          string                 `json:"slug" db:"slug"`
	Description   string                 `json:"description" db:"description"`
	Definition    map[string]interface{} `json:"definition" db:"definition"`
	TriggerType   string                 `json:"trigger_type" db:"trigger_type"`
	TriggerConfig map[string]interface{} `json:"trigger_config" db:"trigger_config"`
	Status        string                 `json:"status" db:"status"`
	OwnerID       uuid.UUID              `json:"owner_id" db:"owner_id"`
	Version       int                    `json:"version" db:"version"`
	NodesCount    int                    `json:"nodes_count" db:"nodes_count"`
	EdgesCount    int                    `json:"edges_count" db:"edges_count"`
	LastRunAt     *time.Time             `json:"last_run_at" db:"last_run_at"`
	NextRunAt     *time.Time             `json:"next_run_at" db:"next_run_at"`
	RunCount      int                    `json:"run_count" db:"run_count"`
	SuccessCount  int                    `json:"success_count" db:"success_count"`
	FailureCount  int                    `json:"failure_count" db:"failure_count"`
	AvgDuration   int                    `json:"avg_duration" db:"avg_duration"`
	Schedule      string                 `json:"schedule" db:"schedule"`
	IsActive      bool                   `json:"is_active" db:"is_active"`
	IsPublished   bool                   `json:"is_published" db:"is_published"`
	ExecutionMode string                 `json:"execution_mode" db:"execution_mode"`
	Timeout       int                    `json:"timeout" db:"timeout"`
	RetryPolicy   map[string]interface{} `json:"retry_policy" db:"retry_policy"`
	Metadata      map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time             `json:"deleted_at,omitempty" db:"deleted_at"`
}

type WorkflowStatus string

const (
	WorkflowStatusDraft    WorkflowStatus = "draft"
	WorkflowStatusActive   WorkflowStatus = "active"
	WorkflowStatusPaused   WorkflowStatus = "paused"
	WorkflowStatusError    WorkflowStatus = "error"
	WorkflowStatusArchived WorkflowStatus = "archived"
)

type WorkflowTriggerType string

const (
	WorkflowTriggerManual   WorkflowTriggerType = "manual"
	WorkflowTriggerSchedule WorkflowTriggerType = "schedule"
	WorkflowTriggerWebhook  WorkflowTriggerType = "webhook"
	WorkflowTriggerEvent    WorkflowTriggerType = "event"
	WorkflowTriggerCron     WorkflowTriggerType = "cron"
)

type ExecutionMode string

const (
	ExecutionModeSync  ExecutionMode = "sync"
	ExecutionModeAsync ExecutionMode = "async"
	ExecutionModeQueue ExecutionMode = "queue"
)

type CreateWorkflowInput struct {
	TenantID      string                 `json:"tenant_id"`
	Name          string                 `json:"name" binding:"required"`
	Slug          string                 `json:"slug"`
	Description   string                 `json:"description"`
	Definition    map[string]interface{} `json:"definition"`
	TriggerType   string                 `json:"trigger_type"`
	TriggerConfig map[string]interface{} `json:"trigger_config"`
	Schedule      string                 `json:"schedule"`
	ExecutionMode string                 `json:"execution_mode"`
	Timeout       int                    `json:"timeout"`
	RetryPolicy   map[string]interface{} `json:"retry_policy"`
}

type UpdateWorkflowInput struct {
	Name          *string                `json:"name,omitempty"`
	Description   *string                `json:"description,omitempty"`
	Definition    map[string]interface{} `json:"definition,omitempty"`
	TriggerType   *string                `json:"trigger_type,omitempty"`
	TriggerConfig map[string]interface{} `json:"trigger_config,omitempty"`
	Status        *string                `json:"status,omitempty"`
	Schedule      *string                `json:"schedule,omitempty"`
	IsActive      *bool                  `json:"is_active,omitempty"`
	IsPublished   *bool                  `json:"is_published,omitempty"`
	ExecutionMode *string                `json:"execution_mode,omitempty"`
	Timeout       *int                   `json:"timeout,omitempty"`
	RetryPolicy   map[string]interface{} `json:"retry_policy,omitempty"`
}

type WorkflowFilter struct {
	Status      string `query:"status"`
	TriggerType string `query:"trigger_type"`
	TenantID    string `query:"tenant_id"`
	OwnerID     string `query:"owner_id"`
	Search      string `query:"search"`
	Page        int    `query:"page"`
	PageSize    int    `query:"page_size"`
}

type WorkflowResponse struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     *uuid.UUID `json:"tenant_id,omitempty"`
	Name         string     `json:"name"`
	Slug         string     `json:"slug"`
	Description  string     `json:"description"`
	TriggerType  string     `json:"trigger_type"`
	Status       string     `json:"status"`
	OwnerID      uuid.UUID  `json:"owner_id"`
	Version      int        `json:"version"`
	LastRunAt    *time.Time `json:"last_run_at"`
	NextRunAt    *time.Time `json:"next_run_at"`
	RunCount     int        `json:"run_count"`
	SuccessCount int        `json:"success_count"`
	FailureCount int        `json:"failure_count"`
	AvgDuration  int        `json:"avg_duration"`
	IsActive     bool       `json:"is_active"`
	IsPublished  bool       `json:"is_published"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (w *Workflow) ToResponse() *WorkflowResponse {
	return &WorkflowResponse{
		ID:           w.ID,
		TenantID:     w.TenantID,
		Name:         w.Name,
		Slug:         w.Slug,
		Description:  w.Description,
		TriggerType:  w.TriggerType,
		Status:       w.Status,
		OwnerID:      w.OwnerID,
		Version:      w.Version,
		LastRunAt:    w.LastRunAt,
		NextRunAt:    w.NextRunAt,
		RunCount:     w.RunCount,
		SuccessCount: w.SuccessCount,
		FailureCount: w.FailureCount,
		AvgDuration:  w.AvgDuration,
		IsActive:     w.IsActive,
		IsPublished:  w.IsPublished,
		CreatedAt:    w.CreatedAt,
		UpdatedAt:    w.UpdatedAt,
	}
}
