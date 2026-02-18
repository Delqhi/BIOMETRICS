# User Feedback

## Overview
In-app user feedback collection system for BIOMETRICS.

## Features
- Star ratings (1-5)
- Text feedback
- Screenshot attachments
- Category tagging
- Upvote/downvote feedback

## Implementation

### Feedback Submission
```typescript
// types/feedback.ts
interface Feedback {
  id: string;
  userId?: string;
  type: 'bug' | 'feature' | 'improvement' | 'general';
  rating?: number;
  title: string;
  description: string;
  screenshotUrl?: string;
  createdAt: Date;
}
```

### API Endpoints
```
POST /api/feedback
Body: { type, rating, title, description, screenshot }

GET /api/feedback
  ?status=pending|reviewed|resolved

PUT /api/feedback/:id/status
Body: { status, response }

POST /api/feedback/:id/upvote
POST /api/feedback/:id/downvote
```

### Screenshot Upload
```
POST /api/feedback/screenshot
Returns: { url: string }
```

## Moderation
- Admin review queue
- Mark as resolved
- Reply to feedback

## Storage
- Feedback: Supabase `feedback` table
- Screenshots: Supabase Storage `/feedback/screenshots/`

## Analytics
- Average rating dashboard
- Feedback by category
- Trending issues
