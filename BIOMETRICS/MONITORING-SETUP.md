# Monitoring Setup Documentation

**Project:** BIOMETRICS  
**Last Updated:** 2026-02-18  
**Maintainer:** DevOps Team

---

## Overview

This document describes the monitoring infrastructure for BIOMETRICS, including Prometheus metrics collection, Grafana dashboards, and custom instrumentation.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        MONITORING ARCHITECTURE                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐                │
│  │  Application │───▶│  Prometheus   │───▶│   Grafana    │                │
│  │  (Metrics)   │    │  (Storage)    │    │  (Dashboards)│                │
│  └──────────────┘    └──────────────┘    └──────────────┘                │
│         │                    │                    │                         │
│         ▼                    ▼                    ▼                        │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐                │
│  │   Node.js    │    │  Alertmanager│    │    PagerDuty │                │
│  │  (Exporter)  │    │  (Alerts)    │    │  (On-call)   │                │
│  └──────────────┘    └──────────────┘    └──────────────┘                │
│                                                                              │
│  ┌──────────────┐    ┌──────────────┐                                   │
│  │  Exporters   │    │   Loki       │                                   │
│  │  (System)   │    │   (Logs)     │                                   │
│  └──────────────┘    └──────────────┘                                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Prometheus Configuration

### prometheus.yml

Create `monitoring/prometheus/prometheus.yml`:

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    cluster: 'biometrics-production'
    environment: 'production'

alerting:
  alertmanagers:
    - static_configs:
        - targets:
            - alertmanager:9093

rule_files:
  - 'alerts/*.yml'

scrape_configs:
  - job_name: 'biometrics-app'
    static_configs:
      - targets: ['app:3000']
    metrics_path: '/api/metrics'
    scrape_interval: 10s
    honor_labels: true
    relabel_configs:
      - source_labels: [__address__]
        target_label: instance
        regex: '([^:]+):\\d+'
        replacement: '${1}'

  - job_name: 'node-exporter'
    static_configs:
      - targets: ['node-exporter:9100']

  - job_name: 'postgres-exporter'
    static_configs:
      - targets: ['postgres-exporter:9187']

  - job_name: 'redis-exporter'
    static_configs:
      - targets: ['redis-exporter:9121']

  - job_name: 'nginx'
    static_configs:
      - targets: ['nginx-exporter:9113']

  - job_name: 'blackbox'
    metrics_path: '/probe'
    params:
      module: [http_2xx]
    static_configs:
      - targets:
          - https://biometrics.example.com
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: blackbox-exporter:9113
```

### Docker Compose for Monitoring

Create `monitoring/docker-compose.yml`:

```yaml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus:v2.48.0
    container_name: prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--storage.tsdb.retention.time=30d'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--web.enable-lifecycle'
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - ./alerts:/etc/prometheus/alerts:ro
      - prometheus_data:/prometheus
    restart: unless-stopped
    networks:
      - monitoring

  alertmanager:
    image: prom/alertmanager:v0.26.0
    container_name: alertmanager
    command:
      - '--config.file=/etc/alertmanager/alertmanager.yml'
      - '--storage.path=/alertmanager'
    ports:
      - "9093:9093"
    volumes:
      - ./alertmanager.yml:/etc/alertmanager/alertmanager.yml:ro
      - alertmanager_data:/alertmanager
    restart: unless-stopped
    networks:
      - monitoring

  grafana:
    image: grafana/grafana:10.2.2
    container_name: grafana
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_PASSWORD}
      - GF_USERS_ALLOW_SIGN_UP=false
      - GF_SERVER_ROOT_URL=https://grafana.biometrics.example.com
      - GF_SMTP_ENABLED=true
      - GF_SMTP_HOST=smtp.gmail.com:587
      - GF_SMTP_USER=${SMTP_USER}
      - GF_SMTP_PASSWORD=${SMTP_PASSWORD}
    ports:
      - "3001:3000"
    volumes:
      - ./grafana/dashboards:/etc/grafana/provisioning/dashboards
      - ./grafana/datasources:/etc/grafana/provisioning/datasources
      - grafana_data:/var/lib/grafana
    restart: unless-stopped
    networks:
      - monitoring

  node-exporter:
    image: prom/node-exporter:v1.6.1
    container_name: node-exporter
    command:
      - '--path.procfs=/host/proc'
      - '--path.sysfs=/host/sys'
      - '--collector.filesystem.mount-points-exclude=^/(sys|proc|dev|host|etc)($$|/)'
    ports:
      - "9100:9100"
    volumes:
      - /proc:/host/proc:ro
      - /sys:/host/sys:ro
      - /:/rootfs:ro
    restart: unless-stopped
    networks:
      - monitoring

  postgres-exporter:
    image: prometheuscommunity/postgres-exporter:v0.15.0
    container_name: postgres-exporter
    environment:
      - DATA_SOURCE_NAME=postgresql://postgres:postgres@db:5432/biometrics?sslmode=disable
    ports:
      - "9187:9187"
    restart: unless-stopped
    networks:
      - monitoring

  redis-exporter:
    image: oliver006/redis_exporter:v1.54.0
    container_name: redis-exporter
    environment:
      - REDIS_ADDR=redis://redis:6379
    ports:
      - "9121:9121"
    restart: unless-stopped
    networks:
      - monitoring

  loki:
    image: grafana/loki:2.9.2
    container_name: loki
    command: -config.file=/etc/loki/local-config.yaml
    ports:
      - "3100:3100"
    volumes:
      - ./loki:/etc/loki
      - loki_data:/loki
    restart: unless-stopped
    networks:
      - monitoring

  promtail:
    image: grafana/promtail:2.9.2
    container_name: promtail
    volumes:
      - ./promtail.yml:/etc/promtail/config.yml:ro
      - /var/log:/var/log
    restart: unless-stopped
    networks:
      - monitoring

volumes:
  prometheus_data:
  grafana_data:
  alertmanager_data:
  loki_data:

networks:
  monitoring:
    driver: bridge
```

## Application Metrics

### Metrics Library Setup

Create `src/lib/metrics.ts`:

```typescript
import { Registry, Counter, Gauge, Histogram, Summary } from 'prom-client';

export const registry = new Registry();

// Add default metrics
import * as client from 'prom-client';
client.collectDefaultMetrics({ register: registry, prefix: 'biometrics_' });

// Custom metrics
export const httpRequestDuration = new Histogram({
  name: 'biometrics_http_request_duration_seconds',
  help: 'Duration of HTTP requests in seconds',
  labelNames: ['method', 'route', 'status_code'],
  buckets: [0.01, 0.05, 0.1, 0.5, 1, 2, 5],
  registers: [registry],
});

export const httpRequestTotal = new Counter({
  name: 'biometrics_http_requests_total',
  help: 'Total number of HTTP requests',
  labelNames: ['method', 'route', 'status_code'],
  registers: [registry],
});

export const activeUsers = new Gauge({
  name: 'biometrics_active_users',
  help: 'Number of currently active users',
  registers: [registry],
});

export const databaseQueryDuration = new Histogram({
  name: 'biometrics_database_query_duration_seconds',
  help: 'Duration of database queries in seconds',
  labelNames: ['query_type', 'table'],
  buckets: [0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1],
  registers: [registry],
});

export const queueSize = new Gauge({
  name: 'biometrics_queue_size',
  help: 'Number of items in processing queue',
  labelNames: ['queue_name'],
  registers: [registry],
});

export const jobDuration = new Histogram({
  name: 'biometrics_job_duration_seconds',
  help: 'Duration of background jobs in seconds',
  labelNames: ['job_name', 'status'],
  buckets: [1, 5, 10, 30, 60, 300, 600],
  registers: [registry],
});
```

### Express Middleware

Create `src/middleware/metrics.ts`:

```typescript
import { Request, Response, NextFunction } from 'express';
import { httpRequestDuration, httpRequestTotal } from '../lib/metrics';

export function metricsMiddleware(req: Request, res: Response, next: NextFunction) {
  const start = Date.now();

  res.on('finish', () => {
    const duration = (Date.now() - start) / 1000;
    const route = req.route?.path || req.path;

    httpRequestDuration.observe(
      {
        method: req.method,
        route,
        status_code: res.statusCode,
      },
      duration
    );

    httpRequestTotal.inc({
      method: req.method,
      route,
      status_code: res.statusCode,
    });
  });

  next();
}
```

### Metrics Endpoint

Create `src/routes/metrics.ts`:

```typescript
import { Router, Request, Response } from 'express';
import { registry } from '../lib/metrics';

const router = Router();

router.get('/metrics', async (req: Request, res: Response) => {
  try {
    res.set('Content-Type', registry.contentType);
    res.end(await registry.metrics());
  } catch (err) {
    res.status(500).end(err);
  }
});

router.get('/metrics/json', async (req: Request, res: Response) => {
  try {
    res.json(await registry.getMetricsAsJSON());
  } catch (err) {
    res.status(500).json({ error: err });
  }
});

export default router;
```

## Grafana Dashboards

### Dashboard Provisioning

Create `monitoring/grafana/datasources/datasources.yml`:

```yaml
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: true
    jsonData:
      timeInterval: "15s"

  - name: Loki
    type: loki
    access: proxy
    url: http://loki:3100
    editable: true
    jsonData:
      maxLines: 1000
```

### Dashboard Configuration

Create `monitoring/grafana/dashboards/dashboard.yml`:

```yaml
apiVersion: 1

providers:
  - name: 'BIOMETRICS Dashboards'
    orgId: 1
    folder: 'BIOMETRICS'
    type: file
    disableDeletion: false
    updateIntervalSeconds: 10
    allowUiUpdates: true
    options:
      path: /etc/grafana/provisioning/dashboards
```

### Sample Dashboard JSON

Create `monitoring/grafana/dashboards/app-overview.json`:

```json
{
  "dashboard": {
    "title": "BIOMETRICS - Application Overview",
    "tags": ["biometrics", "production"],
    "timezone": "browser",
    "panels": [
      {
        "id": 1,
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(rate(biometrics_http_requests_total[5m])) by (method)",
            "legendFormat": "{{method}}"
          }
        ],
        "gridPos": {"x": 0, "y": 0, "w": 12, "h": 8}
      },
      {
        "id": 2,
        "title": "Response Time (p95)",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, sum(rate(biometrics_http_request_duration_seconds_bucket[5m])) by (le))",
            "legendFormat": "p95"
          }
        ],
        "gridPos": {"x": 12, "y": 0, "w": 12, "h": 8}
      },
      {
        "id": 3,
        "title": "Active Users",
        "type": "stat",
        "targets": [
          {
            "expr": "biometrics_active_users",
            "legendFormat": "Active Users"
          }
        ],
        "gridPos": {"x": 0, "y": 8, "w": 6, "h": 4}
      },
      {
        "id": 4,
        "title": "Error Rate",
        "type": "gauge",
        "targets": [
          {
            "expr": "sum(rate(biometrics_http_requests_total{status_code=~\"5..\"}[5m])) / sum(rate(biometrics_http_requests_total[5m])) * 100",
            "legendFormat": "Error Rate %"
          }
        ],
        "gridPos": {"x": 6, "y": 8, "w": 6, "h": 4}
      },
      {
        "id": 5,
        "title": "Queue Size",
        "type": "graph",
        "targets": [
          {
            "expr": "biometrics_queue_size",
            "legendFormat": "{{queue_name}}"
          }
        ],
        "gridPos": {"x": 12, "y": 8, "w": 12, "h": 8}
      }
    ],
    "refresh": "10s",
    "schemaVersion": 38,
    "version": 1
  }
}
```

## Logging with Loki

### Promtail Configuration

Create `monitoring/promtail.yml`:

```yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://loki:3100/loki/api/v1/push

scrape_configs:
  - job_name: application
    static_configs:
      - targets:
          - localhost
        labels:
          job: biometrics-app
          environment: production
          __path__: /var/log/biometrics/*.log
    pipeline_stages:
      - regex:
          expression: '^(?P<timestamp>\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{3}Z) (?P<level>(INFO|WARN|ERROR|DEBUG)) (?P<message>.*)$'
      - labels:
          level:
      - output:
          source: message
```

## System Metrics Collection

### Node.js Application Metrics

```typescript
// In your main app file
import { registry } from './lib/metrics';
import express from 'express';

const app = express();

// Add metrics middleware
app.use(metricsMiddleware);

// Add metrics endpoint
app.get('/api/metrics', async (req, res) => {
  res.set('Content-Type', registry.contentType);
  res.send(await registry.metrics());
});

// Health check with detailed metrics
app.get('/health', (req, res) => {
  const used = process.memoryUsage();
  const cpu = process.cpuUsage();

  res.json({
    status: 'healthy',
    uptime: process.uptime(),
    memory: {
      heapUsed: Math.round(used.heapUsed / 1024 / 1024) + 'MB',
      heapTotal: Math.round(used.heapTotal / 1024 / 1024) + 'MB',
      rss: Math.round(used.rss / 1024 / 1024) + 'MB',
    },
    cpu: {
      user: cpu.user,
      system: cpu.system,
    },
  });
});
```

## Related Documentation

- [ALERTING-CONFIG.md](./ALERTING-CONFIG.md)
- [CI-CD-PIPELINE.md](./CI-CD-PIPELINE.md)
- [BACKUP-STRATEGY.md](./BACKUP-STRATEGY.md)

---

**End of Monitoring Setup Documentation**
