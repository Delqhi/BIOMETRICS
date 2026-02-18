# SOCIAL-MEDIA-API.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The Social Media API enables BIOMETRICS to integrate with major social media platforms for content publishing, analytics, and social login. This document covers the architecture, supported platforms, and API specifications.

## Supported Platforms

| Platform | Features | Status |
|----------|----------|--------|
| Facebook | Login, Share, Analytics | Active |
| Instagram | Media upload, Insights | Active |
| Twitter/X | Post, Analytics, Login | Active |
| LinkedIn | Share, Analytics | Active |
| TikTok | Video upload, Analytics | Beta |
| YouTube | Video upload, Analytics | Beta |

## Architecture

### Components

| Component | Description | Technology |
|-----------|-------------|------------|
| OAuth Manager | Platform authentication | Node.js |
| Content Scheduler | Post scheduling | PostgreSQL + Cron |
| Media Processor | Image/video optimization | FFmpeg + Sharp |
| Analytics Aggregator | Platform analytics | Supabase Functions |
| Webhook Handler | Real-time events | Supabase Webhooks |

### Integration Flow

```
BIOMETRICS App          OAuth Flow             Social Platform
     |                      |                        |
     |--- Request Login ---->|                        |
     |<-- Redirect ---------|                        |
     |                      |--- OAuth Authorize --->|
     |<-- Auth Code --------|                        |
     |                      |--- Exchange Token ---->|
     |<-- Access Token -----|                        |
     |                      |                        |
     |--- API Request ----->|                        |
     |<-- Response ---------|                        |
```

## Database Schema

```sql
-- Connected social accounts
CREATE TABLE social_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id),
    platform VARCHAR(50) NOT NULL,
    platform_user_id VARCHAR(255),
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    token_expires_at TIMESTAMPTZ,
    scope TEXT[],
    profile_data JSONB,
    is_active BOOLEAN DEFAULT true,
    connected_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Scheduled posts
CREATE TABLE social_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id),
    platform VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    media_urls TEXT[],
    scheduled_at TIMESTAMPTZ,
    posted_at TIMESTAMPTZ,
    platform_post_id VARCHAR(255),
    platform_post_url VARCHAR(500),
    status VARCHAR(20) DEFAULT 'pending',
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Post analytics
CREATE TABLE social_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID REFERENCES social_posts(id),
    platform VARCHAR(50) NOT NULL,
    impressions INT DEFAULT 0,
    engagements INT DEFAULT 0,
    likes INT DEFAULT 0,
    comments INT DEFAULT 0,
    shares INT DEFAULT 0,
    clicks INT DEFAULT 0,
    recorded_at TIMESTAMPTZ DEFAULT NOW()
);

-- Webhook deliveries
CREATE TABLE social_webhooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    platform VARCHAR(50) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    processed BOOLEAN DEFAULT false,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## API Endpoints

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/social/connect/:platform | Initiate OAuth |
| GET | /api/social/callback/:platform | OAuth callback |
| DELETE | /api/social/accounts/:id | Disconnect account |
| GET | /api/social/accounts | List connected accounts |

### Publishing

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/social/post | Create post |
| GET | /api/social/posts | List posts |
| GET | /api/social/posts/:id | Get post details |
| PUT | /api/social/posts/:id | Update post |
| DELETE | /api/social/posts/:id | Delete post |
| POST | /api/social/posts/:id/schedule | Schedule post |

### Analytics

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/social/analytics | Get analytics summary |
| GET | /api/social/analytics/:postId | Get post analytics |
| GET | /api/social/analytics/trends | Get trend data |

### Webhooks

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/social/webhooks | Register webhook |
| GET | /api/social/webhooks | List webhooks |
| DELETE | /api/social/webhooks/:id | Remove webhook |

## OAuth Configuration

### Facebook

```javascript
const facebookOAuth = {
  authorizationURL: 'https://www.facebook.com/v18.0/dialog/oauth',
  tokenURL: 'https://graph.facebook.com/v18.0/oauth/access_token',
  scope: [
    'public_profile',
    'email',
    'pages_read_engagement',
    'pages_manage_posts'
  ]
};
```

### Twitter/X

```javascript
const twitterOAuth = {
  authorizationURL: 'https://twitter.com/i/oauth2/authorize',
  tokenURL: 'https://api.twitter.com/2/oauth2/token',
  scope: ['tweet.read', 'tweet.write', 'users.read']
};
```

### LinkedIn

```javascript
const linkedinOAuth = {
  authorizationURL: 'https://www.linkedin.com/oauth/v2/authorization',
  tokenURL: 'https://www.linkedin.com/oauth/v2/accessToken',
  scope: ['r_emailaddress', 'r_liteprofile', 'w_member_social']
};
```

## Rate Limits

| Platform | Requests | Window |
|----------|----------|--------|
| Facebook Graph API | 200 | 1 hour |
| Twitter API v2 | 500K | 1 month |
| LinkedIn API | 100 | 1 day |
| Instagram API | 200 | 1 hour |

## Media Requirements

### Images

| Platform | Max Size | Formats | Dimensions |
|----------|----------|---------|------------|
| Facebook | 8MB | JPG, PNG | 1200x630 |
| Twitter | 5MB | JPG, PNG, GIF | 1200x675 |
| LinkedIn | 5MB | JPG, PNG | 1200x627 |
| Instagram | 8MB | JPG, PNG | 1080x1080 |

### Videos

| Platform | Max Size | Formats | Duration |
|----------|----------|---------|----------|
| Facebook | 4GB | MP4, MOV | 240 min |
| Twitter | 512MB | MP4, MOV | 2:20 |
| LinkedIn | 5GB | MP4 | 10 min |
| Instagram | 4GB | MP4 | 60 sec |

## Error Handling

| Code | Description | Resolution |
|------|-------------|------------|
| TOKEN_EXPIRED | Access token expired | Refresh token or re-authenticate |
| RATE_LIMIT | API rate limit reached | Wait and retry with backoff |
| PERMISSION_DENIED | Insufficient permissions | Request additional scope |
| CONTENT_REJECTED | Post violates platform policy | Modify content |

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
