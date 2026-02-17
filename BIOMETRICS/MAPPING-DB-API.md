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
| users | {API_USERS} | CRU | user/admin policies | planned |
| projects | {API_PROJECTS} | CRUD | project policies | planned |
| assets | {API_ASSETS} | CRUD | asset policies | planned |
| tasks | {API_TASKS} | CRUD | task policies | planned |

## Prüfregeln
1. Keine API ohne dokumentierte Datenbasis
2. Jede schreibende API referenziert Validierungs- und RLS-Regeln
3. Löschpfade sind dokumentiert und nachvollziehbar

## Abnahme-Check
1. Tabellen-zu-API Zuordnung vollständig
2. Operationen pro Endpoint dokumentiert
3. RLS-Bezug je Zuordnung enthalten

---
