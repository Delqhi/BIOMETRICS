package patterns

import (
	"context"
	"time"

	"biometrics-cli/pkg/delegation"
)

type FanInPattern struct {
	resultTasks []*delegation.Task
	finalTask   *delegation.Task
	router      *delegation.DelegationRouter
	engine      *delegation.WorkerPool
	aggregator  *delegation.ResultAggregator
	batchID     string
}

func NewFanInPattern(router *delegation.DelegationRouter, engine *delegation.WorkerPool, batchID string) *FanInPattern {
	aggregator := delegation.NewResultAggregator(delegation.AggregatorConfig{
		Strategy: delegation.MergeStrategyMerge,
		Timeout:  10 * time.Minute,
	})

	return &FanInPattern{
		resultTasks: make([]*delegation.Task, 0),
		router:      router,
		engine:      engine,
		aggregator:  aggregator,
		batchID:     batchID,
	}
}

func (f *FanInPattern) AddResultTask(task *delegation.Task) {
	f.resultTasks = append(f.resultTasks, task)
}

func (f *FanInPattern) SetFinalTask(task *delegation.Task) {
	f.finalTask = task
}

func (f *FanInPattern) Execute(ctx context.Context) (*delegation.TaskResult, error) {
	f.aggregator.SetTotalTasks(f.batchID, len(f.resultTasks))

	for _, task := range f.resultTasks {
		f.engine.Submit(task)
	}

	_, err := f.aggregator.Wait(f.batchID, 10*time.Minute)
	if err != nil {
		return nil, err
	}

	merged, err := f.aggregator.Merge(f.batchID)
	if err != nil {
		return nil, err
	}

	if f.finalTask != nil {
		f.finalTask.SetContext("merged_results", merged)
		f.engine.Submit(f.finalTask)

		resultChan := f.engine.Results()
		select {
		case result := <-resultChan:
			return result, nil
		case <-time.After(5 * time.Minute):
			return nil, context.DeadlineExceeded
		}
	}

	return &delegation.TaskResult{
		TaskID:  f.batchID,
		Success: true,
		Data:    merged,
	}, nil
}
