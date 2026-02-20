# Rendered Assets Directory

## Overview

This directory contains rendered images and visual outputs from 3D models and design mockups. These renders are used for marketing materials, presentations, and documentation.

## Contents

### Marketing Renders

| File | Description | Resolution | Format |
|------|-------------|------------|--------|
| hero-product.png | Main product shot | 4K | PNG |
| feature-face-id.png | Face ID feature | 2K | PNG |
| feature-voice.png | Voice recognition | 2K | PNG |
| device-angle-1.png | Product angle 1 | 4K | PNG |
| device-angle-2.png | Product angle 2 | 4K | PNG |

### Presentation Renders

| File | Description | Resolution |
|------|-------------|------------|
| architecture-diagram.png | System architecture | 1920x1080 |
| data-flow.png | Data flow visualization | 1920x1080 |
| security-model.png | Security model | 1920x1080 |

### Style Guidelines

#### Color Palette
- **Primary**: #2563EB (Blue)
- **Secondary**: #1E293B (Dark slate)
- **Accent**: #10B981 (Green)
- **Background**: #F8FAFC (Light)

#### Typography
- **Headings**: Inter Bold
- **Body**: Inter Regular
- **Code**: JetBrains Mono

### Lighting Setup

For product renders:
- **Key Light**: 45° above, 45° right
- **Fill Light**: Soft, 30% intensity
- **Rim Light**: Backlight for edge definition
- **Environment**: HDRI studio

### Export Settings

#### PNG (High Quality)
- Compression: None
- Color: RGBA
- Bit Depth: 16-bit

#### JPEG (Web)
- Quality: 90%
- Progressive: Yes
- Optimized: Yes

## Versioning

Renders are versioned by date:
```
renders/
├── 2026-01/
│   ├── product-hero-v1.png
│   └── product-hero-v2.png
└── 2026-02/
    └── product-hero-v3.png
```

## Integration

### Website
```html
<img src="/assets/renders/hero-product.png" 
     alt="Biometrics Product" 
     loading="lazy">
```

### Presentations
```markdown
![Architecture Diagram](../renders/architecture-diagram.png)
```

## Maintenance

- Review renders monthly
- Update for new product versions
- Archive old versions

## See Also

- [3D Models Directory](../3d/)
- [Images Directory](../images/)
- [Logos Directory](../logos/)
