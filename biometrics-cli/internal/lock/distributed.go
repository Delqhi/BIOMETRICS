package lock

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type DistributedLock struct {
	mu       sync.RWMutex
	locks    map[string]*lockInfo
	redisURL string
}

type lockInfo struct {
	holder    string
	expiresAt time.Time
	refCount  int
}

func NewDistributedLock(redisURL string) *DistributedLock {
	return &DistributedLock{
		locks:    make(map[string]*lockInfo),
		redisURL: redisURL,
	}
}

func (dl *DistributedLock) Acquire(ctx context.Context, key, holder string, ttl time.Duration) (bool, error) {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	if info, exists := dl.locks[key]; exists {
		if time.Now().Before(info.expiresAt) && info.holder != holder {
			return false, nil
		}
	}

	dl.locks[key] = &lockInfo{
		holder:    holder,
		expiresAt: time.Now().Add(ttl),
		refCount:  1,
	}

	return true, nil
}

func (dl *DistributedLock) Release(key, holder string) error {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	if info, exists := dl.locks[key]; exists {
		if info.holder == holder {
			delete(dl.locks, key)
			return nil
		}
		return fmt.Errorf("lock held by different holder: %s", info.holder)
	}

	return fmt.Errorf("lock not found: %s", key)
}

func (dl *DistributedLock) Extend(key, holder string, ttl time.Duration) error {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	if info, exists := dl.locks[key]; exists {
		if info.holder != holder {
			return fmt.Errorf("lock held by different holder: %s", info.holder)
		}
		info.expiresAt = time.Now().Add(ttl)
		return nil
	}

	return fmt.Errorf("lock not found: %s", key)
}

func (dl *DistributedLock) IsLocked(key string) bool {
	dl.mu.RLock()
	defer dl.mu.RUnlock()

	if info, exists := dl.locks[key]; exists {
		return time.Now().Before(info.expiresAt)
	}
	return false
}

func (dl *DistributedLock) GetHolder(key string) string {
	dl.mu.RLock()
	defer dl.mu.RUnlock()

	if info, exists := dl.locks[key]; exists {
		if time.Now().Before(info.expiresAt) {
			return info.holder
		}
	}
	return ""
}

func (dl *DistributedLock) ForceRelease(key string) error {
	dl.mu.Lock()
	defer dl.mu.Unlock()
	delete(dl.locks, key)
	return nil
}

func (dl *DistributedLock) CleanExpired() int {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	now := time.Now()
	cleaned := 0

	for key, info := range dl.locks {
		if now.After(info.expiresAt) {
			delete(dl.locks, key)
			cleaned++
		}
	}

	return cleaned
}

func (dl *DistributedLock) GetStats() map[string]interface{} {
	dl.mu.RLock()
	defer dl.mu.RUnlock()

	now := time.Now()
	active := 0
	expired := 0

	for _, info := range dl.locks {
		if now.Before(info.expiresAt) {
			active++
		} else {
			expired++
		}
	}

	return map[string]interface{}{
		"total_locks": len(dl.locks),
		"active":      active,
		"expired":     expired,
	}
}
