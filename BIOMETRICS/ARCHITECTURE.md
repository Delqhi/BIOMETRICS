# ARCHITECTURE.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Architekturentscheidungen folgen den globalen Regeln aus `AGENTS-GLOBAL.md`.
- Jede Strukturänderung benötigt Mapping- und Integrationsabgleich.
- Security-by-Default, NLM-First und Nachweisbarkeit sind Pflicht.

Status: ACTIVE  
Version: 1.1 (Qwen 3.5 Brain)  
Stand: Februar 2026

## Zweck
Universelle Architekturvorlage für modulare, skalierbare und auditierbare Systeme.

## Scope
- Frontend: Next.js
- Backend: Go + Supabase
- Integrationen: OpenClaw, n8n, Cloudflare, Vercel (optional)
- **Primary AI Brain:** Qwen 3.5 (NVIDIA NIM)

---

## 1) Systemkontext

### Nutzergruppen
- **Primary:** BIOMETRICS CLI Benutzer
- **Secondary:** KI-Agenten (OpenClaw, n8n)

### Hauptziele
- Vollautomatisierte KI-Dienstleistungen
- Self-Building AI Agent System
- NLM-gestützte Content-Generierung

### Externe Systeme
- Supabase (Datenbank, Auth, Storage)
- n8n (Workflow-Automatisierung)
- OpenClaw (Agenten-Orchestrierung)
- NVIDIA NIM (Qwen 3.5)
- Cloudflare (Netzwerksicherheit)

---

## 2) Architekturprinzipien
1. API-first
2. Modular statt monolithisch
3. Security by default
4. Observability by design
5. Reproduzierbarer Betrieb
6. **AI-Brain-first** (Qwen 3.5 als zentrale推理-Engine)

## 2) Architekturprinzipien
1. API-first
2. Modular statt monolithisch
3. Security by default
4. Observability by design
5. Reproduzierbarer Betrieb

## 3) Modulübersicht (Template)

| Modul | Verantwortung | Eingänge | Ausgänge | Abhängigkeiten |
|---|---|---|---|---|
| qwen-brain | AI推理 & Generierung | prompts, files, images | text, code, analysis | NVIDIA NIM API |
| web-frontend | UI und User Flows | user events | API calls | api-gateway |
| api-gateway | Zugriffsschicht | HTTP requests | responses | services |
| service-core | Business-Logik | API payloads | domain events | supabase |
| content-orchestrator | NLM-Delegation | content tasks | generated assets | nlm-cli |
| workflow-engine | Automationen | triggers | actions | n8n/openclaw |

### Datenfluss mit Qwen Brain

```
User Request → Frontend → API Gateway → Service Core
                                            ↓
                                    [Qwen 3.5 Brain]
                                            ↓
                                    NVIDIA NIM API
                                            ↓
                                    Response Processing
                                            ↓
                                    Supabase (persist)
                                            ↓
                                    User Response
```

## 4) Qwen 3.5 Brain - Modell-Architektur

### Modell-Details

| Attribut | Wert |
|----------|------|
| **Modell-ID** | `qwen/qwen3.5-397b-a17b` |
| **Provider** | NVIDIA NIM |
| **Context Window** | 262.144 Tokens |
| **Output Limit** | 32.768 Tokens |
| **Modalitäten** | Text, Vision, Code, Multimodal |
| **Latenz** | 70-90s (High-Latency Modus) |

### Verfügbare Skills

#### qwen_vision_analysis
- **Zweck:** Bildanalyse und visuelle Erkennung
- **Use Cases:** Produktbild-Qualitätsprüfung, Layout-Analyse
- **Input:** Bilder (PNG, JPG, WebP)
- **Output:** Strukturierte Analyse mit Tags und Metriken

#### qwen_code_generation
- **Zweck:** Full-Stack Code-Generierung
- **Use Cases:** Komponenten, API-Routen, Datenbank-Schema
- **Input:** Natürliche Sprache oder Spezifikation
- **Output:** Fertiger, getesteter Code

#### qwen_document_ocr
- **Zweck:** Texterkennung und Dokumentanalyse
- **Use Cases:** Rechnungsverarbeitung, Vertragsanalyse
- **Input:** PDF, Bilder mit Text
- **Output:** Extrahierter Text, Metadaten, Struktur

#### qwen_video_understanding
- **Zweck:** Video-Inhaltsanalyse
- **Use Cases:** Video-Vorschau, Content-Indexierung
- **Input:** Videos (MP4, MOV, WebM)
- **Output:** Szenenbeschreibung, Key-Frames, Metadaten

#### qwen_conversation
- **Zweck:** Natürliche Konversations-KI
- **Use Cases:** Support-Chat, Produktberatung
- **Input:** Benutzer-Nachrichten, Kontext
- **Output:** Kontextbezogene Antworten

### Capabilities Matrix

| Capability | Status | Integration |
|------------|--------|-------------|
| Text Generation | ✅ Active | API Call |
| Vision Analysis | ✅ Active | API Call + Image Input |
| Code Generation | ✅ Active | API Call |
| OCR/Document | ✅ Active | API Call |
| Video Understanding | ✅ Active | API Call |
| Multimodal | ✅ Active | Combined Skills |
| Tool Calling | ✅ Active | Function Definitions |
| Streaming | ❌ Not Supported | Polling fallback |
1. User interagiert mit Next.js Frontend.
2. Frontend ruft API-Endpunkte auf.
3. Go-Services validieren, orchestrieren und persistieren.
4. Supabase liefert DB/Auth/Storage.
5. Bei Content-Jobs delegiert der Agent via NLM-CLI.

---

## 5) Integrationsschnittstellen

### Qwen 3.5 Integration
- **Endpoint:** `https://integrate.api.nvidia.com/v1`
- **API:** OpenAI-Compatible (`openai-completions`)
- **Auth:** Bearer Token (`NVIDIA_API_KEY`)
- **Timeout:** 120.000ms (Critical for high-latency)
- **Rate Limit:** 40 RPM (Free Tier)

### Andere Integrationen
- OpenClaw: Integrationsauth und Connector-Flows
- n8n: Workflow-Orchestrierung
- Cloudflare: Netz-/Zugangsabsicherung
- Vercel: Frontend-Auslieferung (optional)

## 6) NLM-Architekturanker
- NLM-Artefakte: Video, Infografik, Präsentation, Datentabelle
- Erzeugung ausschließlich via NLM-CLI
- Übernahme nur nach Qualitätsmatrix
- Protokollierung in `MEETING.md`

## 7) Nicht-funktionale Anforderungen
- Performance: **Qwen 3.5 Latenz:** 70-90s, **Timeout:** 120s mandatory
- Verfügbarkeit: **Target:** 99.9%, **Fallback Chain:** Qwen → Kimi → Claude
- Sicherheit: API-Keys in Environment, Cloudflare WAF
- Wartbarkeit: kleine, klar getrennte Module, Qwen-Skills als separate Funktionen

## 8) Risiken

| Risiko | Wahrscheinlichkeit | Auswirkung | Mitigation |
|--------|------------------|------------|------------|
| Qwen 3.5 Latenz (70-90s) | Hoch | Mittel | Async Polling, User Feedback |
| NVIDIA Rate Limit (429) | Mittel | Hoch | Fallback Chain, Retry-Logic |
| API Key Rotation | Niedrig | Hoch | Automatisierte Rotation |

## 9) Entscheidungen (ADR-Index)

| ADR-ID | Thema | Entscheidung | Status |
|---|---|---|---|
| ADR-001 | Frontend-Stack | Next.js | accepted |
| ADR-002 | Backend-Stack | Go + Supabase | accepted |
| ADR-003 | Content-Generierung | NLM-CLI Pflicht | accepted |
| ADR-004 | Primary AI Brain | Qwen 3.5 (NVIDIA NIM) | accepted |
| ADR-005 | Qwen Timeout | 120s Mandatory | accepted |
| ADR-006 | Fallback Chain | Qwen → Kimi → Claude | accepted |

## 10) Verifikation
- Konsistenzcheck mit `COMMANDS.md` und `ENDPOINTS.md`
- Abgleich mit `SECURITY.md` und `INFRASTRUCTURE.md`
- Task-20 Integrationsreview

## Abnahme-Check ARCHITECTURE
1. ✅ Module und Verantwortungen klar
2. ✅ Datenflüsse nachvollziehbar
3. ✅ NLM-Integration verankert
4. ✅ Qwen 3.5 Brain dokumentiert
5. ✅ Risiken dokumentiert
6. ✅ ADR-Index gepflegt

---

## Anhang: Qwen 3.5 API-Referenz

### Request Format
```json
{
  "model": "qwen/qwen3.5-397b-a17b",
  "messages": [
    {
      "role": "user",
      "content": [
        { "type": "text", "text": "Analyze this image" },
        { "type": "image_url", "image_url": { "url": "data:image/jpeg;base64,..." } }
      ]
    }
  ],
  "temperature": 0.7,
  "max_tokens": 32768
}
```

### Error Handling
- **HTTP 429:** Wait 60s + Fallback to Kimi K2.5
- **Timeout:** Retry with exponential backoff
- **Invalid Request:** Log + Return user-friendly error

---

## 6) NLM-Architekturanker

---
