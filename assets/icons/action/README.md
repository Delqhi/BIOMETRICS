# Action Icons

## Overview

This directory contains action icons used for user interface controls. These icons represent common actions like edit, delete, save, and navigate.

## Contents

### CRUD Actions

| File | Action | Usage |
|------|--------|-------|
| create.svg | Create new | Add buttons |
| edit.svg | Edit | Edit buttons |
| delete.svg | Delete | Remove buttons |
| view.svg | View | View buttons |
| copy.svg | Copy | Duplicate |
| download.svg | Download | Export |
| upload.svg | Upload | Import |

### Navigation Actions

| File | Action | Usage |
|------|--------|-------|
| home.svg | Home | Home navigation |
| back.svg | Back | Go back |
| forward.svg | Forward | Go forward |
| refresh.svg | Refresh | Reload |
| external.svg | External | Open link |

### Utility Actions

| File | Action | Usage |
|------|--------|-------|
| settings.svg | Settings | Configuration |
| filter.svg | Filter | Filter data |
| sort.svg | Sort | Sort data |
| search.svg | Search | Search |
| zoom-in.svg | Zoom in | Magnify |
| zoom-out.svg | Zoom out | Reduce |

## Icon Design

### Style
- **Type**: Line icons
- **Stroke**: 2px
- **Corners**: Rounded
- **Size**: 24px default

### States

#### Default
```css
.icon-action {
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
}
```

#### Hover
```css
.icon-action:hover {
  stroke: #2563EB;
}
```

#### Active
```css
.icon-action:active {
  stroke: #1D4ED8;
  transform: scale(0.95);
}
```

#### Disabled
```css
.icon-action:disabled {
  stroke: #CBD5E1;
  cursor: not-allowed;
}
```

## Usage

### Buttons
```html
<button class="btn-action">
  <svg class="icon">...</svg>
  <span>Edit</span>
</button>
```

### Toolbar
```html
<div class="toolbar">
  <button aria-label="Create">
    <svg>...</svg>
  </button>
  <button aria-label="Edit">
    <svg>...</svg>
  </button>
  <button aria-label="Delete">
    <svg>...</svg>
  </button>
</div>
```

## Accessibility

### Keyboard Navigation
```html
<button class="action-btn" aria-label="Edit user">
  <svg aria-hidden="true">...</svg>
</button>
```

### Focus Indicators
```css
.action-btn:focus-visible {
  outline: 2px solid #2563EB;
  outline-offset: 2px;
}
```

## Best Practices

1. Use consistent sizing
2. Maintain visual hierarchy
3. Provide text labels
4. Support keyboard navigation

## See Also

- [Icons Overview](../icons/)
- [Social Icons](./social/)
- [Feature Icons](./feature/)
