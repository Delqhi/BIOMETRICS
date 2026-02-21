# Data Export Workflow Template

## Overview

The Data Export workflow template provides a robust, automated solution for extracting data from various sources and transforming it into multiple output formats. This template is designed to handle complex data migration scenarios, scheduled data backups, and real-time data synchronization requirements.

The workflow supports multiple data sources including databases, APIs, file systems, and cloud storage, with flexible output options spanning CSV, JSON, XML, Parquet, and custom formats. Built on AI-powered transformation logic, it can intelligently map fields, apply transformations, and validate data integrity throughout the export process.

This template is particularly valuable for organizations needing to:
- Extract data for analytics and reporting
- Migrate data between systems
- Create regular backups of critical data
- Share data with external partners
- Archive data for compliance requirements

## Purpose

The primary purpose of the Data Export template is to:

1. **Simplify Data Extraction** - Provide a unified interface for exporting data from diverse sources
2. **Ensure Data Quality** - Validate data integrity before, during, and after export
3. **Support Multiple Formats** - Handle various output formats needed by different consumers
4. **Automate Scheduling** - Enable hands-free scheduled exports for recurring needs
5. **Provide Audit Trails** - Maintain complete logs of all export operations

### Key Use Cases

- **Database Migration** - Export data from legacy systems for migration to new platforms
- **Analytics Preparation** - Extract data for business intelligence and analytics tools
- **Compliance Archiving** - Create auditable data exports for regulatory compliance
- **Partner Data Sharing** - Generate structured exports for external stakeholders
- **Backup and Recovery** - Regular automated backups of critical data stores

## Input Parameters

The Data Export template accepts the following input parameters:

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `source` | string | Yes | - | Data source (database, API, file, etc.) |
| `query` | string | Yes | - | Query or filter for data extraction |
| `format` | string | No | csv | Output format (csv, json, xml, parquet, sql) |
| `destination` | string | Yes | - | Target location (path, bucket, endpoint) |
| `encoding` | string | No | utf-8 | Character encoding for output |
| `compression` | string | No | none | Compression type (none, gzip, zip) |
| `batch_size` | number | No | 1000 | Records per batch for processing |
| `include_headers` | boolean | No | true | Include column headers in CSV/JSON |
| `date_format` | string | No | ISO8601 | Date format for timestamp fields |
| `null_values` | string | No | "" | Value to use for NULL fields |

### Input Examples

```yaml
# Example 1: Basic CSV export
inputs:
  source: postgresql://user:pass@localhost:5432/mydb
  query: SELECT * FROM customers WHERE created_at > '2025-01-01'
  format: csv
  destination: /exports/customers_2025.csv

# Example 2: JSON export with compression
inputs:
  source: mongodb://localhost:27017
  query: orders.find({status: "completed"})
  format: json
  destination: s3://my-bucket/exports/
  compression: gzip
  batch_size: 500

# Example 3: Parquet export for analytics
inputs:
  source: postgresql://analytics-db.internal:5432/analytics
  query: SELECT * FROM events WHERE event_date >= CURRENT_DATE - INTERVAL '30 days'
  format: parquet
  destination: /data-warehouse/events/
  date_format: unix

# Example 4: XML export with custom mapping
inputs:
  source: https://api.example.com/data
  format: xml
  destination: ftp://partner.example.com/incoming/
  encoding: ISO-8859-1
  include_headers: true
```

## Output Results

The template produces comprehensive export results:

| Output | Type | Description |
|--------|------|-------------|
| `file_path` | string | Path to exported file(s) |
| `row_count` | number | Total rows exported |
| `file_size` | number | Size of export file in bytes |
| `duration_seconds` | number | Time taken for export |
| `checksum` | string | MD5/SHA checksum for verification |

### Output Report Structure

```json
{
  "export": {
    "id": "exp_20260219_001",
    "source": "postgresql://localhost:5432/mydb",
    "query": "SELECT * FROM customers",
    "format": "csv",
    "destination": "/exports/customers.csv",
    "timestamp": "2026-02-19T10:30:00Z"
  },
  "statistics": {
    "row_count": 15420,
    "file_size_bytes": 2456789,
    "duration_seconds": 45,
    "records_per_second": 342
  },
  "validation": {
    "null_count": 12,
    "duplicate_count": 0,
    "validation_errors": []
  },
  "checksum": {
    "algorithm": "sha256",
    "value": "a1b2c3d4e5f6..."
  },
  "files": [
    {
      "path": "/exports/customers.csv",
      "rows": 15420,
      "size": 2456789
    }
  ]
}
```

## Workflow Steps

### Step 1: Validate Source Connection

**ID:** `validate-source`  
**Type:** agent  
**Timeout:** 5 minutes  
**Provider:** opencode-zen

Validates connectivity to the data source and verifies credentials.

### Step 2: Export Data

**ID:** `export-data`  
**Type:** agent  
**Timeout:** 30 minutes  
**Provider:** opencode-zen

Executes the data extraction query and streams results to the destination:
- Connects to source database/API
- Executes query with pagination
- Handles large result sets via batching
- Manages connection pooling

### Step 3: Transform Data

**ID:** `transform-data`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Applies data transformations:
- Format conversion
- Field mapping
- Data type conversions
- Encoding adjustments

### Step 4: Validate Export

**ID:** `validate-export`  
**Type:** agent  
**Timeout:** 5 minutes  
**Provider:** opencode-zen

Verifies exported data:
- Row count verification
- Checksum calculation
- Sample validation
- Format verification

### Step 5: Upload to Destination

**ID:** `upload-data`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Transfers the exported file to the final destination:
- Local file storage
- Cloud storage (S3, GCS, Azure Blob)
- FTP/SFTP servers
- API endpoints

## Usage Examples

### CLI Usage

```bash
# Simple CSV export
biometrics workflow run data-export \
  --source "postgresql://user:pass@localhost:5432/mydb" \
  --query "SELECT * FROM customers" \
  --format csv \
  --destination /exports/

# JSON export to S3
biometrics workflow run data-export \
  --source "mongodb://localhost:27017" \
  --query "orders.find({status: 'completed'})" \
  --format json \
  --destination "s3://my-bucket/exports/" \
  --compression gzip

# Incremental export (last 24 hours)
biometrics workflow run data-export \
  --source "postgresql://analytics.internal:5432/db" \
  --query "SELECT * FROM events WHERE created_at > NOW() - INTERVAL '24 hours'" \
  --format parquet \
  --destination "/data-warehouse/"
```

### Programmatic Usage

```go
import "github.com/biometrics/biometrics-cli/pkg/workflows"

engine := workflows.NewWorkflowEngine("./templates")
template, _ := engine.LoadTemplate("data-export")

instance, _ := engine.CreateInstance(template, map[string]interface{}{
    "source":      "postgresql://user:pass@localhost:5432/mydb",
    "query":       "SELECT * FROM customers WHERE region = 'EU'",
    "format":      "csv",
    "destination": "/exports/eu_customers.csv",
    "compression": "gzip",
    "batch_size":  1000,
})

result, err := engine.Execute(context.Background(), instance)
fmt.Printf("Exported %d rows\n", result["row_count"])
```

### API Usage

```bash
curl -X POST http://localhost:8080/api/v1/workflows/run \
  -H "Content-Type: application/json" \
  -d '{
    "template": "data-export",
    "inputs": {
      "source": "postgresql://user:pass@localhost:5432/mydb",
      "query": "SELECT * FROM orders",
      "format": "json",
      "destination": "/exports/orders.json",
      "compression": "gzip"
    }
  }'
```

## Configuration

### Source Configuration

Configure different data sources:

```yaml
# Database sources
source:
  type: postgresql
  connection_string: "postgresql://user:pass@host:5432/db"
  
source:
  type: mongodb
  connection_string: "mongodb://host:27017"
  database: mydb
  collection: events

# API sources
source:
  type: api
  url: https://api.example.com/data
  auth:
    type: bearer
    token: $API_TOKEN
  
# File sources
source:
  type: file
  path: /data/input/*.csv
  format: csv
```

### Destination Configuration

```yaml
# Local storage
destination:
  type: local
  path: /exports/
  
# Cloud storage
destination:
  type: s3
  bucket: my-bucket
  prefix: exports/
  region: us-east-1
  
# Remote server
destination:
  type: ftp
  host: ftp.example.com
  path: /incoming/
  username: $FTP_USER
  password: $FTP_PASS
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `EXPORT_DEFAULT_PATH` | Default output directory | `./exports` |
| `EXPORT_COMPRESSION_LEVEL` | Gzip compression level (1-9) | 6 |
| `EXPORT_BATCH_SIZE` | Default batch size | 1000 |
| `EXPORT_MAX_RETRIES` | Maximum retry attempts | 3 |

## Troubleshooting

### Common Issues

#### Issue: Connection Timeout

**Symptom:** Data source connection fails with timeout

**Solution:**
```yaml
inputs:
  source: "postgresql://user:pass@host:5432/db?connect_timeout=30"
# Increase timeout in connection string
```

#### Issue: Memory Overflow

**Symptom:** Export fails with out-of-memory on large datasets

**Solution:**
```yaml
inputs:
  batch_size: 500  # Reduce batch size
options:
  timeout: 60m    # Increase timeout
```

#### Issue: Invalid Characters

**Symptom:** Export contains garbled characters

**Solution:**
```yaml
inputs:
  encoding: utf-8  # Specify correct encoding
```

#### Issue: Permission Denied

**Symptom:** Cannot write to destination

**Solution:**
```bash
# Check destination permissions
ls -la /exports/
chmod -R 755 /exports/
```

### Debug Mode

Enable detailed logging:

```yaml
options:
  debug: true
  log_level: trace
```

### Getting Help

- Documentation: `/docs/data-export/`
- Slack: `#data-help` channel

## Best Practices

### 1. Use Batch Processing

Always use batching for large datasets to avoid memory issues:
```yaml
inputs:
  batch_size: 1000
```

### 2. Validate Exports

Always verify exported data:
```yaml
steps:
  - id: validate-export
    type: agent
    # validation logic
```

### 3. Compress Large Exports

Reduce storage and transfer time:
```yaml
inputs:
  compression: gzip
```

### 4. Schedule Regular Exports

Automate recurring exports:
```yaml
trigger:
  type: schedule
  cron: "0 2 * * *"  # Daily at 2 AM
```

## Related Templates

- **Backup** (`backup/`) - Full system backup workflows
- **Migration** (`migration/`) - Database migration scripts
- **Data Import** (inverse operation)

---

**Template Version:** 1.0.0  
**Author:** BIOMETRICS Team  
**Category:** Data Operations  
**Tags:** export, data, migration, ETL, backup

*Last Updated: February 2026*
