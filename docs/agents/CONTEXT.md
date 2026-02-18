# CONTEXT.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Kontextdokumentation wird als operative Steuerungsgrundlage gepflegt.
- Wissensintegrität ist append-only; Migration statt Informationsverlust.
- Regel- und Entscheidungsbezüge bleiben stets nachvollziehbar.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Latest Changes (Februar 2026)

### Qwen 3.5 397B Integration
- **Modell:** qwen/qwen3.5-397b-a17b
- **Context:** 262K tokens
- **Output:** 32K tokens
- **Provider:** NVIDIA NIM
- **Status:** ✅ Aktiv
- **Use Case:** Code-Generation (Best-in-Class)

### NVIDIA NIM Konfiguration
- **Endpoint:** https://integrate.api.nvidia.com/v1
- **API:** openai-completions
- **Timeout:** 120000ms (120s) - erforderlich wegen hoher Latenz
- **Rate Limit:** 40 RPM (Free Tier)
- **HTTP 429 Lösung:** 60 Sekunden warten + Fallback nutzen

### Verfügbare NVIDIA NIM Modelle
| Modell | Context | Output | Use Case |
|--------|---------|--------|----------|
| qwen3.5-397b | 262K | 32K | Code (BEST) |
| Qwen2.5-Coder-32B | 128K | 8K | Code (fast) |
| Qwen2.5-Coder-7B | 128K | 8K | Code (fastest) |
| Kimi K2.5 | 1M | 64K | General |

### OpenCode.json Konfiguration (NVIDIA NIM)
```json
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
```

## Zweck
Universeller Kontext-Container für Ziele, Zielgruppen, Scope und Randbedingungen.
Dieses Dokument dient als Orientierung für neue Teammitglieder und definiert die Projektgrundlagen.

## Platzhalter
- {PROJECT_NAME} = BIOMETRICS
- {PRODUCT_TYPE} = KI-Orchestrierungsplattform
- {PRIMARY_AUDIENCE} = Entwickler, die KI-Workflows automatisieren möchten
- {SECONDARY_AUDIENCE} = Unternehmen, die Self-Building AI Agents einsetzen wollen
- {BUSINESS_GOAL} = Vollautomatisierte KI-Entwicklung ermöglichen
- {SUCCESS_METRICS} = Deployment-Zeit < 5min, 100% automatisierte Tests, 24/7 Verfügbarkeit
- {CONSTRAINTS} = Keine externen Kosten (Self-Hosted), DSGVO-konform
- {NON_GOALS} = Keine Enterprise-Lizenzierung, kein proprietäres SaaS

## 1) Problemraum
- **Nutzerproblem:** Entwicklung von KI-Workflows ist manuell, zeitintensiv und fehleranfällig
- **Ist-Zustand:** Separate Tools für n8n, Supabase, OpenClaw ohne zentrale Orchestrierung
- **Soll-Zustand:** Einheitliche CLI mit Self-Building Agent für vollautomatisierte Entwicklung

## 2) Zielbild
- **Produktvision:** Dein KI-Entwicklungspartner - von der Idee zur Produktion in Minuten
- **Kernnutzen:** Automatisierung aller repetitive Entwicklungsaufgaben durch Self-Building Agents
- **Differenzierung:** Vollständig Self-Hosted, DSGVO-konform, keine externen Kosten

## 3) Business-Ziele
1. **Automatisierung:** 90% der Entwicklungsaufgaben durch KI-Agenten automatisieren
2. **Geschwindigkeit:** Deployment-Zeit von Stunden auf Minuten reduzieren
3. **Qualität:** 100% Testabdeckung bei allen neuen Features

## 4) Zielgruppen
### Primär
- **Segment:** Full-Stack Entwickler und DevOps Engineers
- **Jobs-to-be-done:** Automatisierte CI/CD-Pipelines, Self-Healing Infrastruktur, KI-Workflows ohne Handarbeit
- **Risiken/Einwände:** Angst vor Verlust der Kontrolle, Lernkurve, Integration mit bestehenden Systemen

### Sekundär
- **Segment:** CTOs und Tech-Leads, die Entwicklungsteams effizienter machen wollen
- **Jobs-to-be-done:** Kosteneffiziente Automatisierung, Standardisierung von Entwicklungsprozessen
- **Risiken/Einwände:** Investitionsbedenken, Mangel an internem Know-how

## 5) Scope
### In Scope
- **CLI-Tool:** biometrics CLI für Setup und Orchestrierung
- **Agent-Framework:** Self-Building Agents für Code- und Workflow-Generierung
- **Integrationen:** Supabase, n8n, OpenClaw, Cloudflare, Vercel
- **Dokumentation:** Vollständige NLM-Doku für alle Komponenten

### Out of Scope
- Proprietäre SaaS-Produkte
- Externe API-Integrationen mit Kostenpflicht
- Enterprise-spezifische Features (LDAP, SAML)

## 6) Qualitätsmaßstäbe
- Best Practices Februar 2026
- Production-ready statt Demo
- Messbare Qualität statt Behauptungen

## 7) Technikrahmen
- Frontend: Next.js
- Backend: Go + Supabase
- Package Manager: pnpm
- Integrationen: OpenClaw, n8n, Cloudflare, Vercel (falls genutzt)

## 8) Content- und NLM-Kontext
- NLM-CLI Pflicht: Ja
- Asset-Typen: Video, Infografik, Präsentation, Datentabelle
- Qualitätsmatrix aktiv: Ja (13/16, Korrektheit 2/2)

## 9) Risiken
- **Risiko 1:** NVIDIA NIM Rate Limits (40 RPM) → Fallback auf alternative Modelle
- **Risiko 2:** Komplexität der Integrationen → Modularisierung und gute Doku
- **Risiko 3:** Wartbarkeit des Self-Building Agents → Klare Architekturprinzipien

## 10) Annahmen
- **Annahme 1:** Das Team hat grundlegende DevOps-Kenntnisse
- **Annahme 2:** Docker und Node.js sind vorhanden
- **Annahme 3:** Die NVIDIA NIM API bleibt kostenlos nutzbar

## 11) Offene Entscheidungen
- **Entscheidung 1:** Welche Vercel-Features werden genutzt? (Edge Functions vs. Serverless)
- **Entscheidung 2:** Backup-Strategie für Supabase-Daten

## 12) Verifikationslogik
- Akzeptanzkriterien pro Task
- Tests pro Task
- Doku-Sync Pflicht
- Evidenzpflicht für Done

## Abnahme-Check CONTEXT
1. Projektagnostische Platzhalter vorhanden
2. Scope und Non-Goals klar
3. Risiken und Annahmen dokumentiert
4. NLM-CLI Kontext enthalten

---

## Für neue Teammitglieder

Willkommen bei BIOMETRICS! Hier ist dein Startpfad:

### Tag 1: Orientation
1. README.md lesen
2. AGENTS.md und AGENTS-PLAN.md lesen
3. CONTEXT.md (dieses Dokument) vollständig lesen
4. ONBOARDING.md durcharbeiten

### Tag 2: Setup
1. `biometrics-cli` installieren: `cd biometrics-cli && pnpm install && pnpm link --global`
2. Onboarding starten: `biometrics`
3. Ersten Task aus USER-PLAN.md oder AGENTS-PLAN.md übernehmen

### Tag 3-7: Einarbeitung
- In AGENTS-PLAN.md nach verfügbaren Tasks suchen
- Mit kleineren P2-Tasks beginnen
- Fortschritt in MEETING.md dokumentieren
- Bei Fragen: CHANGELOG.md und TROUBLESHOOTING.md prüfen

### Wichtige Commands
```bash
biometrics           # CLI starten
biometrics-onboard   # Vollständiges Onboarding
pnpm dev           # Entwicklungsserver starten
pnpm test          # Tests ausführen
pnpm lint          # Code-Qualität prüfen
```

### Wo finde ich was?
- **Architektur:** ARCHITECTURE.md
- **Security:** SECURITY.md  
- **API-Endpunkte:** ENDPOINTS.md
- **n8n Workflows:** N8N.md
- **Supabase:** SUPABASE.md
- **Troubleshooting:** TROUBLESHOOTING.md

---
