# 🛠️ BIOMETRICS CLI Commands

**Purpose:** Command reference for all CLI operations.

**Status:** ✅ Active  
**Last Updated:** 2026-02-19

---

## 📋 Command Categories

### Agent Management

| Command | Description |
|---------|-------------|
| `biometrics agent start [name]` | Start an agent |
| `biometrics agent stop [name]` | Stop an agent |
| `biometrics agent list` | List all agents |
| `biometrics agent status [name]` | Show agent status |

### Swarm Orchestration

| Command | Description |
|---------|-------------|
| `biometrics swarm start` | Start swarm (5+ agents) |
| `biometrics swarm stop` | Stop all agents |
| `biometrics swarm status` | Show swarm status |
| `biometrics swarm scale [count]` | Scale agent count |

### Project Management

| Command | Description |
|---------|-------------|
| `biometrics project create [name]` | Create new project |
| `biometrics project list` | List all projects |
| `biometrics project delete [name]` | Delete project |
| `biometrics project info [name]` | Show project details |

### Diagnostics

| Command | Description |
|---------|-------------|
| `biometrics doctor` | Full health check |
| `biometrics config show` | Show configuration |
| `biometrics logs [agent]` | View agent logs |
| `biometrics metrics` | Show performance metrics |

---

## 🚀 Usage Examples

### Start Agent Swarm

```bash
# Start 5 agents with different models
biometrics swarm start --agents 5 --models qwen,kimi,minimax
```

### Check Health

```bash
# Full diagnostics
biometrics doctor

# Output:
# ✅ OpenCode: Connected
# ✅ NVIDIA API: Active
# ✅ PostgreSQL: Running
# ✅ Redis: Connected
# ✅ Vault: Accessible
```

### Create Project

```bash
# Create new project with template
biometrics project create my-project --template enterprise
```

---

**For detailed usage:** See [docs/usage.md](../docs/usage.md)
