# OH-MY-OPENCODE Configuration Guide

**Project:** BIOMETRICS  
**Last Updated:** 2026-02-21  
**Status:** Active

---

## Overview

This document describes the oh-my-opencode plugin configuration and custom enhancements for the BIOMETRICS project. OH-MY-OPENCODE extends OpenCode with additional functionality for enterprise workflows.

## Plugin Configuration

### Enabled Plugins

```json
"plugin": [
  "opencode-antigravity-auth@latest",
  "oh-my-opencode",
  "opencode-qwencode-auth"
]
```

## opencode-antigravity-auth

### Description

The Antigravity Auth plugin provides Google OAuth authentication for Gemini models. This enables access to premium Gemini features through the user's Google account.

### Authentication Commands

```bash
# Start OAuth flow (use PRIVATE Gmail!)
opencode auth login

# Remove credentials
opencode auth logout

# Refresh tokens
opencode auth refresh

# Show authentication status
opencode auth status
```

### Important Notes

- **Use private Gmail account** (e.g., aimazing2024@gmail.com)
- **NOT** Google Workspace accounts
- Tokens are stored securely in `~/.config/opencode/antigravity-accounts.json`
- File permissions should be 600: `chmod 600 ~/.config/opencode/antigravity-accounts.json`

### Configuration

The plugin automatically handles:
- OAuth token refresh
- Token storage encryption
- Multi-account support

## oh-my-opencode

### Description

OH-MY-OPENCODE provides enhanced OpenCode experience with additional features including custom commands, improved workflows, and extended functionality.

### Features

#### Enhanced Agent Management

- Improved agent session handling
- Extended context preservation
- Better error handling and recovery

#### Custom Workflows

- Pre-configured task templates
- Automated documentation generation
- Quality gate enforcement

#### Integration Extensions

- Extended MCP server support
- Custom tool wrappers
- Enhanced API coordination

## opencode-qwencode-auth

### Description

Authentication plugin for QwenCode models through NVIDIA NIM provider.

### Authentication Commands

```bash
# Authenticate with NVIDIA
opencode auth add nvidia-nim

# Verify authentication
opencode models | grep nvidia
```

### Configuration

```json
{
  "provider": {
    "nvidia-nim": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "NVIDIA NIM (Qwen 3.5)",
      "options": {
        "baseURL": "https://integrate.api.nvidia.com/v1"
      }
    }
  }
}
```

## Model Authentication Matrix

| Model | Provider | Authentication Required | Plugin |
|-------|----------|------------------------|--------|
| Gemini 3.1 Pro | Google | Yes (OAuth) | opencode-antigravity-auth |
| Gemini 3 Flash | Google | Yes (OAuth) | opencode-antigravity-auth |
| Qwen 3.5 | NVIDIA | Yes (API Key) | opencode-qwencode-auth |
| OpenCode ZEN | OpenCode | No | - |
| XiaoMi MIMO | XiaoMi | API Key | - |
| Streamlake | Streamlake | API Key | - |

## Authentication Workflow

### Step 1: Google Authentication (for Gemini)

```bash
# Login with private Gmail
opencode auth login

# Verify status
opencode auth status

# Should show: "Authenticated: yes"
```

### Step 2: NVIDIA Authentication (for Qwen)

```bash
# Add NVIDIA provider
opencode auth add nvidia-nim

# Verify
opencode models | grep nvidia
# Should show: qwen-3.5-397b available
```

### Step 3: Verify All Models

```bash
# List all available models
opencode models

# Should show all configured providers and models
```

## Environment Setup

### Required API Keys

| Key | Provider | Purpose | How to Get |
|-----|----------|---------|------------|
| GOOGLE_API_KEY | Google | Gemini API | console.cloud.google.com |
| NVIDIA_API_KEY | NVIDIA | Qwen Models | https://build.nvidia.com/ |
| TAVILY_API_KEY | Tavily | Web Search | https://tavily.com/ |

### Setting Up Environment Variables

```bash
# Add to ~/.zshrc or ~/.bashrc
export GOOGLE_API_KEY="AIzaSy..."
export NVIDIA_API_KEY="nvapi-..."
export TAVILY_API_KEY="tvly-..."

# Reload shell
source ~/.zshrc
```

### .env File

Create a `.env` file in the project root:

```bash
# Google
GOOGLE_API_KEY=AIzaSyAVWKxhWCT64Z0VxxmskWzPNTwfWVecC_U

# NVIDIA
NVIDIA_API_KEY=nvapi-xxx

# Tavily
TAVILY_API_KEY=tvly-dev-xxx
```

**Note:** The `.env` file is already in `.gitignore` - do NOT commit it!

## Troubleshooting Authentication

### Google Auth Issues

**Problem:** "Authentication failed" or "Token expired"

**Solution:**
```bash
# Logout and re-authenticate
opencode auth logout
opencode auth login
```

### NVIDIA Auth Issues

**Problem:** "Invalid API key" or "Rate limit exceeded"

**Solution:**
```bash
# Re-authenticate
opencode auth add nvidia-nim

# Wait 60 seconds if rate limited (429 error)
# NVIDIA free tier: 40 requests/minute
```

### Model Not Available

**Problem:** Model shows in list but fails to respond

**Solution:**
1. Check API key is set: `echo $GOOGLE_API_KEY`
2. Verify auth status: `opencode auth status`
3. Try different model variant

## Security Best Practices

### API Key Protection

- **NEVER** commit keys to git
- Use `.env` files (in `.gitignore`)
- Rotate keys periodically
- Use least-privilege accounts

### Token Storage

- File permissions: 600
- Encrypted storage where possible
- Regular token rotation

### Access Control

- Use separate accounts for dev/prod
- Monitor usage for anomalies
- Set up billing alerts

## Related Documentation

- [OPENCODE.md](OPENCODE.md) - Main OpenCode configuration
- [AGENTS-PLAN.md](../AGENTS-PLAN.md) - Agent task planning
- [ARCHITECTURE.md](../ARCHITECTURE.md) - System architecture

---

**End of Document**
