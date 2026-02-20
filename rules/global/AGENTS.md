# GLOBAL AGENTS RULES

**Version:** 1.0 "BIOMETRICS Foundation"  
**Date:** 2026-02-20  
**Status:** ✅ ACTIVE - MUST READ BEFORE EVERY SESSION  
**Source:** ~/.config/opencode/AGENTS.md (V20.0, 3100+ lines, 33 mandates)  

---

## 🎯 PURPOSE

Dieses Dokument definiert die **globalen Regeln für ALLE KI-Agenten** die im BIOMETRICS-Ökosystem arbeiten. Es ist die **Source of Truth** für Agenten-Verhalten, Modell-Zuweisung, und Arbeitsweise.

**Jeder Agent MUSS diese Regeln lesen und befolgen BEVOR er arbeitet.**

---

## 🔑 TOP 10 EXECUTIVE RULES

### 1. **PARALLEL EXECUTION MANDATE**
```typescript
// ❌ FALSCH
task(run_in_background=false, prompt="...") // NIE!

// ✅ RICHTIG
task(run_in_background=true, prompt="...") // IMMER!
```

**Warum:** Agents dürfen NIEMALS aufeinander warten. Parallele Ausführung ist PFLICHT.

### 2. **SEARCH BEFORE CREATE**
```typescript
// ❌ FALSCH
write({ filePath: "new.md", content: "..." }) // Blind erstellen!

// ✅ RICHTIG
glob("**/*.md")           // Erst suchen
read("existing.md")       // Verstehen
write({ filePath: "new.md", content: "..." }) // Dann erstellen
```

**Warum:** Duplikate vermeiden, existierende Strukturen wiederverwenden.

### 3. **VERIFY-THEN-EXECUTE**
```typescript
// ❌ FALSCH
write(...) // Ohne Prüfung

// ✅ RICHTIG
write(...)
lsp_diagnostics(filePath)  // Prüfen
bash("go test ./...")      // Testen
```

**Warum:** Fehler sofort erkennen, nicht erst später.

### 4. **GIT COMMIT DISCIPLINE**
```bash
# Nach JEDER signifikanten Änderung:
git add -A
git commit -m "type: description"
git push origin main
```

**Warum:** Jede Änderung ist gesichert, Rollback möglich.

### 5. **FREE-FIRST PHILOSOPHY**
- ✅ Self-hosted Lösungen bevorzugen
- ✅ Free Tiers nutzen (NVIDIA NIM, OpenCode ZEN)
- ✅ Open Source vor kommerziellen Tools

### 6. **RESOURCE PRESERVATION**
- ❌ NIEMALS OpenCode neu installieren
- ❌ NIEMALS ~/.config/opencode/ löschen
- ❌ NIEMALS Container/Docker-Configs löschen

**Warum:** Konfigurationen sind wertvoll, Verlust = Katastrophe.

### 7. **NO-SCRIPT MANDATE**
- ❌ KEINE manuellen bash-scripts schreiben
- ✅ Agents für ALLES nutzen

**Warum:** Agents sind flexibler, wartbarer.

### 8. **NLM DUPLICATE PREVENTION**
```bash
# VOR jedem Upload:
nlm source list <notebook-id>     # Prüfen
nlm source delete <old-id> -y     # Altes löschen
nlm source add <notebook-id> --file "new.md" --wait  # Neues hinzufügen
```

**Warum:** Duplikate verwirren das NLM.

### 9. **TODO DISCIPLINE**
```typescript
// Bei MULTIPLE Steps (2+):
todowrite([
  { id: "task-1", content: "Step 1", status: "pending" },
  { id: "task-2", content: "Step 2", status: "pending" },
])

// Vor JEDEM Step:
todowrite([{ id: "task-1", status: "in_progress" }])

// Nach JEDEM Step:
todowrite([{ id: "task-1", status: "completed" }])
```

**Warum:** User sieht Fortschritt, Tasks werden nicht vergessen.

### 10. **PERFORMANCE FIRST**
- ✅ Native CDP über Playwright (46x schneller!)
- ✅ Ultra-fast native Workers
- ✅ Connection Pooling (5-10 parallele Connections)

---

## 🚨 CRITICAL MANDATES

### DEQLHI-LOOP (INFINITE WORK MODE)

**PRINZIP:** Nach JEDER abgeschlossenen Task → SOFORT 5 neue Tasks hinzufügen.

```
START
  │
  ▼
┌─────────────────┐
│ Task N Complete │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Git Commit      │ ← JEDE ÄNDERUNG COMMITTEN + PUSHEN
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Update Docs     │ ← ARCHITECTURE.md + AGENTS.md
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Add 5 New Tasks │ ← IMMER 5 NEUE TASKS!
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Next Task N+1   │ ← SOFORT WEITER, KEINE PAUSE!
└─────────────────┘
         │
         └──────────────┐
                        │
                        ▼
                 (Loop continues)
```

**Warum:** Produktivität ist unendlich skalierbar. Kein "Fertig" - nur "Nächster Task".

### PORT SOVEREIGNTY (NO STANDARD PORTS)

**VERBOTENE PORTS:**
- ❌ 3000, 8080, 5432, 6379, 5678, 8000, 9000, 3306, 27017, 9200, 80, 443

**ERLAUBTE PORTS:**
- ✅ 50000-59999 range (unique ports)

**Naming Convention:**
```
{CATEGORY}-{NUMBER}-{NAME}
- agent-XX-    → AI Workers, Orchestrators
- room-XX-     → Infrastructure, Databases
- solver-X.X-  → Money-Making Workers
- builder-X-   → Content Creation
```

**Beispiele:**
- ✅ `agent-01-n8n-manager:8001`
- ✅ `room-03-postgres-master:8103`
- ❌ `sin-zimmer-01-n8n:5678` (falsches Prefix + Standard-Port!)

---

## 📊 AGENT MODEL ASSIGNMENT

### Verfügbare Modelle

| Modell | Provider | Max Parallel | Use Case | Latency |
|--------|----------|-------------|----------|---------|
| **qwen/qwen3.5-397b-a17b** | NVIDIA NIM | **1** | Code, Docs, Architecture | 70-90s |
| **opencode/kimi-k2.5-free** | OpenCode ZEN | **1** | Deep Analysis, Heavy Lifting | 10-20s |
| **opencode/minimax-m2.5-free** | OpenCode ZEN | **1** | Quick Tasks, Configs, MD | 5-10s |

### KRITISCHE REGELN

1. **NIEMALS 2 Agents mit gleichem Modell parallel!**
   - Qwen 3.5: MAX 1 Agent
   - Kimi K2.5: MAX 1 Agent
   - MiniMax M2.5: MAX 1 Agent
   - **MAXIMAL 3 Agents parallel (je 1 pro Modell)**

2. **Workflow:**
   ```typescript
   // ✅ KORREKT (3 verschiedene Modelle):
   task(category="visual-engineering", prompt="...")     // Qwen 3.5
   task(category="deep", model="opencode/kimi-k2.5-free", prompt="...")  // Kimi
   task(category="quick", model="opencode/minimax-m2.5-free", prompt="...") // MiniMax
   
   // ❌ FALSCH (alle gleiches Modell):
   task(category="visual-engineering", prompt="...")  // Qwen 3.5
   task(category="visual-engineering", prompt="...")  // Qwen 3.5 - BLOCKED!
   ```

### Model Selection Guide

| Task Type | Model | Why |
|-----------|-------|-----|
| **Code Implementation** | Qwen 3.5 397B | Beste Code-Qualität |
| **Documentation (MD)** | MiniMax M2.5 | Schnell, 10x parallel |
| **Architecture Design** | Qwen 3.5 397B | Komplexes Denken |
| **Research / Search** | MiniMax M2.5 | Schnell, effizient |
| **Deep Analysis** | Kimi K2.5 | 1M Context Window |
| **Quick Tasks** | MiniMax M2.5 | <10s Latency |
| **Planning** | Qwen 3.5 397B | Strategisches Denken |

---

## 📖 FILE READING REQUIREMENTS

### PFLICHT-DATEIEN (BEVOR DU ARBEITEST)

**JEDER Subagent MUSS vor Arbeitsbeginn lesen:**

1. **~/.config/opencode/AGENTS.md** (diese Datei)
   - Globale Regeln
   - Modell-Zuweisung
   - Mandates

2. **BIOMETRICS/ARCHITECTURE.md**
   - Projekt-Architektur
   - Verzeichnis-Struktur
   - Migration-Plan

3. **BIOMETRICS/AGENTS-PLAN.md** (falls vorhanden)
   - Aktueller Plan
   - Offene Tasks
   - Blocker

4. **BIOMETRICS/lastchanges.md** (falls vorhanden)
   - Was wurde zuletzt gemacht
   - Nächste Schritte

### WIE LESEN

```typescript
// ❌ FALSCH
read("file.md", { limit: 50 })  // Nur 50 Zeilen!

// ✅ RICHTIG
read("file.md")  // KOMPLETT (bis letzte Zeile!)
```

**Warum:** Oberflächliches Lesen führt zu Fehlern. Agents MÜSSEN den kompletten Kontext haben.

---

## 🚫 FORBIDDEN ACTIONS

### NIEMALS TUN

1. ❌ **NIEMALS `run_in_background=false`**
   - Agents dürfen NICHT sequentiell arbeiten

2. ❌ **NIEMALS Dateien erstellen ohne `glob()` oder `ls`**
   - IMMER zuerst prüfen ob Datei existiert

3. ❌ **NIEMALS "fertig" sagen ohne Evidenz**
   - IMMER Dateiinhalt zeigen
   - IMMER Tests machen
   - IMMER `lsp_diagnostics` prüfen

4. ❌ **NIEMALS User-Onboarding überspringen**
   - IMMER mit User Config erstellen
   - IMMER API Keys erklären
   - IMMER gemeinsam testen

5. ❌ **NIEMALS OpenCode neu installieren**
   - Reparatur vor Neuinstallation
   - Configs niemals löschen

6. ❌ **NIEMALS Standard-Ports verwenden**
   - Immer 50000-59999 range
   - Container Naming Convention beachten

7. ❌ **NIEMALS Secrets im Code**
   - Immer Environment Variables
   - `.gitignore` für .env Dateien

8. ❌ **NIEMALS Type Errors suppressen**
   - Kein `as any`, `@ts-ignore`, `@ts-expect-error`

---

## ✅ REQUIRED ACTIONS

### IMMER TUN

1. ✅ **IMMER `glob()` oder `grep()` vor Datei-Erstellung**
   ```typescript
   glob("**/*.md")  // Existiert Datei schon?
   ```

2. ✅ **IMMER `lsp_diagnostics` nach Datei-Änderung**
   ```typescript
   write(...)
   lsp_diagnostics(filePath)  // Fehler prüfen
   ```

3. ✅ **IMMER `git commit` nach Änderung**
   ```bash
   git add -A && git commit -m "type: description" && git push
   ```

4. ✅ **IMMER 5 neue Tasks nach Completion**
   ```typescript
   todowrite([...newTasks])  // 5 neue Tasks hinzufügen
   ```

5. ✅ **IMMER in `lastchanges.md` dokumentieren**
   ```markdown
   ## [YYYY-MM-DD HH:MM] - [AGENT/TASK-ID]
   **Beobachtungen:** [...]
   **Fehler:** [...]
   **Lösungen:** [...]
   **Nächste Schritte:** [...]
   ```

6. ✅ **IMMER visuell prüfen**
   - Screenshots machen
   - Browser-Checks
   - CDP Logs prüfen

7. ✅ **IMMER Best Practices 2026 nutzen**
   - Native CDP über Playwright
   - ripgrep über grep
   - fd über find
   - sd über sed

---

## 🔧 OPENCODE CONFIGURATION

### NVIDIA NIM Setup

```json
{
  "provider": {
    "nvidia-nim": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "NVIDIA NIM (Qwen 3.5)",
      "options": {
        "baseURL": "https://integrate.api.nvidia.com/v1",
        "timeout": 120000  // PFLICHT! (70-90s Latency)
      },
      "models": {
        "qwen-3.5-397b": {
          "id": "qwen/qwen3.5-397b-a17b",
          "limit": {
            "context": 262144,
            "output": 32768
          }
        }
      }
    }
  }
}
```

### WICHTIG

- **Timeout:** 120000ms (120s) - Qwen 3.5 braucht 70-90s!
- **API Field:** `"api": "openai-completions"` (OpenClaw)
- **Rate Limit:** 40 RPM (Free Tier)
- **HTTP 429:** 60 Sekunden warten + Fallbacks

### Quick Commands

```bash
# Testen
openclaw models | grep nvidia
opencode models | grep nvidia

# Gateway
openclaw gateway restart
openclaw doctor

# API Test
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
  https://integrate.api.nvidia.com/v1/models
```

---

## 📝 DOCUMENTATION RULES

### 500-LINE MANDATE

**JEDE Rule-Datei, jeder Guide, jedes Template:**
- ✅ MINDESTENS 500 Zeilen
- ✅ Vollständig, nicht oberflächlich
- ✅ Mit Beispielen, Use Cases, Anti-Patterns

**Warum:** Oberflächliche Dokumentation führt zu Fehlern.

### TRINITY DOCUMENTATION STANDARD

**JEDES Projekt MUSS haben:**

```
project/
├── docs/
│   ├── non-dev/    # Für User (Guides, Tutorials, FAQs)
│   ├── dev/        # Für Coders (API Ref, Architecture)
│   ├── project/    # Für Team (Deployment, Roadmap)
│   └── postman/    # API Collections
├── DOCS.md         # Index & Standards
└── README.md       # Gateway (Document360)
```

### COMMENTING RULES

**Comments NUR wenn:**
- ✅ Complex algorithms
- ✅ Security-related code
- ✅ Performance optimization
- ✅ Regex patterns
- ✅ Mathematical formulas

**KEINE Comments für:**
- ❌ Selbst-erklärenden Code
- ❌ Simple variables
- ❌ Offensichtliche Logik

---

## 🎯 SUCCESS CRITERIA

### Agent Behavior ✅

- [ ] Liest ALLE Pflicht-Dateien vor Arbeit
- [ ] Nutzt korrektes Modell für Task
- [ ] Arbeitet parallel (run_in_background=true)
- [ ] Sucht vor Erstellen (glob/grep)
- [ ] Verifiziert nach Änderung (lsp_diagnostics)
- [ ] Commit nach Änderung (git add/commit/push)
- [ ] Dokumentiert (lastchanges.md)
- [ ] Fügt 5 neue Tasks nach Completion hinzu

### Code Quality ✅

- [ ] Keine Type Errors
- [ ] Keine Secrets im Code
- [ ] Error Handling vorhanden
- [ ] Tests geschrieben
- [ ] Dokumentation aktuell

---

## 🔗 REFERENCES

- **Source of Truth:** `~/.config/opencode/AGENTS.md` (3100+ lines, 33 mandates)
- **Architecture:** `BIOMETRICS/ARCHITECTURE.md`
- **Audit Report:** `BIOMETRICS/audit-report.md`
- **Structure Analysis:** `BIOMETRICS/structure-analysis.md`
- **Rearchitecture Plan:** `BIOMETRICS/BIOMETRICS-REARCHITECTURE-PLAN.md`

---

**LAST UPDATED:** 2026-02-20  
**NEXT REVIEW:** After Phase 2 completion  
**OWNER:** BIOMETRICS Core Team  

---

*Dieses Dokument ist das "Gesetzbuch" für ALLE Agents im BIOMETRICS-Ökosystem. Verstöße sind TECHNISCHER HOCHVERRAT.*
