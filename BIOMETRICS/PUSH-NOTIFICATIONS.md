# Push Notifications

## Overview
Real-time push notification system for BIOMETRICS web application.

## Features
- Browser push notifications via Web Push API
- Device token management
- Notification scheduling
- Delivery tracking

## Implementation

### Service Worker Setup
```javascript
// public/sw.js
self.addEventListener('push', (event) => {
  const data = event.data.json();
  self.registration.showNotification(data.title, {
    body: data.body,
    icon: '/icon.png',
    badge: '/badge.png'
  });
});
```

### Subscription Endpoint
```
POST /api/notifications/subscribe
Body: { endpoint, keys: { p256dh, auth } }
```

### Send Notification
```
POST /api/notifications/send
Body: { userId, title, body, data }
```

## Providers
- Web Push API (native browser)
- Fallback: In-app notifications

## Configuration
```typescript
// config/notifications.ts
export const pushConfig = {
  vapidPublicKey: process.env.VAPID_PUBLIC_KEY,
  ttl: 86400,
  priority: 'high'
};
```
