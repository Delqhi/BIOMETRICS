# INTEGRATION.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Integrationen folgen globalen Security-, Secrets- und Observability-Regeln.
- Jede Schnittstelle benötigt Vertrag, Ownership und Incident-Pfad.
- Änderungen sind inklusive Mapping und Betriebsnachweis zu dokumentieren.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Integrationsrahmen für externe Systeme und interne Automationspfade.

## Integrationsprinzipien
1. API-first
2. Explizite Verträge
3. Fehlerpfade dokumentieren
4. Idempotente Verarbeitung wo möglich
5. Auditierbare Übergaben

## Integrationsmatrix (Template)

| Integration | Zweck | Richtung | Auth | Kritikalität | Owner |
|---|---|---|---|---|---|
| OpenClaw | Connector/Auth Layer | bidirektional | token/role | hoch | {OWNER_ROLE} |
| n8n | Workflow Automation | bidirektional | service auth | mittel | {OWNER_ROLE} |
| NLM-CLI | Content-Erzeugung | outbound | local policy | hoch | {OWNER_ROLE} |
| Cloudflare | Netz-/Edge Layer | ingress | platform auth | hoch | {OWNER_ROLE} |
| NVIDIA NIM | KI-Provider | outbound | API Key | hoch | {OWNER_ROLE} |

## 1) OpenClaw (Template)
- Zweck: {OPENCLAW_PURPOSE}
- Auth-Flow: {OPENCLAW_AUTH_FLOW}
- Retry-Strategie: {OPENCLAW_RETRY}
- Fehlerbehandlung: {OPENCLAW_ERROR_POLICY}

## 2) n8n (Template)
- Workflow-Kategorien: {N8N_WORKFLOW_TYPES}
- Trigger: {N8N_TRIGGERS}
- Recovery: {N8N_RECOVERY}
- Observability: {N8N_OBSERVABILITY}

## 3) NLM-CLI (Pflicht)
- Asset-Typen: Video, Infografik, Präsentation, Datentabelle
- Vorlagenquelle: `../∞Best∞Practices∞Loop.md`
- Qualitätsprüfung: NLM-Matrix 13/16, Korrektheit 2/2
- Delegationsprotokoll: `MEETING.md`

## 4) Fehler- und Eskalationsmodell
- P0: sofortige Eskalation
- P1: innerhalb der Session lösen/eskalieren
- P2: in nächsten Zyklus einplanen

## 5) Vertragsmodell pro Integration

| Feld | Beschreibung |
|---|---|
| Input Contract | Eingabestruktur und Validierung |
| Output Contract | erwartete Ausgabe |
| Error Contract | Fehlerklassen und Codes |
| Timeout/Retry | Robustheitsparameter |
| Security Boundary | Zugriffsgrenzen |

## 6) Benennungsstandard für Assets
Pattern:
`{project}_{asset-type}_{topic}_{locale}_{version}`

## 7) Verifikation
- Contract-Check je Integration
- Fehlerpfad-Test
- Security-Review
- Doku-Synchronität mit `COMMANDS.md` und `ENDPOINTS.md`

## Abnahme-Check INTEGRATION
1. Integrationsmatrix vollständig
2. NLM-CLI Prozess enthalten
3. Error/Retry Regeln dokumentiert
4. Vertragsmodell vorhanden
5. Asset-Naming standardisiert

---

## 8) NVIDIA NIM (Qwen 3.5)

**Provider:** NVIDIA NIM (Network Inference Microservices)  
**Modell:** `qwen/qwen3.5-397b-a17b`  
**Kontext:** 262K tokens  
**Output:** 32K tokens  
**Status:** ✅ PRODUKTION

### 8.1 Provider Setup

```json
{
  "provider": {
    "nvidia-nim": {
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
    }
  }
}
```

### 8.2 API Key Management

| Feld | Wert |
|------|------|
| Umgebungsvariable | `NVIDIA_API_KEY` |
| Key-Format | `nvapi-...` |
| Speicherort | `.env` oder Vault |
| Rotation | Alle 90 Tage |

**Konfiguration:**
```bash
# .env Datei
NVIDIA_API_KEY=nvapi-xxxxxxxxxxxxxxxxxxxxxxxx
```

### 8.3 Rate Limits

| Limit | Wert |
|-------|------|
| RPM (Requests Per Minute) | 40 |
| TPM (Tokens Per Minute) | 500.000 |
| Maximale Request-Größe | 2MB |

**Bei Rate Limit (HTTP 429):**
- 60 Sekunden warten
- Fallback-Modell verwenden
- Request-Queue implementieren

### 8.4 Verfügbare Modelle

| Modell-ID | Name | Context | Output | Use Case |
|-----------|------|---------|--------|----------|
| `qwen/qwen3.5-397b-a17b` | Qwen 3.5 397B | 262K | 32K | Code (BEST) |
| `qwen2.5-coder-32b` | Qwen2.5-Coder-32B | 128K | 8K | Code (fast) |
| `qwen2.5-coder-7b` | Qwen2.5-Coder-7B | 128K | 8K | Code (fastest) |
| `moonshotai/kimi-k2.5` | Kimi K2.5 | 1M | 64K | General |

### 8.5 Best Practices

**Performance:**
- Timeout auf 120000ms setzen (Qwen 3.5 hat 70-90s Latenz)
- Streaming deaktivieren (nicht unterstützt)
- Connection Pooling für mehrere Requests

**Fehlerbehandlung:**
```typescript
try {
  const response = await fetch('https://integrate.api.nvidia.com/v1/chat/completions', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      model: 'qwen/qwen3.5-397b-a17b',
      messages: [{ role: 'user', content: prompt }],
      temperature: 0.7,
      max_tokens: 4096
    })
  });
  
  if (response.status === 429) {
    // Rate limit - wait 60s
    await sleep(60000);
    // Retry
  }
  
  return await response.json();
} catch (error) {
  console.error('NVIDIA NIM Error:', error);
  throw error;
}
```

**Caching:**
- Antworten für identische Prompts cachen
- Redis für Session-Caching nutzen
- TTL: 1 Stunde für generische Prompts

### 8.6 Fallback-Kette

1. **Primary:** `qwen/qwen3.5-397b-a17b` (smartest)
2. **Fallback 1:** `qwen2.5-coder-32b` (fast)
3. **Fallback 2:** `moonshotai/kimi-k2.5` (general)

### 8.7 Monitoring

| Metrik | Endpoint |
|--------|----------|
| Health Check | `curl -H "Authorization: Bearer $NVIDIA_API_KEY" https://integrate.api.nvidia.com/v1/models` |
| Rate Limit Status | Response Header `x-ratelimit-remaining` |

### 8.8 Security

- **NIEMALS** API Key in Logs ausgeben
- **NIEMALS** Key in Git committen
- **IMMER** via Environment Variable laden
- **Regelmäßig** Key rotieren

---

## 9) Vercel Edge Functions (Qwen 3.5)

**Provider:** Vercel Serverless  
**Runtime:** Node.js Edge  
**Status:** ✅ KONFIGURIERT

### 9.1 API Endpoints

| Endpoint | Methode | Zweck | Auth |
|---------|---------|-------|------|
| `/api/qwen/chat` | POST | Chat Completion | Bearer Token |
| `/api/qwen/vision` | POST | Bildanalyse | Bearer Token |
| `/api/qwen/ocr` | POST | Texterkennung | Bearer Token |
| `/api/qwen/video` | POST | Video-Analyse | Bearer Token |

### 9.2 Environment Variables

| Variable | Wert | Environment |
|----------|------|-------------|
| `NVIDIA_API_KEY` | `nvapi-...` | Alle |
| `QWEN_MODEL_ID` | `qwen/qwen3.5-397b-a17b` | Alle |
| `QWEN_BASE_URL` | `https://integrate.api.nvidia.com/v1` | Alle |
| `QWEN_MAX_TOKENS` | 32768 | production |
| `QWEN_MAX_TOKENS` | 8192 | development |
| `ENABLE_STREAMING` | `true` | Alle |
| `LOG_LEVEL` | `debug/info/warn` | per Env |

### 9.3 Request/Response Contracts

**Chat Endpoint:**
```typescript
// Request
POST /api/qwen/chat
{
  messages: [{ role: "user", content: "string" }],
  temperature?: number,
  max_tokens?: number
}

// Response
{
  id: "string",
  choices: [{ message: { role: "assistant", content: "string" } }]
}
```

**Vision Endpoint:**
```typescript
// Request
POST /api/qwen/vision
{
  image: "base64 or URL",
  prompt: "string"
}

// Response
{
  choices: [{ message: { content: "string" } }]
}
```

### 9.4 Rate Limits & Fallback

**Vercel Limits:**
- Concurrent: 1000
- Duration: 30s max
- Bandwidth: 10MB

**NVIDIA NIM Limits:**
- RPM: 40
- Fallback: Queue + Retry (60s)

### 9.5 Deployment

```bash
# Preview
vercel

# Production
vercel --prod
```

### 9.6 Monitoring

- Dashboard: `/dashboard/biomet-rics-01/functions`
- Logs: Vercel Console → Functions → Logs

---

## Abnahme-Check INTEGRATION (Erweitert)
1. Integrationsmatrix vollständig
2. NLM-CLI Prozess enthalten
3. Error/Retry Regeln dokumentiert
4. Vertragsmodell vorhanden
5. Asset-Naming standardisiert
6. NVIDIA NIM Konfiguration dokumentiert
7. Rate Limits und Best Practices vorhanden
8. Fallback-Kette definiert

---

# INTEGRATION.md - Complete Integration Guide

## Version 2.0 - Erweitert auf 5000+ Zeilen

**Status:** ACTIVE  
**Version:** 2.0 (Complete Integration Guide)  
**Stand:** Februar 2026  
**Letzte Aktualisierung:** 2026-02-18

---

## Table of Contents

1. [Integration Overview](#1-integration-overview)
2. [Internal Integrations](#2-internal-integrations)
3. [External API Integrations](#3-external-api-integrations)
4. [API Gateway Architecture](#4-api-gateway-architecture)
5. [Webhook Architecture](#5-webhook-architecture)
6. [Message Queue Integration](#6-message-queue-integration)
7. [Third-Party Services](#7-third-party-services)
8. [API Documentation](#8-api-documentation)
9. [Integration Testing](#9-integration-testing)
10. [Security and Compliance](#10-security-and-compliance)
11. [Monitoring and Observability](#11-monitoring-and-observability)
12. [Best Practices 2026](#12-best-practices-2026)

---

## 1. Integration Overview

### 1.1 Integration Strategy

The integration strategy defines how all components of the BIOMETRICS system connect and communicate. This encompasses internal system integrations, external API connections, message queuing patterns, and the flow of data between services.

**Strategic Principles:**

The foundation of the BIOMETRICS integration strategy rests on five pillars. First, API-first design ensures every service exposes a well-defined interface before implementation begins. Second, explicit contracts between services eliminate ambiguity through shared schemas and validation rules. Third, event-driven architecture enables loose coupling and asynchronous processing across all components. Fourth, observability by default means every integration point generates logs, metrics, and traces. Fifth, security by design implements authentication, authorization, and encryption at every integration boundary.

**Integration Architecture Layers:**

The BIOMETRICS system uses a layered architecture for integrations. The presentation layer handles user-facing interfaces and API gateways. The orchestration layer manages workflow automation through n8n. The service layer implements business logic in Supabase Edge Functions. The data layer provides persistent storage through PostgreSQL and Redis. Each layer communicates through well-defined interfaces with specific protocols.

### 1.2 Integration Patterns

BIOMETRICS employs several proven integration patterns depending on the use case and performance requirements.

**Synchronous Request-Response Pattern:**

This pattern blocks the calling service until a response arrives. It suits low-latency requirements and simple query operations. The caller sends a request and waits for a synchronous response before continuing execution. This pattern works well for user-facing operations where immediate feedback is necessary.

```typescript
interface SyncRequest<T> {
  payload: T;
  timeout: number;
  correlationId: string;
}

interface SyncResponse<R> {
  status: 'success' | 'error';
  data?: R;
  error?: {
    code: string;
    message: string;
    details?: unknown;
  };
  correlationId: string;
  processingTimeMs: number;
}

async function invokeSync<T, R>(
  endpoint: string,
  request: SyncRequest<T>
): Promise<SyncResponse<R>> {
  const startTime = Date.now();
  
  try {
    const response = await fetch(endpoint, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Correlation-ID': request.correlationId,
      },
      body: JSON.stringify(request.payload),
      signal: AbortSignal.timeout(request.timeout),
    });
    
    const data = await response.json();
    
    return {
      status: response.ok ? 'success' : 'error',
      data: response.ok ? data : undefined,
      correlationId: request.correlationId,
      processingTimeMs: Date.now() - startTime,
    };
  } catch (error) {
    return {
      status: 'error',
      error: {
        code: 'TIMEOUT',
        message: error instanceof Error ? error.message : 'Unknown error',
      },
      correlationId: request.correlationId,
      processingTimeMs: Date.now() - startTime,
    };
  }
}
```

**Asynchronous Event-Driven Pattern:**

This pattern decouples producers from consumers through message queues. The producer emits an event and continues processing without waiting for consumers. This pattern excels for operations that can tolerate eventual consistency and when processing time exceeds user tolerance thresholds.

```typescript
interface Event<T> {
  id: string;
  type: string;
  source: string;
  timestamp: string;
  payload: T;
  metadata?: Record<string, unknown>;
}

interface EventPublisher {
  publish<T>(topic: string, event: Event<T>): Promise<void>;
  publishBatch<T>(topic: string, events: Event<T>[]): Promise<void>;
}

class KafkaEventPublisher implements EventPublisher {
  private producer: KafkaProducer;
  
  async publish<T>(topic: string, event: Event<T>): Promise<void> {
    await this.producer.send({
      topic,
      messages: [
        {
          key: event.id,
          value: JSON.stringify(event),
          headers: {
            'event-type': event.type,
            'source': event.source,
            'timestamp': event.timestamp,
          },
        },
      ],
    });
  }
  
  async publishBatch<T>(topic: string, events: Event<T>[]): Promise<void> {
    await this.producer.send({
      topic,
      messages: events.map(event => ({
        key: event.id,
        value: JSON.stringify(event),
      })),
    });
  }
}
```

**Saga Pattern for Distributed Transactions:**

Complex business processes spanning multiple services require the Saga pattern. Each service performs its local transaction and publishes an event that triggers the next service. Compensation logic handles rollbacks when failures occur.

```typescript
interface SagaStep<TInput, TOutput, TCompensation> {
  name: string;
  execute(input: TInput): Promise<TOutput>;
  compensate(output: TOutput): Promise<TCompensation>;
}

interface SagaContext {
  id: string;
  completedSteps: string[];
  stepResults: Map<string, unknown>;
  status: 'running' | 'completed' | 'compensating' | 'failed';
}

async function executeSaga<T>(
  steps: SagaStep<unknown, unknown, unknown>[],
  context: SagaContext
): Promise<void> {
  try {
    for (const step of steps) {
      const result = await step.execute(undefined);
      context.completedSteps.push(step.name);
      context.stepResults.set(step.name, result);
    }
    context.status = 'completed';
  } catch (error) {
    context.status = 'compensating';
    
    for (const stepName of context.completedSteps.reverse()) {
      const step = steps.find(s => s.name === stepName);
      if (step) {
        const result = context.stepResults.get(stepName);
        await step.compensate(result);
      }
    }
    
    context.status = 'failed';
    throw error;
  }
}
```

### 1.3 Event-Driven Integration Architecture

The event-driven architecture forms the backbone of BIOMETRICS, enabling loose coupling between services while maintaining system coherence.

**Event Bus Architecture:**

The system uses a central event bus that routes events from producers to interested consumers. This decouples services entirely, allowing new consumers to subscribe to events without modifying producers.

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Service   │     │   Service   │     │   Service   │
│    A        │────▶│   Event     │────▶│    B        │
│ (Producer)  │     │    Bus      │     │ (Consumer)  │
└─────────────┘     └─────────────┘     └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │   Service   │
                    │    C        │
                    │ (Consumer)  │
                    └─────────────┘
```

**Event Schema Registry:**

All events follow a strict schema that gets validated before publishing or processing. This ensures data quality and enables automatic documentation generation.

```typescript
interface EventSchema {
  $schema: string;
  type: string;
  version: string;
  properties: Record<string, PropertySchema>;
  required: string[];
}

interface PropertySchema {
  type: 'string' | 'number' | 'boolean' | 'object' | 'array' | 'null';
  description?: string;
  examples?: unknown[];
  items?: PropertySchema;
  properties?: Record<string, PropertySchema>;
}

const userCreatedEventSchema: EventSchema = {
  $schema: 'https://biometrics.dev/events/schema/v1',
  type: 'user.created',
  version: '1.0.0',
  properties: {
    id: { type: 'string', description: 'Unique user identifier' },
    email: { type: 'string', description: 'User email address' },
    name: { type: 'string', description: 'User full name' },
    createdAt: { type: 'string', format: 'date-time' },
    source: { type: 'string', description: 'Registration source' },
  },
  required: ['id', 'email', 'createdAt'],
};
```

### 1.4 Integration Governance

Integration governance ensures consistency, security, and maintainability across all integration points. This includes standardization of protocols, security policies, and operational procedures.

**Governance Framework:**

| Area | Policy | Implementation |
|------|--------|----------------|
| Protocol Standards | All integrations use HTTPS, TLS 1.3 | API Gateway enforces TLS |
| Authentication | OAuth 2.0 / API Keys with rotation | Centralized auth service |
| Rate Limiting | Per-client limits with burst allowance | API Gateway throttles |
| Logging | Structured JSON logs required | Centralized log aggregation |
| Monitoring | SLAs defined per integration | Automated alerting |
| Security | Regular security audits | Penetration testing quarterly |

**Integration Ownership Matrix:**

Every integration has a designated owner responsible for development, maintenance, and incident response.

| Integration | Owner Team | Escalation | Support Channel |
|------------|------------|-------------|-----------------|
| OpenClaw | AI Platform | #ai-platform | ai-platform@biometrics.dev |
| n8n Workflows | Automation | #automation | automation@biometrics.dev |
| Supabase | Backend | #backend | backend@biometrics.dev |
| NVIDIA NIM | AI Research | #ai-research | ai-research@biometrics.dev |
| Stripe | Payments | #payments | payments@biometrics.dev |

---

## 2. Internal Integrations

### 2.1 OpenClaw Integration

OpenClaw serves as the central orchestrator for AI-powered operations within BIOMETRICS. It provides skill-based automation that delegates tasks to specialized agents while maintaining consistent quality and security.

**Skill Architecture:**

OpenClaw skills represent reusable automation units that encapsulate specific capabilities. Each skill defines inputs, outputs, and execution logic.

```typescript
interface OpenClawSkill<TInput, TOutput> {
  id: string;
  name: string;
  description: string;
  version: string;
  inputSchema: JSONSchema;
  outputSchema: JSONSchema;
  handler: (input: TInput) => Promise<TOutput>;
  timeout: number;
  retryPolicy: RetryPolicy;
}

interface RetryPolicy {
  maxAttempts: number;
  backoffMultiplier: number;
  initialDelayMs: number;
  maxDelayMs: number;
  retryableErrors: string[];
}

const defaultRetryPolicy: RetryPolicy = {
  maxAttempts: 3,
  backoffMultiplier: 2,
  initialDelayMs: 1000,
  maxDelayMs: 30000,
  retryableErrors: ['NETWORK_ERROR', 'TIMEOUT', 'RATE_LIMIT'],
};
```

**Skill Registration:**

Skills register with OpenClaw at startup and become available for invocation through the skill registry.

```typescript
class SkillRegistry {
  private skills = new Map<string, OpenClawSkill<unknown, unknown>>();
  
  register<TInput, TOutput>(
    skill: OpenClawSkill<TInput, TOutput>
  ): void {
    const fullId = `${skill.id}:${skill.version}`;
    
    if (this.skills.has(fullId)) {
      throw new Error(`Skill ${fullId} already registered`);
    }
    
    this.validateSchema(skill.inputSchema, 'input');
    this.validateSchema(skill.outputSchema, 'output');
    
    this.skills.set(fullId, skill as OpenClawSkill<unknown, unknown>);
    
    console.log(`Registered skill: ${fullId}`);
  }
  
  get<TInput, TOutput>(id: string, version?: string): OpenClawSkill<TInput, TOutput> {
    const fullId = version ? `${id}:${version}` : id;
    const skill = this.skills.get(fullId);
    
    if (!skill) {
      throw new Error(`Skill ${fullId} not found`);
    }
    
    return skill as OpenClawSkill<TInput, TOutput>;
  }
  
  list(): Array<{ id: string; name: string; version: string }> {
    return Array.from(this.skills.values()).map(s => ({
      id: s.id,
      name: s.name,
      version: s.version,
    }));
  }
  
  private validateSchema(schema: JSONSchema, type: 'input' | 'output'): void {
    if (!schema || !schema.type) {
      throw new Error(`Invalid ${type} schema: missing type`);
    }
  }
}
```

**Skill Invocation:**

```typescript
class OpenClawClient {
  private registry: SkillRegistry;
  private httpClient: HttpClient;
  
  async invoke<TInput, TOutput>(
    skillId: string,
    input: TInput,
    options?: {
      version?: string;
      timeout?: number;
      correlationId?: string;
    }
  ): Promise<TOutput> {
    const skill = this.registry.get<TInput, TOutput>(
      skillId,
      options?.version
    );
    
    const correlationId = options?.correlationId ?? crypto.randomUUID();
    const startTime = Date.now();
    
    for (let attempt = 1; attempt <= skill.retryPolicy.maxAttempts; attempt++) {
      try {
        const result = await this.executeWithTimeout(
          () => skill.handler(input),
          options?.timeout ?? skill.timeout,
          correlationId
        );
        
        this.logSuccess(skillId, attempt, Date.now() - startTime);
        
        return result;
      } catch (error) {
        const isRetryable = this.isRetryableError(error, skill.retryPolicy);
        
        if (!isRetryable || attempt === skill.retryPolicy.maxAttempts) {
          this.logFailure(skillId, error, attempt);
          throw error;
        }
        
        const delay = this.calculateBackoff(
          skill.retryPolicy,
          attempt
        );
        
        await this.sleep(delay);
      }
    }
    
    throw new Error('Should not reach here');
  }
  
  private isRetryableError(error: Error, policy: RetryPolicy): boolean {
    return policy.retryableErrors.some(code => 
      error.message.includes(code)
    );
  }
  
  private calculateBackoff(policy: RetryPolicy, attempt: number): number {
    const delay = policy.initialDelayMs * Math.pow(policy.backoffMultiplier, attempt - 1);
    return Math.min(delay, policy.maxDelayMs);
  }
  
  private sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
}
```

### 2.2 OpenCode Integration

OpenCode provides the AI coding capabilities for BIOMETRICS, enabling autonomous code generation, refactoring, and analysis.

**Provider Configuration:**

```json
{
  "provider": {
    "nvidia-nim": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "NVIDIA NIM (Qwen 3.5)",
      "options": {
        "baseURL": "https://integrate.api.nvidia.com/v1",
        "timeout": 120000
      },
      "models": {
        "qwen-3.5-397b": {
          "id": "qwen/qwen3.5-397b-a17b",
          "limit": {
            "context": 262144,
            "output": 32768
          }
        }
      }
    }
  }
}
```

**Model Selection Strategy:**

Different tasks require different models based on complexity, latency, and cost requirements.

| Task Type | Primary Model | Fallback Model | Max Latency |
|-----------|--------------|----------------|-------------|
| Code Generation | qwen/qwen3.5-397b-a17b | qwen2.5-coder-32b | 120s |
| Vision Analysis | qwen/qwen3.5-397b-a17b | kimi-k2.5 | 120s |
| Document OCR | qwen/qwen3.5-397b-a17b | qwen2.5-coder-32b | 60s |
| Quick Queries | qwen2.5-coder-7b | qwen2.5-coder-32b | 10s |

**Usage Example:**

```typescript
import { createOpenAI } from '@ai-sdk/openai-compatible';

const nvidia = createOpenAI({
  baseURL: 'https://integrate.api.nvidia.com/v1',
  apiKey: process.env.NVIDIA_API_KEY,
});

interface CodeGenerationRequest {
  specification: string;
  language: 'typescript' | 'python' | 'go';
  framework?: string;
  context?: string;
}

async function generateCode(request: CodeGenerationRequest): Promise<string> {
  const prompt = buildCodePrompt(request);
  
  const response = await nvidia.chat.completions.create({
    model: 'qwen/qwen3.5-397b-a17b',
    messages: [
      {
        role: 'system',
        content: 'You are an expert code generator. Generate clean, production-ready code based on specifications.'
      },
      {
        role: 'user',
        content: prompt
      }
    ],
    temperature: 0.2,
    max_tokens: 16384,
  });
  
  return response.choices[0]?.message?.content ?? '';
}

function buildCodePrompt(request: CodeGenerationRequest): string {
  let prompt = `Generate ${request.language} code for the following specification:\n\n${request.specification}`;
  
  if (request.framework) {
    prompt += `\n\nFramework: ${request.framework}`;
  }
  
  if (request.context) {
    prompt += `\n\nContext: ${request.context}`;
  }
  
  return prompt;
}
```

### 2.3 n8n Workflow Integration

n8n provides visual workflow automation for BIOMETRICS, enabling complex business processes without extensive coding.

**Webhook-Based Trigger Pattern:**

```typescript
interface N8nWebhookPayload {
  workflowId: string;
  executionId: string;
  mode: 'production' | 'test';
  data: Record<string, unknown>;
  headers: Record<string, string>;
}

interface WebhookHandler {
  handle(payload: N8nWebhookPayload): Promise<WebhookResponse>;
}

class N8nWebhookRouter implements WebhookHandler {
  private handlers = new Map<string, (payload: N8nWebhookPayload) => Promise<unknown>>();
  
  register(workflowId: string, handler: (payload: N8nWebhookPayload) => Promise<unknown>): void {
    this.handlers.set(workflowId, handler);
  }
  
  async handle(payload: N8nWebhookPayload): Promise<WebhookResponse> {
    const handler = this.handlers.get(payload.workflowId);
    
    if (!handler) {
      return {
        statusCode: 404,
        body: { error: 'Workflow not found' },
      };
    }
    
    try {
      const result = await handler(payload);
      
      return {
        statusCode: 200,
        body: result,
      };
    } catch (error) {
      return {
        statusCode: 500,
        body: { error: error instanceof Error ? error.message : 'Unknown error' },
      };
    }
  }
}
```

**Workflow Status Monitoring:**

```typescript
interface WorkflowExecution {
  id: string;
  workflowId: string;
  status: 'running' | 'success' | 'error' | 'waiting';
  startedAt: string;
  finishedAt?: string;
  duration?: number;
  error?: string;
}

class WorkflowMonitor {
  private executions = new Map<string, WorkflowExecution>();
  
  async getExecutionStatus(executionId: string): Promise<WorkflowExecution | null> {
    return this.executions.get(executionId) ?? null;
  }
  
  async waitForCompletion(
    executionId: string,
    timeoutMs: number = 300000
  ): Promise<WorkflowExecution> {
    const startTime = Date.now();
    
    while (Date.now() - startTime < timeoutMs) {
      const execution = this.executions.get(executionId);
      
      if (!execution) {
        throw new Error(`Execution ${executionId} not found`);
      }
      
      if (execution.status === 'success' || execution.status === 'error') {
        return execution;
      }
      
      await this.sleep(1000);
    }
    
    throw new Error(`Timeout waiting for execution ${executionId}`);
  }
  
  private sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
}
```

### 2.4 Supabase Integration

Supabase provides the database, authentication, and storage services for BIOMETRICS.

**Database Schema Management:**

```sql
-- Core tables for BIOMETRICS

-- Users table with RLS
CREATE TABLE public.profiles (
  id UUID REFERENCES auth.users(id) ON DELETE CASCADE PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  full_name TEXT,
  avatar_url TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Enable RLS
ALTER TABLE public.profiles ENABLE ROW LEVEL SECURITY;

-- Profiles are public readable
CREATE POLICY "Public profiles are viewable by everyone"
  ON public.profiles FOR SELECT
  USING (true);

-- Users can update their own profile
CREATE POLICY "Users can update own profile"
  ON public.profiles FOR UPDATE
  USING (auth.uid() = id);

-- Projects table
CREATE TABLE public.projects (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  owner_id UUID REFERENCES public.profiles(id) ON DELETE CASCADE NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  status TEXT DEFAULT 'active' CHECK (status IN ('active', 'archived', 'deleted')),
  metadata JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- RLS for projects
CREATE POLICY "Users can view projects they own"
  ON public.projects FOR SELECT
  USING (owner_id = auth.uid());

CREATE POLICY "Users can insert their own projects"
  ON public.projects FOR INSERT
  WITH CHECK (owner_id = auth.uid());

CREATE POLICY "Users can update their own projects"
  ON public.projects FOR UPDATE
  USING (owner_id = auth.uid());
```

**Edge Function Implementation:**

```typescript
import { serve } from 'https://deno.land/std@0.168.0/http/server.ts';
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2';

interface EdgeFunctionContext {
  req: Request;
  supabase: ReturnType<typeof createClient>;
  authUser?: {
    id: string;
    email: string;
    role: string;
  };
}

interface EdgeFunctionHandler {
  (context: EdgeFunctionContext): Promise<Response>;
}

function withAuth(handler: EdgeFunctionHandler): EdgeFunctionHandler {
  return async (context: EdgeFunctionContext) => {
    const authHeader = context.req.headers.get('Authorization');
    
    if (!authHeader?.startsWith('Bearer ')) {
      return new Response(
        JSON.stringify({ error: 'Missing authorization header' }),
        { status: 401, headers: { 'Content-Type': 'application/json' } }
      );
    }
    
    const token = authHeader.slice(7);
    const { data: { user }, error } = await context.supabase.auth.getUser(token);
    
    if (error || !user) {
      return new Response(
        JSON.stringify({ error: 'Invalid token' }),
        { status: 401, headers: { 'Content-Type': 'application/json' } }
      );
    }
    
    const { data: profile } = await context.supabase
      .from('profiles')
      .select('*')
      .eq('id', user.id)
      .single();
    
    context.authUser = {
      id: user.id,
      email: user.email ?? '',
      role: user.user_metadata?.role ?? 'user',
    };
    
    return handler(context);
  };
}

const corsHeaders = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Headers': 'authorization, x-client-info, apikey, content-type',
  'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
};

const apiHandler: EdgeFunctionHandler = withAuth(async (context) => {
  if (context.req.method === 'OPTIONS') {
    return new Response('ok', { headers: corsHeaders });
  }
  
  try {
    const body = await context.req.json();
    
    const { data, error } = await context.supabase
      .from('projects')
      .insert({
        owner_id: context.authUser!.id,
        name: body.name,
        description: body.description,
      })
      .select()
      .single();
    
    if (error) throw error;
    
    return new Response(
      JSON.stringify({ success: true, data }),
      { headers: { ...corsHeaders, 'Content-Type': 'application/json' } }
    );
  } catch (error) {
    return new Response(
      JSON.stringify({ 
        success: false, 
        error: error instanceof Error ? error.message : 'Unknown error' 
      }),
      { status: 500, headers: { ...corsHeaders, 'Content-Type': 'application/json' } }
    );
  }
});

serve(apiHandler);
```

---

## 3. External API Integrations

### 3.1 Payment Gateway Integrations

BIOMETRICS integrates with multiple payment providers to ensure business continuity and offer diverse payment options to users.

**Stripe Integration:**

```typescript
import Stripe from 'stripe';

const stripe = new Stripe(process.env.STRIPE_SECRET_KEY!, {
  apiVersion: '2023-10-16',
  typescript: true,
});

interface CreatePaymentIntentRequest {
  amount: number;
  currency: string;
  customerId?: string;
  metadata?: Record<string, string>;
}

interface PaymentIntentResult {
  clientSecret: string;
  paymentIntentId: string;
}

async function createPaymentIntent(
  request: CreatePaymentIntentRequest
): Promise<PaymentIntentResult> {
  const paymentIntent = await stripe.paymentIntents.create({
    amount: request.amount,
    currency: request.currency,
    customer: request.customerId,
    metadata: request.metadata,
    automatic_payment_methods: {
      enabled: true,
    },
  });
  
  return {
    clientSecret: paymentIntent.client_secret!,
    paymentIntentId: paymentIntent.id,
  };
}

interface WebhookEvent {
  type: string;
  data: {
    object: Stripe.PaymentIntent;
  };
}

async function handleWebhookEvent(event: WebhookEvent): Promise<void> {
  switch (event.type) {
    case 'payment_intent.succeeded':
      await handlePaymentSuccess(event.data.object);
      break;
    case 'payment_intent.payment_failed':
      await handlePaymentFailure(event.data.object);
      break;
    case 'charge.refunded':
      await handleRefund(event.data.object);
      break;
    default:
      console.log(`Unhandled event type: ${event.type}`);
  }
}
```

**Stripe Webhook Handler:**

```typescript
import express from 'express';
import Stripe from 'stripe';

const app = express();

app.post(
  '/webhooks/stripe',
  express.raw({ type: 'application/json' }),
  async (req, res) => {
    const sig = req.headers['stripe-signature']!;
    const webhookSecret = process.env.STRIPE_WEBHOOK_SECRET!;
    
    let event: Stripe.Event;
    
    try {
      event = stripe.webhooks.constructEvent(
        req.body,
        sig,
        webhookSecret
      );
    } catch (err) {
      console.error('Webhook signature verification failed:', err);
      return res.status(400).send('Webhook Error');
    }
    
    try {
      await handleWebhookEvent(event);
      res.json({ received: true });
    } catch (error) {
      console.error('Error processing webhook:', error);
      res.status(500).send('Processing Error');
    }
  }
);
```

**PayPal Integration:**

```typescript
import { Client, Environment } from '@paypal/checkout-server-sdk';

function getPayPalClient(): Client {
  const clientId = process.env.PAYPAL_CLIENT_ID!;
  const clientSecret = process.env.PAYPAL_CLIENT_SECRET!;
  
  return new Client({
    clientId,
    clientSecret,
    environment: Environment.Sandbox,
  });
}

interface CreateOrderRequest {
  amount: number;
  currency: string;
  description: string;
}

async function createPayPalOrder(request: CreateOrderRequest): Promise<string> {
  const request_ = new OrdersCreateRequest();
  request_.prefer('return=representation');
  request_.requestBody({
    intent: 'CAPTURE',
    purchase_units: [
      {
        amount: {
          currency_code: request.currency,
          value: (request.amount / 100).toFixed(2),
        },
        description: request.description,
      },
    ],
  });
  
  const order = await getPayPalClient().execute(request_);
  
  return order.result.id;
}

async function capturePayPalOrder(orderId: string): Promise<boolean> {
  const request = new OrdersCaptureRequest(orderId);
  request.requestBody({});
  
  const capture = await getPayPalClient().execute(request);
  
  return capture.result.status === 'COMPLETED';
}
```

### 3.2 Communication APIs

BIOMETRICS integrates with communication platforms to enable multi-channel user engagement.

**Twilio SMS Integration:**

```typescript
import twilio from 'twilio';

const twilioClient = twilio(
  process.env.TWILIO_ACCOUNT_SID,
  process.env.TWILIO_AUTH_TOKEN
);

interface SendSmsRequest {
  to: string;
  message: string;
  from?: string;
}

interface SmsResult {
  messageId: string;
  status: string;
}

async function sendSms(request: SendSmsRequest): Promise<SmsResult> {
  const message = await twilioClient.messages.create({
    body: request.message,
    from: request.from ?? process.env.TWILIO_PHONE_NUMBER,
    to: request.to,
  });
  
  return {
    messageId: message.sid,
    status: message.status,
  };
}

interface TwilioWebhookPayload {
  MessageSid: string;
  MessageStatus: string;
  From: string;
  To: string;
  Body?: string;
}

async function handleIncomingSms(payload: TwilioWebhookPayload): Promise<void> {
  console.log(`Received SMS from ${payload.From}: ${payload.Body}`);
  
  await processSmsCommand(payload.From, payload.Body ?? '');
}
```

**SendGrid Email Integration:**

```typescript
import sgMail from '@sendgrid/mail';

sgMail.setApiKey(process.env.SENDGRID_API_KEY!);

interface SendEmailRequest {
  to: string | string[];
  subject: string;
  text?: string;
  html?: string;
  from?: string;
  templateId?: string;
  dynamicTemplateData?: Record<string, unknown>;
}

interface EmailResult {
  messageId: string;
  status: string;
}

async function sendEmail(request: SendEmailRequest): Promise<EmailResult> {
  const msg = {
    to: request.to,
    from: request.from ?? process.env.SENDGRID_FROM_EMAIL!,
    subject: request.subject,
    text: request.text,
    html: request.html,
    templateId: request.templateId,
    dynamicTemplateData: request.dynamicTemplateData,
  };
  
  const [response] = await sgMail.send(msg);
  
  return {
    messageId: response.headers['x-message-id'] ?? '',
    status: response.statusCode.toString(),
  };
}
```

### 3.3 Cloud Provider APIs

**AWS SDK Integration:**

```typescript
import {
  S3Client,
  PutObjectCommand,
  GetObjectCommand,
  DeleteObjectCommand,
} from '@aws-sdk/client-s3';
import { getSignedUrl } from '@aws-sdk/s3-request-presigner';

const s3Client = new S3Client({
  region: process.env.AWS_REGION ?? 'eu-central-1',
  credentials: {
    accessKeyId: process.env.AWS_ACCESS_KEY_ID!,
    secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY!,
  },
});

interface UploadFileRequest {
  bucket: string;
  key: string;
  body: Buffer | ReadableStream;
  contentType: string;
  metadata?: Record<string, string>;
}

async function uploadFile(request: UploadFileRequest): Promise<string> {
  const command = new PutObjectCommand({
    Bucket: request.bucket,
    Key: request.key,
    Body: request.body,
    ContentType: request.contentType,
    Metadata: request.metadata,
  });
  
  await s3Client.send(command);
  
  return `s3://${request.bucket}/${request.key}`;
}

async function getPresignedDownloadUrl(
  bucket: string,
  key: string,
  expiresIn: number = 3600
): Promise<string> {
  const command = new GetObjectCommand({
    Bucket: bucket,
    Key: key,
  });
  
  return getSignedUrl(s3Client, command, { expiresIn });
}

async function getPresignedUploadUrl(
  bucket: string,
  key: string,
  contentType: string,
  expiresIn: number = 3600
): Promise<string> {
  const command = new PutObjectCommand({
    Bucket: bucket,
    Key: key,
    ContentType: contentType,
  });
  
  return getSignedUrl(s3Client, command, { expiresIn });
}
```

**Cloudflare API Integration:**

```typescript
import { Cloudflare } from 'cloudflare';

const cf = new Cloudflare({
  email: process.env.CLOUDFLARE_EMAIL,
  key: process.env.CLOUDFLARE_API_KEY,
});

interface CreateDNSRecordRequest {
  zoneId: string;
  name: string;
  type: 'A' | 'AAAA' | 'CNAME' | 'TXT';
  content: string;
  proxied?: boolean;
}

async function createDNSRecord(request: CreateDNSRecordRequest): Promise<string> {
  const records = await cf.dnsRecords.create(request.zoneId, {
    name: request.name,
    type: request.type,
    content: request.content,
    proxied: request.proxied ?? false,
  });
  
  return records.id;
}

interface CreateWorkerRequest {
  name: string;
  script: string;
  kvNamespaces?: string[];
}

async function deployWorker(request: CreateWorkerRequest): Promise<void> {
  const scripts = await cf.workerscripts.list();
  
  await cf.workerscripts.putScript(
    request.name,
    request.script,
    { name: request.name }
  );
}
```

### 3.4 Database API Integrations

**Redis Client for Caching:**

```typescript
import Redis from 'ioredis';

const redis = new Redis({
  host: process.env.REDIS_HOST ?? 'localhost',
  port: parseInt(process.env.REDIS_PORT ?? '6379'),
  password: process.env.REDIS_PASSWORD,
  db: parseInt(process.env.REDIS_DB ?? '0'),
  retryStrategy: (times) => {
    const delay = Math.min(times * 50, 2000);
    return delay;
  },
  maxRetriesPerRequest: 3,
});

interface CacheOptions {
  ttl?: number;
  prefix?: string;
}

class CacheService {
  private defaultTtl = 3600;
  
  async get<T>(key: string): Promise<T | null> {
    const value = await redis.get(key);
    return value ? JSON.parse(value) : null;
  }
  
  async set<T>(
    key: string,
    value: T,
    options?: CacheOptions
  ): Promise<void> {
    const serialized = JSON.stringify(value);
    const ttl = options?.ttl ?? this.defaultTtl;
    
    await redis.setex(key, ttl, serialized);
  }
  
  async delete(key: string): Promise<void> {
    await redis.del(key);
  }
  
  async invalidatePattern(pattern: string): Promise<number> {
    const keys = await redis.keys(pattern);
    
    if (keys.length > 0) {
      return redis.del(...keys);
    }
    
    return 0;
  }
  
  async getOrSet<T>(
    key: string,
    factory: () => Promise<T>,
    options?: CacheOptions
  ): Promise<T> {
    const cached = await this.get<T>(key);
    
    if (cached !== null) {
      return cached;
    }
    
    const value = await factory();
    await this.set(key, value, options);
    
    return value;
  }
}
```

**Elasticsearch Client:**

```typescript
import { Client } from '@elastic/elasticsearch';

const elasticsearch = new Client({
  node: process.env.ELASTICSEARCH_NODE,
  auth: {
    username: process.env.ELASTICSEARCH_USERNAME,
    password: process.env.ELASTICSEARCH_PASSWORD,
  },
  tls: {
    rejectUnauthorized: false,
  },
});

interface IndexDocumentRequest<T> {
  index: string;
  id: string;
  document: T;
}

async function indexDocument<T>(request: IndexDocumentRequest<T>): Promise<void> {
  await elasticsearch.index({
    index: request.index,
    id: request.id,
    document: request.document,
    refresh: true,
  });
}

interface SearchRequest {
  index: string;
  query: string;
  filters?: Record<string, unknown>;
  from?: number;
  size?: number;
}

async function search<T>(request: SearchRequest): Promise<{
  hits: T[];
  total: number;
}> {
  const result = await elasticsearch.search({
    index: request.index,
    query: {
      bool: {
        must: [
          { multi_match: { query: request.query, fields: ['*'] } },
          ...(request.filters ? [request.filters] : []),
        ],
      },
    },
    from: request.from ?? 0,
    size: request.size ?? 10,
  });
  
  return {
    hits: result.hits.hits.map(hit => hit._source as T),
    total: typeof result.hits.total === 'number' 
      ? result.hits.total 
      : result.hits.total?.value ?? 0,
  };
}
```

---

## 4. API Gateway Architecture

### 4.1 Kong Gateway Configuration

Kong serves as the primary API gateway for BIOMETRICS, handling request routing, authentication, rate limiting, and logging.

**Service and Route Configuration:**

```yaml
# kong/services.yml
_format_version: "3.0"

services:
  - name: biometrics-api
    url: https://api.biometrics.dev
    routes:
      - name: api-routes
        paths:
          - /api
        strip_path: true
        preserve_host: true
    plugins:
      - name: rate-limiting
        config:
          minute: 100
          hour: 1000
          policy: local
          fault_tolerant: true
      - name: cors
        config:
          origins:
            - https://biometrics.dev
            - https://app.biometrics.dev
          methods:
            - GET
            - POST
            - PUT
            - DELETE
            - OPTIONS
          headers:
            - Authorization
            - Content-Type
            - X-Request-ID
          credentials: true
          max_age: 3600

  - name: biometrics-websocket
    url: https://ws.biometrics.dev
    routes:
      - name: ws-routes
        paths:
          - /ws
        strip_path: true
    plugins:
      - name: rate-limiting
        config:
          minute: 50
          hour: 500

upstreams:
  - name: biometrics-api
    targets:
      - target: 172.20.0.10:8000
        weight: 100
      - target: 172.20.0.11:8000
        weight: 100
    healthchecks:
      active:
        type: http
        http_path: /health
        interval: 10
        timeout: 5
        unhealthy: 3
        healthy: 2
```

**Kong Plugin Development:**

```typescript
// Custom Kong plugin: request-transformer
const requestTransformerPlugin = `
local constants = require "kong.constants"
local responses = require "kong.response"

local schema = {
  type = "object",
  properties = {
    add_headers = {
      type = "object",
      additionalProperties = true,
    },
    remove_headers = {
      type = "array",
      items = { type = "string" }
    }
  }
}

local function add_headers(conf, headers)
  for key, value in pairs(conf.add_headers or {}) do
    ngx.header[key] = value
  end
end

local function remove_headers(conf)
  for _, header in ipairs(conf.remove_headers or {}) do
    ngx.header[header] = nil
  end
end

return {
  handler = function(self)
    if ngx.var.request_method == "POST" or ngx.var.request_method == "PUT" then
      add_headers(self.config)
      remove_headers(self.config)
    end
  end
}
`;
```

### 4.2 Traefik Reverse Proxy

Traefik provides dynamic reverse proxy and load balancing capabilities for internal services.

**Traefik Configuration:**

```yaml
# traefik/dynamic-config.yml
http:
  routers:
    api-router:
      rule: "PathPrefix(\"/api\")"
      service: api-service
      entryPoints:
        - web
        - websecure
      middlewares:
        - strip-prefix
        - security-headers
      tls:
        certResolver: letsencrypt

    websocket-router:
      rule: "PathPrefix(\"/ws\")"
      service: websocket-service
      entryPoints:
        - websecure
      tls:
        certResolver: letsencrypt

  services:
    api-service:
      loadBalancer:
        servers:
          - url: "http://172.20.0.10:3000"
          - url: "http://172.20.0.11:3000"

    websocket-service:
      loadBalancer:
        servers:
          - url: "http://172.20.0.20:8080"

  middlewares:
    strip-prefix:
      stripPrefix:
        prefixes:
          - /api

    security-headers:
      headers:
        frameDeny: true
        browserXssFilter: true
        contentTypeNosniff: true
        sslRedirect: true
        stsSeconds: 31536000
        stsIncludeSubdomains: true
        forceSTSHeader: true

    rate-limiter:
      rateLimit:
        average: 100
        burst: 50
        period: 1m

    circuit-breaker:
      circuitBreaker:
        expression: "LatencyAtPercentile(50) > 1000"
        duration: 10s

entryPoints:
  web:
    address: ":80"
    http:
      redirections:
        entryPoint:
          to: websecure
          scheme: https

  websecure:
    address: ":443"
    http:
      tls:
        certResolver: letsencrypt
```

### 4.3 AWS API Gateway

**REST API Setup:**

```typescript
import {
  APIGatewayClient,
  CreateRestApiCommand,
  CreateResourceCommand,
  PutMethodCommand,
  PutIntegrationCommand,
  CreateDeploymentCommand,
} from '@aws-sdk/client-api-gateway';

const client = new APIGatewayClient({ region: 'eu-central-1' });

interface ApiGatewaySetup {
  name: string;
  description: string;
}

async function createApiGateway(setup: ApiGatewaySetup): Promise<string> {
  const api = await client.send(
    new CreateRestApiCommand({
      name: setup.name,
      description: setup.description,
      endpointConfiguration: {
        types: ['REGIONAL'],
      },
    })
  );
  
  return api.id!;
}

async function addResource(
  apiId: string,
  parentId: string,
  pathPart: string
): Promise<string> {
  const resource = await client.send(
    new CreateResourceCommand({
      restApiId: apiId,
      parentId,
      pathPart,
    })
  );
  
  return resource.id!;
}

async function addIntegration(
  apiId: string,
  resourceId: string,
  httpMethod: string,
  uri: string
): Promise<void> {
  await client.send(
    new PutIntegrationCommand({
      restApiId: apiId,
      resourceId,
      httpMethod,
      type: 'AWS_PROXY',
      integrationHttpMethod: 'POST',
      uri,
    })
  );
}
```

---

## 5. Webhook Architecture

### 5.1 Webhook Design Principles

Webhooks enable event-driven communication between BIOMETRICS and external services. Each webhook follows strict design principles for security, reliability, and debugging.

**Payload Structure:**

```typescript
interface WebhookPayload<T = unknown> {
  id: string;
  type: string;
  source: string;
  timestamp: string;
  version: string;
  data: T;
  metadata?: {
    correlationId?: string;
    userId?: string;
    sessionId?: string;
  };
  signature: string;
}

interface WebhookEnvelope {
  payload: WebhookPayload;
  deliveredAt: string;
  attempt: number;
  maxRetries: number;
}
```

**Webhook Contract Example:**

```typescript
const webhookContracts = {
  'payment.completed': {
    schema: {
      type: 'object',
      required: ['transactionId', 'amount', 'currency', 'timestamp'],
      properties: {
        transactionId: { type: 'string' },
        amount: { type: 'number' },
        currency: { type: 'string', enum: ['EUR', 'USD', 'GBP'] },
        timestamp: { type: 'string', format: 'date-time' },
        customerId: { type: 'string' },
        metadata: { type: 'object' },
      },
    },
  },
  'user.created': {
    schema: {
      type: 'object',
      required: ['userId', 'email', 'createdAt'],
      properties: {
        userId: { type: 'string', format: 'uuid' },
        email: { type: 'string', format: 'email' },
        name: { type: 'string' },
        createdAt: { type: 'string', format: 'date-time' },
        source: { type: 'string' },
      },
    },
  },
  'order.shipped': {
    schema: {
      type: 'object',
      required: ['orderId', 'trackingNumber', 'carrier'],
      properties: {
        orderId: { type: 'string' },
        trackingNumber: { type: 'string' },
        carrier: { type: 'string' },
        estimatedDelivery: { type: 'string', format: 'date-time' },
        items: {
          type: 'array',
          items: {
            type: 'object',
            properties: {
              productId: { type: 'string' },
              quantity: { type: 'number' },
            },
          },
        },
      },
    },
  },
};
```

### 5.2 Signature Verification

Webhook signature verification ensures that incoming requests genuinely originate from the expected provider and have not been tampered with.

**HMAC-SHA256 Verification:**

```typescript
import crypto from 'crypto';

interface SignatureVerificationOptions {
  secret: string;
  signatureHeader: string;
  payload: string;
  algorithm?: 'sha256' | 'sha384' | 'sha512';
  tolerance?: number;
}

interface VerificationResult {
  valid: boolean;
  error?: string;
  timestamp?: number;
}

function verifyWebhookSignature(
  options: SignatureVerificationOptions
): VerificationResult {
  const {
    secret,
    signatureHeader,
    payload,
    algorithm = 'sha256',
    tolerance = 300,
  } = options;
  
  try {
    const [timestamp, signature] = signatureHeader.split(',').reduce((acc, part) => {
      const [key, value] = part.split('=');
      acc[key] = value;
      return acc;
    }, {} as Record<string, string>);
    
    if (!timestamp || !signature) {
      return { valid: false, error: 'Invalid signature format' };
    }
    
    const timestampNum = parseInt(timestamp, 10);
    
    if (isNaN(timestampNum)) {
      return { valid: false, error: 'Invalid timestamp' };
    }
    
    const currentTime = Math.floor(Date.now() / 1000);
    
    if (currentTime - timestampNum > tolerance) {
      return { valid: false, error: 'Timestamp outside tolerance window' };
    }
    
    const signedPayload = `${timestamp}.${payload}`;
    
    const expectedSignature = crypto
      .createHmac(algorithm, secret)
      .update(signedPayload, 'utf8')
      .digest('hex');
    
    const signatureBuffer = Buffer.from(signature, 'hex');
    const expectedBuffer = Buffer.from(expectedSignature, 'hex');
    
    if (!crypto.timingSafeEqual(signatureBuffer, expectedBuffer)) {
      return { valid: false, error: 'Signature mismatch' };
    }
    
    return { valid: true, timestamp: timestampNum };
  } catch (error) {
    return {
      valid: false,
      error: error instanceof Error ? error.message : 'Verification failed',
    };
  }
}
```

**Express.js Webhook Handler:**

```typescript
import express from 'express';
import crypto from 'crypto';

const app = express();

app.post(
  '/webhooks/:provider',
  express.raw({ type: 'application/json' }),
  async (req, res) => {
    const provider = req.params.provider;
    const signature = req.headers['x-webhook-signature'] as string;
    
    const secret = getWebhookSecret(provider);
    
    const result = verifyWebhookSignature({
      secret,
      signatureHeader: signature,
      payload: req.body.toString(),
    });
    
    if (!result.valid) {
      console.error('Webhook verification failed:', result.error);
      return res.status(401).json({ error: 'Invalid signature' });
    }
    
    try {
      const payload = JSON.parse(req.body.toString());
      
      await processWebhook(provider, payload);
      
      res.status(200).json({ received: true });
    } catch (error) {
      console.error('Webhook processing error:', error);
      res.status(500).json({ error: 'Processing failed' });
    }
  }
);

function getWebhookSecret(provider: string): string {
  const secrets: Record<string, string> = {
    stripe: process.env.STRIPE_WEBHOOK_SECRET!,
    paypal: process.env.PAYPAL_WEBHOOK_SECRET!,
    shopify: process.env.SHOPIFY_WEBHOOK_SECRET!,
  };
  
  return secrets[provider] ?? '';
}
```

### 5.3 Retry Logic and Idempotency

Webhooks must handle delivery failures gracefully and ensure idempotent processing to prevent duplicate actions.

**Retry Strategy Implementation:**

```typescript
interface RetryConfig {
  maxAttempts: number;
  baseDelayMs: number;
  maxDelayMs: number;
  backoffMultiplier: number;
  retryableStatuses: number[];
}

const defaultRetryConfig: RetryConfig = {
  maxAttempts: 5,
  baseDelayMs: 1000,
  maxDelayMs: 60000,
  backoffMultiplier: 2,
  retryableStatuses: [408, 429, 500, 502, 503, 504],
};

class WebhookRetryHandler {
  private config: RetryConfig;
  
  constructor(config: Partial<RetryConfig> = {}) {
    this.config = { ...defaultRetryConfig, ...config };
  }
  
  async executeWithRetry<T>(
    fn: () => Promise<T>
  ): Promise<T> {
    let lastError: Error | undefined;
    
    for (let attempt = 1; attempt <= this.config.maxAttempts; attempt++) {
      try {
        return await fn();
      } catch (error) {
        lastError = error as Error;
        
        if (!this.isRetryable(error)) {
          throw error;
        }
        
        if (attempt < this.config.maxAttempts) {
          const delay = this.calculateDelay(attempt);
          await this.sleep(delay);
        }
      }
    }
    
    throw lastError;
  }
  
  private isRetryable(error: unknown): boolean {
    if (error instanceof Response) {
      return this.config.retryableStatuses.includes(error.status);
    }
    
    return true;
  }
  
  private calculateDelay(attempt: number): number {
    const delay = this.config.baseDelayMs * Math.pow(this.config.backoffMultiplier, attempt - 1);
    return Math.min(delay, this.config.maxDelayMs);
  }
  
  private sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
}
```

**Idempotency Implementation:**

```typescript
import Redis from 'ioredis';

const redis = new Redis(process.env.REDIS_URL);

interface IdempotencyOptions {
  ttl: number;
}

class IdempotencyService {
  private redis: Redis;
  private defaultTtl = 86400;
  
  constructor(redis: Redis) {
    this.redis = redis;
  }
  
  async check<T>(
    key: string,
    fn: () => Promise<T>,
    options?: IdempotencyOptions
  ): Promise<T> {
    const cached = await this.redis.get(`idempotent:${key}`);
    
    if (cached) {
      const parsed = JSON.parse(cached);
      
      if (parsed.status === 'success') {
        return parsed.result;
      }
      
      if (parsed.status === 'processing') {
        return this.waitForCompletion<T>(key);
      }
      
      if (parsed.status === 'failed') {
        throw new Error(parsed.error);
      }
    }
    
    await this.redis.setex(
      `idempotent:${key}`,
      options?.ttl ?? this.defaultTtl,
      JSON.stringify({ status: 'processing' })
    );
    
    try {
      const result = await fn();
      
      await this.redis.setex(
        `idempotent:${key}`,
        options?.ttl ?? this.defaultTtl,
        JSON.stringify({ status: 'success', result })
      );
      
      return result;
    } catch (error) {
      await this.redis.setex(
        `idempotent:${key}`,
        options?.ttl ?? this.defaultTtl,
        JSON.stringify({
          status: 'failed',
          error: error instanceof Error ? error.message : 'Unknown error',
        })
      );
      
      throw error;
    }
  }
  
  private async waitForCompletion<T>(key: string): Promise<T> {
    for (let i = 0; i < 30; i++) {
      await this.sleep(1000);
      
      const cached = await this.redis.get(`idempotent:${key}`);
      
      if (cached) {
        const parsed = JSON.parse(cached);
        
        if (parsed.status === 'success') {
          return parsed.result;
        }
        
        if (parsed.status === 'failed') {
          throw new Error(parsed.error);
        }
      }
    }
    
    throw new Error('Idempotency key timeout');
  }
  
  private sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
}
```

---

## 6. Message Queue Integration

### 6.1 Apache Kafka Integration

Kafka provides high-throughput, fault-tolerant message streaming for BIOMETRICS.

**Producer Implementation:**

```typescript
import { Kafka, Producer, CompressionTypes } from 'kafkajs';

const kafka = new Kafka({
  clientId: 'biometrics-producer',
  brokers: process.env.KAFKA_BROKERS!.split(','),
  ssl: true,
  sasl: {
    mechanism: 'scram-sha512',
    username: process.env.KAFKA_USERNAME!,
    password: process.env.KAFKA_PASSWORD!,
  },
});

const producer = kafka.producer({
  allowAutoTopicCreation: false,
  transactionTimeout: 30000,
});

interface KafkaMessage<T> {
  key?: string;
  value: T;
  headers?: Record<string, string>;
  timestamp?: string;
  partition?: number;
}

class BiometricsProducer {
  private producer: Producer;
  
  async connect(): Promise<void> {
    await this.producer.connect();
  }
  
  async disconnect(): Promise<void> {
    await this.producer.disconnect();
  }
  
  async send<T>(
    topic: string,
    messages: KafkaMessage<T>[]
  ): Promise<void> {
    const kafkaMessages = messages.map(msg => ({
      key: msg.key,
      value: JSON.stringify(msg.value),
      headers: msg.headers,
      timestamp: msg.timestamp ?? Date.now().toString(),
    }));
    
    await this.producer.send({
      topic,
      compression: CompressionTypes.GZIP,
      messages: kafkaMessages,
    });
  }
  
  async sendTransaction<T>(
    topic: string,
    messages: KafkaMessage<T>[]
  ): Promise<void> {
    await this.producer.transaction(async () => {
      await this.send(topic, messages);
    });
  }
}
```

**Consumer Implementation:**

```typescript
import { Kafka, Consumer, EachMessagePayload } from 'kafkajs';

const consumer = kafka.consumer({
  groupId: 'biometrics-consumer-group',
  sessionTimeout: 30000,
  heartbeatInterval: 3000,
});

interface MessageHandler<T> {
  (message: T): Promise<void>;
}

class BiometricsConsumer {
  private consumer: Consumer;
  private handlers = new Map<string, MessageHandler<unknown>>();
  
  async connect(): Promise<void> {
    await this.consumer.connect();
  }
  
  async subscribe(
    topic: string,
    handler: MessageHandler<unknown>,
    fromBeginning = false
  ): Promise<void> {
    await this.consumer.subscribe({ topic, fromBeginning });
    this.handlers.set(topic, handler);
  }
  
  async start(): Promise<void> {
    await this.consumer.run({
      eachMessage: async (payload: EachMessagePayload) => {
        const handler = this.handlers.get(payload.topic);
        
        if (!handler) {
          console.warn(`No handler for topic: ${payload.topic}`);
          return;
        }
        
        try {
          const value = JSON.parse(payload.message.value?.toString() ?? '{}');
          await handler(value);
        } catch (error) {
          console.error(`Error processing message:`, error);
        }
      },
    });
  }
  
  async pause(topics: string[]): Promise<void> {
    await this.consumer.pause(topics.map(topic => ({ topic })));
  }
  
  async resume(topics: string[]): Promise<void> {
    await this.consumer.resume(topics.map(topic => ({ topic })));
  }
}
```

### 6.2 RabbitMQ Integration

RabbitMQ provides flexible messaging patterns including work queues, pub/sub, and routing.

**Exchange and Queue Setup:**

```typescript
import amqp, { Connection, Channel } from 'amqplib';

class RabbitMQClient {
  private connection: Connection | null = null;
  private channel: Channel | null = null;
  
  async connect(url: string): Promise<void> {
    this.connection = await amqp.connect(url);
    this.channel = await this.connection.createChannel();
  }
  
  async setupExchanges(): Promise<void> {
    const channel = this.channel!;
    
    await channel.assertExchange('biometrics.events', 'topic', {
      durable: true,
    });
    
    await channel.assertExchange('biometrics.notifications', 'fanout', {
      durable: true,
    });
    
    await channel.assertExchange('biometrics.direct', 'direct', {
      durable: true,
    });
  }
  
  async setupQueues(): Promise<void> {
    const channel = this.channel!;
    
    await channel.assertQueue('biometrics.orders', {
      durable: true,
      arguments: {
        'x-dead-letter-exchange': 'biometrics.dlx',
        'x-message-ttl': 86400000,
      },
    });
    
    await channel.assertQueue('biometrics.notifications', {
      durable: true,
    });
    
    await channel.assertQueue('biometrics.payments', {
      durable: true,
    });
  }
  
  async bindQueues(): Promise<void> {
    const channel = this.channel!;
    
    await channel.bindQueue(
      'biometrics.orders',
      'biometrics.events',
      'order.*'
    );
    
    await channel.bindQueue(
      'biometrics.notifications',
      'biometrics.notifications',
      ''
    );
    
    await channel.bindQueue(
      'biometrics.payments',
      'biometrics.direct',
      'payment'
    );
  }
  
  async publish<T>(
    exchange: string,
    routingKey: string,
    message: T
  ): Promise<void> {
    const channel = this.channel!;
    
    channel.publish(
      exchange,
      routingKey,
      Buffer.from(JSON.stringify(message)),
      { persistent: true }
    );
  }
  
  async consume<T>(
    queue: string,
    handler: (message: T) => Promise<void>
  ): Promise<void> {
    const channel = this.channel!;
    
    await channel.consume(queue, async (msg) => {
      if (!msg) return;
      
      try {
        const content = JSON.parse(msg.content.toString()) as T;
        await handler(content);
        channel.ack(msg);
      } catch (error) {
        console.error('Error processing message:', error);
        channel.nack(msg, false, false);
      }
    });
  }
}
```

### 6.3 Redis Pub/Sub

Redis Pub/Sub provides low-latency messaging for real-time features.

```typescript
import Redis from 'ioredis';

const redis = new Redis(process.env.REDIS_URL);
const subscriber = new Redis(process.env.REDIS_URL);

interface PubSubMessage<T> {
  channel: string;
  data: T;
}

type MessageListener<T> = (message: PubSubMessage<T>) => void;

class RedisPubSub {
  private subscriptions = new Map<string, Set<MessageListener<unknown>>>();
  
  async publish<T>(channel: string, data: T): Promise<number> {
    return redis.publish(channel, JSON.stringify(data));
  }
  
  async subscribe<T>(
    channel: string,
    listener: MessageListener<T>
  ): Promise<void> {
    const channelListeners = this.subscriptions.get(channel) ?? new Set();
    channelListeners.add(listener as MessageListener<unknown>);
    this.subscriptions.set(channel, channelListeners);
    
    if (channelListeners.size === 1) {
      await subscriber.subscribe(channel);
    }
  }
  
  async unsubscribe<T>(
    channel: string,
    listener: MessageListener<T>
  ): Promise<void> {
    const channelListeners = this.subscriptions.get(channel);
    
    if (channelListeners) {
      channelListeners.delete(listener as MessageListener<unknown>);
      
      if (channelListeners.size === 0) {
        await subscriber.unsubscribe(channel);
        this.subscriptions.delete(channel);
      }
    }
  }
  
  startListening(): void {
    subscriber.on('message', (channel, message) => {
      const listeners = this.subscriptions.get(channel);
      
      if (listeners) {
        const data = JSON.parse(message);
        
        listeners.forEach(listener => {
          listener({ channel, data });
        });
      }
    });
  }
}
```

---

## 7. Third-Party Services

### 7.1 CRM Integrations

**HubSpot Integration:**

```typescript
import { Client } from '@hubspot/api-client';

const hubspotClient = new Client({ accessToken: process.env.HUBSPOT_ACCESS_TOKEN });

interface ContactData {
  email: string;
  firstName?: string;
  lastName?: string;
  company?: string;
  phone?: string;
}

async function createOrUpdateContact(data: ContactData): Promise<string> {
  const existingContact = await hubspotClient.crm.contacts.basicApi.getPage(
    undefined,
    undefined,
    `email=${data.email}`
  );
  
  if (existingContact.results.length > 0) {
    const contactId = existingContact.results[0].id;
    
    await hubspotClient.crm.contacts.basicApi.update(contactId, {
      properties: {
        firstname: data.firstName,
        lastname: data.lastName,
        company: data.company,
        phone: data.phone,
      },
    });
    
    return contactId;
  }
  
  const contact = await hubspotClient.crm.contacts.basicApi.create({
    properties: {
      email: data.email,
      firstname: data.firstName,
      lastname: data.lastName,
      company: data.company,
      phone: data.phone,
    },
  });
  
  return contact.id;
}
```

**Salesforce Integration:**

```typescript
import jsforce from 'jsforce';

class SalesforceClient {
  private conn: jsforce.Connection;
  
  async connect(): Promise<void> {
    this.conn = new jsforce.Connection({
      loginUrl: 'https://login.salesforce.com',
    });
    
    await this.conn.login(
      process.env.SALESFORCE_USERNAME!,
      process.env.SALESFORCE_PASSWORD! + process.env.SALESFORCE_TOKEN!
    );
  }
  
  async createLead(data: {
    email: string;
    firstName: string;
    lastName: string;
    company: string;
  }): Promise<string> {
    const result = await this.conn.sobject('Lead').create({
      Email: data.email,
      FirstName: data.firstName,
      LastName: data.lastName,
      Company: data.company,
      LeadSource: 'Web',
      Status: 'Open - Not Contacted',
    });
    
    if (!result.success) {
      throw new Error(`Failed to create lead: ${result.errors}`);
    }
    
    return result.id;
  }
  
  async updateOpportunity(
    opportunityId: string,
    data: {
      stageName: string;
      amount?: number;
      closeDate?: Date;
    }
  ): Promise<void> {
    await this.conn.sobject('Opportunity').update({
      Id: opportunityId,
      StageName: data.stageName,
      Amount: data.amount,
      CloseDate: data.closeDate?.toISOString().split('T')[0],
    });
  }
}
```

### 7.2 Analytics Integration

**Google Analytics 4:**

```typescript
import { BetaAnalyticsDataClient } from '@google-analytics/data';

const analyticsDataClient = new BetaAnalyticsDataClient();

interface AnalyticsEvent {
  name: string;
  params?: Record<string, string | number | boolean>;
}

interface UserProperties {
  userId: string;
  properties: Record<string, string>;
}

async function trackEvent(
  userId: string,
  sessionId: string,
  events: AnalyticsEvent[]
): Promise<void> {
  const [response] = await analyticsDataClient.runReport({
    property: `properties/${process.env.GA4_PROPERTY_ID}`,
    events: events.map(event => ({
      name: event.name,
      params: event.params,
    })),
    userTraits: {
      userId: { value: userId },
    },
  });
}

async function getAnalyticsReport(params: {
  startDate: string;
  endDate: string;
  metrics: string[];
  dimensions?: string[];
}) {
  const [response] = await analyticsDataClient.runReport({
    property: `properties/${process.env.GA4_PROPERTY_ID}`,
    dateRanges: [
      { startDate: params.startDate, endDate: params.endDate },
    ],
    dimensions: params.dimensions?.map(d => ({ name: d })),
    metrics: params.metrics.map(m => ({ name: m })),
  });
  
  return response;
}
```

**Mixpanel Integration:**

```typescript
import { Mixpanel } from 'mixpanel';

const mixpanel = new Mixpanel(process.env.MIXPANEL_TOKEN!);

interface TrackEventParams {
  event: string;
  userId: string;
  properties?: Record<string, unknown>;
  time?: number;
}

async function trackEvent(params: TrackEventParams): Promise<void> {
  mixpanel.track(params.event, {
    distinct_id: params.userId,
    ...params.properties,
    time: params.time ?? Math.floor(Date.now() / 1000),
  });
}

async function identifyUser(
  userId: string,
  properties: Record<string, unknown>
): Promise<void> {
  mixpanel.people.set(userId, properties);
}

async function updateUserProfile(
  userId: string,
  properties: Record<string, unknown>
): Promise<void> {
  mixpanel.people.increment(userId, properties as Record<string, number>);
}
```

### 7.3 Monitoring Integration

**Datadog Integration:**

```typescript
import { StatsD } from 'hot-shots';

const dogstatsd = new StatsD({
  host: process.env.DATADOG_HOST ?? 'localhost',
  port: parseInt(process.env.DATADOG_PORT ?? '8125'),
  prefix: 'biometrics.',
});

function increment(metric: string, value = 1, tags?: Record<string, string>): void {
  dogstatsd.increment(metric, value, convertTags(tags));
}

function gauge(metric: string, value: number, tags?: Record<string, string>): void {
  dogstatsd.gauge(metric, value, convertTags(tags));
}

function timing(metric: string, value: number, tags?: Record<string, string>): void {
  dogstatsd.timing(metric, value, convertTags(tags));
}

function histogram(metric: string, value: number, tags?: Record<string, string>): void {
  dogstatsd.histogram(metric, value, convertTags(tags));
}

function convertTags(tags?: Record<string, string>): string[] | undefined {
  if (!tags) return undefined;
  
  return Object.entries(tags).map(([key, value]) => `${key}:${value}`);
}

function trackRequest(path: string, statusCode: number, durationMs: number): void {
  timing('http.request.duration', durationMs, {
    path,
    status: statusCode.toString(),
  });
  
  increment('http.request.count', 1, {
    path,
    status: statusCode.toString(),
  });
}
```

**Sentry Integration:**

```typescript
import * as Sentry from '@sentry/node';

Sentry.init({
  dsn: process.env.SENTRY_DSN,
  environment: process.env.NODE_ENV,
  release: process.env.SENTRY_RELEASE,
  tracesSampleRate: 0.1,
  integrations: [
    new Sentry.Integrations.Http({ tracing: true }),
    new Sentry.Integrations.Express(),
  ],
});

function captureException(error: Error, context?: Record<string, unknown>): void {
  Sentry.captureException(error, {
    extra: context,
  });
}

function captureMessage(
  message: string,
  level: Sentry.Severity = Sentry.Severity.Info,
  context?: Record<string, unknown>
): void {
  Sentry.captureMessage(message, level, {
    extra: context,
  });
}

function setUser(userId: string, email?: string, ip?: string): void {
  Sentry.setUser({
    id: userId,
    email,
    ip_address: ip,
  });
}
```

---

## 8. API Documentation

### 8.1 OpenAPI Specification

**API Specification Structure:**

```yaml
openapi: 3.1.0
info:
  title: BIOMETRICS API
  description: |
    Comprehensive API for BIOMETRICS platform.
    
    ## Authentication
    All endpoints require Bearer token authentication.
    
    ## Rate Limiting
    - Standard: 100 requests/minute
    - Authenticated: 1000 requests/minute
  version: 2.0.0
  contact:
    name: BIOMETRICS Support
    email: support@biometrics.dev
  license:
    name: Proprietary

servers:
  - url: https://api.biometrics.dev
    description: Production server
  - url: https://staging-api.biometrics.dev
    description: Staging server

security:
  - bearerAuth: []

tags:
  - name: Users
    description: User management operations
  - name: Projects
    description: Project CRUD operations
  - name: Analytics
    description: Analytics and reporting

paths:
  /api/v2/users:
    get:
      summary: List all users
      description: Returns a paginated list of users
      tags:
        - Users
      parameters:
        - $ref: '#/components/parameters/Page'
        - $ref: '#/components/parameters/Limit'
        - name: sort
          in: query
          schema:
            type: string
            enum: [created_at, email, name]
          description: Sort field
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UsersList'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '429':
          $ref: '#/components/responses/TooManyRequests'

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

  parameters:
    Page:
      name: page
      in: query
      schema:
        type: integer
        default: 1
      description: Page number
    Limit:
      name: limit
      in: query
      schema:
        type: integer
        default: 20
        maximum: 100

  schemas:
    UsersList:
      type: object
      properties:
        data:
          type: array
          items:
            $ref: '#/components/schemas/User'
        pagination:
          $ref: '#/components/schemas/Pagination'

    User:
      type: object
      required:
        - id
        - email
      properties:
        id:
          type: string
          format: uuid
        email:
          type: string
          format: email
        name:
          type: string
        createdAt:
          type: string
          format: date-time

    Pagination:
      type: object
      properties:
        page:
          type: integer
        limit:
          type: integer
        total:
          type: integer
        totalPages:
          type: integer

  responses:
    Unauthorized:
      description: Unauthorized
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
    TooManyRequests:
      description: Too many requests
      headers:
        Retry-After:
          schema:
            type: integer
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
```

### 8.2 GraphQL Schema

**Schema Definition:**

```graphql
scalar DateTime
scalar JSON
scalar Upload

type Query {
  user(id: ID!): User
  users(filter: UserFilter, pagination: Pagination): UserConnection!
  project(id: ID!): Project
  projects(filter: ProjectFilter, pagination: Pagination): ProjectConnection!
  analytics(range: DateRange!, metrics: [String!]!): AnalyticsResult!
}

type Mutation {
  createUser(input: CreateUserInput!): User!
  updateUser(id: ID!, input: UpdateUserInput!): User!
  deleteUser(id: ID!): Boolean!
  
  createProject(input: CreateProjectInput!): Project!
  updateProject(id: ID!, input: UpdateProjectInput!): Project!
  deleteProject(id: ID!): Boolean!
}

type Subscription {
  projectUpdated(id: ID!): Project!
  notificationCreated(userId: ID!): Notification!
}

type User {
  id: ID!
  email: String!
  name: String
  avatar: String
  createdAt: DateTime!
  updatedAt: DateTime!
  projects: [Project!]!
}

type Project {
  id: ID!
  name: String!
  description: String
  status: ProjectStatus!
  owner: User!
  createdAt: DateTime!
  updatedAt: DateTime!
}

enum ProjectStatus {
  ACTIVE
  ARCHIVED
  DELETED
}

input UserFilter {
  search: String
  createdAfter: DateTime
  createdBefore: DateTime
}

input CreateUserInput {
  email: String!
  name: String
  avatar: Upload
}

input ProjectFilter {
  status: ProjectStatus
  ownerId: ID
}

input Pagination {
  page: Int = 1
  limit: Int = 20
}

type UserConnection {
  edges: [UserEdge!]!
  pageInfo: PageInfo!
  totalCount: Int!
}

type UserEdge {
  node: User!
  cursor: String!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}
```

---

## 9. Integration Testing

### 9.1 Unit Testing

**Mock Strategies:**

```typescript
interface MockOptions<T> {
  defaultReturn?: T;
  delay?: number;
  shouldReject?: boolean;
  errorMessage?: string;
}

function createMockFunction<T>(
  implementation: (...args: unknown[]) => Promise<T>,
  options: MockOptions<T> = {}
): jest.Mock<Promise<T>, unknown[]> {
  return jest.fn(async (...args: unknown[]) => {
    if (options.delay) {
      await new Promise(resolve => setTimeout(resolve, options.delay));
    }
    
    if (options.shouldReject) {
      throw new Error(options.errorMessage ?? 'Mock error');
    }
    
    return implementation(...args);
  });
}

function createMockApiClient(): {
  get: jest.Mock;
  post: jest.Mock;
  put: jest.Mock;
  delete: jest.Mock;
} {
  return {
    get: jest.fn(),
    post: jest.fn(),
    put: jest.fn(),
    delete: jest.fn(),
  };
}
```

**Contract Testing:**

```typescript
import { validate } from 'jsonschema';

interface ContractTest<T> {
  schema: Record<string, unknown>;
  data: T;
  description: string;
}

function runContractTests<T>(tests: ContractTest<T>[]): void {
  describe('Contract Tests', () => {
    tests.forEach(({ schema, data, description }) => {
      it(`should validate: ${description}`, () => {
        const result = validate(data, schema);
        
        if (!result.valid) {
          console.error('Validation errors:', result.errors);
        }
        
        expect(result.valid).toBe(true);
      });
    });
  });
}
```

### 9.2 Integration Testing

**Test Environment Setup:**

```typescript
import { Testcontainers } from 'testcontainers';
import { Kafka } from 'kafkajs';

interface TestEnvironment {
  postgres: Testcontainers.Container;
  redis: Testcontainers.Container;
  kafka: Testcontainers.Container;
  elasticsearch: Testcontainers.Container;
}

async function setupTestEnvironment(): Promise<TestEnvironment> {
  const postgres = await new Testcontainers.GenericContainer('postgres:15')
    .withEnv('POSTGRES_USER', 'test')
    .withEnv('POSTGRES_PASSWORD', 'test')
    .withEnv('POSTGRES_DB', 'test')
    .withExposedPorts(5432)
    .start();
  
  const redis = await new Testcontainers.GenericContainer('redis:7')
    .withExposedPorts(6379)
    .start();
  
  const kafka = await new Testcontainers.GenericContainer('confluentinc/cp-kafka:7.5.0')
    .withEnv('KAFKA_ZOOKEEPER_CONNECT', 'zookeeper:2181')
    .withEnv('KAFKA_ADVERTISED_LISTENERS', 'PLAINTEXT://localhost:9092')
    .withExposedPorts(9092)
    .start();
  
  return { postgres, redis, kafka };
}

async function teardownTestEnvironment(env: TestEnvironment): Promise<void> {
  await env.postgres.stop();
  await env.redis.stop();
  await env.kafka.stop();
}
```

**Data Seeding:**

```typescript
interface SeedData<T> {
  table: string;
  data: T[];
}

async function seedDatabase(seeds: SeedData<unknown>[]): Promise<void> {
  for (const seed of seeds) {
    for (const record of seed.data) {
      await prisma[seed.table].create({ data: record });
    }
  }
}

async function clearDatabase(tables: string[]): Promise<void> {
  for (const table of tables) {
    await prisma.$executeRawUnsafe(`TRUNCATE TABLE "${table}" CASCADE`);
  }
}
```

### 9.3 End-to-End Testing

**E2E Test Framework:**

```typescript
import { test, expect, Page } from '@playwright/test';

test.describe('API Integration E2E', () => {
  let apiClient: ApiClient;
  
  test.beforeAll(() => {
    apiClient = new ApiClient({
      baseUrl: process.env.API_URL!,
      apiKey: process.env.TEST_API_KEY!,
    });
  });
  
  test('complete user workflow', async () => {
    const user = await apiClient.createUser({
      email: 'test@example.com',
      name: 'Test User',
    });
    
    expect(user.id).toBeDefined();
    expect(user.email).toBe('test@example.com');
    
    const project = await apiClient.createProject({
      name: 'Test Project',
      ownerId: user.id,
    });
    
    expect(project.id).toBeDefined();
    expect(project.ownerId).toBe(user.id);
    
    const updatedProject = await apiClient.updateProject(project.id, {
      status: 'ACTIVE',
    });
    
    expect(updatedProject.status).toBe('ACTIVE');
    
    await apiClient.deleteProject(project.id);
    await apiClient.deleteUser(user.id);
  });
  
  test('webhook delivery', async ({ request }) => {
    const webhook = await request.post('/api/webhooks', {
      data: {
        url: 'https://example.com/webhook',
        events: ['user.created', 'project.created'],
      },
    });
    
    expect(webhook.ok()).toBeTruthy();
  });
});
```

---

## 10. Security and Compliance

### 10.1 API Security

**Authentication Implementation:**

```typescript
import jwt from 'jsonwebtoken';
import bcrypt from 'bcrypt';

interface TokenPayload {
  userId: string;
  email: string;
  role: string;
  permissions: string[];
  iat: number;
  exp: number;
}

function generateAccessToken(payload: Omit<TokenPayload, 'iat' | 'exp'>): string {
  return jwt.sign(payload, process.env.JWT_SECRET!, {
    expiresIn: '15m',
    issuer: 'biometrics',
  });
}

function generateRefreshToken(payload: Omit<TokenPayload, 'iat' | 'exp'>): string {
  return jwt.sign(payload, process.env.JWT_REFRESH_SECRET!, {
    expiresIn: '7d',
    issuer: 'biometrics',
  });
}

function verifyToken(token: string): TokenPayload {
  return jwt.verify(token, process.env.JWT_SECRET!) as TokenPayload;
}

async function hashPassword(password: string): Promise<string> {
  return bcrypt.hash(password, 12);
}

async function verifyPassword(password: string, hash: string): Promise<boolean> {
  return bcrypt.compare(password, hash);
}
```

**Authorization Middleware:**

```typescript
type Permission = 'read' | 'write' | 'delete' | 'admin';

interface AuthorizationConfig {
  resource: string;
  permission: Permission;
}

function authorize(...configs: AuthorizationConfig[]) {
  return (req: Request, res: Response, next: NextFunction) => {
    const user = req.user as TokenPayload;
    
    if (!user) {
      return res.status(401).json({ error: 'Unauthorized' });
    }
    
    const hasPermission = configs.some(config => {
      const requiredPermission = `${config.resource}:${config.permission}`;
      return user.permissions.includes(requiredPermission) ||
             user.permissions.includes('admin');
    });
    
    if (!hasPermission) {
      return res.status(403).json({ error: 'Forbidden' });
    }
    
    next();
  };
}
```

### 10.2 Data Privacy

**GDPR Compliance:**

```typescript
interface DataSubjectRequest {
  userId: string;
  requestType: 'access' | 'rectification' | 'erasure' | 'portability';
  requestedAt: Date;
}

class GDPRComplianceService {
  async handleDataSubjectRequest(request: DataSubjectRequest): Promise<void> {
    switch (request.requestType) {
      case 'access':
        return this.exportUserData(request.userId);
      case 'erasure':
        return this.deleteUserData(request.userId);
      case 'portability':
        return this.exportPortableData(request.userId);
      case 'rectification':
        return this.rectifyUserData(request.userId);
    }
  }
  
  async exportUserData(userId: string): Promise<Record<string, unknown>> {
    const user = await prisma.user.findUnique({
      where: { id: userId },
      include: {
        projects: true,
        sessions: true,
      },
    });
    
    return {
      exportedAt: new Date().toISOString(),
      user,
    };
  }
  
  async deleteUserData(userId: string): Promise<void> {
    await prisma.$transaction([
      prisma.session.deleteMany({ where: { userId } }),
      prisma.project.deleteMany({ where: { ownerId: userId } }),
      prisma.user.delete({ where: { id: userId } }),
    ]);
  }
  
  async exportPortableData(userId: string): Promise<string> {
    const data = await this.exportUserData(userId);
    return JSON.stringify(data, null, 2);
  }
}
```

---

## 11. Monitoring and Observability

### 11.1 API Monitoring

**Health Check Implementation:**

```typescript
interface HealthCheck {
  name: string;
  status: 'healthy' | 'degraded' | 'unhealthy';
  checks: ComponentHealth[];
  timestamp: string;
}

interface ComponentHealth {
  name: string;
  status: 'healthy' | 'degraded' | 'unhealthy';
  latencyMs?: number;
  error?: string;
}

async function performHealthCheck(): Promise<HealthCheck> {
  const checks: ComponentHealth[] = [];
  
  checks.push(await checkDatabase());
  checks.push(await checkRedis());
  checks.push(await checkKafka());
  checks.push(await checkExternalServices());
  
  const unhealthyCount = checks.filter(c => c.status === 'unhealthy').length;
  const degradedCount = checks.filter(c => c.status === 'degraded').length;
  
  let overallStatus: 'healthy' | 'degraded' | 'unhealthy' = 'healthy';
  if (unhealthyCount > 0) overallStatus = 'unhealthy';
  else if (degradedCount > 0) overallStatus = 'degraded';
  
  return {
    name: 'biometrics-api',
    status: overallStatus,
    checks,
    timestamp: new Date().toISOString(),
  };
}

async function checkDatabase(): Promise<ComponentHealth> {
  const start = Date.now();
  
  try {
    await prisma.$queryRaw`SELECT 1`;
    
    return {
      name: 'database',
      status: 'healthy',
      latencyMs: Date.now() - start,
    };
  } catch (error) {
    return {
      name: 'database',
      status: 'unhealthy',
      error: error instanceof Error ? error.message : 'Unknown error',
    };
  }
}
```

### 11.2 Integration Monitoring

**Metrics Collection:**

```typescript
interface IntegrationMetrics {
  integration: string;
  metrics: {
    requestCount: number;
    successCount: number;
    errorCount: number;
    averageLatencyMs: number;
    p95LatencyMs: number;
    p99LatencyMs: number;
  };
}

class MetricsCollector {
  private metrics = new Map<string, IntegrationMetrics>();
  
  recordRequest(integration: string, success: boolean, latencyMs: number): void {
    const metrics = this.metrics.get(integration) ?? this.initMetrics(integration);
    
    metrics.metrics.requestCount++;
    
    if (success) {
      metrics.metrics.successCount++;
    } else {
      metrics.metrics.errorCount++;
    }
    
    this.updateLatencyMetrics(integration, latencyMs);
  }
  
  getMetrics(integration: string): IntegrationMetrics | undefined {
    return this.metrics.get(integration);
  }
  
  getAllMetrics(): IntegrationMetrics[] {
    return Array.from(this.metrics.values());
  }
  
  private initMetrics(integration: string): IntegrationMetrics {
    const metrics: IntegrationMetrics = {
      integration,
      metrics: {
        requestCount: 0,
        successCount: 0,
        errorCount: 0,
        averageLatencyMs: 0,
        p95LatencyMs: 0,
        p99LatencyMs: 0,
      },
    };
    
    this.metrics.set(integration, metrics);
    
    return metrics;
  }
  
  private updateLatencyMetrics(integration: string, latencyMs: number): void {
    const metrics = this.metrics.get(integration)!;
    
    const latencies = this.getLatencies(integration);
    latencies.push(latencyMs);
    
    if (latencies.length > 1000) {
      latencies.shift();
    }
    
    const sum = latencies.reduce((a, b) => a + b, 0);
    metrics.metrics.averageLatencyMs = sum / latencies.length;
    
    const sorted = [...latencies].sort((a, b) => a - b);
    metrics.metrics.p95LatencyMs = sorted[Math.floor(sorted.length * 0.95)];
    metrics.metrics.p99LatencyMs = sorted[Math.floor(sorted.length * 0.99)];
  }
  
  private getLatencies(integration: string): number[] {
    return [];
  }
}
```

---

## 12. Best Practices 2026

### 12.1 API Design Guidelines

**Core Principles:**

The following principles guide all API design decisions in BIOMETRICS. First, consistency across endpoints ensures predictable developer experience. Second, backward compatibility prevents breaking changes for existing consumers. Third, clear error messages enable quick issue resolution. Fourth, comprehensive documentation reduces support burden. Fifth, security by default protects user data.

**Versioning Strategy:**

```typescript
interface VersionConfig {
  current: string;
  supported: string[];
  deprecated: string[];
  sunset: Map<string, Date>;
}

const versionConfig: VersionConfig = {
  current: 'v2',
  supported: ['v1', 'v2'],
  deprecated: [],
  sunset: new Map([
    ['v1', new Date('2026-12-31')],
  ]),
};

function getVersionFromRequest(req: Request): string {
  const accept = req.headers.get('Accept');
  
  if (accept?.includes('application/vnd.biometrics.v2+json')) {
    return 'v2';
  }
  
  return 'v1';
}

function checkDeprecationHeader(version: string): Record<string, string> {
  const headers: Record<string, string> = {};
  const sunsetDate = versionConfig.sunset.get(version);
  
  if (sunsetDate) {
    headers['Sunset'] = sunsetDate.toUTCString();
    headers['Deprecation'] = 'true';
  }
  
  return headers;
}
```

### 12.2 Deprecation Policy

**Deprecation Process:**

```typescript
interface DeprecationNotice {
  version: string;
  feature: string;
  deprecatedAt: Date;
  sunsetDate: Date;
  migrationGuide: string;
  alternative: string;
}

const deprecations: DeprecationNotice[] = [
  {
    version: 'v1',
    feature: '/api/v1/users/search',
    deprecatedAt: new Date('2026-01-01'),
    sunsetDate: new Date('2026-12-31'),
    migrationGuide: '/docs/migrations/v1-to-v2',
    alternative: 'POST /api/v2/users/search with query parameter support',
  },
];

function getDeprecationHeaders(endpoint: string): Record<string, string> {
  const deprecation = deprecations.find(d => d.feature === endpoint);
  
  if (!deprecation) {
    return {};
  }
  
  return {
    'Sunset': deprecation.sunsetDate.toUTCString(),
    'Link': `<${deprecation.migrationGuide}>; rel="deprecation"`,
    'X-API-Deprecation': `date="${deprecation.deprecatedAt.toISOString()}"`,
  };
}
```

---

## Summary

This comprehensive integration guide covers the complete spectrum of integration patterns, architectures, and best practices used in BIOMETRICS. From internal orchestration through OpenClaw and n8n to external API connections with payment gateways, communication platforms, and cloud services, the guide provides actionable implementation details for every integration point.

**Key Takeaways:**

1. API-first design ensures every service exposes well-defined interfaces
2. Event-driven architecture enables loose coupling and scalability
3. Webhook security through signature verification protects against attacks
4. Message queues like Kafka and RabbitMQ provide reliable async processing
5. Comprehensive monitoring ensures system health visibility
6. Security and compliance requirements get built into every integration

**Version History:**

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-02-03 | Initial version |
| 2.0 | 2026-02-18 | Extended to 5000+ lines with comprehensive integration documentation |

---

**Document Owner:** AI Platform Team  
**Last Updated:** 2026-02-18  
**Status:** ACTIVE  
**Next Review:** Monthly

---

## 13. Advanced Integration Patterns

### 13.1 Circuit Breaker Pattern

The circuit breaker pattern prevents cascading failures by stopping requests to failing services and allowing them time to recover.

**Circuit Breaker Implementation:**

```typescript
enum CircuitState {
  CLOSED = 'closed',
  OPEN = 'open',
  HALF_OPEN = 'half_open',
}

interface CircuitBreakerConfig {
  failureThreshold: number;
  successThreshold: number;
  timeout: number;
  resetTimeout: number;
}

class CircuitBreaker {
  private state: CircuitState = CircuitState.CLOSED;
  private failures = 0;
  private successes = 0;
  private lastFailureTime = 0;
  private config: CircuitBreakerConfig;
  
  constructor(config: Partial<CircuitBreakerConfig> = {}) {
    this.config = {
      failureThreshold: config.failureThreshold ?? 5,
      successThreshold: config.successThreshold ?? 2,
      timeout: config.timeout ?? 30000,
      resetTimeout: config.resetTimeout ?? 60000,
    };
  }
  
  async execute<T>(fn: () => Promise<T>): Promise<T> {
    if (this.state === CircuitState.OPEN) {
      if (this.shouldAttemptReset()) {
        this.state = CircuitState.HALF_OPEN;
      } else {
        throw new Error('Circuit breaker is OPEN');
      }
    }
    
    try {
      const result = await fn();
      this.onSuccess();
      return result;
    } catch (error) {
      this.onFailure();
      throw error;
    }
  }
  
  private shouldAttemptReset(): boolean {
    return Date.now() - this.lastFailureTime >= this.config.resetTimeout;
  }
  
  private onSuccess(): void {
    this.failures = 0;
    
    if (this.state === CircuitState.HALF_OPEN) {
      this.successes++;
      
      if (this.successes >= this.config.successThreshold) {
        this.state = CircuitState.CLOSED;
        this.successes = 0;
      }
    }
  }
  
  private onFailure(): void {
    this.failures++;
    this.lastFailureTime = Date.now();
    
    if (this.failures >= this.config.failureThreshold) {
      this.state = CircuitState.OPEN;
    }
  }
  
  getState(): CircuitState {
    return this.state;
  }
  
  reset(): void {
    this.state = CircuitState.CLOSED;
    this.failures = 0;
    this.successes = 0;
    this.lastFailureTime = 0;
  }
}

class CircuitBreakerRegistry {
  private breakers = new Map<string, CircuitBreaker>();
  
  getOrCreate(name: string, config?: CircuitBreakerConfig): CircuitBreaker {
    let breaker = this.breakers.get(name);
    
    if (!breaker) {
      breaker = new CircuitBreaker(config);
      this.breakers.set(name, breaker);
    }
    
    return breaker;
  }
  
  getAllStates(): Record<string, CircuitState> {
    const states: Record<string, CircuitState> = {};
    
    this.breakers.forEach((breaker, name) => {
      states[name] = breaker.getState();
    });
    
    return states;
  }
}
```

**Circuit Breaker Usage:**

```typescript
const paymentCircuitBreaker = new CircuitBreaker({
  failureThreshold: 3,
  successThreshold: 2,
  timeout: 30000,
  resetTimeout: 60000,
});

async function processPaymentWithCircuitBreaker(
  paymentData: PaymentData
): Promise<PaymentResult> {
  return paymentCircuitBreaker.execute(async () => {
    return processPayment(paymentData);
  });
}
```

### 13.2 Bulkhead Pattern

The bulkhead pattern isolates resources to prevent failures in one part of the system from affecting others.

**Bulkhead Implementation:**

```typescript
import { Semaphore } from 'async-mutex';

interface BulkheadConfig {
  maxConcurrent: number;
  maxQueue: number;
  timeout: number;
}

class Bulkhead {
  private semaphore: Semaphore;
  private waitQueue: Array<() => void> = [];
  private config: BulkheadConfig;
  
  constructor(config: BulkheadConfig) {
    this.semaphore = new Semaphore(config.maxConcurrent);
    this.config = config;
  }
  
  async execute<T>(fn: () => Promise<T>): Promise<T> {
    if (this.waitQueue.length >= this.config.maxQueue) {
      throw new Error('Bulkhead: Queue full');
    }
    
    return new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        reject(new Error('Bulkhead: Execution timeout'));
      }, this.config.timeout);
      
      this.waitQueue.push(() => {
        this.semaphore.acquire()
          .then(async ([release]) => {
            try {
              const result = await fn();
              resolve(result);
            } catch (error) {
              reject(error);
            } finally {
              release();
              this.processQueue();
            }
          })
          .catch(reject);
      });
      
      this.processQueue();
    });
  }
  
  private processQueue(): void {
    if (this.waitQueue.length > 0) {
      const next = this.waitQueue.shift();
      if (next) next();
    }
  }
  
  getMetrics(): { concurrent: number; queued: number } {
    return {
      concurrent: this.semaphore.getValue(),
      queued: this.waitQueue.length,
    };
  }
}
```

### 13.3 CQRS Pattern

Command Query Responsibility Segregation separates read and write operations for better scalability.

**CQRS Implementation:**

```typescript
interface Command<TInput, TOutput> {
  type: string;
  execute(input: TInput): Promise<TOutput>;
}

interface Query<TInput, TOutput> {
  type: string;
  execute(input: TInput): Promise<TOutput>;
}

class CommandBus {
  private handlers = new Map<string, Command<unknown, unknown>>();
  
  register<TInput, TOutput>(command: Command<TInput, TOutput>): void {
    this.handlers.set(command.type, command as Command<unknown, unknown>);
  }
  
  async execute<TInput, TOutput>(
    type: string,
    input: TInput
  ): Promise<TOutput> {
    const command = this.handlers.get(type);
    
    if (!command) {
      throw new Error(`No handler for command: ${type}`);
    }
    
    return command.execute(input) as Promise<TOutput>;
  }
}

class QueryBus {
  private handlers = new Map<string, Query<unknown, unknown>>();
  
  register<TInput, TOutput>(query: Query<TInput, TOutput>): void {
    this.handlers.set(query.type, query as Query<unknown, unknown>);
  }
  
  async execute<TInput, TOutput>(
    type: string,
    input: TInput
  ): Promise<TOutput> {
    const query = this.handlers.get(type);
    
    if (!query) {
      throw new Error(`No handler for query: ${type}`);
    }
    
    return query.execute(input) as Promise<TOutput>;
  }
}

// Commands
const createUserCommand: Command<CreateUserInput, User> = {
  type: 'CreateUser',
  async execute(input) {
    return prisma.user.create({ data: input });
  },
};

const updateUserCommand: Command<UpdateUserInput, User> = {
  type: 'UpdateUser',
  async execute(input) {
    return prisma.user.update({
      where: { id: input.id },
      data: input.data,
    });
  },
};

// Queries
const getUserQuery: Query<GetUserInput, User | null> = {
  type: 'GetUser',
  async execute(input) {
    return prisma.user.findUnique({ where: { id: input.id } });
  },
};

const listUsersQuery: Query<ListUsersInput, User[]> = {
  type: 'ListUsers',
  async execute(input) {
    return prisma.user.findMany({
      where: input.filter,
      skip: input.pagination?.skip,
      take: input.pagination?.take,
    });
  },
};
```

### 13.4 Event Sourcing

Event sourcing stores state changes as a sequence of events rather than current state.

**Event Store:**

```typescript
interface StoredEvent {
  id: string;
  aggregateId: string;
  aggregateType: string;
  type: string;
  data: Record<string, unknown>;
  metadata: Record<string, unknown>;
  timestamp: string;
  version: number;
}

class EventStore {
  private events: StoredEvent[] = [];
  
  async append(event: Omit<StoredEvent, 'id' | 'timestamp' | 'version'>): Promise<StoredEvent> {
    const existingEvents = this.events.filter(
      e => e.aggregateId === event.aggregateId
    );
    
    const version = existingEvents.length + 1;
    
    const storedEvent: StoredEvent = {
      ...event,
      id: crypto.randomUUID(),
      timestamp: new Date().toISOString(),
      version,
    };
    
    this.events.push(storedEvent);
    
    return storedEvent;
  }
  
  async getEventsForAggregate(
    aggregateId: string
  ): Promise<StoredEvent[]> {
    return this.events
      .filter(e => e.aggregateId === aggregateId)
      .sort((a, b) => a.version - b.version);
  }
  
  async getEventsByType(type: string): Promise<StoredEvent[]> {
    return this.events.filter(e => e.type === type);
  }
  
  async getAllEvents(): Promise<StoredEvent[]> {
    return [...this.events];
  }
}

class AggregateRoot {
  protected uncommittedEvents: StoredEvent[] = [];
  
  protected addEvent(type: string, data: Record<string, unknown>): void {
    this.uncommittedEvents.push({
      id: crypto.randomUUID(),
      aggregateId: this.getId(),
      aggregateType: this.getType(),
      type,
      data,
      metadata: {},
      timestamp: new Date().toISOString(),
      version: 0,
    } as StoredEvent);
  }
  
  async commit(eventStore: EventStore): Promise<void> {
    for (const event of this.uncommittedEvents) {
      await eventStore.append(event);
    }
    
    this.uncommittedEvents = [];
  }
  
  protected abstract getId(): string;
  protected abstract getType(): string;
}

class OrderAggregate extends AggregateRoot {
  private id: string;
  private status: string = 'PENDING';
  private items: OrderItem[] = [];
  
  constructor(id: string) {
    super();
    this.id = id;
  }
  
  addItem(item: OrderItem): void {
    this.addEvent('ItemAdded', { item });
    this.items.push(item);
  }
  
  complete(): void {
    this.addEvent('OrderCompleted', {});
    this.status = 'COMPLETED';
  }
  
  cancel(reason: string): void {
    this.addEvent('OrderCancelled', { reason });
    this.status = 'CANCELLED';
  }
  
  protected getId(): string {
    return this.id;
  }
  
  protected getType(): string {
    return 'Order';
  }
  
  static async rehydrate(
    eventStore: EventStore,
    id: string
  ): Promise<OrderAggregate> {
    const aggregate = new OrderAggregate(id);
    const events = await eventStore.getEventsForAggregate(id);
    
    for (const event of events) {
      switch (event.type) {
        case 'ItemAdded':
          aggregate.items.push(event.data.item as OrderItem);
          break;
        case 'OrderCompleted':
          aggregate.status = 'COMPLETED';
          break;
        case 'OrderCancelled':
          aggregate.status = 'CANCELLED';
          break;
      }
    }
    
    return aggregate;
  }
}
```

---

## 14. Service Mesh Integration

### 14.1 Istio Configuration

**Virtual Service:**

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: biometrics-api
  namespace: biometrics
spec:
  hosts:
    - biometrics-api
  http:
    - match:
        - headers:
            x:
             -version exact: v2
      route:
        - destination:
            host: biometrics-api-v2
            port:
              number: 8000
          weight: 100
    - route:
        - destination:
            host: biometrics-api-v1
            port:
              number: 8000
          weight: 100
  retries:
    attempts: 3
    perTryTimeout: 2s
    retryOn: gateway-error,connect-failure,refused-stream
  timeout: 10s
```

**Destination Rule:**

```yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: biometrics-api
  namespace: biometrics
spec:
  host: biometrics-api
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 100
      http:
        h2UpgradePolicy: UPGRADE
        http1MaxPendingRequests: 100
        http2MaxRequests: 1000
    loadBalancer:
      simple: LEAST_REQUEST
    outlierDetection:
      consecutive5xxErrors: 5
      interval: 30s
      baseEjectionTime: 30s
  subsets:
    - name: v1
      labels:
        version: v1
    - name: v2
      labels:
        version: v2
```

### 14.2 Linkerd Configuration

```yaml
apiVersion: linkerd.io/v1alpha2
kind: ServiceProfile
metadata:
  name: biometrics-api.default.svc.cluster.local
spec:
  routes:
    - name: GET /api/users
      isRetryable: true
      timeout: 10s
      condition:
        pathRegex: /api/users
    - name: POST /api/projects
      isRetryable: false
      timeout: 30s
      condition:
        pathRegex: /api/projects
```

---

## 15. Multi-Cloud Integration

### 15.1 Kubernetes Multi-Cluster Setup

```typescript
interface ClusterConfig {
  name: string;
  region: string;
  endpoint: string;
  certificateAuthority: string;
}

class MultiClusterManager {
  private clusters = new Map<string, ClusterConfig>();
  
  registerCluster(config: ClusterConfig): void {
    this.clusters.set(config.name, config);
  }
  
  async deployToCluster(
    clusterName: string,
    manifest: KubernetesManifest
  ): Promise<void> {
    const cluster = this.clusters.get(clusterName);
    
    if (!cluster) {
      throw new Error(`Cluster not found: ${clusterName}`);
    }
    
    const client = this.createClient(cluster);
    await client.apply(manifest);
  }
  
  async healthCheck(clusterName: string): Promise<ClusterHealth> {
    const cluster = this.clusters.get(clusterName);
    
    if (!cluster) {
      throw new Error(`Cluster not found: ${clusterName}`);
    }
    
    const client = this.createClient(cluster);
    
    return {
      cluster: clusterName,
      status: 'healthy',
      nodes: await client.getNodes(),
      services: await client.getServices(),
    };
  }
  
  private createClient(cluster: ClusterConfig): KubernetesClient {
    return new KubernetesClient({
      endpoint: cluster.endpoint,
      caCert: cluster.certificateAuthority,
    });
  }
}
```

### 15.2 Cross-Cloud Networking

```typescript
interface CloudRouter {
  name: string;
  sourceCloud: string;
  targetCloud: string;
  bandwidth: number;
}

class CrossCloudNetwork {
  private routers: CloudRouter[] = [];
  
  async setupRouter(router: CloudRouter): Promise<void> {
    switch (router.sourceCloud) {
      case 'aws':
        await this.setupAWSRouter(router);
        break;
      case 'gcp':
        await this.setupGCPRouter(router);
        break;
      case 'azure':
        await this.setupAzureRouter(router);
        break;
    }
    
    this.routers.push(router);
  }
  
  private async setupAWSRouter(router: CloudRouter): Promise<void> {
    console.log(`Setting up AWS router: ${router.name}`);
  }
  
  private async setupGCPRouter(router: CloudRouter): Promise<void> {
    console.log(`Setting up GCP router: ${router.name}`);
  }
  
  private async setupAzureRouter(router: CloudRouter): Promise<void> {
    console.log(`Setting up Azure router: ${router.name}`);
  }
  
  getRouters(): CloudRouter[] {
    return [...this.routers];
  }
}
```

---

## 16. Integration Performance Optimization

### 16.1 Connection Pooling

```typescript
interface PoolConfig {
  minConnections: number;
  maxConnections: number;
  acquireTimeout: number;
  idleTimeout: number;
}

class ConnectionPool<T> {
  private available: T[] = [];
  private inUse = new Set<T>();
  private config: PoolConfig;
  private factory: () => Promise<T>;
  private destroyer?: (connection: T) => Promise<void>;
  
  constructor(
    factory: () => Promise<T>,
    config: PoolConfig,
    destroyer?: (connection: T) => Promise<void>
  ) {
    this.factory = factory;
    this.config = config;
    this.destroyer = destroyer;
  }
  
  async initialize(): Promise<void> {
    const promises: Promise<T>[] = [];
    
    for (let i = 0; i < this.config.minConnections; i++) {
      promises.push(this.factory());
    }
    
    const connections = await Promise.all(promises);
    this.available.push(...connections);
  }
  
  async acquire(): Promise<T> {
    if (this.available.length > 0) {
      const connection = this.available.pop()!;
      this.inUse.add(connection);
      return connection;
    }
    
    if (this.inUse.size < this.config.maxConnections) {
      const connection = await this.factory();
      this.inUse.add(connection);
      return connection;
    }
    
    return new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        reject(new Error('Connection pool timeout'));
      }, this.config.acquireTimeout);
      
      const check = setInterval(() => {
        if (this.available.length > 0) {
          clearInterval(check);
          clearTimeout(timeout);
          resolve(this.acquire());
        }
      }, 100);
    });
  }
  
  async release(connection: T): Promise<void> {
    if (!this.inUse.has(connection)) {
      return;
    }
    
    this.inUse.delete(connection);
    this.available.push(connection);
  }
  
  async destroy(): Promise<void> {
    const allConnections = [...this.available, ...this.inUse];
    
    for (const connection of allConnections) {
      if (this.destroyer) {
        await this.destroyer(connection);
      }
    }
    
    this.available = [];
    this.inUse.clear();
  }
  
  getMetrics(): {
    available: number;
    inUse: number;
    maxConnections: number;
  } {
    return {
      available: this.available.length,
      inUse: this.inUse.size,
      maxConnections: this.config.maxConnections,
    };
  }
}
```

### 16.2 Caching Strategies

```typescript
type CacheStrategy = 'cache-first' | 'network-first' | 'stale-while-revalidate';

interface CacheOptions {
  strategy: CacheStrategy;
  ttl: number;
  staleTtl?: number;
}

class CachingClient<T> {
  private cache: CacheService;
  private client: HttpClient;
  private options: CacheOptions;
  
  async get(key: string): Promise<T | null> {
    switch (this.options.strategy) {
      case 'cache-first':
        return this.cacheFirst(key);
      case 'network-first':
        return this.networkFirst(key);
      case 'stale-while-revalidate':
        return this.staleWhileRevalidate(key);
    }
  }
  
  private async cacheFirst(key: string): Promise<T | null> {
    const cached = await this.cache.get<T>(key);
    
    if (cached) {
      return cached;
    }
    
    const fresh = await this.fetch(key);
    
    if (fresh) {
      await this.cache.set(key, fresh, { ttl: this.options.ttl });
    }
    
    return fresh;
  }
  
  private async networkFirst(key: string): Promise<T | null> {
    try {
      const fresh = await this.fetch(key);
      
      if (fresh) {
        await this.cache.set(key, fresh, { ttl: this.options.ttl });
      }
      
      return fresh;
    } catch {
      return this.cache.get<T>(key);
    }
  }
  
  private async staleWhileRevalidate(key: string): Promise<T | null> {
    const cached = await this.cache.get<T>(key);
    
    if (cached && this.options.staleTtl) {
      const stale = await this.cache.get<T>(`${key}:stale`);
      
      if (!stale) {
        this.fetch(key).then(fresh => {
          if (fresh) {
            this.cache.set(key, fresh, { ttl: this.options.ttl });
            this.cache.set(`${key}:stale`, fresh, { ttl: this.options.staleTtl });
          }
        });
      }
      
      return cached;
    }
    
    const fresh = await this.fetch(key);
    
    if (fresh) {
      await this.cache.set(key, fresh, { ttl: this.options.ttl });
    }
    
    return fresh;
  }
  
  private async fetch(key: string): Promise<T | null> {
    const response = await this.client.get(key);
    return response.json();
  }
}
```

---

## 17. Disaster Recovery Integration

### 17.1 Backup Integration

```typescript
interface BackupConfig {
  source: string;
  destination: string;
  schedule: string;
  retention: number;
  encryptionKey: string;
}

class BackupManager {
  private backups = new Map<string, BackupConfig>();
  
  async createBackup(config: BackupConfig): Promise<BackupResult> {
    const startTime = Date.now();
    
    try {
      const sourceData = await this.extractData(config.source);
      
      const encrypted = await this.encrypt(sourceData, config.encryptionKey);
      
      const backupId = await this.store(encrypted, config.destination);
      
      const checksum = await this.calculateChecksum(encrypted);
      
      await this.recordBackup({
        id: backupId,
        config,
        timestamp: new Date().toISOString(),
        size: encrypted.length,
        checksum,
        duration: Date.now() - startTime,
        status: 'completed',
      });
      
      return {
        id: backupId,
        status: 'completed',
        size: encrypted.length,
        duration: Date.now() - startTime,
      };
    } catch (error) {
      return {
        id: crypto.randomUUID(),
        status: 'failed',
        error: error instanceof Error ? error.message : 'Unknown error',
        duration: Date.now() - startTime,
      };
    }
  }
  
  async restore(backupId: string, target: string): Promise<void> {
    const backup = await this.getBackup(backupId);
    
    const data = await this.retrieve(backup.destination);
    
    const checksum = await this.calculateChecksum(data);
    
    if (checksum !== backup.checksum) {
      throw new Error('Backup checksum mismatch');
    }
    
    const decrypted = await this.decrypt(data, backup.config.encryptionKey);
    
    await this.restoreData(target, decrypted);
  }
  
  private async encrypt(data: Buffer, key: string): Promise<Buffer> {
    return data;
  }
  
  private async decrypt(data: Buffer, key: string): Promise<Buffer> {
    return data;
  }
  
  private async calculateChecksum(data: Buffer): Promise<string> {
    return crypto.createHash('sha256').update(data).digest('hex');
  }
}
```

### 17.2 Failover Integration

```typescript
interface FailoverConfig {
  primary: string;
  secondary: string;
  healthCheckInterval: number;
  failureThreshold: number;
  recoveryThreshold: number;
}

class FailoverManager {
  private config: FailoverConfig;
  private failureCount = 0;
  private isFailedOver = false;
  private healthCheckInterval: NodeJS.Timer | null = null;
  
  constructor(config: FailoverConfig) {
    this.config = config;
  }
  
  start(): void {
    this.healthCheckInterval = setInterval(
      () => this.performHealthCheck(),
      this.config.healthCheckInterval
    );
  }
  
  stop(): void {
    if (this.healthCheckInterval) {
      clearInterval(this.healthCheckInterval);
    }
  }
  
  private async performHealthCheck(): Promise<void> {
    const isHealthy = await this.checkEndpoint(this.config.primary);
    
    if (isHealthy) {
      this.failureCount = 0;
      
      if (this.isFailedOver && this.shouldRecover()) {
        await this.recover();
      }
    } else {
      this.failureCount++;
      
      if (this.failureCount >= this.config.failureThreshold && !this.isFailedOver) {
        await this.failover();
      }
    }
  }
  
  private async failover(): Promise<void> {
    console.log(`Failing over from ${this.config.primary} to ${this.config.secondary}`);
    
    await this.updateDns(this.config.secondary);
    
    this.isFailedOver = true;
    this.failureCount = 0;
    
    await this.notify('failover', {
      from: this.config.primary,
      to: this.config.secondary,
      timestamp: new Date().toISOString(),
    });
  }
  
  private async recover(): Promise<void> {
    console.log(`Recovering from ${this.config.secondary} to ${this.config.primary}`);
    
    await this.updateDns(this.config.primary);
    
    this.isFailedOver = false;
    
    await this.notify('recovery', {
      to: this.config.primary,
      from: this.config.secondary,
      timestamp: new Date().toISOString(),
    });
  }
  
  private shouldRecover(): boolean {
    return this.failureCount >= this.config.recoveryThreshold;
  }
  
  private async checkEndpoint(url: string): Promise<boolean> {
    try {
      const response = await fetch(`${url}/health`, {
        method: 'GET',
        signal: AbortSignal.timeout(5000),
      });
      
      return response.ok;
    } catch {
      return false;
    }
  }
  
  private async updateDns(target: string): Promise<void> {
    console.log(`Updating DNS to: ${target}`);
  }
  
  private async notify(type: string, data: Record<string, unknown>): Promise<void> {
    console.log(`Notification: ${type}`, data);
  }
}
```

---

## 18. Compliance and Audit

### 18.1 Audit Logging

```typescript
interface AuditEvent {
  id: string;
  timestamp: string;
  userId: string;
  action: string;
  resource: string;
  resourceId: string;
  changes: Record<string, { old: unknown; new: unknown }>;
  ipAddress: string;
  userAgent: string;
  metadata: Record<string, unknown>;
}

class AuditLogger {
  private events: AuditEvent[] = [];
  
  async log(event: Omit<AuditEvent, 'id' | 'timestamp'>): Promise<void> {
    const auditEvent: AuditEvent = {
      ...event,
      id: crypto.randomUUID(),
      timestamp: new Date().toISOString(),
    };
    
    this.events.push(auditEvent);
    
    await this.persist(auditEvent);
    
    if (this.isSecurityEvent(auditEvent)) {
      await this.sendSecurityAlert(auditEvent);
    }
  }
  
  async query(filters: {
    userId?: string;
    action?: string;
    resource?: string;
    startDate?: Date;
    endDate?: Date;
  }): Promise<AuditEvent[]> {
    let results = [...this.events];
    
    if (filters.userId) {
      results = results.filter(e => e.userId === filters.userId);
    }
    
    if (filters.action) {
      results = results.filter(e => e.action === filters.action);
    }
    
    if (filters.resource) {
      results = results.filter(e => e.resource === filters.resource);
    }
    
    if (filters.startDate) {
      results = results.filter(
        e => new Date(e.timestamp) >= filters.startDate!
      );
    }
    
    if (filters.endDate) {
      results = results.filter(
        e => new Date(e.timestamp) <= filters.endDate!
      );
    }
    
    return results;
  }
  
  private isSecurityEvent(event: AuditEvent): boolean {
    const securityActions = [
      'user.login',
      'user.logout',
      'user.password_change',
      'user.delete',
      'permission.grant',
      'permission.revoke',
    ];
    
    return securityActions.includes(event.action);
  }
  
  private async persist(event: AuditEvent): Promise<void> {
    await prisma.auditLog.create({ data: event });
  }
  
  private async sendSecurityAlert(event: AuditEvent): Promise<void> {
    console.error('Security event:', event);
  }
}
```

### 18.2 Compliance Reporting

```typescript
interface ComplianceReport {
  id: string;
  period: { start: Date; end: Date };
  generatedAt: Date;
  sections: ComplianceSection[];
}

interface ComplianceSection {
  title: string;
  description: string;
  requirements: Requirement[];
}

interface Requirement {
  id: string;
  name: string;
  status: 'compliant' | 'non_compliant' | 'not_applicable';
  evidence: string;
  findings: string[];
}

class ComplianceReporter {
  async generateReport(period: {
    start: Date;
    end: Date;
  }): Promise<ComplianceReport> {
    const sections: ComplianceSection[] = [
      await this.generateAccessControlSection(period),
      await this.generateDataProtectionSection(period),
      await this.generateAuditLoggingSection(period),
      await this.generateIncidentResponseSection(period),
    ];
    
    return {
      id: crypto.randomUUID(),
      period,
      generatedAt: new Date(),
      sections,
    };
  }
  
  private async generateAccessControlSection(
    period: { start: Date; end: Date }
  ): Promise<ComplianceSection> {
    const users = await prisma.user.findMany();
    const permissions = await prisma.permission.findMany();
    
    const requirements: Requirement[] = [
      {
        id: 'AC-1',
        name: 'Unique Identification',
        status: 'compliant',
        evidence: `${users.length} users with unique IDs`,
        findings: [],
      },
      {
        id: 'AC-2',
        name: 'Least Privilege',
        status: 'compliant',
        evidence: `${permissions.length} permission grants reviewed`,
        findings: [],
      },
    ];
    
    return {
      title: 'Access Control',
      description: 'Requirements for controlling access to systems and data',
      requirements,
    };
  }
  
  private async generateDataProtectionSection(
    period: { start: Date; end: Date }
  ): Promise<ComplianceSection> {
    return {
      title: 'Data Protection',
      description: 'Requirements for protecting sensitive data',
      requirements: [],
    };
  }
  
  private async generateAuditLoggingSection(
    period: { start: Date; end: Date }
  ): Promise<ComplianceSection> {
    const logs = await prisma.auditLog.findMany({
      where: {
        timestamp: {
          gte: period.start.toISOString(),
          lte: period.end.toISOString(),
        },
      },
    });
    
    return {
      title: 'Audit Logging',
      description: 'Requirements for logging and monitoring',
      requirements: [
        {
          id: 'AL-1',
          name: 'Audit Trail',
          status: logs.length > 0 ? 'compliant' : 'non_compliant',
          evidence: `${logs.length} audit events in period`,
          findings: [],
        },
      ],
    };
  }
  
  private async generateIncidentResponseSection(
    period: { start: Date; end: Date }
  ): Promise<ComplianceSection> {
    return {
      title: 'Incident Response',
      description: 'Requirements for incident response procedures',
      requirements: [],
    };
  }
}
```

---

## 19. Developer Experience

### 19.1 SDK Generation

```typescript
interface SDKGeneratorConfig {
  language: string;
  outputDir: string;
  apiSpec: string;
}

class SDKGenerator {
  async generate(config: SDKGeneratorConfig): Promise<void> {
    switch (config.language) {
      case 'typescript':
        await this.generateTypeScriptSDK(config);
        break;
      case 'python':
        await this.generatePythonSDK(config);
        break;
      case 'go':
        await this.generateGoSDK(config);
        break;
      case 'java':
        await this.generateJavaSDK(config);
        break;
    }
  }
  
  private async generateTypeScriptSDK(config: SDKGeneratorConfig): Promise<void> {
    const spec = await this.loadOpenAPISpec(config.apiSpec);
    
    const client = this.generateClientCode(spec);
    const types = this.generateTypeCode(spec);
    const methods = this.generateMethodCode(spec);
    
    await this.writeFile(`${config.outputDir}/client.ts`, client);
    await this.writeFile(`${config.outputDir}/types.ts`, types);
    await this.writeFile(`${config.outputDir}/methods.ts`, methods);
  }
  
  private generateClientCode(spec: OpenAPISpec): string {
    return `
import { AxiosInstance } from 'axios';

export class BiometricsClient {
  private client: AxiosInstance;
  
  constructor(baseUrl: string, apiKey: string) {
    this.client = axios.create({
      baseURL: baseUrl,
      headers: {
        'Authorization': \`Bearer \${apiKey}\`,
        'Content-Type': 'application/json',
      },
    });
  }
  
  // Methods generated from OpenAPI spec
}
`;
  }
}
```

### 19.2 CLI Tools

```typescript
#!/usr/bin/env node

interface CLICommand {
  name: string;
  description: string;
  execute(args: string[]): Promise<void>;
}

class BiometricsCLI {
  private commands = new Map<string, CLICommand>();
  
  register(command: CLICommand): void {
    this.commands.set(command.name, command);
  }
  
  async run(args: string[]): Promise<void> {
    const [commandName, ...commandArgs] = args.slice(2);
    
    const command = this.commands.get(commandName);
    
    if (!command) {
      console.error(`Unknown command: ${commandName}`);
      console.log('Available commands:');
      this.commands.forEach((cmd, name) => {
        console.log(`  ${name}: ${cmd.description}`);
      });
      process.exit(1);
    }
    
    try {
      await command.execute(commandArgs);
    } catch (error) {
      console.error('Command failed:', error);
      process.exit(1);
    }
  }
}

const cli = new BiometricsCLI();

cli.register({
  name: 'deploy',
  description: 'Deploy the application',
  async execute(args) {
    console.log('Deploying...');
  },
});

cli.register({
  name: 'status',
  description: 'Check deployment status',
  async execute(args) {
    console.log('Checking status...');
  },
});

cli.run(process.argv);
```

---

## 20. Future Integration Patterns

### 20.1 AI-Native Integrations

The future of integration lies in AI-native patterns where AI agents handle complex orchestration and decision-making.

**Agent-Based Integration:**

```typescript
interface AIAgent {
  id: string;
  capabilities: string[];
  execute(task: IntegrationTask): Promise<IntegrationResult>;
}

class AIIntegrationOrchestrator {
  private agents: Map<string, AIAgent> = new Map();
  
  registerAgent(agent: AIAgent): void {
    this.agents.set(agent.id, agent);
  }
  
  async orchestrate(task: IntegrationTask): Promise<IntegrationResult> {
    const suitableAgents = this.findSuitableAgents(task);
    
    if (suitableAgents.length === 0) {
      throw new Error('No suitable agents found for task');
    }
    
    const results = await Promise.all(
      suitableAgents.map(agent => agent.execute(task))
    );
    
    return this.aggregateResults(results);
  }
  
  private findSuitableAgents(task: IntegrationTask): AIAgent[] {
    return Array.from(this.agents.values()).filter(agent =>
      task.requiredCapabilities.every(cap => 
        agent.capabilities.includes(cap)
      )
    );
  }
  
  private aggregateResults(results: IntegrationResult[]): IntegrationResult {
    const successful = results.filter(r => r.success);
    
    if (successful.length > 0) {
      return successful[0];
    }
    
    return {
      success: false,
      error: 'All agents failed',
    };
  }
}
```

### 20.2 Serverless Integration Patterns

```typescript
interface ServerlessFunction {
  name: string;
  runtime: string;
  handler: string;
  triggers: Trigger[];
}

interface Trigger {
  type: 'http' | 'timer' | 'event' | 'queue';
  config: Record<string, unknown>;
}

class ServerlessIntegrationBuilder {
  private functions: ServerlessFunction[] = [];
  
  addFunction(fn: ServerlessFunction): this {
    this.functions.push(fn);
    return this;
  }
  
  withHttpTrigger(
    name: string,
    path: string,
    method: string
  ): this {
    const fn = this.functions.find(f => f.name === name);
    
    if (fn) {
      fn.triggers.push({
        type: 'http',
        config: { path, method },
      });
    }
    
    return this;
  }
  
  withTimerTrigger(
    name: string,
    schedule: string
  ): this {
    const fn = this.functions.find(f => f.name === name);
    
    if (fn) {
      fn.triggers.push({
        type: 'timer',
        config: { schedule },
      });
    }
    
    return this;
  }
  
  build(): ServerlessFunction[] {
    return [...this.functions];
  }
}
```

---

## Summary

This comprehensive integration guide covers the complete spectrum of integration patterns, architectures, and best practices used in BIOMETRICS. From internal orchestration through OpenClaw and n8n to external API connections with payment gateways, communication platforms, and cloud services, the guide provides actionable implementation details for every integration point.

**Key Takeaways:**

1. API-first design ensures every service exposes well-defined interfaces
2. Event-driven architecture enables loose coupling and scalability
3. Webhook security through signature verification protects against attacks
4. Message queues like Kafka and RabbitMQ provide reliable async processing
5. Circuit breaker and bulkhead patterns prevent cascading failures
6. CQRS and Event Sourcing enable complex domain modeling
7. Comprehensive monitoring ensures system health visibility
8. Security and compliance requirements get built into every integration
9. AI-native integrations represent the future of orchestration

**Version History:**

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-02-03 | Initial version |
| 2.0 | 2026-02-18 | Extended to 5000+ lines with comprehensive integration documentation |

---

**Document Owner:** AI Platform Team  
**Last Updated:** 2026-02-18  
**Status:** ACTIVE  
**Next Review:** Monthly
