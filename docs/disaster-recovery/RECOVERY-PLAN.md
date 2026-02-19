# 🚨 BIOMETRICS DISASTER RECOVERY PLAN

**Version:** 1.0.0  
**Last Updated:** 2026-02-19  
**Status:** Production-Critical  
**RTO:** 4 hours  
**RPO:** 1 hour  

---

## 📋 Table of Contents

1. [Disaster Recovery Overview](#disaster-recovery-overview)
2. [Backup Strategy](#backup-strategy)
3. [Recovery Procedures](#recovery-procedures)
4. [Data Loss Prevention](#data-loss-prevention)
5. [Business Continuity](#business-continuity)
6. [Testing & Validation](#testing--validation)

---

## Disaster Recovery Overview

### Recovery Objectives

| Metric | Target | Description |
|--------|--------|-------------|
| **RTO (Recovery Time Objective)** | 4 hours | Maximum acceptable downtime |
| **RPO (Recovery Point Objective)** | 1 hour | Maximum acceptable data loss |
| **MTTR (Mean Time To Recovery)** | 2 hours | Average recovery time |

### Disaster Scenarios

| Scenario | Severity | Recovery Strategy |
|----------|----------|-------------------|
| Single Pod Failure | Low | Auto-healing (Kubernetes) |
| Node Failure | Medium | Pod rescheduling |
| Database Corruption | High | Point-in-time recovery |
| Region Outage | Critical | Multi-region failover |
| Data Center Loss | Critical | Full restore from backup |

---

## Backup Strategy

### Automated Backups

#### PostgreSQL Database

```yaml
# Backup Schedule
- Frequency: Every 15 minutes (WAL archiving)
- Full Backup: Daily at 02:00 UTC
- Retention: 30 days
- Storage: S3/GCS with versioning

# Backup Commands
pg_basebackup -h biometrics-postgres -U postgres -D /backup/base -Ft -z -P
pg_dump -h biometrics-postgres -U postgres biometrics > /backup/daily.sql
```

#### Redis Cache

```yaml
# Backup Schedule
- Frequency: Every hour (RDB snapshots)
- AOF: Enabled (every second)
- Retention: 7 days
- Storage: Persistent volumes

# Backup Commands
redis-cli BGSAVE
redis-cli BGREWRITEAOF
```

### Backup Locations

| Backup Type | Primary | Secondary | Tertiary |
|-------------|---------|-----------|----------|
| **Database** | PVC | S3 Bucket | Cross-region S3 |
| **Application** | Git Repository | Docker Registry | Artifact Storage |
| **Configuration** | ConfigMaps | Git | Encrypted Storage |
| **Secrets** | Vault | Encrypted S3 | Hardware Security Module |

---

## Recovery Procedures

### Scenario 1: Database Corruption

#### Immediate Actions

```bash
# 1. Stop all write operations
kubectl scale deployment biometrics-opencode-server --replicas=0 -n biometrics

# 2. Identify corruption
kubectl exec -it statefulset/biometrics-postgres -n biometrics -- \
  pg_checksums --check

# 3. Backup current state (for forensics)
kubectl exec -it statefulset/biometrics-postgres -n biometrics -- \
  pg_dump biometrics > corruption_backup_$(date +%Y%m%d_%H%M%S).sql
```

#### Recovery Steps

```bash
# 1. Find latest good backup
BACKUP_DATE=$(aws s3 ls s3://biometrics-backups/postgres/ | tail -1 | awk '{print $1}')

# 2. Download backup
aws s3 cp s3://biometrics-backups/postgres/${BACKUP_DATE}/backup.sql.gz ./

# 3. Restore database
kubectl exec -i statefulset/biometrics-postgres -n biometrics -- \
  psql -U postgres biometrics < backup.sql

# 4. Verify data integrity
kubectl exec -it statefulset/biometrics-postgres -n biometrics -- \
  psql -U postgres -c "SELECT COUNT(*) FROM information_schema.tables;"

# 5. Restart application
kubectl scale deployment biometrics-opencode-server --replicas=3 -n biometrics
```

### Scenario 2: Complete Cluster Loss

#### Multi-Region Failover

```bash
# 1. Activate DR cluster
kubectl config use-context biometrics-dr-cluster

# 2. Restore from cross-region backup
./scripts/restore-cluster.sh --region dr-region --backup latest

# 3. Update DNS
aws route53 change-resource-record-sets \
  --hosted-zone-id ZONE_ID \
  --change-batch file://dns-failover.json

# 4. Verify services
./scripts/verify-cluster-health.sh
```

### Scenario 3: Application Data Loss

#### Recovery from GitOps

```bash
# 1. Find last known good commit
git log --oneline | grep "stable"

# 2. Rollback Helm release
helm rollback biometrics <REVISION> -n biometrics

# 3. Or redeploy from Git
git checkout <COMMIT_HASH>
helm upgrade biometrics ./helm/biometrics -n biometrics
```

---

## Data Loss Prevention

### Prevention Measures

#### Database Protection

```yaml
# Enable Point-in-Time Recovery
postgresql:
  walLevel: logical
  maxReplicationSlots: 10
  maxWALSenders: 10
  archiveMode: "on"
  archiveCommand: "wal-g wal-push"
  
# High Availability
replication:
  enabled: true
  replicas: 2
  synchronousCommit: "on"
```

#### Application Protection

```yaml
# Pod Disruption Budget
pdb:
  enabled: true
  minAvailable: 2
  
# Resource Quotas
quotas:
  requests.cpu: "10"
  requests.memory: "20Gi"
  limits.cpu: "20"
  limits.memory: "40Gi"
```

### Monitoring & Alerts

```yaml
# Critical Alerts
alerts:
  - name: DatabaseDown
    condition: postgres_up == 0
    severity: critical
    action: Page on-call immediately
    
  - name: BackupFailed
    condition: backup_success == 0
    severity: critical
    action: Investigate and re-run backup
    
  - name: ReplicationLag
    condition: replication_lag_seconds > 300
    severity: warning
    action: Monitor and prepare failover
```

---

## Business Continuity

### Failover Strategy

#### Active-Passive (Current)

```
Primary Cluster (us-east-1)
    ↓ [Health Check Failure]
Secondary Cluster (us-west-2)
    ↓ [DNS Update]
Traffic Routed to DR
```

#### Active-Active (Future)

```
Traffic → Load Balancer
    ↓           ↓
Cluster A   Cluster B
    ↓           ↓
Database Replication (Multi-Master)
```

### Communication Plan

| Stakeholder | Notification Method | Timeline |
|-------------|-------------------|----------|
| **Engineering** | PagerDuty | Immediate |
| **Management** | Slack + Email | Within 15 min |
| **Customers** | Status Page | Within 30 min |
| **Public** | Twitter/Website | Within 1 hour |

---

## Testing & Validation

### Quarterly DR Tests

#### Test Checklist

```markdown
## Pre-Test Preparation
- [ ] Schedule maintenance window
- [ ] Notify stakeholders
- [ ] Backup current state
- [ ] Prepare rollback plan

## Test Execution
- [ ] Simulate disaster scenario
- [ ] Execute recovery procedures
- [ ] Measure RTO/RPO
- [ ] Validate data integrity
- [ ] Test application functionality

## Post-Test Activities
- [ ] Document results
- [ ] Identify gaps
- [ ] Update runbooks
- [ ] Train team on lessons learned
```

### Test Scenarios

| Scenario | Frequency | Last Test | Next Test | Result |
|----------|-----------|-----------|-----------|--------|
| Database Restore | Monthly | 2026-01-15 | 2026-02-15 | ✅ PASS |
| Cluster Failover | Quarterly | 2026-01-01 | 2026-04-01 | ✅ PASS |
| Backup Verification | Weekly | 2026-02-12 | 2026-02-19 | ✅ PASS |
| Application Rollback | Monthly | 2026-02-01 | 2026-03-01 | ✅ PASS |

---

## Recovery Scripts

### Automated Recovery

```bash
#!/bin/bash
# disaster-recovery.sh - Automated DR execution

set -euo pipefail

SCENARIO=${1:-database}
BACKUP_DATE=${2:-latest}

case $SCENARIO in
  database)
    echo "Starting database recovery..."
    ./scripts/restore-database.sh --date $BACKUP_DATE
    ;;
  cluster)
    echo "Starting cluster failover..."
    ./scripts/failover-cluster.sh
    ;;
  application)
    echo "Starting application rollback..."
    helm rollback biometrics -n biometrics
    ;;
  *)
    echo "Unknown scenario: $SCENARIO"
    exit 1
    ;;
esac
```

### Verification Script

```bash
#!/bin/bash
# verify-recovery.sh - Post-recovery validation

echo "=== POST-RECOVERY VERIFICATION ==="

# Check database
echo "Checking database..."
kubectl exec -it statefulset/biometrics-postgres -n biometrics -- \
  psql -U postgres -c "SELECT version();"

# Check application
echo "Checking application health..."
curl -f http://biometrics.delqhi.com/global/health

# Check data integrity
echo "Verifying data integrity..."
./scripts/verify-data-integrity.sh

echo "=== VERIFICATION COMPLETE ==="
```

---

## Contact Information

### Emergency Contacts

| Role | Name | Phone | Email |
|------|------|-------|-------|
| **On-Call Engineer** | Rotating | +1-XXX-XXX-XXXX | oncall@delqhi.com |
| **DevOps Lead** | Jeremy | +1-XXX-XXX-XXXX | jeremy@delqhi.com |
| **CTO** | [Name] | +1-XXX-XXX-XXXX | cto@delqhi.com |

### External Support

| Vendor | Support Level | Contact |
|--------|--------------|---------|
| **Kubernetes** | Enterprise | support@k8s.io |
| **PostgreSQL** | Community | pgsql-general@postgresql.org |
| **Cloud Provider** | Premium | support@cloud-provider.com |

---

## Appendix

### Recovery Runbooks

1. **DB-RUNBOOK-001**: PostgreSQL Point-in-Time Recovery
2. **APP-RUNBOOK-001**: Application Rollback Procedures
3. **CLUSTER-RUNBOOK-001**: Multi-Region Failover
4. **BACKUP-RUNBOOK-001**: Backup Verification & Testing

### Related Documentation

- [Production Deployment Guide](./PRODUCTION.md)
- [Monitoring & Alerting](./MONITORING.md)
- [Security Best Practices](./SECURITY.md)
- [Incident Response](./INCIDENT-RESPONSE.md)

---

**Last Reviewed:** 2026-02-19  
**Next Review:** 2026-05-19  
**Owner:** DevOps Team  
