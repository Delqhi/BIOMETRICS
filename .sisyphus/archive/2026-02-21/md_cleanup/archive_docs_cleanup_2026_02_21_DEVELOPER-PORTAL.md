# DEVELOPER-PORTAL.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The Developer Portal provides a self-service platform for developers to access documentation, manage API keys, and test integrations. This document covers portal features and usage.

## Portal Features

| Feature | Description |
|---------|-------------|
| Documentation | Interactive API docs |
| API Keys | Generate and manage keys |
| Sandbox | Test environment |
| Code Samples | SDK examples |
| Status Page | System status |
| Support | Developer support |

## Database Schema

```sql
-- Developer applications
CREATE TABLE developer_apps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    website VARCHAR(500),
    redirect_uris TEXT[],
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- API keys
CREATE TABLE developer_api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id UUID REFERENCES developer_apps(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    key_hash VARCHAR(255) NOT NULL,
    key_prefix VARCHAR(20) NOT NULL,
    scopes TEXT[],
    last_used_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Usage tracking
CREATE TABLE developer_usage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id UUID REFERENCES developer_apps(id) ON DELETE CASCADE,
    endpoint VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    status_code INT,
    response_time_ms INT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## API Endpoints

### Applications

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/developer/apps | List apps |
| POST | /api/developer/apps | Create app |
| GET | /api/developer/apps/:id | Get app |
| PUT | /api/developer/apps/:id | Update app |
| DELETE | /api/developer/apps/:id | Delete app |

### API Keys

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/developer/keys | List keys |
| POST | /api/developer/keys | Create key |
| DELETE | /api/developer/keys/:id | Revoke key |

### Usage

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/developer/usage | Get usage stats |
| GET | /api/developer/usage/logs | API call logs |

## Application Registration

### Create Application

```json
POST /api/developer/apps
{
  "name": "My Integration",
  "description": "Integration with our CRM",
  "website": "https://myapp.com",
  "redirect_uris": [
    "https://myapp.com/callback"
  ]
}
```

### Response

```json
{
  "client_id": "app_abc123",
  "client_secret": "secret_xyz789"
}
```

## API Key Management

### Create Key

```json
POST /api/developer/keys
{
  "name": "Production Key",
  "scopes": ["read:users", "write:scans"],
  "expires_in": 365
}
```

### Key Format

```
bi_sk_live_abc123xyz789
```

| Part | Description |
|------|-------------|
| `bi_` | Prefix |
| `sk_` | Secret key |
| `live_` | Environment |
| `abc123...` | Random identifier |

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
