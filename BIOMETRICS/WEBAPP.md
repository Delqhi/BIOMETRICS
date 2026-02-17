# WEBAPP.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Webapp-Architektur folgt Stack-Lock (Next.js/TS strict) und Governance-Standards.
- Performance-, Security- und Observability-Anforderungen sind verpflichtend.
- API-/DB-Abhängigkeiten bleiben über Mappings konsistent.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Rahmen für Webapp-Produktlogik, Rollen, Kernflows und Qualitätsstandards.

## Grundprinzip
Jede Funktion muss einem klaren Nutzerziel dienen, über Command + Endpoint steuerbar sein und dokumentiert getestet werden.

## Rollenmodell (Template)
- user
- dev
- admin
- agent

## Kernflows (Template)

| Flow-ID | Ziel | Startpunkt | Endzustand | Kritikalität |
|---|---|---|---|---|
| FLOW-001 | {GOAL_1} | {ENTRY_1} | {OUTCOME_1} | P0 |
| FLOW-002 | {GOAL_2} | {ENTRY_2} | {OUTCOME_2} | P1 |

## Funktionsmatrix (Template)

| Funktion | Nutzerwert | Command | Endpoint | Status |
|---|---|---|---|---|
| {FEATURE_1} | {VALUE_1} | {CMD_1} | {API_1} | planned |
| {FEATURE_2} | {VALUE_2} | {CMD_2} | {API_2} | planned |

## UX-Prinzipien
1. konsistente Begriffe
2. klare Handlungsaufforderungen
3. verständliche Fehlerzustände
4. mobile-first Baseline

## NLM-Asset-Einsatz
- komplexe Flows können mit Video/Infografik erklärt werden
- Erzeugung ausschließlich via NLM-CLI
- Freigabe nur nach Qualitätsmatrix

## Qualitäts-Gates
- Flow ist reproduzierbar
- Command/Endpoint Mapping vorhanden
- kritische Pfade testbar
- Doku synchron

## Flow-Kompatibilität (Check)
1. Jeder Kernflow hat mindestens einen Command
2. Jeder genutzte Command hat Endpoint-Gegenstück
3. Rollenrechte in Flow und Endpoint sind konsistent
4. Fehlerpfade sind pro Kernflow beschrieben
5. NLM-Erklärassets unterstützen nur echte Funktionen

## Verifikation
- Kernflow-Walkthrough
- Command-zu-Endpoint Konsistenzcheck
- NLM-Asset-Check bei erklärungsintensiven Flows

## Querverlinkung
- `WEBSITE.md`: öffentlicher Informations- und Content-Kontext
- `COMMANDS.md`: steuerbare Funktionsbefehle
- `ENDPOINTS.md`: technische API-Abbildung
- `BIOMETRICS/MAPPING-COMMANDS-ENDPOINTS.md`: Konsistenzreport

## Abnahme-Check WEBAPP
1. Rollenmodell vorhanden
2. Kernflows definiert
3. Feature-Matrix mit Command/Endpoint vorhanden
4. Qualitäts- und Verifikationslogik dokumentiert

---
