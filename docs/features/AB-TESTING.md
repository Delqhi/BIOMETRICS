# A/B Testing

## Overview
A/B testing framework for BIOMETRICS feature experimentation.

## Features
- Server-side A/B allocation
- Feature flags integration
- Statistical significance tracking
- Goal conversion tracking

## Implementation

### Experiment Definition
```typescript
// types/experiment.ts
interface Experiment {
  id: string;
  name: string;
  variants: Variant[];
  traffic: number; // 0-100%
  goals: string[];
}

interface Variant {
  id: string;
  name: string;
  weight: number;
}
```

### Allocation
```typescript
// lib/ab-testing.ts
function allocateUser(userId: string, experiment: Experiment): string {
  const hash = md5(`${userId}:${experiment.id}`);
  const bucket = parseInt(hash.slice(0, 8), 16) % 100;
  
  if (bucket >= experiment.traffic) return 'control';
  
  let cumulative = 0;
  for (const variant of experiment.variants) {
    cumulative += variant.weight;
    if (bucket < cumulative) return variant.id;
  }
  return 'control';
}
```

### API Endpoints
```
GET /api/experiments
GET /api/experiments/:id
POST /api/experiments/:id/variant
POST /api/experiments/:id/conversion
```

## Tools
- Custom implementation
- Supabase for data storage
- No external A/B testing service required
