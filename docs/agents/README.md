# 🤖 BIOMETRICS Documentation - Agents

**Agent definitions, mandates, and skills documentation.**

---

## 📁 Agent Documentation

| File | Description | Lines |
|------|-------------|-------|
| [AGENTS-MANDATES.md](AGENTS-MANDATES.md) | Executive mandates (V20.0) | 286KB |
| [AGENTS-GLOBAL.md](AGENTS-GLOBAL.md) | Global governance | Original |
| [AGENTS-PLAN.md](AGENTS-PLAN.md) | Agent planning guide | Original |
| [USER-PLAN.md](USER-PLAN.md) | User interaction plan | Original |
| [MCP.md](MCP.md) | MCP server documentation | Original |
| [NOTEBOOKLM.md](NOTEBOOKLM.md) | NotebookLM integration | Original |
| [CONTEXT.md](CONTEXT.md) | Context management | Original |
| [COMMANDS.md](COMMANDS.md) | Available commands | Original |
| [WORKFLOW.md](WORKFLOW.md) | Agent workflows | Original |
| [OH-MY-OPENCODE-AGENTS.md](OH-MY-OPENCODE-AGENTS.md) | Plugin agents | 206 lines |

---

## 🎯 Built-in Agents

### Sisyphus (Primary Orchestrator)
- **Model:** Claude Opus 4.5
- **Role:** Planning, delegation, execution
- **Features:** Background tasks, parallel execution

### Oracle
- **Model:** GPT-5.2
- **Role:** Architecture, debugging, strategy
- **Use:** High-IQ reasoning, deep analysis

### Librarian
- **Model:** GLM 4.7 Free
- **Role:** Multi-repo analysis, docs lookup
- **Use:** GitHub research, implementation examples

### Explore
- **Model:** Grok Code / Gemini Flash
- **Role:** Fast codebase exploration
- **Use:** Pattern matching, file discovery

### Frontend UI/UX Engineer
- **Model:** Gemini 3 Pro
- **Role:** UI generation, design
- **Use:** Creative, beautiful UI code

---

## 🔧 Skills

### Playwright
- Browser automation
- Web scraping
- Testing & screenshots

### Git Master
- Atomic commits
- Rebase/squash
- History search (blame, bisect)

---

## 📋 Agent Mandates (Top 10)

1. **PARALLEL EXECUTION** - Always use `run_in_background=true`
2. **SEARCH BEFORE CREATE** - Use `glob()`, `grep()` first
3. **VERIFY-THEN-EXECUTE** - Check with LSP diagnostics
4. **GIT COMMIT DISCIPLINE** - After every significant change
5. **FREE-FIRST PHILOSOPHY** - Self-hosted, free tiers
6. **RESOURCE PRESERVATION** - Never delete configs
7. **NO-SCRIPT MANDATE** - Use AI agents for everything
8. **NLM DUPLICATE PREVENTION** - List before upload
9. **TODO DISCIPLINE** - Create todos for multi-step tasks
10. **PERFORMANCE FIRST** - Native CDP over Playwright

---

## 🚀 Usage Examples

```typescript
// Delegate to specialized agent
sisyphus_task(category="visual", prompt="Create dashboard")

// Background exploration
delegate_task(category="explore", run_in_background=true)

// Use skill
sisyphus_task(category="quick", load_skills=["git-master"])
```

---

**Last Updated:** 2026-02-18  
**Status:** ✅ Production-Ready
