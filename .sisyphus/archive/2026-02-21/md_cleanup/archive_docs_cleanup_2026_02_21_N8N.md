# N8N.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Workflow-Automation folgt globaler Governance und Secrets-Disziplin.
- Trigger, Fehlerpfade und Retry-Verhalten müssen dokumentiert sein.
- Kritische Flows benötigen Incident- und Recovery-Referenzen.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Betriebsleitfaden für n8n-Workflows in produktionsnahen Umgebungen.

## Prinzipien
1. Workflows sind versioniert
2. Fehlerpfade sind explizit
3. Idempotenz wird berücksichtigt
4. Recovery ist dokumentiert

## Workflow-Katalog (Template)

| Workflow-ID | Zweck | Trigger | Kritikalität | Owner | Status |
|---|---|---|---|---|---|
| WF-001 | {PURPOSE} | {TRIGGER} | P1 | {OWNER_ROLE} | active |

## Trigger-Typen
- webhook
- schedule
- event-driven
- manual

## Input/Output-Contract
Für jeden Workflow dokumentieren:
1. Eingabefelder
2. Validierungsregeln
3. Ausgabefelder
4. Fehlercodes
5. Nebenwirkungen

## Fehlerbehandlung
- Retry bei transienten Fehlern
- Dead-letter/Quarantäne bei dauerhaften Fehlern
- Alarmierung bei kritischen Ausfällen

## Recovery-Plan
1. Fehlerfall erkennen
2. betroffenen Workflow pausieren
3. Ursache isolieren
4. Korrektur deployen
5. sicheren Neustart durchführen

## Observability
- Laufzeitmetriken pro Workflow
- Fehlerquote
- Durchsatz
- Queue-Länge

## NLM-Bezug
- NLM-generierte Asset-Aufgaben können n8n-triggerbar modelliert werden
- Freigabe nur nach Qualitätsmatrix
- Ergebnisprotokoll in `MEETING.md`

## Abnahme-Check N8N
1. Workflow-Katalog vorhanden
2. Input/Output-Contract je Workflow definiert
3. Fehler- und Recoverypfade dokumentiert
4. Observability-Basis enthalten

---

## Qwen 3.5 Workflows (NVIDIA NIM)

Dieses Projekt nutzt Qwen 3.5 (397B) via NVIDIA NIM für anspruchsvolle KI-Workflows.

### Provider-Konfiguration

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

### Workflow 1: AI-Powered Automation

**Use Case:** Intelligente Dokumentenverarbeitung und Entscheidungsfindung

```yaml
Workflow-ID: WF-QWEN-001
Trigger: Webhook (POST /webhook/qwen-automation)
Kritikalität: P1

Input:
- document_text: string
- processing_type: "classify" | "extract" | "summarize"
- context: object

Process:
1. HTTP Request → Document Empfang
2. Code Node → Text preprocessing
3. HTTP Request → Qwen 3.5 API (NVIDIA NIM)
4. Code Node → Response parsing
5. Switch → Route based on processing_type
6. HTTP Response → Clean JSON zurück

Output:
- result: string
- confidence: number
- metadata: object
```

### Workflow 2: Multimodal Processing

**Use Case:** Bildanalyse mit Qwen VL (Vision-Language)

```yaml
Workflow-ID: WF-QWEN-002
Trigger: Webhook (POST /webhook/qwen-vision)
Kritikalität: P1

Input:
- image_url: string
- analysis_type: "describe" | "extract_text" | "detect_objects"
- language: string (default: "de")

Process:
1. HTTP Request → Image URL empfangen
2. HTTP Request → Bild als Base64 laden
3. Code Node → Multimodal Prompt erstellen
4. HTTP Request → Qwen VL API
5. Code Node → Resultat parsen
6. HTTP Response → Strukturiertes Resultat

Output:
- description: string
- extracted_text: string (optional)
- objects: array (optional)
- language: string
```

### Workflow 3: Code Generation Workflow

**Use Case:** Automatische Code-Generierung für Next.js/Supabase

```yaml
Workflow-ID: WF-QWEN-003
Trigger: Webhook (POST /webhook/qwen-code)
Kritikalität: P1

Input:
- specification: string
- stack: "nextjs" | "supabase" | "fullstack"
- components: string[]

Process:
1. HTTP Request → Spezifikation empfangen
2. Code Node → Prompt für Code-Generierung bauen
3. HTTP Request → Qwen 3.5 API
4. Code Node → Code parsen und formatieren
5. IF stack = "supabase"
   - SQL für Supabase extrahieren
6. IF stack = "nextjs"
   - TypeScript/React Code extrahieren
7. HTTP Response → Code + Dateistruktur

Output:
- files: array of {filename, content, language}
- dependencies: string[]
- setup_instructions: string
```

### Workflow 4: Conversation AI

**Use Case:** Kontextbezogene Kundenkommunikation

```yaml
Workflow-ID: WF-QWEN-004
Trigger: Webhook (POST /webhook/qwen-chat)
Kritikalität: P2

Input:
- user_message: string
- conversation_history: array
- context: object (product_info, user_data)

Process:
1. HTTP Request → Nachricht empfangen
2. Code Node → Conversation Context aufbauen
3. HTTP Request → Qwen 3.5 API
4. Code Node → Response validieren
5. HTTP Response → Antwort + suggested_actions

Output:
- response: string
- suggested_actions: array
- sentiment: "positive" | "neutral" | "negative"
```

### Error Handling für Qwen Workflows

```yaml
Retry Strategy:
- transient_errors: 3 retries mit exponential backoff
- rate_limit (429): 60s warten, dann retry
- invalid_response: Fallback auf simpler prompt

Quarantine:
- after 3 failed attempts → to dead-letter queue
- Alerting: PagerDuty/Slack notification
```

### Observability

| Metric | Target |
|--------|--------|
| Response Time P95 | < 30s |
| Success Rate | > 95% |
| Rate Limit Errors | < 5% |

---

5. Observability-Basis enthalten

---

---

## 12) n8n als Heavy Lifting Muscle für AI Skills

n8n Workflows sind die "Automation Muscle" die von AI Skills getriggert werden im Webhook Wrapper Pattern.

**Integration Pattern:**
- OpenClaw Skill trigger n8n Webhook
- n8n führt multi-step Workflow aus
- Clean JSON Response für AI

**Design Principles:**
- Workflows müssen AI-triggbar sein (Webhook endpoint)
- Fehler müssen AI-freundlich zurückgegeben werden
- Idempotenz für wiederholbare Execution

**Meta-Builder Protocol:**
Der Agent kann autonom neue n8n Workflows erstellen und deployen via `deploy_n8n_workflow` Master-Skill.

**Siehe auch:** `WORKFLOW.md` für vollständige Architektur-Dokumentation.
