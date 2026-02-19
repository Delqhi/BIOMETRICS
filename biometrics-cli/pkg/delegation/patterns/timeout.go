package patterns

import (
	"context"
	"time"

	"biometrics-cli/pkg/delegation"
)

type TimeoutPattern struct {
	task    *delegation.Task
	router  *delegation.DelegationRouter
	engine  *delegation.WorkerPool
	timeout time.Duration
}

func NewTimeoutPattern(router *delegation.DelegationRouter, engine *delegation.WorkerPool, timeout time.Duration) *TimeoutPattern {
	return &TimeoutPattern{
		router:  router,
		engine:  engine,
		timeout: timeout,
	}
}

func (t *TimeoutPattern) SetTask(task *delegation.Task) {
	t.task = task
}

func (t *TimeoutPattern) Execute(ctx context.Context) (*delegation.TaskResult, error) {
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	t.engine.Submit(t.task)

	resultChan := t.engine.Results()
	select {
	case result := <-resultChan:
		return result, nil
	case <-ctx.Done():
		return &delegation.TaskResult{
			TaskID:  t.task.ID,
			Success: false,
			Error:   context.DeadlineExceeded,
		}, context.DeadlineExceeded
	}
}
