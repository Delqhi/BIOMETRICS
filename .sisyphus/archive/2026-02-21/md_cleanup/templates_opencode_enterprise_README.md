# Enterprise OpenCode Template

**Version:** 1.0.0  
**Type:** Enterprise Project Template  
**Purpose:** Large-scale Enterprise AI Projects with 26-Pillar Structure  
**Last Updated:** 2026-02-20

---

## Quick Start

```bash
# Clone template
cp -r templates/opencode/enterprise/ my-enterprise-project/
cd my-enterprise-project

# Install dependencies
npm install

# Configure environment
cp .env.example .env
# Edit .env with your credentials

# Start development
npm run dev
```

---

## Project Structure

```
enterprise/
├── AGENTS.md                 # AI Agent rules and policies
├── BLUEPRINT.md              # 22-Pillar enterprise architecture
├── README.md                # Document360 standard documentation
├── package.json              # Dependencies and scripts
├── opencode.json             # OpenCode configuration
├── tsconfig.json             # TypeScript strict config
├── docker-compose.yml        # Full stack orchestration
├── .env.example              # Environment template
├── .github/                  # CI/CD and templates
│   ├── workflows/
│   └── ISSUE_TEMPLATE/
├── src/                      # Source code
│   ├── agents/               # AI Agent implementations
│   ├── tasks/                # Task categories
│   ├── utils/                # Shared utilities
│   ├── config/               # Configuration files
│   ├── middleware/           # Middleware
│   └── services/             # Business logic
├── tests/                    # Test suites
│   ├── unit/
│   ├── integration/
│   └── e2e/
└── docs/                     # Documentation
    ├── non-dev/
    ├── dev/
    ├── project/
    └── postman/
```

---

## Features

- **26-Pillar Structure** - Complete enterprise documentation
- **Multi-Agent System** - Specialized AI agents for different tasks
- **Production Ready** - Security, monitoring, logging built-in
- **TypeScript Strict** - Maximum type safety
- **Docker Ready** - Full containerization
- **CI/CD Included** - GitHub Actions workflows

---

## Configuration

### Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
# Core
PROJECT_NAME=enterprise
ENVIRONMENT=production
LOG_LEVEL=info

# Database
DATABASE_URL=postgresql://user:pass@localhost:5432/db

# Redis
REDIS_URL=redis://localhost:6379

# API Keys
OPENAI_API_KEY=sk-xxx
ANTHROPIC_API_KEY=sk-ant-xxx
```

### OpenCode Configuration

Edit `opencode.json` to configure:
- AI Providers
- Fallback Chains
- MCP Servers
- Agent Models

---

## Available Scripts

| Script | Description |
|--------|-------------|
| `npm run dev` | Start development server |
| `npm run build` | Build for production |
| `npm run test` | Run all tests |
| `npm run test:unit` | Unit tests only |
| `npm run test:integration` | Integration tests |
| `npm run test:e2e` | End-to-end tests |
| `npm run lint` | Lint code |
| `npm run typecheck` | Type check |
| `npm run deploy` | Deploy to production |

---

## Support

- **Documentation:** `/docs/`
- **API Reference:** `/docs/dev/`
- **Troubleshooting:** `/docs/best-practices/`

---

*Enterprise Template v1.0.0 - Built for Scale*
