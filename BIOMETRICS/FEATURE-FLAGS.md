# Feature Flags

## Overview
Feature flag system for BIOMETRICS gradual rollouts.

## Features
- Server-side feature toggles
- User-based targeting
- Percentage rollouts
- A/B test integration

## Implementation

### Flag Configuration
```typescript
// types/feature-flags.ts
interface FeatureFlag {
  key: string;
  name: string;
  description: string;
  enabled: boolean;
  rolloutPercent?: number;
  targetUsers?: string[];
  excludeUsers?: string[];
}
```

### Checking Flags
```typescript
// lib/feature-flags.ts
async function isFeatureEnabled(
  flagKey: string, 
  userId?: string
): Promise<boolean> {
  const flag = await getFlag(flagKey);
  
  if (!flag.enabled) return false;
  if (flag.excludeUsers?.includes(userId)) return false;
  if (flag.targetUsers?.includes(userId)) return true;
  if (flag.rolloutPercent) {
    const hash = md5(`${userId}:${flagKey}`);
    const bucket = parseInt(hash.slice(0, 8), 16) % 100;
    return bucket < flag.rolloutPercent;
  }
  
  return true;
}
```

### API Endpoints
```
GET /api/features
GET /api/features/:key
POST /api/features
PUT /api/features/:key
```

## Management UI
- Admin panel: `/admin/features`
- Toggle flags without deployment
- Rollout percentage adjustment

## Storage
- Flags stored in Supabase `feature_flags` table
- Cached in Redis for performance
