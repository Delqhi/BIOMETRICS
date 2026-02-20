# Dashboard Components

## Overview

This directory contains reusable dashboard UI components used in the analytics and monitoring dashboards. These components are designed for visualization and data presentation.

## Components

### Charts

| Component | Description | Library |
|-----------|-------------|---------|
| line-chart.svg | Line chart visualization | Custom |
| bar-chart.svg | Bar chart visualization | Custom |
| pie-chart.svg | Pie chart visualization | Custom |
| area-chart.svg | Area chart visualization | Custom |

### Widgets

| Component | Description |
|-----------|-------------|
| metric-card.svg | Key metric display |
| status-badge.svg | Status indicator |
| progress-ring.svg | Circular progress |
| data-table.svg | Tabular data |

### Navigation

| Component | Description |
|-----------|-------------|
| sidebar.svg | Sidebar navigation |
| topbar.svg | Top navigation bar |
| breadcrumb.svg | Breadcrumb trail |

### Layouts

| Component | Description |
|-----------|-------------|
| grid-layout.svg | Dashboard grid |
| card-layout.svg | Card container |
| split-view.svg | Split panel |

## Component Guidelines

### Dimensions
- Base unit: 8px
- Margins: 16px, 24px, 32px
- Padding: 8px, 12px, 16px

### Typography
- Headings: Inter Bold
- Body: Inter Regular
- Data: JetBrains Mono

### Colors

| Element | Color |
|---------|-------|
| Primary | #2563EB |
| Success | #10B981 |
| Warning | #F59E0B |
| Error | #EF4444 |
| Text | #1E293B |
| Muted | #64748B |

### Spacing
- Component gap: 16px
- Section gap: 24px
- Page margin: 32px

## Usage

### SVG Components
```html
<svg class="chart">
  <use href="#line-chart" />
</svg>
```

### React Components
```jsx
import { LineChart, MetricCard } from './components';

<Dashboard>
  <MetricCard title="Users" value="1,234" />
  <LineChart data={userData} />
</Dashboard>
```

## Responsive Design

### Breakpoints
| Breakpoint | Width | Columns |
|------------|-------|---------|
| Mobile | < 640px | 1 |
| Tablet | 640-1024px | 2 |
| Desktop | > 1024px | 4 |

### Fluid Scaling
```css
.dashboard-card {
  width: calc(25% - 16px);
  min-width: 200px;
}
```

## Accessibility

### Color Contrast
- Minimum 4.5:1 for text
- Minimum 3:1 for UI elements

### ARIA Labels
```html
<svg role="img" aria-label="User growth chart">
  <title>User Growth</title>
</svg>
```

## Animation

### Transitions
- Duration: 200ms
- Easing: ease-out
- Hover: scale(1.02)

### Loading States
```css
.skeleton {
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: loading 1.5s infinite;
}
```

## Maintenance

- Review components quarterly
- Update for design system changes
- Test responsive behavior

## See Also

- [Images Directory](../images/)
- [Icons Directory](../icons/)
- [Diagrams Directory](../diagrams/)
