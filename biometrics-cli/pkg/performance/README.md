# Performance Monitoring Package

Comprehensive performance monitoring and profiling using Go's built-in pprof infrastructure.

## Features

- **CPU Profiling**: Identify CPU-intensive code paths
- **Memory Profiling**: Track heap allocations and memory usage
- **Block Profiling**: Detect goroutine blocking operations
- **Mutex Profiling**: Find mutex contention points
- **Goroutine Profiling**: Monitor goroutine count and leaks
- **HTTP Profiling Server**: Real-time profiling endpoints
- **Automated Snapshots**: Capture profile snapshots on-demand
- **GC Statistics**: Track garbage collection performance

## Installation

```bash
go get go.uber.org/zap
```

## Usage

### Basic Profiler

```go
package main

import (
    "log"
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
    
    // Start profiling
    if err := profiler.Start(); err != nil {
        log.Fatal(err)
    }
    
    // Your application code here
    
    // Stop and write profiles
    if err := profiler.Stop(); err != nil {
        log.Fatal(err)
    }
}
```

### HTTP Profiling Server

Access profiling endpoints at `http://localhost:6060/debug/pprof/`:

- `/debug/pprof/` - Profile index
- `/debug/pprof/cpu` - CPU profile
- `/debug/pprof/heap` - Heap profile
- `/debug/pprof/block` - Block profile
- `/debug/pprof/goroutine` - Goroutine profile
- `/debug/pprof/mutex` - Mutex profile

### Capture Snapshots

```go
ctx := context.Background()
snapshots, err := profiler.CaptureSnapshot(ctx)
if err != nil {
    log.Fatal(err)
}

// snapshots["memory"] -> path to memory profile
// snapshots["goroutine"] -> path to goroutine profile
// snapshots["gc_stats"] -> path to GC statistics
```

### Analyze Profiles

```bash
# CPU Profile
go tool pprof cpu.prof

# Memory Profile (interactive)
go tool pprof mem.prof

# Generate SVG flame graph
go tool pprof -svg cpu.prof > cpu.svg

# HTTP endpoint
go tool pprof http://localhost:6060/debug/pprof/cpu
```

## Configuration

```go
config := performance.ProfilerConfig{
    EnableCPU:        true,              // Enable CPU profiling
    EnableMemory:     true,              // Enable memory profiling
    EnableBlock:      true,              // Enable block profiling
    EnableMutex:      true,              // Enable mutex profiling
    EnableGoroutine:  true,              // Enable goroutine profiling
    ProfileDir:       "./profiles",      // Directory for profile files
    HTTPAddr:         "localhost:6060",  // HTTP server address
    EnableHTTP:       true,              // Enable HTTP server
    ProfileDuration:  30 * time.Second,  // Auto-profile duration
    BlockProfileRate: 1,                 // Block profile rate
    MutexProfileFraction: 1,             // Mutex profile fraction
}
```

## ProfilingManager

Advanced profiling with manual control:

```go
manager, err := performance.NewProfilingManager(config, logger)

// Start specific profile
err = manager.StartProfile(ctx, performance.ProfileTypeCPU, 30*time.Second)

// Stop specific profile
err = manager.StopProfile(performance.ProfileTypeCPU)

// Capture all profiles
snapshots, err := manager.CaptureAllProfiles(ctx)

// Get metrics
metrics := manager.GetMetrics()
```

## Runtime Metrics

```go
metrics := profiler.ReportMetrics()

fmt.Printf("Goroutines: %d\n", metrics.Goroutines)
fmt.Printf("CPUs: %d\n", metrics.CPUs)
fmt.Printf("Memory: %d KB\n", metrics.Alloc/1024)
fmt.Printf("GC Count: %d\n", metrics.NumGC)
fmt.Printf("GC CPU Fraction: %.2f%%\n", metrics.GCCPUFraction*100)
```

## Best Practices

1. **Enable in Production**: Use HTTP profiling server for production monitoring
2. **Set Profile Rates**: Configure block and mutex profile rates appropriately
3. **Monitor GC**: Track GC statistics for memory optimization
4. **Regular Snapshots**: Capture snapshots during high-load periods
5. **Analyze Trends**: Compare profiles over time to identify regressions

## Files

- `pprof.go` - Main profiler implementation
- `profiling.go` - Profiling manager with advanced features
- `handlers.go` - HTTP handlers and middleware
- `pprof_test.go` - Comprehensive tests

## Examples

See tests in `pprof_test.go` for complete usage examples.

## Related

- [Redis Cache](../cache/) - Distributed caching layer
- [Benchmarks](../../benchmarks/) - Performance benchmarks
