# WHITE-LABEL.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The White-Label module enables full customization of the BIOMETRICS platform for partner organizations. This document covers branding, customization, and subdomain configuration.

## White-Label Features

| Feature | Description |
|---------|-------------|
| Custom Domain | Use own domain |
| Branding | Logo, colors, images |
| Email | Custom email templates |
| Support | Custom help center |

## Database Schema

```sql
-- White-label configurations
CREATE TABLE white_label_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    custom_domain VARCHAR(255),
    ssl_enabled BOOLEAN DEFAULT false,
    logo_url TEXT,
    logo_dark_url TEXT,
    favicon_url TEXT,
    primary_color VARCHAR(7),
    secondary_color VARCHAR(7),
    accent_color VARCHAR(7),
    font_family VARCHAR(100),
    company_name VARCHAR(255),
    company_website VARCHAR(255),
    company_email VARCHAR(255),
    company_address TEXT,
    privacy_policy_url TEXT,
    terms_url TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Custom email templates
CREATE TABLE white_label_emails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    template_name VARCHAR(100) NOT NULL,
    subject VARCHAR(255),
    body_html TEXT,
    body_text TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, template_name)
);
```

## API Endpoints

### Configuration

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/whitelabel/config | Get config |
| PUT | /api/whitelabel/config | Update config |
| POST | /api/whitelabel/domain | Set domain |
| DELETE | /api/whitelabel/domain | Remove domain |

### Email Templates

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/whitelabel/emails | List templates |
| POST | /api/whitelabel/emails | Create template |
| PUT | /api/whitelabel/emails/:id | Update template |

## Configuration Options

### Branding

```typescript
interface WhiteLabelConfig {
  // Logo
  logo: {
    light: string;
    dark: string;
    favicon: string;
  };
  
  // Colors
  colors: {
    primary: '#4F46E5';
    secondary: '#10B981';
    accent: '#F59E0B';
    background: '#FFFFFF';
    text: '#1F2937';
  };
  
  // Fonts
  font: {
    family: 'Inter, sans-serif';
    headings: 'Poppins, sans-serif';
  };
  
  // Company
  company: {
    name: 'Acme Corp';
    website: 'https://acme.com';
    email: 'support@acme.com';
  };
}
```

### Custom Domain

```
Default: app.biometrics.com
White-label: app.acme.com
```

### DNS Configuration

| Record Type | Name | Value |
|-------------|------|-------|
| CNAME | app | cdn.biometrics.com |
| TXT | _acme-challenge | Verification token |

## Email Customization

### Available Templates

| Template | Description |
|----------|-------------|
| welcome | Welcome email |
| reset_password | Password reset |
| verify_email | Email verification |
| invoice | Invoice notification |
| subscription | Subscription alerts |

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
