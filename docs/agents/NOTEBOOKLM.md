# NOTEBOOKLM.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- NLM-First ist verbindlich und vor Implementierung zu berücksichtigen.
- Quellen-/Asset-Hygiene folgt den globalen Duplicate- und Sync-Regeln.
- Jeder Asset-Lebenszyklus muss reproduzierbar dokumentiert sein.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Betriebsleitfaden für NotebookLM-Nutzung über NLM-CLI in Projekten.

## Pflichtsatz
NotebookLM wird vollständig über NLM-CLI genutzt. Alle Outputs werden mit Qualitätsmatrix bewertet und nur verifiziert übernommen.

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

## Universalitäts-Regeln
1. Keine projektspezifischen Secrets in dieser Datei.
2. IDs, Quellen und Artefakte nur mit Platzhaltern dokumentieren.
3. Jeder NLM-Output braucht Zweck, Quelle, Score und Freigabe-Status.

## Platzhalter
- {NLM_WORKSPACE_ID}
- {NLM_NOTEBOOK_ID}
- {ASSET_ID}
- {SOURCE_REF}
- {OWNER_ROLE}

## NLM-CLI Prozess
1. Aufgabe klassifizieren: Video | Infografik | Präsentation | Datentabelle | Bericht | Mindmap | Podcast
2. Passende Vorlage aus `../∞Best∞Practices∞Loop.md` auswählen
3. Quellenrahmen definieren
4. NLM-Generierung durchführen
5. Qualitätsmatrix anwenden
6. Freigabe/Verwerfung dokumentieren
7. Freigegebenes Artefakt in `BIOMETRICS/NLM-ASSETS/...` ablegen und referenzieren

## Quellenregister (Template)

| Source-ID | Typ | Beschreibung | Gültig ab | Owner |
|---|---|---|---|---|
| SRC-001 | doc | {SOURCE_REF} | {DATE} | {OWNER_ROLE} |

## Notebook-Register (Template)

| Notebook-ID | Zweck | Scope | Status | Owner |
|---|---|---|---|---|
| {NLM_NOTEBOOK_ID} | {PURPOSE} | {SCOPE} | ACTIVE | {OWNER_ROLE} |

## Asset-Register (Template)

| Asset-ID | Typ | Thema | Qualitäts-Score | Status | Verwendungsort |
|---|---|---|---|---|---|
| {ASSET_ID} | video | {TOPIC} | {SCORE} | review | {TARGET_DOC} |

## Asset-Typen (Pflichtkatalog)
1. video
2. infographic
3. presentation
4. table
5. report
6. mindmap
7. podcast

## Ablagepfade (kanonisch)
- `BIOMETRICS/NLM-ASSETS/videos/`
- `BIOMETRICS/NLM-ASSETS/infographics/`
- `BIOMETRICS/NLM-ASSETS/presentations/`
- `BIOMETRICS/NLM-ASSETS/tables/`
- `BIOMETRICS/NLM-ASSETS/reports/`
- `BIOMETRICS/NLM-ASSETS/mindmaps/`
- `BIOMETRICS/NLM-ASSETS/podcasts/`

## README-Einbindungspflicht
Für jedes freigegebene NLM-Artefakt:
1. Kurzbeschreibung in `../README.md`
2. Verweis auf Asset-Pfad
3. Zuordnung zu Zielseite/Zielflow

## Qualitätsmatrix
- Korrektheit
- Konsistenz
- Verständlichkeit
- Zielgruppenfit
- Umsetzbarkeit
- Wiederverwendbarkeit
- Evidenzbezug

Freigabe:
- Mindestens 13/16
- Korrektheit 2/2

## Delegationsprotokoll (Template)
```text
Delegation-ID:
Datum:
Asset-Typ:
Verwendete Vorlage:
Quellen:
Output-Score:
Entscheidung: übernommen | verworfen
Begründung:
Nächste Iteration:
Owner:
```

## Verbotene Zustände
1. Übernahme ungeprüfter NLM-Ausgaben
2. Fehlende Quellenbezüge
3. Nicht dokumentierte Verwerfungen

## Abnahme-Check NOTEBOOKLM
1. NLM-CLI Prozess dokumentiert
2. Register für Quellen/Notebooks/Assets vorhanden
3. Qualitätsmatrix und Freigaberegel enthalten
4. Delegationsprotokoll enthalten
5. Mindmap- und Podcast-Regeln enthalten

---

---

## NLM Assets (BIOMETRICS)

### Video Assets

| Asset-ID | Beschreibung | Notebook-Quelle | Status |
|----------|--------------|-----------------|--------|
| VID-001 | Onboarding-Erklärung | biometrics-onboard | ✅ AKTIV |
| VID-002 | CLI Installationsanleitung | biometrics-onboard | ✅ AKTIV |
| VID-003 | Architektur-Übersicht | ARCHITECTURE.md | ✅ AKTIV |

**Ablageort:** `BIOMETRICS/NLM-ASSETS/videos/`

### Infografiken

| Asset-ID | Beschreibung | Notebook-Quelle | Status |
|----------|--------------|-----------------|--------|
| INF-001 | System-Architektur | ARCHITECTURE.md | ✅ AKTIV |
| INF-002 | Datenfluss-Diagramm | INFRASTRUCTURE.md | ✅ AKTIV |
| INF-003 | Installationsprozess | onboarding-flow | ✅ AKTIV |

**Ablageort:** `BIOMETRICS/NLM-ASSETS/infographics/`

### Berichte

| Asset-ID | Beschreibung | Notebook-Quelle | Status |
|----------|--------------|-----------------|--------|
| REP-001 | Sicherheitsanalyse | SECURITY.md | ✅ AKTIV |
| REP-002 | Performance-Report | INFRASTRUCTURE.md | ✅ AKTIV |
| REP-003 | Integrations-Übersicht | INTEGRATION.md | ✅ AKTIV |

**Ablageort:** `BIOMETRICS/NLM-ASSETS/reports/`

### Integration mit biometrics-onboard

Das **biometrics-onboard** Notebook ist das zentrale Notebook für alle Onboarding- und Installations-bezogenen Inhalte.

**Notebook-ID:** `{NLM_NOTEBOOK_ID_ONBOARD}`

**Enthaltene Quellen:**
- `biometrics-cli/README.md`
- `ONBOARDING.md`
- `INSTALLATION.md`
- `TROUBLESHOOTING.md`

**Generierte Artefakte:**
- Video: Installationsanleitung
- Infografik: Installationsprozess
- Mindmap: CLI-Befehlsübersicht

**Sync-Pflicht:** Nach jeder Änderung an CLI oder Onboarding-Dokumentation muss das Notebook synchronisiert werden:

```bash
# 1. Quellen prüfen
nlm source list {NLM_NOTEBOOK_ID_ONBOARD}

# 2. Geänderte Datei hinzufügen
nlm source add {NLM_NOTEBOOK_ID_ONBOARD} --file "biometrics-cli/README.md" --wait
```

---

## 13) NLM-CLI Installation & Setup

### Installation (PNPM ONLY!)

**⚠️ WICHTIG:** Verwende **pnpm** (NICHT npm!) gemäß Stack-Policy.

```bash
# NLM-CLI global installieren (PNPM!)
pnpm install -g nlm

# Installation verifizieren
nlm --version
```

### Authentication Setup

```bash
# Google OAuth durchführen
nlm auth login

# Browser öffnet sich für Google Login
# Mit Google Account anmelden
# Berechtigungen erteilen

# Auth Status prüfen
nlm auth status
```

### Erste Schritte

#### 1. Notebook erstellen
```bash
# Neues Notebook erstellen
nlm notebook create "Mein Projekt Notebook"

# Notebook ID notieren (für spätere Commands)
```

#### 2. Source hinzufügen
```bash
# Source zu Notebook hinzufügen
nlm source add {NOTEBOOK_ID} --file "dokument.md" --wait

# Source muss im aktuellen Verzeichnis existieren
```

#### 3. Query ausführen
```bash
# Einfache Query
nlm query notebook {NOTEBOOK_ID} "Was ist der aktuelle Stand?"

# Deep Research
nlm research start "{THEMA}" --mode deep --notebook-id {NOTEBOOK_ID}
```

### Wichtige Commands

#### Notebook Management
```bash
# Alle Notebooks auflisten
nlm notebook list

# Notebook Details anzeigen
nlm notebook show {NOTEBOOK_ID}

# Notebook löschen
nlm notebook delete {NOTEBOOK_ID}
```

#### Source Management
```bash
# Sources auflisten
nlm source list {NOTEBOOK_ID}

# Source löschen (bei Duplikaten!)
nlm source delete {SOURCE_ID} -y

# Source hinzufügen
nlm source add {NOTEBOOK_ID} --file "dokument.md" --wait
```

#### Query & Research
```bash
# Einfache Query
nlm query notebook {NOTEBOOK_ID} "{FRAGE}"

# Deep Research starten
nlm research start "{THEMA}" --mode deep

# Research Status prüfen
nlm research status {RESEARCH_ID}
```

### Best Practices 2026

#### 1. Duplicate Prevention
**IMMER vor `source add` prüfen:**
```bash
# 1. Sources auflisten
nlm source list {NOTEBOOK_ID}

# 2. Falls Duplikat existiert → LÖSCHEN
nlm source delete {OLD_SOURCE_ID} -y

# 3. DANN neue Source hinzufügen
nlm source add {NOTEBOOK_ID} --file "dokument.md" --wait
```

#### 2. Sync Pflicht
**Nach jeder relevanten Dateiänderung:**
```bash
# Datei zu Notebook synchronisieren
nlm source add {NOTEBOOK_ID} --file "geänderte-datei.md" --wait
```

#### 3. Crash-Tests
**Vor kritischen Entscheidungen:**
```bash
# Notebook nach aktuellem Stand fragen
nlm query notebook {NOTEBOOK_ID} "Was ist der aktuelle Stand zu {THEMA}?"
```

### Troubleshooting

#### Problem: "nlm: command not found"
**Lösung:**
```bash
# PNPM global binaries im PATH?
export PATH="$(pnpm bin -g):$PATH"

# In ~/.zshrc oder ~/.bashrc hinzufügen:
export PATH="$HOME/Library/pnpm/global/bin:$PATH"
```

#### Problem: "Authentication failed"
**Lösung:**
```bash
# Logout und erneuter Login
nlm auth logout
nlm auth login
```

#### Problem: "Notebook not found"
**Lösung:**
```bash
# Notebook ID prüfen
nlm notebook list

# Korrekte ID verwenden (UUID Format)
```

#### Problem: "Duplicate source"
**Lösung:**
```bash
# 1. Sources auflisten
nlm source list {NOTEBOOK_ID}

# 2. Duplikat identifizieren (gleicher Titel)

# 3. Altes löschen
nlm source delete {DUPLICATE_ID} -y

# 4. Neues hinzufügen
nlm source add {NOTEBOOK_ID} --file "dokument.md" --wait
```

### Stack-Policy Compliance

**Gemäß `AGENTS-GLOBAL.md` und `∞Best∞Practices∞Loop.md`:**

- ✅ **PNPM ONLY:** `pnpm install -g nlm` (NIEMALS npm!)
- ✅ **NLM First:** Immer zuerst NLM Query, dann externe Recherche
- ✅ **Duplicate Prevention:** Immer `source list` vor `source add`
- ✅ **Sync Pflicht:** Nach Änderungen sofort synchronisieren
- ✅ **Crash-Tests:** Vor Entscheidungen NLM konsultieren

### Quick Reference

```bash
# Installation
pnpm install -g nlm

# Auth
nlm auth login
nlm auth status

# Notebook
nlm notebook create "Name"
nlm notebook list

# Source
nlm source list {ID}
nlm source add {ID} --file "doc.md" --wait
nlm source delete {ID} -y

# Query
nlm query notebook {ID} "Frage"
nlm research start "Thema" --mode deep
```

---

**Updated:** 2026-02-17  
**Installation:** pnpm install -g nlm  
**Status:** PRODUCTION READY ✅
