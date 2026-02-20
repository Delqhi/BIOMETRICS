// Package shutdown provides graceful shutdown handling for the BIOMETRICS CLI.
// It supports signal handling, timeout management, resource cleanup registration,
// and shutdown hooks for proper application termination.
package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"biometrics-cli/pkg/logging"
)

// ShutdownManager handles graceful application shutdown.
type ShutdownManager struct {
	mu            sync.Mutex
	hooks         []ShutdownHook
	timeout       time.Duration
	ctx           context.Context
	cancel        context.CancelFunc
	shutdownCh    chan struct{}
	shutdownOnce  sync.Once
	isShuttingDown bool
	logger        *logging.Logger
}

// ShutdownHook represents a function to be called during shutdown.
type ShutdownHook struct {
	// Name is a descriptive name for the hook.
	Name string
	// Priority determines the order of execution (higher = earlier).
	Priority int
	// Timeout is the maximum time this hook is allowed to run.
	Timeout time.Duration
	// Fn is the function to execute.
	Fn func(ctx context.Context) error
}

// ShutdownManagerConfig contains configuration for the ShutdownManager.
type ShutdownManagerConfig struct {
	// Timeout is the default timeout for all shutdown operations.
	Timeout time.Duration
	// Logger is the logger to use for shutdown events.
	Logger *logging.Logger
}

// DefaultShutdownTimeout is the default timeout for shutdown operations.
const DefaultShutdownTimeout = 30 * time.Second

// NewShutdownManager creates a new ShutdownManager.
func NewShutdownManager(config *ShutdownManagerConfig) *ShutdownManager {
	if config == nil {
		config = &ShutdownManagerConfig{}
	}

	timeout := config.Timeout
	if timeout <= 0 {
		timeout = DefaultShutdownTimeout
	}

	logger := config.Logger
	if logger == nil {
		logger = logging.Default()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &ShutdownManager{
		hooks:      make([]ShutdownHook, 0),
		timeout:    timeout,
		ctx:        ctx,
		cancel:     cancel,
		shutdownCh: make(chan struct{}),
		logger:     logger,
	}
}

// RegisterHook registers a shutdown hook to be executed during shutdown.
func (m *ShutdownManager) RegisterHook(name string, priority int, timeout time.Duration, fn func(ctx context.Context) error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if timeout <= 0 {
		timeout = m.timeout
	}

	hook := ShutdownHook{
		Name:     name,
		Priority: priority,
		Timeout:  timeout,
		Fn:       fn,
	}

	m.hooks = append(m.hooks, hook)
	m.logger.Debug("Registered shutdown hook",
		logging.String("name", name),
		logging.Int("priority", priority),
		logging.Duration("timeout", timeout),
	)
}

// RegisterHookWithTimeout registers a shutdown hook with a specific timeout.
func (m *ShutdownManager) RegisterHookWithTimeout(name string, fn func(ctx context.Context) error, timeout time.Duration) {
	m.RegisterHook(name, 0, timeout, fn)
}

// RegisterHookDefault registers a shutdown hook with default timeout.
func (m *ShutdownManager) RegisterHookDefault(name string, fn func(ctx context.Context) error) {
	m.RegisterHook(name, 0, 0, fn)
}

// WaitForSignal starts listening for shutdown signals and initiates graceful shutdown.
func (m *ShutdownManager) WaitForSignal(signals ...os.Signal) <-chan struct{} {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, signals...)

	go func() {
		sig := <-sigCh
		m.logger.Info("Received shutdown signal",
			logging.String("signal", sig.String()),
		)
		m.InitiateShutdown()
	}()

	m.logger.Info("Shutdown manager initialized",
		logging.Duration("timeout", m.timeout),
		logging.Int("hooks", len(m.hooks)),
	)

	return m.shutdownCh
}

// InitiateShutdown starts the graceful shutdown process.
func (m *ShutdownManager) InitiateShutdown() {
	m.shutdownOnce.Do(func() {
		m.mu.Lock()
		m.isShuttingDown = true
		m.mu.Unlock()

		m.logger.Info("Initiating graceful shutdown",
			logging.Int("hooks", len(m.hooks)),
			logging.Duration("timeout", m.timeout),
		)

		startTime := time.Now()
		ctx, cancel := context.WithTimeout(m.ctx, m.timeout)
		defer cancel()

		m.executeHooks(ctx)

		elapsed := time.Since(startTime)
		m.logger.Info("Graceful shutdown completed",
			logging.Duration("elapsed", elapsed),
		)

		m.cancel()
		close(m.shutdownCh)
	})
}

// executeHooks executes all registered hooks in priority order.
func (m *ShutdownManager) executeHooks(ctx context.Context) {
	m.mu.Lock()
	hooks := make([]ShutdownHook, len(m.hooks))
	copy(hooks, m.hooks)
	m.mu.Unlock()

	// Sort by priority (higher first)
	for i := 0; i < len(hooks)-1; i++ {
		for j := i + 1; j < len(hooks); j++ {
			if hooks[i].Priority < hooks[j].Priority {
				hooks[i], hooks[j] = hooks[j], hooks[i]
			}
		}
	}

	m.logger.Debug("Executing shutdown hooks",
		logging.Int("count", len(hooks)),
	)

	var wg sync.WaitGroup
	results := make(chan hookResult, len(hooks))

	for _, hook := range hooks {
		if ctx.Err() != nil {
			m.logger.Warn("Shutdown context cancelled, skipping remaining hooks",
				logging.String("error", ctx.Err().Error()),
			)
			break
		}

		wg.Add(1)
		go func(hook ShutdownHook) {
			defer wg.Done()
			m.executeHook(ctx, hook, results)
		}(hook)
	}

	wg.Wait()
	close(results)

	// Log results
	failed := 0
	for result := range results {
		if result.err != nil {
			failed++
			m.logger.Error("Shutdown hook failed",
				logging.String("name", result.name),
				logging.Err(result.err),
				logging.Duration("elapsed", result.elapsed),
			)
		} else {
			m.logger.Debug("Shutdown hook completed",
				logging.String("name", result.name),
				logging.Duration("elapsed", result.elapsed),
			)
		}
	}

	if failed > 0 {
		m.logger.Warn("Some shutdown hooks failed",
			logging.Int("failed", failed),
			logging.Int("total", len(hooks)),
		)
	}
}

// hookResult contains the result of a hook execution.
type hookResult struct {
	name     string
	err      error
	elapsed  time.Duration
}

// executeHook executes a single hook with timeout handling.
func (m *ShutdownManager) executeHook(ctx context.Context, hook ShutdownHook, results chan<- hookResult) {
	startTime := time.Now()

	hookCtx, cancel := context.WithTimeout(ctx, hook.Timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- fmt.Errorf("panic in hook %s: %v", hook.Name, r)
			}
		}()
		done <- hook.Fn(hookCtx)
	}()

	select {
	case err := <-done:
		results <- hookResult{
			name:    hook.Name,
			err:     err,
			elapsed: time.Since(startTime),
		}
	case <-hookCtx.Done():
		results <- hookResult{
			name:    hook.Name,
			err:     fmt.Errorf("hook timed out after %v", hook.Timeout),
			elapsed: time.Since(startTime),
		}
	case <-ctx.Done():
		results <- hookResult{
			name:    hook.Name,
			err:     ctx.Err(),
			elapsed: time.Since(startTime),
		}
	}
}

// IsShuttingDown returns true if shutdown is in progress.
func (m *ShutdownManager) IsShuttingDown() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isShuttingDown
}

// Context returns the shutdown context that is cancelled when shutdown begins.
func (m *ShutdownManager) Context() context.Context {
	return m.ctx
}

// Timeout returns the configured shutdown timeout.
func (m *ShutdownManager) Timeout() time.Duration {
	return m.timeout
}

// WithShutdownContext wraps a context to be cancelled when shutdown begins.
func (m *ShutdownManager) WithShutdownContext(parent context.Context) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	ctx, cancel := context.WithCancel(parent)

	go func() {
		select {
		case <-m.ctx.Done():
			cancel()
		case <-parent.Done():
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx
}

// ShutdownFunc is a convenience function that creates a ShutdownManager,
// registers hooks, waits for signals, and performs shutdown.
func ShutdownFunc(timeout time.Duration, hooks ...ShutdownHook) func() {
	manager := NewShutdownManager(&ShutdownManagerConfig{
		Timeout: timeout,
	})

	for _, hook := range hooks {
		manager.RegisterHook(hook.Name, hook.Priority, hook.Timeout, hook.Fn)
	}

	manager.WaitForSignal()

	return func() {
		manager.InitiateShutdown()
	}
}

// CleanupResource is a helper that creates a shutdown hook for resource cleanup.
func CleanupResource(name string, priority int, cleanup func() error) ShutdownHook {
	return ShutdownHook{
		Name:     name,
		Priority: priority,
		Timeout:  10 * time.Second,
		Fn: func(ctx context.Context) error {
			done := make(chan error, 1)
			go func() {
				done <- cleanup()
			}()

			select {
			case err := <-done:
				return err
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	}
}

// CloseCloser is a helper that creates a shutdown hook to close an io.Closer.
func CloseCloser(name string, priority int, closer interface{ Close() error }) ShutdownHook {
	return CleanupResource(name, priority, closer.Close)
}

// StopFunc is a helper that creates a shutdown hook to call a stop function.
func StopFunc(name string, priority int, stop func() error) ShutdownHook {
	return CleanupResource(name, priority, stop)
}

// Default manager for package-level functions.
var defaultManager = NewShutdownManager(nil)

// DefaultManager returns the default shutdown manager.
func DefaultManager() *ShutdownManager {
	return defaultManager
}

// SetDefaultManager sets the default shutdown manager.
func SetDefaultManager(m *ShutdownManager) {
	defaultManager = m
}

// RegisterHook registers a hook with the default manager.
func RegisterHook(name string, priority int, timeout time.Duration, fn func(ctx context.Context) error) {
	defaultManager.RegisterHook(name, priority, timeout, fn)
}

// WaitForSignal starts signal handling on the default manager.
func WaitForSignal(signals ...os.Signal) <-chan struct{} {
	return defaultManager.WaitForSignal(signals...)
}

// InitiateShutdown initiates shutdown on the default manager.
func InitiateShutdown() {
	defaultManager.InitiateShutdown()
}

// Context returns the default manager's context.
func Context() context.Context {
	return defaultManager.Context()
}

// IsShuttingDown returns true if the default manager is shutting down.
func IsShuttingDown() bool {
	return defaultManager.IsShuttingDown()
}
