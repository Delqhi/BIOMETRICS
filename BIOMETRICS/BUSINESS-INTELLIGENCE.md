# BUSINESS-INTELLIGENCE.md - Metabase BI

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - BI/Analytics  
**Author:** BIOMETRICS Analytics Team

---

## 1. Overview

This document describes the Business Intelligence platform using Metabase for the BIOMETRICS platform, enabling self-service analytics and data visualization.

## 2. Architecture

### 2.1 Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| BI Platform | Metabase | Analytics frontend |
| Database | ClickHouse | Analytics DB |
| Data Warehouse | PostgreSQL | Primary storage |
| Embedding | JWT | Secure embedding |

### 2.2 Setup

```yaml
# docker-compose.yml
services:
  metabase:
    image: metabase/metabase:v0.48.0
    container_name: biometrics-metabase
    environment:
      MB_DB_TYPE: postgres
      MB_DB_DBNAME: metabase
      MB_DB_PORT: 5432
      MB_DB_USER: metabase
      MB_DB_PASS: ${MB_DB_PASSWORD}
      MB_DB_HOST: room-03-postgres-master
      MB_ANALYTICS_STRATEGY: analytics
      JAVA_TIMEZONE: Europe/Berlin
    volumes:
      - ./metabase-data:/metabase-data
    ports:
      - "51302:3000"
```

## 3. Data Sources

### 3.1 Connections

| Database | Purpose | Sync |
|----------|---------|------|
| ClickHouse | Analytics queries | Real-time |
| PostgreSQL | User data | Hourly |
| MongoDB | Logs | Daily |

### 3.2 Configuration

```bash
# Add database in Metabase admin
# Go to: Settings → Admin → Databases → Add database

Database type: PostgreSQL
Name: BIOMETRICS Primary
Host: room-03-postgres-master
Port: 5432
Database name: biometrics
Username: metabase_reader
Password: [from vault]
SSL: true
```

## 4. Collections Structure

### 4.1 Collection Hierarchy

```
📁 BIOMETRICS
├── 📁 Executive Dashboards
│   ├── KPI Overview
│   ├── Revenue Dashboard
│   └── User Growth
├── 📁 Product Analytics
│   ├── User Activity
│   ├── Feature Usage
│   └── Funnel Analysis
├── 📁 Biometric Health
│   ├── Health Trends
│   ├── Device Performance
│   └── Alert Analysis
├── 📁 Marketing
│   ├── Campaign Performance
│   ├── Attribution
│   └── Cohort Analysis
└── 📁 Operations
    ├── System Health
    ├── Support Metrics
    └── SLA Dashboard
```

### 4.2 Permissions

| Collection | Viewers | Editors |
|------------|---------|---------|
| Executive Dashboards | All users | Leadership |
| Product Analytics | All users | Product team |
| Biometric Health | Medical team | Medical team |
| Marketing | Marketing | Marketing |
| Operations | Support | Engineering |

## 5. Dashboard Templates

### 5.1 KPI Dashboard

```yaml
# Dashboard configuration
name: KPI Overview
description: Executive KPIs
refresh: 5 minutes

cards:
  - title: Total Users
    type: metric
    query: "SELECT count(*) FROM users"
    display: number
    goal: 100000
  
  - title: Daily Active Users
    type: metric
    query: "SELECT count(*) FROM users WHERE last_active >= now() - INTERVAL '1 day'"
    display: trend
    comparison: yesterday
  
  - title: Revenue (MTD)
    type: metric
    query: "SELECT sum(amount) FROM subscriptions WHERE status = 'active'"
    display: progress
  
  - title: User Growth Trend
    type: line
    query: "SELECT date_trunc('day', created_at), count(*) FROM users GROUP BY 1"
    x_axis: date
    y_axis: users
```

### 5.2 Biometric Health Dashboard

```yaml
name: Biometric Health Overview
description: Health metrics
refresh: 1 minute

cards:
  - title: Heart Rate Distribution
    type: histogram
    query: "SELECT value FROM biometric_events WHERE event_type = 'heart_rate'"
    bins: 50
    color: #FF6B6B
  
  - title: Average Daily Steps
    type: line
    query: "SELECT date, avg(steps) FROM daily_metrics WHERE metric = 'steps' GROUP BY date"
    x_axis: date
    y_axis: steps
  
  - title: Alert Rate by Device
    type: bar
    query: "SELECT device, count(*) / 1000 as alerts FROM alerts GROUP BY device"
    x_axis: device
    y_axis: alerts
```

## 6. Question Templates

### 6.1 Common Questions

| Question | Query | Visualization |
|----------|-------|---------------|
| Daily Active Users | Count users by last_active | Line chart |
| Revenue by Plan | Sum revenue by subscription_type | Pie chart |
| User Retention | Cohort retention matrix | Heatmap |
| Funnel Conversion | Funnel stages | Funnel |
| Geographic Distribution | Users by country | Map |

### 6.2 SQL Templates

```sql
-- Daily Active Users
SELECT 
    date_trunc('day', timestamp) as date,
    count(DISTINCT user_id) as dau
FROM user_activity
WHERE timestamp >= now() - INTERVAL '30 day'
GROUP BY date_trunc('day', timestamp)
ORDER BY date;

-- Cohort Retention
WITH cohorts AS (
    SELECT 
        user_id,
        date_trunc('month', created_at) as cohort_month
    FROM users
)
SELECT 
    cohort_month,
    date_trunc('month', u.created_at) as activity_month,
    count(DISTINCT u.user_id) as users,
    round(count(DISTINCT u.user_id) * 100.0 / c.cohort_size, 2) as retention_pct
FROM users u
JOIN cohorts c ON u.user_id = c.user_id
LEFT JOIN (
    SELECT 
        cohort_month,
        count(*) as cohort_size
    FROM cohorts
    GROUP BY cohort_month
) cs ON c.cohort_month = cs.cohort_month
WHERE u.created_at >= c.cohort_month
GROUP BY cohort_month, cs.cohort_size
ORDER BY cohort_month, activity_month;
```

## 7. Embedding

### 7.1 Secure Embedding

```javascript
// Embed dashboard with JWT
constMetabase.init({
  metabaseUrl: "https://biometrics.metabase.delqhi.com",
  jwtToken: generateJWT({
    resource: { dashboard: 123 },
    params: { user_id: currentUser.id },
    exp: Math.round(Date.now() / 1000) + 3600
  })
});

function generateJWT(payload) {
  const header = { alg: "HS256", typ: "JWT" };
  const secret = process.env.METABASE_SECRET_KEY;
  
  const encoded = [
    base64urlEncode(JSON.stringify(header)),
    base64urlEncode(JSON.stringify(payload)),
  ];
  
  const signature = hmacSHA256(encoded.join("."), secret);
  return encoded.join(".") + "." + signature;
}
```

### 7.2 Embed Options

```javascript
const embeddingOptions = {
  theme: "transparent",
  bordered: true,
  titled: true,
  hide_download: false,
  hide_parameters: false,
  default_filters: true,
  additional_info: true,
  transforms: ["timeseries"],
};
```

## 8. Alerts

### 8.1 Alert Configuration

| Alert | Condition | Channels |
|-------|-----------|----------|
| DAU drop | DAU < 1000 | Slack + Email |
| Error rate | Errors > 5% | PagerDuty |
| Revenue anomaly | Revenue < 80% avg | Email |
| New users spike | New users > 200% | Slack |

### 8.2 Alert Setup

```yaml
# Metabase pulse configuration
name: Production Alerts
schedule: "0 * * * *"  # Every hour

cards:
  - query: "SELECT count(*) FROM errors WHERE created_at > now() - INTERVAL '1 hour'"
    condition: "> 100"
    action: notify

channels:
  - type: slack
    channel: "#alerts"
  - type: email
    recipients: ["oncall@biometrics.com"]
```

## 9. Permissions Model

### 9.1 Groups

| Group | Permissions |
|-------|------------|
| Admin | All access |
| Leadership | View all, edit dashboards |
| Product | View product, create questions |
| Marketing | View marketing, edit own |
| Support | View support dashboards |
| User | View own data only |

### 9.2 Row-Level Security

```sql
-- Enable row-level permissions
ALTER TABLE user_data ENABLE ROW LEVEL SECURITY;

-- Policy for users viewing own data
CREATE POLICY user_own_data ON user_data
FOR SELECT
USING (user_id = current_setting('app.current_user_id')::uuid);
```

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
