# WORKFLOW.md — Unified Skill Architecture & Meta-Builder Protocol

**Status:** ACTIVE
**Version:** 2.0
**Stand:** Februar 2026
**Purpose:** Zentrale Dokumentation der Self-Building AI Agent Architektur
**Line Count:** 600+ lines

---

## 🧠 NLM CLI COMMANDS

```bash
# Create notebook
nlm notebook create "Title"

# List sources
nlm source list <notebook-id>

# Delete old source (BEFORE adding new!)
nlm source delete <source-id> -y

# Add new source
nlm source add <notebook-id> --file "file.md" --wait
```

**⚠️ DUPLICATE PREVENTION:** ALWAYS run `nlm source list` before `nlm source add`!

---

## 🔄 DEQLHI-LOOP (INFINITE WORK MODE)

- After each completed task → Add 5 new tasks immediately
- Never "done" - only "next task"
- Always document → Every change in files
- Git commit + push after EVERY change
- Parallel execution ALWAYS (run_in_background=true)

### Loop Mechanism:
1. Task N Complete
2. Git Commit + Push
3. Update Docs
4. Add 5 New Tasks
5. Next Task N+1
6. Repeat infinitely

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

## 📋 13. Development Workflows

### 13.1 Git Branch Strategy

#### Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        GIT BRANCH STRATEGY                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│    main (production) ─────────────────────────────────────────────     │
│         │                                                              │
│         │ Merge via PR                                                │
│         ▼                                                              │
│    develop (integration) ──────────────────────────────────────       │
│         │                                                              │
│         │ Merge via PR                                                │
│         ▼                                                              │
│    feature/* ────────────────────────┬───────────────────────────     │
│    bugfix/* ─────────────────────────┤                              │
│    hotfix/* ──────────────────────────┘                              │
│                                                                         │
│    release/* ───────────────────────────────────────────────────       │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Branch Naming Convention

```bash
# Feature Branches
feature/TICKET-ID-short-description
feature/add-user-authentication
feature/BIOMETRICS-123-payment-integration

# Bugfix Branches
bugfix/TICKET-ID-short-description
bugfix/BIOMETRICS-456-login-error

# Hotfix Branches
hotfix/TICKET-ID-short-description
hotfix/BIOMETRICS-789-security-patch

# Release Branches
release/2026.02.15
release/v2.0.0
```

#### Branch Workflow Commands

```bash
# 1. Start new feature
git checkout develop
git pull origin develop
git checkout -b feature/BIOMETRICS-123-new-feature

# 2. Work on feature (commit regularly)
git add .
git commit -m "feat: add initial feature implementation"
git push -u origin feature/BIOMETRICS-123-new-feature

# 3. Keep branch updated with develop
git fetch origin
git rebase origin/develop

# 4. When ready, create PR
gh pr create \
  --title "feat: Add new feature" \
  --body "$(cat <<'EOF'
## Summary
- Added new feature for user authentication
- Implemented OAuth2 flow with Google

## Changes
- New endpoint: /api/auth/google
- Updated user table with provider field
- Added OAuth callback handler

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Screenshots
(if applicable)
EOF
)"

# 5. After PR merge, cleanup
git checkout develop
git pull origin develop
git branch -d feature/BIOMETRICS-123-new-feature
git push origin --delete feature/BIOMETRICS-123-new-feature
```

#### Branch Protection Rules

```yaml
# .github/branch-protection.yml
rules:
  - name: main
    required_reviewers: 1
    required_status_checks:
      - ci/lint
      - ci/typecheck
      - ci/test
      - ci/build
    required_signed_commits: false
    allow_force_pushes: false
    allow_deletions: false
    
  - name: develop
    required_reviewers: 0
    required_status_checks:
      - ci/lint
      - ci/test
    allow_force_pushes: false
    allow_deletions: false
```

---

### 13.2 Code Review Process

#### Review Checklist

```markdown
## Code Review Checklist

### Functionality
- [ ] Does this code do what it's supposed to do?
- [ ] Are all edge cases handled?
- [ ] Is the error handling appropriate?
- [ ] Are there any security vulnerabilities?

### Code Quality
- [ ] Is the code readable and well-structured?
- [ ] Are variable/function names descriptive?
- [ ] Is there any duplicated code?
- [ ] Are there appropriate comments for complex logic?

### Performance
- [ ] Are there any N+1 queries?
- [ ] Is caching considered where appropriate?
- [ ] Are there any memory leaks?
- [ ] Is the algorithm efficient?

### Testing
- [ ] Are there adequate unit tests?
- [ ] Are edge cases covered in tests?
- [ ] Do existing tests still pass?
- [ ] Is test coverage maintained/increased?

### Documentation
- [ ] Is the README updated if needed?
- [ ] Are API endpoints documented?
- [ ] Are breaking changes documented?
```

#### Review Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        CODE REVIEW FLOW                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   Developer              Reviewer              System                   │
│      │                     │                     │                      │
│      │  Pushes branch      │                     │                      │
│      │───────────────────▶│                     │                      │
│      │                     │                     │                      │
│      │                     │  Reviews code      │                      │
│      │                     │  Adds comments     │                      │
│      │                     │───────────────────▶│                     │
│      │                     │                     │                      │
│      │  Receives feedback │                     │                      │
│      │◀───────────────────│                     │                      │
│      │                     │                     │                      │
│      │  Addresses comments│                     │                      │
│      │  (if needed)       │                     │                      │
│      │───────────────────▶│                     │                      │
│      │                     │                     │                      │
│      │               Approve                   │                      │
│      │               Merge                     │                      │
│      │───────────────────▶│───────────────────▶│                      │
│      │                     │                     │                      │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Comment Types

```typescript
// Comment severity levels
enum CommentType {
  BLOCKING = 'blocking',      // Must fix before merge
  REQUIRED = 'required',      // Should fix, but not blocking
  SUGGESTION = 'suggestion', // Nice to have
  QUESTION = 'question',      // Needs clarification
  NITPICK = 'nitpick',       // Minor style issue
}

// Example comments
{
  type: CommentType.BLOCKING,
  line: 45,
  message: 'SQL injection vulnerability - use parameterized queries',
  suggestion: 'Replace string concatenation with parameterized query'
}

{
  type: CommentType.REQUIRED,
  line: 78,
  message: 'Missing error handling for network failure',
  suggestion: 'Add try-catch block and retry logic'
}

{
  type: CommentType.SUGGESTION,
  line: 102,
  message: 'Could use const instead of let',
  suggestion: 'const is more appropriate here since value is never reassigned'
}
```

#### Review Metrics

```typescript
interface ReviewMetrics {
  reviewId: string;
  author: string;
  reviewer: string;
  startTime: Date;
  endTime: Date;
  comments: {
    blocking: number;
    required: number;
    suggestion: number;
    question: number;
  };
  reviewDurationMinutes: number;
  codeLinesReviewed: number;
}

// Track average review time
const reviewMetrics = {
  averageReviewTimeMinutes: 45,
  maxReviewTimeHours: 24,
  minReviewers: 1,
  requiredReviewersForSecurity: 2,
};
```

---

### 13.3 PR Template Usage

#### PR Template

```markdown
## 🎯 Pull Request

**Project:** BIOMETRICS
**Ticket:** [BIOMETRICS-XXX](link)
**Type:** [Feature/Bugfix/Hotfix/Refactor/Docs]

---

### 📝 Summary

<!-- What does this PR do? -->

---

### ✅ Changes

<!-- List all changes made -->

- [ ] Change 1
- [ ] Change 2
- [ ] Change 3

---

### 🧪 Testing

**Test Environment:**
- [ ] Local development
- [ ] Staging
- [ ] Production (if applicable)

**Test Results:**
- [ ] Unit tests passing
- [ ] Integration tests passing
- [ ] E2E tests passing
- [ ] Manual testing completed

**Test Coverage:**
```
<!-- Paste coverage report -->
```

---

### 📸 Screenshots (if applicable)

<!-- Add screenshots of changes -->

---

### ⚠️ Breaking Changes

<!-- List any breaking changes -->

None

---

### 📋 Checklist

- [ ] Code follows style guidelines
- [ ] No console.log statements (use proper logging)
- [ ] No hardcoded secrets
- [ ] Documentation updated
- [ ] No merge conflicts
- [ ] All checks passing

---

### 📎 Related Issues

- Related to [BIOMETRICS-XXX](link)
- Blocks [BIOMETRICS-YYY](link)

---

### 🚀 Deployment

**Deploy Instructions:**
<!-- How to deploy this PR -->

---

### 👥 Reviewers

<!-- Tag reviewers -->

---

### 📝 Notes for Reviewers

<!-- Any additional notes -->
```

#### PR Description Examples

**Feature PR:**
```markdown
## Summary
Adds OAuth2 authentication with Google and GitHub providers.

### Changes
- New table: `auth_providers`
- New API: `/api/auth/[provider]`
- New frontend: Login with OAuth buttons

### Testing
- Tested OAuth flow with Google account
- Tested OAuth flow with GitHub account
- Tested error handling for denied permissions

### Screenshots
![Login Page](screenshot.png)
```

**Bugfix PR:**
```markdown
## Summary
Fixes login error when using special characters in password.

### Root Cause
Password validation regex was incorrectly rejecting special characters.

### Fix
Updated validation regex to allow: !@#$%^&*()

### Testing
- Tested with password containing each special character
- Tested with password containing multiple special characters
- Tested with password containing no special characters
```

---

### 13.4 Merge Strategies

#### Merge vs Rebase vs Squash

```typescript
enum MergeStrategy {
  MERGE = 'merge',       // Preserves history, creates merge commit
  REBASE = 'rebase',     // Linear history, rewrites commits
  SQUASH = 'squash',     // Combines all commits into one
}

// When to use each strategy
const mergeStrategyGuide = {
  [MergeStrategy.MERGE]: {
    when: 'When working on shared feature branches',
    pros: [
      'Preserves complete history',
      'Shows when feature was merged',
      'Easier to track down bugs in history'
    ],
    cons: [
      'Can create cluttered history',
      'Merge commits can be confusing'
    ]
  },
  
  [MergeStrategy.REBASE]: {
    when: 'Before merging feature branch into develop',
    pros: [
      'Creates clean, linear history',
      'Easier to follow the story',
      'No unnecessary merge commits'
    ],
    cons: [
      'Rewrites history',
      'Can be dangerous if not done carefully',
      'Loses context of when branch was created'
    ]
  },
  
  [MergeStrategy.SQUASH]: {
    when: 'When branch has many small commits',
    pros: [
      'Clean history',
      'Each PR = one commit',
      'Easy to revert entire feature'
    ],
    cons: [
      'Loses granular commit history',
      'Loses context of development process'
    ]
  }
};
```

#### Merge Workflow

```bash
# Standard merge (develop -> main)
git checkout main
git pull origin main
git merge develop
git push origin main

# Rebase workflow (feature -> develop)
git checkout feature/my-feature
git fetch origin
git rebase origin/develop
# Resolve conflicts if any
git push --force-with-lease origin feature/my-feature

# Squash merge (feature -> main)
git checkout main
git merge --squash feature/my-feature
git commit -m "feat: Add my feature (closes #123)"
git push origin main
```

#### Conflict Resolution

```bash
# 1. Start rebase
git checkout feature/my-feature
git rebase develop

# 2. Conflicts detected - resolve each file
# Edit files to resolve conflicts
git add resolved-file.ts
git add another-file.ts

# 3. Continue rebase
git rebase --continue

# 4. If things go wrong
git rebase --abort

# 5. Push after resolution
git push --force-with-lease
```

---

## 🚀 14. Deployment Workflows

### 14.1 Staging to Production

#### Deployment Pipeline

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    DEPLOYMENT PIPELINE                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐          │
│   │  Code   │───▶│  Build  │───▶│  Test   │───▶│ Deploy  │          │
│   │  Commit │    │         │    │         │    │ Staging │          │
│   └─────────┘    └─────────┘    └─────────┘    └────┬────┘          │
│                                                      │                │
│                        ┌─────────────────────────────┘                │
│                        ▼                                               │
│                  ┌───────────┐    ┌───────────┐    ┌───────────┐     │
│                  │  Manual   │───▶│  Deploy   │───▶│  Verify   │     │
│                  │  Approval │    │ Production │    │           │     │
│                  └───────────┘    └───────────┘    └───────────┘     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Deployment Checklist

```markdown
## Pre-Deployment Checklist

### Code
- [ ] All tests passing
- [ ] Code coverage maintained
- [ ] No linting errors
- [ ] Security scan passed
- [ ] Dependencies updated

### Testing
- [ ] Unit tests: PASS
- [ ] Integration tests: PASS
- [ ] E2E tests: PASS
- [ ] Performance tests: PASS
- [ ] Manual QA: COMPLETE

### Documentation
- [ ] CHANGELOG updated
- [ ] API docs updated
- [ ] README updated (if needed)

### Configuration
- [ ] Environment variables verified
- [ ] Feature flags configured
- [ ] Rollback plan ready

### Communication
- [ ] Team notified
- [ ] Downtime window communicated
- [ ] On-call team prepared
```

#### Deployment Commands

```bash
# Build application
npm run build

# Run tests
npm run test
npm run test:e2e

# Build Docker image
docker build -t biometrics-app:$VERSION .
docker tag biometrics-app:$VERSION biometrics-app:latest

# Deploy to staging
kubectl apply -f k8s/staging/
kubectl set image deployment/app app=biometrics-app:$VERSION -n staging

# Verify staging
curl -f https://staging.biometrics.example.com/health

# Deploy to production (requires approval)
kubectl apply -f k8s/production/
kubectl set image deployment/app app=biometrics-app:$VERSION -n production

# Verify production
curl -f https://biometrics.example.com/health
```

#### Environment Configuration

```typescript
// config/environments.ts
interface EnvironmentConfig {
  name: 'development' | 'staging' | 'production';
  apiUrl: string;
  databaseUrl: string;
  sentryDsn: string;
  features: {
    enableAnalytics: boolean;
    enableBetaFeatures: boolean;
    maintenanceMode: boolean;
  };
}

export const environments: Record<string, EnvironmentConfig> = {
  development: {
    name: 'development',
    apiUrl: 'http://localhost:3000/api',
    databaseUrl: process.env.DEV_DATABASE_URL,
    sentryDsn: '',
    features: {
      enableAnalytics: false,
      enableBetaFeatures: true,
      maintenanceMode: false,
    },
  },
  
  staging: {
    name: 'staging',
    apiUrl: 'https://staging.biometrics.example.com/api',
    databaseUrl: process.env.STAGING_DATABASE_URL,
    sentryDsn: process.env.SENTRY_DSN_STAGING,
    features: {
      enableAnalytics: true,
      enableBetaFeatures: true,
      maintenanceMode: false,
    },
  },
  
  production: {
    name: 'production',
    apiUrl: 'https://biometrics.example.com/api',
    databaseUrl: process.env.PROD_DATABASE_URL,
    sentryDsn: process.env.SENTRY_DSN_PRODUCTION,
    features: {
      enableAnalytics: true,
      enableBetaFeatures: false,
      maintenanceMode: false,
    },
  },
};
```

---

### 14.2 Rollback Procedures

#### Rollback Decision Tree

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      ROLLBACK DECISION TREE                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│                        Issue Detected                                   │
│                             │                                           │
│                             ▼                                           │
│                   ┌─────────────────┐                                   │
│                   │ Severity Level │                                   │
│                   └────────┬────────┘                                   │
│                            │                                             │
│          ┌─────────────────┼─────────────────┐                          │
│          ▼                 ▼                 ▼                          │
│   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐                 │
│   │  Critical  │   │    High     │   │   Medium    │                 │
│   │  (P0)      │   │    (P1)     │   │    (P2)     │                 │
│   └─────┬───────┘   └─────┬───────┘   └─────┬───────┘                 │
│         │                  │                  │                         │
│         ▼                  ▼                  ▼                         │
│   ┌──────────┐       ┌──────────┐       ┌──────────┐                  │
│   │ ROLLBACK │       │ Investig.│       │ Monitor  │                  │
│   │ IMMEDIATE│       │ + Fix    │       │ + Fix    │                  │
│   └──────────┘       └──────────┘       └──────────┘                  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Rollback Commands

```bash
# Kubernetes rollback
kubectl rollout undo deployment/app
kubectl rollout status deployment/app

# Or rollback to specific revision
kubectl rollout undo deployment/app --to-revision=3

# Docker rollback
docker pull biometrics/app:previous-version
kubectl set image deployment/app app=biometrics/app:previous-version

# Database rollback (if needed)
psql -h database.biometrics.example.com -U biometrics -c \
  "SELECT pg_restore --verbose --clean --if-exists backup.dump"

# Verify rollback
curl -f https://biometrics.example.com/health
```

#### Rollback Checklist

```markdown
## Rollback Checklist

### Immediate Actions
- [ ] Confirm rollback decision
- [ ] Notify on-call team
- [ ] Document issue in incident tracker

### Execution
- [ ] Rollback application code
- [ ] Rollback database (if needed)
- [ ] Verify rollback success

### Post-Rollback
- [ ] Run health checks
- [ ] Verify critical features
- [ ] Monitor error rates
- [ ] Update incident status

### Post-Incident
- [ ] Root cause analysis
- [ ] Document lessons learned
- [ ] Plan fix for re-deployment
```

---

### 14.3 Blue-Green Deployment

#### Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                   BLUE-GREEN DEPLOYMENT                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│                         Load Balancer                                   │
│                              │                                          │
│              ┌───────────────┴───────────────┐                         │
│              ▼                               ▼                          │
│      ┌─────────────┐                 ┌─────────────┐                   │
│      │   BLUE      │                 │   GREEN     │                   │
│      │  (Active)   │                 │  (Standby)  │                   │
│      │             │                 │             │                   │
│      │  v1.0.0     │                 │  v1.1.0     │                   │
│      │             │   ◀──────▶     │             │                   │
│      │  Production │   Switch       │  Testing    │                   │
│      │             │                 │             │                   │
│      └─────────────┘                 └─────────────┘                   │
│                                                                         │
│   Traffic: 100% Blue                    Traffic: 0%                   │
│                                                                         │
│   After Testing:                                                   │
│                                                                         │
│   Traffic: 0% Blue                       Traffic: 100% Green          │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Blue-Green Workflow

```bash
# 1. Current production is BLUE (v1.0.0)
# Deploy GREEN (v1.1.0) for testing

# 2. Deploy to green environment
kubectl apply -f k8s/green/
kubectl set image deployment/app-green app=biometrics-app:v1.1.0

# 3. Run smoke tests on green
./scripts/smoke-tests.sh --env=green

# 4. If tests pass, switch traffic
kubectl patch service app-lb -p '{"spec":{"selector":{"version":"green"}}}'

# 5. Verify traffic switch
curl -I https://biometrics.example.com
# Should return: X-Green-Version: v1.1.0

# 6. If issues, rollback to blue
kubectl patch service app-lb -p '{"spec":{"selector":{"version":"blue"}}}'

# 7. Promote green to production (old blue becomes standby)
kubectl label deployment/app-blue version=blue --overwrite
kubectl label deployment/app-green version=blue --overwrite
```

#### Traffic Switch Script

```typescript
// scripts/traffic-switch.ts
import { k8sClient } from './k8s-client';

interface TrafficSwitchOptions {
  serviceName: string;
  fromVersion: string;
  toVersion: string;
  percentage?: number;
}

async function switchTraffic(options: TrafficSwitchOptions) {
  const { serviceName, fromVersion, toVersion, percentage = 100 } = options;
  
  console.log(`🔄 Switching traffic from ${fromVersion} to ${toVersion}`);
  
  if (percentage < 100) {
    // Partial rollout - use canary
    await applyCanaryTrafficSplit(serviceName, toVersion, percentage);
    console.log(`📊 ${percentage}% traffic to ${toVersion}`);
  } else {
    // Full switch
    await k8sClient.patchNamespacedService(
      serviceName,
      'default',
      {
        spec: {
          selector: { version: toVersion }
        }
      }
    );
    console.log(`✅ 100% traffic now going to ${toVersion}`);
  }
  
  // Verify switch
  await verifyTrafficSwitch(serviceName, toVersion);
}

async function verifyTrafficSwitch(serviceName: string, version: string) {
  const checks = 10;
  let success = true;
  
  for (let i = 0; i < checks; i++) {
    const response = await fetch(`https://biometrics.example.com/health`);
    const versionHeader = response.headers.get('X-Version');
    
    if (versionHeader !== version) {
      success = false;
      console.error(`❌ Health check ${i + 1} failed: expected ${version}, got ${versionHeader}`);
    }
    
    await sleep(1000);
  }
  
  if (success) {
    console.log(`✅ Traffic switch verified`);
  } else {
    throw new Error('Traffic switch verification failed');
  }
}
```

---

### 14.4 Canary Releases

#### Canary Strategy

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      CANARY RELEASE STRATEGY                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   Phase 1: 5% Traffic                                                  │
│   ┌─────────────────────────────────────────────────────────────┐       │
│   │  ██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  │       │
│   │  5% Canary ───────────────────────────────────────────────  │       │
│   │  95% Stable ───────────────────────────────────────────────  │       │
│   └─────────────────────────────────────────────────────────────┘       │
│                                                                         │
│   Phase 2: 20% Traffic (if Phase 1 OK)                                │
│   ┌─────────────────────────────────────────────────────────────┐       │
│   │  ████████████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  │       │
│   │  20% Canary ──────────────────────────────────────────────  │       │
│   │  80% Stable ─────────────────────────────────────────────  │       │
│   └─────────────────────────────────────────────────────────────┘       │
│                                                                         │
│   Phase 3: 50% Traffic (if Phase 2 OK)                                 │
│   ┌─────────────────────────────────────────────────────────────┐       │
│   │  ███████████████████████████████████████████████░░░░░░░░░░░ │       │
│   │  50% Canary ─────────────────────────────────────────────  │       │
│   │  50% Stable ────────────────────────────────────────────  │       │
│   └─────────────────────────────────────────────────────────────┘       │
│                                                                         │
│   Phase 4: 100% Traffic (if Phase 3 OK)                               │
│   ┌─────────────────────────────────────────────────────────────┐       │
│   │  ████████████████████████████████████████████████████████  │       │
│   │  100% Canary (New Version) ─────────────────────────────  │       │
│   └─────────────────────────────────────────────────────────────┘       │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Canary Configuration

```yaml
# k8s/canary.yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: biometrics-app
spec:
  replicas: 10
  strategy:
    canary:
      canaryService: app-canary
      stableService: app-stable
      trafficRouting:
        nginx:
          stableIngress: app-ingress
          additionalIngressAnnotations:
            canary-by-header: X-Canary
      steps:
        - setWeight: 5
        - pause: {duration: 10m}
        - analysis:
            templates:
              - templateName: success-rate
        - setWeight: 20
        - pause: {duration: 15m}
        - analysis:
            templates:
              - templateName: success-rate
        - setWeight: 50
        - pause: {duration: 30m}
        - analysis:
            templates:
              - templateName: success-rate
        - setWeight: 100
      analysis:
        templates:
          - templateName: success-rate
        startingStep: 1
        args:
          - name: service-name
            value: app-canary
---
apiVersion: argoproj.io/v1alpha1
kind: AnalysisTemplate
metadata:
  name: success-rate
spec:
  args:
    - name: service-name
  metrics:
    - name: success-rate
      interval: 1m
      successCondition: result[0] >= 0.99
      failureLimit: 3
      provider:
        prometheus:
          address: http://prometheus:9090
          query: |
            sum(rate(http_requests_total{service="{{args.service-name}}",status=~"2.."}[5m])) 
            / 
            sum(rate(http_requests_total{service="{{args.service-name}}"}[5m]))
```

#### Canary Monitoring

```typescript
// lib/canary-analysis.ts
interface CanaryMetrics {
  canaryVersion: string;
  stableVersion: string;
  metrics: {
    errorRate: { canary: number; stable: number; threshold: number };
    latencyP50: { canary: number; stable: number; threshold: number };
    latencyP99: { canary: number; stable: number; threshold: number };
    successRate: { canary: number; stable: number; threshold: number };
  };
  decision: 'promote' | 'rollback' | 'continue';
}

async function analyzeCanary(canaryWeight: number): Promise<CanaryMetrics> {
  const metrics = await Promise.all([
    getErrorRate('canary'),
    getErrorRate('stable'),
    getLatency('canary', 'p50'),
    getLatency('stable', 'p50'),
    getLatency('canary', 'p99'),
    getLatency('stable', 'p99'),
  ]);
  
  const [errorRate, latencyP50, latencyP99] = metrics;
  
  const analysis: CanaryMetrics = {
    canaryVersion: 'v1.1.0',
    stableVersion: 'v1.0.0',
    metrics: {
      errorRate,
      latencyP50,
      latencyP99,
      successRate: {
        canary: 1 - errorRate.canary,
        stable: 1 - errorRate.stable,
        threshold: 0.95,
      },
    },
    decision: 'continue',
  };
  
  // Decision logic
  if (errorRate.canary > errorRate.stable * 2) {
    analysis.decision = 'rollback';
  } else if (latencyP99.canary > latencyP99.stable * 1.5) {
    analysis.decision = 'rollback';
  } else if (canaryWeight >= 100) {
    analysis.decision = 'promote';
  }
  
  return analysis;
}
```

---

## 🔴 15. Incident Management

### 15.1 On-Call Rotation

#### Rotation Schedule

```typescript
// config/oncall.ts
interface OnCallSchedule {
  team: string;
  schedule: {
    [week: string]: {
      primary: string;
      secondary: string;
      startDate: string;
      endDate: string;
    };
  };
}

const onCallSchedule: OnCallSchedule = {
  team: 'BIOMETRICS-SRE',
  schedule: {
    '2026-W07': {
      primary: 'agent-alpha',
      secondary: 'agent-beta',
      startDate: '2026-02-17',
      endDate: '2026-02-24',
    },
    '2026-W08': {
      primary: 'agent-beta',
      secondary: 'agent-gamma',
      startDate: '2026-02-24',
      endDate: '2026-03-03',
    },
    '2026-W09': {
      primary: 'agent-gamma',
      secondary: 'agent-alpha',
      startDate: '2026-03-03',
      endDate: '2026-03-10',
    },
  },
};

function getCurrentOnCall(): { primary: string; secondary: string } {
  const currentWeek = getWeekNumber(new Date());
  return onCallSchedule.schedule[`2026-W${currentWeek}`];
}
```

#### On-Call Responsibilities

```markdown
## On-Call Responsibilities

### Primary On-Call
- First responder to all alerts
- Acknowledge alerts within 15 minutes
- Initial triage and assessment
- Escalate if unable to resolve within 30 minutes

### Secondary On-Call
- Backup for primary on-call
- Assist if primary is overwhelmed
- Handle non-critical issues
- Take over if primary has an emergency

### Escalation Levels
1. **Primary On-Call** → Initial response
2. **Secondary On-Call** → If primary unavailable
3. **Team Lead** → If unresolved after 30 min
4. **Engineering Manager** → If unresolved after 1 hour
5. **CTO** → Critical production outage

### Response Time SLA
| Severity | First Response | Resolution Target |
|----------|-----------------|-------------------|
| P0/Critical | 15 min | 1 hour |
| P1/High | 30 min | 4 hours |
| P2/Medium | 2 hours | 24 hours |
| P3/Low | 24 hours | 7 days |
```

---

### 15.2 Severity Classification

#### Severity Levels

```typescript
enum Severity {
  P0_CRITICAL = 'P0',
  P1_HIGH = 'P1',
  P2_MEDIUM = 'P2',
  P3_LOW = 'P3',
}

interface IncidentSeverity {
  level: Severity;
  definition: string;
  examples: string[];
  responseTime: number;
  resolutionTarget: number;
}

const severityDefinitions: Record<Severity, IncidentSeverity> = {
  [Severity.P0_CRITICAL]: {
    level: Severity.P0_CRITICAL,
    definition: 'Complete service outage affecting all users',
    examples: [
      'Database completely unavailable',
      'All API endpoints returning 5xx errors',
      'Complete data loss',
      'Security breach with data exposure',
    ],
    responseTime: 15, // minutes
    resolutionTarget: 60, // minutes
  },
  
  [Severity.P1_HIGH]: {
    level: Severity.P1_HIGH,
    definition: 'Major functionality impaired affecting majority of users',
    examples: [
      'Payment processing not working',
      'Login system completely down',
      'API rate limiting affecting all users',
      'High error rate (>10%) on main endpoints',
    ],
    responseTime: 30, // minutes
    resolutionTarget: 240, // minutes
  },
  
  [Severity.P2_MEDIUM]: {
    level: Severity.P2_MEDIUM,
    definition: 'Partial functionality affected, workaround available',
    examples: [
      'Search feature slow but working',
      'Email notifications delayed',
      'Non-critical feature broken',
      'Performance degradation <50%',
    ],
    responseTime: 120, // minutes
    resolutionTarget: 1440, // minutes (24 hours)
  },
  
  [Severity.P3_LOW]: {
    level: Severity.P3_LOW,
    definition: 'Minor issue with minimal user impact',
    examples: [
      'Typo in UI',
      'Minor cosmetic issue',
      'Documentation error',
      'Feature request',
    ],
    responseTime: 1440, // minutes (24 hours)
    resolutionTarget: 10080, // minutes (7 days)
  },
};
```

#### Severity Decision Matrix

```
┌─────────────────────────────────────────────────────────────────────────┐
│                   SEVERITY DECISION MATRIX                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   Impact │   All Users   │   Some Users  │  Few Users    │            │
│   ───────┼───────────────┼───────────────┼───────────────┤            │
│   Time    │               │               │               │            │
│   ───────┼───────────────┼───────────────┼───────────────┤            │
│   <1 hr   │     P0        │      P1       │      P2       │            │
│   1-4 hr  │     P1        │      P1       │      P2       │            │
│   4-24 hr │     P1        │      P2       │      P3       │            │
│   >24 hr  │     P2        │      P2       │      P3       │            │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

### 15.3 Escalation Paths

#### Escalation Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        ESCALATION FLOW                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   Alert Triggered                                                        │
│        │                                                                 │
│        ▼                                                                 │
│   ┌─────────────┐     ┌─────────────┐                                  │
│   │  Primary    │ Yes │  Resolved?  │                                  │
│   │  On-Call    │────▶│             │                                  │
│   └─────────────┘     └──────┬──────┘                                  │
│        │ No                  │ Yes                                      │
│        ▼                     │                                          │
│   ┌─────────────┐            ▼                                          │
│   │  Ack?       │     ┌─────────────┐                                   │
│   │   (15 min)  │     │   Close     │                                   │
│   └──────┬──────┘     │   Incident  │                                   │
│        │ No            └─────────────┘                                   │
│        ▼                                                               │
│   ┌─────────────┐     ┌─────────────┐                                  │
│   │ Secondary   │ Yes │  Resolved?  │                                  │
│   │  On-Call    │────▶│             │                                  │
│   └─────────────┘     └──────┬──────┘                                  │
│        │ No                  │ Yes                                      │
│        ▼                     ▼                                          │
│   ┌─────────────┐     ┌─────────────┐                                   │
│   │  Team Lead  │ Yes │  Resolved?  │                                  │
│   │  (30 min)   │────▶│             │                                  │
│   └─────────────┘     └──────┬──────┘                                  │
│        │ No                  │ Yes                                      │
│        ▼                     ▼                                          │
│   ┌─────────────┐     ┌─────────────┐                                   │
│   │   Eng       │ Yes │  Resolved?  │                                  │
│   │   Manager   │────▶│             │                                  │
│   └─────────────┘     └──────┬──────┘                                  │
│        │ No                  │ Yes                                      │
│        ▼                     ▼                                          │
│   ┌─────────────┐     ┌─────────────┐                                   │
│   │     CTO     │────▶│  Resolved   │                                  │
│   │  (Critical) │     │             │                                  │
│   └─────────────┘     └─────────────┘                                   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Escalation Contacts

```typescript
// config/escalation.ts
interface EscalationContact {
  role: string;
  name: string;
  phone: string;
  email: string;
  availability: '24x7' | 'business_hours';
}

const escalationContacts: EscalationContact[] = [
  {
    role: 'Primary On-Call',
    name: 'System',
    phone: '+49XXXOnCall1',
    email: 'oncall-primary@biometrics.example.com',
    availability: '24x7',
  },
  {
    role: 'Secondary On-Call',
    name: 'System',
    phone: '+49XXXOnCall2',
    email: 'oncall-secondary@biometrics.example.com',
    availability: '24x7',
  },
  {
    role: 'Team Lead - Platform',
    name: 'TBD',
    phone: '+49XXXTeamLead',
    email: 'team-lead-platform@biometrics.example.com',
    availability: '24x7',
  },
  {
    role: 'Engineering Manager',
    name: 'TBD',
    phone: '+49XXXEngManager',
    email: 'eng-manager@biometrics.example.com',
    availability: 'business_hours',
  },
  {
    role: 'CTO',
    name: 'TBD',
    phone: '+49XXXCTO',
    email: 'cto@biometrics.example.com',
    availability: '24x7',
  },
];
```

---

### 15.4 Post-Mortem Process

#### Post-Mortem Template

```markdown
# Incident Post-Mortem

**Incident ID:** INC-2026-001
**Date:** 2026-02-15
**Duration:** 2 hours 34 minutes
**Severity:** P1
**Author:** [Name]

---

## Summary

Brief description of what happened and impact.

---

## Timeline (UTC)

| Time | Event |
|------|-------|
| 10:00 | Alert triggered: High error rate on /api/auth |
| 10:15 | On-call acknowledged alert |
| 10:22 | Initial investigation started |
| 10:45 | Root cause identified: Database connection pool exhausted |
| 11:15 | Fix deployed |
| 11:30 | Service recovered |
| 12:34 | Incident closed |

---

## Impact

- **Users Affected:** ~5,000
- **Error Rate:** 23%
- **Downtime:** 1 hour 30 minutes
- **Revenue Impact:** TBD

---

## Root Cause

Detailed explanation of what caused the incident.

---

## Resolution

How the issue was resolved.

---

## Lessons Learned

### What Went Well
- Alert was triggered quickly
- Team responded promptly
- Documentation was clear

### What Could Be Improved
- Database monitoring could be more detailed
- Runbook for this scenario needed
- Auto-scaling should be configured

---

## Action Items

| Item | Owner | Due Date |
|------|-------|----------|
| Add database connection pool monitoring | @person1 | 2026-02-22 |
| Create runbook for database issues | @person2 | 2026-02-25 |
| Configure auto-scaling for API | @person3 | 2026-03-01 |

---

## Supporting Data

- [Link to Grafana dashboard]
- [Link to error logs]
- [Link to related incidents]
```

#### Post-Mortem Meeting

```typescript
interface PostMortemMeeting {
  scheduledWithin: '48 hours';
  attendees: [
    'On-call who responded',
    'Team lead',
    'Engineering manager',
    'Relevant developers',
  ];
  agenda: [
    'Timeline review (10 min)',
    'Impact discussion (10 min)',
    'Root cause analysis (20 min)',
    'Lessons learned (15 min)',
    'Action items (15 min)',
  ];
  output: [
    'Published post-mortem document',
    'Action items in tracking system',
    'Updated runbooks if needed',
  ];
}
```

---

## 💾 16. Data Workflows

### 16.1 ETL Pipelines

#### ETL Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         ETL PIPELINE                                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐          │
│   │ Source  │───▶│ Extract │───▶│ Transform│───▶│  Load   │          │
│   │         │    │         │    │         │    │         │          │
│   │ API     │    │ Fetch   │    │ Clean   │    │ Insert  │          │
│   │ Database│    │ Parse   │    │ Validate│    │ Update  │          │
│   │ Files   │    │ Normalize│    │ Enrich  │    │ Merge   │          │
│   └─────────┘    └─────────┘    └─────────┘    └────┬────┘          │
│                                                      │                │
│                                                      ▼                │
│                                              ┌───────────┐          │
│                                              │ Destination│          │
│                                              │           │          │
│                                              │ Data Lake │          │
│                                              │ Warehouse │          │
│                                              │ Analytics │          │
│                                              └───────────┘          │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### ETL Job Configuration

```yaml
# config/etl/jobs/user-sync.yaml
name: user-sync
schedule: "0 */6 * * *"  # Every 6 hours
source:
  type: postgres
  connection:
    host: "{{ env.DATABASE_HOST }}"
    port: 5432
    database: biometrics
    schema: public
  query: |
    SELECT 
      id,
      email,
      name,
      created_at,
      updated_at,
      status,
      metadata
    FROM users
    WHERE updated_at > :last_sync_timestamp
    OR created_at > :last_sync_timestamp

transformations:
  - name: normalize_email
    type: transform
    function: lowercase
    field: email
  
  - name: enrich_profile
    type: enrich
    data:
      source: http
      url: https://api.gravatar.com/{{ hashlib.md5(email).hexdigest() }}
      cache_ttl: 86400
      
  - name: validate_fields
    type: validate
    rules:
      - field: email
        type: email
        required: true
      - field: name
        type: string
        required: true
        min_length: 1
        max_length: 100

destination:
  type: bigquery
  connection:
    project: biometrics-prod
    dataset: analytics
    table: users
  mode: merge
  merge_keys:
    - user_id

error_handling:
  on_error: log
  retry_attempts: 3
  retry_delay: 300  # seconds
  dead_letter_queue: etl-errors
```

#### ETL Monitoring

```typescript
// lib/etl/monitor.ts
interface ETLJobMetrics {
  jobName: string;
  startTime: Date;
  endTime: Date;
  duration: number;
  status: 'success' | 'failed' | 'partial';
  recordsExtracted: number;
  recordsTransformed: number;
  recordsLoaded: number;
  errors: {
    count: number;
    messages: string[];
  };
}

async function monitorETLJob(jobName: string): Promise<ETLJobMetrics> {
  const startTime = Date.now();
  
  try {
    const extraction = await extractData(jobName);
    const transformation = await transformData(extraction);
    const loadResult = await loadData(transformation);
    
    const metrics: ETLJobMetrics = {
      jobName,
      startTime: new Date(startTime),
      endTime: new Date(),
      duration: Date.now() - startTime,
      status: loadResult.errors.length > 0 ? 'partial' : 'success',
      recordsExtracted: extraction.count,
      recordsTransformed: transformation.count,
      recordsLoaded: loadResult.count,
      errors: {
        count: loadResult.errors.length,
        messages: loadResult.errors,
      },
    };
    
    await sendMetrics(metrics);
    return metrics;
  } catch (error) {
    const metrics: ETLJobMetrics = {
      jobName,
      startTime: new Date(startTime),
      endTime: new Date(),
      duration: Date.now() - startTime,
      status: 'failed',
      recordsExtracted: 0,
      recordsTransformed: 0,
      recordsLoaded: 0,
      errors: {
        count: 1,
        messages: [error.message],
      },
    };
    
    await sendMetrics(metrics);
    await alertOnCall(metrics);
    throw error;
  }
}
```

---

### 16.2 Data Migration

#### Migration Strategy

```typescript
interface MigrationStrategy {
  name: string;
  description: string;
  downtime: boolean;
  risk: 'low' | 'medium' | 'high';
  steps: string[];
}

const migrationStrategies: Record<string, MigrationStrategy> = {
  copyTable: {
    name: 'Copy Table',
    description: 'Create new table, copy data, switch',
    downtime: true,
    risk: 'low',
    steps: [
      'Create new table with desired schema',
      'Copy data from old table',
      'Verify data integrity',
      'Update application to use new table',
      'Drop old table',
    ],
  },
  
  expandContract: {
    name: 'Expand-Contract',
    description: 'Add new column, migrate data, remove old',
    downtime: false,
    risk: 'medium',
    steps: [
      'Add new column to existing table',
      'Update writes to populate both columns',
      'Backfill new column with transformed data',
      'Verify data consistency',
      'Update reads to use new column',
      'Remove old column',
    ],
  },
  
  parallelRun: {
    name: 'Parallel Run',
    description: 'Run both systems simultaneously',
    downtime: false,
    risk: 'low',
    steps: [
      'Implement new system alongside old',
      'Write to both systems',
      'Read from old system',
      'Compare results in background',
      'Switch reads to new system when stable',
      'Decommission old system',
    ],
  },
  
  featureFlag: {
    name: 'Feature Flag Migration',
    description: 'Use feature flags for gradual rollout',
    downtime: false,
    risk: 'low',
    steps: [
      'Add feature flag to control migration',
      'Implement new logic alongside old',
      'Enable for small percentage',
      'Monitor for errors',
      'Increase percentage gradually',
      'Remove old logic',
      'Remove feature flag',
    ],
  },
};
```

#### Migration Runbook

```markdown
# Data Migration Runbook: User Table Schema Update

## Overview
Migration of users table from schema v1 to v2

## Pre-Migration Checklist
- [ ] Backup created
- [ ] Migration tested on staging
- [ ] Rollback plan tested
- [ ] Communication sent to team
- [ ] Monitoring alerts configured

## Migration Steps

### Step 1: Create Backup
```bash
pg_dump -h database.biometrics.example.com \
  -U biometrics \
  -t users > backup_users_$(date +%Y%m%d).sql
```

### Step 2: Add New Columns
```sql
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS display_name VARCHAR(100),
ADD COLUMN IF NOT EXISTS avatar_url TEXT,
ADD COLUMN IF NOT EXISTS preferences JSONB DEFAULT '{}';
```

### Step 3: Backfill Data
```sql
UPDATE users 
SET 
  display_name = COALESCE(name, email),
  avatar_url = CONCAT('https://api.dicebear.com/7.x/initials/svg?seed=', id),
  preferences = COALESCE(preferences, '{}'::jsonb)
WHERE display_name IS NULL;
```

### Step 4: Verify Data
```sql
SELECT 
  COUNT(*) as total,
  COUNT(display_name) as with_display_name,
  COUNT(avatar_url) as with_avatar
FROM users;
```

### Step 5: Switch Application
- Deploy new application version
- Monitor error rates

### Step 6: Cleanup (after 7 days)
```sql
ALTER TABLE users DROP COLUMN IF EXISTS name;
```

## Rollback
```bash
psql -h database.biometrics.example.com \
  -U biometrics < backup_users_20260215.sql
```

## Monitoring
- Error rate dashboard: [Link]
- Migration progress: [Link]
- User reports: #biometrics-support
```

---

### 16.3 Backup and Restore

#### Backup Strategy

```yaml
# config/backup.yaml
backups:
  database:
    - name: daily-full
      schedule: "0 2 * * *"  # 2 AM daily
      type: full
      retention: 30 days
      destination: s3://biometrics-backups/database/daily/
      
    - name: weekly-full
      schedule: "0 3 * * 0"  # 3 AM Sunday
      type: full
      retention: 90 days
      destination: s3://biometrics-backups/database/weekly/
      
    - name: incremental
      schedule: "0 */4 * * *"  # Every 4 hours
      type: incremental
      retention: 7 days
      destination: s3://biometrics-backups/database/incremental/
      
  files:
    - name: user-uploads
      schedule: "0 */6 * * *"  # Every 6 hours
      source: /data/uploads
      destination: s3://biometrics-backups/uploads/
      retention: 30 days
        
  config:
    - name: application-config
      schedule: "0 * * * *"  # Hourly
      source: /app/config
      destination: s3://biometrics-backups/config/
      retention: 90 days
```

#### Restore Procedures

```bash
#!/bin/bash
# scripts/restore-database.sh

set -e

BACKUP_FILE=$1
TARGET_DB=$2

if [ -z "$BACKUP_FILE" ] || [ -z "$TARGET_DB" ]; then
  echo "Usage: $0 <backup-file> <target-database>"
  exit 1
fi

echo "🔄 Starting restore of $BACKUP_FILE to $TARGET_DB"

# 1. Stop application
echo "Stopping application..."
kubectl scale deployment app --replicas=0 -n production

# 2. Drop existing database
echo "Dropping existing database..."
psql -h $DB_HOST -U $DB_USER -c "DROP DATABASE IF EXISTS $TARGET_DB;"

# 3. Create new database
echo "Creating new database..."
psql -h $DB_HOST -U $DB_USER -c "CREATE DATABASE $TARGET_DB;"

# 4. Restore from backup
echo "Restoring from backup..."
pg_restore -h $DB_HOST -U $DB_USER -d $TARGET_DB \
  --verbose \
  --no-owner \
  --no-privileges \
  "$BACKUP_FILE"

# 5. Verify restore
echo "Verifying restore..."
psql -h $DB_HOST -U $DB_USER -d $TARGET_DB -c "SELECT COUNT(*) FROM users;"

# 6. Start application
echo "Starting application..."
kubectl scale deployment app --replicas=3 -n production

# 7. Verify application
echo "Verifying application health..."
curl -f https://biometrics.example.com/health

echo "✅ Restore complete!"
```

#### Backup Verification

```typescript
// lib/backup/verify.ts
interface BackupVerification {
  backupId: string;
  timestamp: Date;
  status: 'passed' | 'failed';
  checks: {
    name: string;
    passed: boolean;
    details: string;
  }[];
}

async function verifyBackup(backupId: string): Promise<BackupVerification> {
  const checks = [];
  
  // Check 1: File exists
  const exists = await checkFileExists(backupId);
  checks.push({
    name: 'File exists',
    passed: exists,
    details: exists ? `File size: ${await getFileSize(backupId)}` : 'File not found',
  });
  
  // Check 2: Checksum valid
  const checksum = await verifyChecksum(backupId);
  checks.push({
    name: 'Checksum valid',
    passed: checksum.valid,
    details: checksum.details,
  });
  
  // Check 3: Extractable
  const extractable = await verifyExtractable(backupId);
  checks.push({
    name: 'Backup extractable',
    passed: extractable,
    details: extractable ? 'Backup can be extracted' : 'Extraction failed',
  });
  
  // Check 4: Contains expected tables
  const tables = await verifyTables(backupId);
  checks.push({
    name: 'Expected tables present',
    passed: tables.valid,
    details: tables.details,
  });
  
  return {
    backupId,
    timestamp: new Date(),
    status: checks.every(c => c.passed) ? 'passed' : 'failed',
    checks,
  };
}
```

---

### 16.4 Data Quality Checks

#### DQ Check Framework

```typescript
// lib/data-quality/checks.ts
interface DQCheck {
  id: string;
  name: string;
  description: string;
  severity: 'critical' | 'warning' | 'info';
  query: string;
  threshold: number;
}

const dataQualityChecks: DQCheck[] = [
  {
    id: 'dq-001',
    name: 'Null Email Check',
    description: 'Check for users with null email',
    severity: 'critical',
    query: `SELECT COUNT(*) FROM users WHERE email IS NULL`,
    threshold: 0,
  },
  {
    id: 'dq-002',
    name: 'Duplicate Email Check',
    description: 'Check for duplicate email addresses',
    severity: 'critical',
    query: `
      SELECT COUNT(*) FROM (
        SELECT email, COUNT(*) as cnt
        FROM users
        WHERE email IS NOT NULL
        GROUP BY email
        HAVING COUNT(*) > 1
      ) duplicates
    `,
    threshold: 0,
  },
  {
    id: 'dq-003',
    name: 'Stale Data Check',
    description: 'Check for users not updated in over 1 year',
    severity: 'warning',
    query: `
      SELECT COUNT(*) FROM users 
      WHERE updated_at < NOW() - INTERVAL '1 year'
    `,
    threshold: 1000,
  },
  {
    id: 'dq-004',
    name: 'Data Freshness',
    description: 'Check last data load timestamp',
    severity: 'info',
    query: `SELECT MAX(updated_at) as last_update FROM users`,
    threshold: 24, // hours
  },
];

async function runDQCheck(check: DQCheck): Promise<DQCheckResult> {
  const result = await db.query(check.query);
  const value = result.rows[0]?.count || result.rows[0]?.last_update;
  
  return {
    checkId: check.id,
    passed: value <= check.threshold,
    value,
    threshold: check.threshold,
    timestamp: new Date(),
  };
}
```

#### DQ Dashboard

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    DATA QUALITY DASHBOARD                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   Overall Score: ████████████████░░░  94%                              │
│                                                                         │
│   ┌───────────────────────────────────────────────────────────────┐    │
│   │ Critical Issues (2)                                            │    │
│   ├───────────────────────────────────────────────────────────────┤    │
│   │ ❌ dq-001: Null Email Check    - 5 found (threshold: 0)       │    │
│   │ ❌ dq-002: Duplicate Email     - 12 found (threshold: 0)      │    │
│   └───────────────────────────────────────────────────────────────┘    │
│                                                                         │
│   ┌───────────────────────────────────────────────────────────────┐    │
│   │ Warnings (1)                                                  │    │
│   ├───────────────────────────────────────────────────────────────┤    │
│   │ ⚠️ dq-003: Stale Data           - 850 found (threshold: 1000)│    │
│   └───────────────────────────────────────────────────────────────┘    │
│                                                                         │
│   ┌───────────────────────────────────────────────────────────────┐    │
│   │ Info (1)                                                      │    │
│   ├───────────────────────────────────────────────────────────────┤    │
│   │ ℹ️ dq-004: Data Freshness       - 2 hours ago                │    │
│   └───────────────────────────────────────────────────────────────┘    │
│                                                                         │
│   Last Check: 2026-02-15 14:00 UTC                                    │
│   Next Check: 2026-02-15 15:00 UTC                                    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 🔒 17. Security Workflows

### 17.1 Vulnerability Response

#### Vulnerability Handling Process

```
┌─────────────────────────────────────────────────────────────────────────┐
│                  VULNERABILITY RESPONSE PROCESS                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   Vulnerability Discovered                                              │
│             │                                                          │
│             ▼                                                          │
│   ┌─────────────────┐                                                  │
│   │    Triage       │                                                  │
│   │  (24 hours)     │                                                  │
│   └────────┬────────┘                                                  │
│            │                                                           │
│            ▼                                                           │
│   ┌─────────────────────────────────────┐                              │
│   │        Severity Assessment          │                              │
│   │                                     │                              │
│   │   Critical (9-10) ───► Immediate    │                              │
│   │   High (7-8)      ───► 7 days        │                              │
│   │   Medium (4-6)    ───► 30 days       │                              │
│   │   Low (1-3)       ───► 90 days       │                              │
│   └─────────────────────────────────────┘                              │
│            │                                                           │
│            ▼                                                           │
│   ┌─────────────────────────────────────┐                              │
│   │           Fix Development           │                              │
│   │                                     │                              │
│   │   - Create fix                      │                              │
│   │   - Write tests                     │                              │
│   │   - Code review                     │                              │
│   │   - Security review                 │                              │
│   └─────────────────────────────────────┘                              │
│            │                                                           │
│            ▼                                                           │
│   ┌─────────────────────────────────────┐                              │
│   │           Deployment                 │                              │
│   │                                     │                              │
│   │   - Deploy to staging               │                              │
│   │   - Run security tests              │                              │
│   │   - Deploy to production            │                              │
│   └─────────────────────────────────────┘                              │
│            │                                                           │
│            ▼                                                           │
│   ┌─────────────────────────────────────┐                              │
│   │           Post-Mortem                │                              │
│   │                                     │                              │
│   │   - Document vulnerability          │                              │
│   │   - Review fix                      │                              │
│   │   - Update dependencies             │                              │
│   └─────────────────────────────────────┘                              │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Vulnerability Report Template

```markdown
# Vulnerability Report

**ID:** VULN-2026-001
**Date Reported:** 2026-02-15
**Reported By:** [Name]
**Status:** [Open/In Progress/Resolved]

---

## Summary

Brief description of the vulnerability.

---

## Technical Details

**Component:** [e.g., authentication-module]
**CVSS Score:** [e.g., 8.5]
**CVE ID:** [if applicable]

**Vulnerability Type:**
- [ ] SQL Injection
- [ ] XSS
- [ ] CSRF
- [ ] Authentication Bypass
- [ ] Authorization Bypass
- [ ] Information Disclosure
- [ ] Other: __________

---

## Reproduction Steps

1. Step 1
2. Step 2
3. Step 3

---

## Impact

What can an attacker do?

---

## Recommended Fix

Suggested remediation approach.

---

## Timeline

| Date | Action |
|------|--------|
| 2026-02-15 | Vulnerability reported |
| 2026-02-16 | Triaged and assigned |
| [Date] | Fix deployed |
| [Date] | Verified fixed |
```

---

### 17.2 Patch Management

#### Patch Classification

```typescript
interface Patch {
  id: string;
  description: string;
  severity: 'critical' | 'high' | 'medium' | 'low';
  affectedSystems: string[];
  requiresDowntime: boolean;
  releaseDate: Date;
  dependencies: string[];
}

const patchCategories = {
  security: {
    critical: {
      responseTime: '24 hours',
      deploymentWindow: 'Immediate',
      testingRequired: 'Full regression',
      approvalRequired: 'Security team',
    },
    high: {
      responseTime: '7 days',
      deploymentWindow: 'Next release',
      testingRequired: 'Smoke tests',
      approvalRequired: 'Team lead',
    },
  },
  
  bugfix: {
    critical: {
      responseTime: '48 hours',
      deploymentWindow: 'Hotfix',
      testingRequired: 'Targeted tests',
      approvalRequired: 'Team lead',
    },
    normal: {
      responseTime: 'Next sprint',
      deploymentWindow: 'Release window',
      testingRequired: 'Standard',
      approvalRequired: 'None',
    },
  },
};
```

#### Patch Deployment Process

```bash
#!/bin/bash
# scripts/deploy-patch.sh

PATCH_ID=$1
VERSION=$2

if [ -z "$PATCH_ID" ] || [ -z "$VERSION" ]; then
  echo "Usage: $0 <patch-id> <version>"
  exit 1
fi

echo "Deploying patch $PATCH_ID (v$VERSION)"

# 1. Verify patch package
echo "Verifying patch package..."
./scripts/verify-patch.sh "$PATCH_ID"

# 2. Run pre-deployment checks
echo "Running pre-deployment checks..."
./scripts/pre-deploy-checks.sh

# 3. Create database backup
echo "Creating backup..."
./scripts/backup-database.sh

# 4. Deploy patch
echo "Deploying..."
kubectl set image deployment/app app=biometrics-app:$VERSION

# 5. Verify deployment
echo "Verifying deployment..."
./scripts/verify-deployment.sh

# 6. Run smoke tests
echo "Running smoke tests..."
./scripts/smoke-tests.sh

# 7. Monitor for 30 minutes
echo "Monitoring for 30 minutes..."
./scripts/monitor-deployment.sh 30

echo "✅ Patch deployed successfully!"
```

---

### 17.3 Key Rotation

#### Rotation Schedule

```typescript
interface KeyRotationPolicy {
  keyType: string;
  rotationPeriod: number; // days
  gracePeriod: number; // days
  alertBeforeExpiry: number; // days
}

const keyRotationPolicies: KeyRotationPolicy[] = [
  {
    keyType: 'API Keys',
    rotationPeriod: 90,
    gracePeriod: 14,
    alertBeforeExpiry: 7,
  },
  {
    keyType: 'Database Passwords',
    rotationPeriod: 180,
    gracePeriod: 30,
    alertBeforeExpiry: 14,
  },
  {
    keyType: 'Encryption Keys',
    rotationPeriod: 365,
    gracePeriod: 30,
    alertBeforeExpiry: 30,
  },
  {
    keyType: 'OAuth Client Secrets',
    rotationPeriod: 180,
    gracePeriod: 14,
    alertBeforeExpiry: 14,
  },
  {
    keyType: 'JWT Signing Keys',
    rotationPeriod: 90,
    gracePeriod: 7,
    alertBeforeExpiry: 7,
  },
];
```

#### Key Rotation Workflow

```bash
#!/bin/bash
# scripts/rotate-api-key.sh

KEY_NAME=$1

if [ -z "$KEY_NAME" ]; then
  echo "Usage: $0 <key-name>"
  exit 1
fi

echo "Rotating API key: $KEY_NAME"

# 1. Generate new key
echo "Generating new key..."
NEW_KEY=$(./scripts/generate-api-key.sh)
echo "New key generated: ${NEW_KEY:0:8}..."

# 2. Update secrets manager
echo "Updating secrets manager..."
aws secretsmanager put-secret-value \
  --secret-id "$KEY_NAME" \
  --secret-string "{\"key\": \"$NEW_KEY\"}" \
  --version-staging-version-label AWSPENDING

# 3. Update application config (without restart)
kubectl set env deployment/app -n production \
  API_KEY="$NEW_KEY"

# 4. Verify new key works
echo "Verifying new key..."
./scripts/test-api-key.sh "$NEW_KEY"

# 5. Finalize rotation
echo "Finalizing rotation..."
aws secretsmanager.update-secret-version-stage \
  --secret-id "$KEY_NAME" \
  --version-stage AWSPENDING \
  --move-to-version-id $(aws secretsmanager list-secret-version-ids --secret-id "$KEY_NAME" | jq -r '.Versions[-1].VersionId')

# 6. Revoke old key (after grace period)
echo "Old key will be revoked in 14 days"
echo "$KEY_NAME" >> keys_to_revoke.txt

echo "✅ Key rotation complete!"
```

---

### 17.4 Access Reviews

#### Access Review Process

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      ACCESS REVIEW PROCESS                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   Quarterly Access Review                                               │
│             │                                                          │
│             ▼                                                          │
│   ┌─────────────────┐                                                  │
│   │ Generate Report │  (Who has access to what)                       │
│   └────────┬────────┘                                                  │
│            │                                                           │
│            ▼                                                           │
│   ┌─────────────────┐                                                  │
│   │  Send to       │  (Team leads review their teams)                  │
│   │  Team Leads    │                                                  │
│   └────────┬────────┘                                                  │
│            │                                                           │
│            ▼                                                           │
│   ┌─────────────────┐                                                  │
│   │  Review Access │  (14 days to complete)                           │
│   └────────┬────────┘                                                  │
│            │                                                           │
│      ┌─────┴─────┐                                                    │
│      │            │                                                    │
│      ▼            ▼                                                    │
│   Approve     Request Removal                                          │
│      │            │                                                    │
│      │            ▼                                                    │
│      │     ┌─────────────┐                                             │
│      │     │  Remove     │                                             │
│      │     │  Access     │                                             │
│      │     └─────────────┘                                             │
│      │            │                                                    │
│      └──────┬─────┘                                                    │
│             ▼                                                           │
│   ┌─────────────────┐                                                  │
│   │  Sign-off       │                                                  │
│   │  & Document     │                                                  │
│   └─────────────────┘                                                  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Access Review Template

```markdown
# Quarterly Access Review - Q1 2026

**Review Period:** January 1 - March 31, 2026
**Due Date:** April 15, 2026
**Reviewer:** [Team Lead Name]

---

## Systems Reviewed

| System | Total Users | Active | Inactive | Removed |
|--------|-------------|--------|----------|---------|
| Production API | 25 | 20 | 3 | 2 |
| Admin Dashboard | 15 | 12 | 2 | 1 |
| Database | 8 | 6 | 1 | 1 |
| CI/CD | 12 | 10 | 2 | 0 |

---

## Access Changes Needed

| User | System | Action | Reason |
|------|--------|--------|--------|
| john@company.com | Production API | Remove | Left company |
| jane@company.com | Admin | Remove | Role change |
| new-hire@company.com | CI/CD | Add | New team member |

---

## Certification

I certify that:
- [ ] All access listed has been reviewed
- [ ] Access is appropriate for job function
- [ ] No unauthorized access exists
- [ ] Access removals have been completed

**Reviewer Signature:** _________________
**Date:** _________________
```

---

## 📊 18. Monitoring Workflows

### 18.1 Alert Triage

#### Alert Triage Process

```
┌─────────────────────────────────────────────────────────────────────────┐
│                       ALERT TRIAGE PROCESS                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   Alert Received                                                         │
│        │                                                                │
│        ▼                                                                │
│   ┌─────────────────────┐                                              │
│   │  Is this real?      │  (Automated classification)                  │
│   └──────────┬──────────┘                                              │
│              │                                                          │
│       ┌──────┴──────┐                                                   │
│       │             │                                                   │
│       ▼             ▼                                                   │
│    Yes/No        Yes/No                                                │
│       │             │                                                   │
│       ▼             ▼                                                   │
│   ┌────────┐   ┌────────┐                                               │
│   │  Investigate │  Dismiss    │                                       │
│   │  + Fix   │   │ (no action)│                                       │
│   └────────┘   └────────┘                                               │
│       │                                                            │
│       ▼                                                            │
│   ┌─────────────────────────────────────┐                            │
│   │         Resolution                   │                            │
│   │                                     │                            │
│   │  - Fix issue                         │                            │
│   │  - Document in runbook              │                            │
│   │  - Create alert if missing          │                            │
│   └─────────────────────────────────────┘                            │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Alert Classification

```typescript
interface Alert {
  id: string;
  name: string;
  severity: 'critical' | 'warning' | 'info';
  source: string;
  condition: string;
  for: number; // minutes
}

const alertRules: Alert[] = [
  {
    id: 'alert-001',
    name: 'High Error Rate',
    severity: 'critical',
    source: 'prometheus',
    condition: 'rate(http_requests_total{status=~"5.."}[5m]) > 0.05',
    for: 5,
  },
  {
    id: 'alert-002',
    name: 'High Latency P99',
    severity: 'warning',
    source: 'prometheus',
    condition: 'histogram_quantile(0.99, http_request_duration_seconds_bucket) > 2',
    for: 10,
  },
  {
    id: 'alert-003',
    name: 'Database Connection Pool Exhausted',
    severity: 'critical',
    source: 'prometheus',
    condition: 'pg_stat_activity_count > 80',
    for: 2,
  },
  {
    id: 'alert-004',
    name: 'Disk Space Low',
    severity: 'warning',
    source: 'prometheus',
    condition: '(disk_free_bytes / disk_total_bytes) < 0.1',
    for: 30,
  },
  {
    id: 'alert-005',
    name: 'Memory Usage High',
    severity: 'warning',
    source: 'prometheus',
    condition: '(memory_used_bytes / memory_total_bytes) > 0.85',
    for: 15,
  },
];
```

#### Triage Checklist

```markdown
## Alert Triage Checklist

### Initial Assessment
- [ ] Alert acknowledged
- [ ] Severity confirmed
- [ ] Impact understood

### Investigation
- [ ] Checked Grafana dashboard
- [ ] Reviewed recent deployments
- [ ] Checked related alerts
- [ ] Reviewed error logs

### Root Cause
- [ ] Root cause identified
- [ ] Related incidents found

### Resolution
- [ ] Fix applied
- [ ] Verification successful
- [ ] Post-incident actions documented

### Follow-up
- [ ] Runbook updated (if needed)
- [ ] Alert tuned (if needed)
- [ ] Team notified (if needed)
```

---

### 18.2 Dashboard Reviews

#### Daily Dashboard Review

```markdown
# Daily Dashboard Review

**Date:** 2026-02-15
**Reviewer:** [Name]

---

## Key Metrics

| Metric | Current | Target | Status |
|--------|---------|--------|--------|
| Uptime | 99.95% | ≥99.9% | ✅ |
| Error Rate | 0.3% | <1% | ✅ |
| P50 Latency | 120ms | <200ms | ✅ |
| P99 Latency | 850ms | <1s | ✅ |
| CPU Usage | 45% | <80% | ✅ |
| Memory Usage | 62% | <85% | ✅ |

---

## Services Status

| Service | Status | Notes |
|---------|--------|-------|
| API | ✅ Healthy | |
| Web | ✅ Healthy | |
| Database | ✅ Healthy | |
| Cache | ✅ Healthy | |
| Queue | ✅ Healthy | |

---

## Alerts Triggered (24h)

| Alert | Count | Resolved |
|-------|-------|----------|
| High Error Rate | 2 | 2 |
| High Latency | 1 | 1 |
| Disk Space | 0 | - |

---

## Notable Events

- 10:00 - Deployment v2.1.0 (success)
- 14:30 - Database maintenance (no impact)

---

## Action Items

| Item | Owner | Due |
|------|-------|-----|
| None | | |

---

## Notes

[Additional observations or concerns]
```

---

### 18.3 Capacity Planning

#### Capacity Metrics

```typescript
interface CapacityMetrics {
  service: string;
  currentUsage: {
    cpu: number;      // percentage
    memory: number;   // percentage
    storage: number;  // percentage
    requests: number; // per second
  };
  capacity: {
    cpu: number;
    memory: number;
    storage: number;
    requests: number;
  };
  trends: {
    cpuGrowth: number;      // % per month
    memoryGrowth: number;
    requestGrowth: number;
  };
  forecast: {
    cpuExhaustion: Date;
    memoryExhaustion: Date;
  };
}

const capacityThresholds = {
  warning: 75,  // Trigger warning at 75%
  critical: 90, // Trigger critical at 90%
};

function calculateCapacityForecast(metrics: CapacityMetrics): CapacityForecast {
  const monthsUntilCPU = Math.log(capacityThresholds.critical / metrics.currentUsage.cpu) 
    / Math.log(1 + metrics.trends.cpuGrowth / 100);
  
  const cpuExhaustionDate = new Date();
  cpuExhaustionDate.setMonth(cpuExhaustionDate.getMonth() + monthsUntilCPU);
  
  return {
    cpuExhaustion: cpuExhaustionDate,
    recommendedAction: monthsUntilCPU < 3 ? 'scale_up' : 'monitor',
  };
}
```

#### Scaling Recommendations

```markdown
# Capacity Planning Report - February 2026

## Current Capacity

| Service | CPU | Memory | Storage | Requests/s |
|---------|-----|--------|---------|------------|
| API | 45% | 62% | 38% | 1,200 |
| Web | 52% | 58% | 45% | 800 |
| Database | 38% | 71% | 55% | 500 |

## Growth Trends

- API Requests: +15% month-over-month
- Storage: +8% month-over-month
- Database Queries: +12% month-over-month

## Forecast

| Resource | Current Peak | 3-Month Forecast | 6-Month Forecast |
|----------|--------------|-------------------|------------------|
| CPU | 65% | 78% | 95% |
| Memory | 72% | 85% | 102% |
| Storage | 55% | 65% | 80% |

## Recommendations

### Immediate (This Month)
- [ ] Increase database connection pool from 100 to 150

### Short-term (3 Months)
- [ ] Scale API pods from 10 to 15
- [ ] Upgrade database to larger instance

### Long-term (6 Months)
- [ ] Implement read replicas for database
- [ ] Add CDN for static assets
- [ ] Consider database sharding

## Budget Impact

Estimated monthly cost increase: €500
```

---

### 18.4 Performance Tuning

#### Performance Optimization Process

```
┌─────────────────────────────────────────────────────────────────────────┐
│               PERFORMANCE OPTIMIZATION PROCESS                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   Performance Issue Identified                                          │
│             │                                                          │
│             ▼                                                          │
│   ┌─────────────────┐                                                  │
│   │  Baseline       │  (Establish current performance)                │
│   │  Measurement    │                                                  │
│   └────────┬────────┘                                                  │
│            │                                                           │
│            ▼                                                           │
│   ┌─────────────────┐                                                  │
│   │  Analysis       │  (Identify bottlenecks)                         │
│   │                 │                                                  │
│   │  - CPU Profile  │                                                  │
│   │  - Memory       │                                                  │
│   │  - Database    │                                                  │
│   │  - Network     │                                                  │
│   └────────┬────────┘                                                  │
│            │                                                           │
│            ▼                                                           │
│   ┌─────────────────┐                                                  │
│   │  Identify       │  (Root cause)                                   │
│   │  Root Cause     │                                                  │
│   └────────┬────────┘                                                  │
│            │                                                           │
│            ▼                                                           │
│   ┌─────────────────┐                                                  │
│   │  Implement      │  (Fix)                                          │
│   │  Fix            │                                                  │
│   └────────┬────────┘                                                  │
│            │                                                           │
│            ▼                                                          │
│   ┌─────────────────┐                                                  │
│   │  Verify         │  (Compare to baseline)                          │
│   │  Improvement    │                                                  │
│   └────────┬────────┘                                                  │
│            │                                                           │
│            ▼                                                           │
│   ┌─────────────────┐                                                  │
│   │  Document       │  (Add to runbook)                              │
│   │  & Monitor      │                                                  │
│   └─────────────────┘                                                  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Common Optimization Techniques

```typescript
const performanceOptimizations = {
  database: [
    {
      issue: 'N+1 Query Problem',
      detection: 'Many small queries in loop',
      solution: 'Use eager loading or batch queries',
      example: `
        // Before
        for (const user of users) {
          const posts = await db.query('SELECT * FROM posts WHERE user_id = ?', user.id);
        }
        
        // After
        const posts = await db.query('SELECT * FROM posts WHERE user_id IN (?)', users.map(u => u.id));
        const postsByUser = groupBy(posts, 'user_id');
      `,
    },
    {
      issue: 'Missing Index',
      detection: 'Slow queries on large tables',
      solution: 'Add appropriate index',
      example: `
        CREATE INDEX idx_users_email ON users(email);
      `,
    },
    {
      issue: 'Connection Pool Exhaustion',
      detection: 'Connection pool errors under load',
      solution: 'Optimize query performance, increase pool size',
    },
  ],
  
  caching: [
    {
      issue: 'Repeated Expensive Computations',
      detection: 'Same calculation done multiple times',
      solution: 'Cache results with TTL',
    },
    {
      issue: 'Cache Miss Rate High',
      detection: 'Low hit rate on cache',
      solution: 'Review cache key strategy, increase TTL',
    },
  ],
  
  code: [
    {
      issue: 'Synchronous Operations Blocking',
      detection: 'High latency under load',
      solution: 'Use async/await, web workers',
    },
    {
      issue: 'Large Bundle Size',
      detection: 'Slow initial page load',
      solution: 'Code splitting, tree shaking, lazy loading',
    },
  ],
};
```

#### Performance Benchmark Template

```markdown
# Performance Benchmark Report

**Date:** 2026-02-15
**Test:** API Response Time Optimization
**Objective:** Reduce P99 latency by 50%

---

## Baseline

| Metric | Before | Target | After | Improvement |
|--------|--------|--------|-------|-------------|
| P50 | 120ms | 80ms | 75ms | 37.5% |
| P95 | 350ms | 200ms | 180ms | 48.6% |
| P99 | 850ms | 400ms | 380ms | 55.3% |
| Error Rate | 0.3% | <0.5% | 0.2% | 33% |

---

## Changes Made

1. Added Redis caching for user profiles (TTL: 5 min)
2. Optimized database query with composite index
3. Implemented connection pooling

---

## Test Environment

- Load: 1000 RPS for 10 minutes
- Region: eu-central-1
- Instance: production-sized

---

## Conclusion

Target achieved. P99 reduced by 55%, well below 400ms target.
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
