# Search Engine

## Overview
Full-text search implementation for BIOMETRICS.

## Search Providers
1. **Primary**: Supabase pgvector (PostgreSQL full-text search)
2. **Fallback**: FlexSearch (client-side)

## Indexing

### Indexed Content
- User profiles
- Content/Media metadata
- Products
- Categories
- Tags

### Search Types
- Full-text search with ranking
- Fuzzy matching
- Autocomplete suggestions

## API Endpoints

```
GET /api/search
  ?q=query
  &type=user|content|product
  &limit=20
  &offset=0

GET /api/search/suggestions
  ?q=prefix
```

## Implementation
```typescript
// lib/search.ts
async function search(query: string, options: SearchOptions) {
  const results = await supabase
    .from('search_index')
    .select('*')
    .textSearch('searchable', query)
    .limit(options.limit);
  return results;
}
```

## Configuration
```typescript
// config/search.ts
export const searchConfig = {
  minQueryLength: 2,
  maxResults: 100,
  fuzzyThreshold: 0.3,
  highlight: true
};
```
