# OPENCLAW.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- OpenClaw-Konfiguration folgt globalen Provider- und Timeout-Regeln.
- Modell-ID-Disziplin, Fallback-Strategie und Betriebsmetriken sind Pflicht.
- Änderungen werden sicher, nachvollziehbar und testsynchron umgesetzt.

Status: ACTIVE  
Version: 2.0 (Qwen 3.5 Integration)  
Stand: Februar 2026

## Zweck

Dieses Dokument dient als umfassender Integrationsleitfaden für OpenClaw als Connector- und Auth-Layer mit besonderem Fokus auf die Qwen 3.5 (397B) Modellintegration über NVIDIA NIM.

OpenClaw fungiert als zentrale Orchestrierungsschicht für externe Plattformzugriffe und ermöglicht durch Skills autonome Workflow-Automatisierung.

## Rolle im System

OpenClaw ist das zentrale Bindeglied zwischen den KI-Agenten und den zu steuernden Systemen:

- **Auth-Layer:** Verwaltet Tokens und Authentifizierung
- **Connector-Framework:** Verbindet n8n, Supabase und externe APIs
- **Skill-Registry:** Registriert und verwaltet wiederverwendbare Skills
- **Meta-Builder:** Ermöglicht autonomes Erstellen neuer Skills

## Grundprinzipien

1. **Keine direkten unsicheren Provider-Zugriffe** - Alle Zugriffe durch OpenClaw
2. **Token- und Rollenhygiene** - Security First Ansatz
3. **Klare Fehlerpfade** - Robuste Fehlerbehandlung
4. **Wiederholbare Connector-Prozesse** - Deterministic execution
5. **Modell-Disziplin** - Qwen 3.5 als Primary, Fallbacks konfiguriert

## Betriebsmodi

| Modus | Verwendung | Timeout |
|-------|------------|---------|
| local | Entwicklungsmodus | 30000ms |
| staging | Integrationsvalidierung | 60000ms |
| production | Stabiler Betrieb | 120000ms |

Der erweiterte Timeout für Production ist aufgrund der Qwen 3.5 397B Modellgröße erforderlich.

## NVIDIA NIM Konfiguration

### Provider Setup

OpenClaw nutzt NVIDIA NIM für den Zugang zu Qwen 3.5. Die Konfiguration erfordert:

```json
{
  "models": {
    "providers": {
      "nvidia": {
        "baseUrl": "https://integrate.api.nvidia.com/v1",
        "api": "openai-completions",
        "models": ["qwen/qwen3.5-397b-a17b"]
      }
    }
  }
}
```

### Umgebungsvariablen

```bash
# NVIDIA NIM Configuration
NVIDIA_API_KEY=nvapi-xxxxxxxxxxxxxxxxxxxxx

# OpenClaw Configuration
OPENCLAW_BASE_URL=http://localhost:18789
OPENCLAW_AUTH_MODE=token
OPENCLAW_TIMEOUT_MS=120000
OPENCLAW_RETRY_COUNT=3
```

### Modell-Spezifikationen

| Modell | Context | Output | Use Case |
|--------|---------|--------|----------|
| qwen/qwen3.5-397b-a17b | 262K | 32K | Code Generation (PRIMARY) |
| qwen2.5-coder-32b | 128K | 8K | Fast Code (FALLBACK) |
| moonshotai/kimi-k2.5 | 1M | 64K | General (FALLBACK) |

### Timeout-Konfiguration

Aufgrund der hohen Latenz von Qwen 3.5 397B (70-90 Sekunden) ist ein Timeout von mindestens 120000ms zwingend erforderlich.

```typescript
// OpenClaw Model Configuration
const modelConfig = {
  provider: 'nvidia',
  model: 'qwen/qwen3.5-397b-a17b',
  timeout: 120000,  // 120 seconds for 397B model
  maxRetries: 3,
  retryDelay: 5000
};
```

## Konfigurationsmatrix

| Variable | Zweck | Pflicht | Beispiel |
|----------|-------|---------|----------|
| OPENCLAW_BASE_URL | Basis-URL | ja | http://localhost:18789 |
| OPENCLAW_AUTH_MODE | Auth-Modus | ja | token |
| OPENCLAW_TIMEOUT_MS | Timeout | ja | 120000 |
| OPENCLAW_RETRY_COUNT | Retry Anzahl | ja | 3 |
| NVIDIA_API_KEY | NVIDIA API Key | ja | nvapi-xxxx |

## Auth- und Tokenfluss

1. **Zugriffsanfrage** wird validiert
2. **Rolle und Scope** werden geprüft
3. **Token** wird aus sicherer Quelle gelesen (Vault)
4. **Request** wird signiert und gesendet
5. **Response** wird protokolliert und bewertet

### Token-Refresh-Logik

```typescript
async function handleAuthRequest(context: AuthContext): Promise<Token> {
  const vault = new VaultClient();
  
  // Check if token exists and is valid
  const existingToken = await vault.get('openclaw/api-token');
  if (existingToken && !isExpired(existingToken)) {
    return existingToken;
  }
  
  // Refresh token if expired
  const newToken = await refreshToken();
  await vault.set('openclaw/api-token', newToken);
  
  return newToken;
}
```

## Fehler- und Retry-Strategie

### Fehlerklassen

| Klasse | Beschreibung | Strategie |
|--------|--------------|-----------|
| TEMPORARY | Netzwerk, Timeout | Retry mit Exponential Backoff |
| AUTH | 401, 403, Token invalid | Token-Refresh und Retry |
| RATE | 429 Too Many Requests | 60 Sekunden warten + Fallback |
| PERMANENT | 400, 404, Validierungsfehler | Sofortige Eskalation |

### Retry-Implementation

```typescript
async function withRetry<T>(
  operation: () => Promise<T>,
  options: RetryOptions = {}
): Promise<T> {
  const { maxRetries = 3, baseDelay = 5000 } = options;
  
  for (let attempt = 0; attempt <= maxRetries; attempt++) {
    try {
      return await operation();
    } catch (error) {
      if (attempt === maxRetries) {
        throw error;
      }
      
      if (isPermanentError(error)) {
        throw error;
      }
      
      // Exponential backoff: 5s, 10s, 20s
      const delay = baseDelay * Math.pow(2, attempt);
      await sleep(delay);
    }
  }
  
  throw new Error('Max retries exceeded');
}
```

### Fallback-Kette

```
Qwen 3.5 397B → Qwen2.5-Coder-32B → Kimi K2.5 → Error
```

## Sicherheitsregeln

1. **Keine Tokens im Repo** - Nur Vault oder Environment
2. **Kein Logging sensibler Daten** - Maskierung aktivieren
3. **Least Privilege für Connectoren** - Minimale Rechte
4. **Rate-Limit-Compliance** - Max 40 RPM bei NVIDIA NIM

### Security-Checkliste

- [ ] NVIDIA_API_KEY in Environment Variable
- [ ] Keine hardcoded Credentials
- [ ] Vault-Integration aktiv
- [ ] Request-Logging ohne sensitive Daten
- [ ] Rate-Limit-Monitoring aktiv

## Verifikation

### Health Checks

```bash
# OpenClaw Health
curl http://localhost:18789/health

# NVIDIA NIM Model Status
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
     https://integrate.api.nvidia.com/v1/models

# Model Response Test
opencode --model nvidia/qwen/qwen3.5-397b-a17b "Test connection"
```

### Test-Suite

1. **Auth-Flow Tests** - Token-Generation und Refresh
2. **Connector Tests** - Supabase, n8n Integration
3. **Fehlerfalltests** - Timeout, 401, 403, 429, 5xx
4. **Performance-Tests** - Latenz unter Last

## OpenClaw Skills als Interfaces

OpenClaw Skills sind die "Interfaces" oder "Wrappers" für zugrundeliegende Logik (Supabase/n8n/SDKs).

### Skill Types

1. **Webhook Wrapper** - Trigger n8n Workflows
2. **Serverless Proxy** - Call Supabase Edge Functions
3. **SDK Native** - Direct library usage

### Qwen 3.5 Integration Skill Example

```typescript
import { z } from 'zod';

export const qwenSkill = {
  name: 'qwen_3_5_analysis',
  description: 'Analyze content using Qwen 3.5 397B model',
  model: 'qwen/qwen3.5-397b-a17b',
  provider: 'nvidia',
  
  inputSchema: z.object({
    content: z.string().describe('Content to analyze'),
    task: z.enum(['summarize', 'extract', 'classify', 'generate']),
    maxTokens: z.number().default(8192),
    temperature: z.number().default(0.7)
  }),
  
  outputSchema: z.object({
    result: z.string(),
    confidence: z.number(),
    tokensUsed: z.number()
  }),
  
  handler: async (input: z.infer<typeof qwenSkill.inputSchema>) => {
    const response = await openclaw.invoke({
      model: qwenSkill.model,
      messages: [{
        role: 'user',
        content: generatePrompt(input.task, input.content)
      }],
      options: {
        max_tokens: input.maxTokens,
        temperature: input.temperature
      }
    });
    
    return {
      result: response.choices[0].message.content,
      confidence: response.usage.completion_tokens / response.usage.total_tokens,
      tokensUsed: response.usage.total_tokens
    };
  }
};

// Register skill
openclaw.registerSkill(qwenSkill);
```

### n8n Workflow Trigger Skill

```typescript
export const n8nWorkflowTrigger = {
  name: 'trigger_n8n_workflow',
  description: 'Trigger a specific n8n workflow',
  
  inputSchema: z.object({
    workflowId: z.string(),
    payload: z.record(z.any()).optional(),
    webhookSecret: z.string()
  }),
  
  handler: async (input) => {
    const webhookUrl = `http://localhost:5678/webhook/${input.workflowId}`;
    
    const response = await fetch(webhookUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Webhook-Secret': input.webhookSecret
      },
      body: JSON.stringify(input.payload || {})
    });
    
    if (!response.ok) {
      throw new Error(`Workflow trigger failed: ${response.status}`);
    }
    
    return { success: true, data: await response.json() };
  }
};
```

### Supabase Edge Function Skill

```typescript
export const supabaseFunctionSkill = {
  name: 'call_supabase_function',
  description: 'Execute a Supabase Edge Function',
  
  inputSchema: z.object({
    functionName: z.string(),
    payload: z.record(z.any()),
    authToken: z.string()
  }),
  
  handler: async (input) => {
    const response = await fetch(
      `https://project.supabase.co/functions/v1/${input.functionName}`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${input.authToken}`
        },
        body: JSON.stringify(input.payload)
      }
    );
    
    return {
      success: response.ok,
      data: response.ok ? await response.json() : null,
      status: response.status
    };
  }
};
```

### Meta-Builder Skill

```typescript
export const metaBuilderSkill = {
  name: 'register_openclaw_skill',
  description: 'Autonomously register a new OpenClaw skill',
  
  inputSchema: z.object({
    skillName: z.string(),
    skillType: z.enum(['webhook', 'supabase', 'sdk']),
    config: z.record(z.any())
  }),
  
  handler: async (input) => {
    // Use Qwen 3.5 to generate skill code
    const skillCode = await openclaw.invoke({
      model: 'qwen/qwen3.5-397b-a17b',
      messages: [{
        role: 'user',
        content: `Generate OpenClaw skill code for ${input.skillName} of type ${input.skillType}`
      }]
    });
    
    // Register the generated skill
    const skill = eval(skillCode); // In production, use safe parsing
    openclaw.registerSkill(skill);
    
    return { success: true, skillName: input.skillName };
  }
};
```

### Data Processing Skill

```typescript
export const dataProcessingSkill = {
  name: 'process_biometrics_data',
  description: 'Process biometric data with Qwen analysis',
  
  inputSchema: z.object({
    dataType: z.enum(['heart_rate', 'steps', 'sleep', 'workout']),
    rawData: z.string(),
    analysisType: z.enum(['anomaly', 'trend', 'insight'])
  }),
  
  handler: async (input) => {
    const prompt = `
      Analyze the following ${input.dataType} data:
      ${input.rawData}
      
      Provide ${input.analysisType} analysis with actionable insights.
    `;
    
    const result = await openclaw.invoke({
      model: 'qwen/qwen3.5-397b-a17b',
      messages: [{ role: 'user', content: prompt }],
      options: { max_tokens: 16384 }
    });
    
    return {
      analysis: result.choices[0].message.content,
      model: 'qwen/qwen3.5-397b-a17b',
      tokensUsed: result.usage.total_tokens
    };
  }
};
```

## Ultimate Goal

OpenClaw soll sich selbst replizieren durch autonomes Erstellen neuer Skills via Meta-Builder Protocol.

### Master-Skills für Self-Replication

- `deploy_n8n_workflow` - Deploy new n8n workflows
- `deploy_supabase_function` - Deploy Supabase Edge Functions
- `register_openclaw_skill` - Register new skills dynamically

## Best Practices 2026

1. **Immer Timeout erhöhen** - Qwen 3.5 braucht 120s
2. **Fallback-Chain definieren** - Nie nur ein Modell
3. **Retry mit Backoff** - 5s, 10s, 20s
4. **Rate-Limit beachten** - Max 40 RPM bei NVIDIA
5. **Health-Checks regelmäßig** - Monitore aktiv

## Troubleshooting

### Häufige Probleme

| Problem | Ursache | Lösung |
|---------|---------|--------|
| Timeout bei 90s | Modell zu groß | Timeout auf 120000ms erhöhen |
| HTTP 429 | Rate-Limit erreicht | 60s warten + Fallback |
| 401 Unauthorized | Token abgelaufen | Token refreshen |
| Langsame Antworten | Netzwerklatenz | CDN/Region prüfen |

### Diagnose-Befehle

```bash
# Check OpenClaw status
openclaw status

# Check model availability
openclaw models | grep nvidia

# Test connection
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
     https://integrate.api.nvidia.com/v1/models

# View recent logs
openclaw logs --last 50
```

## Abnahme-Check OPENCLAW

1. [ ] Konfigurationsmatrix vorhanden
2. [ ] Qwen 3.5 Integration dokumentiert
3. [ ] NVIDIA NIM Setup korrekt
4. [ ] Timeout auf 120000ms gesetzt
5. [ ] Fallback-Kette definiert
6. [ ] Authfluss dokumentiert
7. [ ] Retry-/Fehlerstrategie definiert
8. [ ] Sicherheitsregeln enthalten
9. [ ] Verifikationsplan vorhanden
10. [ ] 5 Skill-Beispiele dokumentiert

## Siehe auch

- `WORKFLOW.md` - Skill creation architecture
- `INTEGRATION.md` - Full system integration
- `AGENTS.md` - Agent orchestration rules
- NVIDIA NIM Documentation: https://integrate.api.nvidia.com/

---

**Version:** 2.0  
**Letzte Aktualisierung:** 2026-02-18  
**Qwen 3.5 Integration:** qwen/qwen3.5-397b-a17b  
**Status:** ACTIVE
