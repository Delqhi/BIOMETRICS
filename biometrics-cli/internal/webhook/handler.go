package webhook

import (
	"biometrics-cli/internal/metrics"
	"biometrics-cli/internal/state"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Handler struct {
	mu          sync.Mutex
	handlers    map[string]WebhookHandler
	authTokens  map[string]string
	rateLimiter *RateLimiter
	middlewares []Middleware
}

type WebhookHandler func(payload []byte) (interface{}, error)

type Middleware func(http.Handler) http.Handler

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

type WebhookPayload struct {
	Event     string          `json:"event"`
	Agent     string          `json:"agent"`
	SessionID string          `json:"session_id"`
	Data      json.RawMessage `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

var webhookHandler = &Handler{
	handlers:   make(map[string]WebhookHandler),
	authTokens: make(map[string]string),
	rateLimiter: &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    100,
		window:   time.Minute,
	},
	middlewares: make([]Middleware, 0),
}

func New() *Handler {
	return webhookHandler
}

func (h *Handler) Register(event string, handler WebhookHandler) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers[event] = handler
	state.GlobalState.Log("INFO", fmt.Sprintf("Registered webhook handler for event: %s", event))
}

func (h *Handler) SetAuthToken(token string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.authTokens["default"] = token
}

func (h *Handler) AddMiddleware(m Middleware) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.middlewares = append(h.middlewares, m)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	metrics.WebhookRequestsTotal.Inc()

	if r.Method != http.MethodPost {
		h.sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clientIP := r.RemoteAddr
	if h.rateLimiter.isRateLimited(clientIP) {
		h.sendError(w, "rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if !h.validateAuth(authHeader) {
		h.sendError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.sendError(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		h.sendError(w, "invalid json", http.StatusBadRequest)
		return
	}

	if payload.Timestamp.IsZero() {
		payload.Timestamp = time.Now()
	}

	h.mu.Lock()
	handler, exists := h.handlers[payload.Event]
	h.mu.Unlock()

	if !exists {
		h.sendError(w, "event not found", http.StatusNotFound)
		return
	}

	result, err := handler(body)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start).Seconds()
	metrics.WebhookDuration.Observe(duration)

	h.sendSuccess(w, result)
}

func (h *Handler) validateAuth(authHeader string) bool {
	if authHeader == "" {
		return len(h.authTokens) == 0
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return false
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	token := strings.TrimPrefix(parts[1], "Bearer ")
	for _, stored := range h.authTokens {
		if stored == token {
			return true
		}
	}
	return false
}

func (h *Handler) sendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := Response{Success: true, Data: data}
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) sendError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := Response{Success: false, Error: msg}
	json.NewEncoder(w).Encode(resp)
}

func (r *RateLimiter) isRateLimited(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.window)

	requests := r.requests[key]
	var valid []time.Time
	for _, t := range requests {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= r.limit {
		r.requests[key] = valid
		return true
	}

	r.requests[key] = append(valid, now)
	return false
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		state.GlobalState.Log("INFO", fmt.Sprintf("Webhook: %s %s", r.Method, r.URL.Path))
		next.ServeHTTP(w, r)
		state.GlobalState.Log("INFO", fmt.Sprintf("Webhook completed in %v", time.Since(start)))
	})
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.WebhookRequestsTotal.Inc()
		next.ServeHTTP(w, r)
	})
}

func RegisterDefaultHandlers() {
	h := New()
	h.AddMiddleware(LoggingMiddleware)
	h.AddMiddleware(MetricsMiddleware)

	h.Register("task.completed", handleTaskCompleted)
	h.Register("task.failed", handleTaskFailed)
	h.Register("task.started", handleTaskStarted)
	h.Register("task.progress", handleTaskProgress)
	h.Register("agent.started", handleAgentStarted)
	h.Register("agent.stopped", handleAgentStopped)
	h.Register("agent.error", handleAgentError)
	h.Register("session.created", handleSessionCreated)
	h.Register("session.ended", handleSessionEnded)
	h.Register("session.timeout", handleSessionTimeout)
	h.Register("plan.activated", handlePlanActivated)
	h.Register("plan.completed", handlePlanCompleted)
	h.Register("plan.failed", handlePlanFailed)
	h.Register("health.check", handleHealthCheck)
	h.Register("health.degraded", handleHealthDegraded)
	h.Register("model.acquired", handleModelAcquired)
	h.Register("model.released", handleModelReleased)
	h.Register("model.failed", handleModelFailed)
	h.Register("docker.container.started", handleDockerContainerStarted)
	h.Register("docker.container.stopped", handleDockerContainerStopped)
	h.Register("docker.container.failed", handleDockerContainerFailed)
	h.Register("scheduler.job.started", handleSchedulerJobStarted)
	h.Register("scheduler.job.completed", handleSchedulerJobCompleted)
	h.Register("scheduler.job.failed", handleSchedulerJobFailed)
	h.Register("notification.sent", handleNotificationSent)
	h.Register("notification.failed", handleNotificationFailed)
	h.Register("ratelimit.exceeded", handleRateLimitExceeded)
	h.Register("git.commit", handleGitCommit)
	h.Register("git.push", handleGitPush)
	h.Register("git.pull", handleGitPull)
	h.Register("cache.hit", handleCacheHit)
	h.Register("cache.miss", handleCacheMiss)
}

func handleTaskCompleted(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Task completed: %v", data["task_id"]))
	metrics.TasksCompletedTotal.Inc()
	return data, nil
}

func handleTaskFailed(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("ERROR", fmt.Sprintf("Task failed: %v", data["task_id"]))
	metrics.TasksFailedTotal.Inc()
	return data, nil
}

func handleAgentStarted(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Agent started: %v", data["agent"]))
	metrics.AgentsStartedTotal.Inc()
	return data, nil
}

func handleAgentStopped(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Agent stopped: %v", data["agent"]))
	metrics.AgentsStoppedTotal.Inc()
	return data, nil
}

func handleSessionCreated(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Session created: %v", data["session_id"]))
	metrics.SessionsCreatedTotal.Inc()
	return data, nil
}

func handleSessionEnded(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Session ended: %v", data["session_id"]))
	metrics.SessionsEndedTotal.Inc()
	return data, nil
}

func handlePlanActivated(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Plan activated: %v", data["plan_name"]))
	metrics.PlansActivatedTotal.Inc()
	return data, nil
}

func handlePlanCompleted(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Plan completed: %v", data["plan_name"]))
	metrics.PlansCompletedTotal.Inc()
	return data, nil
}

func handleHealthCheck(payload []byte) (interface{}, error) {
	return map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	}, nil
}

func handleModelAcquired(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Model acquired: %v", data["model"]))
	metrics.ModelAcquisitions.WithLabelValues(fmt.Sprintf("%v", data["model"])).Inc()
	return data, nil
}

func handleModelReleased(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Model released: %v", data["model"]))
	return data, nil
}

func handleTaskStarted(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Task started: %v", data["task_id"]))
	metrics.TasksStartedTotal.Inc()
	return data, nil
}

func handleTaskProgress(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Task progress: %v - %v%%", data["task_id"], data["progress"]))
	return data, nil
}

func handleAgentError(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("ERROR", fmt.Sprintf("Agent error: %v - %v", data["agent"], data["error"]))
	return data, nil
}

func handleSessionTimeout(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("WARN", fmt.Sprintf("Session timeout: %v", data["session_id"]))
	return data, nil
}

func handlePlanFailed(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("ERROR", fmt.Sprintf("Plan failed: %v", data["plan_name"]))
	return data, nil
}

func handleHealthDegraded(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("WARN", fmt.Sprintf("Health degraded: %v", data["component"]))
	return data, nil
}

func handleModelFailed(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("ERROR", fmt.Sprintf("Model failed: %v - %v", data["model"], data["error"]))
	return data, nil
}

func handleDockerContainerStarted(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Docker container started: %v", data["container"]))
	metrics.DockerContainerStartsTotal.Inc()
	return data, nil
}

func handleDockerContainerStopped(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Docker container stopped: %v", data["container"]))
	return data, nil
}

func handleDockerContainerFailed(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("ERROR", fmt.Sprintf("Docker container failed: %v - %v", data["container"], data["error"]))
	metrics.DockerContainerStartsFailedTotal.Inc()
	return data, nil
}

func handleSchedulerJobStarted(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Scheduler job started: %v", data["job"]))
	return data, nil
}

func handleSchedulerJobCompleted(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Scheduler job completed: %v", data["job"]))
	return data, nil
}

func handleSchedulerJobFailed(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("ERROR", fmt.Sprintf("Scheduler job failed: %v - %v", data["job"], data["error"]))
	return data, nil
}

func handleNotificationSent(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Notification sent: %v", data["channel"]))
	metrics.NotificationsSentTotal.Inc()
	return data, nil
}

func handleNotificationFailed(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("ERROR", fmt.Sprintf("Notification failed: %v - %v", data["channel"], data["error"]))
	metrics.NotificationsFailedTotal.Inc()
	return data, nil
}

func handleRateLimitExceeded(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("WARN", fmt.Sprintf("Rate limit exceeded: %v", data["key"]))
	return data, nil
}

func handleGitCommit(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Git commit: %v", data["message"]))
	return data, nil
}

func handleGitPush(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Git push: %v", data["branch"]))
	return data, nil
}

func handleGitPull(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Git pull: %v", data["branch"]))
	return data, nil
}

func handleCacheHit(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func handleCacheMiss(payload []byte) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func StartServer(addr string) error {
	RegisterDefaultHandlers()
	state.GlobalState.Log("INFO", "Starting webhook server on "+addr)
	return http.ListenAndServe(addr, webhookHandler)
}
