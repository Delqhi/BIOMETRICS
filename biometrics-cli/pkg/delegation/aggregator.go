package delegation

import (
	"errors"
	"sync"
	"time"
)

type MergeStrategy string

const (
	MergeStrategyConcat MergeStrategy = "concat"
	MergeStrategyMerge  MergeStrategy = "merge"
	MergeStrategyReduce MergeStrategy = "reduce"
)

type AggregatorConfig struct {
	Strategy           MergeStrategy `json:"strategy"`
	Timeout            time.Duration `json:"timeout"`
	MinResults         int           `json:"min_results"`
	ConflictResolution string        `json:"conflict_resolution"`
}

type ResultAggregator struct {
	config     AggregatorConfig
	results    map[string][]*TaskResult
	progress   map[string]int
	totalTasks map[string]int
	mu         sync.RWMutex
	done       map[string]chan struct{}
}

func NewResultAggregator(config AggregatorConfig) *ResultAggregator {
	return &ResultAggregator{
		config:     config,
		results:    make(map[string][]*TaskResult),
		progress:   make(map[string]int),
		totalTasks: make(map[string]int),
		done:       make(map[string]chan struct{}),
	}
}

func (a *ResultAggregator) Collect(batchID string, result *TaskResult) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.results[batchID]; !exists {
		a.results[batchID] = make([]*TaskResult, 0)
		a.done[batchID] = make(chan struct{})
	}

	a.results[batchID] = append(a.results[batchID], result)
	a.progress[batchID]++

	if a.progress[batchID] >= a.totalTasks[batchID] {
		close(a.done[batchID])
	}
}

func (a *ResultAggregator) SetTotalTasks(batchID string, total int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.totalTasks[batchID] = total
}

func (a *ResultAggregator) Wait(batchID string, timeout time.Duration) ([]*TaskResult, error) {
	a.mu.RLock()
	doneChan, exists := a.done[batchID]
	a.mu.RUnlock()

	if !exists {
		return nil, errors.New("batch not found")
	}

	select {
	case <-doneChan:
		return a.GetResults(batchID)
	case <-time.After(timeout):
		return a.GetResults(batchID)
	}
}

func (a *ResultAggregator) GetResults(batchID string) ([]*TaskResult, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	results, exists := a.results[batchID]
	if !exists {
		return nil, errors.New("batch not found")
	}

	return results, nil
}

func (a *ResultAggregator) Merge(batchID string) (interface{}, error) {
	results, err := a.GetResults(batchID)
	if err != nil {
		return nil, err
	}

	switch a.config.Strategy {
	case MergeStrategyConcat:
		return a.mergeConcat(results)
	case MergeStrategyMerge:
		return a.mergeMerge(results)
	case MergeStrategyReduce:
		return a.mergeReduce(results)
	default:
		return nil, errors.New("unknown merge strategy")
	}
}

func (a *ResultAggregator) mergeConcat(results []*TaskResult) (interface{}, error) {
	merged := make([]interface{}, 0, len(results))
	for _, result := range results {
		if result.Success {
			merged = append(merged, result.Data)
		}
	}
	return merged, nil
}

func (a *ResultAggregator) mergeMerge(results []*TaskResult) (interface{}, error) {
	merged := make(map[string]interface{})
	for _, result := range results {
		if result.Success {
			if dataMap, ok := result.Data.(map[string]interface{}); ok {
				for k, v := range dataMap {
					merged[k] = v
				}
			}
		}
	}
	return merged, nil
}

func (a *ResultAggregator) mergeReduce(results []*TaskResult) (interface{}, error) {
	if len(results) == 0 {
		return nil, errors.New("no results to reduce")
	}

	var successCount int
	var lastError error

	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			lastError = result.Error
		}
	}

	return map[string]interface{}{
		"total":      len(results),
		"success":    successCount,
		"failed":     len(results) - successCount,
		"last_error": lastError,
	}, nil
}

func (a *ResultAggregator) GetProgress(batchID string) (int, int) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.progress[batchID], a.totalTasks[batchID]
}

func (a *ResultAggregator) GetErrors(batchID string) []error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	errors := make([]error, 0)
	for _, result := range a.results[batchID] {
		if !result.Success && result.Error != nil {
			errors = append(errors, result.Error)
		}
	}
	return errors
}
