# Health Check Workflow Template

## Overview

The Health Check workflow template provides comprehensive system health monitoring. This template verifies the operational status of services, dependencies, and infrastructure components across your entire stack.

The workflow checks various health indicators including service availability, dependency connectivity, resource usage, and response times. It generates detailed health reports and can trigger alerts for degraded states, making it essential for maintaining operational excellence.

This template is essential for organizations seeking to:
- Monitor service availability 24/7
- Verify system dependencies
- Detect degradation early
- Maintain SLA compliance
- Support incident response
- Enable automated failover

## Purpose

The primary purpose of the Health Check template is to:

1. **Verify Services** - Check if all services are running
2. **Test Dependencies** - Validate external dependencies
3. **Measure Performance** - Check response times
4. **Alert on Issues** - Notify on degraded states
5. **Generate Reports** - Document system health status

### Key Use Cases

- **Pre-deployment Check** - Verify health before deployments
- **Scheduled Monitoring** - Regular health verification
- **Incident Response** - Check health during incidents
- **SLA Validation** - Verify SLA compliance
- **Load Balancer Health** - Backend health for routing
- **Database Health** - Database connectivity checks

## Input Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `targets` | array | Yes | - | Services/endpoints to check |
| `check_types` | array | No | all | Types: availability, dependency, performance |
| `alert_on_failure` | boolean | No | true | Send alerts on failures |
| `timeout` | number | No | 30 | Timeout per check in seconds |

### Input Examples

```yaml
# Example 1: Check all services
inputs:
  targets:
    - name: api
      url: https://api.example.com/health
    - name: database
      type: postgresql
      connection: postgresql://localhost:5432/db
    - name: redis
      type: redis
      connection: redis://localhost:6379

# Example 2: Performance-focused checks
inputs:
  targets:
    - name: api
      url: https://api.example.com
  check_types:
    - performance
    - availability

# Example 3: Database health checks
inputs:
  targets:
    - name: primary-db
      type: postgresql
      connection: postgresql://user:pass@db:5432/main
    - name: replica-db
      type: postgresql
      connection: postgresql://user:pass@db-replica:5432/main
```

## Output Results

The template produces detailed health reports:

| Output | Type | Description |
|--------|------|-------------|
| `overall_status` | string | Overall health status |
| `services` | array | Individual service status |
| `issues` | array | Detected issues |

### Output Report Structure

```json
{
  "health_check": {
    "timestamp": "2026-02-19T10:30:00Z",
    "duration_seconds": 15,
    "overall_status": "healthy"
  },
  "summary": {
    "total": 10,
    "healthy": 9,
    "degraded": 1,
    "unhealthy": 0
  },
  "services": [
    {
      "name": "api",
      "status": "healthy",
      "response_time_ms": 45,
      "checks": {
        "availability": "pass",
        "performance": "pass"
      },
      "details": {
        "status_code": 200,
        "content_valid": true
      }
    },
    {
      "name": "database",
      "status": "degraded",
      "response_time_ms": 2500,
      "checks": {
        "connection": "pass",
        "query_performance": "degraded"
      },
      "issues": [
        {
          "severity": "warning",
          "message": "Slow query detected"
        }
      ]
    }
  ],
  "issues": [
    {
      "severity": "warning",
      "service": "database",
      "message": "Query performance degraded",
      "recommendation": "Check database indexes"
    }
  ]
}
```

## Workflow Steps

### Step 1: Check Availability

**ID:** `check-availability`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Verifies services are running:
- HTTP endpoint checks
- TCP port checks
- Process checks
- Container checks

### Step 2: Check Dependencies

**ID:** `check-dependencies`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Tests external dependencies:
- Database connectivity
- Cache availability
- External API status
- Network connectivity

### Step 3: Measure Performance

**ID:** `measure-performance`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Checks response times:
- API response time
- Query execution time
- Resource usage
- Throughput

### Step 4: Generate Report

**ID:** `generate-report`  
**Type:** agent  
**Timeout:** 5 minutes  
**Provider:** opencode-zen

Creates health status report.

### Step 5: Send Alerts

**ID:** `send-alerts`  
**Type:** agent  
**Timeout:** 5 minutes  
**Provider:** opencode-zen

Notifies on failures:
- Slack notifications
- Email alerts
- PagerDuty integration
- Webhook triggers

## Usage Examples

### CLI Usage

```bash
# Basic health check
biometrics workflow run health-check \
  --targets '[{"name": "api", "url": "https://api.example.com/health"}]'

# Multiple services
biometrics workflow run health-check \
  --targets '[
    {"name": "api", "url": "https://api.example.com/health"},
    {"name": "db", "type": "postgresql", "connection": "postgresql://..."}
  ]'

# Performance focus
biometrics workflow run health-check \
  --targets '[{"name": "api", "url": "https://api.example.com"}]' \
  --check_types '["performance"]'
```

### Programmatic Usage

```go
engine := workflows.NewWorkflowEngine("./templates")
template, _ := engine.LoadTemplate("health-check")

instance, _ := engine.CreateInstance(template, map[string]interface{}{
    "targets": []map[string]interface{}{
        {
            "name": "api",
            "url":  "https://api.example.com/health",
        },
        {
            "name":      "database",
            "type":      "postgresql",
            "connection": "postgresql://user:pass@localhost:5432/db",
        },
    },
    "alert_on_failure": true,
})

result, err := engine.Execute(ctx, instance)
```

## Configuration

### Target Configuration

```yaml
options:
  targets:
    api:
      type: http
      url: https://api.example.com/health
      timeout: 10
      expected_status: 200
      
    database:
      type: postgresql
      connection: postgresql://user:pass@host:5432/db
      timeout: 5
      queries:
        - "SELECT 1"
```

### Thresholds

```yaml
options:
  thresholds:
    response_time_ms: 500
    error_rate_percent: 1
    cpu_percent: 80
    memory_percent: 85
    disk_percent: 90
```

### Alerting

```yaml
options:
  alerts:
    slack:
      webhook_url: https://hooks.slack.com/...
      channel: "#alerts"
      
    email:
      smtp: smtp.example.com
      from: health@example.com
      to: oncall@example.com
```

## Troubleshooting

### Common Issues

#### Issue: Service Unreachable

**Symptom:** Connection failed

**Solution:**
- Check network connectivity
- Verify firewall rules
- Confirm service is running

#### Issue: Slow Response

**Symptom:** Response time exceeds threshold

**Solution:**
- Investigate performance bottlenecks
- Check database queries
- Review resource usage

#### Issue: False Positives

**Symptom:** Healthy service marked unhealthy

**Solution:**
- Adjust thresholds
- Increase timeout
- Fix check logic

### Debug Mode

```yaml
options:
  debug: true
  verbose: true
  include_details: true
```

## Best Practices

### 1. Run Regularly

Schedule health checks:
```yaml
trigger:
  type: schedule
  cron: "*/5 * * * *"  # Every 5 minutes
```

### 2. Set Alerts

Configure notifications for failures:
```yaml
options:
  alerts:
    on_failure: true
    on_degraded: true
```

### 3. Track History

Maintain health metrics over time for trend analysis.

### 4. Use Multiple Check Types

Combine availability, dependency, and performance checks.

### 5. Include Dependencies

Always check downstream dependencies.

## Related Templates

- **Monitoring** (`monitoring/`) - Continuous monitoring
- **Deployment** (`deployment/`) - Deployment with health checks
- **Log Analyzer** (`log-analyzer/`) - Log-based health analysis

---

**Template Version:** 1.0.0  
**Author:** BIOMETRICS Team  
**Category:** Operations  
**Tags:** health, monitoring, availability, uptime, incident-response

*Last Updated: February 2026*
