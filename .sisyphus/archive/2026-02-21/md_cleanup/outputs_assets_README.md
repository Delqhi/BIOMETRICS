# Assets Output Directory

## Overview

This directory contains processed and optimized asset files generated from source materials. These are the final assets ready for deployment.

## Contents

### Optimized Images

| Directory | Description |
|-----------|-------------|
| images/ | Optimized image files |
| thumbnails/ | Thumbnail versions |
| responsive/ | Responsive image sets |

### Processed Graphics

| Directory | Description |
|-----------|-------------|
| icons/ | Processed icon files |
| logos/ | Final logo variants |
| diagrams/ | Rendered diagrams |

### Web Assets

| Directory | Description |
|-----------|-------------|
| sprites/ | Image sprites |
| inlines/ | Inline optimized assets |
| fonts/ | Web font files |

## Optimization

### Image Optimization
```bash
# Optimize images
for f in *.png; do
  optipng "$f"
done

# Convert to WebP
for f in *.jpg; do
  cwebp -q 80 "$f" -o "${f%.jpg}.webp"
done
```

### Sprite Generation
```bash
# Create sprite
sprite-generator --input icons/*.svg --output sprites/sprite.png
```

## Responsive Images

### HTML Usage
```html
<picture>
  <source srcset="image-400.webp 400w,
                  image-800.webp 800w"
          type="image/webp">
  <img src="image-800.jpg" alt="Description">
</picture>
```

## CDN Configuration

### Upload Script
```bash
# Upload to CDN
./scripts/upload-assets.sh --env production
```

### Cache Headers
```yaml
cdn:
  assets:
    - path: /assets/*
      cache: 1y
      gzip: true
```

## Version Control

### Don't Track
```gitignore
# Optimized assets
*.min.js
*.min.css
*.webp

# Generated sprites
sprites/*.png
```

## Organization

### Directory Structure
```
assets/
├── 2026-02/
│   ├── home-hero.webp
│   └── product-card.webp
└── 2026-01/
    └── blog-header.webp
```

## Cleanup

### Remove Unused
```bash
# Find unused assets
biometrics-cli assets find-unused

# Remove old versions
biometrics-cli assets clean --older-than 90d
```

## Quality Control

### Validation
```bash
# Validate images
biometrics-cli assets validate --quality 90

# Check dimensions
biometrics-cli assets check-size --max-width 1920
```

## Maintenance

- Review weekly
- Optimize new assets
- Remove duplicates

## Best Practices

1. **Optimize First**: Optimize before commit
2. **Use CDN**: Serve via CDN
3. **Version**: Use cache busting
4. **Monitor**: Track usage

## See Also

- [Assets Overview](../assets/)
- [Videos Output](./videos/)
- [Outputs Overview](../outputs/)
