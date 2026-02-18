# QUEUE-SYSTEM.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Queue-Systeme ermöglichen asynchrone Verarbeitung und skalieren unabhängig.
- Message-Ordering und Idempotenz sind kritisch.
-监控 und Alerting sind Pflicht.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

## Zweck
Dokumentation der Queue-System-Architektur und -Konfiguration für BIOMETRICS.

## 1) Queue-Übersicht

| Queue Name | Type | Purpose | Priority | Retention |
|------------|------|---------|----------|-----------|
| ai-processing | Redis | KI-Anfragen | HIGH | 24h |
| email-send | Redis | E-Mail-Versand | MEDIUM | 7d |
| webhook-dispatch | Redis | Webhook-Delivery | MEDIUM | 3d |
| image-processing | Redis | Bildoptimierung | LOW | 24h |
| analytics-events | Redis | Event-Tracking | LOW | 7d |
| export-jobs | Redis | Export-Aufgaben | LOW | 3d |

## 2) Architektur

### 2.1) High-Level

```
┌─────────┐    ┌─────────┐    ┌─────────┐
│  API    │───▶│  Queue  │───▶│ Worker  │
│         │    │ (Redis) │    │         │
└─────────┘    └─────────┘    └─────────┘
                  │               │
                  ▼               ▼
            ┌─────────┐     ┌─────────┐
            │  DLQ    │     │ Result  │
            │         │     │ Store   │
            └─────────┘     └─────────┘
```

### 2.2) Producer (API)

```typescript
// Enqueue job
async function enqueueJob(
  queueName: string,
  payload: object,
  options: {
    priority?: 'high' | 'medium' | 'low'
    delay?: number // ms
    attempts?: number
  } = {}
) {
  const job = {
    id: crypto.randomUUID(),
    queue: queueName,
    payload,
    priority: options.priority || 'medium',
    delay: options.delay || 0,
    attempts: options.attempts || 3,
    createdAt: Date.now()
  }
  
  // Add to sorted set with score = timestamp
  await redis.zadd(
    `queue:${queueName}`,
    Date.now() + (job.delay || 0),
    JSON.stringify(job)
  )
  
  return job.id
}
```

### 2.3) Consumer (Worker)

```typescript
// Worker process
async function startWorker(queueName: string, handler: Function) {
  while (true) {
    // Get next job (blocking pop)
    const result = await redis.bzpopmin(`queue:${queueName}`, 0)
    
    if (!result) continue
    
    const job = JSON.parse(result[1])
    
    try {
      const result = await handler(job.payload)
      
      // Store result
      await redis.setex(
        `result:${job.id}`,
        86400, // 24h
        JSON.stringify({ status: 'success', data: result })
      )
      
    } catch (error) {
      // Handle failure
      await handleJobFailure(job, error)
    }
  }
}

async function handleJobFailure(job: object, error: Error) {
  job.attempts--
  
  if (job.attempts > 0) {
    // Retry with backoff
    const delay = Math.pow(2, (3 - job.attempts)) * 1000
    await redis.zadd(
      `queue:${job.queue}`,
      Date.now() + delay,
      JSON.stringify(job)
    )
  } else {
    // Move to DLQ
    await redis.zadd(
      `queue:dlq`,
      Date.now(),
      JSON.stringify({ ...job, error: error.message })
    )
  }
}
```

## 3) Job-Typen

### 3.1) AI Processing Queue

```typescript
interface AIJob {
  type: 'completion' | 'embedding' | 'vision'
  model: string
  input: string | object
  options: {
    temperature?: number
    maxTokens?: number
  }
}
```

### 3.2) Email Queue

```typescript
interface EmailJob {
  to: string
  subject: string
  template: string
  data: object
  priority: 'high' | 'normal' | 'low'
}
```

### 3.3) Webhook Queue

```typescript
interface WebhookJob {
  event: string
  endpoint: string
  payload: object
  signature: string
}
```

## 4) Monitoring

### 4.1) Metrics

| Metric | Description | Alert |
|--------|-------------|-------|
| queue.size | Jobs in Queue | > 1000 |
| queue.wait_time | Avg Wait Time | > 30s |
| worker.processing | Aktive Worker | < 1 |
| job.failed | Fehlgeschlagene Jobs | > 10 |
| job.duration | Avg Processing Time | > 60s |

### 4.2) Dashboard Queries

```promql
# Queue Size
redis_zcard(queue:ai-processing)

# Job Success Rate
sum(rate(jobs_completed_total[5m])) / sum(rate(jobs_total[5m]))

# Avg Wait Time
histogram_quantile(0.95, job_wait_time_seconds_bucket)
```

## 5) Best Practices

### 5.1) Idempotency
```typescript
// Check if job already processed
const isProcessed = await redis.exists(`processed:${jobId}`)
if (isProcessed) {
  return { status: 'skipped', reason: 'already_processed' }
}

// Mark as processed
await redis.setex(`processed:${jobId}`, 86400, '1')
```

### 5.2) Dead Letter Handling
```typescript
// Manual retry from DLQ
async function retryFromDLQ(jobId: string) {
  const job = await redis.get(`dlq:${jobId}`)
  await enqueueJob(job.queue, job.payload, { attempts: job.attempts })
}

// Manual resolve
async function resolveDLQ(jobId: string, resolution: string) {
  await redis.setex(`resolution:${jobId}`, 2592000, resolution) // 30d
}
```

---

## Abnahme-Check QUEUE-SYSTEM
1. Queue-Architektur dokumentiert
2. Producer/Consumer implementiert
3. Monitoring konfiguriert
4. DLQ vorhanden
5. Idempotency sichergestellt

---
