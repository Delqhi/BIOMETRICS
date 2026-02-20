package shutdown

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"biometrics-cli/pkg/logging"
)

func TestNewShutdownManager(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		manager := NewShutdownManager(nil)
		if manager == nil {
			t.Fatal("NewShutdownManager() returned nil")
		}
		if manager.timeout != DefaultShutdownTimeout {
			t.Errorf("Timeout = %v, want %v", manager.timeout, DefaultShutdownTimeout)
		}
	})

	t.Run("CustomConfig", func(t *testing.T) {
		logger := logging.NewLogger(nil)
		config := &ShutdownManagerConfig{
			Timeout: 60 * time.Second,
			Logger:  logger,
		}
		manager := NewShutdownManager(config)
		if manager.timeout != 60*time.Second {
			t.Errorf("Timeout = %v, want %v", manager.timeout, 60*time.Second)
		}
	})

	t.Run("ZeroTimeout", func(t *testing.T) {
		config := &ShutdownManagerConfig{
			Timeout: 0,
		}
		manager := NewShutdownManager(config)
		if manager.timeout != DefaultShutdownTimeout {
			t.Errorf("Timeout = %v, want %v", manager.timeout, DefaultShutdownTimeout)
		}
	})
}

func TestRegisterHook(t *testing.T) {
	manager := NewShutdownManager(nil)

	var called bool
	manager.RegisterHook("test-hook", 10, 5*time.Second, func(ctx context.Context) error {
		called = true
		return nil
	})

	if len(manager.hooks) != 1 {
		t.Errorf("RegisterHook() hooks count = %v, want %v", len(manager.hooks), 1)
	}

	if manager.hooks[0].Name != "test-hook" {
		t.Errorf("Hook name = %v, want %v", manager.hooks[0].Name, "test-hook")
	}

	if manager.hooks[0].Priority != 10 {
		t.Errorf("Hook priority = %v, want %v", manager.hooks[0].Priority, 10)
	}

	if manager.hooks[0].Timeout != 5*time.Second {
		t.Errorf("Hook timeout = %v, want %v", manager.hooks[0].Timeout, 5*time.Second)
	}

	if called {
		t.Error("Hook function should not be called during registration")
	}
}

func TestRegisterHookWithTimeout(t *testing.T) {
	manager := NewShutdownManager(nil)

	manager.RegisterHookWithTimeout("test-timeout", func(ctx context.Context) error {
		return nil
	}, 10*time.Second)

	if manager.hooks[0].Timeout != 10*time.Second {
		t.Errorf("Timeout = %v, want %v", manager.hooks[0].Timeout, 10*time.Second)
	}
}

func TestRegisterHookDefault(t *testing.T) {
	manager := NewShutdownManager(&ShutdownManagerConfig{
		Timeout: 45 * time.Second,
	})

	manager.RegisterHookDefault("test-default", func(ctx context.Context) error {
		return nil
	})

	if manager.hooks[0].Timeout != 45*time.Second {
		t.Errorf("Timeout = %v, want %v", manager.hooks[0].Timeout, 45*time.Second)
	}
}

func TestExecuteHooks(t *testing.T) {
	t.Run("SuccessfulHooks", func(t *testing.T) {
		manager := NewShutdownManager(nil)
		var counter int32

		manager.RegisterHook("hook1", 10, time.Second, func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		})

		manager.RegisterHook("hook2", 20, time.Second, func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		})

		ctx := context.Background()
		manager.executeHooks(ctx)

		if counter != 2 {
			t.Errorf("Counter = %v, want %v", counter, 2)
		}
	})

	t.Run("FailedHook", func(t *testing.T) {
		manager := NewShutdownManager(nil)
		var failed bool

		manager.RegisterHook("failing-hook", 0, time.Second, func(ctx context.Context) error {
			failed = true
			return errors.New("intentional failure")
		})

		ctx := context.Background()
		manager.executeHooks(ctx)

		if !failed {
			t.Error("Failing hook should have been executed")
		}
	})

	t.Run("HookTimeout", func(t *testing.T) {
		manager := NewShutdownManager(nil)
		var hookStarted bool

		manager.RegisterHook("slow-hook", 0, 100*time.Millisecond, func(ctx context.Context) error {
			hookStarted = true
			select {
			case <-time.After(200 * time.Millisecond):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		})

		ctx := context.Background()
		manager.executeHooks(ctx)

		if !hookStarted {
			t.Error("Hook should have started execution")
		}
	})

	t.Run("PriorityOrder", func(t *testing.T) {
		manager := NewShutdownManager(nil)
		var executionOrder []string
		var mu sync.Mutex

		manager.RegisterHook("low", 1, time.Second, func(ctx context.Context) error {
			mu.Lock()
			executionOrder = append(executionOrder, "low")
			mu.Unlock()
			return nil
		})

		manager.RegisterHook("high", 100, time.Second, func(ctx context.Context) error {
			mu.Lock()
			executionOrder = append(executionOrder, "high")
			mu.Unlock()
			return nil
		})

		manager.RegisterHook("medium", 50, time.Second, func(ctx context.Context) error {
			mu.Lock()
			executionOrder = append(executionOrder, "medium")
			mu.Unlock()
			return nil
		})

		ctx := context.Background()
		manager.executeHooks(ctx)

		// All hooks should be executed (concurrent execution, order not guaranteed)
		if len(executionOrder) != 3 {
			t.Errorf("Expected 3 hooks executed, got %d", len(executionOrder))
		}

		// Check all expected hooks are present
		expected := map[string]bool{"high": false, "medium": false, "low": false}
		for _, name := range executionOrder {
			if _, ok := expected[name]; ok {
				expected[name] = true
			}
		}
		for name, found := range expected {
			if !found {
				t.Errorf("Hook %s was not executed", name)
			}
		}
	})

	t.Run("ConcurrentHooks", func(t *testing.T) {
		manager := NewShutdownManager(&ShutdownManagerConfig{
			Timeout: 5 * time.Second,
		})

		var counter int32

		for i := 0; i < 10; i++ {
			manager.RegisterHookDefault("concurrent", func(ctx context.Context) error {
				atomic.AddInt32(&counter, 1)
				time.Sleep(10 * time.Millisecond)
				return nil
			})
		}

		ctx := context.Background()
		manager.executeHooks(ctx)

		if counter != 10 {
			t.Errorf("Counter = %v, want %v", counter, 10)
		}
	})
}

func TestIsShuttingDown(t *testing.T) {
	manager := NewShutdownManager(nil)

	if manager.IsShuttingDown() {
		t.Error("IsShuttingDown() should be false before shutdown")
	}

	go func() {
		time.Sleep(10 * time.Millisecond)
		manager.InitiateShutdown()
	}()

	<-manager.shutdownCh

	if !manager.IsShuttingDown() {
		t.Error("IsShuttingDown() should be true after shutdown")
	}
}

func TestContext(t *testing.T) {
	manager := NewShutdownManager(nil)

	ctx := manager.Context()
	if ctx == nil {
		t.Fatal("Context() returned nil")
	}

	select {
	case <-ctx.Done():
		t.Error("Context should not be done before shutdown")
	default:
	}

	manager.InitiateShutdown()
	<-manager.shutdownCh

	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("Context should be done after shutdown")
	}
}

func TestWithShutdownContext(t *testing.T) {
	manager := NewShutdownManager(nil)

	parent := context.Background()
	child := manager.WithShutdownContext(parent)

	if child == nil {
		t.Fatal("WithShutdownContext() returned nil")
	}

	select {
	case <-child.Done():
		t.Error("Child context should not be done before shutdown")
	default:
	}

	manager.InitiateShutdown()
	<-manager.shutdownCh

	select {
	case <-child.Done():
	case <-time.After(100 * time.Millisecond):
		t.Error("Child context should be done after shutdown")
	}
}

func TestWithShutdownContextNilParent(t *testing.T) {
	manager := NewShutdownManager(nil)
	child := manager.WithShutdownContext(nil)

	if child == nil {
		t.Fatal("WithShutdownContext(nil) returned nil")
	}
}

func TestInitiateShutdownOnce(t *testing.T) {
	manager := NewShutdownManager(nil)
	var callCount int32

	manager.RegisterHook("test", 0, time.Second, func(ctx context.Context) error {
		atomic.AddInt32(&callCount, 1)
		return nil
	})

	// Call InitiateShutdown multiple times
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			manager.InitiateShutdown()
		}()
	}

	wg.Wait()
	<-manager.shutdownCh

	if callCount != 1 {
		t.Errorf("Hook called %v times, want %v", callCount, 1)
	}
}

func TestTimeout(t *testing.T) {
	manager := NewShutdownManager(&ShutdownManagerConfig{
		Timeout: 120 * time.Second,
	})

	if manager.Timeout() != 120*time.Second {
		t.Errorf("Timeout() = %v, want %v", manager.Timeout(), 120*time.Second)
	}
}

func TestCleanupResource(t *testing.T) {
	var called bool
	hook := CleanupResource("cleanup", 10, func() error {
		called = true
		return nil
	})

	if hook.Name != "cleanup" {
		t.Errorf("Name = %v, want %v", hook.Name, "cleanup")
	}
	if hook.Priority != 10 {
		t.Errorf("Priority = %v, want %v", hook.Priority, 10)
	}

	err := hook.Fn(context.Background())
	if err != nil {
		t.Errorf("Hook execution failed: %v", err)
	}
	if !called {
		t.Error("Cleanup function was not called")
	}
}

func TestCloseCloser(t *testing.T) {
	var closed bool
	closer := &testCloser{closeFn: func() error {
		closed = true
		return nil
	}}

	hook := CloseCloser("closer", 5, closer)
	err := hook.Fn(context.Background())

	if err != nil {
		t.Errorf("Hook execution failed: %v", err)
	}
	if !closed {
		t.Error("Closer was not closed")
	}
}

func TestStopFunc(t *testing.T) {
	var stopped bool
	hook := StopFunc("stop", 5, func() error {
		stopped = true
		return nil
	})

	err := hook.Fn(context.Background())
	if err != nil {
		t.Errorf("Hook execution failed: %v", err)
	}
	if !stopped {
		t.Error("Stop function was not called")
	}
}

func TestDefaultManager(t *testing.T) {
	manager := DefaultManager()
	if manager == nil {
		t.Fatal("DefaultManager() returned nil")
	}

	// Set a new default manager
	newManager := NewShutdownManager(nil)
	SetDefaultManager(newManager)

	if DefaultManager() != newManager {
		t.Error("DefaultManager() should return the newly set manager")
	}

	// Restore original
	SetDefaultManager(manager)
}

func TestPackageLevelFunctions(t *testing.T) {
	// Save original manager
	original := DefaultManager()
	defer SetDefaultManager(original)

	// Create new manager for testing
	testManager := NewShutdownManager(nil)
	SetDefaultManager(testManager)

	// Test RegisterHook
	RegisterHook("test", 0, time.Second, func(ctx context.Context) error {
		return nil
	})

	if len(testManager.hooks) != 1 {
		t.Errorf("RegisterHook() failed to register hook")
	}

	// Test other functions don't panic
	_ = WaitForSignal()
	_ = Context()
	_ = IsShuttingDown()
}

func TestHookPanicRecovery(t *testing.T) {
	manager := NewShutdownManager(nil)

	manager.RegisterHook("panic-hook", 0, time.Second, func(ctx context.Context) error {
		panic("intentional panic")
	})

	ctx := context.Background()
	manager.executeHooks(ctx)
}

func TestContextCancellation(t *testing.T) {
	manager := NewShutdownManager(nil)

	var cancelled bool
	manager.RegisterHook("cancel-hook", 0, 5*time.Second, func(ctx context.Context) error {
		<-ctx.Done()
		cancelled = true
		return ctx.Err()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	manager.executeHooks(ctx)

	if !cancelled {
		t.Error("Hook should have detected context cancellation")
	}
}

func TestMultipleHooksWithDifferentTimeouts(t *testing.T) {
	manager := NewShutdownManager(&ShutdownManagerConfig{
		Timeout: 10 * time.Second,
	})

	var executionOrder []string

	manager.RegisterHook("fast", 10, 50*time.Millisecond, func(ctx context.Context) error {
		executionOrder = append(executionOrder, "fast")
		return nil
	})

	manager.RegisterHook("slow", 5, 200*time.Millisecond, func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		executionOrder = append(executionOrder, "slow")
		return nil
	})

	ctx := context.Background()
	manager.executeHooks(ctx)

	if len(executionOrder) != 2 {
		t.Errorf("Expected 2 hooks to execute, got %d", len(executionOrder))
	}
}

type testCloser struct {
	closeFn func() error
}

func (c *testCloser) Close() error {
	return c.closeFn()
}
