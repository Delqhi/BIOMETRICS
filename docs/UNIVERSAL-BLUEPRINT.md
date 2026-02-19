# 🚨 UNIVERSAL BLUEPRINT - COMPLETE OPENCODE SETUP

**Version:** 1.0 "PERFECT REPLICATION"  
**Status:** ✅ PRODUCTION-READY  
**Last Updated:** 2026-02-18  
**Source:** Complete consolidation from:
- `~/.config/opencode/` (Primary config - 286KB AGENTS.md)
- `~/.opencode/` (Legacy config)
- `~/.oh-my-opencode/` (Plugin system)

---

## 📋 TABLE OF CONTENTS

1. [Pre-Installation Requirements](#pre-installation-requirements)
2. [OpenCode Installation](#opencode-installation)
3. [Provider Setup](#provider-setup)
4. [MCP Server Configuration](#mcp-server-configuration)
5. [File System Organization](#file-system-organization)
6. [Secrets Management](#secrets-management)
7. [Verification Checklist](#verification-checklist)
8. [Troubleshooting](#troubleshooting)

---

## 📦 PRE-INSTALLATION REQUIREMENTS

### System Requirements
- **OS:** macOS 13.0+ (Ventura or later)
- **RAM:** Minimum 16GB (32GB recommended for Qwen 3.5 397B)
- **Storage:** 50GB free space minimum
- **Node.js:** v20.0+ 
- **Package Manager:** Homebrew (macOS)

### Required Tools
```bash
# Install Homebrew (if not installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Node.js v20
brew install node@20

# Install pnpm
brew install pnpm

# Install Git
brew install git
```

### Environment Variables
Create `~/.env` with:
```bash
# NVIDIA API Key (FREE tier - 40 RPM)
export NVIDIA_API_KEY="nvapi-xxx"

# OpenCode API Key (if using)
export OPENCODE_API_KEY="xxx"

# Google Antigravity (for Gemini models)
export GOOGLE_APPLICATION_CREDENTIALS="~/.config/opencode/antigravity-accounts.json"
```

---

## 🚀 OPENCODE INSTALLATION

### Step 1: Install OpenCode
```bash
# Install via Homebrew
brew install opencode

# Verify installation
opencode --version  # Should be >= 1.0.150
```

### Step 2: Initialize Configuration
```bash
# Create config directory
mkdir -p ~/.config/opencode

# Initialize opencode.json
cat > ~/.config/opencode/opencode.json << 'EOF'
{
  "$schema": "https://opencode.ai/config.schema.json",
  "provider": {},
  "model": {},
  "plugin": ["oh-my-opencode"]
}
EOF
```

### Step 3: Install Oh-My-OpenCode Plugin
```bash
# Install plugin
bunx oh-my-opencode install --no-tui --claude=no --chatgpt=no --gemini=yes

# Or with npx
npx oh-my-opencode install --no-tui --claude=no --chatgpt=no --gemini=yes
```

---

## 🔌 PROVIDER SETUP

### NVIDIA NIM (Qwen 3.5 397B) - BEST FOR CODE

**⚠️ CRITICAL:** Qwen 3.5 397B has extreme latency (70-90s)!  
**✅ SOLUTION:** Timeout MUST be set to 120000ms (120s) in OpenCode.

#### Configuration in `~/.config/opencode/opencode.json`:
```json
{
  "provider": {
    "nvidia-nim": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "NVIDIA NIM (Qwen 3.5)",
      "options": {
        "baseURL": "https://integrate.api.nvidia.com/v1",
        "timeout": 120000
      },
      "models": {
        "qwen-3.5-397b": {
          "id": "qwen/qwen3.5-397b-a17b",
          "limit": {
            "context": 262144,
            "output": 32768
          }
        }
      }
    }
  }
}
```

#### Authentication:
```bash
# Add NVIDIA provider
opencode auth add nvidia-nim

# Verify models
opencode models | grep nvidia
```

#### Rate Limits:
- **RPM:** 40 requests per minute (FREE tier)
- **HTTP 429:** Wait 60 seconds + use fallbacks

### Google Gemini (Antigravity OAuth)

#### Step 1: Install Plugin
```bash
# Add to ~/.config/opencode/opencode.json
{
  "plugin": [
    "oh-my-opencode",
    "opencode-antigravity-auth@latest"
  ]
}
```

#### Step 2: Authenticate
```bash
opencode auth login
# Select: Google → OAuth with Google (Antigravity)
# Complete sign-in in browser
```

#### Step 3: Configure Models
Add to `~/.config/opencode/opencode.json`:
```json
{
  "provider": {
    "google": {
      "npm": "@ai-sdk/google",
      "models": {
        "antigravity-gemini-3-flash": {
          "id": "gemini-3-flash-preview",
          "name": "Gemini 3 Flash (Antigravity)",
          "limit": {
            "context": 1048576,
            "output": 65536
          }
        },
        "antigravity-gemini-3-pro": {
          "id": "gemini-3-pro-preview",
          "name": "Gemini 3 Pro (Antigravity)",
          "limit": {
            "context": 2097152,
            "output": 65536
          }
        }
      }
    }
  }
}
```

### OpenCode ZEN (FREE - UNCENSORED)

```json
{
  "provider": {
    "opencode-zen": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "OpenCode ZEN (FREE)",
      "options": {
        "baseURL": "https://api.opencode.ai/v1"
      },
      "models": {
        "opencode/minimax-m2.5-free": {
          "name": "Big Pickle (OpenCode ZEN - UNCENSORED)",
          "limit": {
            "context": 200000,
            "output": 128000
          }
        },
        "zen/uncensored": {
          "name": "Uncensored (OpenCode ZEN)",
          "limit": {
            "context": 200000,
            "output": 128000
          }
        },
        "grok-code": {
          "name": "Grok Code (VIA OPENROUTER)",
          "limit": {
            "context": 2000000,
            "output": 131072
          }
        },
        "glm-4.7-free": {
          "name": "GLM 4.7 Free (VIA OPENROUTER)",
          "limit": {
            "context": 1000000,
            "output": 65536
          }
        }
      }
    }
  }
}
```

---

## 🔧 MCP SERVER CONFIGURATION

### Built-in MCP Servers (Enabled by Default)

```json
{
  "mcp": {
    "websearch": {
      "type": "local",
      "command": ["npx", "-y", "@tavily/claude-mcp"],
      "enabled": true
    },
    "context7": {
      "type": "local",
      "command": ["npx", "-y", "@anthropics/context7-mcp"],
      "enabled": true
    },
    "grep_app": {
      "type": "remote",
      "url": "https://mcp.grep.app",
      "enabled": true
    }
  }
}
```

### Docker-based MCP Servers

For Docker containers (e.g., Skyvern, Steel Browser):

1. **Create stdio wrapper** (Docker containers are HTTP APIs, NOT MCP servers!)
2. **Use `type: "local"`** in opencode.json (NOT `type: "remote"`)

Example wrapper (`mcp-wrappers/skyvern-mcp-wrapper.js`):
```javascript
#!/usr/bin/env node
const { Server } = require('@modelcontextprotocol/sdk/server/index.js');
const { StdioServerTransport } = require('@modelcontextprotocol/sdk/server/stdio.js');
const axios = require('axios');

const SKYVERN_URL = process.env.SKYVERN_URL || 'http://localhost:8030';

const server = new Server(
  { name: 'skyvern-mcp', version: '1.0.0' },
  { capabilities: { tools: {} } }
);

// Implement tool handlers...
```

---

## 📁 FILE SYSTEM ORGANIZATION

### Primary Directories
```
/Users/jeremy/
├── .config/opencode/          # PRIMARY CONFIG (Source of Truth)
│   ├── opencode.json          # Main configuration
│   ├── AGENTS.md              # Executive mandates (this document)
│   ├── oh-my-opencode.json    # Plugin config
│   └── antigravity-accounts.json # OAuth tokens
├── .opencode/                 # LEGACY (preserved, not edited)
├── .oh-my-opencode/           # Plugin system
└── dev/
    ├── sin-code/              # MAIN workspace
    ├── SIN-Solver/            # AI automation project
    └── BIOMETRICS/            # This project
```

### Project Documentation Structure
```
BIOMETRICS/
├── docs/
│   ├── setup/                 # Installation guides
│   ├── config/                # Configuration references
│   ├── agents/                # Agent documentation
│   ├── best-practices/        # Mandates & workflows
│   ├── architecture/          # System architecture
│   ├── features/              # Feature documentation
│   ├── advanced/              # Advanced topics
│   ├── data/                  # Data engineering
│   └── devops/                # DevOps & CI/CD
├── biometrics-cli/            # CLI tool
└── README.md                  # Main entry point
```

---

## 🔐 SECRETS MANAGEMENT

### NEVER Commit Secrets
```bash
# Add to .gitignore
echo ".env" >> .gitignore
echo "*.key" >> .gitignore
echo "*.pem" >> .gitignore
echo "credentials.json" >> .gitignore
echo "antigravity-accounts.json" >> .gitignore
```

### Use Environment Variables
```bash
# ~/.zshrc or ~/.bashrc
export NVIDIA_API_KEY="nvapi-xxx"
export OPENCODE_API_KEY="xxx"
export GOOGLE_APPLICATION_CREDENTIALS="~/.config/opencode/antigravity-accounts.json"

# Reload shell
source ~/.zshrc
```

### File Permissions
```bash
chmod 600 ~/.config/opencode/opencode.json
chmod 600 ~/.config/opencode/antigravity-accounts.json
chmod 600 ~/.openclaw/openclaw.json
```

---

## ✅ VERIFICATION CHECKLIST

### OpenCode Verification
```bash
# Check version
opencode --version  # >= 1.0.150

# List models
opencode models

# Check providers
opencode models | grep -E "(nvidia|google|opencode-zen)"
```

### Plugin Verification
```bash
# Check plugin loaded
opencode --version | grep -i "oh-my-opencode"

# Check agents available
opencode --help | grep -A 10 "Agents"
```

### MCP Server Verification
```bash
# Test websearch
curl -X POST http://localhost:8080/mcp/websearch \
  -H "Content-Type: application/json" \
  -d '{"query": "test"}'

# Test context7
curl -X POST http://localhost:8080/mcp/context7 \
  -H "Content-Type: application/json" \
  -d '{"library": "react"}'
```

### Provider Health Check
```bash
# NVIDIA API
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
  https://integrate.api.nvidia.com/v1/models

# Should return list of available models
```

---

## 🐛 TROUBLESHOOTING

### Common Issues

#### 1. Models Not Found
**Problem:** `opencode models` shows no models  
**Solution:**
```bash
# 1. Check provider config
cat ~/.config/opencode/opencode.json | jq '.provider'

# 2. Re-authenticate provider
opencode auth add nvidia-nim

# 3. Restart terminal
exec zsh
```

#### 2. Timeout in Config (MANDATE 0.35)
**Problem:** Qwen 3.5 397B times out  
**Solution:**
```bash
# Check for timeout entries (MUST BE EMPTY!)
grep -r "timeout" ~/.config/opencode/opencode.json

# If found, REMOVE them immediately!
# Only NVIDIA NIM should have timeout: 120000
```

#### 3. HTTP 429 Rate Limit
**Problem:** NVIDIA API returns 429 Too Many Requests  
**Solution:**
- Wait 60 seconds
- Implement exponential backoff
- Use fallback models (Gemini, Grok)

#### 4. Plugin Not Loading
**Problem:** oh-my-opencode features not available  
**Solution:**
```bash
# 1. Check plugin in config
cat ~/.config/opencode/opencode.json | jq '.plugin'

# 2. Reinstall plugin
bunx oh-my-opencode install

# 3. Restart OpenCode
opencode --version
```

#### 5. OAuth Authentication Fails
**Problem:** Google Antigravity OAuth fails  
**Solution:**
```bash
# 1. Clear old tokens
rm ~/.config/opencode/antigravity-accounts.json

# 2. Re-authenticate
opencode auth login

# 3. Use PRIVATE Gmail (NOT Workspace!)
```

---

## 📚 ADDITIONAL RESOURCES

### Documentation
- **Main AGENTS.md:** `docs/agents/AGENTS-MANDATES.md` (286KB)
- **Oh-My-OpenCode:** `docs/agents/OH-MY-OPENCODE-AGENTS.md`
- **Best Practices:** `docs/best-practices/`
- **Architecture:** `docs/architecture/`

### Configuration Files
- **OpenCode:** `~/.config/opencode/opencode.json`
- **Plugin:** `~/.config/opencode/oh-my-opencode.json`
- **OpenClaw:** `~/.openclaw/openclaw.json`

### Elite Guides
- **Antigravity Plugin:** `/Users/jeremy/dev/sin-code/Guides/01-antigravity-plugin-guide.md` (783 lines)
- **Blueprint Template:** `~/.opencode/blueprint-vorlage.md` (500+ lines)

---

## 🎯 QUICK START COMMANDS

```bash
# 1. Install OpenCode
brew install opencode

# 2. Install plugin
bunx oh-my-opencode install --no-tui --claude=no --chatgpt=no --gemini=yes

# 3. Authenticate
opencode auth login

# 4. Configure providers
# Edit ~/.config/opencode/opencode.json

# 5. Verify
opencode models

# 6. Start coding
opencode
```

---

**END OF UNIVERSAL BLUEPRINT**

*"Omniscience is not a goal; it is our technical starting point."*

**Document Statistics:**
- Total Lines: 500+
- Providers: 3 (NVIDIA, Google, OpenCode ZEN)
- MCP Servers: 3 (websearch, context7, grep_app)
- Status: ✅ PRODUCTION-READY
