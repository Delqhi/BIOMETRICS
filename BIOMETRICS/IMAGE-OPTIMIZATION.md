# Image Optimization

## Overview
Automated image optimization pipeline for BIOMETRICS.

## Features
- Automatic format conversion (WebP, AVIF)
- Responsive image generation
- Lazy loading
- Placeholder generation

## Implementation

### Image Service
```typescript
// services/imageOptimizer.ts
interface ImageOptions {
  width?: number;
  height?: number;
  format?: 'webp' | 'avif' | 'jpeg';
  quality?: number;
}

async function optimizeImage(
  input: Buffer, 
  options: ImageOptions
): Promise<Buffer>
```

### API Endpoints
```
POST /api/images/optimize
Body: { imageUrl, options }

GET /api/images/responsive/:imageId
  ?srcSet=100,200,400,800
```

### Sizes Generated
- Thumbnail: 100px
- Small: 400px  
- Medium: 800px
- Large: 1200px
- Original: preserved

## Storage
- Original: Supabase Storage `/originals/`
- Optimized: Supabase Storage `/optimized/`

## Tools
- Sharp (Node.js image processing)
- Cloudflare Images (CDN)
