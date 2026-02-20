// Package health_test provides comprehensive tests for health check functionality
package health_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"biometrics-cli/pkg/health"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHealthChecker_NewHealthChecker tests creating a new health checker
func TestHealthChecker_NewHealthChecker(t *testing.T) {
	config := health.DefaultHealthCheckerConfig()
	config.Version = "2.0.0"

	hc := health.NewHealthChecker(config)
	require.NotNil(t, hc)

	assert.Equal(t, "2.0.0", hc.GetVersion())
	assert.False(t, hc.GetStartTime().IsZero())
}

// TestHealthChecker_DefaultConfig tests default configuration
func TestHealthChecker_DefaultConfig(t *testing.T) {
	config := health.DefaultHealthCheckerConfig()

	assert.NotEmpty(t, config.Version)
	assert.False(t, config.StartTime.IsZero())
	assert.True(t, config.EnableChecks)
}

// TestHealthChecker_HealthEndpoint tests the /health endpoint
func TestHealthChecker_HealthEndpoint(t *testing.T) {
	hc := health.NewHealthChecker(health.DefaultHealthCheckerConfig())
	handler := hc.HealthHandler()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var response health.HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, health.StatusHealthy, response.Status)
	assert.NotEmpty(t, response.Version)
	assert.NotEmpty(t, response.Uptime)
	assert.False(t, response.Timestamp.IsZero())
}

// TestHealthChecker_ReadyEndpoint tests the /ready endpoint
func TestHealthChecker_ReadyEndpoint(t *testing.T) {
	hc := health.NewHealthChecker(health.DefaultHealthCheckerConfig())
	handler := hc.ReadyHandler()

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response health.HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, health.StatusHealthy, response.Status)
}

// TestHealthChecker_LiveEndpoint tests the /live endpoint
func TestHealthChecker_LiveEndpoint(t *testing.T) {
	hc := health.NewHealthChecker(health.DefaultHealthCheckerConfig())
	handler := hc.LiveHandler()

	req := httptest.NewRequest(http.MethodGet, "/live", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response health.HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, health.StatusHealthy, response.Status)
	assert.NotEmpty(t, response.Version)
	assert.NotEmpty(t, response.Uptime)
}

// TestHealthChecker_CustomCheck tests registering custom health checks
func TestHealthChecker_CustomCheck(t *testing.T) {
	hc := health.NewHealthChecker(health.DefaultHealthCheckerConfig())

	hc.RegisterCheck("custom", func(ctx context.Context) health.CheckResult {
		return health.CheckResult{
			Status:  health.StatusHealthy,
			Message: "Custom check passed",
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	hc.HealthHandler().ServeHTTP(w, req)

	var response health.HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Contains(t, response.Checks, "custom")
	assert.Equal(t, health.StatusHealthy, response.Checks["custom"].Status)
}

// TestHealthChecker_CustomCheckUnhealthy tests unhealthy custom check
func TestHealthChecker_CustomCheckUnhealthy(t *testing.T) {
	hc := health.NewHealthChecker(health.DefaultHealthCheckerConfig())

	hc.RegisterCheck("failing", func(ctx context.Context) health.CheckResult {
		return health.CheckResult{
			Status:  health.StatusUnhealthy,
			Message: "Check failed",
			Error:   "simulated error",
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	hc.HealthHandler().ServeHTTP(w, req)

	var response health.HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, health.StatusUnhealthy, response.Status)
	assert.Contains(t, response.Checks, "failing")
}

// TestHealthChecker_SimpleDependencyChecker tests simple dependency checker
func TestHealthChecker_SimpleDependencyChecker(t *testing.T) {
	checker := health.NewSimpleDependencyChecker("test-db", func(ctx context.Context) (bool, error) {
		return true, nil
	}, true)

	assert.Equal(t, "test-db", checker.Name())
	assert.True(t, checker.IsRequired())

	ctx := context.Background()
	status := checker.CheckHealth(ctx)

	assert.True(t, status.Available)
	assert.True(t, status.Healthy)
	assert.GreaterOrEqual(t, status.Latency, time.Duration(0))
	assert.Empty(t, status.Error)
}

// TestHealthChecker_SimpleDependencyCheckerFailing tests failing dependency
func TestHealthChecker_SimpleDependencyCheckerFailing(t *testing.T) {
	checker := health.NewSimpleDependencyChecker("failing-db", func(ctx context.Context) (bool, error) {
		return false, assert.AnError
	}, true)

	ctx := context.Background()
	status := checker.CheckHealth(ctx)

	assert.True(t, status.Available)
	assert.False(t, status.Healthy)
	assert.NotEmpty(t, status.Error)
}

// TestHealthChecker_HTTPDependencyChecker tests HTTP dependency checker
func TestHealthChecker_HTTPDependencyChecker(t *testing.T) {
	// Create test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	checker := health.NewHTTPDependencyChecker("test-service", testServer.URL, 5*time.Second, true)

	assert.Equal(t, "test-service", checker.Name())
	assert.True(t, checker.IsRequired())

	ctx := context.Background()
	status := checker.CheckHealth(ctx)

	assert.True(t, status.Available)
	assert.True(t, status.Healthy)
	assert.Greater(t, status.Latency, time.Duration(0))
	assert.Empty(t, status.Error)
}

// TestHealthChecker_HTTPDependencyCheckerFailing tests failing HTTP dependency
func TestHealthChecker_HTTPDependencyCheckerFailing(t *testing.T) {
	checker := health.NewHTTPDependencyChecker("unreachable", "http://localhost:1", 100*time.Millisecond, true)

	ctx := context.Background()
	status := checker.CheckHealth(ctx)

	assert.False(t, status.Available)
	assert.False(t, status.Healthy)
	assert.NotEmpty(t, status.Error)
}

// TestHealthChecker_HTTPDependencyCheckerWrongStatus tests HTTP dependency with wrong status
func TestHealthChecker_HTTPDependencyCheckerWrongStatus(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer testServer.Close()

	checker := health.NewHTTPDependencyChecker("unhealthy-service", testServer.URL, 5*time.Second, true)

	ctx := context.Background()
	status := checker.CheckHealth(ctx)

	assert.True(t, status.Available)
	assert.False(t, status.Healthy)
	assert.NotEmpty(t, status.Error)
}

// TestHealthChecker_RequiredDependency tests required dependency impact on readiness
func TestHealthChecker_RequiredDependency(t *testing.T) {
	hc := health.NewHealthChecker(health.DefaultHealthCheckerConfig())

	// Add failing required dependency
	checker := health.NewSimpleDependencyChecker("required-db", func(ctx context.Context) (bool, error) {
		return false, assert.AnError
	}, true)

	hc.RegisterDependency(checker)

	handler := hc.ReadyHandler()
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response health.HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, health.StatusUnhealthy, response.Status)
}

// TestHealthChecker_OptionalDependency tests optional dependency doesn't affect readiness
func TestHealthChecker_OptionalDependency(t *testing.T) {
	hc := health.NewHealthChecker(health.DefaultHealthCheckerConfig())

	// Add failing optional dependency
	checker := health.NewSimpleDependencyChecker("optional-service", func(ctx context.Context) (bool, error) {
		return false, assert.AnError
	}, false)

	hc.RegisterDependency(checker)

	handler := hc.ReadyHandler()
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Should still be OK because dependency is optional
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestHealthChecker_RegisterUnregisterDependency tests dependency registration lifecycle
func TestHealthChecker_RegisterUnregisterDependency(t *testing.T) {
	hc := health.NewHealthChecker(health.DefaultHealthCheckerConfig())

	checker := health.NewSimpleDependencyChecker("temp-db", func(ctx context.Context) (bool, error) {
		return true, nil
	}, true)

	hc.RegisterDependency(checker)
	hc.UnregisterDependency("temp-db")

	// Should not affect health after unregister
	handler := hc.ReadyHandler()
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestHealthChecker_SetupHealthEndpoints tests setting up all endpoints
func TestHealthChecker_SetupHealthEndpoints(t *testing.T) {
	hc := health.NewHealthChecker(health.DefaultHealthCheckerConfig())
	mux := http.NewServeMux()

	health.SetupHealthEndpoints(mux, hc)

	// Test all endpoints
	endpoints := []string{"/health", "/ready", "/live"}
	for _, endpoint := range endpoints {
		req := httptest.NewRequest(http.MethodGet, endpoint, nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Endpoint %s should return OK", endpoint)
	}
}

// TestHealthChecker_WithHealthChecks tests health check middleware
func TestHealthChecker_WithHealthChecks(t *testing.T) {
	hc := health.NewHealthChecker(health.DefaultHealthCheckerConfig())

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := health.WithHealthChecks(nextHandler, hc)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	wrapped.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestHealthChecker_Getters tests getter methods
func TestHealthChecker_Getters(t *testing.T) {
	config := health.HealthCheckerConfig{
		Version:   "3.0.0",
		StartTime: time.Now().Add(-time.Hour),
	}

	hc := health.NewHealthChecker(config)

	assert.Equal(t, "3.0.0", hc.GetVersion())
	assert.Equal(t, config.StartTime, hc.GetStartTime())
	assert.GreaterOrEqual(t, hc.GetUptime(), time.Hour)
}

// TestHealthChecker_CheckResult tests check result structure
func TestHealthChecker_CheckResult(t *testing.T) {
	result := health.CheckResult{
		Status:  health.StatusHealthy,
		Message: "All good",
		Latency: 10 * time.Millisecond,
	}

	assert.Equal(t, health.StatusHealthy, result.Status)
	assert.Equal(t, "All good", result.Message)
	assert.Equal(t, 10*time.Millisecond, result.Latency)
}

// TestHealthChecker_DependencyStatus tests dependency status structure
func TestHealthChecker_DependencyStatus(t *testing.T) {
	status := health.DependencyStatus{
		Available: true,
		Healthy:   true,
		Latency:   5 * time.Millisecond,
	}

	assert.True(t, status.Available)
	assert.True(t, status.Healthy)
	assert.Equal(t, 5*time.Millisecond, status.Latency)
	assert.Empty(t, status.Error)
}

// TestHealthChecker_ContextCancellation tests health check with cancelled context
func TestHealthChecker_ContextCancellation(t *testing.T) {
	hc := health.NewHealthChecker(health.DefaultHealthCheckerConfig())

	// Add slow check
	hc.RegisterCheck("slow", func(ctx context.Context) health.CheckResult {
		select {
		case <-ctx.Done():
			return health.CheckResult{
				Status: health.StatusUnhealthy,
				Error:  "context cancelled",
			}
		case <-time.After(100 * time.Millisecond):
			return health.CheckResult{
				Status: health.StatusHealthy,
			}
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	response := hc.Check(ctx)
	// Should handle cancellation gracefully
	assert.NotEmpty(t, response.Timestamp)
}

// TestHealthChecker_MultipleChecks tests multiple concurrent checks
func TestHealthChecker_MultipleChecks(t *testing.T) {
	hc := health.NewHealthChecker(health.DefaultHealthCheckerConfig())

	// Register multiple checks
	for i := 0; i < 5; i++ {
		name := string(rune('a' + i))
		hc.RegisterCheck(name, func(ctx context.Context) health.CheckResult {
			return health.CheckResult{
				Status:  health.StatusHealthy,
				Message: "Check " + name + " passed",
			}
		})
	}

	response := hc.Check(context.Background())

	assert.Equal(t, health.StatusHealthy, response.Status)
	assert.GreaterOrEqual(t, len(response.Checks), 5)
}

// TestHealthChecker_ResponseJSON tests JSON response format
func TestHealthChecker_ResponseJSON(t *testing.T) {
	hc := health.NewHealthChecker(health.DefaultHealthCheckerConfig())
	handler := hc.HealthHandler()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify JSON is valid
	var rawJSON json.RawMessage
	err := json.NewDecoder(w.Body).Decode(&rawJSON)
	require.NoError(t, err)

	// Verify required fields exist
	var response map[string]interface{}
	err = json.Unmarshal(rawJSON, &response)
	require.NoError(t, err)

	assert.Contains(t, response, "status")
	assert.Contains(t, response, "timestamp")
	assert.Contains(t, response, "version")
	assert.Contains(t, response, "uptime")
}

// BenchmarkHealthChecker_HealthHandler benchmarks health endpoint
func BenchmarkHealthChecker_HealthHandler(b *testing.B) {
	hc := health.NewHealthChecker(health.DefaultHealthCheckerConfig())
	handler := hc.HealthHandler()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(w, req)
	}
}

// BenchmarkHealthChecker_Check benchmarks health check
func BenchmarkHealthChecker_Check(b *testing.B) {
	hc := health.NewHealthChecker(health.DefaultHealthCheckerConfig())
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hc.Check(ctx)
	}
}

// BenchmarkHealthChecker_ReadyHandler benchmarks ready endpoint
func BenchmarkHealthChecker_ReadyHandler(b *testing.B) {
	hc := health.NewHealthChecker(health.DefaultHealthCheckerConfig())
	handler := hc.ReadyHandler()

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(w, req)
	}
}
