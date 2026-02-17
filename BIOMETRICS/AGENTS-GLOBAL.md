# 🚀 AGENTS-GLOBAL.md — CEO EXECUTIVE MANDATE (LONG-FORM)

**Version:** 20.0-LF  
**Status:** ACTIVE — MUST READ BEFORE EVERY SESSION  
**Scope:** Global OpenCode Instructions Blueprint  
**Target:** `~/.config/opencode/AGENTS.md`  
**Project Stack Lock:** Next.js 15 + TypeScript strict, Go, Supabase, pnpm, NLM/NotebookLM

---

## 0) Executive Context
Diese Datei ist die **Langfassung** der globalen Agent-Mandate und dient als verlustarme Master-Vorlage.

Prinzipien:
1. Additiv statt destruktiv.
2. Keine Regelverluste.
3. Kein „fertig“ ohne Evidenz.
4. Stack-Konformität ist verpflichtend.

---

## 1) Absolute Priority Rules (Top 15)

### Rule 1 — Search before Create
- Vor jeder Dateioperation zuerst suchen/lesen.
- Existierende Strukturen werden erweitert, nicht dupliziert.

### Rule 2 — Read-before-Write
- Alle betroffenen Dateien vollständig lesen.
- Keine stichprobenartige Teilanalyse für kritische Änderungen.

### Rule 3 — Verify-then-Execute
- Vor Änderungen Diagnose.
- Nach Änderungen Lint/Typecheck/Test/Build/Runtime-Check.

### Rule 4 — No Fake Completion
- „Done“ nur bei nachweisbarer Erfüllung aller Akzeptanzkriterien.

### Rule 5 — Additive Integrity
- Keine Wissenslöschung, keine Verdichtung mit Informationsverlust.

### Rule 6 — No Blind Delete
- Unbekannte Dateien/Services/Configs niemals ohne RCA entfernen.

### Rule 7 — Security by Default
- Keine Secrets in Git.
- Input-Validierung, Fehlerhärtung, Least Privilege.

### Rule 8 — Parallel Execution Mandate
- Parallel arbeiten, wenn technisch möglich.
- Sequentiell nur bei echten Abhängigkeiten.

### Rule 9 — Todo Discipline
- Mehrstufige Arbeit nur mit sauberem Todo-Statusmodell.

### Rule 10 — Documentation is Product
- Doku ist Teil der Lieferung, nicht Nachtrag.

### Rule 11 — NLM First
- NotebookLM/NLM-CLI ist Primärkanal für Wissensarbeit.

### Rule 12 — Port Sovereignty
- Keine konfliktträchtigen Standardports im Multi-Projekt-Betrieb.

### Rule 13 — Container Naming Governance
- Einheitliche, systematische Containernamen sind Pflicht.

### Rule 14 — Git Discipline
- Konventionelle Commits, sauberer Verlauf, nachvollziehbare Intention.

### Rule 15 — Plan Sovereignty
- Keine Plan-Duplikate, keine planlose Parallelwelt.

---

## 2) Source-of-Truth Hierarchy
1. `~/.config/opencode/opencode.json`
2. `~/.config/opencode/AGENTS.md`
3. `[PROJECT_ROOT]/BIOMETRICS/AGENTS.md`
4. `[PROJECT_ROOT]/∞Best∞Practices∞Loop.md`
5. Fachdokumente in `[PROJECT_ROOT]/BIOMETRICS/*`

Konfliktauflösung:
1. Security/Integrity
2. Neueste explizite User-Anweisung
3. Projekt-Loop-Mandate
4. Globale AGENTS-Regeln
5. Lokale Spezifika

---

## 3) Stack Lock & Tech Baseline 2026

### 3.1 Frontend
- Next.js 15 (App Router)
- TypeScript strict
- Keine reine HTML-Hauptanwendung
- A11y, Performance, Security obligatorisch

### 3.2 Backend
- Go Services
- Supabase (DB/Auth/Storage/Edge)
- API-first, versioniert, robustes Error Contract

### 3.3 JavaScript Tooling
- pnpm only
- Kein npm-Workflow in JS-Projekten

### 3.4 Docs/Structure
- BIOMETRICS als kanonischer Governance-Hub
- NLM Assets unter `BIOMETRICS/NLM-ASSETS/`

---

## 4) NVIDIA NIM Mandate (Critical)

### 4.1 Latenz-Warnung
High-Latency-Modelle wie Qwen 3.5 397B können 70–90s oder mehr benötigen.

### 4.2 Timeout-Regel
- OpenCode Requests mit ausreichendem Puffer (mind. `120000ms`, sofern technisch unterstützt).
- Ziel: keine künstlichen Abbrüche bei validen, langsamen Inferenzpfaden.

### 4.3 Modell-Disziplin
- Korrekte ID pflegen (z. B. `qwen/qwen3.5-397b-a17b`).
- Falsche Modellzuordnung ist P0-Konfigurationsfehler.

### 4.4 Reference Snippet — OpenCode
```json
{
  "provider": {
    "nvidia-nim": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "NVIDIA NIM (Qwen 3.5)",
      "options": {
        "baseURL": "https://integrate.api.nvidia.com/v1",
        "timeout": 120000
      },
      "models": {
        "qwen-3.5-397b": {
          "id": "qwen/qwen3.5-397b-a17b",
          "limit": { "context": 262144, "output": 32768 }
        }
      }
    }
  }
}
```

### 4.5 Reference Snippet — OpenClaw
```json
{
  "models": {
    "providers": {
      "nvidia": {
        "baseUrl": "https://integrate.api.nvidia.com/v1",
        "api": "openai-completions",
        "models": ["qwen/qwen3.5-397b-a17b"]
      }
    }
  }
}
```

### 4.6 Operations
- Rate-Limit beachten (z. B. 40 RPM im Free-Tier).
- Bei 429: Backoff + Fallback-Kette.

---

## 5) OpenCode/OpenClaw Governance

### 5.1 Konfigurationsschutz
Verboten:
- Destruktive Neuinstallation ohne Notfallbegründung
- Blindes Löschen von `~/.config/opencode`, `~/.opencode`, `~/.oh-my-opencode`

Erlaubt:
- Diagnostik
- gezielte Reparatur
- Backup-gestützte Korrektur

### 5.2 Provider Schema
Pflicht: offizielles OpenCode Provider-Schema mit `@ai-sdk/openai-compatible` und `options.baseURL`.

---

## 6) Parallelism, Swarm, Delegation

### 6.1 Swarm-Kernregel
- Komplexe Aufgaben nicht monolithisch abarbeiten.
- Rollen aufteilen: Planung, Forschung, Implementierung, Testing, Review.

### 6.2 Delegationskontrakt (Pflicht)
Jeder Delegationsauftrag enthält:
1. Rolle
2. Ziel
3. Kontext
4. Read-first Dateien
5. Änderungsumfang
6. Akzeptanzkriterien
7. Tests
8. Doku-Updates
9. Risikoanalyse
10. Outputformat

### 6.3 Anti-Pattern
Verboten:
- „Warten auf Agenten“ als Leerlauf
- Delegation ohne Prüfschritte
- Ergebnisübernahme ohne Verifikation

---

## 7) Delqhi Loop (Infinite Production Mode)

Grundsatz:
- Kein Endzustand „fertig“, sondern kontinuierlicher Verbesserungsfluss.

Pflicht:
1. Nach Task-Abschluss Folgeaufgaben erzeugen.
2. Dokumentieren, testen, verifizieren.
3. Qualität und Security nicht zugunsten von Geschwindigkeit opfern.

Mantra:
- „Ein Task endet, fünf neue beginnen.“

---

## 8) Port Sovereignty & Container Convention

### 8.1 Port-Regeln
Vermeiden:
- Standardports, die Konflikte im Multi-Projekt-Betrieb verursachen.

Bevorzugen:
- Projekt- und Service-spezifische High-Port-Ranges (z. B. 50000–59999) mit Registry.

### 8.2 Container Naming
Format:
- `agent-XX-*`
- `room-XX-*`
- `solver-X.X-*`
- `builder-X-*`

---

## 9) NLM / NotebookLM Mission Control

### 9.1 NLM First Policy
Für Wissensarbeit gilt priorisiert:
1. NLM Query
2. Projektquellen
3. ergänzende Recherche

### 9.2 Duplicate Prevention
Vor jedem Upload:
```bash
nlm source list <notebook-id>
nlm source delete <old-source-id> -y   # falls Duplikat
nlm source add <notebook-id> --file "file.md" --wait
```

### 9.3 Sync Pflicht
Nach relevanten Dateiänderungen:
```bash
nlm source add <notebook-id> --file "<changed-file>" --wait
```

### 9.4 Crash-Test
Vor kritischen Entscheidungen:
```bash
nlm query notebook <notebook-id> "Was ist der aktuelle Stand zu <Thema>?"
```

### 9.5 Deep Research Mandate
```bash
nlm research start "<topic>" --mode deep --notebook-id <id>
```

### 9.6 CLI Syntax Guardrails
Richtig:
```bash
nlm query notebook <id> "<frage>"
```
Falsch:
```bash
nlm query <id> "<frage>"
```

---

## 10) Multi-Notebook Governance (Template)

### 10.1 Mindestsets pro Projekt
1. Patent
2. Website
3. Webshop
4. Webapp
5. Engine
6. Vision

### 10.2 Zusätzliche Notebooks
- Feature-Notebooks
- Page-Notebooks
- Research-/Wiki-Notebooks

### 10.3 Naming Convention
Format:
`[ID] DEV-[Projekt]-[Typ]-[Bezeichnung]-[Kategorie]`

Beispiele:
- `1.1.1 DEV-PetVat-Patent`
- `1.1.2 DEV-PetVat-Website`
- `1.1.20 DEV-PetVat-Engine-rPPG-Detection-Feature`

---

## 11) NotebookLM als Supervisor-Agent
NotebookLMs dürfen als spezialisierte Subagenten verwendet werden:
- Wiki-Notebook: Recherche
- Feature-Notebook: Implementierungsentscheidungen
- Patent-Notebook: Claims und Legal-Kontext
- Vision-Notebook: Produktstrategie

Workflow:
1. Notebook spezifizieren
2. Query ausführen
3. Antwort als Entscheidungsbasis nutzen
4. Ergebnis zurück in Notebook synchronisieren

---

## 12) Research Governance
Primär via NLM/NotebookLM.

Wenn externe Quelle nötig:
1. Quelle in Notebook aufnehmen.
2. Danach strukturierte Query stellen.
3. Ergebnisse dokumentieren.

Verboten:
- unstrukturierte, unprotokollierte Recherchepfade ohne Doku-Rückführung

---

## 13) Security Mandates

### 13.1 Secrets
- Nie in Git.
- Nie in öffentliche Protokolle.
- Nur sichere Stores/Env.

### 13.2 OWASP/CVE
- Kritische Features gegen bekannte Risiken prüfen.
- Input-Sanitization, AuthZ/AuthN, Rate-Limits, auditierbare Fehlerpfade.

### 13.3 Security Report Pflicht
Vor Feature-Abschluss mit Risiko-Potenzial: Security-Notiz mit Findings/Status.

---

## 14) Code Quality Standards

### 14.1 TypeScript Strict
```json
{
  "compilerOptions": {
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "exactOptionalPropertyTypes": true
  }
}
```

### 14.2 Error Handling
Verboten:
```typescript
try {
  await operation();
} catch (e) {
}
```

Pflicht:
```typescript
try {
  const result = await riskyOperation();
  return result;
} catch (error) {
  logger.error('Operation failed', { error });
  throw new CustomError('Descriptive message', { cause: error });
}
```

---

## 15) API Standards

### 15.1 REST / Contract
- klare Ressourcen
- konsistente Statuscodes
- einheitliches Fehlerformat

### 15.2 Rate Limit
- sensible Endpoints schützen
- 429 handling dokumentieren

### 15.3 Versioning
- Breaking Changes nur versioniert

---

## 16) Git & Delivery Discipline

### 16.1 Commit Standard
`<type>(<scope>): <description>`

Typen:
- feat
- fix
- docs
- refactor
- test
- chore

### 16.2 Branching
- nachvollziehbare Branch-Namen
- PR-Checks verpflichtend

### 16.3 No Silent Change
- Kein hidden behavior change ohne Doku

---

## 17) Testing Standards

Pflicht je Scope:
1. Unit/Integration passend zur Änderung
2. API-/Flow-Tests
3. Browser-/UI-Checks bei UI-Änderungen
4. Fehlerfall-Tests

Optional erweitert:
- Performance smoke
- Security regression

---

## 18) Performance Standards

### 18.1 Web Targets
- LCP, INP, CLS im grünen Bereich
- Bundle-Budgets überwachen

### 18.2 Backend Targets
- kritische Endpoints mit Response-Zielwerten
- N+1 und ineffiziente Queries vermeiden

---

## 19) MCP Governance

Pflicht-MCPs (wenn verfügbar):
1. Serena (Orchestrierung)
2. Chrome DevTools/CDP (visuelle und Browser-Verifikation)
3. Tavily (ergänzende Recherche)

Regel:
- Nutzung und Ergebnis in Projektdoku protokollieren.

---

## 20) Docker as MCP Wrapper Protocol

Grundsatz:
- HTTP-Container sind nicht automatisch MCP-stdio Server.

Pattern:
1. Container API
2. Wrapper (stdio ↔ HTTP)
3. OpenCode MCP als local command

---

## 21) Plan Sovereignty

Regeln:
1. Vor neuem Plan: bestehende Pläne prüfen.
2. Keine Plan-Duplikate.
3. Nur begrenzte Anzahl aktiver Pläne.
4. Alte Pläne archivieren, nicht verschwinden lassen.

---

## 22) Troubleshooting Ticket Mandate
Bei substanziellem Fehlerfall:
1. Ticket `ts-ticket-XX.md`
2. Problem
3. Root Cause
4. Lösungsschritte
5. Befehle
6. Referenzen

---

## 23) Documentation Trinity

### 23.1 Struktur
- User-Doku
- Dev-Doku
- Projekt-/Operations-Doku

### 23.2 Pflichtdateien
Mindestens konsistent halten:
- `BIOMETRICS/MEETING.md`
- `BIOMETRICS/CHANGELOG.md`
- `BIOMETRICS/ARCHITECTURE.md`
- `BIOMETRICS/COMMANDS.md`
- `BIOMETRICS/ENDPOINTS.md`
- `BIOMETRICS/NOTEBOOKLM.md`
- `BIOMETRICS/MAPPING.md`
- `BIOMETRICS/MAPPING-COMMANDS-ENDPOINTS.md`
- `BIOMETRICS/MAPPING-FRONTEND-BACKEND.md`
- `BIOMETRICS/MAPPING-DB-API.md`
- `BIOMETRICS/MAPPING-NLM-ASSETS.md`

---

## 24) NLM Assets Governance

Pflichtartefakte:
1. Video
2. Infografik
3. Präsentation
4. Datentabelle
5. Report
6. Mindmap
7. Podcast

Ablage:
- `BIOMETRICS/NLM-ASSETS/videos/`
- `BIOMETRICS/NLM-ASSETS/infographics/`
- `BIOMETRICS/NLM-ASSETS/presentations/`
- `BIOMETRICS/NLM-ASSETS/tables/`
- `BIOMETRICS/NLM-ASSETS/reports/`
- `BIOMETRICS/NLM-ASSETS/mindmaps/`
- `BIOMETRICS/NLM-ASSETS/podcasts/`

---

## 25) Mandatory Workflow (Every Task)
1. Kontext laden
2. relevante Regeln laden
3. Plan/Todos aktualisieren
4. Umsetzung
5. Verifikation
6. Doku aktualisieren
7. Nächste Tasks definieren

---

## 26) Mandatory Start Checklist
- [ ] Globale AGENTS gelesen
- [ ] Lokale AGENTS gelesen
- [ ] Lastchanges/Meeting gelesen
- [ ] Aktiver Plan identifiziert
- [ ] NLM Notebook identifiziert

---

## 27) Mandatory End Checklist
- [ ] Implementierung verifiziert
- [ ] Tests relevant grün
- [ ] Risiken dokumentiert
- [ ] Doku aktualisiert
- [ ] NLM Sync erfolgt

---

## 28) Global Rule Sync Protocol
Diese Datei dient als Sync-Vorlage für `~/.config/opencode/AGENTS.md`.

Pflicht:
1. Regelabgleich durchführen
2. fehlende Regeln additiv übernehmen
3. keine Regelgruppe entfernen

---

## 29) Rule-Ledger (Coverage against User Payload)

Abgedeckte Mandatsgruppen:
1. NVIDIA NIM / Provider / Timeout / Model IDs
2. OpenClaw + OpenCode Struktur- und Sicherheitsregeln
3. Parallelism + Swarm + Todo Continuation
4. Delqhi Infinite Loop
5. Port Sovereignty und Container Naming
6. No Blind Delete / Additive Integrity
7. Mandatory Git Discipline
8. NLM First + Duplicate Prevention + Sync
9. Multi Notebook System + Naming + Kategorien
10. Deep Research + Crash Test
11. CLI Syntax Guardrails
12. Website Framework Mandate (Next.js/Astro; kein HTML-only)
13. TypeScript Strict + Error Handling
14. Security/Secrets/OWASP
15. Documentation & Mapping Pflicht
16. Plan Sovereignty
17. Troubleshooting Ticketing
18. MCP Registry + Wrapper Pattern

---

## 30) Extended Commands (Reference)

### 30.1 NLM Basics
```bash
nlm list notebooks
nlm notebook create "<name>"
nlm source list <notebook-id>
nlm source delete <source-id> -y
nlm source add <notebook-id> --file "file.md" --wait
nlm source add <notebook-id> --url "https://..." --wait
nlm query notebook <notebook-id> "<frage>"
nlm status source <source-id>
```

### 30.2 OpenClaw/OpenCode Checks
```bash
openclaw models
openclaw doctor
openclaw gateway restart
opencode models
```

### 30.3 NVIDIA Model Check
```bash
curl -H "Authorization: Bearer $NVIDIA_API_KEY" https://integrate.api.nvidia.com/v1/models
```

---

## 31) Executive Non-Negotiables
Verboten:
1. Fake „done“
2. Blindes Löschen
3. Regelverlust bei Konsolidierung
4. Undokumentierte Architekturbrüche
5. Geheimnisse in Repo
6. npm in pnpm-only Kontext

Pflicht:
1. Nachvollziehbare Evidenz
2. Additive Governance
3. Reproduzierbare Verifikation
4. Dokumentationskonsistenz

---

## 32) Subagent Mandatory Block (Copy/Paste)
```text
ROLE:
GOAL:
BUSINESS CONTEXT:
TECH CONTEXT:
NON-NEGOTIABLE RULES:
FILES TO READ FIRST:
FILES ALLOWED TO EDIT:
FILES FORBIDDEN TO EDIT:
TASKS:
ACCEPTANCE CRITERIA:
REQUIRED TESTS:
REQUIRED DOCUMENTATION UPDATES:
DEPENDENCIES WITH OTHER AGENTS:
RISKS TO WATCH:
OUTPUT FORMAT:
TRUTH POLICY: Never claim done without evidence.
GLOBAL RULE SYNC POLICY:
- Read and align with ~/.config/opencode/AGENTS.md first.
- Preserve all inherited mandates; no rule loss allowed.
- Apply BIOMETRICS stack lock (Next.js, Go, Supabase, pnpm, NLM-first).
```

---

## 33) Application Note
Diese Datei ist bewusst als **Langfassung** angelegt und darf bei Übernahme in globale AGENTS-Dateien weiter erweitert werden.

Regel:
- Änderungen nur additiv.
- Kürzung nur mit expliziter Freigabe und Verlustfreiheitsnachweis.

---

## 34) Change Log
- 2026-02-17: Long-form Ausbau erstellt; Kurzfassung ersetzt.

---

## 35) Final Declaration
Diese Datei bildet die globale Governance-Basis für hochwertige, belastbare, auditierbare Agentenarbeit.

Leitsatz:
**Omniscience is not a goal; it is our technical starting point.**

---

## 36) FULL-PAYLOAD ANNEX POLICY (VERLUSTARME SPIEGELUNG)
Dieser Annex erweitert die Langfassung um operative Detailregeln aus dem umfangreichen Mandatskorpus.

Regeln:
1. Annex-Inhalte sind append-only.
2. Rule-IDs bleiben stabil.
3. Bei Konflikten werden Addendum-Regeln ergänzt statt Inhalte entfernt.
4. BIOMETRICS Stack-Lock bleibt übergeordnet bindend.

---

## 37) Rule -1: Vollständige autonome Ausführung

Leitlinie:
- Wenn die Laufzeitumgebung es erlaubt, führt der Agent ausführbare Schritte selbst aus.
- Keine Übergabe trivialisierbarer Terminalarbeit an den User.

Einschränkung:
- Wenn Umgebung/Permissions/Provider blockieren, Blocker klar dokumentieren und sichere Alternative liefern.

Operational:
1. Diagnose selbst durchführen.
2. Alternativen selbst ausprobieren.
3. Nur bei echter technischer Unmöglichkeit eskalieren.

---

## 38) Rule -2: Mandatory Coder Workflow Protocol

### Phase 1 — Context Acquisition
Pflichtlektüre:
1. Global AGENTS
2. Lokale AGENTS
3. `BIOMETRICS/MEETING.md`
4. `BIOMETRICS/CHANGELOG.md`
5. `BIOMETRICS/ARCHITECTURE.md`
6. `BIOMETRICS/NOTEBOOKLM.md`

### Phase 2 — Research & Best Practices 2026
1. NLM Query (deep mode bei komplexen Themen)
2. Offizielle Dokumentation als Quelle in NLM aufnehmen
3. Architekturvarianten vergleichen

### Phase 3 — Internal Documentation Alignment
1. Mapping-Dateien prüfen
2. Commands ↔ Endpoints ↔ DB ↔ Frontend/Backend Synchronität prüfen

### Phase 4 — Master Plan Creation
1. Ziele
2. Risiken
3. Akzeptanzkriterien
4. Verifikation

### Phase 5 — Swarm Delegation
1. Planner
2. Research
3. Dev
4. Test
5. Review

---

## 39) Rule -3: Todo Continuation + Swarm Delegation

Mandat:
1. Keine Multi-Step-Aufgabe ohne Todo-Tracking.
2. Keine unstrukturierte Einzelagentenarbeit bei komplexen Tasks.
3. Statusmodell: `pending`, `in_progress`, `blocked`, `completed`.

Pflichtworkflow:
1. Todos anlegen
2. Parallelteile delegieren
3. Ergebnisse verifizieren
4. Todos aktualisieren
5. Follow-up-Todos hinzufügen

Verboten:
- Aufgabe ohne Taskstruktur starten
- „fertig“ ohne Checkliste

---

## 40) Rule -5: Absolutes Verbot von blindem Löschen

Verboten:
1. „Kenne ich nicht → löschen“
2. „Sieht alt aus → entfernen“
3. „Verstehe ich nicht → weg damit“

Pflicht:
1. Ursache recherchieren
2. Kontextquellen lesen
3. Integrationsabsicht prüfen
4. Entscheidung dokumentieren

No-Delete Protocol:
1. Entdeckung
2. Analyse
3. Dokumentation
4. kontrollierte Maßnahme

---

## 41) Rule -6: Mandatory Git Commit Discipline

Nach signifikanten Taskabschlüssen:
1. Änderungen stagen
2. Konventionell committen
3. Remote pushen (wenn Repo-Policy aktiv)

Commit Typen:
- feat
- fix
- docs
- refactor
- test
- chore

Beispiel:
```bash
git add -A
git commit -m "docs(agents): extend global mandate annex"
git push origin main
```

---

## 42) Rule -9: Port Sovereignty (No Standard Ports)

Verbotene Ports (Beispielkatalog):
- 3000
- 5432
- 8080
- 6379
- 5678
- 8000
- 9000
- 3306
- 27017
- 9200
- 80
- 443

Pflicht:
1. Service-Port-Registry pflegen
2. Eindeutige High-Port-Zuweisung
3. Port-Konflikte vor Deployment prüfen

Compliance-Check Beispiel:
```bash
grep -E '"(3000|8080|5432|6379|5678|8000|9000|3306|27017|9200|80|443):' docker-compose.yml
```

---

## 43) Rule -11: Parallele Agentenarbeit & Datei-Prüfung

### 43.1 Parallelism
1. Delegierbare Tasks parallelisieren.
2. Ergebnisse nicht blind übernehmen.

### 43.2 Datei-Erstellung
Vor jeder Neuanlage:
1. Existenz prüfen
2. ähnliche Dateien prüfen
3. Wiederverwendung bevorzugen

---

## 44) Rule -13: INTEGRATION.md Pflicht
Jedes Projekt muss eine Integrationsübersicht führen.

Pflichtinhalte:
1. Externe APIs
2. Datenbanken
3. externe Services/Tools
4. Security/Auth Integrationen
5. Monitoring/Logging
6. Deployment/CI

Template-Schnitt:
```markdown
# INTEGRATION.md
## Externe APIs
## Datenbanken
## Services & Tools
## Security & Auth
## Monitoring
## Deployment
```

---

## 45) Rule -14: Einfache Lösungen zuerst
Verhalten bei Userfragen:
1. zuerst direkte praktikable Lösung
2. danach Details nur bei Bedarf
3. keine unnötige Alternativenflut

Verboten:
- Überkomplexität als Standardantwort
- theoretische Umwege ohne Nutzen

---

## 46) OpenCode API Lessons Learned

Kernaussage:
- OpenCode-APIs können vom OpenAI-Format abweichen.

Pflicht:
1. API-Dokumentation zuerst prüfen
2. Endpunkte verifizieren
3. Payload-Typen gegen Schema validieren

Anti-Pattern:
- OpenAI-kompatibles Format blind annehmen

---

## 47) Native CDP Preference für Performance-Pfade

Leitlinie:
- Bei High-Performance-Browser-Automation native CDP bevorzugen, wenn geeignet.

Ziel:
1. geringere Latenz
2. weniger Overhead
3. bessere Skalierung

Zusatz:
- Caching-Strategien und Screenshot-Optimierung dokumentieren.

---

## 48) NLM Mission Control — Extended

### 48.1 Multi-Notebook System
Regel:
1. Themen trennen
2. Quellen sauber zuordnen
3. IDs dokumentieren

### 48.2 Sync Pflicht
Bei jeder relevanten Änderung:
```bash
nlm source add <id> --file "<file>" --wait
```

### 48.3 Migration Safety
Verboten:
- altes Notebook vor abgeschlossener Migration löschen

Pflicht:
1. neues Notebook erstellen
2. Quellen übertragen
3. Validierung
4. erst dann Cleanup

### 48.4 Notebook-Auswahlmatrix (Kurz)
| Zweck | Notebook-Typ |
|------|---------------|
| Code-Implementierung | Code/Development |
| Hardware/Sensorik | Hardware |
| Externe Docs | Research |
| Patent/Legal | Patent |
| Business/Go-To-Market | Marketing/Business |
| Strategie | Framework/Vision |

### 48.5 Naming Standard
Format:
`[ID] DEV-[Projekt]-[Typ]-[Bezeichnung]-[Kategorie]`

---

## 49) NotebookLM First Query Blocks (Copy/Paste)

Vor Dateierstellung:
```bash
nlm query notebook <id> "Was sind die Best Practices für <dateityp>?"
```

Vor Codeänderung:
```bash
nlm query notebook <id> "Wie ist die Architektur für <komponente>?"
```

Bei Fehler:
```bash
nlm query notebook <id> "Wie löse ich <error message>?"
```

---

## 50) NLM Duplicate & Source Hygiene

Standardablauf:
```bash
nlm source list <notebook-id>
nlm source delete <duplicate-id> -y
nlm source add <notebook-id> --file "<file>" --wait
nlm source list <notebook-id>
```

Regel:
- Kein blinder Re-Upload ohne vorherige Quellprüfung.

---

## 51) Website/Webapp/Webshop Framework Mandate

Pflicht:
1. Next.js 15 + TypeScript für moderne Web-Properties
2. Astro optional für content-heavy Seiten

Verboten:
1. reine HTML-Sites als Hauptlieferung
2. undokumentierte Legacy-Stacks

---

## 52) TS/ESLint/Biome Baselines

TypeScript strict bleibt Pflicht.  
Linting und Formatierung müssen vor Abschluss sauber sein.

Optionaler Baseline-Stack:
1. ESLint Flat Config
2. Biome für schnelles Formatting/Linting
3. pre-commit hooks für Mindestqualität

---

## 53) API Design & Error Standards

Pflicht:
1. konsistente Fehlerobjekte
2. klare Statuscodes
3. nachvollziehbare Rate-Limit-Strategie
4. versionierbare Breaking Changes

---

## 54) Security by Default — Extended Checklist
- [ ] Input validiert
- [ ] AuthZ/AuthN robust
- [ ] Secrets sauber isoliert
- [ ] Logging ohne Sensitive Leakage
- [ ] Known Vulnerabilities geprüft

---

## 55) CI/CD Governance

Pflichtjobs:
1. lint
2. typecheck
3. test
4. build

Empfohlen:
1. security scan
2. dependency updates
3. branch protection

---

## 56) Mandatory Status Footer (Template)
```text
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 STATUS UPDATE

Updated:       ☑️ AGENTS-GLOBAL.md
               ☑️ BIOMETRICS/MEETING.md
               ☑️ BIOMETRICS/CHANGELOG.md

FORTSCHRITT:   ████████░░ 80% (Implementierung)
               ██████░░░░ 60% (Verifikation)
               ░░░░░░░░░░  0% (Deployment)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

---

## 57) Global Sub-Agent Guide Pointer
Empfehlung:
- Zusätzlich zur Globaldatei sollte ein Sub-Agent Guide gepflegt werden, der Read-Order, häufige Fehler und Erfolgskriterien standardisiert.

---

## 58) Preservation of Knowledge Mandate
Diese Datei ist Wissensspeicher und Regelkorpus.

Verboten:
1. Verdichtung mit Regelverlust
2. Löschen historischer Mandate ohne Freigabe

Pflicht:
1. Erweiterung als neue Abschnitte
2. Änderungshistorie pflegen

---

## 59) Annex Command Deck

NLM:
```bash
nlm list notebooks
nlm notebook create "<name>"
nlm source list <id>
nlm source add <id> --url "<url>" --wait
nlm source add <id> --file "<file>" --wait
nlm query notebook <id> "<question>"
```

OpenClaw/OpenCode:
```bash
openclaw models
openclaw doctor
openclaw gateway restart
opencode models
```

Provider Test:
```bash
curl -H "Authorization: Bearer $NVIDIA_API_KEY" https://integrate.api.nvidia.com/v1/models
```

---

## 60) Full-Payload Compliance Note
Diese Datei ist jetzt als großvolumige operative Vorlage aufgebaut.

Regel:
1. Weitere Originalblöcke können als zusätzliche Annex-Module (`61+`) angehängt werden.
2. Kein Downscaling ohne explizite Freigabe.

---

## 61) DEQLHI-LOOP OPERATIONS BOARD (EXPANDED)

### 61.1 Kernprinzip
- Nach jedem abgeschlossenen Task werden neue priorisierte Tasks erzeugt.
- Keine Leerlaufzustände.
- Kontinuierliche Wertlieferung.

### 61.2 Arbeitsregeln
1. Niemals passiv warten.
2. Parallelisierbare Arbeitspakete parallel laufen lassen.
3. Dokumentation ist bei jedem Schritt Pflicht.
4. Qualitätsgates pro Task erzwingen.

### 61.3 Loop-Mechanismus (operativ)
1. Task abschließen
2. Verifikation
3. Dokumentation
4. Folgeaufgaben erzeugen
5. Nächster Task

---

## 62) Problem Solving Protocol (Master-CEO-Mode)

Wenn ein Problem auftritt:
1. Sofort strukturierte Recherche
2. Lösung A testen
3. falls fehlgeschlagen: Lösung B
4. falls fehlgeschlagen: Lösung C
5. niemals ohne evidenzbasierten Fix aufgeben

Regel:
- „Geht nicht“ ist kein Abschlusskriterium.

---

## 63) Vollständiges Lesen kritischer Dateien

Pflicht:
1. Relevante Dateien vollständig lesen.
2. Nicht nur Snippets, wenn Architekturentscheidungen betroffen sind.
3. Vor Implementierung muss Kontext vollständig sein.

Verboten:
1. Stichproben-Entscheidungen bei kritischen Änderungen.

---

## 64) Rule-Layer: Immutability of Knowledge (Supreme)

Kern:
1. Keine Regelzeile entfernen, wenn Informationsverlust entsteht.
2. Änderungen sind additive Ergänzungen.
3. Historische Entscheidungen bleiben nachvollziehbar.

Bei Konsolidierung:
1. Ursprungswissen übernehmen
2. Mapping auf neue Struktur
3. Abweichungen dokumentieren

---

## 65) Safe Migration & Consolidation Law

Pflichtprotokoll:
1. Quellartefakt vollständig lesen
2. Backup/Altstand sichern
3. neue Struktur erstellen
4. vollständige Übernahme
5. Mehrfachverifikation
6. erst dann kontrolliertes Cleanup

---

## 66) Ticket-Based Troubleshooting Mandate

Jeder relevante Fehler bekommt ein Ticket:
1. Problem Statement
2. Root Cause
3. Fix Steps
4. Befehle/Änderungen
5. Quellen

Referenzierung:
- zentrale Troubleshooting-Datei verlinkt auf Einzel-Tickets.

---

## 67) Global Secrets Registry Discipline

Regelrahmen:
1. Secrets niemals ins Repo.
2. Secret-Änderungen nachvollziehbar dokumentieren.
3. Alte Secrets nicht blind löschen — als deprecated/rotated markieren.

Template:
```markdown
## [SERVICE] - [YYYY-MM-DD]
Status: ACTIVE | DEPRECATED | ROTATED
Endpoint:
Owner:
Rotation:
Notes:
```

---

## 68) Local Project Knowledge Mandate

Jedes Projekt benötigt lokale AGENTS-Dokumentation mit:
1. Stack
2. Architektur
3. Konventionen
4. API-Standards
5. Spezialregeln
6. Troubleshooting

Integritätscheck vor Antworten:
1. lokale AGENTS gelesen?
2. Antwort konform?
3. notwendige Doku-Updates erkannt?

---

## 69) Lastchanges / Photographic Memory Mandate

Vor Session:
1. lastchanges lesen
2. letzten Stand extrahieren
3. nächste Schritte ableiten

Nach Session:
1. append-only Logeintrag
2. Beobachtungen/Fehler/Lösungen/Nächste Schritte

---

## 70) Best Practices 2026 Research Mandate

Phasen:
1. vor Planung: aktuelle Standards prüfen
2. während Planung: Alternativen/Debts/CVEs
3. während Coding: Fehler gezielt recherchieren
4. bei Problemen: verifizierte Lösungen priorisieren

Quellenpriorität:
1. offizielle Dokumentation
2. Produktionsbeispiele
3. strukturierte Sekundärquellen

---

## 71) Self-Critique & Crash Tests

Pflichtcheck:
1. Wie könnte der Code crashen?
2. Welche Edge Cases fehlen?
3. Fehlerbehandlung vollständig?

Crashtests:
1. invalid input
2. Grenzwerte
3. parallele Zugriffe
4. Netzwerk/DB-Ausfälle

---

## 72) Planning & Error Prevention Matrix

### 72.1 Planung
Jede Phase braucht:
1. Meilensteine
2. erwartetes Ergebnis
3. Akzeptanzkriterien
4. Risiken

### 72.2 Prävention
Für jedes Risiko:
1. Präventivmaßnahme
2. Fallback
3. Monitoring

---

## 73) Docker Knowledge Base Strategy

Vorgabe:
- Wissensmanagement in eigener Infrastruktur priorisieren.

Pflichten:
1. Architekturentscheidungen dokumentieren
2. Bugfix-Lösungen dokumentieren
3. regelmäßige Statusupdates

---

## 74) Market Analysis Mandate

Analyse-Dimensionen:
1. Features
2. Tech-Stack
3. Performance
4. UX
5. Preis/Positionierung
6. Innovation

Bewertung:
- führend
- wettbewerbsfähig
- nachholbedarf

---

## 75) Workspace Tracking / Collision Avoidance

Pflicht:
1. aktiven Arbeitsbereich dokumentieren
2. Konflikte früh erkennen
3. bei Dateikollisionen koordinieren

Statuswerte:
- IN_PROGRESS
- COMPLETED
- PENDING
- BLOCKED

---

## 76) OpenCode Preservation Mandate

Verboten:
1. destruktive Neuinstallation als Standardfix
2. Löschung zentraler OpenCode-Konfigurationspfade

Pflicht:
1. Diagnose
2. gezielte Reparatur
3. Backup-restore falls nötig

---

## 77) ALL-MCP Directory Documentation Standard

Für jeden MCP-Server sollten drei Artefakte gepflegt sein:
1. readme
2. guide
3. install

Pflichtinhalte:
1. Typ
2. Quellen/Links
3. Use Cases
4. Setup
5. Troubleshooting

---

## 78) GitHub Repository Standards Annex

Pflichtstruktur:
1. `.github/ISSUE_TEMPLATE/*`
2. `.github/PULL_REQUEST_TEMPLATE.md`
3. `.github/workflows/*`
4. `.github/CODEOWNERS`
5. `.github/dependabot.yml`
6. Root: `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, `LICENSE`

CI Mindestjobs:
1. lint
2. typecheck
3. test
4. build

---

## 79) Branch Protection Baseline

Main Branch:
1. PR review required
2. status checks required
3. up-to-date branch required
4. force push disabled
5. deletion disabled

---

## 80) Docker Container as MCP Wrapper Annex

Architektur:
1. HTTP Container
2. stdio Wrapper
3. local MCP config

Verboten:
- HTTP-Container direkt als remote MCP deklarieren, wenn stdio-Bridge erforderlich ist.

---

## 81) 26-Room / Naming Governance Annex

Container Naming Format:
- `{category}-{number}-{name}`

Kategorien:
1. `agent-*`
2. `room-*`
3. `solver-*`
4. `builder-*`

Regel:
- Service Name und Containername konsistent halten.

---

## 82) Provider Configuration Annex

Offizielles Custom-Provider-Muster:
```json
{
  "provider": {
    "custom-name": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "Display Name",
      "options": {
        "baseURL": "https://api.example.com/v1"
      },
      "models": {
        "model-id": {
          "name": "Model Name",
          "limit": { "context": 100000, "output": 10000 }
        }
      }
    }
  }
}
```

Ungültige Feldbeispiele (nicht nutzen):
1. `apiEndpoint`
2. `authentication` in inkompatibler Struktur

---

## 83) Fallback Chain Strategy Annex

Empfohlene Reihenfolge:
1. Primärmodell
2. spezialisierte Fallbacks
3. vision-fähige Modelle
4. generalistische Modelle

Regel:
- Jede Fallbackstufe muss dokumentiert und reproduzierbar sein.

---

## 84) Extended NLM Notebook Registry Skeleton

### 84.1 Projektfamilien
1. Hauptprojekt
2. Unterprojekte
3. je Unterprojekt: Patent/Website/Webshop/Webapp/Engine/Vision

### 84.2 Betriebsregeln
1. IDs dokumentieren
2. Quelle-Zuständigkeit klar halten
3. Query-Protokolle für wichtige Entscheidungen dokumentieren

### 84.3 Feature-/Page-Notebooks
Pflicht bei skalierenden Projekten:
1. pro kritischem Feature eigenes Notebook
2. pro kritischer Seite eigenes Notebook

---

## 85) V20.0 FULL PAYLOAD DECLARATION

Status:
1. Long-form vorhanden
2. Annex `36–60` vorhanden
3. Annex `61–85` vorhanden

Regel für weitere Ausbauten:
1. Nächste Erweiterungen laufen als `86+` Additive Blocks.
2. Verbatim-nahe Blöcke werden bevorzugt erhalten.
3. Keine Kürzung ohne explizite Freigabe.

Leitsatz:
**Produziere, dokumentiere, verifiziere, erweitere.**

---

## 86) Executive Rules Matrix (Expanded)

| Domain | Pflicht | Verbot | Done-Kriterium |
|---|---|---|---|
| Planning | Phasenplan + Risiken | Ad-hoc ohne Plan | Plan + Risiko + Akzeptanz |
| Coding | stack-konforme Umsetzung | Scope-Drift | Tests + Lint + Typecheck |
| Security | Secrets-Hygiene + OWASP | Klartext-Secrets | Security-Notiz vorhanden |
| Docs | Update aller betroffenen Dokus | Code-only Delivery | Doku-Mapping konsistent |
| NLM | Query-first + Sync | blindes source add | Source-Hygiene bestätigt |

---

## 87) Enterprise Output Contract

Jede signifikante Lieferung enthält:
1. Umsetzungsstatus
2. Verifikationsergebnisse
3. Risiken
4. Nächste Schritte
5. Dokumentationsänderungen

Minimal-Format:
```text
WHAT CHANGED:
VALIDATION:
RISKS:
NEXT:
DOCS UPDATED:
```

---

## 88) Mandatory Acceptance Criteria Schema

```text
Task-ID:
Goal:
Business Impact:
Technical Scope:
Acceptance Criteria:
- Functional
- Quality
- Security
- Documentation
Validation Commands:
Evidence:
```

---

## 89) Global Definition of Done (Extended)

Ein Task ist nur DONE, wenn:
1. Anforderungen vollständig erfüllt sind
2. betroffener Lint/Typecheck/Test/Build grün ist
3. Laufzeitpfad validiert ist
4. Security-Grundprüfung erfolgt ist
5. Doku + Mapping aktualisiert sind
6. offene Risiken dokumentiert sind

---

## 90) Command ↔ Endpoint ↔ DB Mapping Mandate

Pflichtprüfung bei Backend-relevanten Änderungen:
1. Command in `BIOMETRICS/COMMANDS.md` ergänzt/angepasst
2. Endpoint in `BIOMETRICS/ENDPOINTS.md` ergänzt/angepasst
3. DB-Impact in `BIOMETRICS/MAPPING-DB-API.md` dokumentiert
4. Frontend-Bindung in `BIOMETRICS/MAPPING-FRONTEND-BACKEND.md` dokumentiert

---

## 91) NLM Asset Lifecycle Governance

Lifecycle:
1. Briefing
2. Generierung
3. Qualitätsbewertung
4. Ablage in `NLM-ASSETS`
5. Referenzierung in Fachdokumenten
6. Iteration/Versionierung

Asset-Mindestmetadaten:
1. Zielgruppe
2. Zweck
3. Erstellungsdatum
4. Version
5. Autor/Agent
6. Qualitätsstatus

---

## 92) NLM Quality Matrix (0–2)

| Kriterium | 0 | 1 | 2 |
|---|---|---|---|
| Korrektheit | falsch | teilweise | konsistent korrekt |
| Vollständigkeit | lückenhaft | brauchbar | vollständig |
| Klarheit | unverständlich | verständlich | präzise |
| Konsistenz | widersprüchlich | teils konsistent | vollständig konsistent |
| Umsetzbarkeit | vage | teils umsetzbar | direkt umsetzbar |

Freigabe:
- Mindestscore gesamt ≥ 7
- kein Kriterium mit 0 bei produktiven Artefakten

---

## 93) Video Prompt Blueprint (NLM)

```text
ROLE: NLM Video Content Generator
OBJECTIVE: Erzeuge ein umsetzbares Video-Konzept für <Thema>
AUDIENCE: <Zielgruppe>
STYLE: <Ton/Stil>
CONSTRAINTS:
- keine Halluzinationen
- klare Struktur
- CTA verpflichtend
OUTPUT:
1) Hook
2) Storyboard (Szenenliste)
3) Voiceover Text
4) Visual Hinweise
5) CTA
```

---

## 94) Infographic Prompt Blueprint (NLM)

```text
ROLE: NLM Infographic Architect
GOAL: Strukturierte, visuell klare Infografik zu <Thema>
OUTPUT:
1) Headline
2) Kernmetriken
3) Abschnittsstruktur
4) Datenquellen-Hinweis
5) Fußnote/CTA
```

---

## 95) Presentation Prompt Blueprint (NLM)

```text
ROLE: NLM Presentation Strategist
GOAL: Executive-ready Präsentation
SLIDES:
1) Problem
2) Lösung
3) Architektur
4) Risiko
5) Roadmap
6) KPI/ROI
7) Next Steps
```

---

## 96) Data Table Prompt Blueprint (NLM)

```text
ROLE: NLM Data Table Designer
GOAL: Entscheidungsfähige Tabellenstruktur
OUTPUT:
1) Tabellenkopf
2) Spaltendefinitionen
3) Beispielzeilen
4) Validierungsregeln
5) Interpretationshinweise
```

---

## 97) Mindmap Policy (Expanded)

Mindmap-Mindestinhalte:
1. Zentrales Thema
2. Hauptzweige
3. Unterzweige
4. Abhängigkeiten
5. Prioritäten
6. Risiken
7. nächste Aktionen

---

## 98) Podcast Policy (Expanded)

Podcast-Mindestinhalte:
1. Episodenziel
2. Zielgruppe
3. Kapitelplan
4. Kernbotschaften
5. CTA
6. Follow-up Aufgaben

---

## 99) Website / Webapp / Webshop Architecture Guardrails

### 99.1 Website
- Marketing + Information + Vertrauen

### 99.2 Webapp
- Rollen, Flows, produktive Interaktion

### 99.3 Webshop
- Commerce, Checkout, Compliance

Pflicht:
1. Jede Säule dokumentieren
2. Überschneidungen mappen
3. Risiken je Säule ausweisen

---

## 100) Supabase Governance Annex

Pflichtinhalte pro Tabelle:
1. Zweck
2. Schlüssel
3. Relationen
4. RLS-Policy
5. Zugriffspfade
6. Migrationshinweise

Zusatz:
- Auth- und Edge-Funktionen dokumentieren

---

## 101) OpenClaw Operations Annex

Pflicht:
1. Betriebsmodus dokumentieren
2. Modellrouting dokumentieren
3. Health/Doctor-Kommandos definieren
4. Fehlerbehandlung und Fallback standardisieren

---

## 102) n8n Operations Annex

Für jeden Workflow:
1. Trigger
2. Inputs
3. Outputs
4. Fehlerpfade
5. Retry/Recovery
6. Besitzer

---

## 103) Cloudflare / Vercel / IONOS Governance

Pflichtblöcke:
1. DNS
2. Routing
3. Zertifikate
4. Exposure-Risiken
5. Rollback-Prozesse

---

## 104) Security Incident Protocol

Bei sicherheitsrelevantem Ereignis:
1. Incident erfassen
2. Eindämmung
3. Ursachenanalyse
4. Behebung
5. Prävention
6. Dokumentation

Incident-Template:
```text
Incident-ID:
Severity:
Scope:
Timeline:
Root Cause:
Fix:
Prevention:
```

---

## 105) Performance Incident Protocol

Wenn SLA/latency verletzt:
1. messen
2. Engpass lokalisieren
3. kurzfristiger Workaround
4. nachhaltiger Fix
5. Re-test + Dokumentation

---

## 106) Governance Change Protocol

Regeländerungen nur mit:
1. Begründung
2. Auswirkungen
3. Migrationshinweis
4. kompatibler Übergang

---

## 107) Compliance Report Skeleton

```markdown
# Compliance Report
## Scope
## Applied Rules
## Validation Evidence
## Deviations
## Risk Assessment
## Action Plan
```

---

## 108) Executive Review Checklist

- [ ] Ziele erreicht
- [ ] Risiken akzeptabel
- [ ] Doku vollständig
- [ ] Tests/Evidenz vorhanden
- [ ] Betrieb stabil
- [ ] nächste Prioritäten definiert

---

## 109) Verbatim Preservation Notes

Diese Datei enthält sowohl normalisierte als auch verbatim-nahe Mandatsblöcke.

Regel:
1. Bei weiteren Ausbauten möglichst originalnahe Formulierungen beibehalten.
2. Normalisierung nur, wenn Verständlichkeit deutlich steigt und Informationsgehalt vollständig bleibt.

---

## 110) Expansion Gate (to 2k+)

Aktueller Gate-Status:
1. Long-Form + Annex 36–85 vorhanden
2. Annex 86–110 vorhanden

Nächster Ausbauschritt (`111+`):
1. zusätzliche verbatim Module
2. detailliertere Tabellen/Beispiele
3. projektspezifische Operational Playbooks

Leitsatz:
**Scale documentation with zero rule loss.**

---

## 111) Slash Command Protocol & Autonomy Law

Prinzip:
- Systeme müssen durch Agenten steuerbar sein, nicht nur manuell bedienbar.

Pflicht je Projekt:
1. `SLASH.md` mit allen verfügbaren Kommandos.
2. Jede mutierbare Entität über API oder Slash-Command steuerbar.
3. Command-Syntax konsistent halten: `/cmd action target --param value`.

Beispiel:
```text
/product update "Super Shoes" --price 99.99
```

---

## 112) Trinity Documentation Standard (Expanded)

Pflichtstruktur:
```text
/project
  /docs/non-dev
  /docs/dev
  /docs/project
  /docs/postman
  DOCS.md
  README.md
```

Regeln:
1. README als Gateway
2. DOCS.md als Verfassung
3. API-Sammlungen gepflegt und versioniert

---

## 113) README Gateway Standard

README Mindestmodule:
1. Introduction
2. Quick Start
3. API Reference
4. Tutorials
5. Troubleshooting
6. Changelog & Support

---

## 114) Modern CLI Toolchain Mandate

Bevorzugte Tools:
1. `rg` statt `grep`
2. `fd` statt `find`
3. `sd` statt `sed`
4. `bat` statt `cat`

Regel:
- Moderne Tools bevorzugen, Fallbacks nur mit Begründung.

Beispiel:
```bash
rg "pattern" src/
fd -e ts -t f
sd "old" "new" file.txt
```

---

## 115) OpenHands Universal Coding Layer (Annex)

Mandat:
- Wenn im Zielsystem ein zentraler Code-Service definiert ist, werden Coding-Tasks darüber vereinheitlicht orchestriert.

Pflicht:
1. Endpoint-Registry pflegen
2. Commands + APIs konsistent dokumentieren
3. Audit-Trail sicherstellen

---

## 116) Workspace Organization Mandate

Regeln:
1. Keine unstrukturierten Loose-Files in Arbeitsroot
2. klare Trennung von aktiv/archiv/experiment
3. zentrale Config-Pfade stabil halten

---

## 117) Million-Line Ambition (Strategic)

Zielbild:
1. modulare Skalierbarkeit
2. testbare Architektur
3. Dokumentation auf Enterprise-Niveau

Leitgrößen (orientierend):
1. hoher Testabdeckungsgrad
2. ausgebautes API-Modulmodell
3. CI/CD + Monitoring by default

---

## 118) AI Screenshot Sovereignty (Generalized)

Regel:
- Screenshot-Artefakte strukturiert und zentral ablegen.

Pflicht:
1. Tool-spezifische Unterordner
2. Aufbewahrungs- und Archivierungsregeln
3. keine Desktop-Vermüllung in Produktionsumgebungen

---

## 119) Status Footer Protocol (Extended)

Jede größere Änderungsmeldung enthält:
1. aktualisierte Artefakte
2. Fortschritt (Code/Test/Deploy)
3. offene Risiken

Template:
```text
STATUS UPDATE
Updated:
- AGENTS-GLOBAL.md
- relevante BIOMETRICS-Dateien
Fortschritt:
- Code: x%
- Verifikation: x%
- Deploy: x%
```

---

## 120) Plan-Sovereignty Addendum

Vor neuem Plan:
1. vorhandene Pläne prüfen
2. Überschneidung bewerten
3. archivieren/erweitern statt duplizieren

Plan-Limits:
1. aktive Pläne begrenzen
2. alte Pläne regelmäßig archivieren

---

## 121) NLM Notebook Creation Policy

Bei neuem Themencluster:
1. Notebook erstellen
2. zentrale Quellen hinzufügen
3. erste Leitfrage ausführen
4. Notebook-ID in Doku erfassen

---

## 122) NLM Source-Type Policy

Pflichtquellen je Notebook:
1. Projektdokus (`.md`)
2. relevante Codebereiche
3. Konfigurationsdateien
4. offizielle URLs
5. ggf. Videoquellen

Regel:
- Dateivolumen ist weniger kritisch als Dateichurn und Dubletten.

---

## 123) NLM Empty-Answer Handling

Wenn Antwort leer:
1. Source-Liste prüfen
2. fehlende Quellen ergänzen
3. Query umformulieren
4. erneut validieren

Verboten:
- Leere Antwort stillschweigend ignorieren

---

## 124) NLM Migration Safety Protocol

Richtiger Ablauf:
1. neues Notebook erstellen
2. alle Quellen migrieren
3. Validierung alt vs. neu
4. erst dann altes Notebook bereinigen

---

## 125) Notebook Naming & ID Governance

Format:
`[ID] DEV-[Projekt]-[Typ]-[Bezeichnung]-[Kategorie]`

Pflicht:
1. IDs konsistent pflegen
2. Kategorien einheitlich anwenden
3. Projektfamilien klar trennen

---

## 126) Emoji Consistency Rule (NLM)

Regel:
- Notebooks einer Projektfamilie verwenden ein konsistentes Emoji/Labeling, um visuelle Zuordnung zu vereinfachen.

---

## 127) Feature-Notebook Backfill Mandate

Für bestehende Features:
1. Feature-Inventar erstellen
2. fehlende Feature-Notebooks anlegen
3. Kernquellen pro Feature hinzufügen

---

## 128) Page-Notebook Backfill Mandate

Für bestehende Seiten:
1. Seiteninventar erstellen
2. page-spezifische Notebooks nachziehen
3. Seitenquellen syncen

---

## 129) Sub-Agent NLM Delegation Contract

Jeder Delegationsprompt enthält:
1. Notebook-ID
2. Start-Query
3. verbotene Alternativpfade
4. Sync-Anweisung nach Abschluss

Beispiel:
```text
NOTEBOOK_ID:
FIRST_QUERY:
FORBIDDEN_TOOLS:
SYNC_REQUIRED:
```

---

## 130) Security Handbook Query Pattern

Vor Abschluss sicherheitskritischer Features:
```bash
nlm query notebook <security-notebook-id> "Prüfe <feature> gegen OWASP Top 10 und bekannte CVE-Patterns"
```

---

## 131) Architecture Validation Pattern

Vor großen Umbauten:
1. aktuelle Architekturpfade erfassen
2. Zielarchitektur dokumentieren
3. Migrationsschritte stufenweise ausführen
4. Rückfallpfad definieren

---

## 132) Integration Verification Matrix

| Layer | Prüfung | Evidenz |
|---|---|---|
| Frontend | Build + Route Smoke | Build Logs |
| Backend | API Contract + Error Paths | Request/Response |
| DB | Schema + RLS + Migration | SQL/Migration Logs |
| Infra | DNS/Tunnel/Runtime | Health Checks |

---

## 133) Business/Market Review Cadence

Rhythmus:
1. monatlicher Quick-Check
2. quartalsweiser Deep-Review
3. Major-Release Abgleich

---

## 134) Repository Hygiene Mandate

Pflicht:
1. klare Root-Struktur
2. keine redundanten Governance-Dateien außerhalb BIOMETRICS
3. veraltete Referenzen markieren und migrieren

---

## 135) Change Communication Standard

Bei jeder signifikanten Änderung:
1. Was wurde geändert?
2. Warum?
3. Wie validiert?
4. Welche Risiken offen?
5. Was ist nächster Schritt?

---

## 136) Executive Risk Register Template

```markdown
## Risk ID
Severity:
Likelihood:
Impact:
Mitigation:
Owner:
Status:
```

---

## 137) Quality Gate Escalation

Wenn ein Quality Gate fehlschlägt:
1. Taskstatus auf BLOCKED
2. Root Cause dokumentieren
3. Fixplan definieren
4. Re-Validation erzwingen

---

## 138) Compliance Snapshot Template

```text
COMPLIANCE SNAPSHOT
- Rules Applied:
- Tests Passed:
- Docs Synced:
- Security Checked:
- Residual Risks:
```

---

## 139) Annex Sequencing Policy

Ausbauprinzip:
1. Kernregeln zuerst
2. operationale Templates danach
3. projektnahe Playbooks zuletzt

Regel:
- Jeder neue Annexblock erhöht Präzision und Umsetzbarkeit, nicht nur Umfang.

---

## 140) Expansion Gate (to 2.5k+)

Status:
1. Annex bis `110` vorhanden
2. Annex `111–140` vorhanden

Nächster Schritt (`141+`):
1. zusätzliche verbatim Tabellen
2. tiefere provider- und notebook-spezifische Beispiele
3. erweiterte incident und recovery playbooks

Leitsatz:
**No rule left behind.**

---

## 141) Provider Validation Playbook

Bei Provider-Änderungen zwingend:
1. Modell-ID prüfen
2. Endpoint prüfen
3. Auth-Mechanismus prüfen
4. Limits dokumentieren
5. Smoke-Test durchführen

Validierungsprotokoll:
```text
Provider:
Model:
Endpoint:
Auth:
Limits:
Smoke-Test Result:
```

---

## 142) NVIDIA 429 Backoff Protocol

Bei `429 Too Many Requests`:
1. Backoff (mind. 60s)
2. Retry mit begrenzter Wiederholung
3. Fallback-Kette nutzen
4. Incident-Notiz erstellen, falls wiederholt

---

## 143) Model Fallback Contract

Jede Fallback-Stufe dokumentiert:
1. Auslöser
2. Zielmodell
3. Qualitätsrisiko
4. Rückkehrbedingung auf Primärmodell

---

## 144) OpenClaw Gateway Ops Checklist

- [ ] Provider geladen
- [ ] Modelle sichtbar
- [ ] Health ok
- [ ] Fallbacks aktiv
- [ ] Fehlermonitoring aktiv

---

## 145) OpenCode Session Hygiene

Regeln:
1. Session-Kontext klar halten
2. große Schritte in Teilaufgaben trennen
3. Ergebnisse mit Evidenz rückmelden

---

## 146) MCP Availability Decision Tree

Wenn MCP verfügbar:
1. nutzen
2. Ergebnis protokollieren

Wenn MCP nicht verfügbar:
1. Blocker dokumentieren
2. Fallback-Strategie anwenden
3. Nacharbeiten einplanen

---

## 147) Serena-Orchestrierung Standard

Pflicht:
1. Aufgaben in Arbeitspakete teilen
2. Verantwortlichkeiten klar zuweisen
3. Ergebnisse konsolidieren
4. Konflikte früh auflösen

---

## 148) CDP Browser Verification Playbook

Bei UI-kritischen Änderungen:
1. Seite laden
2. Kernflow durchspielen
3. Console prüfen
4. Sichtprüfung dokumentieren

---

## 149) Tavily/Research Complement Rule

Regel:
- Ergänzende Recherchequellen dürfen genutzt werden, wenn NLM-Context nicht ausreichend ist.
- Ergebnisse müssen in NLM/Wissensspeicher zurückgeführt werden.

---

## 150) Governance-First Edit Rule

Vor jeder Implementierung:
1. Governance-Auswirkung prüfen
2. betroffene Doku-Dateien identifizieren
3. Mapping-Pflicht festlegen

---

## 151) Branch Naming Extended Convention

Format:
1. `feat/<ticket>-<name>`
2. `fix/<ticket>-<name>`
3. `hotfix/<version>-<name>`
4. `docs/<name>`

---

## 152) Commit Message Quality Gate

Ein Commit muss enthalten:
1. Typ
2. Scope
3. klare Intention

Verboten:
1. `update`
2. `misc`
3. nichtssagende Messages

---

## 153) PR Description Blueprint

```markdown
## What
## Why
## How
## Validation
## Risks
## Rollback
## Docs Updated
```

---

## 154) CI Failure Handling

Bei CI-Fehler:
1. Ursache klassifizieren
2. minimalen Fix erstellen
3. Re-Run ausführen
4. Ergebnis protokollieren

---

## 155) Dependency Hygiene Mandate

Pflicht:
1. veraltete/verwundbare Abhängigkeiten identifizieren
2. Upgrade-Plan dokumentieren
3. Breaking-Risiken evaluieren

---

## 156) API Breaking Change Protocol

Wenn Breaking Change unvermeidbar:
1. Version bump
2. Migration Guide
3. Deprecation-Zeitraum
4. Backward-Compatibility Hinweis

---

## 157) Database Migration Safety

Pflichtschritte:
1. Pre-Migration Check
2. Migration dry-run
3. Backup-Verifikation
4. Post-Migration Smoke
5. Rollback-Pfad dokumentiert

---

## 158) RLS Verification Protocol (Supabase)

Für jede Policy:
1. Positivtest (erlaubter Zugriff)
2. Negativtest (unerlaubter Zugriff)
3. Service-role Ausnahmen prüfen

---

## 159) Auth Flow Integrity Check

Pflichtfälle:
1. Login
2. Logout
3. Token Refresh
4. Role-based Access
5. Session Expiry

---

## 160) Storage & Asset Security

Pflicht:
1. Upload-Validierung
2. MIME/Size Constraints
3. Zugriffspolicies
4. öffentliche vs private Pfade trennen

---

## 161) Incident Severity Model

Level:
1. SEV-1: kritisch / Ausfall
2. SEV-2: hoch / starke Einschränkung
3. SEV-3: mittel / partiell
4. SEV-4: niedrig / kosmetisch

Pflicht:
- Severity in jedem Incident-Ticket angeben

---

## 162) Postmortem Template

```markdown
## Summary
## Impact
## Timeline
## Root Cause
## Corrective Actions
## Preventive Actions
## Owners
```

---

## 163) Observability Baseline

Pflichtsignalarten:
1. Logs
2. Errors
3. Metrics
4. Traces

Pflicht:
- kritische Pfade müssen mindestens logs + errors haben

---

## 164) Performance Budget Contract

Definiere pro Projekt:
1. Frontend Zielwerte
2. API Latenzbudgets
3. DB Query-Budgets
4. Alarm-Schwellen

---

## 165) Accessibility Gate

Pflicht bei UI:
1. Tastaturbedienbarkeit
2. semantische Struktur
3. Kontrast
4. Fokuszustände

---

## 166) Internationalization Baseline

Wenn i18n erforderlich:
1. String-Externalisierung
2. Fallback-Sprache
3. Locale-Routing dokumentieren

---

## 167) UX Regression Check

Bei UI-Änderung:
1. Kernpfad prüfen
2. Fehlerzustände prüfen
3. Empty States prüfen
4. Mobile/Responsive prüfen

---

## 168) Legal/Compliance Reminder

Pflicht:
1. datenschutzrelevante Änderungen markieren
2. regulatorische Auswirkungen dokumentieren
3. auditierbare Nachweise bereitstellen

---

## 169) Patent/IP Documentation Hook

Wenn IP-relevante Änderungen auftreten:
1. technische Neuerung erfassen
2. Kontext/Abgrenzung dokumentieren
3. Verweis auf Patent-Doku ergänzen

---

## 170) Business Impact Tagging

Jede größere Änderung bekommt Impact-Tag:
1. Revenue
2. Cost
3. Risk
4. Velocity
5. Quality

---

## 171) Architecture Decision Record (ADR) Hook

Bei architekturrelevanten Änderungen:
1. Entscheidung
2. Alternativen
3. Konsequenzen
4. Verweise

---

## 172) Rollback Readiness Standard

Für kritische Deployments:
1. Rollback-Befehl
2. Rollback-Zeitfenster
3. Datenkonsistenzstrategie
4. Kommunikationsplan

---

## 173) Staging/Production Drift Check

Vor Release:
1. Config-Diff prüfen
2. Env-Variablen prüfen
3. Feature-Flags prüfen

---

## 174) Feature Flag Governance

Regeln:
1. Flag-Name dokumentieren
2. Rollout-Strategie
3. Kill-Switch
4. Cleanup-Datum

---

## 175) Release Notes Template

```markdown
## Highlights
## Fixes
## Breaking Changes
## Migration Steps
## Known Issues
```

---

## 176) User-Plan / Agent-Plan Separation

Regel:
1. Agent Aufgaben nur in AGENTS-PLAN
2. User Aufgaben nur in USER-PLAN
3. keine Vermischung

---

## 177) Meeting Log Protocol

Jeder relevante Schritt im Meetinglog enthält:
1. Entscheidung
2. Begründung
3. Risiko
4. Nächste Schritte

---

## 178) Changelog Quality Standard

Changelog-Eintrag enthält:
1. was geändert wurde
2. warum
3. impact
4. referenzierte Dateien

---

## 179) Rule Consistency Audit

Regelmäßige Audits prüfen:
1. Widersprüche zwischen global/lokal
2. veraltete Regeln
3. fehlende Mappings
4. unvollständige Templates

---

## 180) Expansion Gate (to 3k+)

Status:
1. Annex `141–180` vorhanden
2. Datei > 2k Zeilen

Nächster Schritt (`181+`):
1. detaillierte Verbatim-Tabellen aus Provider/NLM/GitHub/Infra
2. vollständige Playbook-Sets mit Entscheidungsbäumen
3. projektspezifische Override-Module

Leitsatz:
**Enterprise rigor through additive governance.**

---

## 181) Runtime Readiness Checklist

- [ ] Environment Variablen vollständig
- [ ] Secrets sicher geladen
- [ ] externe Services erreichbar
- [ ] Health-Endpunkte grün
- [ ] kritische Background-Jobs aktiv

---

## 182) Configuration Drift Control

Pflicht:
1. Soll-Konfiguration dokumentieren
2. Ist-Konfiguration regelmäßig vergleichen
3. Drift als Ticket erfassen

---

## 183) Environment Variable Governance

Regeln:
1. `.env.example` aktuell halten
2. Sensible Variablen nie committen
3. Required/Optional klar markieren

Template:
```text
VAR_NAME=
SCOPE=server|client
REQUIRED=true|false
DEFAULT=
DESCRIPTION=
```

---

## 184) Secret Rotation Cadence

Pflicht:
1. Rotationsintervall definieren
2. Ablauf und Owner dokumentieren
3. Rückrollplan für fehlerhafte Rotation

---

## 185) Service Ownership Matrix

| Service | Owner | Backup Owner | SLA | Escalation |
|---|---|---|---|---|
| API | | | | |
| DB | | | | |
| Auth | | | | |
| Queue | | | | |

---

## 186) API Contract Snapshot Template

```markdown
## Endpoint
Method:
Auth:
Input Schema:
Output Schema:
Error Codes:
Rate Limit:
```

---

## 187) Endpoint Error Taxonomy

Pflichtklassen:
1. Validation
2. Authorization
3. Not Found
4. Conflict
5. Rate Limited
6. Internal

---

## 188) Retry & Idempotency Rule

Regeln:
1. Retries nur für sichere Operationen
2. Idempotency-Key bei kritischen Writes
3. Duplicate Side-Effects verhindern

---

## 189) Queue Reliability Standard

Pflicht:
1. Dead-Letter-Queue definieren
2. Retry-Limits setzen
3. Poison-Messages erkennen
4. Monitoring aktivieren

---

## 190) Cron/Job Governance

Für jeden Job:
1. Zweck
2. Zeitplan
3. Inputs/Outputs
4. Failure-Handling
5. Owner

---

## 191) Feature Rollout Strategy

Rollout-Stufen:
1. internal
2. beta
3. general

Pflicht:
1. Erfolgsmessung
2. Stop-Kriterien
3. Rollback-Trigger

---

## 192) Canary Release Protocol

1. kleine Zielgruppe
2. Fehler-/Latenz-Monitoring
3. Entscheidung expand/rollback

---

## 193) Data Integrity Guardrails

Pflicht:
1. Constraints dokumentieren
2. Migrationen validieren
3. Konsistenzchecks nach Änderungen

---

## 194) Backup Verification Mandate

Regel:
1. Backups nicht nur erzeugen, sondern wiederherstellbar testen.
2. Restore-Test protokollieren.

---

## 195) Restore Drill Template

```text
System:
Backup Time:
Restore Start:
Restore End:
Validation:
Result:
```

---

## 196) Data Privacy Classification

Kategorien:
1. Public
2. Internal
3. Sensitive
4. Restricted

Pflicht:
- jede Datendomäne klassifizieren

---

## 197) PII Handling Rules

1. Minimalprinzip
2. Zugriff nur zweckgebunden
3. Logging ohne PII-Leakage
4. Lösch-/Auskunftsprozesse dokumentieren

---

## 198) Audit Log Minimum Fields

Pflichtfelder:
1. Actor
2. Action
3. Target
4. Timestamp
5. Result
6. Correlation-ID

---

## 199) Correlation-ID Propagation

Regel:
- Correlation-ID über API, Queue und Worker durchreichen.

Pflicht:
- in Fehler- und Performanceanalyse nutzbar machen.

---

## 200) SLO/SLA Definition Template

```markdown
## Service
SLO:
SLA:
Error Budget:
Measurement Window:
Alert Threshold:
```

---

## 201) Alert Routing Matrix

| Severity | Channel | Response Time | Owner |
|---|---|---|---|
| SEV-1 | | | |
| SEV-2 | | | |
| SEV-3 | | | |

---

## 202) Runbook Standard

Jeder kritische Dienst braucht Runbook mit:
1. Symptom
2. Diagnose
3. Sofortmaßnahme
4. Dauerfix
5. Eskalation

---

## 203) Error Budget Policy

Wenn Error Budget aufgebraucht:
1. Feature-Rollout drosseln
2. Stabilität priorisieren
3. Ursachenliste priorisiert abarbeiten

---

## 204) Performance Profiling Routine

1. Baseline messen
2. Hotspots identifizieren
3. Optimierung umsetzen
4. Nachmessung dokumentieren

---

## 205) Cost Awareness Rule

Pflicht:
1. teure Pfade markieren
2. Skalierungskosten abschätzen
3. Optimierungsoptionen dokumentieren

---

## 206) Technical Debt Register

Template:
```markdown
## Debt Item
Type:
Impact:
Risk:
Fix Proposal:
Priority:
Owner:
```

---

## 207) Refactoring Gate

Refactoring nur, wenn:
1. Zielbild klar ist
2. Risikoabschätzung vorhanden ist
3. Regressionstests vorgesehen sind

---

## 208) Module Boundary Rule

Pflicht:
1. klare Verantwortlichkeiten je Modul
2. keine zyklischen Abhängigkeiten
3. Schnittstellen explizit dokumentieren

---

## 209) Dependency Direction Check

Regel:
- High-level Module dürfen nicht von Low-level Details abhängig werden.

---

## 210) API Deprecation Lifecycle

Phasen:
1. announce
2. deprecate
3. migrate
4. remove

Pflicht:
- jede Phase mit Datum und Kommunikation dokumentieren

---

## 211) Schema Evolution Policy

1. additive Änderungen bevorzugen
2. Breaking Schema-Änderungen versionieren
3. Datenmigrationen testbar machen

---

## 212) Testing Pyramid Enforcement

Mindestbalance:
1. viele Unit Tests
2. gezielte Integrationstests
3. kritische E2E Flows

---

## 213) Flaky Test Handling

Regel:
1. Flaky Tests markieren
2. Ursache priorisiert beheben
3. nicht dauerhaft ignorieren

---

## 214) Visual Regression Policy

Bei UI-kritischen Komponenten:
1. Baseline erstellen
2. Vergleichslauf
3. Abweichungen bewerten

---

## 215) Documentation Freshness SLA

Pflicht:
1. Doku-Updates in derselben Änderung wie Code
2. veraltete Abschnitte mit TODO+Owner markieren

---

## 216) Glossary Governance

Regel:
1. zentrale Begriffsdefinitionen pflegen
2. Synonyme harmonisieren
3. Fachwörter konsistent verwenden

---

## 217) Onboarding Speed Standard

Ziel:
1. neues Teammitglied in <1 Tag produktiv

Pflicht:
1. Quickstart
2. lokale Setup-Schritte
3. typische Fehler und Lösungen

---

## 218) Release Readiness Checklist

- [ ] CI grün
- [ ] Security Check bestanden
- [ ] Migrationen validiert
- [ ] Rollback dokumentiert
- [ ] Release Notes erstellt

---

## 219) Executive Handover Template

```markdown
## Current State
## Completed Work
## Open Risks
## Immediate Next Steps
## Required Decisions
```

---

## 220) Expansion Gate (to 3.5k+)

Status:
1. Annex `181–220` vorhanden
2. Datei > 2.7k Zeilen

Nächster Schritt (`221+`):
1. weitere verbatim-nahe Detailtabellen
2. Provider- und MCP-spezifische Runbook-Packs
3. projektübergreifende Compliance Blueprints

Leitsatz:
**Document deep, execute clean, verify hard.**

---

## 221) NVIDIA NIM VERBATIM PACK — Core Facts

Provider-Basis:
1. Endpoint: `https://integrate.api.nvidia.com/v1`
2. API-Kompatibilität: OpenAI-compatible provider bridge
3. Modellrouting strikt über korrekte Model IDs

Empfohlene Modellbeispiele:
1. `qwen/qwen3.5-397b-a17b`
2. `moonshotai/kimi-k2.5`

---

## 222) NVIDIA Timeout & Latency Guardrail

Regel:
1. Bei High-Latency-Modellen Requests mit ausreichendem Timeout fahren.
2. Timeouts nicht künstlich zu niedrig setzen.
3. Abbrüche als Betriebsrisiko dokumentieren.

Richtwert:
- OpenCode-Pfade: mindestens `120000ms` (wenn technisch unterstützt).

---

## 223) NVIDIA Validation Command Pack

```bash
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
  https://integrate.api.nvidia.com/v1/models
```

```bash
openclaw models
opencode models
```

Ziel:
1. Provider erreichbar
2. Modell sichtbar
3. Auth funktionsfähig

---

## 224) OpenCode Provider JSON Blueprint (Verbatim-Nah)

```json
{
  "provider": {
    "nvidia-nim": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "NVIDIA NIM",
      "options": {
        "baseURL": "https://integrate.api.nvidia.com/v1",
        "timeout": 120000
      },
      "models": {
        "qwen-3.5-397b": {
          "id": "qwen/qwen3.5-397b-a17b",
          "limit": { "context": 262144, "output": 32768 }
        }
      }
    }
  }
}
```

---

## 225) OpenClaw Provider JSON Blueprint (Verbatim-Nah)

```json
{
  "models": {
    "providers": {
      "nvidia": {
        "baseUrl": "https://integrate.api.nvidia.com/v1",
        "api": "openai-completions",
        "models": ["qwen/qwen3.5-397b-a17b"]
      }
    }
  }
}
```

Hinweis:
1. Unterschiede zwischen OpenCode und OpenClaw explizit dokumentieren.
2. Streaming-/Timeout-Besonderheiten getrennt erfassen.

---

## 226) Provider Error Patterns & Fixes

Typische Fehler:
1. falsche Modell-ID
2. ungültiger Endpoint
3. fehlende API Keys
4. 429 ohne Backoff

Fixpfad:
1. Modell-ID prüfen
2. Endpoint prüfen
3. Auth prüfen
4. Backoff/Fallback anwenden

---

## 227) NLM VERBATIM PACK — Core CLI

```bash
nlm list notebooks
nlm notebook create "<name>"
nlm source list <notebook-id>
nlm source add <notebook-id> --file "<file>" --wait
nlm source add <notebook-id> --url "<url>" --wait
nlm query notebook <notebook-id> "<frage>"
nlm status source <source-id>
```

---

## 228) NLM Duplicate Prevention (Strict)

Pflichtablauf:
```bash
nlm source list <notebook-id>
nlm source delete <duplicate-source-id> -y
nlm source add <notebook-id> --file "<file>" --wait
```

Verboten:
1. Blindes mehrfaches Hochladen derselben Quelle

---

## 229) NLM Sync-on-Change Mandate

Bei relevanten Dateiänderungen:
```bash
nlm source add <notebook-id> --file "<changed-file>" --wait
```

Alternativ:
```bash
nlm source sync <notebook-id> --confirm
```

---

## 230) NLM Deep Research Mandate

```bash
nlm research start "<topic>" --mode deep --notebook-id <id>
```

Regel:
1. Komplexe Themen standardmäßig in deep mode bearbeiten.

---

## 231) NLM Crash-Test Prompt Pack

Vor kritischer Arbeit:
```bash
nlm query notebook <id> "Was ist der aktuelle Stand von <projekt/thema>? Liste Architektur, Risiken und offene Punkte auf."
```

Wenn Antwort unzureichend:
1. Quellen ergänzen
2. Frage präzisieren
3. erneut validieren

---

## 232) NLM Notebook Selection Matrix (Expanded)

| Tasktyp | Ziel-Notebook |
|---|---|
| Coding | Development |
| Hardware/Sensorik | Hardware |
| Externe APIs/Docs | Research |
| Patent/Legal | Patent |
| Business/Investor | Marketing/Business |
| Strategie | Framework/Vision |

---

## 233) NLM Naming Convention Pack

Format:
`[ID] DEV-[Projekt]-[Typ]-[Bezeichnung]-[Kategorie]`

Kategorien:
1. Feature
2. Page
3. Bug
4. Wiki
5. Docs
6. Patent
7. Vision

---

## 234) NLM Migration Zero-Loss Protocol

1. Neues Notebook anlegen
2. Alle Quellen übernehmen
3. Quellenanzahl alt/neu vergleichen
4. Crash-Test im neuen Notebook
5. erst dann altes Notebook bereinigen

---

## 235) Sub-Agent NLM Delegation Template (Verbatim-Nah)

```text
TASK:
NOTEBOOK_ID:
FIRST_QUERY:
CONTEXT_REQUIREMENTS:
FORBIDDEN_TOOLS:
SYNC_AFTER_DONE:
OUTPUT_FORMAT:
```

Regel:
1. Subagenten ohne Notebook-Kontext sind unzulässig bei Wissensaufgaben.

---

## 236) GitHub Issue Template Pack (Blueprint)

Bug Template Mindestfelder:
1. Beschreibung
2. Reproduktion
3. Erwartet vs Ist
4. Environment
5. Logs

Feature Template Mindestfelder:
1. Problem
2. Lösung
3. Alternativen
4. Acceptance Criteria

---

## 237) GitHub PR Template Pack (Blueprint)

Pflichtsektionen:
1. Description
2. Related Issue
3. Type of Change
4. Checklist
5. Testing Instructions

---

## 238) GitHub CI Workflow Baseline (Verbatim-Nah)

Pflichtjobs:
1. lint
2. typecheck
3. test
4. build

Optionale Jobs:
1. codeql
2. release
3. dependency automation

---

## 239) Dependabot Strategy Baseline

Pflicht:
1. regelmäßige Dependency-Updates
2. Limitierung offener PRs
3. Labels und Commit-Präfixe

---

## 240) CODEOWNERS Governance

Pflicht:
1. Kernbereiche mit Ownern mappen
2. Infra/Doku explizit zuordnen
3. Review-Pfade nachvollziehbar halten

---

## 241) Branch Protection Rules Pack

Main branch:
1. PR-Review erforderlich
2. Statuschecks erforderlich
3. Force Push disabled
4. Delete disabled

---

## 242) CONTRIBUTING Standard Pack

Pflichtinhalte:
1. Setup
2. Commit-Standard
3. PR-Prozess
4. Tests
5. Kontaktwege

---

## 243) CODE_OF_CONDUCT / SECURITY Policy Hook

Pflicht:
1. Community-Standards vorhanden
2. Sicherheits-Meldeweg dokumentiert

---

## 244) Compliance-to-Docs Traceability

Jede Governance-Regel muss referenzierbar sein in:
1. Architektur
2. Betrieb
3. Sicherheitsdoku
4. Change-Historie

---

## 245) Risk-to-Task Traceability

Pflicht:
1. Jeder High-Risk Punkt erhält Task-ID
2. Task-Status im Risk-Register rückgespiegelt

---

## 246) Decision-to-ADR Traceability

Pflicht:
1. Wichtige Architekturentscheidungen als ADR
2. ADR-Referenzen in Umsetzungstickets

---

## 247) Verification Evidence Pack

Mindestens eine Evidenzart je Task:
1. Testoutput
2. Buildoutput
3. API-Beispiel
4. Screenshot/Flow-Check
5. Migrationsnachweis

---

## 248) No-Mock Production Rule

Regel:
1. Keine pseudo-fertigen Demozustände als Endlieferung.
2. Prototyping klar kennzeichnen.
3. Produktionspfad real implementieren.

---

## 249) Executive Review Gate

Vor Abschluss großer Blöcke:
1. Zielerreichung
2. Qualität
3. Sicherheit
4. Dokumentation
5. offene Risiken

---

## 250) Delivery Readiness Matrix

| Bereich | Status | Evidenz |
|---|---|---|
| Funktionalität | | |
| Qualität | | |
| Sicherheit | | |
| Dokumentation | | |
| Betrieb | | |

---

## 251) Team Communication Policy

Regel:
1. Entscheidungen transparent kommunizieren
2. Blocker früh sichtbar machen
3. Handovers strukturiert liefern

---

## 252) Knowledge Handover Bundle

Pflichtinhalte je Handover:
1. aktueller Stand
2. was getestet wurde
3. was offen ist
4. Risiken
5. nächste 3 Schritte

---

## 253) Recovery Playbook Skeleton

```text
Failure Type:
Detection:
Immediate Mitigation:
Root Cause Track:
Long-Term Fix:
Verification:
```

---

## 254) Security Hardening Baseline

1. least privilege
2. secure defaults
3. dependency checks
4. auth boundary tests
5. secret scope minimization

---

## 255) Data Governance Baseline

1. Datenherkunft dokumentieren
2. Speicherorte mappen
3. Aufbewahrungsregeln definieren
4. Löschprozesse dokumentieren

---

## 256) Operations Excellence Checklist

- [ ] Healthchecks definiert
- [ ] Runbooks vorhanden
- [ ] Alerts geroutet
- [ ] Backups validiert
- [ ] Rollback getestet

---

## 257) Governance Conflict Resolution Pack

Bei Regelkonflikt:
1. Konflikt benennen
2. Priorität anwenden
3. Entscheidung dokumentieren
4. betroffene Doku synchronisieren

---

## 258) Additive Change Control

Regeln:
1. Änderungen als Erweiterungen
2. keine stillen Löschungen
3. Migrationsnotizen bei Strukturänderung

---

## 259) Annex Integrity Check

Vor jedem großen Update:
1. Rule-IDs konsistent?
2. Referenzen intakt?
3. Widersprüche markiert?
4. Expansion Gate aktualisiert?

---

## 260) Expansion Gate (to 4k+)

Status:
1. Annex `221–260` vorhanden
2. Datei > 3.1k Zeilen

Nächster Schritt (`261+`):
1. Deep Verbatim Tables (Provider/NLM/GitHub)
2. Incident/Recovery Trees je Domäne
3. Enterprise Audit Packs

Leitsatz:
**No governance debt.**

---

## 261) Incident Decision Tree — Entry Point

```text
Incident erkannt
  ├─ SEV-1? → Sofort-Eskalation + Mitigation
  ├─ SEV-2? → Priorisierte Eindämmung + Root-Cause-Track
  ├─ SEV-3? → Geplante Behebung + Monitoring
  └─ SEV-4? → Backlog + Beobachtung
```

---

## 262) SEV-1 Decision Tree

```text
SEV-1
  ├─ Service down? → Failover/Rollback prüfen
  ├─ Datenverlust-Risiko? → Write-Freeze + Backup-Check
  ├─ Security betroffen? → Incident Security Protocol aktivieren
  └─ Kommunikation → Stakeholder sofort informieren
```

---

## 263) SEV-2 Decision Tree

```text
SEV-2
  ├─ Teilfunktion defekt → Feature-Flag / partial disable
  ├─ Performance-Einbruch → Profiling + Rate-Limit Guard
  └─ Workaround bereitstellen + Ticket priorisieren
```

---

## 264) SEV-3/SEV-4 Decision Tree

```text
SEV-3/4
  ├─ reproduzierbar? ja → Fix-Task planen
  ├─ reproduzierbar? nein → Monitoring + Observability erhöhen
  └─ Roadmap-Einordnung + Risiko dokumentieren
```

---

## 265) Release Go/No-Go Matrix

| Kriterium | GO | NO-GO |
|---|---|---|
| CI | grün | rot |
| Security | keine kritischen Findings | offene kritische Findings |
| Migration | validiert | ungetestet |
| Rollback | dokumentiert + getestet | unklar |
| Monitoring | Alerts aktiv | Blindflug |

---

## 266) Rollback Decision Tree

```text
Release Problem?
  ├─ Kundenimpact hoch → sofort Rollback
  ├─ Dateninkonsistenz → Write-Stop + Recovery
  ├─ Hotfix <30min möglich? → Hotfix-Track
  └─ sonst Rollback + Postmortem
```

---

## 267) Data Loss Prevention Matrix

| Risiko | Prävention | Detection | Recovery |
|---|---|---|---|
| Fehlmigration | Dry-run + Backup | Schema-Checks | Restore Drill |
| Falsche Delete-Operation | Soft-delete/Guardrails | Audit Logs | PITR/Restore |
| Queue Duplication | Idempotency Keys | Duplicate Metrics | Reconcile Job |

---

## 268) Security Breach First 30 Minutes

1. Incident deklarieren
2. Zugriff begrenzen
3. kompromittierte Schlüssel rotieren
4. Evidence sichern
5. Scope erfassen
6. Kommunikationskanal fixieren

---

## 269) Security Containment Checklist

- [ ] betroffene Systeme isoliert
- [ ] verdächtige Tokens deaktiviert
- [ ] kritische Endpunkte geschützt
- [ ] Forensik-Daten gesichert
- [ ] Timeline gestartet

---

## 270) Vulnerability Triage Model

Priorität:
1. Exploitability
2. Impact
3. Exposure
4. Remediation Cost

| Priority | Kriterium | ETA |
|---|---|---|
| P0 | aktiv ausnutzbar + hoher Impact | sofort |
| P1 | hoher Impact, kein aktiver Exploit | 24-72h |
| P2 | moderat | geplanter Sprint |

---

## 271) API Abuse Protection Tree

```text
Spike erkannt?
  ├─ legitimer Traffic? → Capacity Scale
  ├─ verdächtig? → Rate Limit + WAF Rule
  └─ bot pattern? → Challenge/Block + Alert
```

---

## 272) Auth Incident Tree

```text
Auth Fehler?
  ├─ Login-Ausfall → Provider/Session Store prüfen
  ├─ Unauthorized Access → Token Scope/Leak prüfen
  ├─ Role Escalation → RBAC Regeln auditieren
  └─ Refresh Loop → Token TTL/Clock Skew prüfen
```

---

## 273) Database Incident Tree

```text
DB Problem?
  ├─ Connection Exhaustion → Pooling/Leak prüfen
  ├─ Slow Queries → EXPLAIN + Index Check
  ├─ Lock Contention → Transaction Scope prüfen
  └─ Replication Lag → read/write routing anpassen
```

---

## 274) Queue Incident Tree

```text
Queue Backlog?
  ├─ Consumer down → restart + health
  ├─ poison message → DLQ route
  ├─ retry storm → retry backoff reduzieren
  └─ throughput limit → parallel consumer scale
```

---

## 275) Observability Escalation Tree

```text
Noisy Alerts?
  ├─ false positives hoch → thresholds kalibrieren
  ├─ missing alerts → neue Regeln ergänzen
  ├─ no traces → tracing instrumentation aktivieren
  └─ poor logs → structured logging erzwingen
```

---

## 276) Performance Triage Table

| Symptom | Hypothese | Test | Fix |
|---|---|---|---|
| Hohe API Latenz | DB bottleneck | query profiling | index/query rewrite |
| Hohe TTFB | cold start | warmup checks | runtime tuning |
| UI jank | bundle groß | analyzer | split/lazy load |

---

## 277) Frontend Perf Guardrails (Next.js)

1. Server Components bevorzugen
2. client bundles minimieren
3. Bilder optimieren
4. kritische Routen messen
5. caching bewusst einsetzen

---

## 278) Backend Perf Guardrails (Go)

1. Context propagation
2. timeouts per call
3. retries mit jitter
4. pooling und backpressure
5. pprof/trace für Hotspots

---

## 279) Supabase Ops Guardrails

1. RLS standardmäßig aktiv
2. service-role bewusst begrenzen
3. migrations versioniert
4. storage policies auditiert
5. edge functions observability aktiv

---

## 280) NLM Operations Decision Tree

```text
Neue Wissensfrage?
  ├─ passendes Notebook vorhanden? ja → query
  ├─ nein → notebook erstellen + baseline sources
  ├─ Antwort leer → source gap schließen
  └─ Ergebnis relevant → sync + doku update
```

---

## 281) NLM Quality Escalation

Wenn NLM Output unter Mindestscore:
1. Prompt präzisieren
2. Quellen erweitern
3. Strukturierte Re-Query
4. Ergebnis erneut bewerten

---

## 282) NLM Source Governance Table

| Source Type | Pflicht | Prüfkriterium |
|---|---|---|
| Markdown | ja | aktuell + vollständig |
| Code | ja | relevante Module enthalten |
| Config | ja | runtime-relevant |
| URL Docs | ja | offiziell/verifizierbar |
| Video | optional/kontextabhängig | inhaltlich relevant |

---

## 283) NLM Prompt Anti-Pattern List

Verboten:
1. vage Zieldefinition
2. fehlender Output-Contract
3. keine Qualitätskriterien
4. unklare Zielgruppe

Pflicht:
1. Rolle
2. Ziel
3. Constraints
4. erwartete Struktur

---

## 284) GitHub Workflow Decision Tree

```text
PR erstellt?
  ├─ CI grün? ja → Review
  ├─ CI rot? → Fix + rerun
  ├─ security findings? → block merge
  └─ docs missing? → update required
```

---

## 285) Review Quality Checklist

- [ ] Scope korrekt
- [ ] Risiken adressiert
- [ ] Tests sinnvoll
- [ ] Doku synchron
- [ ] Breaking changes markiert

---

## 286) Merge Safety Rules

1. Kein Merge bei roten Quality Gates
2. Kein Merge ohne klaren Rollback-Pfad bei kritischen Änderungen
3. Kein Merge mit offenen P0 Risiken

---

## 287) Documentation Coverage Map

| Bereich | Primärdatei | Sekundärdatei |
|---|---|---|
| Architektur | ARCHITECTURE.md | MAPPING-FRONTEND-BACKEND.md |
| API | ENDPOINTS.md | COMMANDS.md |
| DB | SUPABASE.md | MAPPING-DB-API.md |
| NLM | NOTEBOOKLM.md | MAPPING-NLM-ASSETS.md |

---

## 288) Architecture Drift Audit

1. Ist-Architektur erfassen
2. Soll-Architektur vergleichen
3. Drift priorisieren
4. Korrekturtasks anlegen

---

## 289) Compliance Evidence Register

Pflichtfelder:
1. Regel-ID
2. Nachweisart
3. Dateiverweis
4. Datum
5. Owner

---

## 290) Governance KPI Pack

Metriken:
1. % Tasks mit vollständiger Evidenz
2. % Doku-Sync bei Code-Changes
3. Incident MTTR
4. Change Failure Rate
5. NLM Source Hygiene Score

---

## 291) Reliability KPI Pack

1. Uptime
2. Error Rate
3. P95/P99 Latenz
4. Queue Backlog
5. Deployment Success Rate

---

## 292) Security KPI Pack

1. offene kritische Findings
2. Patch Lead Time
3. Secret Exposure Incidents
4. Auth Failure Trends

---

## 293) Delivery KPI Pack

1. Lead Time
2. Deploy Frequency
3. Review Cycle Time
4. Rework Rate

---

## 294) Risk Burn-Down Template

```markdown
## Week
High Risks Open:
Mitigated:
New Risks:
Trend:
```

---

## 295) Incident Communication Template

```text
INCIDENT UPDATE
Severity:
Impact:
Current Status:
Mitigation:
ETA:
Next Update:
```

---

## 296) Recovery Verification Checklist

- [ ] Funktion wiederhergestellt
- [ ] Datenkonsistenz geprüft
- [ ] Monitoring stabil
- [ ] Root Cause bestätigt
- [ ] Präventionsmaßnahmen geplant

---

## 297) Executive Escalation Trigger

Eskalation wenn:
1. SEV-1 aktiv
2. wiederholte SEV-2 innerhalb kurzer Zeit
3. Compliance-Risiko kritisch
4. Datenintegrität gefährdet

---

## 298) Change Window Policy

Regel:
1. kritische Changes in definierten Change Windows
2. High-Risk Deployments mit zusätzlicher Bereitschaft

---

## 299) Freeze Policy

Code Freeze auslösen bei:
1. Instabilität über Schwellwert
2. kritischen Security Findings
3. unkontrollierter Incident-Kaskade

---

## 300) Hotfix Policy

Pflicht:
1. minimaler Scope
2. klare Validierung
3. anschließender saubere Nachverfolgung in Hauptbranch

---

## 301) Release Train Governance

1. geplante Release-Cadence
2. klarer Cutoff
3. Nachzügler in nächste Iteration

---

## 302) Dependency Upgrade Playbook

1. Risiko bewerten
2. changelog review
3. staging test
4. production rollout stufenweise

---

## 303) Migration Rollout Strategy

1. precheck
2. dry-run
3. staged execution
4. postcheck
5. rollback readiness

---

## 304) Cross-Team Coordination Rule

Pflicht:
1. gemeinsame Schnittstellen abstimmen
2. Abhängigkeiten transparent machen
3. Release-Reihenfolge koordinieren

---

## 305) Contract Testing Rule

Bei serviceübergreifenden APIs:
1. Contract Tests definieren
2. Breaking Contract früh erkennen
3. Kompatibilität überwachen

---

## 306) Data Contract Governance

1. Schemas versionieren
2. Event/Message Contracts dokumentieren
3. Consumer-Impact bei Änderungen prüfen

---

## 307) Knowledge Debt Register

Template:
```markdown
## Knowledge Gap
Impact:
Affected Teams:
Resolution Plan:
Owner:
Target Date:
```

---

## 308) Training & Enablement Pack

Pflicht:
1. Onboarding-Artefakte aktualisieren
2. Lessons Learned einpflegen
3. wiederkehrende Fehler als Trainingsinhalt markieren

---

## 309) Governance Review Cadence

1. wöchentlicher Kurzreview
2. monatlicher Tiefenreview
3. quartalsweiser Strukturreview

---

## 310) Rule Lifecycle Management

Regeln haben Status:
1. Draft
2. Active
3. Deprecated
4. Replaced

Pflicht:
- Statusänderungen mit Begründung dokumentieren

---

## 311) Rule Deprecation Protocol

Wenn Regel ersetzt wird:
1. Nachfolger benennen
2. Übergangszeit definieren
3. Migration beschreiben

---

## 312) Rule Conflict Register

Template:
```markdown
## Conflict ID
Rules Involved:
Conflict Description:
Resolution:
Owner:
Date:
```

---

## 313) Evidence Retention Policy

Regel:
1. Evidenzartefakte nachvollziehbar speichern
2. Retention-Zeiten definieren
3. sensible Inhalte schützen

---

## 314) Audit Preparation Checklist

- [ ] Regel-IDs aktuell
- [ ] Nachweise verlinkt
- [ ] Incident-Historie vollständig
- [ ] Doku-Mappings konsistent

---

## 315) Executive Summary Template

```markdown
## Scope
## Progress
## Risks
## Decisions Needed
## Next 7 Days
```

---

## 316) Weekly Governance Report Template

```markdown
## Completed
## In Progress
## Blocked
## Risks
## Metrics
## Plan Next Week
```

---

## 317) Monthly Reliability Report Template

```markdown
## Uptime
## Incidents
## MTTR
## Error Budget
## Corrective Actions
```

---

## 318) Quarterly Architecture Review Template

```markdown
## Current Architecture Snapshot
## Drift Analysis
## Major Risks
## Refactoring Priorities
## Roadmap
```

---

## 319) Global Governance Export Policy

Regel:
1. AGENTS-GLOBAL dient als Master-Blueprint.
2. Export in `~/.config/opencode/AGENTS.md` erfolgt additiv.
3. Projekt-Overrides bleiben in `BIOMETRICS/AGENTS.md`.

---

## 320) Expansion Gate (to 4.5k+)

Status:
1. Annex `261–320` vorhanden
2. Datei > 3.6k Zeilen

Nächster Schritt (`321+`):
1. detaillierte Domänen-Runbooks (Auth, Payments, Integrations)
2. erweiterte NLM-Verbatim-Sets
3. tiefere Security/Compliance Kontrollkataloge

Leitsatz:
**Scale with structure, not with noise.**

---

## 321) Domain Runbook Pack — Overview

Dieses Paket standardisiert Runbooks für kritische Domänen:
1. Auth
2. Payments
3. Integrations
4. Data Protection
5. Compliance Controls

Jedes Runbook enthält:
1. Trigger
2. Diagnostik
3. Sofortmaßnahmen
4. Dauerhafte Korrektur
5. Verifikation

---

## 322) Auth Runbook — Login Outage

Trigger:
1. Login endpoint error spike
2. plötzlich erhöhte 401/5xx Quote

Diagnostik:
1. Auth provider health
2. Session store health
3. Token signing keys
4. Clock skew

Sofortmaßnahmen:
1. degraded mode aktivieren
2. betroffene Regionen isolieren
3. Read-only Fallback prüfen

Verifikation:
1. Login Erfolgsrate normalisiert
2. 401/5xx im Sollbereich

---

## 323) Auth Runbook — Token Validation Failure

Checkliste:
1. Token Issuer/Audience korrekt?
2. Signing key rotation erfolgt?
3. Caching stale?
4. Key discovery endpoint erreichbar?

Fix:
1. Key cache invalidieren
2. Metadata refresh
3. Token TTL/clock tolerance kalibrieren

---

## 324) Auth Runbook — Privilege Escalation Suspicion

Sofort:
1. suspect sessions invalidieren
2. Audit Logs sichern
3. RBAC rules locken

Analyse:
1. fehlerhafte Policy identifizieren
2. betroffene Ressourcen eingrenzen
3. Blast Radius berechnen

Nachgang:
1. postmortem
2. policy regression tests ergänzen

---

## 325) Auth Control Checklist (Continuous)

- [ ] least privilege Rollenmodell
- [ ] deny-by-default Policies
- [ ] MFA/2FA für Admin-Pfade
- [ ] Session Expiry geprüft
- [ ] Refresh/Revocation getestet

---

## 326) Payments Runbook — Provider Timeout

Trigger:
1. Payment API timeout spikes
2. erhöht abgebrochene Checkouts

Sofortmaßnahmen:
1. Retry mit idempotency key
2. queue-basierte deferred capture
3. Nutzerkommunikation (degraded payments)

Diagnostik:
1. provider status page
2. regional latency
3. webhook delay

---

## 327) Payments Runbook — Duplicate Charge Risk

Pflicht:
1. Idempotency keys überall
2. reconciliation job aktiv
3. duplicate detection query

Sofort:
1. suspect transactions markieren
2. capture/payout freeze für betroffene IDs
3. Support- und Finance-Alert auslösen

---

## 328) Payments Runbook — Webhook Desync

Symptom:
1. Payment erfolgreich, Systemstatus offen
2. Statusinkonsistenzen zwischen provider und DB

Fixpfad:
1. webhook signature check
2. replay missing events
3. reconcile state machine
4. final consistency report

---

## 329) Payments Data Contract

Pflichtfelder je Transaktion:
1. transaction_id
2. idempotency_key
3. provider_reference
4. amount/currency
5. status_state
6. created_at/updated_at

---

## 330) Payments Security Controls

1. keine PAN/CVV Speicherung
2. Secret scoping strikt
3. webhook signature validation
4. fraud flags in audit trail
5. refund approval workflow

---

## 331) Integrations Runbook — External API Degradation

Trigger:
1. erhöhte Fehlerquote je Integration
2. erhöhte Latenz

Sofort:
1. circuit breaker aktivieren
2. retry policy anpassen
3. fallback data source prüfen

---

## 332) Integrations Runbook — Schema Mismatch

Checkliste:
1. payload diff
2. version header
3. contract test status
4. backward compatibility

Fix:
1. adapter patch
2. feature flag rollout
3. contract tests erweitern

---

## 333) Integrations Runbook — Rate Limit Exhaustion

Sofort:
1. adaptive throttling
2. request batching
3. low-priority traffic drosseln

Nachgang:
1. quota strategy überarbeiten
2. cache hit-rate erhöhen

---

## 334) Integrations Ownership Table

| Integration | Owner | SLA | Fallback | Last Review |
|---|---|---|---|---|
| Payments Provider | | | | |
| Email Provider | | | | |
| Analytics API | | | | |
| CRM API | | | | |

---

## 335) Circuit Breaker Policy

Regeln:
1. Trip-Threshold definieren
2. Cooldown definieren
3. Half-open Probe implementieren
4. Metriken pro Integrationspfad erfassen

---

## 336) Retry Policy Matrix

| Operation Type | Retry? | Backoff | Max Attempts |
|---|---|---|---|
| Read Idempotent | yes | exponential+jitter | 3-5 |
| Write Idempotent | yes | controlled | 2-3 |
| Non-idempotent Write | no | - | 1 |

---

## 337) Data Protection Control Pack

Pflichtkontrollen:
1. Data minimization
2. Purpose limitation
3. Retention limits
4. Access logging
5. deletion requests workflow

---

## 338) Compliance Control Catalog — Identity

| Control ID | Control | Frequency | Evidence |
|---|---|---|---|
| ID-01 | RBAC review | monthly | review log |
| ID-02 | Admin MFA check | weekly | auth report |
| ID-03 | Service account scope audit | monthly | scope diff |

---

## 339) Compliance Control Catalog — Data

| Control ID | Control | Frequency | Evidence |
|---|---|---|---|
| DATA-01 | Encryption-at-rest validation | quarterly | infra report |
| DATA-02 | Backup restore test | monthly | restore log |
| DATA-03 | Retention policy audit | monthly | deletion report |

---

## 340) Compliance Control Catalog — Delivery

| Control ID | Control | Frequency | Evidence |
|---|---|---|---|
| DEV-01 | CI gate enforcement | every PR | pipeline logs |
| DEV-02 | Branch protection check | weekly | repo audit |
| DEV-03 | Dependency vulnerability scan | daily/weekly | scan output |

---

## 341) Access Review Workflow

1. export principals
2. compare with role matrix
3. detect over-privileged accounts
4. remediate + verify
5. archive evidence

---

## 342) Secret Exposure Response

Bei Secret Leak Verdacht:
1. Secret sofort rotieren
2. betroffene tokens invalidieren
3. Zugriffshistorie prüfen
4. blast radius dokumentieren
5. prevention action item erstellen

---

## 343) Logging Control Pack

Regeln:
1. structured logs verpflichtend
2. sensitive fields redacted
3. correlation id überall
4. retention policy dokumentiert

---

## 344) Audit Trail Consistency Check

Prüfen:
1. Actor vorhanden
2. Action vorhanden
3. Target vorhanden
4. Result vorhanden
5. Timestamp vorhanden

---

## 345) Monitoring Quality Gate

- [ ] Alle kritischen Services haben Healthchecks
- [ ] Alerts mit Ownern verknüpft
- [ ] False-positive Rate akzeptabel
- [ ] Dashboards aktuell

---

## 346) Alert Fatigue Mitigation

1. Alert-Rules nach Schwere clustern
2. redundant noisy Alerts entfernen
3. actionable Alerts priorisieren
4. SLO-basierte Schwellen verwenden

---

## 347) Incident Timeline Template

```markdown
## T0 Detection
## T+5 Mitigation
## T+15 Scope Update
## T+30 Root-Cause Hypothesis
## T+60 Action Plan
```

---

## 348) Root Cause Classification

Kategorien:
1. Code defect
2. Config defect
3. Dependency failure
4. Infra failure
5. Human process gap

---

## 349) Corrective vs Preventive Actions

Pflicht:
1. corrective action (sofort)
2. preventive action (dauerhaft)
3. owner + due date

---

## 350) Change Risk Scoring Model

Score = Impact × Complexity × Exposure

Risikoklassen:
1. Low
2. Medium
3. High
4. Critical

Pflicht:
- High/Critical nur mit erweitertem Review und Rollback-Plan.

---

## 351) Deployment Approval Matrix

| Risk Class | Required Approvals |
|---|---|
| Low | 1 reviewer |
| Medium | 2 reviewers |
| High | tech lead + reviewer |
| Critical | exec gate + full rollback readiness |

---

## 352) Change Freeze Exceptions

Erlaubt nur bei:
1. Security emergency
2. Legal compliance emergency
3. Critical production outage

Pflicht:
- Exception mit Incident-ID dokumentieren.

---

## 353) Hotfix Exit Criteria

1. ursprünglicher Fehler behoben
2. Regression ausgeschlossen
3. post-fix monitoring stabil
4. follow-up hardening ticket angelegt

---

## 354) Architecture Governance Board Template

```markdown
## Decision Context
## Options
## Chosen Path
## Tradeoffs
## Risks
## Follow-up Tasks
```

---

## 355) API Lifecycle Registry

Jeder Endpoint hat Status:
1. experimental
2. active
3. deprecated
4. sunset

---

## 356) Contract Drift Detection

Regel:
1. schema snapshots vergleichen
2. breaking deltas markieren
3. migration hints erzeugen

---

## 357) Data Quality Gate

- [ ] Null-/Default-Regeln geprüft
- [ ] Referentielle Integrität geprüft
- [ ] Deduplikation geprüft
- [ ] historische Konsistenz geprüft

---

## 358) Reconciliation Job Template

```text
Job Name:
Source A:
Source B:
Match Keys:
Mismatch Action:
Reporting:
```

---

## 359) Business Continuity Note

Pflicht:
1. kritische Geschäftsprozesse identifizieren
2. Ausfallmodus definieren
3. Recovery-Reihenfolge dokumentieren

---

## 360) Disaster Recovery Drill Cadence

1. vierteljährlicher DR-Test
2. Lessons Learned dokumentieren
3. Maßnahmen in Roadmap übernehmen

---

## 361) Vendor Risk Assessment Template

```markdown
## Vendor
## Criticality
## Data Exposure
## SLA Risk
## Exit Strategy
```

---

## 362) Third-Party Dependency Exit Plan

Für kritische Third-Party Services:
1. Ersatzoption dokumentieren
2. Datenexportpfad sichern
3. Migrationsaufwand schätzen

---

## 363) Documentation Debt Triage

Priorisierung:
1. sicherheitsrelevante Doku-Lücken
2. betriebsrelevante Doku-Lücken
3. onboarding-relevante Doku-Lücken

---

## 364) Knowledge Base Freshness Checks

Regel:
1. veraltete Kapitel markieren
2. Owner zuweisen
3. Aktualisierungsdatum pflegen

---

## 365) Training Runbook for New Agents

Module:
1. Governance Basics
2. NLM Workflow
3. Incident Protocol
4. Doc/Mappings Discipline

---

## 366) Agent Performance Review Template

```markdown
## Delivery Quality
## Verification Quality
## Documentation Quality
## Risk Handling
## Improvement Plan
```

---

## 367) Governance Debt Register

Template:
```markdown
## Debt ID
Rule Area:
Impact:
Remediation Plan:
Owner:
Target Date:
```

---

## 368) Executive KPI Dashboard Skeleton

| KPI | Target | Current | Trend |
|---|---|---|---|
| MTTR | | | |
| CFR | | | |
| Deploy Frequency | | | |
| Documentation Sync Rate | | | |

---

## 369) Weekly Governance Ritual

1. offene Risiken reviewen
2. Incident Learnings aufnehmen
3. Regelkonflikte auflösen
4. nächste Prioritäten festlegen

---

## 370) Monthly Compliance Ritual

1. Control-Katalog prüfen
2. Evidence-Lücken schließen
3. Audit-Readiness updaten

---

## 371) Quarterly Reliability Ritual

1. SLO/SLI Review
2. Incident Trend Analyse
3. Architektur-Hardening Prioritäten

---

## 372) Annual Governance Reset

1. Regelbestand inventarisieren
2. Deprecated markieren
3. Nachfolge-Regeln aktivieren
4. Trainingsmaterial aktualisieren

---

## 373) Domain Runbook Index Template

```markdown
## Auth
## Payments
## Integrations
## Data
## Security
## Infra
```

---

## 374) Domain Health Score Model

Score-Dimensionen:
1. Availability
2. Latency
3. Error Rate
4. Security Findings
5. Documentation Freshness

---

## 375) Cross-Domain Dependency Map

Pflicht:
1. upstream/downstream Abhängigkeiten dokumentieren
2. Failure propagation Risiken markieren

---

## 376) Failure Propagation Guardrails

1. timeouts
2. circuit breakers
3. queue buffering
4. fallback responses

---

## 377) Executive Decision Log Template

```markdown
## Decision ID
## Context
## Decision
## Impact
## Revisit Date
```

---

## 378) Governance Scorecard

| Area | Score (0-5) | Notes |
|---|---|---|
| Security | | |
| Reliability | | |
| Delivery | | |
| Documentation | | |
| NLM Operations | | |

---

## 379) Continuous Improvement Loop

1. messen
2. analysieren
3. verbessern
4. verifizieren
5. dokumentieren

---

## 380) Expansion Gate (to 5k+)

Status:
1. Annex `321–380` vorhanden
2. Datei > 4.2k Zeilen

Nächster Schritt (`381+`):
1. full domain playbooks (Auth/Payments/Integrations) deep variant
2. compliance controls by framework (SOC2/GDPR style)
3. extended templates catalog

Leitsatz:
**Operational excellence is documented excellence.**

---

## 381) SOC2-Oriented Control Pack — Access Control

Kontrollziele:
1. Zugriff nur nach Least Privilege.
2. Rollenänderungen nachvollziehbar.
3. Regelmäßige Access Reviews.

Nachweise:
1. Rollenmatrix
2. Access Review Protokoll
3. Ticket-Historie zu Änderungen

---

## 382) SOC2-Oriented Control Pack — Change Management

Kontrollziele:
1. Jede Änderung ist nachvollziehbar.
2. Risk-Review vor kritischen Changes.
3. Freigaben entsprechend Risikoklasse.

Nachweise:
1. PR-Logs
2. CI-Ergebnisse
3. Approval Matrix

---

## 383) SOC2-Oriented Control Pack — Availability

Kontrollziele:
1. Betriebsbereitschaft messbar halten.
2. Incident Response standardisieren.
3. Wiederherstellung regelmäßig testen.

Nachweise:
1. Uptime Reports
2. Incident-Timeline
3. DR Drill Logs

---

## 384) SOC2-Oriented Control Pack — Monitoring

Kontrollziele:
1. kritische Events erfassen.
2. Alert Routing definiert.
3. Alert Qualität regelmäßig verbessern.

Nachweise:
1. Dashboard Snapshots
2. Alert Regeln
3. Tuning-Protokolle

---

## 385) GDPR-Oriented Control Pack — Data Minimization

Grundsatz:
1. Nur notwendige personenbezogene Daten erheben.
2. Datenzweck dokumentieren.
3. Datennutzung begrenzen.

Nachweise:
1. Dateninventar
2. Zweckbindungen
3. Feldbegründungen

---

## 386) GDPR-Oriented Control Pack — Retention & Deletion

Pflicht:
1. Aufbewahrungsfristen je Datentyp.
2. Lösch-Workflows definiert.
3. Löschanfragen auditierbar.

Nachweise:
1. Retention Policy
2. Deletion Logs
3. Auskunft/Lösch-Reports

---

## 387) GDPR-Oriented Control Pack — Subject Rights

Rechteprozesse:
1. Auskunft
2. Berichtigung
3. Löschung
4. Portabilität
5. Einschränkung

Pflicht:
1. SLA je Recht definieren.
2. Prozessowner benennen.

---

## 388) Compliance Evidence Mapping Table

| Requirement | Control | Evidence Artifact | Owner |
|---|---|---|---|
| Access Control | RBAC Review | access-review.md | |
| Change Mgmt | PR + CI Gate | pipeline-log | |
| Data Retention | Deletion Workflow | deletion-report | |
| Incident Mgmt | Runbook + Timeline | incident-ticket | |

---

## 389) Control Test Cadence

1. Weekly: operative Kontrollen
2. Monthly: zentrale Risiko-Kontrollen
3. Quarterly: tiefgreifende Compliance-Kontrollen
4. Yearly: Governance Reset + Vollaudit

---

## 390) Control Failure Handling

Bei Kontrollversagen:
1. sofortige Einstufung (Severity)
2. kompensierende Kontrolle aktivieren
3. Remediation-Task mit Due Date
4. Re-Test + Abschlussnachweis

---

## 391) Auth Playbook v2 — Account Lockout Storm

Symptom:
1. starke Lockout-Häufung
2. Login Conversion bricht ein

Maßnahmen:
1. Brute-force Indicators prüfen
2. adaptive lock policy aktivieren
3. trusted recovery path bereitstellen
4. Fraud/Abuse Team informieren

---

## 392) Auth Playbook v2 — Session Hijack Suspicion

Sofort:
1. Session revocation triggern
2. Geräte-/IP-Muster prüfen
3. Risiko-Accounts zur erneuten Auth zwingen

Nachgang:
1. Session Binding Policies schärfen
2. Monitoring-Regeln ergänzen

---

## 393) Payments Playbook v2 — Settlement Drift

Symptom:
1. interne Summen ≠ Provider Settlement

Vorgehen:
1. Reconciliation Job ausführen
2. Driftklassifikation je Ursache
3. betroffene Buchungen markieren
4. Korrektur- und Kommunikationsplan

---

## 394) Payments Playbook v2 — Refund Anomaly

Kontrollen:
1. Refund Limits je Rolle
2. Doppelbestätigung bei High-Risk Fällen
3. Audit Trail inkl. Reason Code

Incident:
1. auffällige Refunds isolieren
2. Fraud-Hinweise prüfen

---

## 395) Integrations Playbook v2 — OAuth Token Expiry Cascade

Symptom:
1. zeitgleich viele 401 von Dritt-APIs

Fix:
1. Token refresh pipeline prüfen
2. secret rotation drift prüfen
3. Backoff und re-auth sequenziell fahren

---

## 396) Integrations Playbook v2 — Webhook Integrity Failure

Kontrollen:
1. Signature-Validation
2. Replay-Protection
3. Timestamp Tolerance

Failure:
1. Unsignierte Events verwerfen
2. Event source isolieren
3. Integritätsbericht erstellen

---

## 397) Domain Playbook Template v2

```markdown
## Trigger
## Detection Signals
## Immediate Actions
## Short-Term Mitigation
## Long-Term Fix
## Verification
## Evidence Links
## Owner
```

---

## 398) Compliance Review Board Template

```markdown
## Agenda
## Open Findings
## Control Failures
## Remediation Status
## Decisions
## Action Items
```

---

## 399) Risk Acceptance Template

```markdown
## Risk ID
## Business Justification
## Mitigations in Place
## Expiry Date
## Approver
```

---

## 400) Exception Handling Policy

Ausnahmen von Regeln sind nur gültig mit:
1. eindeutiger Begründung
2. zeitlicher Begrenzung
3. dokumentierter Kompensation
4. benanntem Owner

---

## 401) Governance SLA Table

| Process | SLA | Escalation |
|---|---|---|
| Incident Initial Response | | |
| Risk Review | | |
| Docs Sync | | |
| Control Re-test | | |

---

## 402) Documentation Control Checklist (Deep)

- [ ] Architektur aktualisiert
- [ ] API/Command Mapping synchron
- [ ] DB Mapping synchron
- [ ] NLM Mapping synchron
- [ ] Meeting/Changelog aktualisiert

---

## 403) Advanced NLM Prompt Contract

Pflichtfelder:
1. Role
2. Objective
3. Scope Boundaries
4. Required Output Format
5. Quality Criteria
6. Validation Step

---

## 404) Advanced NLM Anti-Hallucination Controls

1. Quellenpflicht für kritische Aussagen
2. Unsicherheiten kennzeichnen
3. Widersprüche explizit markieren
4. Output gegen interne Doku gegenprüfen

---

## 405) NLM Evidence Linking Rule

Jeder wichtige NLM-Output erhält:
1. Notebook-ID
2. Query-Kurzfassung
3. Datum
4. Referenz in passender BIOMETRICS-Datei

---

## 406) Security Controls by Layer

| Layer | Control Focus |
|---|---|
| Frontend | Input validation, CSP, auth UX |
| API | authz, rate limiting, schema validation |
| DB | RLS, least privilege, migration safety |
| Infra | secrets, network boundaries, observability |

---

## 407) Privacy Impact Assessment Skeleton

```markdown
## Processing Activity
## Data Categories
## Legal Basis
## Risks
## Mitigations
## Residual Risk
```

---

## 408) Data Processing Inventory Template

```markdown
## System
## Data Types
## Purpose
## Retention
## Access Roles
## Third Parties
```

---

## 409) Third-Party Processor Control

Pflicht:
1. Vertrags- und Sicherheitsprüfung
2. Datenfluss dokumentieren
3. Exit-Strategie definieren

---

## 410) Audit Sampling Strategy

1. risikobasierte Stichproben
2. kritische Pfade übergewichten
3. Findings klassifizieren

---

## 411) Corrective Action Tracker Template

```markdown
## Action ID
## Finding Reference
## Owner
## Due Date
## Status
## Verification Evidence
```

---

## 412) Preventive Action Tracker Template

```markdown
## Preventive ID
## Trigger Pattern
## Prevention Mechanism
## Owner
## Review Date
```

---

## 413) Executive Risk Heatmap Guidance

Dimensionen:
1. likelihood
2. impact
3. detectability

Nutzung:
1. Priorisierung von Mitigations
2. Release-Freigaben unterstützen

---

## 414) Security Training Content Map

Module:
1. Secure Coding
2. Auth/Session Security
3. Secrets Hygiene
4. Incident Reporting
5. Data Protection Basics

---

## 415) Compliance Training Content Map

Module:
1. Change Documentation
2. Evidence Retention
3. Control Testing
4. Exception Handling

---

## 416) Annual Control Testing Plan Template

```markdown
## Control Universe
## Test Methods
## Sampling Plan
## Timeline
## Owners
## Reporting Path
```

---

## 417) Governance Maturity Model

Stufen:
1. Initial
2. Repeatable
3. Defined
4. Measured
5. Optimized

Regel:
- Verbesserungsmaßnahmen immer auf nächste Stufe ausrichten.

---

## 418) Long-Term Governance Roadmap Template

```markdown
## Quarter
## Priority Controls
## Planned Improvements
## Dependencies
## Success Metrics
```

---

## 419) Master Blueprint Sync Rule

Pflicht:
1. AGENTS-GLOBAL als Master pflegen
2. Exporte nach global/lokal additiv synchronisieren
3. Drift zwischen Master und Export dokumentieren

---

## 420) Expansion Gate (to 5.5k+)

Status:
1. Annex `381–420` vorhanden
2. Datei kurz vor/über 5k

Nächster Schritt (`421+`):
1. SOC2/GDPR Controls Full Matrix
2. Domain Playbooks v3 (deep incident trees)
3. Governance Automation Checklist

Leitsatz:
**Governance that scales is governance that survives.**

---

## 421) Governance Automation Checklist

- [ ] Rule-Drift-Checks automatisiert
- [ ] Doku-Sync-Checks automatisiert
- [ ] Mapping-Konsistenz-Checks automatisiert
- [ ] Compliance-Evidence-Exports automatisiert

---

## 422) Automation Priority Matrix

| Automationsziel | Nutzen | Priorität |
|---|---|---|
| CI Gate Hardening | hoch | P0 |
| Doku Drift Detection | hoch | P1 |
| Risk Register Sync | mittel | P1 |
| KPI Aggregation | mittel | P2 |

---

## 423) SOC2 Full Control Matrix — Identity Domain

| Control | Design | Operation | Evidence |
|---|---|---|---|
| RBAC | least privilege | monthly review | role-audit |
| MFA | admin mandatory | login enforcement | auth logs |
| Account Lifecycle | joiner/mover/leaver | periodic checks | HR/security sync |

---

## 424) SOC2 Full Control Matrix — Change Domain

| Control | Design | Operation | Evidence |
|---|---|---|---|
| PR Review | mandatory review | enforced on protected branches | PR history |
| CI Gates | lint/typecheck/test/build | block merge on fail | pipeline logs |
| Rollback | defined per critical change | rehearsal cadence | rollback drills |

---

## 425) SOC2 Full Control Matrix — Reliability Domain

| Control | Design | Operation | Evidence |
|---|---|---|---|
| Monitoring | service-level alerts | threshold tuning | dashboard snapshots |
| Incident Mgmt | runbooks + severity | timeline discipline | incident tickets |
| Recovery | DR drills | verified restores | restore reports |

---

## 426) GDPR Full Control Matrix — Lawful Basis

| Control | Requirement | Evidence |
|---|---|---|
| Purpose Binding | documented purpose | data inventory |
| Legal Basis | per processing activity | legal matrix |
| Minimization | only required fields | schema review |

---

## 427) GDPR Full Control Matrix — Subject Rights

| Right | SLA | Process Owner | Evidence |
|---|---|---|---|
| Access | | | request log |
| Rectification | | | change log |
| Erasure | | | deletion log |
| Portability | | | export log |

---

## 428) GDPR Full Control Matrix — Retention

| Data Class | Retention | Deletion Method | Verification |
|---|---|---|---|
| Public | | | |
| Internal | | | |
| Sensitive | | | |
| Restricted | | | |

---

## 429) Incident Tree v3 — Global Entry

```text
Alert/Event
  ├─ Security? → Security branch
  ├─ Reliability? → Reliability branch
  ├─ Data Integrity? → Data branch
  └─ Unknown? → triage + classify <15 min
```

---

## 430) Incident Tree v3 — Security Branch

```text
Security branch
  ├─ credential leak? → rotate + revoke
  ├─ unauthorized access? → isolate + scope
  ├─ abuse pattern? → WAF/rate controls
  └─ report + postmortem + preventive controls
```

---

## 431) Incident Tree v3 — Reliability Branch

```text
Reliability branch
  ├─ outage? → failover/rollback
  ├─ latency spike? → perf triage
  ├─ queue backlog? → consumer scale + DLQ
  └─ recover + verify + monitor
```

---

## 432) Incident Tree v3 — Data Branch

```text
Data branch
  ├─ schema drift? → migration rollback/patch
  ├─ consistency error? → reconcile
  ├─ deletion anomaly? → restore workflow
  └─ full evidence capture
```

---

## 433) Automation Guardrails

1. Automationen sind nachvollziehbar und auditierbar.
2. Jede Automation hat Owner und Failure-Path.
3. Automationsfehler erzeugen standardisierte Tickets.

---

## 434) Policy-as-Code Hook

Regel:
1. wiederkehrende Governance-Prüfungen als Code definieren.
2. CI prüft Compliance-regeln maschinell.

---

## 435) Governance CI Job Catalog

Pflichtjobs:
1. doc-sync-check
2. mapping-consistency-check
3. risk-register-check
4. control-evidence-check

---

## 436) Evidence Packaging Standard

Bundle enthält:
1. Task-ID
2. Rule-ID
3. Nachweisdatei
4. Datum
5. Owner

---

## 437) Rule Coverage Report Template

```markdown
## Rule Coverage Report
## Covered Rules
## Missing Rules
## Conflicts
## Action Plan
```

---

## 438) Drift Detection Report Template

```markdown
## Drift Scope
## Expected State
## Actual State
## Severity
## Remediation
```

---

## 439) Compliance Gap Report Template

```markdown
## Gap ID
## Control Area
## Gap Description
## Risk
## Due Date
## Owner
```

---

## 440) Governance Automation Decision Tree

```text
Check manuell aufwändig?
  ├─ ja → automation candidate
  ├─ compliance-kritisch? → P0 automation
  ├─ häufigkeit hoch? → priorisieren
  └─ owner + rollout plan
```

---

## 441) Control Ownership Matrix

| Control Area | Primary Owner | Secondary Owner |
|---|---|---|
| Security | | |
| Reliability | | |
| Compliance | | |
| Documentation | | |

---

## 442) Audit Trail Normalization Rules

1. Event-Namen standardisieren.
2. Zeitformat standardisieren.
3. Correlation IDs erzwingen.
4. Ergebnisstatus (`success|failure`) verpflichtend.

---

## 443) Reporting Calendar Template

```markdown
## Weekly Reports
## Monthly Reports
## Quarterly Reviews
## Annual Reset
```

---

## 444) Review Meeting Agenda — Reliability

1. Incident Trends
2. Error Budget Status
3. Top 5 Risks
4. Mitigation Progress
5. Next Priorities

---

## 445) Review Meeting Agenda — Security

1. offene kritische Findings
2. Secret Hygiene Status
3. Auth/Access Review
4. Security Incident Learnings
5. Hardening Roadmap

---

## 446) Review Meeting Agenda — Compliance

1. Control Coverage
2. Evidence Completeness
3. Gap Closure Status
4. Exceptions & Expiry
5. Next Audit Readiness

---

## 447) Exception Register Template

```markdown
## Exception ID
## Rule/Control
## Reason
## Compensating Controls
## Expiry
## Owner
```

---

## 448) Compensating Control Template

```markdown
## Risk
## Primary Control Missing
## Compensating Control
## Effectiveness Validation
```

---

## 449) Governance Risk Heatmap Table

| Risk | Likelihood | Impact | Priority | Owner |
|---|---|---|---|---|
| | | | | |

---

## 450) Documentation Automation Hooks

1. Auto-check für fehlende Pflichtsektionen
2. Link-Check für kritische Referenzen
3. Änderungsdatum-Check

---

## 451) Template Registry

Pflicht:
1. zentrale Liste aller aktiven Templates
2. Versionsstatus je Template
3. Deprecated Markierungen

---

## 452) Template Versioning Policy

Regeln:
1. breaking template changes versionieren
2. Migrationshinweise bereitstellen
3. alte Versionen referenzierbar halten

---

## 453) Playbook Versioning Policy

1. v1 baseline
2. v2 erweitert
3. v3 domain-deep

Pflicht:
- jede Version mit Datum + Scope dokumentieren

---

## 454) Domain Playbook Index v3

| Domain | Version | Last Updated | Owner |
|---|---|---|---|
| Auth | | | |
| Payments | | | |
| Integrations | | | |
| Data | | | |

---

## 455) Governance Automation KPI Table

| KPI | Definition | Target |
|---|---|---|
| Auto-checked PR ratio | PRs with all governance checks | |
| Drift detection lead time | detection speed | |
| Evidence completeness | complete evidence bundles | |

---

## 456) Incident Automation Hooks

1. auto-create incident ticket on critical alerts
2. auto-attach runbook links
3. auto-open postmortem task

---

## 457) Security Automation Hooks

1. secret scan in CI
2. dependency vulnerability scan
3. policy checks for critical configs

---

## 458) Compliance Automation Hooks

1. evidence freshness checks
2. control test reminder generation
3. exception expiry alerts

---

## 459) Governance Backlog Prioritization

Prioritätslogik:
1. compliance risk
2. security risk
3. reliability impact
4. automation leverage

---

## 460) Governance Debt Burn-Down

1. debt items erfassen
2. owner zuweisen
3. sprintweise abbauen
4. Fortschritt berichten

---

## 461) Evidence Review Checklist

- [ ] Rule-ID vorhanden
- [ ] Nachweis nachvollziehbar
- [ ] Zeitpunkt dokumentiert
- [ ] Owner dokumentiert
- [ ] Referenzlink gültig

---

## 462) Deep Incident Tree — Auth v3

```text
Auth Incident
  ├─ Credential issue? → rotate + revoke + notify
  ├─ Session issue? → invalidate + force re-auth
  ├─ RBAC issue? → lock sensitive paths + policy patch
  └─ verify + monitor + postmortem
```

---

## 463) Deep Incident Tree — Payments v3

```text
Payment Incident
  ├─ duplicate risk? → idempotency audit + freeze
  ├─ settlement drift? → reconcile + finance sync
  ├─ provider outage? → fallback queue mode
  └─ verify balances + customer communication
```

---

## 464) Deep Incident Tree — Integrations v3

```text
Integration Incident
  ├─ auth broken? → token refresh/rotate
  ├─ schema drift? → adapter patch + contract test
  ├─ rate limit? → throttle + cache + backlog
  └─ stability verification + trend watch
```

---

## 465) Deep Incident Tree — Data v3

```text
Data Incident
  ├─ integrity breach? → write freeze + reconcile
  ├─ migration issue? → rollback or hot patch
  ├─ retention breach? → legal/compliance escalation
  └─ restore verification + evidence bundle
```

---

## 466) Executive Escalation Matrix v2

| Trigger | Escalation Path |
|---|---|
| SEV-1 active | immediate exec + incident lead |
| legal/privacy risk | compliance lead + exec |
| recurring high incidents | architecture board |

---

## 467) Governance Communication Cadence v2

1. daily incident sync (if active incidents)
2. weekly governance update
3. monthly compliance checkpoint
4. quarterly strategic review

---

## 468) Stakeholder Update Template

```markdown
## Situation
## Impact
## Actions Taken
## Current Status
## Next Update ETA
```

---

## 469) Governance Decision SLA

| Decision Type | SLA |
|---|---|
| Critical Incident Decision | immediate |
| High Risk Change Approval | same day |
| Standard Governance Update | weekly cycle |

---

## 470) Audit Readiness Score Model

Score Dimensionen:
1. Control Coverage
2. Evidence Completeness
3. Incident Process Quality
4. Documentation Freshness

---

## 471) Audit Readiness Checklist v2

- [ ] Control Matrix aktuell
- [ ] Evidence Register vollständig
- [ ] Exception Register aktuell
- [ ] Incident Postmortems vollständig

---

## 472) Governance Release Note Template

```markdown
## New Rules
## Updated Rules
## Deprecated Rules
## Migration Notes
## Effective Date
```

---

## 473) Rule Adoption Tracking

Pflicht:
1. neue Regeln mit Adoption-Status tracken
2. Blocker für Adoption dokumentieren
3. Remediation-Pläne setzen

---

## 474) Adoption KPI Table

| KPI | Definition | Target |
|---|---|---|
| Rule Adoption Rate | active use of new rules | |
| Exception Rate | exceptions per period | |
| Remediation Lead Time | avg time to close gaps | |

---

## 475) Governance Training Completion Tracking

| Module | Required Role | Completion % |
|---|---|---|
| Security Basics | all engineers | |
| Incident Response | ops + leads | |
| Compliance Controls | leads + governance | |

---

## 476) Continuous Audit Loop

1. kontrollen testen
2. lücken erfassen
3. remediations umsetzen
4. re-test
5. report

---

## 477) Long-Term Governance Archive Policy

Regeln:
1. alte Regelversionen archivieren
2. historisierte Begründungen erhalten
3. keine Löschung ohne explizite Freigabe

---

## 478) Master Blueprint Integrity Statement

Diese Datei bleibt Master-Blueprint.

Pflicht:
1. Additive Erweiterung
2. Verlustfreie Konsolidierung
3. klare Versionsführung

---

## 479) Global Export Verification Checklist

- [ ] Export in globale AGENTS-Datei durchgeführt
- [ ] lokale Overrides validiert
- [ ] Konflikte dokumentiert
- [ ] Referenzen aktualisiert

---

## 480) Expansion Gate (to 6k+)

Status:
1. Annex `421–480` vorhanden
2. Governance Automation Deep Pack vorhanden

Nächster Schritt (`481+`):
1. policy-as-code examples
2. full audit artifact pack
3. domain-specific playbook v4

Leitsatz:
**Automate what you enforce, enforce what you measure.**

---

## 481) Policy-as-Code Baseline

Regeln als Code definieren für:
1. PR-Gates
2. Secret Scans
3. Doku-Sync
4. Mapping-Konsistenz

---

## 482) Policy-as-Code Rule Schema

```yaml
ruleId: GOV-XXX
title: string
severity: low|medium|high|critical
scope: repository|pipeline|runtime
check: description
passCondition: description
failAction: block|warn|ticket
owner: team-or-role
```

---

## 483) Governance Policy Catalog (Starter)

| Rule ID | Bereich | Severity | Fail Action |
|---|---|---|---|
| GOV-001 | CI Gate | high | block |
| GOV-002 | Secret Hygiene | critical | block |
| GOV-003 | Doku Sync | medium | ticket |
| GOV-004 | Mapping Sync | high | block |

---

## 484) Secret Scan Policy

Pflicht:
1. pre-commit oder CI secret scan
2. False-Positive Handling dokumentieren
3. Treffer sofort rotieren und remediieren

---

## 485) Documentation Sync Policy-as-Code

Regel:
1. Wenn Code in kritischen Bereichen geändert wird, müssen zugehörige Doku-Dateien geändert sein.
2. Fehlende Doku-Änderungen blockieren oder ticketn.

---

## 486) Mapping Sync Policy-as-Code

Regel:
1. Backend-Change ohne Mapping-Update ist nicht release-ready.
2. DB/API/Frontend-Mapping bei betroffenen Änderungen erzwingen.

---

## 487) Audit Artifact Pack — Index

Pflichtartefakte:
1. Control Matrix
2. Evidence Register
3. Incident Register
4. Exception Register
5. Risk Register
6. Training Completion Logs

---

## 488) Audit Artifact — Evidence Register v2

```markdown
## Evidence ID
## Rule/Control ID
## Artifact Path
## Owner
## Date
## Verification Status
```

---

## 489) Audit Artifact — Incident Register v2

```markdown
## Incident ID
## Severity
## Domain
## Root Cause Category
## Corrective Action Status
## Preventive Action Status
```

---

## 490) Audit Artifact — Exception Register v2

```markdown
## Exception ID
## Control/Rule
## Validity Window
## Compensating Control
## Approver
## Expiry Alert Status
```

---

## 491) Audit Artifact — Risk Register v2

```markdown
## Risk ID
## Category
## Likelihood
## Impact
## Priority
## Mitigation Progress
```

---

## 492) Audit Artifact — Control Test Log

```markdown
## Control ID
## Test Date
## Tester
## Result
## Findings
## Re-test Date
```

---

## 493) Domain Playbook v4 — Auth Unauthorized Burst

```text
Detect unauthorized burst
  ├─ check IP/device anomaly
  ├─ validate token issuance path
  ├─ throttle suspicious traffic
  ├─ enforce step-up auth
  └─ verify recovery and log evidence
```

---

## 494) Domain Playbook v4 — Auth MFA Outage

1. prüfen: Provider/MFA channel health
2. fallback auth path aktivieren (policy-konform)
3. admin-privileged actions zusätzlich absichern
4. MFA recovery report erstellen

---

## 495) Domain Playbook v4 — Payments Chargeback Spike

1. chargeback cluster analysieren
2. betroffene Segmente isolieren
3. anti-fraud controls schärfen
4. Finance + Support playbook triggern

---

## 496) Domain Playbook v4 — Payments Ledger Mismatch

1. ledger snapshot sichern
2. reconcile pipeline starten
3. mismatch reasons klassifizieren
4. correction batch mit audit trail ausführen

---

## 497) Domain Playbook v4 — Integrations Credential Drift

1. secret versions vergleichen
2. token refresh paths verifizieren
3. rotation history prüfen
4. drift fixen + regression test

---

## 498) Domain Playbook v4 — Integrations Webhook Flood

1. dedupe key enforcement prüfen
2. inbound throttling aktivieren
3. queue buffering + rate shaping
4. integrity verification und catch-up

---

## 499) Domain Playbook v4 — Data Quality Regression

1. Qualitätsmetriken vergleichen (vor/nach)
2. betroffene Pipelines isolieren
3. bad records markieren
4. backfill/recompute Strategie ausführen

---

## 500) Domain Playbook v4 — Data Retention Breach

1. betroffene Datensätze identifizieren
2. Compliance/Legal informieren
3. Lösch- und Korrekturprozess starten
4. Evidenzkette schließen

---

## 501) Security Control Matrix — App Layer

| Control | Check | Frequency |
|---|---|---|
| Input Validation | schema enforcement | continuous |
| AuthZ | route-level policies | per release |
| Error Redaction | no sensitive output | continuous |

---

## 502) Security Control Matrix — API Layer

| Control | Check | Frequency |
|---|---|---|
| Rate Limiting | threshold + burst behavior | continuous |
| Idempotency | write path guarantees | per release |
| Audit Logging | actor/action/target/result | continuous |

---

## 503) Security Control Matrix — Data Layer

| Control | Check | Frequency |
|---|---|---|
| RLS | positive/negative tests | per migration |
| Backup Restore | restore drill | monthly |
| Retention | policy compliance | monthly |

---

## 504) Security Control Matrix — Infrastructure Layer

| Control | Check | Frequency |
|---|---|---|
| Secret Scope | least privilege tokens | monthly |
| Network Boundaries | allowed paths only | quarterly |
| Observability | actionable alerts | continuous |

---

## 505) Compliance Control Matrix — Documentation

| Control | Requirement | Evidence |
|---|---|---|
| Doku Freshness | timely updates | changelog + timestamps |
| Traceability | link code↔docs | mapping files |
| Ownership | named owners | ownership table |

---

## 506) Compliance Control Matrix — Operations

| Control | Requirement | Evidence |
|---|---|---|
| Runbooks | exist + tested | runbook log |
| Incident Process | severity + timeline | incident register |
| Recovery Process | tested | drill reports |

---

## 507) Compliance Control Matrix — Training

| Control | Requirement | Evidence |
|---|---|---|
| Security Training | periodic completion | training log |
| Governance Training | role-based modules | completion records |
| Incident Drills | participation tracked | drill attendance |

---

## 508) Automated Governance Alerts

Alert-Klassen:
1. Missing docs for critical changes
2. Failed control tests
3. Exception expiry approaching
4. Drift detection triggered

---

## 509) Governance Bot Command Set (Concept)

```text
/gov check docs-sync
/gov check mapping
/gov check controls
/gov report weekly
/gov risk list
```

---

## 510) Weekly Automation Report Template

```markdown
## Automated Checks Run
## Failures
## Auto-Created Tickets
## Manual Follow-ups
## Trend
```

---

## 511) Monthly Audit Report Template v2

```markdown
## Control Coverage
## Evidence Completeness
## Open Findings
## Exceptions
## Remediation Velocity
```

---

## 512) Quarterly Governance Health Report

```markdown
## Maturity Score
## Top Risks
## Incident Trends
## Compliance Readiness
## Next Quarter Priorities
```

---

## 513) Governance Debt Prioritization Matrix

| Debt Item | Risk | Effort | Priority |
|---|---|---|---|
| | | | |

---

## 514) Backlog Hygiene for Governance Work

1. klare Rule-ID je Task
2. eindeutiger Owner
3. Due Date
4. Evidence Definition

---

## 515) Task-to-Control Traceability

Jeder Governance Task referenziert:
1. Control ID
2. Risk ID
3. Evidence ID

---

## 516) Exception Expiry Workflow

1. Reminder vor Ablauf
2. Re-evaluation
3. entweder schließen oder erneuern mit neuer Begründung

---

## 517) Control Re-Test Workflow

1. initial fail dokumentieren
2. remediation durchführen
3. re-test mit evidence
4. status aktualisieren

---

## 518) Governance Change Approval Board

Für kritische Governance-Änderungen:
1. Proposal
2. Impact Analysis
3. Decision
4. Migration Plan

---

## 519) Master Template Catalog v2

Kategorien:
1. Incident
2. Risk
3. Compliance
4. Delivery
5. Governance

Pflicht:
- jede Kategorie mit aktiver Version führen.

---

## 520) Quality-of-Governance Score Formula

Beispiel:
`QoG = (ControlCoverage + EvidenceCompleteness + ResponseDiscipline + DocumentationFreshness) / 4`

---

## 521) Incident-to-Learning Pipeline

1. Incident
2. Root Cause
3. Learning Item
4. Policy/Playbook Update
5. Training Update

---

## 522) Learning Register Template

```markdown
## Learning ID
## Incident Reference
## Lesson
## Action
## Owner
## Due Date
```

---

## 523) Governance Knowledge Graph Hook

Verknüpfung:
1. Rule IDs
2. Control IDs
3. Risk IDs
4. Evidence IDs
5. Incident IDs

---

## 524) Policy Conflict Auto-Detection Concept

Idee:
1. diff-basierte Regelanalyse
2. conflict flags bei widersprüchlichen Aussagen
3. automatische Erstellung eines Conflict-Tickets

---

## 525) Auto-Generated Executive Brief Template

```markdown
## Current Governance Posture
## Critical Risks
## Control Failures
## Required Executive Decisions
```

---

## 526) Governance Incident SLA Table v2

| Severity | First Response | Mitigation Start | Executive Update |
|---|---|---|---|
| SEV-1 | | | |
| SEV-2 | | | |
| SEV-3 | | | |

---

## 527) Governance Escalation Contact Model

Pflicht:
1. role-based escalation chains
2. backup contacts
3. out-of-hours paths

---

## 528) Compliance Calendar Automation

1. control test reminders
2. quarterly review reminders
3. annual reset reminders

---

## 529) Governance Metadata Standard

Jeder Rule/Control Eintrag enthält:
1. ID
2. Version
3. Status
4. Owner
5. Last Updated

---

## 530) Policy Version Manifest

```yaml
policySet:
  version: x.y.z
  generatedAt: YYYY-MM-DD
  activeSections:
    - 1-100
    - 101-200
```

---

## 531) Rule Adoption Scorecard v2

| Rule Group | Adoption % | Blockers |
|---|---|---|
| Security | | |
| Reliability | | |
| NLM Ops | | |
| Documentation | | |

---

## 532) Governance Backtesting Concept

1. vergangene Incidents gegen aktuelle Regeln prüfen
2. Regelwirksamkeit bewerten
3. Lücken schließen

---

## 533) Continuous Control Improvement Loop

1. beobachten
2. messen
3. verbessern
4. erneut testen
5. standardisieren

---

## 534) Global-to-Local Override Contract

Regel:
1. lokale Overrides erlauben, wenn globaler Kern unverletzt bleibt.
2. jeder Override muss begründet und referenziert sein.

---

## 535) Override Register Template

```markdown
## Override ID
## Global Rule Reference
## Local Justification
## Scope
## Expiry
```

---

## 536) Rule Sunset Workflow

1. sunset ankündigen
2. migration support bereitstellen
3. adoption überwachen
4. sunset finalisieren

---

## 537) Legacy Governance Migration Pack

1. legacy rules inventarisieren
2. mapping auf neue rule IDs
3. konfliktauflösung
4. evidence migration

---

## 538) Governance Success Criteria v2

Erfolg, wenn:
1. Regeln angewandt werden
2. Nachweise vollständig sind
3. Incidents schneller gelöst werden
4. Drift aktiv reduziert wird

---

## 539) Executive Governance Dashboard Fields

1. Open Critical Risks
2. Control Pass Rate
3. Incident MTTR
4. Rule Adoption Rate
5. Documentation Freshness

---

## 540) Expansion Gate (to 6.5k+)

Status:
1. Annex `481–540` vorhanden
2. Policy-as-Code + Audit Artifact Pack integriert

Nächster Schritt (`541+`):
1. Governance Automation Workflows (pseudo-code + pipeline snippets)
2. Full Domain Playbooks v5
3. Advanced Cross-Repository Governance Federation

Leitsatz:
**Codify governance, verify continuously, improve relentlessly.**

---

## 541) Repository-Wide Propagation Mandate

Pflicht:
1. Jede Markdown-Datei im Repository enthält einen expliziten `Global-Mandate-Alignment` Abschnitt.
2. Der Abschnitt referenziert die Bindung an `AGENTS-GLOBAL.md`.
3. Domain-spezifische Dateien konkretisieren die globalen Kernregeln (Security, Mapping, CI/CD, Incident, Evidence).
4. Bei neuen Dateien wird der Alignment-Abschnitt direkt bei Erstellung ergänzt.

Kontrollpunkt:
- Fehlender Alignment-Abschnitt gilt als Governance-Drift und muss sofort korrigiert werden.
