package circuit

import (
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

type CircuitBreaker struct {
	mu           sync.RWMutex
	name         string
	state        State
	failures     int
	successes    int
	maxFailures  int
	timeout      time.Duration
	halfOpenMax  int
	lastFailure  time.Time
	resetTimeout time.Duration
}

type CircuitBreakerConfig struct {
	Name         string
	MaxFailures  int
	Timeout      time.Duration
	HalfOpenMax  int
	ResetTimeout time.Duration
}

func NewCircuitBreaker(config *CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		name:         config.Name,
		state:        StateClosed,
		failures:     0,
		successes:    0,
		maxFailures:  config.MaxFailures,
		timeout:      config.Timeout,
		halfOpenMax:  config.HalfOpenMax,
		resetTimeout: config.ResetTimeout,
	}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.allowRequest() {
		return &CircuitOpenError{Name: cb.name}
	}

	err := fn()
	cb.recordResult(err)

	return err
}

func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			cb.state = StateHalfOpen
			cb.successes = 0
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFailure = time.Now()

		if cb.state == StateHalfOpen {
			cb.state = StateOpen
			cb.failures = 0
		} else if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
		}
	} else {
		cb.successes++
		cb.failures = 0

		if cb.state == StateHalfOpen && cb.successes >= cb.halfOpenMax {
			cb.state = StateClosed
			cb.successes = 0
		}
	}
}

func (cb *CircuitBreaker) GetState() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

func (cb *CircuitBreaker) GetFailureCount() int {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failures
}

func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.failures = 0
	cb.successes = 0
}

type CircuitOpenError struct {
	Name string
}

func (e *CircuitOpenError) Error() string {
	return "circuit " + e.Name + " is open"
}

type CircuitBreakerRegistry struct {
	mu       sync.RWMutex
	breakers map[string]*CircuitBreaker
}

var registry = &CircuitBreakerRegistry{
	breakers: make(map[string]*CircuitBreaker),
}

func GetOrCreate(name string, config *CircuitBreakerConfig) *CircuitBreaker {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	if cb, exists := registry.breakers[name]; exists {
		return cb
	}

	cb := NewCircuitBreaker(config)
	registry.breakers[name] = cb
	return cb
}

func GetAllBreakers() []*CircuitBreaker {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	breakers := make([]*CircuitBreaker, 0, len(registry.breakers))
	for _, cb := range registry.breakers {
		breakers = append(breakers, cb)
	}
	return breakers
}
