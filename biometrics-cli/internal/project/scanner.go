package project

import (
	"os"
	"path/filepath"
)

// DiscoverProjects scannt das Sisyphus-Verzeichnis nach Projekten mit einer boulder.json
func DiscoverProjects(basePath string) ([]string, error) {
	var projects []string

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			boulderPath := filepath.Join(basePath, entry.Name(), "boulder.json")
			if _, err := os.Stat(boulderPath); err == nil {
				projects = append(projects, entry.Name())
			}
		}
	}
	return projects, nil
}
