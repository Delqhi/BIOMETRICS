package tracker

import (
	"fmt"
	"sync"
)

type ModelTracker struct {
	mu     sync.Mutex
	models map[string]bool
}

func NewModelTracker() *ModelTracker {
	return &ModelTracker{models: make(map[string]bool)}
}

func (mt *ModelTracker) Acquire(model string) error {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	if mt.models[model] {
		return fmt.Errorf("model %s is already in use", model)
	}
	mt.models[model] = true
	return nil
}

func (mt *ModelTracker) Release(model string) {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	delete(mt.models, model)
}
