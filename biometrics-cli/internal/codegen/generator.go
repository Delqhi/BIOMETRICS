package codegen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// CodeTask represents a code generation task
type CodeTask struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Agent       string    `json:"agent"`  // sisyphus, prometheus, oracle
	Status      string    `json:"status"` // pending, running, completed, failed
	Code        string    `json:"code,omitempty"`
	Files       []string  `json:"files,omitempty"`
	Errors      []string  `json:"errors,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

// CodeGenerator handles code generation via OpenCode CLI
type CodeGenerator struct {
	OpenCodePath string
	ProjectDir   string
	Tasks        map[string]*CodeTask
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		OpenCodePath: "opencode", // Assume in PATH
		ProjectDir:   getProjectDir(),
		Tasks:        make(map[string]*CodeTask),
	}
}

// CreateTask creates a new code generation task
func (g *CodeGenerator) CreateTask(title, description, agent string) (*CodeTask, error) {
	task := &CodeTask{
		ID:          fmt.Sprintf("task-%d", time.Now().UnixNano()),
		Title:       title,
		Description: description,
		Agent:       agent,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	g.Tasks[task.ID] = task
	return task, nil
}

// ExecuteTask executes a code generation task via OpenCode CLI
func (g *CodeGenerator) ExecuteTask(taskID string) error {
	task, exists := g.Tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	task.Status = "running"

	// Build opencode command
	cmd := exec.Command(g.OpenCodePath, task.Description, "--agent", task.Agent)
	cmd.Dir = g.ProjectDir
	cmd.Env = os.Environ()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Parse output for generated files
	output := stdout.String()
	if err != nil {
		task.Status = "failed"
		task.Errors = append(task.Errors, stderr.String(), err.Error())
	} else {
		task.Status = "completed"
		task.Code = output
		task.Files = extractFilesFromOutput(output)
		task.CompletedAt = time.Now()
	}

	return err
}

// GetTask returns a task by ID
func (g *CodeGenerator) GetTask(taskID string) (*CodeTask, error) {
	task, exists := g.Tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}
	return task, nil
}

// ListTasks returns all tasks
func (g *CodeGenerator) ListTasks() []*CodeTask {
	tasks := make([]*CodeTask, 0, len(g.Tasks))
	for _, task := range g.Tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// GetActiveTasks returns tasks that are running or pending
func (g *CodeGenerator) GetActiveTasks() []*CodeTask {
	active := make([]*CodeTask, 0)
	for _, task := range g.Tasks {
		if task.Status == "running" || task.Status == "pending" {
			active = append(active, task)
		}
	}
	return active
}

// Helper functions
func getProjectDir() string {
	dir, _ := os.Getwd()
	return dir
}

func extractFilesFromOutput(output string) []string {
	// TODO: Parse opencode output for file paths
	// Look for patterns like "Created: src/file.go"
	return []string{}
}

// TriggerOpenCode sends a request to OpenCode CLI
func (g *CodeGenerator) TriggerOpenCode(prompt string, agent string) (string, error) {
	cmd := exec.Command(g.OpenCodePath, prompt, "--agent", agent)
	cmd.Dir = g.ProjectDir

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("opencode failed: %w", err)
	}

	return string(output), nil
}

// RunCodeGeneration runs code generation with progress tracking
func (g *CodeGenerator) RunCodeGeneration(taskID string, progressChan chan<- string) error {
	task, exists := g.Tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found")
	}

	progressChan <- fmt.Sprintf("Starting task: %s", task.Title)
	progressChan <- fmt.Sprintf("Agent: %s", task.Agent)

	// Execute via OpenCode
	output, err := g.TriggerOpenCode(task.Description, task.Agent)
	if err != nil {
		progressChan <- fmt.Sprintf("ERROR: %v", err)
		return err
	}

	progressChan <- "Code generation completed!"
	progressChan <- fmt.Sprintf("Output length: %d bytes", len(output))

	task.Code = output
	task.Status = "completed"
	task.CompletedAt = time.Now()

	return nil
}

// CodeGenerationRequest represents an API request for code generation
type CodeGenerationRequest struct {
	Prompt  string `json:"prompt"`
	Agent   string `json:"agent"`
	Project string `json:"project,omitempty"`
}

// CodeGenerationResponse represents an API response
type CodeGenerationResponse struct {
	Success bool     `json:"success"`
	TaskID  string   `json:"task_id,omitempty"`
	Code    string   `json:"code,omitempty"`
	Files   []string `json:"files,omitempty"`
	Error   string   `json:"error,omitempty"`
}

// HandleGenerationRequest handles HTTP POST requests for code generation
func (g *CodeGenerator) HandleGenerationRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CodeGenerationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create task
	task, err := g.CreateTask(req.Prompt, req.Prompt, req.Agent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute task (async in production)
	err = g.ExecuteTask(task.ID)
	if err != nil {
		resp := CodeGenerationResponse{
			Success: false,
			Error:   err.Error(),
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := CodeGenerationResponse{
		Success: true,
		TaskID:  task.ID,
		Code:    task.Code,
		Files:   task.Files,
	}
	json.NewEncoder(w).Encode(resp)
}
