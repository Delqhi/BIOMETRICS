# DATA-EXPORT.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The Data Export module enables users to export their data in various formats. This document covers export types, formats, scheduling, and delivery methods.

## Export Types

| Type | Description | Formats |
|------|-------------|---------|
| Full | All user data | JSON, CSV |
| Partial | Selected data types | JSON, CSV, XLSX |
| Scheduled | Recurring exports | JSON, CSV, XLSX, XML |
| Real-time | Streaming exports | JSON |

## Database Schema

```sql
-- Export requests
CREATE TABLE data_exports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id),
    tenant_id UUID REFERENCES tenants(id),
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    format VARCHAR(20) NOT NULL,
    filters JSONB,
    file_url TEXT,
    file_size BIGINT,
    record_count INT,
    download_count INT DEFAULT 0,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

-- Export data types
CREATE TABLE export_data_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    tables JSONB NOT NULL,
    columns JSONB,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Scheduled exports
CREATE TABLE scheduled_exports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id),
    tenant_id UUID REFERENCES tenants(id),
    data_types UUID[] NOT NULL,
    format VARCHAR(20) NOT NULL,
    frequency VARCHAR(20) NOT NULL,
    next_run_at TIMESTAMPTZ NOT NULL,
    last_run_at TIMESTAMPTZ,
    delivery_method VARCHAR(50) DEFAULT 'download',
    delivery_email BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

## API Endpoints

### Exports

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/exports | Create export |
| GET | /api/exports | List exports |
| GET | /api/exports/:id | Get export status |
| GET | /api/exports/:id/download | Download export |
| DELETE | /api/exports/:id | Cancel export |

### Scheduled

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/exports/scheduled | List scheduled |
| POST | /api/exports/scheduled | Create schedule |
| PUT | /api/exports/scheduled/:id | Update schedule |
| DELETE | /api/exports/scheduled/:id | Delete schedule |

## Data Types

### User Data

```json
{
  "profile": ["users", "user_profiles"],
  "biometric": ["biometric_scans", "biometric_templates"],
  "documents": ["documents", "document_versions"],
  "payments": ["payment_methods", "invoices"],
  "subscriptions": ["subscriptions", "subscription_usage"],
  "activity": ["sessions", "audit_logs"]
}
```

## Export Formats

### JSON

```json
{
  "exportDate": "2026-02-18T10:00:00Z",
  "userId": "user-123",
  "data": {
    "profile": [...],
    "biometric": [...]
  }
}
```

### CSV

```csv
id,name,email,created_at
1,John Doe,john@example.com,2026-01-01
```

## Processing

### Async Processing
- Large exports queued
- Progress tracking
- Email notification on complete

### File Storage
- Temporary storage: 7 days
- Signed URLs for download
- Encryption at rest

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
