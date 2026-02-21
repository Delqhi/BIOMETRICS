# SETUP-COMPLETE.md - Post-Clone Verification Checklist

**Last Updated:** 2026-02-18  
**Purpose:** Ensure BIOMETRICS agents work on ANY Mac after git clone  
**Status:** ✅ MANDATORY FOR ALL USERS

---

## 🎯 WHY THIS EXISTS

**Problem:** Agents fail on cloned Macs because setup steps are missed.

**Solution:** This checklist ensures EVERY step is completed in correct order.

---

## 📋 PRE-FLIGHT CHECKLIST

### Before You Start

- [ ] macOS 13+ (Ventura or later)
- [ ] Internet connection stable
- [ ] Terminal app ready (iTerm2 or default)
- [ ] GitHub account logged in
- [ ] ~30 minutes available

---

## 🔴 PHASE 1: SYSTEM PREREQUISITES (10 min)

### 1.1 Install Homebrew

```bash
# Check if already installed
which brew

# If not installed:
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Verify
brew --version
```

- [ ] Homebrew installed

### 1.2 Install Node.js 20+

```bash
brew install node@20
node --version  # Must be v20.x.x
```

- [ ] Node.js 20+ installed

### 1.3 Install pnpm

```bash
brew install pnpm
pnpm --version  # Must be 8.x.x or higher
```

- [ ] pnpm installed

### 1.4 Install Python 3.11+

```bash
brew install python@3.11
python3 --version  # Must be 3.11.x or higher
```

- [ ] Python 3.11+ installed

---

## 🔴 PHASE 2: NVIDIA API KEY (5 min) - CRITICAL!

### 2.1 Register on NVIDIA Build

```bash
# Open browser and go to:
open https://build.nvidia.com/
```

- [ ] Account created on NVIDIA Build
- [ ] Email verified

### 2.2 Generate API Key

```bash
# In browser:
# 1. Dashboard → API Keys
# 2. Click "Generate New Key"
# 3. Copy the key (starts with "nvapi-")
```

- [ ] API Key generated and copied

### 2.3 Add to Shell Config

```bash
# Add to ~/.zshrc (or ~/.bashrc if using bash):
cat >> ~/.zshrc << 'SHELL_EOF'

# NVIDIA API Keys (BIOMETRICS)
export NVIDIA_API_KEY="nvapi-YOUR_KEY_HERE"
export NVIDIA_NIM_API_KEY="nvapi-YOUR_KEY_HERE"
SHELL_EOF

# Replace YOUR_KEY_HERE with actual key:
nano ~/.zshrc  # Edit and paste your key
```

- [ ] NVIDIA_API_KEY added to ~/.zshrc
- [ ] NVIDIA_NIM_API_KEY added to ~/.zshrc
- [ ] Keys saved (no typos!)

### 2.4 Reload Shell (MANDATORY!)

```bash
# THIS STEP IS CRITICAL - DO NOT SKIP!
exec zsh

# Verify keys are loaded:
echo $NVIDIA_API_KEY
echo $NVIDIA_NIM_API_KEY
```

- [ ] Shell reloaded with `exec zsh`
- [ ] Keys visible in echo output

---

## 🔴 PHASE 3: OPENCODE CONFIGURATION (10 min)

### 3.1 Install Opencode

```bash
brew install opencode
opencode --version
```

- [ ] Opencode installed

### 3.2 Authenticate Opencode

```bash
# This opens browser for OAuth:
opencode auth login

# Wait for browser to open
# Complete OAuth flow
# Return to terminal
```

- [ ] OAuth completed
- [ ] Authentication successful

### 3.3 Add Providers

```bash
# Add all required providers:
opencode auth add moonshot-ai
opencode auth add kimi-for-coding
opencode auth add opencode-zen
```

- [ ] moonshot-ai added
- [ ] kimi-for-coding added
- [ ] opencode-zen added

### 3.4 Verify Models

```bash
opencode models

# Expected output (at minimum):
# - moonshotai/kimi-k2.5
# - qwen/qwen3.5-397b-a17b
# - opencode/minimax-m2.5-free
# - opencode/minimax-m2.5-free
# - opencode/minimax-m2.5-free
```

- [ ] At least 5 models visible
- [ ] All three providers represented

### 3.5 Verify opencode.json (MANDATE 0.35)

```bash
# Check file exists:
ls -la ~/.config/opencode/opencode.json

# CRITICAL: Check for timeout violation
grep -r "timeout" ~/.config/opencode/opencode.json

# MUST BE EMPTY! If not:
nano ~/.config/opencode/opencode.json
# Remove ALL lines containing "timeout"

# Validate JSON:
cat ~/.config/opencode/opencode.json | python3 -m json.tool
```

- [ ] opencode.json exists
- [ ] NO timeout entries (MANDATE 0.35 compliant)
- [ ] JSON syntax valid

---

## 🟠 PHASE 4: OPENCLAW CONFIGURATION (5 min)

### 4.1 Install OpenClaw

```bash
# Check if already installed
which openclaw

# If not:
brew install openclaw
openclaw --version
```

- [ ] OpenClaw installed

### 4.2 Configure openclaw.json

```bash
# Check file exists:
ls -la ~/.openclaw/openclaw.json

# Verify NVIDIA_API_KEY in env section:
grep "NVIDIA_API_KEY" ~/.openclaw/openclaw.json

# If missing, edit:
nano ~/.openclaw/openclaw.json
# Add to "env" section:
# "NVIDIA_API_KEY": "nvapi-YOUR_KEY_HERE"
```

- [ ] openclaw.json exists
- [ ] NVIDIA_API_KEY in env section

### 4.3 Health Check

```bash
openclaw doctor --fix

# Should show:
# ✅ All checks passed
```

- [ ] Health check passed

**NOTE:** OpenClaw CAN have timeout in config (Gateway manages it). This is OK and different from OpenCode!

---

## 🟠 PHASE 5: OH-MY-OPENCODE CONFIGURATION (3 min)

### 5.1 Verify oh-my-opencode.json

```bash
# Check file exists:
ls -la ~/.config/opencode/oh-my-opencode.json

# Verify NO timeout entries:
grep -r "timeout" ~/.config/opencode/oh-my-opencode.json

# MUST BE EMPTY!
```

- [ ] oh-my-opencode.json exists
- [ ] NO timeout entries

### 5.2 Verify Agent Models

```bash
cat ~/.config/opencode/oh-my-opencode.json | python3 -m json.tool

# Check that agents have models assigned:
# - sisyphus: moonshotai/kimi-k2.5
# - librarian: opencode/minimax-m2.5-free
# - explore: opencode/minimax-m2.5-free
```

- [ ] Agent models configured correctly

---

## 🟢 PHASE 6: FINAL RESTART & VERIFICATION (2 min)

### 6.1 Complete Terminal Restart

```bash
# This ensures ALL configs are loaded:
exec zsh

# Wait for shell to fully reload
```

- [ ] Terminal restarted

### 6.2 Run Verification Script

```bash
# Create and run verification:
cat > /tmp/verify-biometrics.sh << 'VERIFY_EOF'
#!/bin/bash
echo "=== BIOMETRICS SETUP VERIFICATION ==="
echo ""

# 1. NVIDIA API Key
if [ -z "$NVIDIA_API_KEY" ]; then
  echo "❌ NVIDIA_API_KEY not set"
  exit 1
else
  echo "✅ NVIDIA_API_KEY: ${#NVIDIA_API_KEY} chars"
fi

# 2. Opencode
if ! command -v opencode &> /dev/null; then
  echo "❌ Opencode not installed"
  exit 1
else
  echo "✅ Opencode: $(opencode --version)"
fi

# 3. Models
model_count=$(opencode models 2>/dev/null | wc -l)
if [ "$model_count" -lt 5 ]; then
  echo "❌ Not enough models"
  exit 1
else
  echo "✅ Models: $model_count configured"
fi

# 4. Timeout check
timeout_count=$(grep -r "timeout" ~/.config/opencode/opencode.json 2>/dev/null | wc -l)
if [ "$timeout_count" -gt 0 ]; then
  echo "❌ MANDATE 0.35 VIOLATION"
  exit 1
else
  echo "✅ No timeout in opencode.json"
fi

# 5. OpenClaw
if [ -f ~/.openclaw/openclaw.json ]; then
  echo "✅ OpenClaw configured"
else
  echo "⚠️  OpenClaw not configured (optional)"
fi

echo ""
echo "=== ALL CHECKS PASSED ==="
VERIFY_EOF

chmod +x /tmp/verify-biometrics.sh
/tmp/verify-biometrics.sh
```

- [ ] All verification checks passed

---

## 🎯 POST-SETUP: CLONE BIOMETRICS

```bash
# Now you can clone and use BIOMETRICS:
cd ~/dev  # Or your preferred directory
git clone https://github.com/Delqhi/BIOMETRICS.git
cd BIOMETRICS

# Run the biometrics CLI:
cd biometrics-cli
pnpm install
pnpm link --global
biometrics
```

- [ ] BIOMETRICS cloned
- [ ] CLI installed
- [ ] Onboarding completed

---

## 🚨 TROUBLESHOOTING

### Issue: "NVIDIA_API_KEY not found"

**Solution:**
```bash
# 1. Check ~/.zshrc has export
grep NVIDIA ~/.zshrc

# 2. Reload shell
exec zsh

# 3. Verify
echo $NVIDIA_API_KEY
```

### Issue: "Provider not found"

**Solution:**
```bash
# Re-add providers
opencode auth add moonshot-ai
opencode auth add kimi-for-coding
opencode auth add opencode-zen

# Reload
exec zsh
```

### Issue: "Models not loading"

**Solution:**
```bash
# Check authentication
opencode auth status

# Re-login if needed
opencode auth logout
opencode auth login
```

### Issue: "timeout in config" error

**Solution:**
```bash
# Remove timeout immediately!
nano ~/.config/opencode/opencode.json
# Delete ALL lines with "timeout"

# Validate
cat ~/.config/opencode/opencode.json | python3 -m json.tool
```

---

## ✅ COMPLETION CERTIFICATE

```
╔══════════════════════════════════════════════════╗
║     BIOMETRICS SETUP CERTIFICATE OF COMPLETION   ║
╠══════════════════════════════════════════════════╣
║                                                  ║
║  User: _________________                         ║
║  Date: _________________                         ║
║  Mac:  _________________                         ║
║                                                  ║
║  ✅ System Prerequisites      COMPLETE           ║
║  ✅ NVIDIA API Key           COMPLETE           ║
║  ✅ Opencode Configuration   COMPLETE           ║
║  ✅ OpenClaw Configuration   COMPLETE           ║
║  ✅ Final Verification       COMPLETE           ║
║                                                  ║
║  Status: READY FOR PRODUCTION                    ║
║                                                  ║
╚══════════════════════════════════════════════════╝
```

---

## 📚 NEXT STEPS

1. Read `BIOMETRICS/README.md` for project overview
2. Read `BIOMETRICS/AGENTS.md` for agent usage
3. Read `BIOMETRICS/PROVIDER.md` for detailed provider info
4. Start building with `biometrics` CLI

**Welcome to BIOMETRICS!** 🚀
