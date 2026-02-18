# SESSION-MGMT.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The Session Management module handles user authentication sessions, device management, and session security. This document covers session lifecycle, security features, and management APIs.

## Session Types

| Type | Duration | Use Case |
|------|----------|----------|
| Web | 30 days | Browser sessions |
| Mobile | 1 year | Native apps |
| API | 1 hour | API access tokens |
| Refresh | 1 week | Token refresh |

## Database Schema

```sql
-- User sessions
CREATE TABLE auth_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    refresh_token VARCHAR(255) UNIQUE,
    device_info JSONB,
    ip_address INET,
    user_agent TEXT,
    location JSONB,
    expires_at TIMESTAMPTZ NOT NULL,
    last_activity_at TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Device management
CREATE TABLE user_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    device_id VARCHAR(255) NOT NULL,
    name VARCHAR(100),
    type VARCHAR(50),
    os VARCHAR(50),
    browser VARCHAR(50),
    last_seen TIMESTAMPTZ DEFAULT NOW(),
    is_trusted BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, device_id)
);

-- Session events
CREATE TABLE auth_session_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID REFERENCES auth_sessions(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Token blacklist
CREATE TABLE auth_token_blacklist (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token VARCHAR(255) UNIQUE NOT NULL,
    user_id UUID REFERENCES auth.users(id),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## API Endpoints

### Sessions

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/auth/sessions | List active sessions |
| GET | /api/auth/sessions/:id | Get session details |
| DELETE | /api/auth/sessions/:id | Terminate session |
| DELETE | /api/auth/sessions | Terminate all sessions |

### Devices

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/auth/devices | List devices |
| POST | /api/auth/devices/:id/trust | Trust device |
| DELETE | /api/auth/devices/:id | Remove device |

### Tokens

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/auth/token/refresh | Refresh token |
| POST | /api/auth/token/revoke | Revoke token |

## Token Structure

### Access Token (JWT)

```typescript
interface AccessToken {
  sub: string;           // User ID
  iat: number;           // Issued at
  exp: number;           // Expires at
  type: 'access';
  permissions: string[];
  session_id: string;
}
```

### Refresh Token

```typescript
interface RefreshToken {
  sub: string;           // User ID
  iat: number;
  exp: number;
  type: 'refresh';
  session_id: string;
  token_id: string;
}
```

## Session Lifecycle

```
User Login
    |
    v
Create Session
    |
    v
Generate Tokens --> Access + Refresh
    |
    v
Store Session
    |
    v
Use Access Token
    |
    v
Access Expired?
    |---- No --> Continue
    |
    Yes
    |
    v
Refresh Token
    |
    v
Valid Refresh?
    |---- No --> Require Login
    |
    Yes
    |
    v
Issue New Access Token
    |
    v
Continue
```

## Security Features

### Session Security

| Feature | Description |
|---------|-------------|
| Token rotation | New refresh token each use |
| IP validation | Verify IP consistency |
| Device fingerprint | Validate device |
| Activity tracking | Log all actions |
| Absolute timeout | Max session length |

### Threat Protection

- Token hijacking detection
- Concurrent session limits
- Suspicious activity alerts
- Automatic session termination

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
