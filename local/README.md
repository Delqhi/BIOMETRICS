# Local Projects

**Purpose:** Project-specific configurations and isolated workspaces.

**Status:** Active
**Last Updated:** 2026-02-19
**Version:** 1.0.0
**Compliance:** MANDATE 0.22-0.23

---

## Purpose

This directory contains project-specific configurations:

- **Isolation:** Each project has its own workspace
- **Autonomy:** Projects can have custom rules
- **Traceability:** Clear project boundaries

---

## Directory Structure

```
local/
├── README.md                      # This file
└── projects/                      # Individual projects
    ├── README.md                  # Project index
    └── [project-name]/            # Project workspace
        ├── AGENTS.md              # Local agent rules
        ├── lastchanges.md         # Change log (append-only)
        └── [project-files]       # Project-specific files
```

---

## Enterprise Practices 2026

1. **Project Knowledge Sovereignty:** Each project must maintain its own AGENTS.md with local conventions, naming patterns, and specific rules (MANDATE 0.22).

2. **Photographic Memory:** Every project must maintain lastchanges.md with append-only logging. All sessions, observations, errors, and solutions must be documented chronologically (MANDATE 0.23).

3. **Zero Collision:** Before starting work, agents must check existing project workspaces to avoid conflicts. Use workspace tracking format: `{task};{id}-{path}-{status}`.

4. **Immutability Preservation:** Project files follow the same immutability rules as global. Never delete, only add. Deprecate with labels, never remove.

5. **Session Continuity:** Before starting work, agents must read the project's lastchanges.md and AGENTS.md to understand context and prior decisions.

---

## Files Reference

| # | File | Purpose | Lines |
|---|------|---------|-------|
| 1 | README.md | This file - local projects index | 100+ |
| 2 | projects/README.md | Project index | - |
| 3 | projects/[name]/AGENTS.md | Local agent rules | - |
| 4 | projects/[name]/lastchanges.md | Change log | - |

---

## Usage

### Creating New Projects

```bash
# Create new project directory
mkdir -p local/projects/[project-name]

# Add required files
cd local/projects/[project-name]
touch AGENTS.md lastchanges.md README.md
```

### Project Requirements

Every project MUST have:

1. `AGENTS.md` - Local agent rules
2. `lastchanges.md` - Change log (append-only)
3. `README.md` - Project overview

---

## Project Template

```markdown
# [Project Name]

**Status:** [Active | Inactive | Archived]
**Created:** [YYYY-MM-DD]
**Last Updated:** [YYYY-MM-DD]

## Overview

[Project description]

## Local Rules

[Project-specific agent rules]

## Documentation

- [AGENTS.md](AGENTS.md)
- [lastchanges.md](lastchanges.md)
```

---

## Validation Command

```bash
# Validate local configuration integrity
cd /Users/jeremy/dev/BIOMETRICS/local
find . -name "*.md" -exec wc -l {} \; | awk '{sum+=$1} END {print "Total lines:", sum}'
ls -la projects/
```

---

## Related Documentation

- [Global Config](../global/)
- [CLI Tools](../biometrics-cli/)
- [Architecture](../docs/architecture/)
- [Parent README](../README.md)

---

**Maintainer:** BIOMETRICS Team  
**Contact:** support@biometrics.dev  
**License:** MIT
