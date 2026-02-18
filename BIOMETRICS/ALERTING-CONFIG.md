# Alerting Configuration Documentation

**Project:** BIOMETRICS  
**Last Updated:** 2026-02-18  
**Maintainer:** DevOps Team

---

## Overview

This document describes the alerting infrastructure for BIOMETRICS, including Prometheus alert rules, Alertmanager configuration, and notification channels.

## Alert Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         ALERTING ARCHITECTURE                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐                │
│  │  Prometheus  │───▶│ Alertmanager │───▶│   Slack      │                │
│  │   (Rules)    │    │   (Config)   │    │   (Primary)  │                │
│  └──────────────┘    └──────────────┘    └──────────────┘                │
│         │                    │                    │                         │
│         │                    ▼                    │                         │
│         │             ┌──────────────┐            │                         │
│         │             │   PagerDuty  │            │                         │
│         │             │  (Critical)  │            │                         │
│         │             └──────────────┘            │                         │
│         │                    │                    ▼                         │
│         │                    │            ┌──────────────┐                 │
│         │                    └───────────▶│    Email     │                 │
│         │                             │    (Reports)  │                 │
│         │                             └──────────────┘                 │
│         │                                                               │
│         ▼                                                               │
│  ┌──────────────┐                                                      │
│  │   Telegram   │                                                      │
│  │  (On-call)   │                                                      │
│  └──────────────┘                                                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Prometheus Alert Rules

### Application Alerts

Create `monitoring/alerts/app-alerts.yml`:

```yaml
groups:
  - name: application
    interval: 30s
    rules:
      # High error rate
      - alert: HighErrorRate
        expr: |
          sum(rate(biometrics_http_requests_total{status_code=~"5.."}[5m])) 
          / sum(rate(biometrics_http_requests_total[5m])) > 0.05
        for: 5m
        labels:
          severity: critical
          team: backend
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value | humanizePercentage }} for the last 5 minutes"

      # High response time
      - alert: HighResponseTime
        expr: |
          histogram_quantile(0.95, sum(rate(biometrics_http_request_duration_seconds_bucket[5m])) by (le)) > 2
        for: 5m
        labels:
          severity: warning
          team: backend
        annotations:
          summary: "High response time detected"
          description: "p95 response time is {{ $value | humanizeDuration }}"

      # Request rate anomaly
      - alert: RequestRateAnomaly
        expr: |
          sum(rate(biometrics_http_requests_total[10m])) 
          / sum(rate(biometrics_http_requests_total[1h])) < 0.5
        for: 10m
        labels:
          severity: warning
          team: backend
        annotations:
          summary: "Request rate dropped significantly"
          description: "Current request rate is {{ $value | humanizePercentage }} of the hourly average"

      # Active users anomaly
      - alert: NoActiveUsers
        expr: biometrics_active_users == 0
        for: 15m
        labels:
          severity: critical
          team: backend
        annotations:
          summary: "No active users detected"
          description: "No active users for the last 15 minutes - possible service outage"

      # Queue backup
      - alert: QueueBackup
        expr: biometrics_queue_size > 1000
        for: 10m
        labels:
          severity: warning
          team: backend
        annotations:
          summary: "Queue backup detected"
          description: "Queue {{ $labels.queue_name }} has {{ $value }} pending items"
```

### Infrastructure Alerts

Create `monitoring/alerts/infra-alerts.yml`:

```yaml
groups:
  - name: infrastructure
    interval: 30s
    rules:
      # High CPU usage
      - alert: HighCPUUsage
        expr: 100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
        for: 5m
        labels:
          severity: warning
          team: ops
        annotations:
          summary: "High CPU usage on {{ $labels.instance }}"
          description: "CPU usage is {{ $value | humanizePercentage }}"

      - alert: CriticalCPUUsage
        expr: 100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 95
        for: 2m
        labels:
          severity: critical
          team: ops
        annotations:
          summary: "Critical CPU usage on {{ $labels.instance }}"
          description: "CPU usage is {{ $value | humanizePercentage }} - immediate action required"

      # Memory usage
      - alert: HighMemoryUsage
        expr: (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100 > 85
        for: 5m
        labels:
          severity: warning
          team: ops
        annotations:
          summary: "High memory usage on {{ $labels.instance }}"
          description: "Memory usage is {{ $value | humanizePercentage }}"

      # Disk space
      - alert: DiskSpaceLow
        expr: (node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}) * 100 < 10
        for: 10m
        labels:
          severity: warning
          team: ops
        annotations:
          summary: "Low disk space on {{ $labels.instance }}"
          description: "Disk space is {{ $value | humanizePercentage }}"

      - alert: DiskSpaceCritical
        expr: (node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}) * 100 < 5
        for: 5m
        labels:
          severity: critical
          team: ops
        annotations:
          summary: "Critical disk space on {{ $labels.instance }}"
          description: "Disk space is {{ $value | humanizePercentage }} - immediate action required"

      # Network issues
      - alert: NetworkInterfaceDown
        expr: node_network_up{device!="lo"} == 0
        for: 1m
        labels:
          severity: critical
          team: ops
        annotations:
          summary: "Network interface down"
          description: "Network interface {{ $labels.device }} is down on {{ $labels.instance }}"
```

### Database Alerts

Create `monitoring/alerts/database-alerts.yml`:

```yaml
groups:
  - name: database
    interval: 30s
    rules:
      # PostgreSQL down
      - alert: PostgreSQLDown
        expr: pg_up == 0
        for: 1m
        labels:
          severity: critical
          team: dba
        annotations:
          summary: "PostgreSQL is down"
          description: "PostgreSQL instance {{ $labels.instance }} is not responding"

      # High connections
      - alert: PostgreSQLHighConnections
        expr: pg_stat_database_numbackends / pg_settings_max_connections * 100 > 80
        for: 5m
        labels:
          severity: warning
          team: dba
        annotations:
          summary: "High database connections"
          description: "Database is using {{ $value | humanizePercentage }} of max connections"

      # Slow queries
      - alert: PostgreSQLSlowQueries
        expr: rate(pg_stat_database_statements_mean_exec_time[5m]) > 1000
        for: 5m
        labels:
          severity: warning
          team: dba
        annotations:
          summary: "Slow database queries detected"
          description: "Average query execution time is {{ $value | humanizeDuration }}"

      # Replication lag
      - alert: PostgreSQLReplicationLag
        expr: pg_replication_lag > 10
        for: 5m
        labels:
          severity: warning
          team: dba
        annotations:
          summary: "PostgreSQL replication lag"
          description: "Replication is {{ $value | humanizeDuration }} behind"

      # Redis down
      - alert: RedisDown
        expr: redis_up == 0
        for: 1m
        labels:
          severity: critical
          team: ops
        annotations:
          summary: "Redis is down"
          description: "Redis instance {{ $labels.instance }} is not responding"

      # Redis memory
      - alert: RedisHighMemory
        expr: redis_memory_used_bytes / redis_memory_max_bytes * 100 > 80
        for: 5m
        labels:
          severity: warning
          team: ops
        annotations:
          summary: "Redis high memory usage"
          description: "Redis is using {{ $value | humanizePercentage }} of max memory"
```

### Job/Worker Alerts

Create `monitoring/alerts/jobs-alerts.yml`:

```yaml
groups:
  - name: jobs
    interval: 30s
    rules:
      # Job failures
      - alert: JobFailure
        expr: biometrics_job_duration_seconds_count{status="failed"} > 0
        for: 1m
        labels:
          severity: critical
          team: backend
        annotations:
          summary: "Job {{ $labels.job_name }} failed"
          description: "Job {{ $labels.job_name }} has failed"

      # Job timeout
      - alert: JobTimeout
        expr: biometrics_job_duration_seconds_count > 0
        for: 30m
        labels:
          severity: warning
          team: backend
        annotations:
          summary: "Job {{ $labels.job_name }} running too long"
          description: "Job {{ $labels.job_name }} has been running for over 30 minutes"

      # Failed restarts
      - alert: ServiceRestartFailure
        expr: changes(process_restart_count[5m]) > 2
        for: 1m
        labels:
          severity: warning
          team: ops
        annotations:
          summary: "Service restarting frequently"
          description: "Service has restarted {{ $value }} times in the last 5 minutes"
```

## Alertmanager Configuration

### alertmanager.yml

Create `monitoring/alertmanager.yml`:

```yaml
global:
  resolve_timeout: 5m
  smtp_smarthost: smtp.gmail.com:587
  smtp_from: alerts@biometrics.example.com
  smtp_auth_username: alerts@biometrics.example.com
  smtp_auth_password: ${SMTP_PASSWORD}

route:
  group_by: ['alertname', 'cluster', 'service']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  receiver: 'default-receiver'
  routes:
    # Critical alerts - immediate notification
    - match:
        severity: critical
      receiver: 'critical-receiver'
      repeat_interval: 1h
      continue: true

    # Database alerts - DBA team
    - match:
        team: dba
      receiver: 'dba-receiver'
      group_by: ['alertname', 'database']

    # Backend alerts
    - match:
        team: backend
      receiver: 'backend-receiver'

    # Ops alerts
    - match:
        team: ops
      receiver: 'ops-receiver'

# Inhibition rules - suppress alerts when higher-level alert fires
inhibit_rules:
  # Suppress warning if critical alert is firing
  - source_match:
      severity: 'critical'
    target_match:
      severity: 'warning'
    equal: ['alertname', 'cluster', 'service']

  # Suppress infra alerts if app alerts are firing
  - source_match:
      alertname: 'HighErrorRate'
    target_match_re:
      alertname: 'HighCPU|HighMemory'
    equal: ['instance']

receivers:
  - name: 'default-receiver'
    slack_configs:
      - channel: '#alerts'
        send_resolved: true
        color: '{{ if eq .Status "firing" }}danger{{ else }}good{{ end }}'
        title: '{{ range .Alerts }}{{ .Annotations.summary }}{{ end }}'
        text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'
        footer: 'BIOMETRICS Alerts'
        mrkdwn_in: ['text', 'title']

  - name: 'critical-receiver'
    slack_configs:
      - channel: '#alerts-critical'
        send_resolved: true
        color: 'danger'
        title: '🚨 CRITICAL ALERT'
        text: '{{ range .Alerts }}{{ .Annotations.summary }}\n{{ .Annotations.description }}{{ end }}'
    pagerduty_configs:
      - service_key: ${PAGERDUTY_SERVICE_KEY}
        severity: critical
        details:
          firing_alerts: '{{ range .Alerts }}{{ .Annotations.summary }}: {{ .Annotations.description }}{{ end }}'
    email_configs:
      - to: 'on-call@biometrics.example.com'
        send_resolved: true
        headers:
          subject: '🚨 CRITICAL: {{ range .Alerts }}{{ .Annotations.summary }}{{ end }}'

  - name: 'backend-receiver'
    slack_configs:
      - channel: '#alerts-backend'
        send_resolved: true

  - name: 'dba-receiver'
    email_configs:
      - to: 'dba@biometrics.example.com'
        send_resolved: true

  - name: 'ops-receiver'
    slack_configs:
      - channel: '#alerts-ops'
        send_resolved: true

# Template for custom messages
templates:
  - '/etc/alertmanager/template/*.tmpl'
```

## Notification Templates

### Custom Slack Template

Create `monitoring/alertmanager/template/slack.tmpl`:

```tmpl
{{ define "slack.default.title" }}
[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}] {{ .GroupLabels.alertname }}
{{ end }}

{{ define "slack.default.text" }}
{{ range .Alerts }}
*Alert:* {{ .Annotations.summary }}
*Description:* {{ .Annotations.description }}
*Details:*
  {{ range .Labels.SortedPairs }} • *{{ .Name }}:* `{{ .Value }}`
  {{ end }}
{{ end }}
{{ end }}

{{ define "slack.critical.title" }}
🚨 *CRITICAL* 🚨
{{ end }}

{{ define "slack.critical.text" }}
{{ range .Alerts }}
*Impact:* {{ .Annotations.impact | default "Service degradation" }}
*Action Required:* {{ .Annotations.action | default "Investigate immediately" }}
*Runbook:* {{ .Annotations.runbook_url | default "No runbook" }}
{{ end }}
{{ end }}
```

## Alert Routing Examples

### GitHub Integration for Incidents

```yaml
# Alertmanager webhook receiver for GitHub
receivers:
  - name: 'github-incidents'
    webhook_configs:
      - url: 'https://api.github.com/repos/biometrics/incidents/issues'
        send_resolved: true
        headers:
          Authorization: 'token ${GITHUB_TOKEN}'
```

### Custom Webhook

```yaml
receivers:
  - name: 'custom-webhook'
    webhook_configs:
      - url: 'https://biometrics.example.com/api/webhooks/alerts'
        send_resolved: true
```

## Testing Alerts

### Manual Alert Testing

```bash
# Test alert rule syntax
promtool check rules monitoring/alerts/*.yml

# Test alert configuration
promtool check config monitoring/alertmanager.yml

# Reload rules without restart
curl -X POST http://localhost:9090/-/reload

# Reload alertmanager config
curl -X POST http://alertmanager:9093/-/reload
```

### Sending Test Alerts

```bash
# Using amtool (Alertmanager tool)
amtool alert --config.file=alertmanager.yml \
  --alertmanager.url=http://localhost:9093 \
  add --duration=1m \
  --annotation summary="Test alert" \
  --annotation description="This is a test alert" \
  TestAlert
```

## On-Call Schedule

### PagerDuty Schedule

| Day | Primary | Secondary |
|-----|---------|-----------|
| Monday | oncall-1 | oncall-2 |
| Tuesday | oncall-2 | oncall-3 |
| Wednesday | oncall-3 | oncall-1 |
| Thursday | oncall-1 | oncall-2 |
| Friday | oncall-2 | oncall-3 |
| Weekend | oncall-3 | rotation |

## Related Documentation

- [MONITORING-SETUP.md](./MONITORING-SETUP.md)
- [CI-CD-PIPELINE.md](./CI-CD-PIPELINE.md)
- [BACKUP-STRATEGY.md](./BACKUP-STRATEGY.md)

---

**End of Alerting Configuration Documentation**
