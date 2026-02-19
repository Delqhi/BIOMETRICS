# 🌐 Global Configuration

**Purpose:** Centralized configuration for all BIOMETRICS agents, models, and mandates.

**Status:** ✅ Active  
**Last Updated:** 2026-02-19  
**Version:** 1.0.0

---

## 📁 Directory Structure

```
global/
├── README.md              # This file
├── 01-agents/            # Agent configurations
│   ├── README.md
│   └── [agent-configs].md
├── 02-models/            # AI model definitions
│   ├── README.md
│   └── [model-configs].md
└── 03-mandates/          # Core mandates
    ├── README.md
    └── [mandate-files].md
```

---

## 🎯 Purpose

This directory contains **global configurations** that apply across all BIOMETRICS projects:

- **Agents:** AI agent definitions, roles, and capabilities
- **Models:** AI model configurations, providers, and fallbacks
- **Mandates:** Core rules and best practices (33 mandates)

---

## 🚀 Usage

### For KI Agents

Before starting work, **READ** these files in order:

1. `03-mandates/` - Core rules (MANDATORY)
2. `02-models/` - Model assignments
3. `01-agents/` - Agent capabilities

### For Developers

- **Configuration:** All global settings live here
- **Standards:** Enterprise-wide conventions
- **Compliance:** Mandate enforcement

---

## 📋 Key Files

| File | Purpose | Required |
|------|---------|----------|
| `01-agents/README.md` | Agent overview | ✅ Yes |
| `02-models/README.md` | Model configurations | ✅ Yes |
| `03-mandates/README.md` | Core mandates | ✅ Yes |

---

## 🔗 Related Documentation

- [Architecture](../docs/architecture/)
- [Best Practices](../docs/best-practices/)
- [Setup Guide](../docs/setup/)

---

**Maintained by:** BIOMETRICS Team  
**Contact:** support@biometrics.dev
