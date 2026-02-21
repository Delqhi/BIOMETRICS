# 📊 BIOMETRICS Agent Status Dashboard

Production-ready real-time monitoring dashboard for BIOMETRICS AI Agent Orchestration Platform.

## 🎯 Features

- **Real-Time Agent Monitoring** - Live status updates for all AI agents
- **Performance Metrics** - Request rates, response times, error rates
- **Interactive Charts** - Chart.js powered analytics
- **Cyberpunk Theme** - Bold, futuristic design with animations
- **Responsive Design** - Works on desktop, tablet, and mobile
- **WebSocket Updates** - Sub-second latency for live data

## 🏗️ Architecture

```
Dashboard UI (HTML/CSS/JS)
    ↓
WebSocket Connection
    ↓
BIOMETRICS API (/api/dashboard/data)
    ↓
Prometheus Metrics + Agent Status
```

## 📁 File Structure

```
assets/dashboard/
├── index.html          # Main dashboard page (255 lines)
├── styles.css          # Cyberpunk theme styling (40+ lines)
├── app.js              # Real-time logic (227 lines)
├── components/         # Reusable UI components
│   ├── agent-card.html
│   ├── metrics-panel.html
│   └── alert-banner.html
└── README.md           # This documentation
```

## 🚀 Quick Start

### 1. Access Dashboard

Open in browser:
```
http://localhost:8080/assets/dashboard/index.html
```

### 2. WebSocket Connection

Dashboard automatically connects to:
```javascript
ws://localhost:8080/ws/dashboard
```

### 3. API Endpoints

Dashboard fetches initial data from:
```
GET /api/dashboard/data
```

Response format:
```json
{
  "metrics": {
    "requestRate": 42,
    "avgResponse": 187,
    "errorRate": 0.12,
    "queueSize": 15
  },
  "agents": [
    {
      "id": "sisyphus",
      "name": "Sisyphus",
      "role": "Main Coder",
      "status": "active",
      "tasksCompleted": 127,
      "avgTime": 2340,
      "errors": 2,
      "progress": 75,
      "currentTask": "Refactoring API endpoints"
    }
  ]
}
```

## 🎨 Design System

### Color Palette

| Color | Variable | Hex | Usage |
|-------|----------|-----|-------|
| Primary Cyan | `--color-primary-cyan` | #00f5ff | Accents, highlights |
| Primary Magenta | `--color-primary-magenta` | #ff00ff | Gradients |
| Success | `--color-success` | #00ff88 | Active states |
| Warning | `--color-warning` | #ffaa00 | Idle states |
| Error | `--color-error` | #ff3366 | Error states |
| Background Dark | `--color-bg-dark` | #0a0a0f | Page background |

### Typography

- **Display Font**: Rajdhani (headings, UI text)
- **Mono Font**: JetBrains Mono (metrics, numbers)

### Animations

- **Pulse**: Live indicator breathing effect
- **Grid Move**: Animated background grid
- **Card Hover**: Elevation on hover
- **Slide Down**: Alert banner entrance

## 🔧 Customization

### Theme Colors

Edit CSS variables in `styles.css`:
```css
:root {
    --color-primary-cyan: #00f5ff;
    --color-primary-magenta: #ff00ff;
    /* ... more variables */
}
```

### Chart Configuration

Modify chart options in `app.js`:
```javascript
this.charts.request = new Chart(ctx, {
    type: 'line',
    data: { ... },
    options: { ... }
});
```

### Agent Card Template

Edit template in `index.html`:
```html
<template id="agent-card-template">
    <!-- Custom agent card structure -->
</template>
```

## 📊 Metrics Displayed

### Header Statistics

- **Active Agents** - Currently working agents
- **Tasks Completed** - Total tasks finished
- **Success Rate** - Percentage of successful tasks

### Metric Cards

1. **Request Rate** - Requests per second (req/s)
2. **Avg Response** - Average response time (ms)
3. **Error Rate** - Percentage of failed requests
4. **Queue Size** - Pending tasks in queue

### Agent Status Cards

Each agent card displays:
- Name and role
- Status badge (active/idle/error)
- Tasks completed count
- Average completion time
- Error count
- Current task progress bar
- Current task name

## 🔌 WebSocket Events

### Incoming Events

```javascript
// Metrics update
{
  "type": "metrics",
  "payload": {
    "requestRate": 42,
    "avgResponse": 187,
    "errorRate": 0.12,
    "queueSize": 15
  }
}

// Agent status update
{
  "type": "agents",
  "payload": [
    { "id": "sisyphus", "status": "active", ... }
  ]
}

// Alert notification
{
  "type": "alert",
  "payload": {
    "message": "Agent offline",
    "severity": "error"
  }
}
```

## 🛠️ Development

### Adding New Metrics

1. Add HTML element with unique ID:
```html
<div class="metric-value" id="new-metric">0</div>
```

2. Update in JavaScript:
```javascript
document.getElementById('new-metric').textContent = value;
```

3. Add to metrics update handler:
```javascript
updateMetrics(metrics) {
    // ... existing updates
    document.getElementById('new-metric').textContent = metrics.newMetric;
}
```

### Adding New Agent Status

1. Define CSS class:
```css
.agent-status-badge.custom {
    background: rgba(0, 100, 255, 0.1);
    color: #0066ff;
    border: 1px solid rgba(0, 100, 255, 0.3);
}
```

2. Use in agent data:
```javascript
{ status: 'custom', ... }
```

### Filter Implementation

Agent filters are ready for implementation:
```javascript
document.querySelectorAll('.filter-btn').forEach(btn => {
    btn.addEventListener('click', (e) => {
        const filter = e.target.dataset.filter;
        // Filter agents by status: all, active, idle, error
    });
});
```

## 📱 Responsive Breakpoints

- **Desktop**: > 1024px (full grid layout)
- **Tablet**: 768px - 1024px (adjusted grid)
- **Mobile**: < 768px (single column layout)

## 🔒 Security Considerations

- WebSocket connection should use authentication tokens
- API endpoints require proper authorization
- No sensitive data displayed in dashboard
- CORS headers must be configured properly

## 🚀 Performance

### Optimizations Applied

- **CSS Variables** - Fast theme updates
- **Chart.js 'none' mode** - No animation on updates
- **Debounced updates** - Prevent excessive re-renders
- **Lazy loading** - Charts initialized on demand

### Target Metrics

- **Initial Load**: < 2 seconds
- **Update Latency**: < 100ms
- **FPS**: 60fps during updates
- **Memory**: < 50MB during operation

## 🧪 Testing

### Manual Testing Checklist

- [ ] All metrics display correctly
- [ ] Agent cards render properly
- [ ] WebSocket reconnects on disconnect
- [ ] Charts update in real-time
- [ ] Responsive design works on all devices
- [ ] Animations are smooth
- [ ] No console errors

### Automated Testing

```javascript
// Example test
describe('DashboardApp', () => {
    it('should initialize charts', () => {
        const app = new DashboardApp();
        expect(app.charts.request).toBeDefined();
    });
    
    it('should update metrics', () => {
        const app = new DashboardApp();
        app.updateMetrics({ requestRate: 100 });
        expect(document.getElementById('request-rate').textContent).toBe('100');
    });
});
```

## 📈 Future Enhancements

- [ ] Dark/Light theme toggle
- [ ] Export metrics to CSV
- [ ] Custom date range selection
- [ ] Agent performance comparison
- [ ] Historical data charts
- [ ] Alert configuration UI
- [ ] Multi-language support
- [ ] PWA offline support

## 🤝 Contributing

1. Follow existing code style
2. Add comments for complex logic
3. Test on multiple screen sizes
4. Update this README for new features

## 📄 License

MIT License - Part of BIOMETRICS Platform

---

**Last Updated:** 2026-02-19  
**Version:** 1.0.0  
**Status:** Production Ready
