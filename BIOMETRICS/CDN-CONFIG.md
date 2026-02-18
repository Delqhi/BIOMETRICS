# CDN Configuration

## Overview
Content Delivery Network setup for BIOMETRICS static assets.

## CDN Provider
Cloudflare (integrated with existing infrastructure)

## Configuration

### Zones
- `biometrics.app` - Main application
- `cdn.biometrics.app` - Static assets

### Caching Rules

| Path Pattern | Cache Level | TTL |
|--------------|-------------|-----|
| `/static/*` | Cache Everything | 1 year |
| `/images/*` | Cache Everything | 1 year |
| `/fonts/*` | Cache Everything | 1 year |
| `/api/*` | Bypass | 0 |
| `/*` | Standard | 1 hour |

### Edge Caching
- Static assets: Cloudflare CDN
- HTML: Cache with stale-while-revalidate
- API responses: No cache (dynamic)

## Asset Optimization
- Automatic Brotli compression
- HTTP/3 support
- Image optimization via Cloudflare Images

## Environment Variables
```
CDN_URL=https://cdn.biometrics.app
CLOUDFLARE_ZONE_ID=xxx
CLOUDFLARE_API_TOKEN=xxx
```
