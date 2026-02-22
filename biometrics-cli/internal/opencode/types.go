package opencode

type AgentRequest struct {
	ProjectID string
	Model     string
	Prompt    string
	Category  string
}

type AgentResult struct {
	Success bool
	Output  string
	Error   error
}
