# PROVIDER.md - AI Provider Configuration Guide

**Last Updated:** 2026-02-18  
**Status:** ✅ PRODUCTION READY  
**Mandate:** 0.35 (NO TIMEOUT in OpenCode configs)

---

## 🚨 CRITICAL: WHY AGENTS FAIL ON CLONED MAC

When you clone the BIOMETRICS repository on a new Mac, agents fail because:

1. ❌ **No API Keys configured** - NVIDIA_API_KEY not set
2. ❌ **Opencode not authenticated** - `opencode auth login` not run
3. ❌ **Providers not added** - moonshot-ai, kimi-for-coding, opencode-zen missing
4. ❌ **Shell not reloaded** - Environment variables not loaded
5. ❌ **Config files missing** - opencode.json, oh-my-opencode.json not created

---

## 📋 COMPLETE SETUP CHECKLIST

### Phase 1: System Prerequisites

```bash
# 1. Install Homebrew (if not installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 2. Install Node.js 20+
brew install node@20

# 3. Install pnpm
brew install pnpm

# 4. Install Python 3.11+
brew install python@3.11

# 5. Verify installations
node --version  # v20.x.x
pnpm --version  # 8.x.x
python3 --version  # 3.11.x
```

### Phase 2: NVIDIA API KEY (CRITICAL - DO THIS FIRST!)

```bash
# 1. Register on NVIDIA Build
# Go to: https://build.nvidia.com/
# Create account with email

# 2. Generate API Key
# Dashboard → API Keys → Generate New Key
# Copy the key (starts with "nvapi-")

# 3. Add to shell config
cat >> ~/.zshrc << 'SHELL_EOF'

# NVIDIA API Keys (BIOMETRICS)
export NVIDIA_API_KEY="nvapi-YOUR_KEY_HERE"
export NVIDIA_NIM_API_KEY="nvapi-YOUR_KEY_HERE"
SHELL_EOF

# 4. RELOAD SHELL (MANDATORY!)
exec zsh  # OR: source ~/.zshrc

# 5. Verify
echo $NVIDIA_API_KEY  # Must show your key
```

### Phase 3: Opencode Installation & Configuration

```bash
# 1. Install Opencode
brew install opencode

# 2. Authenticate (opens browser)
opencode auth login

# 3. Add Providers
opencode auth add moonshot-ai
opencode auth add kimi-for-coding
opencode auth add opencode-zen

# 4. Verify models
opencode models

# Expected output should include:
# - moonshotai/kimi-k2.5
# - kimi-for-coding/k2p5
# - opencode-zen/zen/big-pickle
```

### Phase 4: Config File Verification

```bash
# 1. Check opencode.json exists
ls -la ~/.config/opencode/opencode.json

# 2. Verify NO timeout entries (MANDATE 0.35)
grep -r "timeout" ~/.config/opencode/opencode.json
# MUST BE EMPTY! If not, edit and remove timeout lines.

# 3. Validate JSON syntax
cat ~/.config/opencode/opencode.json | python3 -m json.tool

# 4. Check oh-my-opencode.json
ls -la ~/.config/opencode/oh-my-opencode.json
cat ~/.config/opencode/oh-my-opencode.json | python3 -m json.tool
```

### Phase 5: OpenClaw Configuration

```bash
# 1. Check openclaw.json exists
ls -la ~/.openclaw/openclaw.json

# 2. Verify NVIDIA_API_KEY in env section
grep "NVIDIA_API_KEY" ~/.openclaw/openclaw.json
# Should show: "NVIDIA_API_KEY": "nvapi-..."

# 3. Run health check
openclaw doctor --fix

# NOTE: OpenClaw CAN have timeout in config (Gateway manages it)
# This is DIFFERENT from OpenCode!
```

### Phase 6: Terminal Restart (ABSOLUTELY MANDATORY!)

```bash
# After ALL configurations, RESTART TERMINAL:
exec zsh  # This reloads your shell completely

# Verify everything is loaded:
echo "=== VERIFICATION ==="
echo "NVIDIA_API_KEY: ${NVIDIA_API_KEY:0:20}..."
echo "Node: $(node --version)"
echo "Opencode: $(opencode --version)"
echo "OpenClaw: $(openclaw --version 2>/dev/null || 'not installed')"
```

---

## 🔴 COMMON ERRORS & SOLUTIONS

| Error | Cause | Solution |
|-------|-------|----------|
| `NVIDIA_API_KEY not found` | Shell not reloaded | Run `exec zsh` |
| `Provider not found` | Not authenticated | Run `opencode auth add <provider>` |
| `timeout in config` | Violation of MANDATE 0.35 | Remove timeout from opencode.json |
| `Models not loading` | Missing API key | Check ~/.zshrc has NVIDIA_API_KEY |
| `OpenClaw not found` | Not installed | Run `brew install openclaw` |

---

## ✅ FINAL VERIFICATION SCRIPT

```bash
#!/bin/bash
# Save as verify-providers.sh and run

echo "=== PROVIDER VERIFICATION ==="
echo ""

# 1. Check NVIDIA API Key
if [ -z "$NVIDIA_API_KEY" ]; then
  echo "❌ NVIDIA_API_KEY not set"
  echo "   Fix: Add to ~/.zshrc and run 'exec zsh'"
  exit 1
else
  echo "✅ NVIDIA_API_KEY set (${#NVIDIA_API_KEY} chars)"
fi

# 2. Check Opencode
if ! command -v opencode &> /dev/null; then
  echo "❌ Opencode not installed"
  echo "   Fix: brew install opencode"
  exit 1
else
  echo "✅ Opencode installed: $(opencode --version)"
fi

# 3. Check models
model_count=$(opencode models 2>/dev/null | wc -l)
if [ "$model_count" -lt 5 ]; then
  echo "❌ Not enough models configured"
  echo "   Fix: opencode auth add moonshot-ai kimi-for-coding opencode-zen"
  exit 1
else
  echo "✅ $model_count models configured"
fi

# 4. Check timeout violation
timeout_count=$(grep -r "timeout" ~/.config/opencode/opencode.json 2>/dev/null | wc -l)
if [ "$timeout_count" -gt 0 ]; then
  echo "❌ MANDATE 0.35 VIOLATION: timeout found in opencode.json"
  echo "   Fix: Remove timeout entries immediately!"
  exit 1
else
  echo "✅ No timeout in opencode.json (MANDATE 0.35 compliant)"
fi

# 5. Check OpenClaw
if [ -f ~/.openclaw/openclaw.json ]; then
  echo "✅ OpenClaw configured"
  openclaw doctor --fix 2>/dev/null && echo "✅ OpenClaw health check passed"
else
  echo "⚠️  OpenClaw not configured (optional)"
fi

echo ""
echo "=== ALL CHECKS PASSED ==="
echo "Your Mac is ready for BIOMETRICS agents!"
```

---

## 📚 REFERENCES

- **MANDATE 0.35:** No timeout in OpenCode configs
- **BIOMETRICS/README.md:** Post-clone setup instructions
- **BIOMETRICS/TROUBLESHOOTING.md:** Detailed troubleshooting
- **~/.config/opencode/AGENTS.md:** Global agent mandates

---

## 🎯 QUICK REFERENCE

```bash
# One-liner setup (after git clone):
exec zsh && opencode auth login && opencode auth add moonshot-ai kimi-for-coding opencode-zen && opencode models

# Verify everything:
echo $NVIDIA_API_KEY | head -c 20 && echo "..." && opencode --version && openclaw doctor --fix
```

**Remember:** After ANY config change → `exec zsh` → Verify!
