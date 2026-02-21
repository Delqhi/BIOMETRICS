# 🤝 Contributing to BIOMETRICS

<div align="center">

**Welcome to the BIOMETRICS Community!** 🚀

[![Stars](https://img.shields.io/github/stars/Delqhi/BIOMETRICS?style=social)](https://github.com/Delqhi/BIOMETRICS/stargazers)
[![Forks](https://img.shields.io/github/forks/Delqhi/BIOMETRICS?style=social)](https://github.com/Delqhi/BIOMETRICS/network/members)
[![Issues](https://img.shields.io/github/issues/Delqhi/BIOMETRICS?logo=github)](https://github.com/Delqhi/BIOMETRICS/issues)
[![Discord](https://img.shields.io/discord/BIOMETRICS?logo=discord)](https://discord.gg/biometrics)

[Quick Start](#-quick-start) • [Development Setup](#-development-setup) • [Code Style](#-code-style) • [Git Workflow](#-git-workflow) • [PR Guide](#-pull-request-process) • [Good First Issues](#-good-first-issues)

</div>

---

## 📖 Table of Contents

<!-- TOC start (generated with https://github.com/derlin/bitdowntoc) -->

- [🎯 Welcome](#-welcome)
- [⚡ Quick Start](#-quick-start)
- [🛠️ Development Setup](#️-development-setup)
  - [Prerequisites](#prerequisites)
  - [Fork & Clone](#fork--clone)
  - [Install Dependencies](#install-dependencies)
  - [Environment Configuration](#environment-configuration)
  - [Verify Setup](#verify-setup)
- [📋 Code Style Guidelines](#-code-style-guidelines)
  - [Go Backend](#go-backend)
  - [TypeScript Frontend](#typescript-frontend)
  - [Markdown Documentation](#markdown-documentation)
  - [YAML Configuration](#yaml-configuration)
  - [Commit Messages](#commit-messages)
- [🧪 Testing Requirements](#-testing-requirements)
  - [Go Tests](#go-tests)
  - [TypeScript Tests](#typescript-tests)
  - [Coverage Requirements](#coverage-requirements)
  - [Running Tests](#running-tests)
- [🔧 Git Workflow](#-git-workflow)
  - [Branch Naming](#branch-naming)
  - [Commit Guidelines](#commit-guidelines)
  - [Visual Workflow](#visual-workflow)
- [📝 Pull Request Process](#-pull-request-process)
  - [Before Submitting](#before-submitting)
  - [PR Checklist](#pr-checklist)
  - [Review Process](#review-process)
- [🐛 Issue Guidelines](#-issue-guidelines)
  - [Bug Reports](#bug-reports)
  - [Feature Requests](#feature-requests)
- [📚 Documentation](#-documentation)
  - [Documentation Standards](#documentation-standards)
  - [Writing Guides](#writing-guides)
- [🏗️ Architecture Overview](#️-architecture-overview)
- [🎯 Good First Issues](#-good-first-issues)
- [💬 Getting Help](#-getting-help)
- [🎖️ Contributor Recognition](#️-contributor-recognition)

<!-- TOC end -->

---

## 🎯 Welcome

**Thank you for your interest in contributing to BIOMETRICS!** 

BIOMETRICS is an enterprise-grade AI agent orchestration framework built on **33 core mandates**. We welcome contributions from developers of all skill levels, from first-time contributors to seasoned open-source veterans.

### 🌟 Why Contribute?

- **Learn cutting-edge AI orchestration** - Work with state-of-the-art agent swarms
- **Enterprise-quality code** - Follow best practices used by Fortune 500 companies
- **Active community** - Join 100+ developers building the future of AI development
- **Real-world impact** - Your code will be used in production systems
- **Career growth** - Build your portfolio with high-quality open-source work

### 🎯 What You Can Contribute

| Contribution Type | Examples | Difficulty | Time |
|------------------|----------|------------|------|
| **🐛 Bug Fixes** | Fix typos, resolve errors, patch security issues | ⭐ Easy | 1-2 hours |
| **✨ Features** | New CLI commands, API endpoints, UI components | ⭐⭐ Medium | 4-8 hours |
| **📚 Documentation** | Tutorials, guides, API docs, examples | ⭐ Easy | 2-4 hours |
| **🧪 Tests** | Unit tests, integration tests, E2E tests | ⭐⭐ Medium | 3-5 hours |
| **♻️ Refactoring** | Code improvements, performance optimization | ⭐⭐⭐ Hard | 8+ hours |
| **🎨 Design** | UI/UX improvements, diagrams, logos | ⭐⭐ Medium | 4-6 hours |

### 🚀 Quick Navigation

```bash
# New to the project?
→ Read: docs/setup/COMPLETE-SETUP.md

# Ready to code?
→ Jump to: #development-setup

# Looking for tasks?
→ Check: #good-first-issues

# Need help?
→ Ask: #getting-help
```

---

## ⚡ Quick Start

**Get up and running in under 30 minutes!** Follow these steps to make your first contribution.

### Step-by-Step Guide (30 Minutes)

```mermaid
graph LR
    A[Fork Repository] --> B[Clone Fork]
    B --> C[Install Dependencies]
    C --> D[Create Branch]
    D --> E[Make Changes]
    E --> F[Run Tests]
    F --> G[Commit Changes]
    G --> H[Create PR]
    
    style A fill:#3498db,stroke:#2980b9,stroke-width:2px,color:#fff
    style H fill:#2ecc71,stroke:#27ae60,stroke-width:2px,color:#fff
```

#### 1. Fork the Repository (2 minutes)

1. Click the **Fork** button in the top-right corner of GitHub
2. Wait for the fork to complete
3. You'll be redirected to `https://github.com/YOUR_USERNAME/BIOMETRICS`

#### 2. Clone Your Fork (3 minutes)

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/BIOMETRICS.git
cd BIOMETRICS

# Add upstream remote (to sync with main repo)
git remote add upstream https://github.com/Delqhi/BIOMETRICS.git

# Verify remotes
git remote -v
# Expected output:
# origin    https://github.com/YOUR_USERNAME/BIOMETRICS.git (fetch)
# origin    https://github.com/YOUR_USERNAME/BIOMETRICS.git (push)
# upstream  https://github.com/Delqhi/BIOMETRICS.git (fetch)
# upstream  https://github.com/Delqhi/BIOMETRICS.git (push)
```

#### 3. Install Dependencies (10 minutes)

```bash
# Install pnpm (if not already installed)
curl -fsSL https://get.pnpm.io/install.sh | sh -

# Install project dependencies
pnpm install

# Install Go tools (for biometrics-cli)
cd biometrics-cli
go mod download
cd ..

# Verify installation
pnpm run doctor
```

#### 4. Create a Branch (2 minutes)

```bash
# Sync with upstream
git fetch upstream
git checkout main
git merge upstream/main

# Create feature branch
git checkout -b feat/your-feature-name
# Format: feat|fix|docs|refactor|test|chore/description
```

#### 5. Make Changes & Test (10+ minutes)

```bash
# Make your changes in your favorite editor
code .  # VS Code
# or
vim .   # Terminal editor

# Run tests before committing
pnpm test

# Check code formatting
pnpm format
```

#### 6. Commit & Push (3 minutes)

```bash
# Stage changes
git add .

# Commit with conventional commit message
git commit -m "feat: add your feature description"

# Push to your fork
git push origin feat/your-feature-name
```

#### 7. Create Pull Request (5 minutes)

1. Go to `https://github.com/YOUR_USERNAME/BIOMETRICS`
2. Click **"Compare & pull request"**
3. Fill out the PR template
4. Submit for review!

**🎉 Congratulations!** You've submitted your first contribution!

---

## 🛠️ Development Setup

### Prerequisites

Before you start, ensure you have the following installed:

| Tool | Version | Required For | Install Command |
|------|---------|--------------|-----------------|
| **Git** | 2.30+ | Version control | `brew install git` |
| **Node.js** | 20.x+ | Frontend/Tooling | `brew install node@20` |
| **pnpm** | 8.x+ | Package management | `curl -fsSL https://get.pnpm.io/install.sh \| sh -` |
| **Go** | 1.24+ | CLI development | `brew install go` |
| **Python** | 3.10+ | Scripts & automation | `brew install python@3.10` |
| **Docker** | 24.x+ | Containerization | `brew install --cask docker` |

#### Verify Prerequisites

```bash
# Run this script to verify all prerequisites
./scripts/verify-prerequisites.sh

# Expected output:
# ✅ Git 2.40.0
# ✅ Node.js v20.10.0
# ✅ pnpm 8.12.0
# ✅ Go 1.24.3
# ✅ Python 3.10.13
# ✅ Docker 24.0.7
```

### Fork & Clone

#### Detailed Fork Instructions

<div align="center">

![Fork Repository](https://img.shields.io/badge/Step-1-blue?style=for-the-badge)

**Creating Your Fork**

</div>

1. **Navigate to the repository**: Go to [https://github.com/Delqhi/BIOMETRICS](https://github.com/Delqhi/BIOMETRICS)

2. **Click Fork button**: Top-right corner of the GitHub page

3. **Configure fork**:
   - Select your GitHub account
   - Keep "Copy the `main` branch only" checked
   - Click **"Create fork"**

4. **Wait for completion**: GitHub will create your fork (usually takes 5-10 seconds)

<div align="center">

![Clone Repository](https://img.shields.io/badge/Step-2-green?style=for-the-badge)

**Cloning to Your Machine**

</div>

```bash
# Replace YOUR_USERNAME with your GitHub username
git clone https://github.com/YOUR_USERNAME/BIOMETRICS.git
cd BIOMETRICS

# Configure upstream remote (CRITICAL for syncing)
git remote add upstream https://github.com/Delqhi/BIOMETRICS.git

# Verify configuration
git remote -v
```

**Expected Output:**
```
origin    https://github.com/YOUR_USERNAME/BIOMETRICS.git (fetch)
origin    https://github.com/YOUR_USERNAME/BIOMETRICS.git (push)
upstream  https://github.com/Delqhi/BIOMETRICS.git (fetch)
upstream  https://github.com/Delqhi/BIOMETRICS.git (push)
```

### Install Dependencies

#### Frontend Dependencies

```bash
# Install pnpm (if not installed)
curl -fsSL https://get.pnpm.io/install.sh | sh -

# Install all dependencies
pnpm install

# What gets installed:
# - React/Next.js packages
# - Development tools (ESLint, Prettier)
# - Testing libraries (Jest, Testing Library)
# - Build tools (Webpack, Babel)
```

#### Go Dependencies

```bash
# Navigate to CLI directory
cd biometrics-cli

# Download Go modules
go mod download

# Install development tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Verify installation
go version
# Expected: go version go1.24.3 darwin/amd64
```

#### Python Dependencies (Optional)

```bash
# Create virtual environment
python3 -m venv .venv
source .venv/bin/activate  # On macOS/Linux
# or
.venv\Scripts\activate  # On Windows

# Install Python packages
pip install -r requirements.txt
```

### Environment Configuration

#### Create .env File

```bash
# Copy example environment file
cp .env.example .env

# Edit with your configuration
nano .env  # or use your preferred editor
```

#### Required Environment Variables

```bash
# .env - Required Variables

# NVIDIA API Key (Required for AI Agents)
NVIDIA_API_KEY="nvapi-YOUR_API_KEY_HERE"

# OpenCode Authentication
OPENCODE_API_KEY="your-opencode-key"

# Database Configuration (Development)
DATABASE_URL="postgresql://postgres:password@localhost:5432/biometrics_dev"

# Redis Cache
REDIS_URL="redis://localhost:6379"

# JWT Secret (Development only - use strong random string)
JWT_SECRET="your-super-secret-jwt-key-min-32-chars"

# Optional: Sentry Error Tracking
SENTRY_DSN="https://your-sentry-dsn@sentry.io/project-id"
```

#### Getting API Keys

1. **NVIDIA API Key**:
   - Visit [https://build.nvidia.com](https://build.nvidia.com)
   - Sign up for free account
   - Create API key in dashboard
   - Free tier: 40 RPM, sufficient for development

2. **OpenCode Key**:
   - Run `opencode auth login`
   - Follow authentication flow
   - Key stored automatically

### Verify Setup

```bash
# Run comprehensive setup verification
pnpm run doctor

# Expected successful output:
# ✅ Node.js v20.10.0
# ✅ pnpm 8.12.0
# ✅ Go 1.24.3
# ✅ All dependencies installed
# ✅ Environment variables configured
# ✅ Database connection successful
# ✅ Redis connection successful
# ✅ Ready for development!
```

#### Troubleshooting Setup Issues

| Issue | Solution |
|-------|----------|
| `pnpm: command not found` | Run `curl -fsSL https://get.pnpm.io/install.sh \| sh -` then restart terminal |
| `go mod download` fails | Check internet connection, try `go clean -modcache` |
| Database connection error | Ensure PostgreSQL is running: `brew services start postgresql` |
| Port already in use | Change port in .env or stop conflicting service |

---
