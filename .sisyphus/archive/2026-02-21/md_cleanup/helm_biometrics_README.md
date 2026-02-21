# Biometrics Helm Chart - Deployment Guide

## Prerequisites

- Kubernetes 1.25+
- Helm 3.10+
- PV provisioner support in the underlying infrastructure
- NVIDIA API Key (required)
- Ingress controller (nginx recommended)
- cert-manager for TLS (optional but recommended)

## Installation

### 1. Add Helm Repository (if published)

```bash
helm repo add delqhi https://charts.delqhi.com
helm repo update
```

### 2. Install from Local Directory

```bash
# Navigate to helm directory
cd /Users/jeremy/dev/BIOMETRICS/helm

# Install with default values
helm install biometrics ./biometrics --namespace biometrics --create-namespace

# Install with custom values
helm install biometrics ./biometrics \
  --namespace biometrics \
  --create-namespace \
  -f custom-values.yaml

# Install with secrets
helm install biometrics ./biometrics \
  --namespace biometrics \
  --create-namespace \
  --set secrets.nvidiaApiKey="nvapi-YOUR_KEY_HERE" \
  --set secrets.dbPassword="your-secure-password"
```

### 3. Verify Installation

```bash
# Check deployment status
helm status biometrics -n biometrics

# List all resources
kubectl get all -n biometrics

# Check pods are running
kubectl get pods -n biometrics

# Check services
kubectl get svc -n biometrics
```

## Configuration

### Minimal Configuration

```yaml
# minimal-values.yaml
replicaCount: 2

secrets:
  nvidiaApiKey: "nvapi-YOUR_KEY"
  dbPassword: "secure-password"

ingress:
  enabled: false
```

### Production Configuration

```yaml
# production-values.yaml
replicaCount: 5

opencodeServer:
  resources:
    requests:
      memory: 2Gi
      cpu: 1000m
    limits:
      memory: 8Gi
      cpu: 4000m
  hpa:
    enabled: true
    minReplicas: 5
    maxReplicas: 20

postgresql:
  primary:
    persistence:
      size: 50Gi

ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  tls:
    - secretName: biometrics-tls
      hosts:
        - biometrics.example.com
```

### NVIDIA API Key Setup

**Option 1: Via CLI**
```bash
helm install biometrics ./biometrics \
  --set secrets.nvidiaApiKey="nvapi-YOUR_KEY"
```

**Option 2: Via values.yaml**
```yaml
secrets:
  nvidiaApiKey: "nvapi-YOUR_KEY"
```

**Option 3: Via Kubernetes Secret**
```bash
kubectl create secret generic biometrics-secrets \
  --from-literal=nvidia-api-key="nvapi-YOUR_KEY" \
  --from-literal=db-password="secure-password" \
  -n biometrics
```

## Upgrade

```bash
# Upgrade with new values
helm upgrade biometrics ./biometrics \
  --namespace biometrics \
  -f new-values.yaml

# Upgrade and recreate pods
helm upgrade biometrics ./biometrics \
  --namespace biometrics \
  --force \
  --recreate-pods
```

## Uninstall

```bash
# Uninstall chart
helm uninstall biometrics -n biometrics

# Delete namespace (removes everything)
kubectl delete namespace biometrics
```

## Access Services

### Port Forward (Development)

```bash
# OpenCode Server
kubectl port-forward svc/biometrics-opencode-server 8080:8080 -n biometrics

# Grafana
kubectl port-forward svc/biometrics-grafana 3000:3000 -n biometrics

# Prometheus
kubectl port-forward svc/biometrics-prometheus 9090:9090 -n biometrics
```

### Via Ingress (Production)

Once DNS is configured:
- **Application**: https://biometrics.delqhi.com
- **API**: https://api.biometrics.delqhi.com
- **Grafana**: https://grafana.biometrics.delqhi.com
- **Prometheus**: https://prometheus.biometrics.delqhi.com

## Monitoring

### Grafana Dashboards

1. Login to Grafana (admin/admin)
2. Import dashboards from `/dashboards/`
3. Default dashboards included:
   - Application Overview
   - PostgreSQL Metrics
   - Redis Metrics
   - AI Agent Performance

### Prometheus Metrics

Access Prometheus at: http://prometheus.biometrics.delqhi.com

Key metrics:
- `opencode_server_requests_total`
- `opencode_server_request_duration_seconds`
- `ai_agent_tasks_total`
- `ai_agent_task_duration_seconds`

## Troubleshooting

### Check Pod Logs

```bash
# OpenCode Server logs
kubectl logs -f deployment/biometrics-opencode-server -n biometrics

# PostgreSQL logs
kubectl logs -f statefulset/biometrics-postgres -n biometrics

# Redis logs
kubectl logs -f statefulset/biometrics-redis -n biometrics
```

### Debug Commands

```bash
# Describe pod
kubectl describe pod <pod-name> -n biometrics

# Check events
kubectl get events -n biometrics --sort-by='.lastTimestamp'

# Exec into pod
kubectl exec -it <pod-name> -n biometrics -- /bin/sh
```

### Common Issues

**Issue: Pods not starting**
```bash
# Check resource quotas
kubectl describe quota -n biometrics

# Check PVC status
kubectl get pvc -n biometrics
```

**Issue: Database connection errors**
```bash
# Verify secrets
kubectl get secret biometrics-secrets -n biometrics -o yaml

# Test database connection
kubectl exec -it deployment/biometrics-opencode-server -n biometrics -- \
  psql -h biometrics-postgres -U postgres -d biometrics
```

## Backup & Restore

### Backup Database

```bash
# Create backup
kubectl exec -it statefulset/biometrics-postgres -n biometrics -- \
  pg_dump -U postgres biometrics > backup.sql

# Copy to local
kubectl cp biometrics-postgres-0:/var/lib/postgresql/data ./backup -n biometrics
```

### Restore Database

```bash
# Restore from backup
kubectl exec -i statefulset/biometrics-postgres -n biometrics -- \
  psql -U postgres biometrics < backup.sql
```

## Best Practices

1. **Always use secrets** for sensitive data
2. **Enable HPA** for production workloads
3. **Use persistent volumes** for stateful services
4. **Enable monitoring** with Prometheus/Grafana
5. **Configure resource limits** to prevent resource starvation
6. **Use network policies** for security
7. **Regular backups** of PostgreSQL database
8. **Test disaster recovery** procedures regularly

## Support

- **Documentation**: https://github.com/Delqhi/BIOMETRICS/docs
- **Issues**: https://github.com/Delqhi/BIOMETRICS/issues
- **Discord**: https://discord.gg/biometrics
- **Email**: support@delqhi.com
