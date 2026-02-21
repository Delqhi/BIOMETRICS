# Redis Cache Package

Production-ready Redis cache layer with advanced strategies and invalidation.

## Features

- **Connection Pooling**: Efficient connection management
- **Multiple Strategies**: LRU, TTL, Write-Through caching
- **Cache Invalidation**: Immediate, delayed, and lazy policies
- **Metrics & Monitoring**: Hit rate, latency tracking
- **JSON Support**: Automatic serialization/deserialization
- **Key Prefixing**: Namespace isolation
- **Error Handling**: Graceful fallback on cache failures
- **Context Support**: Cancellation and timeouts

## Installation

```bash
go get github.com/go-redis/redis/v8
go get go.uber.org/zap
```

## Usage

### Basic Cache

```go
package main

import (
    "context"
    "log"
    "time"
    "github.com/delqhi/biometrics/pkg/cache"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewDevelopment()
    config := cache.DefaultCacheConfig()
    config.Addr = "localhost:6379"
    
    redisCache, err := cache.NewRedisCache(config, logger)
    if err != nil {
        log.Fatal(err)
    }
    defer redisCache.Close()
    
    ctx := context.Background()
    
    // Set value
    err = redisCache.Set(ctx, "user:123", []byte("John Doe"), 5*time.Minute)
    if err != nil {
        log.Fatal(err)
    }
    
    // Get value
    value, err := redisCache.Get(ctx, "user:123")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("User: %s", string(value))
}
```

### JSON Cache

```go
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Set JSON
user := &User{ID: "123", Name: "John", Email: "john@example.com"}
err := redisCache.SetJSON(ctx, "user:123", user, 10*time.Minute)

// Get JSON
var cached User
err = redisCache.GetJSON(ctx, "user:123", &cached)
```

### Cache Strategies

```go
// LRU Strategy
lru := cache.NewLRUStrategy(redisCache, 1000, 5*time.Minute)
lru.Set(ctx, "key", []byte("value"))

// TTL Strategy
ttl := cache.NewTTLStrategy(redisCache, 10*time.Minute)
ttl.Set(ctx, "key", []byte("value"))

// Write-Through Strategy
wt := cache.NewWriteThroughStrategy(redisCache, 5*time.Minute)
wt.Set(ctx, "key", []byte("value"))
```

### Cache Invalidation

```go
// Create invalidator
invalidator := cache.NewCacheInvalidator(
    redisCache,
    cache.InvalidationImmediate,
    0,
)

// Invalidate single key
err := invalidator.Invalidate(ctx, "user:123")

// Invalidate pattern
err = invalidator.InvalidatePattern(ctx, "user:*")

// Add patterns for bulk invalidation
invalidator.AddPattern("user:*")
invalidator.AddPattern("session:*")
err = invalidator.InvalidateAll(ctx)
```

### Cache Metrics

```go
metrics := redisCache.GetMetrics()

fmt.Printf("Hits: %d\n", metrics.Hits)
fmt.Printf("Misses: %d\n", metrics.Misses)
fmt.Printf("Hit Rate: %.2f%%\n", metrics.GetHitRate()*100)
fmt.Printf("Avg Latency: %v\n", metrics.GetAverageLatency())

// Reset metrics
metrics.Reset()
```

## Configuration

```go
config := cache.CacheConfig{
    Addr:         "localhost:6379",       // Redis address
    Password:     "",                      // Redis password (optional)
    DB:           0,                       // Database number
    PoolSize:     100,                     // Connection pool size
    MinIdleConns: 10,                      // Minimum idle connections
    DialTimeout:  5 * time.Second,         // Connection timeout
    ReadTimeout:  3 * time.Second,         // Read timeout
    WriteTimeout: 3 * time.Second,         // Write timeout
    Prefix:       "biometrics:",           // Key prefix
    DefaultTTL:   5 * time.Minute,         // Default TTL
}
```

## Invalidation Policies

- **Immediate**: Delete immediately when triggered
- **Delayed**: Delete after specified delay
- **Lazy**: No automatic deletion (lazy expiration)

## Best Practices

1. **Use Key Prefixes**: Namespace your keys to avoid collisions
2. **Set Appropriate TTLs**: Balance freshness vs. cache efficiency
3. **Monitor Hit Rates**: Target >80% hit rate for optimal performance
4. **Handle Cache Misses**: Implement fallback logic for misses
5. **Use Connection Pooling**: Configure pool size based on workload
6. **Implement Circuit Breaker**: Fail fast when Redis is unavailable

## Examples

### Task Results Cache

```go
type TaskResult struct {
    ID     string                 `json:"id"`
    Status string                 `json:"status"`
    Data   map[string]interface{} `json:"data"`
}

func GetTaskResult(ctx context.Context, cache *cache.RedisCache, taskID string) (*TaskResult, error) {
    var result TaskResult
    err := cache.GetJSON(ctx, "task:result:"+taskID, &result)
    if err != nil {
        return nil, err
    }
    return &result, nil
}

func CacheTaskResult(ctx context.Context, cache *cache.RedisCache, taskID string, result *TaskResult) error {
    return cache.SetJSON(ctx, "task:result:"+taskID, result, 5*time.Minute)
}
```

### Agent State Cache

```go
func CacheAgentState(ctx context.Context, cache *cache.RedisCache, agentID string, state string) error {
    return cache.Set(ctx, "agent:state:"+agentID, []byte(state), 30*time.Second)
}

func GetAgentState(ctx context.Context, cache *cache.RedisCache, agentID string) (string, error) {
    data, err := cache.Get(ctx, "agent:state:"+agentID)
    if err != nil {
        return "", err
    }
    if data == nil {
        return "", cache.ErrCacheMiss
    }
    return string(data), nil
}
```

## Files

- `redis.go` - Core Redis cache implementation
- `strategies.go` - Cache strategies (LRU, TTL, Write-Through)
- `invalidation.go` - Cache invalidation policies
- `metrics.go` - Cache metrics and monitoring
- `cache_test.go` - Comprehensive tests

## Testing

```bash
# Run tests (requires Redis running on localhost:6379)
go test ./pkg/cache/... -v

# Run with coverage
go test ./pkg/cache/... -coverprofile=coverage.out
```

## Related

- [Performance Monitoring](../performance/) - pprof integration
- [Rate Limiting](../ratelimit/) - Rate limiting with Redis
