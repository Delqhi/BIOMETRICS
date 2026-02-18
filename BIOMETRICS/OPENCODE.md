# OPENCODE.md — Complete Configuration Guide

**Status:** ACTIVE
**Version:** 1.0
**Stand:** 2026-02-17
**Purpose:** Vollständige Konfiguration für OpenCode + OpenClaw mit NVIDIA NIM + Google Antigravity

---

## 🧠 NLM CLI COMMANDS

```bash
# Create notebook
nlm notebook create "Title"

# List sources
nlm source list <notebook-id>

# Delete old source (BEFORE adding new!)
nlm source delete <source-id> -y

# Add new source
nlm source add <notebook-id> --file "file.md" --wait
```

**⚠️ DUPLICATE PREVENTION:** ALWAYS run `nlm source list` before `nlm source add`!

---

## 🔄 DEQLHI-LOOP (INFINITE WORK MODE)

- After each completed task → Add 5 new tasks immediately
- Never "done" - only "next task"
- Always document → Every change in files
- Git commit + push after EVERY change
- Parallel execution ALWAYS (run_in_background=true)

### Loop Mechanism:
1. Task N Complete
2. Git Commit + Push
3. Update Docs
4. Add 5 New Tasks
5. Next Task N+1
6. Repeat infinitely

---

## 🚀 1. Quick Start

### Voraussetzungen
- Node.js >= 20
- pnpm (NICHT npm!)
- Docker (für MCP Server)
- NVIDIA API Key (für Qwen 3.5 397B)

### Installation
```bash
# OpenCode installieren
npm install -g opencode

# OpenClaw installieren (pnpm!)
pnpm install -g @anthropic/openclaw

# NVIDIA API Key setzen
export NVIDIA_API_KEY="nvapi-xxx"
```

### Agent Model Configuration

**Standard Agent Models:**

| Agent | Model | Provider | Use Case |
|-------|-------|----------|----------|
| **build** | `nvidia-nim/qwen-3.5-397b` | NVIDIA NIM | Code implementation |
| **plan** | `nvidia-nim/qwen-3.5-397b` | NVIDIA NIM | Planning & architecture |
| **quick** | `opencode/minimax-m2.5-free` | OpenCode | Fast tasks |
| **deep** | `nvidia-nim/qwen-3.5-397b` | NVIDIA NIM | Complex reasoning |
| **ultrabrain** | `nvidia-nim/qwen-3.5-397b` | NVIDIA NIM | High-level thinking |
| **artistry** | `nvidia-nim/qwen-3.5-397b` | NVIDIA NIM | Creative tasks |
| **visual-engineering** | `nvidia-nim/qwen-3.5-397b` | NVIDIA NIM | Frontend/UI/UX |
| **oracle** | `nvidia-nim/qwen-3.5-397b` | NVIDIA NIM | Architecture review |
| **metis** | `nvidia-nim/qwen-3.5-397b` | NVIDIA NIM | Pre-planning analysis |
| **momus** | `nvidia-nim/qwen-3.5-397b` | NVIDIA NIM | Quality assurance |
| **writing** | `nvidia-nim/qwen-3.5-397b` | NVIDIA NIM | Documentation |
| **explore** | `opencode/minimax-m2.5-free` | OpenCode | Codebase search |
| **librarian** | `opencode/minimax-m2.5-free` | OpenCode | External research |

**Configuration Files:**

1. **`~/.config/opencode/opencode.json`** - Main OpenCode config
2. **`~/.config/opencode/oh-my-opencode.json`** - Agent-specific models

**Example Configuration:**

```json
{
  "agent": {
    "build": { "model": "nvidia-nim/qwen-3.5-397b" },
    "plan": { "model": "nvidia-nim/qwen-3.5-397b" },
    "quick": { "model": "opencode/minimax-m2.5-free" },
    "deep": { "model": "nvidia-nim/qwen-3.5-397b" },
    "ultrabrain": { "model": "nvidia-nim/qwen-3.5-397b" },
    "explore": { "model": "opencode/minimax-m2.5-free" },
    "librarian": { "model": "opencode/minimax-m2.5-free" }
  },
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
          "name": "Qwen 3.5 397B (NVIDIA NIM)",
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

---

## 🔐 2. Google Antigravity OAuth Setup

### Schritt 1: OAuth Flow starten
```bash
opencode auth login
```

### Schritt 2: Browser Authentication
1. Browser öffnet sich automatisch
2. Mit **privater Gmail** anmelden (NICHT Google Workspace!)
3. Berechtigungen erteilen
4. Token wird automatisch gespeichert

### Schritt 3: Verifizieren
```bash
# Auth File prüfen (Permissions: 600)
ls -la ~/.config/opencode/antigravity-accounts.json

# Account anzeigen
cat ~/.config/opencode/antigravity-accounts.json | jq '.accounts[0].email'
```

### Schritt 4: Testen
```bash
# Gemini 3 Flash testen
opencode run "Hello" --model google/antigravity-gemini-3-flash

# Mit Thinking Level
opencode run "Solve this problem" --model google/antigravity-gemini-3-pro:high

# Multimodal (Bild)
opencode run "Describe this image" --model google/antigravity-gemini-3-flash --image ./test.png
```

### Verfügbare Modelle

**Gemini 3 Flash:**
- `google/antigravity-gemini-3-flash:minimal` - Fast, minimal thinking
- `google/antigravity-gemini-3-flash:high` - Slower, deeper reasoning
- Context: 1M tokens
- Output: 64K tokens
- Multimodal: Text, Image, PDF

**Gemini 3 Pro:**
- `google/antigravity-gemini-3-pro:low` - Low thinking budget
- `google/antigravity-gemini-3-pro:high` - High thinking budget (32K tokens)
- Context: 2M tokens
- Output: 64K tokens

**Claude Sonnet 4.5:**
- `google/antigravity-claude-sonnet-4-5-thinking:low` - 8K thinking budget
- `google/antigravity-claude-sonnet-4-5-thinking:max` - 32K thinking budget
- Context: 200K tokens
- Output: 64K tokens

### Rate Limits
- 100 RPM (Requests Per Minute)
- 50,000 RPD (Requests Per Day)
- Auto-refresh enabled

---

## 🎮 3. NVIDIA NIM Configuration

### Schritt 1: API Key setzen
```bash
export NVIDIA_API_KEY="nvapi-xxx"
```

### Schritt 2: OpenCode Config bearbeiten
```bash
nano ~/.config/opencode/opencode.json
```

### Schritt 3: Provider hinzufügen
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

### ⚠️ WICHTIGE REGELN

**KRITISCH:**
- ✅ Model: `qwen/qwen3.5-397b-a17b` (NICHT qwen2.5!)
- ✅ Timeout: `120000` (120 Sekunden!)
- ❌ KEIN `reasoning: true` (verursacht Fehler!)
- ❌ KEIN `stream: true` (nicht unterstützt!)

### Schritt 4: Verifizieren
```bash
# Models auflisten
opencode models | grep nvidia

# Testen
opencode run "Hello" --model nvidia-nim/qwen-3.5-397b

# Code-Beispiel
opencode run "Write a TypeScript function" --model nvidia-nim/qwen-3.5-397b
```

### Rate Limits
- 40 RPM (FREE Tier)
- Bei HTTP 429: 60 Sekunden warten + Fallback

---

## 🦎 4. OpenClaw Configuration

### Schritt 1: Installation
```bash
# MIT PNPM (NICHT NPM!)
pnpm install -g @anthropic/openclaw
```

### Schritt 2: Config bearbeiten
```bash
nano ~/.openclaw/openclaw.json
```

### Schritt 3: NVIDIA Provider hinzufügen
```json
{
  "env": {
    "NVIDIA_API_KEY": "${NVIDIA_API_KEY}"
  },
  "models": {
    "providers": {
      "nvidia": {
        "baseUrl": "https://integrate.api.nvidia.com/v1",
        "api": "openai-completions",
        "models": ["qwen/qwen3.5-397b-a17b"]
      }
    },
    "defaults": {
      "model": {
        "primary": "nvidia/qwen/qwen3.5-397b-a17b",
        "fallbacks": [
          "nvidia/meta/llama-3.3-70b-instruct",
          "nvidia/mistralai/mistral-large-3-675b-instruct-2512"
        ]
      }
    }
  },
  "agents": {
    "defaults": {
      "model": {
        "primary": "nvidia/qwen/qwen3.5-397b-a17b"
      }
    }
  }
}
```

### ⚠️ PFLICHT-FELDER

**KRITISCH:**
- ✅ `api: "openai-completions"` (PFLICHT!)
- ✅ Model: `qwen/qwen3.5-397b-a17b`
- ✅ `NVIDIA_API_KEY` in `env` section

### Schritt 4: Verifizieren
```bash
# Models auflisten
openclaw models | grep nvidia

# Testen
openclaw run "Hello" --model nvidia/qwen/qwen3.5-397b-a17b

# Health Check
openclaw doctor --fix
```

---

## 🔧 5. MCP Server Setup

### Serena MCP (Orchestrierung)
```bash
# Installieren
uv tool install serena

# Starten
serena start-mcp-server

# In opencode.json hinzufügen:
{
  "mcp": {
    "serena": {
      "type": "local",
      "command": ["uvx", "serena", "start-mcp-server"],
      "enabled": true
    }
  }
}
```

### Context7 MCP (Dokumentation)
```bash
# Installieren
npx @anthropics/context7-mcp

# In opencode.json:
{
  "mcp": {
    "context7": {
      "type": "local",
      "command": ["npx", "@anthropics/context7-mcp"],
      "enabled": true
    }
  }
}
```

### Chrome DevTools MCP (Browser Verification)
```bash
# Installieren
npx @anthropics/chrome-devtools-mcp

# In opencode.json:
{
  "mcp": {
    "chrome-devtools": {
      "type": "local",
      "command": ["npx", "@anthropics/chrome-devtools-mcp"],
      "enabled": true
    }
  }
}
```

---

## 🧪 6. Complete Verification

### OpenCode Tests
```bash
# 1. Models auflisten
opencode models

# Erwartet: NVIDIA + Google angezeigt

# 2. NVIDIA testen
opencode run "Write a function" --model nvidia-nim/qwen-3.5-397b

# Erwartet: Code-Generierung funktioniert

# 3. Antigravity testen
opencode run "Explain quantum computing" --model google/antigravity-gemini-3-flash

# Erwartet: Erklärung generiert

# 4. Multimodal testen
opencode run "Describe this image" --model google/antigravity-gemini-3-flash --image ./test.png

# Erwartet: Bildbeschreibung

# 5. Auth Status
opencode auth status

# Erwartet: Auth accounts angezeigt
```

### OpenClaw Tests
```bash
# 1. Models auflisten
openclaw models

# Erwartet: NVIDIA angezeigt

# 2. NVIDIA testen
openclaw run "Write a function" --model nvidia/qwen/qwen3.5-397b-a17b

# Erwartet: Code-Generierung funktioniert

# 3. Health Check
openclaw doctor --fix

# Erwartet: Keine Fehler
```

### Success Criteria
- [ ] `opencode models` zeigt NVIDIA + Google
- [ ] `openclaw models` zeigt NVIDIA
- [ ] NVIDIA Test-Command funktioniert
- [ ] Antigravity Test-Command funktioniert
- [ ] Multimodal funktioniert
- [ ] Health Check ohne Fehler

---

## 🚨 7. Troubleshooting

### Problem: "Bad credentials"
**Lösung:**
```bash
# OAuth neu starten
opencode auth logout
opencode auth login
```

### Problem: "Model not found"
**Lösung:**
- Prüfe Model-ID: `qwen/qwen3.5-397b-a17b` (NICHT `qwen2.5`)
- Config neu laden: `opencode --version`

### Problem: "Timeout"
**Lösung:**
- Timeout auf `120000` setzen (120 Sekunden)
- Qwen 3.5 397B hat hohe Latenz (70-90s)

### Problem: "HTTP 429 Too Many Requests"
**Lösung:**
- NVIDIA: 60 Sekunden warten
- Fallback-Model verwenden
- Rate Limit respektieren (40 RPM)

---

## 📚 8. Referenzen

- **Unified Skill Architecture:** `WORKFLOW.md`
- **Meta-Builder Protocol:** `WORKFLOW.md`
- **Supabase Integration:** `SUPABASE.md`
- **n8n Integration:** `N8N.md`
- **OpenClaw Skills:** `OPENCLAW.md`
- **Agent Behavior:** `AGENTS-GLOBAL.md`

---

## 🔐 9. Security Best Practices

### Secrets Management
- ❌ NIEMALS API Keys in Git committen
- ✅ Environment Variables verwenden
- ✅ `.gitignore` für `.env` Dateien

### File Permissions
```bash
# Auth File schützen
chmod 600 ~/.config/opencode/antigravity-accounts.json
```

### Rate Limiting
- NVIDIA: 40 RPM (FREE Tier)
- Antigravity: 100 RPM / 50K RPD
- Backoff bei 429 Errors

---

**Version:** 1.0  
**Stand:** 2026-02-17  
**Status:** PRODUCTION READY ✅  
**Next:** Siehe `WORKFLOW.md` für Self-Building AI Architecture

---

## 🔌 10. Plugins Installation & Konfiguration

### Oh My OpenCode Plugin

**Purpose:** Enhanced OpenCode experience with additional features, UI improvements, and productivity tools.

**Installation:**
```bash
# Plugin installieren
opencode plugins install oh-my-opencode

# Oder manuell:
mkdir -p ~/.config/opencode/plugins
cd ~/.config/opencode/plugins
git clone https://github.com/opencode/oh-my-opencode.git
```

**Konfiguration in `opencode.json`:**
```json
{
  "plugin": [
    "oh-my-opencode@latest"
  ]
}
```

**Features:**
- ✅ Enhanced UI/UX
- ✅ Auto-completion improvements
- ✅ Better error messages
- ✅ Productivity shortcuts
- ✅ Theme support
- ✅ Command palette

**Verifizieren:**
```bash
opencode plugins list
# Sollte "oh-my-opencode" anzeigen
```

---

### Qwen OAuth Plugin (opencode-qwencode-auth)

**Purpose:** OAuth authentication for Qwen models via official Qwen Code platform.

**Installation:**
```bash
# Plugin installieren
opencode plugins install opencode-qwencode-auth

# Oder via npm:
npm install -g opencode-qwencode-auth
```

**Konfiguration in `opencode.json`:**
```json
{
  "provider": {
    "qwencode": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "Qwen Code (OAuth)",
      "options": {
        "baseURL": "https://chat.qwen.ai/api/v1"
      },
      "models": {
        "qwen-max": {
          "id": "qwen-max",
          "name": "Qwen Max",
          "limit": {
            "context": 32000,
            "output": 8000
          }
        },
        "qwen-plus": {
          "id": "qwen-plus",
          "name": "Qwen Plus",
          "limit": {
            "context": 131000,
            "output": 8000
          }
        },
        "qwen-turbo": {
          "id": "qwen-turbo",
          "name": "Qwen Turbo",
          "limit": {
            "context": 131000,
            "output": 8000
          }
        }
      }
    }
  }
}
```

**OAuth Flow:**
```bash
# OAuth starten
opencode auth login https://chat.qwen.ai

# Browser öffnet sich
# Mit Qwen Account anmelden
# Berechtigungen erteilen
```

**Verifizieren:**
```bash
# Models auflisten
opencode models | grep qwen

# Testen
opencode run "Hello" --model qwencode/qwen-max
```

**Rate Limits:**
- Qwen Max: 10 RPM (FREE)
- Qwen Plus: 20 RPM (FREE)
- Qwen Turbo: 50 RPM (FREE)

---

### Antigravity Plugin (Google OAuth)

**Bereits installiert via:**
```bash
opencode auth login
```

**Plugin ist in `opencode.json`:**
```json
{
  "plugin": [
    "opencode-antigravity-auth@latest",
    "oh-my-opencode"
  ]
}
```

**Auth File Location:**
```bash
~/.config/opencode/antigravity-accounts.json
```

---

## ⚙️ 11. Complete opencode.json Example

**Vollständige Konfiguration mit ALLEN Plugins:**

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
    },
    "google": {
      "npm": "@ai-sdk/google",
      "models": {
        "antigravity-gemini-3-flash": {
          "id": "gemini-3-flash-preview",
          "name": "Gemini 3 Flash (Antigravity)",
          "limit": {
            "context": 1048576,
            "output": 65536
          },
          "modalities": {
            "input": ["text", "image", "pdf"],
            "output": ["text"]
          },
          "variants": {
            "minimal": {
              "thinkingLevel": "minimal"
            },
            "high": {
              "thinkingLevel": "high"
            }
          }
        },
        "antigravity-gemini-3-pro": {
          "id": "gemini-3-pro-preview",
          "name": "Gemini 3 Pro (Antigravity)",
          "limit": {
            "context": 2097152,
            "output": 65536
          },
          "variants": {
            "low": {
              "thinkingLevel": "low"
            },
            "high": {
              "thinkingLevel": "high"
            }
          }
        }
      }
    },
    "qwencode": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "Qwen Code (OAuth)",
      "options": {
        "baseURL": "https://chat.qwen.ai/api/v1"
      },
      "models": {
        "qwen-max": {
          "id": "qwen-max",
          "limit": {
            "context": 32000,
            "output": 8000
          }
        }
      }
    }
  },
  "mcp": {
    "serena": {
      "type": "local",
      "command": ["uvx", "serena", "start-mcp-server"],
      "enabled": true
    },
    "context7": {
      "type": "local",
      "command": ["npx", "@anthropics/context7-mcp"],
      "enabled": true
    }
  },
  "plugin": [
    "opencode-antigravity-auth@latest",
    "oh-my-opencode@latest",
    "opencode-qwencode-auth@latest"
  ]
}
```

---

## 📋 12. Best Practices 2026

### Installation
- ✅ OpenCode: `npm install -g opencode`
- ✅ OpenClaw: `pnpm install -g @anthropic/openclaw` (NICHT npm!)
- ✅ Plugins: Immer `@latest` Tag verwenden

### Configuration
- ✅ Timeout: `120000ms` für High-Latency Modelle (Qwen 3.5 397B)
- ✅ Model IDs: Immer exakte IDs verwenden (`qwen/qwen3.5-397b-a17b`)
- ✅ Provider Schema: Offizielles `@ai-sdk/openai-compatible` Schema

### Security
- ✅ API Keys: Nur in Environment Variables
- ✅ Auth Files: Permissions `600` (owner read/write only)
- ✅ `.gitignore`: `.env`, `*.json` mit Secrets

### Performance
- ✅ Connection Pooling: MCP Server persistent laufen lassen
- ✅ Caching: Response-Caching für wiederholte Queries
- ✅ Fallback Chain: Bei 429 Errors automatisch fallbacken

### Error Handling
- ✅ Leere Catch-Blöcke: VERBOTEN
- ✅ Descriptive Errors: Mit Context und Stacktrace
- ✅ Retry Logic: Exponential backoff bei transient errors

### Testing
- ✅ Vor Production: Alle Models testen
- ✅ Rate Limits: Limits respektieren und monitoren
- ✅ Health Checks: Regelmäßig `opencode doctor` / `openclaw doctor`

---

## 🎯 13. Quick Reference Commands

```bash
# Installation
npm install -g opencode
pnpm install -g @anthropic/openclaw

# Plugins
opencode plugins install oh-my-opencode
opencode plugins install opencode-qwencode-auth
opencode plugins list

# Auth
opencode auth login
opencode auth logout
opencode auth status

# Models
opencode models
openclaw models

# Testing
opencode run "Hello" --model nvidia-nim/qwen-3.5-397b
openclaw run "Hello" --model nvidia/qwen/qwen3.5-397b-a17b

# Health
opencode doctor
openclaw doctor --fix

# MCP
serena start-mcp-server
uvx serena start-mcp-server
```

---

**Updated:** 2026-02-17  
**Status:** PRODUCTION READY ✅  
**All Plugins:** Oh My OpenCode + Qwen OAuth + Antigravity

---

## 🤖 14. Agenten-Konfiguration (OH-MY-OPENCODE)

### Agenten-Modelle Konfiguration

**Alle Agenten sind in `~/.config/opencode/opencode.json` konfiguriert:**

```json
{
  "agent": {
    "sisyphus": {
      "model": {
        "provider": "moonshot-ai",
        "id": "moonshotai/kimi-k2.5"
      }
    },
    "sisyphus-junior": {
      "model": {
        "provider": "kimi-for-coding",
        "id": "kimi-for-coding/k2p5"
      }
    },
    "prometheus": {
      "model": {
        "provider": "kimi-for-coding",
        "id": "kimi-for-coding/k2p5"
      }
    },
    "metis": {
      "model": {
        "provider": "kimi-for-coding",
        "id": "kimi-for-coding/k2p5"
      }
    },
    "momus": {
      "model": {
        "provider": "kimi-for-coding",
        "id": "kimi-for-coding/k2p5"
      }
    },
    "oracle": {
      "model": {
        "provider": "kimi-for-coding",
        "id": "kimi-for-coding/k2p5"
      }
    },
    "frontend-ui-ux-engineer": {
      "model": {
        "provider": "kimi-for-coding",
        "id": "kimi-for-coding/k2p5"
      }
    },
    "document-writer": {
      "model": {
        "provider": "kimi-for-coding",
        "id": "kimi-for-coding/k2p5"
      }
    },
    "multimodal-looker": {
      "model": {
        "provider": "kimi-for-coding",
        "id": "kimi-for-coding/k2p5"
      }
    },
    "atlas": {
      "model": {
        "provider": "kimi-for-coding",
        "id": "kimi-for-coding/k2p5"
      }
    },
    "librarian": {
      "model": {
        "provider": "opencode-zen",
        "id": "zen/big-pickle"
      }
    },
    "explore": {
      "model": {
        "provider": "opencode-zen",
        "id": "zen/big-pickle"
      }
    }
  }
}
```

### Agenten-Übersicht

| Agent | Modell | Provider | Kosten | Purpose |
|-------|--------|----------|--------|---------|
| **sisyphus** | kimi-k2.5 | Moonshot AI | 💰 | Haupt-Agent (Orchestrator) |
| **sisyphus-junior** | k2p5 | Kimi For Coding | 💰 | Junior Developer |
| **prometheus** | k2p5 | Kimi For Coding | 💰 | Planning & Strategy |
| **metis** | k2p5 | Kimi For Coding | 💰 | Pre-Planning Consultant |
| **momus** | k2p5 | Kimi For Coding | 💰 | Code Reviewer |
| **oracle** | k2p5 | Kimi For Coding | 💰 | Architecture Consultant |
| **frontend-ui-ux-engineer** | k2p5 | Kimi For Coding | 💰 | Frontend Specialist |
| **document-writer** | k2p5 | Kimi For Coding | 💰 | Technical Writer |
| **multimodal-looker** | k2p5 | Kimi For Coding | 💰 | Image/PDF Analysis |
| **atlas** | k2p5 | Kimi For Coding | 💰 | Heavy Lifting |
| **librarian** | zen/big-pickle | OpenCode ZEN | 🆓 FREE | Research & Documentation |
| **explore** | zen/big-pickle | OpenCode ZEN | 🆓 FREE | Codebase Exploration |

### Warum diese Verteilung?

1. **Sisyphus (moonshotai/kimi-k2.5)**
   - Premium-Modell für Haupt-Agent
   - Beste Code-Qualität
   - Höchste Zuverlässigkeit

2. **Andere Coding-Agenten (kimi-for-coding/k2p5)**
   - Gutes Modell, kosteneffizient
   - Spezialisiert auf Coding-Tasks
   - Balance zwischen Qualität und Kosten

3. **Recherche-Agenten (zen/big-pickle)**
   - 100% KOSTENLOS
   - Perfekt für Suche und Recherche
   - Uncensored, keine Limits

### Provider Setup

**Alle Provider wurden über `/connect` hinzugefügt:**

```bash
# Moonshot AI
opencode auth add moonshot-ai

# Kimi For Coding
opencode auth add kimi-for-coding

# OpenCode ZEN (FREE)
opencode auth add opencode-zen
```

**Verifizierung:**
```bash
opencode auth list
opencode models
```

---

## 📋 15. Vollständige Agenten-Beschreibungen

### Sisyphus (Haupt-Agent)
**Modell:** `moonshotai/kimi-k2.5`  
**Rolle:** Orchestrator, Haupt-Developer  
**Aufgaben:**
- Gesamt-Orchestrierung
- Code-Implementierung
- Quality Assurance
- Final Reviews

**Stärken:**
- Höchste Code-Qualität
- Zuverlässig bei komplexen Tasks
- Enterprise-grade Output

---

### Sisyphus-Junior
**Modell:** `kimi-for-coding/k2p5`  
**Rolle:** Junior Developer  
**Aufgaben:**
- Einfache Code-Änderungen
- Bugfixes
- Testing
- Dokumentation

**Stärken:**
- Schnell bei einfachen Tasks
- Kosteneffizient
- Lernt von Sisyphus

---

### Prometheus
**Modell:** `kimi-for-coding/k2p5`  
**Rolle:** Planning & Strategy  
**Aufgaben:**
- Task-Planning
- Strategie-Entwicklung
- Roadmap-Erstellung
- Resource-Allocation

**Stärken:**
- Strategisches Denken
- Langfristige Planung
- Risk-Assessment

---

### Metis
**Modell:** `kimi-for-coding/k2p5`  
**Rolle:** Pre-Planning Consultant  
**Aufgaben:**
- Anforderungs-Analyse
- Ambiguity Detection
- Scope-Definition
- Risk-Identification

**Stärken:**
- Findet versteckte Anforderungen
- Erkennt Probleme vor Implementation
- Spart Zeit durch frühe Klärung

---

### Momus
**Modell:** `kimi-for-coding/k2p5`  
**Rolle:** Code Reviewer  
**Aufgaben:**
- Code Reviews
- Quality Checks
- Best Practices Enforcement
- Security Audits

**Stärken:**
- Findet Bugs vor Production
- Erzwingt Coding-Standards
- Security-first Mindset

---

### Oracle
**Modell:** `kimi-for-coding/k2p5`  
**Rolle:** Architecture Consultant  
**Aufgaben:**
- Architecture Reviews
- Complex Problem Solving
- Trade-off Analysis
- Pattern Recognition

**Stärken:**
- Deep Technical Knowledge
- Multi-system Thinking
- Enterprise Architecture

---

### Frontend-UI-UX-Engineer
**Modell:** `kimi-for-coding/k2p5`  
**Rolle:** Frontend Specialist  
**Aufgaben:**
- UI Implementation
- UX Optimization
- Responsive Design
- Accessibility

**Stärken:**
- Modern Frontend Stack
- Design Sense
- User-centric Approach

---

### Document-Writer
**Modell:** `kimi-for-coding/k2p5`  
**Rolle:** Technical Writer  
**Aufgaben:**
- API Documentation
- User Guides
- README Files
- CHANGELOG

**Stärken:**
- Klare Sprache
- Vollständige Doku
- Developer Experience

---

### Multimodal-Looker
**Modell:** `kimi-for-coding/k2p5`  
**Rolle:** Image/PDF Analysis  
**Aufgaben:**
- Image Recognition
- PDF Parsing
- Screenshot Analysis
- Visual QA

**Stärken:**
- Multimodal Understanding
- Detail-oriented
- Visual Problem Solving

---

### Atlas
**Modell:** `kimi-for-coding/k2p5`  
**Rolle:** Heavy Lifting  
**Aufgaben:**
- Large Refactorings
- Complex Migrations
- Bulk Operations
- Data Processing

**Stärken:**
- Ausdauer bei großen Tasks
- Systematisches Vorgehen
- Keine Müdigkeit

---

### Librarian
**Modell:** `zen/big-pickle` (FREE)  
**Rolle:** Research & Documentation  
**Aufgaben:**
- Web Research
- Documentation Search
- Knowledge Management
- Fact Checking

**Stärken:**
- 100% KOSTENLOS
- Uncensored Research
- Deep Dives

---

### Explore
**Modell:** `zen/big-pickle` (FREE)  
**Rolle:** Codebase Exploration  
**Aufgaben:**
- Code Discovery
- Pattern Finding
- Structure Analysis
- Dependency Mapping

**Stärken:**
- 100% KOSTENLOS
- Fast Codebase Understanding
- Pattern Recognition

---

## 🔧 16. Agent Usage Examples

### Sisyphus für Haupt-Development
```bash
opencode run "Implement user authentication with JWT" \
  --agent sisyphus \
  --model moonshotai/kimi-k2.5
```

### Librarian für Research (FREE)
```bash
opencode run "Find best practices for React Server Components 2026" \
  --agent librarian \
  --model zen/big-pickle
```

### Explore für Codebase Analysis (FREE)
```bash
opencode run "Find all authentication middleware in src/" \
  --agent explore \
  --model zen/big-pickle
```

### Oracle für Architecture Review
```bash
opencode run "Review this microservices architecture for scalability issues" \
  --agent oracle \
  --model kimi-for-coding/k2p5
```

### Momus für Code Review
```bash
opencode run "Review this pull request for security vulnerabilities" \
  --agent momus \
  --model kimi-for-coding/k2p5
```

---

## ⚙️ 17. Agent Selection Guide

### Wann welchen Agenten verwenden?

| Task Typ | Agent | Modell | Kosten |
|----------|-------|--------|--------|
| **Haupt-Development** | sisyphus | kimi-k2.5 | 💰💰 |
| **Einfache Fixes** | sisyphus-junior | k2p5 | 💰 |
| **Planning** | prometheus | k2p5 | 💰 |
| **Requirements** | metis | k2p5 | 💰 |
| **Code Review** | momus | k2p5 | 💰 |
| **Architecture** | oracle | k2p5 | 💰💰 |
| **Frontend** | frontend-ui-ux | k2p5 | 💰 |
| **Dokumentation** | document-writer | k2p5 | 💰 |
| **Image Analysis** | multimodal-looker | k2p5 | 💰 |
| **Große Tasks** | atlas | k2p5 | 💰💰 |
| **Research** | librarian | zen/big-pickle | 🆓 |
| **Code Exploration** | explore | zen/big-pickle | 🆓 |

### Kosten-Optimierung

**Strategie:**
1. **Research/Exploration:** Immer FREE Agents (librarian, explore)
2. **Einfache Tasks:** Junior Agents (sisyphus-junior)
3. **Komplexe Tasks:** Premium Agents (sisyphus, oracle)
4. **Reviews:** Spezialisten (momus, metis)

**Beispiel-Workflow:**
```bash
# 1. FREE: Research
opencode run "Research best auth patterns" --agent librarian

# 2. FREE: Explore Codebase
opencode run "Find existing auth code" --agent explore

# 3. 💰: Implementation
opencode run "Implement JWT auth" --agent sisyphus

# 4. 💰: Review
opencode run "Review auth implementation" --agent momus
```

---

## 🌐 18. Provider Configuration Deep Dive

### 18.1 NVIDIA NIM Provider (Ausführlich)

Der NVIDIA NIM Provider ermöglicht Zugang zu verschiedenen KI-Modellen über die NVIDIA Inference Microservices.

#### Konfiguration

```json
{
  "provider": {
    "nvidia-nim": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "NVIDIA NIM",
      "options": {
        "baseURL": "https://integrate.api.nvidia.com/v1",
        "timeout": 120000,
        "headers": {
          "Authorization": "Bearer ${NVIDIA_API_KEY}"
        }
      },
      "models": {
        "qwen-3.5-397b": {
          "id": "qwen/qwen3.5-397b-a17b",
          "name": "Qwen 3.5 397B",
          "limit": {
            "context": 262144,
            "output": 32768
          },
          "supports": {
            "vision": true,
            "function-calling": true,
            "streaming": false
          }
        },
        "qwen2.5-coder-32b": {
          "id": "qwen/qwen2.5-coder-32b-instruct",
          "name": "Qwen 2.5 Coder 32B",
          "limit": {
            "context": 131072,
            "output": 8192
          }
        },
        "llama-3.3-70b": {
          "id": "meta/llama-3.3-70b-instruct",
          "name": "Llama 3.3 70B",
          "limit": {
            "context": 131072,
            "output": 8192
          }
        }
      }
    }
  }
}
```

#### Umgebungsvariablen

```bash
# NVIDIA API Key
export NVIDIA_API_KEY="nvapi-xxxxxxxxxxxxxxxxxxxx"

# Optional: Custom Timeout
export NVIDIA_TIMEOUT=120000
```

#### API-Endpunkte

| Endpoint | Methode | Beschreibung |
|----------|---------|--------------|
| /v1/models | GET | Liste aller verfügbaren Modelle |
| /v1/chat/completions | POST | Chat-Completion erstellen |
| /v1/completions | POST | Text-Completion erstellen |

#### Request-Beispiel

```bash
curl -X POST https://integrate.api.nvidia.com/v1/chat/completions \
  -H "Authorization: Bearer $NVIDIA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "qwen/qwen3.5-397b-a17b",
    "messages": [
      {
        "role": "user",
        "content": "Explain quantum computing"
      }
    ],
    "max_tokens": 1024,
    "temperature": 0.7
  }'
```

#### Fehlerbehandlung

| HTTP Status | Bedeutung | Lösung |
|-------------|-----------|--------|
| 200 | Erfolgreich | Response verarbeiten |
| 400 | Bad Request | Request-Format prüfen |
| 401 | Unauthorized | API-Key prüfen |
| 429 | Rate Limited | 60s warten + retry |
| 500 | Server Error | Retry mit backoff |

---

### 18.2 Google Provider (Antigravity) (Ausführlich)

Der Google Provider nutzt OAuth für den Zugang zu Gemini-Modellen über Antigravity.

#### Konfiguration

```json
{
  "provider": {
    "google": {
      "npm": "@ai-sdk/google",
      "models": {
        "gemini-3-flash": {
          "id": "gemini-3-flash-preview",
          "name": "Gemini 3 Flash",
          "limit": {
            "context": 1048576,
            "output": 65536
          },
          "reasoning": true,
          "modalities": {
            "input": ["text", "image", "pdf", "video"],
            "output": ["text"]
          },
          "variants": {
            "minimal": {
              "thinkingLevel": "minimal",
              "description": "Fast responses"
            },
            "high": {
              "thinkingLevel": "high",
              "description": "Deep reasoning"
            }
          }
        },
        "gemini-3-pro": {
          "id": "gemini-3-pro-preview",
          "name": "Gemini 3 Pro",
          "limit": {
            "context": 2097152,
            "output": 65536
          },
          "reasoning": true,
          "variants": {
            "low": {
              "thinkingLevel": "low",
              "thinkingBudget": 8192
            },
            "high": {
              "thinkingLevel": "high",
              "thinkingBudget": 32768
            }
          }
        },
        "antigravity-gemini-3-flash": {
          "id": "gemini-3-flash-preview",
          "name": "Gemini 3 Flash (Antigravity)",
          "limit": {
            "context": 1048576,
            "output": 65536
          },
          "modalities": {
            "input": ["text", "image", "pdf"],
            "output": ["text"]
          },
          "variants": {
            "minimal": {
              "thinkingLevel": "minimal"
            },
            "high": {
              "thinkingLevel": "high"
            }
          }
        },
        "antigravity-claude-sonnet-4-5-thinking": {
          "id": "claude-sonnet-4-5-20250501",
          "name": "Claude Sonnet 4.5 Thinking",
          "limit": {
            "context": 200000,
            "output": 64000
          },
          "variants": {
            "low": {
              "thinkingConfig": {
                "thinkingBudget": 8192
              }
            },
            "max": {
              "thinkingConfig": {
                "thinkingBudget": 32768
              }
            }
          }
        }
      }
    }
  }
}
```

#### OAuth-Flow

```bash
# Schritt 1: Auth starten
opencode auth login

# Schritt 2: Google auswählen
# > Google
#   > OAuth with Google (Antigravity)

# Schritt 3: Im Browser anmelden
# Automatische Weiterleitung nach erfolgreicher Anmeldung

# Schritt 4: Token speichern
# Token wird automatisch in ~/.config/opencode/antigravity-accounts.json gespeichert
```

#### Multi-Account Support

```bash
# Mehrere Accounts hinzufügen
opencode auth login

# Option: Account hinzufügen auswählen
# Weitere Google-Accounts anmelden

# Automatisches Load Balancing
# Bei Rate Limit wird automatisch zwischen Accounts gewechselt
```

#### Token-Refresh

```bash
# Manueller Refresh
opencode auth refresh

# Automatischer Refresh
# Wird bei Bedarf automatisch durchgeführt
```

---

### 18.3 OpenCode ZEN Provider (Ausführlich)

Der OpenCode ZEN Provider bietet kostenlose Modelle.

#### Konfiguration

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
        "zen/big-pickle": {
          "name": "Big Pickle (UNCENSORED)",
          "limit": {
            "context": 200000,
            "output": 128000
          },
          "description": "Uncensored model for unrestricted generation"
        },
        "zen/uncensored": {
          "name": "Uncensored",
          "limit": {
            "context": 200000,
            "output": 128000
          }
        },
        "zen/advanced": {
          "name": "Advanced",
          "limit": {
            "context": 200000,
            "output": 128000
          }
        },
        "zen/code": {
          "name": "Code",
          "limit": {
            "context": 200000,
            "output": 128000
          },
          "description": "Optimized for code generation"
        },
        "zen/reasoning": {
          "name": "Reasoning",
          "limit": {
            "context": 200000,
            "output": 128000
          },
          "description": "Optimized for reasoning tasks"
        },
        "grok-code": {
          "name": "Grok Code (via OpenRouter)",
          "limit": {
            "context": 2000000,
            "output": 131072
          }
        },
        "glm-4.7-free": {
          "name": "GLM 4.7 Free (via OpenRouter)",
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

#### Nutzung

```bash
# Kostenloses Model nutzen
opencode run "Hello" --model opencode-zen/zen/big-pickle

# Code-spezifisch
opencode run "Write a React component" --model opencode-zen/zen/code
```

---

### 18.4 Streamlake Provider

```json
{
  "provider": {
    "streamlake": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "Streamlake",
      "options": {
        "baseURL": "https://vanchin.streamlake.ai/api/gateway/v1/endpoints/kat-coder-pro-v1/claude-code-proxy"
      },
      "models": {
        "kat-coder-pro-v1": {
          "name": "KAT Coder Pro v1",
          "limit": {
            "context": 2000000,
            "output": 128000
          }
        }
      }
    }
  }
}
```

---

### 18.5 XiaoMi MIMO Provider

```json
{
  "provider": {
    "xiaomi": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "XiaoMi MIMO",
      "options": {
        "baseURL": "https://api.xiaomi.ai/v1"
      },
      "models": {
        "mimo-v2-flash": {
          "name": "MIMO v2 Flash",
          "limit": {
            "context": 1000000,
            "output": 65536
          },
          "modalities": {
            "input": ["text", "image"],
            "output": ["text"]
          }
        },
        "mimo-v2-turbo": {
          "name": "MIMO v2 Turbo",
          "limit": {
            "context": 1500000,
            "output": 100000
          }
        }
      }
    }
  }
}
```

---

## 🧠 19. Model Configuration Deep Dive

### 19.1 Qwen 3.5 397B (Ausführlich)

**Spezifikationen:**

| Attribut | Wert |
|----------|------|
| Modell-ID | qwen/qwen3.5-397b-a17b |
| Provider | NVIDIA NIM |
| Architektur | Mixture of Experts (MoE) |
| Gesamt-Parameter | 397B |
| Aktive Parameter | 17B pro Token |
| Context Window | 262.144 Tokens |
| Output Limit | 32.768 Tokens |
| Trainingsdaten | Up to 2024-12 |
| Sprachen | 201 |
| Lizenz | Apache 2.0 |

**Fähigkeiten:**

- Text-Generierung
- Code-Generierung (Best-in-Class)
- Vision-Analyse
- Tool Calling
- Reasoning
- Multimodal

**Konfigurations-Optionen:**

```json
{
  "model": "qwen/qwen3.5-397b-a17b",
  "messages": [
    {
      "role": "user",
      "content": "Your prompt here"
    }
  ],
  "temperature": 0.7,
  "top_p": 0.9,
  "max_tokens": 8192,
  "stream": false,
  "stop": null,
  "presence_penalty": 0,
  "frequency_penalty": 0
}
```

**Parameter-Erklärung:**

| Parameter | Typ | Default | Beschreibung |
|-----------|-----|---------|--------------|
| temperature | float | 0.7 | Kreativität (0-2) |
| top_p | float | 0.9 | Nucleus sampling |
| max_tokens | int | 8192 | Maximale Output-Länge |
| stream | bool | false | Streaming aktivieren |
| stop | array | null | Stop-Sequenzen |
| presence_penalty | float | 0 | Wiederholungsstrafe |
| frequency_penalty | float | 0 | Häufigkeitsstrafe |

**Use Cases:**

1. **Code-Generierung**
   ```bash
   opencode run "Write a REST API with Express.js" --model nvidia-nim/qwen-3.5-397b
   ```

2. **Refactoring**
   ```bash
   opencode run "Refactor this JavaScript to TypeScript" --model nvidia-nim/qwen-3.5-397b
   ```

3. **Code-Review**
   ```bash
   opencode run "Review this code for security issues" --model nvidia-nim/qwen-3.5-397b
   ```

**Performance-Tipps:**

- Timeout auf 120000ms setzen
- First Token Latency: ~70-90s
- Buffer für lange Outputs einplanen
- Fallback für Rate Limits definieren

---

### 19.2 Gemini 3 Pro/Flash (Ausführlich)

**Spezifikationen Gemini 3 Pro:**

| Attribut | Wert |
|----------|------|
| Modell-ID | gemini-3-pro-preview |
| Provider | Google Antigravity |
| Context Window | 2M Tokens |
| Output Limit | 64K Tokens |
| Reasoning | Ja (bis 32K tokens) |
| Multimodal | Text, Image, PDF, Video |

**Spezifikationen Gemini 3 Flash:**

| Attribut | Wert |
|----------|------|
| Modell-ID | gemini-3-flash-preview |
| Provider | Google Antigravity |
| Context Window | 1M Tokens |
| Output Limit | 64K Tokens |
| Reasoning | Ja |
| Multimodal | Text, Image, PDF |

**Thinking-Levels:**

```json
{
  "thinkingConfig": {
    "thinkingBudget": 8192  // oder 32768 für high
  }
}
```

**Use Cases:**

1. **Schnelle Analysen**
   ```bash
   opencode run "Analyze this data" --model google/antigravity-gemini-3-flash
   ```

2. **Tiefgehende Recherche**
   ```bash
   opencode run "Research quantum computing" --model google/antigravity-gemini-3-pro:high
   ```

3. **Multimodal**
   ```bash
   opencode run "Describe this image" --model google/antigravity-gemini-3-flash --image screenshot.png
   ```

---

### 19.3 Claude Sonnet 4.5 (Antigravity)

**Spezifikationen:**

| Attribut | Wert |
|----------|------|
| Modell-ID | claude-sonnet-4-5-20250501 |
| Provider | Google Antigravity |
| Context Window | 200K Tokens |
| Output Limit | 64K Tokens |
| Thinking | Ja (bis 32K) |

**Thinking-Konfiguration:**

```json
{
  "thinkingConfig": {
    "thinkingBudget": 8192  // low
    // oder
    "thinkingBudget": 32768  // max
  }
}
```

**Use Cases:**

1. **Code Review**
   ```bash
   opencode run "Review this code" --model google/antigravity-claude-sonnet-4-5-thinking:max
   ```

2. **Komplexe Probleme**
   ```bash
   opencode run "Solve this algorithm problem" --model google/antigravity-claude-sonnet-4-5-thinking:max
   ```

---

## 💻 20. CLI Usage Comprehensive Guide

### 20.1 Grundbefehle

#### opencode run

Führt einen einzelnen Prompt aus.

```bash
# Basisnutzung
opencode run "Hello World"

# Mit spezifischem Model
opencode run "Write a function" --model nvidia-nim/qwen-3.5-397b

# Mit Agent
opencode run "Create a React app" --agent sisyphus

# Mit Bild
opencode run "Describe this" --image path/to/image.png

# Session fortsetzen
opencode run "Continue" --continue

# Mit Custom System Prompt
opencode run "Task" --system "You are a helpful assistant"
```

#### opencode chat

Startet einen interaktiven Chat.

```bash
# Chat starten
opencode chat

# Mit spezifischem Model
opencode chat --model google/antigravity-gemini-3-pro

# Session fortsetzen
opencode chat --continue
```

#### opencode code

Öffnet den Editor für Code-Änderungen.

```bash
# Code bearbeiten
opencode code

# Mit spezifischem Model
opencode code --model nvidia-nim/qwen-3.5-397b
```

#### opencode docs

Generiert Dokumentation.

```bash
# Dokumentation erstellen
opencode docs

# Für spezifisches Model
opencode docs --model opencode-zen/zen/advanced
```

---

### 20.2 Modell-Verwaltung

#### opencode models

Liste alle verfügbaren Modelle auf.

```bash
# Alle Modelle
opencode models

# Gefiltert nach Provider
opencode models | grep nvidia
opencode models | grep google
opencode models | grep zen

# Mit Details
opencode models --verbose
```

#### Modelle auswählen

```bash
# Interaktiv
opencode
# > /models

# Direkt
/models google/antigravity-gemini-3-flash
```

---

### 20.3 Auth-Verwaltung

#### opencode auth

Authentifizierungsverwaltung.

```bash
# Login starten
opencode auth login

# Logout
opencode auth logout

# Status anzeigen
opencode auth status

# Token refreshen
opencode auth refresh

# Accounts auflisten
opencode auth list
```

---

### 20.4 Session-Management

#### Sessions auflisten

```bash
# Alle Sessions
opencode sessions list

# Aktive Session
opencode sessions active

# Session-Details
opencode sessions info ses_xxx
```

#### Session fortsetzen

```bash
# Letzte Session fortsetzen
opencode --continue

# Spezifische Session
opencode --session ses_xxx

# Mit Fork
opencode --continue --fork
```

#### Session speichern

```bash
# Session speichern
opencode sessions save

# Mit Namen
opencode sessions save my-project
```

---

### 20.5 Konfiguration

#### opencode config

Konfigurationsverwaltung.

```bash
# Config anzeigen
opencode config view

# Config bearbeiten
opencode config edit

# Config validieren
opencode config validate

# Config zurücksetzen
opencode config reset
```

---

### 20.6 Plugin-Management

#### opencode plugins

Plugin-Verwaltung.

```bash
# Plugins auflisten
opencode plugins list

# Plugin installieren
opencode plugins install <plugin-name>

# Plugin deinstallieren
opencode plugins uninstall <plugin-name>

# Plugin updaten
opencode plugins update <plugin-name>

# Plugin suchen
opencode plugins search <term>
```

---

### 20.7 Agent-Management

#### opencode agent

Agent-Verwaltung.

```bash
# Agenten auflisten
opencode agent list

# Agent starten
opencode agent run <agent-name>

# Agent-Status
opencode agent status <agent-name>

# Agent konfigurieren
opencode agent config <agent-name>
```

---

## 🔌 21. Plugin System Comprehensive

### 21.1 Plugin-Architektur

OpenCode nutzt ein Plugin-System, das die Funktionalität erweitert.

```
┌─────────────────────────────────────────┐
│            OpenCode Core                │
├─────────────────────────────────────────┤
│  - Core CLI                            │
│  - Config Management                   │
│  - Session Management                  │
│  - Model Routing                       │
└─────────────────┬───────────────────────┘
                  │
        ┌─────────┴─────────┐
        │                   │
        ▼                   ▼
┌───────────────┐   ┌───────────────┐
│   Plugins     │   │    MCP        │
│               │   │   Servers     │
├───────────────┤   ├───────────────┤
│ antigravity   │   │   serena      │
│ oh-my-opencode│   │   context7    │
│ qwencode-auth │   │   tavily      │
└───────────────┘   └───────────────┘
```

### 21.2 Plugin-Typen

#### Authentication Plugins

Erweitern Authentifizierungsoptionen.

```json
{
  "plugin": [
    "opencode-antigravity-auth",
    "opencode-qwencode-auth"
  ]
}
```

#### Provider Plugins

Fügen neue Provider hinzu.

```bash
# Third-Party Plugins
opencode plugins install @example/opencode-provider-plugin
```

#### UI Plugins

Erweitern die Benutzeroberfläche.

```json
{
  "plugin": [
    "oh-my-opencode"
  ]
}
```

### 21.3 Plugin-Entwicklung

#### Plugin-Struktur

```
my-plugin/
├── package.json
├── src/
│   ├── index.ts
│   ├── commands/
│   ├── hooks/
│   └── providers/
├── tests/
└── README.md
```

#### Basis-Plugin

```typescript
import { Plugin } from 'opencode';

export default class MyPlugin implements Plugin {
  name = 'my-plugin';
  version = '1.0.0';

  async onLoad() {
    console.log('Plugin geladen!');
  }

  async onCommand(command: string, args: any) {
    if (command === 'my-command') {
      return this.handleMyCommand(args);
    }
  }

  private async handleMyCommand(args: any) {
    // Implementierung
  }
}
```

#### Plugin registrieren

```json
{
  "plugin": [
    "my-plugin@1.0.0"
  ]
}
```

---

## 🔗 22. MCP Integration Comprehensive

### 22.1 MCP-Architektur

MCP (Model Context Protocol) ermöglicht OpenCode die Integration mit externen Tools.

```
┌──────────────┐     MCP      ┌──────────────┐
│   OpenCode   │◄────────────►│  MCP Server  │
│              │   stdio/     │              │
│  - Session   │    HTTP      │  - Tools     │
│  - Messages  │              │  - Resources │
│  - Files     │              │  - Prompts   │
└──────────────┘              └──────────────┘
```

### 22.2 MCP-Server-Typen

#### Local MCP Servers

Laufen lokal auf dem System.

```json
{
  "mcp": {
    "serena": {
      "type": "local",
      "command": ["uvx", "serena", "start-mcp-server"],
      "enabled": true
    }
  }
}
```

#### Remote MCP Servers

Laufen extern (HTTP/WebSocket).

```json
{
  "mcp": {
    "linear": {
      "type": "remote",
      "url": "https://mcp.linear.app/sse",
      "enabled": true
    }
  }
}
```

### 22.3 MCP-Tools

#### Tools definieren

```typescript
// tool-definition.ts
import { Tool } from '@modelcontextprotocol/sdk';

export const myTool: Tool = {
  name: 'my_tool',
  description: 'Beschreibung des Tools',
  inputSchema: {
    type: 'object',
    properties: {
      param1: {
        type: 'string',
        description: 'Erster Parameter'
      },
      param2: {
        type: 'number',
        description: 'Zweiter Parameter'
      }
    },
    required: ['param1']
  }
};
```

#### Tool aufrufen

```bash
# Via OpenCode
opencode run "Nutze das my_tool" --tool my_tool
```

### 22.4 Eigene MCP-Server erstellen

#### Server-Implementierung

```typescript
import { Server } from '@modelcontextprotocol/sdk/server';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio';

const server = new Server(
  { name: 'my-mcp-server', version: '1.0.0' },
  { capabilities: { tools: {} } }
);

server.setRequestHandler(ListToolsRequestSchema, async () => ({
  tools: [{
    name: 'my_tool',
    description: 'Ein nützliches Tool',
    inputSchema: {
      type: 'object',
      properties: {
        input: { type: 'string' }
      }
    }
  }]
}));

server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args } = request.params;
  
  if (name === 'my_tool') {
    return { result: `Verarbeitet: ${args.input}` };
  }
  
  throw new Error(`Unknown tool: ${name}`);
});

const transport = new StdioServerTransport();
server.connect(transport);
```

#### In OpenCode integrieren

```json
{
  "mcp": {
    "my-mcp": {
      "type": "local",
      "command": ["node", "/path/to/mcp-server.js"],
      "enabled": true
    }
  }
}
```

---

## 🔒 23. Authentication Security

### 23.1 Token-Management

#### Token-Speicher

Tokens werden verschlüsselt gespeichert.

```bash
# Token-Datei Berechtigungen
ls -la ~/.config/opencode/
# -rw------- 1 user staff 256 Feb 18 10:00 antigravity-accounts.json
```

#### Token-Refresh

Automatischer Token-Refresh:

```bash
# Manuell
opencode auth refresh

# Automatisch
# Wird bei Ablauf automatisch durchgeführt
```

### 23.2 OAuth-Sicherheit

#### Security-Best-Practices

1. **Private Accounts nutzen**
   - NICHT Google Workspace
   - Private Gmail-Accounts

2. **Multi-Account für Load Balancing**
   - Mehrere Accounts bei Rate Limits
   - Automatisches Failover

3. **Token-Storage**
   - Verschlüsselte Speicherung
   - Sichere Berechtigungen

### 23.3 API-Key-Sicherheit

#### Environment Variables

```bash
# NVIDIA API Key
export NVIDIA_API_KEY="nvapi-xxx"

# In der Shell setzen
# NICHT in Config-Files speichern!
```

#### Vault-Integration

```bash
# Optional: HashiCorp Vault
eval $(vault env nvidia-api)
```

---

## ⚡ 24. Performance Optimization

### 24.1 Caching

#### Response-Caching

```json
{
  "cache": {
    "enabled": true,
    "ttl": 3600,
    "maxSize": 1000
  }
}
```

#### Model-Caching

Einige Modelle können gecached werden.

```bash
# Caching aktivieren
opencode config set cache.enabled true
```

### 24.2 Connection Pooling

#### MCP-Server

```bash
# MCP-Server persistent halten
# NICHT bei jeder Anfrage neu starten

# In opencode.json:
{
  "mcp": {
    "serena": {
      "type": "local",
      "command": ["uvx", "serena", "start-mcp-server"],
      "enabled": true,
      "keepAlive": true
    }
  }
}
```

### 24.3 Rate-Limit-Management

#### Fallback-Strategie

```json
{
  "models": {
    "defaults": {
      "primary": "nvidia/qwen/qwen3.5-397b-a17b",
      "fallbacks": [
        "google/antigravity-gemini-3-flash",
        "opencode-zen/zen/big-pickle"
      ]
    }
  }
}
```

#### Backoff-Konfiguration

```json
{
  "retry": {
    "maxRetries": 3,
    "baseDelay": 1000,
    "maxDelay": 60000,
    "backoffMultiplier": 2
  }
}
```

---

## 🧪 25. Testing Debugging

### 25.1 Testing-Strategien

#### Unit-Testing

```bash
# Tests ausführen
opencode test

# Spezifischer Test
opencode test --file path/to/test.ts

# Mit Coverage
opencode test --coverage
```

#### Integration-Testing

```bash
# Integration-Tests
opencode test --integration

# E2E-Tests
opencode test --e2e
```

### 25.2 Debugging

#### Debug-Mode

```bash
# Debug aktivieren
opencode --debug run "Task"

# Oder
export DEBUG=opencode:*
opencode run "Task"
```

#### Logging

```bash
# Logs anzeigen
opencode logs

# Mit Timestamp
opencode logs --timestamp

# Filter
opencode logs --level error
```

#### Trace-Collection

```bash
# Trace starten
opencode trace start

# Trace stoppen
opencode trace stop

# Trace analysieren
opencode trace analyze trace.json
```

### 25.3 Health-Checks

```bash
# OpenCode Health
opencode doctor

# OpenClaw Health
openclaw doctor --fix

# Spezifische Prüfung
opencode doctor --check models
opencode doctor --check auth
opencode doctor --check mcp
```

---

## 🔮 26. Future Roadmap

### 26.1 Geplante Features

#### Q1 2026

- Erweiterte Agent-Fähigkeiten
- Verbessertes Multi-Model-Routing
- Erweiterte MCP-Integration

#### Q2 2026

- Bessere Code-Generierung
- Verbesserte Vision-Fähigkeiten
- Erweiterte Tool-Integration

#### Q3-Q4 2026

- Enterprise-Features
- Erweiterte Security
- Verbesserte Performance

### 26.2 Neue Provider

Geplante Provider-Integrationen:

- Weitere NVIDIA NIM Modelle
- Mehr OpenCode ZEN Modelle
- Custom Provider Support

### 26.3 Plugin-Ecosystem

- Plugin Marketplace
- Plugin-Entwickler-Dokumentation
- Community-Plugins

---

## ⚙️ 27. Advanced Configuration

### 27.1 Custom Provider

Eigene Provider hinzufügen:

```json
{
  "provider": {
    "my-provider": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "My Custom Provider",
      "options": {
        "baseURL": "https://api.myprovider.com/v1",
        "apiKey": "${MY_PROVIDER_API_KEY}"
      },
      "models": {
        "my-model": {
          "id": "my-model-v1",
          "name": "My Model",
          "limit": {
            "context": 100000,
            "output": 10000
          }
        }
      }
    }
  }
}
```

### 27.2 Custom Agents

Eigene Agenten definieren:

```json
{
  "agent": {
    "my-custom-agent": {
      "model": {
        "provider": "my-provider",
        "id": "my-model"
      },
      "systemPrompt": "Du bist ein spezialisierter Agent für..."
    }
  }
}
```

### 27.3 Environment-spezifische Configs

```bash
# Development
export OPENCODE_ENV=development

# Production
export OPENCODE_ENV=production
```

---

## 💰 28. Cost Optimization Strategies

### 28.1 Model-Selection

| Task | Empfohlenes Model | Kosten |
|------|-------------------|-------|
| Recherche | zen/big-pickle | 🆓 |
| Exploration | zen/big-pickle | 🆓 |
| Einfache Tasks | gemini-3-flash | 💰 |
| Komplexe Tasks | gemini-3-pro | 💰💰 |
| Code-Intensive | qwen-3.5-397b | 💰💰 |

### 28.2 Token-Management

```bash
# Max Tokens begrenzen
opencode run "Task" --max-tokens 1000

# Temperature optimieren
# Niedrig für präzise Tasks
opencode run "Fakten" --temperature 0.1

# Hoch für kreative Tasks
opencode run "Brainstorming" --temperature 0.9
```

### 28.3 Caching-Strategien

```json
{
  "cache": {
    "enabled": true,
    "similarityThreshold": 0.9,
    "maxAge": 86400
  }
}
```

---

## 🏢 29. Enterprise Deployment

### 29.1 Multi-User Setup

```json
{
  "enterprise": {
    "multiUser": true,
    "userManagement": "ldap",
    "sso": {
      "enabled": true,
      "provider": "okta"
    }
  }
}
```

### 29.2 Security-Compliance

```json
{
  "security": {
    "encryption": "AES-256",
    "auditLogging": true,
    "dataResidency": "EU",
    "compliance": ["GDPR", "SOC2"]
  }
}
```

### 29.3 Monitoring

```json
{
  "monitoring": {
    "enabled": true,
    "metricsEndpoint": "https://metrics.company.com",
    "alerting": {
      "enabled": true,
      "channels": ["slack", "email"]
    }
  }
}
```

---

## 🔧 30. Troubleshooting Advanced

### 30.1 Connection-Probleme

#### Problem: "Connection timeout"

**Lösung:**
```bash
# Timeout erhöhen
opencode config set timeout 120000

# Netzwerk prüfen
ping api.provider.com

# DNS prüfen
nslookup api.provider.com
```

#### Problem: "SSL certificate error"

**Lösung:**
```bash
# Zertifikate aktualisieren
update-ca-certificates

# Oder: Custom CA setzen
export REQUESTS_CA_BUNDLE=/path/to/ca.crt
```

### 30.2 Auth-Probleme

#### Problem: "Token expired"

**Lösung:**
```bash
# Token refreshen
opencode auth refresh

# Oder neu anmelden
opencode auth logout
opencode auth login
```

#### Problem: "Invalid credentials"

**Lösung:**
```bash
# Auth neu starten
opencode auth logout
opencode auth login

# Account-Datei prüfen
cat ~/.config/opencode/antigravity-accounts.json | jq .
```

### 30.3 Model-Probleme

#### Problem: "Model not supported"

**Lösung:**
```bash
# Verfügbare Modelle prüfen
opencode models

# Config aktualisieren
opencode config edit
```

#### Problem: "Rate limit exceeded"

**Lösung:**
```bash
# Warten (60 Sekunden)
# Oder Fallback nutzen
opencode run "Task" --fallback
```

### 30.4 MCP-Probleme

#### Problem: "MCP server not responding"

**Lösung:**
```bash
# Server neustarten
pkill -f serena
serena start-mcp-server

# Oder: Health-Check
opencode doctor --check mcp
```

#### Problem: "Tool not found"

**Lösung:**
```bash
# MCP aktivieren
opencode config set mcp.serena.enabled true

# Oder neu installieren
npm install -g @anthropics/context7-mcp
```

---

## 📊 Anhang: Config-Referenz

### Vollständige Config-Optionen

```json
{
  "$schema": "https://opencode.ai/config.json",
  "model": "string",
  "small_model": "string",
  "default_agent": "string",
  "theme": "string",
  "autoupdate": "boolean",
  "plugin": ["string"],
  "provider": {
    "[provider-name]": {
      "npm": "string",
      "name": "string",
      "options": {
        "baseURL": "string",
        "timeout": "number",
        "apiKey": "string",
        "headers": "object"
      },
      "models": {
        "[model-name]": {
          "id": "string",
          "name": "string",
          "limit": {
            "context": "number",
            "output": "number"
          },
          "modalities": "object",
          "variants": "object"
        }
      }
    }
  },
  "mcp": {
    "[mcp-name]": {
      "type": "local|remote",
      "command": ["string"],
      "url": "string",
      "enabled": "boolean",
      "environment": "object"
    }
  },
  "agent": {
    "[agent-name]": {
      "model": {
        "provider": "string",
        "id": "string"
      },
      "systemPrompt": "string"
    }
  },
  "cache": {
    "enabled": "boolean",
    "ttl": "number",
    "maxSize": "number"
  },
  "retry": {
    "maxRetries": "number",
    "baseDelay": "number",
    "maxDelay": "number",
    "backoffMultiplier": "number"
  }
}
```

---

## 📚 Referenzen

### Offizielle Dokumentation
- OpenCode Docs: https://opencode.ai/docs/
- NVIDIA NIM: https://build.nvidia.com/
- Antigravity Plugin: https://github.com/shekohex/opencode-google-antigravity-auth

### Community
- Discord: https://discord.gg/opencode
- GitHub: https://github.com/opencode/opencode
- Twitter: @opencode

### Support
- Issues: https://github.com/opencode/opencode/issues
- Discussions: https://github.com/opencode/opencode/discussions

---

**Updated:** 2026-02-18  
**Status:** COMPLETE ✅  
**All Plugins:** Oh My OpenCode + Qwen OAuth + Antigravity

---

## 🔬 31. OpenCode Server Architecture Deep Dive

### 31.1 Server Components Overview

The OpenCode Server is built on a modular architecture that enables flexible scaling and extensibility. The server consists of several core components that work together to provide a seamless AI coding experience.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        OPENCODE SERVER ARCHITECTURE                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        CLIENT LAYER                                   │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │   │
│  │  │  CLI        │  │  Web UI     │  │  API Client │              │   │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘              │   │
│  │         │                 │                 │                       │   │
│  └─────────┼─────────────────┼─────────────────┼───────────────────────┘   │
│            │                 │                 │                           │
│            ▼                 ▼                 ▼                           │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      API GATEWAY                                    │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │   │
│  │  │  Auth      │  │  Rate Limit │  │  Load      │              │   │
│  │  │  Middleware│  │  Middleware │  │  Balancer  │              │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘              │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│            ┌───────────────────────┼───────────────────────┐              │
│            │                       │                       │              │
│            ▼                       ▼                       ▼              │
│  ┌─────────────────┐   ┌─────────────────┐   ┌─────────────────┐        │
│  │  SESSION       │   │  PROVIDER      │   │  MCP           │        │
│  │  SERVICE       │   │  ORCHESTRATOR │   │  GATEWAY       │        │
│  │                │   │                │   │                │        │
│  │  - Create     │   │  - Model      │   │  - Tools       │        │
│  │  - Resume     │   │    Routing    │   │  - Resources   │        │
│  │  - Context    │   │  - Fallback   │   │  - Prompts     │        │
│  │  - History    │   │  - Retry      │   │                │        │
│  └────────┬────────┘   └────────┬────────┘   └────────┬────────┘        │
│           │                     │                     │                 │
│           └─────────────────────┼─────────────────────┘                 │
│                                 │                                         │
│                                 ▼                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     STORAGE LAYER                                   │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │   │
│  │  │  Sessions  │  │  Models     │  │  Cache     │              │   │
│  │  │  (Postgres)│  │  (Files)    │  │  (Redis)   │              │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘              │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 31.2 Session Management Deep Dive

The session management system is the heart of OpenCode's conversational capabilities. Each session maintains context across multiple interactions, enabling sophisticated multi-turn conversations and complex task execution.

#### Session Lifecycle

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        SESSION LIFECYCLE                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   CREATE ──► ACTIVE ──► PAUSED ──► ACTIVE ──► COMPLETED                    │
│     │          │          │          │          │                            │
│     │          │          │          │          │                            │
│     ▼          ▼          ▼          ▼          ▼                            │
│  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐                      │
│  │New   │  │Inter-│  │Idle   │  │Re-   │  │Task  │                      │
│  │Session│  │active│  │Timeout│  │sume  │  │Done  │                      │
│  └──────┘  └──────┘  └──────┘  └──────┘  └──────┘                      │
│                                                                              │
│  EVENTS:                                                                    │
│  ┌────────────────────────────────────────────────────────────────────┐    │
│  │ 1. Session Created: Initialize context, load system prompt          │    │
│  │ 2. User Message: Add to history, trigger model processing           │    │
│  │ 3. Model Response: Generate, store, return to user                  │    │
│  │ 4. Context Update: Expand context window, manage memory             │    │
│  │ 5. Session Paused: Persist state, release resources                 │    │
│  │ 6. Session Resumed: Restore state, continue conversation            │    │
│  │ 7. Session Completed: Finalize, cleanup, archive                    │    │
│  └────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Session State Structure

```typescript
interface Session {
  // Core identifiers
  id: string;                    // Unique session ID (ses_xxx)
  title: string;                 // User-defined session title
  created_at: Date;              // Session creation timestamp
  updated_at: Date;              // Last interaction timestamp
  
  // Context management
  context: {
    system_prompt: string;       // Base system prompt
    max_tokens: number;          // Maximum context tokens
    used_tokens: number;         // Current token usage
    available_tokens: number;    // Remaining tokens
  };
  
  // Model configuration
  model: {
    provider: string;            // Provider ID (nvidia, google, etc.)
    model_id: string;            // Model identifier
    temperature: number;         // Generation temperature
    max_tokens: number;         // Output token limit
    thinking_level?: string;     // For models with thinking support
  };
  
  // Message history
  messages: Message[];          // Conversation history
  parts: Part[];                // Multi-modal content parts
  
  // State management
  state: 'creating' | 'active' | 'paused' | 'completed' | 'error';
  
  // Metadata
  metadata: {
    user_id?: string;           // Associated user ID
    tags?: string[];            // Session tags
    parent_session?: string;     // Forked from session ID
    statistics?: SessionStats;   // Usage statistics
  };
  
  // MCP integrations
  mcp_context?: {
    active_tools: string[];     // Currently enabled MCP tools
    tool_results: Map<string, any>; // Tool execution results
  };
}

interface Message {
  id: string;                   // Unique message ID
  role: 'user' | 'assistant' | 'system' | 'tool';
  content: string;              // Text content
  parts?: Part[];               // Multi-modal parts
  created_at: Date;             // Message timestamp
  
  // For assistant messages
  model?: string;               // Model used
  usage?: {
    input_tokens: number;
    output_tokens: number;
    total_tokens: number;
  };
  finish_reason?: string;        // Completion reason
  
  // For tool messages
  tool_calls?: ToolCall[];      // Tool invocations
  tool_call_id?: string;        // Reference to tool call
}

interface Part {
  type: 'text' | 'file' | 'image' | 'tool' | 'reasoning' | 'agent';
  
  // Text part
  text?: string;
  
  // File part
  file?: {
    mime: string;               // MIME type (image/jpeg, application/pdf, etc.)
    filename: string;
    url: string;                // Data URL or file URL
  };
  
  // Image part (legacy compatibility)
  image?: {
    url: string;
    detail?: 'low' | 'high';
  };
  
  // Tool part
  tool?: {
    id: string;
    name: string;
    input: Record<string, any>;
  };
  
  // Reasoning part (for thinking models)
  reasoning?: string;
  
  // Agent reference
  agent?: {
    id: string;
    name: string;
  };
}
```

#### Session API Endpoints

```yaml
/api/session:
  POST:
    summary: Create new session
    requestBody:
      content:
        application/json:
          schema:
            type: object
            properties:
              title:
                type: string
                description: Session title
              model:
                $ref: '#/components/schemas/ModelConfig'
              systemPrompt:
                type: string
              metadata:
                type: object
    responses:
      201:
        description: Session created
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Session'

/api/session/{session_id}:
  GET:
    summary: Get session details
    parameters:
      - name: session_id
        in: path
        required: true
        schema:
          type: string
    responses:
      200:
        description: Session details
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Session'

  DELETE:
    summary: Delete session
    parameters:
      - name: session_id
        in: path
        required: true
        schema:
          type: string
    responses:
      204:
        description: Session deleted

/api/session/{session_id}/prompt:
  POST:
    summary: Send prompt to session
    parameters:
      - name: session_id
        in: path
        required: true
        schema:
          type: string
    requestBody:
      content:
        application/json:
          schema:
            type: object
            properties:
              parts:
                type: array
                items:
                  $ref: '#/components/schemas/Part'
              model:
                $ref: '#/components/schemas/ModelConfig'
              stream:
                type: boolean
                default: false
    responses:
      200:
        description: Prompt response
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Message'

/api/session/{session_id}/prompt_async:
  POST:
    summary: Send prompt asynchronously
    parameters:
      - name: session_id
        in: path
        required: true
        schema:
          type: string
    requestBody:
      content:
        application/json:
          schema:
            type: object
            properties:
              parts:
                type: array
                items:
                  $ref: '#/components/schemas/Part'
              model:
                $ref: '#/components/schemas/ModelConfig'
    responses:
      202:
        description: Async prompt accepted
        content:
          application/json:
            schema:
              type: object
              properties:
                task_id:
                  type: string
                status_url:
                  type: string
                  format: uri

/api/session/{session_id}/message:
  GET:
    summary: Get session messages
    parameters:
      - name: session_id
        in: path
        required: true
        schema:
          type: string
      - name: limit
        in: query
        schema:
          type: integer
          default: 50
      - name: before
        in: query
        schema:
          type: string
          description: Message ID for pagination
    responses:
      200:
        description: Message list
        content:
          application/json:
            schema:
              type: array
              items:
                $ref: '#/components/schemas/Message'
```

### 31.3 Provider Orchestration

The provider orchestration layer manages communication with multiple AI providers, handling fallback logic, rate limiting, and error recovery.

#### Provider Interface

```typescript
interface Provider {
  // Provider identification
  id: string;
  name: string;
  
  // Capabilities
  capabilities: {
    streaming: boolean;
    vision: boolean;
    functionCalling: boolean;
    reasoning: boolean;
    multiModal: boolean;
  };
  
  // Model limits
  limits: {
    maxContextTokens: number;
    maxOutputTokens: number;
    maxImagesPerRequest?: number;
  };
  
  // Rate limiting
  rateLimits: {
    requestsPerMinute: number;
    requestsPerDay: number;
    tokensPerMinute?: number;
  };
  
  // Authentication
  auth: {
    type: 'apiKey' | 'oauth' | 'none';
    config: AuthConfig;
  };
  
  // API endpoints
  endpoints: {
    chat?: string;
    completion?: string;
    embedding?: string;
    vision?: string;
  };
  
  // Methods
  initialize(): Promise<void>;
  chat(request: ChatRequest): Promise<ChatResponse>;
  complete(request: CompletionRequest): Promise<CompletionResponse>;
  stream(request: ChatRequest): AsyncGenerator<StreamChunk>;
  validateAuth(): Promise<boolean>;
  refreshAuth(): Promise<void>;
}
```

#### Fallback Chain Implementation

```typescript
class ProviderOrchestrator {
  private providers: Map<string, Provider>;
  private fallbackChains: Map<string, string[]>;
  
  constructor(config: OrchestratorConfig) {
    this.providers = new Map();
    this.fallbackChains = config.fallbackChains;
  }
  
  async executeWithFallback(
    request: ChatRequest,
    primaryProvider: string,
    options: FallbackOptions = {}
  ): Promise<ChatResponse> {
    const chain = this.getFallbackChain(primaryProvider);
    const errors: ProviderError[] = [];
    
    for (const providerId of chain) {
      const provider = this.providers.get(providerId);
      
      if (!provider) {
        errors.push({
          provider: providerId,
          error: new Error('Provider not available'),
          timestamp: new Date()
        });
        continue;
      }
      
      try {
        // Check rate limits
        if (await this.isRateLimited(providerId)) {
          if (options.waitForRateLimit) {
            await this.waitForRateLimitReset(providerId);
          } else {
            continue; // Skip to next provider
          }
        }
        
        // Execute request
        const response = await provider.chat(request);
        
        // Success - return response
        return response;
        
      } catch (error) {
        const providerError = this.classifyError(error, providerId);
        errors.push(providerError);
        
        // Log error for monitoring
        this.logError(providerError);
        
        // Check if we should retry
        if (!this.isRetryable(providerError)) {
          throw new MultiProviderError(errors);
        }
        
        // Continue to next provider in chain
        continue;
      }
    }
    
    // All providers failed
    throw new MultiProviderError(errors);
  }
  
  private classifyError(error: any, providerId: string): ProviderError {
    if (error.status === 401) {
      return {
        provider: providerId,
        error,
        type: 'auth',
        retryable: false
      };
    }
    
    if (error.status === 429) {
      return {
        provider: providerId,
        error,
        type: 'rate_limit',
        retryable: true,
        retryAfter: error.retryAfter || 60
      };
    }
    
    if (error.status === 500 || error.status === 503) {
      return {
        provider: providerId,
        error,
        type: 'server_error',
        retryable: true
      };
    }
    
    if (error.code === 'ETIMEDOUT') {
      return {
        provider: providerId,
        error,
        type: 'timeout',
        retryable: true
      };
    }
    
    return {
      provider: providerId,
      error,
      type: 'unknown',
      retryable: false
    };
  }
  
  private async isRateLimited(providerId: string): Promise<boolean> {
    const provider = this.providers.get(providerId);
    if (!provider) return true;
    
    const usage = await this.getRateLimitUsage(providerId);
    const limits = provider.rateLimits;
    
    return usage.requestsPerMinute >= limits.requestsPerMinute;
  }
}
```

### 31.4 MCP Gateway Architecture

The MCP (Model Context Protocol) Gateway enables OpenCode to integrate with external tools and services, extending its capabilities beyond built-in features.

#### MCP Communication Flow

```
┌─────────────┐     JSON-RPC      ┌─────────────┐     stdio/HTTP    ┌─────────────┐
│   OpenCode  │◄────────────────►│  MCP        │◄────────────────►│  External   │
│   Server    │    2.0           │  Gateway    │                   │  Tool       │
│             │                  │             │                   │             │
│  ┌────────┐│                  │  ┌────────┐ │                   │  ┌────────┐ │
│  │ Session ││                  │  │ Tool   │ │                   │  │  File  │ │
│  │ Handler ││                  │  │ Router │ │                   │  │ System │ │
│  └────────┘│                  │  └────────┘ │                   │  └────────┘ │
│            │                  │  ┌────────┐ │                   │  ┌────────┐ │
│  ┌────────┐│                  │  │ Result │ │                   │  │  API   │ │
│  │ Message ││                  │  │ Cache  │ │                   │  │  (REST│ │
│  │ Parser  ││                  │  └────────┘ │                   │  │ GraphQL│ │
│  └────────┘│                  │             │                   │  └────────┘ │
└─────────────┘                  └─────────────┘                   └─────────────┘

MESSAGE FLOW:
1. User sends message with tool call
2. Session Handler identifies tool call in message
3. Tool Router resolves tool name to MCP server
4. MCP Gateway forwards request to external tool
5. External tool processes request
6. Result returned to MCP Gateway
7. Result cached for reuse
8. Result injected into conversation context
9. Continue conversation with tool result
```

#### MCP Tool Definition Schema

```typescript
// MCP Tool Schema (JSON Schema)
const ToolSchema = {
  type: 'object',
  properties: {
    name: {
      type: 'string',
      description: 'Unique identifier for the tool',
      pattern: '^[a-z][a-z0-9_]*$'
    },
    description: {
      type: 'string',
      description: 'Human-readable description of what the tool does'
    },
    inputSchema: {
      type: 'object',
      properties: {
        type: {
          type: 'string',
          const: 'object'
        },
        properties: {
          type: 'object',
          additionalProperties: {
            type: 'object',
            properties: {
              type: { type: 'string' },
              description: { type: 'string' },
              enum: { type: 'array' },
              default: { },
              minimum: { type: 'number' },
              maximum: { type: 'number' },
              minLength: { type: 'integer' },
              maxLength: { type: 'integer' },
              pattern: { type: 'string' },
              format: { type: 'string' }
            }
          },
          required: { type: 'array', items: { type: 'string' } }
        }
      }
    }
  },
  required: ['name', 'description', 'inputSchema']
};

// Example tool definition
const exampleTool = {
  name: 'filesystem_read',
  description: 'Read contents of a file from the filesystem',
  inputSchema: {
    type: 'object',
    properties: {
      path: {
        type: 'string',
        description: 'Absolute or relative path to the file'
      },
      encoding: {
        type: 'string',
        enum: ['utf-8', 'base64', 'binary'],
        default: 'utf-8'
      },
      maxBytes: {
        type: 'number',
        description: 'Maximum number of bytes to read',
        minimum: 1,
        maximum: 1048576 // 1MB
      },
      lineOffset: {
        type: 'number',
        description: 'Line number to start reading from (1-indexed)',
        minimum: 1
      },
      lineLimit: {
        type: 'number',
        description: 'Maximum number of lines to read',
        minimum: 1
      }
    },
    required: ['path']
  }
};
```

### 31.5 Storage Layer Architecture

The storage layer provides persistent storage for sessions, model configurations, and cached data.

#### Database Schema

```sql
-- Sessions table
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    state TEXT NOT NULL DEFAULT 'active',
    
    -- Context
    system_prompt TEXT,
    max_tokens INTEGER DEFAULT 128000,
    used_tokens INTEGER DEFAULT 0,
    
    -- Model
    provider TEXT NOT NULL,
    model_id TEXT NOT NULL,
    temperature DECIMAL(3,2) DEFAULT 0.7,
    model_max_tokens INTEGER DEFAULT 8192,
    thinking_level TEXT,
    
    -- Metadata
    user_id TEXT,
    tags TEXT[],
    parent_session TEXT,
    FOREIGN KEY (parent_session) REFERENCES sessions(id)
);

-- Messages table
CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    role TEXT NOT NULL,
    content TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- For assistant messages
    model TEXT,
    input_tokens INTEGER DEFAULT 0,
    output_tokens INTEGER DEFAULT 0,
    total_tokens INTEGER DEFAULT 0,
    finish_reason TEXT,
    
    -- Tool references
    tool_call_id TEXT,
    tool_call_name TEXT,
    
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Parts table (for multi-modal content)
CREATE TABLE parts (
    id TEXT PRIMARY KEY,
    message_id TEXT NOT NULL,
    part_index INTEGER NOT NULL,
    type TEXT NOT NULL,
    
    -- Text content
    text_content TEXT,
    
    -- File content
    file_mime TEXT,
    file_filename TEXT,
    file_url TEXT,
    
    -- Tool content
    tool_id TEXT,
    tool_name TEXT,
    tool_input JSONB,
    
    -- Reasoning content
    reasoning_content TEXT,
    
    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE
);

-- Rate limit tracking
CREATE TABLE rate_limits (
    provider TEXT NOT NULL,
    window_start TIMESTAMP WITH TIME ZONE NOT NULL,
    request_count INTEGER DEFAULT 0,
    token_count INTEGER DEFAULT 0,
    
    PRIMARY KEY (provider, window_start)
);

-- Indexes
CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_state ON sessions(state);
CREATE INDEX idx_messages_session ON messages(session_id);
CREATE INDEX idx_messages_created ON messages(created_at);
CREATE INDEX idx_parts_message ON parts(message_id, part_index);
CREATE INDEX idx_rate_limits_provider ON rate_limits(provider, window_start);
```

#### Cache Strategy

```typescript
interface CacheStrategy {
  // Session context cache
  sessionContext: {
    maxAge: number;           // milliseconds
    maxSize: number;          // entries
    evictionPolicy: 'lru' | 'lfu' | 'ttl';
  };
  
  // Model responses cache
  modelResponse: {
    enabled: boolean;
    maxAge: number;
    similarityThreshold: number;  // For semantic caching
    keyGenerator: (request: ChatRequest) => string;
  };
  
  // Tool results cache
  toolResult: {
    enabled: boolean;
    maxAge: number;
    keyGenerator: (toolName: string, input: any) => string;
  };
}

// Redis-based caching implementation
class CacheManager {
  private redis: Redis;
  private strategies: Map<string, CacheStrategy>;
  
  async get<T>(key: string): Promise<T | null> {
    const cached = await this.redis.get(key);
    if (!cached) return null;
    
    return JSON.parse(cached) as T;
  }
  
  async set<T>(key: string, value: T, ttl: number): Promise<void> {
    await this.redis.setex(key, ttl, JSON.stringify(value));
  }
  
  async getOrSet<T>(
    key: string,
    factory: () => Promise<T>,
    ttl: number
  ): Promise<T> {
    const cached = await this.get<T>(key);
    if (cached) return cached;
    
    const value = await factory();
    await this.set(key, value, ttl);
    
    return value;
  }
  
  // Semantic caching for model responses
  async getSimilarResponse(
    request: ChatRequest,
    threshold: number = 0.9
  ): Promise<ChatResponse | null> {
    const requestHash = await this.hashRequest(request);
    
    // Find similar cached requests
    const similarKeys = await this.redis.keys(`cache:response:*`);
    
    for (const key of similarKeys) {
      const cached = await this.get<CachedResponse>(key);
      if (!cached) continue;
      
      const similarity = await this.computeSimilarity(
        requestHash,
        cached.requestHash
      );
      
      if (similarity >= threshold) {
        return cached.response;
      }
    }
    
    return null;
  }
}
```

---

## 🔌 32. Advanced API Reference

### 32.1 REST API Overview

The OpenCode Server provides a comprehensive REST API for programmatic access to all features.

#### Base URL and Authentication

```yaml
baseURL: http://localhost:8080/api/v1

authentication:
  type: bearer_token
  header: Authorization
  format: "Bearer <token>"

  # Token types:
  # - session_token: Short-lived token for API access
  # - personal_token: Long-lived personal access token
  # - oauth_token: OAuth 2.0 access token

contentTypes:
  - application/json
  - application/json; charset=utf-8
```

#### Global Headers

```yaml
headers:
  Required:
    - name: Authorization
      description: Bearer token for authentication
      example: "Bearer sk-oc_xxxxx"
    
    - name: Content-Type
      description: Request body content type
      value: "application/json"

  Optional:
    - name: X-Request-ID
      description: Client-generated request ID for tracing
      example: "req_abc123"
    
    - name: X-Session-ID
      description: Continue existing session
      example: "ses_xyz789"
    
    - name: X-Trace-Context
      description: Distributed tracing context
      example: "00-abc123..."
```

### 32.2 Authentication Endpoints

#### POST /api/v1/auth/token

Generate an API access token.

```yaml
summary: Generate access token
description: Exchange credentials for an API access token

requestBody:
  required: true
  content:
    application/json:
      schema:
        type: object
        properties:
          grant_type:
            type: string
            enum: [password, client_credentials, refresh_token]
            description: OAuth grant type
          username:
            type: string
            description: Username (for password grant)
          password:
            type: string
            description: Password (for password grant)
          client_id:
            type: string
            description: Client ID (for client_credentials grant)
          client_secret:
            type: string
            description: Client secret (for client_credentials grant)
          refresh_token:
            type: string
            description: Refresh token (for refresh_token grant)
          scope:
            type: string
            description: Requested scopes (space-separated)

responses:
  200:
    description: Token generated successfully
    content:
      application/json:
        schema:
          type: object
          properties:
            access_token:
              type: string
              description: Short-lived access token
            token_type:
              type: string
              example: "Bearer"
            expires_in:
              type: integer
              description: Token lifetime in seconds
            refresh_token:
              type: string
              description: Long-lived refresh token
            scope:
              type: string
              description: Granted scopes

  400:
    description: Invalid grant request
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/Error'

  401:
    description: Invalid credentials
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/Error'
```

#### POST /api/v1/auth/token/refresh

Refresh an access token.

```yaml
summary: Refresh access token
description: Exchange refresh token for new access token

requestBody:
  required: true
  content:
    application/json:
      schema:
        type: object
        required: [refresh_token]
        properties:
          refresh_token:
            type: string
            description: Valid refresh token

responses:
  200:
    description: Token refreshed successfully
    content:
      application/json:
        schema:
          type: object
          properties:
            access_token:
              type: string
            expires_in:
              type: integer

  401:
    description: Invalid or expired refresh token
```

#### POST /api/v1/auth/token/revoke

Revoke an access or refresh token.

```yaml
summary: Revoke token
description: Invalidate a token

requestBody:
  required: true
  content:
    application/json:
      schema:
        type: object
        required: [token]
        properties:
          token:
            type: string
            description: Token to revoke
          token_type_hint:
            type: string
            enum: [access_token, refresh_token]
            description: Hint about token type
```

### 32.3 Session Endpoints

#### GET /api/v1/sessions

List all sessions.

```yaml
summary: List sessions
description: Get paginated list of user sessions

parameters:
  - name: limit
    in: query
    schema:
      type: integer
      default: 20
      minimum: 1
      maximum: 100
    description: Maximum number of sessions to return

  - name: before
    in: query
    schema:
      type: string
    description: Cursor for pagination (session ID)

  - name: state
    in: query
    schema:
      type: string
      enum: [active, paused, completed, all]
      default: all
    description: Filter by session state

  - name: tags
    in: query
    schema:
      type: array
      items:
        type: string
    description: Filter by tags

responses:
  200:
    description: Session list
    content:
      application/json:
        schema:
          type: object
          properties:
            sessions:
              type: array
              items:
                $ref: '#/components/schemas/Session'
            next_cursor:
              type: string
              description: Cursor for next page
            has_more:
              type: boolean
```

#### GET /api/v1/sessions/{session_id}

Get session details.

```yaml
summary: Get session
description: Retrieve detailed information about a session

parameters:
  - name: session_id
    in: path
    required: true
    schema:
      type: string
    description: Session ID

responses:
  200:
    description: Session details
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/SessionDetail'

  404:
    description: Session not found
```

#### PATCH /api/v1/sessions/{session_id}

Update session.

```yaml
summary: Update session
description: Update session properties

parameters:
  - name: session_id
    in: path
    required: true
    schema:
      type: string

requestBody:
  content:
    application/json:
      schema:
        type: object
        properties:
          title:
            type: string
          tags:
            type: array
            items:
              type: string
          model:
            $ref: '#/components/schemas/ModelConfig'
          systemPrompt:
            type: string

responses:
  200:
    description: Session updated
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/Session'

  404:
    description: Session not found
```

### 32.4 Message Endpoints

#### POST /api/v1/sessions/{session_id}/prompt

Send a prompt (synchronous).

```yaml
summary: Send prompt (sync)
description: Send a prompt and get immediate response

parameters:
  - name: session_id
    in: path
    required: true
    schema:
      type: string

requestBody:
  required: true
  content:
    application/json:
      schema:
        type: object
        properties:
          parts:
            type: array
            items:
              $ref: '#/components/schemas/Part'
          model:
            $ref: '#/components/schemas/ModelConfig'
          temperature:
            type: number
            minimum: 0
            maximum: 2
            default: 0.7
          max_tokens:
            type: integer
            minimum: 1
            maximum: 100000
          stream:
            type: boolean
            default: false

responses:
  200:
    description: Prompt response
    content:
      application/json:
        schema:
          type: object
          properties:
            message:
              $ref: '#/components/schemas/Message'
            usage:
              $ref: '#/components/schemas/Usage'

  400:
    description: Invalid request
  429:
    description: Rate limited
  500:
    description: Server error
```

#### POST /api/v1/sessions/{session_id}/prompt_async

Send a prompt (asynchronous).

```yaml
summary: Send prompt (async)
description: Send a prompt and poll for response

parameters:
  - name: session_id
    in: path
    required: true
    schema:
      type: string

requestBody:
  required: true
  content:
    application/json:
      schema:
        type: object
        properties:
          parts:
            type: array
            items:
              $ref: '#/components/schemas/Part'
          model:
            $ref: '#/components/schemas/ModelConfig'
          callback_url:
            type: string
            format: uri
            description: Webhook URL for async callback

responses:
  202:
    description: Task accepted
    content:
      application/json:
        schema:
          type: object
          properties:
            task_id:
              type: string
            status_url:
              type: string
              format: uri
```

#### GET /api/v1/tasks/{task_id}

Get async task status.

```yaml
summary: Get task status
description: Poll for async task completion

parameters:
  - name: task_id
    in: path
    required: true
    schema:
      type: string

responses:
  200:
    description: Task status
    content:
      application/json:
        schema:
          type: object
          properties:
            task_id:
              type: string
            status:
              type: string
              enum: [pending, processing, completed, failed]
            result:
              type: object
              description: Present when status is completed
            error:
              type: object
              description: Present when status is failed
            created_at:
              type: string
              format: date-time
            completed_at:
              type: string
              format: date-time
```

### 32.5 Model Management Endpoints

#### GET /api/v1/models

List available models.

```yaml
summary: List models
description: Get all available models from all providers

parameters:
  - name: provider
    in: query
    schema:
      type: string
    description: Filter by provider

  - name: capability
    in: query
    schema:
      type: string
      enum: [vision, function_calling, reasoning, streaming]
    description: Filter by capability

responses:
  200:
    description: Model list
    content:
      application/json:
        schema:
          type: object
          properties:
            models:
              type: array
              items:
                $ref: '#/components/schemas/Model'
            providers:
              type: array
              items:
                $ref: '#/components/schemas/Provider'
```

#### GET /api/v1/models/{provider}/{model_id}

Get model details.

```yaml
summary: Get model details
description: Get detailed information about a specific model

parameters:
  - name: provider
    in: path
    required: true
    schema:
      type: string
  - name: model_id
    in: path
    required: true
    schema:
      type: string

responses:
  200:
    description: Model details
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/ModelDetail'
```

### 32.6 MCP Tool Endpoints

#### GET /api/v1/tools

List available MCP tools.

```yaml
summary: List tools
description: Get all available tools from MCP servers

responses:
  200:
    description: Tool list
    content:
      application/json:
        schema:
          type: object
          properties:
            tools:
              type: array
              items:
                $ref: '#/components/schemas/Tool'
            mcp_servers:
              type: array
              items:
                $ref: '#/components/schemas/MCPServer'
```

#### GET /api/v1/mcp/status

Get MCP server status.

```yaml
summary: MCP status
description: Get status of all connected MCP servers

responses:
  200:
    description: MCP status
    content:
      application/json:
        schema:
          type: object
          properties:
            servers:
              type: array
              items:
                type: object
                properties:
                  name:
                    type: string
                  status:
                    type: string
                    enum: [connected, disconnected, error]
                  last_ping:
                    type: string
                    format: date-time
                  tools_count:
                    type: integer
```

### 32.7 Error Response Schema

```yaml
components:
  schemas:
    Error:
      type: object
      properties:
        error:
          type: object
          properties:
            code:
              type: string
              description: Error code
              enum:
                - invalid_request
                - unauthorized
                - forbidden
                - not_found
                - rate_limited
                - server_error
                - provider_error
            message:
              type: string
              description: Human-readable error message
            details:
              type: object
              description: Additional error details
            provider:
              type: string
              description: Provider that caused the error (if applicable)
            retry_after:
              type: integer
              description: Seconds until retry is possible
              minimum: 0
            request_id:
              type: string
              description: Request ID for tracing

    Usage:
      type: object
      properties:
        input_tokens:
          type: integer
        output_tokens:
          type: integer
        total_tokens:
          type: integer
        providers:
          type: object
          additionalProperties:
            type: object
            properties:
              input_tokens:
                type: integer
              output_tokens:
                type: integer
```

### 32.8 Rate Limit Headers

```yaml
headers:
  X-RateLimit-Limit:
    description: Maximum requests per window
    example: "60"

  X-RateLimit-Remaining:
    description: Remaining requests in current window
    example: "45"

  X-RateLimit-Reset:
    description: Unix timestamp when rate limit resets
    example: "1700000000"

  Retry-After:
    description: Seconds to wait before retrying (only on 429)
    example: "60"
```

---

## 🔗 33. Integration Patterns

### 33.1 Webhook Integration

OpenCode supports webhook integrations for real-time event notifications.

#### Webhook Configuration

```typescript
interface WebhookConfig {
  id: string;
  url: string;                    // HTTPS endpoint
  events: WebhookEvent[];         // Subscribed events
  secret: string;                 // HMAC signature secret
  active: boolean;
  retryPolicy: RetryPolicy;
}

type WebhookEvent = 
  | 'session.created'
  | 'session.completed'
  | 'message.created'
  | 'message.delta'
  | 'task.completed'
  | 'task.failed'
  | 'rate_limit.exceeded';

interface RetryPolicy {
  maxRetries: number;
  initialDelay: number;          // ms
  maxDelay: number;              // ms
  backoffMultiplier: number;
}
```

#### Webhook Payload

```json
{
  "event": "message.created",
  "timestamp": "2026-02-18T10:30:00Z",
  "webhook_id": "wh_abc123",
  "data": {
    "session_id": "ses_xyz789",
    "message_id": "msg_abc123",
    "role": "assistant",
    "content": "Here's your response...",
    "usage": {
      "input_tokens": 150,
      "output_tokens": 50,
      "total_tokens": 200
    }
  },
  "signature": "sha256=abc123def456..."
}
```

#### Webhook Signature Verification

```typescript
import crypto from 'crypto';

function verifyWebhookSignature(
  payload: string,
  signature: string,
  secret: string
): boolean {
  const expectedSignature = crypto
    .createHmac('sha256', secret)
    .update(payload)
    .digest('hex');
  
  return crypto.timingSafeEqual(
    Buffer.from(signature),
    Buffer.from(`sha256=${expectedSignature}`)
  );
}

// Express middleware example
app.post('/webhook', (req, res) => {
  const signature = req.headers['x-opencode-signature'] as string;
  
  if (!verifyWebhookSignature(
    JSON.stringify(req.body),
    signature,
    process.env.WEBHOOK_SECRET
  )) {
    return res.status(401).json({ error: 'Invalid signature' });
  }
  
  // Process webhook
  const { event, data } = req.body;
  
  switch (event) {
    case 'session.completed':
      await handleSessionComplete(data);
      break;
    case 'task.failed':
      await handleTaskFailure(data);
      break;
  }
  
  res.status(200).json({ received: true });
});
```

### 33.2 Database Integration

#### Session Storage (PostgreSQL)

```typescript
import { Pool } from 'pg';

class SessionStore {
  private pool: Pool;
  
  async createSession(session: Session): Promise<void> {
    await this.pool.query(
      `INSERT INTO sessions (id, title, provider, model_id, state)
       VALUES ($1, $2, $3, $4, $5)`,
      [session.id, session.title, session.model.provider, 
       session.model.model_id, session.state]
    );
  }
  
  async getSession(id: string): Promise<Session | null> {
    const result = await this.pool.query(
      `SELECT * FROM sessions WHERE id = $1`,
      [id]
    );
    return result.rows[0] || null;
  }
  
  async appendMessage(sessionId: string, message: Message): Promise<void> {
    await this.pool.query(
      `INSERT INTO messages (id, session_id, role, content)
       VALUES ($1, $2, $3, $4)`,
      [message.id, sessionId, message.role, message.content]
    );
  }
  
  async getMessages(
    sessionId: string, 
    limit: number = 50
  ): Promise<Message[]> {
    const result = await this.pool.query(
      `SELECT * FROM messages 
       WHERE session_id = $1 
       ORDER BY created_at DESC 
       LIMIT $2`,
      [sessionId, limit]
    );
    return result.rows;
  }
}
```

#### Cache Storage (Redis)

```typescript
import Redis from 'ioredis';

class CacheStore {
  private redis: Redis;
  
  // Session context cache
  async getSessionContext(sessionId: string): Promise<Context | null> {
    const key = `session:context:${sessionId}`;
    const data = await this.redis.get(key);
    return data ? JSON.parse(data) : null;
  }
  
  async setSessionContext(
    sessionId: string, 
    context: Context,
    ttl: number = 3600
  ): Promise<void> {
    const key = `session:context:${sessionId}`;
    await this.redis.setex(key, ttl, JSON.stringify(context));
  }
  
  // Rate limiting
  async checkRateLimit(
    provider: string,
    limit: number
  ): Promise<boolean> {
    const key = `ratelimit:${provider}:${Date.now()}`;
    const current = await this.redis.incr(key);
    
    if (current === 1) {
      await this.redis.expire(key, 60); // 1 minute window
    }
    
    return current <= limit;
  }
  
  // Tool result cache
  async cacheToolResult(
    toolName: string,
    input: any,
    result: any,
    ttl: number = 300
  ): Promise<void> {
    const key = `tool:${toolName}:${this.hashInput(input)}`;
    await this.redis.setex(key, ttl, JSON.stringify(result));
  }
  
  async getCachedToolResult(
    toolName: string,
    input: any
  ): Promise<any | null> {
    const key = `tool:${toolName}:${this.hashInput(input)}`;
    const data = await this.redis.get(key);
    return data ? JSON.parse(data) : null;
  }
}
```

### 33.3 External Service Integration

#### Integration with n8n

```yaml
# n8n workflow trigger configuration
# This webhook receives events from OpenCode

nodes:
  - name: "OpenCode Webhook"
    type: "n8n-nodes-base.webhook"
    parameters:
      httpMethod: "POST"
      path: "opencode-events"
      responseMode: "onReceived"
      options:
        rawBody: false

  - name: "Route by Event"
    type: "n8n-nodes-base.switch"
    parameters:
      dataType: "string"
      value1: "{{ $json.event }}"
      rules:
        conditions:
          - value1: "session.completed"
            operation: "equal"
            value2: "session.completed"

  # Handle session completion
  - name: "Process Session Complete"
    type: "n8n-nodes-base.function"
    parameters:
      functionCode: |
        const { session_id, statistics } = $json.data;
        
        // Store statistics to database
        await $db.sessions.update({
          where: { id: session_id },
          data: { 
            completed_at: new Date(),
            total_tokens: statistics.total_tokens
          }
        });
        
        return { processed: true, session_id };

# Webhook URL to configure in OpenCode:
# https://your-n8n-instance.com/webhook/opencode-events
```

#### Integration with Supabase

```typescript
import { createClient } from '@supabase/supabase-js';

class SupabaseIntegration {
  private supabase: SupabaseClient;
  
  // Store session in Supabase
  async storeSession(session: Session): Promise<void> {
    const { error } = await this.supabase
      .from('sessions')
      .insert({
        id: session.id,
        title: session.title,
        provider: session.model.provider,
        model_id: session.model.model_id,
        created_at: session.created_at,
        metadata: session.metadata
      });
    
    if (error) throw error;
  }
  
  // Query historical sessions
  async querySessions(filters: {
    provider?: string;
    startDate?: Date;
    endDate?: Date;
    limit?: number;
  }): Promise<Session[]> {
    let query = this.supabase
      .from('sessions')
      .select('*');
    
    if (filters.provider) {
      query = query.eq('provider', filters.provider);
    }
    if (filters.startDate) {
      query = query.gte('created_at', filters.startDate.toISOString());
    }
    if (filters.endDate) {
      query = query.lte('created_at', filters.endDate.toISOString());
    }
    
    const { data, error } = await query
      .order('created_at', { ascending: false })
      .limit(filters.limit || 20);
    
    if (error) throw error;
    return data;
  }
  
  // Real-time subscriptions
  subscribeToSessions(callback: (session: Session) => void): () => void {
    const subscription = this.supabase
      .channel('sessions')
      .on('postgres_changes', 
        { event: 'INSERT', schema: 'public', table: 'sessions' },
        (payload) => callback(payload.new as Session)
      )
      .subscribe();
    
    return () => this.supabase.removeChannel(subscription);
  }
}
```

### 33.4 Custom Provider Integration

```typescript
// Example: Adding a custom Ollama provider

interface OllamaConfig {
  baseURL: string;
  defaultModel: string;
  timeout: number;
}

class OllamaProvider implements Provider {
  id = 'ollama';
  name = 'Ollama (Local)';
  
  capabilities = {
    streaming: true,
    vision: false,
    functionCalling: false,
    reasoning: false,
    multiModal: false
  };
  
  limits = {
    maxContextTokens: 128000,
    maxOutputTokens: 8192
  };
  
  rateLimits = {
    requestsPerMinute: 60,
    requestsPerDay: Infinity
  };
  
  private config: OllamaConfig;
  
  constructor(config: OllamaConfig) {
    this.config = config;
  }
  
  async chat(request: ChatRequest): Promise<ChatResponse> {
    const response = await fetch(`${this.config.baseURL}/api/chat`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        model: request.model || this.config.defaultModel,
        messages: request.messages,
        stream: false,
        options: {
          temperature: request.temperature,
          num_predict: request.max_tokens,
          stop: request.stop
        }
      }),
      signal: AbortSignal.timeout(this.config.timeout)
    });
    
    if (!response.ok) {
      throw new Error(`Ollama API error: ${response.statusText}`);
    }
    
    const data = await response.json();
    
    return {
      id: data.id,
      model: data.model,
      choices: [{
        message: {
          role: 'assistant',
          content: data.message.content
        },
        finish_reason: data.done ? 'stop' : 'length'
      }],
      usage: {
        prompt_tokens: data.prompt_eval_count || 0,
        completion_tokens: data.eval_count || 0,
        total_tokens: (data.prompt_eval_count || 0) + (data.eval_count || 0)
      }
    };
  }
  
  async *stream(request: ChatRequest): AsyncGenerator<StreamChunk> {
    const response = await fetch(`${this.config.baseURL}/api/chat`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        model: request.model || this.config.defaultModel,
        messages: request.messages,
        stream: true
      }),
      signal: AbortSignal.timeout(this.config.timeout)
    });
    
    if (!response.ok) {
      throw new Error(`Ollama API error: ${response.statusText}`);
    }
    
    const reader = response.body?.getReader();
    if (!reader) throw new Error('No response body');
    
    const decoder = new TextDecoder();
    
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      
      const chunk = decoder.decode(value);
      const lines = chunk.split('\n').filter(Boolean);
      
      for (const line of lines) {
        if (line.trim() === 'done') continue;
        
        try {
          const data = JSON.parse(line);
          yield {
            id: data.id,
            delta: data.message?.content || '',
            finishReason: data.done ? 'stop' : undefined
          };
        } catch {
          // Skip parse errors
        }
      }
    }
  }
  
  async validateAuth(): Promise<boolean> {
    try {
      const response = await fetch(`${this.config.baseURL}/api/tags`);
      return response.ok;
    } catch {
      return false;
    }
  }
}

// Register custom provider
const ollamaProvider = new OllamaProvider({
  baseURL: 'http://localhost:11434',
  defaultModel: 'llama3.1',
  timeout: 60000
});

providerRegistry.register(ollamaProvider);
```

---

## 🚀 34. Deployment & Operations

### 34.1 Docker Deployment

#### Dockerfile

```dockerfile
FROM node:20-slim AS builder

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm ci --only=production

# Copy source
COPY dist/ ./dist/
COPY package*.json ./

# Production image
FROM node:20-slim AS runner

# Security: Run as non-root user
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 opencode

WORKDIR /app

# Copy built artifacts
COPY --from=builder --chown=opencode:nodejs /app/node_modules ./node_modules
COPY --from=builder --chown=opencode:nodejs /app/dist ./dist
COPY --from=builder --chown=opencode:nodejs /app/package*.json ./

# Environment configuration
ENV NODE_ENV=production
ENV PORT=8080

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD node -e "require('http').get('http://localhost:8080/health', (r) => process.exit(r.statusCode === 200 ? 0 : 1))"

# Run as non-root
USER opencode

# Start application
CMD ["node", "dist/main.js"]
```

#### Docker Compose

```yaml
version: '3.8'

services:
  opencode-server:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "53080:8080"  # Port Sovereignty: 8080 → 53080
    environment:
      - NODE_ENV=production
      - DATABASE_URL=postgresql://postgres:password@postgres:5432/opencode
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=${JWT_SECRET}
      - ENCRYPTION_KEY=${ENCRYPTION_KEY}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    volumes:
      - ./data:/app/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "node", "-e", "require('http').get('http://localhost:8080/health', (r) => process.exit(r.statusCode === 200 ? 0 : 1))"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=opencode
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
  redis_data:
```

### 34.2 Kubernetes Deployment

#### Deployment YAML

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: opencode-server
  labels:
    app: opencode-server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: opencode-server
  template:
    metadata:
      labels:
        app: opencode-server
    spec:
      containers:
        - name: opencode-server
          image: opencode/server:latest
          ports:
            - containerPort: 8080
              name: http
          env:
            - name: NODE_ENV
              value: "production"
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: opencode-secrets
                  key: database-url
            - name: REDIS_URL
              valueFrom:
                configMapKeyRef:
                  name: opencode-config
                  key: redis-url
            - name: JWT_SECRET
              valueFrom:
                secretKeyRef:
                  name: opencode-secrets
                  key: jwt-secret
          resources:
            requests:
              memory: "512Mi"
              cpu: "250m"
            limits:
              memory: "1Gi"
              cpu: "1000m"
          livenessProbe:
            httpGet:
              path: /health
              port: http
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health
              port: http
            initialDelaySeconds: 10
            periodSeconds: 5

---
apiVersion: v1
kind: Service
metadata:
  name: opencode-server
spec:
  selector:
    app: opencode-server
  ports:
    - port: 80
      targetPort: 8080
      name: http
  type: ClusterIP

---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: opencode-server-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: opencode-server
  minReplicas: 3
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
```

### 34.3 Monitoring & Observability

#### Metrics Collection

```typescript
import { Registry, Counter, Histogram, Gauge } from 'prom-client';

const registry = new Registry();

// Request metrics
const httpRequestsTotal = new Counter({
  name: 'http_requests_total',
  help: 'Total number of HTTP requests',
  labelNames: ['method', 'endpoint', 'status'],
  registers: [registry]
});

const httpRequestDuration = new Histogram({
  name: 'http_request_duration_seconds',
  help: 'HTTP request duration in seconds',
  labelNames: ['method', 'endpoint'],
  buckets: [0.1, 0.5, 1, 2, 5, 10],
  registers: [registry]
});

// Model usage metrics
const modelRequestsTotal = new Counter({
  name: 'model_requests_total',
  help: 'Total model requests',
  labelNames: ['provider', 'model', 'outcome'],
  registers: [registry]
});

const modelTokensUsed = new Counter({
  name: 'model_tokens_used_total',
  help: 'Total tokens used',
  labelNames: ['provider', 'model', 'type'],
  registers: [registry]
});

const modelLatency = new Histogram({
  name: 'model_latency_seconds',
  help: 'Model inference latency',
  labelNames: ['provider', 'model'],
  buckets: [0.5, 1, 2, 5, 10, 30, 60, 120],
  registers: [registry]
});

// Session metrics
const activeSessions = new Gauge({
  name: 'active_sessions',
  help: 'Number of active sessions',
  registers: [registry]
});

const sessionMessagesTotal = new Counter({
  name: 'session_messages_total',
  help: 'Total messages in sessions',
  labelNames: ['session_state'],
  registers: [registry]
});

// MCP metrics
const mcpToolCallsTotal = new Counter({
  name: 'mcp_tool_calls_total',
  help: 'Total MCP tool calls',
  labelNames: ['server', 'tool', 'outcome'],
  registers: [registry]
});

const mcpToolLatency = new Histogram({
  name: 'mcp_tool_latency_seconds',
  help: 'MCP tool call latency',
  labelNames: ['server', 'tool'],
  buckets: [0.1, 0.5, 1, 2, 5, 10],
  registers: [registry]
});

// Middleware for metrics collection
function metricsMiddleware(req: Request, res: Response, next: NextFunction) {
  const start = Date.now();
  
  res.on('finish', () => {
    const duration = (Date.now() - start) / 1000;
    
    httpRequestsTotal.labels(req.method, req.path, res.statusCode.toString()).inc();
    httpRequestDuration.labels(req.method, req.path).observe(duration);
  });
  
  next();
}

// Metrics endpoint
app.get('/metrics', async (req, res) => {
  res.set('Content-Type', registry.contentType);
  res.end(await registry.metrics());
});
```

#### Health Check Endpoint

```typescript
interface HealthStatus {
  status: 'healthy' | 'degraded' | 'unhealthy';
  timestamp: string;
  version: string;
  uptime: number;
  checks: {
    database: CheckResult;
    redis: CheckResult;
    providers: ProviderCheck[];
    mcp: MCPCheck[];
  };
}

interface CheckResult {
  status: 'ok' | 'error';
  latency_ms?: number;
  error?: string;
}

app.get('/health', async (req, res) => {
  const health: HealthStatus = {
    status: 'healthy',
    timestamp: new Date().toISOString(),
    version: process.env.APP_VERSION || 'unknown',
    uptime: process.uptime(),
    checks: {
      database: await checkDatabase(),
      redis: await checkRedis(),
      providers: await checkProviders(),
      mcp: await checkMCP()
    }
  };
  
  // Determine overall status
  const checks = [
    health.checks.database,
    health.checks.redis,
    ...health.checks.providers.map(p => ({ status: p.status as 'ok' | 'error' })),
    ...health.checks.mcp.map(m => ({ status: m.status as 'ok' | 'error' }))
  ];
  
  const hasError = checks.some(c => c.status === 'error');
  const hasDegraded = checks.some(c => c.status !== 'ok');
  
  health.status = hasError ? 'unhealthy' : hasDegraded ? 'degraded' : 'healthy';
  
  const statusCode = health.status === 'healthy' ? 200 : 
                    health.status === 'degraded' ? 200 : 503;
  
  res.status(statusCode).json(health);
});

async function checkDatabase(): Promise<CheckResult> {
  try {
    const start = Date.now();
    await pool.query('SELECT 1');
    return { status: 'ok', latency_ms: Date.now() - start };
  } catch (error) {
    return { status: 'error', error: error.message };
  }
}

async function checkRedis(): Promise<CheckResult> {
  try {
    const start = Date.now();
    await redis.ping();
    return { status: 'ok', latency_ms: Date.now() - start };
  } catch (error) {
    return { status: 'error', error: error.message };
  }
}
```

### 34.4 Backup & Recovery

#### Backup Strategy

```typescript
import { Pool } from 'pg';
import { createClient } from '@redis/redis';
import fs from 'fs/promises';
import path from 'path';

class BackupManager {
  private pool: Pool;
  private redis: Redis;
  private backupDir: string;
  
  constructor(backupDir: string) {
    this.backupDir = backupDir;
  }
  
  async createBackup(): Promise<BackupManifest> {
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    const backupId = `backup-${timestamp}`;
    
    console.log(`Starting backup: ${backupId}`);
    
    // Backup database
    const dbBackupPath = await this.backupDatabase(backupId);
    
    // Backup Redis cache
    const cacheBackupPath = await this.backupCache(backupId);
    
    // Backup configuration
    const configBackupPath = await this.backupConfig(backupId);
    
    // Create manifest
    const manifest: BackupManifest = {
      id: backupId,
      timestamp: new Date().toISOString(),
      files: [
        { type: 'database', path: dbBackupPath, size: await this.getFileSize(dbBackupPath) },
        { type: 'cache', path: cacheBackupPath, size: await this.getFileSize(cacheBackupPath) },
        { type: 'config', path: configBackupPath, size: await this.getFileSize(configBackupPath) }
      ],
      checksum: await this.calculateChecksum(backupId)
    };
    
    // Save manifest
    const manifestPath = path.join(this.backupDir, `${backupId}-manifest.json`);
    await fs.writeFile(manifestPath, JSON.stringify(manifest, null, 2));
    
    console.log(`Backup completed: ${backupId}`);
    
    return manifest;
  }
  
  async restoreBackup(backupId: string): Promise<void> {
    const manifestPath = path.join(this.backupDir, `${backupId}-manifest.json`);
    const manifest: BackupManifest = JSON.parse(
      await fs.readFile(manifestPath, 'utf-8')
    );
    
    console.log(`Starting restore: ${backupId}`);
    
    // Verify checksum
    const currentChecksum = await this.calculateChecksum(backupId);
    if (currentChecksum !== manifest.checksum) {
      throw new Error('Backup checksum mismatch');
    }
    
    // Restore in order: config -> cache -> database
    for (const file of manifest.files) {
      switch (file.type) {
        case 'config':
          await this.restoreConfig(file.path);
          break;
        case 'cache':
          await this.restoreCache(file.path);
          break;
        case 'database':
          await this.restoreDatabase(file.path);
          break;
      }
    }
    
    console.log(`Restore completed: ${backupId}`);
  }
  
  private async backupDatabase(backupId: string): Promise<string> {
    const outputPath = path.join(this.backupDir, `${backupId}-database.sql`);
    
    // Use pg_dump
    const { stdout } = await exec(
      `pg_dump -h ${process.env.DB_HOST} -U postgres -d opencode`,
      { env: { PGPASSWORD: process.env.DB_PASSWORD } }
    );
    
    await fs.writeFile(outputPath, stdout);
    return outputPath;
  }
  
  private async backupCache(backupId: string): Promise<string> {
    const outputPath = path.join(this.backupDir, `${backupId}-cache.json`);
    
    // Export all keys
    const keys = await this.redis.keys('*');
    const data: Record<string, string> = {};
    
    for (const key of keys) {
      const value = await this.redis.get(key);
      if (value) data[key] = value;
    }
    
    await fs.writeFile(outputPath, JSON.stringify(data));
    return outputPath;
  }
  
  private async backupConfig(backupId: string): Promise<string> {
    const outputPath = path.join(this.backupDir, `${backupId}-config.json`);
    
    const config = {
      providers: await this.getProviderConfigs(),
      mcpServers: await this.getMCPConfigs(),
      settings: await this.getSettings()
    };
    
    await fs.writeFile(outputPath, JSON.stringify(config, null, 2));
    return outputPath;
  }
}
```

---

## 🔒 35. Security Hardening

### 35.1 Authentication & Authorization

#### JWT Token Management

```typescript
import jwt from 'jsonwebtoken';
import crypto from 'crypto';

interface TokenPayload {
  userId: string;
  email: string;
  roles: string[];
  permissions: string[];
  iat: number;
  exp: number;
}

class TokenManager {
  private privateKey: string;
  private publicKey: string;
  private accessTokenTTL = 3600;      // 1 hour
  private refreshTokenTTL = 604800;  // 7 days
  
  constructor() {
    // Generate key pair on initialization
    const { publicKey, privateKey } = crypto.generateKeyPairSync('rsa', {
      modulusLength: 4096,
      publicKeyEncoding: { type: 'spki', format: 'pem' },
      privateKeyEncoding: { type: 'pkcs8', format: 'pem' }
    });
    
    this.privateKey = privateKey;
    this.publicKey = publicKey;
  }
  
  generateAccessToken(user: User): string {
    const payload = {
      userId: user.id,
      email: user.email,
      roles: user.roles,
      permissions: user.permissions
    };
    
    return jwt.sign(payload, this.privateKey, {
      algorithm: 'RS256',
      expiresIn: this.accessTokenTTL
    });
  }
  
  generateRefreshToken(user: User): string {
    const payload = {
      userId: user.id,
      type: 'refresh',
      random: crypto.randomBytes(32).toString('hex')
    };
    
    return jwt.sign(payload, this.privateKey, {
      algorithm: 'RS256',
      expiresIn: this.refreshTokenTTL
    });
  }
  
  verifyToken(token: string): TokenPayload {
    return jwt.verify(token, this.publicKey, {
      algorithms: ['RS256']
    }) as TokenPayload;
  }
  
  refreshAccessToken(refreshToken: string): string {
    const payload = this.verifyToken(refreshToken);
    
    if (payload.type !== 'refresh') {
      throw new Error('Invalid token type');
    }
    
    // Get user from database
    const user = await this.getUserById(payload.userId);
    
    return this.generateAccessToken(user);
  }
}
```

#### Role-Based Access Control (RBAC)

```typescript
interface Role {
  name: string;
  permissions: Permission[];
}

interface Permission {
  resource: string;
  actions: ('create' | 'read' | 'update' | 'delete')[];
}

const ROLES: Role[] = [
  {
    name: 'admin',
    permissions: [
      { resource: '*', actions: ['create', 'read', 'update', 'delete'] }
    ]
  },
  {
    name: 'developer',
    permissions: [
      { resource: 'session', actions: ['create', 'read', 'update'] },
      { resource: 'message', actions: ['create', 'read'] },
      { resource: 'model', actions: ['read'] },
      { resource: 'tool', actions: ['read'] }
    ]
  },
  {
    name: 'user',
    permissions: [
      { resource: 'session', actions: ['create', 'read'] },
      { resource: 'message', actions: ['create', 'read'] }
    ]
  }
];

class RBAC {
  private rolePermissions: Map<string, Permission[]>;
  
  constructor() {
    this.rolePermissions = new Map(
      ROLES.map(role => [role.name, role.permissions])
    );
  }
  
  hasPermission(
    userRoles: string[],
    resource: string,
    action: string
  ): boolean {
    for (const role of userRoles) {
      const permissions = this.rolePermissions.get(role);
      if (!permissions) continue;
      
      for (const perm of permissions) {
        // Wildcard permission
        if (perm.resource === '*') {
          return perm.actions.includes(action as any);
        }
        
        // Exact match
        if (perm.resource === resource && perm.actions.includes(action as any)) {
          return true;
        }
      }
    }
    
    return false;
  }
  
  // Middleware for Express
  requirePermission(resource: string, action: string) {
    return (req: Request, res: Response, next: NextFunction) => {
      const user = req.user as TokenPayload;
      
      if (!user) {
        return res.status(401).json({ error: 'Unauthorized' });
      }
      
      if (!this.hasPermission(user.roles, resource, action)) {
        return res.status(403).json({ error: 'Forbidden' });
      }
      
      next();
    };
  }
}
```

### 35.2 Data Encryption

#### Encryption at Rest

```typescript
import crypto from 'crypto';

class EncryptionService {
  private algorithm = 'aes-256-gcm';
  private keyLength = 32;
  private ivLength = 16;
  private saltLength = 64;
  private tagLength = 16;
  
  private masterKey: Buffer;
  
  constructor() {
    // Derive master key from environment
    const keyMaterial = process.env.ENCRYPTION_KEY!;
    this.masterKey = crypto.pbkdf2Sync(keyMaterial, 'salt', 100000, this.keyLength, 'sha512');
  }
  
  encrypt(plaintext: string): EncryptedData {
    const iv = crypto.randomBytes(this.ivLength);
    const cipher = crypto.createCipheriv(this.algorithm, this.masterKey, iv);
    
    let encrypted = cipher.update(plaintext, 'utf8', 'hex');
    encrypted += cipher.final('hex');
    
    const tag = cipher.getAuthTag();
    
    return {
      iv: iv.toString('hex'),
      data: encrypted,
      tag: tag.toString('hex')
    };
  }
  
  decrypt(encrypted: EncryptedData): string {
    const iv = Buffer.from(encrypted.iv, 'hex');
    const tag = Buffer.from(encrypted.tag, 'hex');
    const decipher = crypto.createDecipheriv(this.algorithm, this.masterKey, iv);
    
    decipher.setAuthTag(tag);
    
    let decrypted = decipher.update(encrypted.data, 'hex', 'utf8');
    decrypted += decipher.final('utf8');
    
    return decrypted;
  }
  
  // Encrypt sensitive fields in objects
  encryptFields<T extends Record<string, any>>(
    obj: T,
    fields: (keyof T)[]
  ): T {
    const encrypted = { ...obj };
    
    for (const field of fields) {
      const value = encrypted[field];
      if (value && typeof value === 'string') {
        const encryptedData = this.encrypt(value);
        (encrypted as any)[field] = JSON.stringify(encryptedData);
      }
    }
    
    return encrypted;
  }
  
  decryptFields<T extends Record<string, any>>(
    obj: T,
    fields: (keyof T)[]
  ): T {
    const decrypted = { ...obj };
    
    for (const field of fields) {
      const value = decrypted[field];
      if (value && typeof value === 'string') {
        try {
          const encryptedData = JSON.parse(value) as EncryptedData;
          (decrypted as any)[field] = this.decrypt(encryptedData);
        } catch {
          // Not encrypted, leave as-is
        }
      }
    }
    
    return decrypted;
  }
}
```

### 35.3 Audit Logging

```typescript
interface AuditLog {
  id: string;
  timestamp: Date;
  userId?: string;
  action: string;
  resource: string;
  resourceId?: string;
  changes?: Record<string, { old: any; new: any }>;
  ipAddress: string;
  userAgent: string;
  outcome: 'success' | 'failure';
  error?: string;
}

class AuditLogger {
  private logger: Console;
  private db?: Pool;
  
  async log(event: AuditLog): Promise<void> {
    const logEntry = {
      ...event,
      timestamp: event.timestamp.toISOString()
    };
    
    // Console output (for development)
    this.logger.info('[AUDIT]', JSON.stringify(logEntry));
    
    // Database storage (for production)
    if (this.db) {
      await this.db.query(
        `INSERT INTO audit_logs 
         (id, timestamp, user_id, action, resource, resource_id, changes, ip_address, user_agent, outcome, error)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
        [
          logEntry.id,
          logEntry.timestamp,
          logEntry.userId,
          logEntry.action,
          logEntry.resource,
          logEntry.resourceId,
          JSON.stringify(logEntry.changes),
          logEntry.ipAddress,
          logEntry.userAgent,
          logEntry.outcome,
          logEntry.error
        ]
      );
    }
  }
  
  // Middleware for Express
  auditMiddleware(req: Request, res: Response, next: NextFunction) {
    const start = Date.now();
    
    res.on('finish', async () => {
      const auditLog: AuditLog = {
        id: crypto.randomUUID(),
        timestamp: new Date(),
        userId: (req.user as TokenPayload)?.userId,
        action: `${req.method} ${req.path}`,
        resource: this.getResourceType(req.path),
        resourceId: this.getResourceId(req.path),
        ipAddress: req.ip,
        userAgent: req.headers['user-agent'] || 'unknown',
        outcome: res.statusCode < 400 ? 'success' : 'failure',
        error: res.statusCode >= 400 ? res.statusMessage : undefined
      };
      
      await this.log(auditLog);
    });
    
    next();
  }
}
```

---

## 📊 Anhang: Erweiterte Config-Referenz

### Vollständige Config-Schema

```json
{
  "$schema": "https://opencode.ai/config/v1/schema.json",
  "type": "object",
  "properties": {
    "version": {
      "type": "string",
      "description": "Config version"
    },
    "model": {
      "type": "string",
      "description": "Default model identifier"
    },
    "small_model": {
      "type": "string",
      "description": "Small model for quick tasks"
    },
    "default_agent": {
      "type": "string",
      "description": "Default agent name"
    },
    "theme": {
      "type": "string",
      "enum": ["light", "dark", "system"],
      "default": "system"
    },
    "autoupdate": {
      "type": "boolean",
      "default": true
    },
    "plugin": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "description": "Enabled plugins"
    },
    "provider": {
      "type": "object",
      "additionalProperties": {
        "type": "object",
        "properties": {
          "npm": {
            "type": "string",
            "description": "NPM package for provider"
          },
          "name": {
            "type": "string",
            "description": "Display name"
          },
          "options": {
            "type": "object",
            "properties": {
              "baseURL": {
                "type": "string",
                "format": "uri"
              },
              "timeout": {
                "type": "number",
                "minimum": 1000,
                "maximum": 300000
              },
              "apiKey": {
                "type": "string"
              },
              "headers": {
                "type": "object",
                "additionalProperties": {
                  "type": "string"
                }
              },
              "maxConcurrentRequests": {
                "type": "number",
                "minimum": 1,
                "maximum": 100
              }
            }
          },
          "models": {
            "type": "object",
            "additionalProperties": {
              "type": "object",
              "properties": {
                "id": {
                  "type": "string"
                },
                "name": {
                  "type": "string"
                },
                "limit": {
                  "type": "object",
                  "properties": {
                    "context": {
                      "type": "number"
                    },
                    "output": {
                      "type": "number"
                    }
                  }
                },
                "modalities": {
                  "type": "object",
                  "properties": {
                    "input": {
                      "type": "array",
                      "items": {
                        "type": "string"
                      }
                    },
                    "output": {
                      "type": "array",
                      "items": {
                        "type": "string"
                      }
                    }
                  }
                },
                "variants": {
                  "type": "object",
                  "additionalProperties": {
                    "type": "object"
                  }
                },
                "supports": {
                  "type": "object",
                  "properties": {
                    "vision": {
                      "type": "boolean"
                    },
                    "function_calling": {
                      "type": "boolean"
                    },
                    "streaming": {
                      "type": "boolean"
                    }
                  }
                }
              }
            }
          }
        }
      }
    },
    "mcp": {
      "type": "object",
      "additionalProperties": {
        "type": "object",
        "properties": {
          "type": {
            "type": "string",
            "enum": ["local", "remote"]
          },
          "command": {
            "type": "array",
            "items": {
              "type": "string"
            }
          },
          "url": {
            "type": "string",
            "format": "uri"
          },
          "enabled": {
            "type": "boolean"
          },
          "environment": {
            "type": "object",
            "additionalProperties": {
              "type": "string"
            }
          },
          "keepAlive": {
            "type": "boolean"
          }
        }
      }
    },
    "agent": {
      "type": "object",
      "additionalProperties": {
        "type": "object",
        "properties": {
          "model": {
            "type": "object",
            "properties": {
              "provider": {
                "type": "string"
              },
              "id": {
                "type": "string"
              }
            }
          },
          "systemPrompt": {
            "type": "string"
          }
        }
      }
    },
    "cache": {
      "type": "object",
      "properties": {
        "enabled": {
          "type": "boolean"
        },
        "ttl": {
          "type": "number",
          "minimum": 0
        },
        "maxSize": {
          "type": "number",
          "minimum": 1
        },
        "similarityThreshold": {
          "type": "number",
          "minimum": 0,
          "maximum": 1
        }
      }
    },
    "retry": {
      "type": "object",
      "properties": {
        "maxRetries": {
          "type": "number",
          "minimum": 0,
          "maximum": 10
        },
        "baseDelay": {
          "type": "number",
          "minimum": 100
        },
        "maxDelay": {
          "type": "number",
          "minimum": 1000
        },
        "backoffMultiplier": {
          "type": "number",
          "minimum": 1
        }
      }
    },
    "rateLimit": {
      "type": "object",
      "properties": {
        "enabled": {
          "type": "boolean"
        },
        "windowMs": {
          "type": "number"
        },
        "maxRequests": {
          "type": "number"
        }
      }
    },
    "security": {
      "type": "object",
      "properties": {
        "encryption": {
          "type": "string",
          "enum": ["AES-256-GCM", "AES-256-CBC"]
        },
        "auditLogging": {
          "type": "boolean"
        },
        "sessionTimeout": {
          "type": "number"
        },
        "maxSessionAge": {
          "type": "number"
        }
      }
    }
  }
}
```

---

**Version:** 1.0  
**Stand:** 2026-02-18  
**Status:** EXPANDED ✅  
**Total Lines:** 5000+

---

**Updated:** 2026-02-17  
**Status:** COMPLETE ✅  
**All Agents:** 12 Agents fully documented with models, costs, and use cases
