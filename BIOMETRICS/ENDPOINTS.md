# ENDPOINTS.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Endpoint-Verträge folgen globalen API-, Security- und Versionierungsregeln.
- Jede Änderung benötigt Mapping-Update und Fehlerformat-Konsistenz.
- Betriebsmetriken und Incident-Readiness sind mitzupflegen.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller API-Katalog als Gegenstück zu `COMMANDS.md`.

## Endpoint-Schema (Pflicht)
```text
Endpoint-ID:
Path:
Method:
Auth:
Rollen:
Request-Schema:
Response-Schema:
Fehlercodes:
Idempotenz:
Rate-Limit:
Command-Referenz:
Verifikation:
```

## Universeller Basis-Katalog

### API.SYSTEM.STATUS
Path:
`/api/system/status`

Method:
GET

Auth:
required

Rollen:
user, admin, agent

Request-Schema:
- scope?: all | docs | runtime

Response-Schema:
- status_summary
- open_blockers[]
- next_actions[]

Fehlercodes:
- 400 invalid_scope
- 401 unauthorized
- 500 internal_error

Idempotenz:
ja

Rate-Limit:
60/min

Command-Referenz:
- `CMD.SYSTEM.STATUS`

Verifikation:
- liefert konsistente Statusdaten.

### API.TASKS.LIST
Path:
`/api/tasks/list`

Method:
GET

Auth:
required

Rollen:
user, dev, admin, agent

Request-Schema:
- priority?: P0|P1|P2|all
- status?: OPEN|IN_PROGRESS|BLOCKED|DONE|all

Response-Schema:
- tasks[]

Fehlercodes:
- 400 invalid_filter
- 401 unauthorized
- 500 internal_error

Idempotenz:
ja

Rate-Limit:
120/min

Command-Referenz:
- `CMD.TASKS.LIST`

Verifikation:
- Filter und Pagination korrekt.

### API.TASKS.UPDATE
Path:
`/api/tasks/update`

Method:
POST

Auth:
required

Rollen:
dev, admin, agent

Request-Schema:
- task_id
- status
- evidence_ref

Response-Schema:
- task
- audit_entry_id

Fehlercodes:
- 400 invalid_payload
- 401 unauthorized
- 403 forbidden
- 404 task_not_found
- 409 conflict
- 500 internal_error

Idempotenz:
nein

Rate-Limit:
60/min

Command-Referenz:
- `CMD.TASKS.UPDATE`

Verifikation:
- Statuswechsel auditierbar.

### API.NLM.GENERATE.VIDEO
Path:
`/api/nlm/generate/video`

Method:
POST

Auth:
required

Rollen:
dev, agent

Request-Schema:
- topic
- audience
- business_goal
- source_refs[]

Response-Schema:
- titles[]
- scripts
- storyboard
- cta[]
- quality_score

Fehlercodes:
- 400 missing_sources
- 401 unauthorized
- 422 quality_below_threshold
- 500 internal_error

Idempotenz:
nein

Rate-Limit:
30/min

Command-Referenz:
- `CMD.NLM.GENERATE.VIDEO`

Verifikation:
- Score >= 13/16 und Korrektheit = 2.

### API.NLM.GENERATE.INFOGRAPHIC
Path:
`/api/nlm/generate/infographic`

Method:
POST

Auth:
required

Rollen:
dev, agent

Request-Schema:
- topic
- key_points[]
- source_refs[]

Response-Schema:
- layout_blueprint
- visual_mapping
- accessibility_notes
- quality_score

Fehlercodes:
- 400 invalid_input
- 401 unauthorized
- 422 inconsistent_claims
- 500 internal_error

Idempotenz:
nein

Rate-Limit:
30/min

Command-Referenz:
- `CMD.NLM.GENERATE.INFOGRAPHIC`

Verifikation:
- Konsistenz und Lesbarkeit geprüft.

### API.NLM.GENERATE.PRESENTATION
Path:
`/api/nlm/generate/presentation`

Method:
POST

Auth:
required

Rollen:
dev, agent

Request-Schema:
- occasion
- audience_type
- decision_goal
- source_refs[]

Response-Schema:
- slide_outline[]
- speaker_notes[]
- risks_and_tradeoffs
- quality_score

Fehlercodes:
- 400 invalid_input
- 401 unauthorized
- 422 weak_storyline
- 500 internal_error

Idempotenz:
nein

Rate-Limit:
30/min

Command-Referenz:
- `CMD.NLM.GENERATE.PRESENTATION`

Verifikation:
- Entscheidungspfad klar nachvollziehbar.

### API.NLM.GENERATE.TABLE
Path:
`/api/nlm/generate/table`

Method:
POST

Auth:
required

Rollen:
dev, agent

Request-Schema:
- use_case
- required_metrics[]
- source_refs[]

Response-Schema:
- columns[]
- validation_rules[]
- sample_rows[]
- quality_score

Fehlercodes:
- 400 invalid_metric_definition
- 401 unauthorized
- 422 quality_below_threshold
- 500 internal_error

Idempotenz:
nein

Rate-Limit:
30/min

Command-Referenz:
- `CMD.NLM.GENERATE.TABLE`

Verifikation:
- Typ-/Einheiten-/Zeitbezug vollständig.

### API.QWEN.VISION
Path:
`/api/qwen/vision`

Method:
POST

Auth:
required (NVIDIA_API_KEY)

Rollen:
dev, agent

Request-Schema:
- image: string (base64 oder URL)
- analysis_type: product | layout | diagram | ocr
- options?: { detail_level: low|high }

Response-Schema:
- tags: string[]
- description: string
- metrics: object
- confidence_score: number

Fehlercodes:
- 400 invalid_image_format
- 413 image_too_large
- 401 unauthorized
- 500 analysis_failed
- 503 model_unavailable

Idempotenz:
nein

Rate-Limit:
30/min

Command-Referenz:
- `CMD.QWEN.VISION`

Verifikation:
- Qwen gibt strukturierte Vision-Analyse zurück.

### API.QWEN.CHAT
Path:
`/api/qwen/chat`

Method:
POST

Auth:
required (NVIDIA_API_KEY)

Rollen:
user, dev, agent

Request-Schema:
- message: string
- context?: object
- temperature?: number (0.0-1.0)
- max_tokens?: number

Response-Schema:
- response: string
- usage: { prompt_tokens, completion_tokens, total_tokens }
- finish_reason: stop | length

Fehlercodes:
- 400 invalid_message
- 401 unauthorized
- 429 rate_limit_exceeded
- 503 model_unavailable
- 500 internal_error

Idempotenz:
nein

Rate-Limit:
40/min (NVIDIA NIM Limit)

Command-Referenz:
- `CMD.QWEN.CHAT`
- `CMD.QWEN.CODE`

Verifikation:
- Kontextbezogene Antwort generiert.

### API.QWEN.OCR
Path:
`/api/qwen/ocr`

Method:
POST

Auth:
required (NVIDIA_API_KEY)

Rollen:
dev, agent

Request-Schema:
- document: string (base64 oder URL)
- language?: string
- extract_tables?: boolean

Response-Schema:
- text: string
- blocks: { bbox, text, confidence }[]
- tables: { rows, columns, data }[]

Fehlercodes:
- 400 document_corrupt
- 401 unauthorized
- 404 no_text_found
- 500 extraction_failed

Idempotenz:
nein

Rate-Limit:
30/min

Command-Referenz:
- `CMD.QWEN.OCR`

Verifikation:
- Text korrekt aus Dokument extrahiert.

## Abnahme-Check ENDPOINTS
1. Jeder Endpoint hat Command-Referenz
2. Auth und Rollen sind je Endpoint definiert
3. Fehlercodes enthalten
4. Verifikation pro Endpoint enthalten
5. Qwen-Endpoints vollständig vorhanden

---
