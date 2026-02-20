// Package performance provides comprehensive performance monitoring and profiling capabilities
// using Go's built-in pprof infrastructure with production-ready enhancements.
package performance

import (
	"context"
	"fmt"
	"net/http"
	pprofHTTP "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"go.uber.org/zap"
)

// ProfilerConfig holds configuration for the profiler
type ProfilerConfig struct {
	// EnableCPU enables CPU profiling
	EnableCPU bool
	// EnableMemory enables memory profiling
	EnableMemory bool
	// EnableBlock enables block profiling
	EnableBlock bool
	// EnableMutex enables mutex profiling
	EnableMutex bool
	// EnableGoroutine enables goroutine profiling
	EnableGoroutine bool
	// ProfileDir is the directory to store profiles
	ProfileDir string
	// HTTPAddr is the address for the HTTP profiling server
	HTTPAddr string
	// EnableHTTP enables the HTTP profiling server
	EnableHTTP bool
	// ProfileDuration is the duration for automatic profiling
	ProfileDuration time.Duration
	// BlockProfileRate sets the fraction of blocking events to profile
	BlockProfileRate int
	// MutexProfileFraction sets the fraction of mutex contention to profile
	MutexProfileFraction int
}

// DefaultProfilerConfig returns a production-ready default configuration
func DefaultProfilerConfig() ProfilerConfig {
	return ProfilerConfig{
		EnableCPU:            true,
		EnableMemory:         true,
		EnableBlock:          true,
		EnableMutex:          true,
		EnableGoroutine:      true,
		ProfileDir:           "./profiles",
		HTTPAddr:             "localhost:6060",
		EnableHTTP:           true,
		ProfileDuration:      30 * time.Second,
		BlockProfileRate:     1,
		MutexProfileFraction: 1,
	}
}

// Profiler manages performance profiling
type Profiler struct {
	config     ProfilerConfig
	logger     *zap.Logger
	cpuFile    *os.File
	memFile    *os.File
	blockFile  *os.File
	mutexFile  *os.File
	httpServer *http.Server
	mu         sync.Mutex
	running    bool
}

// NewProfiler creates a new profiler instance
func NewProfiler(config ProfilerConfig, logger *zap.Logger) (*Profiler, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	// Create profile directory if it doesn't exist
	if err := os.MkdirAll(config.ProfileDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create profile directory: %w", err)
	}

	// Set profile rates
	runtime.SetBlockProfileRate(config.BlockProfileRate)
	runtime.SetMutexProfileFraction(config.MutexProfileFraction)

	profiler := &Profiler{
		config: config,
		logger: logger,
	}

	return profiler, nil
}

// Start starts the profiler and all enabled profiling types
func (p *Profiler) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return fmt.Errorf("profiler already running")
	}

	p.logger.Info("Starting performance profiler",
		zap.Bool("cpu", p.config.EnableCPU),
		zap.Bool("memory", p.config.EnableMemory),
		zap.Bool("block", p.config.EnableBlock),
		zap.Bool("mutex", p.config.EnableMutex),
		zap.Bool("goroutine", p.config.EnableGoroutine),
	)

	// Start CPU profiling if enabled
	if p.config.EnableCPU {
		if err := p.startCPUProfiling(); err != nil {
			return fmt.Errorf("failed to start CPU profiling: %w", err)
		}
	}

	// Start HTTP server if enabled
	if p.config.EnableHTTP {
		if err := p.startHTTPServer(); err != nil {
			return fmt.Errorf("failed to start HTTP profiling server: %w", err)
		}
	}

	p.running = true
	return nil
}

// Stop stops the profiler and writes all profiles
func (p *Profiler) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return nil
	}

	p.logger.Info("Stopping performance profiler")

	// Stop CPU profiling
	if p.cpuFile != nil {
		pprof.StopCPUProfile()
		p.cpuFile.Close()
		p.logger.Info("CPU profile written", zap.String("file", p.cpuFile.Name()))
	}

	// Write memory profile
	if p.config.EnableMemory {
		if err := p.writeMemoryProfile(); err != nil {
			p.logger.Error("Failed to write memory profile", zap.Error(err))
		}
	}

	// Write block profile
	if p.config.EnableBlock {
		if err := p.writeBlockProfile(); err != nil {
			p.logger.Error("Failed to write block profile", zap.Error(err))
		}
	}

	// Write mutex profile
	if p.config.EnableMutex {
		if err := p.writeMutexProfile(); err != nil {
			p.logger.Error("Failed to write mutex profile", zap.Error(err))
		}
	}

	// Write goroutine profile
	if p.config.EnableGoroutine {
		if err := p.writeGoroutineProfile(); err != nil {
			p.logger.Error("Failed to write goroutine profile", zap.Error(err))
		}
	}

	// Stop HTTP server
	if p.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := p.httpServer.Shutdown(ctx); err != nil {
			p.logger.Error("Failed to shutdown HTTP server", zap.Error(err))
		}
	}

	p.running = false
	return nil
}

// startCPUProfiling starts CPU profiling
func (p *Profiler) startCPUProfiling() error {
	filename := filepath.Join(p.config.ProfileDir, "cpu.prof")
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CPU profile file: %w", err)
	}

	p.cpuFile = f

	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		return fmt.Errorf("failed to start CPU profile: %w", err)
	}

	p.logger.Info("CPU profiling started", zap.String("file", filename))
	return nil
}

// writeMemoryProfile writes the memory profile
func (p *Profiler) writeMemoryProfile() error {
	filename := filepath.Join(p.config.ProfileDir, "mem.prof")
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create memory profile file: %w", err)
	}
	defer f.Close()

	if err := pprof.WriteHeapProfile(f); err != nil {
		return fmt.Errorf("failed to write memory profile: %w", err)
	}

	p.logger.Info("Memory profile written", zap.String("file", filename))
	return nil
}

// writeBlockProfile writes the block profile
func (p *Profiler) writeBlockProfile() error {
	filename := filepath.Join(p.config.ProfileDir, "block.prof")
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create block profile file: %w", err)
	}
	defer f.Close()

	profile := pprof.Lookup("block")
	if profile == nil {
		return fmt.Errorf("block profile not available")
	}

	if err := profile.WriteTo(f, 0); err != nil {
		return fmt.Errorf("failed to write block profile: %w", err)
	}

	p.logger.Info("Block profile written", zap.String("file", filename))
	return nil
}

// writeMutexProfile writes the mutex profile
func (p *Profiler) writeMutexProfile() error {
	filename := filepath.Join(p.config.ProfileDir, "mutex.prof")
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create mutex profile file: %w", err)
	}
	defer f.Close()

	profile := pprof.Lookup("mutex")
	if profile == nil {
		return fmt.Errorf("mutex profile not available")
	}

	if err := profile.WriteTo(f, 0); err != nil {
		return fmt.Errorf("failed to write mutex profile: %w", err)
	}

	p.logger.Info("Mutex profile written", zap.String("file", filename))
	return nil
}

// writeGoroutineProfile writes the goroutine profile
func (p *Profiler) writeGoroutineProfile() error {
	filename := filepath.Join(p.config.ProfileDir, "goroutine.prof")
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create goroutine profile file: %w", err)
	}
	defer f.Close()

	profile := pprof.Lookup("goroutine")
	if profile == nil {
		return fmt.Errorf("goroutine profile not available")
	}

	if err := profile.WriteTo(f, 0); err != nil {
		return fmt.Errorf("failed to write goroutine profile: %w", err)
	}

	p.logger.Info("Goroutine profile written", zap.String("file", filename))
	return nil
}

// startHTTPServer starts the HTTP profiling server
func (p *Profiler) startHTTPServer() error {
	mux := http.NewServeMux()

	// Register pprof handlers
	mux.HandleFunc("/debug/pprof/", pprofHTTP.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprofHTTP.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprofHTTP.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprofHTTP.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprofHTTP.Trace)

	mux.Handle("/debug/pprof/heap", pprofHTTP.Handler("heap"))
	mux.Handle("/debug/pprof/allocs", pprofHTTP.Handler("allocs"))
	mux.Handle("/debug/pprof/block", pprofHTTP.Handler("block"))
	mux.Handle("/debug/pprof/goroutine", pprofHTTP.Handler("goroutine"))
	mux.Handle("/debug/pprof/mutex", pprofHTTP.Handler("mutex"))
	mux.Handle("/debug/pprof/threadcreate", pprofHTTP.Handler("threadcreate"))

	p.httpServer = &http.Server{
		Addr:         p.config.HTTPAddr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		p.logger.Info("Starting HTTP profiling server", zap.String("addr", p.config.HTTPAddr))
		if err := p.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			p.logger.Error("HTTP profiling server failed", zap.Error(err))
		}
	}()

	return nil
}

// CaptureSnapshot captures an instant snapshot of all enabled profiles
func (p *Profiler) CaptureSnapshot(ctx context.Context) (map[string]string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	snapshots := make(map[string]string)
	timestamp := time.Now().Format("20060102_150405")

	// Capture memory profile
	if p.config.EnableMemory {
		filename := filepath.Join(p.config.ProfileDir, fmt.Sprintf("mem_snapshot_%s.prof", timestamp))
		f, err := os.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to create memory snapshot: %w", err)
		}

		if err := pprof.WriteHeapProfile(f); err != nil {
			f.Close()
			return nil, fmt.Errorf("failed to write memory snapshot: %w", err)
		}
		f.Close()

		snapshots["memory"] = filename
	}

	// Capture goroutine profile
	if p.config.EnableGoroutine {
		filename := filepath.Join(p.config.ProfileDir, fmt.Sprintf("goroutine_snapshot_%s.prof", timestamp))
		f, err := os.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to create goroutine snapshot: %w", err)
		}

		profile := pprof.Lookup("goroutine")
		if err := profile.WriteTo(f, 0); err != nil {
			f.Close()
			return nil, fmt.Errorf("failed to write goroutine snapshot: %w", err)
		}
		f.Close()

		snapshots["goroutine"] = filename
	}

	// Capture GC stats
	gcStatsFile := filepath.Join(p.config.ProfileDir, fmt.Sprintf("gc_stats_%s.txt", timestamp))
	if err := p.writeGCStats(gcStatsFile); err != nil {
		return nil, fmt.Errorf("failed to write GC stats: %w", err)
	}
	snapshots["gc_stats"] = gcStatsFile

	return snapshots, nil
}

// writeGCStats writes garbage collection statistics
func (p *Profiler) writeGCStats(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	_, err = fmt.Fprintf(f,
		`GC Statistics
=============
Alloc: %d KB
TotalAlloc: %d KB
Sys: %d KB
Lookups: %d
Mallocs: %d
Frees: %d
HeapAlloc: %d KB
HeapSys: %d KB
HeapIdle: %d KB
HeapInuse: %d KB
HeapReleased: %d KB
HeapObjects: %d
NumGC: %d
GCCPUFraction: %.6f
PauseTotalNs: %d ns
`,
		stats.Alloc/1024,
		stats.TotalAlloc/1024,
		stats.Sys/1024,
		stats.Lookups,
		stats.Mallocs,
		stats.Frees,
		stats.HeapAlloc/1024,
		stats.HeapSys/1024,
		stats.HeapIdle/1024,
		stats.HeapInuse/1024,
		stats.HeapReleased/1024,
		stats.HeapObjects,
		stats.NumGC,
		stats.GCCPUFraction,
		stats.PauseTotalNs,
	)

	return err
}

// GetProfileInfo returns information about available profiles
func (p *Profiler) GetProfileInfo() map[string]*pprof.Profile {
	profiles := make(map[string]*pprof.Profile)

	profileNames := []string{
		"goroutine",
		"heap",
		"allocs",
		"block",
		"mutex",
		"threadcreate",
	}

	for _, name := range profileNames {
		if profile := pprof.Lookup(name); profile != nil {
			profiles[name] = profile
		}
	}

	return profiles
}

// ReportMetrics returns current runtime metrics
func (p *Profiler) ReportMetrics() RuntimeMetrics {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	return RuntimeMetrics{
		Goroutines:    runtime.NumGoroutine(),
		CPUs:          runtime.NumCPU(),
		Alloc:         stats.Alloc,
		TotalAlloc:    stats.TotalAlloc,
		Sys:           stats.Sys,
		HeapAlloc:     stats.HeapAlloc,
		HeapObjects:   stats.HeapObjects,
		NumGC:         stats.NumGC,
		GCCPUFraction: stats.GCCPUFraction,
	}
}

// RuntimeMetrics holds current runtime statistics
type RuntimeMetrics struct {
	Goroutines    int
	CPUs          int
	Alloc         uint64
	TotalAlloc    uint64
	Sys           uint64
	HeapAlloc     uint64
	HeapObjects   uint64
	NumGC         uint32
	GCCPUFraction float64
}
