package performance

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestProfiler_StartStop(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false

	profiler, err := NewProfiler(config, logger)
	require.NoError(t, err)

	err = profiler.Start()
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	err = profiler.Stop()
	require.NoError(t, err)
}

func TestProfiler_CaptureSnapshot(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false
	config.EnableCPU = false

	profiler, err := NewProfiler(config, logger)
	require.NoError(t, err)

	ctx := context.Background()
	snapshots, err := profiler.CaptureSnapshot(ctx)
	require.NoError(t, err)

	assert.Contains(t, snapshots, "memory")
	assert.Contains(t, snapshots, "goroutine")
	assert.Contains(t, snapshots, "gc_stats")
}

func TestProfiler_ReportMetrics(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false

	profiler, err := NewProfiler(config, logger)
	require.NoError(t, err)

	metrics := profiler.ReportMetrics()

	assert.Greater(t, metrics.Goroutines, 0)
	assert.Greater(t, metrics.CPUs, 0)
	assert.Greater(t, metrics.Alloc, uint64(0))
}

func TestRegisterPprofHandlers(t *testing.T) {
	mux := http.NewServeMux()
	RegisterPprofHandlers(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/debug/pprof/")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestWithPprofMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	wrapped := WithPprofMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestProfilingManager(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false

	manager, err := NewProfilingManager(config, logger)
	require.NoError(t, err)

	ctx := context.Background()

	err = manager.StartProfile(ctx, ProfileTypeMemory, 0)
	require.NoError(t, err)

	err = manager.StartProfile(ctx, ProfileTypeGoroutine, 0)
	require.NoError(t, err)

	active := manager.GetActiveProfiles()
	assert.Len(t, active, 0)

	snapshots, err := manager.CaptureAllProfiles(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, snapshots)
}

func TestRuntimeMetrics(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false

	manager, err := NewProfilingManager(config, logger)
	require.NoError(t, err)

	metrics := manager.GetMetrics()

	assert.Greater(t, metrics.Goroutines, 0)
	assert.Greater(t, metrics.CPUs, 0)
	assert.Greater(t, metrics.Alloc, uint64(0))
	assert.GreaterOrEqual(t, metrics.NumGC, uint32(0))
}
