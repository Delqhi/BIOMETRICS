# 🚀 BIOMETRICS ONBOARDING GUIDE

**Version:** 1.0 "Infinity Loop Edition"  
**Last Updated:** 2026-02-20  
**Status:** ✅ PRODUCTION-READY  
**Setup Time:** < 15 minutes

---

## 🎯 QUICK START (5 Minutes)

### Prerequisites

```bash
# Check system requirements
uname -a                    # macOS 13.0+ (Ventura or later)
node --version              # v20.0+
go version                  # Go 1.21+
docker --version            # Docker Desktop (optional)
```

### 1. Clone Repository

```bash
git clone https://github.com/Delqhi/BIOMETRICS.git
cd BIOMETRICS
```

### 2. Install Dependencies

```bash
# Node.js dependencies
npm install

# Go dependencies
cd biometrics-cli && go mod download && cd ..

# Python dependencies (optional)
pip install -r requirements.txt
```

### 3. Install OpenCode

```bash
npm install -g opencode
```

### 4. Authenticate Providers

```bash
# NVIDIA NIM (Required for Qwen 3.5 397B)
opencode auth add nvidia-nim
# Follow prompts to enter NVIDIA_API_KEY

# OpenCode ZEN (FREE models)
opencode auth add moonshot-ai
# Follow prompts to enter MOONSHOT_API_KEY
```

### 5. Verify Installation

```bash
# Check available models
opencode models | grep nvidia
opencode models | grep kimi

# Expected output:
# ✅ qwen/qwen3.5-397b-a17b (NVIDIA NIM)
# ✅ opencode/kimi-k2.5-free (OpenCode ZEN)
# ✅ opencode/minimax-m2.5-free (OpenCode ZEN)
```

### 6. Run First Agent

```bash
# Test with a simple task
opencode "List all files in current directory" --agent explore

# Expected: Agent responds with file listing
```

**✅ Congratulations! BIOMETRICS is ready!**

---

## 🔑 API KEY SETUP

### Required API Keys

| Provider | Key | Cost | Get From |
|----------|-----|------|----------|
| **NVIDIA NIM** | `NVIDIA_API_KEY` | ✅ FREE | https://build.nvidia.com |
| **Moonshot AI** | `MOONSHOT_API_KEY` | ✅ FREE | https://platform.moonshot.cn |

### Optional API Keys

| Provider | Key | Purpose | Cost |
|----------|-----|---------|------|
| **Supabase** | `SUPABASE_URL`, `SUPABASE_ANON_KEY` | Database | FREE tier |
| **Cloudflare** | `CLOUDFLARE_API_TOKEN` | Tunnel | FREE |
| **GitHub** | `GITHUB_TOKEN` | Git operations | FREE |

### Setting API Keys

#### Method 1: Environment Variables (Recommended)

```bash
# Add to ~/.zshrc or ~/.bashrc
export NVIDIA_API_KEY="nvapi-your-key-here"
export MOONSHOT_API_KEY="your-moonshot-key-here"

# Reload shell
source ~/.zshrc
```

#### Method 2: .env File

```bash
# Copy example
cp .env.example .env

# Edit .env
nano .env

# Add your keys:
NVIDIA_API_KEY=nvapi-your-key-here
MOONSHOT_API_KEY=your-moonshot-key-here
```

#### Method 3: OpenCode Auth (Automatic)

```bash
# OpenCode will prompt for keys
opencode auth add nvidia-nim
opencode auth add moonshot-ai
```

---

## 🤖 AGENT MODELS

### Available Models

| Model | Provider | Category | Max Parallel | Use Case |
|-------|----------|----------|--------------|----------|
| **Qwen 3.5 397B** | NVIDIA NIM | build, visual-engineering, writing | **1** | Main coding, docs |
| **Kimi K2.5** | OpenCode ZEN | deep | **1** | Heavy lifting, analysis |
| **MiniMax M2.5** | OpenCode ZEN | quick, explore | **1** | Quick tasks, configs |

### Model Assignment Rules

⚠️ **CRITICAL:** Never run 2 agents with the same model in parallel!

```bash
# ✅ CORRECT (3 different models):
opencode "Build REST API" --agent sisyphus &
opencode "Analyze architecture" --agent atlas &
opencode "Create config" --agent librarian &

# ❌ WRONG (same model):
opencode "Task 1" --agent sisyphus &
opencode "Task 2" --agent sisyphus &  # BLOCKED!
```

### Using Specific Models

```bash
# Use Qwen 3.5 397B (best for code)
opencode "Build feature" --model qwen/qwen3.5-397b-a17b

# Use Kimi K2.5 (best for analysis)
opencode "Deep research" --model opencode/kimi-k2.5-free

# Use MiniMax M2.5 (fastest)
opencode "Quick fix" --model opencode/minimax-m2.5-free
```

---

## 📚 ESSENTIAL COMMANDS

### Basic Commands

```bash
# Start agent with task
opencode "Your task description"

# Use specific agent
opencode "Build API" --agent sisyphus

# Use specific model
opencode "Code review" --model qwen/qwen3.5-397b-a17b

# List available agents
opencode agents

# List available models
opencode models
```

### Advanced Commands

```bash
# Run multiple agents in parallel (different models!)
opencode "Task 1" --agent sisyphus &
opencode "Task 2" --agent atlas &
opencode "Task 3" --agent librarian &
wait

# Check agent status
opencode status

# View session history
opencode sessions

# Save session
opencode session-save my-session

# Restore session
opencode session-restore my-session
```

### Orchestrator Commands

```bash
# Start orchestrator (24/7 mode)
./biometrics-cli/orchestrator start

# Stop orchestrator
./biometrics-cli/orchestrator stop

# Check orchestrator status
./biometrics-cli/orchestrator status

# View active sessions
./biometrics-cli/orchestrator sessions

# View model usage
./biometrics-cli/orchestrator models
```

---

## 🏗️ PROJECT STRUCTURE

```
BIOMETRICS/
├── biometrics-cli/          # Go CLI tool
│   ├── pkg/
│   │   ├── orchestrator/   # 24/7 Agent Loop
│   │   │   ├── orchestrator.go
│   │   │   ├── sicher_check.go
│   │   │   └── prompt_generator.go
│   │   ├── agents/         # Agent implementations
│   │   ├── auth/           # Authentication
│   │   └── ...
│   └── main.go
├── docs/                    # Documentation
│   ├── ORCHESTRATOR-MANDATE.md
│   ├── agents/
│   ├── architecture/
│   └── best-practices/
├── rules/                   # Rules and mandates
│   ├── global/
│   └── tools/
├── templates/               # Project templates
│   ├── global/
│   └── opencode/
├── archive/                 # Archived files
│   ├── sprint5-packages/
│   └── reports/
├── AGENTS-PLAN.md          # Infinity Loop tasks
├── AGENTS.md               # Global agent rules
├── ARCHITECTURE.md         # System architecture
├── README.md               # This file
└── .sisyphus/              # Session data
    ├── sessions/
    ├── prompts/
    ├── sicher-checks/
    └── boulder.json
```

---

## 🔥 INFINITY LOOP

### How It Works

```
START (20 Tasks)
  ↓
Task Complete → +5 New Tasks
  ↓
Sicher? Check → Verify work
  ↓
Git Commit → Save changes
  ↓
Next Task → Continue loop
  ↓
REPEAT FOREVER ♾️
```

### Task Lifecycle

1. **Spawn:** Orchestrator creates agent with massive prompt
2. **Execute:** Agent works on task (reads files first!)
3. **Verify:** Sicher? check validates work (6 checks)
4. **Commit:** Git commit with conventional message
5. **Generate:** +5 new tasks added to AGENTS-PLAN.md
6. **Repeat:** Next agent starts immediately

### Monitoring Progress

```bash
# View AGENTS-PLAN.md
cat AGENTS-PLAN.md

# Check active sessions
./biometrics-cli/orchestrator status

# View completed tasks
grep "✅" AGENTS-PLAN.md

# View Sicher? checks
ls -la .sisyphus/sicher-checks/
```

---

## 🛠️ TROUBLESHOOTING

### Common Issues

#### Issue 1: "Model not available"

**Solution:**
```bash
# Check authentication
opencode auth list

# Re-authenticate if needed
opencode auth add nvidia-nim
```

#### Issue 2: "Too many agents with same model"

**Solution:**
```bash
# Wait for running agents to complete
# Or use different models:
opencode "Task" --model opencode/kimi-k2.5-free
```

#### Issue 3: "Sicher? check failed"

**Solution:**
```bash
# View check results
cat .sisyphus/sicher-checks/ses_*.json

# Fix issues mentioned in report
# Re-run agent with corrected instructions
```

#### Issue 4: "Git commit failed"

**Solution:**
```bash
# Check git status
git status

# Configure git if needed
git config --global user.name "Your Name"
git config --global user.email "your@email.com"

# Commit manually
git add .
git commit -m "feat: your change"
```

### Getting Help

- 📚 **Documentation:** `docs/` directory
- 💬 **Discord:** https://discord.gg/biometrics
- 🐛 **Issues:** https://github.com/Delqhi/BIOMETRICS/issues
- 📧 **Email:** support@biometrics.dev

---

## 📊 NEXT STEPS

### After Onboarding

1. **Read AGENTS-PLAN.md** - Understand current tasks
2. **Read ORCHESTRATOR-MANDATE.md** - Learn orchestrator rules
3. **Run first task** - Test with simple command
4. **Monitor sessions** - Watch agents work
5. **Review Sicher? checks** - Verify quality

### Learning Path

| Level | Focus | Resources |
|-------|-------|-----------|
| **Beginner** | Basic commands, setup | This guide, README.md |
| **Intermediate** | Agent management, parallel tasks | docs/agents/ |
| **Advanced** | Orchestrator, custom agents | docs/orchestrator/ |
| **Expert** | Model fine-tuning, extensions | docs/advanced/ |

### Contributing

1. Fork repository
2. Create branch: `git checkout -b feature/your-feature`
3. Make changes
4. Run tests: `npm test`
5. Commit: `git commit -m "feat: your feature"`
6. Push: `git push origin feature/your-feature`
7. Create PR

---

## ✅ ONBOARDING CHECKLIST

Use this checklist to verify your setup:

```bash
# System Requirements
[ ] macOS 13.0+ (Ventura or later)
[ ] Node.js v20.0+
[ ] Go 1.21+
[ ] Git installed

# Installation
[ ] Repository cloned
[ ] npm install completed
[ ] go mod download completed
[ ] OpenCode installed

# Authentication
[ ] NVIDIA NIM authenticated
[ ] Moonshot AI authenticated
[ ] API keys in environment or .env

# Verification
[ ] opencode models shows all 3 models
[ ] First agent ran successfully
[ ] Git configured

# Understanding
[ ] Read AGENTS-PLAN.md
[ ] Read ORCHESTRATOR-MANDATE.md
[ ] Understand model assignment rules
[ ] Know how to run agents

# Ready to Start
[ ] Can spawn agents
[ ] Can monitor sessions
[ ] Can view Sicher? checks
[ ] Understand Infinity Loop
```

**All checked? You're ready! 🚀**

---

## 🎯 QUICK REFERENCE

### Essential Files

| File | Purpose |
|------|---------|
| `AGENTS-PLAN.md` | Current tasks and infinity loop |
| `AGENTS.md` | Global agent rules |
| `ARCHITECTURE.md` | System design |
| `docs/ORCHESTRATOR-MANDATE.md` | Orchestrator workflow |
| `docs/agents/AGENT-MODEL-MAPPING.md` | Model assignment |

### Essential Commands

```bash
# Start agent
opencode "task"

# Check status
opencode status

# View models
opencode models

# View sessions
opencode sessions

# Orchestrator
./biometrics-cli/orchestrator start
```

### Model Limits

| Model | Max Parallel | Category |
|-------|--------------|----------|
| Qwen 3.5 | 1 | build, writing |
| Kimi K2.5 | 1 | deep |
| MiniMax M2.5 | 1 | quick, explore |

**Remember:** MAX 3 agents total, all different models!

---

**Welcome to BIOMETRICS!** 🎉

**Infinity Loop awaits. Let's build the future together.** 🚀

*"Ein Task endet, fünf neue beginnen"*
