# Testing Suite Documentation

**Project:** BIOMETRICS  
**Last Updated:** 2026-02-18  
**Maintainer:** QA Team

---

## Overview

This document describes the testing infrastructure for BIOMETRICS, including unit tests, integration tests, and end-to-end tests using Jest, Vitest, and Playwright.

## Testing Stack

| Tool | Purpose | Version |
|------|---------|---------|
| **Jest** | Unit & Integration Tests | 29.x |
| **Vitest** | Fast unit testing (alternative) | 1.x |
| **Playwright** | E2E Browser Testing | 1.x |
| **MSW** | API Mocking | 2.x |
| **Faker.js** | Test Data Generation | 7.x |

## Project Configuration

### Jest Configuration

Create `jest.config.js`:

```javascript
module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  roots: ['<rootDir>/src', '<rootDir>/tests'],
  testMatch: ['**/__tests__/**/*.ts', '**/?(*.)+(spec|test).ts'],
  transform: {
    '^.+\\.ts$': ['ts-jest', {
      tsconfig: {
        jsx: 'react',
        esModuleInterop: true,
      },
    }],
  },
  collectCoverageFrom: [
    'src/**/*.ts',
    '!src/**/*.d.ts',
    '!src/**/index.ts',
  ],
  coverageDirectory: 'coverage',
  coverageReporters: ['text', 'lcov', 'html'],
  coverageThreshold: {
    global: {
      branches: 80,
      functions: 80,
      lines: 80,
      statements: 80,
    },
  },
  setupFilesAfterEnv: ['<rootDir>/tests/setup.ts'],
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/src/$1',
    '\\.(css|less|scss|sass)$': 'identity-obj-proxy',
  },
  testTimeout: 10000,
  verbose: true,
};
```

### Vitest Configuration

Create `vitest.config.ts`:

```typescript
import { defineConfig } from 'vitest/config';
import path from 'path';

export default defineConfig({
  test: {
    globals: true,
    environment: 'node',
    include: ['tests/**/*.test.ts'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'tests/',
        '**/*.config.*',
        '**/*.d.ts',
      ],
    },
    setupFiles: ['tests/setup.ts'],
    testTimeout: 10000,
    hookTimeout: 30000,
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
});
```

### Playwright Configuration

Create `playwright.config.ts`:

```typescript
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [
    ['html', { outputFolder: 'playwright-report' }],
    ['json', { outputFile: 'test-results/results.json' }],
    ['list'],
  ],
  use: {
    baseURL: process.env.PLAYWRIGHT_BASE_URL || 'http://localhost:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
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
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] },
    },
  ],
  webServer: {
    command: 'npm run start',
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
    timeout: 120000,
  },
});
```

## Unit Testing

### Basic Test Example

Create `tests/unit/user.service.test.ts`:

```typescript
import { UserService } from '@/services/user.service';
import { jest } from '@jest/globals';

// Mock dependencies
const mockUserRepository = {
  findById: jest.fn(),
  findByEmail: jest.fn(),
  create: jest.fn(),
  update: jest.fn(),
  delete: jest.fn(),
};

const mockEmailService = {
  sendWelcomeEmail: jest.fn(),
  sendPasswordReset: jest.fn(),
};

describe('UserService', () => {
  let userService: UserService;

  beforeEach(() => {
    jest.clearAllMocks();
    userService = new UserService(
      mockUserRepository as any,
      mockEmailService as any
    );
  });

  describe('createUser', () => {
    it('should create a new user with valid data', async () => {
      // Arrange
      const userData = {
        email: 'test@example.com',
        password: 'securePassword123',
        firstName: 'John',
        lastName: 'Doe',
      };

      const expectedUser = {
        id: '1',
        ...userData,
        createdAt: new Date(),
      };

      mockUserRepository.findByEmail.mockResolvedValue(null);
      mockUserRepository.create.mockResolvedValue(expectedUser);
      mockEmailService.sendWelcomeEmail.mockResolvedValue(true);

      // Act
      const result = await userService.createUser(userData);

      // Assert
      expect(result).toEqual(expectedUser);
      expect(mockUserRepository.findByEmail).toHaveBeenCalledWith(userData.email);
      expect(mockUserRepository.create).toHaveBeenCalledWith(
        expect.objectContaining({
          email: userData.email,
          firstName: userData.firstName,
        })
      );
      expect(mockEmailService.sendWelcomeEmail).toHaveBeenCalledWith(expectedUser);
    });

    it('should throw error if email already exists', async () => {
      // Arrange
      const userData = {
        email: 'existing@example.com',
        password: 'securePassword123',
      };

      mockUserRepository.findByEmail.mockResolvedValue({
        id: '1',
        email: userData.email,
      });

      // Act & Assert
      await expect(userService.createUser(userData)).rejects.toThrow(
        'User with this email already exists'
      );
    });

    it('should hash password before saving', async () => {
      // Arrange
      const userData = {
        email: 'test@example.com',
        password: 'plainPassword',
      };

      mockUserRepository.findByEmail.mockResolvedValue(null);
      mockUserRepository.create.mockImplementation((data) => 
        Promise.resolve({ id: '1', ...data })
      );

      // Act
      await userService.createUser(userData);

      // Assert
      const createCall = mockUserRepository.create.mock.calls[0][0];
      expect(createCall.password).not.toBe(userData.password);
      expect(createCall.password).toMatch(/^\$2[aby]\$\d+\$/);
    });
  });

  describe('getUserById', () => {
    it('should return user if found', async () => {
      // Arrange
      const userId = '1';
      const expectedUser = { id: userId, email: 'test@example.com' };
      mockUserRepository.findById.mockResolvedValue(expectedUser);

      // Act
      const result = await userService.getUserById(userId);

      // Assert
      expect(result).toEqual(expectedUser);
      expect(mockUserRepository.findById).toHaveBeenCalledWith(userId);
    });

    it('should return null if user not found', async () => {
      // Arrange
      mockUserRepository.findById.mockResolvedValue(null);

      // Act
      const result = await userService.getUserById('nonexistent');

      // Assert
      expect(result).toBeNull();
    });
  });
});
```

## Integration Testing

### API Integration Test

Create `tests/integration/api/auth.test.ts`:

```typescript
import request from 'supertest';
import { app } from '@/app';
import { prisma } from '@/lib/prisma';
import { jest } from '@jest/globals';

describe('Authentication API', () => {
  beforeAll(async () => {
    // Setup test database
    await prisma.user.deleteMany();
  });

  afterAll(async () => {
    // Cleanup
    await prisma.$disconnect();
  });

  beforeEach(async () => {
    jest.clearAllMocks();
  });

  describe('POST /api/auth/register', () => {
    it('should register a new user successfully', async () => {
      // Arrange
      const userData = {
        email: 'newuser@example.com',
        password: 'SecurePass123!',
        firstName: 'Jane',
        lastName: 'Smith',
      };

      // Act
      const response = await request(app)
        .post('/api/auth/register')
        .send(userData);

      // Assert
      expect(response.status).toBe(201);
      expect(response.body).toHaveProperty('token');
      expect(response.body.user).toHaveProperty('email', userData.email);
    });

    it('should return 400 for invalid email', async () => {
      // Arrange
      const invalidData = {
        email: 'not-an-email',
        password: 'password123',
      };

      // Act
      const response = await request(app)
        .post('/api/auth/register')
        .send(invalidData);

      // Assert
      expect(response.status).toBe(400);
      expect(response.body.errors).toContain('Invalid email format');
    });

    it('should return 409 for duplicate email', async () => {
      // Arrange
      const userData = {
        email: 'duplicate@example.com',
        password: 'Password123!',
      };

      // Create first user
      await request(app).post('/api/auth/register').send(userData);

      // Act - Try to create duplicate
      const response = await request(app)
        .post('/api/auth/register')
        .send(userData);

      // Assert
      expect(response.status).toBe(409);
      expect(response.body.message).toContain('already exists');
    });
  });

  describe('POST /api/auth/login', () => {
    beforeEach(async () => {
      // Create test user
      await request(app).post('/api/auth/register').send({
        email: 'login@example.com',
        password: 'Password123!',
      });
    });

    it('should login with correct credentials', async () => {
      // Arrange
      const credentials = {
        email: 'login@example.com',
        password: 'Password123!',
      };

      // Act
      const response = await request(app)
        .post('/api/auth/login')
        .send(credentials);

      // Assert
      expect(response.status).toBe(200);
      expect(response.body).toHaveProperty('token');
    });

    it('should reject invalid password', async () => {
      // Arrange
      const credentials = {
        email: 'login@example.com',
        password: 'WrongPassword!',
      };

      // Act
      const response = await request(app)
        .post('/api/auth/login')
        .send(credentials);

      // Assert
      expect(response.status).toBe(401);
      expect(response.body.message).toBe('Invalid credentials');
    });
  });
});
```

## E2E Testing with Playwright

### E2E Test Example

Create `tests/e2e/dashboard.spec.ts`:

```typescript
import { test, expect } from '@playwright/test';

test.describe('Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/login');
    await page.fill('[data-testid="email"]', 'test@example.com');
    await page.fill('[data-testid="password"]', 'Password123!');
    await page.click('[data-testid="login-button"]');
    await expect(page).toHaveURL('/dashboard');
  });

  test('should display user metrics', async ({ page }) => {
    // Arrange & page.goto('/ Act
    awaitdashboard');

    // Assert
    await expect(page.locator('[data-testid="total-users"]')).toBeVisible();
    await expect(page.locator('[data-testid="active-sessions"]')).toBeVisible();
    await expect(page.locator('[data-testid="revenue"]')).toBeVisible();
  });

  test('should navigate to user profile', async ({ page }) => {
    // Arrange & Act
    await page.goto('/dashboard');
    await page.click('[data-testid="user-row-1"]');

    // Assert
    await expect(page).toHaveURL('/users/1');
    await expect(page.locator('[data-testid="user-name"]')).toBeVisible();
  });

  test('should filter data correctly', async ({ page }) => {
    // Arrange & Act
    await page.goto('/dashboard');
    await page.fill('[data-testid="search-input"]', 'John');
    await page.press('[data-testid="search-input"]', 'Enter');

    // Assert
    await expect(page.locator('[data-testid="results-count"]')).toContainText('5');
  });

  test('should export data to CSV', async ({ page }) => {
    // Arrange
    const downloadPromise = page.waitForEvent('download');

    // Act
    await page.goto('/dashboard');
    await page.click('[data-testid="export-button"]');
    await page.click('[data-testid="export-csv"]');

    // Assert
    const download = await downloadPromise;
    expect(download.suggestedFilename()).toContain('export');
  });
});
```

### Visual Regression Testing

Create `tests/e2e/visual.spec.ts`:

```typescript
import { test, expect } from '@playwright/test';

test.describe('Visual Regression', () => {
  test('homepage visual snapshot', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveScreenshot('homepage.png', {
      maxDiffPixelRatio: 0.1,
    });
  });

  test('dashboard layout', async ({ page }) => {
    await page.goto('/dashboard');
    await expect(page).toHaveScreenshot('dashboard.png', {
      fullPage: true,
    });
  });

  test('mobile responsive layout', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto('/');
    await expect(page).toHaveScreenshot('mobile-homepage.png');
  });
});
```

## Test Utilities

### Test Fixtures

Create `tests/fixtures/user.fixture.ts`:

```typescript
import { faker } from '@faker-js/faker';

export interface UserFixture {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  password: string;
  role: 'admin' | 'user' | 'guest';
  createdAt: Date;
  updatedAt: Date;
}

export function createUserFixture(overrides?: Partial<UserFixture>): UserFixture {
  const defaultUser: UserFixture = {
    id: faker.string.uuid(),
    email: faker.internet.email(),
    firstName: faker.person.firstName(),
    lastName: faker.person.lastName(),
    password: 'TestPassword123!',
    role: 'user',
    createdAt: new Date(),
    updatedAt: new Date(),
  };

  return { ...defaultUser, ...overrides };
}

export function createAdminFixture(): UserFixture {
  return createUserFixture({ role: 'admin' });
}

export function createMultipleUsers(count: number): UserFixture[] {
  return Array.from({ length: count }, () => createUserFixture());
}
```

### Mock Server Setup

Create `tests/mocks/server.ts`:

```typescript
import { setupServer } from 'msw/node';
import { http, HttpResponse } from 'msw';

export const mockServer = setupServer(
  http.get('/api/users', () => {
    return HttpResponse.json({
      users: [
        { id: '1', name: 'John Doe', email: 'john@example.com' },
        { id: '2', name: 'Jane Smith', email: 'jane@example.com' },
      ],
      total: 2,
    });
  }),

  http.post('/api/users', async ({ request }) => {
    const body = await request.json();
    return HttpResponse.json(
      { id: '3', ...body },
      { status: 201 }
    );
  }),

  http.get('/api/users/:id', ({ params }) => {
    return HttpResponse.json({
      id: params.id,
      name: 'John Doe',
      email: 'john@example.com',
    });
  }),

  http.delete('/api/users/:id', ({ params }) => {
    return new HttpResponse(null, { status: 204 });
  })
);

beforeAll(() => mockServer.listen());
afterAll(() => mockServer.close());
afterEach(() => mockServer.resetHandlers());
```

## Running Tests

### NPM Scripts

```json
{
  "scripts": {
    "test": "jest",
    "test:watch": "jest --watch",
    "test:coverage": "jest --coverage",
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui",
    "test:e2e:headed": "playwright test --headed",
    "test:vitest": "vitest run",
    "test:vitest:watch": "vitest"
  }
}
```

### Running Specific Tests

```bash
# Run unit tests only
npm test -- --testPathPattern="tests/unit"

# Run integration tests only
npm test -- --testPathPattern="tests/integration"

# Run specific test file
npm test -- tests/unit/user.service.test.ts

# Run tests with coverage
npm run test:coverage

# Run E2E tests in headed mode
npm run test:e2e:headed

# Run E2E tests with UI
npm run test:e2e:ui
```

## CI Integration

### GitHub Actions Test Job

```yaml
test:
  name: Run Tests
  runs-on: ubuntu-latest
  services:
    postgres:
      image: postgres:15
      env:
        POSTGRES_PASSWORD: test
        POSTGRES_DB: biometrics_test
      ports:
        - 5432:5432
    redis:
      image: redis:7
      ports:
        - 6379:6379
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-node@v4
      with:
        node-version: '20.x'
        cache: 'npm'
    - run: npm ci
    - run: npm run test:coverage
      env:
        DATABASE_URL: postgresql://test:test@localhost:5432/biometrics_test
        REDIS_URL: redis://localhost:6379
```

## Related Documentation

- [CI-CD-PIPELINE.md](./CI-CD-PIPELINE.md)
- [MONITORING-SETUP.md](./MONITORING-SETUP.md)

---

**End of Testing Suite Documentation**
