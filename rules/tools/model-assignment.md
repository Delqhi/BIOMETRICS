# Model Assignment Guide - BIOMETRICS Project

**Version:** 1.0  
**Created:** 2026-02-20  
**Status:** ACTIVE  
**Project:** BIOMETRICS - AI Automation System  

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Available Models](#2-available-models)
3. [Model Characteristics](#3-model-characteristics)
4. [Model Selection Matrix](#4-model-selection-matrix)
5. [Parallel Execution Rules](#5-parallel-execution-rules)
6. [Fallback Chains](#6-fallback-chains)
7. [Cost Optimization](#7-cost-optimization)
8. [Performance Tuning](#8-performance-tuning)
9. [Monitoring \& Metrics](#9-monitoring--metrics)
10. [Implementation Guidelines](#10-implementation-guidelines)
11. [Troubleshooting](#11-troubleshooting)
12. [Appendix](#12-appendix)

---

## 1. Executive Summary

This document serves as the definitive guide for model assignment within the BIOMETRICS project. It establishes clear rules for selecting the appropriate AI model based on task type, performance requirements, cost considerations, and parallel execution constraints.

### Purpose

The primary objectives of this model assignment framework are:

- **Optimize Performance**: Assign the best-suited model for each task type to maximize output quality and efficiency
- **Control Costs**: Leverage FREE models whenever possible while maintaining quality standards
- **Ensure Reliability**: Implement robust fallback chains to handle model failures gracefully
- **Maximize Throughput**: Enable parallel execution within defined limits to increase productivity

### Guiding Principles

1. **FREE First**: Always prefer free models (OpenCode ZEN, NVIDIA NIM) over paid alternatives
2. **Fit for Purpose**: Use specialized models for their strengths rather than forcing one model to handle all tasks
3. **Parallel Efficiency**: Maximize concurrent work within rate limit constraints
4. **Graceful Degradation**: Implement multiple fallback levels to ensure task completion

### Document Maintenance

This document should be reviewed and updated whenever:
- New models are added to the system
- Performance characteristics change
- New use cases are identified
- Pricing models are modified

---

## 2. Available Models

### 2.1 NVIDIA NIM Models

#### Qwen 3.5 397B

| Property | Value |
|----------|-------|
| **Provider** | NVIDIA NIM |
| **Model ID** | `qwen/qwen3.5-397b-a17b` |
| **Context Window** | 262,144 tokens |
| **Output Limit** | 32,768 tokens |
| **Latency** | 70-90 seconds |
| **RPM Limit** | 40 requests/minute |
| **Cost** | FREE (NVIDIA Free Tier) |

**Description**: Qwen 3.5 397B is the most capable code generation model available in our stack. It excels at complex reasoning, architectural design, and sophisticated implementation tasks. However, its high latency (70-90 seconds) and strict rate limits (40 RPM) require careful orchestration.

**Best For**:
- Complex code implementation
- Architecture design
- Multi-file refactoring
- Bug fixing in complex systems
- Security analysis

**Limitations**:
- High latency makes it unsuitable for quick lookups
- Rate limit of 40 RPM requires queuing
- Should not be used for simple documentation tasks

#### Available NVIDIA NIM Models

| Model | Context | Output | Latency | Use Case |
|-------|---------|--------|---------|----------|
| `qwen/qwen3.5-397b-a17b` | 262K | 32K | 70-90s | Code, Architecture |
| `qwen2.5-coder-32b` | 128K | 8K | 10-20s | Fast code completion |
| `qwen2.5-coder-7b` | 128K | 8K | 3-5s | Quick code snippets |
| `meta/llama-3.3-70b-instruct` | 128K | 8K | 15-25s | General purpose |
| `mistralai/mistral-large-3-675b-instruct-2512` | 128K | 32K | 20-30s | Complex reasoning |

### 2.2 OpenCode ZEN Models

#### Kimi K2.5

| Property | Value |
|----------|-------|
| **Provider** | OpenCode ZEN (via NVIDIA NIM) |
| **Model ID** | `moonshotai/kimi-k2.5` |
| **Context Window** | 1,048,576 tokens (1M) |
| **Output Limit** | 65,536 tokens (64K) |
| **Latency** | 10-20 seconds |
| **RPM Limit** | Variable (depends on provider) |
| **Cost** | FREE |

**Description**: Kimi K2.5 offers an exceptional 1M token context window, making it ideal for processing large documents, analyzing extensive codebases, and performing deep research tasks. Its moderate latency (10-20s) provides a good balance between capability and speed.

**Best For**:
- Deep research and analysis
- Large document processing
- Context-heavy conversations
- Multi-file code analysis
- Long-form content generation

**Limitations**:
- Not as specialized for code as Qwen
- Still has latency considerations for real-time applications

#### MiniMax M2.5 Free

| Property | Value |
|----------|-------|
| **Provider** | OpenCode ZEN |
| **Model ID** | `opencode/minimax-m2.5-free` |
| **Context Window** | 200,000 tokens |
| **Output Limit** | 32,000 tokens |
| **Latency** | 5-10 seconds |
| **RPM Limit** | High (can run 10+ in parallel) |
| **Cost** | 100% FREE |

**Description**: MiniMax M2.5 is our workhorse for high-throughput tasks. With blazing fast latency (5-10 seconds) and the ability to run 10+ parallel instances, it's perfect for documentation generation, research tasks, and any operation where speed matters more than maximum capability.

**Best For**:
- Documentation generation
- Research and web searches
- File creation and editing
- Quick lookups
- Parallel processing workloads
- Simple transformations

**Limitations**:
- Not suitable for complex code architecture
- May miss subtle edge cases in sophisticated implementations

### 2.3 Alternative Models

#### Streamlake KAT Coder Pro

| Property | Value |
|----------|-------|
| **Provider** | Streamlake |
| **Model ID** | `kat-coder-pro-v1` |
| **Context Window** | 2,000,000 tokens |
| **Output Limit** | 128,000 tokens |
| **Latency** | 15-25 seconds |
| **Cost** | $0.50/1M input, $1.50/1M output |

**Description**: Streamlake offers exceptional context limits but at a cost. Use as a fallback when other models are unavailable or when processing extremely large codebases.

#### XiaoMi MIMO

| Property | Value |
|----------|-------|
| **Provider** | XiaoMi |
| **Model ID** | `mimo-v2-turbo` |
| **Context Window** | 1,500,000 tokens |
| **Output Limit** | 100,000 tokens |
| **Latency** | 10-15 seconds |
| **Cost** | $0.70/1M input, $2.10/1M output |

**Description**: MIMO offers strong multimodal capabilities and large context windows. Good alternative for complex tasks requiring vision input.

#### OpenRouter Models (Backup)

| Model | Context | Output | Cost | Use Case |
|-------|---------|--------|------|----------|
| `grok-code` | 2M | 131K | FREE | Code generation |
| `glm-4.7-free` | 1M | 65K | FREE | General purpose |

---

## 3. Model Characteristics

### 3.1 Context Window Analysis

The context window determines how much information the model can consider in a single request. Here's a detailed breakdown:

#### Context Window Comparison

```
Model                   Context Window    Practical Limit    Best Use Case
─────────────────────────────────────────────────────────────────────────────
Qwen 3.5 397B          262,144 tokens    ~50,000 tokens    Complex code
Kimi K2.5              1,048,576 tokens   ~200,000 tokens   Research
MiniMax M2.5          200,000 tokens      ~40,000 tokens    Fast docs
Streamlake KAT         2,000,000 tokens   ~300,000 tokens   Large codebases
XiaoMi MIMO            1,500,000 tokens   ~250,000 tokens   Multimodal
```

#### Context Management Strategies

**For Large Context Requirements**:
1. Use Kimi K2.5 for 100K+ token inputs
2. Implement chunking for Streamlake when needed
3. Consider summarization before passing to smaller models

**For Standard Context**:
1. Qwen 3.5 for complex code (keep under 50K tokens)
2. MiniMax for quick tasks (keep under 40K tokens)

### 3.2 Latency Analysis

Latency directly impacts user experience and throughput. Here's the breakdown:

#### Latency by Model

```
Model                   Cold Start    Per-Token    Total (1K)    Total (4K)
──────────────────────────────────────────────────────────────────────────────
Qwen 3.5 397B          60-70s        0.01s        70-80s        85-95s
Kimi K2.5              5-10s         0.003s       8-13s         15-22s
MiniMax M2.5           2-3s          0.002s       4-5s          8-10s
Qwen 2.5 Coder 32B    5-8s          0.003s       8-11s         15-20s
```

#### Latency Optimization Techniques

1. **Pre-warming**: Keep a warm instance for critical paths
2. **Streaming**: Use streaming for perceived latency reduction
3. **Async Processing**: Queue requests and process in background
4. **Smart Routing**: Route to fastest model for simple tasks

### 3.3 Cost Analysis

#### Cost Per Million Tokens

```
Model                   Input Cost      Output Cost      Daily Budget (1000 req)
─────────────────────────────────────────────────────────────────────────────────
Qwen 3.5 397B          FREE            FREE              $0.00
Kimi K2.5              FREE            FREE              $0.00
MiniMax M2.5           FREE            FREE              $0.00
NVIDIA Fallbacks       FREE            FREE              $0.00
Streamlake KAT         $0.50           $1.50             $50-100
XiaoMi MIMO            $0.70           $2.10             $70-140
```

#### Cost Optimization Rules

1. **Always start with FREE models**
2. **Use paid models only as last resort**
3. **Implement request caching to reduce costs**
4. **Monitor token usage per model**

### 3.4 Strengths and Weaknesses Matrix

#### Qwen 3.5 397B

| Strengths | Weaknesses |
|-----------|------------|
| Best code generation quality | High latency (70-90s) |
| Excellent architectural thinking | Rate limited (40 RPM) |
| Strong reasoning capabilities | Expensive when rate limited |
| Good at complex refactoring | Not suitable for simple tasks |
| Security-aware outputs | Requires careful orchestration |

#### Kimi K2.5

| Strengths | Weaknesses |
|-----------|------------|
| 1M token context window | Less specialized for code |
| Fast for its capability | Not as fast as MiniMax |
| Excellent research capabilities | Moderate latency |
| Good for long文档 processing | Context switching overhead |

#### MiniMax M2.5

| Strengths | Weaknesses |
|-----------|------------|
| Blazing fast (5-10s) | Not for complex architecture |
| Can run 10+ parallel | May miss edge cases |
| 100% FREE | Limited reasoning depth |
| Perfect for documentation | Short context (200K) |
| Great for simple transformations | Not for security-critical tasks |

---

## 4. Model Selection Matrix

### 4.1 Primary Selection Table

| Task Type | Best Model | Alternative | Avoid | Max Parallel |
|-----------|-----------|-------------|-------|--------------|
| **Code Implementation** | Qwen 3.5 397B | Qwen 2.5 Coder 32B | MiniMax | 1 |
| **Code Review** | Qwen 3.5 397B | Kimi K2.5 | MiniMax | 1 |
| **Architecture Design** | Qwen 3.5 397B | Kimi K2.5 | MiniMax | 1 |
| **Bug Fixing** | Qwen 3.5 397B | Qwen 2.5 Coder 32B | - | 1 |
| **Security Analysis** | Qwen 3.5 397B | - | MiniMax | 1 |
| **Documentation** | MiniMax M2.5 | Kimi K2.5 | Qwen | 10 |
| **Research** | MiniMax M2.5 | Kimi K2.5 | Qwen | 10 |
| **Web Search** | MiniMax M2.5 | - | Qwen | 10 |
| **File Creation** | MiniMax M2.5 | - | Qwen | 10 |
| **Deep Analysis** | Kimi K2.5 | Qwen 3.5 397B | MiniMax | 1 |
| **Large Document** | Kimi K2.5 | Streamlake | Qwen | 1 |
| **Multi-file Analysis** | Kimi K2.5 | Qwen 3.5 397B | MiniMax | 1 |
| **Simple Transformation** | MiniMax M2.5 | - | Qwen | 10 |
| **API Design** | Qwen 3.5 397B | Kimi K2.5 | MiniMax | 1 |
| **Testing** | Qwen 3.5 397B | MiniMax | - | 1 |

### 4.2 Detailed Use Case Mapping

#### Code Generation Tasks

```
Task: New Feature Implementation
────────────────────────────────
Decision Tree:
├── Complexity: High → Qwen 3.5 397B
├── Complexity: Medium → Qwen 2.5 Coder 32B
└── Complexity: Low → MiniMax M2.5

Task: Bug Fix
────────────────────────────────
Decision Tree:
├── Severity: Critical → Qwen 3.5 397B
├── Severity: High → Qwen 2.5 Coder 32B
└── Severity: Low/Medium → MiniMax M2.5

Task: Code Refactoring
────────────────────────────────
Decision Tree:
├── Scope: Multi-file → Qwen 3.5 397B
├── Scope: Single file → Qwen 2.5 Coder 32B
└── Scope: Simple → MiniMax M2.5
```

#### Documentation Tasks

```
Task: Technical Documentation
────────────────────────────────
Decision Tree:
├── Length: Long-form → MiniMax M2.5 (fast output)
├── Length: Detailed → Kimi K2.5 (for accuracy)
└── Length: Quick → MiniMax M2.5

Task: API Documentation
────────────────────────────────
Decision Tree:
├── Complexity: High → Qwen 3.5 397B
├── Complexity: Standard → MiniMax M2.5
└── Generation: Bulk → MiniMax M2.5 (parallel)

Task: README Creation
────────────────────────────────
Decision Tree:
└── Always → MiniMax M2.5 (fast, good enough)
```

#### Research and Analysis Tasks

```
Task: Web Research
────────────────────────────────
Decision Tree:
└── Always → MiniMax M2.5 (speed critical)

Task: Codebase Analysis
────────────────────────────────
Decision Tree:
├── Size: Large (>100 files) → Kimi K2.5
├── Size: Medium → Qwen 3.5 397B
└── Size: Small → MiniMax M2.5

Task: Security Audit
────────────────────────────────
Decision Tree:
└── Always → Qwen 3.5 397B (quality critical)
```

### 4.3 Decision Flowchart

```
START
  │
  ▼
Is this a CODE task? ────YES───► Is complexity HIGH?
  │                                    │
  NO                                   YES
  │                                    │
  ▼                                    ▼
Is this RESEARCH?               Is this SECURITY?
  │                                    │
  YES                            YES        NO
  │                                    │
  ▼                                 ▼          ▼
Use MiniMax             Use Qwen 3.5     Use Qwen 2.5
(10 parallel)           397B (1 only)    Coder 32B
```

---

## 5. Parallel Execution Rules

### 5.1 Rate Limit Constraints

#### NVIDIA NIM Limits

| Model | RPM Limit | Daily Limit | Concurrent Max |
|-------|-----------|-------------|----------------|
| Qwen 3.5 397B | 40 | 57,600 | 1 |
| Qwen 2.5 Coder 32B | 40 | 57,600 | 2 |
| Llama 3.3 70B | 40 | 57,600 | 2 |
| Mistral Large 3 | 40 | 57,600 | 2 |

#### OpenCode ZEN Limits

| Model | RPM Limit | Daily Limit | Concurrent Max |
|-------|-----------|-------------|----------------|
| MiniMax M2.5 | Unlimited | Unlimited | 10+ |
| Kimi K2.5 | 100 | 144,000 | 3 |
| ZEN Uncensored | Unlimited | Unlimited | 5 |

### 5.2 Parallel Execution Matrix

```
Model                    Max Concurrent    Reason
────────────────────────────────────────────────────
Qwen 3.5 397B           1                 Rate limit
Qwen 2.5 Coder 32B      2                 Rate limit
Llama 3.3 70B           2                 Rate limit
Mistral Large 3         2                 Rate limit
MiniMax M2.5            10                No limit
Kimi K2.5               3                 Moderate limit
ZEN Uncensored         5                 No limit
```

### 5.3 Execution Strategy

#### Rule 1: Never Block on Qwen

```python
# BAD: Blocking wait for Qwen
result = await call_qwen(task)  # 70-90s wait
next_task()

# GOOD: Fire and forget with callback
task_id = await queue_qwen(task)
await register_callback(task_id, handle_result)
continue_with_other_work()
```

#### Rule 2: Batch MiniMax Requests

```python
# GOOD: Batch multiple small tasks
tasks = [create_task(i) for i in range(10)]
results = await asyncio.gather(*tasks)  # Parallel!
```

#### Rule 3: Implement Request Queuing

```python
class ModelQueue:
    def __init__(self, model, max_concurrent):
        self.model = model
        self.semaphore = asyncio.Semaphore(max_concurrent)
        self.queue = asyncio.Queue()
    
    async def submit(self, task):
        async with self.semaphore:
            return await self.model.execute(task)
```

### 5.4 Priority-Based Execution

#### Priority Levels

| Priority | Model Assignment | Use Case |
|----------|-----------------|----------|
| P0 - Critical | Qwen 3.5 (1 max) | Security, Production bugs |
| P1 - High | Qwen 2.5 / Kimi | Feature development |
| P2 - Medium | Any available | General tasks |
| P3 - Low | MiniMax (any) | Documentation, research |

---

## 6. Fallback Chains

### 6.1 Primary Fallback Chains

#### Code Generation Chain

```
PRIMARY: Qwen 3.5 397B
    │
    ├─ FAILURE: Rate limit (429)
    │   └─► WAIT 60s + RETRY
    │
    ├─ FAILURE: Timeout (120s)
    │   └─► RETRY (max 3x)
    │
    ├─ FAILURE: Model unavailable
    │   └─► FALLBACK 1: Qwen 2.5 Coder 32B
    │           │
    │           └─ FAILURE
    │               └─► FALLBACK 2: Llama 3.3 70B
    │                       │
    │                       └─ FAILURE
    │                           └─► FALLBACK 3: ZEN Big Pickle
    │
    └─ CENSORSHIP DETECTED
        └─► SWITCH: ZEN Big Pickle (uncensored)
```

#### Documentation Chain

```
PRIMARY: MiniMax M2.5
    │
    ├─ FAILURE: Timeout
    │   └─► RETRY (max 3x)
    │
    ├─ FAILURE: Unavailable
    │   └─► FALLBACK 1: Kimi K2.5
    │           │
    │           └─ FAILURE
    │               └─► FALLBACK 2: ZEN Uncensored
    │
    └─ QUALITY TOO LOW
        └─► UPGRADE: Kimi K2.5
```

#### Research Chain

```
PRIMARY: MiniMax M2.5 (web search)
    │
    ├─ FAILURE: Search unavailable
    │   └─► FALLBACK 1: Tavily MCP
    │           │
    │           └─ FAILURE
    │               └─► FALLBACK 2: Context7 Docs
    │
    └─ ANALYSIS NEEDED
        └─► CHAIN: MiniMax for search + Kimi for analysis
```

### 6.2 Censorship Handling

#### Detection Triggers

| Trigger Phrase | Action |
|----------------|--------|
| "I cannot generate that" | Switch to ZEN |
| "Content Policy violation" | Switch to ZEN |
| "Not allowed" | Switch to ZEN |
| "Cannot assist with this" | Switch to ZEN |

#### ZEN Fallback Chain

```
DETECTED: Censorship
    │
    ▼
Is it CODE? ────YES───► Use: ZEN Big Pickle
  │
  NO
  │
  ▼
Is it CREATIVE? ────YES───► Use: ZEN Uncensored
  │
  NO
  │
  ▼
Use: ZEN Code
```

### 6.3 Rate Limit Handling

#### HTTP 429 Response

```python
async def handle_rate_limit(response, retry_count):
    if response.status == 429:
        wait_time = int(response.headers.get('Retry-After', 60))
        logger.warning(f"Rate limited, waiting {wait_time}s")
        await asyncio.sleep(wait_time)
        
        if retry_count < 3:
            return await retry_with_backoff(retry_count + 1)
        else:
            return await fallback_model()
```

#### Exponential Backoff

```
Attempt 1: Wait 60s
Attempt 2: Wait 120s  
Attempt 3: Switch to fallback model
```

---

## 7. Cost Optimization

### 7.1 FREE-First Philosophy

#### Priority List

```
PRIORITY 1 (Always): OpenCode ZEN Models
├── MiniMax M2.5 - Primary for docs/research
├── Kimi K2.5 - Large context tasks
└── ZEN Big Pickle - Censorship fallback

PRIORITY 2 (When needed): NVIDIA NIM Free Tier
├── Qwen 3.5 397B - Code architecture
├── Qwen 2.5 Coder 32B - Fast code
└── Llama 3.3 70B - General fallback

PRIORITY 3 (Last resort): Paid Models
├── Streamlake KAT - Large context
├── XiaoMi MIMO - Multimodal
└── NEVER exceed $10/day
```

### 7.2 Cost Monitoring

#### Daily Budget Rules

| Model | Daily Budget | Alert Threshold |
|-------|-------------|----------------|
| All FREE | $0.00 | N/A |
| Streamlake | $20.00 | $15.00 |
| XiaoMi | $15.00 | $10.00 |

#### Cost Tracking Implementation

```python
class CostTracker:
    def __init__(self):
        self.daily_spend = 0
        self.model_costs = {}
    
    async def execute_with_tracking(self, model, task):
        start_time = time.time()
        result = await model.execute(task)
        duration = time.time() - start_time
        
        cost = self.calculate_cost(model, duration)
        self.daily_spend += cost
        self.model_costs[model.name] = cost
        
        if self.daily_spend > ALERT_THRESHOLD:
            await alert_team(f"Near daily budget: ${self.daily_spend}")
        
        return result
    
    def calculate_cost(self, model, duration):
        rates = {
            'streamlake': 0.001 / 60,  # per second
            'xiaomi': 0.0014 / 60,
        }
        return rates.get(model.name, 0) * duration
```

### 7.3 Caching Strategies

#### Request Caching

```python
class RequestCache:
    def __init__(self, ttl=3600):
        self.cache = {}
        self.ttl = ttl
    
    async def get_or_execute(self, key, executor):
        if key in self.cache:
            cached, timestamp = self.cache[key]
            if time.time() - timestamp < self.ttl:
                logger.info(f"Cache hit for {key}")
                return cached
        
        result = await executor()
        self.cache[key] = (result, time.time())
        return result
```

#### Cache Key Strategies

| Task Type | Cache Key | TTL |
|-----------|-----------|-----|
| Documentation | `doc:{file_hash}` | 24h |
| Research | `search:{query_hash}` | 1h |
| Code Review | `review:{file_hash}:{line_hash}` | 4h |
| File Analysis | `analyze:{file_hash}` | 12h |

---

## 8. Performance Tuning

### 8.1 Timeout Configuration

#### Model-Specific Timeouts

| Model | Recommended Timeout | Absolute Maximum |
|-------|--------------------|--------------------|
| Qwen 3.5 397B | 120,000ms (120s) | 180,000ms |
| Kimi K2.5 | 60,000ms (60s) | 90,000ms |
| MiniMax M2.5 | 30,000ms (30s) | 60,000ms |
| Qwen 2.5 Coder 32B | 45,000ms (45s) | 60,000ms |

#### Timeout Implementation

```python
class ModelClient:
    def __init__(self, model_name):
        self.timeouts = {
            'qwen-3.5': 120000,
            'kimi-k2.5': 60000,
            'minimax': 30000,
            'qwen-coder': 45000,
        }
        self.timeout = self.timeouts.get(model_name, 30000)
    
    async def execute(self, task):
        try:
            return await asyncio.wait_for(
                self.model.execute(task),
                timeout=self.timeout / 1000
            )
        except asyncio.TimeoutError:
            logger.error(f"Timeout for {self.model_name}")
            raise
```

### 8.2 Batch Processing

#### Batch Size Recommendations

| Model | Optimal Batch | Max Batch | Notes |
|-------|---------------|-----------|-------|
| MiniMax M2.5 | 10 | 50 | No rate limit |
| Kimi K2.5 | 3 | 5 | Moderate limit |
| Qwen 3.5 397B | 1 | 2 | Strict limit |
| Qwen 2.5 Coder | 2 | 3 | Moderate limit |

#### Batch Processing Implementation

```python
async def process_batch(tasks, model, max_batch=10):
    results = []
    for i in range(0, len(tasks), max_batch):
        batch = tasks[i:i + max_batch]
        batch_results = await asyncio.gather(
            *[model.execute(task) for task in batch],
            return_exceptions=True
        )
        results.extend(batch_results)
        
        # Respect rate limits between batches
        if i + max_batch < len(tasks):
            await asyncio.sleep(1)  # 1 second between batches
    
    return results
```

### 8.3 Connection Pooling

#### Pool Configuration

```python
class ConnectionPool:
    def __init__(self, model_name, pool_size):
        self.pool_size = pool_size
        self.connections = asyncio.Queue(maxsize=pool_size)
        self.active = 0
    
    async def acquire(self):
        if self.active < self.pool_size:
            self.active += 1
            return await self.create_connection()
        
        return await self.connections.get()
    
    async def release(self, connection):
        await self.connections.put(connection)
        self.active -= 1
```

#### Pool Size by Model

| Model | Recommended Pool | Notes |
|-------|-----------------|-------|
| Qwen 3.5 397B | 1 | Don't pool, use queue |
| Kimi K2.5 | 3 | Moderate pooling |
| MiniMax M2.5 | 10 | Aggressive pooling |
| Qwen 2.5 Coder | 2 | Conservative pooling |

---

## 9. Monitoring & Metrics

### 9.1 Key Metrics

#### Performance Metrics

| Metric | Target | Alert Threshold |
|--------|--------|------------------|
| Qwen Latency | < 90s | > 120s |
| Kimi Latency | < 20s | > 30s |
| MiniMax Latency | < 10s | > 15s |
| Success Rate | > 95% | < 90% |
| Cache Hit Rate | > 30% | < 20% |

#### Cost Metrics

| Metric | Target | Alert |
|--------|--------|-------|
| Daily Spend | $0.00 | > $0.00 |
| API Calls | Monitor | Spike detection |
| Token Usage | Track | Per-model limits |

### 9.2 Implementation

#### Metrics Collection

```python
class MetricsCollector:
    def __init__(self):
        self.metrics = {
            'requests': [],
            'latencies': [],
            'errors': [],
            'costs': [],
        }
    
    async def record_request(self, model, latency, success, tokens):
        self.metrics['requests'].append({
            'model': model,
            'latency': latency,
            'success': success,
            'tokens': tokens,
            'timestamp': time.time(),
        })
    
    def get_stats(self, model=None):
        requests = self.metrics['requests']
        if model:
            requests = [r for r in requests if r['model'] == model]
        
        return {
            'total': len(requests),
            'success_rate': sum(1 for r in requests if r['success']) / len(requests),
            'avg_latency': sum(r['latency'] for r in requests) / len(requests),
            'total_tokens': sum(r['tokens'] for r in requests),
        }
```

#### Dashboard Integration

```
Metrics Endpoints:
├── /metrics/latency - Response times by model
├── /metrics/success - Success/failure rates
├── /metrics/costs - Cost tracking
├── /metrics/usage - API usage patterns
└── /metrics/alerts - Active alerts
```

---

## 10. Implementation Guidelines

### 10.1 Model Router Implementation

```python
class ModelRouter:
    def __init__(self):
        self.rules = [
            (is_code_task, QwenRouter()),
            (is_doc_task, MiniMaxRouter()),
            (is_research_task, MiniMaxRouter()),
            (is_analysis_task, KimiRouter()),
        ]
    
    async def route(self, task):
        for condition, router in self.rules:
            if condition(task):
                return await router.route(task)
        
        return await self.default_router.route(task)

def is_code_task(task):
    return any(keyword in task.type for keyword in [
        'implement', 'create', 'refactor', 'fix', 'bug',
        'code', 'function', 'class', 'api'
    ])

def is_doc_task(task):
    return any(keyword in task.type for keyword in [
        'document', 'readme', 'guide', 'docs', 'write'
    ])

def is_research_task(task):
    return any(keyword in task.type for keyword in [
        'research', 'search', 'find', 'lookup', 'analyze'
    ])

def is_analysis_task(task):
    return any(keyword in task.type for keyword in [
        'analysis', 'audit', 'review', 'examine', 'assess'
    ])
```

### 10.2 Health Checks

```python
async def model_health_check(model):
    try:
        start = time.time()
        result = await model.execute({'type': 'ping'})
        latency = time.time() - start
        
        return {
            'model': model.name,
            'healthy': result.success,
            'latency': latency,
            'timestamp': time.time(),
        }
    except Exception as e:
        return {
            'model': model.name,
            'healthy': False,
            'error': str(e),
            'timestamp': time.time(),
        }

async def check_all_models():
    models = [Qwen(), Kimi(), MiniMax()]
    results = await asyncio.gather(*[model_health_check(m) for m in models])
    return {r['model']: r for r in results}
```

### 10.3 Error Handling

```python
class ModelErrorHandler:
    async def execute_with_fallback(self, task, model_chain):
        last_error = None
        
        for model in model_chain:
            try:
                return await model.execute(task)
            except RateLimitError as e:
                last_error = e
                await handle_rate_limit(e)
                continue
            except TimeoutError as e:
                last_error = e
                logger.warning(f"Timeout for {model.name}, trying next")
                continue
            except ModelUnavailableError as e:
                last_error = e
                logger.warning(f"Model {model.name} unavailable")
                continue
        
        raise AllModelsFailedError(f"All models failed: {last_error}")
```

---

## 11. Troubleshooting

### 11.1 Common Issues

#### Issue: Qwen Timeout

```
Symptom: Requests timeout after 120s
Cause: Model overloaded or network issue
Solution:
1. Check NVIDIA status
2. Retry with exponential backoff
3. Fall back to Qwen 2.5 Coder
```

#### Issue: Rate Limit (HTTP 429)

```
Symptom: Too many requests
Cause: Exceeded 40 RPM
Solution:
1. Implement request queuing
2. Wait 60s before retry
3. Use fallback model
```

#### Issue: Censorship Detected

```
Symptom: "Cannot generate that content"
Cause: Content policy triggered
Solution:
1. Detect trigger phrases
2. Switch to ZEN Big Pickle
3. Rephrase request if possible
```

#### Issue: High Latency

```
Symptom: Slow response times
Cause: Model cold start or load
Solution:
1. Pre-warm model
2. Use faster alternative
3. Implement caching
```

### 11.2 Debugging Commands

```bash
# Check model status
opencode models status

# Check NVIDIA API
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
     https://integrate.api.nvidia.com/v1/models

# Check current rate limits
opencode doctor

# View recent errors
tail -f logs/model-errors.log
```

---

## 12. Appendix

### 12.1 Model Configuration Reference

#### opencode.json Snippet

```json
{
  "provider": {
    "nvidia": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "NVIDIA NIM (Qwen 3.5)",
      "options": {
        "baseURL": "https://integrate.api.nvidia.com/v1",
        "timeout": 120000
      },
      "models": {
        "qwen-3.5-397b": {
          "id": "qwen/qwen3.5-397b-a17b",
          "limit": { "context": 262144, "output": 32768 }
        }
      }
    },
    "opencode-zen": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "OpenCode ZEN (FREE)",
      "options": {
        "baseURL": "https://api.opencode.ai/v1"
      },
      "models": {
        "minimax-m2.5-free": {
          "name": "MiniMax M2.5 (OpenCode ZEN)",
          "limit": { "context": 200000, "output": 32000 }
        },
        "kimi-k2.5-free": {
          "name": "Kimi K2.5 (OpenCode ZEN)",
          "limit": { "context": 1048576, "output": 65536 }
        }
      }
    }
  }
}
```

### 12.2 Environment Variables

```bash
# Required
NVIDIA_API_KEY=nvapi-xxxxxxxxxxxxx

# Optional (for monitoring)
METRICS_ENABLED=true
METRICS_ENDPOINT=http://localhost:9090
ALERT_WEBHOOK=https://hooks.slack.com/xxx
```

### 12.3 Quick Reference Card

```
┌─────────────────────────────────────────────────────────────────┐
│                    MODEL SELECTION QUICK REF                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  CODE TASKS                                                     │
│    Complex/Architecture → Qwen 3.5 397B                        │
│    Fast/Simple → Qwen 2.5 Coder 32B                            │
│    Security Critical → Qwen 3.5 397B                           │
│                                                                  │
│  DOCUMENTATION                                                  │
│    Always → MiniMax M2.5 (parallel OK)                         │
│                                                                  │
│  RESEARCH                                                       │
│    Web Search → MiniMax M2.5                                   │
│    Deep Analysis → Kimi K2.5                                    │
│                                                                  │
│  PARALLEL EXECUTION                                            │
│    Max 1: Qwen 3.5                                              │
│    Max 3: Kimi K2.5                                             │
│    Max 10: MiniMax M2.5                                        │
│                                                                  │
│  FALLBACK ORDER                                                 │
│    Qwen → Qwen Coder → Llama → ZEN                             │
│    MiniMax → Kimi → ZEN                                        │
│                                                                  │
│  TIMEOUTS                                                       │
│    Qwen: 120s | Kimi: 60s | MiniMax: 30s                     │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 12.4 Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-02-20 | Initial document creation |

---

**Document Status**: ACTIVE  
**Last Updated**: 2026-02-20  
**Owner**: BIOMETRICS Development Team  

---

*This document is part of the BIOMETRICS project and follows the guidelines established in AGENTS.md. For questions or suggestions, please refer to the project documentation.*
