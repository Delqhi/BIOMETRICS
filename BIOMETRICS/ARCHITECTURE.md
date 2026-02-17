# ARCHITECTURE.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Architekturentscheidungen folgen den globalen Regeln aus `AGENTS-GLOBAL.md`.
- Jede Strukturänderung benötigt Mapping- und Integrationsabgleich.
- Security-by-Default, NLM-First und Nachweisbarkeit sind Pflicht.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universelle Architekturvorlage für modulare, skalierbare und auditierbare Systeme.

## Scope
- Frontend: Next.js
- Backend: Go + Supabase
- Integrationen: OpenClaw, n8n, Cloudflare, Vercel (optional)

## 1) Systemkontext
- Nutzergruppen: {PRIMARY_AUDIENCE}, {SECONDARY_AUDIENCE}
- Hauptziele: {BUSINESS_GOAL}
- Externe Systeme: {EXTERNAL_SYSTEMS}

## 2) Architekturprinzipien
1. API-first
2. Modular statt monolithisch
3. Security by default
4. Observability by design
5. Reproduzierbarer Betrieb

## 3) Modulübersicht (Template)

| Modul | Verantwortung | Eingänge | Ausgänge | Abhängigkeiten |
|---|---|---|---|---|
| web-frontend | UI und User Flows | user events | API calls | api-gateway |
| api-gateway | Zugriffsschicht | HTTP requests | responses | services |
| service-core | Business-Logik | API payloads | domain events | supabase |
| content-orchestrator | NLM-Delegation | content tasks | generated assets | nlm-cli |
| workflow-engine | Automationen | triggers | actions | n8n/openclaw |

## 4) Datenfluss (High-Level)
1. User interagiert mit Next.js Frontend.
2. Frontend ruft API-Endpunkte auf.
3. Go-Services validieren, orchestrieren und persistieren.
4. Supabase liefert DB/Auth/Storage.
5. Bei Content-Jobs delegiert der Agent via NLM-CLI.

## 5) Integrationsschnittstellen
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
- Performance: {PERFORMANCE_TARGETS}
- Verfügbarkeit: {AVAILABILITY_TARGETS}
- Sicherheit: {SECURITY_TARGETS}
- Wartbarkeit: kleine, klar getrennte Module

## 8) Risiken
- {ARCH_RISK_1}
- {ARCH_RISK_2}
- {ARCH_RISK_3}

## 9) Entscheidungen (ADR-Index)

| ADR-ID | Thema | Entscheidung | Status |
|---|---|---|---|
| ADR-001 | Frontend-Stack | Next.js | accepted |
| ADR-002 | Backend-Stack | Go + Supabase | accepted |
| ADR-003 | Content-Generierung | NLM-CLI Pflicht | accepted |

## 10) Verifikation
- Konsistenzcheck mit `COMMANDS.md` und `ENDPOINTS.md`
- Abgleich mit `SECURITY.md` und `INFRASTRUCTURE.md`
- Task-20 Integrationsreview

## Abnahme-Check ARCHITECTURE
1. Module und Verantwortungen klar
2. Datenflüsse nachvollziehbar
3. NLM-Integration verankert
4. Risiken dokumentiert
5. ADR-Index gepflegt

---
