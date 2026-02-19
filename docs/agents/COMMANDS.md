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

---

## 🧠 NLM CLI COMMANDS

```bash
# Create notebook
nlm notebook create "Title"

# List sources
nlm source list <notebook-id>

# Delete old source (BEFORE adding new!)
nlm source delete <source-id> -y

# Add new source
nlm source add <notebook-id> --file "file.md" --wait
```

**⚠️ DUPLICATE PREVENTION:** ALWAYS run `nlm source list` before `nlm source add`!

---

## 🔍 BIOMETRICS CLI COMMANDS

### CMD.BIOMETRICS.CHECK

**Zweck:**
Prüft BIOMETRICS Repository auf Konsistenz und Vollständigkeit.

**Rolle:**
User, Dev, Admin, Agent

**Input-Schema:**
- Keine Parameter (automatische Prüfung)

**Output-Schema:**
- status_summary (Bestanden/Durchgefallen)
- directory_checks (global, local, biometrics-cli, docs)
- config_checks (oh-my-opencode.json)
- model_consistency (qwen/qwen3.5-397b-a17b)
- readme_checks (README.md in allen Verzeichnissen)
- agent_mapping (AGENT-MODEL-MAPPING.md)
- loop_config (∞Best∞Practices∞Loop.md)
- error_count (Anzahl gefundener Fehler)

**Fehlerfälle:**
- directory_missing (Verzeichnis fehlt)
- config_missing (Konfiguration fehlt)
- model_inconsistent (Modell-Namen inkonsistent)
- readme_missing (README fehlt)
- mapping_missing (Agent-Mapping fehlt)
- loop_config_missing (Loop-Konfiguration fehlt)

**Nebenwirkungen:**
- Keine (read-only Prüfung)

**Verifikation:**
- Alle Verzeichnisse existieren
- Alle Konfigurationen vorhanden
- Modell-Namen konsistent
- READMEs in Hauptverzeichnissen
- Agent-Mapping dokumentiert
- Loop-Konfiguration korrekt

**Endpoint-Referenz:**
- `biometrics-cli/bin/biometrics-check` (Bash Script)

**Usage:**
```bash
cd /Users/jeremy/dev/BIOMETRICS
./biometrics-cli/bin/biometrics-check
```

**Example Output:**
```
🔍 BIOMETRICS REPO CHECK
========================

📁 Step 1: Hauptverzeichnisse
─────────────────────────────
  ✓ global/ existiert
  ✓ local/ existiert
  ✓ biometrics-cli/ existiert
  ✓ docs/ existiert

⚙️  Step 3: Konfigurationsdateien
──────────────────────────────────
  ✓ oh-my-opencode.json existiert
  📊 Agents:
    - cosmos-video-edit
    - cosmos-video-gen
    - flux1-image
    ...

✅ CHECK COMPLETE
🎉 ALLE CHECKS BESTANDEN!
```

---

---

## 🔄 DEQLHI-LOOP (INFINITE WORK MODE)

- After each completed task → Add 5 new tasks immediately
- Never "done" - only "next task"
- Always document → Every change in files
- Git commit + push after EVERY change
- Parallel execution ALWAYS (run_in_background=true)

### Loop Mechanism:
1. Task N Complete
2. Git Commit + Push
3. Update Docs
4. Add 5 New Tasks
5. Next Task N+1
6. Repeat infinitely

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
