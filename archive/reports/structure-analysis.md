# BIOMETRICS STRUCTURE ANALYSIS

**Date:** 2026-02-20  
**Status:** PHASE 1.2 - Directory Structure Analysis Complete

---

## Current Structure Overview

### Root-Level Directories

| Directory | Files | Purpose | Status | New Home |
|-----------|-------|---------|--------|----------|
| **docs/** | 200+ | Comprehensive documentation | CHAOTIC | `docs/` (structured) |
| **biometrics-cli/** | 38 | Go CLI implementation | ACTIVE | Keep as `cli/` |
| **BIOMETRICS/** | 18 | Main project (assets, internal) | LEGACY | `archive/biometrics-main/` |
| **assets/** | 50+ | Media assets (images, videos, 3D) | OK | Keep as `assets/` |
| **global/** | 3 subdirs | Global configs for templates | EMPTY | `rules/global/` |
| **local/** | 1 subdir | Local project configs | EMPTY | `configs/local/` |
| **scripts/** | 9 | Helper scripts (Python, Bash) | MIXED | `scripts/` |
| **inputs/** | 3 subdirs | Input data for generation | OK | Keep as `inputs/` |
| **outputs/** | 3 subdirs | Generated output files | OK | Keep as `outputs/` |
| **backups/** | few | Backup files | OK | Keep as `backups/` |
| **skills/** | few | OpenCode skills | OK | Keep as `skills/` |
| **helm/** | few | Kubernetes Helm charts | OK | Keep as `helm/` |
| **.sisyphus/** | few | Sisyphus planning | OK | Keep as `.sisyphus/` |

---

## Detailed Analysis

### 1. **docs/** - 200+ Files - NEEDS STRUCTURE

**Status:** CHAOTIC - No clear organization despite existing subdirectories

**Subdirectories:**
- `advanced/` (26 files) - Advanced topics (AR-VR, Blockchain, DeFi, IoT, etc.)
- `agents/` (14 files) - Agent configurations, mandates, workflows
- `api/` (8 files) - API documentation, OpenAPI spec, Postman collection
- `architecture/` (23 files) - Architecture docs, CI/CD, backup strategies
- `best-practices/` (36 files) - Best practices documentation
- `config/` (13 files) - Configuration documentation
- `features/` (many) - Feature documentation
- `setup/` (7 files) - Setup guides

**Problems:**
1. **Massive file count** - 200+ markdown files without clear hierarchy
2. **Redundant content** - Similar topics spread across multiple directories
3. **Naming inconsistency** - Mixed naming conventions
4. **Old v1 docs** - ARCHITECTURE-v1_old.md (165KB) still present

**Proposed Solution:**
```
docs/
├── 01-getting-started/     # Quick start, setup
├── 02-architecture/       # System architecture
├── 03-api/               # API reference, OpenAPI
├── 04-agents/            # Agent configs, mandates
├── 05-features/          # Feature documentation
├── 06-best-practices/    # Best practices
├── 07-advanced/          # Advanced topics
├── 08-setup/            # Installation guides
└── 09-troubleshooting/   # Problem solving
```

---

### 2. **biometrics-cli/** - Go CLI - KEEP AS MAIN CLI

**Status:** ACTIVE - Working Go implementation with 38 directories

**Key Files:**
- `biometrics` (1.8MB binary) - Compiled CLI
- `biometrics-cli` (2.6MB) - Main CLI binary
- `go.mod`, `go.sum` - Go dependencies
- `Makefile` - Build automation
- `cmd/` - Command implementations
- `pkg/` - 30 packages - Core functionality
- `templates/` (26 dirs) - Template system
- `docs/` - CLI documentation

**Analysis:**
- This is the **MAIN working CLI** - should stay at root level as `cli/`
- Contains actual Go code vs. the legacy `BIOMETRICS/` directory
- Binary files present (should be in `.gitignore`)

**Proposed Action:** Rename to `cli/` at root

---

### 3. **BIOMETRICS/** - Legacy Project - ARCHIVE

**Status:** LEGACY - Old project structure, superceded by biometrics-cli

**Contents:**
- `biometrics/` subdir - Node.js project
- `infografik.png` (6MB) - Large media file
- `logo.png` (792KB)
- `video.mp4` (5MB)
- `praesentation.pdf` (12MB)
- `NLM-ASSETS/` - NotebookLM assets
- `internal/`, `pkg/` - Node.js packages

**Problems:**
1. **Confusing name** - Same as root directory
2. **Large media files** - Should be in assets/
3. **Superseded** - Replaced by biometrics-cli (Go)

**Proposed Action:** Move to `archive/biometrics-main/`

---

### 4. **global/** - Empty Configs - REDUCE

**Status:** EMPTY - Only 3 subdirectories with minimal content

**Contents:**
- `01-agents/` - Agent configs (empty or minimal)
- `02-models/` - Model configs (empty or minimal)
- `03-mandates/` - Mandate configs (empty or minimal)

**Analysis:**
- These are template configs for new projects
- Should be consolidated into `rules/global/`

---

### 5. **local/** - Local Project Configs

**Status:** UNDERUTILIZED

**Contents:**
- `projects/` - Project-specific configs

**Proposed Action:** Move to `configs/local/`

---

### 6. **scripts/** - Helper Scripts - KEEP

**Status:** MIXED - Some useful, some obsolete

**Scripts:**
- `setup.sh` - Main setup script
- `cosmos_video_gen.py` - Video generation (16KB)
- `blue-green-deploy.sh` - Deployment script
- `validate-config.sh` - Config validation
- `nim_engine.py` - NIM engine wrapper
- `sealcams_analysis.py` - Analysis script
- `upload-to-gitlab.sh` - GitLab upload
- `video_quality_check.py` - Quality check

**Proposed Action:** Keep as `scripts/` but clean obsolete ones

---

### 7. **assets/** - Media Assets - KEEP

**Status:** OK - Well organized by type

**Subdirs:**
- `3d/` - 3D models
- `audio/` - Audio files
- `dashboard/` - Dashboard screenshots
- `diagrams/` - Architecture diagrams
- `frames/` - Video frames
- `icons/` - Icon assets
- `images/` - General images
- `logos/` - Logo assets
- `renders/` - Rendered assets
- `videos/` - Video files

**Proposed Action:** Keep as `assets/`

---

### 8. **inputs/** & **outputs/** - I/O Directories

**Status:** OK - Used for generation workflows

**inputs/**
- `brand_assets/` - Brand materials
- `references/` - Reference files

**outputs/**
- `assets/` - Generated assets
- `videos/` - Generated videos

**Proposed Action:** Keep as `inputs/` and `outputs/`

---

### 9. Root-Level Large Files

**Problem Files (should be in .gitignore):**
- `∞Best∞Practices∞Loop.md` (69KB) - Unicode filename, should rename
- `biometrics` (1.8MB binary)
- `biometrics-cli` (2.6MB binary)
- `ratelimit.test` (6.2MB binary)
- `docker-compose.yml` - Should be in biometrics-cli/

---

## Proposed New Structure

```
BIOMETRICS/
├── cli/                           # ⭐ Main Go CLI (from biometrics-cli/)
│   ├── cmd/
│   ├── pkg/
│   ├── templates/
│   ├── docs/
│   ├── go.mod
│   ├── Makefile
│   └── README.md
│
├── rules/                         # Rules, configs, templates
│   ├── global/                    # Global agent/model configs
│   │   ├── agents/
│   │   ├── models/
│   │   └── mandates/
│   ├── tools/                    # Tool configurations
│   └── projects/                 # Project-specific rules
│
├── configs/                       # Configuration files
│   ├── local/                    # Local project configs
│   └── opencode.json             # OpenCode config
│
├── docs/                         # ⭐ Structured documentation
│   ├── 01-getting-started/
│   ├── 02-architecture/
│   ├── 03-api/
│   ├── 04-agents/
│   ├── 05-features/
│   ├── 06-best-practices/
│   ├── 07-advanced/
│   ├── 08-setup/
│   └── 09-troubleshooting/
│
├── templates/                     # Template files
│   ├── global/
│   ├── opencode/
│   └── openclaw/
│
├── assets/                        # Media assets (OK)
├── inputs/                        # Input files (OK)
├── outputs/                       # Generated files (OK)
├── scripts/                       # Helper scripts (OK)
├── skills/                        # OpenCode skills
├── helm/                          # Kubernetes configs
├── backups/
│
├── archive                       # Backup files/                       # ⭐ Archived legacy code
│   └── biometrics-main/          # From BIOMETRICS/
│       ├── biometrics/
│       ├── internal/
│       ├── pkg/
│       ├── NLM-ASSETS/
│       └── media files
│
├── .sisyphus/                    # Planning
├── .github/                      # GitHub config
│
├── README.md                     # Main README
├── CHANGELOG.md
├── docker-compose.yml
├── Makefile
└── requirements.txt
```

---

## Migration Steps

### Phase 1: Create New Structure
1. [ ] Create `rules/global/` from `global/`
2. [ ] Create `configs/` from `local/`
3. [ ] Create `archive/` directory
4. [ ] Create `docs/01-getting-started/` through `09-troubleshooting/`

### Phase 2: Move Active Content
1. [ ] Move `biometrics-cli/` → `cli/`
2. [ ] Move `global/` content → `rules/global/`
3. [ ] Move `local/` content → `configs/local/`
4. [ ] Move `.gitignore`, `requirements.txt`, etc. to root if needed

### Phase 3: Archive Legacy
1. [ ] Move `BIOMETRICS/` → `archive/biometrics-main/`
2. [ ] Move large binary files to archive or .gitignore
3. [ ] Rename `∞Best∞Practices∞Loop.md` to `best-practices-loop.md`

### Phase 4: Restructure docs/
1. [ ] Analyze all 200+ doc files
2. [ ] Categorize into 9 new subdirectories
3. [ ] Remove duplicates and old v1 docs

### Phase 5: Cleanup
1. [ ] Update all import paths and references
2. [ ] Update README.md with new structure
3. [ ] Clean .gitignore (add binaries)
4. [ ] Commit migration

---

## Critical Directories to Preserve

1. **biometrics-cli/** - Working Go CLI implementation
2. **docs/** - All documentation (200+ files)
3. **assets/** - Media assets
4. **scripts/** - Helper scripts
5. **inputs/** - Input data for generation

---

## Summary

| Category | Count | Action |
|----------|-------|--------|
| Keep as-is | 5 | assets, inputs, outputs, scripts, skills |
| Rename/Move | 3 | biometrics-cli → cli, global → rules, local → configs |
| Archive | 1 | BIOMETRICS/ → archive/ |
| Restructure | 1 | docs/ → 9 subdirectories |
| Cleanup | ~5 | Remove binaries, rename files |

**Total directories to process:** 15  
**Estimated effort:** 2-3 hours
