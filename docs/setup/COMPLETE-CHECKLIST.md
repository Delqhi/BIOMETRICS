# 🚨 UNIVERSAL SETUP CHECKLIST - COMPLETE GUIDE

**Version:** 1.0 "PERFECT REPLICATION"  
**Status:** ACTIVE - MUST FOLLOW FOR NEW MAC SETUP  
**Source:** Consolidated from `~/.config/opencode/`, `~/.opencode/`, `~/.oh-my-opencode/`

---

## 📋 PRE-INSTALLATION CHECKLIST

### 1. System Requirements
- [ ] macOS 14.0+ (Sonoma or later)
- [ ] Homebrew installed: `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`
- [ ] Node.js 20+: `brew install node`
- [ ] Git configured: `git config --global user.name` & `git config --global user.email`
- [ ] SSH keys setup: `ssh-keygen -t ed25519`

### 2. Directory Structure Setup
```bash
# Create standard directory structure
mkdir -p ~/dev/sin-code
mkdir -p ~/dev/SIN-Solver
mkdir -p ~/.config/opencode
mkdir -p ~/Bilder/AI-Screenshots/{playwright,skyvern,steel,stagehand,opencode}
```

---

## 🔧 OPENCODE INSTALLATION

### Step 1: Install OpenCode
```bash
brew install opencode
```

### Step 2: Initialize Configuration
```bash
# Create config directory
mkdir -p ~/.config/opencode

# Initialize with default config
cd ~/.config/opencode
```

### Step 3: Provider Authentication
```bash
# NVIDIA NIM (PRIMARY)
opencode auth add nvidia-nim
# Follow OAuth flow

# Google Antigravity (if using Gemini)
opencode auth add google

# OpenCode ZEN (FREE)
# No auth required - free tier
```

### Step 4: Verify Installation
```bash
opencode --version
opencode models
opencode auth list
```

---

## 🦞 OH-MY-OPENCODE INSTALLATION

### Step 1: Install Plugin
```bash
cd ~/.config/opencode
npm init -y
npm install oh-my-opencode@latest
```

### Step 2: Configure oh-my-opencode.json
```json
{
  "$schema": "https://raw.githubusercontent.com/code-yeongyu/oh-my-opencode/master/assets/oh-my-opencode.schema.json",
  "google_auth": false,
  "sisyphus_agent": {
    "disabled": false,
    "planner_enabled": true,
    "replace_plan": true
  },
  "agents": {
    "Delqhi": { "model": "nvidia-nim/qwen-3.5-397b", "category": "general" },
    "sisyphus": { "model": "nvidia-nim/qwen-3.5-397b", "category": "general" },
    "prometheus": { "model": "nvidia-nim/qwen-3.5-397b", "category": "ultrabrain" },
    "metis": { "model": "nvidia-nim/qwen-3.5-397b", "category": "ultrabrain" },
    "librarian": { "model": "opencode/minimax-m2.5-free", "category": "general" },
    "explore": { "model": "opencode/minimax-m2.5-free", "category": "quick" }
  },
  "categories": {
    "ultrabrain": {
      "model": "nvidia-nim/qwen-3.5-397b",
      "temperature": 0.9,
      "thinking": { "type": "enabled", "budgetTokens": 16000 }
    },
    "quick": {
      "model": "opencode/minimax-m2.5-free",
      "temperature": 0.5
    },
    "general": {
      "model": "nvidia-nim/qwen-3.5-397b",
      "temperature": 0.5
    }
  },
  "git_master": {
    "commit_footer": true,
    "include_co_authored_by": true
  }
}
```

---

## 🔌 PROVIDER CONFIGURATION

### NVIDIA NIM Setup (PRIMARY)
**File:** `~/.config/opencode/opencode.json`

```json
{
  "$schema": "https://opencode.ai/config.json",
  "model": "nvidia-nim/qwen-3.5-397b",
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

**⚠️ CRITICAL:**
- Timeout MUST be 120000ms (Qwen 3.5 has 70-90s latency)
- Model ID MUST be `qwen/qwen3.5-397b-a17b` (NOT `qwen2.5`)
- Rate limit: 40 RPM (Free Tier)

### OpenCode ZEN (FREE FALLBACK)
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
        "minimax-m2.5-free": {
          "name": "MiniMax M2.5 Free (OpenCode ZEN - UNCENSORED)",
          "limit": { "context": 200000, "output": 128000 }
        },
        "kimi-k2.5-free": {
          "name": "Kimi K2.5 Free",
          "limit": { "context": 1048576, "output": 65536 }
        }
      }
    }
  }
}
```

---

## 🛠️ MCP SERVER CONFIGURATION

### Essential MCP Servers
**File:** `~/.config/opencode/opencode.json`

```json
{
  "mcp": {
    "context7": {
      "type": "remote",
      "url": "https://mcp.context7.com/mcp",
      "enabled": true
    },
    "gh_grep": {
      "type": "remote",
      "url": "https://mcp.grep.app",
      "enabled": true
    },
    "sin_social": {
      "type": "remote",
      "url": "http://localhost:8213",
      "enabled": true
    },
    "sin_deep_research": {
      "type": "remote",
      "url": "http://localhost:8214",
      "enabled": true
    },
    "sin_video_gen": {
      "type": "remote",
      "url": "http://localhost:8215",
      "enabled": true
    }
  }
}
```

### Local MCP Servers (Install Separately)
```bash
# Serena (Orchestration)
uvx serena start-mcp-server

# Tavily (Web Search)
npx @tavily/claude-mcp

# Canva (Design)
npx @canva/claude-mcp

# Skyvern (Browser Automation)
python -m skyvern.mcp.server
```

---

## 📁 FILE SYSTEM ORGANIZATION

### Primary Directories (CEO-LEVEL)
```
/Users/jeremy/
├── Desktop/                    # CLEAN - Auto-cleaned daily
├── Documents/                  # Important documents only
├── Downloads/                  # Temp, cleaned weekly
├── Bilder/
│   └── AI-Screenshots/        # All AI tool screenshots
│       ├── playwright/
│       ├── skyvern/
│       ├── steel/
│       ├── stagehand/
│       └── opencode/
├── dev/                       # ALL development
│   ├── sin-code/             # Main workspace
│   │   ├── OpenCode/
│   │   ├── Docker/
│   │   ├── Guides/
│   │   └── troubleshooting/
│   ├── SIN-Solver/           # AI automation
│   └── BIOMETRICS/           # Main project
└── .config/opencode/         # PRIMARY CONFIG
    ├── opencode.json         # Main config
    ├── AGENTS.md             # Mandates
    ├── oh-my-opencode.json   # Plugin config
    └── blueprint-vorlage.md  # Blueprint template
```

### Auto-Cleanup Scripts
**File:** `~/Library/LaunchAgents/com.sincode.desktop-cleanup.plist`
```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.sincode.desktop-cleanup</string>
    <key>ProgramArguments</key>
    <array>
        <string>/bin/bash</string>
        <string>-c</string>
        <string>find ~/Desktop -type f -mtime +7 -delete</string>
    </array>
    <key>StartCalendarInterval</key>
    <dict>
        <key>Hour</key>
        <integer>3</integer>
        <key>Minute</key>
        <integer>0</integer>
    </dict>
</dict>
</plist>
```

---

## 🔐 SECRETS MANAGEMENT

### Environment Variables
**File:** `~/.zshrc` (append)
```bash
# NVIDIA API Key
export NVIDIA_API_KEY="nvapi-xxx"

# OpenCode API Key (if needed)
export OPENCODE_API_KEY="xxx"

# Other API keys
export MISTRAL_API_KEY="xxx"
export ANTHROPIC_API_KEY="xxx"
```

### Global Secrets Registry
**File:** `~/dev/environments-jeremy.md` (APPEND-ONLY!)

```markdown
## NVIDIA NIM - 2026-02-16
**Service:** NVIDIA NIM API
**API Key:** nvapi-xxx (in ~/.zshrc)
**Endpoint:** https://integrate.api.nvidia.com/v1
**Rate Limit:** 40 RPM
**Status:** ACTIVE

## OpenCode ZEN - 2026-02-16
**Service:** OpenCode FREE API
**Endpoint:** https://api.opencode.ai/v1
**Status:** ACTIVE (no auth required)
```

**⚠️ CRITICAL RULES:**
- NEVER commit secrets to git
- ALWAYS store in environment variables
- NEVER delete from environments-jeremy.md (APPEND-ONLY)
- Mark deprecated secrets as "ROTATED" but don't delete

---

## ✅ VERIFICATION CHECKLIST

### OpenCode Verification
```bash
# 1. Version check
opencode --version

# 2. Models list
opencode models | grep nvidia

# 3. Auth status
opencode auth list

# 4. Test request
opencode --model nvidia-nim/qwen-3.5-397b "Hello"
```

### MCP Server Verification
```bash
# List active MCP servers
opencode mcp list

# Test context7
curl https://mcp.context7.com/mcp

# Test gh_grep
curl https://mcp.grep.app
```

### Directory Structure Verification
```bash
# Check directories exist
ls -la ~/.config/opencode/
ls -la ~/dev/sin-code/
ls -la ~/Bilder/AI-Screenshots/

# Check file permissions
chmod 600 ~/.config/opencode/opencode.json
chmod 600 ~/.zshrc
```

---

## 🚨 TROUBLESHOOTING

### Common Issues

**Issue:** OpenCode timeout errors  
**Solution:** Increase timeout to 120000ms in opencode.json

**Issue:** Model not found  
**Solution:** Verify model ID is `qwen/qwen3.5-397b-a17b` (not `qwen2.5`)

**Issue:** MCP server not connecting  
**Solution:** Check URL and ensure server is running

**Issue:** Rate limit exceeded (HTTP 429)  
**Solution:** Wait 60 seconds, use fallback model

**Issue:** Config not loading  
**Solution:** Run `opencode auth refresh`

### Recovery Commands
```bash
# Refresh auth
opencode auth refresh

# Clear cache
rm -rf ~/.config/opencode/node_modules/.cache

# Reinstall OpenCode (LAST RESORT!)
brew reinstall opencode

# Restore config from backup
cp ~/.config/opencode/archive/opencode.json.backup-* ~/.config/opencode/opencode.json
```

---

## 📚 NEXT STEPS

After completing this checklist:

1. **Read AGENTS.md** - Understand all 33 mandates
2. **Read blueprint-vorlage.md** - Understand 22 pillars
3. **Setup BIOMETRICS project** - Follow project-specific docs
4. **Configure Docker** - Setup 26-room infrastructure
5. **Setup MCP wrappers** - For Docker container integration

---

**Document Statistics:**
- Total Steps: 50+
- Verification Commands: 15+
- Files Created: 10+
- Estimated Time: 2-3 hours

**Status:** ✅ PRODUCTION-READY FOR NEW MAC SETUP
