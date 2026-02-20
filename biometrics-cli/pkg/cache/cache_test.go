package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupTestCache(t *testing.T) *RedisCache {
	logger, _ := zap.NewDevelopment()
	config := DefaultCacheConfig()
	config.Addr = "localhost:6379"

	cache, err := NewRedisCache(config, logger)
	if err != nil {
		t.Skip("Redis not available, skipping test")
		return nil
	}

	t.Cleanup(func() {
		cache.Close()
	})

	return cache
}

func TestRedisCache_GetSet(t *testing.T) {
	cache := setupTestCache(t)
	if cache == nil {
		return
	}

	ctx := context.Background()
	key := "test:key1"
	value := []byte("test value")

	err := cache.Set(ctx, key, value, time.Minute)
	require.NoError(t, err)

	got, err := cache.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, value, got)
}

func TestRedisCache_GetJSON(t *testing.T) {
	cache := setupTestCache(t)
	if cache == nil {
		return
	}

	ctx := context.Background()
	key := "test:json1"
	data := map[string]interface{}{
		"name":  "test",
		"value": 123,
	}

	err := cache.SetJSON(ctx, key, data, time.Minute)
	require.NoError(t, err)

	var got map[string]interface{}
	err = cache.GetJSON(ctx, key, &got)
	require.NoError(t, err)
	assert.Equal(t, data["name"], got["name"])
	assert.Equal(t, float64(123), got["value"])
}

func TestRedisCache_Delete(t *testing.T) {
	cache := setupTestCache(t)
	if cache == nil {
		return
	}

	ctx := context.Background()
	key := "test:delete1"
	value := []byte("to delete")

	err := cache.Set(ctx, key, value, time.Minute)
	require.NoError(t, err)

	err = cache.Delete(ctx, key)
	require.NoError(t, err)

	got, err := cache.Get(ctx, key)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestRedisCache_Exists(t *testing.T) {
	cache := setupTestCache(t)
	if cache == nil {
		return
	}

	ctx := context.Background()
	key := "test:exists1"

	exists, err := cache.Exists(ctx, key)
	require.NoError(t, err)
	assert.False(t, exists)

	err = cache.Set(ctx, key, []byte("value"), time.Minute)
	require.NoError(t, err)

	exists, err = cache.Exists(ctx, key)
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestCacheMetrics(t *testing.T) {
	metrics := NewCacheMetrics()

	metrics.RecordGet(100*time.Millisecond, true)
	metrics.RecordGet(200*time.Millisecond, false)
	metrics.RecordGet(50*time.Millisecond, true)

	assert.Equal(t, int64(2), metrics.Hits)
	assert.Equal(t, int64(1), metrics.Misses)
	assert.InDelta(t, 0.667, metrics.GetHitRate(), 0.01)
}

func TestLRUStrategy(t *testing.T) {
	cache := setupTestCache(t)
	if cache == nil {
		return
	}

	strategy := NewLRUStrategy(cache, 100, time.Minute)
	ctx := context.Background()

	err := strategy.Set(ctx, "lru:key1", []byte("value1"))
	require.NoError(t, err)

	got, err := strategy.Get(ctx, "lru:key1")
	require.NoError(t, err)
	assert.Equal(t, []byte("value1"), got)
}

func TestCacheInvalidator(t *testing.T) {
	cache := setupTestCache(t)
	if cache == nil {
		return
	}

	invalidator := NewCacheInvalidator(cache, InvalidationImmediate, 0)
	ctx := context.Background()

	err := cache.Set(ctx, "invalidate:key1", []byte("value1"), time.Minute)
	require.NoError(t, err)

	err = invalidator.Invalidate(ctx, "invalidate:key1")
	require.NoError(t, err)

	got, err := cache.Get(ctx, "invalidate:key1")
	require.NoError(t, err)
	assert.Nil(t, got)
}
