# WEBHOOK-RETRY.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Webhook-Delivery muss zuverlässig sein - keine Nachrichten verlieren.
- Retry-Logik mit Exponential Backoff ist Pflicht.
- Dead Letter Queue für fehlgeschlagene Deliveries erforderlich.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

## Zweck
Dokumentation der Webhook-Retry-Strategie und -Implementierung für BIOMETRICS.

## 1) Retry-Konfiguration

| Event Type | Max Retries | Initial Delay | Max Delay | Backoff Multiplier |
|------------|-------------|---------------|-----------|---------------------|
| payment.completed | 5 | 1s | 5min | 2x |
| payment.failed | 5 | 1s | 5min | 2x |
| user.created | 3 | 5s | 30s | 1.5x |
| user.deleted | 3 | 5s | 30s | 1.5x |
| order.created | 5 | 2s | 2min | 2x |
| subscription.updated | 3 | 5s | 30s | 1.5x |

## 2) Retry-Logik (Exponential Backoff)

### 2.1) Algorithmus

```
delay = initialDelay * (multiplier ^ attemptNumber)
maxDelay = min(delay, maxDelay)
jitter = random(0, 1000ms)
finalDelay = delay + jitter
```

### 2.2) Implementierung (TypeScript)

```typescript
interface WebhookConfig {
  maxRetries: number
  initialDelayMs: number
  maxDelayMs: number
  multiplier: number
}

const defaultConfig: WebhookConfig = {
  maxRetries: 5,
  initialDelayMs: 1000,
  maxDelayMs: 300000, // 5 min
  multiplier: 2
}

function calculateDelay(attempt: number, config: WebhookConfig): number {
  const delay = config.initialDelayMs * Math.pow(config.multiplier, attempt)
  const cappedDelay = Math.min(delay, config.maxDelayMs)
  const jitter = Math.floor(Math.random() * 1000)
  return cappedDelay + jitter
}

async function sendWithRetry(
  url: string,
  payload: object,
  config: WebhookConfig = defaultConfig
): Promise<{ success: boolean; attempts: number }> {
  for (let attempt = 0; attempt < config.maxRetries; attempt++) {
    try {
      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Webhook-Signature': generateSignature(payload)
        },
        body: JSON.stringify(payload)
      })
      
      if (response.ok) {
        return { success: true, attempts: attempt + 1 }
      }
      
      // Non-retryable status codes
      if (response.status === 400 || response.status === 401 || response.status === 403) {
        console.error(`Non-retryable error: ${response.status}`)
        break
      }
      
    } catch (error) {
      console.error(`Attempt ${attempt + 1} failed:`, error)
    }
    
    // Wait before next attempt
    if (attempt < config.maxRetries - 1) {
      const delay = calculateDelay(attempt, config)
      console.log(`Retrying in ${delay}ms...`)
      await new Promise(resolve => setTimeout(resolve, delay))
    }
  }
  
  return { success: false, attempts: config.maxRetries }
}
```

## 3) Dead Letter Queue (DLQ)

### 3.1) Schema

```sql
CREATE TABLE webhook_failures (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  webhook_type VARCHAR(100) NOT NULL,
  payload JSONB NOT NULL,
  endpoint VARCHAR(500) NOT NULL,
  attempts INT NOT NULL DEFAULT 0,
  last_error TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  status VARCHAR(20) DEFAULT 'pending' -- pending, retrying, failed, resolved
);

CREATE INDEX idx_webhook_failures_status ON webhook_failures(status);
CREATE INDEX idx_webhook_failures_created ON webhook_failures(created_at);
```

### 3.2) DLQ-Handler

```typescript
async function handleDeadLetter(webhookId: string) {
  const webhook = await db.query(
    'SELECT * FROM webhook_failures WHERE id = $1',
    [webhookId]
  )
  
  if (webhook.attempts >= webhook.maxRetries) {
    // Move to long-term storage or alert
    await notifyAdmin(webhook)
    await archiveForReview(webhook)
  }
}

async function retryFailedWebhooks() {
  const failed = await db.query(
    'SELECT * FROM webhook_failures WHERE status = $1 AND created_at > NOW() - INTERVAL \'24 hours\'',
    ['failed']
  )
  
  for (const webhook of failed) {
    await sendWithRetry(webhook.endpoint, webhook.payload)
  }
}
```

## 4) Monitoring

### 4.1) Metrics

| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| webhook.delivered.total | Gesamt deliverd | - |
| webhook.delivered.success | Erfolgreich | < 95% |
| webhook.delivered.failed | Fehlgeschlagen | > 5% |
| webhook.retry.count | Retry-Versuche | > 10/min |
| webhook.dlq.size | DLQ-Größe | > 100 |

### 4.2) Dashboard

```
Grafana Dashboard: Webhook Health
- Success Rate (%)
- Retry Distribution
- DLQ Size Over Time
- Average Delivery Time
```

## 5) Security

### 5.1) Signature Verification

```typescript
import crypto from 'crypto'

function generateSignature(payload: object, secret: string): string {
  const hmac = crypto.createHmac('sha256', secret)
  hmac.update(JSON.stringify(payload))
  return hmac.digest('hex')
}

function verifySignature(payload: object, signature: string, secret: string): boolean {
  const expected = generateSignature(payload, secret)
  return crypto.timingSafeEqual(Buffer.from(signature), Buffer.from(expected))
}
```

### 5.2) Webhook Headers

| Header | Description |
|--------|-------------|
| X-Webhook-Signature | HMAC-SHA256 |
| X-Webhook-Timestamp | Unix Timestamp |
| X-Webhook-Event | Event Type |
| X-Webhook-Id | Unique ID |

---

## Abnahme-Check WEBHOOK-RETRY
1. Retry-Konfiguration dokumentiert
2. Exponential Backoff implementiert
3. Dead Letter Queue vorhanden
4. Monitoring konfiguriert
5. Security (Signature) implementiert

---
