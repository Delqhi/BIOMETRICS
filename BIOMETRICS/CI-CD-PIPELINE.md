# CI/CD Pipeline Documentation

**Project:** BIOMETRICS  
**Last Updated:** 2026-02-18  
**Maintainer:** DevOps Team

---

## Overview

This document describes the CI/CD pipeline infrastructure for the BIOMETRICS project. The pipeline automates build, test, and deployment processes using GitHub Actions.

## Pipeline Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         CI/CD PIPELINE ARCHITECTURE                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐           │
│  │  PUSH    │───▶│  BUILD   │───▶│  TEST    │───▶│ DEPLOY   │           │
│  │  (main)  │    │  (npm)   │    │ (Jest)   │    │ (Vercel) │           │
│  └──────────┘    └──────────┘    └──────────┘    └──────────┘           │
│                      │              │                                    │
│                      ▼              ▼                                    │
│                ┌──────────┐   ┌──────────┐                             │
│                │  LINT    │   │  E2E     │                             │
│                │(ESLint)  │   │(Playwright│                             │
│                └──────────┘   └──────────┘                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## GitHub Actions Workflow

### Main Workflow File

Create `.github/workflows/ci.yml`:

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

env:
  NODE_VERSION: '20.x'
  NPM_TOKEN: ${{ secrets.NPM_TOKEN }}

jobs:
  lint:
    name: Lint Code
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'

      - name: Install dependencies
        run: npm ci

      - name: Run ESLint
        run: npm run lint
        continue-on-error: true

      - name: Run Prettier check
        run: npm run format:check

  typecheck:
    name: Type Check
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'

      - name: Install dependencies
        run: npm ci

      - name: Run TypeScript compiler
        run: npm run typecheck

  test:
    name: Unit Tests
    runs-on: ubuntu-latest
    needs: [lint, typecheck]
    strategy:
      matrix:
        node-version: [18.x, 20.x, 22.x]
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Setup Node.js ${{ matrix.node-version }}
        uses: actions/setup-node@v4
        with:
          node-version: ${{ matrix.node-version }}
          cache: 'npm'

      - name: Install dependencies
        run: npm ci

      - name: Run tests
        run: npm run test:coverage
        env:
          CI: true
          DATABASE_URL: postgresql://test:test@localhost:5432/biometrics_test
          REDIS_URL: redis://localhost:6379

      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v4
        with:
          token: ${{ secrets.CODECOV_TOKEN }}
          files: ./coverage/lcov.info
          flags: unittests

  e2e:
    name: End-to-End Tests
    runs-on: ubuntu-latest
    needs: [test]
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'

      - name: Install dependencies
        run: npm ci

      - name: Install Playwright browsers
        run: npx playwright install --with-deps

      - name: Start application
        run: npm run start &
        env:
          DATABASE_URL: postgresql://test:test@localhost:5432/biometrics_test
          REDIS_URL: redis://localhost:6379

      - name: Wait for application
        run: sleep 10 && curl -f http://localhost:3000/health || exit 1

      - name: Run Playwright tests
        run: npm run test:e2e
        env:
          PLAYWRIGHT_BASE_URL: http://localhost:3000

      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: playwright-results
          path: test-results/

  build:
    name: Build Application
    runs-on: ubuntu-latest
    needs: [e2e]
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'

      - name: Install dependencies
        run: npm ci

      - name: Build application
        run: npm run build

      - name: Upload build artifacts
        uses: actions/upload-artifact@v4
        with:
          name: build-output
          path: dist/

  deploy:
    name: Deploy to Production
    runs-on: ubuntu-latest
    needs: [build]
    if: github.ref == 'refs/heads/main'
    environment:
      name: production
      url: https://biometrics.example.com
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Download build artifacts
        uses: actions/download-artifact@v4
        with:
          name: build-output
          path: dist/

      - name: Deploy to Vercel
        uses: amondnet/vercel-action@v25
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}
          vercel-org-id: ${{ secrets.VERCEL_ORG_ID }}
          vercel-project-id: ${{ secrets.VERCEL_PROJECT_ID }}
          vercel-args: '--prod'

      - name: Run smoke tests
        run: |
          sleep 5
          curl -f https://biometrics.example.com/health || exit 1

  notify:
    name: Notify on Failure
    runs-on: ubuntu-latest
    needs: [deploy]
    if: failure()
    steps:
      - name: Send Slack notification
        uses: 8398a7/action-slack@v3
        with:
          status: failure
          fields: repo,message,commit,author
        env:
          SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}
```

## Environment-Specific Workflows

### Staging Deployment

Create `.github/workflows/staging.yml`:

```yaml
name: Deploy to Staging

on:
  push:
    branches: [develop]

jobs:
  deploy-staging:
    name: Deploy Staging
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20.x'
          cache: 'npm'

      - name: Install dependencies
        run: npm ci

      - name: Build application
        run: npm run build

      - name: Deploy to Vercel (Staging)
        uses: amondnet/vercel-action@v25
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}
          vercel-org-id: ${{ secrets.VERCEL_ORG_ID }}
          vercel-project-id: ${{ secrets.VERCEL_PROJECT_ID }}
          vercel-args: '--environment=preview'

      - name: Comment deployment URL
        uses: actions/github-script@v7
        with:
          script: |
            const url = '${{ steps.deploy.outputs.url }}';
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: `🚀 Deployed to staging: ${url}`
            })
```

### Security Scanning

Create `.github/workflows/security.yml`:

```yaml
name: Security Scan

on:
  push:
    branches: [main, develop]
  schedule:
    - cron: '0 0 * * 0'  # Weekly scan

jobs:
  dependency-scan:
    name: Dependency Vulnerability Scan
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Run npm audit
        run: npm audit --audit-level=high
        continue-on-error: true

      - name: Run Snyk vulnerability scanner
        uses: snyk/actions/node@master
        env:
          SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}

  codeql:
    name: CodeQL Analysis
    runs-on: ubuntu-latest
    permissions:
      actions: read
      contents: read
      security-events: write
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Initialize CodeQL
        uses: github/codeql-action/init@v3
        with:
          languages: javascript,typescript

      - name: Build
        run: npm ci && npm run build

      - name: Perform CodeQL analysis
        uses: github/codeql-action/analyze@v3

  secrets-scan:
    name: Secrets Detection
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Run GitLeaks
        uses: zricethezav/gitleaks-action@master
        continue-on-error: true
```

## Pipeline Configuration

### Workflow Dispatch Parameters

```yaml
on:
  workflow_dispatch:
    inputs:
      environment:
        description: 'Deployment environment'
        required: true
        default: 'staging'
        type: choice
        options:
          - staging
          - production
      skip_tests:
        description: 'Skip test execution'
        required: false
        type: boolean
        default: false
```

### Concurrency Control

```yaml
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true
```

### Caching Strategy

```yaml
- name: Cache Node modules
  uses: actions/cache@v4
  with:
    path: ~/.npm
    key: ${{ runner.os }}-npm-${{ hashFiles('**/package-lock.json') }}
    restore-keys: |
      ${{ runner.os }}-npm-

- name: Cache Playwright browsers
  uses: actions/cache@v4
  with:
    path: ~/.cache/ms-playwright
    key: ${{ runner.os }}-playwright-${{ hashFiles('**/package-lock.json') }}
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `NODE_VERSION` | Node.js version | Yes |
| `DATABASE_URL` | PostgreSQL connection string | Yes |
| `REDIS_URL` | Redis connection string | Yes |
| `VERCEL_TOKEN` | Vercel API token | Yes (deploy) |
| `CODECOV_TOKEN` | Codecov upload token | No |
| `SLACK_WEBHOOK_URL` | Slack notification URL | No |

## Troubleshooting

### Common Issues

1. **Build fails with memory error**
   ```bash
   # Increase Node.js memory limit
   NODE_OPTIONS="--max-old-space-size=4096" npm run build
   ```

2. **Playwright tests timeout**
   - Increase timeout in `playwright.config.ts`
   - Check browser installation: `npx playwright install`

3. **Cache miss on every run**
   - Verify `package-lock.json` is committed
   - Check cache key includes hash of lock file

## Related Documentation

- [TESTING-SUITE.md](./TESTING-SUITE.md)
- [MONITORING-SETUP.md](./MONITORING-SETUP.md)
- [BACKUP-STRATEGY.md](./BACKUP-STRATEGY.md)

---

**End of CI/CD Pipeline Documentation**
