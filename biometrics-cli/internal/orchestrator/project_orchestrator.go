package orchestrator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"biometrics-cli/internal/models"
)

type ProjectBoulder struct {
	mu             sync.RWMutex
	Project        string                 `json:"project"`
	ActivePlan     string                 `json:"active_plan"`
	PlanName       string                 `json:"plan_name"`
	CurrentSession *models.Session        `json:"current_session,omitempty"`
	Tasks          []*models.Task         `json:"tasks"`
	CompletedTasks []*models.Task         `json:"completed_tasks"`
	LastUpdated    time.Time              `json:"last_updated"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type ProjectOrchestrator struct {
	mu             sync.RWMutex
	projects       map[string]*ProjectBoulder
	basePath       string
	currentProject string
}

func NewProjectOrchestrator(basePath string) *ProjectOrchestrator {
	return &ProjectOrchestrator{
		projects: make(map[string]*ProjectBoulder),
		basePath: basePath,
	}
}

func (po *ProjectOrchestrator) LoadProject(projectName string) (*ProjectBoulder, error) {
	po.mu.Lock()
	defer po.mu.Unlock()

	boulderPath := filepath.Join(po.basePath, "plans", projectName, "boulder.json")

	data, err := os.ReadFile(boulderPath)
	if err != nil {
		if os.IsNotExist(err) {
			boulder := &ProjectBoulder{
				Project:        projectName,
				ActivePlan:     "",
				PlanName:       projectName,
				Tasks:          make([]*models.Task, 0),
				CompletedTasks: make([]*models.Task, 0),
				LastUpdated:    time.Now(),
				Metadata:       make(map[string]interface{}),
			}
			po.projects[projectName] = boulder
			if err := po.SaveProject(projectName); err != nil {
				return nil, err
			}
			return boulder, nil
		}
		return nil, fmt.Errorf("failed to read boulder: %w", err)
	}

	var boulder ProjectBoulder
	if err := json.Unmarshal(data, &boulder); err != nil {
		return nil, fmt.Errorf("failed to parse boulder: %w", err)
	}

	boulder.LastUpdated = time.Now()
	po.projects[projectName] = &boulder

	return &boulder, nil
}

func (po *ProjectOrchestrator) SaveProject(projectName string) error {
	po.mu.RLock()
	boulder, exists := po.projects[projectName]
	po.mu.RUnlock()

	if !exists {
		return fmt.Errorf("project %s not loaded", projectName)
	}

	boulder.mu.RLock()
	defer boulder.mu.RUnlock()

	boulderPath := filepath.Join(po.basePath, "plans", projectName, "boulder.json")
	boulder.LastUpdated = time.Now()

	data, err := json.MarshalIndent(boulder, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal boulder: %w", err)
	}

	if err := os.WriteFile(boulderPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write boulder: %w", err)
	}

	return nil
}

func (po *ProjectOrchestrator) SetCurrentProject(projectName string) error {
	po.mu.Lock()
	defer po.mu.Unlock()

	if _, exists := po.projects[projectName]; !exists {
		if _, err := po.LoadProject(projectName); err != nil {
			return err
		}
	}

	po.currentProject = projectName
	return nil
}

func (po *ProjectOrchestrator) GetCurrentProject() string {
	po.mu.RLock()
	defer po.mu.RUnlock()
	return po.currentProject
}

func (po *ProjectOrchestrator) GetCurrentBoulder() (*ProjectBoulder, error) {
	po.mu.RLock()
	projectName := po.currentProject
	po.mu.RUnlock()

	if projectName == "" {
		return nil, fmt.Errorf("no project selected")
	}

	po.mu.RLock()
	boulder, exists := po.projects[projectName]
	po.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("project %s not loaded", projectName)
	}

	return boulder, nil
}

func (po *ProjectOrchestrator) AddTask(projectName string, task *models.Task) error {
	po.mu.RLock()
	boulder, exists := po.projects[projectName]
	po.mu.RUnlock()

	if !exists {
		var err error
		boulder, err = po.LoadProject(projectName)
		if err != nil {
			return err
		}
	}

	boulder.mu.Lock()
	defer boulder.mu.Unlock()

	task.ID = fmt.Sprintf("%s-%d", projectName, len(boulder.Tasks)+1)
	now := time.Now()
	task.CreatedAt = now
	boulder.Tasks = append(boulder.Tasks, task)

	return po.SaveProject(projectName)
}

func (po *ProjectOrchestrator) CompleteTask(projectName, taskID string) error {
	po.mu.RLock()
	boulder, exists := po.projects[projectName]
	po.mu.RUnlock()

	if !exists {
		return fmt.Errorf("project %s not found", projectName)
	}

	boulder.mu.Lock()
	defer boulder.mu.Unlock()

	for i, task := range boulder.Tasks {
		if task.ID == taskID {
			task.Status = "completed"
			now := time.Now()
			task.CompletedAt = &now

			boulder.CompletedTasks = append(boulder.CompletedTasks, task)
			boulder.Tasks = append(boulder.Tasks[:i], boulder.Tasks[i+1:]...)

			return po.SaveProject(projectName)
		}
	}

	return fmt.Errorf("task %s not found", taskID)
}

func (po *ProjectOrchestrator) GetNextTask(projectName string) (*models.Task, error) {
	po.mu.RLock()
	boulder, exists := po.projects[projectName]
	po.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("project %s not found", projectName)
	}

	boulder.mu.RLock()
	defer boulder.mu.RUnlock()

	if len(boulder.Tasks) == 0 {
		return nil, fmt.Errorf("no tasks available")
	}

	return boulder.Tasks[0], nil
}

func (po *ProjectOrchestrator) ListProjects() []string {
	po.mu.RLock()
	defer po.mu.RUnlock()

	projects := make([]string, 0, len(po.projects))
	for name := range po.projects {
		projects = append(projects, name)
	}

	return projects
}

func (po *ProjectOrchestrator) GetProjectStats(projectName string) map[string]interface{} {
	po.mu.RLock()
	boulder, exists := po.projects[projectName]
	po.mu.RUnlock()

	if !exists {
		return map[string]interface{}{
			"error": "project not found",
		}
	}

	boulder.mu.RLock()
	defer boulder.mu.RUnlock()

	return map[string]interface{}{
		"project":         boulder.Project,
		"active_plan":     boulder.ActivePlan,
		"pending_tasks":   len(boulder.Tasks),
		"completed_tasks": len(boulder.CompletedTasks),
		"last_updated":    boulder.LastUpdated,
		"has_session":     boulder.CurrentSession != nil,
	}
}

func (po *ProjectOrchestrator) RotateProject() (string, error) {
	po.mu.Lock()
	defer po.mu.Unlock()

	if len(po.projects) == 0 {
		return "", fmt.Errorf("no projects available")
	}

	projectNames := make([]string, 0, len(po.projects))
	for name := range po.projects {
		projectNames = append(projectNames, name)
	}

	nextIndex := 0
	if po.currentProject != "" {
		for i, name := range projectNames {
			if name == po.currentProject {
				nextIndex = (i + 1) % len(projectNames)
				break
			}
		}
	}

	po.currentProject = projectNames[nextIndex]
	return po.currentProject, nil
}
