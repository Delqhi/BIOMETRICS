package delegation

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TaskResult struct {
	TaskID   string        `json:"task_id"`
	Success  bool          `json:"success"`
	Data     interface{}   `json:"data"`
	Error    error         `json:"error,omitempty"`
	Duration time.Duration `json:"duration"`
	AgentID  string        `json:"agent_id"`
}

type WorkerPool struct {
	poolSize   int
	taskChan   chan *Task
	resultChan chan *TaskResult
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	router     *DelegationRouter
	mu         sync.RWMutex
}

func NewWorkerPool(poolSize int, router *DelegationRouter) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	wp := &WorkerPool{
		poolSize:   poolSize,
		taskChan:   make(chan *Task, poolSize*2),
		resultChan: make(chan *TaskResult, poolSize*2),
		ctx:        ctx,
		cancel:     cancel,
		router:     router,
	}

	wp.startWorkers()

	return wp
}

func (wp *WorkerPool) startWorkers() {
	for i := 0; i < wp.poolSize; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case task := <-wp.taskChan:
			wp.executeTask(task, id)
		}
	}
}

func (wp *WorkerPool) executeTask(task *Task, workerID int) {
	startTime := time.Now()

	ctx, cancel := context.WithTimeout(wp.ctx, 10*time.Minute)
	defer cancel()

	agentID, err := wp.router.Route(task)
	if err != nil {
		wp.resultChan <- &TaskResult{
			TaskID:   task.ID,
			Success:  false,
			Error:    err,
			Duration: time.Since(startTime),
		}
		return
	}

	task.SetStatus(TaskStatusRunning)

	result := wp.executeWithAgent(ctx, task, agentID, workerID)
	result.Duration = time.Since(startTime)
	result.AgentID = agentID

	if result.Success {
		wp.router.RecordSuccess(agentID)
		task.SetStatus(TaskStatusCompleted)
	} else {
		wp.router.RecordFailure(agentID)
		if task.IncrementRetry() {
			task.SetStatus(TaskStatusPending)
			wp.Submit(task)
			return
		}
		task.SetStatus(TaskStatusFailed)
	}

	wp.resultChan <- result
}

func (wp *WorkerPool) executeWithAgent(ctx context.Context, task *Task, agentID string, workerID int) *TaskResult {
	select {
	case <-ctx.Done():
		return &TaskResult{
			TaskID:  task.ID,
			Success: false,
			Error:   ctx.Err(),
		}
	default:
		fmt.Printf("[Worker %d] Executing task %s on agent %s\n", workerID, task.ID, agentID)

		time.Sleep(100 * time.Millisecond)

		return &TaskResult{
			TaskID:  task.ID,
			Success: true,
			Data:    fmt.Sprintf("Task %s completed by agent %s", task.ID, agentID),
		}
	}
}

func (wp *WorkerPool) Submit(task *Task) {
	wp.taskChan <- task
}

func (wp *WorkerPool) Results() <-chan *TaskResult {
	return wp.resultChan
}

func (wp *WorkerPool) Shutdown() {
	wp.cancel()
	wp.wg.Wait()
	close(wp.taskChan)
	close(wp.resultChan)
}

func (wp *WorkerPool) GetStats() map[string]interface{} {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	return map[string]interface{}{
		"pool_size":    wp.poolSize,
		"task_queue":   len(wp.taskChan),
		"result_queue": len(wp.resultChan),
	}
}
