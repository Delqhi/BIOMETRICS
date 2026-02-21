# Log Analyzer Workflow Template

## Overview

The Log Analyzer workflow template provides intelligent log analysis and pattern detection. This template processes application logs to identify errors, anomalies, performance issues, and security threats across your entire infrastructure.

The workflow uses AI-powered analysis to understand log patterns, detect anomalies, and provide actionable insights. It can analyze various log formats and sources, from application logs to system logs, supporting JSON, plain text, syslog, and custom formats.

This template is essential for organizations seeking to:
- Identify application errors quickly
- Detect performance issues in real-time
- Spot security threats proactively
- Understand user behavior patterns
- Maintain overall system health
- Reduce MTTR (Mean Time To Recovery)

## Purpose

The primary purpose of the Log Analyzer template is to:

1. **Detect Errors** - Automatically identify errors and exceptions
2. **Find Anomalies** - Detect unusual patterns in logs
3. **Analyze Performance** - Identify performance bottlenecks
4. **Security Monitoring** - Spot potential security issues
5. **Generate Insights** - Provide actionable recommendations

### Key Use Cases

- **Error Detection** - Find and categorize errors across services
- **Performance Analysis** - Identify slow requests and bottlenecks
- **Security Auditing** - Detect suspicious activities and threats
- **User Behavior** - Understand usage patterns and trends
- **Root Cause Analysis** - Investigate production incidents
- **Compliance Logging** - Ensure audit trail requirements

## Input Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `log_source` | string | Yes | - | Path or URL to logs |
| `log_format` | string | No | auto | Log format: json, plain, syslog |
| `time_range` | string | No | 24h | Time range to analyze |
| `analysis_type` | array | No | all | Types: errors, performance, security, patterns |
| `output_format` | string | No | json | Output format: json, html, markdown |

### Input Examples

```yaml
# Example 1: Analyze all log types
inputs:
  log_source: /var/log/app.log
  time_range: 24h

# Example 2: Security-focused analysis
inputs:
  log_source: s3://logs/security/
  log_format: json
  analysis_type:
    - security
    - anomalies

# Example 3: Error analysis with time range
inputs:
  log_source: https://logs.example.com/app
  time_range: 7d
  analysis_type:
    - errors

# Example 4: Performance analysis
inputs:
  log_source: /var/log/nginx/access.log
  log_format: plain
  time_range: 1h
  analysis_type:
    - performance
```

## Output Results

The template produces comprehensive analysis reports:

| Output | Type | Description |
|--------|------|-------------|
| `report` | object | Detailed analysis report |
| `errors` | array | List of identified errors |
| `anomalies` | array | Detected anomalies |
| `recommendations` | array | Actionable recommendations |

### Output Report Structure

```json
{
  "analysis": {
    "timestamp": "2026-02-19T10:30:00Z",
    "log_source": "/var/log/app.log",
    "log_entries": 15420,
    "time_range": "24h",
    "duration_seconds": 45
  },
  "findings": {
    "errors": 45,
    "warnings": 123,
    "anomalies": 5,
    "security_events": 2,
    "performance_issues": 8
  },
  "top_errors": [
    {
      "count": 15,
      "message": "Connection timeout to database",
      "locations": ["api/service.go:145", "worker/processor.go:89"],
      "first_seen": "2026-02-19T08:30:00Z",
      "last_seen": "2026-02-19T09:45:00Z"
    },
    {
      "count": 12,
      "message": "NullPointerException in user service",
      "locations": ["src/services/UserService.java:234"],
      "severity": "high"
    }
  ],
  "anomalies": [
    {
      "type": "unusual_traffic",
      "description": "Traffic spike detected from IP range",
      "confidence": 0.92,
      "details": "5x normal traffic from 192.168.1.0/24"
    }
  ],
  "security_events": [
    {
      "severity": "critical",
      "event": "Multiple failed login attempts",
      "source_ip": "10.0.0.55",
      "attempts": 47,
      "time_window": "10 minutes"
    }
  ],
  "recommendations": [
    {
      "priority": "high",
      "issue": "Database connection pool exhausted",
      "action": "Increase pool size in config",
      "estimated_effort": "15 minutes"
    }
  ],
  "metrics": {
    "requests_per_minute": 1250,
    "error_rate_percent": 0.8,
    "avg_response_time_ms": 145,
    "p99_response_time_ms": 890
  }
}
```

## Workflow Steps

### Step 1: Collect Logs

**ID:** `collect-logs`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Gathers logs from specified source:
- Reads log files
- Fetches from remote sources
- Handles large files with streaming
- Supports multiple formats

### Step 2: Parse Logs

**ID:** `parse-logs`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Parses and structures log entries:
- Detects log format
- Extracts fields
- Normalizes timestamps
- Handles parsing errors

### Step 3: Analyze Errors

**ID:** `analyze-errors`  
**Type:** agent  
**Timeout:** 15 minutes  
**Provider:** opencode-zen

Identifies and categorizes errors:
- Groups similar errors
- Identifies root causes
- Counts occurrences
- Determines severity

### Step 4: Detect Anomalies

**ID:** `detect-anomalies`  
**Type:** agent  
**Timeout:** 15 minutes  
**Provider:** opencode-zen

Finds unusual patterns:
- Traffic anomalies
- Error rate spikes
- Response time changes
- Unusual user behavior

### Step 5: Security Analysis

**ID:** `security-analysis`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Evaluates security events:
- Failed login attempts
- Suspicious IP addresses
- Access pattern anomalies
- Potential breaches

### Step 6: Generate Report

**ID:** `generate-report`  
**Type:** agent  
**Timeout:** 5 minutes  
**Provider:** opencode-zen

Creates comprehensive analysis report.

## Usage Examples

### CLI Usage

```bash
# Basic log analysis
biometrics workflow run log-analyzer \
  --log_source /var/log/app.log \
  --time_range 24h

# Security-focused analysis
biometrics workflow run log-analyzer \
  --log_source s3://logs/ \
  --analysis_type '["security", "errors"]'

# Performance analysis
biometrics workflow run log-analyzer \
  --log_source /var/log/nginx/access.log \
  --time_range 1h \
  --analysis_type '["performance"]'
```

### Programmatic Usage

```go
engine := workflows.NewWorkflowEngine("./templates")
template, _ := engine.LoadTemplate("log-analyzer")

instance, _ := engine.CreateInstance(template, map[string]interface{}{
    "log_source":   "/var/log/app.log",
    "log_format":   "json",
    "time_range":   "24h",
    "analysis_type": []string{"errors", "performance", "security"},
})

result, err := engine.Execute(ctx, instance)
```

## Configuration

### Filters

```yaml
options:
  filters:
    - level: error, warning  # Only analyze errors/warnings
    - exclude_health_checks: true
    - exclude_readiness_checks: true
    - min_error_count: 1
```

### Thresholds

```yaml
options:
  thresholds:
    error_rate_percent: 1.0
    response_time_p99_ms: 1000
    anomaly_confidence: 0.8
```

### Alerting

```yaml
options:
  alerts:
    - type: slack
      channel: "#logs-alerts"
      on:
        - critical_errors
        - security_events
    - type: email
      recipients: ["oncall@company.com"]
      on:
        - errors
```

## Troubleshooting

### Common Issues

#### Issue: Large Log Files

**Symptom:** Memory issues or timeout

**Solution:** Process in smaller time ranges:
```yaml
inputs:
  time_range: 1h  # Instead of 24h
```

#### Issue: Unknown Format

**Symptom:** Parse errors

**Solution:** Specify log_format explicitly:
```yaml
inputs:
  log_format: json
```

#### Issue: Missing Logs

**Symptom:** No log entries found

**Solution:** Verify log_source path and permissions.

### Debug Mode

```yaml
options:
  debug: true
  verbose: true
  include_raw_logs: true
```

## Best Practices

### 1. Run Regularly

Schedule periodic log analysis:
```yaml
trigger:
  type: schedule
  cron: "0 * * * *"  # Hourly
```

### 2. Set Up Alerts

Configure notifications for critical findings:
```yaml
options:
  alerts:
    on_critical_errors: true
    on_security_events: true
```

### 3. Track Trends

Compare analysis over time to identify patterns.

### 4. Correlate Events

Use multiple log sources to correlate events.

### 5. Preserve Raw Logs

Keep raw logs for deeper investigation when needed.

## Related Templates

- **Health Check** (`health-check/`) - System health monitoring
- **Monitoring** (`monitoring/`) - Continuous monitoring setup
- **Security Audit** (`security-audit/`) - Detailed security analysis

---

**Template Version:** 1.0.0  
**Author:** BIOMETRICS Team  
**Category:** Operations  
**Tags:** logging, analysis, monitoring, errors, security, debugging

*Last Updated: February 2026*
