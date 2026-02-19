# Backup Templates

## Overview

The `backup` templates directory provides reusable templates for implementing backup solutions across the biometrics system. These templates ensure consistent, reliable backup procedures for databases, files, and system state.

## Purpose

These templates serve multiple critical functions:

- **Data Protection**: Prevent data loss from hardware failures, accidents, or security incidents
- **Disaster Recovery**: Enable system restoration after catastrophic events
- **Compliance**: Meet regulatory requirements for data retention
- **Testing**: Provide known-good data states for testing environments
- **Migration**: Facilitate data movement between environments

## Template Categories

### 1. Database Backup Templates

#### PostgreSQL Backup

```yaml
apiVersion: v1
kind: CronJob
metadata:
  name: {{ .Values.backup.postgres.name }}
  namespace: {{ .Values.namespace }}
spec:
  schedule: {{ .Values.backup.postgres.schedule }}
  successfulJobsHistoryLimit: 3
  failedJobsHistoryLimit: 3
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: {{ .Values.backup.image }}
            command:
            - /scripts/backup-postgres.sh
            env:
            - name: DB_HOST
              value: {{ .Values.database.host }}
            - name: DB_NAME
              value: {{ .Values.database.name }}
            - name: BACKUP_PATH
              value: {{ .Values.backup.storage.path }}
            volumeMounts:
            - name: scripts
              mountPath: /scripts
            - name: backup-storage
              mountPath: /backups
          volumes:
          - name: scripts
            configMap:
              name: {{ .Values.backup.configMap }}
          - name: backup-storage
            persistentVolumeClaim:
              claimName: {{ .Values.backup.pvc }}
          restartPolicy: OnFailure
```

#### Redis Backup

```yaml
apiVersion: v1
kind: CronJob
metadata:
  name: {{ .Values.backup.redis.name }}
spec:
  schedule: {{ .Values.backup.redis.schedule }}
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: redis-backup
            command:
            - redis-cli
            - SAVE
            env:
            - name: REDIS_HOST
              value: {{ .Values.redis.host }}
          restartPolicy: OnFailure
```

### 2. File System Backup Templates

#### Persistent Volume Backup

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: {{ .Values.backup.pvc.name }}-backup
spec:
  template:
    spec:
      containers:
      - name: velero
        image: velero/velero:latest
        command:
        - velero
        - backup
        - create
        - {{ .Values.backup.pvc.backupName }}
        - --include-namespaces
        - {{ .Values.namespace }}
        - --snapshot-volumes
        volumes:
        - name: credentials
          secret:
            secretName: velero-credentials
      restartPolicy: OnFailure
```

#### Directory Sync Backup

```yaml
apiVersion: v1
kind: CronJob
metadata:
  name: directory-sync-backup
spec:
  schedule: "0 3 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: sync
            image: alpine/rsync:latest
            command:
            - rsync
            - -avz
            - --delete
            - /data/
            - /backup/data/
            volumeMounts:
            - name: source
              mountPath: /data
            - name: destination
              mountPath: /backup
          volumes:
          - name: source
            persistentVolumeClaim:
              claimName: {{ .Values.data.pvc }}
          - name: destination
            persistentVolumeClaim:
              claimName: {{ .Values.backup.pvc }}
          restartPolicy: OnFailure
```

### 3. Application State Backup

#### Configuration Backup

```yaml
apiVersion: v1
kind: CronJob
metadata:
  name: config-backup
spec:
  schedule: "0 */6 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: {{ .Values.backup.image }}
            command:
            - /scripts/backup-configs.sh
            env:
            - name: CONFIG_PATH
              value: {{ .Values.config.path }}
            - name: GIT_REPO
              value: {{ .Values.backup.git.repo }}
            volumeMounts:
            - name: config-map
              mountPath: {{ .Values.config.path }}
          restartPolicy: OnFailure
```

### 4. Kubernetes Resources Backup

#### Resource Export

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: k8s-resources-backup
spec:
  template:
    spec:
      containers:
      - name: export
        image: bitnami/kubectl:latest
        command:
        - /bin/sh
        - -c
        - |
          kubectl get all,configmap,secret,ingress,service -n $NAMESPACE -o yaml > /backups/resources.yaml
        env:
        - name: NAMESPACE
          value: {{ .Values.namespace }}
        volumeMounts:
        - name: backup
          mountPath: /backups
      volumes:
      - name: backup
        persistentVolumeClaim:
          claimName: {{ .Values.backup.pvc }}
      restartPolicy: OnFailure
```

### 5. Incremental Backup

#### Point-in-Time Backup

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: incremental-backup
spec:
  schedule: "0 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: {{ .Values.backup.image }}
            command:
            - /scripts/incremental-backup.sh
            env:
            - name: INCREMENTAL
              value: "true"
            - name: BASELINE
              value: {{ .Values.backup.baseline }}
          restartPolicy: OnFailure
```

## Configuration

### Values Schema

```yaml
backup:
  enabled: true
  storage:
    path: /backups
    pvc: backup-storage
    capacity: 100Gi
  schedule:
    postgres: "0 2 * * *"
    redis: "*/15 * * * *"
    files: "0 3 * * *"
    config: "0 */6 * * *"
  retention:
    days: 30
    weekly: 12
    monthly: 12
    yearly: 7
  encryption:
    enabled: true
    keySecret: backup-encryption-key
  compression:
    enabled: true
    algorithm: gzip
```

### Storage Backend Configuration

#### S3 Compatible Storage

```yaml
storage:
  backend: s3
  s3:
    bucket: biometrics-backups
    region: us-west-2
    endpoint: https://s3.amazonaws.com
    prefix: backups/
```

#### NFS Storage

```yaml
storage:
  backend: nfs
  nfs:
    server: nfs.example.com
    path: /backups/biometrics
```

## Restoration Procedures

### Database Restoration

```bash
# List available backups
kubectl exec -it backup-pod -- ls /backups/postgres/

# Restore specific backup
kubectl exec -it backup-pod -- \
  /scripts/restore-postgres.sh /backups/postgres/backup-20240115.sql
```

### File Restoration

```bash
# Restore specific directory
kubectl exec -it backup-pod -- \
  rsync -avz /backups/data/ /restored/data/
```

## Monitoring

### Backup Status

Track backup success:

```yaml
metrics:
  - name: backup_duration_seconds
    type: gauge
    help: Duration of backup operation
  - name: backup_size_bytes
    type: gauge
    size: Size of backup
  - name: backup_last_success
    type: gauge
    help: Timestamp of last successful backup
```

### Alerts

```yaml
alerts:
  - name: BackupFailed
    expr: backup_last_success < (time() - 86400)
    severity: critical
    annotations:
      summary: "Backup has not completed in 24 hours"
```

## Testing

### Backup Verification

Test backup integrity:

```bash
# Verify backup exists
test -f /backups/postgres/backup-$(date +%Y%m%d).sql

# Verify backup is not empty
test -s /backups/postgres/backup-$(date +%Y%m%d).sql

# Test restoration in isolation
kubectl apply -f test-restore.yaml
```

### Recovery Testing

Regular recovery drills:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: backup-recovery-test
spec:
  schedule: "0 4 * * 0"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: test
            command:
            - /scripts/test-restore.sh
          restartPolicy: OnFailure
```

## Best Practices

1. **Regular Scheduling**: Run backups frequently enough to minimize data loss
2. **Verification**: Test restoration procedures regularly
3. **Offsite Copies**: Maintain backups in separate geographic locations
4. **Encryption**: Always encrypt backups at rest and in transit
5. **Monitoring**: Alert on backup failures immediately
6. **Retention**: Follow compliance requirements for retention periods
7. **Documentation**: Document restoration procedures in detail

## Related Documentation

- [Disaster Recovery Plan](../docs/disaster-recovery.md)
- [Storage Configuration](../docs/storage.md)
- [Security Policies](../docs/security.md)
