# BIOMETRICS ARCHITECTURE

**Version:** 2.0 "Complete Rearchitecture"  
**Date:** 2026-02-20  
**Status:** ✅ APPROVED FOR IMPLEMENTATION  
**Based on:** Audit Report (254 files analyzed)  

---

## 🎯 VISION

BIOMETRICS ist das **zentrale Rules & Templates Repository** für KI-Agenten-Systeme. Es definiert das "Gesetzbuch" das ALLE Agents befolgen müssen.

**Drei Hauptsäulen:**
1. **Rules** - Globale Regeln für alle Agents (Source of Truth: ~/.config/opencode/AGENTS.md)
2. **Templates** - Projekt-Vorlagen für schnelle Replikation
3. **CLI** - Bubbletea TUI für Onboarding + Project Setup

---

## 🏗️ NEUE STRUKTUR

### Complete Directory Tree

```
BIOMETRICS/
│
├── 📜 rules/                          # DAS GESETZBUCH (höchste Priorität)
│   ├── global/                        # Globale Regeln für ALLE Agents
│   │   ├── AGENTS.md                  # Hauptregeln (extrahiert aus ~/.config/opencode/AGENTS.md)
│   │   ├── coding-standards.md        # TypeScript, Go, Error Handling
│   │   ├── documentation-rules.md     # 500-line mandate, Trinity docs
│   │   ├── security-mandates.md       # Secrets, Git, Permissions
│   │   └── git-workflow.md            # Commits, Branches, PRs
│   │
│   ├── tools/                         # Tool-spezifische Regeln
│   │   ├── opencode-rules.md          # OpenCode usage, models, providers
│   │   ├── openclaw-rules.md          # OpenClaw config, agents, MCPs
│   │   ├── mcp-server-rules.md        # MCP integration, wrapper pattern
│   │   └── model-assignment.md        # Wann welches Modell? (Qwen/Kimi/MiniMax)
│   │
│   └── projects/                      # Projekt-spezifische Regeln
│       ├── sin-solver-rules.md        # SIN-Solver spezifisch
│       ├── delqhi-rules.md            # Delqhi Platform spezifisch
│       └── [projekt]-rules.md
│
├── 🏗️ templates/                      # PROJEKT-VORLAGEN
│   ├── global/                        # Globale Templates
│   │   ├── AGENTS.md                  # Template für AGENTS.md (500+ lines)
│   │   ├── BLUEPRINT.md               # 22-Säulen Blueprint Template
│   │   ├── README.md                  # Document360 Standard
│   │   ├── docker-compose.yml         # Modular architecture
│   │   ├── package.json               # TypeScript strict mode
│   │   └── tsconfig.json              # Strict configuration
│   │
│   ├── opencode/                      # OpenCode Projekt-Templates
│   │   ├── standard/                  # Standard OpenCode project
│   │   ├── minimal/                   # Minimal setup
│   │   └── enterprise/                # Full enterprise (26 pillars)
│   │
│   └── openclaw/                      # OpenClaw Projekt-Templates
│       ├── standard/                  # Standard OpenClaw project
│       └── enterprise/                # OpenClaw enterprise
│
├── ⚙️ configs/                        # TOOL-KONFIGURATIONEN
│   ├── opencode/
│   │   ├── opencode.json              # Master config mit allen Providern
│   │   ├── provider-configs/          # Google, Streamlake, XiaoMi, ZEN
│   │   └── model-presets/             # Coding, Research, Writing presets
│   │
│   ├── openclaw/
│   │   ├── openclaw.json              # Agent defaults, MCPs
│   │   └── agent-presets/             # Sisyphus, Prometheus, Atlas, etc.
│   │
│   └── mcp-servers/
│       ├── local-mcps.json            # Serena, Tavily, Canva, etc.
│       └── remote-mcps.json           # Docker-based MCPs
│
├── 🎓 onboarding/                     # ONBOARDING PROZESS
│   ├── checklist.md                   # 100+ Schritte Checkliste
│   ├── api-keys/                      # API Key Setup Guides
│   │   ├── nvidia-nim.md
│   │   ├── openrouter.md
│   │   └── [provider].md
│   ├── accounts/                      # Account Setup
│   │   ├── github.md
│   │   ├── docker-hub.md
│   │   └── cloudflare.md
│   └── tools/                         # Tool Installation
│       ├── opencode-setup.md
│       ├── openclaw-setup.md
│       └── docker-setup.md
│
├── 💻 cli/                            # CLI TOOL (Bubbletea TUI)
│   ├── cmd/
│   │   ├── onboarding.go              # Onboarding command
│   │   ├── project.go                 # Project setup command
│   │   └── agent-loop.go              # Agent loop command (Zukunft)
│   ├── tui/                           # Bubbletea UI Components
│   │   ├── dashboard.go               # Haupt-Dashboard
│   │   ├── onboarding-wizard.go       # Onboarding Wizard
│   │   └── project-wizard.go          # Project Wizard
│   ├── internal/
│   │   ├── config/                    # Config loading
│   │   ├── templates/                 # Template rendering
│   │   └── utils/                     # Utilities
│   └── main.go                        # Entry point
│
├── 📚 docs/                           # DOKUMENTATION (restrukturiert)
│   ├── architecture/                  # System-Architektur
│   │   ├── agent-loop.md
│   │   ├── orchestrator-design.md
│   │   └── persistent-queue.md
│   ├── guides/                        # How-To Guides
│   │   ├── setup-guide.md
│   │   ├── agent-delegation.md
│   │   └── troubleshooting.md
│   └── best-practices/                # Best Practices 2026
│       ├── mandates.md
│       ├── workflows.md
│       └── testing.md
│
├── 🔧 scripts/                        # HILFSSKRPTE
│   ├── migrate-old-docs.sh            # Migration script
│   ├── validate-rules.sh              # Rules validation
│   ├── generate-project.sh            # Project generation test
│   └── cleanup.sh                     # Cleanup obsolete files
│
├── 🗄️ archive/                        # ALTE STRUKTUR (Sprint 5 abgebrochen)
│   ├── biometrics-cli/                # Alte CLI (wird ersetzt durch cli/)
│   ├── old-packages/                  # Sprint 5 Packages (sinnlos)
│   │   ├── circuitbreaker/            # → ARCHIVE (kein Use Case)
│   │   ├── vault/                     # → ARCHIVE (kein Use Case)
│   │   ├── websocket/                 # → ARCHIVE (kein Use Case)
│   │   └── ...
│   └── deprecated-docs/               # Alte docs (nicht mehr relevant)
│
├── 📁 assets/                         # MEDIA (keep as-is)
│   ├── 3d/
│   ├── audio/
│   ├── dashboard/
│   ├── diagrams/
│   ├── icons/
│   ├── images/
│   ├── logos/
│   └── videos/
│
├── 📥 inputs/                         # INPUT FILES (keep as-is)
│   └── references/
│
├── 📤 outputs/                        # GENERATED FILES (keep as-is)
│   └── assets/
│
├── 📄 README.md                       # HAUPTDOKUMENTATION (Document360)
├── CHANGELOG.md                       # Änderungen
├── CONTRIBUTING.md                    # Contribution Guidelines
└── .gitignore                         # Git ignore rules
```

---

## 📊 MIGRATION SUMMARY

### Von → Nach

| Alt | Neu | Status |
|-----|-----|--------|
| `docs/` (200+ files chaotisch) | `docs/` (restrukturiert in 3 subdirs) | MIGRATE |
| `biometrics-cli/` | `cli/` (umbenannt, bereinigt) | RENAME |
| `BIOMETRICS/` (nested) | `archive/biometrics-main/` | ARCHIVE |
| `global/` | `rules/global/` | MOVE |
| `local/` | `configs/local/` | MOVE |
| `docs/best-practices/` | `docs/best-practices/` | KEEP (bereits gut) |
| `docs/agents/` | `rules/tools/` | MOVE + MERGE |
| `docs/architecture/` | `docs/architecture/` | KEEP |
| `scripts/` | `scripts/` | KEEP (bereinigen) |
| `assets/` | `assets/` | KEEP |
| `inputs/` | `inputs/` | KEEP |
| `outputs/` | `outputs/` | KEEP |

### Sprint 5 Packages (ARCHIVE)

**14 Packages ohne Use Case:**
- circuitbreaker, completion, encoding, encryption, envconfig
- featureflags, migration, cert, metrics, ratelimit
- plugin, tracing, vault, websocket, errors

**Alle werden nach `archive/old-packages/` verschoben.**

---

## 🔄 MIGRATION PHASEN

### Phase 1: Foundation (DONE)
- ✅ Audit Report erstellt
- ✅ Structure Analysis erstellt
- ✅ Source of Truth Extract erstellt
- ⏳ ARCHITECTURE.md (diese Datei)

### Phase 2: Rules (NEXT)
- [ ] rules/global/AGENTS.md (aus ~/.config/opencode/AGENTS.md extrahieren)
- [ ] rules/global/coding-standards.md
- [ ] rules/global/documentation-rules.md
- [ ] rules/global/security-mandates.md
- [ ] rules/global/git-workflow.md
- [ ] rules/tools/opencode-rules.md
- [ ] rules/tools/openclaw-rules.md
- [ ] rules/tools/mcp-server-rules.md
- [ ] rules/tools/model-assignment.md

### Phase 3: Templates
- [ ] templates/global/AGENTS.md
- [ ] templates/global/BLUEPRINT.md
- [ ] templates/global/README.md
- [ ] templates/global/docker-compose.yml
- [ ] templates/global/package.json
- [ ] templates/global/tsconfig.json
- [ ] templates/opencode/standard/
- [ ] templates/opencode/minimal/
- [ ] templates/opencode/enterprise/
- [ ] templates/openclaw/standard/
- [ ] templates/openclaw/enterprise/

### Phase 4: Configs
- [ ] configs/opencode/opencode.json
- [ ] configs/opencode/provider-configs/
- [ ] configs/opencode/model-presets/
- [ ] configs/openclaw/openclaw.json
- [ ] configs/openclaw/agent-presets/
- [ ] configs/mcp-servers/local-mcps.json
- [ ] configs/mcp-servers/remote-mcps.json

### Phase 5: Onboarding
- [ ] onboarding/checklist.md (100+ Schritte)
- [ ] onboarding/api-keys/nvidia-nim.md
- [ ] onboarding/api-keys/openrouter.md
- [ ] onboarding/api-keys/[provider].md
- [ ] onboarding/accounts/github.md
- [ ] onboarding/accounts/docker-hub.md
- [ ] onboarding/accounts/cloudflare.md
- [ ] onboarding/tools/opencode-setup.md
- [ ] onboarding/tools/openclaw-setup.md
- [ ] onboarding/tools/docker-setup.md

### Phase 6: CLI TUI
- [ ] cli/main.go
- [ ] cli/cmd/onboarding.go
- [ ] cli/cmd/project.go
- [ ] cli/cmd/agent-loop.go (stub)
- [ ] cli/tui/dashboard.go
- [ ] cli/tui/onboarding-wizard.go
- [ ] cli/tui/project-wizard.go
- [ ] cli/internal/config/loader.go
- [ ] cli/internal/templates/generator.go

### Phase 7: Agent Loop Architecture
- [ ] docs/architecture/agent-loop.md
- [ ] docs/architecture/orchestrator-design.md
- [ ] docs/architecture/persistent-queue.md
- [ ] docs/architecture/status-files.md
- [ ] docs/architecture/delegation-pattern.md

### Phase 8: Migration & Cleanup
- [ ] scripts/migrate-old-docs.sh ausführen
- [ ] docs/ bereinigen (200+ → ~50)
- [ ] archive/ erstellen (Sprint 5 Packages)
- [ ] README.md (Document360 Standard)
- [ ] .gitignore aktualisieren
- [ ] Validation Scripts erstellen

---

## 📏 DATEI-METRIKEN

### Aktuelle Statistik (vor Migration)
| Kategorie | Anzahl | Größe |
|-----------|--------|-------|
| **Total Files** | ~6,219 | ~500MB |
| Markdown (.md) | ~100 | ~2MB |
| Go (.go) | ~100 | ~1MB |
| YAML (.yaml/.yml) | ~45 | ~500KB |
| JSON (.json) | ~5 | ~50KB |
| Python (.py) | ~4 | ~100KB |
| Shell (.sh) | ~9 | ~50KB |
| Binary (png, pdf, mp4) | ~50 | ~450MB |

### Nach Migration (geplant)
| Kategorie | Anzahl | Change |
|-----------|--------|--------|
| **rules/** | 9 files | NEU |
| **templates/** | 11 files | NEU |
| **configs/** | 7 files | NEU |
| **onboarding/** | 10 files | NEU |
| **cli/** | 9 files | RENAME + CLEANUP |
| **docs/** | ~50 files | REDUCED (200+ → 50) |
| **archive/** | ~60 files | NEU (Sprint 5 + old) |
| **scripts/** | ~10 files | CLEANUP |

---

## 🎯 SUCCESS CRITERIA

### Architecture Complete ✅
- [x] Audit Report erstellt (254 files analyzed)
- [x] Structure Analysis erstellt
- [x] Source of Truth Extract erstellt
- [x] ARCHITECTURE.md erstellt (diese Datei)
- [ ] Migration durchgeführt
- [ ] CLI TUI funktioniert
- [ ] Alle Tests grün

### Quality Gates
- ✅ Jede Rule-Datei 500+ Zeilen
- ✅ Alle Templates kopierfertig
- ✅ Alle Configs valide JSON
- ✅ CLI baut ohne Fehler
- ✅ TUI startet zeigt Dashboard
- ✅ Onboarding-Wizard funktioniert
- ✅ Project-Wizard erstellt Projekte

---

## 🔗 REFERENCES

- **Audit Report:** `/Users/jeremy/dev/BIOMETRICS/audit-report.md`
- **Structure Analysis:** `/Users/jeremy/dev/BIOMETRICS/structure-analysis.md`
- **Source of Truth Extract:** `/Users/jeremy/dev/BIOMETRICS/source-of-truth-extract.md`
- **Rearchitecture Plan:** `/Users/jeremy/dev/BIOMETRICS/BIOMETRICS-REARCHITECTURE-PLAN.md`
- **Original AGENTS.md:** `~/.config/opencode/AGENTS.md` (3100+ lines, 33 mandates)

---

**NEXT:** Migration beginnen mit Phase 2 (Rules)
