# 🧠 AI Model Configurations

**Purpose:** AI model definitions, providers, and fallback chains.

**Status:** ✅ Active  
**Last Updated:** 2026-02-19

---

## 📊 Model Overview

### Primary Models

| Model | Provider | Context | Output | Latency | Cost |
|-------|----------|---------|--------|---------|------|
| **Qwen 3.5 397B** | NVIDIA NIM | 262K | 32K | 70-90s | Free |
| **Kimi K2.5** | Moonshot AI | 128K | 8K | 5-10s | Free |
| **MiniMax M2.5** | MiniMax | 1M | 64K | 10-20s | Free |

### Fallback Models

| Model | Provider | Use Case |
|-------|----------|----------|
| **OpenCode ZEN** | OpenCode | Documentation, Research |
| **Qwen2.5-Coder-32B** | NVIDIA NIM | Fast coding tasks |

---

## ⚙️ Configuration

### NVIDIA NIM (Qwen 3.5)

```json
{
  "provider": "nvidia-nim",
  "model": "qwen/qwen3.5-397b-a17b",
  "baseURL": "https://integrate.api.nvidia.com/v1",
  "timeout": 120000,
  "apiKey": "NVIDIA_API_KEY"
}
```

### Fallback Chain

```
Qwen 3.5 → Kimi K2.5 → MiniMax M2.5 → OpenCode ZEN
```

---

## 🔗 References

- [Provider Setup](../../docs/config/)
- [Agent Mapping](../../docs/agents/AGENT-MODEL-MAPPING.md)
