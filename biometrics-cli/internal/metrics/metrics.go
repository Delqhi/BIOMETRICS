package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CyclesTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_orchestrator_cycles_total",
		Help: "The total number of processed cycles",
	})
	ChaosEventsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_orchestrator_chaos_events_total",
		Help: "The total number of simulated chaos events",
	}, []string{"type"})
	CycleDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "biometrics_orchestrator_cycle_duration_seconds",
		Help:    "Duration of cycles in seconds",
		Buckets: prometheus.DefBuckets,
	})
	ModelAcquisitions = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_orchestrator_model_acquisitions_total",
		Help: "The total number of model acquisitions by model name",
	}, []string{"model"})
	SerenaSessionsCleaned = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_orchestrator_serena_sessions_cleaned_total",
		Help: "Total number of Serena sessions cleaned up",
	})
	HealingFailures = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_selfhealing_failures_total",
		Help: "Total number of self-healing failures by component",
	}, []string{"component"})
	HealingSuccesses = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_selfhealing_successes_total",
		Help: "Total number of self-healing successes by component",
	}, []string{"component"})
	ActiveSessions = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "biometrics_active_sessions",
		Help: "Number of currently active OpenCode sessions",
	})
)
