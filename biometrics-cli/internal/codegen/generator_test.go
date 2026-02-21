package codegen

import (
	"testing"
	"time"
)

func TestNewCodeGenerator(t *testing.T) {
	g := NewCodeGenerator()

	if g == nil {
		t.Fatal("Expected non-nil generator")
	}

	if g.Tasks == nil {
		t.Error("Expected Tasks map to be initialized")
	}

	if g.OpenCodePath != "opencode" {
		t.Errorf("Expected OpenCodePath to be 'opencode', got %s", g.OpenCodePath)
	}
}

func TestCreateTask(t *testing.T) {
	g := NewCodeGenerator()

	task, err := g.CreateTask("Test Task", "Test Description", "sisyphus")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if task == nil {
		t.Fatal("Expected non-nil task")
	}

	if task.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got %s", task.Title)
	}

	if task.Description != "Test Description" {
		t.Errorf("Expected description 'Test Description', got %s", task.Description)
	}

	if task.Agent != "sisyphus" {
		t.Errorf("Expected agent 'sisyphus', got %s", task.Agent)
	}

	if task.Status != "pending" {
		t.Errorf("Expected status 'pending', got %s", task.Status)
	}

	if task.ID == "" {
		t.Error("Expected non-empty task ID")
	}

	if task.CreatedAt.IsZero() {
		t.Error("Expected non-zero CreatedAt")
	}
}

func TestGetTask(t *testing.T) {
	g := NewCodeGenerator()

	// Create task
	task1, _ := g.CreateTask("Task 1", "Desc 1", "sisyphus")

	// Get task
	task2, err := g.GetTask(task1.ID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if task2.ID != task1.ID {
		t.Errorf("Expected ID %s, got %s", task1.ID, task2.ID)
	}

	// Get non-existent task
	_, err = g.GetTask("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
}

func TestListTasks(t *testing.T) {
	g := NewCodeGenerator()

	// Create multiple tasks
	g.CreateTask("Task 1", "Desc 1", "sisyphus")
	g.CreateTask("Task 2", "Desc 2", "prometheus")
	g.CreateTask("Task 3", "Desc 3", "oracle")

	tasks := g.ListTasks()

	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}
}

func TestGetActiveTasks(t *testing.T) {
	g := NewCodeGenerator()

	// Create tasks with different statuses
	task1, _ := g.CreateTask("Task 1", "Desc 1", "sisyphus")
	task2, _ := g.CreateTask("Task 2", "Desc 2", "prometheus")

	// Set one to completed
	g.Tasks[task2.ID].Status = "completed"

	active := g.GetActiveTasks()

	if len(active) != 1 {
		t.Errorf("Expected 1 active task, got %d", len(active))
	}

	if active[0].ID != task1.ID {
		t.Errorf("Expected active task %s, got %s", task1.ID, active[0].ID)
	}
}

func TestTaskStatusTransitions(t *testing.T) {
	g := NewCodeGenerator()

	task, _ := g.CreateTask("Task", "Desc", "sisyphus")

	// Initial status should be pending
	if task.Status != "pending" {
		t.Errorf("Expected initial status 'pending', got %s", task.Status)
	}

	// Simulate running
	g.Tasks[task.ID].Status = "running"
	if g.Tasks[task.ID].Status != "running" {
		t.Error("Failed to set status to running")
	}

	// Simulate completed
	g.Tasks[task.ID].Status = "completed"
	g.Tasks[task.ID].CompletedAt = time.Now()

	if g.Tasks[task.ID].Status != "completed" {
		t.Error("Failed to set status to completed")
	}

	if g.Tasks[task.ID].CompletedAt.IsZero() {
		t.Error("Expected non-zero CompletedAt")
	}
}

func TestConcurrentTaskCreation(t *testing.T) {
	g := NewCodeGenerator()

	// Create tasks concurrently
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			g.CreateTask("Task "+string(rune(id)), "Desc", "sisyphus")
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	tasks := g.ListTasks()
	if len(tasks) != 10 {
		t.Errorf("Expected 10 tasks, got %d", len(tasks))
	}
}
