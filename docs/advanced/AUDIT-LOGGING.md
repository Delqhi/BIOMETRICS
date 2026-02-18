# AUDIT-LOGGING.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The Audit Logging module provides comprehensive logging of all system activities for compliance, security, and debugging purposes. This document covers logging categories, retention, and analysis.

## Log Categories

| Category | Description | Retention |
|----------|-------------|-----------|
| Authentication | Login, logout, 2FA | 7 years |
| Data Access | Read, export | 3 years |
| Data Modification | Create, update, delete | 7 years |
| Security | Blocks, alerts | 7 years |
| Admin | Configuration changes | 5 years |
| API | All API calls | 1 year |
| Payments | Transactions | 10 years |

## Database Schema

```sql
-- Audit logs
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    user_id UUID REFERENCES auth.users(id),
    action VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    resource_type VARCHAR(100),
    resource_id VARCHAR(255),
    changes JSONB,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Audit indexes
CREATE INDEX idx_audit_tenant ON audit_logs(tenant_id, created_at);
CREATE INDEX idx_audit_user ON audit_logs(user_id, created_at);
CREATE INDEX idx_audit_action ON audit_logs(action);
CREATE INDEX idx_audit_category ON audit_logs(category);

-- Log archives
CREATE TABLE audit_archives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    year_month VARCHAR(7) NOT NULL,
    storage_location VARCHAR(500),
    record_count INT,
    size_bytes BIGINT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/audit/logs | Query audit logs |
| GET | /api/audit/logs/:id | Get log entry |
| GET | /api/audit/export | Export logs |
| GET | /api/audit/summary | Activity summary |

## Log Entry Structure

```typescript
interface AuditLogEntry {
  id: string;
  tenantId: string;
  userId: string;
  action: string;
  category: 'auth' | 'data' | 'security' | 'admin' | 'api' | 'payment';
  resourceType: string;
  resourceId: string;
  changes?: {
    field: string;
    old: any;
    new: any;
  }[];
  ipAddress: string;
  userAgent: string;
  metadata?: Record<string, any>;
  timestamp: Date;
}
```

## Event Types

### Authentication Events

| Event | Category | Description |
|-------|----------|-------------|
| user.login | auth | User logged in |
| user.logout | auth | User logged out |
| user.login_failed | auth | Login failed |
| user.password_changed | auth | Password changed |
| user.2fa_enabled | auth | 2FA enabled |
| user.2fa_disabled | auth | 2FA disabled |

### Data Events

| Event | Category | Description |
|-------|----------|-------------|
| data.created | data | Record created |
| data.updated | data | Record updated |
| data.deleted | data | Record deleted |
| data.exported | data | Data exported |
| data.viewed | data | Record viewed |

### Security Events

| Event | Category | Description |
|-------|----------|-------------|
| security.block | security | IP/user blocked |
| security.unblock | security | IP/user unblocked |
| security.rate_limit | security | Rate limited |
| security.suspicious | security | Suspicious activity |

## Retention Policy

| Data Type | Retention | Storage |
|-----------|-----------|---------|
| Last 90 days | Hot storage | PostgreSQL |
| 91 days - 1 year | Warm storage | S3 |
| 1-7 years | Cold storage | Glacier |
| 7+ years | Archive | Tape |

## Compliance

### GDPR
- All data access logged
- User consent tracked
- Right to erasure supported

### SOX
- Financial data changes
- Admin actions
- 7-year retention

### PCI-DSS
- Payment events
- Card data access
- Authentication events

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
