package scheduler

import (
	"biometrics-cli/internal/metrics"
	"biometrics-cli/internal/state"
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID          string
	Name        string
	Schedule    string
	Handler     JobHandler
	Enabled     bool
	LastRun     time.Time
	NextRun     time.Time
	RunCount    int
	FailCount   int
	MaxFailures int
	Timeout     time.Duration
	RetryCount  int
	Tags        []string
}

type JobHandler func(ctx context.Context) error

type Scheduler struct {
	mu        sync.RWMutex
	jobs      map[string]*Job
	running   bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
	executors map[string]chan struct{}
}

var (
	defaultScheduler = &Scheduler{
		jobs:      make(map[string]*Job),
		stopChan:  make(chan struct{}),
		executors: make(map[string]chan struct{}),
	}
	Sched = defaultScheduler
)

func (s *Scheduler) Register(job *Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[job.ID]; exists {
		return fmt.Errorf("job %s already exists", job.ID)
	}

	if job.MaxFailures == 0 {
		job.MaxFailures = 3
	}
	if job.Timeout == 0 {
		job.Timeout = 5 * time.Minute
	}

	job.NextRun = calculateNextRun(job.Schedule, time.Now())
	s.jobs[job.ID] = job
	state.GlobalState.Log("INFO", fmt.Sprintf("Registered job: %s (%s)", job.Name, job.Schedule))

	return nil
}

func (s *Scheduler) Unregister(jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[jobID]; !exists {
		return fmt.Errorf("job %s not found", jobID)
	}

	if s.running {
		if stopChan, ok := s.executors[jobID]; ok {
			close(stopChan)
			delete(s.executors, jobID)
		}
	}

	delete(s.jobs, jobID)
	state.GlobalState.Log("INFO", fmt.Sprintf("Unregistered job: %s", jobID))

	return nil
}

func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	state.GlobalState.Log("INFO", "Scheduler started")

	for {
		select {
		case <-s.stopChan:
			s.wg.Wait()
			state.GlobalState.Log("INFO", "Scheduler stopped")
			return
		default:
			s.tick()
			time.Sleep(10 * time.Second)
		}
	}
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	close(s.stopChan)
	s.running = false

	for jobID, stopChan := range s.executors {
		close(stopChan)
		delete(s.executors, jobID)
	}

	state.GlobalState.Log("INFO", "Scheduler stopping")
}

func (s *Scheduler) tick() {
	now := time.Now()

	s.mu.RLock()
	var toRun []*Job
	for _, job := range s.jobs {
		if !job.Enabled {
			continue
		}
		if now.After(job.NextRun) || now.Equal(job.NextRun) {
			toRun = append(toRun, job)
		}
	}
	s.mu.RUnlock()

	for _, job := range toRun {
		s.runJob(job)
	}

	s.mu.Lock()
	for _, job := range s.jobs {
		if job.Enabled {
			job.NextRun = calculateNextRun(job.Schedule, time.Now())
		}
	}
	s.mu.Unlock()
}

func (s *Scheduler) runJob(job *Job) {
	job.RunCount++
	job.LastRun = time.Now()

	metrics.SchedulerJobsRunTotal.WithLabelValues(job.Name).Inc()

	stopChan := make(chan struct{})
	s.mu.Lock()
	s.executors[job.ID] = stopChan
	s.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), job.Timeout)
	defer cancel()

	go func() {
		select {
		case <-stopChan:
			cancel()
			state.GlobalState.Log("INFO", fmt.Sprintf("Job %s cancelled", job.Name))
			return
		case <-ctx.Done():
			return
		}
	}()

	state.GlobalState.Log("INFO", fmt.Sprintf("Running job: %s", job.Name))

	err := job.Handler(ctx)
	if err != nil {
		job.FailCount++
		state.GlobalState.Log("ERROR", fmt.Sprintf("Job %s failed: %v", job.Name, err))
		metrics.SchedulerJobsFailedTotal.WithLabelValues(job.Name).Inc()

		if job.FailCount >= job.MaxFailures {
			job.Enabled = false
			state.GlobalState.Log("ERROR", fmt.Sprintf("Job %s disabled due to max failures", job.Name))
		}

		if job.RetryCount > 0 && job.FailCount < job.MaxFailures {
			go s.retryJob(job)
		}
	} else {
		job.FailCount = 0
		state.GlobalState.Log("INFO", fmt.Sprintf("Job %s completed successfully", job.Name))
		metrics.SchedulerJobsSuccessTotal.WithLabelValues(job.Name).Inc()
	}

	s.mu.Lock()
	delete(s.executors, job.ID)
	s.mu.Unlock()
}

func (s *Scheduler) retryJob(job *Job) {
	for i := 0; i < job.RetryCount; i++ {
		select {
		case <-s.stopChan:
			return
		case <-time.After(time.Duration(i+1) * 30 * time.Second):
			if !job.Enabled {
				return
			}
			job.NextRun = time.Now()
			s.mu.Lock()
			job.NextRun = time.Now()
			s.mu.Unlock()
			state.GlobalState.Log("INFO", fmt.Sprintf("Retrying job: %s (attempt %d/%d)", job.Name, i+1, job.RetryCount))
			return
		}
	}
}

func (s *Scheduler) GetJob(id string) (*Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, exists := s.jobs[id]
	if !exists {
		return nil, fmt.Errorf("job %s not found", id)
	}
	return job, nil
}

func (s *Scheduler) ListJobs() []*Job {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobs := make([]*Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

func (s *Scheduler) EnableJob(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, exists := s.jobs[id]
	if !exists {
		return fmt.Errorf("job %s not found", id)
	}

	job.Enabled = true
	job.NextRun = calculateNextRun(job.Schedule, time.Now())
	state.GlobalState.Log("INFO", fmt.Sprintf("Job %s enabled", job.Name))

	return nil
}

func (s *Scheduler) DisableJob(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, exists := s.jobs[id]
	if !exists {
		return fmt.Errorf("job %s not found", id)
	}

	job.Enabled = false
	state.GlobalState.Log("INFO", fmt.Sprintf("Job %s disabled", job.Name))

	return nil
}

func calculateNextRun(schedule string, from time.Time) time.Time {
	switch schedule {
	case "@every 1m":
		return from.Add(1 * time.Minute)
	case "@every 5m":
		return from.Add(5 * time.Minute)
	case "@every 15m":
		return from.Add(15 * time.Minute)
	case "@every 30m":
		return from.Add(30 * time.Minute)
	case "@every 1h":
		return from.Add(1 * time.Hour)
	case "@every 2h":
		return from.Add(2 * time.Hour)
	case "@every 6h":
		return from.Add(6 * time.Hour)
	case "@every 12h":
		return from.Add(12 * time.Hour)
	case "@daily":
		return time.Date(from.Year(), from.Month(), from.Day()+1, 0, 0, 0, 0, from.Location())
	case "@weekly":
		daysUntilNext := time.Sunday - from.Weekday()
		if daysUntilNext == 0 {
			daysUntilNext = 7
		}
		next := from.AddDate(0, 0, int(daysUntilNext))
		return time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, from.Location())
	default:
		return from.Add(1 * time.Hour)
	}
}

func RegisterDefaultJobs() {
	Sched.Register(&Job{
		ID:          "health-check",
		Name:        "Health Check",
		Schedule:    "@every 5m",
		Enabled:     true,
		MaxFailures: 3,
		Handler:     healthCheckJob,
		Tags:        []string{"system", "health"},
	})

	Sched.Register(&Job{
		ID:          "cleanup-logs",
		Name:        "Cleanup Logs",
		Schedule:    "@daily",
		Enabled:     true,
		MaxFailures: 1,
		Handler:     cleanupLogsJob,
		Tags:        []string{"maintenance"},
	})

	Sched.Register(&Job{
		ID:          "model-pool-check",
		Name:        "Model Pool Check",
		Schedule:    "@every 15m",
		Enabled:     true,
		MaxFailures: 5,
		Handler:     modelPoolCheckJob,
		Tags:        []string{"models"},
	})

	Sched.Register(&Job{
		ID:          "metrics-rotate",
		Name:        "Metrics Rotation",
		Schedule:    "@every 1h",
		Enabled:     true,
		MaxFailures: 1,
		Handler:     metricsRotateJob,
		Tags:        []string{"metrics"},
	})
}

func healthCheckJob(ctx context.Context) error {
	state.GlobalState.Log("INFO", "Running health check job")
	return nil
}

func cleanupLogsJob(ctx context.Context) error {
	state.GlobalState.Log("INFO", "Running cleanup logs job")
	return nil
}

func modelPoolCheckJob(ctx context.Context) error {
	state.GlobalState.Log("INFO", "Running model pool check job")
	return nil
}

func metricsRotateJob(ctx context.Context) error {
	state.GlobalState.Log("INFO", "Running metrics rotation job")
	return nil
}

func Start() {
	RegisterDefaultJobs()
	Sched.Start()
}

func Stop() {
	Sched.Stop()
}
