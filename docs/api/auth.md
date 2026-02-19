# BIOMETRICS API Authentication

Complete guide to authenticating with the BIOMETRICS API.

## Overview

BIOMETRICS uses Bearer token authentication with NVIDIA API keys. All API endpoints require authentication except for the health check endpoint.

## Getting Your API Key

### 1. NVIDIA NIM Account

1. Visit [NVIDIA NIM](https://build.nvidia.com/)
2. Sign up for a free account
3. Navigate to **API Keys** in your dashboard
4. Click **Generate New Key**
5. Copy your key (format: `nvapi-xxxxxxxxxxxx`)

### 2. Store Securely

**Environment Variable (Recommended):**

```bash
# Add to ~/.zshrc or ~/.bashrc
export NVIDIA_API_KEY="nvapi-YOUR_KEY"

# Reload shell
source ~/.zshrc
```

**Never commit API keys to git!**

Add to `.gitignore`:
```
.env
*.key
*.pem
credentials.json
```

## Authentication Methods

### Method 1: Authorization Header (Recommended)

```bash
curl -H "Authorization: Bearer nvapi-YOUR_KEY" \
  https://api.biometrics.dev/v1/agents
```

### Method 2: Environment Variable

```bash
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
  https://api.biometrics.dev/v1/agents
```

### Method 3: SDK Configuration

**JavaScript/TypeScript:**
```typescript
import { BiometricsClient } from '@biometrics/sdk';

// From environment variable
const client = new BiometricsClient(process.env.NVIDIA_API_KEY);

// Or direct key (NOT recommended for production)
const client = new BiometricsClient('nvapi-YOUR_KEY');
```

**Python:**
```python
from biometrics import BiometricsClient
import os

# From environment variable
client = BiometricsClient(os.environ['NVIDIA_API_KEY'])

# Or direct key (NOT recommended)
client = BiometricsClient('nvapi-YOUR_KEY')
```

**Go:**
```go
import (
    "github.com/biometrics/go-sdk"
    "os"
)

// From environment variable
client := biometrics.NewClient(os.Getenv("NVIDIA_API_KEY"))

// Or direct key (NOT recommended)
client := biometrics.NewClient("nvapi-YOUR_KEY")
```

## Testing Your Key

### Health Check (No Auth Required)

```bash
curl https://api.biometrics.dev/v1/health
```

**Expected Response:**
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime": 86400,
  "timestamp": "2026-02-19T10:30:00Z"
}
```

### List Models (Auth Required)

```bash
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
  https://api.biometrics.dev/v1/models
```

**Expected Response:**
```json
{
  "models": [
    {
      "id": "qwen/qwen3.5-397b-a17b",
      "name": "Qwen 3.5 397B",
      "provider": "nvidia-nim",
      "contextLimit": 262144,
      "outputLimit": 32768,
      "status": "available"
    }
  ],
  "total": 6
}
```

## Common Authentication Errors

### 401 Unauthorized

**Error:**
```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid or missing API key"
  }
}
```

**Causes:**
- Missing `Authorization` header
- Invalid API key format
- Expired or revoked key

**Solution:**
```bash
# Check if key is set
echo $NVIDIA_API_KEY

# Verify format (should start with nvapi-)
echo $NVIDIA_API_KEY | grep "^nvapi-"

# Test with explicit key
curl -H "Authorization: Bearer nvapi-YOUR_KEY" \
  https://api.biometrics.dev/v1/models
```

### 403 Forbidden

**Error:**
```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "API key does not have permission for this resource"
  }
}
```

**Causes:**
- Key lacks required permissions
- Accessing enterprise-only endpoint with free tier key

**Solution:**
- Verify key permissions in NVIDIA dashboard
- Upgrade to enterprise tier if needed

### 429 Rate Limit Exceeded

**Error:**
```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests. Please wait 60 seconds.",
    "details": {
      "retryAfter": 60,
      "limit": 40,
      "window": "1m"
    }
  }
}
```

**Causes:**
- Exceeded 40 RPM (Free Tier)
- Too many parallel requests

**Solution:**
```typescript
// Implement exponential backoff
async function makeRequestWithRetry(maxRetries = 3) {
  for (let i = 0; i < maxRetries; i++) {
    try {
      return await client.agents.list();
    } catch (error) {
      if (error.code === 'RATE_LIMIT_EXCEEDED') {
        const waitTime = Math.pow(2, i) * 1000; // 1s, 2s, 4s
        await sleep(waitTime);
        continue;
      }
      throw error;
    }
  }
}
```

## Security Best Practices

### 1. Use Environment Variables

**✅ GOOD:**
```bash
export NVIDIA_API_KEY="nvapi-YOUR_KEY"
```

**❌ BAD:**
```typescript
// Never hardcode in source code!
const client = new BiometricsClient('nvapi-YOUR_KEY');
```

### 2. Use .env Files Locally

Create `.env` file:
```
NVIDIA_API_KEY=nvapi-YOUR_KEY
```

Load in code:
```typescript
import 'dotenv/config';
const client = new BiometricsClient(process.env.NVIDIA_API_KEY);
```

Add to `.gitignore`:
```
.env
.env.local
.env.*.local
```

### 3. Rotate Keys Regularly

**Recommended:** Every 90 days

**Steps:**
1. Generate new key in NVIDIA dashboard
2. Update environment variable
3. Test with new key
4. Revoke old key

### 4. Use Separate Keys Per Environment

| Environment | Key | Purpose |
|-------------|-----|---------|
| Development | `nvapi-dev-xxx` | Local testing |
| Staging | `nvapi-staging-xxx` | Pre-production |
| Production | `nvapi-prod-xxx` | Live environment |

### 5. Monitor Key Usage

Check usage in NVIDIA dashboard:
- Request count
- Error rate
- Rate limit hits
- Unusual activity

## OAuth 2.0 (Coming Soon)

BIOMETRICS will support OAuth 2.0 for enterprise customers in Q2 2026.

**Features:**
- Refresh tokens
- Scoped permissions
- Team access management
- Audit logging

## Multi-Factor Authentication (Enterprise)

Enterprise accounts can enable MFA for API key generation:

1. Go to **Settings → Security**
2. Enable **MFA for API Keys**
3. Configure TOTP authenticator
4. Generate keys with MFA verification

## Troubleshooting

### Key Not Working

**Check:**
1. Key format (must start with `nvapi-`)
2. No extra spaces or quotes
3. Key not expired or revoked
4. Correct environment variable name

**Debug:**
```bash
# Print key (first 10 chars only)
echo ${NVIDIA_API_KEY:0:10}...

# Test with curl
curl -v -H "Authorization: Bearer $NVIDIA_API_KEY" \
  https://api.biometrics.dev/v1/models
```

### SSL/TLS Errors

**Error:**
```
curl: (60) SSL certificate problem
```

**Solution:**
```bash
# Update CA certificates
# macOS
brew update && brew upgrade ca-certificates

# Ubuntu/Debian
sudo apt-get update && sudo apt-get install --reinstall ca-certificates

# Verify SSL
curl -v https://api.biometrics.dev/v1/health
```

### Timeout Errors

**Error:**
```
curl: (28) Operation timed out
```

**Solution:**
```bash
# Increase timeout
curl --max-time 120 \
  -H "Authorization: Bearer $NVIDIA_API_KEY" \
  https://api.biometrics.dev/v1/models

# Check network
ping api.biometrics.dev
```

## SDK Authentication Examples

### JavaScript/TypeScript

```typescript
import { BiometricsClient } from '@biometrics/sdk';
import 'dotenv/config';

// Initialize client
const client = new BiometricsClient(process.env.NVIDIA_API_KEY);

// Test authentication
try {
  const models = await client.models.list();
  console.log(`Authenticated! Found ${models.total} models`);
} catch (error) {
  console.error('Authentication failed:', error.message);
}
```

### Python

```python
from biometrics import BiometricsClient
import os
from dotenv import load_dotenv

load_dotenv()

# Initialize client
client = BiometricsClient(os.environ['NVIDIA_API_KEY'])

# Test authentication
try:
    models = client.models.list()
    print(f"Authenticated! Found {models.total} models")
except Exception as e:
    print(f"Authentication failed: {e}")
```

### Go

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "github.com/biometrics/go-sdk"
)

func main() {
    // Initialize client
    client := biometrics.NewClient(os.Getenv("NVIDIA_API_KEY"))
    
    // Test authentication
    models, err := client.Models.List(context.Background())
    if err != nil {
        log.Fatalf("Authentication failed: %v", err)
    }
    
    fmt.Printf("Authenticated! Found %d models\n", models.Total)
}
```

## Compliance

### GDPR

BIOMETRICS is GDPR compliant:
- No personal data stored in API keys
- Keys can be revoked at any time
- Data processing agreements available for enterprise

### SOC 2

BIOMETRICS maintains SOC 2 Type II certification:
- Regular security audits
- Access logging and monitoring
- Incident response procedures

## Support

**Authentication Issues:**
- Email: security@biometrics.dev
- Discord: #authentication channel
- GitHub Issues: [Authentication Label](https://github.com/Delqhi/BIOMETRICS/issues?q=label:authentication)

**Emergency (Production Down):**
- Email: emergency@biometrics.dev
- Response time: < 1 hour (Enterprise only)

---

**Last Updated:** February 2026 | **Version:** 1.0.0
