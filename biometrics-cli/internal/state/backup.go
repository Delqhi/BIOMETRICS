package state

import (
	"biometrics-cli/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type BackupManager struct {
	mu            sync.RWMutex
	backupDir     string
	maxBackups    int
	compression   bool
	encryption    bool
	encryptionKey []byte
}

type BackupMetadata struct {
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Size        int64     `json:"size"`
	Checksum    string    `json:"checksum"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
}

type Backup struct {
	Metadata BackupMetadata     `json:"metadata"`
	State    *OrchestratorState `json:"state"`
	Sessions []SessionState     `json:"sessions"`
	Tasks    []TaskState        `json:"tasks"`
	Metrics  interface{}        `json:"metrics"`
}

type OrchestratorState struct {
	ActiveAgents    map[string]AgentState `json:"active_agents"`
	ModelPool       ModelPoolState        `json:"model_pool"`
	TaskQueue       []TaskState           `json:"task_queue"`
	CircuitBreakers map[string]CBState    `json:"circuit_breakers"`
	Timestamp       time.Time             `json:"timestamp"`
}

type AgentState struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Status    string                 `json:"status"`
	Model     string                 `json:"model"`
	StartedAt time.Time              `json:"started_at"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type ModelPoolState struct {
	AvailableModels    []string          `json:"available_models"`
	ActiveAcquisitions map[string]string `json:"active_acquisitions"`
}

type CBState struct {
	Name     string `json:"name"`
	State    int    `json:"state"`
	Failures int    `json:"failures"`
}

type SessionState struct {
	ID        string                 `json:"id"`
	AgentID   string                 `json:"agent_id"`
	Status    string                 `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type TaskState struct {
	ID          string                 `json:"id"`
	AgentID     string                 `json:"agent_id"`
	Priority    int                    `json:"priority"`
	Status      string                 `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
}

func NewBackupManager(backupDir string, maxBackups int) *BackupManager {
	os.MkdirAll(backupDir, 0755)
	return &BackupManager{
		backupDir:   backupDir,
		maxBackups:  maxBackups,
		compression: true,
		encryption:  false,
	}
}

func (bm *BackupManager) CreateBackup(state *OrchestratorState, description string) (string, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	backup := &Backup{
		Metadata: BackupMetadata{
			ID:          fmt.Sprintf("backup-%d", time.Now().Unix()),
			Timestamp:   time.Now(),
			Version:     "1.0.0",
			Description: description,
		},
		State: state,
	}

	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal backup: %w", err)
	}

	filename := fmt.Sprintf("%s/%s.json", bm.backupDir, backup.Metadata.ID)
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write backup file: %w", err)
	}

	info, _ := os.Stat(filename)
	backup.Metadata.Size = info.Size()

	bm.cleanupOldBackups()

	return backup.Metadata.ID, nil
}

func (bm *BackupManager) RestoreBackup(id string) (*Backup, error) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	filename := fmt.Sprintf("%s/%s.json", bm.backupDir, id)
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup file: %w", err)
	}

	var backup Backup
	if err := json.Unmarshal(data, &backup); err != nil {
		return nil, fmt.Errorf("failed to unmarshal backup: %w", err)
	}

	return &backup, nil
}

func (bm *BackupManager) ListBackups() ([]BackupMetadata, error) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	files, err := filepath.Glob(fmt.Sprintf("%s/backup-*.json", bm.backupDir))
	if err != nil {
		return nil, err
	}

	backups := make([]BackupMetadata, 0, len(files))
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var backup Backup
		if err := json.Unmarshal(data, &backup); err == nil {
			backups = append(backups, backup.Metadata)
		}
	}

	return backups, nil
}

func (bm *BackupManager) DeleteBackup(id string) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	filename := fmt.Sprintf("%s/%s.json", bm.backupDir, id)
	return os.Remove(filename)
}

func (bm *BackupManager) cleanupOldBackups() {
	files, err := filepath.Glob(fmt.Sprintf("%s/backup-*.json", bm.backupDir))
	if err != nil {
		return
	}

	if len(files) <= bm.maxBackups {
		return
	}

	type fileInfo struct {
		path    string
		modTime time.Time
	}

	infos := make([]fileInfo, 0, len(files))
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			continue
		}
		infos = append(infos, fileInfo{f, info.ModTime()})
	}

	for i := 0; i < len(infos)-bm.maxBackups; i++ {
		os.Remove(infos[i].path)
	}
}

func (bm *BackupManager) AutoBackup(ctx context.Context, interval time.Duration, stateProvider func() *OrchestratorState) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				state := stateProvider()
				bm.CreateBackup(state, "auto-backup")
			}
		}
	}()
}

func (bm *BackupManager) GetBackupDir() string {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return bm.backupDir
}

func CaptureCurrentState(agents map[string]*models.Agent, taskQueue []models.Task) *OrchestratorState {
	state := &OrchestratorState{
		ActiveAgents:    make(map[string]AgentState),
		ModelPool:       ModelPoolState{},
		TaskQueue:       make([]TaskState, 0),
		CircuitBreakers: make(map[string]CBState),
		Timestamp:       time.Now(),
	}

	for id, agent := range agents {
		state.ActiveAgents[id] = AgentState{
			ID:        agent.ID,
			Name:      agent.Name,
			Status:    agent.Status,
			Model:     agent.Model,
			StartedAt: agent.StartedAt,
		}
	}

	for _, task := range taskQueue {
		state.TaskQueue = append(state.TaskQueue, TaskState{
			ID:       task.ID,
			AgentID:  task.AgentID,
			Priority: task.Priority,
			Status:   task.Status,
		})
	}

	return state
}
