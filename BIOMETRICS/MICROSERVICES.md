# MICROSERVICES.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Architekturentscheidungen folgen den globalen Regeln aus `AGENTS-GLOBAL.md`.
- Jede Strukturänderung benötigt Mapping- und Integrationsabgleich.
- Security-by-Default, NLM-First und Nachweisbarkeit sind Pflicht.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

---

## 1) Overview

This document describes the microservices architecture for BIOMETRICS. The system is decomposed into independent, loosely coupled services that communicate via well-defined APIs and messaging patterns.

### Core Services
- **Auth Service**: Authentication and authorization
- **User Service**: User management
- **Biometric Service**: Biometric operations
- **Session Service**: Session management
- **Notification Service**: Push and email notifications
- **Analytics Service**: Metrics and reporting
- **API Gateway**: Request routing and composition

---

## 2) Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          BIOMETRICS MICROSERVICES                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐ │
│  │    Web      │    │   Mobile    │    │    API      │    │   WebSocket │ │
│  │    App      │    │    App      │    │   Gateway   │    │   Gateway   │ │
│  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘    └──────┬──────┘ │
│         │                  │                  │                  │         │
│         └──────────────────┴────────┬─────────┴──────────────────┘         │
│                                     │                                        │
│                            ┌────────▼────────┐                              │
│                            │   API Gateway    │                              │
│                            │    (Kong/Nginx)  │                              │
│                            └────────┬────────┘                              │
│                                     │                                        │
│         ┌──────────────────────────┼──────────────────────────┐             │
│         │                          │                          │             │
│  ┌──────▼──────┐          ┌───────▼───────┐          ┌───────▼───────┐    │
│  │    Auth     │          │    Biometric   │          │   Analytics   │    │
│  │   Service   │          │    Service     │          │    Service    │    │
│  └──────┬──────┘          └───────┬───────┘          └───────┬───────┘    │
│         │                          │                          │             │
│         │              ┌───────────┼───────────┐              │             │
│         │              │           │           │              │             │
│  ┌──────▼──────┐ ┌─────▼────┐ ┌───▼────┐ ┌────▼─────┐ ┌─────▼─────┐     │
│  │  PostgreSQL │ │ Redis    │ │ MinIO  │ │ Kafka    │ │Elasticsearch│     │
│  │   (Users)   │ │ (Cache)  │ │(Files) │ │ (Events) │ │  (Logs)    │     │
│  └─────────────┘ └──────────┘ └────────┘ └──────────┘ └───────────┘     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3) Service Definitions

### Auth Service

```yaml
# services/auth-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: auth-service
  labels:
    app: biometrics
    service: auth
spec:
  selector:
    app: biometrics
    service: auth
  ports:
    - name: http
      port: 3000
      targetPort: 3000
  replicas: 3
  resources:
    requests:
      memory: "256Mi"
      cpu: "250m"
    limits:
      memory: "512Mi"
      cpu: "500m"
```

```go
// services/auth/main.go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    jwtSecret  []byte
    userRepo   UserRepository
    tokenRepo  TokenRepository
}

func (s *AuthService) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{Error: err.Error()})
        return
    }

    user, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        c.JSON(401, ErrorResponse{Error: "invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        c.JSON(401, ErrorResponse{Error: "invalid credentials"})
        return
    }

    accessToken, err := s.generateToken(user.ID, "access")
    if err != nil {
        c.JSON(500, ErrorResponse{Error: "token generation failed"})
        return
    }

    refreshToken, err := s.generateToken(user.ID, "refresh")
    if err != nil {
        c.JSON(500, ErrorResponse{Error: "token generation failed"})
        return
    }

    // Store refresh token
    s.tokenRepo.Store(user.ID, refreshToken)

    c.JSON(200, LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    900,
        User:         user,
    })
}

func (s *AuthService) generateToken(userID string, tokenType string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "type":    tokenType,
        "exp":     time.Now().Add(15 * time.Minute).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.jwtSecret)
}
```

### Biometric Service

```go
// services/biometric/main.go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/neo4j/neo4j-go-driver/v5"
)

type BiometricService struct {
    neo4jDriver neo4j.Driver
    cache       *redis.Client
    kafka       *kafka.Producer
}

type BiometricRecord struct {
    ID           string                 `json:"id"`
    UserID       string                 `json:"userId"`
    Type         string                 `json:"type"`
    Template     string                 `json:"template"`
    QualityScore float64                `json:"qualityScore"`
    Metadata     map[string]interface{} `json:"metadata"`
    EnrolledAt   time.Time              `json:"enrolledAt"`
    LastUsed     *time.Time             `json:"lastUsed"`
    IsActive     bool                   `json:"isActive"`
}

func (s *BiometricService) EnrollBiometric(c *gin.Context) {
    var req EnrollRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{Error: err.Error()})
        return
    }

    userID := c.GetString("user_id")

    // Store in Neo4j
    session := s.neo4jDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
    defer session.Close()

    _, err := session.Run(`
        CREATE (r:BiometricRecord {
            id: $id,
            userId: $userId,
            type: $type,
            template: $template,
            qualityScore: $qualityScore,
            enrolledAt: datetime(),
            isActive: true
        })
        RETURN r
    `, map[string]interface{}{
        "id":           uuid.New().String(),
        "userId":       userID,
        "type":         req.Type,
        "template":     req.Template,
        "qualityScore": req.QualityScore,
    })

    if err != nil {
        c.JSON(500, ErrorResponse{Error: "enrollment failed"})
        return
    }

    // Publish event to Kafka
    s.kafka.Produce("biometric.enrolled", map[string]interface{}{
        "userId":   userID,
        "type":     req.Type,
        "timestamp": time.Now().Unix(),
    })

    c.JSON(201, EnrollResponse{Success: true})
}

func (s *BiometricService) VerifyBiometric(c *gin.Context) {
    var req VerifyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{Error: err.Error()})
        return
    }

    // Retrieve template from Neo4j
    // Verify against provided data
    // Return result with confidence score

    c.JSON(200, VerifyResponse{
        Success:   true,
        Confidence: 0.95,
    })
}
```

---

## 4) API Gateway

### Kong Configuration

```yaml
# kong.yaml
_format_version: "3.0"
_services:
  - name: auth-service
    url: http://auth-service:3000
    routes:
      - name: auth-login
        paths:
          - /api/v1/auth/login
        methods:
          - POST
        strip_path: true
      - name: auth-register
        paths:
          - /api/v1/auth/register
        methods:
          - POST
        strip_path: true
      - name: auth-refresh
        paths:
          - /api/v1/auth/refresh
        methods:
          - POST
        strip_path: true
    plugins:
      - name: jwt
        config:
          key_claim_name: kid
      - name: rate-limiting
        config:
          minute: 100
          policy: local

  - name: biometric-service
    url: http://biometric-service:3000
    routes:
      - name: biometric-enroll
        paths:
          - /api/v1/biometric/enroll
        methods:
          - POST
        strip_path: true
      - name: biometric-verify
        paths:
          - /api/v1/biometric/verify
        methods:
          - POST
        strip_path: true
        plugins:
          - name: jwt
    plugins:
      - name: cors
      - name: rate-limiting
        config:
          minute: 1000

  - name: analytics-service
    url: http://analytics-service:3000
    routes:
      - name: analytics-metrics
        paths:
          - /api/v1/analytics
        methods:
          - GET
        plugins:
          - name: jwt

_consumers:
  - username: api-consumer
    plugins:
      - name: rate-limiting
        config:
          minute: 10000

_jwt_secrets:
  - consumer: api-consumer
    key: biometrics-api-key
    algorithm: RS256
    rsa_public_key: |
      -----BEGIN PUBLIC KEY-----
      ...
      -----END PUBLIC KEY-----
```

---

## 5) Service Communication

### gRPC Definition

```protobuf
// proto/biometric.proto
syntax = "proto3";

package biometric;

service BiometricService {
    rpc EnrollBiometric(EnrollRequest) returns (EnrollResponse);
    rpc VerifyBiometric(VerifyRequest) returns (VerifyResponse);
    rpc DeleteBiometric(DeleteRequest) returns (DeleteResponse);
    rpc GetBiometricRecord(GetRequest) returns (BiometricRecord);
    rpc ListBiometricRecords(ListRequest) returns (stream BiometricRecord);
    rpc StreamVerification(stream VerifyFrame) returns (stream VerificationResult);
}

message EnrollRequest {
    string user_id = 1;
    BiometricType type = 2;
    string template = 3;
    map<string, string> metadata = 4;
}

message EnrollResponse {
    string record_id = 1;
    bool success = 2;
    double quality_score = 3;
}

message VerifyRequest {
    string record_id = 1;
    string verification_data = 2;
}

message VerifyResponse {
    bool success = 1;
    double confidence = 2;
    string error = 3;
}

message DeleteRequest {
    string record_id = 1;
}

message DeleteResponse {
    bool success = 1;
}

message GetRequest {
    string record_id = 1;
}

message BiometricRecord {
    string id = 1;
    string user_id = 2;
    BiometricType type = 3;
    string template = 4;
    double quality_score = 5;
    int64 enrolled_at = 6;
    int64 last_used = 7;
    bool is_active = 8;
}

message ListRequest {
    string user_id = 1;
    int32 limit = 2;
    int32 offset = 3;
}

enum BiometricType {
    FINGERPRINT = 0;
    FACE = 1;
    IRIS = 2;
    VOICE = 3;
    MULTI = 4;
}
```

---

## 6) Event-Driven Architecture

### Kafka Producers

```go
// internal/kafka/producer.go
package kafka

import (
    "context"
    "encoding/json"
    "time"
)

type Producer struct {
    producer *kafka.Producer
    topics   map[string]string
}

func NewProducer(brokers []string) *Producer {
    config := &kafka.ConfigMap{
        "bootstrap.servers":   strings.Join(brokers, ","),
        "acks":                 "all",
        "retries":              3,
        "retry.backoff.ms":     100,
        "enable.idempotence":   true,
    }

    producer, err := kafka.NewProducer(config)
    if err != nil {
        panic(err)
    }

    return &Producer{
        producer: producer,
        topics: map[string]string{
            "biometric.enrolled": "biometric.events",
            "biometric.verified": "biometric.events",
            "session.created":   "session.events",
            "session.ended":      "session.events",
            "user.registered":    "user.events",
            "notification.send":  "notification.events",
        },
    }
}

func (p *Producer) Publish(ctx context.Context, eventType string, data interface{}) error {
    event := Event{
        Type:      eventType,
        Data:      data,
        Timestamp: time.Now().Unix(),
    }

    payload, err := json.Marshal(event)
    if err != nil {
        return err
    }

    topic := p.topics[eventType]
    if topic == "" {
        topic = "default.events"
    }

    return p.producer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Key:            []byte(eventType),
        Value:          payload,
    }, nil)
}

func (p *Producer) Close() {
    p.producer.Flush(5000)
    p.producer.Close()
}
```

### Kafka Consumers

```go
// internal/kafka/consumer.go
package kafka

type Consumer struct {
    consumer *kafka.Consumer
    handler  MessageHandler
}

type MessageHandler interface {
    HandleBiometricEnrolled(data []byte) error
    HandleBiometricVerified(data []byte) error
    HandleSessionCreated(data []byte) error
    HandleNotificationSend(data []byte) error
}

func NewConsumer(brokers []string, groupID string, handler MessageHandler) *Consumer {
    config := &kafka.ConfigMap{
        "bootstrap.servers":  strings.Join(brokers, ","),
        "group.id":           groupID,
        "auto.offset.reset":  "earliest",
        "enable.auto.commit": false,
    }

    consumer, err := kafka.NewConsumer(config)
    if err != nil {
        panic(err)
    }

    return &Consumer{
        consumer: consumer,
        handler:  handler,
    }
}

func (c *Consumer) Start(ctx context.Context) error {
    topics := []string{
        "biometric.events",
        "session.events",
        "user.events",
        "notification.events",
    }

    if err := c.consumer.SubscribeTopics(topics, nil); err != nil {
        return err
    }

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            msg, err := c.consumer.ReadMessage(5 * time.Second)
            if err != nil {
                continue
            }

            c.processMessage(msg)
            c.consumer.CommitMessage(msg)
        }
    }
}

func (c *Consumer) processMessage(msg *kafka.Message) {
    var event Event
    if err := json.Unmarshal(msg.Value, &event); err != nil {
        return
    }

    switch event.Type {
    case "biometric.enrolled":
        c.handler.HandleBiometricEnrolled(msg.Value)
    case "biometric.verified":
        c.handler.HandleBiometricVerified(msg.Value)
    case "session.created":
        c.handler.HandleSessionCreated(msg.Value)
    case "notification.send":
        c.handler.HandleNotificationSend(msg.Value)
    }
}
```

---

## 7) Service Mesh

### Istio Configuration

```yaml
# istio/virtual-service-biometric.yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: biometric-service
  namespace: biometrics
spec:
  hosts:
    - biometric-service
  http:
    - match:
        - uri:
            prefix: /api/v1/biometric
      route:
        - destination:
            host: biometric-service
            port:
              number: 3000
      retries:
        attempts: 3
        perTryTimeout: 2s
        retryOn: gateway-error,connect-failure,refused-stream
      timeout: 10s
      corsPolicy:
        allowOrigins:
          - exact: "https://biometrics.app"
          - exact: "https://admin.biometrics.app"
        allowMethods:
          - GET
          - POST
          - DELETE
        allowHeaders:
          - Authorization
          - Content-Type
        exposeHeaders:
          - X-Request-ID
        maxAge: 86400s

---
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: biometric-service
  namespace: biometrics
spec:
  host: biometric-service
  trafficPolicy:
    connectionPool:
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
```

---

## 8) Distributed Tracing

### OpenTelemetry

```go
// internal/telemetry/tracer.go
package telemetry

import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/metric"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/otel/sdk/trace/semconv"
)

func InitTracer(serviceName string) (*otel.TracerProvider, error) {
    ctx := context.Background()

    // Create resource
    res, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceNameKey.String(serviceName),
            semconv.ServiceVersionKey.String("1.0.0"),
        ),
    )
    if err != nil {
        return nil, err
    }

    // OTLP exporter
    traceExporter, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint("otel-collector:4317"),
        otlptracegrpc.WithInsecure(),
    )
    if err != nil {
        return nil, err
    }

    // Jaeger exporter
    je, err := jaeger.New(jaeger.WithAgentEndpoint("jaeger:6831"))
    if err != nil {
        return nil, err
    }

    tp := trace.NewTracerProvider(
        trace.WithBatcher(traceExporter),
        trace.WithBatcher(je),
        trace.WithResource(res),
    )

    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    ))

    return tp, nil
}

func InitMeter(serviceName string) (*metric.MeterProvider, error) {
    exporter, err := otlpmetricgrpc.New(context.Background())
    if err != nil {
        return nil, err
    }

    mp := metric.NewMeterProvider(
        metric.WithReader(metric.NewPeriodicReader(exporter)),
    )

    return mp, nil
}
```

---

## 9) Service Discovery

### Consul Integration

```go
// internal/service-discovery/consul.go
package discovery

import (
    "github.com/hashicorp/consul/api"
)

type ConsulRegistry struct {
    client *api.Client
}

func NewConsulRegistry(addr string) (*ConsulRegistry, error) {
    config := api.DefaultConfig()
    config.Address = addr

    client, err := api.NewClient(config)
    if err != nil {
        return nil, err
    }

    return &ConsulRegistry{client: client}, nil
}

func (r *ConsulRegistry) Register(serviceID, serviceName, host string, port int) error {
    registration := &api.AgentServiceRegistration{
        ID:      serviceID,
        Name:    serviceName,
        Address: host,
        Port:    port,
        Check: &api.AgentServiceCheck{
            HTTP:     "http://" + host + ":" + strconv.Itoa(port) + "/health",
            Interval: "10s",
            Timeout:  "5s",
        },
    }

    return r.client.Agent().ServiceRegister(registration)
}

func (r *ConsulRegistry) Deregister(serviceID string) error {
    return r.client.Agent().ServiceDeregister(serviceID)
}

func (r *ConsulRegistry) Discover(serviceName string) ([]*api.ServiceEntry, error) {
    return r.client.Health().Service(serviceName, "", true, nil)
}
```

---

## 10) Circuit Breaker

### Resilience Pattern

```go
// internal/resilience/circuitbreaker.go
package resilience

import (
    "time"
)

type CircuitBreaker struct {
    failures      int
    maxFailures   int
    timeout       time.Duration
    state         State
    lastFailure   time.Time
    onStateChange func(State)
}

type State int

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        maxFailures: maxFailures,
        timeout:     timeout,
        state:       StateClosed,
    }
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    if cb.state == StateOpen {
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = StateHalfOpen
        } else {
            return ErrCircuitOpen
        }
    }

    err := fn()

    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()

        if cb.failures >= cb.maxFailures {
            cb.state = StateOpen
            if cb.onStateChange != nil {
                cb.onStateChange(cb.state)
            }
        }
    } else {
        cb.failures = 0
        if cb.state == StateHalfOpen {
            cb.state = StateClosed
            if cb.onStateChange != nil {
                cb.onStateChange(cb.state)
            }
        }
    }

    return err
}
```

---

## 11) Configuration Management

### Consul KV

```go
// internal/config/consul.go
package config

import (
    "github.com/hashicorp/consul/api"
)

type Config struct {
    consul *api.KV
}

func NewConfig(consulAddr string) (*Config, error) {
    client, err := api.NewClient(&api.Config{Address: consulAddr})
    if err != nil {
        return nil, err
    }

    return &Config{consul: client.KV()}, nil
}

func (c *Config) Get(key string) (string, error) {
    pair, _, err := c.consul.Get(key, nil)
    if err != nil {
        return "", err
    }

    if pair == nil {
        return "", nil
    }

    return string(pair.Value), nil
}

func (c *Config) Set(key, value string) error {
    _, err := c.consul.Put(&api.KVPair{
        Key:   key,
        Value: []byte(value),
    }, nil)

    return err
}

func (c *Config) GetPrefix(prefix string) (map[string]string, error) {
    pairs, _, err := c.consul.List(prefix, nil)
    if err != nil {
        return nil, err
    }

    result := make(map[string]string)
    for _, pair := range pairs {
        result[pair.Key] = string(pair.Value)
    }

    return result, nil
}
```

---

## 12) Health Checks

### Service Health

```go
// internal/health/health.go
package health

import (
    "context"
    "database/sql"
    "time"
)

type HealthChecker struct {
    checks []Check
}

type Check func(ctx context.Context) error

func NewHealthChecker() *HealthChecker {
    return &HealthChecker{
        checks: make([]Check, 0),
    }
}

func (hc *HealthChecker) AddCheck(name string, check Check) {
    hc.checks = append(hc.checks, check)
}

func (hc *HealthChecker) Check(ctx context.Context) HealthReport {
    report := HealthReport{
        Status:    StatusHealthy,
        Timestamp: time.Now(),
        Components: make(map[string]ComponentHealth),
    }

    for _, check := range hc.checks {
        start := time.Now()
        err := check(ctx)
        duration := time.Since(start)

        component := ComponentHealth{
            Status:    StatusHealthy,
            LatencyMs: duration.Milliseconds(),
        }

        if err != nil {
            component.Status = StatusUnhealthy
            component.Error = err.Error()
            report.Status = StatusUnhealthy
        }

        report.Components[name] = component
    }

    return report
}

type HealthReport struct {
    Status    string                     `json:"status"`
    Timestamp time.Time                  `json:"timestamp"`
    Components map[string]ComponentHealth `json:"components"`
}

type ComponentHealth struct {
    Status    string `json:"status"`
    LatencyMs int64  `json:"latencyMs"`
    Error     string `json:"error,omitempty"`
}
```

---

## 13) Testing Strategy

### Integration Tests

```go
// services/biometric/integration_test.go
package biometric_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/testcontainers/testcontainers-go"
)

func TestBiometricEnroll(t *testing.T) {
    // Start testcontainers
    ctx := context.Background()
    
    neo4j, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: testcontainers.ContainerRequest{
            Image: "neo4j:5",
            Env: map[string]string{
                "NEO4J_AUTH": "neo4j/password",
            },
            ExposedPorts: []string{"7687"},
        },
    })
    assert.NoError(t, err)
    defer neo4j.Terminate(ctx)

    // Run service tests
    // ...
}
```

---

## 14) Deployment

### Kubernetes Manifests

```yaml
# k8s/biometric-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: biometric-service
  namespace: biometrics
spec:
  replicas: 3
  selector:
    matchLabels:
      app: biometrics
      service: biometric
  template:
    metadata:
      labels:
        app: biometrics
        service: biometric
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9090"
    spec:
      containers:
        - name: biometric-service
          image: biometrics/biometric-service:latest
          ports:
            - containerPort: 3000
          env:
            - name: NEO4J_URI
              valueFrom:
                configMapKeyRef:
                  name: biometric-config
                  key: neo4j_uri
            - name: REDIS_URL
              valueFrom:
                configMapKeyRef:
                  name: biometric-config
                  key: redis_url
          resources:
            requests:
              memory: "512Mi"
              cpu: "500m"
            limits:
              memory: "1Gi"
              cpu: "1000m"
          livenessProbe:
            httpGet:
              path: /health
              port: 3000
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /ready
              port: 3000
            initialDelaySeconds: 5
            periodSeconds: 5
```

---

## 15) Monitoring & Observability

### Metrics

```go
// internal/metrics/service.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    RequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
        },
        []string{"method", "path", "status"},
    )

    ActiveConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_connections",
        },
    )

    BiometricOperations = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "biometric_operations_total",
        },
        []string{"operation", "status"},
    )
)
```

---

Status: APPROVED  
Version: 1.0  
Last Updated: Februar 2026
