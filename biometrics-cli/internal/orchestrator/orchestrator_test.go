package orchestrator

import (
	"testing"
	"time"
)

func TestOrchestratorInit(t *testing.T) {
	Init()

	if len(DefaultOrchestrator.agents) == 0 {
		t.Error("Expected agents to be initialized")
	}

	if DefaultOrchestrator.agents["sisyphus"] == nil {
		t.Error("Expected sisyphus agent to be registered")
	}

	if DefaultOrchestrator.modelTracker == nil {
		t.Error("Expected model tracker to be initialized")
	}
}

func TestGetAgentForTask(t *testing.T) {
	Init()

	agent := DefaultOrchestrator.GetAgentForTask("coding")
	if agent == nil {
		t.Error("Expected agent for coding task")
	}

	agent = DefaultOrchestrator.GetAgentForTask("research")
	if agent == nil {
		t.Error("Expected agent for research task")
	}
}

func TestAutoCreateTodos(t *testing.T) {
	Init()

	initialLen := len(DefaultOrchestrator.todos)
	DefaultOrchestrator.AutoCreateTodos()

	if len(DefaultOrchestrator.todos) == initialLen {
		t.Error("Expected todos to be created")
	}
}

func TestCompleteTask(t *testing.T) {
	Init()

	DefaultOrchestrator.todos = []TodoTask{}

	task := TodoTask{
		ID:        "test-task-1",
		Title:     "Test Task",
		Status:    "running",
		CreatedAt: time.Now(),
		Agent:     "sisyphus",
	}
	DefaultOrchestrator.todos = append(DefaultOrchestrator.todos, task)

	DefaultOrchestrator.CompleteTask("test-task-1")

	found := false
	for _, t := range DefaultOrchestrator.todos {
		if t.ID == "test-task-1" && t.Status == "completed" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected task to be completed")
	}
}

func TestFailTask(t *testing.T) {
	Init()

	DefaultOrchestrator.todos = []TodoTask{}

	task := TodoTask{
		ID:        "test-task-2",
		Title:     "Test Task",
		Status:    "running",
		CreatedAt: time.Now(),
		Agent:     "sisyphus",
	}
	DefaultOrchestrator.todos = append(DefaultOrchestrator.todos, task)

	DefaultOrchestrator.FailTask("test-task-2", "test error")

	found := false
	for _, t := range DefaultOrchestrator.todos {
		if t.ID == "test-task-2" && t.Status == "failed" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected task to be failed")
	}
}

func TestGetStats(t *testing.T) {
	Init()

	stats := DefaultOrchestrator.GetStats()

	if _, ok := stats["total_tasks"]; !ok {
		t.Error("Expected total_tasks in stats")
	}

	if _, ok := stats["cycle_count"]; !ok {
		t.Error("Expected cycle_count in stats")
	}
}
