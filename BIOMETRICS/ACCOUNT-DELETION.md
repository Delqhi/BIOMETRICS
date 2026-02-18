# ACCOUNT-DELETION.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The Account Deletion module handles user account termination with GDPR compliance and data retention rules. This document covers deletion workflows, grace periods, and data handling.

## Deletion Types

| Type | Description | Timeline |
|------|-------------|----------|
| Immediate | Hard delete | Instant |
| Scheduled | Queue for deletion | 30 days |
| Anonymized | Remove PII, keep data | Instant |

## Database Schema

```sql
-- Deletion requests
CREATE TABLE account_deletions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id),
    tenant_id UUID REFERENCES tenants(id),
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    scheduled_deletion_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    reason TEXT,
    feedback TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Data retention rules
CREATE TABLE data_retention_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    data_type VARCHAR(100) NOT NULL,
    retention_period_days INT NOT NULL,
    delete_after_days INT NOT NULL,
    legal_basis TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Deletion tasks
CREATE TABLE deletion_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deletion_id UUID REFERENCES account_deletions(id) ON DELETE CASCADE,
    table_name VARCHAR(100) NOT NULL,
    record_count INT,
    status VARCHAR(50) DEFAULT 'pending',
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/account/request-deletion | Request deletion |
| GET | /api/account/deletion-status | Get status |
| POST | /api/account/cancel-deletion | Cancel request |
| DELETE | /api/account/data/:type | Delete specific data |

## Deletion Workflow

```
User Request
    |
    v
Validate Request --> Send confirmation email
    |
    v
Grace Period (30 days)
    |
    v
Notification --> User can cancel
    |
    v
Execute Deletion
    |
    v
Verify --> Send confirmation
```

## Data Handling

### Immediate Deletion

| Data Type | Action |
|-----------|--------|
| User profile | Delete |
| Authentication | Delete |
| Personal data | Delete |
| Payment history | Anonymize |

### Retention Exceptions

| Data Type | Retention | Reason |
|-----------|-----------|--------|
| Financial records | 10 years | Legal requirement |
| Audit logs | 7 years | Compliance |
| Tax documents | 7 years | Legal |

### Anonymization

```typescript
const anonymizeUser = (user) => {
  return {
    id: user.id,
    name: 'DELETED_USER',
    email: `deleted_${user.id}@biometrics.invalid`,
    phone: null,
    created_at: user.created_at,
    deleted_at: new Date(),
    is_anonymized: true
  };
};
```

## User Notifications

### Email Sequence

| Day | Action |
|-----|--------|
| 0 | Confirmation request |
| 1 | Deletion scheduled |
| 7 | One week warning |
| 28 | Final warning |
| 30 | Deletion complete |

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
