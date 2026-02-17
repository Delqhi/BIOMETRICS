# biometrics-onboard

**Unified onboarding CLI for complete BIOMETRICS setup.** One command installs and configures EVERYTHING:

- ✅ **GitLab** (media storage project for videos, PDFs, images)
- ✅ **NLM CLI** (NotebookLM command-line interface)
- ✅ **OpenCode** (AI coding assistant)
- ✅ **OpenClaw** (AI orchestration)
- ✅ **Google Antigravity** (OAuth authentication)
- ✅ **WhatsApp Integration** (Business API)
- ✅ **Telegram Integration** (Bot API)
- ✅ **Gmail Integration** (OAuth2)
- ✅ **Twitter/X Integration** (OAuth2)
- ✅ **ClawdBot** (Social media automation)

## Installation

### Option 1: Global Installation (Recommended)

```bash
# Clone the repository
git clone https://github.com/Delqhi/biometrics-onboard.git
cd biometrics-onboard

# Install dependencies with pnpm
pnpm install

# Link globally
pnpm link --global
```

### Option 2: Direct Usage

```bash
# Clone and run directly
git clone https://github.com/Delqhi/biometrics-onboard.git
cd biometrics-onboard
pnpm install
pnpm start
```

## Usage

After installation, simply run:

```bash
biometrics-onboard
```

### 🎯 Unified Onboarding Flow

The CLI guides you through a **single, unified setup process**:

1. **GitLab Setup** - Create media storage project (with API key help link)
2. **NVIDIA API Key** - Enter your key for Qwen 3.5 397B access
3. **Social Media Integrations**:
   - ✅ WhatsApp Business API (optional)
   - ✅ Telegram Bot API (optional)
   - ✅ Gmail OAuth (optional)
   - ✅ Twitter/X API (optional)
4. **OpenCode Installation** - Install via Homebrew (optional)
5. **OpenClaw Installation** - Install via pnpm (optional)
6. **NLM CLI** - Auto-install + browser authentication
7. **Google Antigravity** - Auto-install + browser authentication
8. **ClawdBot Integration** - Connect all social channels
9. **Verification** - Test all installations and integrations

### 🔗 API Key Help Links

**Don't have API keys?** The CLI shows direct links to create them:

- **GitLab:** https://gitlab.com/-/profile/personal_access_tokens
- **NVIDIA:** https://build.nvidia.com/explore/discover
- **WhatsApp:** https://developers.facebook.com/apps/creation/
- **Telegram:** https://core.telegram.org/bots/features#botfather
- **Gmail:** https://console.cloud.google.com/apis/credentials
- **Twitter:** https://developer.twitter.com/en/portal/dashboard

Just click the link, create your API key, and paste it into the CLI!

## What Gets Configured

### OpenCode (`~/.config/opencode/opencode.json`)
```json
{
  "provider": {
    "google": {
      "npm": "@ai-sdk/google",
      "models": {
        "gemini-2.5-pro": {
          "id": "gemini-2.5-pro",
          "name": "Gemini 2.5 Pro"
        }
      }
    },
    "nvidia": {
      "npm": "@ai-sdk/openai-compatible",
      "options": {
        "baseURL": "https://integrate.api.nvidia.com/v1"
      },
      "models": {
        "qwen-3.5-397b": {
          "id": "qwen/qwen3.5-397b-a17b",
          "name": "Qwen 3.5 397B"
        }
      }
    }
  }
}
```

### OpenClaw (`~/.openclaw/openclaw.json`)
```json
{
  "env": {
    "NVIDIA_API_KEY": "your-api-key"
  },
  "models": {
    "providers": {
      "nvidia": {
        "baseUrl": "https://integrate.api.nvidia.com/v1",
        "api": "openai-completions",
        "models": []
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

## Requirements

- **Node.js** >= 18.0.0
- **pnpm** (for package management and OpenClaw installation)
- **Homebrew** (for OpenCode installation on macOS)
- **Internet connection** (for API authentication)

## Manual Installation Commands

If the CLI fails, you can install manually:

```bash
# NLM CLI
pnpm add -g nlm-cli
nlm auth login

# OpenCode (macOS)
brew install opencode

# OpenClaw
pnpm add -g @delqhi/openclaw

# Google Antigravity Plugin
opencode plugin add opencode-antigravity-auth
opencode auth login
```

## Troubleshooting

### "Command not found: biometrics-onboard"

Make sure the global pnpm bin directory is in your PATH:

```bash
# Add to ~/.zshrc or ~/.bashrc
export PATH="$HOME/.local/share/pnpm:$PATH"
```

### "OpenCode installation failed"

OpenCode requires Homebrew on macOS:

```bash
# Install Homebrew first
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Then retry
biometrics-onboard
```

### "NLM CLI authentication failed"

Make sure your browser opens and you complete the OAuth flow. If not:

```bash
# Manual authentication
nlm auth login
```

## Development

```bash
# Clone repository
git clone https://github.com/Delqhi/biometrics-onboard.git
cd biometrics-onboard

# Install dependencies
pnpm install

# Run in development mode
pnpm start

# Test the CLI
pnpm test
```

## License

MIT

## Support

For issues or questions, please open an issue on GitHub or contact the BIOMETRICS team.
