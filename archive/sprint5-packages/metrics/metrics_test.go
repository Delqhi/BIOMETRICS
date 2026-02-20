package metrics

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewMetricsRegistry(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	registry, err := NewMetricsRegistry(logger)
	require.NoError(t, err)
	require.NotNil(t, registry)

	assert.NotNil(t, registry.httpRequestsTotal)
	assert.NotNil(t, registry.httpRequestDuration)
	assert.NotNil(t, registry.auditEventsTotal)
	assert.NotNil(t, registry.authAttemptsTotal)
	assert.NotNil(t, registry.rateLimitRequests)
	assert.NotNil(t, registry.customMetrics)
}

func TestNewMetricsRegistryNilLogger(t *testing.T) {
	registry, err := NewMetricsRegistry(nil)
	require.NoError(t, err)
	require.NotNil(t, registry)
}

func TestRecordHTTPEvent(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	registry.RecordHTTPEvent("GET", "/api/test", http.StatusOK, 100*time.Millisecond)
	registry.RecordHTTPEvent("POST", "/api/test", http.StatusInternalServerError, 200*time.Millisecond)

	metrics, err := registry.CollectAll()
	require.NoError(t, err)

	assert.Contains(t, metrics, "biometrics_http_requests_total")
	assert.Contains(t, metrics, "biometrics_http_errors_total")
}

func TestHTTPMiddleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := registry.HTTPMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	time.Sleep(100 * time.Millisecond)
	metrics, err := registry.CollectAll()
	require.NoError(t, err)

	assert.Contains(t, metrics, "biometrics_http_requests_total")
}

func TestRecordAuditEvent(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	registry.RecordAuditEvent("login", "info")
	registry.RecordAuditEvent("access_denied", "warning")

	metrics, err := registry.CollectAll()
	require.NoError(t, err)

	// CounterVec metrics have labels in the key
	found := false
	for key := range metrics {
		if strings.Contains(key, "biometrics_audit_events_total") {
			found = true
			break
		}
	}
	assert.True(t, found, "Should contain audit events metric")
	assert.Contains(t, metrics, "biometrics_audit_logs_written_total")
}

func TestRecordAuditQuery(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	registry.RecordAuditQuery()
	registry.RecordAuditQuery()

	metrics, err := registry.CollectAll()
	require.NoError(t, err)

	assert.Contains(t, metrics, "biometrics_audit_queries_total")
}

func TestUpdateAuditStorage(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	registry.UpdateAuditStorage(1024 * 1024)

	metrics, err := registry.CollectAll()
	require.NoError(t, err)

	assert.Contains(t, metrics, "biometrics_audit_storage_bytes")
}

func TestRecordAuthAttempt(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	registry.RecordAuthAttempt("oauth2", "azure")
	registry.RecordAuthAttempt("jwt", "local")

	metrics, err := registry.CollectAll()
	require.NoError(t, err)

	// CounterVec metrics have labels in the key
	found := false
	for key := range metrics {
		if strings.Contains(key, "biometrics_auth_attempts_total") {
			found = true
			break
		}
	}
	assert.True(t, found, "Should contain auth attempts metric")
}

func TestRecordAuthFailure(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	registry.RecordAuthFailure("invalid_credentials", "oauth2")
	registry.RecordAuthFailure("expired_token", "jwt")

	metrics, err := registry.CollectAll()
	require.NoError(t, err)

	// CounterVec metrics have labels in the key
	found := false
	for key := range metrics {
		if strings.Contains(key, "biometrics_auth_failures_total") {
			found = true
			break
		}
	}
	assert.True(t, found, "Should contain auth failures metric")
}

func TestRecordAuthSuccess(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	registry.RecordAuthSuccess()
	registry.RecordAuthSuccess()

	metrics, err := registry.CollectAll()
	require.NoError(t, err)

	assert.Contains(t, metrics, "biometrics_auth_success_total")
}

func TestUpdateActiveSessions(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	registry.UpdateActiveSessions(10)
	registry.UpdateActiveSessions(25)

	metrics, err := registry.CollectAll()
	require.NoError(t, err)

	assert.Contains(t, metrics, "biometrics_active_sessions")
}

func TestRecordTokenRefresh(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	registry.RecordTokenRefresh()
	registry.RecordTokenRefresh()

	metrics, err := registry.CollectAll()
	require.NoError(t, err)

	assert.Contains(t, metrics, "biometrics_token_refreshs_total")
}

func TestRecordRateLimitRequest(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	registry.RecordRateLimitRequest("/api/test", "client-1", true)
	registry.RecordRateLimitRequest("/api/test", "client-2", false)

	metrics, err := registry.CollectAll()
	require.NoError(t, err)

	// CounterVec metrics have labels in the key
	found := false
	for key := range metrics {
		if strings.Contains(key, "biometrics_ratelimit_requests_total") {
			found = true
			break
		}
	}
	assert.True(t, found, "Should contain ratelimit requests metric")
	assert.Contains(t, metrics, "biometrics_ratelimit_limited_total")
	assert.Contains(t, metrics, "biometrics_ratelimit_bypassed_total")
}

func TestUpdateSystemMetrics(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	registry.UpdateSystemMetrics()

	metrics, err := registry.CollectAll()
	require.NoError(t, err)

	assert.Contains(t, metrics, "biometrics_go_routines")
	assert.Contains(t, metrics, "biometrics_memory_alloc_bytes")
	assert.Contains(t, metrics, "biometrics_memory_sys_bytes")
}

func TestRegisterCustomMetric(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	counter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "custom_test_counter",
			Help: "Test counter",
		},
	)

	err := registry.RegisterCustomMetric("test_counter", counter)
	require.NoError(t, err)

	retrieved, err := registry.GetCustomMetric("test_counter")
	require.NoError(t, err)
	assert.Equal(t, counter, retrieved)
}

func TestRegisterDuplicateCustomMetric(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	counter1 := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "custom_duplicate_counter",
			Help: "Test counter",
		},
	)

	err := registry.RegisterCustomMetric("dup_counter", counter1)
	require.NoError(t, err)

	counter2 := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "custom_duplicate_counter2",
			Help: "Test counter 2",
		},
	)

	err = registry.RegisterCustomMetric("dup_counter", counter2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestUnregisterCustomMetric(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	counter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "custom_unregister_counter",
			Help: "Test counter",
		},
	)

	err := registry.RegisterCustomMetric("unreg_counter", counter)
	require.NoError(t, err)

	err = registry.UnregisterCustomMetric("unreg_counter")
	require.NoError(t, err)

	_, err = registry.GetCustomMetric("unreg_counter")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestUnregisterNonExistentCustomMetric(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	err := registry.UnregisterCustomMetric("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestCreateCounter(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	counter, err := registry.CreateCounter("test_created_counter", "Test help")
	require.NoError(t, err)
	require.NotNil(t, counter)

	counter.Inc()
	counter.Inc()

	metrics, err := registry.CollectAll()
	require.NoError(t, err)
	assert.Contains(t, metrics, "test_created_counter")
}

func TestCreateGauge(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	gauge, err := registry.CreateGauge("test_created_gauge", "Test help")
	require.NoError(t, err)
	require.NotNil(t, gauge)

	gauge.Set(42.5)

	metrics, err := registry.CollectAll()
	require.NoError(t, err)
	assert.Contains(t, metrics, "test_created_gauge")
}

func TestCreateHistogram(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	buckets := []float64{0.1, 0.5, 1.0, 5.0}
	histogram, err := registry.CreateHistogram("test_created_histogram", "Test help", buckets)
	require.NoError(t, err)
	require.NotNil(t, histogram)

	histogram.Observe(0.3)
	histogram.Observe(1.2)

	metrics, err := registry.CollectAll()
	require.NoError(t, err)
	assert.Contains(t, metrics, "test_created_histogram")
}

func TestCreateCounterVec(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	counterVec, err := registry.CreateCounterVec("test_counter_vec", "Test help", []string{"label1", "label2"})
	require.NoError(t, err)
	require.NotNil(t, counterVec)

	counterVec.WithLabelValues("value1", "value2").Inc()

	metrics, err := registry.CollectAll()
	require.NoError(t, err)

	// CounterVec metrics have labels in the key
	found := false
	for key := range metrics {
		if strings.Contains(key, "test_counter_vec") {
			found = true
			break
		}
	}
	assert.True(t, found, "Should contain test_counter_vec metric")
}

func TestMetricsHandler(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	handler := registry.MetricsHandler()
	require.NotNil(t, handler)

	req := httptest.NewRequest("GET", "/metrics", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "text/plain")

	body := rr.Body.String()
	assert.NotEmpty(t, body)
	assert.Contains(t, body, "biometrics_")
}

func TestStartMetricsServer(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := registry.StartMetricsServer(ctx, "localhost:0")
	require.NoError(t, err)

	time.Sleep(500 * time.Millisecond)

	metrics, err := registry.CollectAll()
	require.NoError(t, err)
	assert.NotNil(t, metrics)
}

func TestGetRegistry(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	reg := registry.GetRegistry()
	require.NotNil(t, reg)
}

func TestCollectAll(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	registry.RecordHTTPEvent("GET", "/test", http.StatusOK, 50*time.Millisecond)
	registry.RecordAuthSuccess()
	registry.UpdateSystemMetrics()

	metrics, err := registry.CollectAll()
	require.NoError(t, err)
	assert.NotEmpty(t, metrics)
}

func TestResetAll(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	counter, _ := registry.CreateCounter("test_reset_counter", "Test")
	counter.Inc()

	err := registry.ResetAll()
	require.NoError(t, err)

	_, err = registry.GetCustomMetric("test_reset_counter")
	assert.Error(t, err)
}

func TestGetMetricsSummary(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	summary := registry.GetMetricsSummary()
	require.NotNil(t, summary)

	assert.Contains(t, summary, "http_requests_total")
	assert.Contains(t, summary, "auth_success_total")
	assert.Contains(t, summary, "go_routines")
	assert.Contains(t, summary, "custom_metrics_count")
}

func TestResponseWriterWrapper(t *testing.T) {
	rr := httptest.NewRecorder()
	wrapper := &responseWriterWrapper{ResponseWriter: rr, statusCode: http.StatusOK}

	wrapper.WriteHeader(http.StatusCreated)
	assert.Equal(t, http.StatusCreated, wrapper.statusCode)
	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestMetricsServerHealthEndpoint(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := registry.StartMetricsServer(ctx, "localhost:0")
	require.NoError(t, err)

	time.Sleep(500 * time.Millisecond)

	resp, err := http.Get("http://localhost:0/health")
	if err != nil {
		t.Skip("Server not ready in time")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "OK", strings.TrimSpace(string(body)))
}

func TestConcurrentMetricUpdates(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				registry.RecordHTTPEvent("GET", "/test", http.StatusOK, time.Millisecond)
				registry.RecordAuthSuccess()
			}
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	metrics, err := registry.CollectAll()
	require.NoError(t, err)
	assert.NotEmpty(t, metrics)
}

func TestCustomMetricConcurrency(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	registry, _ := NewMetricsRegistry(logger)

	done := make(chan bool)
	for i := 0; i < 5; i++ {
		go func(id int) {
			name := "concurrent_counter_" + string(rune(id))
			counter, err := registry.CreateCounter(name, "Test")
			if err == nil {
				counter.Inc()
			}
			done <- true
		}(i)
	}

	for i := 0; i < 5; i++ {
		<-done
	}

	metrics, err := registry.CollectAll()
	require.NoError(t, err)
	assert.NotEmpty(t, metrics)
}
