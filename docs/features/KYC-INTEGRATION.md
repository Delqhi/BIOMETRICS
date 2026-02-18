# KYC-INTEGRATION.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The KYC (Know Your Customer) Integration module provides identity verification services for regulatory compliance. This document covers verification workflows, document types, and compliance requirements.

## Verification Levels

| Level | Requirements | Use Case |
|-------|--------------|----------|
| Basic | Email, Phone | Low-risk |
| Standard | ID Document | Standard users |
| Enhanced | ID + Selfie + Address | High-risk |
| Premium | Full verification | Enterprise |

## Supported Documents

### Identity Documents

| Country | Document Types |
|---------|----------------|
| EU | Passport, ID Card, Driver's License |
| US | Passport, State ID, Driver's License |
| UK | Passport, Biometric ID |
| Global | Passport (all countries) |

### Proof of Address

| Document | Validity |
|----------|----------|
| Utility Bill | 3 months |
| Bank Statement | 3 months |
| Government Letter | 6 months |
| Rental Agreement | Current |

## Database Schema

```sql
-- KYC applications
CREATE TABLE kyc_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id),
    tenant_id UUID REFERENCES tenants(id),
    level VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    provider VARCHAR(100),
    provider_application_id VARCHAR(255),
    expires_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    rejected_at TIMESTAMPTZ,
    rejection_reason TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- KYC documents
CREATE TABLE kyc_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID REFERENCES kyc_applications(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    country VARCHAR(2),
    document_number VARCHAR(100),
    expiry_date DATE,
    front_image_url TEXT,
    back_image_url TEXT,
    selfie_image_url TEXT,
    verification_status VARCHAR(50) DEFAULT 'pending',
    rejection_reason TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- KYC checks
CREATE TABLE kyc_checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID REFERENCES kyc_applications(id) ON DELETE CASCADE,
    check_type VARCHAR(50) NOT NULL,
    result JSONB NOT NULL,
    passed BOOLEAN NOT NULL,
    details JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Watchlist screenings
CREATE TABLE kyc_watchlists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID REFERENCES kyc_applications(id) ON DELETE CASCADE,
    watchlist_type VARCHAR(50) NOT NULL,
    search_name VARCHAR(255),
    match_found BOOLEAN DEFAULT false,
    risk_score INT,
    matches JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## API Endpoints

### Applications

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/kyc/applications | Start KYC process |
| GET | /api/kyc/applications | List applications |
| GET | /api/kyc/applications/:id | Get application status |
| POST | /api/kyc/applications/:id/upload | Upload document |

### Verification

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/kyc/:id/status | Check verification status |
| POST | /api/kyc/:id/verify-document | Verify document |
| POST | /api/kyc/:id/verify-selfie | Verify selfie |
| POST | /api/kyc/:id/verify-liveness | Liveness check |

### Webhooks

| Event | Description |
|-------|-------------|
| kyc.approved | Verification approved |
| kyc.rejected | Verification rejected |
| kyc.review | Manual review needed |
| kyc.expired | Application expired |

## Verification Flow

```
Start KYC
    |
    v
Select Level --> Standard
    |
    v
Upload ID Document
    |
    v
Document Verification --> AI Check
    |                     |
    v                     v
Selfie with ID      Human Review (if needed)
    |
    v
Address Verification (if enhanced)
    |
    v
Watchlist Screening
    |
    v
Decision --> Approved/Rejected
```

## Provider Integration

### Onfido

```typescript
const onfidoConfig = {
  apiToken: process.env.ONFIDO_API_TOKEN,
  webhookSigningKey: process.env.ONFIDO_WEBHOOK_KEY,
  checks: {
    document: true,
    facialSimilarity: true,
    identity: true,
    watchlist: true
  }
};
```

### Stripe Identity

```typescript
const stripeIdentityConfig = {
  apiKey: process.env.STRIPE_SECRET_KEY,
  options: {
    require_matching_selfie: true,
    require_live_capture: true,
    document_types: ['passport', 'driving_license', 'national_id']
  }
};
```

## Compliance

### Data Retention

| Data Type | Retention |
|-----------|-----------|
| Documents | 7 years |
| Selfie images | 5 years |
| Verification results | 7 years |
| Watchlist checks | 10 years |

### Privacy

- GDPR compliant
- Data encryption at rest
- Right to erasure support
- Consent management

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
