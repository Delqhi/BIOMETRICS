# MEETING.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Entscheidungen aus Meetings müssen Rule-/Control-Bezug enthalten.
- Offene Risiken, Exceptions und Follow-ups sind explizit zu tracken.
- Beschlüsse ohne Owner, Termin und Evidence gelten als unvollständig.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Laufendes Entscheidungs- und Abstimmungsprotokoll zwischen Orchestrator, Subagenten und User.

## Universalitäts-Regeln
1. Keine projektspezifischen Secrets im Protokoll.
2. Keine sensiblen personenbezogenen Daten.
3. Entscheidungen immer mit Kontext und Auswirkung dokumentieren.

## Eintragsregeln
Jede Arbeitsrunde erzeugt einen neuen Eintrag.  
Jeder Eintrag muss diese Felder enthalten:

```text
Meeting-ID:
Zeitpunkt:
Teilnehmer:
Kontext:
Ziel der Runde:
Entscheidungen:
Begründung:
Risiken:
Konflikte/Blocker:
Gelöste Punkte:
Offene Punkte:
Nächste Schritte:
Betroffene Dateien:
Referenz auf Task-IDs:
```

## Decision Log (Universal Template)

### Eintrag 001
Meeting-ID: M-001  
Zeitpunkt: {DATE_TIME}  
Teilnehmer: {OWNER_ROLE}, {AGENT_ROLE}

Kontext:
Initiale Rahmenklärung für universellen Orchestrator-Betrieb.

Ziel der Runde:
Startfähige Standards und Prioritäten bestätigen.

Entscheidungen:
1. NLM-CLI Nutzung als Pflicht gesetzt.
2. 20-Task-Loop als Standardprozess gesetzt.
3. Evidenzpflicht für Done-Meldungen bestätigt.

Begründung:
Verhindert Interpretationsspielraum und Schein-Fortschritt.

Risiken:
- Unklare Delegation an Subagenten führt zu Qualitätsverlust.

Konflikte/Blocker:
- Keine.

Gelöste Punkte:
- Basis-Governance und Kontrollmechanismen definiert.

Offene Punkte:
- Konkrete Projektausprägung mit echten Parametern.

Nächste Schritte:
1. AGENTS-PLAN LOOP-001 ausführen.
2. Cross-Doc Konsistenzprüfung durchführen.

Betroffene Dateien:
- `../∞Best∞Practices∞Loop.md`
- `AGENTS-PLAN.md`

Referenz auf Task-IDs:
- LOOP-001-T01 bis LOOP-001-T20

### Eintrag 002
Meeting-ID: M-002  
Zeitpunkt: 2026-02-17  
Teilnehmer: Orchestrator, Agent

Kontext:
Erweiterung des universellen Betriebsrahmens um Architektur-, Sicherheits-, Integrations- und Troubleshooting-Dokumente.

Ziel der Runde:
Fehlende Core-Dokumente im gleichen projektagnostischen Standard erstellen.

Entscheidungen:
1. Fünf Kern-Dateien werden universell mit Platzhaltern erstellt.
2. NLM-CLI Pflicht in allen relevanten Dokumenten konsistent verankern.
3. Abnahme über Fehlerprüfung und Konsistenzprüfung durchführen.

Begründung:
Vollständiger, wiederverwendbarer Grundrahmen reduziert Interpretationsspielraum für Subagenten.

Risiken:
- Projektspezifische Werte müssen später bewusst und konsistent eingesetzt werden.

Konflikte/Blocker:
- Keine.

Gelöste Punkte:
- Architektur-, Security-, Supabase-, Integrations- und Troubleshooting-Baselines vorhanden.

Offene Punkte:
- Weitere Pflichtdateien aus der Gesamtliste sukzessive ergänzen.

Nächste Schritte:
1. Zusätzliche Pflichtdateien im Universal-Standard erstellen.
2. Cross-Doc Konsistenzmatrix ausführen.
3. LOOP-002 vorbereiten.

Betroffene Dateien:
- `ARCHITECTURE.md`
- `SECURITY.md`
- `SUPABASE.md`
- `INTEGRATION.md`
- `TROUBLESHOOTING.md`

Referenz auf Task-IDs:
- LOOP-001-T14
- LOOP-001-T15
- LOOP-001-T18
- LOOP-001-T19
- LOOP-001-T20

### Eintrag 003
Meeting-ID: M-003  
Zeitpunkt: 2026-02-17  
Teilnehmer: Orchestrator, Agent

Kontext:
Erstellung der nächsten universellen Pflichtdateien für Betriebsreife und Integrationsfähigkeit.

Ziel der Runde:
NOTEBOOKLM, INFRASTRUCTURE, CI-CD-SETUP, OPENCLAW, N8N und WEBSITE im gleichen Standard ergänzen.

Entscheidungen:
1. Alle neuen Dateien bleiben strikt projektagnostisch.
2. NLM-CLI Pflicht und Qualitätsmatrix werden in den relevanten Dateien verankert.
3. Website erhält verpflichtende Video-/Asset-Felder für NLM-Integration.

Begründung:
Der Rahmen soll universell in jedem Projekt nutzbar sein, ohne Interpretationsspielraum für Subagenten.

Risiken:
- Platzhalter müssen in realen Projekten diszipliniert und konsistent konkretisiert werden.

Konflikte/Blocker:
- Keine.

Gelöste Punkte:
- Weitere zentrale Governance-Dateien sind als Universal-Templates verfügbar.

Offene Punkte:
- Restliche optionale Pflichtdokumente in gleicher Qualität ergänzen.

Nächste Schritte:
1. Cross-Doc Mapping für neue Dateien prüfen.
2. LOOP-002 Aufgaben auf neue Artefakte ausrichten.
3. Ergänzende Betriebsdokumente nachziehen.

Betroffene Dateien:
- `NOTEBOOKLM.md`
- `INFRASTRUCTURE.md`
- `CI-CD-SETUP.md`
- `OPENCLAW.md`
- `N8N.md`
- `WEBSITE.md`

Referenz auf Task-IDs:
- LOOP-001-T09
- LOOP-001-T10
- LOOP-001-T11
- LOOP-001-T14
- LOOP-001-T20

### Eintrag 004
Meeting-ID: M-004  
Zeitpunkt: 2026-02-17  
Teilnehmer: Orchestrator, Agent

Kontext:
Erstellung weiterer universeller Produkt-, Governance- und Compliance-Dokumente zur Vervollständigung des Basissystems.

Ziel der Runde:
WEBAPP, WEBSHOP, ENGINE, CLOUDFLARE, GITHUB, CONTRIBUTING, CODE_OF_CONDUCT und LICENSE ergänzen.

Entscheidungen:
1. Alle Dokumente bleiben strikt template-basiert und projektagnostisch.
2. Optionale Bereiche (Shop/Engine) enthalten klare NOT_APPLICABLE-Hinweise.
3. Governance-Dokumente werden auf denselben Qualitäts- und Nachweisstandard gehoben.

Begründung:
Subagenten sollen in jedem Projekt ohne Kontextverlust starten können.

Risiken:
- Ohne konsequente Platzhalterbefüllung kann operative Unschärfe entstehen.

Konflikte/Blocker:
- Keine.

Gelöste Punkte:
- Wesentliche Dokumentlücken im universellen Pflichtkatalog wurden geschlossen.

Offene Punkte:
- Restliche optionale Dokumente gemäß Gesamtliste ergänzen.

Nächste Schritte:
1. Cross-Doc Konsistenz für neue Artefakte prüfen.
2. LOOP-002 mit Fokus auf Harmonisierung und Mappings vorbereiten.
3. Zusätzliche Infrastruktur-/Deployment-Dokumente nachziehen.

Betroffene Dateien:
- `WEBAPP.md`
- `WEBSHOP.md`
- `ENGINE.md`
- `CLOUDFLARE.md`
- `GITHUB.md`
- `CONTRIBUTING.md`
- `CODE_OF_CONDUCT.md`
- `LICENSE`

Referenz auf Task-IDs:
- LOOP-001-T07
- LOOP-001-T08
- LOOP-001-T12
- LOOP-001-T17
- LOOP-001-T18
- LOOP-001-T20

### Eintrag 005
Meeting-ID: M-005  
Zeitpunkt: 2026-02-17  
Teilnehmer: Orchestrator, Agent

Kontext:
Schließung weiterer offener Pflichtartefakte und Verbesserung der Navigierbarkeit für Subagenten.

Ziel der Runde:
VERCEL, IONOS, BLUEPRINT, vercel.json sowie Mapping-Report erstellen und Website/Webapp/README harmonisieren.

Entscheidungen:
1. Fehlende Plattform- und Blueprint-Dokumente als universelle Templates ergänzt.
2. `BIOMETRICS/MAPPING-COMMANDS-ENDPOINTS.md` als expliziter Konsistenzreport eingeführt.
3. `../README.md` als zentralen Universal-Index mit Startreihenfolge ausgebaut.

Begründung:
Subagenten sollen ohne Suchaufwand in die richtige Dokumentkette einsteigen können.

Risiken:
- Projektteams müssen Platzhalter in Betriebsübernahmen konsequent konkretisieren.

Konflikte/Blocker:
- Keine.

Gelöste Punkte:
- Offene Dokumentlücken weiter reduziert.
- Crosslinks zwischen Website/Webapp/Command/API explizit gesetzt.

Offene Punkte:
- Restliche optionale Dokumente aus der Masterliste nachziehen.

Nächste Schritte:
1. Weitere Pflichtdokumente ergänzen.
2. LOOP-002 auf Konsistenz- und Betriebsfokus ausrichten.
3. Task-20 Gesamtprüfung erneut durchführen.

Betroffene Dateien:
- `VERCEL.md`
- `IONOS.md`
- `BLUEPRINT.md`
- `vercel.json`
- `BIOMETRICS/MAPPING-COMMANDS-ENDPOINTS.md`
- `WEBAPP.md`
- `WEBSITE.md`
- `../README.md`

Referenz auf Task-IDs:
- LOOP-001-T04
- LOOP-001-T05
- LOOP-001-T16
- LOOP-001-T19
- LOOP-001-T20

### Eintrag 006
Meeting-ID: M-006  
Zeitpunkt: 2026-02-17  
Teilnehmer: Orchestrator, Agent

Kontext:
Fortsetzung der Universalisierung mit LOOP-002, vertieften Betriebschecklisten und Basisartefakten.

Ziel der Runde:
LOOP-002 vollständig in AGENTS-PLAN verankern, Website/Webapp/Webshop/Plattformdokus schärfen und package/requirments Basis ergänzen.

Entscheidungen:
1. LOOP-002 mit exakt 20 universellen Tasks angelegt.
2. Vercel- und IONOS-Betriebschecklisten ergänzt.
3. Journey-/Flow-/Checkout-Kompatibilitätschecks in Website/Webapp/Webshop ergänzt.
4. `package.json` und `requirements.txt` als Basisartefakte hinzugefügt.

Begründung:
Die nächste Stufe reduziert operativen Interpretationsspielraum und verbessert Startfähigkeit für neue Subagenten.

Risiken:
- Platzhalter- und Projektparameter müssen in realen Projekten diszipliniert gepflegt werden.

Konflikte/Blocker:
- Keine.

Gelöste Punkte:
- LOOP-002 ist ausführbar dokumentiert.
- Plattform- und Journey-Checks sind operationalisiert.

Offene Punkte:
- Weitere optionale Dokumente aus der Masterliste ergänzen.

Nächste Schritte:
1. LOOP-002 Task 01 bis 05 priorisieren.
2. Cross-Doc Audit gemäß LOOP-002-T01 ausführen.
3. Task-20 Gesamtprüfung für LOOP-002 vorbereiten.

Betroffene Dateien:
- `AGENTS-PLAN.md`
- `VERCEL.md`
- `IONOS.md`
- `WEBSITE.md`
- `WEBAPP.md`
- `WEBSHOP.md`
- `../README.md`
- `package.json`
- `requirements.txt`

Referenz auf Task-IDs:
- LOOP-002-T01 bis LOOP-002-T20

### Eintrag 007
Meeting-ID: M-007  
Zeitpunkt: 2026-02-17  
Teilnehmer: Orchestrator, Agent

Kontext:
Ergänzung von MCP-Governance, kanonischer Verzeichnisstrategie und erweiterten NLM-Artefaktregeln.

Ziel der Runde:
`BIOMETRICS/` als Standardpfad verankern, `MCP.md` einführen und NLM um Mindmaps/Podcasts/Reports erweitern.

Entscheidungen:
1. Kanonischer Governance-Hauptordner auf `BIOMETRICS/` festgelegt.
2. `BIOMETRICS/MCP.md` als globales MCP-Book im Projekt eingeführt.
3. NLM-Artefaktkatalog um Mindmap, Podcast und Bericht erweitert.
4. README-Einbindungspflicht für freigegebene NLM-Artefakte verschärft.

Begründung:
Verbessert Nutzbarkeit für Menschen, reduziert Suchaufwand und erhöht Subagenten-Klarheit.

Risiken:
- Legacy-Dateien im Root bleiben bis geplanter Migration parallel bestehen.

Konflikte/Blocker:
- Keine.

Gelöste Punkte:
- MCP- und Asset-Governance klar operationalisiert.

Offene Punkte:
- Geplante Migration der Legacy-Dokumente nach `BIOMETRICS/` terminieren.

Nächste Schritte:
1. Migrationsplan Root -> `BIOMETRICS/` als Task ergänzen.
2. Erste reale NLM-Artefakte in Asset-Ordner aufnehmen.
3. LOOP-002 Konsistenzprüfung gegen neue Regeln durchführen.

Betroffene Dateien:
- `../∞Best∞Practices∞Loop.md`
- `NOTEBOOKLM.md`
- `../README.md`
- `BIOMETRICS/MCP.md`

Referenz auf Task-IDs:
- LOOP-002-T01
- LOOP-002-T03
- LOOP-002-T04
- LOOP-002-T17
- LOOP-002-T20

### Eintrag 008
Meeting-ID: M-008  
Zeitpunkt: 2026-02-17  
Teilnehmer: Orchestrator, Agent

Kontext:
Korrektur der Verzeichnisvorgaben nach User-Feedback.

Ziel der Runde:
Altverzeichnis ersetzen durch `BIOMETRICS/` als Hauptverzeichnis und `NLM-ASSETS/` als Asset-Bereich.

Entscheidungen:
1. Kanonischer Hauptordner ist `BIOMETRICS/`.
2. NLM-Artefakte liegen unter `BIOMETRICS/NLM-ASSETS/`.
3. `MCP.md` liegt kanonisch in `BIOMETRICS/MCP.md`.

Begründung:
Der User hat eine klare, einheitliche Zielstruktur festgelegt.

Risiken:
- Legacy-Dateien außerhalb `BIOMETRICS/` müssen schrittweise migriert werden.

Konflikte/Blocker:
- Keine.

Gelöste Punkte:
- Kernregeln und Referenzen auf neue Pfade angepasst.

Offene Punkte:
- Vollständige Bestandsmigration aller bereits vorhandenen Dateien in den Zielordner abschließen.

Nächste Schritte:
1. Migrationsplan für bestehende Root-Dateien ausführen.
2. Cross-Doc Pfadprüfung nach Migration.
3. Task-20 Gesamtprüfung erneut laufen lassen.

Betroffene Dateien:
- `../∞Best∞Practices∞Loop.md`
- `NOTEBOOKLM.md`
- `../README.md`
- `BIOMETRICS/MCP.md`

Referenz auf Task-IDs:
- LOOP-002-T01
- LOOP-002-T19
- LOOP-002-T20

### Eintrag 009
Meeting-ID: M-009  
Zeitpunkt: 2026-02-17  
Teilnehmer: Orchestrator, Agent

Kontext:
User-Korrektur zur Zielstruktur und Erweiterung des systemweiten Mapping-Ansatzes.

Ziel der Runde:
`BIOMETRICS/` und `BIOMETRICS/NLM-ASSETS/` endgültig als Standard fixieren und zusätzliche Mapping-Dateien einführen.

Entscheidungen:
1. Altverzeichnisse werden nicht mehr verwendet.
2. MCP-Book liegt nur noch unter `BIOMETRICS/MCP.md`.
3. Mapping-System wird auf Frontend/Backend/DB/NLM erweitert.
4. Website-Querverlinkung auf BIOMETRICS-Mappings umgestellt.

Begründung:
Klarere Hauptstruktur und stärkere Konsistenzsicherung über alle Schichten.

Risiken:
- Vollständige physische Migration aller Legacy-Rootdateien bleibt als separater Schritt notwendig.

Konflikte/Blocker:
- Löschen leerer Alt-Verzeichnisse ist toolseitig eingeschränkt.

Gelöste Punkte:
- Neue Mapping-Dateien in BIOMETRICS angelegt.
- Falsche MCP-/Mapping-Altdateien gelöscht.

Offene Punkte:
- Restliche Root-Legacy-Dateien nach `BIOMETRICS/` migrieren.

### Eintrag 011
Meeting-ID: M-011  
Zeitpunkt: 2026-02-17  
Teilnehmer: Orchestrator, Migration Subagent

Kontext:
Vollmigrationsrunde zur Kanonisierung aller agentenseitigen Dokumentpfade nach `BIOMETRICS/`.

Ziel der Runde:
Root-basierte aktive Referenzen entfernen, BIOMETRICS als zentrale Doku-Navigation festlegen und Migrationsstatus transparent dokumentieren.

Entscheidungen:
1. `../README.md` führt alle aktiven Projektdokumente unter `BIOMETRICS/`.
2. `AGENTS.md` enthält explizite Universalitätsformulierung für Website/Webshop/Webapp/Engine und weitere Projekttypen.
3. Hauptprompt enthält verbindliche BIOMETRICS-Erstellpflicht für agentenseitige Dateien.

Begründung:
Ein einziger kanonischer Doku-Pfad reduziert Fehlablagen und macht Subagenten-Übergaben konsistent.

Risiken:
- Leere Altordner können toolseitig bestehen bleiben und müssen bei verfügbarer Shell final entfernt werden.

Konflikte/Blocker:
- Kein funktionsfähiger Shell-/Task-Dateisystemanbieter in dieser Session für Massen-Moves.

Gelöste Punkte:
- Aktive Pfadreferenzen auf BIOMETRICS umgestellt.
- Altbenennungen in den Kernprotokollen bereinigt.

Offene Punkte:
- Physische Massenverschiebung sämtlicher Root-Dokumentdateien in einem Shell-fähigen Schritt durchführen.

Nächste Schritte:
1. Shell-fähige Session starten und Root-Dateien nach `BIOMETRICS/` verschieben.
2. Leere Altordner final löschen.
3. Vollständige Nachprüfung inkl. Fehlercheck erneut dokumentieren.

Betroffene Dateien:
- `AGENTS.md`
- `../README.md`
- `../∞Best∞Practices∞Loop.md`
- `MEETING.md`
- `CHANGELOG.md`

Referenz auf Task-IDs:
- LOOP-002-T19
- LOOP-002-T20

### Eintrag 013
Meeting-ID: M-013  
Zeitpunkt: 2026-02-18  
Teilnehmer: Orchestrator, Agent

Kontext:
BIOMETRICS-Projekt erweitert um Qwen 3.5 NVIDIA NIM Integration für leistungsstarke KI-Operationen.

Ziel der Runde:
OpenClaw und OpenCode mit Qwen 3.5 NVIDIA NIM konfigurieren und 50+ Tasks für die Erweiterung erstellen.

Entscheidungen:
1. OpenClaw mit NVIDIA NIM Provider konfiguriert (Modell: qwen/qwen3.5-397b-a17b).
2. OpenCode mit NVIDIA NIM Provider konfiguriert (Timeout: 120000ms).
3. 50+ neue Tasks im Todo-System erstellt für BIOMETRICS-Erweiterungen.
4. AGENTS.md um Qwen 3.5 Skills erweitert.

Begründung:
Qwen 3.5 397B bietet überlegene Code- und Reasoning-Fähigkeiten für komplexe BIOMETRICS-Automatisierung.

Risiken:
- NVIDIA NIM Free Tier hat Rate Limit von 40 RPM.

Konflikte/Blocker:
- Keine.

Gelöste Punkte:
- NVIDIA NIM Integration vollständig konfiguriert.
- Beide Tools (OpenClaw + OpenCode) nutzen jetzt Qwen 3.5.

Offene Punkte:
- Weitere NVIDIA NIM Modelle evaluieren (Kimi K2.5, Llama 3.3).

Nächste Schritte:
1. Erste Tests mit Qwen 3.5 durchführen.
2. Weitere NVIDIA NIM Modelle konfigurieren.
3. BIOMETRICS-Erweiterungen basierend auf 50+ Tasks umsetzen.

Betroffene Dateien:
- `~/.openclaw/openclaw.json`
- `~/.config/opencode/opencode.json`
- `BIOMETRICS/AGENTS.md`
- `BIOMETRICS/MEETING.md`

Referenz auf Task-IDs:
- BIOMETRICS-TASK-001 bis BIOMETRICS-TASK-050+

---

## NLM Delegationsprotokoll (Pflicht)
Jede Delegation an NLM wird separat erfasst:

```text
NLM-Delegation-ID:
Datum/Zeit:
Asset-Typ: Video | Infografik | Präsentation | Datentabelle
Verwendete Vorlage:
Input-Quellen:
Ergebnis-Score:
Übernommen:
Verworfen:
Grund für Verwerfung:
Nächste Iteration:
Owner:
```

## Konflikt- und Eskalationslog

```text
Eskalations-ID:
Severity: P0|P1|P2
Beschreibung:
Auswirkung:
Sofortmaßnahme:
Langfristmaßnahme:
Owner:
Status:
```

## Abnahme-Check MEETING
1. Jede Session dokumentiert.
2. Entscheidungen mit Begründung versehen.
3. Risiken und offene Punkte explizit.
4. Task-Referenzen vorhanden.
5. NLM-Delegationen nachvollziehbar.
- 2026-02-17 12:06: Migration Root-Dateien nach BIOMETRICS abgeschlossen. Migriert: 34; Nicht möglich: 0.

---
