package cache

import (
	"context"
	"time"
)

type CacheStrategy interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte) error
}

type LRUStrategy struct {
	cache   *RedisCache
	maxSize int
	ttl     time.Duration
}

func NewLRUStrategy(cache *RedisCache, maxSize int, ttl time.Duration) *LRUStrategy {
	return &LRUStrategy{
		cache:   cache,
		maxSize: maxSize,
		ttl:     ttl,
	}
}

func (l *LRUStrategy) Get(ctx context.Context, key string) ([]byte, error) {
	return l.cache.Get(ctx, key)
}

func (l *LRUStrategy) Set(ctx context.Context, key string, value []byte) error {
	return l.cache.Set(ctx, key, value, l.ttl)
}

type TTLStrategy struct {
	cache *RedisCache
	ttl   time.Duration
}

func NewTTLStrategy(cache *RedisCache, ttl time.Duration) *TTLStrategy {
	return &TTLStrategy{
		cache: cache,
		ttl:   ttl,
	}
}

func (t *TTLStrategy) Get(ctx context.Context, key string) ([]byte, error) {
	return t.cache.Get(ctx, key)
}

func (t *TTLStrategy) Set(ctx context.Context, key string, value []byte) error {
	return t.cache.Set(ctx, key, value, t.ttl)
}

type WriteThroughStrategy struct {
	cache *RedisCache
	ttl   time.Duration
}

func NewWriteThroughStrategy(cache *RedisCache, ttl time.Duration) *WriteThroughStrategy {
	return &WriteThroughStrategy{
		cache: cache,
		ttl:   ttl,
	}
}

func (w *WriteThroughStrategy) Get(ctx context.Context, key string) ([]byte, error) {
	return w.cache.Get(ctx, key)
}

func (w *WriteThroughStrategy) Set(ctx context.Context, key string, value []byte) error {
	return w.cache.Set(ctx, key, value, w.ttl)
}
