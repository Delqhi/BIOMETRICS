# Helm Templates Directory

## Overview

This directory contains Kubernetes Helm chart templates for deploying the biometrics application. These templates provide a configurable, production-ready deployment.

## Contents

### Template Files

| File | Description |
|------|-------------|
| deployment.yaml | Main deployment template |
| service.yaml | Service configuration |
| ingress.yaml | Ingress configuration |
| configmap.yaml | Config maps |
| secrets.yaml | Secrets template |
| hpa.yaml | Horizontal pod autoscaler |

### Helper Templates

| File | Description |
|------|-------------|
| _helpers.tpl | Template helpers |
| _labels.tpl | Label definitions |
| _annotations.tpl | Annotation definitions |

## Chart Structure

```
helm/biometrics/
├── Chart.yaml
├── values.yaml
├── values-prod.yaml
├── values-staging.yaml
├── templates/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   ├── configmap.yaml
│   ├── secrets.yaml
│   ├── hpa.yaml
│   └── _helpers.tpl
└── README.md
```

## Usage

### Installation
```bash
# Install chart
helm install biometrics ./helm/biometrics

# Install with values
helm install biometrics ./helm/biometrics -f values-prod.yaml

# Dry run
helm install biometrics ./helm/biometrics --dry-run --debug
```

### Upgrade
```bash
# Upgrade release
helm upgrade biometrics ./helm/biometrics

# Rollback
helm rollback biometrics 1
```

## Configuration

### Values.yaml
```yaml
replicaCount: 3

image:
  repository: biometrics/api
  tag: latest
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 8080

resources:
  limits:
    cpu: 1000m
    memory: 1Gi
  requests:
    cpu: 500m
    memory: 512Mi

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
```

### Environment Values
```bash
# Production
helm install biometrics ./helm/biometrics \
  -f values-prod.yaml

# Staging
helm install biometrics ./helm/biometrics \
  -f values-staging.yaml
```

## Templates

### Deployment Template
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "biometrics.fullname" . }}
  labels:
    {{- include "biometrics.labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      {{- include "biometrics.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "biometrics.selectorLabels" . | nindent 8 }}
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          ports:
            - name: http
              containerPort: {{ .Values.service.port }}
```

## Best Practices

### Security
- Run as non-root
- Use secrets for sensitive data
- Enable RBAC

### Performance
- Set resource limits
- Configure HPA
- Use readiness probes

### Monitoring
- Add Prometheus annotations
- Configure logging
- Set up alerts

## Maintenance

### Update Dependencies
```bash
helm dependency update
```

### Package Chart
```bash
helm package .
```

## See Also

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Helm Documentation](https://helm.sh/docs/)
- [Deployment Guide](../docs/deployment.md)
