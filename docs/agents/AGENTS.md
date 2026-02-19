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

---

## 🚨 NO-TIMEOUT POLICY

- NIEMALS timeout in opencode.json konfigurieren
- NIEMALS timeout in oh-my-opencode.json konfigurieren
- Modelle brauchen unterschiedlich lange (Qwen 3.5: 70-90s) - das ist OK!
- Timeouts führen zu Abbrüchen
- NIEMALS eintragen!

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

---

## 🎯 OH-MY-OPENCODE CATEGORIES (DEQLHI-LOOP)

**PFlicht:** Bei JEDER delegate_task() muss die richtige Category verwendet werden!

### Category-Liste (OFFIZIELL):

| Category | Model | Wann verwenden |
|----------|-------|----------------|
| **visual-engineering** | `kimi-for-coding/k2p5` | Frontend, UI/UX, Design, Styling, Animation |
| **ultrabrain** | `kimi-for-coding/k2p5` | ECHT schwere Logik-Aufgaben, nicht Step-by-Step |
| **deep** | `opencode-zen/kimi-k2.5-free` | Goal-oriented problem-solving, hairy problems |
| **artistry** | `kimi-for-coding/k2p5` | Unconventional, creative approaches |
| **quick** | `kimi-for-coding/k2p5` | Triviale Tasks, einzelne Files, Typos, simple fixes |
| **unspecified-low** | `kimi-for-coding/k2p5` | Low-effort Tasks |
| **unspecified-high** | `kimi-for-coding/k2p5` | High-effort Tasks |
| **writing** | `kimi-for-coding/k2p5` | Docs, Prose, technisches Schreiben |
| **general** | `kimi-for-coding/k2p5` | Allgemeine Tasks |

### ⚠️ KRITISCHE REGELN:

1. **Category IMMER angeben:**
   ```typescript
   // ✅ RICHTIG
   task(category="quick", load_skills=[], prompt="Fix typo")
   
   // ❌ FALSCH - Keine Category!
   task(load_skills=[], prompt="Fix typo")
   ```

2. **Model-Auswahl:**
   - `kimi-for-coding/k2p5` = DEFAULT für fast alles
   - `opencode-zen/kimi-k2.5-free` = NUR für deep reasoning (FREE!)
   - `moonshotai/kimi-k2.5` = Premium für Sisyphus only

3. **Skills hinzufügen wenn relevant:**
   ```typescript
   // Mit Skill
   task(category="visual-engineering", load_skills=["playwright"], prompt="...")
   
   // Ohne Skill (wenn nicht gebraucht)
   task(category="quick", load_skills=[], prompt="...")
   ```

### Available Skills:

| Skill | Wann |
|-------|------|
| `playwright` | Browser automation, testing, screenshots |
| `frontend-ui-ux` | Designer-style UI ohne Mockup |
| `git-master` | Git operations, commits, rebase |
| `dev-browser` | Navigation, forms, scraping |

### DEQLHI-LOOP Workflow (RICHTIG):

```typescript
// 1. Research (1-2 Agents) - IMMER parallel!
task(category="deep", load_skills=[], run_in_background=true, prompt="Research X")
task(category="writing", load_skills=[], run_in_background=true, prompt="Find docs")

// 2. Plan (1 Agent) - Nach Research
task(category="ultrabrain", load_skills=[], prompt="Create plan based on research")

// 3. Review (1 Agent)
task(category="ultrabrain", load_skills=[], prompt="Review plan quality")

// 4. Implement (MAX 3 Agents) - Parallel!
task(category="quick", load_skills=["git-master"], run_in_background=true, prompt="...")
task(category="visual-engineering", load_skills=["playwright"], run_in_background=true, prompt="...")
```

---

## Qwen 3.5 Skills

Dieses Projekt nutzt Qwen 3.5 (NVIDIA NIM) für spezialisierte KI-Aufgaben. Die folgenden Skills sind verfügbar:

### qwen_vision_analysis
Bildanalyse und visuelle Erkennung für Produktbilder, Grafiken und Diagramme.
- **Use Case:** Produktbild-Qualitätsprüfung, Layout-Analyse
- **Input:** Bilder (PNG, JPG, WebP)
- **Output:** Strukturierte Analyse mit Tags und Metriken
- **API:** `POST /api/qwen/vision` (Vercel Edge Function)

### qwen_code_generation
Full-Stack Code-Generierung mit Next.js, Go und Supabase.
- **Use Case:** Komponenten, API-Routen, Datenbank-Schema
- **Input:** Natürliche Sprache oder Spezifikation
- **Output:** Fertiger, getesteter Code
- **API:** `POST /api/qwen/chat` (Vercel Edge Function)

### qwen_document_ocr
Texterkennung und Dokumentanalyse aus gescannten Dokumenten und PDFs.
- **Use Case:** Rechnungsverarbeitung, Vertragsanalyse
- **Input:** PDF, Bilder mit Text
- **Output:** Extrahierter Text, Metadaten, Struktur
- **API:** `POST /api/qwen/ocr` (Vercel Edge Function)

### qwen_video_understanding
Video-Inhaltsanalyse für帧-Extraction und Szenenbeschreibung.
- **Use Case:** Video-Vorschau, Content-Indexierung
- **Input:** Videos (MP4, MOV, WebM)
- **Output:** Szenenbeschreibung, Key-Frames, Metadaten
- **API:** `POST /api/qwen/video` (Vercel Edge Function)

### qwen_conversation
Natürliche Konversations-KI für Kundenservice und Chat-Interaktionen.
- **Use Case:** Support-Chat, Produktberatung
- **Input:** Benutzer-Nachrichten, Kontext
- **Output:** Kontextbezogene Antworten
- **API:** `POST /api/qwen/chat` (Vercel Edge Function)

**Deployment:** Alle Skills laufen über Vercel Edge Functions mit NVIDIA NIM Backend.

## Abnahme-Check AGENTS
1. Regeln klar und widerspruchsfrei
2. NLM-CLI Pflicht enthalten
3. Rollen und Übergabeformat enthalten
4. Eskalationspfad enthalten

---
