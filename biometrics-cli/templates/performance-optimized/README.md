# Performance-Optimized Agent Workflows

This directory contains performance-optimized workflow templates for BIOMETRICS.

## Performance Guidelines

### Model Selection for Performance

| Task Type | Recommended Model | Expected Latency |
|-----------|------------------|------------------|
| Fast checks | minimax-m2.5-free | <5s |
| Code exploration | minimax-m2.5-free | <10s |
| Documentation | minimax-m2.5-free | <15s |
| Complex coding | qwen/qwen3.5-397b-a17b | 70-90s |
| Deep analysis | kimi-k2.5-free | 30-60s |

### Optimization Strategies

1. **Parallel Execution**: Use MiniMax for parallel tasks (can run 10x)
2. **Caching**: Cache results to avoid redundant work
3. **Batch Processing**: Group similar tasks together
4. **Model Routing**: Route to fastest available model

### Quick Start

```bash
# Run performance-optimized workflow
biometrics run --template performance-review
```
