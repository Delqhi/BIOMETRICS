# Logos Directory

## Overview

This directory contains all logo assets for the biometrics brand. Logos are available in various formats and sizes for different use cases.

## Contents

### Primary Logo

| File | Description | Format |
|------|-------------|--------|
| biometrics-logo.svg | Primary logo (vector) | SVG |
| biometrics-logo.png | Primary logo (raster) | PNG |
| biometrics-logo-dark.png | Dark background version | PNG |

### Wordmark

| File | Description | Format |
|------|-------------|--------|
| biometrics-wordmark.svg | Wordmark (horizontal) | SVG |
| biometrics-wordmark-stacked.svg | Stacked version | SVG |

### Icon

| File | Description | Format |
|------|-------------|--------|
| biometrics-icon.svg | Square icon (vector) | SVG |
| biometrics-icon-16.png | 16x16 favicon | PNG |
| biometrics-icon-32.png | 32x32 favicon | PNG |
| biometrics-icon-192.png | 192x192 PWA icon | PNG |
| biometrics-icon-512.png | 512x512 PWA icon | PNG |

### Variations

| File | Description |
|------|-------------|
| logo-white.svg | White (for dark backgrounds) |
| logo-black.svg | Black (for light backgrounds) |
| logo-mono.svg | Single color (print) |

## Logo Guidelines

### Clear Space
Maintain minimum clear space around logo:
- Height: 1x (based on 'B' height)
- All sides: 1x clear space

### Minimum Sizes
- Print: 20mm width
- Digital: 32px width
- Favicon: 16x16px

### Incorrect Usage
- ❌ Don't stretch or distort
- ❌ Don't change colors
- ❌ Don't add effects
- ❌ Don't place on busy backgrounds

## Color Values

### Primary Colors
| Name | Hex | RGB |
|------|-----|-----|
| Blue | #2563EB | 37, 99, 235 |
| Dark | #1E293B | 30, 41, 59 |

### Secondary Colors
| Name | Hex | RGB |
|------|-----|-----|
| Slate | #64748B | 100, 116, 139 |
| White | #FFFFFF | 255, 255, 255 |

## File Formats

### When to Use SVG
- Web (responsive)
- Print materials
- Large format
- Any scalable use

### When to Use PNG
- Legacy systems
- Email signatures
- Simple placements

### When to Use ICO
- Windows favicon
- App icons

## Tools

### Export Script
```bash
# Export all sizes
./scripts/export-logos.sh
```

### Favicon Generation
```bash
# Generate favicons
convert biometrics-icon.svg -resize 16x16 favicon-16.png
convert biometrics-icon.svg -resize 32x32 favicon-32.png
```

## Version Control

```bash
git lfs track "*.png"
git lfs track "*.svg"
git lfs track "*.ico"
```

## Usage Examples

### HTML
```html
<link rel="icon" type="image/png" href="/assets/logos/biometrics-icon-32.png">
<img src="/assets/logos/biometrics-logo.svg" alt="Biometrics Logo">
```

### CSS
```css
.logo {
  background-image: url('/assets/logos/biometrics-logo.svg');
}
```

## Maintenance

- Review logo usage annually
- Update for brand changes
- Archive old versions

## See Also

- [SVG Logos](./svg/)
- [Images Directory](../images/)
- [Icons Directory](../icons/)
