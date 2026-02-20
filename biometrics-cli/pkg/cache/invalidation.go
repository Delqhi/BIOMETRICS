package cache

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type InvalidationPolicy string

const (
	InvalidationImmediate InvalidationPolicy = "immediate"
	InvalidationDelayed   InvalidationPolicy = "delayed"
	InvalidationLazy      InvalidationPolicy = "lazy"
)

type CacheInvalidator struct {
	cache    *RedisCache
	policy   InvalidationPolicy
	delay    time.Duration
	patterns []string
}

func NewCacheInvalidator(cache *RedisCache, policy InvalidationPolicy, delay time.Duration) *CacheInvalidator {
	return &CacheInvalidator{
		cache:    cache,
		policy:   policy,
		delay:    delay,
		patterns: make([]string, 0),
	}
}

func (ci *CacheInvalidator) AddPattern(pattern string) {
	ci.patterns = append(ci.patterns, pattern)
}

func (ci *CacheInvalidator) Invalidate(ctx context.Context, key string) error {
	switch ci.policy {
	case InvalidationImmediate:
		return ci.cache.Delete(ctx, key)
	case InvalidationDelayed:
		go func() {
			time.Sleep(ci.delay)
			_ = ci.cache.Delete(context.Background(), key)
		}()
		return nil
	case InvalidationLazy:
		return nil
	default:
		return fmt.Errorf("unknown invalidation policy: %s", ci.policy)
	}
}

func (ci *CacheInvalidator) InvalidatePattern(ctx context.Context, pattern string) error {
	keys, err := ci.findKeysByPattern(ctx, pattern)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if err := ci.cache.Delete(ctx, key); err != nil {
			return err
		}
	}

	return nil
}

func (ci *CacheInvalidator) findKeysByPattern(ctx context.Context, pattern string) ([]string, error) {
	fullPattern := ci.cache.config.Prefix + pattern
	return ci.scanKeys(ctx, fullPattern)
}

func (ci *CacheInvalidator) scanKeys(ctx context.Context, pattern string) ([]string, error) {
	var keys []string
	cursor := uint64(0)
	pattern = strings.ReplaceAll(pattern, "*", "*")

	for {
		result, cursor, err := ci.cache.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}

		keys = append(keys, result...)

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

func (ci *CacheInvalidator) InvalidateAll(ctx context.Context) error {
	for _, pattern := range ci.patterns {
		if err := ci.InvalidatePattern(ctx, pattern); err != nil {
			return err
		}
	}
	return nil
}
