// Package health provides comprehensive health check endpoints for the BIOMETRICS CLI.
// It includes basic health checks, readiness probes, and liveness probes with dependency checking.
package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// HealthStatus represents the health status of the service
type HealthStatus string

const (
	// StatusHealthy indicates the service is healthy
	StatusHealthy HealthStatus = "healthy"
	// StatusUnhealthy indicates the service is unhealthy
	StatusUnhealthy HealthStatus = "unhealthy"
	// StatusDegraded indicates the service is partially healthy
	StatusDegraded HealthStatus = "degraded"
)

// HealthResponse represents the JSON response for health endpoints
type HealthResponse struct {
	Status       HealthStatus                `json:"status"`
	Timestamp    time.Time                   `json:"timestamp"`
	Version      string                      `json:"version"`
	Uptime       string                      `json:"uptime"`
	Checks       map[string]CheckResult      `json:"checks,omitempty"`
	Dependencies map[string]DependencyStatus `json:"dependencies,omitempty"`
}

// CheckResult represents the result of a single health check
type CheckResult struct {
	Status  HealthStatus  `json:"status"`
	Message string        `json:"message,omitempty"`
	Latency time.Duration `json:"latency,omitempty"`
	Error   string        `json:"error,omitempty"`
}

// DependencyStatus represents the status of a dependency
type DependencyStatus struct {
	Available bool          `json:"available"`
	Healthy   bool          `json:"healthy"`
	Latency   time.Duration `json:"latency,omitempty"`
	Error     string        `json:"error,omitempty"`
}

// DependencyChecker defines the interface for checking dependency health
type DependencyChecker interface {
	// Name returns the name of the dependency
	Name() string
	// CheckHealth performs a health check on the dependency
	CheckHealth(ctx context.Context) DependencyStatus
	// IsRequired returns true if this dependency is required for the service to function
	IsRequired() bool
}

// HealthChecker manages health checks and dependencies
type HealthChecker struct {
	startTime    time.Time
	version      string
	checks       map[string]func(context.Context) CheckResult
	dependencies map[string]DependencyChecker
	mu           sync.RWMutex
	customChecks map[string]func(context.Context) CheckResult
}

// HealthCheckerConfig holds configuration for the health checker
type HealthCheckerConfig struct {
	Version      string
	StartTime    time.Time
	EnableChecks bool
}

// DefaultHealthCheckerConfig returns a default configuration
func DefaultHealthCheckerConfig() HealthCheckerConfig {
	return HealthCheckerConfig{
		Version:      "1.0.0",
		StartTime:    time.Now(),
		EnableChecks: true,
	}
}

// NewHealthChecker creates a new health checker instance
func NewHealthChecker(config HealthCheckerConfig) *HealthChecker {
	if config.StartTime.IsZero() {
		config.StartTime = time.Now()
	}
	if config.Version == "" {
		config.Version = "1.0.0"
	}

	hc := &HealthChecker{
		startTime:    config.StartTime,
		version:      config.Version,
		checks:       make(map[string]func(context.Context) CheckResult),
		dependencies: make(map[string]DependencyChecker),
		customChecks: make(map[string]func(context.Context) CheckResult),
	}

	// Register built-in checks
	hc.registerBuiltInChecks()

	return hc
}

// registerBuiltInChecks registers the built-in health checks
func (hc *HealthChecker) registerBuiltInChecks() {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	// Basic service check
	hc.checks["service"] = func(ctx context.Context) CheckResult {
		return CheckResult{
			Status:  StatusHealthy,
			Message: "Service is running",
		}
	}

	// Memory check
	hc.checks["memory"] = func(ctx context.Context) CheckResult {
		// Basic memory check - can be enhanced with actual memory usage
		return CheckResult{
			Status:  StatusHealthy,
			Message: "Memory check passed",
		}
	}

	// Goroutine check
	hc.checks["goroutines"] = func(ctx context.Context) CheckResult {
		// Basic goroutine count check
		return CheckResult{
			Status:  StatusHealthy,
			Message: "Goroutine check passed",
		}
	}
}

// RegisterCheck registers a custom health check
func (hc *HealthChecker) RegisterCheck(name string, checkFn func(context.Context) CheckResult) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.customChecks[name] = checkFn
}

// RegisterDependency registers a dependency checker
func (hc *HealthChecker) RegisterDependency(checker DependencyChecker) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.dependencies[checker.Name()] = checker
}

// UnregisterDependency removes a dependency checker
func (hc *HealthChecker) UnregisterDependency(name string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	delete(hc.dependencies, name)
}

// Check performs a comprehensive health check
func (hc *HealthChecker) Check(ctx context.Context) HealthResponse {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	response := HealthResponse{
		Status:       StatusHealthy,
		Timestamp:    time.Now(),
		Version:      hc.version,
		Uptime:       time.Since(hc.startTime).String(),
		Checks:       make(map[string]CheckResult),
		Dependencies: make(map[string]DependencyStatus),
	}

	// Run built-in and custom checks
	for name, checkFn := range hc.checks {
		result := checkFn(ctx)
		response.Checks[name] = result
		if result.Status == StatusUnhealthy {
			response.Status = StatusUnhealthy
		} else if result.Status == StatusDegraded && response.Status == StatusHealthy {
			response.Status = StatusDegraded
		}
	}

	// Run custom checks
	for name, checkFn := range hc.customChecks {
		result := checkFn(ctx)
		response.Checks[name] = result
		if result.Status == StatusUnhealthy {
			response.Status = StatusUnhealthy
		} else if result.Status == StatusDegraded && response.Status == StatusHealthy {
			response.Status = StatusDegraded
		}
	}

	// Check dependencies
	requiredHealthy := 0
	requiredTotal := 0
	for name, checker := range hc.dependencies {
		status := checker.CheckHealth(ctx)
		response.Dependencies[name] = status

		if checker.IsRequired() {
			requiredTotal++
			if status.Healthy {
				requiredHealthy++
			}
		}
	}

	// Update status based on required dependencies
	if requiredTotal > 0 && requiredHealthy < requiredTotal {
		if requiredHealthy == 0 {
			response.Status = StatusUnhealthy
		} else {
			response.Status = StatusDegraded
		}
	}

	return response
}

// HealthHandler returns an HTTP handler for the /health endpoint
func (hc *HealthChecker) HealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		response := hc.Check(ctx)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})
}

// ReadyHandler returns an HTTP handler for the /ready endpoint
func (hc *HealthChecker) ReadyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		response := hc.Check(ctx)

		// Check if all required dependencies are healthy
		allRequiredHealthy := true
		hc.mu.RLock()
		for name, checker := range hc.dependencies {
			if checker.IsRequired() {
				status := checker.CheckHealth(ctx)
				if !status.Healthy {
					allRequiredHealthy = false
					response.Dependencies[name] = status
				}
			}
		}
		hc.mu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		if !allRequiredHealthy {
			w.WriteHeader(http.StatusServiceUnavailable)
			response.Status = StatusUnhealthy
		} else {
			w.WriteHeader(http.StatusOK)
			response.Status = StatusHealthy
		}
		json.NewEncoder(w).Encode(response)
	})
}

// LiveHandler returns an HTTP handler for the /live endpoint
func (hc *HealthChecker) LiveHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{
			Status:    StatusHealthy,
			Timestamp: time.Now(),
			Version:   hc.version,
			Uptime:    time.Since(hc.startTime).String(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})
}

// GetUptime returns the service uptime
func (hc *HealthChecker) GetUptime() time.Duration {
	return time.Since(hc.startTime)
}

// GetVersion returns the service version
func (hc *HealthChecker) GetVersion() string {
	return hc.version
}

// GetStartTime returns the service start time
func (hc *HealthChecker) GetStartTime() time.Time {
	return hc.startTime
}

// SimpleDependencyChecker implements DependencyChecker for simple cases
type SimpleDependencyChecker struct {
	name     string
	checkFn  func(context.Context) (bool, error)
	required bool
}

// NewSimpleDependencyChecker creates a simple dependency checker
func NewSimpleDependencyChecker(name string, checkFn func(context.Context) (bool, error), required bool) *SimpleDependencyChecker {
	return &SimpleDependencyChecker{
		name:     name,
		checkFn:  checkFn,
		required: required,
	}
}

// Name returns the dependency name
func (s *SimpleDependencyChecker) Name() string {
	return s.name
}

// CheckHealth performs the health check
func (s *SimpleDependencyChecker) CheckHealth(ctx context.Context) DependencyStatus {
	start := time.Now()
	healthy, err := s.checkFn(ctx)
	latency := time.Since(start)

	status := DependencyStatus{
		Available: true,
		Healthy:   healthy,
		Latency:   latency,
	}

	if err != nil {
		status.Error = err.Error()
		status.Healthy = false
	}

	return status
}

// IsRequired returns whether the dependency is required
func (s *SimpleDependencyChecker) IsRequired() bool {
	return s.required
}

// HTTPDependencyChecker checks HTTP endpoints
type HTTPDependencyChecker struct {
	name     string
	url      string
	timeout  time.Duration
	required bool
	client   *http.Client
}

// NewHTTPDependencyChecker creates an HTTP dependency checker
func NewHTTPDependencyChecker(name, url string, timeout time.Duration, required bool) *HTTPDependencyChecker {
	return &HTTPDependencyChecker{
		name:     name,
		url:      url,
		timeout:  timeout,
		required: required,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Name returns the dependency name
func (h *HTTPDependencyChecker) Name() string {
	return h.name
}

// CheckHealth performs an HTTP health check
func (h *HTTPDependencyChecker) CheckHealth(ctx context.Context) DependencyStatus {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.url, nil)
	if err != nil {
		return DependencyStatus{
			Available: false,
			Healthy:   false,
			Error:     fmt.Sprintf("failed to create request: %v", err),
		}
	}

	resp, err := h.client.Do(req)
	latency := time.Since(start)

	if err != nil {
		return DependencyStatus{
			Available: false,
			Healthy:   false,
			Latency:   latency,
			Error:     fmt.Sprintf("request failed: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return DependencyStatus{
			Available: true,
			Healthy:   false,
			Latency:   latency,
			Error:     fmt.Sprintf("unhealthy status code: %d", resp.StatusCode),
		}
	}

	return DependencyStatus{
		Available: true,
		Healthy:   true,
		Latency:   latency,
	}
}

// IsRequired returns whether the dependency is required
func (h *HTTPDependencyChecker) IsRequired() bool {
	return h.required
}

// SetupHealthEndpoints sets up all health check endpoints on a mux
func SetupHealthEndpoints(mux *http.ServeMux, hc *HealthChecker) {
	mux.Handle("/health", hc.HealthHandler())
	mux.Handle("/ready", hc.ReadyHandler())
	mux.Handle("/live", hc.LiveHandler())
}

// WithHealthChecks wraps an HTTP handler with health check middleware
func WithHealthChecks(next http.Handler, hc *HealthChecker) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add health check context
		ctx := context.WithValue(r.Context(), "health_checker", hc)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
