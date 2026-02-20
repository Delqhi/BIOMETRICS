# Monitoring Template

## Overview

The `monitoring` template provides comprehensive system monitoring and observability for the biometrics system. This template enables real-time metrics collection, alerting, and dashboard visualization.

## Components

### Metrics Collection
- **System Metrics**: CPU, memory, disk, network
- **Application Metrics**: Request latency, throughput, errors
- **Business Metrics**: Authentications, enrollments, verifications

### Alerting
- **Threshold Alerts**: Based on metric thresholds
- **Anomaly Detection**: ML-based anomaly detection
- **Rate Alerts**: Sudden spikes or drops

### Dashboards
- **System Overview**: Infrastructure health
- **Application Performance**: Service metrics
- **Business Analytics**: User activity trends

## Configuration

```yaml
monitoring:
  enabled: true
  interval: 30s
  
  exporters:
    prometheus:
      enabled: true
      port: 9090
    statsd:
      enabled: true
      host: localhost
      port: 8125
  
  alerts:
    email:
      enabled: true
      recipients:
        - ops@example.com
    slack:
      enabled: true
      webhook_url: "https://hooks.slack.com/..."
```

## Metrics Reference

### System Metrics

| Metric | Type | Description |
|--------|------|-------------|
| system_cpu_percent | gauge | CPU usage percentage |
| system_memory_percent | gauge | Memory usage percentage |
| system_disk_percent | gauge | Disk usage percentage |
| system_network_bytes | counter | Network bytes transferred |

### Application Metrics

| Metric | Type | Description |
|--------|------|-------------|
| app_requests_total | counter | Total requests |
| app_request_duration_ms | histogram | Request duration |
| app_errors_total | counter | Total errors |
| app_active_connections | gauge | Active connections |

### Business Metrics

| Metric | Type | Description |
|--------|------|-------------|
| auth_attempts_total | counter | Authentication attempts |
| auth_success_total | counter | Successful authentications |
| auth_failure_total | counter | Failed authentications |
| enrollments_total | counter | New enrollments |

## Alert Rules

### Critical Alerts
```yaml
alerts:
  - name: HighCPU
    condition: cpu_percent > 90
    severity: critical
    duration: 5m
    message: "CPU usage above 90%"
    
  - name: HighErrorRate
    condition: error_rate > 0.05
    severity: critical
    duration: 2m
    message: "Error rate above 5%"
```

### Warning Alerts
```yaml
  - name: HighMemory
    condition: memory_percent > 80
    severity: warning
    duration: 10m
    message: "Memory usage above 80%"
    
  - name: SlowResponse
    condition: p95_latency_ms > 1000
    severity: warning
    duration: 5m
    message: "P95 latency above 1s"
```

## Dashboards

### System Overview Dashboard
- CPU/Memory/Disk gauges
- Network traffic graph
- Active connections
- Process list

### Application Dashboard
- Request rate (requests/second)
- Latency distribution (P50, P95, P99)
- Error rate
- Active users

### Business Dashboard
- Authentication success rate
- Enrollment trends
- Geographic distribution
- Peak usage times

## Prometheus Queries

### CPU Usage
```promql
100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)
```

### Request Rate
```promql
rate(http_requests_total[5m])
```

### Error Rate
```promql
rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])
```

### P95 Latency
```promql
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
```

## Integration

### Prometheus Configuration
```yaml
scrape_configs:
  - job_name: 'biometrics'
    static_configs:
      - targets: ['localhost:9090']
```

### Grafana Integration
Import dashboard JSON from `./dashboards/grafana/`:

```bash
curl -X POST https://grafana.example.com/api/dashboards/import \
  -H "Authorization: Bearer $GRAFANA_TOKEN" \
  -d @dashboards/grafana/biometrics-overview.json
```

### AlertManager Integration
```yaml
route:
  group_by: ['alertname']
  receiver: 'default'
  
receivers:
  - name: 'default'
    email_configs:
      - to: 'alerts@example.com'
```

## Logging

Enable debug logging for troubleshooting:
```yaml
logging:
  level: debug
  format: json
```

## Performance

| Metric | Value |
|--------|-------|
| Metrics Collection | < 100ms |
| Alert Evaluation | < 500ms |
| Dashboard Load | < 1s |
| Storage Retention | 30 days |

## Troubleshooting

### Metrics Not Appearing
1. Check exporter is running
2. Verify network connectivity
3. Check firewall rules

### Alerts Not Firing
1. Verify alert rules are loaded
2. Check alert manager configuration
3. Review notification settings

## Maintenance

### Backup Metrics
```bash
# Export metrics
curl -s http://localhost:9090/api/v1/query_range \
  -G -d 'query=up' \
  -d 'start=2026-01-01T00:00:00Z' \
  -d 'end=2026-01-31T23:59:59Z' \
  > metrics_backup.json
```

### Clean Old Data
```bash
# Clean data older than 30 days
biometrics-cli monitor clean --older-than 30d
```

## See Also

- [CLI Commands](../cmd/README.md)
- [API Documentation](../docs/api.md)
- [Alert Configuration](./alerts/README.md)
