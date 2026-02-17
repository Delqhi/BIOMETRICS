# OPENCODE.md — Complete Configuration Guide

**Status:** ACTIVE  
**Version:** 1.0  
**Stand:** 2026-02-17  
**Purpose:** Vollständige Konfiguration für OpenCode + OpenClaw mit NVIDIA NIM + Google Antigravity

---

## 🚀 1. Quick Start

### Voraussetzungen
- Node.js >= 20
- pnpm (NICHT npm!)
- Docker (für MCP Server)

### Installation
```bash
# OpenCode installieren
npm install -g opencode

# OpenClaw installieren (pnpm!)
pnpm install -g @anthropic/openclaw
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

**Updated:** 2026-02-17  
**Status:** COMPLETE ✅  
**All Agents:** 12 Agents fully documented with models, costs, and use cases
