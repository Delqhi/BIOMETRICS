# OpenCode Rules & Best Practices

**Projekt:** BIOMETRICS  
**Version:** 1.0  
**Letzte Aktualisierung:** 2026-02-20  
**Status:** ACTIVE  

---

## Inhaltsverzeichnis

1. [OpenCode Installation](#1-opencode-installation)
2. [Provider Setup](#2-provider-setup)
3. [Model Selection](#3-model-selection)
4. [Agent Usage](#4-agent-usage)
5. [Task Delegation](#5-task-delegation)
6. [Configuration](#6-configuration)
7. [Best Practices](#7-best-practices)
8. [Troubleshooting](#8-troubleshooting)

---

## 1. OpenCode Installation

### 1.1 Grundinstallation

OpenCode ist der primäre AI-Code-Assistant für dieses Projekt. Die Installation erfolgt via Homebrew:

```bash
# Installation via Homebrew
brew install opencode

# Oder via npm (Fallback)
npm install -g @opencodeai/cli
```

### 1.2 Versionsanforderungen

**Mindestversion:** 1.0.150

Die folgenden Features erfordern Version 1.0.150 oder höher:
- Erweiterte Model-Konfiguration
- MCP-Server Support
- Verbesserte Task-Delegation
- Streaming Support

```bash
# Version verifizieren
opencode --version

# Sollte ausgeben: opencode version 1.0.150+
```

### 1.3 Initialisierung

Nach der Installation muss OpenCode initialisiert werden:

```bash
# Initialisierung starten
opencode init

# Authentifizierung prüfen
opencode auth status

# Verfügbare Modelle anzeigen
opencode models
```

### 1.4 Verifizierungsbefehle

```bash
# Health-Check
opencode doctor

# Model-Status
opencode models status

# Konfiguration anzeigen
opencode config show
```

---

## 2. Provider Setup

### 2.1 NVIDIA NIM (Empfohlen - FREE)

NVIDIA NIM bietet kostenlose High-Performance-Modelle. Konfiguration in `~/.config/opencode/opencode.json`:

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
          "limit": { "context": 262144, "output": 32768 }
        },
        "qwen2.5-coder-32b": {
          "id": "Qwen/Qwen2.5-Coder-32B-Instruct",
          "limit": { "context": 131072, "output": 8192 }
        },
        "qwen2.5-coder-7b": {
          "id": "Qwen/Qwen2.5-Coder-7B-Instruct",
          "limit": { "context": 131072, "output": 8192 }
        }
      }
    }
  }
}
```

**Umgebungsvariable:**
```bash
export NVIDIA_API_KEY="nvapi-ihre-api-key-hier"
```

### 2.2 OpenCode ZEN (100% FREE)

OpenCode ZEN bietet unzensierte Modelle ohne API-Kosten:

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
        "zen-big-pickle": {
          "name": "Big Pickle (UNCENSORED)",
          "limit": { "context": 200000, "output": 128000 }
        },
        "zen-uncensored": {
          "name": "Uncensored",
          "limit": { "context": 200000, "output": 128000 }
        },
        "zen-code": {
          "name": "Code (ZEN)",
          "limit": { "context": 200000, "output": 128000 }
        }
      }
    }
  }
}
```

### 2.3 Google Antigravity (Gemini)

Für Gemini-Modelle via Antigravity-Plugin:

```json
{
  "provider": {
    "google": {
      "npm": "@ai-sdk/google",
      "models": {
        "antigravity-gemini-3-flash": {
          "id": "gemini-3-flash-preview",
          "name": "Gemini 3 Flash (Antigravity)",
          "limit": { "context": 1048576, "output": 65536 }
        },
        "antigravity-gemini-3-pro": {
          "id": "gemini-3-pro-preview",
          "name": "Gemini 3 Pro (Antigravity)",
          "limit": { "context": 2097152, "output": 65536 }
        }
      }
    }
  }
}
```

**Authentifizierung:**
```bash
opencode auth login --provider google
# Private Gmail verwenden (nicht Workspace!)
```

### 2.4 Streamlake

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
          "limit": { "context": 2000000, "output": 128000 }
        }
      }
    }
  }
}
```

### 2.5 XiaoMi MIMO

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
          "limit": { "context": 1000000, "output": 65536 }
        },
        "mimo-v2-turbo": {
          "name": "MIMO v2 Turbo",
          "limit": { "context": 1500000, "output": 100000 }
        }
      }
    }
  }
}
```

### 2.6 Komplette Provider-Konfiguration

Die vollständige `opencode.json` sollte alle Provider enthalten:

```json
{
  "provider": {
    "nvidia-nim": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "NVIDIA NIM",
      "options": {
        "baseURL": "https://integrate.api.nvidia.com/v1",
        "timeout": 120000
      },
      "models": {
        "qwen-3.5-397b": {
          "id": "qwen/qwen3.5-397b-a17b",
          "limit": { "context": 262144, "output": 32768 }
        }
      }
    },
    "opencode-zen": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "OpenCode ZEN",
      "options": {
        "baseURL": "https://api.opencode.ai/v1"
      },
      "models": {
        "kimi-k2.5-free": {
          "name": "Kimi K2.5 (ZEN)",
          "limit": { "context": 1048576, "output": 65536 }
        },
        "minimax-m2.5-free": {
          "name": "MiniMax M2.5 (ZEN)",
          "limit": { "context": 200000, "output": 8192 }
        }
      }
    }
  },
  "defaults": {
    "model": {
      "primary": "nvidia-nim/qwen-3.5-397b",
      "fallbacks": [
        "opencode-zen/minimax-m2.5-free",
        "opencode-zen/kimi-k2.5-free"
      ]
    }
  }
}
```

---

## 3. Model Selection

### 3.1 Model-Übersicht

| Modell | Provider | Context | Output | Use Case |
|--------|----------|---------|--------|----------|
| **Qwen 3.5 397B** | NVIDIA NIM | 262K | 32K | Code-Implementation |
| **Kimi K2.5** | OpenCode ZEN | 1M | 64K | Deep Analysis |
| **MiniMax M2.5** | OpenCode ZEN | 200K | 8K | Recherche, MD-Erstellung |
| **Gemini 3 Pro** | Google | 2M | 65K | Multimodal |
| **KAT Coder Pro** | Streamlake | 2M | 128K | Code-Generation |

### 3.2 Wann welches Modell

**Qwen 3.5 397B (NVIDIA NIM):**
- Code-Implementation und Refactoring
- Komplexe Algorithmen
- Architektur-Entscheidungen
- Bug-Fixing in kritischen Systemen

```bash
# Qwen 3.5 verwenden
opencode --model nvidia-nim/qwen-3.5-397b "Implementiere OAuth2 Flow"
```

**Kimi K2.5 (OpenCode ZEN):**
- Deep Research und Analysen
- Komplexe Datenanalysen
- Multimodale Aufgaben (Bilder + Text)
- Lange Kontext-Verarbeitung

```bash
# Kimi K2.5 verwenden
opencode --model opencode-zen/kimi-k2.5-free "Analysiere diesen Datensatz"
```

**MiniMax M2.5 (OpenCode ZEN):**
- Schnelle Recherche
- Dokumentation erstellen
- Parallele Aufgaben (bis zu 10 parallel)
- Text-Generation

```bash
# MiniMax für parallele Recherche
opencode --model opencode-zen/minimax-m2.5-free "Recherchiere Best Practices"
```

### 3.3 Model Limits

**Wichtig:** Context-Limits nicht überschreiten!

```javascript
// Context-Berechnung Beispiel
const maxTokens = 262144;  // Qwen 3.5 Context
const systemPrompt = 2000;  // System-Prompt
const availableForInput = maxTokens - systemPrompt;  // 260K

// Bei langen Dokumenten: Chunking verwenden
// Max 260K Tokens pro Request
```

### 3.4 Fallback Chain

Immer Fallback-Modelle konfigurieren:

```json
{
  "defaults": {
    "model": {
      "primary": "nvidia-nim/qwen-3.5-397b",
      "fallbacks": [
        "opencode-zen/kimi-k2.5-free",
        "opencode-zen/minimax-m2.5-free"
      ]
    }
  }
}
```

---

## 4. Agent Usage

### 4.1 Verfügbare Agenten

| Agent | Modell | Primary Use |
|-------|--------|-------------|
| **sisyphus** | Qwen 3.5 397B | Code-Implementation |
| **prometheus** | Qwen 3.5 397B | Planung |
| **oracle** | Qwen 3.5 397B | Architektur-Review |
| **atlas** | Kimi K2.5 | Heavy Lifting |
| **librarian** | MiniMax M2.5 | Recherche |
| **explore** | MiniMax M2.5 | Code Discovery |

### 4.2 Sisyphus (Main Coder)

Sisyphus ist der primäre Code-Implementierungs-Agent.

```typescript
// Typischer Einsatz
const agent = await call_omo_agent({
  subagent_type: "sisyphus",
  prompt: `
Implementiere einen OAuth2 Authorization Code Flow mit PKCE.

Anforderungen:
- Authorization Code mit PKCE
- Token-Refresh
- Sichere Session-Verwaltung
- TypeScript, strict mode

Output: Vollständige Implementierung mit Tests
  `,
  run_in_background: true
});
```

### 4.3 Prometheus (Planner)

Prometheus erstellt detaillierte Implementierungspläne.

```typescript
const plan = await call_omo_agent({
  subagent_type: "prometheus",
  prompt: `
Erstelle einen 10-Phasen Master-Plan für:
"Biometrische Authentifizierung mit Face ID"

Folgende Aspekte abdecken:
1. Architektur-Design
2. Security-Anforderungen
3. Datenbank-Schema
4. API-Endpoints
5. Frontend-Integration
6. Testing-Strategie
7. Deployment
8. Monitoring
9. Dokumentation
10. Wartung

Pro Phase: Beschreibung, Deliverables, Abhängigkeiten
  `,
  run_in_background: true
});
```

### 4.4 Oracle (Architecture)

Oracle analysiert und bewertet Architektur-Entscheidungen.

```typescript
const review = await call_omo_agent({
  subagent_type: "oracle",
  prompt: `
Review folgende Architektur-Entscheidung:

"Verwendung von JWT für Session-Management inkl. Refresh-Token"

Bitte analysieren:
1. Security-Aspekte
2. Performance-Implications
3. Skalierbarkeit
4. Wartbarkeit
5. Alternativen

Mit konkreten Verbesserungsvorschlägen
  `,
  run_in_background: true
});
```

### 4.5 Atlas (Heavy Lifting)

Atlas übernimmt rechenintensive oder komplexe Aufgaben.

```typescript
const heavyTask = await call_omo_agent({
  subagent_type: "atlas",
  prompt: `
Implementiere ein komplettes YOLO-Training-Pipeline:

1. Data Preprocessing
2. Model Training mit Validation
3. Model Export
4. Performance Benchmarking

Verwende:
- ultralytics YOLO
- PyTorch
- Redis Cache
- Kubernetes-ready
  `,
  run_in_background: true
});
```

### 4.6 Librarian (Documentation)

Librarian erstellt und pflegt Dokumentation.

```typescript
const docs = await call_omo_agent({
  subagent_type: "librarian",
  prompt: `
Erstelle eine umfassende API-Dokumentation für:
"Biometrics Authentication API"

Enthalten:
1. Overview
2. Authentication
3. Endpoints (alle CRUD-Operationen)
4. Request/Response Examples
5. Error Codes
6. Rate Limiting
7. Best Practices

Format: Markdown mit OpenAPI-Referenzen
  `,
  run_in_background: true
});
```

### 4.7 Explore (Code Discovery)

Explore findet und analysiert Code in der Codebase.

```typescript
const discovery = await call_omo_agent({
  subagent_type: "explore",
  prompt: `
Finde alle Implementierungen von:
"Face Recognition Service"

In:
- /Users/jeremy/dev/BIOMETRICS/

Analysiere:
1. Wo ist der Service implementiert?
2. Welche Interfaces werden verwendet?
3. Welche Dependencies gibt es?
4. Wo wird der Service verwendet?

Output: Vollständiger Report mit Dateipfaden
  `,
  run_in_background: true
});
```

---

## 5. Task Delegation

### 5.1 delegate_task() Syntax

Die Hauptmethode für Agenten-Delegation:

```typescript
// Grundstruktur
await delegate_task({
  category: "solver",        // Agent-Kategorie
  subagent: "sisyphus",     // Spezifischer Agent
  prompt: "...",            // Aufgabenbeschreibung
  run_in_background: true,  // IMMER true!
  model: "qwen-3.5-397b"   // Optional: Modell überschreiben
});
```

### 5.2 run_in_background=true (PFLICHT!)

**WICHTIG:** `run_in_background` muss IMMER `true` sein!

```typescript
// ✅ RICHTIG
await delegate_task({
  prompt: "Implementiere Feature X",
  run_in_background: true  // Parallel ausführen
});

// ❌ VERBOTEN
await delegate_task({
  prompt: "Implementiere Feature X",
  run_in_background: false  // NICHT erlaubt!
});
```

### 5.3 Category Selection

| Category | Agent | Use Case |
|----------|-------|----------|
| solver | sisyphus | Code-Implementation |
| planner | prometheus | Planung |
| architecture | oracle | Architektur-Review |
| heavy | atlas | Komplexe Aufgaben |
| docs | librarian | Dokumentation |
| explore | explore | Code-Suche |

### 5.4 Prompt Structure (6 Sections)

Jeder Prompt sollte 6 Abschnitte haben:

```typescript
const prompt = `
## AUFGABE
[Beschreibung der Hauptaufgabe]

## ANFORDERUNGEN
- Anforderung 1
- Anforderung 2
- Anforderung 3

## TECHNISCHE DETAILS
- Tech Stack: TypeScript, Node.js
- Framework: Express.js
- Database: PostgreSQL

## ERWARTETES OUTPUT
- [ ] Deliverable 1
- [ ] Deliverable 2
- [ ] Deliverable 3

## QUALITÄTSKRITERIEN
- Tests erforderlich
- TypeScript strict mode
- Dokumentation aktualisieren

## TIMELINE
- Geschätzte Zeit: 2 Stunden
- Meilensteine: 3
`;
```

### 5.5 Parallele Delegation

Mehrere Agenten parallel starten:

```typescript
// Parallele Aufgaben
const [task1, task2, task3] = await Promise.all([
  delegate_task({ category: "explore", prompt: "...", run_in_background: true }),
  delegate_task({ category: "librarian", prompt: "...", run_in_background: true }),
  delegate_task({ category: "solver", prompt: "...", run_in_background: true })
]);

// Ergebnisse sammeln
const results = await Promise.all([
  background_output({ task_id: task1.task_id }),
  background_output({ task_id: task2.task_id }),
  background_output({ task_id: task3.task_id })
]);
```

### 5.6 Maximale Parallelität

**Regel:** Maximal 3 Agenten parallel pro Hauptaufgabe!

```typescript
// ✅ Max 3 parallel
const tasks = await Promise.all([
  delegate_task({ category: "explore", prompt: "...", run_in_background: true }),
  delegate_task({ category: "librarian", prompt: "...", run_in_background: true }),
  delegate_task({ category: "solver", prompt: "...", run_in_background: true })
]);

// ❌ Mehr als 3 - VERBOTEN
// Dies führt zu Rate-Limit-Problemen!
```

---

## 6. Configuration

### 6.1 Hauptkonfigurationsdatei

Location: `~/.config/opencode/opencode.json`

```json
{
  "version": "1.0",
  "provider": {
    "nvidia-nim": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "NVIDIA NIM",
      "options": {
        "baseURL": "https://integrate.api.nvidia.com/v1",
        "timeout": 120000
      },
      "models": {
        "qwen-3.5-397b": {
          "id": "qwen/qwen3.5-397b-a17b",
          "limit": { "context": 262144, "output": 32768 }
        }
      }
    },
    "opencode-zen": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "OpenCode ZEN",
      "options": {
        "baseURL": "https://api.opencode.ai/v1"
      },
      "models": {
        "minimax-m2.5-free": {
          "name": "MiniMax M2.5",
          "limit": { "context": 200000, "output": 8192 }
        },
        "kimi-k2.5-free": {
          "name": "Kimi K2.5",
          "limit": { "context": 1048576, "output": 65536 }
        }
      }
    }
  },
  "defaults": {
    "model": {
      "primary": "nvidia-nim/qwen-3.5-397b",
      "fallbacks": [
        "opencode-zen/kimi-k2.5-free",
        "opencode-zen/minimax-m2.5-free"
      ]
    },
    "temperature": 0.7,
    "max_tokens": 8192
  },
  "mcp": {
    "enabled": true
  },
  "plugins": [
    "opencode-antigravity-auth"
  ]
}
```

### 6.2 Environment-Variablen

```bash
# NVIDIA API
export NVIDIA_API_KEY="nvapi-..."

# OpenCode Auth
export OPENCODE_AUTH_TOKEN="..."

# Proxy (falls benötigt)
export HTTP_PROXY="http://localhost:8080"
export HTTPS_PROXY="http://localhost:8080"
```

### 6.3 Model Presets

Vordefinierte Model-Konfigurationen:

```json
{
  "presets": {
    "code": {
      "model": "nvidia-nim/qwen-3.5-397b",
      "temperature": 0.3,
      "max_tokens": 32768
    },
    "creative": {
      "model": "opencode-zen/kimi-k2.5-free",
      "temperature": 0.9,
      "max_tokens": 16384
    },
    "fast": {
      "model": "opencode-zen/minimax-m2.5-free",
      "temperature": 0.5,
      "max_tokens": 4096
    }
  }
}
```

### 6.4 Plugin System

Plugins erweitern OpenCode-Funktionalität:

```bash
# Plugin installieren
opencode plugin install opencode-antigravity-auth

# Plugin-Liste anzeigen
opencode plugin list

# Plugin entfernen
opencode plugin remove opencode-antigravity-auth
```

### 6.5 Antigravity Auth Setup

```bash
# Google OAuth starten
opencode auth login --provider google

# Status prüfen
opencode auth status

# Token refresh
opencode auth refresh
```

---

## 7. Best Practices

### 7.1 Parallel Execution Rules

**Maximale Parallelität:** 3 Agenten gleichzeitig

```typescript
// RICHTIG: Max 3 parallel
const results = await Promise.all([
  delegate_task({ category: "explore", run_in_background: true }),
  delegate_task({ category: "librarian", run_in_background: true }),
  delegate_task({ category: "solver", run_in_background: true })
]);

// ❌ FALSCH: Mehr als 3
// Dies führt zu:
// - Rate Limiting
// - Provider-Fehlern
// - Verschlechterter Qualität
```

### 7.2 Model Assignment Rules

| Aufgabe | Modell | Warum |
|---------|--------|-------|
| Code schreiben | Qwen 3.5 397B | Beste Code-Qualität |
| Recherche | MiniMax M2.5 | Schnell, parallel möglich |
| Analyse | Kimi K2.5 | Großer Context |
| Dokumentation | MiniMax M2.5 | Schnell |

```typescript
// Code-Aufgabe → Qwen
await delegate_task({
  category: "solver",
  model: "nvidia-nim/qwen-3.5-397b",
  prompt: "..."
});

// Recherche → MiniMax
await delegate_task({
  category: "explore",
  model: "opencode-zen/minimax-m2.5-free",
  prompt: "..."
});
```

### 7.3 Timeout Configuration

**Qwen 3.5 erfordert 120s Timeout!**

```json
{
  "provider": {
    "nvidia-nim": {
      "options": {
        "baseURL": "https://integrate.api.nvidia.com/v1",
        "timeout": 120000  // 120 Sekunden - PFLICHT!
      }
    }
  }
}
```

```typescript
// Timeout in delegate_task
await delegate_task({
  prompt: "...",
  timeout: 120000  // 120 Sekunden
});
```

### 7.4 Error Handling

Immer mit Fallback arbeiten:

```typescript
try {
  const result = await delegate_task({
    category: "solver",
    prompt: "Implementiere X"
  });
} catch (error) {
  // Fallback auf anderes Modell
  const fallback = await delegate_task({
    category: "solver",
    model: "opencode-zen/kimi-k2.5-free",
    prompt: "Implementiere X"
  });
}
```

### 7.5 Session Management

```typescript
// Neue Session erstellen
const session = await opencode.session.create({
  title: "Biometrics Feature Implementation"
});

// Session fortsetzen
await opencode.session.continue(session.id);

// Session beenden
await opencode.session.end(session.id);
```

### 7.6 Context Management

```typescript
// Kontext komprimieren wenn nötig
const compressedContext = await compressContext({
  messages: chatHistory,
  maxTokens: 100000
});

// Chunking für große Dokumente
const chunks = await chunkDocument({
  content: largeDocument,
  maxChunkSize: 50000
});
```

---

## 8. Troubleshooting

### 8.1 Common Errors

**Error: "Model not found"**
```bash
# Lösung: Modell-Liste aktualisieren
opencode models update

# Prüfen ob Modell verfügbar
opencode models | grep qwen
```

**Error: "Rate limit exceeded"**
```bash
# Lösung: 60 Sekunden warten
# Dann mit Fallback erneut versuchen
```

**Error: "Authentication failed"**
```bash
# Lösung: Authentifizierung erneuern
opencode auth login --provider nvidia
opencode auth status
```

### 8.2 Rate Limits (NVIDIA)

**Limit:** 40 Requests pro Minute (Free Tier)

```typescript
// Rate Limit Handling
const RATE_LIMIT_DELAY = 60000; // 1 Minute

async function withRateLimit(fn) {
  try {
    return await fn();
  } catch (error) {
    if (error.status === 429) {
      console.log("Rate limit reached, waiting...");
      await sleep(RATE_LIMIT_DELAY);
      return withRateLimit(fn); // Retry
    }
    throw error;
  }
}
```

### 8.3 HTTP 429 Handling

```typescript
// 429 Error Handling mit Exponential Backoff
async function requestWithBackoff(url, options, retries = 3) {
  for (let i = 0; i < retries; i++) {
    const response = await fetch(url, options);
    
    if (response.status === 429) {
      const waitTime = Math.pow(2, i) * 1000;
      console.log(`Rate limited, waiting ${waitTime}ms...`);
      await sleep(waitTime);
      continue;
    }
    
    return response;
  }
  throw new Error("Max retries exceeded");
}
```

### 8.4 Provider Fallbacks

Immer Fallback-Provider konfigurieren:

```typescript
const providers = [
  { name: "nvidia-nim", model: "qwen-3.5-397b" },
  { name: "opencode-zen", model: "kimi-k2.5-free" },
  { name: "opencode-zen", model: "minimax-m2.5-free" }
];

async function requestWithFallback(prompt) {
  for (const provider of providers) {
    try {
      return await delegate_task({
        model: `${provider.name}/${provider.model}`,
        prompt: prompt
      });
    } catch (error) {
      console.log(`Provider ${provider.name} failed, trying next...`);
      continue;
    }
  }
  throw new Error("All providers failed");
}
```

### 8.5 Debugging

```bash
# Debug-Modus aktivieren
opencode --debug

# Logs anzeigen
opencode logs

# Konfiguration prüfen
opencode config validate

# Doctor (Auto-Repair)
opencode doctor --fix
```

### 8.6 Performance Optimization

```typescript
// Connection Pooling
const pool = new ConnectionPool({
  maxConnections: 10,
  timeout: 30000
});

// Caching
const cache = new LRUCache({
  maxSize: 1000,
  ttl: 3600000 // 1 Stunde
});
```

### 8.7 Health Checks

```bash
# Alle Provider prüfen
opencode doctor

# Spezifischen Provider prüfen
opencode models status --provider nvidia-nim

# Netzwerk-Check
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
     https://integrate.api.nvidia.com/v1/models
```

### 8.8 Quick Fixes

| Problem | Lösung |
|---------|--------|
| Langsames Modell | MiniMax für schnelle Tasks verwenden |
| Rate Limits | Fallback-Provider konfigurieren |
| Auth-Fehler | `opencode auth refresh` |
| Modell nicht gefunden | `opencode models update` |
| Konfigurationsfehler | `opencode config validate` |

---

## Appendix A: Command Reference

### Grundbefehle

```bash
# Hilfe
opencode --help

# Version
opencode --version

# Modelle anzeigen
opencode models

# Aktuelles Modell
opencode model current

# Auth Status
opencode auth status

# Konfiguration
opencode config show
```

### Erweiterte Befehle

```bash
# Mit spezifischem Modell starten
opencode --model nvidia-nim/qwen-3.5-397b

# Mit Temperature
opencode --temperature 0.5

# Mit System-Prompt
opencode --system "Du bist ein erfahrener TypeScript Entwickler"

# Session erstellen
opencode session new "Mein Projekt"

# Session fortsetzen
opencode session continue <session-id>
```

---

## Appendix B: Model Comparison

### Code Quality (Qwen 3.5 vs Others)

| Kriterium | Qwen 3.5 397B | Kimi K2.5 | MiniMax M2.5 |
|-----------|---------------|-----------|--------------|
| Code-Syntax | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| TypeScript | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| Architektur | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐ |
| Dokumentation | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| Geschwindigkeit | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

### Context Handling

| Modell | Max Context | Optimal für |
|--------|-------------|-------------|
| Qwen 3.5 397B | 262K | Mid-size Projekte |
| Kimi K2.5 | 1M | Große Codebases |
| MiniMax M2.5 | 200K | Schnelle Lookups |

---

## Appendix C: Security Best Practices

### API Keys

```bash
# NIEMALS in Config-Files speichern!
# Environment-Variablen verwenden

# ~/.zshrc oder ~/.bashrc
export NVIDIA_API_KEY="nvapi-..."
```

### Rate Limiting

```typescript
// Client-seitiges Rate Limiting
const rateLimiter = new RateLimiter({
  maxRequests: 40,
  windowMs: 60000 // 1 Minute
});

await rateLimiter.schedule(() => makeRequest());
```

---

## Appendix D: Integration Examples

### Mit n8n

```typescript
// n8n Webhook für OpenCode
const n8nWebhook = async (req, res) => {
  const { prompt, model } = req.body;
  
  const result = await delegate_task({
    category: "solver",
    model: model || "nvidia-nim/qwen-3.5-397b",
    prompt: prompt
  });
  
  res.json({ result });
};
```

### Mit GitHub Actions

```yaml
# .github/workflows/opencode.yml
name: OpenCode Review

on: [pull_request]

jobs:
  review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run OpenCode Review
        run: |
          opencode --model nvidia-nim/qwen-3.5-397b \
            "Review this PR for security issues"
```

---

## Appendix E: Fehlercodes

| Code | Bedeutung | Lösung |
|------|-----------|--------|
| 401 | Unauthorized | Auth erneuern |
| 403 | Forbidden | Rechte prüfen |
| 404 | Not Found | Modell-Name prüfen |
| 429 | Rate Limited | Warten + Fallback |
| 500 | Server Error | Retry mit Fallback |
| 503 | Service Unavailable | Später erneut |

---

## Appendix F: Checkliste

Vor jedem OpenCode-Einsatz:

- [ ] Modell ausgewählt (Qwen für Code, MiniMax für Recherche)
- [ ] Timeout auf 120s für Qwen gesetzt
- [ ] Fallback-Provider konfiguriert
- [ ] max. 3 parallele Agenten
- [ ] run_in_background=true
- [ ] Prompt strukturiert (6 Sections)
- [ ] Error Handling implementiert

---

## Appendix G: Glossar

| Begriff | Bedeutung |
|---------|-----------|
| **Provider** | API-Anbieter (NVIDIA, OpenCode ZEN, etc.) |
| **Model** | KI-Modell (Qwen, Kimi, MiniMax) |
| **Agent** | Spezialisierter KI-Assistent |
| **Delegate** | Aufgabe an Agent delegieren |
| **Fallback** | Reserve-Provider bei Fehlern |
| **Context** | Token-Limit für Input |
| **Output** | Token-Limit für Response |

---

## Version History

| Version | Datum | Änderungen |
|---------|-------|------------|
| 1.0 | 2026-02-20 | Initiale Version |

---

**ENDE DES DOKUMENTS**
