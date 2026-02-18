# COMMANDS.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Befehle folgen dem globalen Sicherheits- und Ausführungsprotokoll.
- Kritische Commands benötigen Verifikations- und Rollback-Hinweise.
- Command-Änderungen müssen Mapping und Doku synchron aktualisieren.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Katalog steuerbarer Funktionen für KI/Agentensteuerung.

## Command-Schema (Pflicht)
```text
Command:
Zweck:
Rolle:
Input-Schema:
Output-Schema:
Fehlerfälle:
Nebenwirkungen:
Verifikation:
Endpoint-Referenz:
```

## Universeller Basis-Katalog

### CMD.SYSTEM.STATUS
Zweck:
System- und Doku-Status abrufen.

Rolle:
User, Admin, Agent

Input-Schema:
- scope: all | docs | runtime

Output-Schema:
- status_summary
- open_blockers
- next_actions

Fehlerfälle:
- unknown_scope

Nebenwirkungen:
- keine

Verifikation:
- Ergebnis enthält Status und offene Blocker.

Endpoint-Referenz:
- `API.SYSTEM.STATUS`

### CMD.TASKS.LIST
Zweck:
Tasks nach Priorität und Status abrufen.

Rolle:
User, Dev, Admin, Agent

Input-Schema:
- priority: P0|P1|P2|all
- status: OPEN|IN_PROGRESS|BLOCKED|DONE|all

Output-Schema:
- task_list[]

Fehlerfälle:
- invalid_filter

Nebenwirkungen:
- keine

Verifikation:
- Filter greifen korrekt.

Endpoint-Referenz:
- `API.TASKS.LIST`

### CMD.TASKS.UPDATE
Zweck:
Task-Status aktualisieren.

Rolle:
Dev, Admin, Agent

Input-Schema:
- task_id
- status
- evidence_ref

Output-Schema:
- updated_task

Fehlerfälle:
- task_not_found
- invalid_status

Nebenwirkungen:
- schreibt Protokoll in Meeting/Changelog

Verifikation:
- Status und Evidenz sichtbar.

Endpoint-Referenz:
- `API.TASKS.UPDATE`

### CMD.NLM.GENERATE.VIDEO
Zweck:
NLM-Videoartefakt anhand Standardvorlage generieren.

Rolle:
Dev, Agent

Input-Schema:
- topic
- audience
- business_goal
- source_refs[]

Output-Schema:
- titles[]
- scripts
- storyboard
- cta[]

Fehlerfälle:
- missing_sources
- quality_below_threshold

Nebenwirkungen:
- erzeugt Delegationsprotokoll

Verifikation:
- Score >= 13/16 und Korrektheit 2/2.

Endpoint-Referenz:
- `API.NLM.GENERATE.VIDEO`

### CMD.NLM.GENERATE.INFOGRAPHIC
Zweck:
NLM-Infografik-Spezifikation erzeugen.

Rolle:
Dev, Agent

Input-Schema:
- topic
- key_points[]
- source_refs[]

Output-Schema:
- layout_blueprint
- visual_mapping
- accessibility_notes

Fehlerfälle:
- inconsistent_claims

Nebenwirkungen:
- erzeugt Delegationsprotokoll

Verifikation:
- Konsistenzcheck bestanden.

Endpoint-Referenz:
- `API.NLM.GENERATE.INFOGRAPHIC`

### CMD.NLM.GENERATE.PRESENTATION
Zweck:
NLM-Präsentationsstruktur inkl. Sprechernotizen erzeugen.

Rolle:
Dev, Agent

Input-Schema:
- occasion
- audience_type
- decision_goal
- source_refs[]

Output-Schema:
- slide_outline[]
- speaker_notes[]
- risks_and_tradeoffs

Fehlerfälle:
- weak_storyline

Nebenwirkungen:
- erzeugt Delegationsprotokoll

Verifikation:
- Entscheidungsvorbereitung klar.

Endpoint-Referenz:
- `API.NLM.GENERATE.PRESENTATION`

### CMD.NLM.GENERATE.TABLE
Zweck:
NLM-Datentabellen-Spezifikation erzeugen.

Rolle:
Dev, Agent

Input-Schema:
- use_case
- required_metrics[]
- source_refs[]

Output-Schema:
- columns[]
- validation_rules[]
- sample_rows[]

Fehlerfälle:
- invalid_metric_definition

Nebenwirkungen:
- erzeugt Delegationsprotokoll

Verifikation:
- Typ-/Einheitenkonsistenz gegeben.

Endpoint-Referenz:
- `API.NLM.GENERATE.TABLE`

### CMD.QWEN.VISION
Zweck:
Bildanalyse und visuelle Erkennung mit Qwen 3.5 Vision.

Rolle:
Dev, Agent

Input-Schema:
- image_url oder base64
- analysis_type: product | layout | diagram | ocr
- options?: { detail_level: low|high }

Output-Schema:
- tags[]
- description
- metrics{}
- confidence_score

Fehlerfälle:
- invalid_image_format
- image_too_large
- analysis_failed

Nebenwirkungen:
- keine

Verifikation:
- Qwen gibt strukturierte Analyse zurück.

Endpoint-Referenz:
- `API.QWEN.VISION`

### CMD.QWEN.CODE
Zweck:
Full-Stack Code-Generierung mit Qwen 3.5.

Rolle:
Dev, Agent

Input-Schema:
- prompt
- language: typescript | go | python
- framework?: nextjs | supabase | generic
- context?: {}

Output-Schema:
- code
- file_path
- dependencies[]

Fehlerfälle:
- prompt_too_long
- unsupported_language
- generation_failed

Nebenwirkungen:
- keine

Verifikation:
- Code ist syntaktisch korrekt.

Endpoint-Referenz:
- `API.QWEN.CHAT`

### CMD.QWEN.CHAT
Zweck:
Natürliche Konversation mit Qwen 3.5 für Chat und Textaufgaben.

Rolle:
User, Dev, Agent

Input-Schema:
- message
- context?: {}
- temperature?: 0.0-1.0
- max_tokens?: number

Output-Schema:
- response
- usage{}
- finish_reason

Fehlerfälle:
- invalid_message
- rate_limit_exceeded
- model_unavailable

Nebenwirkungen:
- keine

Verifikation:
- Antwort ist kontextbezogen.

Endpoint-Referenz:
- `API.QWEN.CHAT`

### CMD.QWEN.OCR
Zweck:
Texterkennung aus Dokumenten und PDFs mit Qwen Vision.

Rolle:
Dev, Agent

Input-Schema:
- document_url oder base64
- language?: de|en|fr
- extract_tables?: boolean

Output-Schema:
- text
- blocks[]
- tables[]

Fehlerfälle:
- document_corrupt
- no_text_found
- extraction_failed

Nebenwirkungen:
- keine

Verifikation:
- Text korrekt extrahiert.

Endpoint-Referenz:
- `API.QWEN.OCR`

## Abnahme-Check COMMANDS
1. Jeder Command hat Schema
2. Endpoint-Referenz vorhanden
3. Fehlerfälle und Verifikation enthalten
4. NLM-Commands vollständig vorhanden
5. Qwen-Commands vollständig vorhanden

---
