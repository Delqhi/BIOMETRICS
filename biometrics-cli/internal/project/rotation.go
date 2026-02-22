package project

import (
	"sync"
)

// RotationScheduler verwaltet die Round-Robin Projektverwaltung.
type RotationScheduler struct {
	projects []string
	current  int
	mu       sync.Mutex
}

// NewRotationScheduler erstellt einen neuen RotationScheduler.
func NewRotationScheduler(projects []string) *RotationScheduler {
	return &RotationScheduler{
		projects: projects,
		current:  0,
	}
}

// NextProject gibt das nächste Projekt im Round-Robin Verfahren zurück.
func (r *RotationScheduler) NextProject() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.projects) == 0 {
		return ""
	}

	project := r.projects[r.current]
	r.current = (r.current + 1) % len(r.projects)

	return project
}
