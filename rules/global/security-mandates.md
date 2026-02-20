# SECURITY MANDATES - BIOMETRICS PROJECT

**Version:** 1.0  
**Status:** ACTIVE  
**Effective Date:** 2026-02-20  
**Classification:** CONFIDENTIAL  

---

## TABLE OF CONTENTS

1. [Executive Summary](#1-executive-summary)
2. [Secrets Management](#2-secrets-management)
3. [Git Security](#3-git-security)
4. [Input Validation](#4-input-validation)
5. [Authentication and Authorization](#5-authentication-and-authorization)
6. [Data Protection](#6-data-protection)
7. [Infrastructure Security](#7-infrastructure-security)
8. [Monitoring and Auditing](#8-monitoring-and-auditing)
9. [OWASP Top 10](#9-owasp-top-10)
10. [Security Review Checklist](#10-security-review-checklist)
11. [Incident Response](#11-incident-response)
12. [References](#12-references)

---

## 1. EXECUTIVE SUMMARY

This document establishes the mandatory security requirements for the BIOMETRICS project. All developers, security professionals, and operations teams MUST adhere to these mandates without exception. The security of biometric data is paramount, as it represents highly sensitive personal information that cannot be changed if compromised (unlike passwords).

### 1.1 Purpose

The purpose of this security mandates document is to:

- Establish a comprehensive security framework for all BIOMETRICS project activities
- Define mandatory security controls and best practices
- Provide concrete guidance for secure development and deployment
- Ensure compliance with relevant regulations (GDPR, CCPA, BSI)
- Protect against known attack vectors and vulnerabilities

### 1.2 Scope

These mandates apply to:

- All source code in the BIOMETRICS repository
- All infrastructure configurations (Docker, Kubernetes, cloud resources)
- All development, staging, and production environments
- All team members with access to BIOMETRICS systems
- All third-party integrations and dependencies

### 1.3 Guiding Principles

1. **Defense in Depth**: Multiple layers of security controls
2. **Least Privilege**: Minimal permissions required for tasks
3. **Fail Secure**: Default to secure state on failure
4. **Zero Trust**: Never trust, always verify
5. **Privacy by Design**: Data protection built into architecture

### 1.4 Data Classification

| Classification | Description | Examples | Handling |
|---------------|-------------|----------|----------|
| **CRITICAL** | Biometric templates, raw biometric data | Fingerprint images, face scans, iris templates | Encrypt at rest, encrypt in transit, access logging, MFA required |
| **HIGH** | Personal identifiers, PII | Name, address, email, phone | Encrypt at rest, access logging |
| **MEDIUM** | Technical metadata | Timestamps, IP addresses, user agents | Access logging |
| **LOW** | Aggregated anonymized data | Usage statistics | No special handling |

---

## 2. SECRETS MANAGEMENT

### 2.1 Core Principle

**NEVER store secrets in source code, configuration files, or version control.** This is the most fundamental security rule. Secrets include API keys, passwords, tokens, certificates, and any credentials that grant access to systems or data.

### 2.2 Secrets Inventory

All secrets used in the BIOMETRICS project MUST be documented in the secrets inventory:

```yaml
# secrets-inventory.yaml
secrets:
  - name: DATABASE_PASSWORD
    type: password
    rotation_period: 90_days
    last_rotated: 2026-01-15
    owners:
      - security-team@biometrics.local
    environments:
      - production
      - staging
    
  - name: JWT_SECRET_KEY
    type: key
    rotation_period: 180_days
    last_rotated: 2025-12-01
    owners:
      - auth-team@biometrics.local
    environments:
      - production
      - staging
      - development
      
  - name: ENCRYPTION_KEY_BIOMETRIC
    type: key
    rotation_period: 365_days
    last_rotated: 2025-06-01
    owners:
      - security-team@biometrics.local
    environments:
      - production
```

### 2.3 Environment Variables

All secrets MUST be managed through environment variables in production and staging environments. Development environments MAY use `.env` files with restricted access.

#### 2.3.1 Environment Variable Standards

```bash
# CORRECT: Descriptive, prefixed naming
export BIOMETRICS_DB_PASSWORD="secure_random_string_here"
export BIOMETRICS_JWT_SECRET="jwt_secret_key_here"
export BIOMETRICS_ENCRYPTION_KEY="aes256_key_here"

# INCORRECT: Ambiguous or generic naming
export DB_PASS="password"      # Too vague
export SECRET="value"          # Too generic
export KEY="12345"            # Unclear purpose
```

#### 2.3.2 Environment Validation

All applications MUST validate required environment variables at startup:

```typescript
// src/config/validate-env.ts
import { z } from 'zod';

const envSchema = z.object({
  // Database
  BIOMETRICS_DB_HOST: z.string().ip(),
  BIOMETRICS_DB_PORT: z.coerce.number().min(1).max(65535),
  BIOMETRICS_DB_NAME: z.string().min(1).max(64),
  BIOMETRICS_DB_USER: z.string().min(1).max(64),
  BIOMETRICS_DB_PASSWORD: z.string().min(32).max(256),
  
  // Authentication
  BIOMETRICS_JWT_SECRET: z.string().min(64).max(256),
  BIOMETRICS_JWT_EXPIRY: z.coerce.number().min(300).max(86400),
  BIOMETRICS_REFRESH_TOKEN_SECRET: z.string().min(64).max(256),
  
  // Encryption
  BIOMETRICS_ENCRYPTION_KEY: z.string().length(64), // 256-bit hex key
  BIOMETRICS_ENCRYPTION_ALGORITHM: z.enum(['AES-256-GCM', 'AES-256-CBC']),
  
  // External Services
  BIOMETRICS_AWS_ACCESS_KEY_ID: z.string().min(20).max(20),
  BIOMETRICS_AWS_SECRET_ACCESS_KEY: z.string().min(40).max(40),
  BIOMETRICS_AWS_REGION: z.string().length(2),
  
  // Application
  NODE_ENV: z.enum(['development', 'staging', 'production']),
  LOG_LEVEL: z.enum(['debug', 'info', 'warn', 'error']),
});

type Env = z.infer<typeof envSchema>;

export function validateEnvironment(): Env {
  const result = envSchema.safeParse(process.env);
  
  if (!result.success) {
    const errors = result.error.flatten().fieldErrors;
    console.error('Environment validation failed:', errors);
    throw new Error(`Environment validation failed: ${JSON.stringify(errors)}`);
  }
  
  return result.data;
}

export const env = validateEnvironment();
```

### 2.4 Vault Integration

For production environments, HashiCorp Vault MUST be used for secrets management.

#### 2.4.1 Vault Architecture

```typescript
// src/lib/vault/client.ts
import Vault from 'node-vault';

interface VaultConfig {
  endpoint: string;
  token: string;
  namespace?: string;
}

interface SecretVersion {
  version: number;
  created_time: string;
  data: {
    [key: string]: string;
  };
}

export class VaultClient {
  private client: Vault.client;
  private readonly mountPoint: string;
  private readonly secretPath: string;

  constructor(config: VaultConfig) {
    this.client = Vault({
      endpoint: config.endpoint,
      token: config.token,
      namespace: config.namespace,
    });
    this.mountPoint = 'biometrics';
    this.secretPath = 'production';
  }

  async getSecret(key: string): Promise<string> {
    try {
      const secret = await this.client.read(
        `${this.mountPoint}/data/${this.secretPath}`
      ) as SecretVersion;
      
      if (!secret.data.data[key]) {
        throw new Error(`Secret key not found: ${key}`);
      }
      
      return secret.data.data[key];
    } catch (error) {
      console.error(`Failed to retrieve secret ${key}:`, error);
      throw new Error(`Secret retrieval failed for key: ${key}`);
    }
  }

  async setSecret(key: string, value: string): Promise<void> {
    try {
      await this.client.write(
        `${this.mountPoint}/data/${this.secretPath}`,
        { [key]: value }
      );
      console.log(`Secret ${key} updated successfully`);
    } catch (error) {
      console.error(`Failed to set secret ${key}:`, error);
      throw new Error(`Secret update failed for key: ${key}`);
    }
  }

  async deleteSecret(key: string): Promise<void> {
    try {
      await this.client.delete(
        `${this.mountPoint}/data/${this.secretPath}`
      );
      console.log(`Secret ${key} deleted successfully`);
    } catch (error) {
      console.error(`Failed to delete secret ${key}:`, error);
      throw new Error(`Secret deletion failed for key: ${key}`);
    }
  }
}
```

#### 2.4.2 Vault Authentication

```yaml
# Kubernetes Vault Authentication
apiVersion: apps/v1
kind: Deployment
metadata:
  name: biometrics-api
spec:
  template:
    spec:
      serviceAccountName: biometrics-vault
      containers:
        - name: api
          env:
            - name: VAULT_ADDR
              value: "https://vault.biometrics.local:8200"
            - name: VAULT_ROLE
              value: "biometrics-api"
            - name: VAULT_AUTH_PATH
              value: "auth/kubernetes/login"
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: biometrics-vault
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: biometrics-vault
subjects:
  - kind: ServiceAccount
    name: biometrics-vault
roleRef:
  kind: ClusterRole
  name: system:auth-delegator
  apiGroup: rbac.authorization.k8s.io
```

### 2.5 Git Ignore Configuration

All secrets and sensitive files MUST be excluded from version control:

```gitignore
# Environment Files
.env
.env.local
.env.*.local
.env.production
.env.staging
*.env

# Secrets and Credentials
*.pem
*.key
*.crt
*.p12
*.pfx
credentials.json
service-account.json
*.credentials

# Terraform
*.tfstate
*.tfstate.*
.terraform/
terraform.tfvars

# Kubernetes
*.kubeconfig
kubeconfig*

# Application
config/secrets.yaml
config/credentials.yaml
secrets/
keys/
```

### 2.6 Secret Rotation Policy

#### 2.6.1 Rotation Schedule

| Secret Type | Rotation Period | Approval Required |
|-------------|----------------|-------------------|
| Database Password | 90 days | Security Team |
| API Keys (External) | 90 days | Security Team |
| JWT Secrets | 180 days | Auth Team |
| Encryption Keys | 365 days | Security Team |
| SSH Keys | 180 days | DevOps Team |
| TLS Certificates | 90 days | Auto-renew |

#### 2.6.2 Rotation Procedure

```bash
#!/bin/bash
# scripts/rotate-secrets.sh

set -euo pipefail

echo "Starting secret rotation..."
echo "Timestamp: $(date -u +%Y-%m-%dT%H:%M:%SZ)"

# Generate new database password
NEW_DB_PASSWORD=$(openssl rand -base64 32)
echo "Generated new database password"

# Update Vault
vault kv put biometrics/production/database \
  password="$NEW_DB_PASSWORD" \
  rotated_by="${USER:-automation}" \
  rotation_date="$(date -u +%Y-%m-%d)"

# Update Kubernetes secrets
kubectl create secret generic biometrics-db-creds \
  --from-literal=password="$NEW_DB_PASSWORD" \
  --dry-run=client -o yaml | kubectl apply -f -

# Notify relevant teams
curl -X POST "${SLACK_WEBHOOK_URL}" \
  -H 'Content-Type: application/json' \
  -d "{\"text\":\"Secret rotation completed for BIOMETRICS database credentials\"}"

echo "Secret rotation completed successfully"
```

---

## 3. GIT SECURITY

### 3.1 Signed Commits

All commits to protected branches MUST be signed with GPG or SSH keys. This ensures commit authenticity and prevents tampering.

#### 3.1.1 GPG Setup

```bash
# Generate GPG key
gpg --full-generate-key
# Algorithm: RSA and RSA (default)
# Key size: 4096 bits
# Validity: 5 years
# Name: Your Name
# Email: your.email@biometrics.local

# List keys
gpg --list-secret-keys --keyid-format=long

# Configure Git
git config --global user.signingkey YOUR_KEY_ID
git config --global commit.gpgsign true
git config --global gpg.program gpg

# Export public key for GitHub
gpg --armor --export YOUR_KEY_ID
```

#### 3.1.2 GitHub Configuration

```yaml
# .github/CODEOWNERS
# Security Team owns all security-related code
/rules/** @biometrics/security-team
/src/auth/** @biometrics/auth-team
/src/encryption/** @biometrics/security-team
/config/** @biometrics/devops-team

# Default owners
* @biometrics/default-reviewers
```

### 3.2 Branch Protection Rules

All production branches MUST have protection rules enforced:

```yaml
# GitHub Branch Protection (via API or Admin UI)
branch_protection:
  - pattern: main
    required_status_checks:
      - context: ci/lint
      - context: ci/typecheck  
      - context: ci/test
      - context: ci/security-scan
      - context: ci/dependency-review
    required_reviews: 2
    dismiss_stale_reviews: true
    require_code_owner_reviews: true
    required_signoff_on_required_reviewers: true
    restrict_pushes:
      - biometrics-admins
    allow_force_pushes: false
    allow_deletions: false
    
  - pattern: develop
    required_status_checks:
      - context: ci/lint
      - context: ci/test
    required_reviews: 1
    dismiss_stale_reviews: true
    require_code_owner_reviews: false
    allow_force_pushes: false
    restrict_pushes:
      - biometrics-developers
```

### 3.3 Commit Message Standards

All commits MUST follow conventional commits format:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting
- `refactor`: Code restructure
- `test`: Tests
- `chore`: Maintenance
- `security`: Security-related changes

Example:
```
security(auth): add rate limiting to login endpoint

- Implemented rate limiting of 5 attempts per 15 minutes
- Added failed login attempt tracking
- Added account lockout after 10 failed attempts

Fixes: CVE-2026-0123
```

### 3.4 Code Review Requirements

#### 3.4.1 Review Checklist

```markdown
## Security Code Review Checklist

### Authentication
- [ ] Authentication logic properly validates credentials
- [ ] Passwords properly hashed (bcrypt/argon2)
- [ ] MFA enforcement for privileged accounts
- [ ] Session tokens are cryptographically random
- [ ] Session timeout properly enforced

### Authorization
- [ ] RBAC properly implemented
- [ ] Permission checks on all protected endpoints
- [ ] No privilege escalation vulnerabilities
- [ ] Resource ownership verified

### Input Validation
- [ ] All inputs validated and sanitized
- [ ] SQL injection prevention (parameterized queries)
- [ ] XSS prevention (output encoding)
- [ ] CSRF tokens implemented
- [ ] File upload validation

### Data Protection
- [ ] Sensitive data encrypted at rest
- [ ] TLS for all data in transit
- [ ] No sensitive data in logs
- [ ] No sensitive data in error messages

### Error Handling
- [ ] No stack traces in production
- [ ] Generic error messages for users
- [ ] Proper error logging
- [ ] Fail-secure defaults
```

### 3.5 CI/CD Security Scanning

#### 3.5.1 GitHub Actions Workflow

```yaml
name: Security Scan

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  dependency-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
      
      - name: Install dependencies
        run: npm ci
      
      - name: Run npm audit
        run: npm audit --audit-level=moderate
        
      - name: Check for vulnerabilities
        run: |
          npm audit --json > audit-results.json
          if [ $(jq '.metadata.vulnerability_count' audit-results.json) -gt 0 ]; then
            echo "Vulnerabilities found!"
            exit 1
          fi
      
      - name: Dependency Review
        uses: actions/dependency-review-action@v4
        with:
          fail-on-severity: moderate

  secret-scanning:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Scan for secrets
        uses: trufflesecurity/trufflehog@main
        with:
          base: ${{ github.event.repository.default_branch }}
          head: HEAD
      
      - name: GitLeaks
        uses: gitleaks/gitleaks-action@v2
        env:
          GITLEAKS_LICENSE: ${{ secrets.GITLEAKS_LICENSE }}

  code-scanning:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Initialize CodeQL
        uses: github/codeql-action/init@v3
        with:
          languages: javascript,typescript
          queries: security-extended
      
      - name: Build
        run: npm ci && npm run build
      
      - name: Perform CodeQL Analysis
        uses: github/codeql-action/analyze@v3
        with:
          category: "/language:javascript-typescript"

  container-scan:
    runs-on: ubuntu-latest
    if: github.event_name == 'push'
    steps:
      - uses: actions/checkout@v4
      
      - name: Build Docker image
        run: docker build -t biometrics:${{ github.sha }} .
      
      - name: Scan Docker image
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: 'biometrics:${{ github.sha }}'
          format: sarif
          output: 'trivy-results.sarif'
      
      - name: Upload results to GitHub
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: 'trivy-results.sarif'
```

---

## 4. INPUT VALIDATION

### 4.1 Core Principle

**All input is malicious until proven otherwise.** Every piece of data entering the application from external sources MUST be validated, sanitized, and treated as potentially dangerous. This includes user form data, API parameters, HTTP headers, cookies, and data from databases.

### 4.2 SQL Injection Prevention

#### 4.2.1 Never Use String Concatenation

```typescript
// ❌ INSECURE: Never do this!
const query = `SELECT * FROM users WHERE id = ${userId}`;
db.query(query);

// ❌ INSECURE: Don't use template literals with user input
const sql = `SELECT * FROM users WHERE email = '${email}'`;

// ❌ INSECURE: Don't use string replacement
const query = "SELECT * FROM users WHERE id = " + userId;
```

#### 4.2.2 Always Use Parameterized Queries

```typescript
// ✅ SECURE: Parameterized queries
const query = 'SELECT * FROM users WHERE id = $1';
const result = await db.query(query, [userId]);

// ✅ SECURE: Multiple parameters
const sql = `
  SELECT u.*, r.name as role_name 
  FROM users u 
  JOIN roles r ON u.role_id = r.id 
  WHERE u.email = $1 AND u.active = $2
`;
const result = await db.query(sql, [email, true]);

// ✅ SECURE: Named parameters
const query = `
  INSERT INTO audit_log (user_id, action, timestamp)
  VALUES (@userId, @action, @timestamp)
`;
const result = await db.query(query, {
  userId: userId,
  action: 'LOGIN',
  timestamp: new Date(),
});
```

#### 4.2.3 ORM Best Practices

```typescript
// ✅ SECURE: Using ORM with proper escaping
const user = await db.user.findFirst({
  where: {
    email: email, // ORM handles escaping
    active: true,
  },
});

// ✅ SECURE: Parameterized raw queries
const result = await db.$queryRaw`
  SELECT * FROM biometric_templates 
  WHERE user_id = ${userId} 
  AND created_at > NOW() - INTERVAL '30 days'
`;

// ❌ INSECURE: Never pass raw SQL with user input
const result = await db.$queryRaw(
  `SELECT * FROM users WHERE name = '${name}'`
);
```

### 4.3 XSS Prevention

#### 4.3.1 Output Encoding

```typescript
// ✅ SECURE: Using a library for HTML encoding
import { encode } from 'html-entities';

function safeDisplay(userInput: string): string {
  return encode(userInput);
}

// Usage in templates
const safeContent = safeDisplay(userProvidedContent);
document.getElementById('output').innerHTML = safeContent;
```

#### 4.3.2 Content Security Policy

```typescript
// Express.js CSP middleware
import helmet from 'helmet';

app.use(helmet.contentSecurityPolicy({
  directives: {
    defaultSrc: ["'self'"],
    scriptSrc: ["'self'", "'unsafe-inline'"], // Only if necessary
    styleSrc: ["'self'", "'unsafe-inline'"],
    imgSrc: ["'self'", 'data:', 'https:'],
    fontSrc: ["'self'"],
    objectSrc: ["'none'"],
    mediaSrc: ["'self'"],
    frameSrc: ["'none'"],
    connectSrc: ["'self'", 'https://api.biometrics.local'],
    baseUri: ["'self'"],
    formAction: ["'self'"],
    frameAncestors: ["'none'"],
  },
}));

// Report violations
app.use(helmet.contentSecurityPolicy({
  directives: {
    // ... other directives
    reportUri: '/api/security/csp-report',
  },
}));
```

#### 4.3.3 React XSS Prevention

```typescript
// ✅ SECURE: Using textContent instead of innerHTML
function SafeDisplay({ content }: { content: string }) {
  const element = document.createElement('div');
  element.textContent = content; // Safe - automatically escaped
  return element;
}

// ✅ SECURE: Using libraries that auto-escape
import DOMPurify from 'dompurify';

function SanitizedDisplay({ html }: { html: string }) {
  const sanitized = DOMPurify.sanitize(html, {
    ALLOWED_TAGS: ['b', 'i', 'em', 'strong', 'p', 'br'],
    ALLOWED_ATTR: [],
  });
  return <div dangerouslySetInnerHTML={{ __html: sanitized }} />;
}

// ❌ INSECURE: Never use dangerouslySetInnerHTML with user input
function UnsafeDisplay({ html }: { html: string }) {
  return <div dangerouslySetInnerHTML={{ __html: html }} />;
}
```

### 4.4 CSRF Protection

#### 4.4.1 CSRF Token Implementation

```typescript
// src/middleware/csrf.ts
import { Request, Response, NextFunction } from 'express';
import crypto from 'crypto';

interface CSRFConfig {
  cookieName: string;
  headerName: string;
  tokenLength: number;
}

const csrfConfig: CSRFConfig = {
  cookieName: 'XSRF-TOKEN',
  headerName: 'X-XSRF-TOKEN',
  tokenLength: 32,
};

export function generateCSRFToken(): string {
  return crypto.randomBytes(csrfConfig.tokenLength).toString('hex');
}

export function csrfMiddleware(req: Request, res: Response, next: NextFunction) {
  // Generate token for new sessions
  if (!req.cookies[csrfConfig.cookieName]) {
    const token = generateCSRFToken();
    res.cookie(csrfConfig.cookieName, token, {
      httpOnly: false, // Must be accessible to JavaScript
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'strict',
      maxAge: 3600000, // 1 hour
    });
    req.csrfToken = token;
  } else {
    req.csrfToken = req.cookies[csrfConfig.cookieName];
  }
  
  // Validate token on state-changing requests
  if (['POST', 'PUT', 'DELETE', 'PATCH'].includes(req.method)) {
    const requestToken = req.headers[csrfConfig.headerName] as string;
    
    if (!requestToken || requestToken !== req.cookies[csrfConfig.cookieName]) {
      return res.status(403).json({
        error: 'CSRF token validation failed',
        code: 'CSRF_INVALID_TOKEN',
      });
    }
  }
  
  next();
}

declare global {
  namespace Express {
    interface Request {
      csrfToken?: string;
    }
  }
}
```

#### 4.4.2 Double Submit Cookie Pattern

```typescript
// For APIs that don't use cookies
export function validateCSRFToken(req: Request, res: Response): boolean {
  const cookieToken = req.cookies[csrfConfig.cookieName];
  const headerToken = req.headers[csrfConfig.headerName] as string;
  
  if (!cookieToken || !headerToken) {
    return false;
  }
  
  // Constant-time comparison to prevent timing attacks
  return crypto.timingSafeEqual(
    Buffer.from(cookieToken),
    Buffer.from(headerToken)
  );
}

// Usage in routes
app.post('/api/data', 
  validateCSRFToken, 
  (req, res) => {
    // Process request
  }
);
```

### 4.5 Rate Limiting

```typescript
// src/middleware/rate-limit.ts
import rateLimit from 'express-rate-limit';
import RedisStore from 'rate-limit-redis';
import Redis from 'ioredis';
import { Request, Response } from 'express';

const redis = new Redis(process.env.REDIS_URL || 'redis://localhost:6379');

// General API rate limiter
export const apiLimiter = rateLimit({
  store: new RedisStore({
    sendCommand: (...args: string[]) => redis.call(...args),
    prefix: 'rl:api:',
  }),
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 1000, // 1000 requests per window
  standardHeaders: true,
  legacyHeaders: false,
  message: {
    error: 'Too many requests',
    code: 'RATE_LIMIT_EXCEEDED',
    retryAfter: 900, // seconds
  },
});

// Authentication endpoints - stricter limits
export const authLimiter = rateLimit({
  store: new RedisStore({
    sendCommand: (...args: string[]) => redis.call(...args),
    prefix: 'rl:auth:',
  }),
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 5, // 5 attempts per window
  standardHeaders: true,
  legacyHeaders: false,
  skipSuccessfulRequests: true, // Only count failed attempts
  message: {
    error: 'Too many authentication attempts',
    code: 'AUTH_RATE_LIMIT_EXCEEDED',
    retryAfter: 900,
  },
});

// Biometric operations - strict limits due to sensitivity
export const biometricLimiter = rateLimit({
  store: new RedisStore({
    sendCommand: (...args: string[]) => redis.call(...args),
    prefix: 'rl:biometric:',
  }),
  windowMs: 60 * 1000, // 1 minute
  max: 10, // 10 requests per minute
  standardHeaders: true,
  legacyHeaders: false,
  message: {
    error: 'Rate limit exceeded for biometric operations',
    code: 'BIOMETRIC_RATE_LIMIT',
    retryAfter: 60,
  },
});

// IP-based limiter for extreme cases
export const ipLimiter = rateLimit({
  store: new RedisStore({
    sendCommand: (...args: string[]) => redis.call(...args),
    prefix: 'rl:ip:',
  }),
  windowMs: 60 * 1000, // 1 minute
  max: 100, // 100 requests per minute per IP
  standardHeaders: true,
  legacyHeaders: false,
  keyGenerator: (req: Request) => req.ip || 'unknown',
});
```

---

## 5. AUTHENTICATION AND AUTHORIZATION

### 5.1 JWT Best Practices

#### 5.1.1 JWT Structure

```typescript
// src/lib/jwt/types.ts
interface JWTPayload {
  sub: string; // User ID
  email: string;
  role: string;
  permissions: string[];
  iat: number; // Issued at
  exp: number; // Expiration
  iss: string; // Issuer
  aud: string; // Audience
}

interface BiometricJWTPayload extends JWTPayload {
  biometricId: string;
  biometricType: 'fingerprint' | 'face' | 'iris';
  sessionId: string;
}
```

#### 5.1.2 JWT Generation

```typescript
// src/lib/jwt/generate.ts
import jwt from 'jsonwebtoken';
import { randomBytes } from 'crypto';

interface TokenConfig {
  secret: string;
  algorithm: 'RS256' | 'ES256';
  accessTokenExpiry: string;
  refreshTokenExpiry: string;
}

export class JWTGenerator {
  private accessTokenSecret: string;
  private refreshTokenSecret: string;
  private algorithm: 'RS256' | 'ES256';

  constructor(config: TokenConfig) {
    this.accessTokenSecret = config.secret;
    this.refreshTokenSecret = config.secret + '_refresh';
    this.algorithm = config.algorithm;
  }

  generateAccessToken(payload: Omit<JWTPayload, 'iat' | 'exp'>): string {
    const token = jwt.sign(payload, this.accessTokenSecret, {
      algorithm: this.algorithm,
      expiresIn: '15m', // Short-lived access tokens
      issuer: 'biometrics-api',
      audience: 'biometrics-app',
      jwtid: randomBytes(16).toString('hex'), // Unique token ID
    });
    
    return token;
  }

  generateRefreshToken(payload: { sub: string; email: string }): string {
    const token = jwt.sign(payload, this.refreshTokenSecret, {
      algorithm: this.algorithm,
      expiresIn: '7d', // Longer-lived refresh tokens
      issuer: 'biometrics-api',
      audience: 'biometrics-app',
      jwtid: randomBytes(16).toString('hex'),
    });
    
    return token;
  }

  generateTokenPair(
    user: User
  ): { accessToken: string; refreshToken: string } {
    const payload = {
      sub: user.id,
      email: user.email,
      role: user.role,
      permissions: user.permissions,
    };

    return {
      accessToken: this.generateAccessToken(payload),
      refreshToken: this.generateRefreshToken(payload),
    };
  }
}
```

#### 5.1.3 JWT Validation

```typescript
// src/lib/jwt/verify.ts
import jwt, { JsonWebTokenError, TokenExpiredError } from 'jsonwebtoken';

interface JWTVerificationResult {
  valid: boolean;
  payload?: JWTPayload;
  error?: {
    code: string;
    message: string;
  };
}

export class JWTVerifier {
  constructor(private secret: string) {}

  verifyAccessToken(token: string): JWTVerificationResult {
    try {
      const payload = jwt.verify(token, this.secret, {
        algorithms: ['RS256', 'ES256'],
        issuer: 'biometrics-api',
        audience: 'biometrics-app',
        complete: true,
      }) as JWTPayload;

      // Additional security checks
      if (!payload.permissions || !Array.isArray(payload.permissions)) {
        return {
          valid: false,
          error: {
            code: 'INVALID_TOKEN',
            message: 'Token missing permissions',
          },
        };
      }

      return { valid: true, payload };
    } catch (error) {
      if (error instanceof TokenExpiredError) {
        return {
          valid: false,
          error: {
            code: 'TOKEN_EXPIRED',
            message: 'Access token has expired',
          },
        };
      }

      if (error instanceof JsonWebTokenError) {
        return {
          valid: false,
          error: {
            code: 'INVALID_TOKEN',
            message: 'Invalid token',
          },
        };
      }

      return {
        valid: false,
        error: {
          code: 'TOKEN_ERROR',
          message: 'Token verification failed',
        },
      };
    }
  }
}
```

### 5.2 OAuth2 Flows

#### 5.2.1 Authorization Code Flow

```typescript
// src/routes/oauth.ts
import { Router, Request, Response } from 'express';
import axios from 'axios';
import crypto from 'crypto';

const router = Router();

interface OAuthState {
  redirectUri: string;
  codeVerifier?: string;
  nonce?: string;
}

// In-memory state storage (use Redis in production)
const oauthStates = new Map<string, OAuthState>();

router.get('/oauth/authorize', async (req: Request, res: Response) => {
  const {
    client_id,
    redirect_uri,
    response_type,
    scope,
    state,
  } = req.query;

  // Validate client_id
  const client = await validateClient(client_id as string);
  if (!client) {
    return res.status(400).json({ error: 'invalid_client' });
  }

  // Generate state
  const stateId = crypto.randomBytes(32).toString('hex');
  oauthStates.set(stateId, {
    redirectUri: redirect_uri as string,
  });

  // Build authorization URL
  const authUrl = new URL('https://auth.biometrics.local/oauth/authorize');
  authUrl.searchParams.set('client_id', client_id as string);
  authUrl.searchParams.set('redirect_uri', redirect_uri as string);
  authUrl.searchParams.set('response_type', response_type as string);
  authUrl.searchParams.set('scope', scope as string);
  authUrl.searchParams.set('state', stateId);

  res.redirect(authUrl.toString());
});

router.get('/oauth/callback', async (req: Request, res: Response) => {
  const { code, state, error } = req.query;

  if (error) {
    return res.redirect(`${state?.toString()}?error=${error}`);
  }

  // Validate state
  const oauthState = oauthStates.get(state as string);
  if (!oauthState) {
    return res.status(400).json({ error: 'invalid_state' });
  }

  // Exchange code for tokens
  const tokenResponse = await axios.post(
    'https://auth.biometrics.local/oauth/token',
    {
      grant_type: 'authorization_code',
      code,
      client_id: process.env.OAUTH_CLIENT_ID,
      client_secret: process.env.OAUTH_CLIENT_SECRET,
      redirect_uri: oauthState.redirectUri,
    }
  );

  // Store tokens securely
  const { access_token, refresh_token, expires_in } = tokenResponse.data;
  
  // Redirect back to application with auth code
  res.redirect(`${oauthState.redirectUri}?code=${access_token}`);
});
```

### 5.3 Session Management

```typescript
// src/lib/session/manager.ts
import Redis from 'ioredis';
import crypto from 'crypto';

interface SessionData {
  userId: string;
  email: string;
  role: string;
  createdAt: number;
  lastActivity: number;
  ipAddress: string;
  userAgent: string;
  biometricVerified: boolean;
}

export class SessionManager {
  private redis: Redis;
  private readonly SESSION_PREFIX = 'session:';
  private readonly SESSION_TTL = 3600; // 1 hour

  constructor(redisUrl: string) {
    this.redis = new Redis(redisUrl);
  }

  async createSession(data: Omit<SessionData, 'createdAt' | 'lastActivity'>): Promise<string> {
    const sessionId = crypto.randomBytes(32).toString('hex');
    
    const sessionData: SessionData = {
      ...data,
      createdAt: Date.now(),
      lastActivity: Date.now(),
    };

    await this.redis.setex(
      `${this.SESSION_PREFIX}${sessionId}`,
      this.SESSION_TTL,
      JSON.stringify(sessionData)
    );

    return sessionId;
  }

  async getSession(sessionId: string): Promise<SessionData | null> {
    const data = await this.redis.get(`${this.SESSION_PREFIX}${sessionId}`);
    return data ? JSON.parse(data) : null;
  }

  async refreshSession(sessionId: string): Promise<void> {
    await this.redis.expire(`${this.SESSION_PREFIX}${sessionId}`, this.SESSION_TTL);
  }

  async destroySession(sessionId: string): Promise<void> {
    await this.redis.del(`${this.SESSION_PREFIX}${sessionId}`);
  }

  async destroyAllUserSessions(userId: string): Promise<void> {
    const keys = await this.redis.keys(`${this.SESSION_PREFIX}*`);
    const userSessions = keys.filter(async (key) => {
      const data = await this.redis.get(key);
      return data && JSON.parse(data).userId === userId;
    });
    
    await Promise.all(
      userSessions.map(key => this.redis.del(key))
    );
  }
}
```

### 5.4 Role-Based Access Control (RBAC)

```typescript
// src/lib/rbac/types.ts
enum Permission {
  // User management
  USER_READ = 'user:read',
  USER_CREATE = 'user:create',
  USER_UPDATE = 'user:update',
  USER_DELETE = 'user:delete',
  
  // Biometric operations
  BIOMETRIC_ENROLL = 'biometric:enroll',
  BIOMETRIC_VERIFY = 'biometric:verify',
  BIOMETRIC_DELETE = 'biometric:delete',
  BIOMETRIC_EXPORT = 'biometric:export',
  
  // Admin operations
  ADMIN_DASHBOARD = 'admin:dashboard',
  ADMIN_AUDIT = 'admin:audit',
  ADMIN_SETTINGS = 'admin:settings',
}

enum Role {
  USER = 'user',
  ADMIN = 'admin',
  SECURITY_OFFICER = 'security_officer',
  AUDITOR = 'auditor',
  SYSTEM = 'system',
}

const RolePermissions: Record<Role, Permission[]> = {
  [Role.USER]: [
    Permission.USER_READ,
    Permission.BIOMETRIC_ENROLL,
    Permission.BIOMETRIC_VERIFY,
  ],
  [Role.ADMIN]: [
    Permission.USER_READ,
    Permission.USER_CREATE,
    Permission.USER_UPDATE,
    Permission.BIOMETRIC_ENROLL,
    Permission.BIOMETRIC_VERIFY,
    Permission.BIOMETRIC_DELETE,
    Permission.ADMIN_DASHBOARD,
  ],
  [Role.SECURITY_OFFICER]: [
    Permission.USER_READ,
    Permission.BIOMETRIC_ENROLL,
    Permission.BIOMETRIC_VERIFY,
    Permission.BIOMETRIC_DELETE,
    Permission.BIOMETRIC_EXPORT,
    Permission.ADMIN_DASHBOARD,
    Permission.ADMIN_AUDIT,
  ],
  [Role.AUDITOR]: [
    Permission.USER_READ,
    Permission.ADMIN_AUDIT,
  ],
  [Role.SYSTEM]: [
    Permission.USER_CREATE,
    Permission.USER_UPDATE,
    Permission.BIOMETRIC_ENROLL,
    Permission.BIOMETRIC_VERIFY,
    Permission.ADMIN_AUDIT,
  ],
};
```

#### 5.4.1 RBAC Middleware

```typescript
// src/middleware/rbac.ts
import { Request, Response, NextFunction } from 'express';
import { Permission, RolePermissions } from '../lib/rbac/types';

export function requirePermission(...permissions: Permission[]) {
  return (req: Request, res: Response, next: NextFunction) => {
    const user = req.user;
    
    if (!user) {
      return res.status(401).json({
        error: 'Unauthorized',
        code: 'NO_AUTHENTICATION',
      });
    }

    const userPermissions = RolePermissions[user.role] || [];
    
    const hasAllPermissions = permissions.every(
      perm => userPermissions.includes(perm)
    );

    if (!hasAllPermissions) {
      return res.status(403).json({
        error: 'Forbidden',
        code: 'INSUFFICIENT_PERMISSIONS',
        required: permissions,
        userRole: user.role,
      });
    }

    next();
  };
}

export function requireRole(...roles: string[]) {
  return (req: Request, res: Response, next: NextFunction) => {
    const user = req.user;
    
    if (!user) {
      return res.status(401).json({
        error: 'Unauthorized',
        code: 'NO_AUTHENTICATION',
      });
    }

    if (!roles.includes(user.role)) {
      return res.status(403).json({
        error: 'Forbidden',
        code: 'INSUFFICIENT_ROLE',
        required: roles,
        userRole: user.role,
      });
    }

    next();
  };
}
```

---

## 6. DATA PROTECTION

### 6.1 Encryption at Rest

#### 6.1.1 Database Encryption

```sql
-- PostgreSQL TDE (Transparent Data Encryption)
-- Enable at database level
ALTER SYSTEM SET shared_preload_libraries = 'pg_tde';
SELECT pg_reload_conf();

-- For specific columns with application-level encryption
CREATE EXTENSION pgcrypto;

-- Encrypted column example
CREATE TABLE biometric_templates (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  template_data BYTEA ENCRYPTED, -- Application-level encrypted
  template_hash VARCHAR(255),     -- For verification without decryption
  biometric_type VARCHAR(50),
  created_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### 6.1.2 Application-Level Encryption

```typescript
// src/lib/encryption/biometric.ts
import crypto from 'crypto';

interface EncryptionConfig {
  algorithm: 'aes-256-gcm';
  keyLength: 32;
  ivLength: 16;
  tagLength: 16;
  authTagLength: number;
}

const config: EncryptionConfig = {
  algorithm: 'aes-256-gcm',
  keyLength: 32,
  ivLength: 16,
  tagLength: 16,
  authTagLength: 128,
};

export class BiometricEncryption {
  private key: Buffer;

  constructor(keyHex: string) {
    if (keyHex.length !== 64) {
      throw new Error('Encryption key must be 256 bits (64 hex characters)');
    }
    this.key = Buffer.from(keyHex, 'hex');
  }

  encrypt(plaintext: Buffer): EncryptedData {
    // Generate random IV
    const iv = crypto.randomBytes(config.ivLength);
    
    // Create cipher
    const cipher = crypto.createCipheriv(
      config.algorithm,
      this.key,
      iv,
      {
        authTagLength: config.authTagLength / 8,
      }
    );

    // Encrypt
    const encrypted = Buffer.concat([
      cipher.update(plaintext),
      cipher.final(),
    ]);

    // Get auth tag
    const authTag = cipher.getAuthTag();

    return {
      encrypted,
      iv,
      authTag,
      algorithm: config.algorithm,
    };
  }

  decrypt(data: EncryptedData): Buffer {
    const decipher = crypto.createDecipheriv(
      data.algorithm,
      this.key,
      data.iv,
      {
        authTagLength: config.authTagLength / 8,
      }
    );

    decipher.setAuthTag(data.authTag);

    return Buffer.concat([
      decipher.update(data.encrypted),
      decipher.final(),
    ]);
  }
}

interface EncryptedData {
  encrypted: Buffer;
  iv: Buffer;
  authTag: Buffer;
  algorithm: string;
}
```

### 6.2 Encryption in Transit

#### 6.2.1 TLS Configuration

```typescript
// src/config/tls.ts
import https from 'https';
import fs from 'fs';

interface TLSConfig {
  key: Buffer;
  cert: Buffer;
  ca?: Buffer;
  minVersion: string;
  ciphers: string[];
  honorCipherOrder: boolean;
  sessionTickets: boolean;
  sessionTimeout: number;
}

export function createHTTPSServer(app: Express.Application): https.Server {
  const config: TLSConfig = {
    key: fs.readFileSync('/etc/ssl/private/biometrics.key'),
    cert: fs.readFileSync('/etc/ssl/certs/biometrics.crt'),
    ca: fs.readFileSync('/etc/ssl/certs/ca-bundle.crt'),
    minVersion: 'TLSv1.2',
    ciphers: [
      'ECDHE-ECDSA-AES256-GCM-SHA384',
      'ECDHE-RSA-AES256-GCM-SHA384',
      'ECDHE-ECDSA-CHACHA20-POLY1305',
      'ECDHE-RSA-CHACHA20-POLY1305',
      'ECDHE-ECDSA-AES128-GCM-SHA256',
      'ECDHE-RSA-AES128-GCM-SHA256',
    ].join(':'),
    honorCipherOrder: true,
    sessionTickets: false,
    sessionTimeout: 600,
  };

  return https.createServer(config, app);
}
```

#### 6.2.2 Certificate Pinning

```typescript
// src/middleware/certificate-pinning.ts
import { Request, Response, NextFunction } from 'express';

interface PinnedCertificate {
  publicKeyHash: string; // SHA-256 hash of public key
  issuer: string;
}

const PINNED_CERTIFICATES: PinnedCertificate[] = [
  {
    publicKeyHash: 'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=',
    issuer: 'DigiCert',
  },
];

export function certificatePinning(req: Request, res: Response, next: NextFunction) {
  // In production, implement proper certificate pinning
  // This is a simplified example
  const cert = req.socket.getPeerCertificate();
  
  if (!cert || !cert.pubkey) {
    return res.status(495).json({
      error: 'SSL Certificate Error',
      code: 'CERTIFICATE_INVALID',
    });
  }

  const publicKeyHash = crypto
    .createHash('sha256')
    .update(cert.pubkey)
    .digest('base64');

  const isPinned = PINNED_CERTIFICATES.some(
    pin => pin.publicKeyHash === publicKeyHash
  );

  if (!isPinned && process.env.NODE_ENV === 'production') {
    return res.status(495).json({
      error: 'SSL Certificate Error',
      code: 'CERTIFICATE_NOT_PINNED',
    });
  }

  next();
}
```

### 6.3 GDPR Compliance

```typescript
// src/lib/gdpr/compliance.ts
interface GDPRRequirements {
  dataMinimization: boolean;
  purposeLimitation: boolean;
  storageLimitation: boolean;
  consentManagement: boolean;
  dataPortability: boolean;
  rightToErasure: boolean;
}

export class GDPRCompliance {
  private readonly RETENTION_PERIODS = {
    biometric_template: '7_years', // Per BSI requirements
    session_data: '90_days',
    audit_logs: '10_years',
    transaction_logs: '7_years',
  };

  async processDataSubjectRequest(
    requestType: 'access' | 'rectification' | 'erasure' | 'portability',
    userId: string,
    data: any
  ): Promise<any> {
    switch (requestType) {
      case 'access':
        return this.exportUserData(userId);
      
      case 'rectification':
        return this.correctUserData(userId, data);
      
      case 'erasure':
        return this.eraseUserData(userId);
      
      case 'portability':
        return this.exportUserDataPortable(userId);
      
      default:
        throw new Error('Invalid request type');
    }
  }

  private async exportUserData(userId: string): Promise<ExportData> {
    const user = await db.user.findUnique({ where: { id: userId } });
    const biometrics = await db.biometricTemplate.findMany({
      where: { userId },
    });
    const auditLogs = await db.auditLog.findMany({
      where: { userId },
      take: 1000, // Limit
    });

    return {
      personalData: user,
      biometricData: biometrics,
      auditLogs,
      exportedAt: new Date().toISOString(),
    };
  }

  private async eraseUserData(userId: string): Promise<ErasureResult> {
    // Note: Biometric templates may need to be retained for legal reasons
    // Implement based on specific legal requirements
    
    const result = await db.$transaction([
      db.user.update({
        where: { id: userId },
        data: {
          email: `deleted-${userId}@biometrics.local`,
          name: 'Deleted User',
          deletedAt: new Date(),
        },
      }),
      db.biometricTemplate.updateMany({
        where: { userId },
        data: {
          active: false,
          deletedAt: new Date(),
        },
      }),
      db.auditLog.create({
        data: {
          userId,
          action: 'GDPR_ERASURE',
          timestamp: new Date(),
        },
      }),
    ]);

    return {
      success: true,
      erasedAt: new Date().toISOString(),
      retainedData: ['biometric_templates_legal'],
    };
  }
}
```

---

## 7. INFRASTRUCTURE SECURITY

### 7.1 Docker Security Best Practices

```dockerfile
# Dockerfile - SECURE
FROM node:20-alpine AS builder

# Create non-root user
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci --only=production

# Copy source
COPY --chown=nodejs:nodejs . .

# Switch to non-root user
USER nodejs

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD node -e "require('http').get('http://localhost:3000/health', (r) => process.exit(r.statusCode === 200 ? 0 : 1))"

# Run as non-root
CMD ["node", "dist/index.js"]
```

#### 7.1.1 Docker Security Options

```yaml
# docker-compose.yml - SECURE
services:
  biometrics-api:
    build: .
    image: biometrics/api:latest
    container_name: biometrics-api
    restart: unless-stopped
    
    security_opt:
      - no-new-privileges:true
      - seccomp:unconfined  # Only if needed
    
    cap_drop:
      - ALL
    
    # Read-only filesystem
    read_only: true
    tmpfs:
      - /tmp
      - /run
    
    # Resource limits
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '0.5'
          memory: 512M
    
    # Health check
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:3000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    
    networks:
      - biometrics-network
    
    secrets:
      - db_password
      - jwt_secret

networks:
  biometrics-network:
    driver: bridge
    driver_opts:
      com.docker.network.bridge.name: br-biometrics
      com.docker.network.enable_ipv6: "false"

secrets:
  db_password:
    file: ./secrets/db_password.txt
  jwt_secret:
    file: ./secrets/jwt_secret.txt
```

### 7.2 Network Segmentation

```yaml
# kubernetes-network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: biometrics-api-policy
  namespace: biometrics
spec:
  podSelector:
    matchLabels:
      app: biometrics-api
  policyTypes:
    - Ingress
    - Egress
  ingress:
    # Allow from ingress controller
    - from:
        - namespaceSelector:
            matchLabels:
              name: ingress-nginx
      ports:
        - protocol: TCP
          port: 3000
    # Allow from auth service
    - from:
        - podSelector:
            matchLabels:
              app: biometrics-auth
      ports:
        - protocol: TCP
          port: 3000
  egress:
    # Allow DNS
    - to:
        - namespaceSelector: {}
          podSelector:
            matchLabels:
              k8s-app: kube-dns
      ports:
        - protocol: UDP
          port: 53
    # Allow to database
    - to:
        - podSelector:
            matchLabels:
              app: biometrics-db
      ports:
        - protocol: TCP
          port: 5432
    # Allow to Redis
    - to:
        - podSelector:
            matchLabels:
              app: biometrics-cache
      ports:
        - protocol: TCP
          port: 6379
    # Allow to internal services
    - to:
        - podSelector:
            matchLabels:
              tier: internal
      ports:
        - protocol: TCP
          port: 8080
    # Block all other egress
    - to: []
```

### 7.3 Firewall Rules

```bash
#!/bin/bash
# scripts/configure-firewall.sh

set -euo pipefail

echo "Configuring firewall rules..."

# Flush existing rules
iptables -F
iptables -X

# Default policies
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT ACCEPT

# Allow loopback
iptables -A INPUT -i lo -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT

# Allow established connections
iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT

# Allow SSH (limit to specific IPs)
iptables -A INPUT -p tcp --dport 22 -s 10.0.0.0/8 -j ACCEPT
iptables -A INPUT -p tcp --dport 22 -s 172.16.0.0/12 -j ACCEPT

# Allow HTTP/HTTPS
iptables -A INPUT -p tcp --dport 80 -j ACCEPT
iptables -A INPUT -p tcp --dport 443 -j ACCEPT

# Allow health checks
iptables -A INPUT -p tcp --dport 3000 -s 10.0.0.0/8 -j ACCEPT
iptables -A INPUT -p tcp --dport 3000 -s 172.16.0.0/12 -j ACCEPT

# Rate limiting
iptables -A INPUT -p tcp --dport 22 -m state --state NEW -m recent --set
iptables -A INPUT -p tcp --dport 22 -m state --state NEW -m recent --update --seconds 60 --hitcount 4 -j DROP

# Log dropped packets
iptables -A INPUT -m limit --limit 5/min -j LOG --log-prefix "iptables-dropped: " --log-level 4

# Save rules
iptables-save > /etc/iptables/rules.v4

echo "Firewall configuration complete"
```

---

## 8. MONITORING AND AUDITING

### 8.1 Audit Logging

```typescript
// src/lib/audit/logger.ts
import { createLogger, format, transports } from 'winston';
import { v4 as uuidv4 } from 'uuid';

enum AuditAction {
  // Authentication
  LOGIN = 'AUTH_LOGIN',
  LOGOUT = 'AUTH_LOGOUT',
  LOGIN_FAILED = 'AUTH_LOGIN_FAILED',
  PASSWORD_CHANGED = 'AUTH_PASSWORD_CHANGED',
  MFA_ENABLED = 'AUTH_MFA_ENABLED',
  MFA_DISABLED = 'AUTH_MFA_DISABLED',
  
  // Authorization
  PERMISSION_GRANTED = 'AUTH_PERMISSION_GRANTED',
  PERMISSION_REVOKED = 'AUTH_PERMISSION_REVOKED',
  ROLE_CHANGED = 'AUTH_ROLE_CHANGED',
  
  // Biometric
  BIOMETRIC_ENROLLED = 'BIOMETRIC_ENROLLED',
  BIOMETRIC_VERIFIED = 'BIOMETRIC_VERIFIED',
  BIOMETRIC_DELETED = 'BIOMETRIC_DELETED',
  BIOMETRIC_EXPORTED = 'BIOMETRIC_EXPORTED',
  
  // Data Access
  DATA_ACCESSED = 'DATA_ACCESSED',
  DATA_EXPORTED = 'DATA_EXPORTED',
  DATA_DELETED = 'DATA_DELETED',
}

interface AuditEntry {
  id: string;
  timestamp: string;
  action: AuditAction;
  userId: string;
  userEmail: string;
  userRole: string;
  ipAddress: string;
  userAgent: string;
  resource: string;
  resourceId: string;
  result: 'success' | 'failure';
  metadata: Record<string, any>;
}

export class AuditLogger {
  private logger;

  constructor() {
    this.logger = createLogger({
      level: 'info',
      format: format.combine(
        format.timestamp(),
        format.json()
      ),
      transports: [
        new transports.File({ 
          filename: 'audit.log',
          dirname: '/var/log/biometrics',
        }),
        new transports.Console(),
      ],
    });
  }

  log(entry: Omit<AuditEntry, 'id' | 'timestamp'>): void {
    const fullEntry: AuditEntry = {
      ...entry,
      id: uuidv4(),
      timestamp: new Date().toISOString(),
    };

    this.logger.info({
      type: 'audit',
      ...fullEntry,
    });
  }

  // Authentication logging
  logLogin(params: {
    userId: string;
    email: string;
    role: string;
    ipAddress: string;
    userAgent: string;
    success: boolean;
    mfaUsed?: boolean;
  }): void {
    this.log({
      action: params.success ? AuditAction.LOGIN : AuditAction.LOGIN_FAILED,
      userId: params.userId,
      userEmail: params.email,
      userRole: params.role,
      ipAddress: params.ipAddress,
      userAgent: params.userAgent,
      resource: 'auth',
      resourceId: params.userId,
      result: params.success ? 'success' : 'failure',
      metadata: {
        mfaUsed: params.mfaUsed,
      },
    });
  }

  // Biometric operations logging
  logBiometricOperation(params: {
    userId: string;
    email: string;
    role: string;
    action: AuditAction;
    biometricId: string;
    biometricType: string;
    ipAddress: string;
    success: boolean;
  }): void {
    this.log({
      action: params.action,
      userId: params.userId,
      userEmail: params.email,
      userRole: params.role,
      ipAddress: params.ipAddress,
      userAgent: '', // Populate from request
      resource: 'biometric',
      resourceId: params.biometricId,
      result: params.success ? 'success' : 'failure',
      metadata: {
        biometricType: params.biometricType,
      },
    });
  }
}
```

### 8.2 Security Event Monitoring

```yaml
# prometheus-alerts.yaml
groups:
  - name: security
    interval: 30s
    rules:
      - alert: HighFailedLogins
        expr: |
          sum(rate(biometrics_auth_login_failed_total[5m])) 
          > 10
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "High number of failed login attempts"
          description: "More than 10 failed logins per minute detected"
          
      - alert: UnusualAccessPattern
        expr: |
          sum(rate(biometrics_api_requests_total{status=~"5.."}[5m])) 
          / sum(rate(biometrics_api_requests_total[5m])) > 0.05
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Unusual access pattern detected"
          description: "High error rate may indicate an attack"
          
      - alert: RateLimitExceeded
        expr: |
          sum(rate(biometrics_rate_limit_exceeded_total[5m])) 
          > 100
        for: 1m
        labels:
          severity: warning
        annotations:
          summary: "Rate limit exceeded frequently"
          description: "More than 100 rate limit violations per minute"
          
      - alert: SuspiciousIP
        expr: |
          topk(10, sum by (ip) (rate(biometrics_api_requests_total[5m]))) 
          > 1000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Suspicious IP address"
          description: "IP address making excessive requests"
          
      - alert: DatabaseConnectionSpike
        expr: |
          biometrics_db_connections_active > 80
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "Database connection spike"
          description: "Active connections exceed 80% of limit"
```

---

## 9. OWASP TOP 10

### 9.1 A01:2021 - Broken Access Control

```typescript
// ❌ INSECURE: Missing authorization checks
app.get('/api/users/:id', async (req, res) => {
  const user = await db.user.findUnique({
    where: { id: req.params.id },
  });
  res.json(user); // Anyone can access any user!
});

// ✅ SECURE: Proper authorization
app.get(
  '/api/users/:id',
  requireAuth,
  requirePermission(Permission.USER_READ),
  async (req, res) => {
    // Users can only access their own data unless admin
    if (req.user.role !== 'admin' && req.params.id !== req.user.id) {
      return res.status(403).json({
        error: 'Access denied',
        code: 'FORBIDDEN',
      });
    }
    
    const user = await db.user.findUnique({
      where: { id: req.params.id },
      select: {
        id: true,
        email: true,
        name: true,
        role: true,
        // Exclude sensitive fields
      },
    });
    
    res.json(user);
  }
);
```

**Prevention:**
- Implement proper authorization at every layer
- Deny by default
- Log access control failures
- Rate limit API access

### 9.2 A02:2021 - Cryptographic Failures

```typescript
// ❌ INSECURE: Weak encryption algorithm
const cipher = crypto.createCipher('aes-256', key); // Deprecated!
const encrypted = cipher.update(data, 'utf8', 'hex');

// ✅ SECURE: Strong encryption
const cipher = crypto.createCipheriv('aes-256-gcm', key, iv);
const encrypted = Buffer.concat([
  cipher.update(data, 'utf8'),
  cipher.final(),
]);
const authTag = cipher.getAuthTag();
```

**Prevention:**
- Use strong encryption algorithms (AES-256-GCM)
- Ensure proper key management
- Disable weak protocols (TLS 1.0, TLS 1.1)
- Use secure random number generators

### 9.3 A03:2021 - Injection

```typescript
// ❌ INSECURE: SQL Injection
const query = `SELECT * FROM users WHERE id = ${req.params.id}`;

// ✅ SECURE: Parameterized query
const query = 'SELECT * FROM users WHERE id = $1';
const result = await db.query(query, [req.params.id]);

// ✅ SECURE: ORM
const user = await db.user.findUnique({
  where: { id: req.params.id },
});
```

**Prevention:**
- Use parameterized queries
- Use ORMs or stored procedures
- Validate and sanitize input
- Use LIMIT and other SQL controls

### 9.4 A04:2021 - Insecure Design

```typescript
// ❌ INSECURE: No rate limiting on password reset
app.post('/api/auth/reset-password', async (req, res) => {
  await sendPasswordResetEmail(req.body.email); // Unlimited requests!
  res.json({ success: true });
});

// ✅ SECURE: Rate limiting and design
const passwordResetLimiter = rateLimit({
  windowMs: 3600000, // 1 hour
  max: 5, // 5 requests per hour
  message: {
    error: 'Too many password reset requests',
    retryAfter: 3600,
  },
});

app.post(
  '/api/auth/reset-password',
  passwordResetLimiter,
  async (req, res) => {
    // Add additional safeguards
    await validateEmailNotBlocked(req.body.email);
    await sendPasswordResetEmail(req.body.email);
    res.json({ success: true });
  }
);
```

**Prevention:**
- Threat modeling during design
- Secure design patterns
- Rate limiting
- Account lockout policies

### 9.5 A05:2021 - Security Misconfiguration

```yaml
# ❌ INSECURE: Default configuration
server:
  port: 3000
  debug: true  # Exposes sensitive info!

# ✅ SECURE: Production configuration
server:
  port: ${SERVER_PORT}
  debug: ${DEBUG}  # Should be false in production!
  trustProxy: true
  cors:
    origin: ${ALLOWED_ORIGINS}
    credentials: true
  helmet:
    contentSecurityPolicy: true
    hsts: true
  logging:
    level: ${LOG_LEVEL}
    redact: ['password', 'token', 'secret']
```

**Prevention:**
- Disable debug mode in production
- Secure headers (Helmet.js)
- Disable unnecessary features
- Regular security audits

### 9.6 A06:2021 - Vulnerable Components

```json
// package.json - Audit dependencies
{
  "scripts": {
    "security:audit": "npm audit --audit-level=high",
    "security:update": "npm outdated --depth=0",
    "security:scan": "trufflehog filesystem ."
  },
  "devDependencies": {
    "npm-audit": "^3.0.0"
  }
}
```

**Prevention:**
- Regular dependency scanning
- Keep components updated
- Remove unused dependencies
- Monitor CVEs

### 9.7 A07:2021 - Authentication Failures

```typescript
// ❌ INSECURE: Weak password policy
const isValid = password.length >= 6;

// ✅ SECURE: Strong password policy
const passwordSchema = z.string()
  .min(12, 'Password must be at least 12 characters')
  .regex(/[A-Z]/, 'Password must contain uppercase letter')
  .regex(/[a-z]/, 'Password must contain lowercase letter')
  .regex(/[0-9]/, 'Password must contain number')
  .regex(/[^A-Za-z0-9]/, 'Password must contain special character');

// ✅ SECURE: Secure authentication
async function authenticate(email: string, password: string) {
  // Rate limiting
  await checkRateLimit(email);
  
  // Get user with hash
  const user = await getUserWithPasswordHash(email);
  
  // Verify password with constant-time comparison
  const valid = await bcrypt.compare(password, user.passwordHash);
  
  if (!valid) {
    await recordFailedAttempt(email);
    throw new AuthenticationError('Invalid credentials');
  }
  
  // Check account lockout
  if (await isAccountLocked(email)) {
    throw new AuthenticationError('Account locked');
  }
  
  // Generate secure session
  return generateSession(user);
}
```

**Prevention:**
- Strong password policies
- MFA implementation
- Account lockout policies
- Secure session management

### 9.8 A08:2021 - Software and Data Integrity Failures

```yaml
# CI/CD Pipeline Security
jobs:
  build:
    steps:
      - name: Verify checksums
        run: |
          echo "${{ secrets.CHECKSUM }}" | sha256sum -c -
      
      - name: Verify signatures
        run: |
          cosign verify --key ${{ secrets.COSIGN_KEY }} image
      
      - name: SBOM Generation
        uses: cyclonedx/cyclonedx-npm@v1
        
      - name: Dependency Review
        uses: actions/dependency-review-action@v4
```

**Prevention:**
- Verify integrity of dependencies
- Use signed containers
- Implement SBOM generation
- CI/CD pipeline security

### 9.9 A09:2021 - Security Logging Failures

```typescript
// ❌ INSECURE: No logging
app.post('/api/login', async (req, res) => {
  const user = await authenticate(req.body);
  res.json({ token: user.token });
});

// ✅ SECURE: Comprehensive logging
app.post('/api/login', async (req, res) => {
  const startTime = Date.now();
  
  try {
    const user = await authenticate(req.body);
    
    auditLogger.log({
      action: AuditAction.LOGIN,
      userId: user.id,
      result: 'success',
    });
    
    res.json({ token: user.token });
  } catch (error) {
    auditLogger.log({
      action: AuditAction.LOGIN_FAILED,
      userId: req.body.email,
      result: 'failure',
      metadata: { error: error.message },
    });
    
    res.status(401).json({ error: 'Authentication failed' });
  } finally {
    logger.info('Login request completed', {
      duration: Date.now() - startTime,
      ip: req.ip,
    });
  }
});
```

**Prevention:**
- Log authentication events
- Log access control failures
- Use structured logging
- Protect log integrity

### 9.10 A10:2021 - Server-Side Request Forgery (SSRF)

```typescript
// ❌ INSECURE: No URL validation
app.get('/api/fetch', async (req, res) => {
  const response = await fetch(req.query.url as string);
  const data = await response.json();
  res.json(data);
});

// ✅ SECURE: Strict URL validation
import { URL } from 'url';

const ALLOWED_DOMAINS = ['api.biometrics.local', 'internal.biometrics.local'];
const ALLOWED_PROTOCOLS = ['https'];

app.get('/api/fetch', async (req, res) => {
  const targetUrl = req.query.url as string;
  
  try {
    const parsedUrl = new URL(targetUrl);
    
    // Check protocol
    if (!ALLOWED_PROTOCOLS.includes(parsedUrl.protocol)) {
      throw new Error('Invalid protocol');
    }
    
    // Check hostname
    if (!ALLOWED_DOMAINS.includes(parsedUrl.hostname)) {
      throw new Error('Invalid hostname');
    }
    
    // Block private IPs
    const ip = await dns.resolve(parsedUrl.hostname);
    if (isPrivateIP(ip)) {
      throw new Error('Private IP not allowed');
    }
    
    const response = await fetch(parsedUrl.toString());
    const data = await response.json();
    res.json(data);
  } catch (error) {
    res.status(400).json({ error: 'Invalid request' });
  }
});

function isPrivateIP(ip: string): boolean {
  const parts = ip.split('.').map(Number);
  return (
    parts[0] === 10 ||
    (parts[0] === 172 && parts[1] >= 16 && parts[1] <= 31) ||
    (parts[0] === 192 && parts[1] === 168)
  );
}
```

**Prevention:**
- URL validation
- Disable HTTP redirections
- Use allowlists
- Block private IP addresses

---

## 10. SECURITY REVIEW CHECKLIST

### 10.1 Pre-Commit Security Checklist

- [ ] No secrets in code
- [ ] All inputs validated
- [ ] Parameterized queries used
- [ ] Output encoding applied
- [ ] CSRF tokens implemented
- [ ] Rate limiting applied
- [ ] Proper error handling
- [ ] No sensitive data in logs

### 10.2 Pre-Deployment Security Checklist

- [ ] All dependencies scanned
- [ ] No high/critical vulnerabilities
- [ ] TLS configured
- [ ] Security headers enabled
- [ ] Rate limiting enabled
- [ ] Audit logging enabled
- [ ] Health checks configured
- [ ] Container scanned
- [ ] Secrets properly configured
- [ ] Network policies applied

### 10.3 Production Security Checklist

- [ ] Monitoring active
- [ ] Alerts configured
- [ ] Log retention configured
- [ ] Backup verified
- [ ] Incident response plan tested
- [ ] Access review completed
- [ ] Certificate rotation scheduled
- [ ] Penetration testing completed

---

## 11. INCIDENT RESPONSE

### 11.1 Incident Classification

| Severity | Description | Response Time | Examples |
|----------|-------------|---------------|----------|
| **Critical** | Active breach, data exposed | Immediate | Database leak, ransomware |
| **High** | Potential breach, immediate action needed | 1 hour | Suspicious activity, failed attack attempt |
| **Medium** | Security policy violation | 4 hours | Unauthorized access attempt |
| **Low** | Minor security issue | 24 hours | Configuration drift |

### 11.2 Response Procedure

```bash
#!/bin/bash
# scripts/incident-response.sh

INCIDENT_ID=$(uuidgen)
INCIDENT_SEVERITY=$1
INCIDENT_TYPE=$2

echo "Starting incident response..."
echo "Incident ID: $INCIDENT_ID"
echo "Severity: $INCIDENT_SEVERITY"
echo "Type: $INCIDENT_TYPE"

# Step 1: Isolate affected systems
echo "[1/6] Isolating affected systems..."
./scripts/isolate-systems.sh

# Step 2: Preserve evidence
echo "[2/6] Preserving evidence..."
./scripts/preserve-evidence.sh $INCIDENT_ID

# Step 3: Notify security team
echo "[3/6] Notifying security team..."
curl -X POST "${SLACK_WEBHOOK_URL}" \
  -H 'Content-Type: application/json' \
  -d "{\"text\":\"🚨 INCIDENT: $INCIDENT_ID - $INCIDENT_SEVERITY\"}"

# Step 4: Begin investigation
echo "[4/6] Beginning investigation..."
./scripts/investigate.sh $INCIDENT_ID

# Step 5: Contain threat
echo "[5/6] Containing threat..."
./scripts/contain-threat.sh

# Step 6: Document and report
echo "[6/6] Documenting incident..."
./scripts/generate-report.sh $INCIDENT_ID

echo "Incident response initiated. See $INCIDENT_ID-report.md"
```

---

## 12. REFERENCES

### 12.1 Standards

- OWASP Top 10: https://owasp.org/www-project-top-ten/
- NIST SP 800-53: https://csrc.nist.gov/publications/detail/sp/800-53/rev-5/final
- ISO 27001: https://www.iso.org/isoiec-27001-information-security.html
- GDPR: https://gdpr.eu/

### 12.2 CVEs and Vulnerabilities

- CVE Database: https://cve.mitre.org/
- NVD: https://nvd.nist.gov/

### 12.3 Tools

- npm audit: https://docs.npmjs.com/cli/v8/commands/npm-audit
- OWASP ZAP: https://www.zaproxy.org/
- Burp Suite: https://portswigger.net/burp
- SonarQube: https://www.sonarsource.com/products/sonarqube/

---

**Document Control**

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-02-20 | Security Team | Initial release |

---

*This document is confidential and intended for internal use only. Distribution outside the organization requires written approval from the Security Team.*
