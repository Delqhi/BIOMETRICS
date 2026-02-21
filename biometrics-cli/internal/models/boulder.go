package models

import "time"

type Boulder struct {
	ActivePlan string   `json:"active_plan"`
	StartedAt  string   `json:"started_at"`
	SessionIDs []string `json:"session_ids"`
	PlanName   string   `json:"plan_name"`
	Agent      string   `json:"agent"`
}

type Session struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Task struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	Priority    int                    `json:"priority"`
	Project     string                 `json:"project"`
	CreatedAt   time.Time              `json:"created_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}
