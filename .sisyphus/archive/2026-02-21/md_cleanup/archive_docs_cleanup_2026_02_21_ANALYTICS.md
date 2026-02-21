# Analytics

## Overview
Analytics tracking and reporting for BIOMETRICS.

## Tracked Events

### User Events
- Page views
- Session duration
- User actions
- Sign ups
- Conversions

### Business Metrics
- Revenue (via Stripe)
- Active users
- Retention rates
- Feature usage

## Implementation

### Event Tracking
```typescript
// lib/analytics.ts
interface AnalyticsEvent {
  name: string;
  properties?: Record<string, any>;
  timestamp: Date;
  userId?: string;
}

async function track(event: AnalyticsEvent) {
  await supabase.from('analytics_events').insert({
    event_name: event.name,
    properties: event.properties,
    user_id: event.userId,
    created_at: event.timestamp
  });
}
```

### API Endpoints
```
POST /api/analytics/track

GET /api/analytics/dashboard
GET /api/analytics/users
GET /api/analytics/events
GET /api/analytics/revenue
```

## Dashboards
- Supabase Analytics (built-in)
- Custom dashboard: `/admin/analytics`

## Privacy
- GDPR compliant
- IP anonymization
- Cookie consent banner
- Data retention: 2 years
