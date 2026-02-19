# 🤖 Agent Configurations

**Purpose:** AI agent definitions, roles, and capabilities.

**Status:** ✅ Active  
**Last Updated:** 2026-02-19

---

## 📋 Available Agents

| Agent | Role | Model | Max Parallel |
|-------|------|-------|--------------|
| **Sisyphus** | Main Coder | Qwen 3.5 397B | 1 |
| **Prometheus** | Planner | Qwen 3.5 397B | 1 |
| **Oracle** | Architect | Qwen 3.5 397B | 1 |
| **Atlas** | Heavy Lifting | Kimi K2.5 | 1 |
| **Librarian** | Documentation | OpenCode ZEN | Unlimited |
| **Explore** | Code Discovery | OpenCode ZEN | Unlimited |

---

## ⚠️ Critical Rules

1. **MAX 1 agent per model** at any time
2. **MAX 3 agents parallel** total (different models)
3. **NEVER** use same model for multiple agents

---

## 🔧 Configuration

See [../../docs/agents/AGENT-MODEL-MAPPING.md](../../docs/agents/AGENT-MODEL-MAPPING.md)
