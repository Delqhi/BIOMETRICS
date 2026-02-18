# TWO-FACTOR-AUTH.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The Two-Factor Authentication (2FA) module provides multiple authentication methods to secure user accounts. This document covers supported methods, implementation, and security features.

## Supported Methods

| Method | Type | Security Level | User Experience |
|--------|------|----------------|-----------------|
| TOTP | App-based | High | Good |
| SMS | Phone-based | Medium | Good |
| Email | Backup | Low | Fair |
| WebAuthn | Hardware Key | Very High | Excellent |
 | Recovery | Low| Backup Codes | N/A |

## Database Schema

```sql
-- 2FA methods per user
CREATE TABLE auth_2fa_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    method VARCHAR(50) NOT NULL,
    identifier VARCHAR(255),
    encrypted_secret TEXT,
    verified_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Backup codes
CREATE TABLE auth_backup_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    code_hash VARCHAR(255) NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2FA settings
CREATE TABLE auth_2fa_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    require_for_login BOOLEAN DEFAULT false,
    require_for_transactions BOOLEAN DEFAULT false,
    method_preference VARCHAR(50) DEFAULT 'totp',
    trusted_devices JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2FA attempts
CREATE TABLE auth_2fa_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id),
    method VARCHAR(50) NOT NULL,
    success BOOLEAN NOT NULL,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## API Endpoints

### Setup

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/auth/2fa/methods | List enabled methods |
| POST | /api/auth/2fa/setup/totp | Setup TOTP |
| POST | /api/auth/2fa/setup/sms | Setup SMS |
| POST | /api/auth/2fa/setup/webauthn | Setup WebAuthn |
| POST | /api/auth/2fa/setup/backup-codes | Generate backup codes |

### Verification

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/auth/2fa/verify | Verify 2FA code |
| POST | /api/auth/2fa/verify-backup | Verify backup code |
| POST | /api/auth/2fa/resend | Resend code |

### Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| DELETE | /api/auth/2fa/methods/:id | Remove method |
| POST | /api/auth/2fa/disable | Disable 2FA |
| GET | /api/auth/2fa/trusted-devices | List trusted devices |

## TOTP Implementation

### Setup Flow

```
1. Generate secret key
   --> Base32 encoded

2. Display QR code
   --> otpauth://totp/{issuer}:{email}?secret={secret}&issuer={issuer}

3. User scans with app
   --> Google Authenticator, Authy, etc.

4. Verify first code
   --> Validate and save
```

### Secret Generation

```typescript
import speakeasy from 'speakeasy';

const secret = speakeasy.generateSecret({
  name: 'BIOMETRICS',
  length: 32
});

// Result
{
  ascii: 'random32charstring',
  base32: 'JBSWY3DPEHPK3PXP',
  otpauth_url: 'otpauth://totp/BIOMETRICS:user@example.com?secret=JBSWY3DPEHPK3PXP'
}
```

### Verification

```typescript
const verified = speakeasy.totp.verify({
  secret: userSecret,
  encoding: 'base32',
  token: userToken,
  window: 1  // Allow 1 step tolerance
});
```

## WebAuthn Implementation

### Registration

```typescript
// Server generates options
const options = await webAuthnServer.generateRegistrationOptions({
  rpName: 'BIOMETRICS',
  rpID: 'biometrics.com',
  userID: userId,
  userName: userEmail
});

// Client creates credential
const credential = await navigator.credentials.create({
  publicKey: options
});

// Server verifies
await webAuthnServer.verifyRegistrationResponse(credential);
```

### Authentication

```typescript
// Server generates options
const options = await webAuthnServer.generateAuthenticationOptions({
  allowCredentials: userCredentials
});

// Client gets assertion
const assertion = await navigator.credentials.get({
  publicKey: options
});

// Server verifies
await webAuthnServer.verifyAuthenticationResponse(assertion);
```

## Security Features

### Rate Limiting

| Action | Limit | Window |
|--------|-------|--------|
| TOTP verify | 10 | 5 min |
| SMS send | 5 | 1 hour |
| Backup code | 10 | 10 min |

### Lockout

- 10 failed attempts → 15 min lockout
- Account flagged for review

### Trust Devices

- Remember device for 30 days
- Require 2FA on new devices

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
