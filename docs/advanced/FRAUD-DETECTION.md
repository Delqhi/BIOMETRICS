# FRAUD-DETECTION.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The Fraud Detection module identifies and prevents fraudulent activities using machine learning, rule-based checks, and behavioral analysis. This document covers detection methods, risk scoring, and response actions.

## Detection Methods

| Method | Description | Effectiveness |
|--------|-------------|---------------|
| Velocity Checks | Rapid repeated actions | High |
| Device Fingerprinting | Identify suspicious devices | Medium |
| Geolocation | Detect impossible travel | High |
| Behavioral Biometrics | User behavior patterns | Very High |
| ML Scoring | AI-powered risk analysis | Very High |
| Blocklists | Known fraud patterns | Medium |

## Risk Factors

### Account Risk
- New account with high-value orders
- Account age vs order value mismatch
- Multiple accounts from same device
- Suspicious account activity

### Transaction Risk
- Unusual purchase amount
- Multiple failed payments
- High-risk merchant category
- Card not present (CNP) transactions

### Behavioral Risk
- Unusual navigation patterns
- Keyboard/mouse dynamics
- Session replay attempts
- Bot-like behavior

## Database Schema

```sql
-- Fraud rules
CREATE TABLE fraud_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    conditions JSONB NOT NULL,
    action VARCHAR(50) NOT NULL,
    action_params JSONB,
    risk_weight INT DEFAULT 10,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Fraud assessments
CREATE TABLE fraud_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    risk_score INT NOT NULL,
    risk_level VARCHAR(20) NOT NULL,
    factors JSONB NOT NULL,
    recommendations JSONB,
    decision VARCHAR(50) NOT NULL,
    reviewed_by UUID REFERENCES auth.users(id),
    reviewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Blocked entities
CREATE TABLE fraud_blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    reason TEXT NOT NULL,
    block_type VARCHAR(50) NOT NULL,
    expires_at TIMESTAMPTZ,
    blocked_by UUID REFERENCES auth.users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Device fingerprints
CREATE TABLE device_fingerprints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fingerprint VARCHAR(255) UNIQUE NOT NULL,
    user_id UUID REFERENCES auth.users(id),
    device_info JSONB,
    risk_score INT DEFAULT 0,
    is_blocked BOOLEAN DEFAULT false,
    first_seen TIMESTAMPTZ DEFAULT NOW(),
    last_seen TIMESTAMPTZ DEFAULT NOW()
);
```

## API Endpoints

### Assessment

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/fraud/check | Run fraud check |
| GET | /api/fraud/assessments | List assessments |
| GET | /api/fraud/assessments/:id | Get assessment |

### Rules

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/fraud/rules | List rules |
| POST | /api/fraud/rules | Create rule |
| PUT | /api/fraud/rules/:id | Update rule |
| DELETE | /api/fraud/rules/:id | Delete rule |

### Blocks

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/fraud/blocks | List blocks |
| POST | /api/fraud/blocks | Create block |
| DELETE | /api/fraud/blocks/:id | Remove block |

## Risk Scoring

### Score Range

| Score | Risk Level | Action |
|-------|------------|--------|
| 0-30 | Low | Approve |
| 31-60 | Medium | Review |
| 61-80 | High | Manual Review |
| 81-100 | Critical | Block |

### Score Factors

```typescript
interface RiskFactor {
  category: string;
  name: string;
  weight: number;
  value: number;
  description: string;
}

// Example factors
const factors: RiskFactor[] = [
  { category: 'velocity', name: 'failed_payments', weight: 20, value: 3, description: '3 failed payments' },
  { category: 'geolocation', name: 'impossible_travel', weight: 30, value: 1, description: 'Impossible travel detected' },
  { category: 'device', name: 'new_device', weight: 15, value: 1, description: 'New device for user' },
  { category: 'behavior', name: 'bot_pattern', weight: 25, value: 1, description: 'Bot-like behavior detected' }
];
```

## Response Actions

| Action | Description | Trigger |
|--------|-------------|---------|
| ALLOW | Auto approve | Score < 30 |
| REVIEW | Manual review | Score 30-80 |
| BLOCK | Immediate block | Score > 80 |
| CAPTCHA | Require CAPTCHA | Suspicious |
| MFA | Require 2FA | High risk |
| HOLD | Delay processing | Review pending |

## Webhooks

| Event | Description |
|-------|-------------|
| fraud.high_risk | High risk detected |
| fraud.blocked | Entity blocked |
| fraud.review_required | Manual review needed |
| fraud.cleared | Risk cleared |

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
