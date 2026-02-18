# GITOPS-ARGOCD.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Architekturentscheidungen folgen den globalen Regeln aus `AGENTS-GLOBAL.md`.
- Jede Strukturänderung benötigt Mapping- und Integrationsabgleich.
- Security-by-Default, NLM-First und Nachweisbarkeit sind Pflicht.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

---

## 1) Overview

This document describes the GitOps implementation using ArgoCD for BIOMETRICS. ArgoCD provides declarative, Git-based continuous delivery with automated application deployment to Kubernetes clusters.

### Core Features
- Git-based declarative deployments
- Automated sync and drift detection
- Multi-cluster management
- Progressive delivery (Canary, Blue-Green)
- Secret management integration
- Visual UI and CLI

---

## 2) Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        BIOMETRICS GITOPS ARCHITECTURE                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                          GIT REPOSITORIES                            │   │
│  │                                                                       │   │
│  │  ┌──────────────────┐    ┌──────────────────┐    ┌────────────────┐ │   │
│  │  │   K8s Manifests │    │  Helm Charts     │    │ Kustomize     │ │   │
│  │  │   (YAML)        │    │  (Templates)     │    │ (Overlays)    │ │   │
│  │  └────────┬─────────┘    └────────┬─────────┘    └───────┬────────┘ │   │
│  │           │                       │                      │          │   │
│  │           └───────────────────────┼──────────────────────┘          │   │
│  │                                   ▼                                   │   │
│  │                    ┌────────────────────────┐                        │   │
│  │                    │    GitHub/GitLab       │                        │   │
│  │                    │    Webhooks            │                        │   │
│  │                    └────────────────────────┘                        │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                      │                                       │
│                                      ▼                                       │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      ARGOCD SERVER                                  │   │
│  │                                                                       │   │
│  │   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │   │
│  │   │ Application  │  │   Sync       │  │   Health     │            │   │
│  │   │ Controller   │  │  Controller  │  │  Assessment  │            │   │
│  │   └──────────────┘  └──────────────┘  └──────────────┘            │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                      │                                       │
│                                      ▼                                       │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                       TARGET CLUSTERS                              │   │
│  │                                                                       │   │
│  │  ┌────────────────┐    ┌────────────────┐    ┌────────────────┐      │   │
│  │  │   Production  │    │   Staging     │    │   Development │      │   │
│  │  │   (prod-01)   │    │   (staging)   │    │   (dev)       │      │   │
│  │  └────────────────┘    └────────────────┘    └────────────────┘      │   │
│  │                                                                       │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3) ArgoCD Installation

### Installation YAML

```yaml
# argocd/install.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: argocd
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: argocd-application-controller
  namespace: argocd
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: argocd-application-controller
rules:
  - apiGroups: ["*"]
    resources: ["*"]
    verbs: ["*"]
  - nonResourceURLs: ["*"]
    verbs: ["*"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: argocd-application-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: argocd-application-controller
subjects:
  - kind: ServiceAccount
    name: argocd-application-controller
    namespace: argocd
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: argocd-server
  namespace: argocd
spec:
  selector:
    matchLabels:
      app: argocd-server
  template:
    spec:
      containers:
        - name: argocd-server
          image: argoproj/argocd:latest
          ports:
            - containerPort: 8080
            - containerPort: 8083
          args:
            - /usr/local/bin/argocd-server
            - --staticassets
            - /shared
            - --repo-server
            - argocd-repo-server:8081
            - --insecure
          env:
            - name: ARGOCD_AUTH_TOKEN
              valueFrom:
                secretKeyRef:
                  name: argocd-secret
                  key: auth-token
          volumeMounts:
            - name: static-assets
              mountPath: /shared
      serviceAccountName: argocd-server
```

---

## 4) Application Definitions

### Application CRDs

```yaml
# argocd/applicationset.yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: biometrics-microservices
  namespace: argocd
spec:
  generators:
    - matrix:
        generators:
          - git:
              repoURL: https://github.com/biometrics/k8s-manifests
              revision: main
              directories:
                - path: services/*
          list:
            elements:
              - cluster: prod-cluster
                url: https://kubernetes.default.svc
              - cluster: staging-cluster
                url: https://staging.k8s.example.com
  template:
    metadata:
      name: '{{path.basename}}-{{cluster}}'
    spec:
      project: biometrics
      source:
        repoURL: https://github.com/biometrics/k8s-manifests
        targetRevision: main
        path: '{{path}}'
        kustomize:
          images:
            - biometrics/{{path.basename}}:{{chip}}
      destination:
        server: '{{url}}'
        namespace: biometrics
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
```

---

## 5) Progressive Delivery

### Canary Deployments

```yaml
# argocd/canary-service.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: auth-service-canary
  namespace: argocd
spec:
  project: biometrics
  source:
    repoURL: https://github.com/biometrics/k8s-manifests
    targetRevision: main
    path: services/auth-service
    kustomize:
      images:
        - biometrics/auth-service:v2.0.0
  destination:
    server: https://kubernetes.default.svc
    namespace: biometrics
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
  strategy:
    type: Rolling
    rolling:
      maxSurge: "25%"
      maxUnavailable: "25%"
```

---

## 6) Multi-Cluster Setup

### Cluster Registration

```yaml
# argocd/cluster-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: prod-cluster
  namespace: argocd
  labels:
    argocd.argoproj.io/secret-type: cluster
type: Opaque
stringData:
  name: prod-cluster
  server: https://kubernetes.default.svc
  config: |
    {
      "bearerToken": "<token>",
      "tlsClientConfig": {
        "insecure": false,
        "caData": "<base64-ca-cert>"
      }
    }
```

---

## 7) Helm Integration

### Helm Charts

```yaml
# argocd/helm-release.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: metrics-server
  namespace: argocd
spec:
  project: infrastructure
  source:
    repoURL: https://kubernetes-sigs.github.io/metrics-server
    chart: metrics-server
    targetRevision: 3.8.0
    helm:
      valueFiles:
        - values.yaml
      parameters:
        - name: replicas
          value: "2"
        - name: args[0]
          value: "--kubelet-preferred-address-types=InternalIP"
      releaseName: metrics-server
  destination:
    server: https://kubernetes.default.svc
    namespace: monitoring
  syncPolicy:
    automated:
      prune: true
```

---

## 8) Kustomize Overlays

### Environment Overlays

```yaml
# kustomize/overlays/production/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: biometrics-prod
resources:
  - ../../base
images:
  - name: biometrics/app
    newTag: v2.1.0
patches:
  - patch: |-
      - op: replace
        path: /spec/replicas
        value: 10
    target:
      kind: Deployment
      name: app
replicas:
  - name: app
    count: 10
configMapGenerator:
  - name: app-config
    literals:
      - ENV=production
      - LOG_LEVEL=info
```

---

## 9) Secret Management

### Vault Integration

```yaml
# argocd/vault-plugin.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: auth-service-secure
  namespace: argocd
spec:
  source:
    plugin:
      name: argocd-vault-plugin-replace
      env:
        - name: VAULT_ADDR
          value: https://vault.biometrics.app
        - name: VAULT_ROLE
          value: biometrics-app
    repoURL: https://github.com/biometrics/k8s-manifests
    path: services/auth-service
  destination:
    server: https://kubernetes.default.svc
    namespace: biometrics
```

---

## 10) Sync Strategies

### Sync Options

```yaml
# argocd/sync-options.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: full-sync-app
  namespace: argocd
spec:
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
      allowEmpty: false
    syncOptions:
      - CreateNamespace=true
      - PrunePropagationPolicy=foreground
      - PruneLast=true
      - RespectIgnoreDifferences=true
    retry:
      limit: 5
      backoff:
        duration: 5s
        factor: 2
        maxDuration: 3m
```

---

## 11) Webhook Configuration

### Git Webhooks

```yaml
# argocd/webhook-config.yaml
apiVersion: v1
kind: Secret
metadata:
  name: argocd-webhook
  namespace: argocd
type: Opaque
stringData:
  github.secret: github-webhook-secret
  gitlab.secret: gitlab-webhook-secret
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-cm
  namespace: argocd
data:
  resource.customizations: |
    argoproj.io/Application:
      health.lua: |
        hs = {}
        hs.status = "Progressing"
        hs.message = ""
        if obj.status ~= nil then
          if obj.status.health ~= nil then
            hs.status = obj.status.health.status
            if obj.status.health.message ~= nil then
              hs.message = obj.status.health.message
            end
          end
        end
        return hs
```

---

## 12) Notifications

### ArgoCD Notifications

```yaml
# argocd/notifications.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-notifications-cm
  namespace: argocd
data:
  service.slack: |
    webhook:
      url: https://hooks.slack.com/services/xxx
  template.app-sync-status: |
    message: |
      Application {{.app.metadata.name}} sync is {{.app.status.operationState.phase}}.
    slack:
      attachments: |
        [{
          "title": "{{.app.metadata.name}}",
          "color": "#FF0000",
          "fields": [
            {"title": "Status", "value": "{{.app.status.operationState.phase}}", "short": true},
            {"title": "Repository", "value": "{{.app.spec.source.repoURL}}", "short": true}
          ]
        }]
  trigger.on-sync-failed: |
    - when: app.status.operationState.phase == 'Failed'
      send:
        - app-sync-status
        - slack
  trigger.on-sync-succeeded: |
    - when: app.status.operationState.phase == 'Succeeded'
      send:
        - app-sync-status
```

---

## 13) RBAC Configuration

### Access Control

```yaml
# argocd/rbac-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-rbac-cm
  namespace: argocd
data:
  policy.default: |
    p, role:readonly, applications, get, */*, allow
    p, role:readonly, projects, get, */*, allow
    p, role:readonly, applications, sync, */*, allow
  policy.csv: |
    g, argocd-admins, role:admin
    g, devops-team, role:admin
    g, developers, role:deployer
    g, viewers, role:readonly
    
    p, role:deployer, applications, *, */*, allow
    p, role:deployer, applications, sync, */*, allow
    p, role:deployer, applications, rollback, */*, allow
```

---

## 14) Backup & Disaster Recovery

### Backup Strategy

```bash
#!/bin/bash
# backup-argocd.sh

# Backup all ArgoCD resources
kubectl get applications -n argocd -o yaml > applications-backup.yaml
kubectl get projects -n argocd -o yaml > projects-backup.yaml
kubectl get applicationsets -n argocd -o yaml > applicationsets-backup.yaml

# Backup with version control
git add .
git commit -m "Backup: $(date)"
git push origin main
```

---

## 15) CI/CD Pipeline

### GitHub Actions

```yaml
# .github/workflows/deploy.yaml
name: Deploy to Kubernetes
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup kubectl
        uses: azure/k8s-set-context@v1
        with:
          kubeconfig: ${{ secrets.KUBECONFIG }}
      
      - name: Update ArgoCD Application
        run: |
          kubectl patch application auth-service \
            -n argocd \
            --type merge \
            -p '{"spec":{"source":{"targetRevision":"main"}}}'
      
      - name: Wait for Sync
        run: |
          argocd app wait auth-service --timeout 300
```

---

## 16) Monitoring & Dashboards

### Grafana Integration

```yaml
# argocd/dashboard-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-grafana-dashboard
  namespace: monitoring
  labels:
    grafana_dashboard: "1"
data:
  argocd-dashboard.json: |
    {
      "dashboard": {
        "title": "ArgoCD Performance",
        "panels": [
          {
            "title": "Sync Duration",
            "type": "graph",
            "targets": [
              {
                "expr": "histogram_quantile(0.99, argocd_app_sync_duration_seconds_bucket)"
              }
            ]
          },
          {
            "title": "Reconciliation Count",
            "type": "graph",
            "targets": [
              {
                "expr": "rate(argocd_app_reconcile_total[5m])"
              }
            ]
          }
        ]
      }
    }
```

---

## 17) Troubleshooting Commands

### Common Commands

```bash
# Sync application manually
argocd app sync auth-service

# Wait for sync
argocd app wait auth-service

# Rollback
argocd app rollback auth-service 1

# View resource tree
argocd app resources auth-service

# View logs
argocd app logs auth-service

# Debug sync
argocd app sync auth-service --dry-run --debug
```

---

Status: APPROVED  
Version: 1.0  
Last Updated: Februar 2026
