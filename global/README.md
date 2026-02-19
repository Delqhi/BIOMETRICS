# Global Configuration

**Purpose:** Centralized configuration for all BIOMETRICS agents, models, and mandates.

**Status:** Active
**Last Updated:** 2026-02-19
**Version:** 1.0.0
**Compliance:** MANDATE 0.0-0.36

---

## Purpose

This directory contains global configurations that apply across all BIOMETRICS projects:

- **Agents:** AI agent definitions, roles, and capabilities
- **Models:** AI model configurations, providers, and fallbacks
- **Mandates:** Core rules and best practices (33 mandates)

---

## Directory Structure

```
global/
├── README.md                      # This file
├── 01-agents/                    # Agent configurations
│   ├── README.md                 # Agent overview
│   └── [agent-configs].md       # Individual agent configs
├── 02-models/                    # AI model definitions
│   ├── README.md                 # Model configurations
│   └── [model-configs].md       # Individual model configs
└── 03-mandates/                 # Core mandates
    ├── README.md                 # Mandate index
    └── [mandate-files].md       # Individual mandates
```

---

## Enterprise Practices 2026

1. **Immutability of Knowledge:** No existing line may be deleted or overwritten with less information. All modifications must be additive enhancements.

2. **Parallel Agent Execution:** Agents must never wait for each other. Use `run_in_background=true` for all delegate_task calls. Maximum 3 agents with different models in parallel.

3. **File Existence Verification:** Before creating any file, agents must verify existence using glob/ls. Never create duplicates. Always read existing files completely before modification.

4. **Zero-Defect Validation:** Every deliverable requires verification evidence. Tests must pass, screenshots must be captured, and commit must follow conventional commits format.

5. **Documentation Sovereignty:** Every change must be documented in lastchanges.md (append-only), AGENTS.md (project-specific), and README.md (if applicable).

---

## Files Reference

| # | File | Purpose | Lines |
|---|------|---------|-------|
| 1 | README.md | This file - global configuration index | 100+ |
| 2 | 01-agents/README.md | Agent definitions and capabilities | - |
| 3 | 02-models/README.md | Model configurations and providers | - |
| 4 | 03-mandates/README.md | Core mandates index | - |

---

## Usage

### For AI Agents

Before starting work, agents must read these files in order:

1. `03-mandates/README.md` - Core rules (MANDATORY)
2. `02-models/README.md` - Model assignments
3. `01-agents/README.md` - Agent capabilities

### For Developers

- **Configuration:** All global settings reside here
- **Standards:** Enterprise-wide conventions enforced
- **Compliance:** MANDATE 0.0-0.36 enforced

---

## Model Assignments

| Agent | Model | Provider | Max Parallel |
|-------|-------|----------|--------------|
| sisyphus | qwen/qwen3.5-397b-a17b | NVIDIA NIM | 1 |
| prometheus | qwen/qwen3.5-397b-a17b | NVIDIA NIM | 1 |
| oracle | qwen/qwen3.5-397b-a17b | NVIDIA NIM | 1 |
| atlas | kimi-k2.5 | Moonshot AI | 1 |
| librarian | zen/big-pickle | OpenCode ZEN | 1 |
| explore | zen/big-pickle | OpenCode ZEN | 1 |

---

## Validation Command

```bash
# Validate global configuration integrity
cd /Users/jeremy/dev/BIOMETRICS/global
find . -name "*.md" -exec wc -l {} \; | awk '{sum+=$1} END {print "Total lines:", sum}'
ls -la 01-agents/ 02-models/ 03-mandates/
```

---

## Related Documentation

- [Architecture](../docs/architecture/)
- [Best Practices](../docs/best-practices/)
- [Setup Guide](../docs/setup/)
- [Parent README](../README.md)

---

**Maintainer:** BIOMETRICS Team  
**Contact:** support@biometrics.dev  
**License:** MIT
