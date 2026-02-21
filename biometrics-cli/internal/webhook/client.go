package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type WebhookEvent struct {
	Type     string                 `json:"type"`
	Source   string                 `json:"source"`
	Agent    string                 `json:"agent"`
	Data     map[string]interface{} `json:"data"`
	Metadata map[string]string      `json:"metadata"`
	Time     time.Time              `json:"time"`
}

type WebhookClient struct {
	URL     string
	Secret  string
	Client  *http.Client
	Retries int
	Timeout time.Duration
}

func NewWebhookClient(url, secret string) *WebhookClient {
	return &WebhookClient{
		URL:    url,
		Secret: secret,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		Retries: 3,
		Timeout: 30 * time.Second,
	}
}

func (c *WebhookClient) Send(event *WebhookEvent) error {
	event.Time = time.Now()

	payload, err := encodeEvent(event)
	if err != nil {
		return fmt.Errorf("failed to encode event: %w", err)
	}

	signature := c.sign(payload)

	req, err := http.NewRequest("POST", c.URL, strings.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Signature", signature)
	req.Header.Set("X-Webhook-Timestamp", fmt.Sprintf("%d", event.Time.Unix()))

	var lastErr error
	for i := 0; i < c.Retries; i++ {
		resp, err := c.Client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}

		lastErr = fmt.Errorf("webhook returned status %d", resp.StatusCode)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return lastErr
}

func (c *WebhookClient) sign(payload string) string {
	mac := hmac.New(sha256.New, []byte(c.Secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func encodeEvent(event *WebhookEvent) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`{"type":"%s","source":"%s","agent":"%s","data":`,
		event.Type, event.Source, event.Agent))

	dataStr := "{"
	first := true
	for k, v := range event.Data {
		if !first {
			dataStr += ","
		}
		dataStr += fmt.Sprintf(`"%v":"%v"`, k, v)
		first = false
	}
	dataStr += "}"
	sb.WriteString(dataStr)
	sb.WriteString(fmt.Sprintf(`,"metadata":{},"time":"%s"}`, event.Time.Format(time.RFC3339)))

	return sb.String(), nil
}

func VerifySignature(payload, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

func SendTaskStarted(taskID, agent string) error {
	event := &WebhookEvent{
		Type:   "task.started",
		Source: "biometrics",
		Agent:  agent,
		Data: map[string]interface{}{
			"task_id": taskID,
		},
	}
	return sendWebhookEvent(event)
}

func SendTaskCompleted(taskID, agent, result string) error {
	event := &WebhookEvent{
		Type:   "task.completed",
		Source: "biometrics",
		Agent:  agent,
		Data: map[string]interface{}{
			"task_id": taskID,
			"result":  result,
		},
	}
	return sendWebhookEvent(event)
}

func SendTaskFailed(taskID, agent, error string) error {
	event := &WebhookEvent{
		Type:   "task.failed",
		Source: "biometrics",
		Agent:  agent,
		Data: map[string]interface{}{
			"task_id": taskID,
			"error":   error,
		},
	}
	return sendWebhookEvent(event)
}

func SendAgentStarted(agent, model string) error {
	event := &WebhookEvent{
		Type:   "agent.started",
		Source: "biometrics",
		Agent:  agent,
		Data: map[string]interface{}{
			"model": model,
		},
	}
	return sendWebhookEvent(event)
}

func SendAgentStopped(agent string) error {
	event := &WebhookEvent{
		Type:   "agent.stopped",
		Source: "biometrics",
		Agent:  agent,
		Data:   map[string]interface{}{},
	}
	return sendWebhookEvent(event)
}

func sendWebhookEvent(event *WebhookEvent) error {
	webhookURL := getWebhookURL()
	if webhookURL == "" {
		return nil
	}

	client := NewWebhookClient(webhookURL, getWebhookSecret())
	return client.Send(event)
}

func getWebhookURL() string {
	return ""
}

func getWebhookSecret() string {
	return ""
}
