# SUBSCRIPTION-MGMT.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The Subscription Management module handles all aspects of recurring billing, plan management, and subscription lifecycle operations. This document covers pricing models, billing cycles, and administrative functions.

## Subscription Plans

### Plan Structure

| Plan | Price | Features | Users |
|------|-------|----------|-------|
| Starter | €29/mo | Basic features | 1-3 |
| Professional | €79/mo | Advanced features | 4-10 |
| Enterprise | €199/mo | Full features | 11-50 |
| Custom | Custom | Custom features | Unlimited |

### Plan Features

```typescript
interface Plan {
  id: string;
  name: string;
  price: number;
  currency: string;
  billingPeriod: 'monthly' | 'yearly';
  features: {
    biometricScans: number;
    apiCalls: number;
    storage: number;
    teamMembers: number;
    analytics: boolean;
    integrations: boolean;
    priority: 'low' | 'medium' | 'high';
  };
  limits: {
    maxProjects: number;
    maxDevices: number;
    maxApiKeys: number;
  };
}
```

## Database Schema

```sql
-- Plans table
CREATE TABLE subscription_plans (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price_monthly DECIMAL(10,2),
    price_yearly DECIMAL(10,2),
    currency VARCHAR(3) DEFAULT 'EUR',
    features JSONB NOT NULL,
    limits JSONB NOT NULL,
    is_active BOOLEAN DEFAULT true,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Subscription table
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id),
    tenant_id UUID REFERENCES tenants(id),
    plan_id VARCHAR(50) REFERENCES subscription_plans(id),
    status VARCHAR(50) DEFAULT 'trialing',
    billing_cycle VARCHAR(20) DEFAULT 'monthly',
    current_period_start TIMESTAMPTZ NOT NULL,
    current_period_end TIMESTAMPTZ NOT NULL,
    cancel_at_period_end BOOLEAN DEFAULT false,
    canceled_at TIMESTAMPTZ,
    trial_start TIMESTAMPTZ,
    trial_end TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Plan changes
CREATE TABLE subscription_changes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id UUID REFERENCES subscriptions(id) ON DELETE CASCADE,
    from_plan_id VARCHAR(50) REFERENCES subscription_plans(id),
    to_plan_id VARCHAR(50) REFERENCES subscription_plans(id),
    change_type VARCHAR(50) NOT NULL,
    effective_date TIMESTAMPTZ NOT NULL,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Usage tracking
CREATE TABLE subscription_usage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id UUID REFERENCES subscriptions(id) ON DELETE CASCADE,
    metric VARCHAR(100) NOT NULL,
    count INT DEFAULT 0,
    limit_value INT,
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(subscription_id, metric, period_start)
);
```

## API Endpoints

### Plans

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/subscriptions/plans | List available plans |
| GET | /api/subscriptions/plans/:id | Get plan details |

### Subscriptions

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/subscriptions | List user subscriptions |
| GET | /api/subscriptions/current | Get current subscription |
| POST | /api/subscriptions | Create subscription |
| PUT | /api/subscriptions/:id | Update subscription |
| DELETE | /api/subscriptions/:id | Cancel subscription |
| POST | /api/subscriptions/:id/change-plan | Change plan |
| POST | /api/subscriptions/:id/reactivate | Reactivate |

### Usage

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/subscriptions/usage | Get current usage |
| GET | /api/subscriptions/usage/history | Usage history |
| GET | /api/subscriptions/usage/forecast | Usage forecast |

## Billing Cycles

### Monthly Billing
- Billed on same day each month
- Prorated charges for mid-cycle changes
- Auto-renewal enabled

### Yearly Billing
- 20% discount vs monthly
- Billed annually
- Prorated upgrades

### Trial Period
- 14-day free trial
- No payment required to start
- Automatic conversion unless cancelled

## Lifecycle Events

| Event | Trigger | Action |
|-------|---------|--------|
| Subscription Created | Trial starts | Create subscription record |
| Payment Succeeded | Payment confirmed | Extend period |
| Payment Failed | Payment declined | Retry, then suspend |
| Trial Ending | 3 days before | Send reminder email |
| Subscription Cancel | User cancels | Mark cancel_at_period_end |
| Plan Change Requested | User requests | Schedule change |
| Subscription Expired | Period ended | Grace period, then suspend |

## Usage Limits

### Soft Limits
- Warning at 80% usage
- Email notification
- In-app banner

### Hard Limits
- Block at 100%
- Upgrade prompt
- 7-day grace period

### Overage Pricing
- Additional scans: €0.10 each
- Additional API calls: €0.001 each
- Additional storage: €0.10/GB/month

## Webhooks

| Event | Description |
|-------|-------------|
| subscription.created | New subscription |
| subscription.updated | Subscription changed |
| subscription.canceled | Subscription cancelled |
| subscription.renewed | Auto-renewal |
| subscription.paused | Payment failed |
| subscription.resumed | Payment recovered |
| usage.limit_reached | At 100% limit |

## Analytics

### Metrics Tracked
- MRR (Monthly Recurring Revenue)
- ARR (Annual Recurring Revenue)
- Churn rate
- LTV (Lifetime Value)
- Trial conversion rate
- ARPU (Average Revenue Per User)

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
