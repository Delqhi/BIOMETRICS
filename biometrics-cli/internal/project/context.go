package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"biometrics-cli/internal/models"
)

type ProjectContext struct {
	ID       string
	BasePath string
	Boulder  *models.Boulder
}

// LoadProjectContext lädt die projektspezifische boulder.json
func LoadProjectContext(projectID string) (*ProjectContext, error) {
	basePath := filepath.Join("/Users/jeremy/.sisyphus/plans", projectID)
	boulderPath := filepath.Join(basePath, "boulder.json")

	data, err := os.ReadFile(boulderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read boulder for project %s: %w", projectID, err)
	}

	var b models.Boulder
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("failed to parse boulder for project %s: %w", projectID, err)
	}

	return &ProjectContext{
		ID:       projectID,
		BasePath: basePath,
		Boulder:  &b,
	}, nil
}
