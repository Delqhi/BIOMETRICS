package models

type Boulder struct {
	ActivePlan string   `json:"active_plan"`
	StartedAt  string   `json:"started_at"`
	SessionIDs []string `json:"session_ids"`
	PlanName   string   `json:"plan_name"`
	Agent      string   `json:"agent"`
}
