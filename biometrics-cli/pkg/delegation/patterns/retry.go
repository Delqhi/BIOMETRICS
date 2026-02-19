package patterns

import (
	"context"
	"time"

	"biometrics-cli/pkg/delegation"
)

type RetryPattern struct {
	task       *delegation.Task
	router     *delegation.DelegationRouter
	engine     *delegation.WorkerPool
	maxRetries int
	backoff    time.Duration
}

func NewRetryPattern(router *delegation.DelegationRouter, engine *delegation.WorkerPool, maxRetries int, backoff time.Duration) *RetryPattern {
	return &RetryPattern{
		router:     router,
		engine:     engine,
		maxRetries: maxRetries,
		backoff:    backoff,
	}
}

func (r *RetryPattern) SetTask(task *delegation.Task) {
	r.task = task
}

func (r *RetryPattern) Execute(ctx context.Context) (*delegation.TaskResult, error) {
	var lastError error

	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if attempt > 0 {
			time.Sleep(r.backoff * time.Duration(attempt))
		}

		r.engine.Submit(r.task)

		resultChan := r.engine.Results()
		select {
		case result := <-resultChan:
			if result.Success {
				return result, nil
			}
			lastError = result.Error
		case <-time.After(5 * time.Minute):
			lastError = context.DeadlineExceeded
		}
	}

	return &delegation.TaskResult{
		TaskID:  r.task.ID,
		Success: false,
		Error:   lastError,
	}, lastError
}
