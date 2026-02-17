# MAPPING-FRONTEND-BACKEND.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Frontend↔Backend-Verträge sind als verbindliche Kontrollpunkte zu pflegen.
- Jede Änderung erfordert zeitnahes Mapping-Update mit Ownership.
- Inkonsistenzen werden als Delivery-Risiko behandelt.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Abgleich zwischen Frontend-Flows und Backend-Funktionen.

## Mapping-Tabelle

| Frontend-Flow | Command | Endpoint | Backend-Service | Status |
|---|---|---|---|---|
| {FLOW_1} | {CMD_1} | {API_1} | {SERVICE_1} | planned |
| {FLOW_2} | {CMD_2} | {API_2} | {SERVICE_2} | planned |

## Prüfregeln
1. Kein Flow ohne API-Anbindung
2. Kein Endpoint ohne Frontend-Nutzen in Kernpfaden
3. Fehlerfälle müssen UI-seitig behandelbar sein

## Abnahme-Check
1. Kernflows vollständig gemappt
2. Service-Zuordnung vorhanden
3. Fehlerpfad pro Kernflow dokumentiert

---
