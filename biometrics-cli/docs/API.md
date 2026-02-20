# BIOMETRICS CLI API Reference

**Version:** 2.0.0  
**Base URL:** `http://localhost:8080/api/v1`  
**Protocol:** HTTP/1.1, HTTP/2  
**Format:** JSON  

---

## Table of Contents

1. [Authentication](#authentication)
2. [Health Endpoints](#health-endpoints)
3. [Metrics](#metrics)
4. [Audit API](#audit-api)
5. [Rate Limiting](#rate-limiting)
6. [Error Handling](#error-handling)
7. [SDK Examples](#sdk-examples)

---

## Authentication

### mTLS (Mutual TLS)

All API endpoints require mutual TLS authentication by default.

**Required Headers:**
```
X-Client-Cert: <base64-encoded-client-certificate>
X-Client-Key: <base64-encoded-client-key>
```

**Example:**
```bash
curl --cert client.crt --key client.key \
  https://api.biometrics.local/api/v1/health
```

### OAuth2 Bearer Token

Alternative authentication via OAuth2 bearer tokens.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Example:**
```bash
curl -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIs..." \
  https://api.biometrics.local/api/v1/audit/events
```

---

## Health Endpoints

### GET /health

Basic health check endpoint. Returns service status.

**Request:**
```http
GET /api/v1/health HTTP/1.1
Host: api.biometrics.local
Accept: application/json
```

**Response:**
```json
{
  "status": "healthy",
  "version": "2.0.0",
  "uptime": "24h30m15s",
  "timestamp": "2026-02-20T15:30:00Z"
}
```

**Status Codes:**
- `200 OK` - Service is healthy
- `503 Service Unavailable` - Service is unhealthy

---

### GET /ready

Readiness probe for Kubernetes. Checks all dependencies.

**Dependencies Checked:**
- Database connection
- Redis connection
- Storage availability

**Request:**
```http
GET /api/v1/ready HTTP/1.1
Host: api.biometrics.local
Accept: application/json
```

**Response:**
```json
{
  "status": "ready",
  "checks": {
    "database": {
      "status": "healthy",
      "latency_ms": 5
    },
    "redis": {
      "status": "healthy",
      "latency_ms": 2
    },
    "storage": {
      "status": "healthy",
      "available_bytes": 107374182400
    }
  }
}
```

**Status Codes:**
- `200 OK` - All dependencies ready
- `503 Service Unavailable` - One or more dependencies unavailable

---

### GET /live

Liveness probe for Kubernetes. Returns 200 if process is alive.

**Request:**
```http
GET /api/v1/live HTTP/1.1
Host: api.biometrics.local
```

**Response:**
```json
{
  "status": "alive",
  "pid": 12345
}
```

**Status Codes:**
- `200 OK` - Process is alive
- `500 Internal Server Error` - Deadlocked or unresponsive

---

## Metrics

### GET /metrics

Prometheus-compatible metrics endpoint.

**Request:**
```http
GET /metrics HTTP/1.1
Host: api.biometrics.local
Accept: text/plain
```

**Response:**
```
# HELP biometrics_http_requests_total Total number of HTTP requests
# TYPE biometrics_http_requests_total counter
biometrics_http_requests_total{method="GET",path="/api/v1/health",status="200"} 1523

# HELP biometrics_http_request_duration_seconds HTTP request duration
# TYPE biometrics_http_request_duration_seconds histogram
biometrics_http_request_duration_seconds_bucket{method="GET",le="0.1"} 1450
biometrics_http_request_duration_seconds_bucket{method="GET",le="0.5"} 1500

# HELP biometrics_audit_events_total Total audit events logged
# TYPE biometrics_audit_events_total counter
biometrics_audit_events_total{type="auth_success"} 5420
biometrics_audit_events_total{type="auth_failure"} 23

# HELP biometrics_rate_limit_requests_total Rate limited requests
# TYPE biometrics_rate_limit_requests_total counter
biometrics_rate_limit_requests_total{action="allowed"} 9823
biometrics_rate_limit_requests_total{action="rejected"} 45
```

---

## Audit API

### GET /audit/events

Query audit events with filters.

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `start_time` | RFC3339 | Filter events after this time |
| `end_time` | RFC3339 | Filter events before this time |
| `event_type` | string | Filter by event type |
| `actor` | string | Filter by actor ID |
| `limit` | int | Maximum results (default: 100, max: 1000) |
| `offset` | int | Pagination offset |

**Request:**
```http
GET /api/v1/audit/events?start_time=2026-02-19T00:00:00Z&limit=50 HTTP/1.1
Host: api.biometrics.local
Authorization: Bearer <token>
Accept: application/json
```

**Response:**
```json
{
  "events": [
    {
      "id": 12345,
      "timestamp": "2026-02-20T14:30:00Z",
      "event_type": "auth_success",
      "actor": "user-123",
      "action": "login",
      "resource": "auth-system",
      "metadata": {
        "ip": "192.168.1.100",
        "method": "oauth2",
        "session": "sess-abc123"
      }
    }
  ],
  "total_count": 1523,
  "has_more": true
}
```

**Status Codes:**
- `200 OK` - Events retrieved successfully
- `400 Bad Request` - Invalid query parameters
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions

---

### POST /audit/events

Log a new audit event.

**Request Body:**
```json
{
  "event_type": "data_access",
  "actor": "user-456",
  "action": "read",
  "resource": "medical-records/rec-001",
  "metadata": {
    "fields": ["name", "dob", "diagnosis"],
    "purpose": "treatment"
  }
}
```

**Response:**
```json
{
  "id": 12346,
  "timestamp": "2026-02-20T15:00:00Z",
  "hash": "sha256:abc123..."
}
```

**Status Codes:**
- `201 Created` - Event logged successfully
- `400 Bad Request` - Invalid event data
- `401 Unauthorized` - Missing authentication

---

### GET /audit/export

Export audit events to a file.

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `start_time` | RFC3339 | Export events after this time |
| `end_time` | RFC3339 | Export events before this time |
| `format` | string | Export format: `json`, `csv`, `xml` |

**Request:**
```http
GET /api/v1/audit/export?format=json&start_time=2026-02-01T00:00:00Z HTTP/1.1
Host: api.biometrics.local
Authorization: Bearer <token>
Accept: application/json
```

**Response:**
```json
{
  "download_url": "/api/v1/audit/export/download/abc123",
  "expires_at": "2026-02-20T16:00:00Z",
  "size_bytes": 1048576
}
```

---

### GET /audit/stats

Get audit statistics.

**Response:**
```json
{
  "total_events": 15234,
  "events_by_type": {
    "auth_success": 5420,
    "auth_failure": 23,
    "data_access": 8912,
    "data_modify": 456,
    "security_alert": 5
  },
  "events_by_actor": {
    "user-123": 1523,
    "user-456": 892,
    "system": 12819
  },
  "storage_size_bytes": 52428800,
  "avg_events_per_day": 507.8
}
```

---

## Rate Limiting

### Rate Limit Headers

All responses include rate limit headers:

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1708444800
X-RateLimit-RetryAfter: 0
```

### Rate Limit Exceeded

When rate limited, the API returns:

```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
X-RateLimit-Remaining: 0
X-RateLimit-RetryAfter: 60

{
  "error": "rate_limit_exceeded",
  "message": "Rate limit exceeded. Retry after 60 seconds.",
  "retry_after": 60
}
```

---

## Error Handling

### Error Response Format

All errors follow this format:

```json
{
  "error": "error_code",
  "message": "Human-readable message",
  "details": {
    "field": "validation error details"
  },
  "request_id": "req-abc123",
  "timestamp": "2026-02-20T15:00:00Z"
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `invalid_request` | 400 | Malformed request |
| `unauthorized` | 401 | Missing authentication |
| `forbidden` | 403 | Insufficient permissions |
| `not_found` | 404 | Resource not found |
| `conflict` | 409 | Resource conflict |
| `rate_limit_exceeded` | 429 | Too many requests |
| `internal_error` | 500 | Server error |

---

## SDK Examples

### Go SDK

```go
package main

import (
    "context"
    "fmt"
    
    "biometrics-cli/pkg/audit"
    "biometrics-cli/pkg/auth"
)

func main() {
    ctx := context.Background()
    
    auditor, _ := audit.NewAuditor(nil)
    defer auditor.Stop()
    
    auditor.LogAuthentication("user-123", true, "oauth2", "192.168.1.1")
    
    events, _ := auditor.Query(ctx, &audit.AuditQuery{
        Limit:  10,
        Actors: []string{"user-123"},
    })
    
    fmt.Printf("Found %d events\n", len(events.Events))
}
```

### cURL Examples

```bash
# Health check
curl https://api.biometrics.local/api/v1/health

# Query audit events
curl -H "Authorization: Bearer $TOKEN" \
  "https://api.biometrics.local/api/v1/audit/events?limit=10"

# Export audit logs
curl -H "Authorization: Bearer $TOKEN" \
  "https://api.biometrics.local/api/v1/audit/export?format=json" \
  -o audit-export.json
```

---

## Versioning

The API uses URL-based versioning: `/api/v1/`, `/api/v2/`, etc.

Current version: **v1**

**Version Header:**
```http
API-Version: 2026-02-20
```

**Deprecation:**
Deprecated endpoints include:
```http
Deprecation: true
Sunset: Sat, 01 Mar 2026 00:00:00 GMT
Link: </api/v2/events>; rel="successor-version"
```

---

## Rate Limits by Tier

| Tier | Requests/Second | Burst |
|------|-----------------|-------|
| Free | 10 | 20 |
| Pro | 100 | 200 |
| Enterprise | 1000 | 2000 |

---

**Document Version:** 1.0.0  
**Last Updated:** February 2026  
**Contact:** api@biometrics.local
