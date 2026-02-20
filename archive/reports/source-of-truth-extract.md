# SOURCE OF TRUTH EXTRACT

**Generated:** 2026-02-20  
**Source Files:** 4 primary documents from ~/.config/opencode and /Users/jeremy/dev/BIOMETRICS/

---

## Key Mandates (from ~/.config/opencode/AGENTS.md)

### Top 10 Executive Rules:
1. **PARALLEL EXECUTION** - `run_in_background=false` NEVER, `run_in_background=true` ALWAYS
2. **SEARCH BEFORE CREATE** - `glob()`, `grep()` first, NEVER blind file creation
3. **VERIFY-THEN-EXECUTE** - `lsp_diagnostics`, bash checks ALWAYS
4. **GIT COMMIT DISCIPLINE** - After every significant change
5. **FREE-FIRST PHILOSOPHY** - Self-hosted, free tiers, open source
6. **RESOURCE PRESERVATION** - Delete OpenCode, configs, containers NEVER
7. **NO-SCRIPT MANDATE** - Manual bash scripts NEVER, AI agents ALWAYS
8. **NLM DUPLICATE PREVENTION** - `nlm source list` before upload, `nlm source delete` old versions
9. **TODO DISCIPLINE** - Create todos for multi-step tasks
10. **PERFORMANCE FIRST** - Native CDP over Playwright

### Critical Additional Mandates:
- **DEQLHI-LOOP** - Nach jeder Task: 5 neue Tasks hinzufügen, Git Commit, Docs updaten
- **PORT SOVEREIGNTY** - Keine Standard-Ports (3000, 5432, etc.), 50000-59999 range
- **NO-TIMEOUT POLICY** - NIEMALS timeout in opencode.json eintragen (Qwen braucht 70-90s)

---

## Agent Model Assignment

| Model | Provider | Categories | Max Parallel | Use Case |
|-------|----------|------------|--------------|----------|
| `qwen/qwen3.5-397b-a17b` | NVIDIA NIM | build, visual-engineering, writing, general | **1** | Haupt-Code, komplexe Tasks |
| `opencode/kimi-k2.5-free` | OpenCode ZEN | deep | **1** | Heavy Lifting, Setup |
| `opencode/minimax-m2.5-free` | OpenCode ZEN | quick, explore, librarian | **1** | Triviale Tasks, Code Discovery |

**MAXIMUM 3 AGENTS PARALLEL** (je 1 pro Modell)

---

## Forbidden Actions

- ❌ `run_in_background=false` verwenden
- ❌ Dateien erstellen ohne `glob()`/`grep()` vorher
- ❌ Vertrauen ohne Verifikation (`lsp_diagnostics`)
- ❌ Git Commit überspringen
- ❌ OpenCode/Configs/Container löschen
- ❌ Manuelle bash-Scripts schreiben (statt AI agents)
- ❌ NLM source add ohne vorher `nlm source list`
- ❌ Timeout in opencode.json eintragen
- ❌ Standard-Ports verwenden (3000, 5432, 8080, etc.)
- ❌ "Fertig" sagen ohne Evidenz (Tests, Screenshots)

---

## Required Actions

- ✅ IMMER `run_in_background=true` bei delegate_task()
- ✅ IMMER zuerst lesen mit `glob()`/`grep()`/`read()`
- ✅ IMMER verifizieren mit `lsp_diagnostics`, bash checks
- ✅ IMMER Git Commit nach jeder Änderung
- ✅ IMMER TodoWrite mit 5 neuen Tasks nach Abschluss
- ✅ IMMER lastchanges.md und AGENTS.md aktualisieren
- ✅ IMMER visuelle Prüfung (Screenshots, Browser-Checks)
- ✅ IMMER Best Practices 2026 (CEO-Elite Niveau)
- ✅ IMMER Problem → SOFORT Internet-Recherche
- ✅ IMMER Subagentenmassive Prompts geben

---

## File Reading Requirements

### Before Every Task:
1. **~/.config/opencode/AGENTS.md** - Globale Agenten-Regeln (Source of Truth)
2. **docs/agents/AGENTS.md** - Lokale Projekt-Regeln
3. **docs/ORCHESTRATOR-MANDATE.md** - Orchestrator-spezifische Regeln
4. **docs/architecture/ARCHITECTURE.md** - Projekt-Architektur
5. **CHANGELOG.md** - Letzte Änderungen
6. **SETUP-CHECKLISTE.md** - Setup-Status

### Pflichtlektüre für Subagenten (aus docs/agents/AGENTS.md):
- JEDER Subagent MUSS zuerst AGENTS.md und ARCHITECTURE.md komplett lesen (bis zur letzten Zeile!)
- Mit `read()` Tool die Dateien bis zur letzten Zeile lesen

---

## Duplicates Found

### Keine Duplikate:
- `docs/agents/AGENTS.md` (263 Zeilen) ist **KEIN Duplikat** von ~/.config/opencode/AGENTS.md
- Es ist eine **lokale Ergänzung** mit projekt-spezifischen Regeln
- **docs/best-practices/AGENTS.md existiert NICHT** (nur docs/best-practices/ mit anderen MD-Dateien)

### Verwandte Dateien:
| Datei | Größe | Zweck |
|-------|-------|-------|
| `docs/agents/AGENTS-MANDATES.md` | 310KB | Original AGENTS-Mandates (外部) |
| `docs/agents/AGENTS-GLOBAL.md` | 134KB | Global governance |
| `docs/ORCHESTRATOR-MANDATE.md` | 307 Zeilen | Orchestrator-spezifisch |
| `docs/UNIVERSAL-BLUEPRINT.md` | 524 Zeilen | Setup-Anleitung |

---

## Recommendations for rules/global/AGENTS.md

### Übernehmen:
1. **DEQLHI-LOOP** - Infinite Work Mode mit 5-Tasks-Regel
2. **Parallel Execution** - MAX 3 Agents (1 pro Modell)
3. **Search Before Create** - glob/grep IMMER zuerst
4. **File Reading Protocol** - Subagenten müssen AGENTS.md + ARCHITECTURE.md lesen
5. **Model Assignment Tabelle** - Klare Zuweisung
6. **NLM CLI Commands** - Duplicate Prevention
7. **No-Timeout Policy** - Qwen 3.5 braucht 70-90s

### Anpassen:
1. **Projekt-spezifische Pfade** - BIOMETRICS-Pfade eintragen
2. **Stack-Policy** - Next.js + Go + Supabase
3. **OH-MY-OPENCODE Categories** - An Projekt anpassen

### Neu schreiben:
1. **Kurze, fokussierte Version** - Nicht 5000+ Zeilen
2. **BIOMETRICS-spezifische Beispiele** - Statt generischer Regeln
3. **Integration mit lokalen Dateien** - Auf CHANGELOG.md, MEETING.md verweisen

---

## Summary

Die wichtigsten Dateien sind:
1. **~/.config/opencode/AGENTS.md** - Source of Truth (33+ Mandate)
2. **docs/agents/AGENTS.md** - Lokale Ergänzung (263 Zeilen)
3. **docs/ORCHESTRATOR-MANDATE.md** - Orchestrator-Regeln (307 Zeilen)
4. **docs/UNIVERSAL-BLUEPRINT.md** - Setup-Guide (524 Zeilen)

**Empfehlung:** Eine kurze rules/global/AGENTS.md erstellen (ca. 100-150 Zeilen) mit den wichtigsten Regeln und Verweisen auf die vollständigen Dokumente.
