# CLI Benchmark Suite

**Purpose:** Performance benchmarking and testing for the biometrics CLI application

## Overview

This directory contains Go-based benchmark tests used to measure and validate the performance of the biometrics CLI commands. Benchmarks are essential for maintaining optimal performance as the codebase grows.

## Files

| File | Purpose |
|------|---------|
| `benchmark_suite.go` | Main benchmark test suite with comprehensive performance tests |
| `benchmark_test.go` | Individual benchmark test cases |

## Running Benchmarks

### Run All Benchmarks
```bash
cd /Users/jeremy/dev/BIOMETRICS/biometrics-cli
go test -bench=. ./benchmarks/ -benchmem
```

### Run Specific Benchmark
```bash
go test -bench=BenchmarkName ./benchmarks/ -benchmem
```

### Run with Memory Profiling
```bash
go test -bench=. ./benchmarks/ -benchmem -cpuprofile=cpu.out -memprofile=mem.out
```

## Benchmark Categories

### Command Performance
- CLI startup time
- Command execution latency
- Output generation speed

### Memory Usage
- Heap allocations
- Memory footprint
- Garbage collection impact

### Concurrent Operations
- Parallel command execution
- Rate limiting effectiveness
- Worker pool performance

## Interpreting Results

### Key Metrics

| Metric | Good | Warning | Critical |
|--------|------|---------|----------|
| ns/op | <1000 | 1000-5000 | >5000 |
| B/op | <1024 | 1024-4096 | >4096 |
| allocs/op | <10 | 10-50 | >50 |

### Example Output
```
BenchmarkCommandExecution-8           10000             120 ns/op            48 B/op          3 allocs/op
```

## Integration with CI/CD

Benchmarks run automatically on:
- Every push to main
- Weekly performance regression check
- Pre-release validation

## Best Practices

1. **Run benchmarks locally** before committing significant changes
2. **Track trends** over time using benchmark history
3. **Set baselines** for critical performance indicators
4. **Alert on regression** if performance degrades >10%

## Related Documentation

- [Performance Testing Guide](../docs/performance.md)
- [CLI Architecture](../docs/architecture.md)
- [Optimization Guidelines](../docs/optimization.md)

## Maintenance

- Update benchmarks when adding new commands
- Review benchmark results weekly
- Adjust thresholds based on hardware capabilities
