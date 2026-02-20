// Package ratelimit_test provides comprehensive tests for rate limiting functionality
package ratelimit_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"biometrics-cli/pkg/ratelimit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTokenBucketLimiter_Basic tests basic token bucket functionality
func TestTokenBucketLimiter_Basic(t *testing.T) {
	limit := ratelimit.Limit{
		Requests:   10,
		PerSeconds: 2.0,
		Burst:      5,
	}

	limiter := ratelimit.NewTokenBucketLimiter(limit, time.Minute)
	require.NotNil(t, limiter)

	ctx := context.Background()
	key := "test-key-1"

	// Should allow burst requests
	allowed := 0
	for i := 0; i < 5; i++ {
		if limiter.Allow(ctx, key) {
			allowed++
		}
	}
	assert.Equal(t, 5, allowed, "Should allow burst of 5 requests")

	// Next requests should be rate limited
	time.Sleep(100 * time.Millisecond)
	if limiter.Allow(ctx, key) {
		allowed++
	}
	assert.GreaterOrEqual(t, allowed, 5, "Should have allowed at least burst")
}

// TestTokenBucketLimiter_BurstHandling tests burst capacity
func TestTokenBucketLimiter_BurstHandling(t *testing.T) {
	limit := ratelimit.Limit{
		Requests:   100,
		PerSeconds: 10.0,
		Burst:      20,
	}

	limiter := ratelimit.NewTokenBucketLimiter(limit, time.Minute)
	ctx := context.Background()
	key := "burst-test"

	// Exhaust burst
	burstAllowed := 0
	for i := 0; i < 25; i++ {
		if limiter.Allow(ctx, key) {
			burstAllowed++
		}
	}

	assert.Equal(t, 20, burstAllowed, "Should allow exactly burst capacity")
}

// TestTokenBucketLimiter_Concurrent tests thread safety with concurrent requests
func TestTokenBucketLimiter_Concurrent(t *testing.T) {
	limit := ratelimit.Limit{
		Requests:   1000,
		PerSeconds: 100.0,
		Burst:      100,
	}

	limiter := ratelimit.NewTokenBucketLimiter(limit, 10*time.Second)
	ctx := context.Background()
	key := "concurrent-test"

	var wg sync.WaitGroup
	var mu sync.Mutex
	allowed := 0

	// Launch 50 concurrent goroutines
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				if limiter.Allow(ctx, key) {
					mu.Lock()
					allowed++
					mu.Unlock()
				}
				time.Sleep(time.Millisecond)
			}
		}()
	}

	wg.Wait()
	assert.Greater(t, allowed, 0, "Should allow some concurrent requests")
	assert.LessOrEqual(t, allowed, 500, "Should not exceed total requests")
}

// TestTokenBucketLimiter_Cleanup tests cleanup of old entries
func TestTokenBucketLimiter_Cleanup(t *testing.T) {
	limit := ratelimit.Limit{
		Requests:   10,
		PerSeconds: 1.0,
		Burst:      5,
	}

	limiter := ratelimit.NewTokenBucketLimiter(limit, 100*time.Millisecond)
	ctx := context.Background()

	// Create multiple keys
	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	for _, key := range keys {
		limiter.Allow(ctx, key)
	}

	// Wait for cleanup
	time.Sleep(200 * time.Millisecond)

	// Keys should still exist (accessed recently)
	for _, key := range keys {
		allowed := limiter.Allow(ctx, key)
		assert.True(t, allowed, "Key should still be accessible")
	}
}

// TestTokenBucketLimiter_Reset tests resetting limiter state
func TestTokenBucketLimiter_Reset(t *testing.T) {
	limit := ratelimit.Limit{
		Requests:   10,
		PerSeconds: 1.0,
		Burst:      5,
	}

	limiter := ratelimit.NewTokenBucketLimiter(limit, time.Minute)
	ctx := context.Background()
	key := "reset-test"

	// Use up burst
	for i := 0; i < 5; i++ {
		limiter.Allow(ctx, key)
	}

	// Reset
	limiter.Reset(key)

	// Should allow burst again
	allowed := limiter.Allow(ctx, key)
	assert.True(t, allowed, "Should allow after reset")
}

// TestTokenBucketLimiter_GetLimit tests retrieving limit configuration
func TestTokenBucketLimiter_GetLimit(t *testing.T) {
	expectedLimit := ratelimit.Limit{
		Requests:   50,
		PerSeconds: 5.0,
		Burst:      10,
	}

	limiter := ratelimit.NewTokenBucketLimiter(expectedLimit, time.Minute)
	key := "limit-test"

	limit := limiter.GetLimit(key)
	assert.Equal(t, expectedLimit.Requests, limit.Requests)
	assert.Equal(t, expectedLimit.PerSeconds, limit.PerSeconds)
	assert.Equal(t, expectedLimit.Burst, limit.Burst)
}

// TestTokenBucketLimiter_Wait tests Wait functionality
func TestTokenBucketLimiter_Wait(t *testing.T) {
	limit := ratelimit.Limit{
		Requests:   10,
		PerSeconds: 10.0,
		Burst:      2,
	}

	limiter := ratelimit.NewTokenBucketLimiter(limit, time.Minute)
	ctx := context.Background()
	key := "wait-test"

	// Exhaust burst
	limiter.Allow(ctx, key)
	limiter.Allow(ctx, key)

	// Wait should block until token available
	start := time.Now()
	err := limiter.Wait(ctx, key)
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Greater(t, elapsed, 0*time.Millisecond, "Wait should block")
}

// TestTokenBucketLimiter_WaitContextCancellation tests Wait with context cancellation
func TestTokenBucketLimiter_WaitContextCancellation(t *testing.T) {
	limit := ratelimit.Limit{
		Requests:   10,
		PerSeconds: 0.1, // Very slow refill
		Burst:      1,
	}

	limiter := ratelimit.NewTokenBucketLimiter(limit, time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	key := "wait-cancel-test"

	// Exhaust burst
	limiter.Allow(ctx, key)

	// Wait should return context error
	err := limiter.Wait(ctx, key)
	assert.Error(t, err, "Wait should return error on context cancellation")
}

// TestSlidingWindowLimiter_Basic tests basic sliding window functionality
func TestSlidingWindowLimiter_Basic(t *testing.T) {
	limit := 10
	window := time.Second

	limiter := ratelimit.NewSlidingWindowLimiter(limit, window)
	ctx := context.Background()
	key := "sliding-test"

	// Should allow up to limit
	allowed := 0
	for i := 0; i < 15; i++ {
		if limiter.Allow(ctx, key) {
			allowed++
		}
	}

	assert.Equal(t, limit, allowed, "Should allow exactly limit requests")
}

// TestSlidingWindowLimiter_WindowReset tests window reset after duration
func TestSlidingWindowLimiter_WindowReset(t *testing.T) {
	limit := 5
	window := 200 * time.Millisecond

	limiter := ratelimit.NewSlidingWindowLimiter(limit, window)
	ctx := context.Background()
	key := "sliding-reset-test"

	// Exhaust limit
	for i := 0; i < limit; i++ {
		limiter.Allow(ctx, key)
	}

	// Should be limited
	assert.False(t, limiter.Allow(ctx, key), "Should be rate limited")

	// Wait for window to pass (add buffer for timing)
	time.Sleep(300 * time.Millisecond)

	// Should allow again
	assert.True(t, limiter.Allow(ctx, key), "Should allow after window reset")
}

// TestSlidingWindowLimiter_Concurrent tests thread safety
func TestSlidingWindowLimiter_Concurrent(t *testing.T) {
	limit := 100
	window := time.Second

	limiter := ratelimit.NewSlidingWindowLimiter(limit, window)
	ctx := context.Background()
	key := "sliding-concurrent-test"

	var wg sync.WaitGroup
	var mu sync.Mutex
	allowed := 0

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				if limiter.Allow(ctx, key) {
					mu.Lock()
					allowed++
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()
	assert.Equal(t, limit, allowed, "Should allow exactly limit requests")
}

// TestSlidingWindowLimiter_GetLimit tests retrieving limit configuration
func TestSlidingWindowLimiter_GetLimit(t *testing.T) {
	limit := 20
	window := 2 * time.Second

	limiter := ratelimit.NewSlidingWindowLimiter(limit, window)
	key := "sliding-limit-test"

	config := limiter.GetLimit(key)
	assert.Equal(t, limit, config.Requests)
	assert.Equal(t, limit, config.Burst)
	assert.InDelta(t, float64(limit)/window.Seconds(), config.PerSeconds, 0.01)
}

// TestSlidingWindowLimiter_Reset tests resetting limiter state
func TestSlidingWindowLimiter_Reset(t *testing.T) {
	limit := 5
	window := time.Second

	limiter := ratelimit.NewSlidingWindowLimiter(limit, window)
	ctx := context.Background()
	key := "sliding-reset-test"

	// Exhaust limit
	for i := 0; i < limit; i++ {
		limiter.Allow(ctx, key)
	}

	// Reset
	limiter.Reset(key)

	// Should allow again
	assert.True(t, limiter.Allow(ctx, key), "Should allow after reset")
}

// TestSlidingWindowLimiter_Wait tests Wait functionality
func TestSlidingWindowLimiter_Wait(t *testing.T) {
	limit := 2
	window := 300 * time.Millisecond

	limiter := ratelimit.NewSlidingWindowLimiter(limit, window)
	ctx := context.Background()
	key := "sliding-wait-test"

	// Exhaust limit
	limiter.Allow(ctx, key)
	limiter.Allow(ctx, key)

	// Wait should block until window passes
	start := time.Now()
	err := limiter.Wait(ctx, key)
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, elapsed, 100*time.Millisecond, "Wait should block")
}

// TestRateLimitMiddleware_Basic tests basic middleware functionality
func TestRateLimitMiddleware_Basic(t *testing.T) {
	limit := ratelimit.Limit{
		Requests:   10,
		PerSeconds: 10.0,
		Burst:      5,
	}

	limiter := ratelimit.NewTokenBucketLimiter(limit, time.Minute)
	keyFunc := func(ctx context.Context) string {
		if key, ok := ctx.Value("key").(string); ok {
			return key
		}
		return "default"
	}

	middleware := ratelimit.NewRateLimitMiddleware(limiter, keyFunc)
	require.NotNil(t, middleware)

	ctx := context.WithValue(context.Background(), "key", "middleware-test")
	called := false

	handler := middleware.Handle(func(ctx context.Context) error {
		called = true
		return nil
	})

	err := handler(ctx)
	require.NoError(t, err)
	assert.True(t, called, "Handler should be called")
}

// TestRateLimitMiddleware_LimitExceeded tests middleware when limit exceeded
func TestRateLimitMiddleware_LimitExceeded(t *testing.T) {
	limit := ratelimit.Limit{
		Requests:   10,
		PerSeconds: 10.0,
		Burst:      2,
	}

	limiter := ratelimit.NewTokenBucketLimiter(limit, time.Minute)
	keyFunc := func(ctx context.Context) string {
		return "limit-test"
	}

	middleware := ratelimit.NewRateLimitMiddleware(limiter, keyFunc)

	limitExceededCalled := false
	middleware.SetOnLimitExceeded(func(ctx context.Context, key string) {
		limitExceededCalled = true
	})

	ctx := context.Background()

	// Exhaust burst
	limiter.Allow(ctx, "limit-test")
	limiter.Allow(ctx, "limit-test")

	handler := middleware.Handle(func(ctx context.Context) error {
		return nil
	})

	err := handler(ctx)
	assert.Error(t, err, "Should return error when limit exceeded")
	assert.True(t, limitExceededCalled, "Callback should be called")
}

// TestRateLimitMiddleware_IPBased tests IP-based rate limiting
func TestRateLimitMiddleware_IPBased(t *testing.T) {
	limit := ratelimit.Limit{
		Requests:   10,
		PerSeconds: 10.0,
		Burst:      5,
	}

	limiter := ratelimit.NewTokenBucketLimiter(limit, time.Minute)

	ipKeyFunc := func(ctx context.Context) string {
		if ip, ok := ctx.Value("ip").(string); ok {
			return "ip:" + ip
		}
		return "ip:unknown"
	}

	middleware := ratelimit.NewRateLimitMiddleware(limiter, ipKeyFunc)

	// Test with different IPs
	ips := []string{"192.168.1.1", "192.168.1.2", "10.0.0.1"}
	for _, ip := range ips {
		ctx := context.WithValue(context.Background(), "ip", ip)
		called := false

		handler := middleware.Handle(func(ctx context.Context) error {
			called = true
			return nil
		})

		err := handler(ctx)
		require.NoError(t, err)
		assert.True(t, called, "Handler should be called for IP %s", ip)
	}
}

// TestConfig_DefaultConfig tests default configuration
func TestConfig_DefaultConfig(t *testing.T) {
	config := ratelimit.DefaultConfig()

	assert.Greater(t, config.RequestsPerSecond, float64(0), "RequestsPerSecond should be positive")
	assert.Greater(t, config.Burst, 0, "Burst should be positive")
	assert.Greater(t, config.CleanupInterval, time.Duration(0), "CleanupInterval should be positive")
	assert.Greater(t, config.WindowDuration, time.Duration(0), "WindowDuration should be positive")
}

// TestRateLimiter_Interface tests that implementations satisfy the interface
func TestRateLimiter_Interface(t *testing.T) {
	limit := ratelimit.Limit{
		Requests:   10,
		PerSeconds: 1.0,
		Burst:      5,
	}

	var _ ratelimit.RateLimiter = ratelimit.NewTokenBucketLimiter(limit, time.Minute)
	var _ ratelimit.RateLimiter = ratelimit.NewSlidingWindowLimiter(10, time.Second)
}

// BenchmarkTokenBucketLimiter_Allow benchmarks token bucket Allow
func BenchmarkTokenBucketLimiter_Allow(b *testing.B) {
	limit := ratelimit.Limit{
		Requests:   1000,
		PerSeconds: 100.0,
		Burst:      100,
	}
	limiter := ratelimit.NewTokenBucketLimiter(limit, time.Minute)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Allow(ctx, "bench-key")
	}
}

// BenchmarkSlidingWindowLimiter_Allow benchmarks sliding window Allow
func BenchmarkSlidingWindowLimiter_Allow(b *testing.B) {
	limiter := ratelimit.NewSlidingWindowLimiter(1000, time.Second)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Allow(ctx, "bench-key")
	}
}

// BenchmarkTokenBucketLimiter_Concurrent benchmarks concurrent access
func BenchmarkTokenBucketLimiter_Concurrent(b *testing.B) {
	limit := ratelimit.Limit{
		Requests:   10000,
		PerSeconds: 1000.0,
		Burst:      1000,
	}
	limiter := ratelimit.NewTokenBucketLimiter(limit, time.Minute)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow(ctx, "parallel-key")
		}
	})
}
