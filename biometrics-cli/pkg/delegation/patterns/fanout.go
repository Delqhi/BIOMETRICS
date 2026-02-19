package patterns

import (
	"context"
	"sync"
	"time"

	"biometrics-cli/pkg/delegation"
)

type FanOutPattern struct {
	tasks      []*delegation.Task
	router     *delegation.DelegationRouter
	engine     *delegation.WorkerPool
	aggregator *delegation.ResultAggregator
	batchID    string
}

func NewFanOutPattern(router *delegation.DelegationRouter, engine *delegation.WorkerPool, batchID string) *FanOutPattern {
	aggregator := delegation.NewResultAggregator(delegation.AggregatorConfig{
		Strategy: delegation.MergeStrategyConcat,
		Timeout:  10 * time.Minute,
	})

	return &FanOutPattern{
		tasks:      make([]*delegation.Task, 0),
		router:     router,
		engine:     engine,
		aggregator: aggregator,
		batchID:    batchID,
	}
}

func (f *FanOutPattern) AddTask(task *delegation.Task) {
	f.tasks = append(f.tasks, task)
}

func (f *FanOutPattern) Execute(ctx context.Context) ([]*delegation.TaskResult, error) {
	f.aggregator.SetTotalTasks(f.batchID, len(f.tasks))

	var wg sync.WaitGroup

	for _, task := range f.tasks {
		wg.Add(1)
		go func(t *delegation.Task) {
			defer wg.Done()
			f.engine.Submit(t)

			resultChan := f.engine.Results()
			select {
			case result := <-resultChan:
				f.aggregator.Collect(f.batchID, result)
			case <-time.After(5 * time.Minute):
			}
		}(task)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return f.aggregator.GetResults(f.batchID)
	case <-time.After(10 * time.Minute):
		return f.aggregator.GetResults(f.batchID)
	}
}

func (f *FanOutPattern) Merge() (interface{}, error) {
	return f.aggregator.Merge(f.batchID)
}
