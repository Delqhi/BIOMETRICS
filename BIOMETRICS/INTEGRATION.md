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
