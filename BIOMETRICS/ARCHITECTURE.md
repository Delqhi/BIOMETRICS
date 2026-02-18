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

---

## 3) Modulübersicht

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

## 11) Executive Summary - Projektvision und Architekturprinzipien

### 11.1 Projektvision

BIOMETRICS ist ein vollautomatisiertes KI-Dienstleistungssystem, das moderne Cloud-Architektur mit fortschrittlicher künstlicher Intelligenz verbindet. Das System ermöglicht es Benutzern, über eine intuitive CLI-Schnittstelle komplexe KI-gestützte Aufgaben auszuführen, ohne tiefgreifendes technisches Wissen zu benötigen. Die Vision umfasst die vollständige Automatisierung von Geschäftsprozessen durch selbstbauende KI-Agenten, die kontinuierlich neue Fähigkeiten erlernen und integrierte Workflows optimieren können.

Das Kernprinzip basiert auf dem Zusammenspiel mehrerer Technologieebenen: einer robusten Go-basierten Backend-Architektur für Performanz und Zuverlässigkeit, einemNext.js-Frontend für optimale Benutzererfahrung, und einem intelligenten KI-Orchestrierungssystem, das verschiedene AI-Modelle nahtlos integriert. Der Primary AI Brain basiert auf Qwen 3.5 (NVIDIA NIM) und ermöglicht fortschrittliche Reasoning-Fähigkeiten für komplexe Aufgaben wie Code-Generierung, Bildanalyse, Dokumentenverarbeitung und multimodale Interaktionen.

Die Architektur folgt dem Prinzip der serviceorientierten Modularität, wobei jedes Modul unabhängig skaliert, entwickelt und deployed werden kann. Dies gewährleistet maximale Flexibilität bei der Systemweiterentwicklung und ermöglicht es, neue KI-Fähigkeiten zu integrieren, ohne bestehende Funktionalität zu beeinträchtigen. Die Trennung von Frontend, Backend und KI-Orchestrierungsschicht bildet das Fundament für eine zukunftssichere Architektur.

### 11.2 Architekturprinzipien im Detail

Die Architektur von BIOMETRICS basiert auf sechs fundamentalen Prinzipien, die jede Designentscheidung leiten:

**Prinzip 1: API-First Design**
Alle Systemkomponenten kommunizieren primär über definierte API-Schnittstellen. Dies gewährleistet klare Verträge zwischen Modulen, ermöglicht unabhängige Entwicklung und Testing, und schafft die Grundlage für zukünftige Integrationen. Jedes Modul expose seine Funktionalität über gut dokumentierte HTTP- oder gRPC-Endpunkte, die konsistente Antwortformate und Fehlerbehandlung bieten.

Das API-First-Prinzip manifestiert sich in mehreren Architekturschichten: Der API-Gateway fungiert als zentrale Eingangsschicht und routet Anfragen an passende Backend-Services. Jeder Service implementiert seine eigene API mit spezifischen Contract-Definitionen, die in OpenAPI/Swagger dokumentiert sind. Die Konsistenz der API-Versionierung wird durch automatisierte Tests sichergestellt, die Breaking Changes frühzeitig erkennen.

```go
// API Contract Beispiel für User Service
type UserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Password string `json:"password" validate:"required,min=8"`
}

type UserResponse struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type ErrorResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}
```

**Prinzip 2: Modularität über Monolithik**
Das System ist in klar abgegrenzte Module unterteilt, die jeweils eine spezifische Verantwortung haben. Diese Module können unabhängig voneinander entwickelt, getestet und deployed werden. Die Modularität erstreckt sich über alle Ebenen der Architektur: von der Package-Struktur im Code bis hin zu separaten Docker-Containern und Kubernetes-Pods.

Die Vorteile dieses Ansatzes manifestieren sich in verschiedenen Aspekten der Systemwartung und -entwicklung. Neue Funktionen können als eigenständige Module hinzugefügt werden, ohne bestehende Komponenten zu modifizieren. Fehlerisolierung stellt sicher, dass Ausfälle in einem Modul nicht automatisch das gesamte System beeinträchtigen. Teamarbeit wird durch die Möglichkeit paralleler Entwicklung verschiedener Module optimiert.

**Prinzip 3: Security by Default**
Sicherheit ist kein nachträglich hinzugefügtes Feature, sondern integraler Bestandteil jeder Architekturschicht. Von der Authentifizierung über Autorisierung bis hin zur Datenverschlüsselung folgen alle Komponenten dem Zero-Trust-Prinzip. Jede Anfrage wird authentifiziert und autorisiert, unabhängig davon, ob sie von internen oder externen Quellen stammt.

Die Implementierung von Security by Default umfasst mehrere Schutzebenen: Eingabevalidierung auf allen Ebenen verhindert Injection-Angriffe. Rollenbasierte Zugriffskontrollen (RBAC) gewährleisten, dass Benutzer nur auf autorisierte Ressourcen zugreifen können. Transport Layer Security (TLS) verschlüsselt alle Datenübertragungen. Secrets werden niemals im Code gespeichert, sondern ausschließlich in spezialisierten Vault-Lösungen verwaltet.

```go
// Security Middleware Beispiel
func SecurityMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // CORS Policy
        w.Header().Set("Access-Control-Allow-Origin", config.AllowedOrigins)
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
        
        // Security Headers
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        
        // Rate Limiting
        if !rateLimiter.Allow(r.RemoteAddr) {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

**Prinzip 4: Observability by Design**
Das System generiert umfassende Telemetriedaten, die Echtzeitüberwachung, Troubleshooting und Kapazitätsplanung ermöglichen. Jedes Modul implementiert standardisierte Logging-, Metrics- und Tracing-Schnittstellen. Die gesammelten Daten werden in zentralen Observability-Plattformen aggregiert und analysiert, um aussagekräftige Einblicke in das Systemverhalten zu gewähren.

Logging folgt dem strukturierten JSON-Format mit kontextuellen Feldern für Korrelation und Filterung. Metriken erfassen quantitative Aspekte wie Request-Latenzen, Fehlerraten und Ressourcennutzung. Distributed Tracing ermöglicht die Nachverfolgung von Anfragen über Servicegrenzen hinweg. Alle Observability-Daten werden mit Trace-IDs korreliert, um vollständige Request-Pfade zu rekonstruieren.

**Prinzip 5: Reproduzierbarer Betrieb**
Jede Änderung am System durchläuft einen standardisierten Deployment-Prozess, der vollständige Reproduzierbarkeit gewährleistet. Infrastructure as Code (IaC) stellt sicher, dass Umgebungen konsistent aufgesetzt werden können. Containerisierung isoliert Anwendungen und eliminiert Umgebungsabhängigkeiten. Konfigurationsmanagement ermöglicht die kontrollierte Verteilung von Einstellungen über verschiedene Umgebungen hinweg.

Die Implementierung umfasst versionskontrollierte Infrastrukturdefinitionen, automatisierte Build- und Test-Pipelines, kontrollierte Deployment-Prozesse mit Approval-Workflows, und umfassende Rollback-Mechanismen für den Fehlerfall. Jede Komponente ist so konzipiert, dass sie in jeder Umgebung identisch funktioniert.

**Prinzip 6: AI-Brain-First**
Die Integration von Künstlicher Intelligenz ist das zentrale Differenzierungsmerkmal von BIOMETRICS. Das System nutzt Qwen 3.5 (NVIDIA NIM) als primäre Reasoning-Engine für komplexe Aufgaben, unterstützt durch spezialisierte KI-Modelle für verschiedene Modalitäten. Die AI-Orchestrierungsschicht ermöglicht nahtlose Interaktion mit verschiedenen KI-Providern und optimiert automatisch die Modellauswahl basierend auf Aufgabenanforderungen.

Die AI-First-Philosophie durchdringt alle Architekturschichten: Das Frontend integriert KI-Chat-Interaktionen, das Backend delegiert komplexe Verarbeitungsaufgaben an KI-Modelle, und die Workflow-Engine automatisiert KI-gestützte Entscheidungsprozesse. Diese Integration ermöglicht selbstlernende Geschäftsprozesse, die kontinuierlich durch KI-Analysen optimiert werden.

### 11.3 Technologie-Stack-Entscheidungen

Die Wahl des Technologie-Stacks basiert auf mehreren Kernkriterien: Performanz, Skalierbarkeit, Entwicklerproduktivität und Ökosystem-Reife. Jede Entscheidung ist durch konkrete Anforderungen und Trade-off-Analysen fundiert.

**Frontend: Next.js**
Next.js bietet das optimale Gleichgewicht zwischen Entwicklerproduktivität und Performanz. Das Framework ermöglicht serverseitiges Rendering für SEO-Optimierung und schnelle First-Contentful-Paint-Zeiten, während Client-seitige Interaktionen für dynamische User-Experiences sorgen. Das App-Router-Architekturmuster fördert modulare Codeorganisation. Das große Ökosystem an Plugins und die aktive Community gewährleisten langfristige Wartbarkeit.

Die Entscheidung für Next.js basiert auf folgenden Faktoren: Native Unterstützung für React 18+ Features, eingebaute Bildoptimierung und Font-Loading, API-Routen für Backend-Funktionalität, und umfangreiche Deployment-Optionen (Vercel, Docker, etc.). Die TypeScript-Unterstützung ermöglicht typsichere Entwicklung mit hervorragender IDE-Integration.

**Backend: Go + Supabase**
Go wurde als primäre Backend-Sprache gewählt aufgrund seiner herausragenden Performanz, einfachen Parallelisierung und geringen Ressourcennutzung. Die statische Typisierung und kompakte Syntax fördern wartbaren Code. Das umfangreiche Standardlibrary reduziert externe Abhängigkeiten. Go-Programme kompilieren zu einzelnen Binärdateien, was Deployment und Distribution vereinfacht.

Supabase ergänzt Go als Datenbank- und Authentifizierungsschicht mit PostgreSQL-Kompatibilität, eingebautem Auth-System, Realtime-Subscriptions und Row Level Security. Die Kombination ermöglicht schnelle Prototypenentwicklung bei gleichzeitiger Produktionsreife. Supabase Edge Functions ermöglichen serverless Computing für burstfähige Workloads.

**Primary AI: Qwen 3.5 (NVIDIA NIM)**
Qwen 3.5 wurde als Primary AI Brain ausgewählt aufgrund seiner herausragenden Reasoning-Fähigkeiten, großen Context-Window (262K Tokens) und multimodaler Unterstützung. NVIDIA NIM bietet optimierte Inference-Performanz und einfache Integration über OpenAI-kompatible APIs. Die 40 RPM Rate Limit im Free Tier ist für die meisten Anwendungsfälle ausreichend, mit Fallback-Optionen bei Bedarf.

Die Integration erfolgt über OpenAI-kompatible Endpoints, was eine einfache Abstraktionsschicht ermöglicht. Für spezielle Anwendungsfälle werden zusätzliche Modelle wie Kimi K2.5 für generale Aufgaben oder spezialisierte OCR/ Vision-Modelle integriert.

---

## 12) System Overview - High-Level Architektur und Datenfluss

### 12.1 High-Level Architekturdiagramm

Die Architektur von BIOMETRICS folgt einem mehrstufigen Schichtenmodell, das klare Trennung von Verantwortlichkeiten gewährleistet. Das folgende Diagramm illustriert die Hauptkomponenten und ihre Interaktionen:

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              PRESENTATION LAYER                                  │
│  ┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────────┐   │
│  │   Next.js Frontend  │  │    CLI Interface    │  │   Admin Dashboard   │   │
│  │   (Web Application) │  │  (Terminal Client)  │  │   (Management)     │   │
│  └──────────┬──────────┘  └──────────┬──────────┘  └──────────┬──────────┘   │
└─────────────┼─────────────────────────┼─────────────────────────┼──────────────┘
              │                         │                         │
              └─────────────────────────┼─────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              API GATEWAY LAYER                                   │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │                         API Gateway (Go)                                  │  │
│  │  • Rate Limiting    • Authentication    • Request Routing              │  │
│  │  • Load Balancing    • SSL Termination   • API Versioning              │  │
│  └──────────────────────────────────┬───────────────────────────────────────┘  │
└─────────────────────────────────────┼──────────────────────────────────────────┘
                                      │
              ┌───────────────────────┬─┴───────────────────────┬─────────────────┐
              │                       │                       │                 │
              ▼                       ▼                       ▼                 ▼
┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────────┐ ┌─────┐
│   Auth Service      │  │   Core Services     │  │   AI Orchestrator   │ │ n8n │
│   (Go + JWT)       │  │   (Go + Business)   │  │   (AI Brain)        │ │     │
│                     │  │                     │  │                     │ │     │
│ • User Management   │  │ • User Service      │  │ • Qwen 3.5          │ │     │
│ • Role Management  │  │ • Content Service   │  │ • Kimi K2.5         │ │     │
│ • Token Validation │  │ • Integration Svc   │  │ • Fallback Chain    │ │     │
│ • Session Mgmt     │  │ • Workflow Service  │  │ • Prompt Engineering│ │     │
└─────────┬───────────┘  └──────────┬──────────┘  └──────────┬──────────┘ └──┬──┘
          │                         │                       │               │
          └─────────────────────────┼───────────────────────┼───────────────┘
                                    │                       │
                                    ▼                       ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              DATA LAYER                                          │
│  ┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────────┐  │
│  │   PostgreSQL        │  │      Redis           │  │      S3/Storage     │  │
│  │   (Supabase)        │  │   (Cache/Session)   │  │   (Media Assets)    │  │
│  │                     │  │                     │  │                     │  │
│  │ • User Data         │  │ • Session Cache     │  │ • Generated Assets  │  │
│  │ • Business Data     │  │ • API Cache         │  │ • User Uploads     │  │
│  │ • Audit Logs       │  │ • Rate Limit Data  │  │ • NLM Output       │  │
│  └─────────────────────┘  └─────────────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 12.2 Component Map und Verantwortlichkeiten

Die folgende Tabelle dokumentiert alle Hauptkomponenten mit ihren spezifischen Verantwortlichkeiten und Abhängigkeiten:

| Komponente | Verantwortlichkeit | Technologie | Abhängigkeiten | SLA |
|------------|------------------|-------------|----------------|-----|
| API Gateway | Request-Routing, Rate Limiting, SSL | Go + Nginx | Alle Backend-Services | 99.9% |
| Auth Service | Benutzer-Authentifizierung, Autorisierung | Go + JWT | PostgreSQL, Redis | 99.95% |
| User Service | Benutzerverwaltung, Profile | Go + gRPC | PostgreSQL, Auth Service | 99.9% |
| Content Service | Content-Erstellung, -Speicherung | Go + Supabase | PostgreSQL, S3, AI Orchestrator | 99.9% |
| Integration Service | Externe API-Verbindungen | Go | Various External APIs | 99.5% |
| Workflow Service | Workflow-Automatisierung | Go + n8n | n8n, PostgreSQL | 99.5% |
| AI Orchestrator | KI-Anfragen, Modellauswahl | Go + Python | NVIDIA NIM, Kimi API | 99.0% |
| Qwen Brain | Reasoning, Code-Gen, Vision | NVIDIA NIM | - | 99.0% |
| PostgreSQL | Primärdatenspeicherung | Supabase/PostgreSQL | - | 99.99% |
| Redis | Caching, Sessions | Redis | - | 99.9% |
| S3 Storage | Blob Storage | Supabase Storage | - | 99.99% |

### 12.3 Datenfluss-Architektur

Der Datenfluss durch das System folgt einem vorhersagbaren Muster, das Monitoring und Troubleshooting vereinfacht. Jede Anfrage durchläuft mehrere Verarbeitungsschichten, wobei jede Schicht spezifische Verantwortlichkeiten hat.

**Request-Response Zyklus:**

```
1. CLIENT REQUEST
   │
   ▼
2. API GATEWAY (Go)
   ├── Validate Request Schema
   ├── Check Authentication Token
   ├── Apply Rate Limiting
   ├── Route to Backend Service
   │
   ▼
3. AUTH SERVICE (Go)
   ├── Validate JWT Signature
   ├── Extract User Claims
   ├── Check Authorization Rules
   ├── Return Auth Context
   │
   ▼
4. BUSINESS SERVICE (Go)
   ├── Validate Business Rules
   ├── Apply Input Sanitization
   ├── Execute Business Logic
   ├── Emit Domain Events
   │
   ▼
5. DATA LAYER
   ├── PostgreSQL: Persistent Storage
   ├── Redis: Cache Read/Write
   ├── S3: File Storage
   │
   ▼
6. AI ORCHESTRATOR (Go + Python)
   ├── Select Optimal AI Model
   ├── Prepare Prompt
   ├── Execute AI Request
   ├── Process Response
   │
   ▼
7. RESPONSE FORMATION
   ├── Serialize Response
   ├── Apply Caching
   ├── Add Monitoring Headers
   │
   ▼
8. API GATEWAY
   ├── Apply Response Transformations
   ├── Add Security Headers
   ├── Return to Client
```

### 12.4 Event-Driven Architektur

Für asynchrone Verarbeitungsaufgaben implementiert das System eine Event-Driven Architektur, die lose Kopplung und hohe Skalierbarkeit ermöglicht:

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Producer    │────▶│   Message    │────▶│  Consumer    │
│  Service     │     │   Broker     │     │  Service     │
└──────────────┘     │   (Redis)    │     └──────────────┘
                     └──────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │   Dead       │
                    │   Letter     │
                    │   Queue      │
                    └──────────────┘
```

**Event Types:**

| Event | Producer | Consumer | Payload |
|-------|----------|----------|---------|
| user.created | Auth Service | Email Service, Analytics | UserID, Email, Name |
| content.generated | AI Orchestrator | Storage Service, Notification | ContentID, Type, URL |
| workflow.triggered | API Gateway | Workflow Engine | WorkflowID, TriggerData |
| payment.completed | Payment Gateway | Order Service, Inventory | OrderID, Amount, Status |
| integration.webhook | External API | Processing Service | WebhookEvent, Payload |

```go
// Event Producer Example
type EventProducer struct {
    redis *redis.Client
}

func (ep *EventProducer) Publish(ctx context.Context, eventType string, payload interface{}) error {
    event := Event{
        Type:      eventType,
        Payload:   payload,
        Timestamp: time.Now().UTC(),
        TraceID:   GetTraceID(ctx),
    }
    
    data, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %w", err)
    }
    
    return ep.redis.Publish(ctx, eventType, data).Err()
}

// Event Consumer Example
type EventConsumer struct {
    redis       *redis.Client
    handlers    map[string]EventHandler
}

func (ec *EventConsumer) Start(ctx context.Context) error {
    pubsub := ec.redis.Subscribe(ctx, "events.*")
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case msg := <-pubsub.Channel():
            var event Event
            if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
                log.Error("failed to unmarshal event", "error", err)
                continue
            }
            
            if handler, ok := ec.handlers[event.Type]; ok {
                if err := handler.Handle(ctx, event); err != nil {
                    log.Error("event handler failed", "type", event.Type, "error", err)
                }
            }
        }
    }
}
```

### 12.5 Externe Systemintegrationen

BIOMETRICS integriert verschiedene externe Systeme über standardisierte Schnittstellen. Die folgende Übersicht dokumentiert alle Integrationen mit ihren spezifischen Konfigurationen:

**NVIDIA NIM (Qwen 3.5):**

```
Endpoint: https://integrate.api.nvidia.com/v1
Auth: Bearer Token (NVIDIA_API_KEY)
Timeout: 120000ms
Rate Limit: 40 RPM
Models:
  - qwen/qwen3.5-397b-a17b (Primary)
  - qwen2.5-coder-32b (Fast fallback)
  - moonshotai/kimi-k2.5 (General fallback)
```

**Supabase:**

```
PostgreSQL: postgresql://user:pass@host:5432/db
Auth: JWT + API Keys
Storage: S3-compatible API
Realtime: WebSocket subscriptions
Edge Functions: Deno runtime
```

**n8n Workflow Engine:**

```
API: http://n8n:5678/api/v1
Auth: Basic Auth / OAuth
Webhooks: /webhook/{workflow-id}
Executions: /executions/{id}
```

**Cloudflare:**

```
WAF: Rulesets für IP, Country, Bot
CDN: Global asset distribution
Workers: Edge computing
Tunnel: Secure tunnel to services
```

---

## 13) Modular Architecture - Go Module Struktur

### 13.1 Go Module Organisation

Die modulare Go-Architektur folgt bewährten Prinzipien für große Codebasen. Die Struktur fördert klare Abgrenzungen zwischen Modulen, minimiert zirkuläre Abhängigkeiten und ermöglicht unabhängige Entwicklung und Testing.

**Verzeichnisstruktur:**

```
biometrics/
├── cmd/                          # Entry Points (executables)
│   ├── api/                      # Main API Server
│   │   ├── main.go
│   │   └── main_test.go
│   ├── worker/                   # Background Worker
│   │   ├── main.go
│   │   └── worker_test.go
│   ├── cli/                      # CLI Application
│   │   ├── main.go
│   │   └── commands/
│   └── migrate/                  # Database Migrations
│       ├── main.go
│       └── migrations/
│
├── internal/                     # Private Application Code (nicht importierbar)
│   ├── auth/                    # Authentication & Authorization
│   │   ├── service.go
│   │   ├── middleware.go
│   │   ├── token.go
│   │   ├── oauth.go
│   │   └── auth_test.go
│   │
│   ├── database/                # Database Layer
│   │   ├── connection.go       # PostgreSQL connection
│   │   ├── repository.go       # Repository pattern
│   │   ├── migrations/         # Migration files
│   │   └── seeds/              # Seed data
│   │
│   ├── cache/                   # Caching Layer
│   │   ├── redis.go            # Redis client
│   │   ├── cache.go            # Cache interface
│   │   └── strategies.go       # Caching strategies
│   │
│   ├── api/                     # HTTP API Layer
│   │   ├── router.go           # Main router
│   │   ├── handlers/           # HTTP Handlers
│   │   │   ├── user.go
│   │   │   ├── content.go
│   │   │   └── integration.go
│   │   ├── middleware/         # HTTP Middleware
│   │   │   ├── logging.go
│   │   │   ├── recovery.go
│   │   │   └── cors.go
│   │   └── validators/         # Request validation
│   │       ├── user_validator.go
│   │       └── content_validator.go
│   │
│   ├── workers/                 # Background Workers
│   │   ├── processor.go        # Main processor
│   │   ├── tasks/              # Task definitions
│   │   │   ├── email.go
│   │   │   ├── notification.go
│   │   │   └── cleanup.go
│   │   └── scheduler.go        # Task scheduler
│   │
│   ├── config/                  # Configuration Management
│   │   ├── config.go           # Config struct
│   │   ├── env.go             # Environment variables
│   │   └── secrets.go         # Secret management
│   │
│   ├── events/                  # Event System
│   │   ├── producer.go        # Event producer
│   │   ├── consumer.go        # Event consumer
│   │   └── handlers.go       # Event handlers
│   │
│   └── services/               # Business Logic Services
│       ├── user_service.go
│       ├── content_service.go
│       ├── integration_service.go
│       └── workflow_service.go
│
├── pkg/                         # Public Libraries (importierbar)
│   ├── models/                 # Domain Models
│   │   ├── user.go
│   │   ├── content.go
│   │   ├── workflow.go
│   │   └── integration.go
│   │
│   ├── utils/                  # Utility Functions
│   │   ├── crypto.go          # Cryptography helpers
│   │   ├── validation.go      # Validation helpers
│   │   ├── conversion.go      # Type conversion
│   │   └── time.go            # Time utilities
│   │
│   ├── middleware/             # Reusable Middleware
│   │   ├── auth.go
│   │   ├── logging.go
│   │   └── ratelimit.go
│   │
│   └── errors/                # Error Handling
│       ├── errors.go
│       ├── codes.go
│       └── handling.go
│
├── api/                        # OpenAPI/Swagger Definitions
│   ├── openapi.yaml
│   ├── generated/             # Generated code
│   └── docs/                  # API Documentation
│
├── deployments/               # Deployment Configurations
│   ├── docker/               # Dockerfiles
│   │   ├── api.dockerfile
│   │   ├── worker.dockerfile
│   │   └── nginx.dockerfile
│   │
│   ├── k8s/                   # Kubernetes manifests
│   │   ├── base/
│   │   ├── overlays/
│   │   └── components/
│   │
│   └── terraform/             # Terraform modules
│       ├── main.tf
│       ├── variables.tf
│       └── outputs.tf
│
├── scripts/                   # Build & Deployment Scripts
│   ├── build.sh
│   ├── deploy.sh
│   └── migrate.sh
│
├── test/                      # Test Utilities
│   ├── fixtures/
│   ├── mocks/
│   └── testutil/
│
├── go.mod                     # Go module definition
├── go.sum                     # Go checksums
├── Makefile                  # Build targets
└── README.md                 # Project documentation
```

### 13.2 Package Organisation und Abhängigkeiten

Die Package-Struktur in Go folgt dem Prinzip der maximalen Kohäsion und minimalen Kopplung. Jedes Package hat eine klar definierte Verantwortung und暴露t nur notwendige Funktionen und Typen.

**Abhängigkeitsgraph (vereinfacht):**

```
cmd/api
    └── internal/api
            ├── internal/services
            ├── internal/auth
            ├── internal/database
            ├── internal/cache
            ├── internal/config
            ├── pkg/models
            ├── pkg/utils
            └── pkg/errors

cmd/worker
    └── internal/workers
            ├── internal/database
            ├── internal/cache
            ├── internal/events
            ├── internal/config
            └── pkg/models
```

**Package-Design Regeln:**

1. **Ein 入, 出 pro Package:** Jedes Package sollte klare Einstiegs- und Austrittspunkte haben
2. **Keine zirkulären Abhängigkeiten:** Import-Graph muss ein Directed Acyclic Graph (DAG) sein
3. **Fremdsprache: Englisch:** Alle Paket-, Funktions- und Variablennamen auf Englisch
4. **Kleine, fokussierte Packages:** Ein Package sollte < 10 Dateien und < 1000 Zeilen haben

### 13.3 Interface Design

Interfaces in Go definieren Verträge zwischen Komponenten und ermöglichen flexible Implementierungen. Das folgende Beispiel zeigt das Interface-Design für die verschiedenen Service-Schichten:

```go
// Repository Interfaces - Data Access Layer
type UserRepository interface {
    Create(ctx context.Context, user *User) (*User, error)
    GetByID(ctx context.Context, id string) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) (*User, error)
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter UserFilter) ([]*User, error)
    Count(ctx context.Context, filter UserFilter) (int64, error)
}

type ContentRepository interface {
    Create(ctx context.Context, content *Content) (*Content, error)
    GetByID(ctx context.Context, id string) (*Content, error)
    Update(ctx context.Context, content *Content) (*Content, error)
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter ContentFilter) ([]*Content, error)
    Search(ctx context.Context, query string) ([]*Content, error)
}

// Service Interfaces - Business Logic Layer
type AuthService interface {
    Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
    Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
    Logout(ctx context.Context, token string) error
    RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
    ValidateToken(ctx context.Context, token string) (*Claims, error)
}

type UserService interface {
    GetProfile(ctx context.Context, userID string) (*UserProfile, error)
    UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*UserProfile, error)
    ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error
    ListUsers(ctx context.Context, filter UserFilter) ([]*User, error)
}

type ContentService interface {
    CreateContent(ctx context.Context, req CreateContentRequest) (*Content, error)
    GetContent(ctx context.Context, id string) (*Content, error)
    UpdateContent(ctx context.Context, id string, req UpdateContentRequest) (*Content, error)
    DeleteContent(ctx context.Context, id string) error
    PublishContent(ctx context.Context, id string) error
    ListContent(ctx context.Context, filter ContentFilter) ([]*Content, error)
    GenerateWithAI(ctx context.Context, req AIGenerationRequest) (*Content, error)
}

// Cache Interfaces
type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
}

type SessionStore interface {
    Create(ctx context.Context, session *Session) error
    Get(ctx context.Context, sessionID string) (*Session, error)
    Update(ctx context.Context, session *Session) error
    Delete(ctx context.Context, sessionID string) error
    Expire(ctx context.Context, sessionID string, ttl time.Duration) error
}
```

### 13.4 Dependency Injection

Dependency Injection (DI) ermöglicht lockere Kopplung und erleichtert Testing. BIOMETRICS nutzt Constructor Injection für Abhängigkeiten:

```go
// Service Construction mit Dependency Injection
type UserService struct {
    userRepo    UserRepository
    cache       Cache
    eventProd   EventProducer
    logger      Logger
    config      *Config
}

func NewUserService(
    userRepo UserRepository,
    cache Cache,
    eventProd EventProducer,
    logger Logger,
    config *Config,
) *UserService {
    return &UserService{
        userRepo:  userRepo,
        cache:     cache,
        eventProd: eventProd,
        logger:    logger,
        config:    config,
    }
}

// Mit Functional Options für flexible Konfiguration
type UserServiceOption func(*UserService)

func WithCache(cache Cache) UserServiceOption {
    return func(s *UserService) {
        s.cache = cache
    }
}

func WithLogger(logger Logger) UserServiceOption {
    return func(s *UserService) {
        s.logger = logger
    }
}

// Usage
userService := NewUserService(
    userRepo,
    redisCache,
    eventProducer,
    logger,
    config,
    WithCache(redisCache),
    WithLogger(logger),
)
```

**Wire - Code Generation für Dependency Injection:**

```go
// +build wireinject

package main

import (
    "github.com/google/wire"
    "biometrics/internal/auth"
    "biometrics/internal/database"
    "biometrics/internal/cache"
    "biometrics/internal/services"
)

func InitializeApp(cfg *config.Config) (*App, error) {
    wire.Build(
        database.NewPostgreSQLConnection,
        cache.NewRedisClient,
        auth.NewAuthService,
        services.NewUserService,
        services.NewContentService,
        NewApp,
    )
    return nil, nil
}
```

### 13.5 Konfigurationsmanagement

Die Konfigurationsverwaltung in BIOMETRICS folgt dem 12-Factor App Prinzip und unterstützt verschiedene Umgebungen:

```go
// config/config.go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    Auth     AuthConfig     `mapstructure:"auth"`
    AI       AIConfig       `mapstructure:"ai"`
    Logging  LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
    Mode string `mapstructure:"mode"` // development, production
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    User     string `mapstructure:"user"`
    Password string `mapstructure:"password"`
    Database string `mapstructure:"database"`
    SSLMode  string `mapstructure:"ssl_mode"`
}

type AIConfig struct {
    Provider     string `mapstructure:"provider"`
    Model        string `mapstructure:"model"`
    APIKey       string `mapstructure:"api_key"`
    Endpoint     string `mapstructure:"endpoint"`
    Timeout      int    `mapstructure:"timeout"` // milliseconds
    MaxRetries   int    `mapstructure:"max_retries"`
    FallbackChain []string `mapstructure:"fallback_chain"`
}

// Konfiguration laden
func LoadConfig(path string) (*Config, error) {
    viper.SetConfigFile(path)
    viper.SetEnvPrefix("BIOMETRICS")
    viper.AutomaticEnv()
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }
    
    return &cfg, nil
}
```

**Environment-Variablen:**

```bash
# .env Datei
BIOMETRICS_SERVER_HOST=0.0.0.0
BIOMETRICS_SERVER_PORT=8080
BIOMETRICS_SERVER_MODE=development

BIOMETRICS_DATABASE_HOST=localhost
BIOMETRICS_DATABASE_PORT=5432
BIOMETRICS_DATABASE_USER=biometrics
BIOMETRICS_DATABASE_PASSWORD=secret
BIOMETRICS_DATABASE_NAME=biometrics_prod
BIOMETRICS_DATABASE_SSL_MODE=require

BIOMETRICS_REDIS_HOST=localhost
BIOMETRICS_REDIS_PORT=6379
BIOMETRICS_REDIS_PASSWORD=
BIOMETRICS_REDIS_DB=0

BIOMETRICS_AUTH_JWT_SECRET=your-jwt-secret-key
BIOMETRICS_AUTH_JWT_EXPIRY=24h
BIOMETRICS_AUTH_REFRESH_TOKEN_EXPIRY=7d

BIOMETRICS_AI_PROVIDER=nvidia-nim
BIOMETRICS_AI_MODEL=qwen/qwen3.5-397b-a17b
BIOMETRICS_AI_API_KEY=nvapi-xxxxx
BIOMETRICS_AI_ENDPOINT=https://integrate.api.nvidia.com/v1
BIOMETRICS_AI_TIMEOUT=120000
```

---

## 14) Microservices Design - Service Boundaries und Kommunikation

### 14.1 Service Boundary Definition

Die Microservices-Architektur von BIOMETRICS definiert klare Service-Grenzen basierend auf Geschäftsdomänen und funktionalen Verantwortlichkeiten. Jeder Service kapselt seine Daten und Geschäftslogik, während er über gut definierte APIs kommuniziert.

**Service-Übersicht:**

| Service | Verantwortung | Daten | Kommunikation |
|---------|-------------|-------|---------------|
| api-gateway | Request-Routing, Auth, Rate Limiting | Keine | HTTP/gRPC |
| auth-service | Authentifizierung, Autorisierung | Users, Roles, Sessions | HTTP/gRPC |
| user-service | Benutzerverwaltung | User Profiles, Preferences | HTTP/gRPC |
| content-service | Content-Erstellung, -Speicherung | Contents, Assets | HTTP/gRPC + Events |
| integration-service | Externe API-Verbindungen | API Credentials, Logs | HTTP |
| workflow-service | Workflow-Orchestrierung | Workflows, Executions | HTTP + Events |
| ai-orchestrator | KI-Anfragen-Routing | AI Requests, Responses | HTTP + Async |

### 14.2 API Gateway Pattern

Das API Gateway fungiert als zentrale Eingangsschicht für alle Client-Anfragen und implementiert cross-cutting concerns:

```go
// api-gateway/main.go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
    "biometrics/internal/api/middleware"
    "biometrics/internal/api/handlers"
    "biometrics/internal/auth"
    "biometrics/internal/config"
)

type APIGateway struct {
    router        *gin.Engine
    server        *http.Server
    config        *config.Config
    authClient    auth.Client
    redisClient   *redis.Client
}

func NewAPIGateway(cfg *config.Config) (*APIGateway, error) {
    gin.SetMode(gin.ReleaseMode)
    router := gin.New()
    
    // Middleware Chain
    router.Use(gin.Recovery())
    router.Use(middleware.Logger())
    router.Use(middleware.SecurityHeaders())
    router.Use(middleware.RequestID())
    router.Use(middleware.Timeout(30 * time.Second))
    
    // Initialize dependencies
    authClient := auth.NewClient(cfg.Auth)
    redisClient := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
        Password: cfg.Redis.Password,
        DB:       cfg.Redis.DB,
    })
    
    gateway := &APIGateway{
        router:      router,
        config:       cfg,
        authClient:   authClient,
        redisClient:  redisClient,
    }
    
    gateway.setupRoutes()
    return gateway, nil
}

func (g *APIGateway) setupRoutes() {
    // Health check
    g.router.GET("/health", handlers.HealthCheck())
    g.router.GET("/ready", handlers.ReadinessCheck())
    
    // Rate limited routes
    api := g.router.Group("/api/v1")
    api.Use(middleware.RateLimiter(g.redisClient))
    api.Use(middleware.Auth(g.authClient))
    {
        // User routes
        api.GET("/users", handlers.ListUsers())
        api.GET("/users/:id", handlers.GetUser())
        api.PUT("/users/:id", handlers.UpdateUser())
        api.DELETE("/users/:id", handlers.DeleteUser())
        
        // Content routes
        api.GET("/content", handlers.ListContent())
        api.GET("/content/:id", handlers.GetContent())
        api.POST("/content", handlers.CreateContent())
        api.PUT("/content/:id", handlers.UpdateContent())
        api.DELETE("/content/:id", handlers.DeleteContent())
        api.POST("/content/:id/publish", handlers.PublishContent())
        
        // AI routes
        api.POST("/ai/generate", handlers.GenerateWithAI())
        api.POST("/ai/analyze", handlers.AnalyzeWithAI())
        
        // Integration routes
        api.GET("/integrations", handlers.ListIntegrations())
        api.POST("/integrations/:provider/connect", handlers.ConnectIntegration())
        api.DELETE("/integrations/:provider/disconnect", handlers.DisconnectIntegration())
        
        // Workflow routes
        api.GET("/workflows", handlers.ListWorkflows())
        api.POST("/workflows", handlers.CreateWorkflow())
        api.POST("/workflows/:id/trigger", handlers.TriggerWorkflow())
    }
    
    // Webhook routes (no auth required)
    webhooks := g.router.Group("/webhooks")
    {
        webhooks.POST("/n8n/*path", handlers.N8NWebhook())
        webhooks.POST("/stripe/*path", handlers.StripeWebhook())
        webhooks.POST("/integration/*path", handlers.IntegrationWebhook())
    }
}

func (g *APIGateway) Start() error {
    g.server = &http.Server{
        Addr:    fmt.Sprintf("%s:%d", g.config.Server.Host, g.config.Server.Port),
        Handler: g.router,
    }
    
    log.Printf("Starting API Gateway on %s:%d", g.config.Server.Host, g.config.Server.Port)
    
    go func() {
        if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()
    
    return nil
}

func (g *APIGateway) Stop(ctx context.Context) error {
    log.Println("Shutting down API Gateway...")
    
    if err := g.server.Shutdown(ctx); err != nil {
        return fmt.Errorf("server shutdown failed: %w", err)
    }
    
    if err := g.redisClient.Close(); err != nil {
        log.Printf("Redis close error: %v", err)
    }
    
    return nil
}

func main() {
    cfg, err := config.LoadConfig("config.yaml")
    if err != nil {
        log.Fatalf("Config load failed: %v", err)
    }
    
    gateway, err := NewAPIGateway(cfg)
    if err != nil {
        log.Fatalf("Gateway initialization failed: %v", err)
    }
    
    if err := gateway.Start(); err != nil {
        log.Fatalf("Gateway start failed: %v", err)
    }
    
    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := gateway.Stop(ctx); err != nil {
        log.Fatalf("Gateway stop failed: %v", err)
    }
}
```

### 14.3 Service Discovery

Für die Service-zu-Service-Kommunikation implementiert BIOMETRICS Service Discovery, das dynamische Service-Registrierung und -Auflösung ermöglicht:

```go
// service-discovery/service_discovery.go
package discovery

import (
    "context"
    "fmt"
    "sync"
    "time"
    
    "github.com/redis/go-redis/v9"
)

type Service struct {
    Name      string
    Host      string
    Port      int
    Version   string
    Metadata  map[string]string
    TTL       time.Duration
}

type ServiceDiscovery struct {
    redis     *redis.Client
    services   map[string][]*Service
    mu        sync.RWMutex
    heartbeat  time.Duration
}

func NewServiceDiscovery(redisHost string, redisPort int) (*ServiceDiscovery, error) {
    client := redis.NewClient(&redis.Options{
        Addr: fmt.Sprintf("%s:%d", redisHost, redisPort),
    })
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("redis connection failed: %w", err)
    }
    
    return &ServiceDiscovery{
        redis:     client,
        services:  make(map[string][]*Service),
        heartbeat: 10 * time.Second,
    }, nil
}

func (sd *ServiceDiscovery) Register(ctx context.Context, svc *Service) error {
    key := fmt.Sprintf("service:%s:%s:%d", svc.Name, svc.Host, svc.Port)
    
    data, err := json.Marshal(svc)
    if err != nil {
        return fmt.Errorf("marshal failed: %w", err)
    }
    
    // Register with TTL (auto-expire if heartbeat stops)
    if err := sd.redis.Set(ctx, key, data, svc.TTL).Err(); err != nil {
        return fmt.Errorf("register failed: %w", err)
    }
    
    // Add to local cache
    sd.mu.Lock()
    defer sd.mu.Unlock()
    
    // Remove existing instance and add new one
    sd.services[svc.Name] = append(
        filterServices(sd.services[svc.Name], svc.Host, svc.Port),
        svc,
    )
    
    return nil
}

func (sd *ServiceDiscovery) Deregister(ctx context.Context, name, host string, port int) error {
    key := fmt.Sprintf("service:%s:%s:%d", name, host, port)
    
    if err := sd.redis.Del(ctx, key).Err(); err != nil {
        return fmt.Errorf("deregister failed: %w", err)
    }
    
    sd.mu.Lock()
    defer sd.mu.Unlock()
    
    sd.services[name] = filterServices(sd.services[name], host, port)
    
    return nil
}

func (sd *ServiceDiscovery) Discover(ctx context.Context, name string) ([]*Service, error) {
    // Try local cache first
    sd.mu.RLock()
    if services, ok := sd.services[name]; ok && len(services) > 0 {
        sd.mu.RUnlock()
        return services, nil
    }
    sd.mu.RUnlock()
    
    // Fetch from Redis
    pattern := fmt.Sprintf("service:%s:*", name)
    keys, err := sd.redis.Keys(ctx, pattern).Result()
    if err != nil {
        return nil, fmt.Errorf("discovery failed: %w", err)
    }
    
    var services []*Service
    if len(keys) == 0 {
        return nil, ErrServiceNotFound
    }
    
    for _, key := range keys {
        data, err := sd.redis.Get(ctx, key).Result()
        if err != nil {
            continue
        }
        
        var svc Service
        if err := json.Unmarshal([]byte(data), &svc); err == nil {
            services = append(services, &svc)
        }
    }
    
    // Update local cache
    sd.mu.Lock()
    sd.services[name] = services
    sd.mu.Unlock()
    
    return services, nil
}

func (sd *ServiceDiscovery) GetRandomService(ctx context.Context, name string) (*Service, error) {
    services, err := sd.Discover(ctx, name)
    if err != nil {
        return nil, err
    }
    
    if len(services) == 0 {
        return nil, ErrServiceNotFound
    }
    
    // Simple round-robin
    sd.mu.Lock()
    idx := sd.roundRobin[name]
    sd.roundRobin[name] = (idx + 1) % len(services)
    sd.mu.Unlock()
    
    return services[idx], nil
}

// StartHeartbeat starts the heartbeat mechanism for a service
func (sd *ServiceDiscovery) StartHeartbeat(ctx context.Context, svc *Service) {
    ticker := time.NewTicker(sd.heartbeat)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            sd.Deregister(ctx, svc.Name, svc.Host, svc.Port)
            return
        case <-ticker.C:
            if err := sd.Register(ctx, svc); err != nil {
                log.Printf("Heartbeat failed: %v", err)
            }
        }
    }
}
```

### 14.4 Load Balancing

Das Load Balancing in BIOMETRICS erfolgt auf mehreren Ebenen: DNS-basiert für geografische Verteilung, Client-seitig für Service-zu-Service-Kommunikation, und Envoy-basiert für Ingress-Traffic:

```go
// loadbalancer/client.go
package loadbalancer

import (
    "context"
    "math/rand"
    "sync"
    "time"
)

type LoadBalancer interface {
    Select(ctx context.Context) (string, error)
}

type RoundRobin struct {
    instances []string
    current   uint32
    mu        sync.Mutex
}

func NewRoundRobin(instances []string) *RoundRobin {
    return &RoundRobin{
        instances: instances,
        current:   0,
    }
}

func (lb *RoundRobin) Select(ctx context.Context) (string, error) {
    if len(lb.instances) == 0 {
        return "", ErrNoInstances
    }
    
    lb.mu.Lock()
    idx := lb.current % uint32(len(lb.instances))
    lb.current++
    lb.mu.Unlock()
    
    return lb.instances[idx], nil
}

type WeightedRoundRobin struct {
    instances []WeightedInstance
    current   int
    mu        sync.Mutex
}

type WeightedInstance struct {
    Address string
    Weight   int
    Current  int
}

func (lb *WeightedRoundRobin) Select(ctx context.Context) (string, error) {
    if len(lb.instances) == 0 {
        return "", ErrNoInstances
    }
    
    lb.mu.Lock()
    defer lb.mu.Unlock()
    
    // Find instance with lowest current weight
    selected := 0
    minWeight := lb.instances[0].Current
    
    for i, inst := range lb.instances {
        if inst.Current < minWeight {
            minWeight = inst.Current
            selected = i
        }
    }
    
    // Increment selected instance
    lb.instances[selected].Current += lb.instances[selected].Weight
    
    return lb.instances[selected].Address, nil
}

type LeastConnections struct {
    instances map[string]*ConnectionTracker
    mu        sync.Mutex
}

type ConnectionTracker struct {
    Address      string
    ActiveConns  int64
    LastUpdated time.Time
}

func (lb *LeastConnections) Select(ctx context.Context) (string, error) {
    if len(lb.instances) == 0 {
        return "", ErrNoInstances
    }
    
    lb.mu.Lock()
    defer lb.mu.Unlock()
    
    var selected string
    var minConns int64 = math.MaxInt64
    
    for addr, tracker := range lb.instances {
        if tracker.ActiveConns < minConns {
            minConns = tracker.ActiveConns
            selected = addr
        }
    }
    
    lb.instances[selected].ActiveConns++
    lb.instances[selected].LastUpdated = time.Now()
    
    return selected, nil
}

func (lb *LeastConnections) Release(addr string) {
    lb.mu.Lock()
    defer lb.mu.Unlock()
    
    if tracker, ok := lb.instances[addr]; ok && tracker.ActiveConns > 0 {
        tracker.ActiveConns--
    }
}
```

### 14.5 gRPC Kommunikation

Für service-interne Kommunikation nutzt BIOMETRICS gRPC aufgrund seiner Performanz und starken Typisierung:

```go
// proto/user.proto
syntax = "proto3";

package biometrics;

option go_package = "biometrics/gen/go/user/v1";

import "google/protobuf/timestamp.proto";
import "google/protobuf/empty.proto";

service UserService {
    rpc CreateUser(CreateUserRequest) returns (User);
    rpc GetUser(GetUserRequest) returns (User);
    rpc UpdateUser(UpdateUserRequest) returns (User);
    rpc DeleteUser(DeleteUserRequest) returns (google.protobuf.Empty);
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
    rpc SearchUsers(SearchUsersRequest) returns (SearchUsersResponse);
}

message User {
    string id = 1;
    string email = 2;
    string name = 3;
    UserRole role = 4;
    google.protobuf.Timestamp created_at = 5;
    google.protobuf.Timestamp updated_at = 6;
    map<string, string> metadata = 7;
}

enum UserRole {
    USER_ROLE_UNSPECIFIED = 0;
    USER_ROLE_ADMIN = 1;
    USER_ROLE_USER = 2;
    USER_ROLE_GUEST = 3;
}

message CreateUserRequest {
    string email = 1;
    string name = 2;
    string password = 3;
    UserRole role = 4;
}

message GetUserRequest {
    string id = 1;
}

// ... more message definitions
```

```go
// grpc/client.go
package grpc

import (
    "context"
    "crypto/tls"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/resolver"
    "google.golang.org/grpc/health/grpc_health_v1"
)

type Client struct {
    conn    *grpc.ClientConn
    client  UserServiceClient
    health  grpc_health_v1.HealthClient
}

func NewClient(addr string, opts ...Option) (*Client, error) {
    options := &Options{
        timeout: 10 * time.Second,
        insecure: false,
    }
    
    for _, opt := range opts {
        opt(options)
    }
    
    var creds credentials.TransportCredentials
    if options.insecure {
        creds = insecure.NewCredentials()
    } else {
        creds = credentials.NewTLS(&tls.Config{
            MinVersion: tls.VersionTLS12,
        })
    }
    
    conn, err := grpc.Dial(
        addr,
        grpc.WithTransportCredentials(creds),
        grpc.WithUnaryInterceptor(UnaryClientInterceptor(options.logger)),
        grpc.WithStreamInterceptor(StreamClientInterceptor(options.logger)),
        grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
    )
    
    if err != nil {
        return nil, fmt.Errorf("dial failed: %w", err)
    }
    
    return &Client{
        conn:   conn,
        client: NewUserServiceClient(conn),
        health: grpc_health_v1.NewHealthClient(conn),
    }, nil
}

func (c *Client) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    return c.client.CreateUser(ctx, req)
}

func (c *Client) GetUser(ctx context.Context, req *GetUserRequest) (*User, error) {
    return c.client.GetUser(ctx, req)
}

func (c *Client) Close() error {
    return c.conn.Close()
}
```

---

## 15) Database Architecture - PostgreSQL Schema und Caching Strategy

### 15.1 Datenbank-Schema Design

Das PostgreSQL-Schema von BIOMETRICS ist für Skalierbarkeit und Performance optimiert. Die folgende Übersicht zeigt die Hauptentitäten und ihre Beziehungen:

```sql
-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "btree_gin";

-- Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    email_verified_at TIMESTAMP WITH TIME ZONE,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT users_email_check CHECK (
        email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
    )
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- User Profiles (1:1 with users)
CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    avatar_url TEXT,
    bio TEXT,
    timezone VARCHAR(50) DEFAULT 'UTC',
    locale VARCHAR(10) DEFAULT 'en',
    preferences JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Sessions Table
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    refresh_token_hash VARCHAR(255),
    ip_address INET,
    user_agent TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT sessions_token_unique UNIQUE (token_hash)
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at) WHERE deleted_at IS NULL;

-- Contents Table
CREATE TABLE contents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    body TEXT,
    content_type VARCHAR(50) NOT NULL DEFAULT 'text',
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    metadata JSONB DEFAULT '{}',
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_contents_user_id ON contents(user_id);
CREATE INDEX idx_contents_status ON contents(status);
CREATE INDEX idx_contents_type ON contents(content_type);
CREATE INDEX idx_contents_published_at ON contents(published_at DESC) 
    WHERE status = 'published';

-- Full-text search index for contents
CREATE INDEX idx_contents_fts ON contents USING gin(
    to_tsvector('english', title || ' ' || COALESCE(body, ''))
);

-- Content Assets Table
CREATE TABLE content_assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content_id UUID NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    asset_type VARCHAR(50) NOT NULL,
    url TEXT NOT NULL,
    mime_type VARCHAR(100),
    size_bytes BIGINT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_content_assets_content_id ON content_assets(content_id);

-- Integrations Table
CREATE TABLE integrations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    access_token TEXT,
    refresh_token TEXT,
    token_expires_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB DEFAULT '{}',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT integrations_user_provider_unique UNIQUE (user_id, provider)
);

CREATE INDEX idx_integrations_user_id ON integrations(user_id);
CREATE INDEX idx_integrations_provider ON integrations(provider);
CREATE INDEX idx_integrations_status ON integrations(status);

-- Workflows Table
CREATE TABLE workflows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    definition JSONB NOT NULL,
    trigger_type VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_run_at TIMESTAMP WITH TIME ZONE,
    run_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workflows_user_id ON workflows(user_id);
CREATE INDEX idx_workflows_trigger_type ON workflows(trigger_type);
CREATE INDEX idx_workflows_is_active ON workflows(is_active);

-- Workflow Executions Table
CREATE TABLE workflow_executions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    input_data JSONB,
    output_data JSONB,
    error_message TEXT,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workflow_executions_workflow_id ON workflow_executions(workflow_id);
CREATE INDEX idx_workflow_executions_status ON workflow_executions(status);
CREATE INDEX idx_workflow_executions_created_at ON workflow_executions(created_at DESC);

-- AI Requests Table (for auditing and analytics)
CREATE TABLE ai_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    model VARCHAR(100) NOT NULL,
    prompt_tokens INTEGER,
    completion_tokens INTEGER,
    total_tokens INTEGER,
    latency_ms INTEGER,
    status VARCHAR(50) NOT NULL,
    error_message TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ai_requests_user_id ON ai_requests(user_id);
CREATE INDEX idx_ai_requests_model ON ai_requests(model);
CREATE INDEX idx_ai_requests_created_at ON ai_requests(created_at DESC);

-- Audit Logs Table
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID,
    changes JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- Function: updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for updated_at
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_profiles_updated_at 
    BEFORE UPDATE ON user_profiles 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_contents_updated_at 
    BEFORE UPDATE ON contents 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_integrations_updated_at 
    BEFORE UPDATE ON integrations 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_workflows_updated_at 
    BEFORE UPDATE ON workflows 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### 15.2 Redis Cache Strategy

Das Caching-System nutzt Redis für mehrere Zwecke und implementiert verschiedene Caching-Strategien:

```go
// cache/redis.go
package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/redis/go-redis/v9"
)

type RedisCache struct {
    client *redis.Client
    defaultTTL time.Duration
}

func NewRedisCache(host string, port int, password string, db int) *RedisCache {
    client := redis.NewClient(&redis.Options{
        Addr:         fmt.Sprintf("%s:%d", host, port),
        Password:     password,
        DB:           db,
        PoolSize:     100,
        MinIdleConns: 10,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    })
    
    return &RedisCache{
        client:     client,
        defaultTTL: 15 * time.Minute,
    }
}

func (rc *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
    result, err := rc.client.Get(ctx, key).Bytes()
    if err == redis.Nil {
        return nil, ErrCacheMiss
    }
    return result, err
}

func (rc *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    if ttl == 0 {
        ttl = rc.defaultTTL
    }
    return rc.client.Set(ctx, key, value, ttl).Err()
}

func (rc *RedisCache) Delete(ctx context.Context, key string) error {
    return rc.client.Del(ctx, key).Err()
}

func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
    n, err := rc.client.Exists(ctx, key).Result()
    return n > 0, err
}

func (rc *RedisCache) GetObject(ctx context.Context, key string, dest interface{}) error {
    data, err := rc.Get(ctx, key)
    if err != nil {
        return err
    }
    return json.Unmarshal(data, dest)
}

func (rc *RedisCache) SetObject(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return rc.Set(ctx, key, data, ttl)
}

// Cache Keys
const (
    KeyUserPrefix       = "user:"
    KeyContentPrefix    = "content:"
    KeySessionPrefix    = "session:"
    KeyRateLimitPrefix  = "ratelimit:"
    KeyAIResponsePrefix = "ai:response:"
)

func UserCacheKey(id string) string           { return KeyUserPrefix + id }
func ContentCacheKey(id string) string        { return KeyContentPrefix + id }
func SessionCacheKey(token string) string    { return KeySessionPrefix + token }
func RateLimitCacheKey(key string) string    { return KeyRateLimitPrefix + key }
func AIResponseCacheKey(hash string) string  { return KeyAIResponsePrefix + hash }
```

**Caching Strategien:**

| Strategie | Anwendungsfall | TTL | Invalidation |
|-----------|---------------|-----|--------------|
| Cache-Aside | User Profiles | 15 min | On update |
| Write-Through | Session Data | 24h | On logout |
| Write-Behind | Analytics | 5 min | Scheduled |
| Time-Based | API Responses | 5-15 min | TTL expiry |
| Stale-While-Refresh | AI Responses | 1h | Background refresh |

```go
// cache/strategies.go
package cache

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
)

// StaleWhileRevalidate returns stale cache data while refreshing in background
func (rc *RedisCache) GetStaleWhileRevalidate(
    ctx context.Context,
    key string,
    staleTTL time.Duration,
    fetchFunc func() ([]byte, error),
) ([]byte, error) {
    // Try to get from cache
    data, err := rc.client.Get(ctx, key).Bytes()
    if err == nil {
        // Cache hit - check if stale
        ttl, _ := rc.client.TTL(ctx, key).Result()
        if ttl > 0 && ttl < staleTTL {
            // Data is stale - trigger background refresh
            go func() {
                if freshData, err := fetchFunc(); err == nil {
                    rc.client.Set(ctx, key, freshData, rc.defaultTTL)
                }
            }()
        }
        return data, nil
    }
    
    // Cache miss - fetch synchronously
    if err != redis.Nil {
        return nil, err
    }
    
    return fetchFunc()
}

// GenerateCacheKey creates a hash-based cache key for AI prompts
func GenerateCacheKey(prompt string, params map[string]interface{}) string {
    h := sha256.New()
    h.Write([]byte(prompt))
    
    if params != nil {
        for k, v := range params {
            h.Write([]byte(k))
            h.Write([]byte(fmt.Sprintf("%v", v)))
        }
    }
    
    return "ai:" + hex.EncodeToString(h.Sum(nil))[:16]
}
```

### 15.3 Daten-Migration und Backup

Das Migrationssystem ermöglicht sichere Schema-Evolutionen:

```go
// database/migration.go
package database

import (
    "context"
    "embed"
    "fmt"
    "io"
    "sort"
    "time"
    
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    "github.com/golang-migrate/migrate/v4/source/iofs"
)

type MigrationManager struct {
    m *migrate.Migrate
}

func NewMigrationManager(db *sql.DB, migrations embed.FS) (*MigrationManager, error) {
    driver, err := postgres.WithInstance(db, &postgres.Config{
        MigrationsTable: "schema_migrations",
    })
    if err != nil {
        return nil, fmt.Errorf("create driver failed: %w", err)
    }
    
    source, err := iofs.New(migrations, "migrations")
    if err != nil {
        return nil, fmt.Errorf("create source failed: %w", err)
    }
    
    m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
    if err != nil {
        return nil, fmt.Errorf("create migrate failed: %w", err)
    }
    
    return &MigrationManager{m: m}, nil
}

func (mm *MigrationManager) Up(ctx context.Context) error {
    if err := mm.m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("migration up failed: %w", err)
    }
    return nil
}

func (mm *MigrationManager) Down(ctx context.Context) error {
    if err := mm.m.Down(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("migration down failed: %w", err)
    }
    return nil
}

func (mm *MigrationManager) Force(version int) error {
    if err := mm.m.Force(version); err != nil {
        return fmt.Errorf("migration force failed: %w", err)
    }
    return nil
}

func (mm *MigrationManager) Version() (version uint, dirty bool, err error) {
    return mm.m.Version()
}

// Backup functionality
type BackupManager struct {
    db     *sql.DB
    s3     *s3.Client
    bucket string
}

func (bm *BackupManager) CreateBackup(ctx context.Context) (string, error) {
    timestamp := time.Now().Format("2006-01-02T15-04-05Z")
    filename := fmt.Sprintf("biometrics-backup-%s.sql", timestamp)
    
    // Create SQL dump
    cmd := exec.CommandContext(ctx, "pg_dump",
        "-h", os.Getenv("DB_HOST"),
        "-p", os.Getenv("DB_PORT"),
        "-U", os.Getenv("DB_USER"),
        "-d", os.Getenv("DB_NAME"),
        "--clean",
        "--if-exists",
    )
    
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("pg_dump failed: %w", err)
    }
    
    // Upload to S3
    reader := bytes.NewReader(output)
    _, err = bm.s3.PutObject(ctx, &s3.PutObjectInput{
        Bucket: aws.String(bm.bucket),
        Key:   aws.String(fmt.Sprintf("backups/%s", filename)),
        Body:  reader,
    })
    
    if err != nil {
        return "", fmt.Errorf("s3 upload failed: %w", err)
    }
    
    return filename, nil
}

func (bm *BackupManager) RestoreBackup(ctx context.Context, filename string) error {
    // Download from S3
    result, err := bm.s3.GetObject(ctx, &s3.GetObjectInput{
        Bucket: aws.String(bm.bucket),
        Key:    aws.String(fmt.Sprintf("backups/%s", filename)),
    })
    if err != nil {
        return fmt.Errorf("s3 download failed: %w", err)
    }
    defer result.Body.Close()
    
    // Execute SQL
    cmd := exec.CommandContext(ctx, "psql",
        "-h", os.Getenv("DB_HOST"),
        "-p", os.Getenv("DB_PORT"),
        "-U", os.Getenv("DB_USER"),
        "-d", os.Getenv("DB_NAME"),
    )
    cmd.Stdin = result.Body
    
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("psql restore failed: %w, output: %s", err, string(output))
    }
    
    return nil
}
```

---

## 16) Security Architecture - Zero Trust Design und Authentifizierung

### 16.1 Zero Trust Architektur

Die Sicherheitsarchitektur von BIOMETRICS basiert auf dem Zero-Trust-Prinzip, das keine implizite Vertrauensstellung annimmt und jede Anfrage kritisch überprüft:

```go
// security/zero_trust.go
package security

import (
    "context"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "math/big"
    "time"
)

type ZeroTrustConfig struct {
    RequireMFA           bool
    SessionTimeout       time.Duration
    MaxLoginAttempts     int
    LockoutDuration      time.Duration
    PasswordMinLength    int
    PasswordRequireUpper bool
    PasswordRequireLower bool
    PasswordRequireDigit bool
    PasswordRequireSpecial bool
}

type SecurityManager struct {
    config         *ZeroTrustConfig
    passwordHasher PasswordHasher
    tokenManager   TokenManager
    mfaService    MFAService
    auditLogger   AuditLogger
}

func NewSecurityManager(cfg *ZeroTrustConfig) *SecurityManager {
    return &SecurityManager{
        config:         cfg,
        passwordHasher: NewBcryptHasher(12),
        tokenManager:   NewJWTManager(cfg.SessionTimeout),
        mfaService:     NewTOTPMFAService(),
        auditLogger:    NewAuditLogger(),
    }
}

// VerifyRequest implements Zero Trust verification for every request
func (sm *SecurityManager) VerifyRequest(ctx context.Context, req *Request) (*AuthContext, error) {
    // 1. Verify identity
    claims, err := sm.tokenManager.ValidateToken(ctx, req.Token)
    if err != nil {
        sm.auditLogger.Log(ctx, AuditEvent{
            Action:   "TOKEN_VALIDATION_FAILED",
            Resource: req.Path,
            Metadata: map[string]string{"error": err.Error()},
        })
        return nil, ErrUnauthorized
    }
    
    // 2. Verify device (if MFA enabled)
    if sm.config.RequireMFA && !claims.MFAVerified {
        if !sm.mfaService.Verify(ctx, claims.UserID, req.MFACode) {
            return nil, ErrMFARequired
        }
    }
    
    // 3. Verify authorization
    if !sm.authorize(claims, req) {
        sm.auditLogger.Log(ctx, AuditEvent{
            Action:   "AUTHORIZATION_FAILED",
            UserID:   claims.UserID,
            Resource: req.Path,
            Metadata: map[string]string{"method": req.Method},
        })
        return nil, ErrForbidden
    }
    
    // 4. Verify request integrity
    if !sm.verifyIntegrity(req) {
        return nil, ErrIntegrityCheckFailed
    }
    
    // 5. Check for anomalies
    if sm.detectAnomaly(ctx, claims, req) {
        sm.auditLogger.Log(ctx, AuditEvent{
            Action:   "ANOMALY_DETECTED",
            UserID:   claims.UserID,
            Resource: req.Path,
            Metadata: map[string]string{"ip": req.ClientIP},
        })
        // Could trigger additional verification or block
    }
    
    return &AuthContext{
        UserID:  claims.UserID,
        Email:   claims.Email,
        Role:    claims.Role,
        MFA:     claims.MFAVerified,
        Device:  req.DeviceID,
        ClientIP: req.ClientIP,
    }, nil
}

func (sm *SecurityManager) authorize(claims *Claims, req *Request) bool {
    // Load user permissions
    permissions := loadUserPermissions(claims.Role)
    
    // Check if permission exists for this resource + method
    key := req.Method + ":" + req.Path
    return permissions[key]
}
```

### 16.2 Authentication Flow

Der Authentifizierungsprozess implementiert Multi-Faktor-Authentifizierung und sichere Session-Verwaltung:

```go
// auth/service.go
package auth

import (
    "context"
    "errors"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    userRepo   UserRepository
    sessionRepo SessionRepository
    cache      Cache
    config     *AuthConfig
}

type AuthConfig struct {
    JWTSecret          string
    JWTExpiry          time.Duration
    RefreshTokenExpiry time.Duration
    MaxSessions        int
    RequireMFA         bool
}

type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required,min=2"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
    MFA      string `json:"mfa,omitempty"`
}

type AuthResponse struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    time.Time `json:"expires_at"`
    User         *User     `json:"user"`
}

func (as *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
    // Validate password strength
    if err := as.validatePassword(req.Password); err != nil {
        return nil, err
    }
    
    // Check if user exists
    existing, _ := as.userRepo.GetByEmail(ctx, req.Email)
    if existing != nil {
        return nil, ErrUserAlreadyExists
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(req.Password),
        bcrypt.DefaultCost,
    )
    if err != nil {
        return nil, ErrPasswordHashFailed
    }
    
    // Create user
    user := &User{
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
        Name:         req.Name,
        Role:         "user",
    }
    
    createdUser, err := as.userRepo.Create(ctx, user)
    if err != nil {
        return nil, err
    }
    
    // Generate tokens
    return as.generateAuthResponse(ctx, createdUser)
}

func (as *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
    // Get user
    user, err := as.userRepo.GetByEmail(ctx, req.Email)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            // Prevent timing attacks
            bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
            return nil, ErrInvalidCredentials
        }
        return nil, err
    }
    
    // Verify password
    if err := bcrypt.CompareHashAndPassword(
        []byte(user.PasswordHash),
        []byte(req.Password),
    ); err != nil {
        // Log failed attempt
        as.logFailedLogin(ctx, user.ID)
        return nil, ErrInvalidCredentials
    }
    
    // Check MFA if enabled
    if user.MFAEnabled {
        if req.MFA == "" {
            return &AuthResponse{
                MFARequired: true,
            }, nil
        }
        
        if !as.verifyMFA(ctx, user.ID, req.MFA) {
            return nil, ErrInvalidMFA
        }
    }
    
    // Check for account lock
    if as.isAccountLocked(ctx, user.ID) {
        return nil, ErrAccountLocked
    }
    
    // Generate tokens
    response, err := as.generateAuthResponse(ctx, user)
    if err != nil {
        return nil, err
    }
    
    // Clear failed login attempts
    as.clearFailedLogins(ctx, user.ID)
    
    // Update last login
    as.userRepo.UpdateLastLogin(ctx, user.ID)
    
    return response, nil
}

func (as *AuthService) generateAuthResponse(ctx context.Context, user *User) (*AuthResponse, error) {
    // Generate access token
    accessClaims := &Claims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(as.config.JWTExpiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).
        SignedString([]byte(as.config.JWTSecret))
    if err != nil {
        return nil, err
    }
    
    // Generate refresh token
    refreshToken, err := generateSecureToken(32)
    if err != nil {
        return nil, err
    }
    
    // Store session
    session := &Session{
        UserID:         user.ID,
        RefreshToken:   refreshToken,
        AccessToken:    accessToken,
        ExpiresAt:      time.Now().Add(as.config.RefreshTokenExpiry),
        IPAddress:      getClientIP(ctx),
        UserAgent:      getUserAgent(ctx),
    }
    
    if err := as.sessionRepo.Create(ctx, session); err != nil {
        return nil, err
    }
    
    return &AuthResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresAt:    session.ExpiresAt,
        User:         user.WithoutPassword(),
    }, nil
}

func (as *AuthService) validatePassword(password string) error {
    if len(password) < 8 {
        return ErrPasswordTooShort
    }
    
    hasUpper := false
    hasLower := false
    hasDigit := false
    hasSpecial := false
    
    for _, char := range password {
        switch {
        case char >= 'A' && char <= 'Z':
            hasUpper = true
        case char >= 'a' && char <= 'z':
            hasLower = true
        case char >= '0' && char <= '9':
            hasDigit = true
        default:
            hasSpecial = true
        }
    }
    
    if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
        return ErrPasswordTooWeak
    }
    
    return nil
}
```

### 16.3 Authorization Matrix

Die Autorisierung basiert auf Rollen und Berechtigungen mit fein granulierter Zugriffskontrolle:

```go
// auth/rbac.go
package auth

type Permission string

const (
    // User permissions
    PermissionUserRead   Permission = "user:read"
    PermissionUserWrite  Permission = "user:write"
    PermissionUserDelete Permission = "user:delete"
    
    // Content permissions
    PermissionContentRead   Permission = "content:read"
    PermissionContentWrite Permission = "content:write"
    PermissionContentDelete Permission = "content:delete"
    PermissionContentPublish Permission = "content:publish"
    
    // AI permissions
    PermissionAIGenerate Permission = "ai:generate"
    PermissionAIAnalyze  Permission = "ai:analyze"
    
    // Admin permissions
    PermissionAdmin Users    Permission = "admin:users"
    PermissionAdminSystem Permission = "admin:system"
    PermissionAdminAudit Permission = "admin:audit"
)

type Role struct {
    Name        string
    Permissions []Permission
}

var Roles = map[string]Role{
    "guest": {
        Name: "guest",
        Permissions: []Permission{
            PermissionContentRead,
        },
    },
    "user": {
        Name: "user",
        Permissions: []Permission{
            PermissionUserRead,
            PermissionContentRead,
            PermissionContentWrite,
            PermissionContentDelete,
            PermissionAIGenerate,
            PermissionAIAnalyze,
        },
    },
    "admin": {
        Name: "admin",
        Permissions: []Permission{
            PermissionUserRead,
            PermissionUserWrite,
            PermissionUserDelete,
            PermissionContentRead,
            PermissionContentWrite,
            PermissionContentDelete,
            PermissionContentPublish,
            PermissionAIGenerate,
            PermissionAIAnalyze,
            PermissionAdminUsers,
            PermissionAdminAudit,
        },
    },
    "superadmin": {
        Name: "superadmin",
        Permissions: []Permission{
            PermissionUserRead,
            PermissionUserWrite,
            PermissionUserDelete,
            PermissionContentRead,
            PermissionContentWrite,
            PermissionContentDelete,
            PermissionContentPublish,
            PermissionAIGenerate,
            PermissionAIAnalyze,
            PermissionAdminUsers,
            PermissionAdminSystem,
            PermissionAdminAudit,
        },
    },
}

type RBACService struct {
    roleRepo RoleRepository
    cache    Cache
}

func (rs *RBACService) HasPermission(ctx context.Context, userID string, perm Permission) (bool, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("perm:%s:%s", userID, perm)
    cached, err := rs.cache.Get(ctx, cacheKey)
    if err == nil {
        return string(cached) == "1", nil
    }
    
    // Get user role
    user, err := rs.userRepo.GetByID(ctx, userID)
    if err != nil {
        return false, err
    }
    
    role, ok := Roles[user.Role]
    if !ok {
        return false, ErrRoleNotFound
    }
    
    // Check permission
    hasPermission := false
    for _, p := range role.Permissions {
        if p == perm {
            hasPermission = true
            break
        }
    }
    
    // Cache result (short TTL)
    rs.cache.Set(ctx, cacheKey, []byte(map[bool]string{true: "1", false: "0"}[hasPermission]), 60*time.Second)
    
    return hasPermission, nil
}

func (rs *RBACService) RequirePermission(perm Permission) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")
        
        hasPermission, err := rs.HasPermission(c.Request.Context(), userID, perm)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
                Code:    "INTERNAL_ERROR",
                Message: "Failed to check permissions",
            })
            return
        }
        
        if !hasPermission {
            c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
                Code:    "FORBIDDEN",
                Message: "Insufficient permissions",
            })
            return
        }
        
        c.Next()
    }
}
```

### 16.4 Encryption Strategy

Datenverschlüsselung erfolgt auf mehreren Ebenen:

```go
// security/encryption.go
package security

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "io"
)

type EncryptionService struct {
    aesKey []byte
    rsaPublicKey *rsa.PublicKey
    rsaPrivateKey *rsa.PrivateKey
}

// AES-GCM for symmetric encryption (data at rest)
func (es *EncryptionService) EncryptAES(plaintext []byte) ([]byte, error) {
    block, err := aes.NewCipher(es.aesKey)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    
    return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (es *EncryptionService) DecryptAES(ciphertext []byte) ([]byte, error) {
    block, err := aes.NewCipher(es.aesKey)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }
    
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}

// RSA for asymmetric encryption (key exchange)
func (es *EncryptionService) EncryptRSA(plaintext []byte) ([]byte, error) {
    return rsa.EncryptOAEP(
        sha256.New(),
        rand.Reader,
        es.rsaPublicKey,
        plaintext,
        nil,
    )
}

func (es *EncryptionService) DecryptRSA(ciphertext []byte) ([]byte, error) {
    return rsa.DecryptOAEP(
        sha256.New(),
        rand.Reader,
        es.rsaPrivateKey,
        ciphertext,
        nil,
    )
}

// Hash for password storage and integrity
func (es *EncryptionService) HashSHA256(data []byte) string {
    hash := sha256.Sum256(data)
    return base64.StdEncoding.EncodeToString(hash[:])
}

// Key derivation for password-based encryption
func (es *EncryptionService) DeriveKey(password string, salt []byte) []byte {
    // Use Argon2id for password-based key derivation
    return argon2.IDKey(
        []byte(password),
        salt,
        64*1024,
        3,
        4,
        32,
    )
}
```

---

## 17) Scalability Design - Horizontale Skalierung und Performance

### 17.1 Horizontal Scaling Strategie

BIOMETRICS ist für horizontale Skalierung ausgelegt, wobei jeder Service unabhängig skaliert werden kann:

```yaml
# kubernetes/autoscaling.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: api-gateway-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: api-gateway
  minReplicas: 2
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 30
      policies:
      - type: Percent
        value: 100
        periodSeconds: 15
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: ai-orchestrator-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: ai-orchestrator
  minReplicas: 1
  maxReplicas: 10
  metrics:
  - type: Pods
    pods:
      metric:
        name: ai_requests_pending
      target:
        type: AverageValue
        averageValue: "10"
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Pods
        value: 2
        periodSeconds: 60
```

### 17.2 Auto-Scaling Regeln

```go
// scaling/auto_scaler.go
package scaling

import (
    "context"
    "time"
    
    "github.com/prometheus/client_golang/api"
    "github.com/prometheus/client_golang/api/prometheus/v1"
    "github.com/prometheus/common/model"
)

type AutoScaler struct {
    prometheusAPI v1.API
    k8sClient     K8sClient
    rules        []ScalingRule
}

type ScalingRule struct {
    MetricName      string
    TargetValue     float64
    ScaleUpDelta    int32
    ScaleDownDelta  int32
    CooldownPeriod  time.Duration
    MinReplicas     int32
    MaxReplicas     int32
}

var DefaultScalingRules = []ScalingRule{
    {
        MetricName:     "cpu_utilization",
        TargetValue:    70.0,
        ScaleUpDelta:   2,
        ScaleDownDelta: 1,
        CooldownPeriod: 5 * time.Minute,
        MinReplicas:    2,
        MaxReplicas:    20,
    },
    {
        MetricName:     "memory_utilization",
        TargetValue:    80.0,
        ScaleUpDelta:   2,
        ScaleDownDelta: 1,
        CooldownPeriod: 5 * time.Minute,
        MinReplicas:    2,
        MaxReplicas:    20,
    },
    {
        MetricName:     "http_requests_per_second",
        TargetValue:    1000.0,
        ScaleUpDelta:   3,
        ScaleDownDelta: 1,
        CooldownPeriod: 3 * time.Minute,
        MinReplicas:    2,
        MaxReplicas:    30,
    },
    {
        MetricName:     "queue_length",
        TargetValue:    100.0,
        ScaleUpDelta:   2,
        ScaleDownDelta: 1,
        CooldownPeriod: 2 * time.Minute,
        MinReplicas:    1,
        MaxReplicas:    10,
    },
}

func (as *AutoScaler) EvaluateScaling(ctx context.Context, deployment string) error {
    for _, rule := range as.rules {
        value, err := as.getMetricValue(ctx, rule.MetricName, deployment)
        if err != nil {
            continue // Skip on error, don't scale
        }
        
        currentReplicas, err := as.k8sClient.GetReplicas(ctx, deployment)
        if err != nil {
            return err
        }
        
        // Determine scaling action
        var newReplicas int32
        if value > rule.TargetValue {
            newReplicas = currentReplicas + rule.ScaleUpDelta
        } else if value < rule.TargetValue*0.5 {
            newReplicas = currentReplicas - rule.ScaleDownDelta
        }
        
        // Apply bounds
        if newReplicas < rule.MinReplicas {
            newReplicas = rule.MinReplicas
        }
        if newReplicas > rule.MaxReplicas {
            newReplicas = rule.MaxReplicas
        }
        
        // Apply scaling if needed
        if newReplicas != currentReplicas {
            if err := as.k8sClient.ScaleDeployment(ctx, deployment, newReplicas); err != nil {
                return err
            }
        }
    }
    
    return nil
}

func (as *AutoScaler) getMetricValue(ctx context.Context, metricName, deployment string) (float64, error) {
    query := ""
    switch metricName {
    case "cpu_utilization":
        query = fmt.Sprintf(
            "avg(rate(container_cpu_usage_seconds_total{namespace='default', pod=~'%s-.*'}[5m])) * 100",
            deployment,
        )
    case "memory_utilization":
        query = fmt.Sprintf(
            "avg(container_memory_usage_bytes{namespace='default', pod=~'%s-.*'}) / avg(container_spec_memory_limit_bytes{namespace='default', pod=~'%s-.*'}) * 100",
            deployment, deployment,
        )
    case "http_requests_per_second":
        query = fmt.Sprintf(
            "sum(rate(http_requests_total{service='%s'}[5m]))",
            deployment,
        )
    case "queue_length":
        query = fmt.Sprintf(
            "sum(redis_queue_length{service='%s'})",
            deployment,
        )
    }
    
    result, warnings, err := as.prometheusAPI.Query(ctx, query, time.Now())
    if err != nil {
        return 0, err
    }
    
    if len(warnings) > 0 {
        // Log warnings
    }
    
    vector := result.(model.Vector)
    if len(vector) == 0 {
        return 0, nil
    }
    
    return float64(vector[0].Value), nil
}
```

### 17.3 Performance Targets und Bottleneck-Analyse

**Performance-Ziele:**

| Metric | Target | Maximum | Critical |
|--------|--------|---------|----------|
| API Latency (p50) | 50ms | 100ms | 200ms |
| API Latency (p95) | 100ms | 200ms | 500ms |
| API Latency (p99) | 200ms | 500ms | 1000ms |
| AI Request Latency | 30s | 60s | 120s |
| Throughput | 1000 RPS | 5000 RPS | 10000 RPS |
| Availability | 99.9% | 99.95% | 99.99% |
| Error Rate | 0.1% | 0.5% | 1% |
| Time to Recovery | 5 min | 15 min | 30 min |

**Bottleneck-Analyse und Optimierungen:**

```go
// performance/optimization.go
package performance

import (
    "context"
    "sync"
    "time"
)

// Connection Pool Optimierung
type OptimizedPool struct {
    conns    chan *Connection
    mu       sync.Mutex
    active   int
    maxConns int
    minConns int
}

func NewOptimizedPool(min, max int) *OptimizedPool {
    pool := &OptimizedPool{
        conns:    make(chan *Connection, max),
        maxConns: max,
        minConns: min,
    }
    
    // Initialize minimum connections
    for i := 0; i < min; i++ {
        pool.conns <- pool.createConnection()
    }
    
    return pool
}

func (p *OptimizedPool) Get(ctx context.Context) (*Connection, error) {
    select {
    case conn := <-p.conns:
        if conn.IsHealthy() {
            return conn, nil
        }
        // Connection unhealthy, create new
        p.mu.Lock()
        p.active--
        p.mu.Unlock()
        return p.createConnection(), nil
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
        // Pool exhausted, create new if under max
        p.mu.Lock()
        defer p.mu.Unlock()
        if p.active < p.maxConns {
            p.active++
            return p.createConnection(), nil
        }
        return nil, ErrPoolExhausted
    }
}

// Query Optimization mit Prepared Statements
type QueryOptimizer struct {
    preparedStatements map[string]*PreparedStatement
    mu                sync.RWMutex
}

func (qo *QueryOptimizer) Prepare(name, query string) error {
    qo.mu.Lock()
    defer qo.mu.Unlock()
    
    stmt, err := db.Prepare(query)
    if err != nil {
        return err
    }
    
    qo.preparedStatements[name] = &PreparedStatement{
        Statement: stmt,
        Query:     query,
        Stats:     &QueryStats{},
    }
    
    return nil
}

func (qo *QueryOptimizer) ExecuteWithStats(ctx context.Context, name string, args ...interface{}) ([]Row, error) {
    start := time.Now()
    
    qo.mu.RLock()
    stmt, ok := qo.preparedStatements[name]
    qo.mu.RUnlock()
    
    if !ok {
        return nil, ErrStatementNotPrepared
    }
    
    rows, err := stmt.Statement.QueryContext(ctx, args...)
    if err != nil {
        return nil, err
    }
    
    // Record stats
    duration := time.Since(start)
    stmt.Stats.Record(duration, err == nil)
    
    return scanRows(rows)
}
```

---

## 18) Observability - Logging, Metrics, Tracing

### 18.1 Logging Strategie

```go
// observability/logger.go
package observability

import (
    "context"
    "encoding/json"
    "os"
    "time"
    
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

type Logger struct {
    *zap.SugaredLogger
    level     zapcore.Level
    fields    map[string]interface{}
}

type LogEntry struct {
    Timestamp time.Time              `json:"timestamp"`
    Level     string                 `json:"level"`
    Message   string                 `json:"message"`
    Fields    map[string]interface{} `json:"fields,omitempty"`
    TraceID   string                 `json:"trace_id,omitempty"`
    SpanID    string                 `json:"span_id,omitempty"`
}

func NewLogger(env string) (*Logger, error) {
    var config zap.Config
    
    if env == "production" {
        config = zap.NewProductionConfig()
        config.EncoderConfig.TimeKey = "timestamp"
        config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    } else {
        config = zap.NewDevelopmentConfig()
        config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    }
    
    // Output to stdout
    config.OutputPaths = []string{"stdout"}
    config.ErrorOutputPaths = []string{"stderr"}
    
    logger, err := config.Build()
    if err != nil {
        return nil, err
    }
    
    return &Logger{
        SugaredLogger: logger.Sugar(),
        level:         config.Level.Level,
        fields:        make(map[string]interface{}),
    }, nil
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
    newLogger := &Logger{
        SugaredLogger: l.SugaredLogger,
        level:         l.level,
        fields:        make(map[string]interface{}),
    }
    
    // Add trace information from context
    if traceID := GetTraceID(ctx); traceID != "" {
        newLogger.fields["trace_id"] = traceID
    }
    if spanID := GetSpanID(ctx); spanID != "" {
        newLogger.fields["span_id"] = spanID
    }
    
    return newLogger
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
    newLogger := &Logger{
        SugaredLogger: l.SugaredLogger,
        level:         l.level,
        fields:        copyFields(l.fields),
    }
    newLogger.fields[key] = value
    return newLogger
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
    l.log(zapcore.InfoLevel, msg, keysAndValues...)
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
    l.log(zapcore.ErrorLevel, msg, keysAndValues...)
}

func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
    l.log(zapcore.WarnLevel, msg, keysAndValues...)
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
    l.log(zapcore.DebugLevel, msg, keysAndValues...)
}

func (l *Logger) log(level zapcore.Level, msg string, keysAndValues ...interface{}) {
    // Convert keysAndValues to structured fields
    fields := make([]zap.Field, 0, len(l.fields)+len(keysAndValues)/2)
    
    for k, v := range l.fields {
        fields = append(fields, zap.Any(k, v))
    }
    
    for i := 0; i < len(keysAndValues); i += 2 {
        fields = append(fields, zap.Any(keysAndValues[i].(string), keysAndValues[i+1]))
    }
    
    switch level {
    case zapcore.InfoLevel:
        l.SugaredLogger.Info(msg, fields...)
    case zapcore.ErrorLevel:
        l.SugaredLogger.Error(msg, fields...)
    case zapcore.WarnLevel:
        l.SugaredLogger.Warn(msg, fields...)
    case zapcore.DebugLevel:
        l.SugaredLogger.Debug(msg, fields...)
    }
}
```

### 18.2 Metrics Collection

```go
// observability/metrics.go
package observability

import (
    "strconv"
    "time"
    
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // HTTP Metrics
    HTTPRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
        },
        []string{"method", "path"},
    )
    
    HTTPRequestsInFlight = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "http_requests_in_flight",
            Help: "Number of HTTP requests currently being processed",
        },
    )
    
    // Business Metrics
    ActiveUsers = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_users",
            Help: "Number of currently active users",
        },
    )
    
    ContentCreatedTotal = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "content_created_total",
            Help: "Total number of content items created",
        },
    )
    
    AIRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "ai_requests_total",
            Help: "Total number of AI requests",
        },
        []string{"model", "status"},
    )
    
    AIRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "ai_request_duration_seconds",
            Help:    "AI request duration in seconds",
            Buckets: []float64{1, 5, 10, 30, 60, 120, 300},
        },
        []string{"model"},
    )
    
    // Infrastructure Metrics
    DBQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "db_query_duration_seconds",
            Help:    "Database query duration in seconds",
            Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
        },
        []string{"query_type"},
    )
    
    RedisOperationDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "redis_operation_duration_seconds",
            Help:    "Redis operation duration in seconds",
            Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01},
        },
        []string{"operation"},
    )
)

// Middleware for HTTP metrics
func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        HTTPRequestsInFlight.Inc()
        defer HTTPRequestsInFlight.Dec()
        
        // Wrap response writer to capture status
        wrapped := &statusResponseWriter{ResponseWriter: w, statusCode: 200}
        
        next.ServeHTTP(wrapped, r)
        
        duration := time.Since(start).Seconds()
        
        HTTPRequestsTotal.WithLabelValues(
            r.Method,
            r.URL.Path,
            strconv.Itoa(wrapped.statusCode),
        ).Inc()
        
        HTTPRequestDuration.WithLabelValues(
            r.Method,
            r.URL.Path,
        ).Observe(duration)
    })
}

type statusResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (w *statusResponseWriter) WriteHeader(statusCode int) {
    w.statusCode = statusCode
    w.ResponseWriter.WriteHeader(statusCode)
}
```

### 18.3 Distributed Tracing

```go
// observability/tracing.go
package observability

import (
    "context"
    "fmt"
    "sync"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/otel/sdk/trace/samplers"
    "go.opentelemetry.io/otel/trace"
)

var (
    tracerProvider *trace.TracerProvider
    propagator    = propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    )
)

func InitTracer(serviceName, jaegerEndpoint string) error {
    exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint(jaegerEndpoint),
    ))
    if err != nil {
        return fmt.Errorf("jaeger exporter: %w", err)
    }
    
    res, err := resource.New(context.Background(),
        resource.WithAttributes(
            attribute.String("service.name", serviceName),
            attribute.String("service.version", "1.0.0"),
        ),
    )
    if err != nil {
        return fmt.Errorf("resource: %w", err)
    }
    
    tp, err := trace.NewTracerProvider(
        trace.WithBatcher(exp),
        trace.WithResource(res),
        trace.WithSampler(samplers.ProbabilityThreshold(0.1)),
    )
    if err != nil {
        return fmt.Errorf("tracer provider: %w", err)
    }
    
    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(propagator)
    
    tracerProvider = tp
    return nil
}

func ShutdownTracer(ctx context.Context) error {
    if tracerProvider != nil {
        return tracerProvider.Shutdown(ctx)
    }
    return nil
}

func Tracer(name string) trace.Tracer {
    return otel.Tracer(name)
}

func StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
    return otel.Tracer("biometrics").Start(ctx, name,
        trace.WithAttributes(attrs...),
    )
}

// Context propagation helpers
func InjectTraceContext(ctx context.Context, carrier propagation.TextMapCarrier) {
    propagator.Inject(ctx, carrier)
}

func ExtractTraceContext(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
    return propagator.Extract(ctx, carrier)
}

// Example usage in HTTP handler
func TracingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := ExtractTraceContext(r.Context(), r.Header)
        
        ctx, span := StartSpan(ctx, r.URL.Path,
            attribute.String("http.method", r.Method),
            attribute.String("http.url", r.URL.String()),
        )
        defer span.End()
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

## 19) Deployment Architecture - Kubernetes und CI/CD

### 19.1 Container Strategie

```dockerfile
# api/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /api \
    -ldflags="-s -w -X main.Version=${VERSION}" \
    ./cmd/api

# Final stage
FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /api .
COPY config.yaml .

ENV PORT=8080

EXPOSE 8080

USER 1000:1000

CMD ["./api"]
```

```yaml
# kubernetes/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
  labels:
    app: api-gateway
    version: v1
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-gateway
  template:
    metadata:
      labels:
        app: api-gateway
        version: v1
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
        prometheus.io/path: "/metrics"
    spec:
      serviceAccountName: api-gateway
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        fsGroup: 1000
      containers:
      - name: api-gateway
        image: biometrics/api-gateway:latest
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: grpc
        env:
        - name: PORT
          value: "8080"
        - name: CONFIG_PATH
          value: "/app/config.yaml"
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        resources:
          requests:
            cpu: "100m"
            memory: "128Mi"
          limits:
            cpu: "500m"
            memory: "512Mi"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          failureThreshold: 3
        volumeMounts:
        - name: config
          mountPath: /app/config.yaml
          subPath: config.yaml
      volumes:
      - name: config
        configMap:
          name: api-gateway-config
```

### 19.2 CI/CD Pipeline

```yaml
# .github/workflows/ci-cd.yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.21'
        
    - name: Run linter
      run: |
        go install golang.org/x/lint/golint@latest
        golint -set_exit_status ./...
        
    - name: Run static analysis
      run: |
        go install golang.org/x/tools/cmd/stringer@latest
        go vet ./...
        
    - name: Check code formatting
      run: |
        gofmt -l .
        if [ -n "$(gofmt -l .)" ]; then
          echo "Code is not formatted"
          exit 1
        fi

  test:
    name: Test
    runs-on: ubuntu-latest
    services:
postgres:
image: postgres:15-alpine
env:
POSTGRES_USER: test
POSTGRES_PASSWORD: test
POSTGRES_DB: test
ports:
- 51003:5432 # Port Sovereignty Compliance (Rule -9): 5432→51003
options: >-
--health-cmd pg_isready
--health-interval 10s
--health-timeout 5s
--health-retries 5
redis:
image: redis:7-alpine
ports:
- 51004:6379 # Port Sovereignty Compliance (Rule -9): 6379→51004
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.21'
        
    - name: Run tests with coverage
      run: |
        go test -v -race -coverprofile=coverage.out ./...
        
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        files: ./coverage.out

  build:
    name: Build
    runs-on: ubuntu-latest
    needs: [lint, test]
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v3
      
    - name: Login to Container Registry
      uses: docker/login-action@v3
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}
        
    - name: Extract metadata
      id: meta
      uses: docker/metadata-action@v5
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
        tags: |
          type=ref,event=branch
          type=sha,prefix={{branch}}-
          type=raw,value={{current_date}}
          
    - name: Build and push API
      uses: docker/build-push-action@v5
      with:
        context: .
        file: api/Dockerfile
        push: true
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}
        cache-from: type=gha
        cache-to: type=gha,mode=max

  deploy:
    name: Deploy
    runs-on: ubuntu-latest
    needs: [build]
    if: github.ref == 'refs/heads/main'
    environment:
      name: production
    steps:
    - name: Deploy to Kubernetes
      uses: azure/k8s-set-context@v1
      with:
        kubeconfig: ${{ secrets.KUBECONFIG }}
        
    - name: Update deployment
      run: |
        kubectl set image deployment/api-gateway \
          api-gateway=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }}
        
    - name: Verify deployment
      run: |
        kubectl rollout status deployment/api-gateway
        kubectl get pods -l app=api-gateway
```

### 19.3 Rollback Strategie

```bash
#!/bin/bash
# rollback.sh

set -e

DEPLOYMENT=$1
REVISION=${2:-1}

echo "Rolling back $DEPLOYMENT to revision $REVISION"

# Get current revision
CURRENT_REVISION=$(kubectl rollout history deployment/$DEPLOYMENT | awk 'NR>1 {print $1}' | tail -1)
echo "Current revision: $CURRENT_REVISION"

if [ "$REVISION" == "previous" ]; then
    kubectl rollout undo deployment/$DEPLOYMENT
else
    kubectl rollout undo deployment/$DEPLOYMENT --to-revision=$REVISION
fi

# Wait for rollback
kubectl rollout status deployment/$DEPLOYMENT --timeout=300s

# Verify
kubectl get pods -l app=$DEPLOYMENT

echo "Rollback completed successfully"
```

---

## 20) Disaster Recovery - RTO/RPO und Business Continuity

### 20.1 RTO/RPO Targets

| Tier | RTO | RPO | Datenverlust | Kosten pro Downtime |
|------|-----|-----|--------------|---------------------|
| Critical (AI Core) | 15 min | 1 min | < 1 min Transaktionsdaten | €10,000/Stunde |
| High (API, Auth) | 30 min | 5 min | < 5 min Transaktionsdaten | €5,000/Stunde |
| Medium (Workers) | 2 hours | 15 min | < 15 min Daten | €1,000/Stunde |
| Low (Analytics) | 4 hours | 1 hour | < 1 Stunde Daten | €100/Stunde |

### 20.2 Failover Mechanismen

```go
// disaster/failover.go
package disaster

import (
    "context"
    "fmt"
    "time"
    
    "github.com/redis/go-redis/v9"
)

type FailoverManager struct {
    redis         *redis.Client
    serviceName   string
    primaryHost   string
    replicaHost   string
    healthCheck   HealthChecker
    switchChan    chan<- FailoverEvent
}

type FailoverEvent struct {
    Service   string    `json:"service"`
    FromHost  string    `json:"from_host"`
    ToHost    string    `json:"to_host"`
    Reason    string    `json:"reason"`
    Timestamp time.Time `json:"timestamp"`
}

func (fm *FailoverManager) Start(ctx context.Context) error {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
            if err := fm.checkAndFailover(ctx); err != nil {
                fmt.Printf("Failover check failed: %v\n", err)
            }
        }
    }
}

func (fm *FailoverManager) checkAndFailover(ctx context.Context) error {
    // Check primary health
    primaryHealthy := fm.healthCheck.Check(fm.primaryHost)
    replicaHealthy := fm.healthCheck.Check(fm.replicaHost)
    
    // Get current primary from Redis
    currentPrimary, err := fm.redis.Get(ctx, fmt.Sprintf("failover:%s:primary", fm.serviceName)).Result()
    if err != nil && err != redis.Nil {
        return err
    }
    
    // Determine if failover is needed
    shouldFailover := false
    reason := ""
    
    if currentPrimary == fm.primaryHost && !primaryHealthy && replicaHealthy {
        shouldFailover = true
        reason = "primary_unhealthy"
    }
    
    if !shouldFailover {
        return nil
    }
    
    // Perform failover
    newPrimary := fm.replicaHost
    if currentPrimary == fm.replicaHost {
        newPrimary = fm.primaryHost
    }
    
    // Update Redis
    pipe := fm.redis.Pipeline()
    pipe.Set(ctx, fmt.Sprintf("failover:%s:primary", fm.serviceName), newPrimary, 0)
    pipe.Set(ctx, fmt.Sprintf("failover:%s:last_switch", fm.serviceName), time.Now().Unix(), 0)
    _, err = pipe.Exec(ctx)
    if err != nil {
        return fmt.Errorf("failover update failed: %w", err)
    }
    
    // Notify
    fm.switchChan <- FailoverEvent{
        Service:   fm.serviceName,
        FromHost:  currentPrimary,
        ToHost:    newPrimary,
        Reason:    reason,
        Timestamp: time.Now(),
    }
    
    return nil
}

func (fm *FailoverManager) GetPrimary(ctx context.Context) (string, error) {
    primary, err := fm.redis.Get(ctx, fmt.Sprintf("failover:%s:primary", fm.serviceName)).Result()
    if err == redis.Nil {
        return fm.primaryHost, nil
    }
    return primary, err
}
```

### 20.3 Datenrettung

```bash
#!/bin/bash
# restore.sh

set -e

BACKUP_FILE=$1
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

echo "Starting restore from $BACKUP_FILE"

# Stop services
echo "Stopping services..."
kubectl scale deployment api-gateway --replicas=0
kubectl scale deployment ai-orchestrator --replicas=0

# Restore database
echo "Restoring database..."
docker exec biometrics-postgres-1 psql -U biometrics -c "DROP DATABASE IF EXISTS biometrics_restore_${TIMESTAMP}"
docker exec biometrics-postgres-1 psql -U biometrics -c "CREATE DATABASE biometrics_restore_${TIMESTAMP}"
docker exec -i biometrics-postgres-1 psql -U biometrics biometrics_restore_${TIMESTAMP} < ${BACKUP_FILE}

# Verify restore
echo "Verifying restore..."
docker exec biometrics-postgres-1 psql -U biometrics biometrics_restore_${TIMESTAMP} -c "SELECT count(*) FROM users"

# Rename databases
echo "Swapping databases..."
docker exec biometrics-postgres-1 psql -U biometrics -c "ALTER DATABASE biometrics RENAME TO biometrics_backup_${TIMESTAMP}"
docker exec biometrics-postgres-1 psql -U biometrics -c "ALTER DATABASE biometrics_restore_${TIMESTAMP} RENAME TO biometrics"

# Restart services
echo "Restarting services..."
kubectl scale deployment api-gateway --replicas=3
kubectl scale deployment ai-orchestrator --replicas=2

echo "Restore completed successfully"
```

---

## 21) Cost Optimization - Ressourcen und Cloud-Kosten

### 21.1 Ressourcenallokation

```yaml
# kubernetes/resource-quotas.yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: biometrics-quota
spec:
  hard:
    requests.cpu: "32"
    requests.memory: 64Gi
    limits.cpu: "64"
    limits.memory: 128Gi
    persistentvolumeclaims: "10"
    services.loadbalancers: "2"
    services.nodeports: "0"
---
apiVersion: v1
kind: LimitRange
metadata:
  name: biometrics-limits
spec:
  limits:
  - max:
      cpu: "4"
      memory: "8Gi"
    min:
      cpu: "50m"
      memory: "64Mi"
    default:
      cpu: "500m"
      memory: "512Mi"
    defaultRequest:
      cpu: "100m"
      memory: "128Mi"
    type: Container
```

### 21.2 Kosten-Optimierungsstrategien

```go
// cost/optimizer.go
package cost

import (
    "context"
    "time"
    
    "github.com/prometheus/client_golang/api"
    "github.com/prometheus/client_golang/api/prometheus/v1"
    "github.com/prometheus/common/model"
)

type CostOptimizer struct {
    prometheusAPI v1.API
    cloudProvider CloudProvider
}

type CostReport struct {
    Service          string  `json:"service"`
    CurrentCost      float64 `json:"current_cost"`
    ProjectedMonthly float64 `json:"projected_monthly"`
    Recommendations  []Recommendation `json:"recommendations"`
}

type Recommendation struct {
    Type        string  `json:"type"`
    Description string  `json:"description"`
    Savings     float64 `json:"savings"`
    Effort      string  `json:"effort"`
}

func (co *CostOptimizer) Analyze(ctx context.Context) ([]CostReport, error) {
    var reports []CostReport
    
    // Analyze compute costs
    computeReport, err := co.analyzeComputeCosts(ctx)
    if err != nil {
        return nil, err
    }
    reports = append(reports, computeReport)
    
    // Analyze storage costs
    storageReport, err := co.analyzeStorageCosts(ctx)
    if err != nil {
        return nil, err
    }
    reports = append(reports, storageReport)
    
    // Analyze network costs
    networkReport, err := co.analyzeNetworkCosts(ctx)
    if err != nil {
        return nil, err
    }
    reports = append(reports, networkReport)
    
    return reports, nil
}

func (co *CostOptimizer) analyzeComputeCosts(ctx context.Context) (CostReport, error) {
    query := `sum(rate(container_cpu_usage_seconds_total[24h])) by (pod)`
    result, _, err := co.prometheusAPI.Query(ctx, query, time.Now())
    if err != nil {
        return CostReport{}, err
    }
    
    vector := result.(model.Vector)
    
    // Calculate current cost based on CPU usage
    // Assume $0.04 per CPU-hour
    var totalCPUHours float64
    for _, sample := range vector {
        cpuCores := float64(sample.Value) / 3600
        totalCPUHours += cpuCores
    }
    
    dailyCost := totalCPUHours * 0.04
    
    report := CostReport{
        Service:          "compute",
        CurrentCost:      dailyCost,
        ProjectedMonthly: dailyCost * 30,
        Recommendations: []Recommendation{
            {
                Type:        "rightsize",
                Description: "3 pods have consistently low CPU usage (<20%)",
                Savings:     dailyCost * 0.15 * 30,
                Effort:      "low",
            },
            {
                Type:        "spot",
                Description: "Use spot instances for non-critical workers",
                Savings:     dailyCost * 0.6 * 30,
                Effort:      "medium",
            },
        },
    }
    
    return report, nil
}
```

---

## 22) Future Roadmap - Phase 1-3 Pläne

### 22.1 Phase 1: Foundation (Q1 2026)

**Monate 1-3: Kerninfrastruktur**

| Woche | Ziel | Meilenstein |
|-------|------|-------------|
| 1-2 | Go API Framework aufsetzen | Basis-API mit CRUD-Operationen |
| 3-4 | PostgreSQL Schema implementieren | Migration erfolgreich |
| 5-6 | Auth-Service implementieren | JWT + MFA funktioniert |
| 7-8 | Redis Caching integrieren | Cache-Hit-Rate > 80% |
| 9-10 | API Gateway aufsetzen | Routing + Rate Limiting |
| 11-12 | Kubernetes Deployment | Produktionsreife |

**Erfolgskriterien Phase 1:**
- [ ] API Latency p95 < 200ms
- [ ] 99.9% Verfügbarkeit
- [ ] Alle Tests bestanden
- [ ] Dokumentation vollständig

### 22.2 Phase 2: AI Integration (Q2 2026)

**Monate 4-6: KI-Fähigkeiten**

| Woche | Ziel | Meilenstein |
|-------|------|-------------|
| 13-14 | Qwen 3.5 Integration | First AI Request erfolgreich |
| 15-16 | Prompt Engineering Pipeline | Template-System |
| 17-18 | Multi-Model Support | Kimi + Claude Fallback |
| 19-20 | AI Response Caching | Kostensenkung |
| 21-22 | Fine-tuning Pipeline | Custom Models |
| 23-24 | AI Analytics Dashboard | Usage Tracking |

**Erfolgskriterien Phase 2:**
- [ ] AI Latency < 60s (p95)
- [ ] Fallback-Chain funktioniert
- [ ] Kosten um 50% reduziert vs. Phase 1

### 22.3 Phase 3: Scale (Q3-Q4 2026)

**Monate 7-12: Skalierung und Optimierung**

| Quartal | Ziel | Meilenstein |
|---------|------|-------------|
| Q3 | Multi-Region Deployment | EU + US Regionen |
| Q3 | Advanced Security | SOC2 Vorbereitung |
| Q3 | Advanced Analytics | Business Intelligence |
| Q4 | Enterprise Features | SSO + Audit Logs |
| Q4 | Platform API | Partner-Ökosystem |
| Q4 | Automation Suite | No-Code Integrationen |

**Erfolgskriterien Phase 3:**
- [ ] 99.99% Verfügbarkeit
- [ ] Multi-Region Failover < 30s
- [ ] 100+ Enterprise Customers
- [ ] SOC2 Zertifizierung

### 22.4 Technology Evolution

**Geplante Technologie-Updates:**

| Zeitpunkt | Technologie | Motivation | Aufwand |
|-----------|-------------|------------|---------|
| Q2 2026 | Go 1.22 | Performance | 1 Sprint |
| Q3 2026 | Kubernetes 1.30 | Features | 2 Sprints |
| Q4 2026 | PostgreSQL 17 | Performance | 1 Sprint |
| Q1 2027 | Rust Services | Performance-Critical | 4 Sprints |
| Q2 2027 | WebAssembly | Edge Computing | 3 Sprints |

---

## 23) Anhang: Architecture Decision Records (ADRs)

### ADR-001: Frontend-Stack

**Status:** Accepted

**Entscheidung:** Next.js als primäres Frontend-Framework

**Kontext:** Das Projekt erfordert ein performantes, SEO-fähiges Frontend mit guter Entwicklerproduktivität.

**Entscheidung:** Next.js 14+ mit App Router

**Konsequenzen:**
- Positive: SSR für SEO, React 18 Features, große Community
- Negative: Serverless-Kosten bei hohem Traffic

### ADR-002: Backend-Stack

**Status:** Accepted

**Entscheidung:** Go + Supabase als Backend-Stack

**Kontext:** Der Backend muss performant, skalierbar und einfach zu warten sein.

**Konsequenzen:**
- Positive: Go ist schnell und einfach, Supabase beschleunigt Entwicklung
- Negative: Weniger Flexibilität als Custom-SQL

### ADR-003: Content-Generierung

**Status:** Accepted

**Entscheidung:** NLM-CLI Pflicht für alle Content-Artefakte

**Kontext:** Konsistente Qualität und Nachvollziehbarkeit erforderlich.

**Konsequenzen:**
- Positive: Einheitliche Qualität, bessere Kontrolle
- Negative: Extra Schritt im Workflow

### ADR-004: Primary AI Brain

**Status:** Accepted

**Entscheidung:** Qwen 3.5 (NVIDIA NIM) als Primary AI Brain

**Kontext:** Beste Reasoning-Fähigkeiten im Kosten-Nutzen-Verhältnis.

**Konsequenzen:**
- Positive: 262K Context, multimodal
- Negative: 70-90s Latenz

### ADR-005: Timeout-Konfiguration

**Status:** Accepted

**Entscheidung:** 120s Timeout für AI-Requests

**Kontext:** Qwen 3.5 benötigt 70-90s für vollständige Responses.

**Konsequenzen:**
- Positive: Keine Timeouts bei langen Anfragen
- Negative: User muss länger warten

### ADR-006: Fallback-Kette

**Status:** Accepted

**Entscheidung:** Qwen → Kimi → Claude Fallback-Kette

**Kontext:** Resilienz bei Provider-Ausfällen erforderlich.

**Konsequenzen:**
- Positive: Hohe Verfügbarkeit
- Negative: Potentiell inkonsistente Antworten

### ADR-007: Caching-Strategie

**Status:** Accepted

**Entscheidung:** Redis für API + AI Response Caching

**Kontext:** Performanz und Kostensenkung durch Caching.

**Konsequenzen:**
- Positive: 80%+ Cache-Hit-Rate möglich
- Negative: Komplexität bei Cache-Invalidation

### ADR-008: Deployment-Plattform

**Status:** Accepted

**Entscheidung:** Kubernetes auf Google Cloud Platform

**Kontext:** Skalierbarkeit und Verwaltungsaufwand.

**Konsequenzen:**
- Positive: Automatische Skalierung, gute GCP-Integration
- Negative: Komplexität, Kosten bei Peak

---

## 27) Detaillierte API-Spezifikationen

### 27.1 REST API Endpoints - User Management

| Method | Endpoint | Beschreibung | Auth | Rate Limit |
|--------|----------|-------------|------|------------|
| POST | /api/v1/users | Create new user | Admin | 10/min |
| GET | /api/v1/users | List users | Admin | 60/min |
| GET | /api/v1/users/:id | Get user by ID | User/Admin | 120/min |
| PUT | /api/v1/users/:id | Update user | User/Admin | 30/min |
| DELETE | /api/v1/users/:id | Delete user | Admin | 5/min |
| POST | /api/v1/users/:id/verify | Verify user email | - | 3/min |
| POST | /api/v1/users/:id/lock | Lock user account | Admin | 10/min |
| POST | /api/v1/users/:id/unlock | Unlock user account | Admin | 10/min |

### 27.2 REST API Endpoints - Content Management

| Method | Endpoint | Beschreibung | Auth | Rate Limit |
|--------|----------|-------------|------|------------|
| POST | /api/v1/content | Create content | User | 30/min |
| GET | /api/v1/content | List content | User | 120/min |
| GET | /api/v1/content/:id | Get content | User | 120/min |
| PUT | /api/v1/content/:id | Update content | User | 30/min |
| DELETE | /api/v1/content/:id | Delete content | User | 20/min |
| POST | /api/v1/content/:id/publish | Publish content | User | 10/min |
| POST | /api/v1/content/:id/unpublish | Unpublish content | User | 10/min |
| GET | /api/v1/content/search | Search content | User | 60/min |

### 27.3 REST API Endpoints - AI Services

| Method | Endpoint | Beschreibung | Auth | Rate Limit |
|--------|----------|-------------|------|------------|
| POST | /api/v1/ai/generate | Generate with AI | User | 10/min |
| POST | /api/v1/ai/analyze | Analyze with AI | User | 10/min |
| POST | /api/v1/ai/chat | Chat with AI | User | 30/min |
| POST | /api/v1/ai/vision | Vision analysis | User | 20/min |
| GET | /api/v1/ai/models | List available models | User | 60/min |
| GET | /api/v1/ai/usage | Get AI usage stats | User | 30/min |

### 27.4 REST API Endpoints - Integrations

| Method | Endpoint | Beschreibung | Auth | Rate Limit |
|--------|----------|-------------|------|------------|
| GET | /api/v1/integrations | List integrations | User | 60/min |
| POST | /api/v1/integrations/:provider/connect | Connect integration | User | 5/min |
| DELETE | /api/v1/integrations/:provider/disconnect | Disconnect integration | User | 5/min |
| GET | /api/v1/integrations/:provider/status | Get integration status | User | 60/min |
| POST | /api/v1/integrations/:provider/sync | Sync integration data | User | 3/min |

### 27.5 REST API Endpoints - Workflows

| Method | Endpoint | Beschreibung | Auth | Rate Limit |
|--------|----------|-------------|------|------------|
| GET | /api/v1/workflows | List workflows | User | 60/min |
| POST | /api/v1/workflows | Create workflow | User | 10/min |
| GET | /api/v1/workflows/:id | Get workflow | User | 60/min |
| PUT | /api/v1/workflows/:id | Update workflow | User | 10/min |
| DELETE | /api/v1/workflows/:id | Delete workflow | User | 5/min |
| POST | /api/v1/workflows/:id/trigger | Trigger workflow | User | 20/min |
| GET | /api/v1/workflows/:id/executions | List executions | User | 60/min |
| GET | /api/v1/workflows/:id/executions/:exec_id | Get execution | User | 60/min |

---

## 28) Detaillierte Datenmodelle

### 28.1 User Model

```go
// pkg/models/user.go
package models

import (
    "time"
    
    "github.com/google/uuid"
)

type User struct {
    ID              uuid.UUID  `json:"id" db:"id"`
    Email           string     `json:"email" db:"email"`
    PasswordHash    string     `json:"-" db:"password_hash"`
    Name            string     `json:"name" db:"name"`
    Role            UserRole   `json:"role" db:"role"`
    EmailVerified   bool       `json:"email_verified" db:"email_verified"`
    MFAEnabled      bool       `json:"mfa_enabled" db:"mfa_enabled"`
    MFA secret      string     `json:"-" db:"mfa_secret"`
    LastLoginAt     *time.Time `json:"last_login_at" db:"last_login_at"`
    FailedLogins    int        `json:"-" db:"failed_logins"`
    LockedUntil     *time.Time `json:"locked_until" db:"locked_until"`
    CreatedAt       time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
    DeletedAt       *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type UserRole string

const (
    RoleSuperAdmin UserRole = "superadmin"
    RoleAdmin      UserRole = "admin"
    RoleUser       UserRole = "user"
    RoleGuest      UserRole = "guest"
)

type UserProfile struct {
    UserID      uuid.UUID `json:"user_id" db:"user_id"`
    AvatarURL   *string   `json:"avatar_url" db:"avatar_url"`
    Bio         *string   `json:"bio" db:"bio"`
    Timezone    string    `json:"timezone" db:"timezone"`
    Locale      string    `json:"locale" db:"locale"`
    Preferences JSON      `json:"preferences" db:"preferences"`
    Metadata    JSON      `json:"metadata" db:"metadata"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type UserFilter struct {
    Page       int
    PageSize   int
    Search     *string
    Role       *UserRole
    Verified   *bool
    MFAEnabled *bool
    CreatedAfter *time.Time
    CreatedBefore *time.Time
    SortBy     string
    SortOrder  string
}

func (u *User) WithoutPassword() *User {
    copy := *u
    copy.PasswordHash = ""
    return &copy
}

func (u *User) IsLocked() bool {
    if u.LockedUntil == nil {
        return false
    }
    return time.Now().Before(*u.LockedUntil)
}
```

### 28.2 Content Model

```go
// pkg/models/content.go
package models

import (
    "time"
    
    "github.com/google/uuid"
)

type Content struct {
    ID          uuid.UUID   `json:"id" db:"id"`
    UserID      uuid.UUID   `json:"user_id" db:"user_id"`
    Title       string      `json:"title" db:"title"`
    Body        *string     `json:"body" db:"body"`
    Excerpt     *string     `json:"excerpt" db:"excerpt"`
    ContentType ContentType `json:"content_type" db:"content_type"`
    Status      ContentStatus `json:"status" db:"status"`
    Featured    bool        `json:"featured" db:"featured"`
    ViewCount   int         `json:"view_count" db:"view_count"`
    LikeCount   int         `json:"like_count" db:"like_count"`
    ShareCount  int         `json:"share_count" db:"share_count"`
    Metadata    JSON        `json:"metadata" db:"metadata"`
    SEO         *SEOData    `json:"seo,omitempty" db:"-"`
    PublishedAt *time.Time `json:"published_at,omitempty" db:"published_at"`
    CreatedAt   time.Time   `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
    DeletedAt   *time.Time  `json:"deleted_at,omitempty" db:"deleted_at"`
}

type ContentType string

const (
    ContentTypeArticle  ContentType = "article"
    ContentTypePage     ContentType = "page"
    ContentTypePost    ContentType = "post"
    ContentTypeProduct ContentType = "product"
)

type ContentStatus string

const (
    ContentStatusDraft     ContentStatus = "draft"
    ContentStatusReview    ContentStatus = "review"
    ContentStatusPublished ContentStatus = "published"
    ContentStatusArchived  ContentStatus = "archived"
)

type SEOData struct {
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Keywords    []string `json:"keywords"`
    OGImage     string   `json:"og_image"`
    Canonical   string   `json:"canonical"`
    NoIndex     bool     `json:"no_index"`
}

type ContentFilter struct {
    Page        int
    PageSize    int
    UserID      *uuid.UUID
    Type        *ContentType
    Status      *ContentStatus
    Featured    *bool
    Search      *string
    Tags        []string
    PublishedAfter  *time.Time
    PublishedBefore *time.Time
    SortBy      string
    SortOrder   string
}
```

### 28.3 Integration Model

```go
// pkg/models/integration.go
package models

import (
    "time"
    
    "github.com/google/uuid"
)

type Integration struct {
    ID              uuid.UUID      `json:"id" db:"id"`
    UserID          uuid.UUID      `json:"user_id" db:"user_id"`
    Provider        IntegrationProvider `json:"provider" db:"provider"`
    AccessToken     string         `json:"-" db:"access_token"`
    RefreshToken    string         `json:"-" db:"refresh_token"`
    TokenType       string         `json:"token_type" db:"token_type"`
    ExpiresAt       *time.Time     `json:"expires_at,omitempty" db:"expires_at"`
    Scope           string         `json:"scope" db:"scope"`
    Metadata        JSON           `json:"metadata" db:"metadata"`
    Status          IntegrationStatus `json:"status" db:"status"`
    LastSyncAt     *time.Time     `json:"last_sync_at" db:"last_sync_at"`
    CreatedAt       time.Time      `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

type IntegrationProvider string

const (
    ProviderGoogle     IntegrationProvider = "google"
    ProviderGitHub     IntegrationProvider = "github"
    ProviderSlack      IntegrationProvider = "slack"
    ProviderStripe     IntegrationProvider = "stripe"
    ProviderOpenAI    IntegrationProvider = "openai"
    ProviderNotion    IntegrationProvider = "notion"
    ProviderFigma     IntegrationProvider = "figma"
    ProviderAirtable  IntegrationProvider = "airtable"
)

type IntegrationStatus string

const (
    IntegrationStatusActive    IntegrationStatus = "active"
    IntegrationStatusExpired   IntegrationStatus = "expired"
    IntegrationStatusRevoked   IntegrationStatus = "revoked"
    IntegrationStatusError    IntegrationStatus = "error"
)

type IntegrationConfig struct {
    Provider       IntegrationProvider
    AuthURL        string
    TokenURL       string
    Scopes         []string
    RedirectURI    string
    ClientID       string
    ClientSecret   string
}
```

### 28.4 Workflow Model

```go
// pkg/models/workflow.go
package models

import (
    "time"
    
    "github.com/google/uuid"
)

type Workflow struct {
    ID          uuid.UUID       `json:"id" db:"id"`
    UserID      uuid.UUID       `json:"user_id" db:"user_id"`
    Name        string          `json:"name" db:"name"`
    Description *string         `json:"description" db:"description"`
    Definition  WorkflowDefinition `json:"definition" db:"definition"`
    TriggerType WorkflowTrigger  `json:"trigger_type" db:"trigger_type"`
    TriggerConfig JSON           `json:"trigger_config" db:"trigger_config"`
    IsActive    bool            `json:"is_active" db:"is_active"`
    LastRunAt   *time.Time      `json:"last_run_at" db:"last_run_at"`
    RunCount    int             `json:"run_count" db:"run_count"`
    SuccessCount int            `json:"success_count" db:"success_count"`
    ErrorCount  int             `json:"error_count" db:"error_count"`
    CreatedAt   time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

type WorkflowTrigger string

const (
    TriggerManual    WorkflowTrigger = "manual"
    TriggerScheduled WorkflowTrigger = "scheduled"
    TriggerWebhook   WorkflowTrigger = "webhook"
    TriggerEvent     WorkflowTrigger = "event"
)

type WorkflowDefinition struct {
    Nodes     []WorkflowNode     `json:"nodes"`
    Edges    []WorkflowEdge     `json:"edges"`
    Variables map[string]string  `json:"variables"`
}

type WorkflowNode struct {
    ID       string          `json:"id"`
    Type     string          `json:"type"`
    Position Position        `json:"position"`
    Data     JSON            `json:"data"`
}

type WorkflowEdge struct {
    ID        string `json:"id"`
    Source    string `json:"source"`
    Target    string `json:"target"`
    Condition *string `json:"condition,omitempty"`
}

type Position struct {
    X float64 `json:"x"`
    Y float64 `json:"y"`
}
```

---

## 29) Error Codes und Fehlerbehandlung

### 29.1 HTTP Status Codes

| Code | Name | Beschreibung | Verwendung |
|------|------|--------------|------------|
| 200 | OK | Erfolgreiche Anfrage | GET, PUT, DELETE erfolgreich |
| 201 | Created | Ressource erstellt | POST erfolgreich |
| 204 | No Content | Kein Inhalt | DELETE erfolgreich |
| 400 | Bad Request | Ungültige Anfrage | Validierungsfehler |
| 401 | Unauthorized | Nicht authentifiziert | Fehlende/ungültiges Token |
| 403 | Forbidden | Keine Berechtigung | Fehlende Rechte |
| 404 | Not Found | Ressource nicht gefunden | ID existiert nicht |
| 409 | Conflict | Konflikt | Duplicate, Version mismatch |
| 422 | Unprocessable Entity | Verarbeitungsfehler | Business-Logik Fehler |
| 429 | Too Many Requests | Rate Limit überschritten | Rate Limit |
| 500 | Internal Server Error | Serverfehler | Unerwarteter Fehler |
| 502 | Bad Gateway | Externer Dienst Fehler | Upstream Fehler |
| 503 | Service Unavailable | Dienst nicht verfügbar | Wartung/Überlastung |
| 504 | Gateway Timeout | Timeout | Upstream Timeout |

### 29.2 Application Error Codes

```go
// pkg/errors/codes.go
package errors

type ErrorCode string

const (
    // Authentication errors (AUTH_*)
    ErrCodeInvalidCredentials ErrorCode = "AUTH_INVALID_CREDENTIALS"
    ErrCodeTokenExpired      ErrorCode = "AUTH_TOKEN_EXPIRED"
    ErrCodeTokenInvalid      ErrorCode = "AUTH_TOKEN_INVALID"
    ErrCodeMFARequired       ErrorCode = "AUTH_MFA_REQUIRED"
    ErrCodeMFAInvalid       ErrorCode = "AUTH_MFA_INVALID"
    ErrCodeAccountLocked    ErrorCode = "AUTH_ACCOUNT_LOCKED"
    ErrCodeAccountDisabled  ErrorCode = "AUTH_ACCOUNT_DISABLED"
    ErrCodeEmailNotVerified ErrorCode = "AUTH_EMAIL_NOT_VERIFIED"
    
    // Authorization errors (AUTHZ_*)
    ErrCodeForbidden        ErrorCode = "AUTHZ_FORBIDDEN"
    ErrCodeInsufficientPerms ErrorCode = "AUTHZ_INSUFFICIENT_PERMISSIONS"
    ErrCodeResourceOwnership ErrorCode = "AUTHZ_RESOURCE_OWNERSHIP"
    
    // Validation errors (VAL_*)
    ErrCodeValidationFailed ErrorCode = "VAL_VALIDATION_FAILED"
    ErrCodeInvalidFormat    ErrorCode = "VAL_INVALID_FORMAT"
    ErrCodeRequiredField    ErrorCode = "VAL_REQUIRED_FIELD"
    ErrCodeFieldTooLong    ErrorCode = "VAL_FIELD_TOO_LONG"
    ErrCodeFieldTooShort   ErrorCode = "VAL_FIELD_TOO_SHORT"
    
    // Resource errors (RES_*)
    ErrCodeNotFound        ErrorCode = "RES_NOT_FOUND"
    ErrCodeAlreadyExists   ErrorCode = "RES_ALREADY_EXISTS"
    ErrCodeAlreadyDeleted  ErrorCode = "RES_ALREADY_DELETED"
    ErrCodeConflict       ErrorCode = "RES_CONFLICT"
    
    // AI errors (AI_*)
    ErrCodeAIRequestFailed  ErrorCode = "AI_REQUEST_FAILED"
    ErrCodeAIQuotaExceeded ErrorCode = "AI_QUOTA_EXCEEDED"
    ErrCodeAIModelNotFound ErrorCode = "AI_MODEL_NOT_FOUND"
    ErrCodeAITimeout       ErrorCode = "AI_TIMEOUT"
    ErrCodeAIInvalidResponse ErrorCode = "AI_INVALID_RESPONSE"
    
    // Integration errors (INT_*)
    ErrCodeIntegrationFailed   ErrorCode = "INT_INTEGRATION_FAILED"
    ErrCodeIntegrationNotFound  ErrorCode = "INT_INTEGRATION_NOT_FOUND"
    ErrCodeIntegrationExpired   ErrorCode = "INT_INTEGRATION_EXPIRED"
    ErrCodeIntegrationRevoked   ErrorCode = "INT_INTEGRATION_REVOKED"
    
    // Rate limit errors (RATE_*)
    ErrCodeRateLimitExceeded ErrorCode = "RATE_LIMIT_EXCEEDED"
    ErrCodeRateLimitBudget   ErrorCode = "RATE_LIMIT_BUDGET"
    
    // Internal errors (INT_*)
    ErrCodeInternalError    ErrorCode = "INT_INTERNAL_ERROR"
    ErrCodeDatabaseError   ErrorCode = "INT_DATABASE_ERROR"
    ErrCodeCacheError      ErrorCode = "INT_CACHE_ERROR"
    ErrCodeExternalService ErrorCode = "INT_EXTERNAL_SERVICE_ERROR"
)
```

### 29.3 Error Response Format

```json
{
  "error": {
    "code": "AUTH_INVALID_CREDENTIALS",
    "message": "Invalid email or password",
    "details": {
      "field": "password",
      "reason": "must be at least 8 characters"
    },
    "trace_id": "abc123def456",
    "timestamp": "2026-02-18T12:34:56Z"
  }
}
```

```go
// pkg/errors/response.go
package errors

import (
    "encoding/json"
    "time"
    
    "github.com/google/uuid"
)

type ErrorResponse struct {
    Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
    Code      ErrorCode              `json:"code"`
    Message   string                 `json:"message"`
    Details   map[string]interface{} `json:"details,omitempty"`
    TraceID   string                 `json:"trace_id"`
    Timestamp time.Time              `json:"timestamp"`
}

func NewErrorResponse(code ErrorCode, message string) *ErrorResponse {
    return &ErrorResponse{
        Error: ErrorDetail{
            Code:      code,
            Message:   message,
            TraceID:   uuid.New().String(),
            Timestamp: time.Now().UTC(),
        },
    }
}

func (e *ErrorResponse) WithDetails(details map[string]interface{}) *ErrorResponse {
    e.Error.Details = details
    return e
}

func (e *ErrorResponse) ToJSON() []byte {
    data, _ := json.Marshal(e)
    return data
}
```

---

## 30) Testing Strategy

### 30.1 Test Pyramid

```
         ┌─────────────┐
         │     E2E     │  ← 10%
         │   Tests     │
        ┌──────────────┐
        │  Integration │  ← 20%
        │   Tests     │
       ┌───────────────┐
       │  Unit Tests  │  ← 70%
       │              │
      └───────────────┘
```

### 30.2 Unit Tests

```go
// internal/auth/service_test.go
package auth

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *User) (*User, error) {
    args := m.Called(ctx, user)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
    // Arrange
    ctx := context.Background()
    mockRepo := new(MockUserRepository)
    authService := NewAuthService(mockRepo, &Config{JWTExpiry: time.Hour})
    
    mockRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, ErrNotFound)
    mockRepo.On("Create", ctx, mock.AnythingOfType("*User")).Return(&User{
        ID:    "test-id",
        Email: "test@example.com",
        Name:  "Test User",
    }, nil)
    
    req := RegisterRequest{
        Email:    "test@example.com",
        Password: "password123",
        Name:     "Test User",
    }
    
    // Act
    response, err := authService.Register(ctx, req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, response.AccessToken)
    assert.Equal(t, "test@example.com", response.User.Email)
    mockRepo.AssertExpectations(t)
}

func TestRegister_EmailExists(t *testing.T) {
    // Arrange
    ctx := context.Background()
    mockRepo := new(MockUserRepository)
    authService := NewAuthService(mockRepo, &Config{JWTExpiry: time.Hour})
    
    existingUser := &User{
        ID:    "existing-id",
        Email: "test@example.com",
    }
    mockRepo.On("GetByEmail", ctx, "test@example.com").Return(existingUser, nil)
    
    req := RegisterRequest{
        Email:    "test@example.com",
        Password: "password123",
        Name:     "Test User",
    }
    
    // Act
    response, err := authService.Register(ctx, req)
    
    // Assert
    assert.Error(t, err)
    assert.Equal(t, ErrUserAlreadyExists, err)
    assert.Nil(t, response)
    mockRepo.AssertExpectations(t)
}

func TestRegister_WeakPassword(t *testing.T) {
    // Arrange
    ctx := context.Background()
    mockRepo := new(MockUserRepository)
    authService := NewAuthService(mockRepo, &Config{JWTExpiry: time.Hour})
    
    req := RegisterRequest{
        Email:    "test@example.com",
        Password: "weak",
        Name:     "Test User",
    }
    
    // Act
    response, err := authService.Register(ctx, req)
    
    // Assert
    assert.Error(t, err)
    assert.Equal(t, ErrPasswordTooWeak, err)
    assert.Nil(t, response)
}
```

### 30.3 Integration Tests

```go
// integration/api_test.go
package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "biometrics/internal/api"
    "biometrics/internal/config"
)

func TestAPI_HealthEndpoint(t *testing.T) {
    // Setup
    cfg := &config.Config{}
    app := api.NewApp(cfg)
    
    // Execute
    req, _ := http.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()
    app.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "healthy", response["status"])
}

func TestAPI_AuthFlow(t *testing.T) {
    // Setup
    cfg := &config.Config{}
    app := api.NewApp(cfg)
    
    // Test register
    registerReq := map[string]string{
        "email":    "test@example.com",
        "password": "password123",
        "name":     "Test User",
    }
    body, _ := json.Marshal(registerReq)
    
    req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    app.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
    
    // Test login
    loginReq := map[string]string{
        "email":    "test@example.com",
        "password": "password123",
    }
    body, _ = json.Marshal(loginReq)
    
    req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w = httptest.NewRecorder()
    app.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    // Extract token
    var loginResp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &loginResp)
    token := loginResp["access_token"].(string)
    
    // Test authenticated endpoint
    req, _ = http.NewRequest("GET", "/api/v1/users/me", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    w = httptest.NewRecorder()
    app.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

---

## 31) Monitoring und Alerting Regeln

### 31.1 Prometheus Alerting Rules

```yaml
# prometheus/alerts.yaml
groups:
- name: biometrics-alerts
  interval: 30s
  rules:
  # High Error Rate
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "High error rate detected"
      description: "Error rate is {{ $value }} per second"
      
  # High Latency
  - alert: HighLatency
    expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High API latency"
      description: "p95 latency is {{ $value }}s"
      
  # Database Connection Pool Exhausted
  - alert: DatabaseConnectionPoolExhausted
    expr: sum(rate(db_connection_pool_exhausted_total[5m])) > 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Database connection pool exhausted"
      
  # AI Service Down
  - alert: AIServiceDown
    expr: up{job="ai-orchestrator"} == 0
    for: 2m
    labels:
      severity: critical
    annotations:
      summary: "AI orchestrator service is down"
      
  # High Memory Usage
  - alert: HighMemoryUsage
    expr: (container_memory_usage_bytes / container_spec_memory_limit_bytes) > 0.9
    for: 10m
    labels:
      severity: warning
    annotations:
      summary: "High memory usage"
      description: "Memory usage is {{ $value | humanizePercentage }}"
      
  # Disk Space Low
  - alert: DiskSpaceLow
    expr: (node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}) < 0.1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Disk space is running low"
      
  # Failed Logins
  - alert: HighFailedLogins
    expr: rate(auth_failed_logins_total[15m]) > 10
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High number of failed login attempts"
```

### 31.2 Grafana Dashboard Konfiguration

```json
{
  "dashboard": {
    "title": "BIOMETRICS Overview",
    "tags": ["biometrics", "production"],
    "timezone": "UTC",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total[5m])) by (method, status)",
            "legendFormat": "{{method}} - {{status}}"
          }
        ]
      },
      {
        "title": "Latency p50/p95/p99",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.50, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))",
            "legendFormat": "p50"
          },
          {
            "expr": "histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))",
            "legendFormat": "p95"
          },
          {
            "expr": "histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))",
            "legendFormat": "p99"
          }
        ]
      },
      {
        "title": "Active Users",
        "type": "stat",
        "targets": [
          {
            "expr": "active_users"
          }
        ]
      },
      {
        "title": "AI Request Latency",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, sum(rate(ai_request_duration_seconds_bucket[5m])) by (le, model))",
            "legendFormat": "{{model}}"
          }
        ]
      },
      {
        "title": "Database Connections",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(pg_stat_database_numbackends) by (datname)"
          }
        ]
      },
      {
        "title": "Cache Hit Rate",
        "type": "gauge",
        "targets": [
          {
            "expr": "rate(redis_keyspace_hits_total[5m]) / (rate(redis_keyspace_hits_total[5m]) + rate(redis_keyspace_misses_total[5m]))"
          }
        ]
      }
    ]
  }
}
```

---

## 32) Sicherheits-Checkliste

### 32.1 Entwicklung

- [ ] Keine Secrets im Code
- [ ] Keine hardcodierten Credentials
- [ ] Environment Variables für Konfiguration
- [ ] Input Validation auf allen Ebenen
- [ ] Output Encoding
- [ ] Prepared Statements für SQL
- [ ] Parameterized Queries
- [ ] Security Headers in Responses
- [ ] CORS Policy definiert
- [ ] Rate Limiting implementiert
- [ ] Logging von Security-Events

### 32.2 Authentifizierung

- [ ] Starke Passwort-Policy
- [ ] Passwort-Hashing mit bcrypt/Argon2
- [ ] JWT mit kurzer Gültigkeit
- [ ] Refresh Token Rotation
- [ ] MFA-Option verfügbar
- [ ] Account Lockout nach Fehlversuchen
- [ ] Session Timeout
- [ ] Sichere Cookie-Attribute

### 32.3 Autorisierung

- [ ] RBAC implementiert
- [ ] Principle of Least Privilege
- [ ] Resource-Level Authorization
- [ ] Audit Logging

### 32.4 Infrastruktur

- [ ] TLS 1.2+ everywhere
- [ ] Starke Cipher Suites
- [ ] Zertifikatsvalidierung
- [ ] Firewall-Regeln
- [ ] Network Segmentation
- [ ] Regular Security Updates
- [ ] Vulnerability Scanning

---

## 33) Performance-Optimierung

### 33.1 Database Optimization

```sql
-- Index Analyse
EXPLAIN ANALYZE 
SELECT * FROM users 
WHERE email = 'test@example.com' 
AND deleted_at IS NULL;

-- Composite Index für häufige Queries
CREATE INDEX idx_users_email_status_created 
ON users(email, deleted_at, created_at DESC);

-- Partial Index für aktive Daten
CREATE INDEX idx_contents_published 
ON contents(published_at DESC) 
WHERE status = 'published';

-- Function-based Index für Case-insensitive Suche
CREATE INDEX idx_users_email_lower 
ON users(LOWER(email));
```

### 33.2 Caching-Strategien

| Cache Type | TTL | Invalidation | Use Case |
|-----------|-----|--------------|----------|
| CDN Assets | 1 year | Cache-Busting | Static JS/CSS |
| API Response | 5-15 min | Time-based | GET Endpoints |
| User Session | 24h | Logout/Expire | Auth Tokens |
| Database Query | 1-5 min | Write-through | Frequently accessed |
| AI Response | 1 hour | Hash-based | Duplicate prompts |

### 33.3 Frontend Optimization

- [ ] Code Splitting
- [ ] Lazy Loading
- [ ] Image Optimization
- [ ] Font Subsetting
- [ ] Gzip/Brotli Compression
- [ ] Service Worker
- [ ] Critical CSS Inline
- [ ] Preload Key Resources

---

## 34) Compliance und Datenschutz

### 34.1 DSGVO Anforderungen

| Anforderung | Implementierung |
|-------------|-----------------|
| Right to Access | API Endpoint für Datenexport |
| Right to Rectification | Profile Edit API |
| Right to Erasure | Delete Account mit Cascade |
| Data Portability | JSON/CSV Export |
| Consent Management | Consent Tabelle |
| Data Breach Notification | Incident Response Plan |

### 34.2 Audit Trail

```sql
-- Audit Logs Schema
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Automatische Protokollierung Trigger
CREATE OR REPLACE FUNCTION audit_trigger()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO audit_logs (user_id, action, resource_type, resource_id, old_values, new_values)
    VALUES (
        current_setting('app.user_id', true)::UUID,
        TG_OP,
        TG_TABLE_NAME,
        NEW.id,
        CASE WHEN TG_OP = 'UPDATE' THEN to_jsonb(OLD) ELSE NULL END,
        CASE WHEN TG_OP IN ('INSERT', 'UPDATE') THEN to_jsonb(NEW) ELSE NULL END
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

---

## 35) Referenz-Tabellen

### 35.1 Service Ports

| Service | Internal Port | External Port | Protocol |
|---------|--------------|---------------|----------|
| API Gateway | 8080 | 443 | HTTP/HTTPS |
| Auth Service | 8081 | - | gRPC |
| User Service | 8082 | - | gRPC |
| Content Service | 8083 | - | gRPC |
| AI Orchestrator | 8084 | - | gRPC |
| n8n | 5678 | 443 | HTTP |
| PostgreSQL | 5432 | - | TCP |
| Redis | 6379 | - | TCP |
| Prometheus | 9090 | - | HTTP |
| Grafana | 3000 | 443 | HTTP |

### 35.2 Environment Variablen

| Variable | Required | Description |
|----------|----------|-------------|
| DATABASE_URL | Yes | PostgreSQL Connection String |
| REDIS_URL | Yes | Redis Connection String |
| JWT_SECRET | Yes | JWT Signing Key |
| NVIDIA_API_KEY | Yes | NVIDIA NIM API Key |
| N8N_URL | Yes | n8n Instance URL |
| CLOUDFLARE_TOKEN | No | Cloudflare API Token |
| S3_BUCKET | No | S3 Bucket for Storage |
| SENTRY_DSN | No | Sentry Error Tracking |

---

## 36) Version History

| Version | Datum | Änderungen | Autor |
|---------|-------|------------|-------|
| 1.0 | Feb 2026 | Initiale Version | A1.1 |
| 1.1 | Feb 2026 | Qwen 3.5 Integration | A1.1 |
| 1.2 | Feb 2026 | Erweiterung auf 5000+ Zeilen | A1.1 |

---

**Ende der ARCHITECTURE.md**
