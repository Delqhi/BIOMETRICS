// Package ratelimit provides thread-safe rate limiting implementations
// including token bucket and sliding window algorithms with IP-based limiting.
package ratelimit

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter defines the interface for rate limiting implementations
type RateLimiter interface {
	Allow(ctx context.Context, key string) bool
	Wait(ctx context.Context, key string) error
	GetLimit(key string) Limit
	Reset(key string)
	Stats(key string) LimiterStats
}

// Limit represents rate limiting configuration
type Limit struct {
	Requests   int     // Number of requests
	PerSeconds float64 // Requests per second
	Burst      int     // Maximum burst size
}

// LimiterStats provides statistics about rate limiter state
type LimiterStats struct {
	Remaining  int           // Remaining requests in current window
	ResetAt    time.Time     // When the limit resets
	RetryAfter time.Duration // Time to wait before retrying
}

// Config holds configuration for rate limiters
type Config struct {
	RequestsPerSecond float64       // Base rate limit
	Burst             int           // Burst allowance
	WindowDuration    time.Duration // Window for sliding window
	CleanupInterval   time.Duration // Cleanup interval for stale entries
	EnableIPBased     bool          // Enable IP-based limiting
	SlidingWindow     bool          // Use sliding window instead of token bucket
}

// DefaultConfig returns a reasonable default configuration
func DefaultConfig() Config {
	return Config{
		RequestsPerSecond: 10.0,
		Burst:             20,
		WindowDuration:    time.Minute,
		CleanupInterval:   time.Minute,
		EnableIPBased:     true,
		SlidingWindow:     false,
	}
}

// TokenBucketLimiter implements token bucket rate limiting with per-key limiters
type TokenBucketLimiter struct {
	limiters        sync.Map
	defaultLimit    Limit
	mu              sync.RWMutex
	cleanupInterval time.Duration
	closed          chan struct{}
}

type bucket struct {
	limiter    *rate.Limiter
	lastAccess time.Time
	hitCount   int64
}

func NewTokenBucketLimiter(defaultLimit Limit, cleanupInterval time.Duration) *TokenBucketLimiter {
	rl := &TokenBucketLimiter{
		defaultLimit:    defaultLimit,
		cleanupInterval: cleanupInterval,
		closed:          make(chan struct{}),
	}

	go rl.cleanupLoop()

	return rl
}

func (rl *TokenBucketLimiter) Allow(ctx context.Context, key string) bool {
	limiter := rl.getLimiter(key)
	return limiter.Allow()
}

func (rl *TokenBucketLimiter) Wait(ctx context.Context, key string) error {
	limiter := rl.getLimiter(key)
	return limiter.Wait(ctx)
}

func (rl *TokenBucketLimiter) GetLimit(key string) Limit {
	return rl.defaultLimit
}

func (rl *TokenBucketLimiter) Reset(key string) {
	rl.limiters.Delete(key)
}

// Stats returns statistics about the rate limiter state
func (rl *TokenBucketLimiter) Stats(key string) LimiterStats {
	if stored, ok := rl.limiters.Load(key); ok {
		b := stored.(*bucket)
		remaining := b.limiter.Tokens()
		return LimiterStats{
			Remaining:  int(remaining),
			ResetAt:    b.lastAccess.Add(time.Second),
			RetryAfter: time.Until(b.lastAccess.Add(time.Second)),
		}
	}
	return LimiterStats{
		Remaining:  rl.defaultLimit.Burst,
		ResetAt:    time.Now(),
		RetryAfter: 0,
	}
}

func (rl *TokenBucketLimiter) getLimiter(key string) *rate.Limiter {
	if stored, ok := rl.limiters.Load(key); ok {
		b := stored.(*bucket)
		b.lastAccess = time.Now()
		return b.limiter
	}

	limiter := rate.NewLimiter(rate.Limit(rl.defaultLimit.PerSeconds), rl.defaultLimit.Burst)
	b := &bucket{
		limiter:    limiter,
		lastAccess: time.Now(),
	}

	rl.limiters.Store(key, b)
	return limiter
}

func (rl *TokenBucketLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanup()
	}
}

func (rl *TokenBucketLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	rl.limiters.Range(func(key, value interface{}) bool {
		b := value.(*bucket)
		if now.Sub(b.lastAccess) > 5*time.Minute {
			rl.limiters.Delete(key)
		}
		return true
	})
}

type SlidingWindowLimiter struct {
	windows sync.Map
	limit   int
	window  time.Duration
	mu      sync.RWMutex
}

type window struct {
	counts      map[int64]int
	lastCleanup time.Time
}

func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		limit:  limit,
		window: window,
	}
}

func (sw *SlidingWindowLimiter) Allow(ctx context.Context, key string) bool {
	now := time.Now()
	windowKey := now.UnixNano() / int64(sw.window)

	stored, _ := sw.windows.LoadOrStore(key, &window{
		counts:      make(map[int64]int),
		lastCleanup: now,
	})

	w := stored.(*window)

	sw.mu.Lock()
	if now.Sub(w.lastCleanup) > sw.window {
		for ts := range w.counts {
			if ts < windowKey-1 {
				delete(w.counts, ts)
			}
		}
		w.lastCleanup = now
	}

	total := 0
	for _, count := range w.counts {
		total += count
	}

	if total >= sw.limit {
		sw.mu.Unlock()
		return false
	}

	w.counts[windowKey]++
	sw.mu.Unlock()

	return true
}

func (sw *SlidingWindowLimiter) Wait(ctx context.Context, key string) error {
	for !sw.Allow(ctx, key) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
	return nil
}

func (sw *SlidingWindowLimiter) GetLimit(key string) Limit {
	return Limit{
		Requests:   sw.limit,
		PerSeconds: float64(sw.limit) / sw.window.Seconds(),
		Burst:      sw.limit,
	}
}

func (sw *SlidingWindowLimiter) Reset(key string) {
	sw.windows.Delete(key)
}

// Stats returns statistics about the sliding window limiter state
func (sw *SlidingWindowLimiter) Stats(key string) LimiterStats {
	stored, ok := sw.windows.Load(key)
	if !ok {
		return LimiterStats{
			Remaining:  sw.limit,
			ResetAt:    time.Now().Add(sw.window),
			RetryAfter: 0,
		}
	}

	w := stored.(*window)
	total := 0
	for _, count := range w.counts {
		total += count
	}

	remaining := sw.limit - total
	if remaining < 0 {
		remaining = 0
	}

	return LimiterStats{
		Remaining:  remaining,
		ResetAt:    time.Now().Add(sw.window),
		RetryAfter: sw.window,
	}
}

type RateLimitMiddleware struct {
	limiter         RateLimiter
	keyFunc         func(ctx context.Context) string
	onLimitExceeded func(ctx context.Context, key string)
}

func NewRateLimitMiddleware(limiter RateLimiter, keyFunc func(ctx context.Context) string) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter:         limiter,
		keyFunc:         keyFunc,
		onLimitExceeded: func(ctx context.Context, key string) {},
	}
}

func (m *RateLimitMiddleware) SetOnLimitExceeded(fn func(ctx context.Context, key string)) {
	m.onLimitExceeded = fn
}

func (m *RateLimitMiddleware) Handle(next func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		key := m.keyFunc(ctx)

		if !m.limiter.Allow(ctx, key) {
			m.onLimitExceeded(ctx, key)
			return ErrRateLimitExceeded
		}

		return next(ctx)
	}
}
