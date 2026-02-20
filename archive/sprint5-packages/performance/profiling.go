package performance

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sync"
	"time"

	"go.uber.org/zap"
)

type ProfileType string

const (
	ProfileTypeCPU       ProfileType = "cpu"
	ProfileTypeMemory    ProfileType = "memory"
	ProfileTypeBlock     ProfileType = "block"
	ProfileTypeMutex     ProfileType = "mutex"
	ProfileTypeGoroutine ProfileType = "goroutine"
	ProfileTypeThread    ProfileType = "threadcreate"
	ProfileTypeAllocs    ProfileType = "allocs"
	ProfileTypeTrace     ProfileType = "trace"
)

type ProfilingManager struct {
	config      ProfilerConfig
	logger      *zap.Logger
	profilesDir string
	mu          sync.RWMutex
	profiles    map[ProfileType]*ActiveProfile
}

type ActiveProfile struct {
	Type      ProfileType
	File      *os.File
	StartTime time.Time
	Duration  time.Duration
}

func NewProfilingManager(config ProfilerConfig, logger *zap.Logger) (*ProfilingManager, error) {
	profilesDir := config.ProfileDir
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create profiles directory: %w", err)
	}

	return &ProfilingManager{
		config:      config,
		logger:      logger,
		profilesDir: profilesDir,
		profiles:    make(map[ProfileType]*ActiveProfile),
	}, nil
}

func (pm *ProfilingManager) StartProfile(ctx context.Context, profileType ProfileType, duration time.Duration) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.profiles[profileType]; exists {
		return fmt.Errorf("profile %s already running", profileType)
	}

	filename := filepath.Join(pm.profilesDir, fmt.Sprintf("%s_%s.prof", profileType, time.Now().Format("20060102_150405")))
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create profile file: %w", err)
	}

	var startErr error
	switch profileType {
	case ProfileTypeCPU:
		startErr = pprof.StartCPUProfile(f)
	case ProfileTypeMemory:
		startErr = pprof.WriteHeapProfile(f)
		f.Close()
		pm.logger.Info("Memory profile captured", zap.String("file", filename))
		return nil
	case ProfileTypeBlock:
		profile := pprof.Lookup("block")
		if profile == nil {
			f.Close()
			return fmt.Errorf("block profile not available")
		}
		startErr = profile.WriteTo(f, 0)
		f.Close()
		pm.logger.Info("Block profile captured", zap.String("file", filename))
		return nil
	case ProfileTypeMutex:
		profile := pprof.Lookup("mutex")
		if profile == nil {
			f.Close()
			return fmt.Errorf("mutex profile not available")
		}
		startErr = profile.WriteTo(f, 0)
		f.Close()
		pm.logger.Info("Mutex profile captured", zap.String("file", filename))
		return nil
	case ProfileTypeGoroutine:
		profile := pprof.Lookup("goroutine")
		if profile == nil {
			f.Close()
			return fmt.Errorf("goroutine profile not available")
		}
		startErr = profile.WriteTo(f, 0)
		f.Close()
		pm.logger.Info("Goroutine profile captured", zap.String("file", filename))
		return nil
	case ProfileTypeTrace:
		startErr = trace.Start(f)
	default:
		f.Close()
		return fmt.Errorf("unknown profile type: %s", profileType)
	}

	if startErr != nil {
		f.Close()
		return fmt.Errorf("failed to start %s profile: %w", profileType, startErr)
	}

	pm.profiles[profileType] = &ActiveProfile{
		Type:      profileType,
		File:      f,
		StartTime: time.Now(),
		Duration:  duration,
	}

	pm.logger.Info("Profile started",
		zap.String("type", string(profileType)),
		zap.String("file", filename),
		zap.Duration("duration", duration),
	)

	if duration > 0 {
		go pm.autoStop(profileType, duration)
	}

	return nil
}

func (pm *ProfilingManager) autoStop(profileType ProfileType, duration time.Duration) {
	<-time.After(duration)
	if err := pm.StopProfile(profileType); err != nil {
		pm.logger.Error("Failed to auto-stop profile", zap.Error(err))
	}
}

func (pm *ProfilingManager) StopProfile(profileType ProfileType) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	active, exists := pm.profiles[profileType]
	if !exists {
		return fmt.Errorf("profile %s not running", profileType)
	}

	switch profileType {
	case ProfileTypeCPU:
		pprof.StopCPUProfile()
	case ProfileTypeTrace:
		trace.Stop()
	}

	if active.File != nil {
		active.File.Close()
	}

	delete(pm.profiles, profileType)

	pm.logger.Info("Profile stopped",
		zap.String("type", string(profileType)),
		zap.Duration("duration", time.Since(active.StartTime)),
	)

	return nil
}

func (pm *ProfilingManager) StopAllProfiles() error {
	pm.mu.Lock()
	profileTypes := make([]ProfileType, 0, len(pm.profiles))
	for pt := range pm.profiles {
		profileTypes = append(profileTypes, pt)
	}
	pm.mu.Unlock()

	var lastErr error
	for _, pt := range profileTypes {
		if err := pm.StopProfile(pt); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func (pm *ProfilingManager) GetActiveProfiles() []ProfileType {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	types := make([]ProfileType, 0, len(pm.profiles))
	for pt := range pm.profiles {
		types = append(types, pt)
	}
	return types
}

func (pm *ProfilingManager) CaptureAllProfiles(ctx context.Context) (map[string]string, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	snapshots := make(map[string]string)
	timestamp := time.Now().Format("20060102_150405")

	profileTypes := []ProfileType{
		ProfileTypeMemory,
		ProfileTypeGoroutine,
		ProfileTypeBlock,
		ProfileTypeMutex,
		ProfileTypeAllocs,
	}

	for _, pt := range profileTypes {
		filename := filepath.Join(pm.profilesDir, fmt.Sprintf("%s_%s.prof", pt, timestamp))
		f, err := os.Create(filename)
		if err != nil {
			pm.logger.Warn("Failed to create profile file", zap.String("type", string(pt)), zap.Error(err))
			continue
		}

		var profile *pprof.Profile
		switch pt {
		case ProfileTypeMemory:
			profile = pprof.Lookup("heap")
		case ProfileTypeGoroutine:
			profile = pprof.Lookup("goroutine")
		case ProfileTypeBlock:
			profile = pprof.Lookup("block")
		case ProfileTypeMutex:
			profile = pprof.Lookup("mutex")
		case ProfileTypeAllocs:
			profile = pprof.Lookup("allocs")
		case ProfileTypeThread:
			profile = pprof.Lookup("threadcreate")
		}

		if profile != nil {
			if err := profile.WriteTo(f, 0); err != nil {
				f.Close()
				pm.logger.Warn("Failed to write profile", zap.String("type", string(pt)), zap.Error(err))
				continue
			}
		}
		f.Close()

		snapshots[string(pt)] = filename
	}

	gcStatsFile := filepath.Join(pm.profilesDir, fmt.Sprintf("gc_stats_%s.txt", timestamp))
	if err := pm.writeGCStats(gcStatsFile); err != nil {
		pm.logger.Warn("Failed to write GC stats", zap.Error(err))
	} else {
		snapshots["gc_stats"] = gcStatsFile
	}

	return snapshots, nil
}

func (pm *ProfilingManager) writeGCStats(filename string) error {
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
PauseNs (last 256): %v
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
		stats.PauseNs[:],
	)

	return err
}

func (pm *ProfilingManager) GetMetrics() RuntimeMetrics {
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
