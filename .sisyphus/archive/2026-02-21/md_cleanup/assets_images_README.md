# Images Directory

## Overview

This directory contains all general image assets used throughout the biometrics project. These include screenshots, photographs, icons, and miscellaneous visual elements.

## Contents

### Screenshots

| Directory | Description |
|-----------|-------------|
| screenshots/ | Application screenshots |
| screenshots/web/ | Web interface screenshots |
| screenshots/mobile/ | Mobile app screenshots |
| screenshots/cli/ | CLI tool screenshots |

### Photographs

| Directory | Description |
|-----------|-------------|
| photos/ | Team and event photos |
| photos/team/ | Team member photos |
| photos/events/ | Conference/event photos |

### UI Elements

| Directory | Description |
|-----------|-------------|
| ui/ | UI components and elements |
| ui/buttons/ | Button states |
| ui/forms/ | Form element styles |
| ui/icons/ | Custom icons |

## Image Standards

### Format Guidelines

| Use Case | Format | Quality |
|----------|--------|---------|
| Photographs | JPEG | 85% |
| Screenshots | PNG | Lossless |
| Icons | SVG | Vector |
| Photos | WebP | 80% |

### Resolution Guidelines

| Type | Minimum | Recommended |
|------|---------|-------------|
| Hero Images | 1920x1080 | 3840x2160 |
| Thumbnails | 300x200 | 600x400 |
| Icons | 24x24 | 64x64 (vector) |

## Naming Convention

Format: `{category}-{description}-{version}.{ext}`

Examples:
- `screenshot-login-v2.png`
- `photo-team-2026.jpg`
- `icon-user-active.svg`

## Optimization

### Automated Optimization
```bash
# Optimize images
biometrics-cli assets optimize \
  --input ./images \
  --output ./images-optimized
```

### WebP Conversion
```bash
# Convert to WebP
for f in *.png; do
  cwebp -q 80 "$f" -o "${f%.png}.webp"
done
```

## Accessibility

### Alt Text Requirements
All images must include descriptive alt text:
```html
<img src="screenshot-login.png" 
     alt="Biometrics login page with email and password fields">
```

### Contrast
- Minimum contrast ratio: 4.5:1
- Text on images: WCAG AA compliant

## Version Control

Track large files in Git LFS:
```bash
git lfs track "*.png"
git lfs track "*.jpg"
git lfs track "*.webp"
```

## Tools

### Image Editing
- **GIMP**: Free image editor
- **Figma**: UI design
- **Sketch**: Mac design tool

### Optimization
- **ImageMagick**: Batch processing
- **Squoosh**: Web optimization
- **TinyPNG**: PNG compression

## Maintenance

### Regular Tasks
- Compress new images
- Remove unused files
- Update outdated screenshots

### Cleanup Schedule
```bash
# Monthly cleanup
biometrics-cli assets clean --older-than 90d
```

## See Also

- [Renders Directory](../renders/)
- [Logos Directory](../logos/)
- [Diagrams Directory](../diagrams/)
