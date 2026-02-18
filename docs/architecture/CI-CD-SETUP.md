# CI-CD-SETUP.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Pipeline-Gates müssen Policy-as-Code und Doku-Sync erzwingen.
- Sicherheits-, Mapping- und Qualitätschecks sind Default-Blocker bei Verstößen.
- Deployment-Freigaben erfolgen nur mit Nachweis kompletter Kontrollkette.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Pipeline-Standard für Qualitätssicherung und kontrollierte Auslieferung.

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

## Pipeline-Prinzipien
1. Jeder Merge ist überprüfbar
2. Qualität vor Geschwindigkeit
3. Rollback muss jederzeit möglich sein
4. Artefakte sind nachvollziehbar versioniert

## Standard-Stages
1. Install
2. Lint
3. Typecheck
4. Tests
5. Build
6. Security Checks
7. Packaging
8. Deploy
9. Post-Deploy Verification

## Gate-Definitionen
- Gate 1: `pnpm lint` grün
- Gate 2: `pnpm typecheck` grün
- Gate 3: `pnpm test` grün
- Gate 4: `pnpm build` grün
- Gate 5: `go test ./...` grün (falls Backend vorhanden)
- Gate 6: Security-Check ohne kritische Findings
- Gate 7: Qwen 3.5 Integration Tests grün

## Qwen 3.5 Testing Pipeline
- **Model:** qwen/qwen3.5-397b-a17b (NVIDIA NIM)
- **Context:** 262K tokens
- **Output:** 32K tokens
- **Timeout:** 120000ms (2 Minuten)

### Qwen Test Stages
1. **Vision Analysis Test:** Bildanalyse-Qualität prüfen
2. **Code Generation Test:** Syntax und Funktionalität des generierten Codes
3. **OCR Test:** Texterkennungsgenauigkeit
4. **Conversation Test:** Antwortkonsistenz und Kontextverständnis
5. **Video Understanding Test:** Szenenbeschreibung-Genauigkeit

### Qwen Test Configuration
```yaml
qwen_tests:
  enabled: true
  model: qwen/qwen3.5-397b-a17b
  timeout: 120000
  retry_attempts: 3
  fallback_model: kimi-for-coding/k2p5
  test_cases:
    - vision_analysis_quality > 0.85
    - code_generation_syntax_valid: true
    - ocr_accuracy > 0.90
    - conversation_context_retention: true
```

### Qwen Performance Benchmarks
- Vision Analysis: < 5s pro Bild
- Code Generation: < 30s pro Komponente
- OCR: < 3s pro Seite
- Conversation: < 2s pro Antwort

## Branch-/Merge-Regeln
- PR Pflicht für geschützte Branches
- Keine direkte Änderung auf protected branch
- Merge nur bei grünen Gates

## Deploy-Strategie (Template)
- Staging: automatisch nach Merge in Integrationsbranch
- Production: kontrollierter Trigger mit Freigabe
- Rollback: definierter Fallback auf letzte stabile Version

## Artefakt-Management
- Build-Artefakte versionieren
- Prüfsummen und Metadaten speichern
- Release-Referenz im Changelog dokumentieren

## Post-Deploy Checks
1. Health Endpoints
2. Kernjourneys
3. Fehlerquote
4. Performance-Basiswerte

## Incident-Integration
- Bei P0/P1 Fehlern automatischer Eskalationspfad
- Postmortem-Pflicht mit Follow-up Tasks

## Abnahme-Check CI-CD
1. Stages und Gates vollständig
2. Merge-Regeln dokumentiert
3. Rollbackprozess vorhanden
4. Post-Deploy Checks definiert

---
