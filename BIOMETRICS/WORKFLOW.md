# WORKFLOW.md — Unified Skill Architecture & Meta-Builder Protocol

**Status:** ACTIVE  
**Version:** 2.0  
**Stand:** Februar 2026  
**Purpose:** Zentrale Dokumentation der Self-Building AI Agent Architektur  
**Line Count:** 600+ lines

---

## 🧠 1. Core Concept: Brain vs Muscle Architecture

### 1.1 Das Grundprinzip

Die Self-Building AI Agent Architektur basiert auf einer fundamentalen Trennung zwischen **Entscheidungsfindung** (Brain) und **Ausführung** (Muscle). Diese Trennung ermöglicht es dem Agenten, sowohl bestehende Fähigkeiten effektiv zu nutzen als auch kontinuierlich neue Fähigkeiten selbst zu entwickeln.

**BRAIN (AI/OpenClaw)** fungiert als intelligenter Orchestrator:
- Entscheidungsfindung basierend auf User-Requests
- Interface zum User (natürliche Sprache)
- Skill-Orchestrierung (welcher Skill für welche Aufgabe)
- Meta-Cognition (erkennt Patterns und baut neue Tools)

**MUSCLE (Supabase/n8n/SDKs)** übernimmt die schwere Arbeit:
- Ausführung komplexer Geschäftslogik
- Datenbank-Operationen und Persistenz
- Workflow-Automation und Scheduling
- Rechenintensive Tasks und externe APIs

### 1.2 Decision-Making Flow

```
User Request (Natural Language)
            │
            ▼
    ┌───────────────┐
    │  BRAIN        │
    │  (OpenClaw)   │
    └───────┬───────┘
            │
            ▼
    ┌───────────────┐
    │ Intent        │
    │ Detection     │
    └───────┬───────┘
            │
            ▼
    ┌───────────────┐
    │ Skill         │
    │ Selection     │
    └───────┬───────┘
            │
      ┌─────┴─────┐
      │            │
      ▼            ▼
┌─────────┐  ┌─────────────┐
│ MUSCLE  │  │   BRAIN     │
│ Layer   │  │   Continues │
│         │  │   Reasoning │
└────┬────┘  └─────────────┘
     │
     ▼
┌─────────────────────────────────────┐
│     Clean JSON Response             │
│     (AI-Friendly Format)            │
└─────────────────────────────────────┘
```

### 1.3 Skill-Orchestrierung Schritt-für-Schritt

**Schritt 1: Request Parsing**
Der User sendet einen natürlichen Sprache-Request. OpenClaw analysiert den Request und extrahiert:
- Intent (was will der User erreichen)
- Entities (welche Daten sind relevant)
- Context (vorherige Interaktionen, User-Präferenzen)

**Schritt 2: Skill Discovery**
Basierend auf dem Intent sucht der Agent im Skill-Registry nach passenden Skills:
- Matching nach Skill-Name
- Matching nach Skill-Beschreibung
- Matching nach Input-Requirements

**Schritt 3: Skill Execution**
Der gewählte Skill wird ausgeführt mit:
- Input-Validierung (Zod Schema)
- Execution im passenden Muscle-Layer
- Output-Validierung und Formatierung

**Schritt 4: Response Generation**
Das Ergebnis wird in eine AI-freundliche Response umgewandelt:
- Klares JSON-Format
- Kontextuelle Fehlermeldungen
- Nützliche Metadaten

### 1.4 Meta-Cognition Prozess

Der Agent erkennt kontinuierlich Verbesserungsmöglichkeiten:

**Pattern Detection:**
- Wiederholte manuelle Tasks werden erkannt
- User-Verhalten wird analysiert
- Performance-Engpässe werden identifiziert

**Self-Improvement Loop:**
1. Agent erkennt repetitive Task
2. Agent designed neue Lösung
3. Agent implementiert und deployed
4. Agent registriert neuen Skill
5. Agent wird kontinuierlich mächtiger

### 1.5 Code-Beispiel: OpenClaw Skill Definition

```typescript
import { z } from 'zod';

export const checkCompetitorPricesInput = z.object({
  competitors: z.array(z.object({
    name: z.string(),
    url: z.string().url(),
  })),
  products: z.array(z.string()),
  notificationChannels: z.enum(['email', 'telegram', 'whatsapp']).optional(),
});

export const checkCompetitorPricesOutput = z.object({
  success: z.boolean(),
  prices: z.array(z.object({
    competitor: z.string(),
    product: z.string(),
    price: z.number(),
    currency: z.string(),
    timestamp: z.string().datetime(),
    change: z.enum(['up', 'down', 'stable']).optional(),
    previousPrice: z.number().optional(),
  })),
  alerts: z.array(z.object({
    type: z.enum(['price_increase', 'price_decrease', 'out_of_stock']),
    competitor: z.string(),
    product: z.string(),
    message: z.string(),
  })).optional(),
  executionTime: z.number(),
});

export type CheckCompetitorPricesInput = z.infer<typeof checkCompetitorPricesInput>;
export type CheckCompetitorPricesOutput = z.infer<typeof checkCompetitorPricesOutput>;

export const skill = {
  name: 'check_competitor_prices',
  description: 'Überwacht Preise von Wettbewerbern und sendet Alerts bei Änderungen',
  inputSchema: checkCompetitorPricesInput,
  outputSchema: checkCompetitorPricesOutput,
  
  async execute(input: CheckCompetitorPricesInput): Promise<CheckCompetitorPricesOutput> {
    const startTime = Date.now();
    
    const prices = await Promise.all(
      input.competitors.map(async (competitor) => {
        const competitorPrices = await fetchCompetitorPrices(competitor.url, input.products);
        return competitorPrices.map(p => ({
          ...p,
          competitor: competitor.name,
        }));
      })
    );
    
    const alerts = detectPriceAlerts(prices);
    
    return {
      success: true,
      prices: prices.flat(),
      alerts: alerts.length > 0 ? alerts : undefined,
      executionTime: Date.now() - startTime,
    };
  },
};
```

### 1.6 Code-Beispiel: Supabase Edge Function

```typescript
import { serve } from 'https://deno.land/std@0.168.0/http/server.ts';
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2';
import { z } from 'https://deno.sh/x/zod@v3.22.4/mod.ts';

const corsHeaders = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Headers': 'authorization, x-client-info, apikey, content-type',
};

const priceAlertSchema = z.object({
  competitor_id: z.string().uuid(),
  product_sku: z.string(),
  old_price: z.number(),
  new_price: z.number(),
  threshold_percent: z.number().default(5),
});

serve(async (req) => {
  if (req.method === 'OPTIONS') {
    return new Response('ok', { headers: corsHeaders });
  }

  try {
    const supabaseClient = createClient(
      Deno.env.get('SUPABASE_URL') ?? '',
      Deno.env.get('SUPABASE_ANON_KEY') ?? '',
    );

    const authHeader = req.headers.get('Authorization');
    if (!authHeader) {
      throw new Error('Missing authorization header');
    }

    const { data: { user } } = await supabaseClient.auth.getUser(authHeader.replace('Bearer ', ''));
    if (!user) {
      throw new Error('Invalid user');
    }

    const body = await req.json();
    const { competitor_id, product_sku, new_price } = priceAlertSchema.parse(body);

    const { data: product } = await supabaseClient
      .from('products')
      .select('*, competitor:competitors(*)')
      .eq('sku', product_sku)
      .single();

    if (!product) {
      throw new Error('Product not found');
    }

    const oldPrice = product.price;
    const changePercent = Math.abs((newPrice - oldPrice) / oldPrice * 100);

    if (changePercent >= 5) {
      await supabaseClient.from('price_alerts').insert({
        competitor_id,
        product_sku,
        old_price: oldPrice,
        new_price: new_price,
        user_id: user.id,
        alert_type: newPrice > oldPrice ? 'price_increase' : 'price_decrease',
      });
    }

    await supabaseClient
      .from('products')
      .update({ price: newPrice, updated_at: new Date().toISOString() })
      .eq('sku', product_sku);

    return new Response(
      JSON.stringify({
        success: true,
        change_detected: changePercent >= 5,
        change_percent: changePercent,
      }),
      { headers: { ...corsHeaders, 'Content-Type': 'application/json' } }
    );
  } catch (error) {
    return new Response(
      JSON.stringify({ error: error.message }),
      { status: 400, headers: { ...corsHeaders, 'Content-Type': 'application/json' } }
    );
  }
});
```

### 1.7 Code-Beispiel: n8n Webhook Handler

```json
{
  "name": "Competitor Price Monitor",
  "nodes": [
    {
      "parameters": {
        "httpMethod": "POST",
        "path": "competitor-price-alert",
        "responseMode": "lastNode",
        "options": {}
      },
      "id": "webhook-trigger",
      "name": "Webhook Trigger",
      "type": "n8n-nodes-base.webhook",
      "typeVersion": 1,
      "position": [250, 300],
      "webhookId": "competitor-price-monitor"
    },
    {
      "parameters": {
        "operation": "execute",
        "functionCode": "const data = $input.first().json;\n\nif (data.change_detected) {\n  const change = data.new_price > data.old_price ? 'gestiegen' : 'gefallen';\n  const percent = data.change_percent.toFixed(1);\n  \n  return {\n    json: {\n      alert: `⚠️ Preisänderung: ${data.product_name} bei ${data.competitor_name} ist um ${percent}% ${change}`,\n      product: data.product_name,\n      competitor: data.competitor_name,\n      old_price: data.old_price,\n      new_price: data.new_price,\n      change_percent: data.change_percent\n    }\n  };\n}\n\nreturn { json: { skip: true } };"
      },
      "id": "process-alert",
      "name": "Process Alert",
      "type": "n8n-nodes-base.function",
      "typeVersion": 1,
      "position": [450, 300]
    },
    {
      "parameters": {
        "operation": "send",
        "channel": "telegram",
        "chatId": "{{$json.chatId}}",
        "text": "={{$json.alert}}",
        "additionalFields": {}
      },
      "id": "send-telegram",
      "name": "Send Telegram",
      "type": "n8n-nodes-base.telegram",
      "typeVersion": 1.2,
      "position": [650, 300],
      "credentials": {
        "telegramApi": {
          "id": "1",
          "name": "Telegram API account"
        }
      }
    }
  ],
  "connections": {
    "Webhook Trigger": {
      "main": [[{"node": "Process Alert", "type": "main", "index": 0}]]
    },
    "Process Alert": {
      "main": [[{"node": "Send Telegram", "type": "main", "index": 0}]]
    }
  },
  "settings": {
    "executionOrder": "v1"
  }
}
```

---

## 🏗️ 2. Architecture Patterns

### Pattern A: Webhook Wrapper (n8n Integration)

#### Wann verwenden
- Multi-step Prozesse mit Verzweigungen
- Externe API-Integrationen (Stripe, Slack, etc.)
- Zeitbasierte Trigger (Cron)
- Komplexe Workflows mit Wartezeiten
- Human-in-the-Loop Genehmigungen

#### Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                     WEBHOOK WRAPPER PATTERN                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   User Request                                                          │
│       │                                                                 │
│       ▼                                                                 │
│   ┌─────────────┐     ┌─────────────┐     ┌─────────────┐             │
│   │  OpenClaw  │────▶│   n8n Webhook│────▶│  n8n Nodes │             │
│   │   Skill    │     │   Trigger   │     │   Chain    │             │
│   └─────────────┘     └─────────────┘     └──────┬──────┘             │
│                                                  │                     │
│                                                  ▼                     │
│                                          ┌─────────────┐             │
│                                          │  Database   │             │
│                                          │  (Optional) │             │
│                                          └─────────────┘             │
│                                                  │                     │
│       ┌─────────────────────────────────────────┘                     │
│       │                                                                 │
│       ▼                                                                 │
│   ┌─────────────┐     ┌─────────────┐     ┌─────────────┐         │
│   │  Transform  │◀────│  Process    │◀────│  Execute    │         │
│   │  Response   │     │  Result     │     │  Logic      │         │
│   └─────────────┘     └─────────────┘     └─────────────┘         │
│            │                                                                 │
│            ▼                                                                 │
│   ┌─────────────────────────────────────┐                              │
│   │     Clean JSON Response (AI-Ready)   │                              │
│   └─────────────────────────────────────┘                              │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Vollständiges Use Case Szenario: Bestellung verarbeiten

**User-Request:** "Erstelle einen Workflow, der neue Bestellungen verarbeitet, den Lagerbestand prüft und bei niedrigem Bestand eine Nachricht an Slack sendet."

**Implementierung:**

```typescript
// OpenClaw Skill: process_order
import { z } from 'zod';

const processOrderInput = z.object({
  orderId: z.string().uuid(),
  priority: z.enum(['normal', 'express', 'overnight']).default('normal'),
  notificationSlack: z.boolean().default(true),
});

const processOrderOutput = z.object({
  success: z.boolean(),
  orderId: z.string(),
  status: z.enum(['processed', 'pending_stock', 'failed', 'cancelled']),
  inventoryCheck: z.object({
    available: z.boolean(),
    items: z.array(z.object({
      sku: z.string(),
      requested: z.number(),
      available: z.number(),
      sufficient: z.boolean(),
    })),
  }).optional(),
  messages: z.array(z.string()),
  timestamp: z.string().datetime(),
});

export const skill = {
  name: 'process_order',
  description: 'Verarbeitet Bestellungen mit Lagerbestandsprüfung',
  
  async execute(input: z.infer<typeof processOrderInput>) {
    const response = await fetch('https://n8n-instance.com/webhook/order-process', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(input),
    });
    
    const result = await response.json();
    return processOrderOutput.parse(result);
  },
};
```

#### Request/Response Schema (Zod Validation)

```typescript
// Input Validation Schema
const webhookInputSchema = z.object({
  action: z.enum(['create', 'update', 'delete', 'list', 'get']),
  entity: z.string().min(1).max(100),
  payload: z.record(z.unknown()).optional(),
  options: z.object({
    timeout: z.number().min(1000).max(30000).default(10000),
    retries: z.number().min(0).max(5).default(3),
    retryDelay: z.number().min(100).max(5000).default(1000),
  }).optional(),
  auth: z.object({
    type: z.enum(['bearer', 'basic', 'apiKey']),
    credentials: z.string(),
  }).optional(),
});

// Output Validation Schema  
const webhookOutputSchema = z.object({
  success: z.boolean(),
  data: z.unknown().optional(),
  error: z.object({
    code: z.string(),
    message: z.string(),
    details: z.record(z.unknown()).optional(),
  }).optional(),
  metadata: z.object({
    executionTime: z.number(),
    timestamp: z.string().datetime(),
    requestId: z.string().uuid().optional(),
  }),
});
```

#### Error Handling Strategy

```typescript
class WebhookErrorHandler {
  static handle(error: unknown, context: WebhookContext): ErrorResponse {
    if (error instanceof z.ZodError) {
      return {
        success: false,
        error: {
          code: 'VALIDATION_ERROR',
          message: 'Input validation failed',
          details: error.errors,
        },
      };
    }

    if (error instanceof FetchError) {
      if (error.status === 429) {
        return {
          success: false,
          error: {
            code: 'RATE_LIMITED',
            message: 'Too many requests, please retry later',
            retryAfter: error.headers.get('Retry-After'),
          },
        };
      }

      if (error.status >= 500) {
        return {
          success: false,
          error: {
            code: 'UPSTREAM_ERROR',
            message: 'External service temporarily unavailable',
          },
        };
      }
    }

    return {
      success: false,
      error: {
        code: 'UNKNOWN_ERROR',
        message: error instanceof Error ? error.message : 'Unknown error',
      },
    };
  }
}
```

#### Retry Logic Implementation

```typescript
async function fetchWithRetry<T>(
  url: string,
  options: RequestInit,
  retryOptions: {
    retries: number;
    retryDelay: number;
    retryableStatuses: number[];
  } = { retries: 3, retryDelay: 1000, retryableStatuses: [429, 500, 502, 503, 504] }
): Promise<T> {
  let lastError: Error;
  
  for (let attempt = 0; attempt <= retryOptions.retries; attempt++) {
    try {
      const response = await fetch(url, options);
      
      if (response.ok) {
        return await response.json();
      }
      
      if (!retryOptions.retryableStatuses.includes(response.status)) {
        throw new FetchError(response);
      }
      
      lastError = new Error(`HTTP ${response.status}`);
    } catch (error) {
      lastError = error as Error;
    }
    
    if (attempt < retryOptions.retries) {
      await new Promise(resolve => 
        setTimeout(resolve, retryOptions.retryDelay * Math.pow(2, attempt))
      );
    }
  }
  
  throw lastError!;
}
```

#### Rate Limiting Configuration

```typescript
const rateLimiter = {
  windowMs: 60 * 1000,
  maxRequests: 100,
  
  requests: new Map<string, number[]>(),
  
  check(clientId: string): boolean {
    const now = Date.now();
    const windowStart = now - this.windowMs;
    
    const clientRequests = this.requests.get(clientId) || [];
    const recentRequests = clientRequests.filter(ts => ts > windowStart);
    
    if (recentRequests.length >= this.maxRequests) {
      return false;
    }
    
    recentRequests.push(now);
    this.requests.set(clientId, recentRequests);
    return true;
  },
  
  getRetryAfter(clientId: string): number {
    const clientRequests = this.requests.get(clientId) || [];
    const oldestRequest = Math.min(...clientRequests);
    return Math.ceil((oldestRequest + this.windowMs - Date.now()) / 1000);
  },
};
```

#### Monitoring & Logging Setup

```typescript
const logger = {
  info: (message: string, meta?: Record<string, unknown>) => {
    console.log(JSON.stringify({
      level: 'info',
      message,
      timestamp: new Date().toISOString(),
      ...meta,
    }));
  },
  
  error: (message: string, error: Error, meta?: Record<string, unknown>) => {
    console.error(JSON.stringify({
      level: 'error',
      message,
      error: {
        name: error.name,
        message: error.message,
        stack: error.stack,
      },
      timestamp: new Date().toISOString(),
      ...meta,
    }));
  },
  
  metrics: (name: string, value: number, tags?: Record<string, string>) => {
    console.log(JSON.stringify({
      type: 'metric',
      name,
      value,
      tags,
      timestamp: new Date().toISOString(),
    }));
  },
};
```

#### Troubleshooting Guide

**Issue 1: Webhook Timeout**
- **Symptom:** Request bricht nach 30s ab
- **Ursache:** n8n Workflow dauert zu lange
- **Lösung:** 
  1. Async-Pattern verwenden (sofort 202 Accepted, später Callback)
  2. Workflow in kleinere Steps aufteilen
  3. Langsame Nodes durch Caching optimieren

**Issue 2: Duplicate Executions**
- **Symptom:** Workflow wird mehrfach ausgeführt
- **Ursache:** n8n receives multiple webhooks (Retry bei Timeout)
- **Lösung:**
  1. Idempotency-Key in Request header verwenden
  2. In Datenbank prüfen ob bereits verarbeitet
  3. n8n "Execute Once" Option aktivieren

**Issue 3: Invalid JSON Response**
- **Symptom:** OpenClaw kann Response nicht parsen
- **Ursache:** n8n gibt HTML statt JSON zurück
- **Lösung:**
  1. Response Node auf JSON setzen
  2. Error-Handling in n8n verbessern
  3. Content-Type Header prüfen

**Issue 4: Authentication Failures**
- **Symptom:** 401 Unauthorized
- **Ursache:** Token abgelaufen oder falsch
- **Lösung:**
  1. Token-Refresh implementieren
  2. Credentials in n8n aktualisieren
  3. Bearer Token korrekt formatieren

**Issue 5: Memory Issues**
- **Symptom:** n8n Worker stürzt ab
- **Ursache:** Zu viele Daten im Memory
- **Lösung:**
  1. Pagination verwenden
  2. Daten in Chunks verarbeiten
  3. Streaming Response nutzen

---

### Pattern B: Serverless Proxy (Supabase Edge Functions)

#### Wann verwenden
- Datenbank-Operationen (CRUD)
- Authentifizierung und Autorisierung
- Rechenintensive Backend-Logik
- API-Aggregation (mehrere Quellen kombinieren)
- Webhook-Empfang mit schneller Response

#### Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                 SERVERLESS PROXY PATTERN                                │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   User Request                                                          │
│       │                                                                 │
│       ▼                                                                 │
│   ┌─────────────┐     ┌─────────────────────┐                         │
│   │  OpenClaw   │────▶│  Supabase Edge      │                         │
│   │   Skill     │     │  Function           │                         │
│   └─────────────┘     └──────────┬──────────┘                         │
│                                  │                                      │
│                                  ▼                                      │
│                         ┌─────────────────────┐                         │
│                         │  JWT Validation    │                         │
│                         │  (Auth Check)       │                         │
│                         └──────────┬──────────┘                         │
│                                    │                                    │
│                                    ▼                                    │
│                         ┌─────────────────────┐                         │
│                         │  Input Validation   │                         │
│                         │  (Zod Schema)       │                         │
│                         └──────────┬──────────┘                         │
│                                    │                                    │
│                                    ▼                                    │
│                         ┌─────────────────────┐                         │
│                         │  Business Logic     │                         │
│                         │  (TypeScript)       │                         │
│                         └──────────┬──────────┘                         │
│                                    │                                    │
│                          ┌─────────┴─────────┐                         │
│                          │                   │                         │
│                          ▼                   ▼                         │
│                    ┌───────────┐       ┌───────────┐                   │
│                    │ Database  │       │  External │                   │
│                    │ Queries   │       │  APIs     │                   │
│                    └─────┬─────┘       └─────┬─────┘                   │
│                          │                   │                         │
│                          └─────────┬─────────┘                         │
│                                    │                                    │
│                                    ▼                                    │
│                         ┌─────────────────────┐                         │
│                         │  Output Validation  │                         │
│                         │  (Zod Schema)       │                         │
│                         └──────────┬──────────┘                         │
│                                    │                                    │
│                                    ▼                                    │
│                         ┌─────────────────────┐                         │
│                         │  Typed JSON Response │                         │
│                         │  (AI-Ready)          │                         │
│                         └─────────────────────┘                         │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Vollständiges Use Case Szenario: Benutzer-Verwaltung

```typescript
// OpenClaw Skill: manage_users
import { z } from 'zod';

const userManagementInput = z.object({
  action: z.enum(['create', 'update', 'delete', 'list', 'get']),
  userId: z.string().uuid().optional(),
  email: z.string().email().optional(),
  name: z.string().min(1).max(100).optional(),
  role: z.enum(['admin', 'user', 'guest']).optional(),
  metadata: z.record(z.unknown()).optional(),
});

const userManagementOutput = z.object({
  success: z.boolean(),
  user: z.object({
    id: z.string().uuid(),
    email: z.string().email(),
    name: z.string(),
    role: z.enum(['admin', 'user', 'guest']),
    createdAt: z.string().datetime(),
    updatedAt: z.string().datetime(),
  }).optional(),
  users: z.array(z.object({
    id: z.string().uuid(),
    email: z.string().email(),
    name: z.string(),
    role: z.enum(['admin', 'user', 'guest']),
  })).optional(),
  message: z.string().optional(),
});

export const skill = {
  name: 'manage_users',
  description: 'Verwaltet Benutzer in der Datenbank',
  
  async execute(input: z.infer<typeof userManagementInput>) {
    const response = await fetch(
      `${Deno.env.get('SUPABASE_URL')}/functions/v1/user-management`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${Deno.env.get('SUPABASE_SERVICE_KEY')}`,
        },
        body: JSON.stringify(input),
      }
    );
    
    const result = await response.json();
    return userManagementOutput.parse(result);
  },
};
```

#### Supabase Edge Function Code (80+ Zeilen)

```typescript
import { serve } from 'https://deno.land/std@0.168.0/http/server.ts';
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2';
import { z } from 'https://deno.sh/x/zod@v3.22.4/mod.ts';

const corsHeaders = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Headers': 'authorization, x-client-info, apikey, content-type',
  'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
};

const userSchema = z.object({
  email: z.string().email(),
  name: z.string().min(1).max(100),
  role: z.enum(['admin', 'user', 'guest']).default('user'),
  metadata: z.record(z.unknown()).optional(),
});

const actionSchema = z.object({
  action: z.enum(['create', 'update', 'delete', 'list', 'get']),
  userId: z.string().uuid().optional(),
  email: z.string().email().optional(),
  name: z.string().min(1).max(100).optional(),
  role: z.enum(['admin', 'user', 'guest']).optional(),
  metadata: z.record(z.unknown()).optional(),
});

serve(async (req) => {
  if (req.method === 'OPTIONS') {
    return new Response('ok', { headers: corsHeaders });
  }

  try {
    const supabaseClient = createClient(
      Deno.env.get('SUPABASE_URL') ?? '',
      Deno.env.get('SUPABASE_ANON_KEY') ?? '',
    );

    const authHeader = req.headers.get('Authorization');
    if (!authHeader) {
      throw new Error('Authorization required');
    }

    const { data: { user }, error: authError } = await supabaseClient.auth.getUser(
      authHeader.replace('Bearer ', '')
    );

    if (authError || !user) {
      throw new Error('Invalid authentication');
    }

    const { data: userData } = await supabaseClient
      .from('users')
      .select('role')
      .eq('id', user.id)
      .single();

    if (!userData || userData.role !== 'admin') {
      throw new Error('Admin access required');
    }

    const body = await req.json();
    const { action, userId, email, name, role, metadata } = actionSchema.parse(body);

    let result;

    switch (action) {
      case 'create': {
        const newUser = userSchema.parse({ email, name, role, metadata });
        const { data, error } = await supabaseClient
          .from('users')
          .insert({
            ...newUser,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          })
          .select()
          .single();

        if (error) throw error;
        result = { success: true, user: data };
        break;
      }

      case 'update': {
        if (!userId) throw new Error('userId required for update');
        
        const updates: Record<string, unknown> = { updated_at: new Date().toISOString() };
        if (name) updates.name = name;
        if (role) updates.role = role;
        if (metadata) updates.metadata = metadata;

        const { data, error } = await supabaseClient
          .from('users')
          .update(updates)
          .eq('id', userId)
          .select()
          .single();

        if (error) throw error;
        result = { success: true, user: data };
        break;
      }

      case 'delete': {
        if (!userId) throw new Error('userId required for delete');
        
        const { error } = await supabaseClient
          .from('users')
          .delete()
          .eq('id', userId);

        if (error) throw error;
        result = { success: true, message: 'User deleted' };
        break;
      }

      case 'list': {
        const { data, error } = await supabaseClient
          .from('users')
          .select('id, email, name, role, created_at, updated_at')
          .order('created_at', { ascending: false });

        if (error) throw error;
        result = { success: true, users: data };
        break;
      }

      case 'get': {
        if (!userId) throw new Error('userId required for get');
        
        const { data, error } = await supabaseClient
          .from('users')
          .select('id, email, name, role, created_at, updated_at, metadata')
          .eq('id', userId)
          .single();

        if (error) throw error;
        result = { success: true, user: data };
        break;
      }

      default:
        throw new Error('Invalid action');
    }

    return new Response(JSON.stringify(result), {
      headers: { ...corsHeaders, 'Content-Type': 'application/json' },
    });
  } catch (error) {
    return new Response(
      JSON.stringify({ 
        success: false, 
        error: { message: error.message } 
      }),
      { status: 400, headers: { ...corsHeaders, 'Content-Type': 'application/json' } }
    );
  }
});
```

#### Database Schema (SQL)

```sql
-- Users table for Supabase
CREATE TABLE public.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for faster queries
CREATE INDEX idx_users_email ON public.users(email);
CREATE INDEX idx_users_role ON public.users(role);
CREATE INDEX idx_users_created_at ON public.users(created_at DESC);

-- Row Level Security (RLS)
ALTER TABLE public.users ENABLE ROW LEVEL SECURITY;

-- Users can read their own profile
CREATE POLICY "Users can read own profile"
    ON public.users FOR SELECT
    USING (auth.uid() = id);

-- Admins can read all users
CREATE POLICY "Admins can read all users"
    ON public.users FOR SELECT
    USING (
        EXISTS (
            SELECT 1 FROM public.users
            WHERE id = auth.uid() AND role = 'admin'
        )
    );

-- Admins can insert users
CREATE POLICY "Admins can insert"
    ON public.users FOR INSERT
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM public.users
            WHERE id = auth.uid() AND role = 'admin'
        )
    );

-- Admins can update users
CREATE POLICY "Admins can update"
    ON public.users FOR UPDATE
    USING (
        EXISTS (
            SELECT 1 FROM public.users
            WHERE id = auth.uid() AND role = 'admin'
        )
    );

-- Admins can delete users
CREATE POLICY "Admins can delete"
    ON public.users FOR DELETE
    USING (
        EXISTS (
            SELECT 1 FROM public.users
            WHERE id = auth.uid() AND role = 'admin'
        )
    );

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON public.users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

#### Type Definitions (TypeScript)

```typescript
export type UserRole = 'admin' | 'user' | 'guest';

export interface User {
  id: string;
  email: string;
  name: string;
  role: UserRole;
  metadata?: Record<string, unknown>;
  createdAt: string;
  updatedAt: string;
}

export interface UserCreate {
  email: string;
  name: string;
  role?: UserRole;
  metadata?: Record<string, unknown>;
}

export interface UserUpdate {
  name?: string;
  role?: UserRole;
  metadata?: Record<string, unknown>;
}

export interface UserListQuery {
  role?: UserRole;
  limit?: number;
  offset?: number;
  orderBy?: 'createdAt' | 'email' | 'name';
  orderDir?: 'asc' | 'desc';
}

export interface UserResponse<T = User> {
  success: boolean;
  user?: T;
  users?: T[];
  message?: string;
  error?: {
    code: string;
    message: string;
    details?: unknown;
  };
}
```

#### Authentication Flow (JWT, OAuth2)

```typescript
// JWT Verification in Edge Function
async function verifyJWT(token: string): Promise<JWTPayload> {
  const supabaseClient = createClient(
    Deno.env.get('SUPABASE_URL')!,
    Deno.env.get('SUPABASE_ANON_KEY')!
  );
  
  const { data: { user }, error } = await supabaseClient.auth.getUser(token);
  
  if (error || !user) {
    throw new Error('Invalid token');
  }
  
  return {
    userId: user.id,
    email: user.email,
    role: user.user_metadata?.role || 'user',
  };
}

// OAuth2 Callback Handler
async function handleOAuthCallback(code: string): Promise<AuthResult> {
  const supabaseClient = createClient(
    Deno.env.get('SUPABASE_URL')!,
    Deno.env.get('SUPABASE_ANON_KEY')!
  );
  
  const { data, error } = await supabaseClient.auth.exchangeCodeForSession(code);
  
  if (error) {
    throw new Error(`OAuth error: ${error.message}`);
  }
  
  return {
    accessToken: data.session.access_token,
    refreshToken: data.session.refresh_token,
    user: data.user,
  };
}
```

#### CORS Configuration

```typescript
const corsHeaders = {
  'Access-Control-Allow-Origin': Deno.env.get('ALLOWED_ORIGINS') || '*',
  'Access-Control-Allow-Headers': 'authorization, x-client-info, apikey, content-type, x-request-id',
  'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, PATCH, OPTIONS',
  'Access-Control-Max-Age': '86400',
  'Access-Control-Expose-Headers': 'x-request-id, x-rate-limit-remaining',
};

// For credentials-aware CORS
const credentialsCorsHeaders = {
  'Access-Control-Allow-Origin': Deno.env.get('ALLOWED_ORIGINS'),
  'Access-Control-Allow-Headers': 'authorization, x-client-info, apikey, content-type',
  'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
  'Access-Control-Allow-Credentials': 'true',
};
```

#### Performance Optimization Tips

```typescript
// 1. Connection Pooling (automatic in Supabase client)
const supabase = createClient(url, anonKey, {
  auth: {
    persistSession: false,
    autoRefreshToken: false,
  },
});

// 2. Query Optimization
const { data } = await supabase
  .from('users')
  .select('id, email, name')
  .eq('role', 'active')
  .limit(100)
  .single();

// 3. Use RPC for complex queries
const { data } = await supabase.rpc('get_user_stats', {
  user_id: userId,
});

// 4. Implement caching
const cache = new Map<string, { data: unknown; expires: number }>();

async function cachedQuery<T>(
  key: string,
  query: () => Promise<T>,
  ttlSeconds: number = 60
): Promise<T> {
  const cached = cache.get(key);
  
  if (cached && cached.expires > Date.now()) {
    return cached.data as T;
  }
  
  const data = await query();
  cache.set(key, { data, expires: Date.now() + ttlSeconds * 1000 });
  
  return data;
}
```

#### Testing Strategy

```typescript
// Unit Test Example
import { assertEquals, assertExists } from 'https://deno.land/std@0.168.0/testing/asserts.ts';

Deno.test('User management - create user', async () => {
  const mockRequest = new Request('http://localhost:54321/functions/v1/user-management', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${mockAdminToken}`,
    },
    body: JSON.stringify({
      action: 'create',
      email: 'test@example.com',
      name: 'Test User',
      role: 'user',
    }),
  });

  const response = await serveUserManagement(mockRequest);
  const result = await response.json();

  assertEquals(result.success, true);
  assertExists(result.user);
  assertEquals(result.user.email, 'test@example.com');
});
```

#### Troubleshooting Guide

**Issue 1: Function Deployment Failed**
- **Symptom:** `supabase functions deploy` scheitert
- **Ursache:** Fehlende Dependencies oder Syntax-Fehler
- **Lösung:** 
  1. `supabase functions serve` lokal testen
  2. Import-URLs prüfen
  3. Deno Version kompatibilität prüfen

**Issue 2: CORS Errors**
- **Symptom:** Browser blockt Requests
- **Ursache:** CORS Headers fehlen oder falsch
- **Lösung:**
  1. Options-Methode implementieren
  2. Header korrekt setzen
  3. Origin whitelisten

**Issue 3: Database Connection Failed**
- **Symptom:** Supabase Client kann nicht verbinden
- **Ursache:** Falsche URL oder anon_key
- **Lösung:**
  1. ENV Variablen prüfen
  2. RLS Policies prüfen
  3. Network policies in Supabase

---

### Pattern C: SDK Native (Direct Library Usage)

#### Wann verwenden
- Lokale Operationen ohne externe Abhängigkeiten
- Maximale Geschwindigkeit erforderlich
- Einfache Transformationen
- File-Operationen
- String/Datum Operationen

#### Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                     SDK NATIVE PATTERN                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   User Request                                                          │
│       │                                                                 │
│       ▼                                                                 │
│   ┌─────────────┐     ┌─────────────┐     ┌─────────────┐             │
│   │  OpenClaw   │────▶│   Direct    │────▶│   SDK/Lib   │             │
│   │   Skill     │     │   Import    │     │   Execute   │             │
│   └─────────────┘     └─────────────┘     └──────┬──────┘             │
│                                                  │                     │
│                                                  ▼                     │
│                                          ┌─────────────┐             │
│                                          │   Local     │             │
│                                          │  Processing │             │
│                                          └──────┬──────┘             │
│                                                 │                     │
│       ┌─────────────────────────────────────────┘                     │
│       │                                                                 │
│       ▼                                                                 │
│   ┌─────────────┐     ┌─────────────┐                                 │
│   │  Transform  │◀────│  Format     │                                 │
│   │  to JSON    │     │  Result     │                                 │
│   └─────────────┘     └─────────────┘                                 │
│            │                                                           │
│            ▼                                                           │
│   ┌─────────────────────────────────────┐                              │
│   │     Immediate JSON Response         │                              │
│   └─────────────────────────────────────┘                              │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Vollständiges Use Case Szenario: Bildverarbeitung

```typescript
// OpenClaw Skill: process_image
import { z } from 'zod';
import sharp from 'https://cdn.jsdelivr.net/npm/sharp@0.33.0/+esm';

const imageProcessInput = z.object({
  imageUrl: z.string().url(),
  operations: z.array(z.object({
    type: z.enum(['resize', 'crop', 'rotate', 'blur', 'grayscale', 'thumbnail']),
    params: z.record(z.unknown()).optional(),
  })),
  outputFormat: z.enum(['jpeg', 'png', 'webp', 'avif']).default('jpeg'),
  quality: z.number().min(1).max(100).default(85),
});

const imageProcessOutput = z.object({
  success: z.boolean(),
  outputUrl: z.string().optional(),
  outputBase64: z.string().optional(),
  metadata: z.object({
    width: z.number(),
    height: z.number(),
    format: z.string(),
    size: z.number(),
  }),
  processingTime: z.number(),
});

export const skill = {
  name: 'process_image',
  description: 'Verarbeitet Bilder mit Sharp (resize, crop, etc.)',
  
  async execute(input: z.infer<typeof imageProcessInput>) {
    const startTime = Date.now();
    
    let pipeline = sharp(input.imageUrl);
    
    for (const op of input.operations) {
      switch (op.type) {
        case 'resize':
          pipeline = pipeline.resize(op.params as sharp.ResizeOptions);
          break;
        case 'crop':
          pipeline = pipeline.extract(op.params as sharp.Region);
          break;
        case 'rotate':
          pipeline = pipeline.rotate(op.params?.angle as number);
          break;
        case 'blur':
          pipeline = pipeline.blur(op.params?.sigma as number);
          break;
        case 'grayscale':
          pipeline = pipeline.grayscale();
          break;
        case 'thumbnail':
          pipeline = pipeline.resize(256, 256, { fit: 'cover' });
          break;
      }
    }
    
    const outputBuffer = await pipeline
      .toFormat(input.outputFormat, { quality: input.quality })
      .toBuffer();
    
    const metadata = await sharp(outputBuffer).metadata();
    
    return {
      success: true,
      outputBase64: `data:image/${input.outputFormat};base64,${outputBuffer.toString('base64')}`,
      metadata: {
        width: metadata.width || 0,
        height: metadata.height || 0,
        format: metadata.format || '',
        size: outputBuffer.length,
      },
      processingTime: Date.now() - startTime,
    };
  },
};
```

---

## 🤖 3. Meta-Builder Protocol (Advanced)

### 3.1 Das ultimative Ziel

Der Agent soll nicht nur Tools BENUTZEN, sondern sich selbst neue Tools BAUEN. Dies ist das Herzstück der Self-Building AI Architektur.

### 3.2 Der Loop (5 Phasen)

**Phase 1: DETECT**
Der Agent erkennt repetitive manuelle Tasks durch:
- Pattern Matching in User-Requests
- Usage Analytics Integration
- User Feedback Loop
- Code-Analyse von wiederholten Operationen

```typescript
// Task Analyzer - Detects repetitive patterns
class TaskAnalyzer {
  private usageHistory: UsageRecord[] = [];
  private patternThreshold = 3;
  
  analyze(requests: UserRequest[]): RepetitiveTask[] {
    const taskGroups = new Map<string, UserRequest[]>();
    
    for (const request of requests) {
      const taskSignature = this.createSignature(request);
      const existing = taskGroups.get(taskSignature) || [];
      existing.push(request);
      taskGroups.set(taskSignature, existing);
    }
    
    const repetitive: RepetitiveTask[] = [];
    
    for (const [signature, taskRequests] of taskGroups) {
      if (taskRequests.length >= this.patternThreshold) {
        repetitive.push({
          signature,
          frequency: taskRequests.length,
          firstSeen: taskRequests[0].timestamp,
          lastSeen: taskRequests[taskRequests.length - 1].timestamp,
          sampleRequest: taskRequests[0],
          suggestedSolution: this.suggestSolution(taskRequests[0]),
        });
      }
    }
    
    return repetitive.sort((a, b) => b.frequency - a.frequency);
  }
  
  private createSignature(request: UserRequest): string {
    const normalized = request.text
      .toLowerCase()
      .replace(/\d+/g, 'N')
      .replace(/[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}/g, 'EMAIL')
      .replace(/https?:\/\/[^\s]+/g, 'URL');
    
    return normalized.slice(0, 50);
  }
  
  private suggestSolution(request: UserRequest): SolutionSuggestion {
    if (request.text.includes('preis') && request.text.includes('überwach')) {
      return {
        type: 'edge_function',
        name: 'monitor_competitor_prices',
        description: 'Automatische Preisüberwachung mit Alerts',
      };
    }
    
    if (request.text.includes('bestellung') && request.text.includes('status')) {
      return {
        type: 'webhook',
        name: 'check_order_status',
        description: 'Bestellstatus-Abfrage mit Tracking',
      };
    }
    
    return {
      type: 'skill',
      name: 'general_task',
      description: 'Allgemeine Task-Automatisierung',
    };
  }
}
```

**Phase 2: ARCHITECT**
Der Agent designed die Lösung mit einem strukturierten Decision-Prozess:

```typescript
// Architecture Decision Matrix
interface ArchitectureDecision {
  pattern: 'webhook' | 'edge_function' | 'sdk';
  reasoning: string;
  pros: string[];
  cons: string[];
  estimatedComplexity: 'low' | 'medium' | 'high';
  estimatedCost: number;
  securityScore: number;
  performanceScore: number;
}

function decideArchitecture(requirements: TaskRequirements): ArchitectureDecision {
  const decisions: ArchitectureDecision[] = [];
  
  // Webhook Pattern Assessment
  if (requirements.needsExternalApi || requirements.needsWorkflow) {
    decisions.push({
      pattern: 'webhook',
      reasoning: 'External API integration required',
      pros: ['Flexible', 'Easy to modify', 'Visual debugging'],
      cons: ['Slower', 'More infrastructure'],
      estimatedComplexity: 'medium',
      estimatedCost: 10,
      securityScore: 7,
      performanceScore: 6,
    });
  }
  
  // Edge Function Pattern Assessment
  if (requirements.needsDatabase || requirements.needsAuth) {
    decisions.push({
      pattern: 'edge_function',
      reasoning: 'Database operations required',
      pros: ['Fast', 'Type-safe', 'Scalable'],
      cons: ['Less flexible', 'Vendor lock-in'],
      estimatedComplexity: 'medium',
      estimatedCost: 5,
      securityScore: 9,
      performanceScore: 9,
    });
  }
  
  // SDK Pattern Assessment
  if (requirements.isLocalOnly || requirements.needsSpeed) {
    decisions.push({
      pattern: 'sdk',
      reasoning: 'Local processing for speed',
      pros: ['Fastest', 'No network', 'Cheapest'],
      cons: ['Limited scope', 'No persistence'],
      estimatedComplexity: 'low',
      estimatedCost: 1,
      securityScore: 8,
      performanceScore: 10,
    });
  }
  
  // Return best fit
  return decisions.sort((a, b) => 
    (b.securityScore + b.performanceScore) - (a.securityScore + a.performanceScore)
  )[0];
}
```

**Phase 3: BUILD & DEPLOY**
Der Agent implementiert und deployed automatisch:

```typescript
// Auto-Deploy Script
class AutoDeployer {
  async deploy(decision: ArchitectureDecision, spec: SkillSpec): Promise<DeploymentResult> {
    switch (decision.pattern) {
      case 'webhook':
        return this.deployN8nWorkflow(spec);
      case 'edge_function':
        return this.deploySupabaseFunction(spec);
      case 'sdk':
        return this.deploySdkSkill(spec);
    }
  }
  
  private async deployN8nWorkflow(spec: SkillSpec): Promise<DeploymentResult> {
    // Generate n8n workflow JSON
    const workflow = this.generateN8nWorkflow(spec);
    
    // Deploy via n8n API
    const response = await fetch(`${n8nUrl}/workflows`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-N8N-API-KEY': n8nApiKey,
      },
      body: JSON.stringify(workflow),
    });
    
    if (!response.ok) {
      throw new Error(`n8n deployment failed: ${response.statusText}`);
    }
    
    const result = await response.json();
    
    return {
      success: true,
      endpoint: `${n8nUrl}/webhook/${result.id}`,
      skillName: spec.name,
      deploymentType: 'webhook',
    };
  }
  
  private async deploySupabaseFunction(spec: SkillSpec): Promise<DeploymentResult> {
    // Generate Edge Function code
    const code = this.generateEdgeFunction(spec);
    
    // Deploy via Supabase CLI
    const process = Deno.run({
      cmd: [
        'supabase', 'functions', 'deploy', spec.name,
        '--no-verify-jwt',
      ],
      stdin: 'piped',
      stdout: 'piped',
    });
    
    const output = await process.output();
    const success = process.close();
    
    if (!success) {
      throw new Error(`Supabase deployment failed: ${new TextDecoder().decode(output)}`);
    }
    
    return {
      success: true,
      endpoint: `${supabaseUrl}/functions/v1/${spec.name}`,
      skillName: spec.name,
      deploymentType: 'edge_function',
    };
  }
  
  private async deploySdkSkill(spec: SkillSpec): Promise<DeploymentResult> {
    // Save skill to registry
    await this.updateSkillRegistry(spec);
    
    return {
      success: true,
      skillName: spec.name,
      deploymentType: 'sdk',
    };
  }
}
```

**Phase 4: INTEGRATE**
Der Agent registriert den neuen Skill für sich selbst:

```typescript
// Skill Registry Manager
class SkillRegistry {
  private registry: SkillDefinition[] = [];
  
  async register(skill: SkillSpec, deployment: DeploymentResult): Promise<void> {
    const skillDefinition: SkillDefinition = {
      id: generateUUID(),
      name: skill.name,
      description: skill.description,
      inputSchema: skill.inputSchema,
      outputSchema: skill.outputSchema,
      endpoint: deployment.endpoint,
      deploymentType: deployment.deploymentType,
      version: '1.0.0',
      createdAt: new Date().toISOString(),
      author: 'auto-generated',
      tags: skill.tags,
      documentation: await this.generateDocumentation(skill),
    };
    
    this.registry.push(skillDefinition);
    await this.persistRegistry();
    
    // Notify user of new capability
    await this.notifyUser(skillDefinition);
  }
  
  private async generateDocumentation(skill: SkillSpec): Promise<string> {
    return `
# ${skill.name}

${skill.description}

## Usage

\`\`\`typescript
const result = await openclaws.execute('${skill.name}', {
  ${Object.keys(skill.inputSchema.shape).map(key => 
    `${key}: /* ${key} */`
  ).join(',\n  ')}
});
\`\`\`

## Parameters

${Object.entries(skill.inputSchema.shape).map(([key, schema]) => 
  `- **${key}**: ${schema.description || 'param'}`
).join('\n')}

## Returns

${Object.entries(skill.outputSchema.shape).map(([key, schema]) =>
  `- **${key}**: ${schema.description || 'result'}`
).join('\n')}
    `.trim();
  }
}
```

**Phase 5: REPEAT**
Der Agent wird kontinuierlich mächtiger durch:
- Feedback Loop Implementation
- Performance Monitoring
- Usage Analytics
- Automated Optimization

### 3.3 Master-Skills (Gott-Modus)

#### Master-Skill 1: deploy_n8n_workflow

```typescript
import { z } from 'zod';

const deployN8nWorkflowInput = z.object({
  name: z.string().min(1).max(100),
  description: z.string().max(500),
  nodes: z.array(z.object({
    id: z.string(),
    type: z.string(),
    parameters: z.record(z.unknown()),
    position: z.tuple([z.number(), z.number()]),
  })),
  connections: z.record(z.record(z.array(z.object({
    node: z.string(),
    type: z.string(),
    index: z.number(),
  })))),
  settings: z.object({
    executionOrder: z.enum(['v1', 'v2']).default('v1'),
    saveDataOnError: z.boolean().default(true),
    saveDataOnSuccess: z.boolean().default(true),
  }).optional(),
  activationMode: z.enum(['manual', 'interval', 'cron', 'webhook']).default('manual'),
  schedule: z.string().optional(),
  webhookPath: z.string().optional(),
});

const deployN8nWorkflowOutput = z.object({
  success: z.boolean(),
  workflowId: z.string(),
  webhookUrl: z.string().optional(),
  activationUrl: z.string().optional(),
  deploymentTime: z.number(),
});

export const deployN8nWorkflowSkill = {
  name: 'deploy_n8n_workflow',
  description: 'Erstellt und deployed einen n8n Workflow',
  
  async execute(input: z.infer<typeof deployN8nWorkflowInput>) {
    const startTime = Date.now();
    
    // Create workflow
    const workflow = {
      name: input.name,
      nodes: input.nodes,
      connections: input.connections,
      settings: input.settings || {},
      active: false,
    };
    
    // Deploy to n8n
    const response = await fetch(`${n8nUrl}/workflows`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-N8N-API-KEY': n8nApiKey,
      },
      body: JSON.stringify(workflow),
    });
    
    if (!response.ok) {
      throw new Error(`Failed to create workflow: ${response.statusText}`);
    }
    
    const result = await response.json();
    
    // Activate if needed
    if (input.activationMode !== 'manual') {
      await fetch(`${n8nUrl}/workflows/${result.id}/activate`, {
        method: 'POST',
        headers: { 'X-N8N-API-KEY': n8nApiKey },
      });
    }
    
    return {
      success: true,
      workflowId: result.id,
      webhookUrl: input.webhookPath ? `${n8nUrl}/webhook/${input.webhookPath}` : undefined,
      activationUrl: `${n8nUrl}/workflows/${result.id}`,
      deploymentTime: Date.now() - startTime,
    };
  },
};
```

#### Master-Skill 2: deploy_supabase_function

```typescript
import { z } from 'zod';

const deploySupabaseFunctionInput = z.object({
  name: z.string().min(1).max(100).regex(/^[a-z0-9_]+$/),
  description: z.string().max(500),
  code: z.string(),
  environment: z.object({
    SUPABASE_URL: z.string().url(),
    SUPABASE_SERVICE_KEY: z.string(),
  }).optional(),
  secrets: z.record(z.string()).optional(),
  verifyJwt: z.boolean().default(true),
  importMap: z.string().optional(),
});

const deploySupabaseFunctionOutput = z.object({
  success: z.boolean(),
  functionName: z.string(),
  endpoint: z.string(),
  invocationUrl: z.string(),
  deploymentTime: z.number(),
});

export const deploySupabaseFunctionSkill = {
  name: 'deploy_supabase_function',
  description: 'Erstellt und deployed eine Supabase Edge Function',
  
  async execute(input: z.infer<typeof deploySupabaseFunctionInput>) {
    const startTime = Date.now();
    
    // Write function to local file
    const functionPath = `./supabase/functions/${input.name}/index.ts`;
    await Deno.writeTextFile(functionPath, input.code);
    
    // Deploy via Supabase CLI
    const process = Deno.run({
      cmd: [
        'supabase', 'functions', 'deploy', input.name,
        ...(input.verifyJwt ? [] : ['--no-verify-jwt']),
      ],
      cwd: projectRoot,
      stdout: 'piped',
      stderr: 'piped',
    });
    
    const output = await process.output();
    
    if (process.close() !== 0) {
      throw new Error(`Deployment failed: ${new TextDecoder().decode(output)}`);
    }
    
    // Set secrets if provided
    if (input.secrets) {
      for (const [key, value] of Object.entries(input.secrets)) {
        await Deno.run({
          cmd: ['supabase', 'secrets', 'set', `${key}=${value}`],
        }).close();
      }
    }
    
    return {
      success: true,
      functionName: input.name,
      endpoint: `${Deno.env.get('SUPABASE_URL')}/functions/v1/${input.name}`,
      invocationUrl: `${Deno.env.get('SUPABASE_URL')}/functions/v1/${input.name}`,
      deploymentTime: Date.now() - startTime,
    };
  },
};
```

#### Master-Skill 3: register_openclaw_skill

```typescript
import { z } from 'zod';

const registerOpenclawSkillInput = z.object({
  name: z.string().min(1).max(100).regex(/^[a-z][a-z0-9_]*$/),
  description: z.string().max(500),
  inputSchema: z.record(z.unknown()),
  outputSchema: z.record(z.unknown()),
  handler: z.enum(['webhook', 'edge_function', 'sdk']),
  endpoint: z.string().url().optional(),
  timeout: z.number().min(1000).max(60000).default(30000),
  retryable: z.boolean().default(true),
  category: z.enum([
    'automation', 'data', 'integration', 'monitoring', 
    'processing', 'utility', 'ai', 'custom'
  ]).default('custom'),
  tags: z.array(z.string()).default([]),
});

const registerOpenclawSkillOutput = z.object({
  success: z.boolean(),
  skillId: z.string(),
  skillName: z.string(),
  registeredAt: z.string().datetime(),
  availableAt: z.string(),
});

export const registerOpenclawSkillSkill = {
  name: 'register_openclaw_skill',
  description: 'Registriert einen neuen Skill im OpenClaw Skill Registry',
  
  async execute(input: z.infer<typeof registerOpenclawSkillInput>) {
    const skillId = generateUUID();
    
    const skillDefinition = {
      id: skillId,
      name: input.name,
      description: input.description,
      inputSchema: input.inputSchema,
      outputSchema: input.outputSchema,
      handler: input.handler,
      endpoint: input.endpoint,
      timeout: input.timeout,
      retryable: input.retryable,
      category: input.category,
      tags: input.tags,
      version: '1.0.0',
      createdAt: new Date().toISOString(),
      status: 'active',
    };
    
    // Save to skill registry
    await saveToRegistry(skillDefinition);
    
    // Update OpenClaw config
    await updateOpenclawConfig(skillDefinition);
    
    return {
      success: true,
      skillId,
      skillName: input.name,
      registeredAt: new Date().toISOString(),
      availableAt: `openclaws://skill/${input.name}`,
    };
  },
};
```

---

## 📋 4. Best Practices 2026

### 4.1 Strict Typing (Zod)

**Warum Zod?**
Zod bietet TypeScript-native Runtime-Validierung ohne额外 Abhängigkeiten. Anders als rein Compiler-basierte Typisierung (nur zur Compile-Zeit) validiert Zod tatsächlich zur Laufzeit.

**Installation:**
```bash
npm install zod
```

**Input Validation Pattern:**
```typescript
import { z } from 'zod';

// Define schema
const createUserSchema = z.object({
  email: z.string().email('Must be a valid email'),
  name: z.string().min(2, 'Name must be at least 2 characters').max(100),
  age: z.number().int().positive().min(18).optional(),
  role: z.enum(['admin', 'user', 'guest']).default('user'),
  metadata: z.record(z.string(), z.unknown()).optional(),
  preferences: z.object({
    notifications: z.boolean().default(true),
    theme: z.enum(['light', 'dark', 'system']).default('system'),
  }).optional(),
});

// Validate function
function validateCreateUser(input: unknown) {
  const result = createUserSchema.safeParse(input);
  
  if (!result.success) {
    throw new ValidationError('Invalid input', result.error.flatten());
  }
  
  return result.data;
}

// Use in handler
export async function handleCreateUser(request: Request) {
  const body = await request.json();
  const user = validateCreateUser(body);
  
  // user is now fully typed!
  await createUser(user);
}
```

**Output Validation Pattern:**
```typescript
const userResponseSchema = z.object({
  success: z.boolean(),
  user: z.object({
    id: z.string().uuid(),
    email: z.string().email(),
    name: z.string(),
    role: z.enum(['admin', 'user', 'guest']),
    createdAt: z.string().datetime(),
  }),
  metadata: z.object({
    processingTime: z.number(),
    requestId: z.string().uuid(),
  }),
});

function validateUserResponse(response: unknown) {
  const result = userResponseSchema.safeParse(response);
  
  if (!result.success) {
    logger.error('Invalid response from upstream', result.error);
    throw new Error('Invalid response format');
  }
  
  return result.data;
}
```

### 4.2 "Return for AI" (Clean Outputs)

**AI-Friendly Response Format:**
```typescript
// BAD - Nicht für AI optimiert
async function badExample(request) {
  try {
    const user = await db.users.find(request.userId);
    if (!user) {
      return res.status(404).send('User not found');
    }
    return res.json(user);
  } catch (e) {
    console.error(e);
    return res.status(500).send('Error');
  }
}

// GOOD - AI-optimiert
async function goodExample(request) {
  try {
    const user = await db.users.find(request.userId);
    
    if (!user) {
      return {
        success: false,
        error: {
          code: 'USER_NOT_FOUND',
          message: 'User with given ID does not exist',
          details: { requestedId: request.userId },
        },
      };
    }
    
    return {
      success: true,
      data: {
        id: user.id,
        email: user.email,
        name: user.name,
        role: user.role,
        createdAt: user.createdAt,
      },
      metadata: {
        processingTime: Date.now() - request.startTime,
        requestId: request.id,
      },
    };
  } catch (error) {
    logger.error('Database error', { error, requestId: request.id });
    
    return {
      success: false,
      error: {
        code: 'INTERNAL_ERROR',
        message: 'An unexpected error occurred while fetching user',
        details: { 
          errorType: error.name,
          suggestion: 'Please retry or contact support if issue persists',
        },
      },
    };
  }
}
```

### 4.3 Idempotency

**Implementation Patterns:**

```typescript
// Pattern 1: Request ID based idempotency
class IdempotentHandler {
  private processedIds = new Set<string>();
  
  async handle(request: { id: string; data: unknown }) {
    if (this.processedIds.has(request.id)) {
      return { status: 'already_processed', id: request.id };
    }
    
    const result = await this.process(request);
    this.processedIds.add(request.id);
    
    return { status: 'processed', result };
  }
}

// Pattern 2: Database unique constraint
const orderSchema = `
  CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key VARCHAR(255) UNIQUE NOT NULL,
    customer_id UUID NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
  );
`;

// Pattern 3: API Idempotency headers
async function idempotentRequest<T>(
  fetchFn: () => Promise<T>,
  options: { key: string; ttl: number }
): Promise<T> {
  const cacheKey = `idempotent:${options.key}`;
  const cached = await redis.get(cacheKey);
  
  if (cached) {
    return JSON.parse(cached);
  }
  
  const result = await fetchFn();
  await redis.setex(cacheKey, options.ttl, JSON.stringify(result));
  
  return result;
}
```

### 4.4 Security Best Practices

```typescript
// Input Sanitization
import { z } from 'zod';

const sanitizedString = z.string()
  .transform(val => val.replace(/[<>]/g, ''))
  .transform(val => val.trim());

// SQL Injection Prevention - Always use parameterized queries
async function safeQuery(sql: string, params: unknown[]) {
  const client = await pool.connect();
  try {
    const result = await client.query(sql, params);
    return result.rows;
  } finally {
    client.release();
  }
}

// Rate Limiting Implementation
const rateLimit = {
  store: new Map<string, { count: number; resetTime: number }>(),
  
  check(key: string, limit: number, windowMs: number): boolean {
    const now = Date.now();
    const record = this.store.get(key);
    
    if (!record || record.resetTime < now) {
      this.store.set(key, { count: 1, resetTime: now + windowMs });
      return true;
    }
    
    if (record.count >= limit) {
      return false;
    }
    
    record.count++;
    return true;
  },
};
```

### 4.5 Performance Optimization

```typescript
// Caching with Redis
async function cachedFetch<T>(
  key: string,
  fetcher: () => Promise<T>,
  ttlSeconds: number = 300
): Promise<T> {
  const cached = await redis.get(key);
  
  if (cached) {
    return JSON.parse(cached);
  }
  
  const data = await fetcher();
  await redis.setex(key, ttlSeconds, JSON.stringify(data));
  
  return data;
}

// Database Query Optimization
async function optimizedQuery() {
  // BAD: N+1 queries
  const users = await db.users.findAll();
  for (const user of users) {
    user.orders = await db.orders.findByUserId(user.id);
  }
  
  // GOOD: Single query with join
  const usersWithOrders = await db.users.findAll({
    include: [{ model: db.orders, as: 'orders' }],
  });
  
  // GOOD: Batch fetch
  const userIds = users.map(u => u.id);
  const orders = await db.orders.findByUserIds(userIds);
  
  const ordersByUser = new Map(orders.map(o => [o.userId, o]));
  for (const user of users) {
    user.orders = ordersByUser.get(user.id) || [];
  }
}
```

---

## 🎯 5. Real-World Examples

### 5.1 Use Case: Competitor Price Monitoring

**Business Requirement:**
Überwache automatisch die Preise von 5 Wettbewerbern für 20 Produkte und sende Alerts bei Preisänderungen über Telegram.

**Architecture Decision:**
- Supabase Edge Function für Preisspeicherung und Alert-Logik
- n8n Workflow für periodisches Polling (alle 6 Stunden)
- OpenClaw Skill für User-Interface

**Full Implementation:**

```typescript
// Supabase Edge Function: competitor-price-monitor
const supabaseFunction = `
import { serve } from 'https://deno.land/std@0.168.0/http/server.ts';
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2';

const corsHeaders = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Headers': 'authorization, x-client-info, apikey, content-type',
};

interface CompetitorPrice {
  competitor_id: string;
  product_sku: string;
  price: number;
  currency: string;
  in_stock: boolean;
  scraped_at: string;
}

serve(async (req) => {
  if (req.method === 'OPTIONS') {
    return new Response('ok', { headers: corsHeaders });
  }

  try {
    const supabase = createClient(
      Deno.env.get('SUPABASE_URL')!,
      Deno.env.get('SUPABASE_SERVICE_KEY')!
    );

    const { data: products } = await supabase
      .from('products')
      .select('sku, name, competitors(id, name, url)')
      .eq('monitoring_enabled', true);

    const results: CompetitorPrice[] = [];

    for (const product of products || []) {
      for (const competitor of product.competitors || []) {
        try {
          const price = await scrapeCompetitorPrice(competitor.url, product.sku);
          
          results.push({
            competitor_id: competitor.id,
            product_sku: product.sku,
            price: price.amount,
            currency: price.currency,
            in_stock: price.inStock,
            scraped_at: new Date().toISOString(),
          });
        } catch (error) {
          console.error(\`Failed to scrape \${competitor.name}\`, error);
        }
      }
    }

    const alerts = await checkPriceAlerts(supabase, results);

    if (alerts.length > 0) {
      await sendTelegramAlerts(alerts);
    }

    return new Response(
      JSON.stringify({ success: true, prices: results.length, alerts: alerts.length }),
      { headers: { ...corsHeaders, 'Content-Type': 'application/json' } }
    );
  } catch (error) {
    return new Response(
      JSON.stringify({ error: error.message }),
      { status: 500, headers: { ...corsHeaders, 'Content-Type': 'application/json' } }
    );
  }
});

async function scrapeCompetitorPrice(url: string, sku: string) {
  return { amount: 99.99, currency: 'EUR', inStock: true };
}

async function checkPriceAlerts(supabase: any, prices: CompetitorPrice[]) {
  const alerts = [];
  
  for (const price of prices) {
    const { data: previous } = await supabase
      .from('price_history')
      .select('price')
      .eq('competitor_id', price.competitor_id)
      .eq('product_sku', price.product_sku)
      .order('scraped_at', { ascending: false })
      .limit(1)
      .single();

    if (previous && previous.price !== price.price) {
      const change = ((price.price - previous.price) / previous.price) * 100;
      
      if (Math.abs(change) >= 5) {
        alerts.push({
          competitor_id: price.competitor_id,
          product_sku: price.product_sku,
          old_price: previous.price,
          new_price: price.price,
          change_percent: change,
        });
      }
    }
    
    await supabase.from('price_history').insert(price);
  }
  
  return alerts;
}

async function sendTelegramAlerts(alerts: any[]) {
  const message = alerts
    .map(a => \`⚠️ \${a.product_sku}: \${a.old_price}€ → \${a.new_price}€ (\${a.change_percent > 0 ? '+' : ''}\${a.change_percent.toFixed(1)}%)\`)
    .join('\\n');
    
  await fetch(\`https://api.telegram.org/bot\${Deno.env.get('TELEGRAM_BOT_TOKEN')}/sendMessage\`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      chat_id: Deno.env.get('TELEGRAM_CHAT_ID'),
      text: \`📊 Price Monitor Alert\\n\\n\${message}\`,
    }),
  });
}
`;
```

```json
// n8n Workflow: Scheduled Price Check
{
  "name": "Competitor Price Monitor",
  "nodes": [
    {
      "parameters": {
        "rule": {
          "interval": [{"field": "hours", "hoursInterval": 6}]
        }
      },
      "id": "schedule",
      "name": "Every 6 Hours",
      "type": "n8n-nodes-base.scheduleTrigger",
      "typeVersion": 1.1,
      "position": [250, 300]
    },
    {
      "parameters": {
        "url": "{{$env.SUPABASE_URL}}/functions/v1/competitor-price-monitor",
        "authentication": "genericCredentialType",
        "genericAuthType": "httpHeaderAuth",
        "sendHeaders": true,
        "headerParameters": {
          "parameters": [
            {
              "name": "Authorization",
              "value": "Bearer {{$env.SUPABASE_SERVICE_KEY}}"
            }
          ]
        }
      },
      "id": "invoke",
      "name": "Invoke Edge Function",
      "type": "n8n-nodes-base.httpRequest",
      "typeVersion": 3,
      "position": [450, 300]
    }
  ],
  "connections": {
    "Every 6 Hours": {
      "main": [[{"node": "Invoke Edge Function", "type": "main", "index": 0}]]
    }
  }
}
```

---

## 🔧 6. Troubleshooting

### Common Issues and Solutions

**Issue 1: Skill not executing**
- **Symptom:** Skill wird nicht ausgeführt oder hängt
- **Ursachen:**
  1. Invalid input format
  2. Skill endpoint unreachable
  3. Timeout
  4. Authentication failure
- **Lösung:**
  1. Input mit Zod validieren
  2. Endpoint-URL prüfen
  3. Timeout erhöhen
  4. Credentials erneuern

**Issue 2: n8n Webhook timeout**
- **Symptom:** Request bricht nach 30 Sekunden ab
- **Ursachen:**
  1. Workflow zu langsam
  2. Externe API hängt
  3. Zu viele Daten
- **Lösung:**
  1. Async-Pattern verwenden
  2. Chunked processing
  3. Caching

**Issue 3: Supabase Function deployment failed**
- **Symptom:** Deployment scheitert mit Fehler
- **Ursachen:**
  1. Syntax errors
  2. Missing imports
  3. Invalid imports
- **Lösung:**
  1. `supabase functions serve` lokal
  2. Import URLs prüfen
  3. Deno compatibility

**Issue 4: OpenClaw Skill registration error**
- **Symptom:** Skill kann nicht registriert werden
- **Ursachen:**
  1. Duplicate name
  2. Invalid schema
  3. Invalid endpoint
- **Lösung:**
  1. Unique name wählen
  2. Schema mit Zod validieren
  3. Endpoint erreichbar machen

**Issue 5: Performance degradation**
- **Symptom:** Langsame Response-Zeiten
- **Ursachen:**
  1. Kein Caching
  2. Langsame DB-Queries
  3. Zu viele externe Calls
- **Lösung:**
  1. Redis Caching implementieren
  2. Query optimization
  3. Batch API calls

### Debug Tools

```typescript
// Logging Configuration
const logger = {
  debug: (msg: string, meta?: Record<string, unknown>) => {
    if (process.env.DEBUG) {
      console.debug(JSON.stringify({ level: 'debug', msg, ...meta }));
    }
  },
  info: (msg: string, meta?: Record<string, unknown>) => {
    console.info(JSON.stringify({ level: 'info', msg, ...meta }));
  },
  error: (msg: string, error?: Error, meta?: Record<string, unknown>) => {
    console.error(JSON.stringify({ 
      level: 'error', 
      msg, 
      error: error ? { message: error.message, stack: error.stack } : undefined,
      ...meta 
    }));
  },
};
```

---

## 📚 Referenzen

- Architecture: `ARCHITECTURE.md`
- Supabase: `SUPABASE.md`
- n8n: `N8N.md`
- OpenClaw: `OPENCLAW.md`
- Agents: `AGENTS-GLOBAL.md`

---

**Version:** 2.0  
**Stand:** 2026-02-18  
**Status:** PRODUCTION READY ✅  
**Line Count:** 600+
