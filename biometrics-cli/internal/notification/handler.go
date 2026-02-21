package notification

import (
	"biometrics-cli/internal/metrics"
	"biometrics-cli/internal/state"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Notification struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Priority  string                 `json:"priority"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Retries   int                    `json:"retries"`
	Status    string                 `json:"status"`
}

type Channel interface {
	Send(n *Notification) error
	GetName() string
}

type Handler struct {
	mu       sync.RWMutex
	channels map[string]Channel
	queue    chan *Notification
	workers  int
	stopChan chan struct{}
	wg       sync.WaitGroup
}

var (
	defaultHandler = &Handler{
		channels: make(map[string]Channel),
		queue:    make(chan *Notification, 1000),
		workers:  3,
		stopChan: make(chan struct{}),
	}
	HandlerInstance = defaultHandler
)

func New(workers int) *Handler {
	if workers > 0 {
		defaultHandler.workers = workers
	}
	return defaultHandler
}

func (h *Handler) RegisterChannel(channel Channel) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.channels[channel.GetName()] = channel
	state.GlobalState.Log("INFO", fmt.Sprintf("Registered notification channel: %s", channel.GetName()))
}

func (h *Handler) UnregisterChannel(name string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.channels, name)
	state.GlobalState.Log("INFO", fmt.Sprintf("Unregistered notification channel: %s", name))
}

func (h *Handler) Send(n *Notification) error {
	if n.Timestamp.IsZero() {
		n.Timestamp = time.Now()
	}
	if n.ID == "" {
		n.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	metrics.NotificationsSentTotal.Inc()

	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.channels) == 0 {
		state.GlobalState.Log("WARN", "No notification channels registered")
		return fmt.Errorf("no channels registered")
	}

	var lastErr error
	for _, channel := range h.channels {
		if err := channel.Send(n); err != nil {
			state.GlobalState.Log("ERROR", fmt.Sprintf("Failed to send via %s: %v", channel.GetName(), err))
			lastErr = err
			metrics.NotificationsFailedTotal.Inc()
		}
	}

	return lastErr
}

func (h *Handler) SendAsync(n *Notification) {
	select {
	case h.queue <- n:
	default:
		state.GlobalState.Log("WARN", "Notification queue full, dropping notification")
		metrics.NotificationsDroppedTotal.Inc()
	}
}

func (h *Handler) Start() {
	for i := 0; i < h.workers; i++ {
		h.wg.Add(1)
		go h.worker(i)
	}
	state.GlobalState.Log("INFO", fmt.Sprintf("Started %d notification workers", h.workers))
}

func (h *Handler) Stop() {
	close(h.stopChan)
	h.wg.Wait()
	state.GlobalState.Log("INFO", "Notification workers stopped")
}

func (h *Handler) worker(id int) {
	defer h.wg.Done()

	for {
		select {
		case <-h.stopChan:
			return
		case n := <-h.queue:
			if err := h.Send(n); err != nil {
				if n.Retries < 3 {
					n.Retries++
					time.Sleep(time.Duration(n.Retries) * time.Second)
					h.queue <- n
				}
			}
		}
	}
}

type DiscordChannel struct {
	webhookURL string
	botToken   string
}

func NewDiscordChannel(webhookURL, botToken string) *DiscordChannel {
	return &DiscordChannel{
		webhookURL: webhookURL,
		botToken:   botToken,
	}
}

func (c *DiscordChannel) GetName() string {
	return "discord"
}

func (c *DiscordChannel) Send(n *Notification) error {
	color := 0x00ff00
	switch n.Priority {
	case "high":
		color = 0xff0000
	case "medium":
		color = 0xffa500
	case "low":
		color = 0x0000ff
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       n.Title,
				"description": n.Message,
				"color":       color,
				"timestamp":   n.Timestamp.Format(time.RFC3339),
				"fields":      []map[string]interface{}{},
			},
		},
	}

	if len(n.Data) > 0 {
		for k, v := range n.Data {
			payload["embeds"].([]map[string]interface{})[0]["fields"] = append(
				payload["embeds"].([]map[string]interface{})[0]["fields"].([]map[string]interface{}),
				map[string]interface{}{
					"name":  k,
					"value": fmt.Sprintf("%v", v),
				},
			)
		}
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", c.webhookURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("discord API error: %d", resp.StatusCode)
	}

	return nil
}

type SlackChannel struct {
	webhookURL string
}

func NewSlackChannel(webhookURL string) *SlackChannel {
	return &SlackChannel{
		webhookURL: webhookURL,
	}
}

func (c *SlackChannel) GetName() string {
	return "slack"
}

func (c *SlackChannel) Send(n *Notification) error {
	emoji := ":white_check_mark:"
	switch n.Priority {
	case "high":
		emoji = ":fire:"
	case "medium":
		emoji = ":warning:"
	case "low":
		emoji = ":information_source:"
	}

	payload := map[string]interface{}{
		"text": fmt.Sprintf("%s *%s*", emoji, n.Title),
		"blocks": []map[string]interface{}{
			{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*%s*\n%s", n.Title, n.Message),
				},
			},
		},
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", c.webhookURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("slack API error: %d", resp.StatusCode)
	}

	return nil
}

type EmailChannel struct {
	smtpHost string
	smtpPort int
	username string
	password string
	from     string
	to       []string
}

func NewEmailChannel(smtpHost, smtpPort, username, password, from string, to []string) *EmailChannel {
	port := 587
	fmt.Sscanf(smtpPort, "%d", &port)
	return &EmailChannel{
		smtpHost: smtpHost,
		smtpPort: port,
		username: username,
		password: password,
		from:     from,
		to:       to,
	}
}

func (c *EmailChannel) GetName() string {
	return "email"
}

func (c *EmailChannel) Send(n *Notification) error {
	state.GlobalState.Log("INFO", fmt.Sprintf("Email notification: %s - %s", n.Title, n.Message))
	return nil
}

type LogChannel struct{}

func NewLogChannel() *LogChannel {
	return &LogChannel{}
}

func (c *LogChannel) GetName() string {
	return "log"
}

func (c *LogChannel) Send(n *Notification) error {
	switch n.Priority {
	case "high":
		state.GlobalState.Log("ERROR", fmt.Sprintf("[NOTIFICATION] %s: %s", n.Title, n.Message))
	case "medium":
		state.GlobalState.Log("WARN", fmt.Sprintf("[NOTIFICATION] %s: %s", n.Title, n.Message))
	default:
		state.GlobalState.Log("INFO", fmt.Sprintf("[NOTIFICATION] %s: %s", n.Title, n.Message))
	}
	return nil
}

func Notify(title, message, priority string) error {
	n := &Notification{
		Type:     "general",
		Title:    title,
		Message:  message,
		Priority: priority,
	}
	return HandlerInstance.Send(n)
}

func NotifyTaskComplete(taskID, result string) error {
	n := &Notification{
		Type:     "task",
		Title:    "Task Completed",
		Message:  fmt.Sprintf("Task %s: %s", taskID, result),
		Priority: "low",
		Data: map[string]interface{}{
			"task_id": taskID,
			"result":  result,
		},
	}
	return HandlerInstance.Send(n)
}

func NotifyError(component, err string) error {
	n := &Notification{
		Type:     "error",
		Title:    "Error Detected",
		Message:  fmt.Sprintf("Component %s: %s", component, err),
		Priority: "high",
		Data: map[string]interface{}{
			"component": component,
			"error":     err,
		},
	}
	return HandlerInstance.Send(n)
}

func NotifyAgentStarted(agent, sessionID string) error {
	n := &Notification{
		Type:     "agent",
		Title:    "Agent Started",
		Message:  fmt.Sprintf("Agent %s started (session: %s)", agent, sessionID),
		Priority: "low",
		Data: map[string]interface{}{
			"agent":      agent,
			"session_id": sessionID,
		},
	}
	return HandlerInstance.Send(n)
}

func NotifyPlanCompleted(planName string) error {
	n := &Notification{
		Type:     "plan",
		Title:    "Plan Completed",
		Message:  fmt.Sprintf("Plan '%s' has been completed", planName),
		Priority: "medium",
		Data: map[string]interface{}{
			"plan_name": planName,
		},
	}
	return HandlerInstance.Send(n)
}
