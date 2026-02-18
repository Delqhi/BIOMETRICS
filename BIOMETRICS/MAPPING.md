# MAPPING.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Mapping-Pflege ist mandatory und gilt als Teil der Definition of Done.
- Jede strukturelle Änderung muss hier referenziert und geprüft werden.
- Drift zwischen Domänen wird sofort als Risiko behandelt.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Zentrale Mapping-Governance für Konsistenz zwischen Produkt, Frontend, Backend, Datenbank, Commands, Endpoints und NLM-Assets.

## Pflichtregel
Jede relevante Änderung wird gegen die Mapping-Dateien geprüft. Ohne Mapping-Konsistenz ist eine Änderung nicht done.

## Mapping-Module
1. `MAPPING-COMMANDS-ENDPOINTS.md`
2. `MAPPING-FRONTEND-BACKEND.md`
3. `MAPPING-DB-API.md`
4. `MAPPING-NLM-ASSETS.md`

## Prüfzyklus
1. Änderung identifizieren
2. betroffenes Mapping-Modul aktualisieren
3. Konsistenzcheck durchführen
4. Ergebnisse in `MEETING.md` und `CHANGELOG.md` dokumentieren

## Abnahme-Check MAPPING
1. Alle Mapping-Module vorhanden
2. Zuordnungstabellen ausgefüllt
3. Offene Deltas explizit markiert
4. Prüfstatus dokumentiert

## Qwen 3.5 Integration Mapping

### Command-zu-Endpoint Mapping

| Command | Endpoint | Method | Auth |
|---------|----------|--------|------|
| CMD.QWEN.VISION | API.QWEN.VISION | POST | NVIDIA_API_KEY |
| CMD.QWEN.CODE | API.QWEN.CHAT | POST | NVIDIA_API_KEY |
| CMD.QWEN.CHAT | API.QWEN.CHAT | POST | NVIDIA_API_KEY |
| CMD.QWEN.OCR | API.QWEN.OCR | POST | NVIDIA_API_KEY |

### Skill-zu-Endpoint Mapping

| Skill | Endpoint | Use Case |
|-------|----------|----------|
| qwen_vision_analysis | API.QWEN.VISION | Bildanalyse |
| qwen_code_generation | API.QWEN.CHAT | Code-Generierung |
| qwen_document_ocr | API.QWEN.OCR | Texterkennung |
| qwen_video_understanding | API.QWEN.VISION | Video-Analyse |
| qwen_conversation | API.QWEN.CHAT | Chat-Interaktion |

### Provider Configuration

| Provider | Endpoint | Model ID | Context | Output |
|----------|----------|----------|---------|--------|
| NVIDIA NIM | https://integrate.api.nvidia.com/v1 | qwen/qwen3.5-397b-a17b | 262K | 32K |

### Environment Variables

| Variable | Zweck | Required |
|----------|-------|----------|
| NVIDIA_API_KEY | Auth für NIM Endpoint | Ja |
| QWEN_MODEL_ID | Modell-Identifier | Nein (Default: qwen/qwen3.5-397b-a17b) |
| QWEN_BASE_URL | Custom Endpoint | Nein |

### Deployment Mapping

| Komponente | Platform | Status |
|------------|----------|--------|
| Edge Functions | Vercel | Required |
| NVIDIA NIM | Cloud (NVIDIA) | Required |
| API Gateway | Vercel Edge | Required |

---
