# OPENCODE Configuration Guide

**Project:** BIOMETRICS  
**Last Updated:** 2026-02-21  
**Status:** Active

---

## Overview

This document describes the OpenCode configuration for BIOMETRICS project, including agent configurations, model assignments, and provider setup.

## Configuration File

**Location:** `.opencode/opencode.json`

### Model Configuration

| Setting | Value |
|---------|-------|
| Default Model | `google/gemini-3.1-pro-preview-customtools` |
| Small Model | `google/gemini-3.1-pro-preview-customtools` |
| Default Agent | `sisyphus` |

### Plugins

```json
"plugin": [
  "opencode-antigravity-auth@latest",
  "oh-my-opencode",
  "opencode-qwencode-auth"
]
```

---

## Agent Configurations

### Core Agents

| Agent | Model | Purpose |
|-------|-------|---------|
| **sisyphus** | `google/gemini-3.1-pro-preview-customtools` | Main Coder |
| **sisyphus-junior** | `google/gemini-3.1-pro-preview-customtools` | Junior Coder |
| **atlas** | `google/gemini-3.1-pro-preview` | Heavy Lifting |
| **Delqhi** | `google/gemini-3.1-pro-preview-customtools` | Main Orchestrator |

### Specialized Agents

| Agent | Model | Purpose |
|-------|-------|---------|
| **general** | `google/gemini-3.1-pro-preview-customtools` | General tasks |
| **plan** | `google/gemini-3.1-pro-preview` | Strategic Planning |
| **build** | `google/gemini-3.1-pro-preview-customtools` | Code Building |
| **explore** | `google/gemini-3-flash-preview` | Code Discovery |
| **librarian** | `google/gemini-3-flash-preview` | Documentation |
| **oracle** | `google/gemini-3.1-pro-preview-customtools` | Architecture Review |
| **metis** | `google/gemini-3-flash-preview` | Planning Support |
| **momus** | `google/gemini-3-flash-preview` | Review/Quality |
| **artistry** | `google/gemini-3.1-pro-preview` | Creative Tasks |
| **deep** | `google/gemini-3.1-pro-preview-customtools` | Deep Analysis |
| **ultrabrain** | `google/gemini-3.1-pro-preview` | Complex Reasoning |
| **visual-engineering** | `google/gemini-3.1-pro-preview` | UI/UX Design |
| **quick** | `google/gemini-3-flash-preview` | Quick Tasks |
| **unspecified-low** | `google/gemini-3.1-pro-preview-customtools` | Low Complexity |
| **unspecified-high** | `google/gemini-3.1-pro-preview-customtools` | High Complexity |
| **writing** | `google/gemini-3.1-pro-preview` | Documentation |
| **td-agent** | `google/gemini-3.1-pro-preview` | Media Tasks |
| **document-writer** | `google/gemini-3.1-pro-preview` | Doc Generation |
| **frontend-ui-ux-engineer** | `google/gemini-3.1-pro-preview` | Frontend Development |

---

## Provider Configuration

### Google Gemini Models

| Model ID | Name | Context | Output | Use Case |
|----------|------|---------|--------|----------|
| `gemini-3.1-pro-preview` | Gemini 3.1 Pro (Standard) | 2M | 64K | Thinker, Planning |
| `gemini-3.1-pro-preview-customtools` | Gemini 3.1 Pro (Worker) | 2M | 64K | Code, Tools |
| `gemini-3-flash-preview` | Gemini 3 Flash | 1M | 64K | Fast, Quick Tasks |
| `gemini-3-pro-preview` | Gemini 3 Pro | 2M | 64K | Advanced Reasoning |
| `antigravity-gemini-3-flash` | Gemini 3 Flash (Antigravity) | 1M | 64K | With OAuth |
| `antigravity-gemini-3-pro` | Gemini 3 Pro (Antigravity) | 1M | 64K | With OAuth |

### Alternative Providers

| Provider | Model | Context | Output |
|----------|-------|---------|--------|
| **NVIDIA NIM** | `qwen-3.5-397b` | 262K | 32K |
| **OpenCode ZEN** | `zen/big-pickle` | 200K | 128K |
| **XiaoMi** | `mimo-v2-flash` | 1M | 64K |
| **Streamlake** | `kat-coder-pro-v1` | 2M | 128K |

---

## Environment Variables

**Location:** `.env`

```bash
# Google API Key (for Gemini models)
GOOGLE_API_KEY=AIzaSyAVWKxhWCT64Z0VxxmskWzPNTwfWVecC_U

# NVIDIA API Key (for Qwen models)
NVIDIA_API_KEY=nvapi-xxx

# Tavily API Key (for web search)
TAVILY_API_KEY=tvly-dev-xxx
```

---

## MCP Servers

### Local MCP Servers

| Server | Command | Status |
|--------|---------|--------|
| **serena** | `uvx serena start-mcp-server` | ✅ Enabled |
| **tavily** | `npx @tavily/claude-mcp` | ✅ Enabled |
| **canva** | `npx @canva/claude-mcp` | ✅ Enabled |
| **context7** | `npx @anthropics/context7-mcp` | ✅ Enabled |
| **skyvern** | `python -m skyvern.mcp.server` | ✅ Enabled |
| **chrome-devtools** | `npx @anthropics/chrome-devtools-mcp` | ✅ Enabled |
| **singularity** | `node singularity.js mcp` | ✅ Enabled |

### Remote MCP Servers

| Server | URL | Status |
|--------|-----|--------|
| **linear** | `https://mcp.linear.app/sse` | ✅ Enabled |
| **gh_grep** | `https://mcp.grep.app` | ✅ Enabled |
| **sin_social** | `https://sin-social.delqhi.com` | ✅ Enabled |
| **sin_deep_research** | `https://sin-research.delqhi.com` | ✅ Enabled |
| **sin_video_gen** | `https://sin-video.delqhi.com` | ✅ Enabled |
| **sin_plugins** | `https://sin-plugins.delqhi.com` | ✅ Enabled |
| **sin_api_coordinator** | `https://api-coordinator.delqhi.com` | ✅ Enabled |
| **sin_clawdbot** | `https://clawdbot-ceo.delqhi.com` | ✅ Enabled |
| **sin_survey_worker** | `https://survey.delqhi.com` | ✅ Enabled |
| **sin_captcha_worker** | `https://captcha.delqhi.com` | ✅ Enabled |
| **sin_website_worker** | `https://website-worker.delqhi.com` | ✅ Enabled |
| **sin_browser_agent_zero** | `https://agent-zero.delqhi.com` | ✅ Enabled |
| **sin_browser_steel** | `https://steel.delqhi.com` | ✅ Enabled |
| **sin_browser_skyvern** | `https://skyvern.delqhi.com` | ✅ Enabled |
| **sin_browser_stagehand** | `https://stagehand.delqhi.com` | ✅ Enabled |

---

## Model Selection Guidelines

### When to Use Gemini 3.1 Pro (Standard)

- Strategic planning and architecture
- Complex reasoning tasks
- Document writing
- UI/UX design
- Creative tasks

**Model:** `google/gemini-3.1-pro-preview`

### When to Use Gemini 3.1 Pro (Custom Tools)

- Code implementation
- Tool-heavy tasks
- Agent orchestration
- Heavy lifting
- Complex debugging

**Model:** `google/gemini-3.1-pro-preview-customtools`

### When to Use Gemini 3 Flash

- Quick tasks
- Simple code fixes
- Fast exploration
- Documentation lookup

**Model:** `google/gemini-3-flash-preview`

---

## Important Notes

### API Key Security

- **NEVER** commit API keys to git
- Store in `.env` file (already in .gitignore)
- Use environment variables in opencode.json

### Model Selection

- **Worker agents** use `customtools` variant for tool access
- **Thinker agents** use standard variant for reasoning
- **Quick agents** use Flash for speed

### Rate Limits

- Gemini API: 60 requests/minute (varies by tier)
- NVIDIA NIM: 40 requests/minute (free tier)
- Tavily: 1000 requests/month (free tier)

---

## Troubleshooting

### JSON Validation Errors

If you encounter JSON errors in opencode.json:
1. Run: `cat .opencode/opencode.json | python3 -m json.tool > /dev/null`
2. Fix any syntax errors reported
3. Ensure no trailing commas

### API Key Issues

If models are not working:
1. Verify API key in `.env` file
2. Test with: `opencode auth status`
3. Re-authenticate if needed: `opencode auth add google`

### MCP Server Connection

If MCP servers fail to connect:
1. Check if server is running: `ps aux | grep [server-name]`
2. Restart server if needed
3. Check logs for error messages

---

## Related Documentation

- [OH-MY-OPENCODE.md](OH-MY-OPENCODE.md) - Custom OpenCode enhancements
- [AGENTS-PLAN.md](../AGENTS-PLAN.md) - Agent task planning
- [ARCHITECTURE.md](../ARCHITECTURE.md) - System architecture

---

**End of Document**
