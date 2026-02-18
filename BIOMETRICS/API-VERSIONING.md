# API-VERSIONING.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The API Versioning module manages API versions, deprecations, and backwards compatibility. This document covers versioning strategies, migration guides, and lifecycle policies.

## Versioning Strategy

### URL-Based Versioning

```
/v1/users
/v2/users
/v3/users
```

### Header-Based Versioning

```
Accept: application/vnd.biometrics.v2+json
```

## Version Lifecycle

| Stage | Duration | Support |
|-------|----------|---------|
| Beta | 3 months | Community |
| Stable | 12 months | Full |
| Legacy | 6 months | Security only |
| Deprecated | 6 months | Read-only |

## Database Schema

```sql
-- API versions
CREATE TABLE api_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version VARCHAR(20) NOT NULL,
    status VARCHAR(50) NOT NULL,
    release_date DATE NOT NULL,
    sunset_date DATE,
    docs_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Version deprecations
CREATE TABLE api_deprecations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version VARCHAR(20) NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    deprecation_date DATE NOT NULL,
    sunset_date DATE NOT NULL,
    migration_guide TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- API changelog
CREATE TABLE api_changelog (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version VARCHAR(20) NOT NULL,
    release_date DATE NOT NULL,
    changes JSONB NOT NULL,
    breaking_changes JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## Deprecation Timeline

```
v1 Released: 2025-01-01
v1 Stable: 2025-04-01
v1 Legacy: 2026-04-01
v1 Deprecated: 2026-10-01
v1 Sunset: 2027-04-01

v2 Released: 2026-01-01
v2 Stable: 2026-04-01
```

## Migration Guide

### v1 to v2 Changes

| v1 Endpoint | v2 Endpoint |
|-------------|-------------|
| GET /users | GET /v2/users |
| POST /scans | POST /v2/biometrics/scans |
| GET /templates | GET /v2/templates |

### Response Format Changes

```json
// v1
{
  "users": [...],
  "count": 100
}

// v2
{
  "data": [...],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 100
  }
}
```

## Version Detection

### Automatic Versioning

```typescript
// Express middleware
app.use((req, res, next) => {
  const version = req.headers['accept-version'] || 'v1';
  req.apiVersion = version;
  next();
});
```

### Version Response Headers

```
X-API-Version: v2
X-API-Deprecated: true
X-API-Sunset: 2026-10-01
X-API-Migration-Guide: https://docs.biometrics.com/migration/v1-v2
```

## Deprecation Notifications

### Headers

```
Deprecation: true
Sunset: Sat, 01 Oct 2026 00:00:00 GMT
Link: <https://api.biometrics.com/v2/users>; rel="successor-version"
```

### Email Notifications

| Timeline | Notification |
|----------|--------------|
| 6 months before | Deprecation announced |
| 3 months before | Migration required |
| 1 month before | Final warning |

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
