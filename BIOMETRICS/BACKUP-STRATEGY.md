# Backup Strategy Documentation

**Project:** BIOMETRICS  
**Last Updated:** 2026-02-18  
**Maintainer:** DevOps Team

---

## Overview

This document describes the backup and disaster recovery strategy for BIOMETRICS, including automated backups to cloud storage, restore procedures, and testing protocols.

## Backup Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        BACKUP ARCHITECTURE                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐            │
│  │  Production  │────▶│   Backup      │────▶│    S3/GCS    │            │
│  │  Database    │     │   Server     │     │   Storage    │            │
│  └──────────────┘     └──────────────┘     └──────────────┘            │
│         │                     │                     │                      │
│         ▼                     ▼                     ▼                      │
│  PostgreSQL              Automated              Cross-region               │
│  (wal-g)                Scripts               Replication                │
│                                                                              │
│  ┌──────────────┐     ┌──────────────┐                                 │
│  │   File       │────▶│   Restic     │                                 │
│  │   Storage    │     │   Backups    │                                 │
│  └──────────────┘     └──────────────┘                                 │
│         │                     │                                            │
│         ▼                     ▼                                            │
│  /uploads               Encrypted                                         │
│  /config                Remote Storage                                     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Database Backups

### PostgreSQL Configuration

Create `backup/postgres/backup.sh`:

```bash
#!/bin/bash

# Configuration
export DATABASE_URL="postgresql://postgres:${POSTGRES_PASSWORD}@localhost:5432/biometrics"
export S3_BUCKET="s3://biometrics-backups/postgres"
export RETENTION_DAYS=30
export BACKUP_NAME="biometrics-db-$(date +%Y%m%d-%H%M%S)"

# WAL-G binary path
WALG=/usr/local/bin/wal-g

# Create backup
echo "[$(date)] Starting PostgreSQL backup..."
${WALG} backup-push ${DATABASE_URL}

if [ $? -eq 0 ]; then
    echo "[$(date)] Backup completed successfully"
    
    # List backups to verify
    ${WALG} backup-list
    
    # Cleanup old backups (older than RETENTION_DAYS)
    ${WALG} delete before FIND_COMPRESSED_BACKUP 2d --confirm
    
    echo "[$(date)] Backup cleanup completed"
else
    echo "[$(date)] Backup failed!"
    exit 1
fi

# Verify latest backup
LATEST=$(${WALG} backup-list --pretty --limit 1 | tail -1 | awk '{print $1}')
echo "[$(date)] Latest backup: ${LATEST}"

# Send backup notification
curl -X POST "https://hooks.example.com/backup" \
  -H "Content-Type: application/json" \
  -d "{\"status\": \"success\", \"backup\": \"${BACKUP_NAME}\", \"timestamp\": \"$(date -Iseconds)\"}"
```

### WAL-G Configuration

Create `backup/postgres/wal-g.env`:

```bash
# S3 Configuration
AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
AWS_REGION=eu-central-1
WALG_S3_PREFIX=s3://biometrics-backups/postgres
WALG_COMPRESSION_METHOD=lz4
WALG_BACKUP_COMPRESSION_METHOD=lz4

# Retention
WALG_DELTA_MAX=7
WALG_RETENTION_FULL_BACKUPS=7
WALG_RETENTION_WAL=7

# Performance
WALG_UPLOAD_CONCURRENCY=4
WALG_UPLOAD_DISKConcurrency=1
WALG_DOWNLOAD_CONCURRENCY=4
```

## Application Backups

### File System Backup Script

Create `backup/files/backup-files.sh`:

```bash
#!/bin/bash

# Configuration
SOURCE_DIRS=("/app/uploads" "/app/config" "/app/.env")
BACKUP_ROOT="/tmp/backups"
S3_BUCKET="s3://biometrics-backups/files"
RETENTION_DAYS=30

# Create timestamp
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
BACKUP_NAME="files-${TIMESTAMP}"
BACKUP_PATH="${BACKUP_ROOT}/${BACKUP_NAME}"

# Create backup directory
mkdir -p ${BACKUP_PATH}

echo "[$(date)] Starting file backup..."

# Backup each directory
for DIR in "${SOURCE_DIRS[@]}"; do
    DIR_NAME=$(basename ${DIR})
    if [ -d "${DIR}" ]; then
        echo "[$(date)] Backing up ${DIR}..."
        tar -czf "${BACKUP_PATH}/${DIR_NAME}.tar.gz" -C $(dirname ${DIR}) ${DIR_NAME}
    else
        echo "[$(date)] Warning: ${DIR} does not exist, skipping"
    fi
done

# Create manifest
cat > ${BACKUP_PATH}/manifest.json <<EOF
{
  "timestamp": "${TIMESTAMP}",
  "hostname": "$(hostname)",
  "directories": $(echo "${SOURCE_DIRS[@]}" | jq -R . | jq -s .),
  "backup_size": $(du -sb ${BACKUP_PATH} | cut -f1)
}
EOF

# Encrypt backup
echo "[$(date)] Encrypting backup..."
tar -czf - -C ${BACKUP_ROOT} ${BACKUP_NAME} | \
    gpg --encrypt --recipient ${GPG_RECIPIENT} | \
    aws s3 cp - ${S3_BUCKET}/${BACKUP_NAME}.tar.gz.gpg

if [ $? -eq 0 ]; then
    echo "[$(date)] File backup completed successfully"
    
    # Cleanup local backup
    rm -rf ${BACKUP_PATH}
    
    # Cleanup old remote backups
    aws s3 ls ${S3_BUCKET} | \
        while read line; do
            fileDate=$(echo $line | awk '{print $1}')
            fileName=$(echo $line | awk '{print $4}')
            if [ $(date -d "$fileDate" +%s) -lt $(date -d "-${RETENTION_DAYS} days" +%s) ]; then
                echo "[$(date)] Deleting old backup: ${fileName}"
                aws s3 rm ${S3_BUCKET}/${fileName}
            fi
        done
else
    echo "[$(date)] File backup failed!"
    exit 1
fi
```

### Restic Backup Configuration

Create `backup/restic/config.sh`:

```bash
#!/bin/bash

# Restic configuration
export RESTIC_PASSWORD="your-encryption-password"
export RESTIC_REPOSITORY="s3:https://s3.eu-central-1.amazonaws.com/biometrics-backups/restic"
export AWS_ACCESS_KEY_ID="AKIAIOSFODNN7EXAMPLE"
export AWS_SECRET_ACCESS_KEY="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

# Backup paths
BACKUP_PATHS=("/app/uploads" "/app/logs" "/app/config")

# Retention policy
RETENTION_OPTIONS="--keep-daily 7 --keep-weekly 4 --keep-monthly 6"

# Initialize repository (first time only)
# restic init

# Perform backup
restic backup ${BACKUP_PATHS} \
    --tag production \
    --tag automated \
    ${RETENTION_OPTIONS}

# Prune old backups
restic prune --max-unused 10%

# Check repository integrity
restic check

# List backups
restic snapshots --latest 5
```

## Cloud Storage Configuration

### AWS S3 Bucket Policy

Create `infrastructure/s3-bucket-policy.json`:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "EnforceSSL",
      "Effect": "Deny",
      "Principal": "*",
      "Action": "s3:*",
      "Resource": [
        "arn:aws:s3:::biometrics-backups",
        "arn:aws:s3:::biometrics-backups/*"
      ],
      "Condition": {
        "Bool": {
          "aws:SecureTransport": "false"
        }
      }
    },
    {
      "Sid": "AllowBackupServerWrite",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::123456789012:role/backup-server"
      },
      "Action": [
        "s3:PutObject",
        "s3:PutObjectAcl"
      ],
      "Resource": "arn:aws:s3:::biometrics-backups/*"
    },
    {
      "Sid": "AllowReadForRestore",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::123456789012:role/restore-server"
      },
      "Action": [
        "s3:GetObject",
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::biometrics-backups",
        "arn:aws:s3:::biometrics-backups/*"
      ]
    },
    {
      "Sid": "DenyDeleteBackup",
      "Effect": "Deny",
      "Principal": "*",
      "Action": "s3:DeleteObject",
      "Resource": "arn:aws:s3:::biometrics-backups/*",
      "Condition": {
        "StringLike": {
          "s3:ObjectTag/Retention": "active"
        }
      }
    }
  ]
}
```

### Lifecycle Configuration

Create `infrastructure/s3-lifecycle.json`:

```json
{
  "Rules": [
    {
      "ID": "MoveToGlacierAfter30Days",
      "Status": "Enabled",
      "Filter": {
        "Prefix": "postgres/"
      },
      "Transitions": [
        {
          "Days": 30,
          "StorageClass": "GLACIER"
        },
        {
          "Days": 90,
          "StorageClass": "DEEP_ARCHIVE"
        }
      ],
      "Expiration": {
        "Days": 365
      }
    },
    {
      "ID": "DeleteOldFilesAfter90Days",
      "Status": "Enabled",
      "Filter": {
        "Prefix": "temp/"
      },
      "Expiration": {
        "Days": 90
      }
    }
  ]
}
```

## Docker Compose for Backups

Create `backup/docker-compose.backup.yml`:

```yaml
version: '3.8'

services:
  backup-postgres:
    image: wal-g/wal-g:v2.0.1
    container_name: backup-postgres
    environment:
      - AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
      - AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
      - AWS_REGION=${AWS_REGION}
      - WALG_S3_PREFIX=${WALG_S3_PREFIX}
      - POSTGRES_CONNECTION_STRING=postgresql://postgres:${POSTGRES_PASSWORD}@db:5432/biometrics
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    networks:
      - backup
    restart: unless-stopped
    command: >
      sh -c "
        while true; do
          /usr/local/bin/wal-g backup-push $${POSTGRES_CONNECTION_STRING} &&
          sleep 3600;
        done
      "

  backup-files:
    image: restic/restic:0.16.1
    container_name: backup-files
    environment:
      - RESTIC_PASSWORD=${RESTIC_PASSWORD}
      - RESTIC_REPOSITORY=${RESTIC_REPOSITORY}
      - AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
      - AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
    volumes:
      - app-data:/data
    networks:
      - backup
    restart: unless-stopped
    command: >
      sh -c "
        restic backup /data/uploads /data/config --tag production &&
        restic forget --keep-daily 7 --keep-weekly 4 --keep-monthly 6 --prune &&
        sleep 3600
      "

  cron-scheduler:
    image: alpine:latest
    container_name: cron-scheduler
    volumes:
      - ./cron-scripts:/scripts
    networks:
      - backup
    restart: unless-stopped
    command: >
      sh -c "
        apk add --no-cache docker-cli curl &&
        while true; do
          sleep 3600
        done
      "

volumes:
  app-data:

networks:
  backup:
    driver: bridge
```

## Backup Verification

### Test Restore Script

Create `backup/scripts/test-restore.sh`:

```bash
#!/bin/bash

set -e

# Configuration
TEST_DB="biometrics_restore_test"
BACKUP_NAME=${1:-latest}
TIMESTAMP=$(date +%Y%m%d-%H%M%S)

echo "[$(date)] Starting backup restoration test..."

# Create test database
echo "[$(date)] Creating test database: ${TEST_DB}"
psql -c "DROP DATABASE IF EXISTS ${TEST_DB};"
psql -c "CREATE DATABASE ${TEST_DB};"

# Restore from backup
echo "[$(date)] Restoring backup: ${BACKUP_NAME}"
wal-g backup-fetch ${TEST_DB} ${BACKUP_NAME}

# Run integrity checks
echo "[$(date)] Running integrity checks..."
psql -d ${TEST_DB} -c "SELECT count(*) FROM pg_class;" > /dev/null
psql -d ${TEST_DB} -c "ANALYZE;" > /dev/null

# Compare row counts with production
echo "[$(date)] Comparing with production..."
PROD_COUNT=$(psql -t -c "SELECT count(*) FROM users;" biometrics)
TEST_COUNT=$(psql -t -c "SELECT count(*) FROM users;" ${TEST_DB})

if [ "${PROD_COUNT}" = "${TEST_COUNT}" ]; then
    echo "[$(date)] ✓ Row counts match"
else
    echo "[$(date)] ✗ Row count mismatch: prod=${PROD_COUNT}, test=${TEST_COUNT}"
    exit 1
fi

# Test specific tables
for TABLE in "users" "sessions" "biometrics_data"; do
    PROD_CHECK=$(psql -t -c "SELECT count(*) FROM ${TABLE};" biometrics 2>/dev/null || echo "0")
    TEST_CHECK=$(psql -t -c "SELECT count(*) FROM ${TABLE};" ${TEST_DB} 2>/dev/null || echo "0")
    
    if [ "${PROD_CHECK}" = "${TEST_CHECK}" ]; then
        echo "[$(date)] ✓ Table ${TABLE}: ${TEST_CHECK} rows"
    else
        echo "[$(date)] ✗ Table ${TABLE} mismatch: prod=${PROD_CHECK}, test=${TEST_CHECK}"
    fi
done

# Cleanup test database
echo "[$(date)] Cleaning up test database..."
psql -c "DROP DATABASE IF EXISTS ${TEST_DB};"

echo "[$(date)] Backup restoration test completed successfully"
```

### Daily Verification Cron

```bash
# Add to crontab
0 6 * * * /app/backup/scripts/test-restore.sh >> /var/log/backup-test.log 2>&1
```

## Disaster Recovery Plan

### Recovery Procedures

1. **Database Recovery**
   ```bash
   # Stop application
   docker-compose stop app
   
   # Restore from latest backup
   wal-g backup-fetch biometrics latest --restore
   
   # Verify data
   psql -c "SELECT count(*) FROM users;" biometrics
   
   # Start application
   docker-compose start app
   ```

2. **Full System Recovery**
   ```bash
   # Restore database
   ./restore-database.sh latest
   
   # Restore files
   ./restore-files.sh latest
   
   # Verify services
   curl -f http://localhost:3000/health
   
   # Check logs
   tail -f /var/log/biometrics/app.log
   ```

### Recovery Time Objectives

| Component | RTO (Recovery Time Objective) | RPO (Recovery Point Objective) |
|-----------|------------------------------|-------------------------------|
| Database | 30 minutes | 1 hour |
| File Storage | 1 hour | 24 hours |
| Application | 15 minutes | N/A |
| Configuration | 15 minutes | 24 hours |

## Backup Monitoring

### Prometheus Metrics for Backups

```yaml
# Add to prometheus.yml
- job_name: 'backup-exporter'
  static_configs:
    - targets: ['backup-metrics:8080']
```

```go
// Simple backup metrics HTTP handler
package main

import (
    "os/exec"
    "encoding/json"
    "net/http"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    lastBackupTime = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "backup_last_timestamp",
            Help: "Timestamp of last backup",
        },
        []string{"type", "status"},
    )
    backupDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "backup_duration_seconds",
            Help:    "Duration of backup operations",
            Buckets: prometheus.DefBuckets,
        },
        []string{"type"},
    )
)

func main() {
    prometheus.MustRegister(lastBackupTime, backupDuration)
    
    http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
        // Update metrics from backup status
        lastBackupTime.WithLabelValues("postgres", "success").Set(float64(time.Now().Unix()))
    })
    
    http.Handle("/metrics", promhttp.Handler())
    http.ListenAndServe(":8080", nil)
}
```

## Related Documentation

- [MONITORING-SETUP.md](./MONITORING-SETUP.md)
- [ALERTING-CONFIG.md](./ALERTING-CONFIG.md)
- [CI-CD-PIPELINE.md](./CI-CD-PIPELINE.md)

---

**End of Backup Strategy Documentation**
