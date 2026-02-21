# PERFORMANCE.md - Biometrics CLI Performance Documentation

**Project:** Biometrics CLI  
**Version:** 1.0.0  
**Date:** February 2026  
**Status:** Sprint 1 Feature Documentation

---

## Table of Contents

1. [Performance Overview](#performance-overview)
2. [pprof Integration Guide](#pprof-integration-guide)
3. [Redis Caching Setup](#redis-caching-setup)
4. [Performance Benchmarks](#performance-benchmarks)
5. [Optimization Strategies](#optimization-strategies)
6. [Monitoring and Alerting](#monitoring-and-alerting)
7. [Example Configurations](#example-configurations)

---

## Performance Overview

The Biometrics CLI implements comprehensive performance monitoring and optimization features as part of Sprint 1. This documentation covers the tools and techniques available for analyzing, monitoring, and optimizing application performance.

### Core Performance Features

The performance architecture consists of several interconnected components:

- **pprof Integration:** Built-in Go profiling support for CPU, memory, goroutine, and mutex analysis
- **Redis Caching:** Distributed caching layer with intelligent invalidation strategies
- **Runtime Metrics:** Real-time collection of Go runtime statistics
- **Automated Profiling:** Scheduled profiling capture for trend analysis
- **HTTP Endpoints:** Built-in profiling HTTP server for production debugging

### Performance Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                   BIOMETRICS PERFORMANCE ARCHITECTURE              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │                  PROFILING LAYER                             │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐        │  │
│  │  │   CPU    │ │  Memory  │ │Goroutine │ │  Mutex  │        │  │
│  │  │Profiler  │ │ Profiler │ │ Profiler │ │Profiler │        │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘        │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                              │                                       │
│                              ▼                                       │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │                    CACHING LAYER                              │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐        │  │
│  │  │ LRU      │ │  Redis   │ │ Write-   │ │  Cache  │        │  │
│  │  │ Cache    │ │  Backend │ │  Behind  │ │Invalidatn│        │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘        │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                              │                                       │
│                              ▼                                       │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │                  MONITORING LAYER                             │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐        │  │
│  │  │ Runtime  │ │   HTTP   │ │  Custom  │ │ GC Stats │        │  │
│  │  │ Metrics  │ │ Handler  │ │ Counters │ │          │        │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘        │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## pprof Integration Guide

The Biometrics CLI includes comprehensive pprof integration for profiling Go applications. pprof is a tool for visualization and analysis of profiling data.

### Basic Profiler Usage

```go
package main

import (
    "log"
    "time"
    
    "github.com/delqhi/biometrics/pkg/performance"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewDevelopment()
    
    config := performance.DefaultProfilerConfig()
    profiler, err := performance.NewProfiler(config, logger)
    if err != nil {
        log.Fatal(err)
    }
    
    // Start all profilers
    if err := profiler.Start(); err != nil {
        log.Fatal(err)
    }
    
    // Your application code here
    doWork()
    
    // Stop and write profiles
    if err := profiler.Stop(); err != nil {
        log.Fatal(err)
    }
}

func doWork() {
    // Simulate work
    time.Sleep(100 * time.Millisecond)
}
```

### Profiling Manager (Advanced)

For more control over profiling, use the `ProfilingManager`:

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/delqhi/biometrics/pkg/performance"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewDevelopment()
    
    config := performance.ProfilerConfig{
        EnableCPU:        true,
        EnableMemory:     true,
        EnableBlock:      true,
        EnableMutex:      true,
        EnableGoroutine:  true,
        ProfileDir:       "./profiles",
        HTTPAddr:         "localhost:6060",
        EnableHTTP:       true,
        ProfileDuration:  30 * time.Second,
    }
    
    manager, err := performance.NewProfilingManager(config, logger)
    if err != nil {
        panic(err)
    }
    
    // Start CPU profiling for 30 seconds
    ctx := context.Background()
    err = manager.StartProfile(ctx, performance.ProfileTypeCPU, 30*time.Second)
    if err != nil {
        panic(err)
    }
    
    // Do your work here
    performApplicationLogic()
    
    // Stop CPU profiling
    err = manager.StopProfile(performance.ProfileTypeCPU)
    if err != nil {
        panic(err)
    }
    
    // Capture all profiles
    snapshots, err := manager.CaptureAllProfiles(ctx)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Captured profiles: %+v\n", snapshots)
    
    // Get runtime metrics
    metrics := manager.GetMetrics()
    fmt.Printf("Goroutines: %d\n", metrics.Goroutines)
    fmt.Printf("Heap Alloc: %d bytes\n", metrics.HeapAlloc)
    fmt.Printf("GC Count: %d\n", metrics.NumGC)
}

func performApplicationLogic() {
    // Application code to profile
}
```

### HTTP Profiling Server

For production debugging, use the built-in HTTP profiling endpoints:

```go
package main

import (
    "net/http"
    
    "github.com/delqhi/biometrics/pkg/performance"
)

func main() {
    mux := http.NewServeMux()
    
    // Register pprof handlers
    performance.RegisterPprofHandlers(mux)
    
    // Add custom middleware
    handler := performance.WithPprofMiddleware(mux)
    
    // Start server
    http.ListenAndServe("localhost:6060", handler)
}
```

### Available Profiling Endpoints

Once running, access these endpoints:

| Endpoint | Description |
|----------|-------------|
| `/debug/pprof/` | Index of all profiles |
| `/debug/pprof/cpu` | CPU profile (30 second default) |
| `/debug/pprof/heap` | Heap memory profile |
| `/debug/pprof/allocs` | Allocation profile |
| `/debug/pprof/block` | Block profile (goroutine blocking) |
| `/debug/pprof/goroutine` | Goroutine profile |
| `/debug/pprof/mutex` | Mutex profile |
| `/debug/pprof/trace` | Execution trace |

### Analyzing Profiles

Command-line analysis:

```bash
# Analyze CPU profile
go tool pprof profiles/cpu_20260220_025632.prof

# Generate flame graph (requires graphviz)
go tool pprof -svg profiles/cpu_20260220_025632.prof > cpu_flamegraph.svg

# View top functions
go tool pprof -top profiles/cpu_20260220_025632.prof

# Compare profiles
go tool pprof -base profiles/cpu_baseline.prof profiles/cpu_current.prof

# Web interface
go tool pprof -http=:8080 profiles/cpu_20260220_025632.prof
```

### Profile Types Explained

- **CPU Profile:** Shows which functions consume the most CPU time
- **Memory Profile:** Shows heap allocations and memory usage
- **Block Profile:** Shows where goroutines block on synchronization
- **Mutex Profile:** Shows where mutex contention occurs
- **Goroutine Profile:** Shows the state of all goroutines

---

## Redis Caching Setup

The Biometrics CLI includes a comprehensive Redis caching layer with multiple caching strategies and intelligent invalidation.

### Basic Redis Configuration

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/delqhi/biometrics/pkg/cache"
)

func main() {
    config := &cache.RedisConfig{
        Addr:         "localhost:6379",
        Password:     "",
        DB:           0,
        PoolSize:     10,
        MinIdleConns: 5,
        DialTimeout:   5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    }
    
    client, err := cache.NewRedisClient(config)
    if err != nil {
        panic(err)
    }
    defer client.Close()
    
    // Test connection
    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        panic(err)
    }
    
    fmt.Println("Redis connected successfully")
}
```

### Cache Operations

```go
// Basic cache operations
func cacheOperations(client *cache.RedisClient) error {
    ctx := context.Background()
    
    // Set a value with expiration
    err := client.Set(ctx, "user:123", `{"name":"John"}`, 24*time.Hour).Err()
    if err != nil {
        return err
    }
    
    // Get a value
    val, err := client.Get(ctx, "user:123").Result()
    if err != nil {
        return err
    }
    fmt.Printf("User: %s\n", val)
    
    // Delete a key
    client.Del(ctx, "user:123")
    
    // Check if key exists
    exists, err := client.Exists(ctx, "user:123").Result()
    if err != nil {
        return err
    }
    fmt.Printf("Exists: %d\n", exists)
    
    return nil
}
```

### Cache Strategies

```go
// Configure different caching strategies
func configureStrategies() {
    // LRU (Least Recently Used)
    lruConfig := cache.StrategyConfig{
        Type:      cache.StrategyLRU,
        MaxSize:   1000,
        TTL:       1 * time.Hour,
    }
    
    // LFU (Least Frequently Used)
    lfuConfig := cache.StrategyConfig{
        Type:      cache.StrategyLFU,
        MaxSize:   5000,
        TTL:       24 * time.Hour,
    }
    
    // TTL-based
    ttlConfig := cache.StrategyConfig{
        Type:      cache.StrategyTTL,
        MaxSize:   100,
        TTL:       30 * time.Minute,
    }
    
    // Write-behind
    behindConfig := cache.StrategyConfig{
        Type:           cache.StrategyWriteBehind,
        MaxSize:        1000,
        TTL:            1 * time.Hour,
        AsyncWrite:     true,
        WriteBatchSize: 100,
    }
}
```

### Cache Invalidation

```go
// Advanced cache invalidation
func cacheInvalidation(client *cache.RedisClient) error {
    ctx := context.Background()
    
    // Invalidate by pattern
    err := client.InvalidateByPattern(ctx, "user:*").Err()
    if err != nil {
        return err
    }
    
    // Invalidate by tag
    client.AddTag(ctx, "user:123", "users", "active")
    client.AddTag(ctx, "user:456", "users", "active")
    
    // Invalidate all keys with tag
    err = client.InvalidateByTag(ctx, "users").Err()
    if err != nil {
        return err
    }
    
    // Invalidate by prefix
    err = client.InvalidateByPrefix(ctx, "session:").Err()
    if err != nil {
        return err
    }
    
    // Manual invalidation
    client.Invalidate(ctx, "user:123")
    
    return nil
}
```

### Distributed Cache

```go
// Create distributed cache with Redis
func distributedCache() (*cache.DistributedCache, error) {
    redisConfig := &cache.RedisConfig{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    }
    
    return cache.NewDistributedCache(redisConfig, cache.CacheOptions{
        DefaultTTL:     1 * time.Hour,
        KeyPrefix:     "biometrics:",
        EnableLRU:      true,
        MaxConnections: 100,
    })
}

// Use distributed cache
func useCache(dist *cache.DistributedCache) error {
    ctx := context.Background()
    
    // User cache
    userCache := dist.WithPrefix("user:")
    
    // Get or set
    cached, err := userCache.GetOrSet(ctx, "123", func() (string, time.Duration, error) {
        // Fetch from database
        user := fetchUserFromDB("123")
        return user, 1 * time.Hour, nil
    })
    if err != nil {
        return err
    }
    
    fmt.Printf("Cached user: %s\n", cached)
    return nil
}

func fetchUserFromDB(id string) string {
    // Simulate DB fetch
    return `{"id":"123","name":"John"}`
}
```

---

## Performance Benchmarks

The Biometrics CLI includes comprehensive benchmarks to measure and track performance.

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. -benchmem ./...

# Run specific benchmark
go test -bench=Profiler -benchmem ./pkg/performance/...

# Run with CPU profile
go test -bench=. -cpuprofile=cpu.prof ./...

# Run with memory profile
go test -bench=. -memprofile=mem.prof ./...

# Run with allocs profile
go test -bench=. -benchtime=5s ./...
```

### Benchmark Results (Sample)

| Operation | Operations/sec | ns/op | Allocations |
|-----------|---------------|-------|-------------|
| Cache Get (hit) | 1,250,000 | 800 | 2 |
| Cache Set | 950,000 | 1,050 | 4 |
| Profile Start | 50,000 | 20,000 | 15 |
| Metrics Collection | 100,000 | 10,000 | 8 |
| Redis Connect | 5,000 | 200,000 | 45 |

### Profile Analysis

```go
// Example profile analysis
func analyzeProfile(profilePath string) {
    f, err := os.Open(profilePath)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    
    pprof := pprof.NewProfile("cpu")
    if err := pprof.ReadFrom(f); err != nil {
        panic(err)
    }
    
    // Top 10 functions by cumulative time
    top := pprof.Top(10)
    for i, sample := range top {
        fmt.Printf("%d. %s: %d samples\n", i+1, sample.Value, sample.Value)
    }
}
```

---

## Optimization Strategies

### Memory Optimization

```go
// Optimize memory allocation
func optimizeMemory() {
    // Use sync.Pool for frequently allocated objects
    var bufferPool = sync.Pool{
        New: func() interface{} {
            return make([]byte, 4096)
        },
    }
    
    // Get from pool
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf)
    
    // Use buffered channels
    ch := make(chan int, 1000)
    
    // Pre-allocate slices
    items := make([]Item, 0, 1000)
}
```

### CPU Optimization

```go
// CPU optimization techniques
func optimizeCPU() {
    // Use worker pools
    workers := runtime.NumCPU()
    var wg sync.WaitGroup
    jobs := make(chan Job, workers)
    
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                processJob(job)
            }
        }()
    }
    
    // Send jobs
    for _, job := range jobsList {
        jobs <- job
    }
    close(jobs)
    wg.Wait()
}
```

### Cache Optimization

```go
// Cache optimization
func optimizeCache() {
    // Use appropriate TTL values
    // - Frequently changing data: short TTL (minutes)
    // - Static data: long TTL (hours/days)
    // - User-specific: session-based TTL
    
    // Implement cache-aside pattern
    func GetUser(id string) (*User, error) {
        // Check cache first
        cached, err := cache.Get("user:" + id)
        if err == nil {
            return unmarshalUser(cached)
        }
        
        // Fetch from database
        user, err := db.GetUser(id)
        if err != nil {
            return nil, err
        }
        
        // Store in cache
        cache.Set("user:"+id, marshalUser(user), 15*time.Minute)
        
        return user, nil
    }
}
```

### Concurrency Optimization

```go
// Optimize goroutine usage
func optimizeConcurrency() {
    // Limit concurrent operations
    semaphore := make(chan struct{}, 10)
    
    for _, task := range tasks {
        semaphore <- struct{}{}
        go func(t Task) {
            defer func() { <-semaphore }()
            processTask(t)
        }(task)
    }
    
    // Use errgroup for coordinated operations
    var g errgroup.Group
    results := make([]Result, len(items))
    
    for i, item := range items {
        i, item := i, item
        g.Go(func() error {
            results[i] = processItem(item)
            return nil
        })
    }
    
    if err := g.Wait(); err != nil {
        // Handle error
    }
}
```

---

## Monitoring and Alerting

### Runtime Metrics

```go
// Collect runtime metrics
func collectMetrics() {
    profiler := performance.NewProfiler(performance.DefaultProfilerConfig(), logger)
    
    // Get current metrics
    metrics := profiler.ReportMetrics()
    
    // Custom metrics collection
    customMetrics := map[string]float64{
        "requests_total":    0,
        "requests_success":  0,
        "requests_failure":  0,
        "response_time_ms":  0,
    }
    
    // Expose via HTTP
    http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
        metrics := profiler.ReportMetrics()
        fmt.Fprintf(w, "goroutines %d\n", metrics.Goroutines)
        fmt.Fprintf(w, "heap_alloc %d\n", metrics.HeapAlloc)
        fmt.Fprintf(w, "gc_count %d\n", metrics.NumGC)
    })
}
```

### Health Checks

```go
// Health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
    profiler := performance.NewProfiler(performance.DefaultProfilerConfig(), logger)
    metrics := profiler.ReportMetrics()
    
    health := map[string]interface{}{
        "status":      "healthy",
        "goroutines": metrics.Goroutines,
        "memory": map[string]int64{
            "heap_alloc": metrics.HeapAlloc,
            "total_alloc": metrics.TotalAlloc,
            "sys":         metrics.Sys,
        },
        "gc": map[string]uint32{
            "count": metrics.NumGC,
        },
    }
    
    // Check thresholds
    if metrics.HeapAlloc > 1<<30 { // 1GB
        health["status"] = "degraded"
    }
    
    json.NewEncoder(w).Encode(health)
}
```

### Alerting Rules

Configure alerts based on metrics:

| Metric | Warning Threshold | Critical Threshold |
|--------|------------------|-------------------|
| Memory Usage | > 70% | > 90% |
| GC Frequency | > 10/sec | > 50/sec |
| Goroutines | > 1000 | > 5000 |
| Response Time | > 500ms | > 2s |
| Error Rate | > 1% | > 5% |

---

## Example Configurations

### Development Configuration

```yaml
# config-dev.yaml
performance:
  profiling:
    enabled: true
    profile_dir: ./profiles
    http_addr: localhost:6060
    cpu_enabled: true
    memory_enabled: true
    
  cache:
    redis:
      enabled: false
    local:
      enabled: true
      max_size: 100
      ttl: 10m
      
  metrics:
    enabled: true
    interval: 10s
```

### Production Configuration

```yaml
# config-prod.yaml
performance:
  profiling:
    enabled: true
    profile_dir: /var/lib/biometrics/profiles
    http_addr: localhost:6060
    cpu_enabled: true
    memory_enabled: true
    auto_profile: true
    profile_interval: 1h
    
  cache:
    redis:
      enabled: true
      addr: redis.internal:6379
      password: ${REDIS_PASSWORD}
      pool_size: 50
      min_idle_conns: 10
      dial_timeout: 5s
      read_timeout: 3s
      write_timeout: 3s
    local:
      enabled: true
      max_size: 10000
      ttl: 30m
      
  metrics:
    enabled: true
    interval: 5s
    export_prometheus: true
    prometheus_port: 9090
    
  alerting:
    enabled: true
    memory_threshold: 1073741824  # 1GB
    goroutine_threshold: 5000
    gc_threshold: 50
```

### Docker Compose Configuration

```yaml
# docker-compose.yml
version: '3.8'

services:
  biometrics:
    build: .
    ports:
      - "8080:8080"
      - "6060:6060"
    environment:
      - REDIS_ADDR=redis:6379
      - PROFILING_ENABLED=true
    depends_on:
      - redis
      
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes

volumes:
  redis_data:
```

---

## Related Documentation

- [Security Documentation](./SECURITY.md)
- [API Reference](./docs/api/)
- [Deployment Guide](./docs/deployment/)
- [Cache Package](./pkg/cache/README.md)
- [Benchmarks](./benchmarks/)

---

*Document Version: 1.0.0*  
*Last Updated: February 2026*  
*Compliant with Enterprise Practices Feb 2026*
