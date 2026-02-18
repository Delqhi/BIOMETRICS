# MAPPING-DB-API.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- DB↔API-Mapping ist ein verpflichtender Kontrollpunkt vor Releases.
- Schema-/Contract-Änderungen werden sofort synchronisiert.
- Drift oder Lücken werden als Blocker behandelt.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Abgleich zwischen Datenmodell (Supabase) und API-Verträgen.

## Mapping-Tabelle

| Tabelle | Endpoint(s) | Operationen | RLS-Bezug | Status |
|---|---|---|---|---|
| users | GET /api/users, POST /api/users, GET /api/users/:id, PATCH /api/users/:id | CRU | user/admin policies | active |
| profiles | GET /api/profiles, PATCH /api/profiles | RU | user policies | active |
| projects | GET /api/projects, POST /api/projects, GET /api/projects/:id, PATCH /api/projects/:id, DELETE /api/projects/:id | CRUD | project policies | active |
| assets | GET /api/assets, POST /api/assets, GET /api/assets/:id, DELETE /api/assets/:id | CRD | asset policies | active |
| tasks | GET /api/tasks, POST /api/tasks, PATCH /api/tasks/:id, DELETE /api/tasks/:id | CRUD | task policies | active |
| products | GET /api/products, POST /api/products, PATCH /api/products/:id, DELETE /api/products/:id | CRUD | public read, admin write | active |
| orders | GET /api/orders, POST /api/orders, GET /api/orders/:id | CR (own only) | user order policies | active |
| sessions | GET /api/sessions, DELETE /api/sessions/:id | RL (own only) | session policies | active |

## Prüfregeln
1. Keine API ohne dokumentierte Datenbasis
2. Jede schreibende API referenziert Validierungs- und RLS-Regeln
3. Löschpfade sind dokumentiert und nachvollziehbar

## Abnahme-Check
1. Tabellen-zu-API Zuordnung vollständig
2. Operationen pro Endpoint dokumentiert
3. RLS-Bezug je Zuordnung enthalten

---
