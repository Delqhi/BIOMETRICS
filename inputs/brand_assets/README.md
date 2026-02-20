# Brand Assets Directory

## Overview

This directory contains official brand assets, including logos, color palettes, typography guidelines, and visual identity elements for the biometrics brand.

## Contents

### Logos

| Directory | Description |
|-----------|-------------|
| logos/ | Logo files in various formats |
| logos/vector/ | Vector logo files |
| logos/raster/ | Raster logo files |

### Color Palette

| File | Description |
|------|-------------|
| colors.json | Color definitions (JSON) |
| colors.css | CSS custom properties |
| colors.scss | SCSS variables |

### Typography

| File | Description |
|------|-------------|
| fonts/ | Font files |
| typography.css | Typography styles |
| type-scale.md | Type scale guide |

## Logo Files

### Primary Logos
```
logos/
├── primary/
│   ├── biometrics-logo.svg
│   ├── biometrics-logo.png
│   └── biometrics-logo-dark.png
├── icon/
│   ├── biometrics-icon.svg
│   └── biometrics-icon.ico
└── monochrome/
    ├── biometrics-mono-black.svg
    └── biometrics-mono-white.svg
```

## Colors

### Primary Colors
```json
{
  "primary": {
    "blue": "#2563EB",
    "dark": "#1E293B"
  },
  "secondary": {
    "slate": "#64748B"
  },
  "accent": {
    "green": "#10B981"
  }
}
```

### CSS Variables
```css
:root {
  --color-primary: #2563EB;
  --color-primary-dark: #1D4ED8;
  --color-secondary: #1E293B;
  --color-accent: #10B981;
}
```

## Typography

### Font Families
- **Primary**: Inter
- **Monospace**: JetBrains Mono
- **Display**: (for headlines)

### Type Scale
| Level | Size | Weight |
|-------|------|--------|
| H1 | 48px | Bold |
| H2 | 36px | Bold |
| H3 | 24px | SemiBold |
| Body | 16px | Regular |
| Small | 14px | Regular |

## Usage Guidelines

### Logo Usage
- Minimum size: 32px
- Clear space: 1x height
- Don't stretch or distort
- Use approved color variants

### Color Usage
- Use primary blue for CTAs
- Use secondary for text
- Use accent for success states
- Never use unapproved colors

## Templates

### Presentation Template
```
templates/
├── presentation.pptx
├── slide-master.pptx
└── handout.docx
```

### Document Template
```
templates/
├── letterhead.docx
├── invoice.docx
└── contract.docx
```

## Maintenance

- Review annually
- Update for rebranding
- Archive old versions

## Best Practices

1. **Consistency**: Use approved assets only
2. **Quality**: Use vector when possible
3. **Accessibility**: Ensure contrast
4. **Version**: Use current versions

## See Also

- [Inputs Overview](../inputs/)
- [References](./references/)
- [Logos Overview](../assets/logos/)
