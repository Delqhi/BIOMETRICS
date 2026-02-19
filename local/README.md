# 📂 Local Projects

**Purpose:** Project-specific configurations and isolated workspaces.

**Status:** ✅ Active  
**Last Updated:** 2026-02-19  
**Version:** 1.0.0

---

## 📁 Directory Structure

```
local/
├── README.md              # This file
└── projects/             # Individual projects
    ├── README.md
    └── [project-name]/
        ├── AGENTS.md
        ├── lastchanges.md
        └── [project-files]
```

---

## 🎯 Purpose

This directory contains **project-specific configurations**:

- **Isolation:** Each project has its own workspace
- **Autonomy:** Projects can have custom rules
- **Traceability:** Clear project boundaries

---

## 🚀 Usage

### Creating New Projects

```bash
# Create new project directory
mkdir -p local/projects/[project-name]

# Add required files
cd local/projects/[project-name]
touch AGENTS.md lastchanges.md README.md
```

### Project Requirements

Every project **MUST** have:

1. ✅ `AGENTS.md` - Local agent rules
2. ✅ `lastchanges.md` - Change log (append-only)
3. ✅ `README.md` - Project overview

---

## 📋 Project Template

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

## 🔗 Related Documentation

- [Global Config](../global/)
- [CLI Tools](../biometrics-cli/)
- [Architecture](../docs/architecture/)

---

**Maintained by:** BIOMETRICS Team  
**Contact:** support@biometrics.dev
