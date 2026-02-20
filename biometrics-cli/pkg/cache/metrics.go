package cache

import (
	"sync"
	"sync/atomic"
	"time"
)

type CacheMetrics struct {
	Hits         int64
	Misses       int64
	Sets         int64
	Deletes      int64
	Errors       int64
	TotalLatency int64
	mu           sync.RWMutex
}

func NewCacheMetrics() *CacheMetrics {
	return &CacheMetrics{}
}

func (m *CacheMetrics) RecordGet(duration time.Duration, hit bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if hit {
		atomic.AddInt64(&m.Hits, 1)
	} else {
		atomic.AddInt64(&m.Misses, 1)
	}
	atomic.AddInt64(&m.TotalLatency, int64(duration))
}

func (m *CacheMetrics) RecordSet(duration time.Duration, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	atomic.AddInt64(&m.Sets, 1)
	if !success {
		atomic.AddInt64(&m.Errors, 1)
	}
}

func (m *CacheMetrics) RecordDelete(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	atomic.AddInt64(&m.Deletes, 1)
}

func (m *CacheMetrics) GetHitRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total := m.Hits + m.Misses
	if total == 0 {
		return 0
	}
	return float64(m.Hits) / float64(total)
}

func (m *CacheMetrics) GetAverageLatency() time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total := m.Hits + m.Misses + m.Sets
	if total == 0 {
		return 0
	}
	return time.Duration(m.TotalLatency / total)
}

func (m *CacheMetrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Hits = 0
	m.Misses = 0
	m.Sets = 0
	m.Deletes = 0
	m.Errors = 0
	m.TotalLatency = 0
}

func (m *CacheMetrics) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return "CacheMetrics{" +
		"Hits=" + string(rune(m.Hits)) +
		", Misses=" + string(rune(m.Misses)) +
		", Sets=" + string(rune(m.Sets)) +
		", Deletes=" + string(rune(m.Deletes)) +
		", Errors=" + string(rune(m.Errors)) +
		", HitRate=" + string(rune(m.GetHitRate())) +
		"}"
}
