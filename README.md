# 🧬 BIOMETRICS

<div align="center">

**Next-Generation Biometric Authentication Platform**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Status: Active](https://img.shields.io/badge/Status-Active-success)](.)
[![Version: 1.0.0](https://img.shields.io/badge/Version-1.0.0-blue)](.)
[![Best Practices: Feb 2026](https://img.shields.io/badge/Best%20Practices-Feb%202026-orange)](.)

[🚀 Quick Start](#-quick-start) • [📚 Documentation](#-documentation) • [🤖 Agents](#-agents) • [🏗️ Architecture](#-architecture) • [📝 Setup Guide](docs/setup/COMPLETE-SETUP.md)

</div>

---

## 🚨 IMPORTANT: SETUP BEFORE CLONING!

> ⚠️ **BEFORE** you clone this repository, you **MUST** complete the setup steps below. Without proper configuration, the AI agents will **NOT** work!

### 📋 Quick Setup Checklist

```bash
# 1. Get NVIDIA API Key from https://build.nvidia.com/
# 2. Add to ~/.zshrc:
echo 'export NVIDIA_API_KEY="nvapi-YOUR_KEY_HERE"' >> ~/.zshrc

# 3. Reload shell:
exec zsh

# 4. Install OpenCode:
npm install -g opencode

# 5. Authenticate providers:
opencode auth add nvidia-nim
opencode auth add moonshot-ai
opencode auth add kimi-for-coding

# 6. Verify setup:
opencode models | grep nvidia
```

**👉 [FULL SETUP INSTRUCTIONS](docs/setup/COMPLETE-SETUP.md)** ← **START HERE!**

---

## 🚀 Quick Start

### Prerequisites

- ✅ Node.js 20+ installed
- ✅ NVIDIA API Key obtained
- ✅ OpenCode CLI installed
- ✅ All providers authenticated

### Installation

```bash
git clone https://github.com/Delqhi/BIOMETRICS.git
cd BIOMETRICS
npm install
npm run doctor
```

---

## 📚 Documentation

| Category | Files | Description |
|----------|-------|-------------|
| [🛠️ Setup](docs/setup/) | 5 | Installation & configuration |
| [⚙️ Config](docs/config/) | 11 | Provider configurations |
| [🤖 Agents](docs/agents/) | 12 | Agent guides & skills |
| [📖 Best Practices](docs/best-practices/) | 18 | Mandates & workflows |
| [🏗️ Architecture](docs/architecture/) | 26 | System design & APIs |
| [✨ Features](docs/features/) | 32 | Product capabilities |
| [🔬 Advanced](docs/advanced/) | 27 | Blockchain, AI, IoT |

**🔥 Essential:**
- [📋 Universal Blueprint](docs/UNIVERSAL-BLUEPRINT.md)
- [🎬 Video Tutorials](docs/tutorials/)
- [💻 Interactive Examples](docs/examples/)
- [📝 Quizzes](docs/quizzes/)

---

## 🤖 AI Agents

| Agent | Role | Model |
|-------|------|-------|
| **Sisyphus** | Main Coder | Qwen 3.5 397B |
| **Prometheus** | Planning | Qwen 3.5 397B |
| **Oracle** | Architecture | Qwen 3.5 397B |

```bash
opencode "Implement X" --agent sisyphus
```

---

## 📊 Stats

- **Files:** 161+
- **Docs:** 9,606+ lines
- **Coverage:** 95%
- **Setup:** ~15 min

---

## ✅ Best Practices Feb 2026

- [x] MANDATE 0.35: NO timeouts
- [x] MANDATE 0.36: DEQLHI-LOOP
- [x] TypeScript strict
- [x] TDD workflow
- [x] Git commit every change

---

<div align="center">

**Made with ❤️ by BIOMETRICS Team**

</div>
