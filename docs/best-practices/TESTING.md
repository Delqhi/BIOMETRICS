# TESTING.md

## Übersicht

Dieses Dokument definiert die Teststrategie für das BIOMETRICS-Projekt. Es umfasst Unit-Tests, Integrationstests, E2E-Tests, Crashtests und Performance-Benchmarks.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

## Test-Strategie

### Test-Pyramide

```
        /\
       /E2E\
      /------\
     /Integration\
    /------------\
   /   Unit Tests \
  /----------------\
```

Die Test-Pyramide zeigt das Verhältnis der Testarten:

- **Unit Tests:** 70% - Schnelle, isolierte Tests für einzelne Funktionen
- **Integration Tests:** 20% - Tests für das Zusammenspiel von Komponenten
- **E2E Tests:** 10% - Vollständige User-Journey-Tests

### Test-Frameworks

| Bereich | Framework | Sprache |
|---------|-----------|---------|
| Unit | Vitest | TypeScript |
| Integration | Supertest | TypeScript |
| E2E | Playwright | TypeScript |
| CLI | Built-in | TypeScript |

## Unit Tests

### Struktur

```
biometrics-cli/
├── src/
│   ├── __tests__/
│   │   ├── commands/
│   │   ├── utils/
│   │   └── services/
│   └── ...
└── vitest.config.ts
```

### Beispiel: CLI Command Test

```typescript
import { describe, it, expect, beforeEach } from 'vitest';
import { execSync } from 'child_process';

describe('biometrics CLI', () => {
  beforeEach(() => {
    // Reset test environment
  });

  it('should show version', () => {
    const output = execSync('biometrics --version', { encoding: 'utf-8' });
    expect(output).toContain('1.0.0');
  });

  it('should run onboarding', () => {
    const output = execSync('biometrics --help', { encoding: 'utf-8' });
    expect(output).toContain('Usage:');
  });
});
```

### Beispiel: Utility Test

```typescript
import { describe, it, expect } from 'vitest';
import { parseConfig, validateEnvironment } from '../src/utils/config';

describe('config utilities', () => {
  it('should parse valid config', () => {
    const config = parseConfig({ apiKey: 'test', timeout: 5000 });
    expect(config.apiKey).toBe('test');
    expect(config.timeout).toBe(5000);
  });

  it('should throw on invalid config', () => {
    expect(() => parseConfig({})).toThrow('Missing required field: apiKey');
  });
});
```

## Integration Tests

### Supabase Integration

```typescript
import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import { createClient } from '@supabase/supabase-js';

const supabase = createClient(
  process.env.SUPABASE_URL!,
  process.env.SUPABASE_ANON_KEY!
);

describe('Supabase Integration', () => {
  it('should connect to database', async () => {
    const { data, error } = await supabase.from('users').select('count');
    expect(error).toBeNull();
    expect(data).toBeDefined();
  });

  it('should insert and retrieve user', async () => {
    const testUser = { email: 'test@example.com', name: 'Test' };
    
    const { data: inserted } = await supabase
      .from('users')
      .insert(testUser)
      .select()
      .single();
    
    expect(inserted.email).toBe(testUser.email);
    
    // Cleanup
    await supabase.from('users').delete().eq('id', inserted.id);
  });
});
```

### n8n Webhook Integration

```typescript
import { describe, it, expect } from 'vitest';

describe('n8n Webhook Integration', () => {
  const WEBHOOK_URL = process.env.N8N_WEBHOOK_URL;

  it('should trigger workflow', async () => {
    const response = await fetch(WEBHOOK_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ test: true })
    });

    expect(response.ok).toBe(true);
  });

  it('should handle workflow errors', async () => {
    const response = await fetch(WEBHOOK_URL, {
      method: 'POST',
      body: JSON.stringify({ invalid: 'data' })
    });

    expect(response.status).toBeGreaterThanOrEqual(400);
  });
});
```

## E2E Tests

### Playwright Konfiguration

```typescript
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  timeout: 30000,
  retries: 2,
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
  ],
});
```

### Beispiel: User Journey Test

```typescript
import { test, expect } from '@playwright/test';

test('complete user onboarding journey', async ({ page }) => {
  // 1. Visit landing page
  await page.goto('/');
  await expect(page.locator('h1')).toContainText('BIOMETRICS');

  // 2. Click CTA
  await page.click('text=Get Started');
  await expect(page).toHaveURL('/onboarding');

  // 3. Fill form
  await page.fill('[name="email"]', 'test@example.com');
  await page.fill('[name="name"]', 'Test User');
  await page.click('button[type="submit"]');

  // 4. Verify dashboard
  await expect(page).toHaveURL('/dashboard');
  await expect(page.locator('.welcome')).toContainText('Test User');
});
```

## Crashtests

Crashtests sind kritische Tests, die sicherstellen, dass das System auch unter extremen Bedingungen stabil bleibt.

### Test 1: CLI Installation

```bash
# Test: biometrics CLI Installation
biometrics --version

# Erwartet: Versionsnummer wird ausgegeben
# Exit Code: 0
```

### Test 2: API Health Check

```bash
# Test: API Health Endpoint
curl -f http://localhost:8080/health

# Erwartet: {"status":"ok"}
# Exit Code: 0
```

### Test 3: Qwen 3.5 Integration Test

```bash
# Test: OpenCode Model Integration
opencode --model nvidia/qwen/qwen3.5-397b-a17b "Test"

# Erwartet: Gültige Antwort
# Exit Code: 0
# Timeout: 120s
```

### Test 4: Datenbank-Verbindung

```bash
# Test: Supabase Connection
psql "$DATABASE_URL" -c "SELECT 1;"

# Erwartet: 1 row
# Exit Code: 0
```

### Test 5: n8n Availability

```bash
# Test: n8n Web UI
curl -f http://localhost:5678

# Erwartet: HTML Response
# Exit Code: 0
```

### Test 6: Rate Limit Handling

```bash
# Test: Rate Limit (40 RPM bei NVIDIA NIM)
for i in {1..50}; do
  curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
       https://integrate.api.nvidia.com/v1/models
done

# Erwartet: Max 40 erfolgreich, dann 429
```

### Test 7: Timeout Handling

```bash
# Test: Request Timeout
timeout 5 curl http://localhost:8080/slow-endpoint

# Erwartet: Timeout nach 5s
# Exit Code: 124 (timeout)
```

### Test 8: Memory Leak Detection

```bash
# Test: Memory Usage
for i in {1..100}; do
  biometrics process-data --input data.json
done

# Erwartet: Stabiler Memory-Verbrauch
# Tools: top, htop
```

### Test 9: Concurrent Requests

```bash
# Test: Parallel Load
for i in {1..10}; do
  curl -f http://localhost:8080/api/health &
done
wait

# Erwartet: Alle Requests erfolgreich
```

### Test 10: Invalid Input Handling

```bash
# Test: Invalid JSON
echo "invalid json" | biometrics parse

# Erwartet: Fehlermeldung
# Exit Code: != 0

# Test: Missing Required Fields
biometrics create-user --name "Test"

# Erwartet: Validation Error
# Exit Code: != 0
```

## Performance Benchmarks

### Latenz-Anforderungen

| Operation | Ziel | Maximum |
|-----------|------|---------|
| CLI Startup | < 1s | 3s |
| API Health | < 100ms | 500ms |
| Database Query | < 200ms | 1s |
| Qwen 3.5 Inference | < 90s | 120s |

### Load Testing

```bash
# Test: 100 concurrent requests
autocannon -c 100 -d 10 http://localhost:8080/health

# Erwartet: < 5s für alle Requests
```

### Benchmark-Skript

```typescript
import { performance } from 'perf_hooks';

async function benchmark(fn: () => Promise<void>, name: string) {
  const start = performance.now();
  await fn();
  const duration = performance.now() - start;
  console.log(`${name}: ${duration}ms`);
  return duration;
}

// Usage
await benchmark(() => cli.parse(['--version']), 'CLI Parse');
await benchmark(() => api.health(), 'API Health');
```

## CI/CD Integration

### GitHub Actions Workflow

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          
      - name: Install dependencies
        run: pnpm install
        
      - name: Run Unit Tests
        run: pnpm test:unit
        
      - name: Run Integration Tests
        run: pnpm test:integration
        env:
          SUPABASE_URL: ${{ secrets.SUPABASE_URL }}
          SUPABASE_KEY: ${{ secrets.SUPABASE_KEY }}
          
      - name: Run E2E Tests
        run: pnpm test:e2e
        
      - name: Upload Coverage
        uses: codecov/codecov-action@v4
```

### Test Script Commands

```json
{
  "scripts": {
    "test": "vitest",
    "test:unit": "vitest run",
    "test:integration": "vitest run --config vitest.integration.config.ts",
    "test:e2e": "playwright test",
    "test:coverage": "vitest run --coverage",
    "test:crashtest": "./scripts/crashtest.sh"
  }
}
```

## Test-Berichterstattung

### Coverage-Anforderungen

| Bereich | Minimum | Ziel |
|---------|---------|------|
| Unit Tests | 60% | 80% |
| Critical Paths | 90% | 95% |

### Test-Reports

```bash
# HTML Coverage Report
pnpm test:coverage

# JUnit XML (CI/CD)
pnpm test:unit -- --reporter=junit > test-results.xml

# Allure Report
allure serve test-results/
```

## Abnahme-Check TESTING

1. [ ] Unit Tests für alle Module vorhanden
2. [ ] Integration Tests für externe Services
3. [ ] E2E Tests für kritische User-Journeys
4. [ ] Crashtests dokumentiert und ausführbar
5. [ ] Performance Benchmarks definiert
6. [ ] CI/CD Pipeline integriert
7. [ ] Coverage-Reporting aktiv

---

**Version:** 1.0  
**Letzte Aktualisierung:** 2026-02-18  
**Status:** ACTIVE
