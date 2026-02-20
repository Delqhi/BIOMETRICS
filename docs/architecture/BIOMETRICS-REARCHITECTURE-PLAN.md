# BIOMETRICS REARCHITECTURE PLAN

**Created:** 2026-02-20  
**Status:** ✅ COMPLETE - READY FOR USER APPROVAL  
**Session:** ses_3836b4ccaffei8C2JF4aTSVtqe  

---

## 🎯 VISION

**BIOMETRICS** ist das zentrale **Rules & Templates Repository** für KI-Agenten-Systeme. Es definiert das "Gesetzbuch" das ALLE Agents befolgen müssen - globale Regeln, Projekt-Templates, Tool-Konfigurationen (OpenCode, OpenClaw, MCP, Models).

Die **CLI mit TUI** (Bubbletea) orchestriert:
- **Onboarding** (API Keys, Accounts, Tools)
- **Project Setup** (neue Projekte nach Templates)
- **Agent Loop** (24/7 autonome Agenten - Zukunft)

**Kernprinzip:** Rules-first → Templates-second → CLI-third → Agent Loop

---

## 🏗️ NEUE STRUKTUR

```
BIOMETRICS/
├── rules/                    # 📜 DAS GESETZBUCH (höchste Priorität)
│   ├── global/              # Globale Regeln für ALLE Agents
│   │   ├── AGENTS.md
│   │   ├── coding-standards.md
│   │   ├── documentation-rules.md
│   │   ├── security-mandates.md
│   │   └── git-workflow.md
│   ├── tools/               # Tool-spezifische Regeln
│   │   ├── opencode-rules.md
│   │   ├── openclaw-rules.md
│   │   ├── mcp-server-rules.md
│   │   └── model-assignment.md
│   └── projects/            # Projekt-spezifische Regeln
│
├── templates/               # 🏗️ PROJEKT-VORLAGEN
│   ├── global/              # Globale Templates
│   ├── opencode/            # OpenCode Projekt-Templates
│   └── openclaw/            # OpenClaw Projekt-Templates
│
├── configs/                 # ⚙️ TOOL-KONFIGURATIONEN
│   ├── opencode/
│   ├── openclaw/
│   └── mcp-servers/
│
├── onboarding/              # 🎓 ONBOARDING PROZESS
│   ├── checklist.md
│   ├── api-keys/
│   ├── accounts/
│   └── tools/
│
├── cli/                     # 💻 CLI TOOL (Bubbletea TUI)
│   ├── cmd/
│   ├── tui/
│   └── internal/
│
├── docs/                    # 📚 DOKUMENTATION (gerettet)
├── archive/                 # 🗄️ ALTE STRUKTUR (Sprint 5)
└── README.md
```

---

## 📋 10-PHASEN PLAN

### Phase 1: AUDIT & ARCHITECTURE
- **Ziel:** Vollständige Analyse + neue Architektur definieren
- **Deliverables:** audit-report.md, ARCHITECTURE.md, migration-plan.md
- **Modell:** MiniMax M2.5 (schnell, MD-Dateien)

### Phase 2: RULES - GLOBAL RULES
- **Ziel:** Das "Gesetzbuch" für alle Agents
- **Deliverables:** 5 Rule-Files (AGENTS.md, coding, docs, security, git)
- **Modell:** MiniMax M2.5

### Phase 3: RULES - TOOL RULES
- **Ziel:** Tool-spezifische Regeln
- **Deliverables:** 4 Tool-Rule-Files (OpenCode, OpenClaw, MCP, Models)
- **Modell:** MiniMax M2.5

### Phase 4: TEMPLATES - GLOBAL
- **Ziel:** Globale Templates für Projekt-Replikation
- **Deliverables:** 6 Templates (AGENTS, BLUEPRINT, README, docker, package, tsconfig)
- **Modell:** MiniMax M2.5

### Phase 5: TEMPLATES - PROJECT
- **Ziel:** OpenCode/OpenClaw Projekt-Templates
- **Deliverables:** 5 Projekt-Templates (standard, minimal, enterprise)
- **Modell:** MiniMax M2.5

### Phase 6: CONFIGS
- **Ziel:** Zentrale Tool-Konfigurationen
- **Deliverables:** opencode.json, openclaw.json, MCP configs, model presets
- **Modell:** Qwen 3.5 (komplexe JSON Configs)

### Phase 7: ONBOARDING
- **Ziel:** Schritt-für-Schritt Onboarding Guides
- **Deliverables:** 100+ Schritte Checkliste, Provider Guides, Tool Setups
- **Modell:** MiniMax M2.5

### Phase 8: CLI TUI
- **Ziel:** CLI Tool mit Bubbletea TUI
- **Deliverables:** main.go, onboarding.go, project.go, TUI components
- **Modell:** Qwen 3.5 (Go Code)

### Phase 9: AGENT LOOP ARCHITECTURE
- **Ziel:** Architektur für 24/7 autonome Agenten
- **Deliverables:** Architektur-Docs, Orchestrator Design, Queue Design
- **Modell:** Qwen 3.5 (komplexes Design)

### Phase 10: MIGRATION & VALIDATION
- **Ziel:** Alte Struktur migrieren + validieren
- **Deliverables:** Migration Scripts, docs/, archive/, README, Tests
- **Modell:** MiniMax M2.5

---

## ⚡ MODELL-ZUWEISUNG

| Phase | Modell | Warum |
|-------|--------|-------|
| 1-5 | MiniMax M2.5 | MD-Dateien, schnell, 10x parallel |
| 6 | Qwen 3.5 | Komplexe JSON Configs |
| 7 | MiniMax M2.5 | Onboarding Guides |
| 8-9 | Qwen 3.5 | **Code + Architektur** |
| 10 | MiniMax M2.5 | Migration + Docs |

**Regel:** Max 1 Qwen 3.5 gleichzeitig (Rate Limit!), bis zu 10 MiniMax parallel

---

## 🚨 KRITISCHE ENTSCHEIDUNGEN

### Sprint 5 Status: **ABBRECHEN**
- ❌ CircuitBreaker, Vault, WebSocket, etc. → **SINNLLOS** (keine Use Cases)
- ✅ Werden nach `archive/` verschoben
- ✅ Nicht erweitern, nicht fertigstellen

### Docs Status: **SICHTEN + MIGRIEREN**
- 📊 200+ alte MD-Dateien existieren
- ✅ Brauchbares wird migriert
- ❌ Unsinn (BLOCKCHAIN.md, WEBSHOP.md) → archive/

### Rules Status: **BASIS AUF AGENTS.md V19.2**
- 📜 ~/.config/opencode/AGENTS.md (286KB, 3100+ Zeilen) ist Source of Truth
- ✅ Daraus werden rules/global/AGENTS.md extrahiert
- ✅ Andere Rules basieren darauf

---

## ✅ SUCCESS CRITERIA

### Phase 1-3 (Rules Foundation)
- [ ] Audit-Report (alle 273 Dateien kategorisiert)
- [ ] 9 Rule-Files (5 global + 4 tools)
- [ ] Jede Datei 500+ Zeilen

### Phase 4-6 (Templates + Configs)
- [ ] 6 globale Templates kopierfertig
- [ ] 5 Projekt-Templates
- [ ] Alle Configs valide JSON

### Phase 7 (Onboarding)
- [ ] Checkliste 100+ Schritte
- [ ] Jeder Provider hat Setup-Guide
- [ ] Alle Guides getestet

### Phase 8 (CLI TUI)
- [ ] CLI baut ohne Fehler
- [ ] TUI zeigt Dashboard
- [ ] Onboarding-Wizard funktioniert
- [ ] Project-Wizard erstellt Projekte

### Phase 9 (Agent Loop)
- [ ] Architektur-Dokument 500+ Zeilen
- [ ] Orchestrator Design implementierbar
- [ ] Queue Design skaliert

### Phase 10 (Migration)
- [ ] 200+ docs gesichtet
- [ ] Brauchbare migriert
- [ ] Sprint 5 archiviert
- [ ] README Document360-konform
- [ ] Validation Scripts grün

---

## 🎯 NÄCHSTE SCHRITTE

**JETZT:** User muss Plan genehmigen

**DANN:** Phase 1 starten mit:
```bash
delegate_task(category="research", run_in_background=true, prompt="Phase 1: AUDIT & ARCHITECTURE - siehe BIOMETRICS-REARCHITECTURE-PLAN.md")
```

**WICHTIG:** 
- Sprint 5 TODOs (CEO-044 bis CEO-048) werden **pausiert**
-KEINE neuen Packages bevor Fundament steht
-IMMER Plan konsultieren bevor gearbeitet wird

---

**PLAN STATUS:** ✅ COMPLETE - WAITING FOR USER APPROVAL
