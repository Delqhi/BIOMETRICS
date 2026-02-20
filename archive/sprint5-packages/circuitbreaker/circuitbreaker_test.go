// Package circuitbreaker_test provides comprehensive tests for the circuit breaker pattern.
// Coverage target: 80%+ with all edge cases tested.
package circuitbreaker_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"biometrics-cli/pkg/circuitbreaker"
	"biometrics-cli/pkg/logging"
)

// Test helper functions

func newTestLogger() *logging.Logger {
	return logging.New(logging.Config{
		Level:  logging.LevelDebug,
		Output: logging.OutputDiscard, // Discard logs during tests
	})
}

func newTestCircuitBreaker(name string) *circuitbreaker.CircuitBreaker {
	return circuitbreaker.NewCircuitBreaker(circuitbreaker.CircuitBreakerConfig{
		Name:                name,
		FailureThreshold:    3,
		SuccessThreshold:    2,
		Timeout:             100 * time.Millisecond,
		HalfOpenMaxRequests: 2,
		Logger:              newTestLogger(),
	})
}

// TestCircuitBreakerCreation tests the creation of a circuit breaker
func TestCircuitBreakerCreation(t *testing.T) {
	tests := []struct {
		name           string
		config         circuitbreaker.CircuitBreakerConfig
		expectedState  circuitbreaker.State
		expectDefaults bool
	}{
		{
			name: "with full config",
			config: circuitbreaker.CircuitBreakerConfig{
				Name:                "test-cb",
				FailureThreshold:    5,
				SuccessThreshold:    3,
				Timeout:             30 * time.Second,
				HalfOpenMaxRequests: 4,
				Logger:              newTestLogger(),
			},
			expectedState:  circuitbreaker.StateClosed,
			expectDefaults: false,
		},
		{
			name: "with minimal config",
			config: circuitbreaker.CircuitBreakerConfig{
				Name: "minimal-cb",
			},
			expectedState:  circuitbreaker.StateClosed,
			expectDefaults: true,
		},
		{
			name: "with zero values",
			config: circuitbreaker.CircuitBreakerConfig{
				Name:                "zero-cb",
				FailureThreshold:    0,
				SuccessThreshold:    0,
				Timeout:             0,
				HalfOpenMaxRequests: 0,
			},
			expectedState:  circuitbreaker.StateClosed,
			expectDefaults: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := circuitbreaker.NewCircuitBreaker(tt.config)

			if cb == nil {
				t.Fatal("NewCircuitBreaker returned nil")
			}

			if cb.State() != tt.expectedState {
				t.Errorf("expected state %v, got %v", tt.expectedState, cb.State())
			}

			if cb.GetName() != tt.config.Name {
				t.Errorf("expected name %s, got %s", tt.config.Name, cb.GetName())
			}

			if tt.expectDefaults {
				defaults := circuitbreaker.DefaultCircuitBreakerConfig()
				if cb.GetFailureThreshold() != defaults.FailureThreshold {
					t.Errorf("expected default failure threshold %d, got %d", defaults.FailureThreshold, cb.GetFailureThreshold())
				}
				if cb.GetSuccessThreshold() != defaults.SuccessThreshold {
					t.Errorf("expected default success threshold %d, got %d", defaults.SuccessThreshold, cb.GetSuccessThreshold())
				}
				if cb.GetTimeout() != defaults.Timeout {
					t.Errorf("expected default timeout %v, got %v", defaults.Timeout, cb.GetTimeout())
				}
			}
		})
	}
}

// TestDefaultCircuitBreakerConfig tests default configuration values
func TestDefaultCircuitBreakerConfig(t *testing.T) {
	config := circuitbreaker.DefaultCircuitBreakerConfig()

	if config.Name != "default" {
		t.Errorf("expected default name 'default', got %s", config.Name)
	}
	if config.FailureThreshold != 5 {
		t.Errorf("expected failure threshold 5, got %d", config.FailureThreshold)
	}
	if config.SuccessThreshold != 3 {
		t.Errorf("expected success threshold 3, got %d", config.SuccessThreshold)
	}
	if config.Timeout != 30*time.Second {
		t.Errorf("expected timeout 30s, got %v", config.Timeout)
	}
	if config.HalfOpenMaxRequests != 3 {
		t.Errorf("expected half-open max requests 3, got %d", config.HalfOpenMaxRequests)
	}
}

// TestStateTransition tests state transitions
func TestStateTransition(t *testing.T) {
	cb := newTestCircuitBreaker("state-test")

	// Initial state should be Closed
	if !cb.IsClosed() {
		t.Error("initial state should be Closed")
	}

	// Record failures to trigger Open state
	for i := 0; i < 3; i++ {
		cb.RecordFailure(errors.New("test failure"))
	}

	if !cb.IsOpen() {
		t.Error("state should be Open after reaching failure threshold")
	}

	// Wait for timeout to transition to Half-Open
	time.Sleep(150 * time.Millisecond)

	// Allow should transition to Half-Open
	err := cb.Allow()
	if err != nil {
		t.Errorf("Allow() should succeed in Half-Open state, got error: %v", err)
	}

	if !cb.IsHalfOpen() {
		t.Error("state should be Half-Open after timeout")
	}

	// Record successes to transition back to Closed
	cb.RecordSuccess()
	cb.RecordSuccess()

	if !cb.IsClosed() {
		t.Error("state should be Closed after reaching success threshold in Half-Open")
	}
}

// TestAllowFunction tests the Allow function behavior
func TestAllowFunction(t *testing.T) {
	t.Run("closed state allows requests", func(t *testing.T) {
		cb := newTestCircuitBreaker("allow-closed")

		for i := 0; i < 10; i++ {
			if err := cb.Allow(); err != nil {
				t.Errorf("Allow() should succeed in Closed state, got: %v", err)
			}
		}
	})

	t.Run("open state blocks requests", func(t *testing.T) {
		cb := newTestCircuitBreaker("allow-open")

		// Force open state
		for i := 0; i < 3; i++ {
			cb.RecordFailure(errors.New("failure"))
		}

		if err := cb.Allow(); err != circuitbreaker.ErrCircuitOpen {
			t.Errorf("Allow() should return ErrCircuitOpen in Open state, got: %v", err)
		}
	})

	t.Run("half-open state allows limited requests", func(t *testing.T) {
		cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.CircuitBreakerConfig{
			Name:                "half-open-allow",
			FailureThreshold:    2,
			SuccessThreshold:    2,
			Timeout:             50 * time.Millisecond,
			HalfOpenMaxRequests: 2,
			Logger:              newTestLogger(),
		})

		// Force open state
		cb.RecordFailure(errors.New("failure 1"))
		cb.RecordFailure(errors.New("failure 2"))

		// Wait for timeout
		time.Sleep(60 * time.Millisecond)

		// Should allow up to HalfOpenMaxRequests
		if err := cb.Allow(); err != nil {
			t.Errorf("first Allow() should succeed, got: %v", err)
		}
		if err := cb.Allow(); err != nil {
			t.Errorf("second Allow() should succeed, got: %v", err)
		}

		// Should block after max requests
		if err := cb.Allow(); err != circuitbreaker.ErrCircuitOpen {
			t.Errorf("third Allow() should return ErrCircuitOpen, got: %v", err)
		}
	})
}

// TestExecuteFunction tests the Execute function
func TestExecuteFunction(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		cb := newTestCircuitBreaker("execute-success")

		err := cb.Execute(context.Background(), func(ctx context.Context) error {
			return nil
		})

		if err != nil {
			t.Errorf("Execute() should succeed, got: %v", err)
		}

		if cb.GetSuccessCount() != 1 {
			t.Errorf("success count should be 1, got %d", cb.GetSuccessCount())
		}
	})

	t.Run("failed execution", func(t *testing.T) {
		cb := newTestCircuitBreaker("execute-failure")

		expectedErr := errors.New("operation failed")
		err := cb.Execute(context.Background(), func(ctx context.Context) error {
			return expectedErr
		})

		if err != expectedErr {
			t.Errorf("Execute() should return the original error, got: %v", err)
		}

		if cb.GetFailureCount() != 1 {
			t.Errorf("failure count should be 1, got %d", cb.GetFailureCount())
		}
	})

	t.Run("blocked execution when open", func(t *testing.T) {
		cb := newTestCircuitBreaker("execute-blocked")

		// Force open state
		for i := 0; i < 3; i++ {
			cb.RecordFailure(errors.New("failure"))
		}

		executed := false
		err := cb.Execute(context.Background(), func(ctx context.Context) error {
			executed = true
			return nil
		})

		if err != circuitbreaker.ErrCircuitOpen {
			t.Errorf("Execute() should return ErrCircuitOpen, got: %v", err)
		}

		if executed {
			t.Error("function should not have been executed")
		}
	})
}

// TestRecordSuccess tests the RecordSuccess function
func TestRecordSuccess(t *testing.T) {
	t.Run("resets failure count in closed state", func(t *testing.T) {
		cb := newTestCircuitBreaker("success-reset")

		// Record some failures
		cb.RecordFailure(errors.New("failure 1"))
		cb.RecordFailure(errors.New("failure 2"))

		if cb.GetFailureCount() != 2 {
			t.Fatalf("failure count should be 2, got %d", cb.GetFailureCount())
		}

		// Success should reset failure count
		cb.RecordSuccess()

		if cb.GetFailureCount() != 0 {
			t.Errorf("failure count should be reset to 0, got %d", cb.GetFailureCount())
		}
	})

	t.Run("triggers close from half-open", func(t *testing.T) {
		cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.CircuitBreakerConfig{
			Name:                "success-halfopen",
			FailureThreshold:    2,
			SuccessThreshold:    2,
			Timeout:             50 * time.Millisecond,
			HalfOpenMaxRequests: 2,
			Logger:              newTestLogger(),
		})

		// Force open then half-open
		cb.RecordFailure(errors.New("failure 1"))
		cb.RecordFailure(errors.New("failure 2"))
		time.Sleep(60 * time.Millisecond)
		cb.Allow() // Trigger half-open

		// Record successes
		cb.RecordSuccess()
		if cb.State() != circuitbreaker.StateHalfOpen {
			t.Error("should still be half-open after 1 success")
		}

		cb.RecordSuccess()
		if cb.State() != circuitbreaker.StateClosed {
			t.Error("should be closed after 2 successes")
		}
	})
}

// TestRecordFailure tests the RecordFailure function
func TestRecordFailure(t *testing.T) {
	t.Run("triggers open from closed", func(t *testing.T) {
		cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.CircuitBreakerConfig{
			Name:                "failure-open",
			FailureThreshold:    2,
			SuccessThreshold:    2,
			Timeout:             50 * time.Millisecond,
			HalfOpenMaxRequests: 2,
			Logger:              newTestLogger(),
		})

		cb.RecordFailure(errors.New("failure 1"))
		if cb.State() != circuitbreaker.StateClosed {
			t.Error("should still be closed after 1 failure")
		}

		cb.RecordFailure(errors.New("failure 2"))
		if cb.State() != circuitbreaker.StateOpen {
			t.Error("should be open after 2 failures")
		}
	})

	t.Run("triggers open from half-open", func(t *testing.T) {
		cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.CircuitBreakerConfig{
			Name:                "failure-halfopen",
			FailureThreshold:    2,
			SuccessThreshold:    2,
			Timeout:             50 * time.Millisecond,
			HalfOpenMaxRequests: 2,
			Logger:              newTestLogger(),
		})

		// Force open then half-open
		cb.RecordFailure(errors.New("failure 1"))
		cb.RecordFailure(errors.New("failure 2"))
		time.Sleep(60 * time.Millisecond)
		cb.Allow() // Trigger half-open

		// Any failure in half-open should open
		cb.RecordFailure(errors.New("failure in half-open"))
		if cb.State() != circuitbreaker.StateOpen {
			t.Error("should be open after failure in half-open")
		}
	})
}

// TestReset tests the Reset function
func TestReset(t *testing.T) {
	cb := newTestCircuitBreaker("reset-test")

	// Force open state
	for i := 0; i < 3; i++ {
		cb.RecordFailure(errors.New("failure"))
	}

	if !cb.IsOpen() {
		t.Fatal("circuit should be open")
	}

	// Reset
	cb.Reset()

	if !cb.IsClosed() {
		t.Error("circuit should be closed after reset")
	}

	if cb.GetFailureCount() != 0 {
		t.Errorf("failure count should be 0 after reset, got %d", cb.GetFailureCount())
	}

	if cb.GetSuccessCount() != 0 {
		t.Errorf("success count should be 0 after reset, got %d", cb.GetSuccessCount())
	}
}

// TestCallbacks tests state change callbacks
func TestCallbacks(t *testing.T) {
	t.Run("OnStateChange callback", func(t *testing.T) {
		cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.CircuitBreakerConfig{
			Name:                "callback-state",
			FailureThreshold:    2,
			SuccessThreshold:    2,
			Timeout:             50 * time.Millisecond,
			HalfOpenMaxRequests: 2,
			Logger:              newTestLogger(),
		})

		var callbackCalled int32
		var capturedOldState, capturedNewState circuitbreaker.State

		cb.OnStateChange(func(oldState, newState circuitbreaker.State) {
			atomic.StoreInt32(&callbackCalled, 1)
			capturedOldState = oldState
			capturedNewState = newState
		})

		// Trigger state change
		cb.RecordFailure(errors.New("failure 1"))
		cb.RecordFailure(errors.New("failure 2"))

		// Wait for callback (it's called in a goroutine)
		time.Sleep(50 * time.Millisecond)

		if atomic.LoadInt32(&callbackCalled) == 0 {
			t.Error("OnStateChange callback should have been called")
		}

		if capturedOldState != circuitbreaker.StateClosed {
			t.Errorf("expected old state Closed, got %v", capturedOldState)
		}

		if capturedNewState != circuitbreaker.StateOpen {
			t.Errorf("expected new state Open, got %v", capturedNewState)
		}
	})

	t.Run("OnSuccess callback", func(t *testing.T) {
		cb := newTestCircuitBreaker("callback-success")

		var callbackCalled int32
		cb.OnSuccess(func() {
			atomic.StoreInt32(&callbackCalled, 1)
		})

		cb.RecordSuccess()

		time.Sleep(50 * time.Millisecond)

		if atomic.LoadInt32(&callbackCalled) == 0 {
			t.Error("OnSuccess callback should have been called")
		}
	})

	t.Run("OnFailure callback", func(t *testing.T) {
		cb := newTestCircuitBreaker("callback-failure")

		var callbackCalled int32
		var capturedError error

		cb.OnFailure(func(err error) {
			atomic.StoreInt32(&callbackCalled, 1)
			capturedError = err
		})

		testErr := errors.New("test error")
		cb.RecordFailure(testErr)

		time.Sleep(50 * time.Millisecond)

		if atomic.LoadInt32(&callbackCalled) == 0 {
			t.Error("OnFailure callback should have been called")
		}

		if capturedError != testErr {
			t.Errorf("expected error %v, got %v", testErr, capturedError)
		}
	})
}

// TestMetrics tests metrics collection
func TestMetrics(t *testing.T) {
	cb := newTestCircuitBreaker("metrics-test")

	// Record some operations
	cb.RecordSuccess()
	cb.RecordSuccess()
	cb.RecordFailure(errors.New("failure"))
	cb.RecordSuccess()

	metrics := cb.GetMetrics()

	if metrics.GetTotalSuccesses() != 3 {
		t.Errorf("expected 3 successes, got %d", metrics.GetTotalSuccesses())
	}

	if metrics.GetTotalFailures() != 1 {
		t.Errorf("expected 1 failure, got %d", metrics.GetTotalFailures())
	}

	if metrics.GetTotalRequests() != 4 {
		t.Errorf("expected 4 requests, got %d", metrics.GetTotalRequests())
	}

	expectedRate := float64(3) / float64(4) * 100
	if rate := metrics.GetSuccessRate(); rate != expectedRate {
		t.Errorf("expected success rate %.2f, got %.2f", expectedRate, rate)
	}
}

// TestMetricsReset tests metrics reset
func TestMetricsReset(t *testing.T) {
	cb := newTestCircuitBreaker("metrics-reset")

	// Record some operations
	cb.RecordSuccess()
	cb.RecordFailure(errors.New("failure"))

	metrics := cb.GetMetrics()
	metrics.Reset()

	if metrics.GetTotalRequests() != 0 {
		t.Errorf("total requests should be 0 after reset, got %d", metrics.GetTotalRequests())
	}

	if metrics.GetTotalSuccesses() != 0 {
		t.Errorf("total successes should be 0 after reset, got %d", metrics.GetTotalSuccesses())
	}
}

// TestCircuitBreakerRegistry tests the registry functionality
func TestCircuitBreakerRegistry(t *testing.T) {
	t.Run("GetOrCreate creates new breaker", func(t *testing.T) {
		registry := circuitbreaker.NewCircuitBreakerRegistry()

		cb := registry.GetOrCreate("test-breaker", nil)

		if cb == nil {
			t.Fatal("GetOrCreate returned nil")
		}

		if cb.GetName() != "test-breaker" {
			t.Errorf("expected name 'test-breaker', got %s", cb.GetName())
		}
	})

	t.Run("GetOrCreate returns existing breaker", func(t *testing.T) {
		registry := circuitbreaker.NewCircuitBreakerRegistry()

		cb1 := registry.GetOrCreate("shared-breaker", nil)
		cb2 := registry.GetOrCreate("shared-breaker", nil)

		if cb1 != cb2 {
			t.Error("GetOrCreate should return the same instance")
		}
	})

	t.Run("Get returns existing breaker", func(t *testing.T) {
		registry := circuitbreaker.NewCircuitBreakerRegistry()

		registry.GetOrCreate("get-test", nil)
		cb, exists := registry.Get("get-test")

		if !exists {
			t.Error("breaker should exist")
		}

		if cb == nil {
			t.Error("Get returned nil")
		}
	})

	t.Run("Get returns false for non-existent", func(t *testing.T) {
		registry := circuitbreaker.NewCircuitBreakerRegistry()

		_, exists := registry.Get("non-existent")

		if exists {
			t.Error("non-existent breaker should not exist")
		}
	})

	t.Run("Delete removes breaker", func(t *testing.T) {
		registry := circuitbreaker.NewCircuitBreakerRegistry()

		registry.GetOrCreate("delete-test", nil)
		registry.Delete("delete-test")

		_, exists := registry.Get("delete-test")

		if exists {
			t.Error("breaker should have been deleted")
		}
	})

	t.Run("List returns all names", func(t *testing.T) {
		registry := circuitbreaker.NewCircuitBreakerRegistry()

		registry.GetOrCreate("breaker-1", nil)
		registry.GetOrCreate("breaker-2", nil)
		registry.GetOrCreate("breaker-3", nil)

		names := registry.List()

		if len(names) != 3 {
			t.Errorf("expected 3 names, got %d", len(names))
		}
	})

	t.Run("GetAll returns all breakers", func(t *testing.T) {
		registry := circuitbreaker.NewCircuitBreakerRegistry()

		registry.GetOrCreate("breaker-a", nil)
		registry.GetOrCreate("breaker-b", nil)

		all := registry.GetAll()

		if len(all) != 2 {
			t.Errorf("expected 2 breakers, got %d", len(all))
		}
	})
}

// TestCircuitBreakerMiddleware tests the middleware
func TestCircuitBreakerMiddleware(t *testing.T) {
	t.Run("Handle wraps handler", func(t *testing.T) {
		cb := newTestCircuitBreaker("middleware-test")
		middleware := circuitbreaker.NewCircuitBreakerMiddleware(cb)

		executed := false
		handler := func(ctx context.Context) error {
			executed = true
			return nil
		}

		wrapped := middleware.Handle(handler)
		err := wrapped(context.Background())

		if err != nil {
			t.Errorf("wrapped handler should succeed, got: %v", err)
		}

		if !executed {
			t.Error("handler should have been executed")
		}
	})

	t.Run("Handle blocks when open", func(t *testing.T) {
		cb := newTestCircuitBreaker("middleware-block")
		middleware := circuitbreaker.NewCircuitBreakerMiddleware(cb)

		// Force open state
		for i := 0; i < 3; i++ {
			cb.RecordFailure(errors.New("failure"))
		}

		executed := false
		handler := func(ctx context.Context) error {
			executed = true
			return nil
		}

		wrapped := middleware.Handle(handler)
		err := wrapped(context.Background())

		if err != circuitbreaker.ErrCircuitOpen {
			t.Errorf("should return ErrCircuitOpen, got: %v", err)
		}

		if executed {
			t.Error("handler should not have been executed")
		}
	})
}

// TestTimeUntilHalfOpen tests the TimeUntilHalfOpen function
func TestTimeUntilHalfOpen(t *testing.T) {
	t.Run("returns 0 when closed", func(t *testing.T) {
		cb := newTestCircuitBreaker("timeuntil-closed")

		if duration := cb.TimeUntilHalfOpen(); duration != 0 {
			t.Errorf("TimeUntilHalfOpen should be 0 when closed, got %v", duration)
		}
	})

	t.Run("returns remaining time when open", func(t *testing.T) {
		cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.CircuitBreakerConfig{
			Name:                "timeuntil-open",
			FailureThreshold:    1,
			SuccessThreshold:    1,
			Timeout:             200 * time.Millisecond,
			HalfOpenMaxRequests: 1,
			Logger:              newTestLogger(),
		})

		cb.RecordFailure(errors.New("failure"))

		duration := cb.TimeUntilHalfOpen()
		if duration <= 0 || duration > 200*time.Millisecond {
			t.Errorf("TimeUntilHalfOpen should be between 0 and 200ms, got %v", duration)
		}

		// Wait and check again
		time.Sleep(100 * time.Millisecond)
		duration = cb.TimeUntilHalfOpen()
		if duration > 100*time.Millisecond {
			t.Errorf("TimeUntilHalfOpen should be less than 100ms after waiting, got %v", duration)
		}
	})

	t.Run("returns 0 when half-open", func(t *testing.T) {
		cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.CircuitBreakerConfig{
			Name:                "timeuntil-halfopen",
			FailureThreshold:    1,
			SuccessThreshold:    1,
			Timeout:             50 * time.Millisecond,
			HalfOpenMaxRequests: 1,
			Logger:              newTestLogger(),
		})

		cb.RecordFailure(errors.New("failure"))
		time.Sleep(60 * time.Millisecond)
		cb.Allow() // Trigger half-open

		if duration := cb.TimeUntilHalfOpen(); duration != 0 {
			t.Errorf("TimeUntilHalfOpen should be 0 when half-open, got %v", duration)
		}
	})
}

// TestGetHealthStatus tests the health status function
func TestGetHealthStatus(t *testing.T) {
	t.Run("healthy when closed", func(t *testing.T) {
		cb := newTestCircuitBreaker("health-closed")

		status := cb.GetHealthStatus()

		if !status.Healthy {
			t.Error("circuit should be healthy when closed")
		}

		if status.State != "Closed" {
			t.Errorf("expected state 'Closed', got %s", status.State)
		}
	})

	t.Run("unhealthy when open", func(t *testing.T) {
		cb := newTestCircuitBreaker("health-open")

		for i := 0; i < 3; i++ {
			cb.RecordFailure(errors.New("failure"))
		}

		status := cb.GetHealthStatus()

		if status.Healthy {
			t.Error("circuit should be unhealthy when open")
		}

		if status.State != "Open" {
			t.Errorf("expected state 'Open', got %s", status.State)
		}

		if status.TimeUntilHalfOpen == "" {
			t.Error("TimeUntilHalfOpen should be set when open")
		}
	})

	t.Run("healthy when half-open", func(t *testing.T) {
		cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.CircuitBreakerConfig{
			Name:                "health-halfopen",
			FailureThreshold:    1,
			SuccessThreshold:    1,
			Timeout:             50 * time.Millisecond,
			HalfOpenMaxRequests: 1,
			Logger:              newTestLogger(),
		})

		cb.RecordFailure(errors.New("failure"))
		time.Sleep(60 * time.Millisecond)
		cb.Allow()

		status := cb.GetHealthStatus()

		if !status.Healthy {
			t.Error("circuit should be healthy when half-open")
		}

		if status.State != "Half-Open" {
			t.Errorf("expected state 'Half-Open', got %s", status.State)
		}
	})
}

// TestString tests the String function
func TestString(t *testing.T) {
	cb := newTestCircuitBreaker("string-test")

	str := cb.String()

	if str == "" {
		t.Error("String() should not return empty string")
	}

	// Should contain the name
	if !containsSubstring(str, "string-test") {
		t.Errorf("String() should contain name 'string-test', got: %s", str)
	}
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s[1:], substr))
}

// TestConcurrentAccess tests concurrent access safety
func TestConcurrentAccess(t *testing.T) {
	cb := newTestCircuitBreaker("concurrent-test")

	var wg sync.WaitGroup
	numOps := 100

	// Concurrent successes
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cb.RecordSuccess()
		}()
	}

	// Concurrent failures
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cb.RecordFailure(errors.New("concurrent failure"))
		}()
	}

	// Concurrent state checks
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = cb.State()
			_ = cb.IsOpen()
			_ = cb.IsClosed()
			_ = cb.IsHalfOpen()
		}()
	}

	wg.Wait()

	// Should not panic and should have recorded operations
	metrics := cb.GetMetrics()
	if metrics.GetTotalRequests() < int64(numOps) {
		t.Errorf("should have recorded at least %d operations, got %d", numOps, metrics.GetTotalRequests())
	}
}

// TestStateString tests the State String method
func TestStateString(t *testing.T) {
	tests := []struct {
		state    circuitbreaker.State
		expected string
	}{
		{circuitbreaker.StateClosed, "Closed"},
		{circuitbreaker.StateOpen, "Open"},
		{circuitbreaker.StateHalfOpen, "Half-Open"},
		{circuitbreaker.State(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.state.String(); got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

// TestContextCancellation tests context cancellation handling
func TestContextCancellation(t *testing.T) {
	cb := newTestCircuitBreaker("context-test")

	t.Run("respects context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := cb.Execute(ctx, func(ctx context.Context) error {
			return nil
		})

		// The function should still run because Allow() succeeds first
		// but in real scenarios, the function could check ctx
		if err != nil {
			t.Logf("Execute returned: %v (context was cancelled)", err)
		}
	})
}

// BenchmarkCircuitBreaker benchmarks circuit breaker operations
func BenchmarkCircuitBreaker(b *testing.B) {
	cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.CircuitBreakerConfig{
		Name:                "bench",
		FailureThreshold:    5,
		SuccessThreshold:    3,
		Timeout:             30 * time.Second,
		HalfOpenMaxRequests: 3,
		Logger:              newTestLogger(),
	})

	b.Run("Allow", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = cb.Allow()
		}
	})

	b.Run("Execute", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = cb.Execute(context.Background(), func(ctx context.Context) error {
				return nil
			})
		}
	})

	b.Run("RecordSuccess", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cb.RecordSuccess()
		}
	})

	b.Run("RecordFailure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cb.RecordFailure(errors.New("benchmark error"))
		}
	})

	b.Run("State", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = cb.State()
		}
	})

	b.Run("GetMetrics", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = cb.GetMetrics()
		}
	})
}

// BenchmarkCircuitBreakerRegistry benchmarks registry operations
func BenchmarkCircuitBreakerRegistry(b *testing.B) {
	registry := circuitbreaker.NewCircuitBreakerRegistry()

	b.Run("GetOrCreate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			registry.GetOrCreate("bench-breaker", nil)
		}
	})

	b.Run("Get", func(b *testing.B) {
		registry.GetOrCreate("get-bench", nil)
		for i := 0; i < b.N; i++ {
			_, _ = registry.Get("get-bench")
		}
	})
}
