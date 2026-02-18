# MULTI-TENANT.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The Multi-Tenant module provides isolated workspaces for multiple organizations on a single platform. This document covers tenant isolation, resource allocation, and administrative controls.

## Tenant Models

| Model | Description | Use Case |
|-------|-------------|----------|
| Isolated | Full separation | Enterprise |
| Shared | Shared resources | SaaS |
| Hybrid | Mix of both | Custom |

## Database Schema

```sql
-- Tenants
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    type VARCHAR(50) DEFAULT 'standard',
    status VARCHAR(50) DEFAULT 'active',
    plan_id VARCHAR(50),
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Tenant users
CREATE TABLE tenant_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    invited_by UUID REFERENCES auth.users(id),
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, user_id)
);

-- Tenant settings
CREATE TABLE tenant_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    key VARCHAR(100) NOT NULL,
    value JSONB,
    is_public BOOLEAN DEFAULT false,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, key)
);

-- Tenant quotas
CREATE TABLE tenant_quotas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    resource VARCHAR(100) NOT NULL,
    limit_value INT NOT NULL,
    used_value INT DEFAULT 0,
    period VARCHAR(20) DEFAULT 'monthly',
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, resource)
);
```

## API Endpoints

### Tenants

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/tenants | List tenants |
| POST | /api/tenants | Create tenant |
| GET | /api/tenants/:id | Get tenant |
| PUT | /api/tenants/:id | Update tenant |
| DELETE | /api/tenants/:id | Delete tenant |

### Tenant Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/tenants/:id/users | List users |
| POST | /api/tenants/:id/users/invite | Invite user |
| PUT | /api/tenants/:id/users/:userId | Update role |
| DELETE | /api/tenants/:id/users/:userId | Remove user |

### Settings

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/tenants/:id/settings | Get settings |
| PUT | /api/tenants/:id/settings | Update settings |

## Tenant Isolation

### Row-Level Security

```sql
-- Enable RLS
ALTER TABLE biometric_scans ENABLE ROW LEVEL SECURITY;

-- Tenant policy
CREATE POLICY tenant_isolation ON biometric_scans
  USING (tenant_id = current_setting('app.tenant_id')::uuid);
```

### Data Isolation

| Isolation Level | Implementation |
|-----------------|----------------|
| Row-level | PostgreSQL RLS |
| Schema | Separate schemas |
| Database | Separate databases |
| Instance | Separate servers |

## Tenant Roles

| Role | Permissions |
|------|-------------|
| Owner | Full access |
| Admin | Manage users/settings |
| Member | Use features |
| Viewer | Read-only |

## Quota Management

### Default Quotas

| Resource | Starter | Professional | Enterprise |
|----------|---------|--------------|------------|
| Users | 3 | 10 | Unlimited |
| Storage | 1GB | 10GB | Unlimited |
| API Calls | 10K | 100K | Unlimited |
| Scans | 100 | 1000 | Unlimited |

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
