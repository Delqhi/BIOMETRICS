# MAPPING-COMMANDS-ENDPOINTS.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Mapping ist verbindliche Kontrollfläche gemäß `AGENTS-GLOBAL.md`.
- Änderungen an Commands/Endpoints erfordern sofortige Synchronisierung.
- Ungeklärte Deltas blockieren Release-Readiness.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Abgleich zwischen Commands und API-Endpunkten.

## Mapping-Tabelle

| Command | Endpoint | Rollen konsistent | Fehlercodes konsistent | Status |
|---|---|---|---|---|
| CMD.SYSTEM.STATUS | API.SYSTEM.STATUS | ja | ja | aligned |
| CMD.TASKS.LIST | API.TASKS.LIST | ja | ja | aligned |
| CMD.TASKS.UPDATE | API.TASKS.UPDATE | ja | ja | aligned |
| CMD.NLM.GENERATE.VIDEO | API.NLM.GENERATE.VIDEO | ja | ja | aligned |
| CMD.NLM.GENERATE.INFOGRAPHIC | API.NLM.GENERATE.INFOGRAPHIC | ja | ja | aligned |
| CMD.NLM.GENERATE.PRESENTATION | API.NLM.GENERATE.PRESENTATION | ja | ja | aligned |
| CMD.NLM.GENERATE.TABLE | API.NLM.GENERATE.TABLE | ja | ja | aligned |

## Offene Deltas
- {DELTA_1}
- {DELTA_2}

## Abnahme-Check
1. Jeder produktive Command hat Endpoint
2. Rollen sind deckungsgleich
3. Fehlersemantik ist konsistent

---
