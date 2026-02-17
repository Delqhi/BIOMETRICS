# TROUBLESHOOTING.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Troubleshooting wird ticketbasiert, reproduzierbar und evidenzgestützt geführt.
- Root Cause, Corrective und Preventive Actions sind Pflichtbestandteile.
- Learnings fließen in Regeln, Playbooks und Training zurück.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universelle Fehlerdiagnose- und Behebungsleitfäden für Betrieb und Entwicklung.

## Diagnoseprinzip
1. Reproduzieren
2. Isolieren
3. Root-Cause finden
4. Fix minimal-invasiv umsetzen
5. Regression verhindern
6. Doku aktualisieren

## Fehlerkatalog (Template)

### Fall 01: Build schlägt fehl
Symptom:
- Build-Prozess endet mit Fehler.

Wahrscheinliche Ursachen:
- inkonsistente Abhängigkeiten
- Typfehler
- ungültige Konfiguration

Diagnose:
- `pnpm build`
- `pnpm typecheck`

Lösung:
1. Fehlerursache priorisieren
2. fixen
3. erneut builden

Verifikation:
- Build erfolgreich

### Fall 02: Tests schlagen fehl
Symptom:
- Test-Suite nicht grün.

Wahrscheinliche Ursachen:
- Regression im Code
- geänderte Schnittstellen

Diagnose:
- `pnpm test`
- `go test ./...`

Lösung:
1. failing tests isolieren
2. root cause beheben
3. relevante Tests erneut ausführen

Verifikation:
- betroffene Tests grün

### Fall 03: API-Fehler 401/403
Symptom:
- Zugriff verweigert.

Wahrscheinliche Ursachen:
- falsche Rolle
- fehlende Auth-Konfiguration

Diagnose:
- Auth-Kontext prüfen
- Rollenmapping prüfen

Lösung:
1. Rollen-/Policy-Konfiguration korrigieren
2. Auth-Flow verifizieren

Verifikation:
- Endpoint mit korrekter Rolle erreichbar

### Fall 04: NLM-Output unbrauchbar
Symptom:
- inkonsistent, unklar oder faktisch falsch.

Wahrscheinliche Ursachen:
- unklarer Prompt
- fehlende Quellen
- fehlende Qualitätsprüfung

Diagnose:
- Prompt gegen Vorlagen prüfen
- Quellabdeckung prüfen

Lösung:
1. Prompt präzisieren
2. zweite Iteration erzeugen
3. Qualitätsmatrix anwenden

Verifikation:
- Score >= 13/16 und Korrektheit 2/2

### Fall 05: Dokumente widersprechen sich
Symptom:
- `COMMANDS.md` und `ENDPOINTS.md` inkonsistent.

Wahrscheinliche Ursachen:
- asynchrone Updates
- fehlende Cross-Doc Prüfung

Diagnose:
- Mapping-Check durchführen

Lösung:
1. Primärquelle festlegen
2. Gegenstück synchronisieren
3. in Changelog dokumentieren

Verifikation:
- 1:1 Mapping konsistent

## Eskalations-Playbook
- P0: sofortige Eskalation und Schadensbegrenzung
- P1: innerhalb Session beheben/eskalieren
- P2: nächsten Zyklus planen

## Postmortem-Vorlage
```text
Incident-ID:
Severity:
Timeline:
Root Cause:
Impact:
Fix:
Prävention:
Owner:
Follow-up Tasks:
```

## Abnahme-Check TROUBLESHOOTING
1. Root-Cause Ansatz beschrieben
2. NLM-Fehlerfall enthalten
3. API/Auth-Fälle enthalten
4. Verifikationsschritte je Fall vorhanden
5. Eskalationspfad vorhanden

---
