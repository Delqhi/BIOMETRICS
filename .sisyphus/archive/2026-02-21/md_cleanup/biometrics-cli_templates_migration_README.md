# Migration Templates

## Overview

The `migration` templates directory contains comprehensive templates for migrating data, configurations, and systems within the biometrics infrastructure. These templates support various migration scenarios including version upgrades, environment transitions, and system consolidations.

## Migration Types

### 1. Database Migrations

Database schema and data migrations handle evolving data structures:

#### Schema Migration

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: schema-migration-{{ .Values.version }}
  namespace: {{ .Values.namespace }}
spec:
  template:
    spec:
      initContainers:
      - name: backup-pre
        image: {{ .Values.backup.image }}
        command: ["/scripts/pre-migration-backup.sh"]
        env:
        - name: DB_NAME
          value: {{ .Values.database.name }}
      containers:
      - name: migrate
        image: {{ .Values.migration.image }}
        command: ["/scripts/migrate.sh"]
        env:
        - name: MIGRATION_PATH
          value: /migrations/{{ .Values.version }}
        - name: DB_HOST
          value: {{ .Values.database.host }}
        volumeMounts:
        - name: migrations
          mountPath: /migrations
      volumes:
      - name: migrations
        configMap:
          name: migration-scripts-{{ .Values.version }}
      restartPolicy: OnFailure
```

#### Data Migration

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: data-migration-{{ .Values.version }}
spec:
  template:
    spec:
      containers:
      - name: migrate-data
        image: {{ .Values.migration.image }}
        command:
        - /scripts/migrate-data.sh
        env:
        - name: SOURCE_DB
          value: {{ .Values.migration.source }}
        - name: TARGET_DB
          value: {{ .Values.migration.target }}
        - name: BATCH_SIZE
          value: "1000"
        - name: PARALLEL_JOBS
          value: "4"
```

### 2. Storage Migration

Data migration between storage systems:

#### PVC Migration

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Values.newPVC.name }}
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Values.newPVC.size }}
  storageClassName: {{ .Values.newPVC.storageClass }}
---
apiVersion: batch/v1
kind: Job
metadata:
  name: pvc-migration
spec:
  template:
    spec:
      containers:
      - name: migrate
        image: alpine/rsync:latest
        command:
        - rsync
        - -avh
        - --progress
        - /source/
        - /target/
        volumeMounts:
        - name: source
          mountPath: /source
        - name: target
          mountPath: /target
      volumes:
      - name: source
        persistentVolumeClaim:
          claimName: {{ .Values.sourcePVC }}
      - name: target
        persistentVolumeClaim:
          claimName: {{ .Values.newPVC.name }}
      restartPolicy: OnFailure
```

### 3. Service Migration

Migrating services between environments or infrastructure:

#### Service Parallel Run

```yaml
apiVersion: v1
kind: Service
metadata:
 .Values.service.name  name: {{ }}-new
spec:
  selector:
    version: new
  ports:
  - port: 80
    targetPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .Values.service.name }}-migration
spec:
  selector:
    version: migration
  ports:
  - port: 80
    targetPort: 8080
```

### 4. Configuration Migration

Moving and updating configurations:

#### ConfigMap Migration

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: config-migration
spec:
  template:
    spec:
      containers:
      - name: migrate
        image: bitnami/kubectl:latest
        command:
        - /bin/sh
        - -c
        - |
          kubectl get configmap -n $OLD_NS -o yaml | \
          sed 's/namespace: $OLD_NS/namespace: $NEW_NS/' | \
          kubectl apply -f -
        env:
        - name: OLD_NS
          value: {{ .Values.namespaces.old }}
        - name: NEW_NS
          value: {{ .Values.namespaces.new }}
```

### 5. Data Transformation Migration

Data format conversions during migration:

#### Format Conversion

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: transform-migration
spec:
  template:
    spec:
      containers:
      - name: transform
        image: {{ .Values.transform.image }}
        command:
        - /scripts/transform.sh
        env:
        - name: INPUT_FORMAT
          value: {{ .Values.transform.inputFormat }}
        - name: OUTPUT_FORMAT
          value: {{ .Values.transform.outputFormat }}
        - name: SCHEMA_VERSION
          value: {{ .Values.transform.schemaVersion }}
```

## Migration Phases

### Phase 1: Planning

- Assess current state
- Define target state
- Identify dependencies
- Estimate timeline
- Risk assessment

### Phase 2: Preparation

- Create backups
- Set up monitoring
- Prepare rollback procedures
- Test in staging
- Document procedures

### Phase 3: Execution

- Run pre-migration checks
- Execute migration scripts
- Monitor progress
- Handle errors
- Verify completion

### Phase 4: Validation

- Verify data integrity
- Test functionality
- Compare configurations
- Performance testing
- User acceptance testing

### Phase 5: Cleanup

- Remove old resources
- Update documentation
- Archive old data
- Decommission legacy systems

## Configuration

### Migration Parameters

```yaml
migration:
  version: "2.0.0"
  
  source:
    database:
      host: old-db.example.com
      port: 5432
      name: biometrics_old
    storage:
      type: nfs
      path: /data/old
  
  target:
    database:
      host: new-db.example.com
      port: 5432
      name: biometrics_new
    storage:
      type: s3
      bucket: biometrics-data
  
  schedule:
    startTime: "02:00"
    maxDuration: "8h"
  
  safety:
    preBackup: true
    dryRun: true
    rollbackEnabled: true
    maxRetries: 3
  
  performance:
    batchSize: 1000
    parallelJobs: 4
    compression: true
```

### Pre-migration Checklist

```yaml
checklist:
  - name: Backup Verification
    status: pending
    
  - name: Stakeholder Notification
    status: pending
    
  - name: Monitoring Activation
    status: pending
    
  - name: Rollback Testing
    status: pending
    
  - name: Data Validation
    status: pending
```

## Rollback Procedures

### Automatic Rollback

```yaml
rollback:
  enabled: true
  trigger:
    - error_rate > 0.05
    - migration_duration > 8h
    - data_loss_detected
  actions:
    - stop_migration
    - restore_backup
    - notify_team
```

### Manual Rollback

```bash
# Stop current migration
kubectl delete job migration-job

# Restore from backup
kubectl exec -it backup-pod -- /scripts/restore.sh backup-20240115

# Verify restoration
kubectl exec -it db-pod -- psql -c "SELECT COUNT(*) FROM users"
```

## Monitoring

### Migration Metrics

```yaml
metrics:
  - name: migration_progress
    type: gauge
    help: Percentage of migration completed
    
  - name: migration_duration_seconds
    type: histogram
    help: Duration of migration
    
  - name: migration_errors_total
    type: counter
    help: Number of migration errors
```

### Progress Tracking

```bash
# View migration progress
kubectl logs job/migration-job -f

# Check migration status
kubectl get migration {{ .Values.migration.name }} -o yaml
```

## Validation

### Data Integrity Checks

```yaml
integrity:
  checks:
    - type: row_count
      source: "SELECT COUNT(*) FROM {{ .Values.table }}"
      target: "SELECT COUNT(*) FROM {{ .Values.newTable }}"
      
    - type: checksum
      algorithm: sha256
      tables:
        - users
        - biometric_data
        - access_logs
        
    - type: referential_integrity
      foreignKeys:
        - user_id
        - session_id
```

### Performance Validation

```yaml
performance:
  benchmarks:
    - query: "SELECT * FROM biometric_data WHERE user_id = ?"
      maxLatency: 100ms
      
    - query: "SELECT COUNT(*) FROM logs WHERE created_at > ?"
      maxLatency: 500ms
```

## Best Practices

1. **Always Backup First**: Never migrate without a verified backup
2. **Test Thoroughly**: Validate in non-production first
3. **Incremental Approach**: Migrate in phases when possible
4. **Monitor Everything**: Track all migration metrics
5. **Plan Rollback**: Always have a tested rollback plan
6. **Communicate**: Keep stakeholders informed
7. **Document**: Record all steps and decisions

## Related Documentation

- [Database Schema](../docs/schema.md)
- [Data Transformation](../docs/transformations.md)
- [Rollback Procedures](../docs/rollback.md)
