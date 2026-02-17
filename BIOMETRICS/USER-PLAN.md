# USER-PLAN.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Plan Sovereignty gilt verbindlich für Priorisierung und Abarbeitung.
- Todos benötigen Ownership, Status und Evidence-Erwartung.
- Planänderungen werden nachvollziehbar und konfliktfrei dokumentiert.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Diese Datei enthält ausschließlich Aufgaben, die der User selbst erledigen muss und die nicht sinnvoll automatisiert werden können.

## Universalitäts-Regeln
1. Keine projektspezifischen Secrets speichern.
2. Keine sensiblen Zugangsdaten in Klartext.
3. Alle Einträge bleiben projektagnostisch mit Platzhaltern.

## Platzhalter
- {PROJECT_NAME}
- {ENV_NAME}
- {DOMAIN}
- {PROVIDER}
- {ACCOUNT_OWNER}
- {DATE_DUE}

## Statusmodell
- OPEN
- IN_PROGRESS
- BLOCKED
- DONE

## Prioritätsmodell
- P0: Blockiert Go-Live oder Sicherheit
- P1: Wichtig für Stabilität und Betrieb
- P2: Optimierung

## Eintragsvorlage
```text
User-Task-ID:
Titel:
Priorität:
Status:
Warum User-Aufgabe:
Voraussetzungen:
Schritte:
Erwartetes Ergebnis:
Verifikation:
Risiko bei Nicht-Erledigung:
Fälligkeitsdatum:
Owner:
```

## Universal User-Task Backlog (Startset)

### UT-001
Titel: Accounts für notwendige Plattformen bereitstellen  
Priorität: P0  
Status: OPEN

Warum User-Aufgabe:
Rechtliche Inhaberschaft und Abrechnung liegen beim User.

Voraussetzungen:
- Entscheidung über Zielplattformen.

Schritte:
1. Plattform-Accounts unter Unternehmensidentität erstellen.
2. Rollen und Zugriffe dokumentieren.
3. Wiederherstellungswege aktivieren.

Erwartetes Ergebnis:
Zugänge existieren, Recovery ist aktiv.

Verifikation:
- Login-Test je Plattform.
- Recovery-E-Mail/2FA bestätigt.

Risiko bei Nicht-Erledigung:
Deployment und Betrieb blockiert.

Fälligkeitsdatum:
{DATE_DUE}

Owner:
{ACCOUNT_OWNER}

### UT-002
Titel: DNS-/Domain-Verantwortung bestätigen  
Priorität: P0  
Status: OPEN

Warum User-Aufgabe:
Domain-Verwaltung ist organisatorisch und vertraglich userseitig.

Voraussetzungen:
- Verfügbarkeit der gewünschten Domain.

Schritte:
1. Eigentum/Transfer der Domain sicherstellen.
2. Zugriff für technische Umsetzung bereitstellen.
3. Verantwortliche Kontaktperson festlegen.

Erwartetes Ergebnis:
Domain administrativ verfügbar und delegierbar.

Verifikation:
- DNS-Panel Zugriff bestätigt.

Risiko bei Nicht-Erledigung:
Keine saubere produktive Erreichbarkeit.

Fälligkeitsdatum:
{DATE_DUE}

Owner:
{ACCOUNT_OWNER}

### UT-003
Titel: Rechtliche Texte bereitstellen  
Priorität: P0  
Status: OPEN

Warum User-Aufgabe:
Rechtliche Inhalte müssen vom verantwortlichen Betreiber freigegeben werden.

Voraussetzungen:
- Juristische Ansprechpartner vorhanden.

Schritte:
1. Impressum/Datenschutz/Nutzungsbedingungen bereitstellen.
2. Gültigkeit für Zielmärkte bestätigen.
3. Freigabestand dokumentieren.

Erwartetes Ergebnis:
Freigegebene rechtliche Inhalte liegen vor.

Verifikation:
- Freigabevermerk mit Datum.

Risiko bei Nicht-Erledigung:
Compliance-Risiko und Betriebsverzögerung.

Fälligkeitsdatum:
{DATE_DUE}

Owner:
{ACCOUNT_OWNER}

### UT-004
Titel: API-/Provider-Verträge und Limits klären  
Priorität: P1  
Status: OPEN

Warum User-Aufgabe:
Kosten- und Vertragsentscheidungen liegen beim User.

Voraussetzungen:
- Auswahl der Provider.

Schritte:
1. Tarife und Limits prüfen.
2. Budgetgrenzen definieren.
3. Vertragsdokumente zentral hinterlegen.

Erwartetes Ergebnis:
Planbare Kosten- und Kapazitätsbasis.

Verifikation:
- Vertragsstatus dokumentiert.

Risiko bei Nicht-Erledigung:
Kostenüberraschungen und Ausfälle durch Limits.

Fälligkeitsdatum:
{DATE_DUE}

Owner:
{ACCOUNT_OWNER}

### UT-005
Titel: Produktionsfreigabe (Go/No-Go) erteilen  
Priorität: P0  
Status: OPEN

Warum User-Aufgabe:
Business-Risikoentscheidung ist Managementverantwortung.

Voraussetzungen:
- Abschlussreport aus Task 20 liegt vor.

Schritte:
1. Risiken und Restaufwand prüfen.
2. Freigabeentscheidung dokumentieren.
3. Rolloutfenster bestätigen.

Erwartetes Ergebnis:
Verbindliche Go/No-Go Entscheidung.

Verifikation:
- Schriftliche Freigabe mit Datum.

Risiko bei Nicht-Erledigung:
Unklare Verantwortlichkeit und verzögerter Launch.

Fälligkeitsdatum:
{DATE_DUE}

Owner:
{ACCOUNT_OWNER}

## Abnahme-Check USER-PLAN
1. Nur userseitige Aufgaben enthalten.
2. Jede Aufgabe ist verifizierbar.
3. Jede Aufgabe hat Risiko und Fälligkeit.
4. Keine technischen Implementierungsaufgaben enthalten.

---
