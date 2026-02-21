package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// RevenueEarningsTotal tracks total earnings in USD
	RevenueEarningsTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "revenue_earnings_total",
			Help: "Total earnings generated in USD",
		},
	)

	// RevenueTasksCompleted tracks number of completed revenue tasks
	RevenueTasksCompleted = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "revenue_tasks_completed_total",
			Help: "Total number of revenue-generating tasks completed",
		},
		[]string{"type"}, // captcha, survey, usertesting
	)

	// RevenueTaskDuration tracks task completion time
	RevenueTaskDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "revenue_task_duration_seconds",
			Help:    "Duration of revenue tasks in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"type", "status"},
	)

	// RevenueAPIErrors tracks API call failures
	RevenueAPIErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "revenue_api_errors_total",
			Help: "Total number of revenue API errors",
		},
		[]string{"endpoint", "error_type"},
	)

	// RevenueActiveWorkers tracks number of active revenue workers
	RevenueActiveWorkers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "revenue_active_workers",
			Help: "Number of currently active revenue workers",
		},
	)

	// RevenueAverageEarningsPerTask tracks average earnings
	RevenueAverageEarningsPerTask = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "revenue_average_earnings_per_task",
			Help: "Average earnings per completed task in USD",
		},
	)
)

// RegisterRevenueMetrics registers all revenue-related metrics
func RegisterRevenueMetrics() {
	prometheus.MustRegister(
		RevenueEarningsTotal,
		RevenueTasksCompleted,
		RevenueTaskDuration,
		RevenueAPIErrors,
		RevenueActiveWorkers,
		RevenueAverageEarningsPerTask,
	)
}
