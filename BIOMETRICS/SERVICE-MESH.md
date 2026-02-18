# SERVICE-MESH.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Architekturentscheidungen folgen den globalen Regeln aus `AGENTS-GLOBAL.md`.
- Jede Strukturänderung benötigt Mapping- und Integrationsabgleich.
- Security-by-Default, NLM-First und Nachweisbarkeit sind Pflicht.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

---

## 1) Overview

This document describes the service mesh architecture using Istio for BIOMETRICS. The service mesh provides traffic management, security, and observability for all microservices communication.

### Core Features
- Traffic routing and load balancing
- Service discovery
- Mutual TLS encryption
- Circuit breaking
- Rate limiting
- Distributed tracing
- Metrics collection
- Access control policies

---

## 2) Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          BIOMETRICS SERVICE MESH                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         CONTROL PLANE                                │   │
│  │                    (Istiod - Port 15012)                            │   │
│  │   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐         │   │
│  │   │   CA     │  │   Config │  │  Discovery│  │  Agent   │         │   │
│  │   │  Server  │  │  Server  │  │  Server  │  │  Watcher │         │   │
│  │   └──────────┘  └──────────┘  └──────────┘  └──────────┘         │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         DATA PLANE                                   │   │
│  │                    (Envoy Proxies - Sidecars)                      │   │
│                                                                         │   │
│  ┌───────┐  ┌───────┐  ┌───────┐  ┌───────┐  ┌───────┐              │   │
│  │Envoy  │  │Envoy  │  │Envoy  │  │Envoy  │  │Envoy  │              │   │
│  │Sidecar│  │Sidecar│  │Sidecar│  │Sidecar│  │Sidecar│              │   │
│  └───────┘  └───────┘  └───────┘  └───────┘  └───────┘              │   │
│     │         │         │         │         │                        │   │
│  ┌──▼────┐ ┌──▼────┐ ┌──▼────┐ ┌──▼────┐ ┌──▼────┐                 │   │
│  │ Auth  │ │Biom.  │ │Session│ │Notif. │ │Analytics│                │   │
│  │Svc    │ │Svc    │ │Svc    │ │Svc    │ │Svc    │                 │   │
│  └───────┘ └───────┘ └───────┘ └───────┘ └───────┘                 │   │
│                                                                         │   │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                       GATEWAY                                    │   │
│  │                 (Ingress/Egress Gateway)                        │   │
│  │   ┌──────────┐  ┌──────────┐  ┌──────────┐                     │   │
│  │   │Ingress  │  │  Egress  │  │  Metrics │                     │   │
│  │   │Gateway  │  │ Gateway  │  │Collector │                     │   │
│  │   └──────────┘  └──────────┘  └──────────┘                     │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3) Istio Installation

### Installation YAML

```yaml
# istio-operator.yaml
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  name: biometrics-istio
  namespace: istio-system
spec:
  profile: default
  meshConfig:
    enableAutoMtls: true
    defaultConfig:
      proxyMetadata:
        ISTIO_META_DNS_CAPTURE: "true"
        ISTIO_META_DNS_AUTO_ALLOCATE: "true"
      tracing:
        sampling: 10
        zipkin:
          address: jaeger-collector.istio-system:9411
    localityLbSetting:
      enabled: true
  components:
    ingressGateways:
      - name: istio-ingressgateway
        enabled: true
        k8s:
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 2000m
              memory: 1024Mi
          hpaSpec:
            minReplicas: 2
            maxReplicas: 10
            metrics:
              - type: Resource
                resource:
                  name: cpu
                  targetAverageUtilization: 70
    egressGateways:
      - name: istio-egressgateway
        enabled: true
        k8s:
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
  values:
    global:
      meshID: biometrics-mesh
      multiCluster:
        clusterName: biometrics-cluster
    prometheus:
      enabled: true
      retention: 30d
    kiali:
      enabled: true
      dashboard:
        user: admin
        passphrase: admin
```

---

## 4) Traffic Management

### Virtual Services

```yaml
# istio/virtual-service-auth.yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: auth-service
  namespace: biometrics
spec:
  hosts:
    - auth-service
    - auth.biometrics.app
  gateways:
    - istio-ingressgateway/biometrics-gateway
  http:
    - match:
        - uri:
            prefix: "/api/v1/auth/login"
        - uri:
            prefix: "/api/v1/auth/register"
      route:
        - destination:
            host: auth-service
            port:
              number: 3000
            subset: v1
      retries:
        attempts: 3
        perTryTimeout: 2s
        retryOn: gateway-error,connect-failure,refused-stream
      timeout: 10s
    
    - match:
        - uri:
            prefix: "/api/v1/auth/refresh"
      route:
        - destination:
            host: auth-service
            port:
              number: 3000
            subset: v1
      # No retries for refresh tokens - prevent token reuse
      timeout: 5s
    
    - match:
        - uri:
            prefix: "/api/v1/auth/logout"
      route:
        - destination:
            host: auth-service
            port:
              number: 3000
            subset: v1
      # Immediate response for logout
      timeout: 3s

---
# Destination rules for auth service
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: auth-service
  namespace: biometrics
spec:
  host: auth-service
  trafficPolicy:
    connectionPool:
      http:
        h2UpgradePolicy: UPGRADE
        http1MaxPendingRequests: 100
        http2MaxRequests: 1000
        maxRequestsPerConnection: 100
    loadBalancer:
      simple: LEAST_REQUEST
      localityLbSetting:
        enabled: true
        failover:
          - from: *
            to: *
    outlierDetection:
      consecutive5xxErrors: 5
      interval: 30s
      baseEjectionTime: 30s
      maxEjectionPercent: 50
      minHealthPercent: 30
  subsets:
    - name: v1
      labels:
        version: v1
    - name: v2
      labels:
        version: v2
```

---

## 5) Service Discovery

### Service Entries

```yaml
# istio/service-entry-external.yaml
apiVersion: networking.istio.io/v1beta1
kind: ServiceEntry
metadata:
  name: external-supabase
  namespace: biometrics
spec:
  hosts:
    - db.supabase.project
  location: MESH_INTERNAL
  ports:
    - number: 5432
      name: postgres
      protocol: TCP
  resolution: DNS
  endpoints:
    - address: db.supabase.project
      ports:
        postgres: 5432

---
apiVersion: networking.istio.io/v1beta1
kind: ServiceEntry
metadata:
  name: external-kafka
  namespace: biometrics
spec:
  hosts:
    - kafka-cluster
  location: MESH_INTERNAL
  ports:
    - number: 9092
      name: kafka
      protocol: TCP
  resolution: STATIC
  endpoints:
    - address: 10.0.0.50
      ports:
        kafka: 9092
    - address: 10.0.0.51
      ports:
        kafka: 9092
    - address: 10.0.0.52
      ports:
        kafka: 9092
```

---

## 6) Security

### mTLS Configuration

```yaml
# istio/peer-authentication.yaml
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: biometrics
spec:
  mtls:
    mode: STRICT

---
# Authorization policies
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: auth-service-authorization
  namespace: biometrics
spec:
  selector:
    matchLabels:
      app: auth-service
  action: ALLOW
  rules:
    - from:
        - source:
            principals:
              - "cluster.local/ns/biometrics/sa/api-gateway"
              - "cluster.local/ns/biometrics/sa/web-app"
      to:
        - operation:
            methods: ["POST", "GET"]
            paths: ["/api/v1/auth/*"]

---
# JWT authentication for biometric service
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: biometric-service-authz
  namespace: biometrics
spec:
  selector:
    matchLabels:
      app: biometric-service
  action: ALLOW
  rules:
    - from:
        - source:
            principals:
              - "cluster.local/ns/biometrics/sa/auth-service"
              - "cluster.local/ns/biometrics/sa/api-gateway"
        - source:
            requestMatchers:
              - headers:
                  authorization:
                    prefix: "Bearer "
      to:
        - operation:
            methods: ["GET", "POST"]
            paths: ["/api/v1/biometric/*"]
    - to:
        - operation:
            methods: ["GET"]
            paths: ["/health", "/ready"]
      when:
        - key: request.headers[authorization]
          notPresent: true

---
# Deny all by default
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: deny-all
  namespace: biometrics
spec:
  {}
```

---

## 7) Rate Limiting

### Envoy Filter

```yaml
# istio/rate-limit-filter.yaml
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: rate-limit-filter
  namespace: biometrics
spec:
  workloadSelector:
    labels:
      app: api-gateway
  configPatches:
    - applyTo: HTTP_FILTER
      match:
        context: SIDECAR_INBOUND
        listener:
          filterChain:
            filter:
              name: envoy.filters.network.http_connection_manager
      patch:
        operation: INSERT_BEFORE
        value:
          name: envoy.filters.http.local_ratelimit
          typed_config:
            "@type": type.googleapis.com/udpa.type.v1.TypedStruct
            type_url: type.googleapis.com/envoy.extensions.filters.http.local_ratelimit.v3.LocalRateLimit
            value:
              stat_prefix: http_local_rate_limiter
              token_bucket:
                max_tokens: 10000
                tokens_per_fill: 1000
                fill_interval: 1s
              filter_enabled:
                runtime_key: local_rate_limit_enabled
                default_value:
                  numerator: 100
                  denominator: HUNDRED

---
# Global rate limit
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: global-rate-limit
  namespace: biometrics
spec:
  workloadSelector:
    labels:
      app: api-gateway
  configPatches:
    - applyTo: CLUSTER
      match:
        context: SIDECAR_OUTBOUND
        cluster:
          name: "outbound|8086||rate-limiter.rate-limit.svc.cluster.local"
      patch:
        operation: ADD
        value:
          name: rate_limit_cluster
          type: STRICT_DNS
          lb_policy: ROUND_ROBIN
          hosts:
            - socket_address:
                address: rate-limiter.rate-limit.svc.cluster.local
                port_value: 8086
```

---

## 8) Circuit Breaker

### Outlier Detection

```yaml
# istio/circuit-breaker.yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: biometric-service-cb
  namespace: biometrics
spec:
  host: biometric-service
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 100
      http:
        h2UpgradePolicy: UPGRADE
        http1MaxPendingRequests: 100
        http2MaxRequests: 1000
    loadBalancer:
      simple: LEAST_REQUEST
    outlierDetection:
      consecutive5xxErrors: 5
      interval: 30s
      baseEjectionTime: 30s
      maxEjectionPercent: 50
      minHealthPercent: 30
      consecutiveGatewayErrors: 3
      consecutiveLocalOriginFailures: 3
      percentage:
        value: 10

---
# Retry policy
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: analytics-service-retry
  namespace: biometrics
spec:
  hosts:
    - analytics-service
  http:
    - route:
        - destination:
            host: analytics-service
      retries:
        attempts: 5
        perTryTimeout: 3s
        retryOn: 5xx,reset,connect-failure,refused-stream
        retryRemoteLocalities: true
        baseEjectionTime: 30s
        maxEjectionTime: 300s
        maxRetries: 10
```

---

## 9) Ingress Gateway

### Gateway Configuration

```yaml
# istio/ingress-gateway.yaml
apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:
  name: biometrics-gateway
  namespace: istio-system
spec:
  selector:
    istio: ingressgateway
  servers:
    - port:
        number: 80
        name: http
        protocol: HTTP
      tls:
        httpsRedirect: true
    - port:
        number: 443
        name: https
        protocol: HTTP2
      tls:
        mode: SIMPLE
        credentialName: biometrics-tls
      hosts:
        - "biometrics.app"
        - "api.biometrics.app"
        - "admin.biometrics.app"
    - port:
        number: 8443
        name: grpc
        protocol: GRPC
      tls:
        mode: SIMPLE
        credentialName: biometrics-grpc-tls
      hosts:
        - "grpc.biometrics.app"

---
# TLS configuration
apiVersion: networking.istio.io/v1beta1
kind: TLSRoute
metadata:
  name: biometrics-tls
  namespace: istio-system
spec:
  hosts:
    - biometrics.app
  match:
    - port: 443
  route:
    - destination:
        host: web-app-service
        port:
          number: 80

---
# HTTP redirection
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: http-redirect
  namespace: istio-system
spec:
  hosts:
    - biometrics.app
  gateways:
    - biometrics-gateway
  http:
    - match:
        - uri:
            prefix: "/"
      redirect:
        uri: "/"
        authority: "www.biometrics.app"
        redirectCode: 301
```

---

## 10) Observability

### Telemetry Configuration

```yaml
# istio/telemetry.yaml
apiVersion: telemetry.istio.io/v1alpha1
kind: Telemetry
metadata:
  name: biometrics-telemetry
  namespace: biometrics
spec:
  tracing:
    - providers:
        - name: jaeger
        - name: zipkin
      randomSamplingPercentage: 10
      customTags:
        environment:
          literal:
            value: production
        customer_id:
          header:
            name: x-customer-id
  metrics:
    - providers:
        - name: prometheus
  accessLogging:
    - providers:
        - name: envoy
      filter:
        expression: >
          ConnectionTermination != true &&
          Request_Type == 'original_request'
```

---

## 11) Sidecar Configuration

### Traffic Interception

```yaml
# istio/sidecar.yaml
apiVersion: networking.istio.io/v1beta1
kind: Sidecar
metadata:
  name: biometric-sidecar
  namespace: biometrics
spec:
  workloadSelector:
    labels:
      app: biometric-service
  egress:
    - hosts:
        - "./"
        - "istio-system/*"
        - "biometrics/kafka-cluster"
        - "biometrics/redis-cluster"
        - "biometrics/postgres-cluster"
  ingress:
    - port:
        number: 3000
        protocol: HTTP
        name: grpc
      defaultEndpoint: 127.0.0.1:3000

---
# Egress only for specific services
apiVersion: networking.istio.io/v1beta1
kind: Sidecar
metadata:
  name: analytics-sidecar
  namespace: biometrics
spec:
  workloadSelector:
    labels:
      app: analytics-service
  egress:
    - hosts:
        - "biometrics/biometric-service"
        - "biometrics/session-service"
        - "istio-system/*"
    - port:
        number: 9090
        protocol: HTTP
      hosts:
        - "*/*"
        - "prometheus/*"
```

---

## 12) Request Authentication

### JWT Verification

```yaml
# istio/request-authentication.yaml
apiVersion: security.istio.io/v1beta1
kind: RequestAuthentication
metadata:
  name: jwt-authentication
  namespace: biometrics
spec:
  selector:
    matchLabels:
      app: api-gateway
  jwtRules:
    - issuer: "biometrics.auth"
      audiences:
        - "biometrics-api"
      forwardOriginalToken: true
      # Specify claim to route by (e.g., tenant)
      claimToHeaders:
        - header: x-tenant-id
          claim: tenant_id
      fromHeaders:
        - name: Authorization
          prefix: "Bearer "
      fromParams:
        - access_token
      # JWKS endpoint
      jwksUri: https://auth.biometrics.app/.well-known/jwks.json

---
# Per-route JWT validation
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: protected-service
  namespace: biometrics
spec:
  hosts:
    - protected-service
  http:
    - match:
        - headers:
            authorization:
              present: true
      route:
        - destination:
            host: protected-service
      # Allow authenticated requests
    - match:
        - uri:
            prefix: "/health"
      route:
        - destination:
            host: protected-service
      # Allow unauthenticated health checks
    - route:
        - destination:
            host: protected-service
      # Default: require authentication
```

---

## 13) Network Resilience

### Connection Pooling

```yaml
# istio/connection-pools.yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: session-service-pool
  namespace: biometrics
spec:
  host: session-service
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 200
        connectTimeout: 10s
      http:
        h2UpgradePolicy: UPGRADE
        http1MaxPendingRequests: 200
        http2MaxRequests: 1000
        maxRequestsPerConnection: 50
    loadBalancer:
      simple: LEAST_REQUEST
      warmup: 30s
    outlierDetection:
      consecutive5xxErrors: 10
      interval: 30s
      baseEjectionTime: 30s
      maxEjectionPercent: 80

---
# Database connection pool (TCP level)
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: postgres-pool
  namespace: biometrics
spec:
  host: postgres-cluster
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 50
        connectTimeout: 5s
        tcpKeepalive:
          interval: 30s
          probes: 3
          time: 60s
```

---

## 14) Monitoring Dashboards

### Grafana Integration

```yaml
# istio/monitoring-dashboards.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: istio-grafana-dashboard
  namespace: istio-system
  labels:
    grafana_dashboard: "1"
data:
  biometrics-dashboard.json: |
    {
      "dashboard": {
        "title": "BIOMETRICS Service Mesh",
        "panels": [
          {
            "title": "Request Rate",
            "type": "graph",
            "targets": [
              {
                "expr": "rate(istio_requests_total{service=~\".*\"}[5m])",
                "legendFormat": "{{service}}"
              }
            ]
          },
          {
            "title": "Error Rate",
            "type": "graph",
            "targets": [
              {
                "expr": "rate(istio_requests_total{service=~\".*\",response_code=~\"5..\"}[5m])",
                "legendFormat": "{{service}} - 5xx"
              }
            ]
          },
          {
            "title": "Latency P99",
            "type": "graph",
            "targets": [
              {
                "expr": "histogram_quantile(0.99, rate(istio_request_duration_milliseconds_bucket{service=~\".*\"}[5m]))",
                "legendFormat": "{{service}}"
              }
            ]
          }
        ]
      }
    }
```

---

## 15) Troubleshooting

### Debug Commands

```bash
# Check Envoy proxy status
istioctl proxy-status

# View Envoy config
istioctl proxy-config endpoint <pod-name> -n biometrics

# Check mTLS status
istioctl authz check <pod-name> -n biometrics

# Debug traffic
istioctl x tap <pod-name> -n biometrics

# Check service entries
kubectl get serviceentries -n biometrics

# View virtual services
kubectl get virtualservices -n biometrics -o yaml

# Analyze configuration
istioctl analyze -n biometrics

# Check metrics
istioctl metrics biometrics

# Tail logs
kubectl logs -n istio-system -l istio=ingressgateway -f
```

---

Status: APPROVED  
Version: 1.0  
Last Updated: Februar 2026
