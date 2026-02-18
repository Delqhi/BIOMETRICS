# AGENTS.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Lokale Agentenregeln sind eine konkrete Ausprägung von `AGENTS-GLOBAL.md`.
- Delegations-, Todo- und Evidence-Disziplin sind zwingend.
- Abweichungen sind nur als dokumentierte Overrides zulässig.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Projektlokale Arbeitsregeln für Orchestrator und Subagenten. Diese Datei ist universell und projektagnostisch.
Die Regeln gelten explizit universell für Website, Webshop, Webapp, Engine und weitere Projekttypen.

## Grundprinzipien
1. Erst lesen, dann schreiben.
2. Keine Done-Meldung ohne Evidenz.
3. Keine Duplikatdateien, bestehende Struktur erweitern.
4. Keine Kommentare in Code-Dateien, außer in Markdown.
5. NLM immer vollumfänglich über NLM-CLI nutzen.
6. Promptvorlagen aus `../∞Best∞Practices∞Loop.md` verpflichtend nutzen.
7. Jede Änderung in `MEETING.md` und `CHANGELOG.md` dokumentieren.

## Stack-Policy
- Frontend: Next.js
- Backend: Go + Supabase
- JS-Paketmanager: pnpm

## Rollen
### Orchestrator
- priorisiert
- delegiert
- prüft Qualität
- validiert Evidenz
- steuert Task-20 Abschluss

### Subagent
- arbeitet in klarem Scope
- nutzt NLM-CLI bei Content-Artefakten
- liefert strukturierte Übergabe
- meldet Blocker frühzeitig

## Pflichtformat Subagenten-Auftrag
```text
ROLE:
GOAL:
CONTEXT:
READ FIRST:
EDIT ONLY:
DO NOT EDIT:
TASKS:
ACCEPTANCE CRITERIA:
REQUIRED TESTS:
REQUIRED DOC UPDATES:
RISKS:
OUTPUT FORMAT:
```

## NLM Pflichtsatz
Du musst NotebookLM vollständig über NLM-CLI nutzen, den passenden Vorlagenprompt verwenden, das Ergebnis gegen die NLM-Qualitätsmatrix bewerten und nur verifizierte, konsistente Inhalte übernehmen.

## Übergabeformat (Pflicht)
1. Was wurde geändert
2. Welche Dateien wurden geändert
3. Welche Prüfungen liefen
4. Welche Risiken bleiben
5. Nächste 3 Schritte

## Eskalation
- P0: sofort
- P1: innerhalb der Session
- P2: in nächsten 20er-Loop einplanen

## Qwen 3.5 Skills

Dieses Projekt nutzt Qwen 3.5 (NVIDIA NIM) für spezialisierte KI-Aufgaben. Die folgenden Skills sind verfügbar:

### qwen_vision_analysis
Bildanalyse und visuelle Erkennung für Produktbilder, Grafiken und Diagramme.
- **Use Case:** Produktbild-Qualitätsprüfung, Layout-Analyse
- **Input:** Bilder (PNG, JPG, WebP)
- **Output:** Strukturierte Analyse mit Tags und Metriken

### qwen_code_generation
Full-Stack Code-Generierung mit Next.js, Go und Supabase.
- **Use Case:** Komponenten, API-Routen, Datenbank-Schema
- **Input:** Natürliche Sprache oder Spezifikation
- **Output:** Fertiger, getesteter Code

### qwen_document_ocr
Texterkennung und Dokumentanalyse aus gescannten Dokumenten und PDFs.
- **Use Case:** Rechnungsverarbeitung, Vertragsanalyse
- **Input:** PDF, Bilder mit Text
- **Output:** Extrahierter Text, Metadaten, Struktur

### qwen_video_understanding
Video-Inhaltsanalyse für帧-Extraction und Szenenbeschreibung.
- **Use Case:** Video-Vorschau, Content-Indexierung
- **Input:** Videos (MP4, MOV, WebM)
- **Output:** Szenenbeschreibung, Key-Frames, Metadaten

### qwen_conversation
Natürliche Konversations-KI für Kundenservice und Chat-Interaktionen.
- **Use Case:** Support-Chat, Produktberatung
- **Input:** Benutzer-Nachrichten, Kontext
- **Output:** Kontextbezogene Antworten

## Abnahme-Check AGENTS
1. Regeln klar und widerspruchsfrei
2. NLM-CLI Pflicht enthalten
3. Rollen und Übergabeformat enthalten
4. Eskalationspfad enthalten

---
