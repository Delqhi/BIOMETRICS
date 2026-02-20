// Package circuitbreaker implements the circuit breaker pattern for fault tolerance.
// It provides three states (Closed, Open, Half-Open) with configurable failure thresholds,
// recovery timeouts, metrics collection, and state change callbacks.
package circuitbreaker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"biometrics-cli/pkg/logging"
)

// State represents the current state of the circuit breaker
type State int

const (
	// StateClosed - Circuit is closed, requests flow normally
	StateClosed State = iota
	// StateOpen - Circuit is open, requests are blocked
	StateOpen
	// StateHalfOpen - Circuit is half-open, limited requests allowed to test recovery
	StateHalfOpen
)

// String returns the string representation of the state
func (s State) String() string {
	switch s {
	case StateClosed:
		return "Closed"
	case StateOpen:
		return "Open"
	case StateHalfOpen:
		return "Half-Open"
	default:
		return "Unknown"
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	mu sync.RWMutex

	// Configuration
	name                string
	failureThreshold    int
	successThreshold    int
	timeout             time.Duration
	halfOpenMaxRequests int

	// State
	state            State
	failures         int
	successes        int
	lastFailureTime  time.Time
	lastStateChange  time.Time
	halfOpenRequests int

	// Metrics
	metrics *CircuitBreakerMetrics

	// Callbacks
	onStateChange []func(oldState, newState State)
	onSuccess     []func()
	onFailure     []func(error)

	// Logger
	logger *logging.Logger
}

// CircuitBreakerConfig holds configuration for the circuit breaker
type CircuitBreakerConfig struct {
	// Name is the identifier for this circuit breaker
	Name string
	// FailureThreshold is the number of failures before opening the circuit
	FailureThreshold int
	// SuccessThreshold is the number of successes in half-open state before closing
	SuccessThreshold int
	// Timeout is the duration the circuit stays open before transitioning to half-open
	Timeout time.Duration
	// HalfOpenMaxRequests is the maximum number of requests allowed in half-open state
	HalfOpenMaxRequests int
	// Logger is the logger to use for circuit breaker events
	Logger *logging.Logger
}

// DefaultCircuitBreakerConfig returns a reasonable default configuration
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		Name:                "default",
		FailureThreshold:    5,
		SuccessThreshold:    3,
		Timeout:             30 * time.Second,
		HalfOpenMaxRequests: 3,
	}
}

// CircuitBreakerMetrics holds metrics for the circuit breaker
type CircuitBreakerMetrics struct {
	mu sync.RWMutex

	// Total counts
	totalRequests     int64
	totalSuccesses    int64
	totalFailures     int64
	totalTimeouts     int64
	totalStateChanges int64

	// State change counts
	closedToOpenCount     int64
	openToHalfOpenCount   int64
	halfOpenToClosedCount int64
	halfOpenToOpenCount   int64

	// Timing
	lastRequestTime time.Time
	lastSuccessTime time.Time
	lastFailureTime time.Time
	averageLatency  time.Duration
	totalLatency    time.Duration
}

// NewCircuitBreaker creates a new circuit breaker instance
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	if config.FailureThreshold <= 0 {
		config.FailureThreshold = DefaultCircuitBreakerConfig().FailureThreshold
	}
	if config.SuccessThreshold <= 0 {
		config.SuccessThreshold = DefaultCircuitBreakerConfig().SuccessThreshold
	}
	if config.Timeout <= 0 {
		config.Timeout = DefaultCircuitBreakerConfig().Timeout
	}
	if config.HalfOpenMaxRequests <= 0 {
		config.HalfOpenMaxRequests = DefaultCircuitBreakerConfig().HalfOpenMaxRequests
	}
	if config.Logger == nil {
		config.Logger = logging.Default()
	}

	cb := &CircuitBreaker{
		name:                config.Name,
		failureThreshold:    config.FailureThreshold,
		successThreshold:    config.SuccessThreshold,
		timeout:             config.Timeout,
		halfOpenMaxRequests: config.HalfOpenMaxRequests,
		state:               StateClosed,
		logger:              config.Logger,
		metrics:             &CircuitBreakerMetrics{},
		onStateChange:       make([]func(State, State), 0),
		onSuccess:           make([]func(), 0),
		onFailure:           make([]func(error), 0),
		lastStateChange:     time.Now(),
	}

	cb.logger.Info("Circuit breaker created",
		logging.String("name", config.Name),
		logging.Int("failure_threshold", config.FailureThreshold),
		logging.Int("success_threshold", config.SuccessThreshold),
		logging.Duration("timeout", config.Timeout),
		logging.Int("half_open_max_requests", config.HalfOpenMaxRequests),
	)

	return cb
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	startTime := time.Now()

	if err := cb.Allow(); err != nil {
		cb.metrics.recordRejected()
		return err
	}

	err := fn(ctx)
	latency := time.Since(startTime)
	cb.metrics.recordRequest(latency)

	if err != nil {
		cb.RecordFailure(err)
	} else {
		cb.RecordSuccess()
	}

	return err
}

// Allow checks if a request is allowed through the circuit breaker
func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.metrics.recordRequest(time.Duration(0))

	switch cb.state {
	case StateClosed:
		return nil

	case StateOpen:
		// Check if timeout has elapsed
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.transitionTo(StateHalfOpen)
			return nil
		}
		return ErrCircuitOpen

	case StateHalfOpen:
		// Allow limited requests in half-open state
		if cb.halfOpenRequests < cb.halfOpenMaxRequests {
			cb.halfOpenRequests++
			return nil
		}
		return ErrCircuitOpen

	default:
		return ErrCircuitOpen
	}
}

// RecordSuccess records a successful operation
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.metrics.recordSuccess()
	cb.successes++

	if cb.state == StateHalfOpen {
		if cb.successes >= cb.successThreshold {
			cb.transitionTo(StateClosed)
		}
	} else if cb.state == StateClosed {
		// Reset failure count on success
		cb.failures = 0
	}

	cb.logger.Debug("Circuit breaker success recorded",
		logging.String("name", cb.name),
		logging.String("state", cb.state.String()),
		logging.Int("successes", cb.successes),
		logging.Int("failures", cb.failures),
	)

	// Call success callbacks
	for _, callback := range cb.onSuccess {
		go callback()
	}
}

// RecordFailure records a failed operation
func (cb *CircuitBreaker) RecordFailure(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.metrics.recordFailure()
	cb.failures++
	cb.lastFailureTime = time.Now()

	if cb.state == StateHalfOpen {
		// Any failure in half-open state opens the circuit
		cb.transitionTo(StateOpen)
	} else if cb.state == StateClosed {
		if cb.failures >= cb.failureThreshold {
			cb.transitionTo(StateOpen)
		}
	}

	cb.logger.Warn("Circuit breaker failure recorded",
		logging.String("name", cb.name),
		logging.String("state", cb.state.String()),
		logging.Int("successes", cb.successes),
		logging.Int("failures", cb.failures),
		logging.Err(err),
	)

	// Call failure callbacks
	for _, callback := range cb.onFailure {
		go callback(err)
	}
}

// Reset manually resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state != StateClosed {
		cb.logger.Info("Circuit breaker manually reset",
			logging.String("name", cb.name),
			logging.String("from_state", cb.state.String()),
		)
		cb.transitionTo(StateClosed)
	}

	cb.failures = 0
	cb.successes = 0
	cb.halfOpenRequests = 0
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// StateString returns the current state as a string
func (cb *CircuitBreaker) StateString() string {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state.String()
}

// IsOpen returns true if the circuit is open
func (cb *CircuitBreaker) IsOpen() bool {
	return cb.State() == StateOpen
}

// IsClosed returns true if the circuit is closed
func (cb *CircuitBreaker) IsClosed() bool {
	return cb.State() == StateClosed
}

// IsHalfOpen returns true if the circuit is half-open
func (cb *CircuitBreaker) IsHalfOpen() bool {
	return cb.State() == StateHalfOpen
}

// GetMetrics returns a copy of the current metrics
func (cb *CircuitBreaker) GetMetrics() CircuitBreakerMetrics {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	cb.metrics.mu.RLock()
	defer cb.metrics.mu.RUnlock()

	return *cb.metrics
}

// OnStateChange registers a callback for state changes
func (cb *CircuitBreaker) OnStateChange(callback func(oldState, newState State)) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.onStateChange = append(cb.onStateChange, callback)
}

// OnSuccess registers a callback for successful operations
func (cb *CircuitBreaker) OnSuccess(callback func()) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.onSuccess = append(cb.onSuccess, callback)
}

// OnFailure registers a callback for failed operations
func (cb *CircuitBreaker) OnFailure(callback func(error)) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.onFailure = append(cb.onFailure, callback)
}

// transitionTo transitions to a new state
func (cb *CircuitBreaker) transitionTo(newState State) {
	oldState := cb.state
	cb.state = newState
	cb.lastStateChange = time.Now()
	cb.metrics.recordStateChange(oldState, newState)

	// Reset counters based on new state
	switch newState {
	case StateClosed:
		cb.failures = 0
		cb.successes = 0
	case StateHalfOpen:
		cb.successes = 0
		cb.halfOpenRequests = 0
	case StateOpen:
		cb.successes = 0
		cb.halfOpenRequests = 0
	}

	cb.logger.Info("Circuit breaker state changed",
		logging.String("name", cb.name),
		logging.String("old_state", oldState.String()),
		logging.String("new_state", newState.String()),
	)

	// Call state change callbacks
	for _, callback := range cb.onStateChange {
		go callback(oldState, newState)
	}
}

// GetFailureThreshold returns the configured failure threshold
func (cb *CircuitBreaker) GetFailureThreshold() int {
	return cb.failureThreshold
}

// GetSuccessThreshold returns the configured success threshold
func (cb *CircuitBreaker) GetSuccessThreshold() int {
	return cb.successThreshold
}

// GetTimeout returns the configured timeout
func (cb *CircuitBreaker) GetTimeout() time.Duration {
	return cb.timeout
}

// GetName returns the circuit breaker name
func (cb *CircuitBreaker) GetName() string {
	return cb.name
}

// GetLastFailureTime returns the time of the last failure
func (cb *CircuitBreaker) GetLastFailureTime() time.Time {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.lastFailureTime
}

// GetLastStateChange returns the time of the last state change
func (cb *CircuitBreaker) GetLastStateChange() time.Time {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.lastStateChange
}

// GetFailureCount returns the current failure count
func (cb *CircuitBreaker) GetFailureCount() int {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failures
}

// GetSuccessCount returns the current success count
func (cb *CircuitBreaker) GetSuccessCount() int {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.successes
}

// TimeUntilHalfOpen returns the duration until the circuit will transition to half-open
func (cb *CircuitBreaker) TimeUntilHalfOpen() time.Duration {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	if cb.state != StateOpen {
		return 0
	}

	elapsed := time.Since(cb.lastFailureTime)
	if elapsed >= cb.timeout {
		return 0
	}

	return cb.timeout - elapsed
}

// CircuitBreakerMetrics methods
func (m *CircuitBreakerMetrics) recordRequest(latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalRequests++
	m.lastRequestTime = time.Now()

	if latency > 0 {
		m.totalLatency += latency
		if m.totalRequests > 0 {
			m.averageLatency = m.totalLatency / time.Duration(m.totalRequests)
		}
	}
}

func (m *CircuitBreakerMetrics) recordSuccess() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalSuccesses++
	m.lastSuccessTime = time.Now()
}

func (m *CircuitBreakerMetrics) recordFailure() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalFailures++
	m.lastFailureTime = time.Now()
}

func (m *CircuitBreakerMetrics) recordRejected() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalTimeouts++
}

func (m *CircuitBreakerMetrics) recordStateChange(oldState, newState State) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalStateChanges++

	switch {
	case oldState == StateClosed && newState == StateOpen:
		m.closedToOpenCount++
	case oldState == StateOpen && newState == StateHalfOpen:
		m.openToHalfOpenCount++
	case oldState == StateHalfOpen && newState == StateClosed:
		m.halfOpenToClosedCount++
	case oldState == StateHalfOpen && newState == StateOpen:
		m.halfOpenToOpenCount++
	}
}

// GetTotalRequests returns total number of requests
func (m *CircuitBreakerMetrics) GetTotalRequests() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.totalRequests
}

// GetTotalSuccesses returns total number of successes
func (m *CircuitBreakerMetrics) GetTotalSuccesses() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.totalSuccesses
}

// GetTotalFailures returns total number of failures
func (m *CircuitBreakerMetrics) GetTotalFailures() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.totalFailures
}

// GetTotalTimeouts returns total number of rejected requests
func (m *CircuitBreakerMetrics) GetTotalTimeouts() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.totalTimeouts
}

// GetSuccessRate returns the success rate as a percentage
func (m *CircuitBreakerMetrics) GetSuccessRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.totalRequests == 0 {
		return 0.0
	}

	return float64(m.totalSuccesses) / float64(m.totalRequests) * 100.0
}

// GetFailureRate returns the failure rate as a percentage
func (m *CircuitBreakerMetrics) GetFailureRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.totalRequests == 0 {
		return 0.0
	}

	return float64(m.totalFailures) / float64(m.totalRequests) * 100.0
}

// GetAverageLatency returns the average request latency
func (m *CircuitBreakerMetrics) GetAverageLatency() time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.averageLatency
}

// GetLastRequestTime returns the time of the last request
func (m *CircuitBreakerMetrics) GetLastRequestTime() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastRequestTime
}

// GetStateChangeCount returns total number of state changes
func (m *CircuitBreakerMetrics) GetStateChangeCount() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.totalStateChanges
}

// Reset resets all metrics
func (m *CircuitBreakerMetrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalRequests = 0
	m.totalSuccesses = 0
	m.totalFailures = 0
	m.totalTimeouts = 0
	m.totalStateChanges = 0
	m.closedToOpenCount = 0
	m.openToHalfOpenCount = 0
	m.halfOpenToClosedCount = 0
	m.halfOpenToOpenCount = 0
	m.averageLatency = 0
	m.totalLatency = 0
}

// Errors
var (
	// ErrCircuitOpen is returned when the circuit is open and requests are blocked
	ErrCircuitOpen = errors.New("circuit breaker is open")
	// ErrCircuitHalfOpen is returned when the circuit is half-open and request limit is reached
	ErrCircuitHalfOpen = errors.New("circuit breaker is half-open, request limit reached")
)

// CircuitBreakerRegistry manages multiple circuit breakers
type CircuitBreakerRegistry struct {
	mu            sync.RWMutex
	breakers      map[string]*CircuitBreaker
	defaultConfig CircuitBreakerConfig
}

// NewCircuitBreakerRegistry creates a new circuit breaker registry
func NewCircuitBreakerRegistry() *CircuitBreakerRegistry {
	return &CircuitBreakerRegistry{
		breakers:      make(map[string]*CircuitBreaker),
		defaultConfig: DefaultCircuitBreakerConfig(),
	}
}

// GetOrCreate gets an existing circuit breaker or creates a new one
func (r *CircuitBreakerRegistry) GetOrCreate(name string, config *CircuitBreakerConfig) *CircuitBreaker {
	r.mu.Lock()
	defer r.mu.Unlock()

	if cb, exists := r.breakers[name]; exists {
		return cb
	}

	if config == nil {
		cfg := r.defaultConfig
		cfg.Name = name
		config = &cfg
	} else {
		config.Name = name
	}

	cb := NewCircuitBreaker(*config)
	r.breakers[name] = cb
	return cb
}

// Get retrieves a circuit breaker by name
func (r *CircuitBreakerRegistry) Get(name string) (*CircuitBreaker, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cb, exists := r.breakers[name]
	return cb, exists
}

// Delete removes a circuit breaker from the registry
func (r *CircuitBreakerRegistry) Delete(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.breakers, name)
}

// List returns all circuit breaker names
func (r *CircuitBreakerRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.breakers))
	for name := range r.breakers {
		names = append(names, name)
	}
	return names
}

// GetAll returns all circuit breakers
func (r *CircuitBreakerRegistry) GetAll() map[string]*CircuitBreaker {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]*CircuitBreaker, len(r.breakers))
	for name, cb := range r.breakers {
		result[name] = cb
	}
	return result
}

// SetDefaultConfig sets the default configuration for new circuit breakers
func (r *CircuitBreakerRegistry) SetDefaultConfig(config CircuitBreakerConfig) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.defaultConfig = config
}

// CircuitBreakerMiddleware creates a middleware function for request handling
type CircuitBreakerMiddleware struct {
	cb *CircuitBreaker
}

// NewCircuitBreakerMiddleware creates a new circuit breaker middleware
func NewCircuitBreakerMiddleware(cb *CircuitBreaker) *CircuitBreakerMiddleware {
	return &CircuitBreakerMiddleware{cb: cb}
}

// Handle wraps a handler function with circuit breaker logic
func (m *CircuitBreakerMiddleware) Handle(handler func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		return m.cb.Execute(ctx, handler)
	}
}

// WithRetry wraps a handler with retry logic when circuit is half-open
func (m *CircuitBreakerMiddleware) WithRetry(maxRetries int, retryDelay time.Duration) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		var lastErr error

		for attempt := 0; attempt <= maxRetries; attempt++ {
			if err := m.cb.Allow(); err != nil {
				if attempt < maxRetries {
					select {
					case <-time.After(retryDelay):
						continue
					case <-ctx.Done():
						return ctx.Err()
					}
				}
				return err
			}

			lastErr = nil
			// Handler will be called by the caller
			break
		}

		return lastErr
	}
}

// String returns a string representation of the circuit breaker
func (cb *CircuitBreaker) String() string {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return fmt.Sprintf("CircuitBreaker{name=%s, state=%s, failures=%d, successes=%d, threshold=%d}",
		cb.name, cb.state, cb.failures, cb.successes, cb.failureThreshold)
}

// HealthStatus returns the health status of the circuit breaker
type HealthStatus struct {
	Name              string  `json:"name"`
	State             string  `json:"state"`
	Healthy           bool    `json:"healthy"`
	FailureCount      int     `json:"failure_count"`
	SuccessCount      int     `json:"success_count"`
	FailureThreshold  int     `json:"failure_threshold"`
	TimeUntilHalfOpen string  `json:"time_until_half_open,omitempty"`
	SuccessRate       float64 `json:"success_rate"`
}

// GetHealthStatus returns the current health status
func (cb *CircuitBreaker) GetHealthStatus() HealthStatus {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	status := HealthStatus{
		Name:             cb.name,
		State:            cb.state.String(),
		Healthy:          cb.state == StateClosed || cb.state == StateHalfOpen,
		FailureCount:     cb.failures,
		SuccessCount:     cb.successes,
		FailureThreshold: cb.failureThreshold,
		SuccessRate:      cb.metrics.GetSuccessRate(),
	}

	if cb.state == StateOpen {
		status.TimeUntilHalfOpen = cb.TimeUntilHalfOpen().String()
	}

	return status
}
