# CHAT-WIDGET.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The Chat Widget provides real-time customer support and messaging capabilities integrated directly into the BIOMETRICS web application. This document describes the implementation, features, and configuration options.

## Architecture

### Components

| Component | Description | Technology |
|-----------|-------------|------------|
| Widget Frontend | Embedded chat interface | React + TypeScript |
| WebSocket Server | Real-time messaging | Supabase Realtime |
| Message Queue | Async message processing | PostgreSQL |
| Notification Service | Push notifications | Firebase/Supabase |
| AI Assistant | Bot responses | OpenAI integration |

### System Flow

```
User Browser          Supabase              AI Service           Support Team
    |                    |                      |                     |
    |--- Open Widget --->|                      |                     |
    |                    |                      |                     |
    |<-- Connection -----|                      |                     |
    |                    |                      |                     |
    |--- Send Message -->|                      |                     |
    |                    |--- Process AI ------>|                     |
    |                    |<-- AI Response ------|                     |
    |<-- Message --------|                      |                     |
    |                    |                      |                     |
    |--- Human Request -->|                      |                     |
    |                    |--------------------->|                     |
    |                    |<----------------------|                     |
    |<-- Human Reply -----|                      |                     |
```

## Database Schema

```sql
-- Conversations table
CREATE TABLE chat_conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id),
    status VARCHAR(50) DEFAULT 'active',
    assigned_to UUID REFERENCES auth.users(id),
    priority VARCHAR(20) DEFAULT 'normal',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    closed_at TIMESTAMPTZ
);

-- Messages table
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID REFERENCES chat_conversations(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES auth.users(id),
    sender_type VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text',
    metadata JSONB,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Quick replies / Canned responses
CREATE TABLE chat_quick_replies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trigger_keyword VARCHAR(100) UNIQUE,
    response_text TEXT NOT NULL,
    category VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Widget settings per tenant
CREATE TABLE chat_widget_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    primary_color VARCHAR(20) DEFAULT '#000000',
    position VARCHAR(20) DEFAULT 'bottom-right',
    greeting_message TEXT,
    is_active BOOLEAN DEFAULT true,
    operating_hours JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

## Widget Configuration

### Installation

```html
<!-- Add to your website -->
<script>
  window.BIOMETRICS_CHAT_CONFIG = {
    tenantId: 'your-tenant-id',
    position: 'bottom-right',
    primaryColor: '#4F46E5',
    greeting: 'Hello! How can we help you today?',
    logoUrl: 'https://your-logo.com/logo.png'
  };
</script>
<script src="https://chat.biometrics.com/widget.js" async></script>
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| tenantId | string | required | Your unique tenant identifier |
| position | string | 'bottom-right' | Widget position |
| primaryColor | string | '#000000' | Brand color |
| greeting | string | null | Welcome message |
| logoUrl | string | null | Custom logo |
| launcherText | string | 'Chat' | Button text |
| onlineHours | object | 24/7 | Operating hours |

## API Endpoints

### Conversations

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/chat/conversations | List user conversations |
| GET | /api/chat/conversations/:id | Get conversation details |
| POST | /api/chat/conversations | Start new conversation |
| PUT | /api/chat/conversations/:id | Update conversation |
| POST | /api/chat/conversations/:id/close | Close conversation |

### Messages

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/chat/conversations/:id/messages | Get messages |
| POST | /api/chat/conversations/:id/messages | Send message |
| PUT | /api/chat/messages/:id/read | Mark as read |

### Settings

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/chat/settings | Get widget settings |
| PUT | /api/chat/settings | Update settings |

## Features

### 1. Real-time Messaging
- WebSocket connection for instant messages
- Typing indicators
- Read receipts
- Message status (sent, delivered, read)

### 2. AI Assistant
- Automatic responses for common questions
- Intent recognition
- Escalation to human agents
- Learning from conversations

### 3. File Sharing
- Image uploads
- Document attachments
- Drag and drop support
- File type restrictions

### 4. Notification System
- Browser push notifications
- Email notifications for offline users
- Mobile app notifications
- Sound alerts (optional)

### 5. Analytics
- Message volume tracking
- Response time metrics
- User satisfaction scores
- Popular topics analysis

### 6. Multi-language Support
- Automatic language detection
- Translation integration
- RTL language support

## WebSocket Events

### Client → Server

| Event | Payload | Description |
|-------|---------|-------------|
| message | `{ conversationId, content, type }` | Send message |
| typing | `{ conversationId, isTyping }` | Typing indicator |
| join | `{ conversationId }` | Join conversation |
| leave | `{ conversationId }` | Leave conversation |

### Server → Client

| Event | Payload | Description |
|-------|---------|-------------|
| message | `{ message }` | New message received |
| typing | `{ userId, isTyping }` | User typing |
| read | `{ messageId }` | Message read |
| agent_joined | `{ agent }` | Agent joined |
| agent_left | `{ agent }` | Agent left |

## Security

- End-to-end encryption for sensitive data
- RLS policies on all tables
- Rate limiting on message sending
- Input sanitization
- XSS protection

## Performance

- Message pagination (50 per page)
- Lazy loading of older messages
- Connection pooling
- CDN for static assets

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
