# 🐚 SHELL SETUP - ZSH CONFIGURATION

**Version:** 1.0 "COMPLETE ENVIRONMENT"  
**Source:** Consolidated from user environment

---

## 📁 ESSENTIAL DIRECTORIES

```bash
# Create standard directory structure
mkdir -p ~/.config/opencode
mkdir -p ~/dev/sin-code
mkdir -p ~/dev/SIN-Solver
mkdir -p ~/dev/BIOMETRICS
mkdir -p ~/Bilder/AI-Screenshots/{playwright,skyvern,steel,stagehand,opencode}
mkdir -p ~/dev/environments-jeremy.md
```

---

## 🔧 .ZSHRC CONFIGURATION

**File:** `~/.zshrc`

```bash
# ═══════════════════════════════════════════════════════════
# OPENCODE & AI ENVIRONMENT
# ═══════════════════════════════════════════════════════════

# NVIDIA API Key (PRIMARY)
export NVIDIA_API_KEY="nvapi-your-key-here"

# OpenCode API Key (if needed)
export OPENCODE_API_KEY="xxx"

# Mistral API Key
export MISTRAL_API_KEY="your-mistral-key"

# Anthropic API Key
export ANTHROPIC_API_KEY="your-anthropic-key"

# ═══════════════════════════════════════════════════════════
# PATH CONFIGURATION
# ═══════════════════════════════════════════════════════════

# Homebrew (Apple Silicon)
export PATH="/opt/homebrew/bin:$PATH"

# Node.js global binaries
export PATH="$HOME/.npm-global/bin:$PATH"

# Python
export PATH="$HOME/.local/bin:$PATH"

# ═══════════════════════════════════════════════════════════
# ALIASES
# ═══════════════════════════════════════════════════════════

# OpenCode shortcuts
alias oc='opencode'
alias oc-auth='opencode auth list'
alias oc-models='opencode models'
alias oc-login='opencode auth login'

# Git shortcuts
alias gs='git status'
alias ga='git add'
alias gc='git commit'
alias gp='git push'
alias gl='git log --oneline -10'
alias gco='git checkout'

# Docker shortcuts
alias dps='docker ps'
alias dlogs='docker logs'
alias dexec='docker exec -it'

# Directory shortcuts
alias dev='cd ~/dev'
alias sin='cd ~/dev/sin-code'
alias bio='cd ~/dev/BIOMETRICS'
alias config='cd ~/.config/opencode'

# ═══════════════════════════════════════════════════════════
# PROMPT CUSTOMIZATION
# ═══════════════════════════════════════════════════════════

# Git branch in prompt
autoload -Uz vcs_info
precmd() { vcs_info }
setopt prompt_subst
PROMPT='%F{green}%n@%m%f:%F{blue}%~%f%F{yellow}$(vcs_info_msg_0_)%f$ '
ZSH_THEME_GIT_PROMPT_PREFIX=' (%F{magenta}'
ZSH_THEME_GIT_PROMPT_SUFFIX='%f)'

# ═══════════════════════════════════════════════════════════
# COMPLETION
# ═══════════════════════════════════════════════════════════

# Enable completion
autoload -Uz compinit
compinit

# Case insensitive completion
zstyle ':completion:*' matcher-list 'm:{a-zA-Z}={A-Za-z}'

# ═══════════════════════════════════════════════════════════
# HISTORY
# ═══════════════════════════════════════════════════════════

# Larger history
HISTSIZE=10000
SAVEHIST=10000
HISTFILE=~/.zsh_history

# Share history between sessions
setopt SHARE_HISTORY
setopt EXTENDED_HISTORY

# ═══════════════════════════════════════════════════════════
# EDITOR
# ═══════════════════════════════════════════════════════════

# Use vim mode (optional)
# bindkey -v

# Use nano as default editor
export EDITOR=nano
export VISUAL=nano

# ═══════════════════════════════════════════════════════════
# LAUNCH AGENTS (AUTO-CLEANUP)
# ═══════════════════════════════════════════════════════════

# Desktop cleanup (daily at 3 AM)
# launchctl load ~/Library/LaunchAgents/com.sincode.desktop-cleanup.plist

# AI Screenshot cleanup (daily at 4 AM)
# launchctl load ~/Library/LaunchAgents/com.sincode.ai-screenshot-cleanup.plist

# ═══════════════════════════════════════════════════════════
# CUSTOM FUNCTIONS
# ═══════════════════════════════════════════════════════════

# Quick git commit
gca() {
  git add -A
  git commit -m "$1"
  git push
}

# Create and enter directory
mkcd() {
  mkdir -p "$1" && cd "$1"
}

# Find files by name
ff() {
  find . -name "*$1*" -type f
}

# Search code (ripgrep)
search() {
  rg "$1" --type "$2"
}

# ═══════════════════════════════════════════════════════════
# NVM (NODE VERSION MANAGER)
# ═══════════════════════════════════════════════════════════

export NVM_DIR="$HOME/.nvm"
[ -s "/opt/homebrew/opt/nvm/nvm.sh" ] && \. "/opt/homebrew/opt/nvm/nvm.sh"
[ -s "/opt/homebrew/opt/nvm/etc/bash_completion.d/nvm" ] && \. "/opt/homebrew/opt/nvm/etc/bash_completion.d/nvm"

# ═══════════════════════════════════════════════════════════
# PYTHON
# ═══════════════════════════════════════════════════════════

export PYENV_ROOT="$HOME/.pyenv"
command -v pyenv >/dev/null || export PATH="$PYENV_ROOT/bin:$PATH"
eval "$(pyenv init -)"

# ═══════════════════════════════════════════════════════════
# MISCELLANEOUS
# ═══════════════════════════════════════════════════════════

# Disable auto-correct
unsetopt CORRECT

# Enable color
export CLICOLOR=1
export LSCOLORS=GxFxCxDxBxegedabagaced

# Set terminal title
echo -ne "\033]0;Terminal\007"
```

---

## 🔐 SECRETS REGISTRY

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

## Mistral AI - 2026-02-16
**Service:** Mistral API
**API Key:** [REDACTED] (in ~/.zshrc)
**Endpoint:** https://api.mistral.ai/v1
**Status:** ACTIVE

## Anthropic - 2026-02-16
**Service:** Anthropic API
**API Key:** [REDACTED] (in ~/.zshrc)
**Endpoint:** https://api.anthropic.com/v1
**Status:** ACTIVE
```

**⚠️ CRITICAL:**
- NEVER delete from this file (APPEND-ONLY)
- Mark rotated secrets as "ROTATED" but don't delete
- Store actual keys in ~/.zshrc (not in this file)

---

## 🛠️ LAUNCH AGENTS

### Desktop Cleanup
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
    <key>StandardOutPath</key>
    <string>/tmp/desktop-cleanup.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/desktop-cleanup.err</string>
</dict>
</plist>
```

### AI Screenshot Cleanup
**File:** `~/Library/LaunchAgents/com.sincode.ai-screenshot-cleanup.plist`

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.sincode.ai-screenshot-cleanup</string>
    <key>ProgramArguments</key>
    <array>
        <string>/bin/bash</string>
        <string>-c</string>
        <string>find ~/Bilder/AI-Screenshots -type f -mtime +7 -exec mv {} ~/Bilder/AI-Screenshots/archive/ \;</string>
    </array>
    <key>StartCalendarInterval</key>
    <dict>
        <key>Hour</key>
        <integer>4</integer>
        <key>Minute</key>
        <integer>0</integer>
    </dict>
</dict>
</plist>
```

### Load Launch Agents
```bash
launchctl load ~/Library/LaunchAgents/com.sincode.desktop-cleanup.plist
launchctl load ~/Library/LaunchAgents/com.sincode.ai-screenshot-cleanup.plist
```

---

## ✅ VERIFICATION

```bash
# Reload zshrc
source ~/.zshrc

# Check environment variables
echo $NVIDIA_API_KEY
echo $OPENCODE_API_KEY

# Test aliases
oc --version
gs

# Check directories
ls -la ~/dev/
ls -la ~/Bilder/AI-Screenshots/

# Check launch agents
launchctl list | grep sincode
```

---

## 🚨 TROUBLESHOOTING

### Issue: Aliases not working
**Solution:** Run `source ~/.zshrc`

### Issue: Environment variables not set
**Solution:** Check ~/.zshrc syntax, then `source ~/.zshrc`

### Issue: Launch agents not loading
**Solution:** 
```bash
launchctl unload ~/Library/LaunchAgents/com.sincode.*.plist
launchctl load ~/Library/LaunchAgents/com.sincode.*.plist
```

---

**Status:** ✅ PRODUCTION-READY  
**Last Updated:** 2026-02-18
