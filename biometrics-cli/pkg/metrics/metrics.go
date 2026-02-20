package metrics

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// MetricsRegistry holds all Prometheus metrics for the biometrics CLI
type MetricsRegistry struct {
	// HTTP Metrics
	httpRequestsTotal    prometheus.Counter
	httpRequestDuration  *prometheus.HistogramVec
	httpRequestsInFlight prometheus.Gauge
	httpErrorsTotal      prometheus.Counter

	// Audit Metrics
	auditEventsTotal  *prometheus.CounterVec
	auditQueriesTotal prometheus.Counter
	auditLogsWritten  prometheus.Counter
	auditStorageBytes prometheus.Gauge

	// Auth Metrics
	authAttemptsTotal  *prometheus.CounterVec
	authFailuresTotal  *prometheus.CounterVec
	authSuccessTotal   prometheus.Counter
	activeSessions     prometheus.Gauge
	tokenRefreshsTotal prometheus.Counter

	// Rate Limit Metrics
	rateLimitRequests *prometheus.CounterVec
	rateLimitLimited  prometheus.Counter
	rateLimitBypassed prometheus.Counter

	// System Metrics
	goRoutines       prometheus.Gauge
	memoryAllocBytes prometheus.Gauge
	memorySysBytes   prometheus.Gauge
	gcPauseNs        prometheus.Histogram

	// Custom metrics registry
	customMetrics map[string]prometheus.Collector
	customMu      sync.RWMutex

	// Logger
	logger *zap.Logger

	// Registry
	reg *prometheus.Registry
}

// NewMetricsRegistry creates a new metrics registry with all standard metrics
func NewMetricsRegistry(logger *zap.Logger) (*MetricsRegistry, error) {
	if logger == nil {
		var err error
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("failed to create logger: %w", err)
		}
	}

	reg := prometheus.NewRegistry()

	m := &MetricsRegistry{
		logger:        logger,
		reg:           reg,
		customMetrics: make(map[string]prometheus.Collector),
	}

	// Initialize HTTP metrics
	m.httpRequestsTotal = promauto.With(reg).NewCounter(
		prometheus.CounterOpts{
			Name: "biometrics_http_requests_total",
			Help: "Total number of HTTP requests",
		},
	)

	m.httpRequestDuration = promauto.With(reg).NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "biometrics_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "endpoint", "status"},
	)

	m.httpRequestsInFlight = promauto.With(reg).NewGauge(
		prometheus.GaugeOpts{
			Name: "biometrics_http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)

	m.httpErrorsTotal = promauto.With(reg).NewCounter(
		prometheus.CounterOpts{
			Name: "biometrics_http_errors_total",
			Help: "Total number of HTTP errors",
		},
	)

	// Initialize Audit metrics
	m.auditEventsTotal = promauto.With(reg).NewCounterVec(
		prometheus.CounterOpts{
			Name: "biometrics_audit_events_total",
			Help: "Total number of audit events by type",
		},
		[]string{"event_type", "severity"},
	)

	m.auditQueriesTotal = promauto.With(reg).NewCounter(
		prometheus.CounterOpts{
			Name: "biometrics_audit_queries_total",
			Help: "Total number of audit log queries",
		},
	)

	m.auditLogsWritten = promauto.With(reg).NewCounter(
		prometheus.CounterOpts{
			Name: "biometrics_audit_logs_written_total",
			Help: "Total number of audit logs written to storage",
		},
	)

	m.auditStorageBytes = promauto.With(reg).NewGauge(
		prometheus.GaugeOpts{
			Name: "biometrics_audit_storage_bytes",
			Help: "Current size of audit log storage in bytes",
		},
	)

	// Initialize Auth metrics
	m.authAttemptsTotal = promauto.With(reg).NewCounterVec(
		prometheus.CounterOpts{
			Name: "biometrics_auth_attempts_total",
			Help: "Total number of authentication attempts by method",
		},
		[]string{"method", "provider"},
	)

	m.authFailuresTotal = promauto.With(reg).NewCounterVec(
		prometheus.CounterOpts{
			Name: "biometrics_auth_failures_total",
			Help: "Total number of authentication failures by reason",
		},
		[]string{"reason", "method"},
	)

	m.authSuccessTotal = promauto.With(reg).NewCounter(
		prometheus.CounterOpts{
			Name: "biometrics_auth_success_total",
			Help: "Total number of successful authentications",
		},
	)

	m.activeSessions = promauto.With(reg).NewGauge(
		prometheus.GaugeOpts{
			Name: "biometrics_active_sessions",
			Help: "Current number of active user sessions",
		},
	)

	m.tokenRefreshsTotal = promauto.With(reg).NewCounter(
		prometheus.CounterOpts{
			Name: "biometrics_token_refreshs_total",
			Help: "Total number of token refresh operations",
		},
	)

	// Initialize Rate Limit metrics
	m.rateLimitRequests = promauto.With(reg).NewCounterVec(
		prometheus.CounterOpts{
			Name: "biometrics_ratelimit_requests_total",
			Help: "Total number of rate-limited requests by endpoint",
		},
		[]string{"endpoint", "client_id"},
	)

	m.rateLimitLimited = promauto.With(reg).NewCounter(
		prometheus.CounterOpts{
			Name: "biometrics_ratelimit_limited_total",
			Help: "Total number of requests that were rate limited",
		},
	)

	m.rateLimitBypassed = promauto.With(reg).NewCounter(
		prometheus.CounterOpts{
			Name: "biometrics_ratelimit_bypassed_total",
			Help: "Total number of requests that bypassed rate limiting",
		},
	)

	// Initialize System metrics
	m.goRoutines = promauto.With(reg).NewGauge(
		prometheus.GaugeOpts{
			Name: "biometrics_go_routines",
			Help: "Current number of goroutines",
		},
	)

	m.memoryAllocBytes = promauto.With(reg).NewGauge(
		prometheus.GaugeOpts{
			Name: "biometrics_memory_alloc_bytes",
			Help: "Current memory allocation in bytes",
		},
	)

	m.memorySysBytes = promauto.With(reg).NewGauge(
		prometheus.GaugeOpts{
			Name: "biometrics_memory_sys_bytes",
			Help: "Total memory in bytes from OS",
		},
	)

	m.gcPauseNs = promauto.With(reg).NewHistogram(
		prometheus.HistogramOpts{
			Name:    "biometrics_gc_pause_nanoseconds",
			Help:    "GC pause time in nanoseconds",
			Buckets: prometheus.ExponentialBuckets(1000, 2, 15),
		},
	)

	logger.Info("Metrics registry initialized successfully")
	return m, nil
}

// HTTPMiddleware returns middleware for request tracking
func (m *MetricsRegistry) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.httpRequestsInFlight.Inc()
		defer m.httpRequestsInFlight.Dec()

		start := time.Now()
		wrapped := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start).Seconds()
		status := wrapped.statusCode

		m.httpRequestsTotal.Inc()
		m.httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, fmt.Sprintf("%d", status)).Observe(duration)

		if status >= 400 {
			m.httpErrorsTotal.Inc()
		}
	})
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// RecordHTTPEvent records an HTTP request metric
func (m *MetricsRegistry) RecordHTTPEvent(method, endpoint string, status int, duration time.Duration) {
	m.httpRequestsTotal.Inc()
	m.httpRequestDuration.WithLabelValues(method, endpoint, fmt.Sprintf("%d", status)).Observe(duration.Seconds())
	if status >= 400 {
		m.httpErrorsTotal.Inc()
	}
}

// RecordAuditEvent records an audit event
func (m *MetricsRegistry) RecordAuditEvent(eventType, severity string) {
	m.auditEventsTotal.WithLabelValues(eventType, severity).Inc()
	m.auditLogsWritten.Inc()
}

// RecordAuditQuery records an audit log query
func (m *MetricsRegistry) RecordAuditQuery() {
	m.auditQueriesTotal.Inc()
}

// UpdateAuditStorage updates the audit storage size metric
func (m *MetricsRegistry) UpdateAuditStorage(bytes int64) {
	m.auditStorageBytes.Set(float64(bytes))
}

// RecordAuthAttempt records an authentication attempt
func (m *MetricsRegistry) RecordAuthAttempt(method, provider string) {
	m.authAttemptsTotal.WithLabelValues(method, provider).Inc()
}

// RecordAuthFailure records an authentication failure
func (m *MetricsRegistry) RecordAuthFailure(reason, method string) {
	m.authFailuresTotal.WithLabelValues(reason, method).Inc()
}

// RecordAuthSuccess records a successful authentication
func (m *MetricsRegistry) RecordAuthSuccess() {
	m.authSuccessTotal.Inc()
}

// UpdateActiveSessions updates the active sessions count
func (m *MetricsRegistry) UpdateActiveSessions(count int) {
	m.activeSessions.Set(float64(count))
}

// RecordTokenRefresh records a token refresh operation
func (m *MetricsRegistry) RecordTokenRefresh() {
	m.tokenRefreshsTotal.Inc()
}

// RecordRateLimitRequest records a rate-limited request
func (m *MetricsRegistry) RecordRateLimitRequest(endpoint, clientID string, limited bool) {
	m.rateLimitRequests.WithLabelValues(endpoint, clientID).Inc()
	if limited {
		m.rateLimitLimited.Inc()
	} else {
		m.rateLimitBypassed.Inc()
	}
}

// UpdateSystemMetrics updates system-level metrics
func (m *MetricsRegistry) UpdateSystemMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.goRoutines.Set(float64(runtime.NumGoroutine()))
	m.memoryAllocBytes.Set(float64(memStats.Alloc))
	m.memorySysBytes.Set(float64(memStats.Sys))
	m.gcPauseNs.Observe(float64(memStats.PauseNs[memStats.NumGC%256]))
}

// RegisterCustomMetric registers a custom Prometheus metric
func (m *MetricsRegistry) RegisterCustomMetric(name string, collector prometheus.Collector) error {
	m.customMu.Lock()
	defer m.customMu.Unlock()

	if _, exists := m.customMetrics[name]; exists {
		return fmt.Errorf("custom metric %s already exists", name)
	}

	if err := m.reg.Register(collector); err != nil {
		return fmt.Errorf("failed to register custom metric %s: %w", name, err)
	}

	m.customMetrics[name] = collector
	m.logger.Debug("Custom metric registered", zap.String("name", name))
	return nil
}

// UnregisterCustomMetric unregisters a custom Prometheus metric
func (m *MetricsRegistry) UnregisterCustomMetric(name string) error {
	m.customMu.Lock()
	defer m.customMu.Unlock()

	collector, exists := m.customMetrics[name]
	if !exists {
		return fmt.Errorf("custom metric %s not found", name)
	}

	if !m.reg.Unregister(collector) {
		return fmt.Errorf("failed to unregister custom metric %s", name)
	}

	delete(m.customMetrics, name)
	m.logger.Debug("Custom metric unregistered", zap.String("name", name))
	return nil
}

// GetCustomMetric retrieves a custom metric by name
func (m *MetricsRegistry) GetCustomMetric(name string) (prometheus.Collector, error) {
	m.customMu.RLock()
	defer m.customMu.RUnlock()

	collector, exists := m.customMetrics[name]
	if !exists {
		return nil, fmt.Errorf("custom metric %s not found", name)
	}

	return collector, nil
}

// CreateCounter creates and registers a new counter metric
func (m *MetricsRegistry) CreateCounter(name, help string) (prometheus.Counter, error) {
	counter := promauto.With(m.reg).NewCounter(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
	)

	m.customMu.Lock()
	m.customMetrics[name] = counter
	m.customMu.Unlock()

	return counter, nil
}

// CreateGauge creates and registers a new gauge metric
func (m *MetricsRegistry) CreateGauge(name, help string) (prometheus.Gauge, error) {
	gauge := promauto.With(m.reg).NewGauge(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
	)

	m.customMu.Lock()
	m.customMetrics[name] = gauge
	m.customMu.Unlock()

	return gauge, nil
}

// CreateHistogram creates and registers a new histogram metric
func (m *MetricsRegistry) CreateHistogram(name, help string, buckets []float64) (prometheus.Histogram, error) {
	histogram := promauto.With(m.reg).NewHistogram(
		prometheus.HistogramOpts{
			Name:    name,
			Help:    help,
			Buckets: buckets,
		},
	)

	m.customMu.Lock()
	m.customMetrics[name] = histogram
	m.customMu.Unlock()

	return histogram, nil
}

// CreateCounterVec creates and registers a new counter vector metric
func (m *MetricsRegistry) CreateCounterVec(name, help string, labelNames []string) (*prometheus.CounterVec, error) {
	counterVec := promauto.With(m.reg).NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
		labelNames,
	)

	m.customMu.Lock()
	m.customMetrics[name] = counterVec
	m.customMu.Unlock()

	return counterVec, nil
}

// MetricsHandler returns an HTTP handler for Prometheus metrics endpoint
func (m *MetricsRegistry) MetricsHandler() http.Handler {
	return promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
		Registry:          m.reg,
	})
}

// StartMetricsServer starts a dedicated HTTP server for metrics
func (m *MetricsRegistry) StartMetricsServer(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", m.MetricsHandler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:         addr,
		Handler:      m.HTTPMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		m.logger.Info("Starting metrics server", zap.String("addr", addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			m.logger.Error("Metrics server failed", zap.Error(err))
		}
	}()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			m.logger.Error("Metrics server shutdown error", zap.Error(err))
		}
	}()

	return nil
}

// GetRegistry returns the underlying Prometheus registry
func (m *MetricsRegistry) GetRegistry() prometheus.Gatherer {
	return m.reg
}

// CollectAll collects all metrics and returns them as a map
func (m *MetricsRegistry) CollectAll() (map[string]interface{}, error) {
	m.UpdateSystemMetrics()

	metrics := make(map[string]interface{})

	metricFamilies, err := m.reg.Gather()
	if err != nil {
		return nil, fmt.Errorf("failed to gather metrics: %w", err)
	}

	for _, mf := range metricFamilies {
		for _, metric := range mf.Metric {
			var value float64
			if metric.Gauge != nil {
				value = metric.Gauge.GetValue()
			} else if metric.Counter != nil {
				value = metric.Counter.GetValue()
			}

			labelMap := make(map[string]string)
			for _, label := range metric.Label {
				labelMap[label.GetName()] = label.GetValue()
			}

			key := mf.GetName()
			if len(labelMap) > 0 {
				key = fmt.Sprintf("%s_%v", key, labelMap)
			}

			metrics[key] = value
		}
	}

	return metrics, nil
}

// ResetAll resets all metrics (useful for testing)
func (m *MetricsRegistry) ResetAll() error {
	m.customMu.Lock()
	defer m.customMu.Unlock()

	// Unregister all custom metrics
	for name, collector := range m.customMetrics {
		m.reg.Unregister(collector)
		delete(m.customMetrics, name)
	}

	m.logger.Info("All metrics reset")
	return nil
}

// GetMetricsSummary returns a summary of key metrics
func (m *MetricsRegistry) GetMetricsSummary() map[string]interface{} {
	m.UpdateSystemMetrics()

	return map[string]interface{}{
		"http_requests_total":  m.httpRequestsTotal,
		"http_errors_total":    m.httpErrorsTotal,
		"audit_events_total":   m.auditEventsTotal,
		"auth_attempts_total":  m.authAttemptsTotal,
		"auth_success_total":   m.authSuccessTotal,
		"rate_limit_limited":   m.rateLimitLimited,
		"go_routines":          m.goRoutines,
		"memory_alloc_bytes":   m.memoryAllocBytes,
		"active_sessions":      m.activeSessions,
		"custom_metrics_count": len(m.customMetrics),
	}
}
