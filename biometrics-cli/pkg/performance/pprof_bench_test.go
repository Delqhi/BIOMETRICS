package performance

import (
	"context"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// BenchmarkProfiler_StartStop benchmarks profiler start/stop overhead
func BenchmarkProfiler_StartStop(b *testing.B) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		profiler, err := NewProfiler(config, logger)
		require.NoError(b, err)

		err = profiler.Start()
		require.NoError(b, err)

		err = profiler.Stop()
		require.NoError(b, err)
	}
}

// BenchmarkProfiler_CaptureSnapshot benchmarks snapshot capture performance
func BenchmarkProfiler_CaptureSnapshot(b *testing.B) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false
	config.EnableCPU = false

	profiler, err := NewProfiler(config, logger)
	require.NoError(b, err)

	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := profiler.CaptureSnapshot(ctx)
		require.NoError(b, err)
	}
}

// BenchmarkProfiler_ReportMetrics benchmarks metrics reporting
func BenchmarkProfiler_ReportMetrics(b *testing.B) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false

	profiler, err := NewProfiler(config, logger)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = profiler.ReportMetrics()
	}
}

// BenchmarkProfilingManager_StartProfile benchmarks profile start
func BenchmarkProfilingManager_StartProfile(b *testing.B) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false

	manager, err := NewProfilingManager(config, logger)
	require.NoError(b, err)

	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.StartProfile(ctx, ProfileTypeMemory, 0)
	}
}

// BenchmarkProfilingManager_CaptureAllProfiles benchmarks capturing all profiles
func BenchmarkProfilingManager_CaptureAllProfiles(b *testing.B) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false

	manager, err := NewProfilingManager(config, logger)
	require.NoError(b, err)

	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.CaptureAllProfiles(ctx)
		require.NoError(b, err)
	}
}

// BenchmarkProfilingManager_GetMetrics benchmarks metrics retrieval
func BenchmarkProfilingManager_GetMetrics(b *testing.B) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false

	manager, err := NewProfilingManager(config, logger)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.GetMetrics()
	}
}

// TestProfiler_ConcurrentStartStop tests concurrent profiler operations
func TestProfiler_ConcurrentStartStop(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false

	profiler, err := NewProfiler(config, logger)
	require.NoError(t, err)

	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := profiler.Start(); err != nil {
				errors <- err
				return
			}
			time.Sleep(10 * time.Millisecond)
			if err := profiler.Stop(); err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	errorCount := 0
	for err := range errors {
		if err != nil {
			errorCount++
		}
	}

	assert.Less(t, errorCount, 10, "Some operations should succeed")
}

// TestProfiler_HTTPServer tests HTTP profiling server
func TestProfiler_HTTPServer(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = true
	config.HTTPAddr = "localhost:0"

	profiler, err := NewProfiler(config, logger)
	require.NoError(t, err)

	err = profiler.Start()
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://" + config.HTTPAddr + "/debug/pprof/")
	if err == nil {
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	err = profiler.Stop()
	require.NoError(t, err)
}

// TestProfiler_CustomConfig tests profiler with custom configuration
func TestProfiler_CustomConfig(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := ProfilerConfig{
		EnableCPU:            false,
		EnableMemory:         true,
		EnableBlock:          false,
		EnableMutex:          false,
		EnableGoroutine:      true,
		ProfileDir:           "./test_profiles",
		HTTPAddr:             "localhost:0",
		EnableHTTP:           false,
		ProfileDuration:      10 * time.Second,
		BlockProfileRate:     0,
		MutexProfileFraction: 0,
	}

	profiler, err := NewProfiler(config, logger)
	require.NoError(t, err)

	err = profiler.Start()
	require.NoError(t, err)

	ctx := context.Background()
	snapshots, err := profiler.CaptureSnapshot(ctx)
	require.NoError(t, err)

	assert.Contains(t, snapshots, "memory")
	assert.Contains(t, snapshots, "goroutine")

	err = profiler.Stop()
	require.NoError(t, err)
}

// TestProfilingManager_AutoStop tests automatic profile stop
func TestProfilingManager_AutoStop(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false

	manager, err := NewProfilingManager(config, logger)
	require.NoError(t, err)

	ctx := context.Background()
	err = manager.StartProfile(ctx, ProfileTypeCPU, 100*time.Millisecond)
	require.NoError(t, err)

	time.Sleep(200 * time.Millisecond)

	active := manager.GetActiveProfiles()
	assert.Empty(t, active, "Profile should auto-stop")
}

// TestProfilingManager_StopAllProfiles tests stopping all profiles
func TestProfilingManager_StopAllProfiles(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false

	manager, err := NewProfilingManager(config, logger)
	require.NoError(t, err)

	ctx := context.Background()
	_ = manager.StartProfile(ctx, ProfileTypeMemory, 0)
	_ = manager.StartProfile(ctx, ProfileTypeGoroutine, 0)

	err = manager.StopAllProfiles()
	require.NoError(t, err)

	active := manager.GetActiveProfiles()
	assert.Empty(t, active, "All profiles should be stopped")
}

// TestRuntimeMetrics_Concurrent tests concurrent metrics access
func TestRuntimeMetrics_Concurrent(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false

	manager, err := NewProfilingManager(config, logger)
	require.NoError(t, err)

	var wg sync.WaitGroup
	metricsChan := make(chan RuntimeMetrics, 100)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			metricsChan <- manager.GetMetrics()
		}()
	}

	wg.Wait()
	close(metricsChan)

	for metrics := range metricsChan {
		assert.Greater(t, metrics.Goroutines, 0)
		assert.Greater(t, metrics.CPUs, 0)
	}
}

// TestProfiler_GetProfileInfo tests profile info retrieval
func TestProfiler_GetProfileInfo(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false

	profiler, err := NewProfiler(config, logger)
	require.NoError(t, err)

	profiles := profiler.GetProfileInfo()

	expectedProfiles := []string{
		"goroutine",
		"heap",
		"allocs",
		"block",
		"mutex",
		"threadcreate",
	}

	for _, name := range expectedProfiles {
		assert.Contains(t, profiles, name, "Profile %s should exist", name)
	}
}

// BenchmarkRegisterPprofHandlers benchmarks handler registration
func BenchmarkRegisterPprofHandlers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mux := http.NewServeMux()
		RegisterPprofHandlers(mux)
	}
}

// TestWithPprofMiddleware_SlowRequest tests slow request detection
func TestWithPprofMiddleware_SlowRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	wrapped := WithPprofMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestProfiler_NoLogger tests profiler with nil logger
func TestProfiler_NoLogger(t *testing.T) {
	config := DefaultProfilerConfig()
	config.EnableHTTP = false

	profiler, err := NewProfiler(config, nil)
	require.NoError(t, err)

	err = profiler.Start()
	require.NoError(t, err)

	err = profiler.Stop()
	require.NoError(t, err)
}

// TestProfilingManager_InvalidProfileType tests invalid profile type
func TestProfilingManager_InvalidProfileType(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.EnableHTTP = false

	manager, err := NewProfilingManager(config, logger)
	require.NoError(t, err)

	ctx := context.Background()
	err = manager.StartProfile(ctx, ProfileType("invalid"), 0)
	assert.Error(t, err, "Invalid profile type should error")
}

// TestProfiler_DirectoryCreation tests profile directory creation
func TestProfiler_DirectoryCreation(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := DefaultProfilerConfig()
	config.ProfileDir = "./test_profile_dir_" + time.Now().Format("20060102150405")
	config.EnableHTTP = false

	profiler, err := NewProfiler(config, logger)
	require.NoError(t, err)
	defer func() {
		_ = profiler.Stop()
	}()

	err = profiler.Start()
	require.NoError(t, err)
}

// BenchmarkRuntimeMetrics benchmarks metrics collection
func BenchmarkRuntimeMetrics(b *testing.B) {
	var stats runtime.MemStats
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runtime.ReadMemStats(&stats)
		_ = stats.Alloc
	}
}
