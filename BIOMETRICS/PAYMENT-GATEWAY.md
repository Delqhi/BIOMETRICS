# PAYMENT-GATEWAY.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The Payment Gateway module provides unified payment processing across multiple payment providers. This document describes the architecture, supported providers, and integration specifications.

## Supported Providers

| Provider | Type | Currencies | Status |
|----------|------|-----------|--------|
| Stripe | Card, Apple Pay, Google Pay | 135+ | Active |
| PayPal | Wallet, Card | 25+ | Active |
| Mollie | iDEAL, Bancontact, Klarna | EUR | Active |
| Adyen | Cards, Wallets | 150+ | Active |

## Architecture

### Components

| Component | Description | Technology |
|-----------|-------------|------------|
| Payment Router | Provider selection | Node.js |
| Token Vault | Secure token storage | Supabase + Encryption |
| Webhook Handler | Async payment events | Supabase Functions |
| Refund Engine | Partial/full refunds | Provider APIs |
| Dispute Manager | Chargeback handling | Provider APIs |

### Payment Flow

```
User               BIOMETRICS              Payment Provider           Bank
  |                    |                        |                       |
  |--- Select Item -->|                        |                       |
  |                    |                        |                       |
  |--- Checkout ----->|                        |                       |
  |                    |--- Create Payment --->|                       |
  |                    |<-- Redirect ---------|                       |
  |<-- Payment Page --|                        |                       |
  |                    |                        |                       |
  |--- Enter Details ->|                        |                       |
  |                    |--- Process Payment --->|                       |
  |                    |<-- Success/Fail ------|                       |
  |<-- Confirmation --|                        |                       |
  |                    |                        |                       |
```

## Database Schema

```sql
-- Payment methods
CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id),
    provider VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL,
    last_four VARCHAR(4),
    brand VARCHAR(50),
    expiry_month INT,
    expiry_year INT,
    is_default BOOLEAN DEFAULT false,
    provider_customer_id VARCHAR(255),
    provider_payment_method_id VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Transactions
CREATE TABLE payment_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id),
    tenant_id UUID REFERENCES tenants(id),
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    provider_transaction_id VARCHAR(255),
    provider_charge_id VARCHAR(255),
    status VARCHAR(50) DEFAULT 'pending',
    type VARCHAR(50) NOT NULL,
    description TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    error_message TEXT
);

-- Subscriptions
CREATE TABLE payment_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id),
    tenant_id UUID REFERENCES tenants(id),
    plan_id VARCHAR(255) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    provider_subscription_id VARCHAR(255),
    status VARCHAR(50) DEFAULT 'active',
    current_period_start TIMESTAMPTZ,
    current_period_end TIMESTAMPTZ,
    cancel_at_period_end BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    canceled_at TIMESTAMPTZ
);

-- Invoices
CREATE TABLE payment_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id),
    tenant_id UUID REFERENCES tenants(id),
    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    subscription_id UUID REFERENCES payment_subscriptions(id),
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status VARCHAR(50) DEFAULT 'draft',
    paid_at TIMESTAMPTZ,
    due_date TIMESTAMPTZ,
    items JSONB NOT NULL,
    pdf_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## API Endpoints

### Payments

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/payments/create | Create payment intent |
| POST | /api/payments/confirm | Confirm payment |
| POST | /api/payments/cancel | Cancel payment |
| GET | /api/payments/:id | Get payment details |
| POST | /api/payments/:id/refund | Process refund |

### Payment Methods

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/payments/methods | List payment methods |
| POST | /api/payments/methods | Add payment method |
| PUT | /api/payments/methods/:id | Update payment method |
| DELETE | /api/payments/methods/:id | Remove payment method |
| POST | /api/payments/methods/:id/set-default | Set default |

### Subscriptions

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/payments/subscriptions | List subscriptions |
| POST | /api/payments/subscriptions | Create subscription |
| PUT | /api/payments/subscriptions/:id | Update subscription |
| DELETE | /api/payments/subscriptions/:id | Cancel subscription |
| POST | /api/payments/subscriptions/:id/reactivate | Reactivate |

### Invoices

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/payments/invoices | List invoices |
| GET | /api/payments/invoices/:id | Get invoice |
| GET | /api/payments/invoices/:id/pdf | Download PDF |

## Webhook Events

| Event | Description | Processing |
|-------|-------------|------------|
| payment_intent.succeeded | Payment successful | Update transaction status |
| payment_intent.failed | Payment failed | Update status, notify user |
| customer.subscription.created | Subscription started | Create subscription record |
| customer.subscription.updated | Subscription changed | Update subscription |
| customer.subscription.deleted | Subscription cancelled | Mark as cancelled |
| invoice.paid | Invoice paid | Mark invoice as paid |
| charge.refunded | Charge refunded | Update transaction |

## Security

- PCI DSS compliance via provider tokens
- No card data touches our servers
- Encrypted storage for sensitive data
- Webhook signature verification
- Rate limiting on payment endpoints
- 3D Secure support

## Error Handling

| Error Code | Description | Action |
|------------|-------------|--------|
| card_declined | Card was declined | Show error, suggest retry |
| expired_card | Card expired | Request new card |
| insufficient_funds | Not enough funds | Notify user |
| processing_error | Provider error | Retry with backoff |
| lost_card | Card reported lost | Block, notify user |

## Fee Structure

| Provider | Transaction Fee | Setup Fee | Monthly Fee |
|----------|-----------------|-----------|-------------|
| Stripe | 1.4% + €0.25 | €0 | €0 |
| PayPal | 1.9% + €0.35 | €0 | €0 |
| Mollie | €0.29 | €0 | €0 |
| Adyen | 1.25% | €0 | €0 |

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
