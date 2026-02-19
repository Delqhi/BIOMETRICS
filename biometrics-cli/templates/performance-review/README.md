# Performance Review Workflow Template

## Overview

The Performance Review workflow template provides comprehensive performance analysis and optimization recommendations for software applications. This template leverages AI-powered analysis to identify performance bottlenecks, memory leaks, inefficient algorithms, and other performance-related issues in codebases.

The workflow is designed to analyze both frontend and backend performance characteristics, examining areas such as database queries, algorithm complexity, memory usage, network requests, and rendering performance. It produces detailed reports with specific, actionable recommendations for improvement.

This template is essential for organizations seeking to:
- Optimize application performance
- Reduce infrastructure costs
- Improve user experience
- Identify scaling bottlenecks
- Meet performance SLAs

## Purpose

The primary purpose of the Performance Review template is to:

1. **Identify Bottlenecks** - Detect performance issues in code and infrastructure
2. **Analyze Complexity** - Evaluate algorithmic complexity and efficiency
3. **Profile Resources** - Analyze memory, CPU, and network usage patterns
4. **Recommend Solutions** - Provide specific, actionable optimization recommendations
5. **Benchmark Progress** - Track performance metrics over time

### Key Use Cases

- **Pre-deployment Review** - Ensure new code meets performance standards
- **Performance Optimization** - Identify and fix performance regressions
- **Scaling Analysis** - Prepare for increased load
- **Incident Response** - Investigate performance-related outages
- **Architecture Review** - Evaluate performance implications of design decisions

## Input Parameters

The Performance Review template accepts the following input parameters:

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `repository` | string | Yes | - | Git repository URL |
| `branch` | string | No | main | Branch to analyze |
| `focus_areas` | array | No | all | Areas to focus on (database, algorithm, memory, network, frontend) |
| `baseline` | object | No | - | Baseline metrics for comparison |
| `test_data` | string | No | - | Path to test data for benchmarking |
| `load_profile` | object | No | - | Expected load profile for scaling analysis |

### Input Examples

```yaml
# Example 1: Full performance review
inputs:
  repository: https://github.com/org/app
  branch: main
  focus_areas:
    - database
    - algorithm
    - memory
    - network

# Example 2: Database-focused review
inputs:
  repository: https://github.com/org/app
  branch: feature/new-api
  focus_areas:
    - database

# Example 3: With baseline comparison
inputs:
  repository: https://github.com/org/app
  branch: main
  baseline:
    response_time_p50: 120
    response_time_p99: 500
    throughput: 1000
    error_rate: 0.01

# Example 4: Frontend performance
inputs:
  repository: https://github.com/org/frontend
  focus_areas:
    - frontend
    - network
```

## Output Results

The template produces comprehensive performance analysis reports:

| Output | Type | Description |
|--------|------|-------------|
| `report` | object | Detailed performance report |
| `score` | number | Overall performance score (0-100) |
| `bottlenecks` | array | List of identified bottlenecks |
| `recommendations` | array | Prioritized optimization recommendations |
| `benchmarks` | object | Benchmark results |

### Output Report Structure

```json
{
  "report": {
    "timestamp": "2026-02-19T10:30:00Z",
    "repository": "https://github.com/org/app",
    "branch": "main",
    "analysis_duration_seconds": 345
  },
  "score": 72,
  "findings": {
    "critical": 3,
    "high": 8,
    "medium": 15,
    "low": 22
  },
  "bottlenecks": [
    {
      "id": "PERF-001",
      "severity": "critical",
      "category": "database",
      "title": "N+1 Query Problem in User Service",
      "location": "src/services/user_service.go:145",
      "impact": "Each user request triggers 50+ database queries",
      "current_metrics": {
        "queries_per_request": 52,
        "avg_query_time_ms": 45,
        "total_time_ms": 2340
      }
    }
  ],
  "recommendations": [
    {
      "id": "PERF-001-OPT",
      "priority": 1,
      "title": "Implement DataLoader for Batch Loading",
      "effort": "medium",
      "estimated_improvement": "85% reduction in query count",
      "implementation": {
        "steps": [
          "1. Install dataloader package",
          "2. Create userDataLoader in request context",
          "3. Replace direct queries with loader.Batch",
          "4. Add tests for batched loading"
        ]
      }
    }
  ],
  "benchmarks": {
    "api_response_time": {
      "p50_ms": 145,
      "p95_ms": 320,
      "p99_ms": 890
    },
    "database_queries": {
      "avg_per_request": 12,
      "p99_per_request": 45
    },
    "memory_usage": {
      "heap_mb": 256,
      "rss_mb": 512
    }
  },
  "comparison": {
    "baseline": {
      "response_time_p50": 120,
      "current": 145,
      "change_percent": 20.8
    }
  }
}
```

## Workflow Steps

### Step 1: Fetch Code

**ID:** `fetch-code`  
**Type:** agent  
**Timeout:** 5 minutes  
**Provider:** opencode-zen

Retrieves code from the repository for analysis.

### Step 2: Analyze Database Queries

**ID:** `analyze-queries`  
**Type:** agent  
**Timeout:** 15 minutes  
**Provider:** opencode-zen

Identifies database query issues:
- N+1 query problems
- Missing indexes
- Inefficient joins
- Unoptimized queries

### Step 3: Analyze Algorithm Complexity

**ID:** `analyze-algorithms`  
**Type:** agent  
**Timeout:** 15 minutes  
**Provider:** opencode-zen

Evaluates algorithmic efficiency:
- Time complexity analysis
- Space complexity analysis
- Inefficient loops
- Suboptimal data structures

### Step 4: Analyze Memory Usage

**ID:** `analyze-memory`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Detects memory issues:
- Memory leaks
- Excessive allocations
- Unbounded caches
- Large object copying

### Step 5: Network Analysis

**ID:** `analyze-network`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Evaluates network performance:
- Unnecessary API calls
- Missing request caching
- Large payload transfers
- Connection pooling

### Step 6: Generate Report

**ID:** `generate-report`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Compiles comprehensive performance report.

## Usage Examples

### CLI Usage

```bash
# Full performance review
biometrics workflow run performance-review \
  --repository https://github.com/org/app \
  --branch main

# Database-focused review
biometrics workflow run performance-review \
  --repository https://github.com/org/app \
  --branch feature/new-api \
  --focus_areas '["database"]'

# With baseline comparison
biometrics workflow run performance-review \
  --repository https://github.com/org/app \
  --branch main \
  --baseline '{"response_time_p50": 120, "response_time_p99": 500}'
```

### Programmatic Usage

```go
import "github.com/biometrics/biometrics-cli/pkg/workflows"

engine := workflows.NewWorkflowEngine("./templates")
template, _ := engine.LoadTemplate("performance-review")

instance, _ := engine.CreateInstance(template, map[string]interface{}{
    "repository": "https://github.com/org/app",
    "branch":    "main",
    "focus_areas": []string{"database", "algorithm", "memory"},
})

result, err := engine.Execute(context.Background(), instance)
fmt.Printf("Performance Score: %d\n", result["score"])
```

## Configuration

### Analysis Configuration

```yaml
inputs:
  focus_areas:
    - database      # Query optimization
    - algorithm     # Complexity analysis
    - memory        # Memory leak detection
    - network       # API call optimization
    - frontend      # Frontend performance
```

### Threshold Configuration

```yaml
options:
  thresholds:
    response_time_p99_ms: 500
    queries_per_request_max: 10
    memory_mb_max: 512
    cpu_percent_max: 80
```

## Troubleshooting

### Common Issues

#### Issue: Analysis Timeout

**Symptom:** Workflow times out on large codebase

**Solution:**
```yaml
options:
  timeout: 60m  # Increase timeout

inputs:
  focus_areas:  # Limit to specific areas
    - database
```

#### Issue: Missing Dependencies

**Symptom:** Cannot analyze without build artifacts

**Solution:**
```bash
# Ensure project is built
cd repository
npm install  # or pip install, go mod download, etc.
```

## Best Practices

### 1. Run Regularly

Schedule periodic performance reviews to catch regressions early.

### 2. Compare to Baseline

Always compare against baseline metrics to track improvements.

### 3. Focus on Critical Areas

Prioritize database and algorithm optimizations for maximum impact.

### 4. Track Progress

Maintain historical performance metrics to measure improvement over time.

## Related Templates

- **Code Review** (`code-review/`) - General code quality checks
- **Security Audit** (`security-audit/`) - Security-focused analysis
- **Deployment** (`deployment/`) - Deployment with performance validation

---

**Template Version:** 1.0.0  
**Author:** BIOMETRICS Team  
**Category:** Performance  
**Tags:** performance, optimization, bottleneck, profiling, benchmarking

*Last Updated: February 2026*
