# Rate Limiting Package

**Purpose:** Rate limiting utilities for API and service protection

## Overview

This package provides rate limiting capabilities to prevent abuse and ensure fair usage of resources across the biometrics CLI.

## Features

- Token bucket algorithm
- Leaky bucket algorithm
- Sliding window limits
- Distributed rate limiting
- Custom limiters

## Usage

### Basic Rate Limiter

```go
import "github.com/delqhi/biometrics/pkg/ratelimit"

limiter := ratelimit.NewLimiter(ratelimit.Config{
    Rate:       100,  // requests per second
    Burst:      10,   // burst capacity
})

// Acquire permit
err := limiter.Acquire(ctx)
if err != nil {
    return err
}
```

### With Callbacks

```go
limiter := ratelimit.NewLimiter(ratelimit.Config{
    Rate:   60,
    Burst:  10,
    OnWait: func(waitTime time.Duration) {
        log.Printf("Rate limited, waiting %v", waitTime)
    },
    OnReject: func() {
        metrics.Inc("rate_limit_rejected")
    },
})
```

### Token Bucket

```go
limiter := ratelimit.NewTokenBucket(ctx, ratelimit.TokenBucketConfig{
    Capacity:    1000,
    RefillRate:  100,  // tokens per second
})
```

### Distributed Rate Limiting

```go
limiter := ratelimit.NewDistributed(ctx, ratelimit.DistributedConfig{
    Backend:     "redis",
    RedisClient: redisClient,
    Key:         "rate_limit:user:123",
    Rate:        100,
    Burst:       20,
})
```

## Configuration

```yaml
ratelimit:
  enabled: true
  global:
    rate: 1000
    burst: 100
  per_user:
    rate: 100
    burst: 10
  per_ip:
    rate: 50
    burst: 5
```

## Middleware Integration

```go
// HTTP middleware
handler := ratelimit.Middleware(limiter, next)

// CLI middleware
limiter := ratelimit.NewCLI(ctx, config)
err := limiter.Acquire(ctx)
```

## Monitoring

```go
// Get current state
state := limiter.State()
fmt.Printf("Available: %d, WaitTime: %v\n", state.Available, state.WaitTime)

// Prometheus metrics
prometheus.MustRegister(ratelimit.NewCollector("biometrics"))
```

## Related Documentation

- [Security Configuration](../docs/security.md)
- [Monitoring](../docs/monitoring.md)
- [Configuration Reference](../docs/configuration.md)
