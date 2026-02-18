# CHANGELOG.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Änderungen müssen Rule-/Control-relevant nachvollziehbar eingetragen werden.
- Sicherheits-, Mapping- und Governance-Änderungen sind explizit zu markieren.
- Keine stille Änderung ohne Audit-fähigen Verlauf.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Nachvollziehbare Historie aller relevanten Änderungen an Struktur, Regeln, Doku und Implementierung.

## Universalitäts-Regeln
1. Keine Secrets oder Tokens eintragen.
2. Einträge müssen reproduzierbar und verifizierbar sein.
3. Jede Änderung referenziert betroffene Dateien.

## Formatstandard

```text
## [YYYY-MM-DD HH:MM] Change-ID
Typ: Architecture | Feature | Reliability | Security | Docs | Ops
Zusammenfassung:
Hintergrund:
Geänderte Dateien:
Prüfungen:
Ergebnis:
Risiken/Offen:
Follow-ups:
Referenzen (Task-IDs/Meetings):
```

## Einträge

## [2026-02-17 00:00] CHG-001
Typ: Docs
Zusammenfassung:
Initiale Anlage einer universellen Changelog-Struktur.

Hintergrund:
Erforderlich für auditierbare, projektagnostische Nachverfolgbarkeit.

Geänderte Dateien:
- `CHANGELOG.md`

Prüfungen:
- Strukturprüfung gegen Changelog-Template.

Ergebnis:
- Erfolgreich.

Risiken/Offen:
- Konkrete Projekteinträge müssen fortlaufend ergänzt werden.

Follow-ups:
- Je Task ein Changelog-Eintrag ergänzen.

Referenzen (Task-IDs/Meetings):
- LOOP-001-T20
- M-001

## [2026-02-17 00:01] CHG-002
Typ: Docs
Zusammenfassung:
Universeller Betriebsrahmen mit AGENTS-PLAN, USER-PLAN, MEETING erstellt.

Hintergrund:
Subagenten sollen ohne Interpretationsspielraum arbeiten.

Geänderte Dateien:
- `AGENTS-PLAN.md`
- `USER-PLAN.md`
- `MEETING.md`
- `../∞Best∞Practices∞Loop.md`

Prüfungen:
- Dokumentkonsistenz
- Basis-Fehlerprüfung

Ergebnis:
- Erfolgreich.

Risiken/Offen:
- Projektkonkrete Werte für Platzhalter noch offen.

Follow-ups:
- LOOP-002 planen nach Task-20 Gesamtprüfung.

Referenzen (Task-IDs/Meetings):
- LOOP-001-T01 bis LOOP-001-T20
- M-001

## [2026-02-17 00:02] CHG-003
Typ: Docs
Zusammenfassung:
Universelle Kern-Dokumente für Architektur, Sicherheit, Supabase, Integration und Troubleshooting erstellt.

Hintergrund:
Der Betriebsrahmen sollte ohne Interpretationsspielraum weiter vervollständigt werden.

Geänderte Dateien:
- `ARCHITECTURE.md`
- `SECURITY.md`
- `SUPABASE.md`
- `INTEGRATION.md`
- `TROUBLESHOOTING.md`
- `MEETING.md`
- `CHANGELOG.md`

Prüfungen:
- Datei-Fehlerprüfung
- Struktur-/Konsistenzprüfung der neuen Templates

Ergebnis:
- Erfolgreich.

Risiken/Offen:
- Platzhalter sind noch projektspezifisch zu befüllen.

Follow-ups:
- Nächste Pflichtdateien im Universal-Standard ergänzen.
- Cross-Doc Mapping zwischen Commands/Endpoints/Architecture vertiefen.

Referenzen (Task-IDs/Meetings):
- LOOP-001-T14
- LOOP-001-T15
- LOOP-001-T18
- LOOP-001-T19
- LOOP-001-T20
- M-002

## [2026-02-17 00:03] CHG-004
Typ: Docs
Zusammenfassung:
Weitere universelle Pflichtdateien für NLM-Betrieb, Infrastruktur, CI/CD, OpenClaw, n8n und Website-Standards erstellt.

Hintergrund:
Der universelle Betriebsrahmen wurde um fehlende Kernartefakte erweitert.

Geänderte Dateien:
- `NOTEBOOKLM.md`
- `INFRASTRUCTURE.md`
- `CI-CD-SETUP.md`
- `OPENCLAW.md`
- `N8N.md`
- `WEBSITE.md`
- `MEETING.md`
- `CHANGELOG.md`

Prüfungen:
- Datei-Fehlerprüfung
- Strukturkonsistenzprüfung

Ergebnis:
- Erfolgreich.

Risiken/Offen:
- Platzhalter müssen projektspezifisch kontrolliert gefüllt werden.

Follow-ups:
- Weitere optionale Pflichtdateien im gleichen Standard ergänzen.
- LOOP-002 auf Basis erweiterter Dokumentlandschaft planen.

Referenzen (Task-IDs/Meetings):
- LOOP-001-T09
- LOOP-001-T10
- LOOP-001-T11
- LOOP-001-T14
- LOOP-001-T20
- M-003

## [2026-02-17 00:04] CHG-005
Typ: Docs
Zusammenfassung:
Universelle Produkt-, Governance- und Compliance-Dateien ergänzt (Webapp, Webshop, Engine, Cloudflare, GitHub, Contribution, Conduct, License).

Hintergrund:
Der universelle Basiskatalog wurde weiter vervollständigt, um Subagentenstart ohne Interpretationsspielraum zu ermöglichen.

Geänderte Dateien:
- `WEBAPP.md`
- `WEBSHOP.md`
- `ENGINE.md`
- `CLOUDFLARE.md`
- `GITHUB.md`
- `CONTRIBUTING.md`
- `CODE_OF_CONDUCT.md`
- `LICENSE`
- `MEETING.md`
- `CHANGELOG.md`

Prüfungen:
- Datei-Fehlerprüfung
- Strukturkonsistenzprüfung

Ergebnis:
- Erfolgreich.

Risiken/Offen:
- Platzhalter und optionale NOT_APPLICABLE-Entscheidungen müssen projektweise diszipliniert gepflegt werden.

Follow-ups:
- Restliche Pflichtdateien aus dem Gesamtkatalog ergänzen.
- LOOP-002 mit Cross-Doc Harmonisierung ausrollen.

Referenzen (Task-IDs/Meetings):
- LOOP-001-T07
- LOOP-001-T08
- LOOP-001-T12
- LOOP-001-T17
- LOOP-001-T18
- LOOP-001-T20
- M-004

## [2026-02-17 00:05] CHG-006
Typ: Docs
Zusammenfassung:
Weitere offene Pflichtartefakte ergänzt und Dokument-Navigation harmonisiert.

Hintergrund:
Der universelle Rahmen wurde um Plattform-Templates, Blueprint-Vorlage und Command/Endpoint-Mapping erweitert.

Geänderte Dateien:
- `VERCEL.md`
- `IONOS.md`
- `BLUEPRINT.md`
- `vercel.json`
- `BIOMETRICS/MAPPING-COMMANDS-ENDPOINTS.md`
- `WEBAPP.md`
- `WEBSITE.md`
- `../README.md`
- `MEETING.md`
- `CHANGELOG.md`

Prüfungen:
- Datei-Fehlerprüfung
- Struktur- und Linkkonsistenzprüfung

Ergebnis:
- Erfolgreich.

Risiken/Offen:
- Platzhalterbefüllung bleibt projektspezifische Pflichtaufgabe.

Follow-ups:
- Restliche optionale Pflichtdateien aus dem Gesamtkatalog ergänzen.
- LOOP-002 inklusive Task-20 Gesamtprüfung vorbereiten.

Referenzen (Task-IDs/Meetings):
- LOOP-001-T04
- LOOP-001-T05
- LOOP-001-T16
- LOOP-001-T19
- LOOP-001-T20
- M-005

## [2026-02-17 00:06] CHG-007
Typ: Docs
Zusammenfassung:
LOOP-002 mit 20 Tasks ergänzt, Betriebschecklisten vertieft und Basisartefakte package/requirements angelegt.

Hintergrund:
Der universelle Rahmen wurde für die nächste Ausführungsrunde operationalisiert.

Geänderte Dateien:
- `AGENTS-PLAN.md`
- `VERCEL.md`
- `IONOS.md`
- `WEBSITE.md`
- `WEBAPP.md`
- `WEBSHOP.md`
- `../README.md`
- `package.json`
- `requirements.txt`
- `MEETING.md`
- `CHANGELOG.md`

Prüfungen:
- Datei-Fehlerprüfung
- Struktur-/Konsistenzprüfung

Ergebnis:
- Erfolgreich.

Risiken/Offen:
- Projektparametrisierung und Platzhalterbefüllung bleiben operative Pflicht.

Follow-ups:
- LOOP-002 aktiv abarbeiten.
- Cross-Doc Audit und Task-20 Verifikation durchführen.

Referenzen (Task-IDs/Meetings):
- LOOP-002-T01 bis LOOP-002-T20
- M-006

## [2026-02-17 00:07] CHG-008
Typ: Docs
Zusammenfassung:
MCP-Governance und kanonische Verzeichnisstrategie ergänzt, NLM-Regeln um Mindmap/Podcast/Report erweitert.

Hintergrund:
Subagenten sollen mit klarer MCP- und Asset-Struktur arbeiten, inklusive stärkerer README-Aktualisierungspflicht.

Geänderte Dateien:
- `../∞Best∞Practices∞Loop.md`
- `NOTEBOOKLM.md`
- `../README.md`
- `BIOMETRICS/MCP.md`

Neue Verzeichnisse:
- `BIOMETRICS/NLM-ASSETS/videos/`
- `BIOMETRICS/NLM-ASSETS/infographics/`
- `BIOMETRICS/NLM-ASSETS/presentations/`
- `BIOMETRICS/NLM-ASSETS/reports/`
- `BIOMETRICS/NLM-ASSETS/tables/`
- `BIOMETRICS/NLM-ASSETS/mindmaps/`
- `BIOMETRICS/NLM-ASSETS/podcasts/`

Prüfungen:
- Datei-Fehlerprüfung
- Struktur-/Konsistenzprüfung

Ergebnis:
- Erfolgreich.

Risiken/Offen:
- Legacy-Root-Struktur bis geplanter Migration weiterhin vorhanden.

Follow-ups:
- Migrationsplan zu `BIOMETRICS/` ausarbeiten.
- Erste reale NLM-Artefakte im neuen Asset-Pfad referenzieren.

Referenzen (Task-IDs/Meetings):
- LOOP-002-T01
- LOOP-002-T03
- LOOP-002-T04
- LOOP-002-T17
- LOOP-002-T20
- M-007

## [2026-02-17 00:08] CHG-009
Typ: Docs
Zusammenfassung:
Verzeichnisstandard auf `BIOMETRICS/` und `BIOMETRICS/NLM-ASSETS/` korrigiert; MCP-Standardpfad auf `BIOMETRICS/MCP.md` gesetzt.

Hintergrund:
Userseitige Präzisierung der gewünschten Zielstruktur.

Geänderte Dateien:
- `../∞Best∞Practices∞Loop.md`
- `NOTEBOOKLM.md`
- `../README.md`
- `MEETING.md`
- `CHANGELOG.md`
- `BIOMETRICS/MCP.md`

Neue Verzeichnisse:
- `BIOMETRICS/NLM-ASSETS/videos/`
- `BIOMETRICS/NLM-ASSETS/infographics/`
- `BIOMETRICS/NLM-ASSETS/presentations/`
- `BIOMETRICS/NLM-ASSETS/reports/`
- `BIOMETRICS/NLM-ASSETS/tables/`
- `BIOMETRICS/NLM-ASSETS/mindmaps/`
- `BIOMETRICS/NLM-ASSETS/podcasts/`

Prüfungen:
- Datei-Fehlerprüfung
- Pfad-Referenzprüfung

Ergebnis:
- Erfolgreich.

Risiken/Offen:
- Vollständige Migration bestehender Legacy-Dateien in den Zielordner steht noch aus.

Follow-ups:
- Root-Dateien schrittweise nach `BIOMETRICS/` migrieren.
- Cross-Doc Pfadprüfung nach vollständiger Migration wiederholen.

Referenzen (Task-IDs/Meetings):
- LOOP-002-T01
- LOOP-002-T19
- LOOP-002-T20
- M-008

## [2026-02-17 00:09] CHG-010
Typ: Docs
Zusammenfassung:
Mapping-Framework systemweit erweitert (Frontend/Backend, DB/API, NLM-Assets) und Website-Referenzen auf BIOMETRICS umgestellt.

Hintergrund:
Userwunsch nach häufigeren und breiteren Konsistenz-Mappings über alle Schichten.

Geänderte Dateien:
- `../∞Best∞Practices∞Loop.md`
- `WEBSITE.md`
- `BIOMETRICS/MAPPING.md`
- `BIOMETRICS/MAPPING-COMMANDS-ENDPOINTS.md`
- `BIOMETRICS/MAPPING-FRONTEND-BACKEND.md`
- `BIOMETRICS/MAPPING-DB-API.md`
- `BIOMETRICS/MAPPING-NLM-ASSETS.md`
- `MEETING.md`
- `CHANGELOG.md`

Gelöschte Dateien:
- `MCP.md` (Altpfad)
- `BIOMETRICS/MAPPING-COMMANDS-ENDPOINTS.md` (frühere Altbenennung)

Prüfungen:
- Datei-Fehlerprüfung
- Pfad- und Referenzprüfung

Ergebnis:
- Erfolgreich.

Risiken/Offen:
- Leere Legacy-Ordner können toolseitig nicht immer vollständig entfernt werden.

Follow-ups:
- Verbleibende Root-Dateien nach `BIOMETRICS/` migrieren.
- Nach Migration erneute Vollkonsistenzprüfung durchführen.

Referenzen (Task-IDs/Meetings):
- LOOP-002-T01
- LOOP-002-T04
- LOOP-002-T05
- LOOP-002-T19
- LOOP-002-T20
- M-009

## [2026-02-17 00:10] CHG-011
Typ: Docs
Zusammenfassung:
Aktive Referenzen auf alte Root-Mappingdatei bereinigt und vollständig auf BIOMETRICS-Mapping umgestellt.

Hintergrund:
Konsistenzanforderung für klare Erstorientierung von Agenten.

Geänderte Dateien:
- `AGENTS-PLAN.md`
- `WEBAPP.md`
- `../README.md`
- `MEETING.md`
- `CHANGELOG.md`

Prüfungen:
- Referenzsuche auf alte Mapping-Datei
- Datei-Fehlerprüfung

Ergebnis:
- Erfolgreich.

Risiken/Offen:
- Physische Vollmigration der Legacy-Rootdateien nach `BIOMETRICS/` noch offen.

Follow-ups:
- Vollmigration planen und durchführen.
- Nachlaufende Cross-Doc Prüfung wiederholen.

Referenzen (Task-IDs/Meetings):
- LOOP-002-T19
- LOOP-002-T20
- M-012

## [2026-02-18 00:00] CHG-014
Typ: Docs
Zusammenfassung:
WORKFLOW.md auf 2520 Zeilen erweitert, Qwen 3.5 NVIDIA NIM Integration dokumentiert, OpenCode/OpenClaw Konfiguration vollständig abgeschlossen, biometrics-onboard CLI erstellt.

Hintergrund:
Projektweiter Qualitätsausbau mit Fokus auf Dokumentationstiefe und KI-Integration.

Geänderte Dateien:
- `BIOMETRICS/WORKFLOW.md` (erweitert auf 2520 Zeilen)
- `~/.config/opencode/opencode.json` (OpenCode Konfiguration)
- `~/.openclaw/openclaw.json` (OpenClaw Konfiguration)
- `BIOMETRICS/biometrics-cli/bin/biometrics-onboard.js` (CLI erstellt)

Neue Integrationen:
- Qwen 3.5 NVIDIA NIM (qwen/qwen3.5-397b-a17b)
- OpenCode NVIDIA NIM Provider konfiguriert
- OpenClaw NVIDIA NIM Provider konfiguriert
- biometrics-onboard CLI Befehl vollständig funktionsfähig

Prüfungen:
- Datei-Fehlerprüfung
- CLI-Funktionalitätstest
- NVIDIA NIM Konnektivitätstest

Ergebnis:
- Erfolgreich.

Risiken/Offen:
- NVIDIA NIM Rate Limits (40 RPM) beachten bei hoher Last.

Follow-ups:
- Weitere NVIDIA NIM Modelle evaluieren.
- Performance-Benchmarks mit Qwen 3.5 durchführen.

Referenzen (Task-IDs/Meetings):
- Task-07
- NVIDIA NIM Integration
- OpenCode/OpenClaw Setup

---

## Abnahme-Check CHANGELOG
1. Jeder substanzielle Schritt protokolliert.
2. Geänderte Dateien stets aufgeführt.
3. Prüfungen und Ergebnis enthalten.
4. Follow-ups und Risiken enthalten.

## [2026-02-17 00:20] CHG-012
Typ: Docs
Zusammenfassung:
Vollmigrations-Policy auf BIOMETRICS finalisiert, Root-Referenzen auf BIOMETRICS umgestellt und Altverzeichnis-Referenzen bereinigt.

Hintergrund:
Uservorgabe verlangt BIOMETRICS als kanonischen Zielpfad für agentenseitig erstellte Dokumente und konsistente Referenzierung.

Geänderte Dateien:
- `AGENTS.md`
- `../README.md`
- `../∞Best∞Practices∞Loop.md`
- `MEETING.md`
- `CHANGELOG.md`

Prüfungen:
- Referenzsuche auf Altverzeichnisnamen in `*.md`
- Referenzsuche auf Alt-Mappingdatei in `*.md`
- Fehlercheck geänderter Dateien

Ergebnis:
- Erfolgreich für Policy- und Referenzkonsistenz.

Risiken/Offen:
- Physische Massenverschiebung aller Root-Dokumentdateien ist wegen Shell-/Task-Dateisystem-Blocker in dieser Session noch offen.
- Leere Altordner können bis zur Shell-basierten Bereinigung bestehen bleiben.

Follow-ups:
- Shell-fähige Session nutzen und Root-Dateien nach `BIOMETRICS/` verschieben.
- Verbleibende leere Altordner entfernen.
- Vollständigen Nachlauf-Fehlercheck dokumentieren.

Referenzen (Task-IDs/Meetings):
- LOOP-002-T19
- LOOP-002-T20
- M-011
- 2026-02-17 12:06: Migration Root-Dateien nach BIOMETRICS abgeschlossen. Migriert: 34; Nicht möglich: 0.

## [2026-02-17 00:21] CHG-013
Typ: Docs
Zusammenfassung:
Hauptprompt um absolute BIOMETRICS-Pfadvorgabe und Repo-Namenskonvention `BIOMETRICS` ergänzt; README auf pending rename aktualisiert.

Hintergrund:
Uservorgabe: Erstagenten müssen sofort den korrekten Zielpfad und den gewünschten Repo-Namen erkennen.

Geänderte Dateien:
- `../∞Best∞Practices∞Loop.md`
- `../README.md`
- `BIOMETRICS/MEETING.md`
- `BIOMETRICS/CHANGELOG.md`

Prüfungen:
- Datei-Fehlerprüfung

Ergebnis:
- Erfolgreich.

Risiken/Offen:
- Direkte Repo-Umbenennung in dieser Session blockiert (ENOPRO bei Terminalzugriff).
- Leere Altordner `PROMPTOPS/` noch vorhanden.

Follow-ups:
- `gh repo rename BIOMETRICS --yes` in shell-fähiger Session ausführen.
- `PROMPTOPS/` nach erfolgreicher Bereinigung entfernen.

Referenzen (Task-IDs/Meetings):
- LOOP-002-T19
- LOOP-002-T20
- M-012

---
