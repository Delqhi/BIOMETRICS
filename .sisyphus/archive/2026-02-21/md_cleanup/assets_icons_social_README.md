# Social Icons

## Overview

This directory contains social media icons used for linking to social profiles, sharing buttons, and marketing materials.

## Contents

### Platform Icons

| File | Platform | Usage |
|------|----------|-------|
| twitter.svg | Twitter/X | Social links |
| facebook.svg | Facebook | Social links |
| instagram.svg | Instagram | Social links |
| linkedin.svg | LinkedIn | Professional network |
| github.svg | GitHub | Code repository |
| youtube.svg | YouTube | Video content |
| discord.svg | Discord | Community |
| slack.svg | Slack | Team communication |

### Share Icons

| File | Description |
|------|-------------|
| share.svg | Generic share |
| share-twitter.svg | Twitter share |
| share-facebook.svg | Facebook share |
| share-linkedin.svg | LinkedIn share |

## Icon Standards

### Size
- **Default**: 24px
- **Small**: 16px (inline)
- **Large**: 32px (buttons)

### Style
- **Fill**: Solid fill
- **Stroke**: None
- **Color**: Brand colors per platform

## Usage

### Links
```html
<a href="https://twitter.com/biometrics" aria-label="Follow us on Twitter">
  <svg class="social-icon">...</svg>
</a>
```

### Share Buttons
```html
<a href="https://twitter.com/intent/tweet?text=..." class="share-twitter">
  <svg class="share-icon">...</svg>
</a>
```

## Brand Colors

| Platform | Color |
|----------|-------|
| Twitter/X | #000000 |
| Facebook | #1877F2 |
| Instagram | #E4405F |
| LinkedIn | #0A66C2 |
| GitHub | #181717 |
| YouTube | #FF0000 |
| Discord | #5865F2 |

## Accessibility

### Focus States
```css
.social-icon:focus {
  outline: 2px solid #2563EB;
  outline-offset: 2px;
}
```

### Screen Readers
```html
<a href="..." aria-label="Follow on Twitter">
  <svg aria-hidden="true">...</svg>
</a>
```

## Maintenance

- Update for rebranding
- Add new platforms as needed
- Remove deprecated services

## See Also

- [Icons Overview](../icons/)
- [Action Icons](./action/)
- [Feature Icons](./feature/)
