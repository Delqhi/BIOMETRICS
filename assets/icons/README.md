# Icons Directory

## Overview

This directory contains icon assets used throughout the biometrics project. Icons are available in multiple formats and sizes for different use cases.

## Contents

### General Icons

| File | Description | Format |
|------|-------------|--------|
| check.svg | Checkmark icon | SVG |
| close.svg | Close icon | SVG |
| menu.svg | Menu icon | SVG |
| search.svg | Search icon | SVG |
| settings.svg | Settings icon | SVG |
| user.svg | User icon | SVG |
| lock.svg | Lock icon | SVG |
| unlock.svg | Unlock icon | SVG |

### Navigation Icons

| File | Description | Format |
|------|-------------|--------|
| arrow-left.svg | Left arrow | SVG |
| arrow-right.svg | Right arrow | SVG |
| arrow-up.svg | Up arrow | SVG |
| arrow-down.svg | Down arrow | SVG |
| chevron-left.svg | Chevron left | SVG |
| chevron-right.svg | Chevron right | SVG |

### Status Icons

| File | Description | Format |
|------|-------------|--------|
| success.svg | Success status | SVG |
| warning.svg | Warning status | SVG |
| error.svg | Error status | SVG |
| info.svg | Info status | SVG |
| loading.svg | Loading spinner | SVG |

### Feature Icons

| File | Description | Format |
|------|-------------|--------|
| face-id.svg | Face ID feature | SVG |
| fingerprint.svg | Fingerprint | SVG |
| voice.svg | Voice recognition | SVG |
| iris.svg | Iris scan | SVG |

## Icon Guidelines

### Sizes
| Size | Use Case |
|------|----------|
| 16px | Inline text |
| 24px | Toolbar |
| 32px | Buttons |
| 48px | Large UI |
| 64px | Marketing |

### Style
- **Stroke**: 2px consistent
- **Corners**: Rounded (2px radius)
- **Colors**: CurrentColor (CSS controlled)

## Usage

### SVG Inline
```html
<svg class="icon" width="24" height="24" viewBox="0 0 24 24">
  <path d="M12 2L2 7v10l10 5 10-5V7L12 2z"/>
</svg>
```

### CSS Styling
```css
.icon {
  width: 24px;
  height: 24px;
  fill: currentColor;
}

.icon-success {
  color: #10B981;
}

.icon-error {
  color: #EF4444;
}
```

## Icon Fonts

### Generation
```bash
# Generate icon font
icon-font-generator src/*.svg -o fonts/ --name biometrics-icons
```

### Usage
```css
.biometrics-icon::before {
  content: "\e900";
}
```

## Optimization

### SVG Optimization
```bash
svgo icon.svg -o icon-optimized.svg --preset icons
```

## Accessibility

### ARIA Labels
```html
<svg aria-hidden="true" focusable="false">
  <use href="#icon-check"/>
</svg>

<button>
  <svg aria-hidden="true">...</svg>
  <span>Save</span>
</button>
```

## Version Control

```bash
git lfs track "*.svg"
git lfs track "*.png"
```

## Maintenance

- Review icons quarterly
- Update for design system
- Remove unused icons

## See Also

- [Social Icons](./social/)
- [Action Icons](./action/)
- [Feature Icons](./feature/)
- [Logos Directory](../logos/)
