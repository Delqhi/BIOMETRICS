# AGENTS-PLAN.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Agentenplanung muss globalen Regelkern und lokale Ausführung koppeln.
- Delegationsblöcke sind standardisiert, verifizierbar und revisionssicher.
- Keine Maßnahme ohne klaren Kontroll- und Übergabepfad.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Universalitäts-Regel
Dieser Plan ist absichtlich projektagnostisch. Er enthält nur wiederverwendbare Strukturen, Platzhalter und Qualitätsregeln.

Pflicht:
1. Keine projektspezifischen Secrets in dieser Datei.
2. Keine hardcodierten Domains, IDs oder Kundendaten.
3. Jede konkrete Projektinstanz ersetzt Platzhalter sauber und vollständig.

## Platzhalter-Konvention
- {PROJECT_NAME}
- {PRIMARY_AUDIENCE}
- {BUSINESS_GOAL}
- {CORE_FEATURE_SET}
- {CONTENT_DOMAIN}
- {TARGET_LANGUAGES}
- {CHANNELS}
- {COMPLIANCE_SCOPE}
- {OWNER_ROLE}

## Betriebsregeln
1. Erst lesen, dann bearbeiten.
2. NLM immer vollumfänglich via NLM-CLI nutzen.
3. Für Website-/Webapp-Erklärpfade Video-Einsatz prüfen und dokumentieren.
4. NLM-Outputs nur nach Qualitätsprüfung übernehmen.
5. Jede Änderung in `MEETING.md` und `CHANGELOG.md` protokollieren.
6. Keine Done-Behauptung ohne Evidenz.
7. Jede Aufgabe hat Akzeptanzkriterien, Tests und Doku-Update.

## Qualitätskriterien (global)
- Korrektheit
- Konsistenz
- Zielgruppenfit
- Umsetzbarkeit
- Wiederverwendbarkeit
- Evidenzbezug

Freigabe:
- Mindestscore: 13/16 (NLM-Matrix)
- Korrektheit muss 2/2 sein

## Zyklus
- Zyklus-ID: LOOP-001
- Umfang: 20 Tasks
- Modus: Universal NLM-Ready
- Abschluss: Task 20 = All-in-One Verification + neue 20 Tasks

## Task-Board (20 Tasks)

### Task 01
Task-ID: LOOP-001-T01  
Titel: Universalen Kontext-Rahmen definieren  
Kategorie: Architektur  
Priorität: P0

Ziel:
Eine robuste, projektagnostische Kontextstruktur vorbereiten.

Read First:
- `CONTEXT.md` (falls vorhanden)
- `ARCHITECTURE.md` (falls vorhanden)

Edit:
- `CONTEXT.md`

Akzeptanzkriterien:
1. Platzhaltermodell vollständig.
2. Zielgruppe und Business-Ziel als Templates dokumentiert.
3. Keine projektspezifischen Details fest verdrahtet.

Tests:
- Konsistenzcheck mit Platzhaltern.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- Abschnitt „Universal Context Template“ vorhanden.

### Task 02
Task-ID: LOOP-001-T02  
Titel: NLM-CLI Betriebsstandard festschreiben  
Kategorie: Enablement  
Priorität: P0

Ziel:
Verbindliche NLM-CLI Nutzung für alle Agenten sicherstellen.

Read First:
- `../∞Best∞Practices∞Loop.md`

Edit:
- `AGENTS.md` (falls vorhanden)
- `COMMANDS.md` (falls vorhanden)

Akzeptanzkriterien:
1. NLM-CLI Pflicht klar dokumentiert.
2. Delegationsregeln enthalten.
3. Fallback-Regel bei NLM-Fehlern dokumentiert.

Tests:
- Regelset-Vollständigkeit gegen Checkliste prüfen.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- Abschnitt „NLM CLI Pflichtbetrieb“ ergänzt.

### Task 03
Task-ID: LOOP-001-T03  
Titel: Video-Policy für Websites definieren  
Kategorie: Feature  
Priorität: P0

Ziel:
Universelle Policy für Videoeinsatz auf Websites/Webapps festlegen.

Read First:
- `WEBSITE.md` (falls vorhanden)
- `WEBAPP.md` (falls vorhanden)

Edit:
- `WEBSITE.md`
- `WEBAPP.md`

Akzeptanzkriterien:
1. Regel „Video prüfen pro Kernseite“ enthalten.
2. NLM-Delegation für Skript/Storyboard dokumentiert.
3. Integrationsstatus pro Seite vorgesehen.

Tests:
- Check gegen 5 Pflichtfelder pro Seite.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- Video-Policy Tabelle angelegt.

### Task 04
Task-ID: LOOP-001-T04  
Titel: NLM Promptvorlagen-Katalog harmonisieren  
Kategorie: Dokumentation  
Priorität: P0

Ziel:
Video/Infografik/Präsentation/Tabelle Vorlagen einheitlich standardisieren.

Read First:
- `../∞Best∞Practices∞Loop.md`

Edit:
- `../∞Best∞Practices∞Loop.md`

Akzeptanzkriterien:
1. Alle 4 Vorlagen folgen gleicher Struktur.
2. Pflichtblöcke (Ziel, Quellen, Qualitätscheck, Output) vorhanden.
3. Keine widersprüchlichen Begriffe.

Tests:
- Vorlagen-Kreuzprüfung mit Struktur-Check.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- Vergleichsmatrix „Template Alignment“.

### Task 05
Task-ID: LOOP-001-T05  
Titel: Universal Command-to-Endpoint Mapping  
Kategorie: Architektur  
Priorität: P0

Ziel:
Sicherstellen, dass steuerbare Funktionen per Command + Endpoint abbildbar sind.

Read First:
- `COMMANDS.md` (falls vorhanden)
- `ENDPOINTS.md` (falls vorhanden)

Edit:
- `COMMANDS.md`
- `ENDPOINTS.md`

Akzeptanzkriterien:
1. Mapping-Tabelle vorhanden.
2. Fehlende Gegenstücke markiert.
3. Auth-Anforderung je Endpoint dokumentiert.

Tests:
- 1:1 Mapping Check.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- Mapping-Report.

### Task 06
Task-ID: LOOP-001-T06  
Titel: NLM Qualitäts-Scoring operationalisieren  
Kategorie: Reliability  
Priorität: P1

Ziel:
NLM-Ausgaben mit klarer Bewertungsroutine freigeben oder verwerfen.

Read First:
- `../∞Best∞Practices∞Loop.md`

Edit:
- `AGENTS-PLAN.md`
- `MEETING.md` (falls vorhanden)

Akzeptanzkriterien:
1. Scorecard pro NLM-Artefakt vorhanden.
2. Freigabeschwelle dokumentiert.
3. Reject-Workflow definiert.

Tests:
- Trockenlauf mit Beispielartefakt.

Doku-Updates:
- `CHANGELOG.md`

Evidenz:
- Scorecard ausgefüllt.

### Task 07
Task-ID: LOOP-001-T07  
Titel: Universal Security-Layer für NLM-Content  
Kategorie: Security  
Priorität: P0

Ziel:
Sicherheits- und Compliance-Prüfung für generierte Inhalte standardisieren.

Read First:
- `SECURITY.md` (falls vorhanden)

Edit:
- `SECURITY.md`
- `INTEGRATION.md` (falls vorhanden)

Akzeptanzkriterien:
1. Kein Overclaim ohne Evidenz.
2. Keine sensiblen Daten in Content-Artefakten.
3. Review-Pfad dokumentiert.

Tests:
- Security-Checklist auf Musteroutput anwenden.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- NLM Security Guardrails Abschnitt.

### Task 08
Task-ID: LOOP-001-T08  
Titel: Datentabellen-Norm für Entscheidungen  
Kategorie: Feature  
Priorität: P1

Ziel:
Einheitliche Tabellenstandards für KPI- und Entscheidungsdaten schaffen.

Read First:
- `SUPABASE.md` (falls vorhanden)
- `ENDPOINTS.md` (falls vorhanden)

Edit:
- `SUPABASE.md`
- `ENDPOINTS.md`

Akzeptanzkriterien:
1. Spaltenkatalog-Standard vorhanden.
2. Typ/Einheit/Zeitbezug verpflichtend.
3. Qualitätswarnungen definiert.

Tests:
- Schema-Validierungscheck anhand Muster.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- „Data Table Standard“ Abschnitt.

### Task 09
Task-ID: LOOP-001-T09  
Titel: Präsentations-Storyline Standardisieren  
Kategorie: Enablement  
Priorität: P1

Ziel:
Ein universelles Executive-Deck-Template für Entscheidungen etablieren.

Read First:
- `CONTEXT.md` (falls vorhanden)
- `ARCHITECTURE.md` (falls vorhanden)

Edit:
- `ONBOARDING.md` (falls vorhanden)
- `WEBSITE.md` (falls vorhanden)

Akzeptanzkriterien:
1. Folienlogik klar und reproduzierbar.
2. Risiko- und Trade-off-Folien enthalten.
3. FAQ-Block enthalten.

Tests:
- Storyline-Check auf Konsistenz.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- Template-Outline vorhanden.

### Task 10
Task-ID: LOOP-001-T10  
Titel: Infografik-Informationshierarchie festlegen  
Kategorie: UX  
Priorität: P1

Ziel:
Infografiken schnell verständlich und konsistent machen.

Read First:
- `WEBSITE.md` (falls vorhanden)
- `WEBAPP.md` (falls vorhanden)

Edit:
- `WEBSITE.md`
- `WEBAPP.md`

Akzeptanzkriterien:
1. Kernaussagen auf max. 5 begrenzt.
2. Visual Mapping dokumentiert.
3. Accessibility-Hinweise enthalten.

Tests:
- 30-Sekunden-Lesbarkeitsprüfung.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- Infografik-Blueprint.

### Task 11
Task-ID: LOOP-001-T11  
Titel: NLM Delegationslog verpflichtend einführen  
Kategorie: Reliability  
Priorität: P1

Ziel:
Jede NLM-Nutzung nachvollziehbar protokollieren.

Read First:
- `MEETING.md` (falls vorhanden)

Edit:
- `MEETING.md`

Akzeptanzkriterien:
1. Anlass, Vorlage, Score, Übernahmegrad erfasst.
2. Verworfenes mit Grund protokolliert.
3. Wiederauffindbarkeit gewährleistet.

Tests:
- Ein Probeeintrag vollständig ausgefüllt.

Doku-Updates:
- `CHANGELOG.md`

Evidenz:
- Delegationsprotokoll-Template.

### Task 12
Task-ID: LOOP-001-T12  
Titel: Universaler Reject-and-Refine Workflow  
Kategorie: Debug  
Priorität: P1

Ziel:
Schwache NLM-Ausgaben systematisch verbessern statt ad hoc neu zu generieren.

Read First:
- `../∞Best∞Practices∞Loop.md`

Edit:
- `TROUBLESHOOTING.md` (falls vorhanden)
- `AGENTS.md` (falls vorhanden)

Akzeptanzkriterien:
1. Fehlerklassen definiert.
2. Prompt-Schärfungsschritte dokumentiert.
3. Vergleich zweier Iterationen vorgesehen.

Tests:
- Simulierter Fehlerfall durchlaufen.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- Reject-and-Refine Playbook.

### Task 13
Task-ID: LOOP-001-T13  
Titel: Universaler Asset-Namensstandard  
Kategorie: Architektur  
Priorität: P2

Ziel:
Dateibenennung für NLM-Artefakte über Projekte konsistent halten.

Read First:
- `INTEGRATION.md` (falls vorhanden)

Edit:
- `INTEGRATION.md`
- `WEBSITE.md` (falls vorhanden)

Akzeptanzkriterien:
1. Benennungsschema dokumentiert.
2. Versionierungssuffixe definiert.
3. Dateityp-Regeln enthalten.

Tests:
- Namensbeispiele gegen Regeln prüfen.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- Naming Conventions Abschnitt.

### Task 14
Task-ID: LOOP-001-T14  
Titel: NLM Artefakt-Lifecycle definieren  
Kategorie: Operations  
Priorität: P1

Ziel:
Lebenszyklus von Erstellung bis Archivierung standardisieren.

Read First:
- `INFRASTRUCTURE.md` (falls vorhanden)

Edit:
- `INFRASTRUCTURE.md`
- `INTEGRATION.md` (falls vorhanden)

Akzeptanzkriterien:
1. Zustände: draft/review/approved/retired definiert.
2. Verantwortliche je Zustand definiert.
3. Archivierungsregeln beschrieben.

Tests:
- Lifecycle auf Beispielasset anwenden.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- Asset Lifecycle Tabelle.

### Task 15
Task-ID: LOOP-001-T15  
Titel: Universal KPI-Set für Content-Assets  
Kategorie: Performance  
Priorität: P1

Ziel:
Messbare Wirkung von Video/Infografik/Präsentation/Tabelle etablieren.

Read First:
- `WEBSITE.md` (falls vorhanden)
- `WEBAPP.md` (falls vorhanden)

Edit:
- `WEBSITE.md`
- `WEBAPP.md`
- `CONTEXT.md` (falls vorhanden)

Akzeptanzkriterien:
1. KPI pro Asset-Typ definiert.
2. Baseline und Zielwert als Platzhalter vorhanden.
3. Review-Intervall definiert.

Tests:
- KPI-Liste gegen Asset-Typen prüfen.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- KPI-Grid.

### Task 16
Task-ID: LOOP-001-T16  
Titel: Universal Prompt-Library Index erstellen  
Kategorie: Enablement  
Priorität: P1

Ziel:
Schneller Zugriff auf freigegebene NLM-Promptvorlagen.

Read First:
- `../∞Best∞Practices∞Loop.md`

Edit:
- `NOTEBOOKLM.md` (falls vorhanden)
- `ONBOARDING.md` (falls vorhanden)

Akzeptanzkriterien:
1. Vorlagen indexiert.
2. Einsatzfall je Vorlage dokumentiert.
3. Qualitätswarnungen je Vorlage enthalten.

Tests:
- Index-Navigation auf Vollständigkeit prüfen.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- Prompt-Library Tabelle.

### Task 17
Task-ID: LOOP-001-T17  
Titel: Universal Rollenschnitt für Agenten und NLM  
Kategorie: Security  
Priorität: P1

Ziel:
Verantwortlichkeiten und Rechte bei NLM-Content sauber trennen.

Read First:
- `AGENTS.md` (falls vorhanden)
- `SECURITY.md` (falls vorhanden)

Edit:
- `AGENTS.md`
- `SECURITY.md`

Akzeptanzkriterien:
1. Rollenmodell dokumentiert.
2. Freigabeinstanz je Asset-Typ benannt.
3. Least-Privilege berücksichtigt.

Tests:
- Rollenrechte gegen Prozessschritte prüfen.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- RACI-Matrix für NLM-Prozesse.

### Task 18
Task-ID: LOOP-001-T18  
Titel: Universaler Compliance-Check für NLM Artefakte  
Kategorie: Compliance  
Priorität: P1

Ziel:
Rechtliche und regulatorische Risiken bei generierten Inhalten minimieren.

Read First:
- `SECURITY.md`
- `CODE_OF_CONDUCT.md` (falls vorhanden)

Edit:
- `SECURITY.md`
- `CODE_OF_CONDUCT.md`

Akzeptanzkriterien:
1. Compliance-Checkliste pro Asset-Typ vorhanden.
2. Eskalationsweg bei Verstoß dokumentiert.
3. Abnahme-Pflicht vor Veröffentlichung beschrieben.

Tests:
- Checkliste auf Probeoutput anwenden.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- Compliance-Check-Template.

### Task 19
Task-ID: LOOP-001-T19  
Titel: Cross-Doc Konsistenzprüfung durchführen  
Kategorie: Reliability  
Priorität: P0

Ziel:
Sicherstellen, dass alle Dokumente denselben Standard abbilden.

Read First:
- `../∞Best∞Practices∞Loop.md`
- alle vorhandenen Pflichtdokumente

Edit:
- `AGENTS-PLAN.md`

Akzeptanzkriterien:
1. Widerspruchsliste erstellt.
2. P0-Inkonsistenzen aufgelöst oder eskaliert.
3. Offene Punkte priorisiert.

Tests:
- Konsistenzmatrix 10/10 geprüft.

Doku-Updates:
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- Cross-Doc Audit-Protokoll.

### Task 20
Task-ID: LOOP-001-T20  
Titel: All-in-One Verification und LOOP-002 erzeugen  
Kategorie: Abschluss  
Priorität: P0

Ziel:
Vollprüfung durchführen und nächsten 20er-Zyklus ausrollen.

Read First:
- alle im Zyklus betroffenen Dokumente

Edit:
- `AGENTS-PLAN.md`
- `MEETING.md`
- `CHANGELOG.md`

Akzeptanzkriterien:
1. Integrationscheck abgeschlossen.
2. NLM-Artefakte gegen Qualitätsmatrix validiert.
3. Offene Risiken priorisiert.
4. LOOP-002 mit 20 neuen Tasks angelegt.

Tests:
- Vollständige Gate-Prüfung.

Doku-Updates:
- `AGENTS-PLAN.md`
- `MEETING.md`
- `CHANGELOG.md`

Evidenz:
- Abschlussreport + neue Taskliste.

## Abschlussreport-Vorlage (Task 20)
1. Umgesetzt
2. Geänderte Dateien
3. Prüfungen und Ergebnisse
4. Risiken und offene Punkte
5. Nächste 20 Tasks

## LOOP-002 Placeholder
- Wird in Task 20 erzeugt.
- Muss wieder exakt 20 Tasks enthalten.
- Muss universal und projektagnostisch bleiben.

## LOOP-002 (Universal, 20 Tasks)
- Zyklus-ID: LOOP-002
- Fokus: Cross-Doc Harmonisierung, Betriebsreife, Verifikationshärte

### Task 01
Task-ID: LOOP-002-T01
Titel: Cross-Doc Konsistenzmatrix aktualisieren
Kategorie: Reliability
Priorität: P0
Read First: `ARCHITECTURE.md`, `COMMANDS.md`, `ENDPOINTS.md`, `MAPPING-COMMANDS-ENDPOINTS.md`
Edit: `MAPPING-COMMANDS-ENDPOINTS.md`
Akzeptanzkriterien: Mapping vollständig und widerspruchsfrei
Tests: 1:1 Mapping-Check

### Task 02
Task-ID: LOOP-002-T02
Titel: README Navigationsqualität schärfen
Kategorie: Dokumentation
Priorität: P1
Read First: `../README.md`, `ONBOARDING.md`
Edit: `../README.md`
Akzeptanzkriterien: klare Startpfade für User/Dev/Admin
Tests: Link- und Vollständigkeitscheck

### Task 03
Task-ID: LOOP-002-T03
Titel: NLM Prompt-Library im Betrieb verankern
Kategorie: Enablement
Priorität: P1
Read First: `NOTEBOOKLM.md`, `../∞Best∞Practices∞Loop.md`
Edit: `NOTEBOOKLM.md`, `ONBOARDING.md`
Akzeptanzkriterien: alle 4 NLM-Assettypen mit Nutzungsroute dokumentiert
Tests: Library-Index Check

### Task 04
Task-ID: LOOP-002-T04
Titel: Website Journey-Konsistenz prüfen
Kategorie: UX
Priorität: P0
Read First: `WEBSITE.md`, `WEBAPP.md`
Edit: `WEBSITE.md`
Akzeptanzkriterien: jede Seite hat klares Ziel, CTA und Folgepfad
Tests: Journey-Flow Review

### Task 05
Task-ID: LOOP-002-T05
Titel: Webapp Flow-Kompatibilität prüfen
Kategorie: UX
Priorität: P0
Read First: `WEBAPP.md`, `ENDPOINTS.md`
Edit: `WEBAPP.md`
Akzeptanzkriterien: Kernflows an Commands/Endpoints gekoppelt
Tests: Flow-to-API Check

### Task 06
Task-ID: LOOP-002-T06
Titel: Webshop Betriebsfähigkeit absichern
Kategorie: Feature
Priorität: P1
Read First: `WEBSHOP.md`, `SECURITY.md`
Edit: `WEBSHOP.md`
Akzeptanzkriterien: Checkout-Risiken und Prüfpfade dokumentiert
Tests: Checkout-Review Checkliste

### Task 07
Task-ID: LOOP-002-T07
Titel: Security Kontrollmatrix erweitern
Kategorie: Security
Priorität: P0
Read First: `SECURITY.md`, `INTEGRATION.md`
Edit: `SECURITY.md`
Akzeptanzkriterien: Kontrollen für API, NLM, Integrationen vollständig
Tests: Security-Matrix Review

### Task 08
Task-ID: LOOP-002-T08
Titel: Supabase RLS Template schärfen
Kategorie: Security
Priorität: P0
Read First: `SUPABASE.md`, `ENDPOINTS.md`
Edit: `SUPABASE.md`
Akzeptanzkriterien: RLS pro Kernbereich eindeutig beschrieben
Tests: RLS-Policy Vollständigkeitscheck

### Task 09
Task-ID: LOOP-002-T09
Titel: CI/CD Gate-Härtung dokumentieren
Kategorie: Reliability
Priorität: P1
Read First: `CI-CD-SETUP.md`, `GITHUB.md`
Edit: `CI-CD-SETUP.md`, `GITHUB.md`
Akzeptanzkriterien: Gate-Regeln und Rollback klar und konfliktfrei
Tests: Pipeline-Review gegen Gate-Liste

### Task 10
Task-ID: LOOP-002-T10
Titel: Infrastruktur Recovery-Drill ergänzen
Kategorie: Operations
Priorität: P1
Read First: `INFRASTRUCTURE.md`, `TROUBLESHOOTING.md`
Edit: `INFRASTRUCTURE.md`, `TROUBLESHOOTING.md`
Akzeptanzkriterien: Recovery-Ablauf vollständig beschrieben
Tests: Restore-Checkliste vorhanden

### Task 11
Task-ID: LOOP-002-T11
Titel: OpenClaw Fehlerpfad vertiefen
Kategorie: Integration
Priorität: P1
Read First: `OPENCLAW.md`, `INTEGRATION.md`
Edit: `OPENCLAW.md`
Akzeptanzkriterien: Auth- und Retry-Fehlerfälle vollständig
Tests: Fehlerklassentabelle vorhanden

### Task 12
Task-ID: LOOP-002-T12
Titel: n8n Workflow-Qualitätsgates ergänzen
Kategorie: Integration
Priorität: P1
Read First: `N8N.md`, `INTEGRATION.md`
Edit: `N8N.md`
Akzeptanzkriterien: Trigger/Input/Output/Recovery pro Workflow beschrieben
Tests: Workflow-Checkliste

### Task 13
Task-ID: LOOP-002-T13
Titel: Vercel Betriebscheckliste ausbauen
Kategorie: Operations
Priorität: P2
Read First: `VERCEL.md`, `vercel.json`
Edit: `VERCEL.md`
Akzeptanzkriterien: Pre-/Post-Deploy und Rollback-Checks vollständig
Tests: Checklisten-Review

### Task 14
Task-ID: LOOP-002-T14
Titel: IONOS DNS Betriebscheckliste ausbauen
Kategorie: Operations
Priorität: P2
Read First: `IONOS.md`, `CLOUDFLARE.md`
Edit: `IONOS.md`
Akzeptanzkriterien: DNS/TLS-Checkliste vollständig
Tests: DNS-TLS-Checklist Review

### Task 15
Task-ID: LOOP-002-T15
Titel: Blueprint Delta-Tracking operationalisieren
Kategorie: Architektur
Priorität: P1
Read First: `BLUEPRINT.md`, `ARCHITECTURE.md`
Edit: `BLUEPRINT.md`
Akzeptanzkriterien: Soll-Ist-Deltas mit Prioritäten erfasst
Tests: Delta-Matrix Vollständigkeitscheck

### Task 16
Task-ID: LOOP-002-T16
Titel: Contribution Workflow schärfen
Kategorie: Enablement
Priorität: P1
Read First: `CONTRIBUTING.md`, `GITHUB.md`
Edit: `CONTRIBUTING.md`
Akzeptanzkriterien: Beitragspfad und Pflichtchecks eindeutig
Tests: PR-Template Konsistenzcheck

### Task 17
Task-ID: LOOP-002-T17
Titel: Onboarding Quickstart synchronisieren
Kategorie: Enablement
Priorität: P1
Read First: `ONBOARDING.md`, `../README.md`
Edit: `ONBOARDING.md`
Akzeptanzkriterien: Schnellstart mit aktuellen Dokumenten konsistent
Tests: Rolle-zu-Pfad Prüfung

### Task 18
Task-ID: LOOP-002-T18
Titel: User-Plan Priorisierung schärfen
Kategorie: Dokumentation
Priorität: P2
Read First: `USER-PLAN.md`, `AGENTS-PLAN.md`
Edit: `USER-PLAN.md`
Akzeptanzkriterien: User-Aufgaben klar priorisiert und verifizierbar
Tests: Prioritäts- und Nachweischeck

### Task 19
Task-ID: LOOP-002-T19
Titel: Gesamtkonsistenz final prüfen
Kategorie: Reliability
Priorität: P0
Read First: alle Kern-Dokumente
Edit: `AGENTS-PLAN.md`
Akzeptanzkriterien: P0-Widersprüche null oder eskaliert
Tests: Cross-Doc Audit

### Task 20
Task-ID: LOOP-002-T20
Titel: All-in-One Verification und LOOP-003 erzeugen
Kategorie: Abschluss
Priorität: P0
Read First: gesamter aktueller Stand
Edit: `AGENTS-PLAN.md`, `MEETING.md`, `CHANGELOG.md`
Akzeptanzkriterien: Vollprüfung dokumentiert, Risiken priorisiert, nächste 20 Tasks erstellt
Tests: Gesamt-Gate-Check

---
