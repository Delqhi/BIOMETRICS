# Kubernetes Helm Charts

**Purpose:** Kubernetes deployment manifests using Helm for the BIOMETRICS project

## Overview

This directory contains Helm charts for deploying BIOMETRICS services to Kubernetes clusters. Helm provides templated Kubernetes manifests with versioning and rollback capabilities.

## Charts

### Main Chart: `biometrics/`

The primary chart for deploying the BIOMETRICS application.

```
biometrics/
├── Chart.yaml          # Chart metadata
├── values.yaml         # Default configuration
├── values-prod.yaml   # Production overrides
├── values-staging.yaml # Staging overrides
├── templates/         # Kubernetes manifests
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   ├── configmap.yaml
│   └── secrets.yaml
└── charts/            # Sub-charts
```

## Quick Start

### Install Chart
```bash
cd /Users/jeremy/dev/BIOMETRICS/helm
helm install biometrics biometrics/ -n biometrics --create-namespace
```

### Upgrade
```bash
helm upgrade biometrics biometrics/ -n biometrics
```

### Uninstall
```bash
helm uninstall biometrics -n biometrics
```

## Configuration

### Values Files

| File | Environment | Description |
|------|-------------|-------------|
| `values.yaml` | All | Base configuration |
| `values-staging.yaml` | Staging | Staging overrides |
| `values-prod.yaml` | Production | Production overrides |

### Key Configuration Options

```yaml
# values.yaml
replicaCount: 3

image:
  repository: delqhi/biometrics
  tag: latest
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 8080

resources:
  limits:
    cpu: 2000m
    memory: 2Gi
  requests:
    cpu: 500m
    memory: 512Mi

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
```

## Dependencies

### External Dependencies

| Dependency | Version | Purpose |
|------------|---------|---------|
| PostgreSQL | 14+ | Primary database |
| Redis | 7+ | Caching and sessions |
| Vault | 1.12+ | Secrets management |

### Installing Dependencies
```bash
helm dependency build biometrics/
```

## Ingress Configuration

### TLS Setup
```yaml
ingress:
  enabled: true
  className: nginx
  tls:
    - hosts:
        - biometrics.delqhi.com
      secretName: biometrics-tls
```

## Monitoring

### Prometheus Integration
```yaml
prometheus:
  enabled: true
  scrapeInterval: 30s
```

### Grafana Dashboards
Dashboards are automatically provisioned when `prometheus.enabled: true`.

## Troubleshooting

### Pod Not Starting
```bash
kubectl describe pod -n biometrics <pod-name>
kubectl logs -n biometrics <pod-name>
```

### Service Not Accessible
```bash
kubectl get svc -n biometrics
kubectl get endpoints -n biometrics
```

### Resource Issues
```bash
kubectl top pods -n biometrics
kubectl describe resources -n biometrics
```

## CI/CD Integration

### GitHub Actions
```yaml
- name: Deploy to Kubernetes
  uses: azure/k8s-set-context@v1
  with:
    kubeconfig: ${{ secrets.KUBE_CONFIG }}
- name: Deploy chart
  run: |
    helm upgrade biometrics ./helm/biometrics \
      --namespace biometrics \
      --install \
      --wait
```

## Maintenance

### Updating Chart Version
```bash
helm package biometrics/
helm repo index .
```

### Rolling Back
```bash
helm rollback biometrics 1 -n biometrics
```

## Related Documentation

- [Kubernetes Architecture](../architecture/kubernetes.md)
- [Deployment Guide](../deployment/)
- [Monitoring Setup](../monitoring/)
- [Security Configuration](../security/)
