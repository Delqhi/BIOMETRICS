# [PROJECT_NAME] - AGENTS.md

**Version:** [VERSION_NUMBER]  
**Date:** [DATE]  
**Status:** ✅ ACTIVE - MUST READ BEFORE EVERY SESSION  
**Template Origin:** `/Users/jeremy/dev/BIOMETRICS/templates/global/AGENTS.md`

---

## 🎯 PURPOSE

Dieses Dokument definiert die **projektspezifischen Regeln für KI-Agenten** im [PROJECT_NAME]-Projekt. Es basiert auf den globalen Regeln in `/Users/jeremy/dev/BIOMETRICS/rules/global/AGENTS.md` und muss für jedes neue Projekt angepasst werden.

**Jeder Agent MUSS diese Regeln lesen und befolgen BEVOR er arbeitet.**

---

## 📋 PROJECT OVERVIEW

### Basic Information

| Field | Value | Notes |
|-------|-------|-------|
| **Project Name** | [PROJECT_NAME] | E.g., "BIOMETRICS", "simone-webshop-01" |
| **Project Type** | [PROJECT_TYPE] | web-app, api, cli-tool, automation |
| **Tech Stack** | [TECH_STACK] | Primary technologies used |
| **Architecture** | [ARCHITECTURE] | monolith, microservices, serverless |
| **Primary Language** | [LANGUAGE] | TypeScript, Python, Go, etc. |
| **License** | [LICENSE] | MIT, Apache, Proprietary, etc. |

### Tech Stack Details

```
[PROJECT_NAME]/
├── Frontend:              [FRONTEND_FRAMEWORK]  (e.g., Next.js 14, React 18)
├── Backend:               [BACKEND_FRAMEWORK]   (e.g., Express, FastAPI, Gin)
├── Database:              [DATABASE_TYPE]       (e.g., PostgreSQL 15, MongoDB)
├── Cache:                 [CACHE_TYPE]          (e.g., Redis 7)
├── Queue:                 [QUEUE_TYPE]          (e.g., RabbitMQ, Kafka)
├── Auth:                  [AUTH_TYPE]           (e.g., JWT, OAuth2, NextAuth)
├── Storage:               [STORAGE_TYPE]        (e.g., S3, local, PostgreSQL)
├── Container:             [CONTAINER_TYPE]      (e.g., Docker, Podman)
├── Orchestration:         [ORCHESTRATION]       (e.g., docker-compose, Kubernetes)
└── CI/CD:                [CI_CD]               (e.g., GitHub Actions, GitLab CI)
```

### Architecture Pattern

**DIAGRAM:** (Optional - Add architecture diagram here)

```
┌─────────────────────────────────────────────────────────────────┐
│                    [PROJECT_NAME] ARCHITECTURE                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────┐     ┌─────────────────┐                  │
│  │   Client/App    │────▶│   Load Balancer │                  │
│  └─────────────────┘     └────────┬────────┘                  │
│                                     │                           │
│                                     ▼                           │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    API Gateway / Router                  │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                     │                           │
│         ┌───────────────────────────┼───────────────────────┐  │
│         │                           │                       │  │
│         ▼                           ▼                       ▼  │
│  ┌─────────────┐            ┌─────────────┐         ┌─────────────┐ │
│  │  Service A  │            │  Service B  │         │  Service C  │ │
│  └──────┬──────┘            └──────┬──────┘         └──────┬──────┘ │
│         │                         │                       │        │
│         └─────────────────────────┼───────────────────────┘        │
│                                   │                                │
│                                   ▼                                │
│  ┌─────────────────────────────────────────────────────────┐        │
│  │              Database / Cache / Storage                  │        │
│  └─────────────────────────────────────────────────────────┘        │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Deployment Information

| Environment | URL | Database | Notes |
|-------------|-----|----------|-------|
| **Development** | [DEV_URL] | [DEV_DB] | Local development |
| **Staging** | [STAGING_URL] | [STAGING_DB] | Pre-production testing |
| **Production** | [PROD_URL] | [PROD_DB] | Live environment |

---

## 🗂️ FOLDER STRUCTURE

### Standard Directory Layout

```
[PROJECT_NAME]/
├── .github/                    # GitHub workflows & templates
│   ├── ISSUE_TEMPLATE/
│   ├── workflows/
│   └── PULL_REQUEST_TEMPLATE.md
├── .sisyphus/                  # Sisyphus planning system
│   ├── plans/                  # Active plans
│   └── archive/                # Completed plans
├── docs/                       # Project documentation
│   ├── dev/                    # Developer documentation
│   ├── non-dev/                # User documentation
│   └── postman/                # API collections
├── src/                        # Main source code
│   ├── [module-a]/             # Feature module A
│   ├── [module-b]/             # Feature module B
│   ├── shared/                 # Shared utilities
│   └── types/                  # TypeScript types
├── tests/                      # Test files (optional)
├── scripts/                    # Build/deploy scripts
├── config/                     # Configuration files
├── docker/                     # Docker configurations
├── helm/                       # Kubernetes Helm charts (optional)
├── assets/                     # Static assets
├── inputs/                     # Input data
├── outputs/                    # Generated output
├── logs/                       # Log files
├── backups/                    # Database backups
├── [ENTRY_POINT]              # Main entry point
├── [CONFIG_FILE]              # Configuration file
├── docker-compose.yml          # Docker Compose config
├── Dockerfile                  # Docker image
├── package.json                # Node.js dependencies (if applicable)
├── go.mod                      # Go modules (if applicable)
├── requirements.txt            # Python dependencies (if applicable)
├── Makefile                    # Build commands
├── .env.example                # Environment template
├── .gitignore                  # Git ignore rules
├── README.md                   # Project readme
├── CHANGELOG.md                # Version history
├── LICENSE                     # License file
└── AGENTS.md                  # This file
```

### Module Structure

Each feature module should follow this pattern:

```
[module-name]/
├── README.md                   # Module overview
├── [module-name].ts            # Main module file
├── [module-name].test.ts       # Unit tests
├── [module-name].e2e.test.ts   # E2E tests
├── types.ts                    # Module-specific types
├── constants.ts               # Module constants
├── utils.ts                    # Module utilities
├── config.ts                  # Module configuration
└── docs/                      # Module documentation
    └── API.md                 # Module API docs
```

---

## 🔧 NAMING CONVENTIONS

### File Naming

| Type | Convention | Example |
|------|------------|---------|
| **TypeScript/JS** | kebab-case | `user-service.ts`, `auth-middleware.js` |
| **Python** | snake_case | `user_service.py`, `auth_middleware.py` |
| **Go** | snake_case | `user_service.go`, `auth_middleware.go` |
| **Components** | PascalCase | `UserProfile.tsx`, `LoginForm.vue` |
| **Tests** | [name].test.ts | `user-service.test.ts` |
| **Config** | kebab-case | `app-config.json`, `nginx.conf` |
| **Docker** | kebab-case | `Dockerfile`, `docker-compose.yml` |

### Variable Naming

| Type | Convention | Example |
|------|------------|---------|
| **Variables** | camelCase | `userName`, `totalCount` |
| **Constants** | UPPER_SNAKE_CASE | `MAX_RETRY_COUNT`, `API_BASE_URL` |
| **Functions** | camelCase | `getUserById()`, `calculateTotal()` |
| **Classes** | PascalCase | `UserService`, `AuthMiddleware` |
| **Interfaces** | PascalCase | `IUser`, `IApiResponse` |
| **Types** | PascalCase | `UserType`, `ApiResponseType` |
| **Enums** | PascalCase | `UserRole`, `ApiStatus` |
| **Boolean** | is/has/can/should | `isActive`, `hasPermission`, `canEdit` |

### Database Naming

| Type | Convention | Example |
|------|------------|---------|
| **Tables** | snake_case, plural | `users`, `order_items` |
| **Columns** | snake_case | `created_at`, `user_id` |
| **Primary Keys** | `id` | `id` |
| **Foreign Keys** | `[table]_id` | `user_id`, `order_id` |
| **Indexes** | `idx_[table]_[column]` | `idx_users_email` |
| **Constraints** | `[table]_[constraint]` | `users_pkey` |

### Git Branch Naming

| Type | Convention | Example |
|------|------------|---------|
| **Feature** | `feature/[ticket]-[short-description]` | `feature/123-add-user-auth` |
| **Bugfix** | `bugfix/[ticket]-[short-description]` | `bugfix/456-fix-login-error` |
| **Hotfix** | `hotfix/[ticket]-[short-description]` | `hotfix/789-security-patch` |
| **Release** | `release/v[version]` | `release/v1.2.0` |
| **Docs** | `docs/[ticket]-[short-description]` | `docs/101-update-readme` |

### Commit Message Convention

Follow Conventional Commits:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `style:` - Code style (formatting)
- `refactor:` - Code refactoring
- `perf:` - Performance improvement
- `test:` - Tests
- `chore:` - Maintenance

**Examples:**
```
feat(auth): add JWT token refresh endpoint
fix(api): resolve null pointer in user service
docs(readme): update installation instructions
refactor(users): extract validation to separate module
```

---

## 💻 CODING STANDARDS

### TypeScript Standards

```json
// tsconfig.json
{
  "compilerOptions": {
    "strict": true,
    "noImplicitAny": true,
    "strictNullChecks": true,
    "strictFunctionTypes": true,
    "strictBindCallApply": true,
    "strictPropertyInitialization": true,
    "noImplicitThis": true,
    "alwaysStrict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,
    "outDir": "./dist"
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist", "**/*.test.ts"]
}
```

### Code Style Rules

**DO:**
- ✅ Use meaningful variable names
- ✅ Keep functions small (< 50 lines)
- ✅ Write JSDoc comments for public APIs
- ✅ Handle errors properly with try-catch
- ✅ Use async/await over raw promises
- ✅ Prefer const over let
- ✅ Use template literals over string concatenation

**DON'T:**
- ❌ Use `any` type without justification
- ❌ Leave commented-out code
- ❌ Hardcode configuration values
- ❌ Use magic numbers
- ❌ Create functions with > 5 parameters
- ❌ Nest more than 3 levels deep
- ❌ Ignore ESLint/Prettier rules

### Error Handling

```typescript
// ✅ CORRECT: Proper error handling
async function fetchUser(id: string): Promise<User> {
  try {
    const response = await fetch(`/api/users/${id}`);
    if (!response.ok) {
      throw new ApiError(`Failed to fetch user: ${response.statusText}`, {
        statusCode: response.status,
        code: 'USER_FETCH_FAILED'
      });
    }
    return response.json();
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }
    logger.error('Unexpected error fetching user', { id, error });
    throw new ApiError('Internal server error', {
      statusCode: 500,
      code: 'INTERNAL_ERROR'
    });
  }
}

// ❌ WRONG: Empty catch
async function fetchUser(id: string): Promise<User> {
  try {
    const response = await fetch(`/api/users/${id}`);
    return response.json();
  } catch (e) {
    // DON'T DO THIS!
  }
}
```

### Logging Standards

```typescript
import { logger } from './utils/logger';

// Log levels: error, warn, info, debug, trace
logger.error('Operation failed', { error, context });
logger.warn('Potential issue', { warning });
logger.info('Operation completed', { result });
logger.debug('Debug information', { details });
logger.trace('Detailed trace', { trace });
```

---

## 🔌 API STANDARDS

### REST API Design

| Method | Endpoint | Description |
|--------|----------|-------------|
| **GET** | `/api/[resource]` | List all resources |
| **GET** | `/api/[resource]/:id` | Get single resource |
| **POST** | `/api/[resource]` | Create new resource |
| **PUT** | `/api/[resource]/:id` | Update resource (full) |
| **PATCH** | `/api/[resource]/:id` | Update resource (partial) |
| **DELETE** | `/api/[resource]/:id` | Delete resource |

### API Response Format

```typescript
// Success Response
interface ApiSuccess<T> {
  success: true;
  data: T;
  meta?: {
    page?: number;
    limit?: number;
    total?: number;
  };
}

// Error Response
interface ApiError {
  success: false;
  error: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
}

// Combined Response Type
type ApiResponse<T> = ApiSuccess<T> | ApiError;
```

### Authentication

| Method | Use Case | Header |
|--------|----------|--------|
| **Bearer Token** | API requests | `Authorization: Bearer <token>` |
| **API Key** | Server-to-server | `X-API-Key: <key>` |
| **Basic Auth** | Simple auth | `Authorization: Basic <base64>` |

### Rate Limiting

| Tier | Requests | Window |
|------|----------|--------|
| **Anonymous** | 100 | per minute |
| **Authenticated** | 1000 | per minute |
| **Premium** | 10000 | per minute |

---

## 🗄️ DATABASE STANDARDS

### Connection Management

```typescript
// ✅ CORRECT: Connection pooling
const pool = new Pool({
  host: process.env.DB_HOST,
  port: parseInt(process.env.DB_PORT || '5432'),
  database: process.env.DB_NAME,
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  max: 20,                    // Max connections
  idleTimeoutMillis: 30000,   // Close idle clients after 30s
  connectionTimeoutMillis: 2000 // Return error after 2s
});

// ❌ WRONG: No connection pooling
const client = new Client({
  host: process.env.DB_HOST,
  // ...
});
await client.connect();
await client.query('SELECT * FROM users');
await client.end();
```

### Query Building

```typescript
// ✅ CORRECT: Parameterized queries (prevent SQL injection)
const result = await pool.query(
  'SELECT * FROM users WHERE id = $1 AND status = $2',
  [userId, 'active']
);

// ❌ WRONG: String concatenation (SQL injection vulnerability)
const result = await pool.query(
  `SELECT * FROM users WHERE id = ${userId} AND status = 'active'`
);
```

### Migrations

```bash
# Create migration
npm run migrate:create -- add_user_email_verification

# Run migrations
npm run migrate:up

# Rollback
npm run migrate:down
```

---

## 🧪 TESTING STANDARDS

### Test Structure

```typescript
describe('[ModuleName]', () => {
  describe('[Function/Method Name]', () => {
    it('should [expected behavior]', async () => {
      // Arrange
      const input = { /* test data */ };
      
      // Act
      const result = await functionUnderTest(input);
      
      // Assert
      expect(result).toEqual({ /* expected output */ });
    });

    it('should throw error when [invalid input]', async () => {
      // Arrange
      const invalidInput = { /* invalid data */ };
      
      // Act & Assert
      await expect(functionUnderTest(invalidInput))
        .rejects.toThrow(ErrorType);
    });
  });
});
```

### Test Coverage Requirements

| Type | Minimum Coverage |
|------|------------------|
| **Unit Tests** | 80% |
| **Integration Tests** | 70% |
| **E2E Tests** | Critical paths only |

### Test Naming Convention

```
[Method] [Condition] [Expected Result]

Examples:
- should return user when valid ID provided
- should throw ValidationError when email is invalid
- should return empty array when no results found
```

---

## 📦 DEPENDENCY MANAGEMENT

### Version Constraints

| Symbol | Meaning | Example |
|--------|---------|---------|
| `^` | Compatible | `^1.2.3` = `1.x.x` |
| `~` | Patch compatible | `~1.2.3` = `1.2.x` |
| `==` | Exact version | `==1.2.3` |
| `>=` | Minimum | `>=1.2.3` |
| `<=` | Maximum | `<=2.0.0` |

### Security Auditing

```bash
# Run security audit
npm audit
npm audit fix

# Check for vulnerabilities
npm outdated
```

### Lock Files

- **npm:** `package-lock.json` (commit!)
- **yarn:** `yarn.lock` (commit!)
- **pip:** `requirements.txt` + ` Pipfile.lock`
- **go:** `go.sum` (commit!)

---

## 🔒 SECURITY STANDARDS

### Secrets Management

**NEVER commit secrets to git!**

```bash
# .gitignore
.env
.env.local
.env.*.local
*.pem
*.key
credentials.json
secrets/
```

**Use environment variables:**
```typescript
const apiKey = process.env.API_KEY;
if (!apiKey) {
  throw new Error('API_KEY environment variable is required');
}
```

### Input Validation

```typescript
// ✅ CORRECT: Validate all inputs
import { z } from 'zod';

const UserSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8).max(100),
  name: z.string().min(1).max(100),
});

function createUser(data: unknown) {
  const validated = UserSchema.parse(data); // Throws if invalid
  // ... proceed with validated data
}
```

### Security Headers

```typescript
// Security middleware
app.use((req, res, next) => {
  res.setHeader('X-Content-Type-Options', 'nosniff');
  res.setHeader('X-Frame-Options', 'DENY');
  res.setHeader('X-XSS-Protection', '1; mode=block');
  res.setHeader('Strict-Transport-Security', 'max-age=31536000; includeSubDomains');
  next();
});
```

---

## 📊 MONITORING & LOGGING

### Log Levels

| Level | Use Case |
|-------|----------|
| **error** | System errors requiring immediate attention |
| **warn** | Potential issues that should be investigated |
| **info** | Normal operational events |
| **debug** | Debugging information |
| **trace** | Detailed trace information |

### Metrics to Track

| Metric | Description |
|--------|-------------|
| **Response Time** | API response time (p50, p95, p99) |
| **Error Rate** | Percentage of failed requests |
| **Throughput** | Requests per second |
| **CPU Usage** | Server CPU utilization |
| **Memory Usage** | Server memory utilization |
| **Database Connections** | Active DB connections |

### Health Checks

```typescript
app.get('/health', async (req, res) => {
  const health = {
    status: 'healthy',
    uptime: process.uptime(),
    timestamp: Date.now(),
    checks: {
      database: await checkDatabase(),
      redis: await checkRedis(),
      external: await checkExternalServices(),
    },
  };
  
  const isHealthy = health.checks.database && 
                    health.checks.redis && 
                    health.checks.external;
  
  res.status(isHealthy ? 200 : 503).json(health);
});
```

---

## 🚀 DEPLOYMENT STANDARDS

### Docker Best Practices

```dockerfile
# ✅ CORRECT: Multi-stage build for smaller images
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build

FROM node:20-alpine AS runner
WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
USER node
CMD ["node", "dist/index.js"]

# ❌ WRONG: No multi-stage, large base image
FROM node:20
WORKDIR /app
COPY . .
RUN npm install
CMD ["node", "index.js"]
```

### Environment-Specific Config

| Environment | Features |
|-------------|----------|
| **Development** | Debug logging, hot reload, mock services |
| **Staging** | Production-like, test data, full logging |
| **Production** | Minimal logging, error tracking, caching |

---

## 🔧 PROJECT-SPECIFIC RULES

### Custom Rules for [PROJECT_NAME]

**[INSERT PROJECT-SPECIFIC RULES HERE]**

Example rules:
- Maximum file size for uploads: 10MB
- Allowed file types: jpg, png, pdf
- API rate limit: 100 req/min
- Session timeout: 30 minutes
- Cache TTL: 1 hour

### Known Limitations

| Limitation | Description | Workaround |
|------------|-------------|-------------|
| [LIMIT_1] | [Description] | [Workaround] |
| [LIMIT_2] | [Description] | [Workaround] |

### Performance Targets

| Metric | Target |
|--------|--------|
| API Response Time (p95) | < 200ms |
| Page Load Time | < 3s |
| Database Query Time | < 100ms |
| Build Time | < 5min |

---

## 🛠️ TROUBLESHOOTING

### Common Issues

#### Issue 1: [Common Error]

**Symptoms:**
- [Symptom description]

**Cause:**
- [Root cause]

**Solution:**
```bash
# Fix command
[command]
```

#### Issue 2: [Common Error]

**Symptoms:**
- [Symptom description]

**Cause:**
- [Root cause]

**Solution:**
```bash
# Fix command
[command]
```

### Debug Commands

```bash
# View logs
docker logs [container_name]
kubectl logs [pod_name]

# Check status
curl http://localhost:[PORT]/health

# Database query
docker exec -it [container] psql -U [user] -d [db]

# Clear cache
redis-cli FLUSHDB
```

---

## 📚 REFERENCES

### Global Rules

| Document | Location | Purpose |
|----------|----------|---------|
| **Global AGENTS.md** | `/Users/jeremy/dev/BIOMETRICS/rules/global/AGENTS.md` | Source of Truth for all agent rules |
| **Coding Standards** | `/Users/jeremy/dev/BIOMETRICS/rules/global/coding-standards.md` | Detailed coding rules |
| **Documentation Rules** | `/Users/jeremy/dev/BIOMETRICS/rules/global/documentation-rules.md` | Documentation standards |
| **Git Workflow** | `/Users/jeremy/dev/BIOMETRICS/rules/global/git-workflow.md` | Git branching & commit rules |
| **Security Mandates** | `/Users/jeremy/dev/BIOMETRICS/rules/global/security-mandates.md` | Security requirements |

### Project Documentation

| Document | Location | Purpose |
|----------|----------|---------|
| **ARCHITECTURE.md** | `./ARCHITECTURE.md` | System architecture |
| **README.md** | `./README.md` | Project overview |
| **CHANGELOG.md** | `./CHANGELOG.md` | Version history |
| **SETUP.md** | `./docs/setup/SETUP.md` | Setup instructions |

### External Resources

| Resource | URL | Purpose |
|----------|-----|---------|
| **[Framework Docs]** | [URL] | [Purpose] |
| **[API Reference]** | [URL] | [Purpose] |
| **[Style Guide]** | [URL] | [Purpose] |

---

## 📝 QUICK REFERENCE CARD

```
┌─────────────────────────────────────────────────────────────────┐
│              [PROJECT_NAME] QUICK REFERENCE                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  📋 KEY COMMANDS:                                               │
│    dev:        [DEV_COMMAND]                                    │
│    build:      [BUILD_COMMAND]                                  │
│    test:       [TEST_COMMAND]                                   │
│    lint:       [LINT_COMMAND]                                   │
│    deploy:     [DEPLOY_COMMAND]                                 │
│                                                                  │
│  🔧 CONFIGURATION:                                              │
│    API URL:     [API_URL]                                        │
│    DB:         [DB_CONNECTION]                                  │
│    Cache:      [CACHE_CONNECTION]                               │
│                                                                  │
│  📦 DEPENDENCIES:                                               │
│    Frontend:   [FRONTEND_VERSION]                               │
│    Backend:    [BACKEND_VERSION]                                 │
│    Database:   [DB_VERSION]                                     │
│                                                                  │
│  🚨 EMERGENCY CONTACTS:                                         │
│    Lead Dev:   [EMAIL]                                          │
│    On-Call:    [EMAIL]                                          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## ✅ COMPLIANCE CHECKLIST

Before submitting code, verify:

- [ ] Code follows naming conventions
- [ ] No hardcoded secrets
- [ ] Error handling implemented
- [ ] Unit tests added (if applicable)
- [ ] Documentation updated (if needed)
- [ ] No console.log statements left
- [ ] ESLint/Prettier passes
- [ ] TypeScript compiles without errors
- [ ] Security vulnerabilities fixed

---

## 📅 CHANGELOG

| Version | Date | Changes |
|---------|------|---------|
| [VERSION] | [DATE] | Initial template |

---

**Template Version:** 1.0  
**Last Updated:** February 2026  
**Based On:** BIOMETRICS Global Rules v1.0

---

*This template is maintained in `/Users/jeremy/dev/BIOMETRICS/templates/global/AGENTS.md`*
