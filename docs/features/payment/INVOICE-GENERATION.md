# INVOICE-GENERATION.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The Invoice Generation module creates, manages, and delivers professional invoices to customers. This document covers invoice creation, customization, delivery, and archival.

## Invoice Types

| Type | Description | Trigger |
|------|-------------|---------|
| Subscription | Recurring billing | Subscription cycle |
| One-time | Single purchase | Direct payment |
| Usage | Overage charges | Usage threshold |
| Proforma | Pre-billing quote | Quote acceptance |

## Database Schema

```sql
-- Invoice templates
CREATE TABLE invoice_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    name VARCHAR(100) NOT NULL,
    logo_url TEXT,
    company_name VARCHAR(255),
    company_address TEXT,
    company_email VARCHAR(255),
    company_phone VARCHAR(50),
    tax_id VARCHAR(50),
    bank_details JSONB,
    default_payment_terms TEXT,
    footer_text TEXT,
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Invoices
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    user_id UUID REFERENCES auth.users(id),
    subscription_id UUID REFERENCES payment_subscriptions(id),
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'draft',
    issue_date DATE NOT NULL,
    due_date DATE NOT NULL,
    paid_date DATE,
    currency VARCHAR(3) NOT NULL,
    subtotal DECIMAL(12,2) NOT NULL,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    total DECIMAL(12,2) NOT NULL,
    amount_paid DECIMAL(12,2) DEFAULT 0,
    amount_due DECIMAL(12,2) NOT NULL,
    notes TEXT,
    terms TEXT,
    template_id UUID REFERENCES invoice_templates(id),
    pdf_url TEXT,
    sent_at TIMESTAMPTZ,
    viewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Invoice line items
CREATE TABLE invoice_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID REFERENCES invoices(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    quantity DECIMAL(10,2) NOT NULL,
    unit_price DECIMAL(12,2) NOT NULL,
    tax_rate DECIMAL(5,2),
    amount DECIMAL(12,2) NOT NULL,
    period_start DATE,
    period_end DATE,
    metadata JSONB,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Invoice payments
CREATE TABLE invoice_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID REFERENCES invoices(id) ON DELETE CASCADE,
    payment_id UUID REFERENCES payment_transactions(id),
    amount DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## API Endpoints

### Invoices

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/invoices | List invoices |
| GET | /api/invoices/:id | Get invoice |
| POST | /api/invoices | Create invoice |
| PUT | /api/invoices/:id | Update invoice |
| DELETE | /api/invoices/:id | Delete draft invoice |
| POST | /api/invoices/:id/send | Send invoice |
| POST | /api/invoices/:id/void | Void invoice |
| GET | /api/invoices/:id/pdf | Download PDF |
| POST | /api/invoices/:id/remind | Send reminder |

### Templates

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/invoices/templates | List templates |
| POST | /api/invoices/templates | Create template |
| PUT | /api/invoices/templates/:id | Update template |
| DELETE | /api/invoices/templates/:id | Delete template |

## Invoice Numbering

### Format

```
{prefix}{year}{month}{sequence}
```

### Examples

| Type | Format | Example |
|------|--------|---------|
| Invoice | INV-2026020001 | INV-2026020001 |
| Proforma | PRO-2026020001 | PRO-2026020001 |
| Credit Note | CN-2026020001 | CN-2026020001 |

### Sequence Reset
- Monthly: Sequence resets each month
- Yearly: Sequence resets each year
- Continuous: No reset

## PDF Generation

### Technology
- Puppeteer for HTML to PDF
- Custom Handlebars templates
- CDN storage for PDFs

### Template Variables

```handlebars
{{invoice_number}}
{{issue_date}}
{{due_date}}
{{customer.name}}
{{customer.email}}
{{customer.address}}
{{line_items[]}}
{{subtotal}}
{{tax_amount}}
{{total}}
{{company.logo}}
{{company.name}}
{{company.address}}
{{payment_details}}
```

## Delivery Methods

| Method | Description | Setup |
|--------|-------------|-------|
| Email | PDF attached to email | SMTP |
| Download | Direct link | CDN |
| Mail | Physical copy | Print service |
| API | JSON/PDF via API | Webhook |

## Payment Terms

### Common Terms

| Term | Net Due |
|------|---------|
| Net 7 | 7 days |
| Net 14 | 14 days |
| Net 30 | 30 days |
| Net 60 | 60 days |
| Due on Receipt | Immediate |
| 50% Advance | 50% now |

## Reminder Schedule

| Reminder | Timing |
|----------|--------|
| First | 3 days before due |
| Second | Due date |
| Third | 3 days overdue |
| Final | 7 days overdue |

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
