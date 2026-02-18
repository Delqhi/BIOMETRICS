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

## Platzhalter
- {PROJECT_NAME}
- {PRODUCT_TYPE}
- {PRIMARY_AUDIENCE}
- {SECONDARY_AUDIENCE}
- {BUSINESS_GOAL}
- {SUCCESS_METRICS}
- {CONSTRAINTS}
- {NON_GOALS}

## 1) Problemraum
- Nutzerproblem: {USER_PROBLEM}
- Ist-Zustand: {CURRENT_STATE}
- Soll-Zustand: {TARGET_STATE}

## 2) Zielbild
- Produktvision: {PRODUCT_VISION}
- Kernnutzen: {CORE_VALUE}
- Differenzierung: {DIFFERENTIATION}

## 3) Business-Ziele
1. {GOAL_1}
2. {GOAL_2}
3. {GOAL_3}

## 4) Zielgruppen
### Primär
- Segment: {PRIMARY_AUDIENCE}
- Jobs-to-be-done: {PRIMARY_JOBS}
- Risiken/Einwände: {PRIMARY_RISKS}

### Sekundär
- Segment: {SECONDARY_AUDIENCE}
- Jobs-to-be-done: {SECONDARY_JOBS}
- Risiken/Einwände: {SECONDARY_RISKS}

## 5) Scope
### In Scope
- {IN_SCOPE_1}
- {IN_SCOPE_2}
- {IN_SCOPE_3}

### Out of Scope
- {OUT_SCOPE_1}
- {OUT_SCOPE_2}

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
- {RISK_1}
- {RISK_2}
- {RISK_3}

## 10) Annahmen
- {ASSUMPTION_1}
- {ASSUMPTION_2}

## 11) Offene Entscheidungen
- {DECISION_1}
- {DECISION_2}

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
