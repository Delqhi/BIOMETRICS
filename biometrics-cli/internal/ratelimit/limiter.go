package ratelimit

import (
	"biometrics-cli/internal/metrics"
	"biometrics-cli/internal/state"
	"fmt"
	"sync"
	"time"
)

type Limiter struct {
	mu            sync.RWMutex
	buckets       map[string]*Bucket
	defaultRate   float64
	defaultBurst  int
	defaultWindow time.Duration
}

type Bucket struct {
	mu          sync.Mutex
	tokens      float64
	maxTokens   float64
	refillRate  float64
	lastRefill  time.Time
	refillEvery time.Duration
}

type Rule struct {
	Key         string
	Rate        float64
	Burst       int
	Window      time.Duration
	Description string
}

var defaultLimiter = &Limiter{
	buckets:       make(map[string]*Bucket),
	defaultRate:   100,
	defaultBurst:  10,
	defaultWindow: time.Second,
}

var LimiterInstance = defaultLimiter

func New(rate float64, burst int, window time.Duration) *Limiter {
	return &Limiter{
		buckets:       make(map[string]*Bucket),
		defaultRate:   rate,
		defaultBurst:  burst,
		defaultWindow: window,
	}
}

func (l *Limiter) Allow(key string) (bool, error) {
	bucket, _ := l.getOrCreateBucket(key)

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(bucket.lastRefill)
	bucket.tokens += elapsed.Seconds() * bucket.refillRate
	if bucket.tokens > bucket.maxTokens {
		bucket.tokens = bucket.maxTokens
	}
	bucket.lastRefill = now

	if bucket.tokens >= 1 {
		bucket.tokens -= 1
		metrics.RateLimitAllowedTotal.WithLabelValues(key).Inc()
		return true, nil
	}

	metrics.RateLimitRejectedTotal.WithLabelValues(key).Inc()
	state.GlobalState.Log("WARN", fmt.Sprintf("Rate limit exceeded for key: %s", key))
	return false, fmt.Errorf("rate limit exceeded")
}

func (l *Limiter) AllowN(key string, n int) (bool, error) {
	bucket, _ := l.getOrCreateBucket(key)

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(bucket.lastRefill)
	bucket.tokens += elapsed.Seconds() * bucket.refillRate
	if bucket.tokens > bucket.maxTokens {
		bucket.tokens = bucket.maxTokens
	}
	bucket.lastRefill = now

	if bucket.tokens >= float64(n) {
		bucket.tokens -= float64(n)
		metrics.RateLimitAllowedTotal.WithLabelValues(key).Inc()
		return true, nil
	}

	metrics.RateLimitRejectedTotal.WithLabelValues(key).Inc()
	state.GlobalState.Log("WARN", fmt.Sprintf("Rate limit exceeded for key: %s", key))
	return false, fmt.Errorf("rate limit exceeded")
}

func (l *Limiter) getOrCreateBucket(key string) (*Bucket, bool) {
	l.mu.RLock()
	bucket, exists := l.buckets[key]
	l.mu.RUnlock()

	if exists {
		return bucket, true
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if bucket, exists := l.buckets[key]; exists {
		return bucket, true
	}

	bucket = &Bucket{
		tokens:      float64(l.defaultBurst),
		maxTokens:   float64(l.defaultBurst),
		refillRate:  l.defaultRate,
		lastRefill:  time.Now(),
		refillEvery: l.defaultWindow,
	}

	l.buckets[key] = bucket
	return bucket, false
}

func (l *Limiter) SetRate(key string, rate float64, burst int, window time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.buckets[key] = &Bucket{
		tokens:      float64(burst),
		maxTokens:   float64(burst),
		refillRate:  rate,
		lastRefill:  time.Now(),
		refillEvery: window,
	}

	state.GlobalState.Log("INFO", fmt.Sprintf("Set rate limit for %s: %.2f req/s, burst %d", key, rate, burst))
}

func (l *Limiter) GetRate(key string) (float64, int, time.Duration) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	bucket, exists := l.buckets[key]
	if !exists {
		return l.defaultRate, l.defaultBurst, l.defaultWindow
	}

	return bucket.refillRate, int(bucket.maxTokens), bucket.refillEvery
}

func (l *Limiter) RemoveRule(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.buckets, key)
	state.GlobalState.Log("INFO", fmt.Sprintf("Removed rate limit rule for: %s", key))
}

func (l *Limiter) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.buckets = make(map[string]*Bucket)
	state.GlobalState.Log("INFO", "Cleared all rate limit rules")
}

func (l *Limiter) GetStats() map[string]interface{} {
	l.mu.RLock()
	defer l.mu.RUnlock()

	stats := make(map[string]interface{})
	for key, bucket := range l.buckets {
		bucket.mu.Lock()
		stats[key] = map[string]interface{}{
			"tokens":      bucket.tokens,
			"max_tokens":  bucket.maxTokens,
			"refill_rate": bucket.refillRate,
		}
		bucket.mu.Unlock()
	}

	return stats
}

func (l *Limiter) AddRules(rules []Rule) {
	for _, rule := range rules {
		l.SetRate(rule.Key, rule.Rate, rule.Burst, rule.Window)
	}
}

type MultiLimiter struct {
	mu       sync.RWMutex
	limiters []*Limiter
}

func NewMultiLimiter(limiters ...*Limiter) *MultiLimiter {
	return &MultiLimiter{
		limiters: limiters,
	}
}

func (m *MultiLimiter) Allow(key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, limiter := range m.limiters {
		allowed, err := limiter.Allow(key)
		if err != nil || !allowed {
			return allowed, err
		}
	}

	return true, nil
}

func (m *MultiLimiter) AllowN(key string, n int) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, limiter := range m.limiters {
		allowed, err := limiter.AllowN(key, n)
		if err != nil || !allowed {
			return allowed, err
		}
	}

	return true, nil
}

func Allow(key string) (bool, error) {
	return LimiterInstance.Allow(key)
}

func AllowN(key string, n int) (bool, error) {
	return LimiterInstance.AllowN(key, n)
}

func SetRate(key string, rate float64, burst int, window time.Duration) {
	LimiterInstance.SetRate(key, rate, burst, window)
}

func RemoveRule(key string) {
	LimiterInstance.RemoveRule(key)
}

func GetStats() map[string]interface{} {
	return LimiterInstance.GetStats()
}
