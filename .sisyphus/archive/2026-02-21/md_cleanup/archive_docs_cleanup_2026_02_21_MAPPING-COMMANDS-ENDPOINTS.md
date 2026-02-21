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
| CMD.QWEN.VISION | POST /api/qwen/vision | ja | ja | aligned |
| CMD.QWEN.CODE | POST /api/qwen/chat | ja | ja | aligned |
| CMD.QWEN.OCR | POST /api/qwen/ocr | ja | ja | aligned |
| CMD.QWEN.VIDEO | POST /api/qwen/video | ja | ja | aligned |
| CMD.QWEN.CONVERSATION | POST /api/qwen/chat | ja | ja | aligned |
| CMD.AUTH.LOGIN | POST /api/auth/login | ja | ja | aligned |
| CMD.AUTH.LOGOUT | POST /api/auth/logout | ja | ja | aligned |
| CMD.PRODUCTS.LIST | GET /api/products | ja | ja | aligned |
| CMD.PRODUCTS.CREATE | POST /api/products | ja | ja | aligned |
| CMD.ORDERS.CREATE | POST /api/orders | ja | ja | aligned |
| CMD.ORDERS.LIST | GET /api/orders | ja | ja | aligned |

## Offene Deltas
- {none - alle Qwen 3.5 Skills gemappt}

## Abnahme-Check
1. Jeder produktive Command hat Endpoint
2. Rollen sind deckungsgleich
3. Fehlersemantik ist konsistent

---
