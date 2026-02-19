# 🔌 PROVIDER SETUP GUIDE - COMPLETE CONFIGURATION

**Version:** 1.0 "ALL PROVIDERS"  
**Source:** Consolidated from `~/.config/opencode/opencode.json`, `~/.opencode/providers/`

---

## 🎯 PROVIDER OVERVIEW

| Provider | Type | Status | Use Case |
|----------|------|--------|----------|
| **NVIDIA NIM** | Primary | ✅ Active | Code generation (Qwen 3.5 397B) |
| **OpenCode ZEN** | Fallback | ✅ Active | FREE uncensored models |
| **Google Antigravity** | Optional | ⚠️ OAuth | Gemini models |
| **Moonshot AI** | Optional | ✅ Active | Kimi K2.5 |
| **Kimi For Coding** | Optional | ✅ Active | Code-specialized Kimi |

---

## 1️⃣ NVIDIA NIM (PRIMARY)

### Installation
```bash
# Add provider
opencode auth add nvidia-nim
```

### Configuration
**File:** `~/.config/opencode/opencode.json`

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
          "name": "Qwen 3.5 397B",
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

### Environment Variables
**File:** `~/.zshrc`
```bash
export NVIDIA_API_KEY="nvapi-your-key-here"
```

### Verification
```bash
# Test connection
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
  https://integrate.api.nvidia.com/v1/models

# List models in OpenCode
opencode models | grep nvidia

# Test request
opencode --model nvidia-nim/qwen-3.5-397b "Hello"
```

### ⚠️ CRITICAL NOTES
- **Timeout:** MUST be 120000ms (Qwen 3.5 has 70-90s latency)
- **Model ID:** MUST be `qwen/qwen3.5-397b-a17b` (NOT `qwen2.5`)
- **Rate Limit:** 40 RPM (Free Tier)
- **HTTP 429:** Wait 60 seconds, then retry

---

## 2️⃣ OPENCODE ZEN (FREE FALLBACK)

### Configuration
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
          }
        },
        "kimi-k2.5-free": {
          "name": "Kimi K2.5 Free",
          "limit": {
            "context": 1048576,
            "output": 65536
          }
        },
        "glm-4.7-free": {
          "name": "GLM 4.7 Free",
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

### Features
- ✅ 100% FREE
- ✅ No API key required
- ✅ Uncensored models
- ✅ High context limits

### Usage
```bash
# Use FREE model
opencode --model opencode/kimi-k2.5-free "Hello"

# Use uncensored model
opencode --model opencode/minimax-m2.5-free "Hello"
```

---

## 3️⃣ GOOGLE ANTIGRAVITY (GEMINI)

### Installation
```bash
# Add provider
opencode auth add google

# Start OAuth flow
opencode auth login
```

### Configuration
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

### ⚠️ IMPORTANT
- Use PRIVATE Gmail (NOT Google Workspace)
- OAuth flow required
- Tokens stored in `~/.config/opencode/antigravity-accounts.json`

### Verification
```bash
opencode auth status
opencode models | grep antigravity
```

---

## 4️⃣ MOONSHOT AI (KIMI)

### Installation
```bash
opencode auth add moonshot-ai
```

### Configuration
```json
{
  "provider": {
    "moonshot-ai": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "Moonshot AI",
      "options": {
        "baseURL": "https://api.moonshot.cn/v1"
      },
      "models": {
        "kimi-k2.5": {
          "name": "Kimi K2.5",
          "limit": {
            "context": 1048576,
            "output": 65536
          }
        }
      }
    }
  }
}
```

---

## 5️⃣ KIMI FOR CODING

### Installation
```bash
opencode auth add kimi-for-coding
```

### Configuration
```json
{
  "provider": {
    "kimi-for-coding": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "Kimi For Coding",
      "options": {
        "baseURL": "https://api.kimi-for-coding.ai/v1"
      },
      "models": {
        "k2p5": {
          "name": "K2.5 for Coding",
          "limit": {
            "context": 200000,
            "output": 65536
          }
        }
      }
    }
  }
}
```

---

## 🔄 FALLBACK CHAIN STRATEGY

### Recommended Order
1. **Primary:** `nvidia-nim/qwen-3.5-397b` (smartest)
2. **Fallback 1:** `opencode/kimi-k2.5-free` (FREE)
3. **Fallback 2:** `opencode/minimax-m2.5-free` (uncensored)
4. **Fallback 3:** `moonshot-ai/kimi-k2.5` (backup)

### Implementation
```typescript
// External fallback logic (not in opencode.json)
const fallbackChain = [
  'nvidia-nim/qwen-3.5-397b',
  'opencode/kimi-k2.5-free',
  'opencode/minimax-m2.5-free',
  'moonshot-ai/kimi-k2.5'
];

async function callWithFallback(prompt: string) {
  for (const model of fallbackChain) {
    try {
      return await callModel(model, prompt);
    } catch (error) {
      console.warn(`Model ${model} failed, trying next...`);
    }
  }
  throw new Error('All models failed');
}
```

---

## 📊 MODEL COMPARISON

| Model | Context | Output | Speed | Cost | Best For |
|-------|---------|--------|-------|------|----------|
| Qwen 3.5 397B | 262K | 32K | Slow (70-90s) | FREE | Code (BEST) |
| Kimi K2.5 | 1M | 64K | Medium | FREE | General |
| Gemini 3 Pro | 2M | 64K | Fast | FREE | Multimodal |
| Big Pickle | 200K | 128K | Fast | FREE | Uncensored |

---

## 🛠️ TROUBLESHOOTING

### Issue: Model not found
```bash
# Verify model ID
opencode models | grep <model-name>

# Check provider auth
opencode auth list

# Refresh auth
opencode auth refresh
```

### Issue: Timeout errors
**Solution:** Increase timeout in opencode.json
```json
{
  "options": {
    "timeout": 120000
  }
}
```

### Issue: Rate limit (HTTP 429)
**Solution:** Wait 60 seconds, use fallback model

### Issue: OAuth expired
```bash
# Re-authenticate
opencode auth logout
opencode auth login
```

---

## ✅ VERIFICATION CHECKLIST

- [ ] All providers added: `opencode auth list`
- [ ] Models visible: `opencode models`
- [ ] Test requests successful for each provider
- [ ] Fallback chain configured
- [ ] Environment variables set: `echo $NVIDIA_API_KEY`
- [ ] Timeout configured (120000ms for NVIDIA)
- [ ] MCP servers connected

---

**Status:** ✅ PRODUCTION-READY  
**Last Updated:** 2026-02-18  
**Source:** `~/.config/opencode/opencode.json`, `~/.opencode/providers/`
