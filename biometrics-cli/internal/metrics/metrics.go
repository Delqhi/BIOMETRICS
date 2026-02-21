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

	WebhookRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_webhook_requests_total",
		Help: "Total number of webhook requests",
	})
	WebhookDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "biometrics_webhook_duration_seconds",
		Help:    "Duration of webhook requests in seconds",
		Buckets: prometheus.DefBuckets,
	})
	WebhookQueueDepth = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "biometrics_webhook_queue_depth",
		Help: "Current depth of webhook event queue",
	})

	WorkStealingTasksStolen = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_work_stealing_tasks_stolen_total",
		Help: "Total number of tasks stolen by work stealer",
	})
	WorkStealingTasksEnqueued = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_work_stealing_tasks_enqueued_total",
		Help: "Total number of tasks enqueued by agent",
	}, []string{"agent"})
	WorkStealingRebalances = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_work_stealing_rebalances_total",
		Help: "Total number of work stealing rebalances",
	})

	DurableCheckpointsCreated = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_durable_checkpoints_created_total",
		Help: "Total number of durable checkpoints created",
	}, []string{"agent"})
	DurableCheckpointsCompleted = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_durable_checkpoints_completed_total",
		Help: "Total number of durable checkpoints completed",
	}, []string{"agent"})
	DurableCheckpointsFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_durable_checkpoints_failed_total",
		Help: "Total number of durable checkpoints failed",
	}, []string{"agent"})
	DurableStepsReplayed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_durable_steps_replayed_total",
		Help: "Total number of durable steps replayed",
	}, []string{"agent"})

	HeartbeatsRegistered = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_heartbeats_registered_total",
		Help: "Total number of registered heartbeats",
	}, []string{"agent"})
	HeartbeatsReceived = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_heartbeats_received_total",
		Help: "Total number of received heartbeats",
	}, []string{"agent"})
	HeartbeatTasksDone = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_heartbeat_tasks_done_total",
		Help: "Total number of tasks done by heartbeat agents",
	}, []string{"agent"})
	HeartbeatTimeouts = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_heartbeat_timeouts_total",
		Help: "Total number of heartbeat timeouts",
	}, []string{"agent"})
	HeartbeatStuckAgents = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_heartbeat_stuck_agents_total",
		Help: "Total number of stuck agents detected",
	}, []string{"agent"})

	TasksCompletedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_tasks_completed_total",
		Help: "Total number of completed tasks",
	})
	TasksFailedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_tasks_failed_total",
		Help: "Total number of failed tasks",
	})
	TasksStartedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_tasks_started_total",
		Help: "Total number of started tasks",
	})

	AgentsStartedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_agents_started_total",
		Help: "Total number of started agents",
	})
	AgentsStoppedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_agents_stopped_total",
		Help: "Total number of stopped agents",
	})

	SessionsCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_sessions_created_total",
		Help: "Total number of created sessions",
	})
	SessionsEndedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_sessions_ended_total",
		Help: "Total number of ended sessions",
	})

	PlansActivatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_plans_activated_total",
		Help: "Total number of activated plans",
	})
	PlansCompletedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_plans_completed_total",
		Help: "Total number of completed plans",
	})

	SchedulerJobsRunTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_scheduler_jobs_run_total",
		Help: "Total number of scheduler jobs run",
	}, []string{"job"})
	SchedulerJobsFailedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_scheduler_jobs_failed_total",
		Help: "Total number of failed scheduler jobs",
	}, []string{"job"})
	SchedulerJobsSuccessTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_scheduler_jobs_success_total",
		Help: "Total number of successful scheduler jobs",
	}, []string{"job"})

	NotificationsSentTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_notifications_sent_total",
		Help: "Total number of sent notifications",
	})
	NotificationsFailedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_notifications_failed_total",
		Help: "Total number of failed notifications",
	})
	NotificationsDroppedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_notifications_dropped_total",
		Help: "Total number of dropped notifications",
	})

	DockerContainersRunning = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "biometrics_docker_containers_running",
		Help: "Number of running Docker containers",
	})
	DockerContainerStartsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_docker_container_starts_total",
		Help: "Total number of container starts",
	})
	DockerContainerStartsFailedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_docker_container_starts_failed_total",
		Help: "Total number of failed container starts",
	})
	DockerContainerStopsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_docker_container_stops_total",
		Help: "Total number of container stops",
	})
	DockerContainerStopsFailedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_docker_container_stops_failed_total",
		Help: "Total number of failed container stops",
	})
	DockerContainerRestartsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_docker_container_restarts_total",
		Help: "Total number of container restarts",
	})
	DockerContainersRemovedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_docker_containers_removed_total",
		Help: "Total number of removed containers",
	})
	DockerImagePullsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_docker_image_pulls_total",
		Help: "Total number of image pulls",
	})
	DockerImagePullsFailedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_docker_image_pulls_failed_total",
		Help: "Total number of failed image pulls",
	})
	DockerNetworksCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_docker_networks_created_total",
		Help: "Total number of created networks",
	})
	DockerNetworksRemovedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_docker_networks_removed_total",
		Help: "Total number of removed networks",
	})

	GitCommitsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_git_commits_total",
		Help: "Total number of git commits",
	})
	GitPullsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_git_pulls_total",
		Help: "Total number of git pulls",
	})
	GitPushesTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_git_pushes_total",
		Help: "Total number of git pushes",
	})
	GitFetchesTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_git_fetches_total",
		Help: "Total number of git fetches",
	})

	RateLimitAllowedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_rate_limit_allowed_total",
		Help: "Total number of allowed rate limit requests",
	}, []string{"key"})
	RateLimitRejectedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_rate_limit_rejected_total",
		Help: "Total number of rejected rate limit requests",
	}, []string{"key"})

	TasksCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_tasks_created_total",
		Help: "Total number of created tasks",
	})

	AgentDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "biometrics_agent_duration_seconds",
		Help:    "Duration of agent executions in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"agent", "status"})

	TaskQueueSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "biometrics_task_queue_size",
		Help: "Current size of task queue",
	})

	CacheHitsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_cache_hits_total",
		Help: "Total number of cache hits",
	})
	CacheMissesTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_cache_misses_total",
		Help: "Total number of cache misses",
	})
	CacheSizeBytes = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "biometrics_cache_size_bytes",
		Help: "Current size of cache in bytes",
	})

	SessionActiveCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "biometrics_sessions_active",
		Help: "Number of currently active sessions",
	})
	SessionAverageDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "biometrics_session_duration_seconds",
		Help:    "Average session duration in seconds",
		Buckets: prometheus.DefBuckets,
	})

	WebhookEventsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_webhook_events_total",
		Help: "Total number of webhook events by type",
	}, []string{"event_type", "status"})

	SchedulerJobDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "biometrics_scheduler_job_duration_seconds",
		Help: "Duration of scheduler jobs in seconds",
	}, []string{"job"})

	ProjectsRegistered = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_projects_registered_total",
		Help: "Total number of registered projects",
	}, []string{"project"})
	ProjectsUnregistered = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_projects_unregistered_total",
		Help: "Total number of unregistered projects",
	}, []string{"project"})
	ProjectRebalanced = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_project_rebalanced_total",
		Help: "Total number of project rebalances",
	})

	ScheduledTasksTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_scheduled_tasks_total",
		Help: "Total number of scheduled tasks",
	}, []string{"project", "task_type"})
	ScheduledTasksCompleted = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_scheduled_tasks_completed_total",
		Help: "Total number of completed scheduled tasks",
	}, []string{"project", "task_type"})
	ScheduledTasksFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_scheduled_tasks_failed_total",
		Help: "Total number of failed scheduled tasks",
	}, []string{"project", "task_type"})
)
