# TAX-CALCULATION.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The Tax Calculation module handles automated tax calculations, compliance, and reporting across multiple jurisdictions. This document covers tax rates, nexus tracking, and reporting features.

## Supported Tax Types

| Tax Type | Regions | Examples |
|----------|---------|----------|
| VAT | EU, UK, Australia | 15-27% |
| GST | Canada, NZ, Singapore | 5-15% |
| Sales Tax | US States | 0-10% |
| Consumption | Japan, Korea | 10% |

## Database Schema

```sql
-- Tax jurisdictions
CREATE TABLE tax_jurisdictions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_code VARCHAR(2) NOT NULL,
    region_code VARCHAR(50),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    rates JSONB NOT NULL,
    is_active BOOLEAN DEFAULT true,
    effective_from DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Tax nexus (where you collect tax)
CREATE TABLE tax_nexus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    country_code VARCHAR(2) NOT NULL,
    region_code VARCHAR(50),
    nexus_type VARCHAR(50) NOT NULL,
    registration_number VARCHAR(100),
    registered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, country_code, region_code)
);

-- Tax settings per tenant
CREATE TABLE tax_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    default_country VARCHAR(2),
    default_rate DECIMAL(5,2),
    include_tax_in_prices BOOLEAN DEFAULT true,
    tax_calculation_mode VARCHAR(20) DEFAULT 'inclusive',
    collected_threshold DECIMAL(12,2),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Tax transactions
CREATE TABLE tax_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL,
    transaction_type VARCHAR(50) NOT NULL,
    country_code VARCHAR(2) NOT NULL,
    region_code VARCHAR(50),
    tax_rate DECIMAL(5,2) NOT NULL,
    subtotal DECIMAL(12,2) NOT NULL,
    tax_amount DECIMAL(12,2) NOT NULL,
    taxJurisdiction_id UUID REFERENCES tax_jurisdictions(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## API Endpoints

### Tax Rates

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/tax/rates | List tax rates |
| GET | /api/tax/rates/:country | Get rates by country |
| POST | /api/tax/calculate | Calculate tax |

### Nexus

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/tax/nexus | List nexus registrations |
| POST | /api/tax/nexus | Register nexus |
| DELETE | /api/tax/nexus/:id | Remove nexus |

### Reports

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/tax/reports/summary | Tax summary |
| GET | /api/tax/reports/by-region | Breakdown by region |
| GET | /api/tax/reports/by-period | Breakdown by period |

## Tax Calculation

### Input

```typescript
interface TaxCalculationInput {
  amount: number;
  countryCode: string;
  regionCode?: string;
  customerCountry: string;
  customerRegion?: string;
  customerVatNumber?: string;
  transactionType: 'sale' | 'refund';
}
```

### Output

```typescript
interface TaxCalculationResult {
  subtotal: number;
  taxRate: number;
  taxAmount: number;
  total: number;
  taxJurisdiction: string;
  isReverseCharge: boolean;
  exemptReason?: string;
}
```

## VAT Handling

### EU VAT
- Standard rates: 17-27%
- Reduced rates: 5-15%
- Zero rated: 0%
- VAT number validation (VIES)

### VAT Number Validation
```typescript
// Verify VAT number
GET /api/tax/vat-validate?country=DE&number=123456789
```

### Reverse Charge
- B2B within EU with valid VAT
- Customer pays tax in their country

## Compliance

### Collecting Countries
| Country | Threshold | Rate |
|---------|-----------|------|
| Germany | €10,000 | 19% |
| France | €10,000 | 20% |
| UK | £0 | 20% |
| Spain | €10,000 | 21% |
| Italy | €10,000 | 22% |

### Filing Requirements

| Country | Frequency | Due Date |
|---------|-----------|----------|
| Germany | Monthly/Quarterly | 10th |
| France | Monthly | 24th |
| UK | Quarterly | 7th |
| OSS (EU) | Quarterly | 20th |

## Tax Reports

### Sales by Tax Rate
```json
{
  "period": "2026-01",
  "taxes": [
    {
      "rate": 19,
      "netAmount": 10000,
      "taxAmount": 1900
    },
    {
      "rate": 7,
      "netAmount": 5000,
      "taxAmount": 350
    }
  ]
}
```

### Tax Liability Report
```json
{
  "period": "2026-01",
  "totalCollected": 2250,
  "totalPaid": 500,
  "netLiability": 1750,
  "dueDate": "2026-02-10"
}
```

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
