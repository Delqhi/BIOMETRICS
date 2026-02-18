# REALTIME-ANALYTICS.md - ClickHouse Analytics

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Data Engineering  
**Author:** BIOMETRICS Analytics Team

---

## 1. Overview

This document describes the real-time analytics architecture using ClickHouse for the BIOMETRICS platform, enabling sub-second query performance on large datasets.

## 2. Architecture

### 2.1 Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| Database | ClickHouse | OLAP database |
| Ingestion | Kafka | Stream data |
| Connector | ClickHouse Kafka | Kafka to CH |
| Query API | HTTP | Query interface |
| Visualization | Metabase | Dashboards |

### 2.2 Architecture Diagram

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Biometric  │────►│   Kafka     │────►│ ClickHouse  │
│  Sensors    │     │   Cluster   │     │  Cluster    │
└─────────────┘     └─────────────┘     └──────┬──────┘
                                                │
                   ┌─────────────┐              │
                   │  Metabase   │◄─────────────┘
                   │  Dashboards │
                   └─────────────┘
```

## 3. ClickHouse Setup

### 3.1 Docker Compose

```yaml
version: '3.8'

services:
  clickhouse:
    image: clickhouse/clickhouse-server:23.8
    container_name: biometrics-clickhouse
    ports:
      - "8123:8123"
      - "9000:9000"
    environment:
      CLICKHOUSE_DB: biometrics
      CLICKHOUSE_USER: analytics
      CLICKHOUSE_PASSWORD: ${CLICKHOUSE_PASSWORD}
    volumes:
      - ./data:/var/lib/clickhouse
      - ./config:/etc/clickhouse-server/config.d
    ulimits:
      nofile:
        soft: 262144
        hard: 262144

  zookeeper:
    image: zookeeper:3.9
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
```

### 3.2 Configuration

```xml
<!-- config/config.xml -->
<?xml version="1.0"?>
<clickhouse>
    <logger>
        <level>information</level>
        <log>/var/log/clickhouse-server/clickhouse-server.log</log>
        <errorlog>/var/log/clickhouse-server/clickhouse-server.err.log</errorlog>
    </logger>
    
    <listen_host>::</listen_host>
    <http_port>8123</http_port>
    <tcp_port>9000</tcp_port>
    
    <max_concurrent_queries>100</max_concurrent_queries>
    <max_server_memory_usage>0.8</max_server_memory_usage>
    
    <query_log>
        <database>system</database>
        <table>query_log</table>
        <flush_interval_milliseconds>7500</flush_interval_milliseconds>
    </query_log>
    
    <metric_log>
        <database>system</database>
        <table>metric_log</table>
        <flush_interval_milliseconds>7500</flush_interval_milliseconds>
    </metric_log>
</clickhouse>
```

## 4. Data Model

### 4.1 Tables

```sql
-- Biometric Events (MergeTree)
CREATE TABLE biometric_events (
    event_id UUID,
    user_id UUID,
    event_type Enum8(
        'heart_rate' = 1,
        'blood_pressure' = 2,
        'temperature' = 3,
        'steps' = 4,
        'sleep' = 5,
        'workout' = 6
    ),
    value Float32,
    unit String,
    timestamp DateTime,
    device_id String,
    location Nullable(String),
    metadata JSON
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (user_id, timestamp, event_type)
TTL timestamp + INTERVAL 1 YEAR
SETTINGS index_granularity = 8192;

-- Aggregated Hourly (SummingMergeTree)
CREATE TABLE hourly_metrics (
    user_id UUID,
    event_type Enum8('heart_rate' = 1, 'blood_pressure' = 2, 'temperature' = 3, 'steps' = 4),
    hour DateTime,
    count UInt64,
    sum_value Float64,
    min_value Float32,
    max_value Float32,
    avg_value Float64
) ENGINE = SummingMergeTree()
PARTITION BY toYYYYMM(hour)
ORDER BY (user_id, event_type, hour);

-- User Activity ( ReplacingMergeTree)
CREATE TABLE user_activity (
    user_id UUID,
    session_id UUID,
    page String,
    action String,
    timestamp DateTime,
    duration UInt32,
    device String,
    browser String
) ENGINE = ReplacingMergeTree()
ORDER BY (user_id, timestamp)
SETTINGS version_column = 'timestamp';
```

### 4.2 Materialized Views

```sql
-- Auto-aggregation for hourly metrics
CREATE MATERIALIZED VIEW hourly_metrics_mv
TO hourly_metrics
AS SELECT
    user_id,
    event_type,
    toStartOfHour(timestamp) as hour,
    count() as count,
    sum(value) as sum_value,
    min(value) as min_value,
    max(value) as max_value,
    avg(value) as avg_value
FROM biometric_events
GROUP BY user_id, event_type, hour;

-- User session summary
CREATE MATERIALIZED VIEW session_summary_mv
TO session
AS_metrics SELECT
    user_id,
    session_id,
    min(timestamp) as session_start,
    max(timestamp) as session_end,
    count() as page_views,
    uniqExact(page) as unique_pages,
    sum(duration) as total_duration
FROM user_activity
GROUP BY user_id, session_id;
```

## 5. Data Ingestion

### 5.1 Kafka Connect

```json
{
  "name": "clickhouse-sink-biometric",
  "config": {
    "connector.class": "com.clickhouse.kafka.connect.ClickHouseSinkConnector",
    "topics": "biometric-events,user-activity",
    "clickhouse.url": "jdbc:clickhouse://clickhouse:8123",
    "clickhouse.database": "biometrics",
    "clickhouse.username": "analytics",
    "clickhouse.password": "${secrets:clickhouse-password}",
    "tasks.max": "4",
    "batch.size": "1000",
    "flush.interval.ms": "5000",
    "errors.tolerance": "all"
  }
}
```

### 5.2 Direct Insert

```python
from clickhouse_driver import Client

class BiometricInserter:
    """Insert biometric data into ClickHouse"""
    
    def __init__(self):
        self.client = Client(
            host='clickhouse',
            port=9000,
            database='biometrics',
            user='analytics',
            password='password'
        )
    
    def insert_event(self, event: BiometricEvent):
        """Insert single event"""
        
        query = """
        INSERT INTO biometric_events 
        (event_id, user_id, event_type, value, unit, timestamp, device_id, metadata)
        VALUES
        """
        
        data = [
            (
                str(event.event_id),
                str(event.user_id),
                event.event_type,
                event.value,
                event.unit,
                event.timestamp,
                event.device_id,
                event.metadata
            )
        ]
        
        self.client.execute(query, data)
    
    def batch_insert(self, events: List[BiometricEvent]):
        """Batch insert events"""
        
        query = """
        INSERT INTO biometric_events 
        (event_id, user_id, event_type, value, unit, timestamp, device_id, metadata)
        VALUES
        """
        
        data = [
            (
                str(e.event_id),
                str(e.user_id),
                e.event_type,
                e.value,
                e.unit,
                e.timestamp,
                e.device_id,
                e.metadata
            )
            for e in events
        ]
        
        self.client.execute(query, data)
```

## 6. Queries

### 6.1 Common Queries

```sql
-- Daily heart rate stats per user
SELECT
    user_id,
    toDate(timestamp) as date,
    count() as readings,
    min(value) as min_hr,
    max(value) as max_hr,
    avg(value) as avg_hr,
    quantile(0.5)(value) as median_hr,
    quantile(0.95)(value) as p95_hr
FROM biometric_events
WHERE event_type = 'heart_rate'
    AND timestamp >= now() - INTERVAL 7 DAY
GROUP BY user_id, date
ORDER BY date DESC;

-- User activity funnel
SELECT
    page,
    count() as views,
    uniqExact(user_id) as unique_users,
    avg(duration) as avg_duration
FROM user_activity
WHERE timestamp >= now() - INTERVAL 1 DAY
GROUP BY page
ORDER BY views DESC;

-- Real-time dashboard query
SELECT
    toStartOfMinute(timestamp) as minute,
    event_type,
    count() as events,
    avg(value) as avg_value
FROM biometric_events
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY minute, event_type
ORDER BY minute;
```

### 6.2 Performance Optimization

```sql
-- Use materialized views for common aggregations
-- Instead of:
SELECT toStartOfHour(timestamp), avg(value)
FROM biometric_events
WHERE user_id = 'xxx'
GROUP BY toStartOfHour(timestamp);

-- Use:
SELECT *
FROM hourly_metrics
WHERE user_id = 'xxx'
    AND hour >= now() - INTERVAL 30 DAY;
```

## 7. API

### 7.1 Query API

```python
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Optional

app = FastAPI(title="Analytics API")

class QueryRequest(BaseModel):
    query: str
    parameters: Optional[dict] = None

@app.post("/api/query")
def execute_query(request: QueryRequest):
    """Execute ClickHouse query"""
    
    try:
        result = clickhouse_client.execute(
            request.query,
            request.parameters or {}
        )
        return {"data": result}
    except Exception as e:
        raise HTTPException(status_code=400, detail=str(e))

@app.get("/api/analytics/heart-rate/{user_id}")
def get_heart_rate_stats(user_id: str, days: int = 7):
    """Get heart rate statistics for user"""
    
    query = """
    SELECT
        toDate(timestamp) as date,
        min(value) as min,
        max(value) as max,
        avg(value) as avg,
        quantile(0.5)(value) as median
    FROM biometric_events
    WHERE user_id = {user_id:UUID}
        AND event_type = 'heart_rate'
        AND timestamp >= now() - INTERVAL {days:Int32} DAY
    GROUP BY date
    ORDER BY date DESC
    """
    
    return clickhouse_client.execute(query, {
        'user_id': user_id,
        'days': days
    })
```

## 8. Monitoring

### 8.1 Query Performance

```sql
-- Slow queries
SELECT
    query_duration_ms,
    query,
    read_rows,
    read_bytes,
    memory_usage
FROM system.query_log
WHERE type = 'ExceptionWhileProcessing'
    AND event_time >= now() - INTERVAL 1 HOUR
ORDER BY query_duration_ms DESC
LIMIT 10;
```

### 8.2 Metrics

| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| Query duration | P95 query time | > 5s |
| Memory usage | Server memory | > 80% |
| Merge latency | Background merge | > 30s |
| Insert latency | Insert response | > 1s |

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
