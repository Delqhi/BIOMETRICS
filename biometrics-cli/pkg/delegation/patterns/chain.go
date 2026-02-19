package patterns

import (
	"context"
	"time"

	"biometrics-cli/pkg/delegation"
)

type ChainPattern struct {
	steps  []*delegation.Task
	router *delegation.DelegationRouter
	engine *delegation.WorkerPool
}

func NewChainPattern(router *delegation.DelegationRouter, engine *delegation.WorkerPool) *ChainPattern {
	return &ChainPattern{
		steps:  make([]*delegation.Task, 0),
		router: router,
		engine: engine,
	}
}

func (c *ChainPattern) AddStep(task *delegation.Task) {
	c.steps = append(c.steps, task)
}

func (c *ChainPattern) Execute(ctx context.Context) (*delegation.TaskResult, error) {
	var lastResult *delegation.TaskResult

	for i, step := range c.steps {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if i > 0 && lastResult != nil {
			step.SetContext("previous_result", lastResult.Data)
		}

		c.engine.Submit(step)

		resultChan := c.engine.Results()
		select {
		case result := <-resultChan:
			lastResult = result
			if !result.Success {
				return lastResult, result.Error
			}
		case <-time.After(5 * time.Minute):
			return nil, context.DeadlineExceeded
		}
	}

	return lastResult, nil
}
