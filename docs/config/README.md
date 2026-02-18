# ⚙️ BIOMETRICS Documentation - Configuration

**Provider and system configuration references.**

---

## 📁 Configuration Files

| File | Description | Source |
|------|-------------|--------|
| [PROVIDER.md](PROVIDER.md) | Provider overview | Original |
| [OPENCODE.md](OPENCODE.md) | OpenCode configuration | Original |
| [OPENCLAW.md](OPENCLAW.md) | OpenClaw configuration | Original |
| [SUPABASE.md](SUPABASE.md) | Supabase database config | Original |
| [CLOUDFLARE.md](CLOUDFLARE.md) | Cloudflare tunnel config | Original |
| [VERCEL.md](VERCEL.md) | Vercel deployment config | Original |
| [IONOS.md](IONOS.md) | IONOS hosting config | Original |
| [PWA-CONFIG.md](PWA-CONFIG.md) | Progressive Web App config | Original |
| [ALERTING-CONFIG.md](ALERTING-CONFIG.md) | Alerting system config | Original |
| [CDN-CONFIG.md](CDN-CONFIG.md) | CDN configuration | Original |

---

## 🔌 Provider Configuration

### NVIDIA NIM (Qwen 3.5 397B)
- **Best for:** Code generation, complex reasoning
- **Context:** 262K tokens
- **Rate Limit:** 40 RPM (FREE tier)
- **Timeout:** 120000ms (MANDATORY!)

### Google Gemini (Antigravity)
- **Best for:** Creative tasks, UI/UX, multimodal
- **Context:** 1M+ tokens
- **Authentication:** OAuth required
- **Models:** Flash, Pro, Claude via Antigravity

### OpenCode ZEN (FREE)
- **Best for:** Uncensored generation, fallback
- **Cost:** 100% FREE
- **Models:** Big Pickle, Grok Code, GLM 4.7 Free

---

## 📝 Configuration Best Practices

1. **NEVER commit secrets** - Use environment variables
2. **NO timeout entries** - Except NVIDIA NIM (120s)
3. **Test after changes** - Always verify with `opencode models`
4. **Backup configs** - Before making changes

---

## 🔧 Verification Commands

```bash
# Check OpenCode config
opencode models

# Check for timeout entries (MUST BE EMPTY!)
grep -r "timeout" ~/.config/opencode/opencode.json

# Verify providers
opencode models | grep -E "(nvidia|google|opencode-zen)"
```

---

**Last Updated:** 2026-02-18  
**Status:** ✅ Production-Ready
