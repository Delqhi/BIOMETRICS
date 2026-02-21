# SVG Logos Directory

## Overview

This directory contains vector logo files in SVG format. SVG logos provide infinite scalability without quality loss and are the preferred format for web and print use.

## Contents

### Source Files

| File | Description | Optimization |
|------|-------------|--------------|
| biometrics-base.svg | Base logo (no styles) | Minified |
| biometrics-filled.svg | Filled version | Minified |
| biometrics-outline.svg | Outline version | Minified |

### Component Parts

| File | Description |
|------|-------------|
| biometrics-symbol.svg | Symbol only |
| biometrics-text.svg | Text only |
| biometrics-badge.svg | Badge/mark version |

### Export Presets

| File | Description | Viewbox |
|------|-------------|---------|
| logo-icon.svg | Icon preset | 0 0 32 32 |
| logo-horizontal.svg | Horizontal layout | 0 0 200 50 |
| logo-stacked.svg | Stacked layout | 0 0 100 100 |

## SVG Guidelines

### Optimization
All SVG files are optimized:
- Remove metadata
- Minify paths
- Reduce precision to 2 decimal places
- Remove unused defs

### Validation
```bash
# Validate SVG syntax
xmllint --noout *.svg

# Check for issues
svgo --folder . --disable --pretty
```

## Usage

### Inline SVG
```html
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 200 50">
  <path fill="#2563EB" d="..."/>
</svg>
```

### As Image
```html
<img src="biometrics-logo.svg" alt="Biometrics">
```

### As Background
```css
.logo {
  background-image: url('biometrics-logo.svg');
}
```

## Styling

### CSS Custom Properties
```css
:root {
  --logo-primary: #2563EB;
  --logo-secondary: #1E293B;
}
```

### Color Override
```css
.logo-blue path {
  fill: #2563EB;
}
```

## Responsive Usage

### Small (Mobile)
```html
<img src="logo-icon.svg" srcset="logo-horizontal.svg 200w">
```

### Large (Desktop)
```html
<img src="logo-horizontal.svg" srcset="logo-stacked.svg 400w">
```

## Accessibility

### ARIA Labels
```html
<svg role="img" aria-label="Biometrics Logo">
  <title>Biometrics Logo</title>
  ...
</svg>
```

### Focus
```html
<a href="/">
  <svg tabindex="0">...</svg>
</a>
```

## Browser Support

- Chrome 4+
- Firefox 3+
- Safari 3.1+
- Edge 12+
- IE 9+ (basic)

## Tools

### Creation
- **Illustrator**: Primary tool
- **Figma**: Alternative design

### Optimization
```bash
# Optimize all SVGs
npx svgo -f . --pretty

# Specific optimization
svgo input.svg -o output.svg --preset default
```

## Maintenance

- Source files in design tools
- Export optimized versions here
- Version in filename if needed

## See Also

- [Logos Overview](../logos/)
- [Icons Directory](../icons/)
- [Images Directory](../images/)
