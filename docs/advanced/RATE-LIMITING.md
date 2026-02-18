# RATE-LIMITING.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Rate Limiting schützt Services vor Überlastung und Missbrauch.
- Limits müssen dokumentiert und überwacht werden.
- Fallback-Mechanismen bei Überschreitung sind Pflicht.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

## Zweck
Dokumentation der Rate-Limiting-Strategie und -Konfiguration für BIOMETRICS.

## 1) Rate-Limit-Übersicht

| Endpoint | Limit | Window | Burst | Strategy |
|----------|-------|--------|-------|----------|
| /api/auth/login | 5 | 15 min | 2 | IP + User |
| /api/auth/register | 3 | 1 hour | 1 | IP |
| /api/ai/completion | 60 | 1 hour | 10 | API Key |
| /api/upload | 10 | 1 hour | 2 | API Key |
| /api/webhook | 100 | 1 min | 20 | IP |
| /* (default) | 100 | 1 min | 10 | IP |

## 2) Implementierung

### 2.1) Supabase (Edge Functions)

```typescript
// rate-limit.ts
import { createClient } from '@supabase/supabase-js'

const RATE_LIMIT_WINDOW = 60 // seconds
const RATE_LIMIT_MAX = 100

export async function rateLimit(
  identifier: string,
  limit: number = RATE_LIMIT_MAX,
  window: number = RATE_LIMIT_WINDOW
): Promise<{ success: boolean; remaining: number; reset: number }> {
  const supabase = createClient(process.env.SUPABASE_URL!, process.env.SUPABASE_KEY!)
  
  const { data, error } = await supabase.rpc('rate_limit', {
    p_identifier: identifier,
    p_limit: limit,
    p_window: window
  })
  
  if (error) throw error
  return data
}
```

### 2.2) Database Function (PostgreSQL)

```sql
CREATE OR REPLACE FUNCTION rate_limit(
  p_identifier TEXT,
  p_limit INT,
  p_window INT
)
RETURNS JSONB AS $$
DECLARE
  v_key TEXT := 'rate_limit:' || p_identifier;
  v_current_count INT;
  v_remaining INT;
  v_reset_time TIMESTAMP;
BEGIN
  -- Get current count
  SELECT COALESCE((SELECT count::int FROM redis.get(v_key)), 0)
  INTO v_current_count;
  
  v_remaining := p_limit - v_current_count;
  v_reset_time := NOW() + (p_window || ' seconds')::interval;
  
  IF v_current_count >= p_limit THEN
    RETURN jsonb_build_object(
      'success', false,
      'remaining', 0,
      'reset', EXTRACT(EPOCH FROM v_reset_time)
    );
  END IF;
  
  -- Increment counter
  PERFORM redis.incr(v_key);
  PERFORM redis.expire(v_key, p_window);
  
  RETURN jsonb_build_object(
    'success', true,
    'remaining', v_remaining,
    'reset', EXTRACT(EPOCH FROM v_reset_time)
  );
END;
$$ LANGUAGE plpgsql;
```

### 2.3) Vercel (Middleware)

```typescript
// middleware.ts
import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

const rateLimits = {
  '/api/auth/login': { limit: 5, window: 15 * 60 }, // 5 per 15 min
  '/api/ai/': { limit: 60, window: 60 * 60 }, // 60 per hour
}

export function middleware(request: NextRequest) {
  const path = request.nextUrl.pathname
  const ip = request.ip || request.headers.get('x-forwarded-for') || 'unknown'
  
  for (const [pattern, config] of Object.entries(rateLimits)) {
    if (path.startsWith(pattern)) {
      const key = `${pattern}:${ip}`
      const result = rateLimit(key, config.limit, config.window)
      
      if (!result.success) {
        return NextResponse.json(
          { error: 'Rate limit exceeded', retryAfter: result.reset },
          { status: 429, headers: { 'Retry-After': String(result.reset) } }
        )
      }
    }
  }
  
  return NextResponse.next()
}
```

## 3) Response-Headers

| Header | Beschreibung | Beispiel |
|--------|--------------|----------|
| X-RateLimit-Limit | Maximale Requests | 100 |
| X-RateLimit-Remaining | Verbleibende Requests | 95 |
| X-RateLimit-Reset | Reset-Zeit (Unix) | 1700000000 |
| Retry-After | Sekunden bis Reset | 45 |

## 4) Error-Handling

### 4.1) Client-Side (JavaScript)

```typescript
async function apiCall(url: string, options: RequestInit) {
  const response = await fetch(url, options)
  
  if (response.status === 429) {
    const retryAfter = response.headers.get('Retry-After') || 60
    console.log(`Rate limited. Retrying in ${retryAfter} seconds...`)
    
    await new Promise(resolve => setTimeout(resolve, retryAfter * 1000))
    return apiCall(url, options) // Retry
  }
  
  return response
}
```

### 4.2) Backend-Response

```json
{
  "error": "Rate limit exceeded",
  "message": "Too many requests. Please try again later.",
  "limit": 60,
  "remaining": 0,
  "reset": 1700000000
}
```

## 5) Monitoring

| Metric | Alert Threshold | Dashboard |
|--------|-----------------|-----------|
| Rate Limited Requests | > 10% des Traffic | Grafana |
| 429 Responses | > 5/min | PagerDuty |
| High Traffic IPs | > 1000 req/min | Slack |

## 6) Exemptions

| Endpoint | Reason | Approved By |
|----------|--------|-------------|
| /api/health | Monitoring | Auto |
| /api/webhook/verify | External Service | CTO |

---

## Abnahme-Check RATE-LIMITING
1. Limits für alle Endpunkte definiert
2. Implementierung dokumentiert
3. Response-Headers korrekt
4. Client-Retry-Logic vorhanden
5. Monitoring konfiguriert

---
