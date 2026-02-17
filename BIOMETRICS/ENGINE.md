# ENGINE.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Engine-Regeln folgen globalen Betriebs-, Sicherheits- und Qualitätsvorgaben.
- Ausführungslogik muss testbar, observierbar und dokumentiert bleiben.
- Kritische Änderungen brauchen Incident- und Recovery-Anschlussfähigkeit.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Rahmen für eine projektspezifische Core-Engine oder zentrale Verarbeitungslogik.

## Hinweis
Falls keine Engine-Komponente existiert, Status auf `NOT_APPLICABLE` setzen und begründen.

## Engine-Verantwortung (Template)
- Inputverarbeitung
- Regelanwendung
- Ergebnisorchestrierung
- Fehlerhärtung

## Architektur-Schnittstellen

| Schnittstelle | Eingabe | Ausgabe | Fehlerverhalten |
|---|---|---|---|
| API -> Engine | payload | response | validation errors |
| Engine -> Storage | normalized data | write status | retry policy |
| Engine -> Integrationslayer | job context | integration result | escalation |

## Laufzeitgrenzen
- erwartete Lastprofile: {LOAD_PROFILE}
- Timeout-Strategie: {TIMEOUT_POLICY}
- Retry-Strategie: {RETRY_POLICY}

## Qualitätsziele
1. deterministisches Verhalten in Kernpfaden
2. klar dokumentierte Fehlerzustände
3. reproduzierbare Ergebnisse

## Observability
- zentrale Metriken
- Fehlerklassifikation
- Durchsatz-/Latenzwerte

## Verifikation
- Core-Path Tests
- Fehlerpfadtests
- Lastprofil-Sanity-Checks

## Abnahme-Check ENGINE
1. Verantwortungen klar
2. Schnittstellen dokumentiert
3. Laufzeit- und Qualitätsziele vorhanden
4. Verifikation definiert

---
