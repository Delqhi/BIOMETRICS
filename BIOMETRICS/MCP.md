# MCP.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- MCP-Betrieb folgt globalen Sicherheits- und Wrapper-Regeln.
- Konfigurationen bleiben nachvollziehbar, testbar und ohne Blindänderungen.
- Troubleshooting ist ticket- und evidenzbasiert zu führen.

Status: ACTIVE  
Version: 1.1 (Universal)  
Stand: Februar 2026

## Zweck
Globales MCP-Betriebshandbuch im Projektkontext: welche MCP-Server existieren, wie sie installiert, aktiviert und sicher genutzt werden.

## Pfadstandard
`BIOMETRICS/MCP.md` ist der kanonische Speicherort für das MCP-Book dieses Projekts.

## Grundregeln
1. MCP-Nutzung ist dokumentationspflichtig.
2. Vor Einsatz immer Zweck, Risiko und erwartetes Ergebnis klären.
3. Serverzugriffe mit Least-Privilege konfigurieren.
4. Änderungen und Nutzung in `BIOMETRICS/MEETING.md` protokollieren.

## Server-Register (Template)

| Server | Zweck | Installation | Aktivierung | Primäre Use Cases | Sicherheitsniveau | Status |
|---|---|---|---|---|---|---|
| Serena MCP | Codekontext, Projektaktivierung | {SERENA_INSTALL} | {SERENA_ACTIVATE} | Refactoring, Strukturarbeit | hoch | active |
| CDP MCP | Browser/Live-UI Analyse | {CDP_INSTALL} | {CDP_ACTIVATE} | UI-Checks, Live-Testflows | hoch | active |
| Tavily MCP | Recherche/Knowledge | {TAVILY_INSTALL} | {TAVILY_ACTIVATE} | Wissensrecherche, Vergleich | mittel | active |
| Custom MCP A | {CUSTOM_A_PURPOSE} | {CUSTOM_A_INSTALL} | {CUSTOM_A_ACTIVATE} | {CUSTOM_A_USE} | {CUSTOM_A_SEC} | planned |

## Installations-Standard
Für jeden Server dokumentieren:
1. Voraussetzungen
2. Installationsschritte
3. Konfigurationsparameter
4. Verbindungsprüfung
5. Rollback bei Fehlinstallation

## Aktivierungs-Standard
Für jeden Server dokumentieren:
1. Aktivierungsbefehl
2. Auth- und Rechtekontext
3. Health-Check
4. Deaktivierungsbefehl

## Nutzungsmatrix

| Aufgabe | Primärer MCP | Optionaler MCP | Warum |
|---|---|---|---|
| Code-Änderungen im Repo | Serena MCP | CDP MCP | Struktur + Liveprüfung |
| Browser-/UI-Probleme | CDP MCP | Serena MCP | reproduzierbarer UI-Test |
| Externe Recherche | Tavily MCP | Serena MCP | Quellenorientierung |
| NLM Asset Governance | Serena MCP | CDP MCP | Doku + Integrationscheck |

## Sicherheitsregeln
1. Keine Secrets in Prompt- oder Doku-Texten.
2. Keine unkontrollierten Schreibzugriffe.
3. Jede MCP-gestützte Änderung braucht Evidenz.
4. Bei P0-Risiko sofort eskalieren.

## Troubleshooting (Template)

| Problem | Ursache (vermutet) | Diagnose | Maßnahme | Eskalation |
|---|---|---|---|---|
| Server nicht erreichbar | Fehlkonfiguration | Health-Check | Konfiguration prüfen | P1 |
| Auth schlägt fehl | Token/Rechte | Auth-Log prüfen | Rechte korrigieren | P0 |
| Unklare Ergebnisse | falscher Scope | Scope-Review | Auftrag schärfen | P1 |

## Änderungslog

| Datum | Änderung | Betroffene Server | Owner |
|---|---|---|---|
| {DATE} | Initialer Universal-Stand | Serena/CDP/Tavily | {OWNER_ROLE} |

## Abnahme-Check MCP
1. Server-Register vorhanden
2. Installation/Aktivierung dokumentiert
3. Nutzungsmatrix vorhanden
4. Sicherheitsregeln enthalten
5. Troubleshooting vorhanden

---
