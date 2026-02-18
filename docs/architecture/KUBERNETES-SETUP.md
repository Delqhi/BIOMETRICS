# KUBERNETES-SETUP.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Architekturentscheidungen folgen den globalen Regeln aus `AGENTS-GLOBAL.md`.
- Jede Strukturänderung benötigt Mapping- und Integrationsabgleich.
- Security-by-Default, NLM-First und Nachweisbarkeit sind Pflicht.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

---

## 1) Overview

This document describes the Kubernetes cluster setup for BIOMETRICS. The Kubernetes infrastructure provides container orchestration, auto-scaling, self-healing, and declarative deployment management for all microservices.

### Core Components
- **Control Plane**: etcd, API Server, Controller Manager, Scheduler
- **Worker Nodes**: Kubelet, Kube Proxy, Container Runtime
- **Addons**: Ingress Controller, Service Mesh, Monitoring, Logging

---

## 2) Cluster Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     BIOMETRICS KUBERNETES CLUSTER                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         CONTROL PLANE                                │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │   etcd     │  │ API Server  │  │ Controller │  │ Scheduler  │ │   │
│  │  │ (3 nodes)  │  │  (3 nodes)  │  │  Manager   │  │ (3 nodes)  │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  │       │                │                │                │        │   │
│  │       └────────────────┴────────────────┴────────────────┘        │   │
│  │                            LB (HAProxy)                             │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         WORKER NODES                                 │   │
│  │                                                                       │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │  node-01   │  │  node-02   │  │  node-03   │  │  node-04   │ │   │
│  │  │             │  │             │  │             │  │             │ │   │
│  │  │ ┌─────────┐│  │ ┌─────────┐│  │ ┌─────────┐│  │ ┌─────────┐│ │   │
│  │  │ │  Pod    ││  │ │  Pod    ││  │ │  Pod    ││  │ │  Pod    ││ │   │
│  │  │ │ (Auth)  ││  │ │(Biom.)  ││  │ │(Session)││  │ │(Notif.) ││ │   │
│  │  │ └─────────┘│  │ └─────────┘│  │ └─────────┘│  │ └─────────┘│ │   │
│  │  │ ┌─────────┐│  │ ┌─────────┐│  │ ┌─────────┐│  │ ┌─────────┐│ │   │
│  │  │ │  Pod    ││  │ │  Pod    ││  │ │  Pod    ││  │ │  Pod    ││ │   │
│  │  │ │ (Ingress)│ │ │(Analytics││  │ │(Redis)  ││  │ │(Kafka)  ││ │   │
│  │  │ └─────────┘│  │ └─────────┘│  │ └─────────┘│  │ └─────────┘│ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  │                                                                       │   │
│  │  ┌─────────────────────────────────────────────────────────────┐   │   │
│  │  │              STORAGE (Longhorn/Ceph)                        │   │   │
│  │  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐        │   │   │
│  │  │  │ PVC-01   │ │ PVC-02   │ │ PVC-03   │ │ PVC-04   │        │   │   │
│  │  │  │(Postgres)│ │(Redis)   │ │(MinIO)   │ │(Kafka)   │        │   │   │
│  │  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘        │   │   │
│  │  └─────────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3) Cluster Setup

### Kubespray Configuration

```yaml
# inventory/biometrics/group_vars/k8s_cluster/k8s-cluster.yml
---
# Kubernetes version
kube_version: "1.30.0"

# Container runtime
container_manager: containerd

# Network plugin
kube_network_plugin: calico
kube_network_plugin_cidr: 10.244.0.0/16

# Kubernetes API
kube_api_server_bind_address: "0.0.0.0"
kube_api_server_port: 6443
kube_controller_manager_bind_address: 0.0.0.0
kube_scheduler_bind_address: 0.0.0.0

# etcd
etcd_deployment_type: kubeadm

# Worker nodes
kube_node_labels:
  node-role.kubernetes.io/worker: ""

# Addons
kube_proxy_mode: iptables
kube_reserved: true
kube_reserved_cgroups: /system.slice/kubelet.service
system_reserved: true
system_reserved_cgroups: /system.slice

# Autoscaling
enable_clusterAutoscaler: true
clusterAutoscaler_expander: priority

# Metrics
kube_metrics_server_enabled: true

# Pod security
podsecuritypolicy_enabled: true
```

---

## 4) Node Pools

### Node Configuration

```yaml
# k8s/nodepools.yaml
apiVersion: v1
kind: NodePool
metadata:
  name: general-purpose
  namespace: biometrics
spec:
  minSize: 3
  maxSize: 10
  scalingType: Ok
  template:
    spec:
      machineType: n2-standard-4
      labels:
        workload: general-purpose
        tier: application
      taints:
        - key: dedicated
          value: application
          effect: NoSchedule

---
apiVersion: v1
kind: NodePool
metadata:
  name: memory-optimized
  namespace: biometrics
spec:
  minSize: 2
  maxSize: 5
  scalingType: Ok
  template:
    spec:
      machineType: n2-highmem-8
      labels:
        workload: database
        tier: data
      taints:
        - key: dedicated
          value: database
          effect: NoSchedule

---
apiVersion: v1
kind: NodePool
metadata:
  name: gpu-accelerated
  namespace: biometrics
spec:
  minSize: 0
  maxSize: 3
  scalingType: Ok
  template:
    spec:
      machineType: n1-standard-8
      accelerators:
        - type: nvidia-tesla-v100
          count: 1
      labels:
        workload: ml
        tier: inference
      taints:
        - key: nvidia.com/gpu
          value: present
          effect: NoSchedule
```

---

## 5) Namespace Configuration

### Namespaces

```yaml
# k8s/namespaces.yaml
---
apiVersion: v1
kind: Namespace
metadata:
  name: biometrics
  labels:
    istio-injection: enabled
    monitoring: enabled
---
apiVersion: v1
kind: Namespace
metadata:
  name: biometrics-prod
  labels:
    istio-injection: enabled
    monitoring: enabled
    environment: production
---
apiVersion: v1
kind: Namespace
metadata:
  name: biometrics-staging
  labels:
    istio-injection: enabled
    monitoring: enabled
    environment: staging
---
apiVersion: v1
kind: Namespace
metadata:
  name: biometrics-dev
  labels:
    istio-injection: enabled
    environment: development
---
apiVersion: v1
kind: ResourceQuota
metadata:
  name: biometrics-quota
  namespace: biometrics
spec:
  hard:
    requests.cpu: "32"
    requests.memory: 64Gi
    limits.cpu: "64"
    limits.memory: 128Gi
    pods: "100"
    services: "20"
    secrets: "50"
    configmaps: "20"
---
apiVersion: v1
kind: LimitRange
metadata:
  name: biometrics-limits
  namespace: biometrics
spec:
  limits:
    - max:
        cpu: "8"
        memory: "16Gi"
      min:
        cpu: "100m"
        memory: "128Mi"
      default:
        cpu: "1"
        memory: "1Gi"
      defaultRequest:
        cpu: "200m"
        memory: "256Mi"
      type: Container
```

---

## 6) Deployments

### Service Deployment

```yaml
# k8s/deployments/auth-service.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
  namespace: biometrics
  labels:
    app: auth-service
    version: v1
spec:
  replicas: 3
  selector:
    matchLabels:
      app: auth-service
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
  template:
    metadata:
      labels:
        app: auth-service
        version: v1
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9090"
        prometheus.io/path: "/metrics"
    spec:
      serviceAccountName: auth-service
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        fsGroup: 1000
      containers:
        - name: auth-service
          image: biometrics/auth-service:v1.2.3
          imagePullPolicy: Always
          ports:
            - name: http
              containerPort: 3000
              protocol: TCP
            - name: grpc
              containerPort: 3001
              protocol: TCP
            - name: metrics
              containerPort: 9090
              protocol: TCP
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: auth-secrets
                  key: database-url
            - name: JWT_SECRET
              valueFrom:
                secretKeyRef:
                  name: auth-secrets
                  key: jwt-secret
            - name: REDIS_URL
              valueFrom:
                configMapKeyRef:
                  name: auth-config
                  key: redis-url
          resources:
            requests:
              cpu: "250m"
              memory: "512Mi"
            limits:
              cpu: "1000m"
              memory: "2Gi"
          livenessProbe:
            httpGet:
              path: /health
              port: 3000
            initialDelaySeconds: 30
            periodSeconds: 10
            timeoutSeconds: 3
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /ready
              port: 3000
            initialDelaySeconds: 5
            periodSeconds: 5
            timeoutSeconds: 2
            failureThreshold: 3
          startupProbe:
            httpGet:
              path: /health
              port: 3000
            initialDelaySeconds: 0
            periodSeconds: 5
            timeoutSeconds: 3
            failureThreshold: 30
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
            - weight: 100
              podAffinityTerm:
                labelSelector:
                  matchExpressions:
                    - key: app
                      operator: In
                      values:
                        - auth-service
                topologyKey: kubernetes.io/hostname
      tolerations:
        - key: "dedicated"
          operator: "Equal"
          value: "application"
          effect: "NoSchedule"
```

---

## 7) Services

### Service Definitions

```yaml
# k8s/services/auth-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: auth-service
  namespace: biometrics
  labels:
    app: auth-service
spec:
  type: ClusterIP
  ports:
    - name: http
      port: 80
      targetPort: 3000
      protocol: TCP
    - name: grpc
      port: 3001
      targetPort: 3001
      protocol: TCP
  selector:
    app: auth-service

---
apiVersion: v1
kind: Service
metadata:
  name: auth-service-headless
  namespace: biometrics
  labels:
    app: auth-service
spec:
  type: ClusterIP
  clusterIP: None
  ports:
    - name: http
      port: 3000
      targetPort: 3000
      protocol: TCP
  selector:
    app: auth-service

---
# External service for Supabase
apiVersion: v1
kind: Service
metadata:
  name: supabase-external
  namespace: biometrics
spec:
  type: ExternalName
  externalName: db.supabase.project
  ports:
    - port: 5432
      targetPort: 5432
      name: postgres
```

---

## 8) ConfigMaps & Secrets

### Configuration

```yaml
# k8s/configmaps/auth-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: auth-config
  namespace: biometrics
data:
  REDIS_URL: "redis://redis-cluster:6379"
  KAFKA_BROKERS: "kafka-cluster:9092"
  JAEGER_ENDPOINT: "http://jaeger-collector:14268/api/traces"
  LOG_LEVEL: "info"
  RATE_LIMIT_PER_MINUTE: "1000"
  SESSION_TIMEOUT_MINUTES: "30"

---
# k8s/secrets/auth-secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: auth-secrets
  namespace: biometrics
type: Opaque
stringData:
  database-url: "postgres://user:password@postgres:5432/auth"
  jwt-secret: "your-super-secret-jwt-key-change-in-production"
  encryption-key: "your-encryption-key-change-in-production"
  stripe-api-key: "sk_test_xxx"
```

---

## 9) Ingress Configuration

### Ingress Controller

```yaml
# k8s/ingress/nginx-ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: biometrics-ingress
  namespace: biometrics
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/proxy-body-size: "50m"
    nginx.ingress.kubernetes.io/proxy-connect-timeout: "30"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "300"
    nginx.ingress.kubernetes.io/proxy-send-timeout: "300"
    nginx.ingress.kubernetes.io/rate-limit: "100"
    nginx.ingress.kubernetes.io/rate-limit-window: "1m"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/affinity: "cookie"
    nginx.ingress.kubernetes.io/session-cookie-name: "route"
    nginx.ingress.kubernetes.io/session-cookie-hash: "sha1"
spec:
  tls:
    - hosts:
        - biometrics.app
        - api.biometrics.app
        - admin.biometrics.app
      secretName: biometrics-tls
  rules:
    - host: biometrics.app
      http:
        paths:
          - path: /api
            pathType: Prefix
            backend:
              service:
                name: api-gateway
                port:
                  number: 80
          - path: /
            pathType: Prefix
            backend:
              service:
                name: web-frontend
                port:
                  number: 80
    - host: api.biometrics.app
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: api-gateway
                port:
                  number: 80
```

---

## 10) Persistent Storage

### PVC Definitions

```yaml
# k8s/storage/postgres-pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-data
  namespace: biometrics
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: fast-storage
  resources:
    requests:
      storage: 100Gi

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: redis-data
  namespace: biometrics
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: fast-storage
  resources:
    requests:
      storage: 10Gi

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: minio-data
  namespace: biometrics
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: standard-storage
  resources:
    requests:
      storage: 500Gi
```

---

## 11) Horizontal Pod Autoscaler

### HPA Configuration

```yaml
# k8s/autoscaling/auth-hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: auth-service-hpa
  namespace: biometrics
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: auth-service
  minReplicas: 3
  maxReplicas: 20
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
    - type: Pods
      pods:
        metric:
          name: http_requests_per_second
        target:
          type: AverageValue
          averageValue: "1000"
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
        - type: Percent
          value: 10
          periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
        - type: Percent
          value: 100
          periodSeconds: 15
        - type: Pods
          value: 4
          periodSeconds: 15
      selectPolicy: Max
```

---

## 12) Pod Disruption Budgets

### PDB Configuration

```yaml
# k8s/pdb/auth-pdb.yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: auth-service-pdb
  namespace: biometrics
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: auth-service

---
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: biometric-service-pdb
  namespace: biometrics
spec:
  maxUnavailable: 1
  selector:
    matchLabels:
      app: biometric-service
```

---

## 13) Network Policies

### Network Security

```yaml
# k8s/network-policies/default-deny.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
  namespace: biometrics
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress

---
# k8s/network-policies/auth-network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: auth-service-network-policy
  namespace: biometrics
spec:
  podSelector:
    matchLabels:
      app: auth-service
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: api-gateway
        - podSelector:
            matchLabels:
              app: ingress-controller
      ports:
        - protocol: TCP
          port: 3000
  egress:
    - to:
        - podSelector:
            matchLabels:
              app: postgres
      ports:
        - protocol: TCP
          port: 5432
    - to:
        - podSelector:
            matchLabels:
              app: redis
      ports:
        - protocol: TCP
          port: 6379
    - to:
        - namespaceSelector: {}
          podSelector:
            matchLabels:
              k8s-app: kube-dns
      ports:
        - protocol: UDP
          port: 53
```

---

## 14) RBAC Configuration

### Role-Based Access

```yaml
# k8s/rbac/service-account.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: auth-service
  namespace: biometrics
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: auth-service-role
  namespace: biometrics
rules:
  - apiGroups: [""]
    resources: ["configmaps", "secrets"]
    resourceNames: ["auth-config", "auth-secrets"]
    verbs: ["get", "watch", "list"]
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: auth-service-rolebinding
  namespace: biometrics
subjects:
  - kind: ServiceAccount
    name: auth-service
    namespace: biometrics
roleRef:
  kind: Role
  name: auth-service-role
  apiGroup: rbac.authorization.k8s.io

---
# Cluster-wide RBAC for metrics
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: prometheus-clusterrolebinding
subjects:
  - kind: ServiceAccount
    name: prometheus
    namespace: monitoring
roleRef:
  kind: ClusterRole
  name: prometheus
  apiGroup: rbac.authorization.k8s.io
```

---

## 15) GitOps with ArgoCD

### Application Definition

```yaml
# argocd/auth-service-app.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: auth-service
  namespace: argocd
spec:
  project: biometrics
  source:
    repoURL: https://github.com/biometrics/k8s-manifests
    targetRevision: main
    path: services/auth-service
    directory:
      recurse: true
      jsonnet:
        tlas:
          - name: imageTag
            value: v1.2.3
  destination:
    server: https://kubernetes.default.svc
    namespace: biometrics
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
      allowEmpty: false
    syncOptions:
      - CreateNamespace=true
      - PrunePropagationPolicy=foreground
      - PruneLast=true
    retry:
      limit: 5
      backoff:
        duration: 5s
        factor: 2
        maxDuration: 3m
```

---

## 16) Monitoring & Alerting

### Prometheus Rules

```yaml
# k8s/monitoring/prometheus-rules.yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: biometrics-alerts
  namespace: monitoring
spec:
  groups:
    - name: biometrics.rules
      rules:
        - alert: HighErrorRate
          expr: |
            sum(rate(istio_requests_total{namespace="biometrics",response_code=~"5.."}[5m]))
            /
            sum(rate(istio_requests_total{namespace="biometrics"}[5m])) > 0.05
          for: 5m
          labels:
            severity: critical
          annotations:
            summary: "High error rate detected"
            description: "Error rate is {{ $value | humanizePercentage }}"
        - alert: HighLatency
          expr: |
            histogram_quantile(0.99,
              sum(rate(istio_request_duration_milliseconds_bucket{namespace="biometrics"}[5m]))
              by (le, service)) > 1000
          for: 5m
          labels:
            severity: warning
          annotations:
            summary: "High latency detected"
            description: "P99 latency is {{ $value }}ms for {{ $labels.service }}"
        - alert: PodNotReady
          expr: |
            kube_pod_status_phase{namespace="biometrics",phase="Ready"} == 0
          for: 10m
          labels:
            severity: warning
          annotations:
            summary: "Pod not ready"
            description: "Pod {{ $labels.pod }} in {{ $labels.namespace }} is not ready"
```

---

Status: APPROVED  
Version: 1.0  
Last Updated: Februar 2026
