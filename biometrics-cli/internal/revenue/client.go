package revenue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// RevenueTask represents a money-generating task
type RevenueTask struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // captcha, survey, usertesting
	Platform    string    `json:"platform"`
	Status      string    `json:"status"`
	Earnings    float64   `json:"earnings"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
}

// CaptchaSolution represents a solved captcha
type CaptchaSolution struct {
	TaskID     string  `json:"task_id"`
	Solution   string  `json:"solution"`
	Confidence float64 `json:"confidence"`
	TimeMs     int64   `json:"time_ms"`
}

// RevenueClient handles revenue-generating API calls
type RevenueClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewRevenueClient creates a new revenue client
func NewRevenueClient() *RevenueClient {
	return &RevenueClient{
		BaseURL: os.Getenv("REVENUE_API_URL"),
		APIKey:  os.Getenv("REVENUE_API_KEY"),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SubmitCaptchaSolution submits a solved captcha for payment
func (c *RevenueClient) SubmitCaptchaSolution(solution *CaptchaSolution) error {
	payload, err := json.Marshal(solution)
	if err != nil {
		return fmt.Errorf("marshal solution: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/api/v1/captcha/submit", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("submit solution: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("submit failed (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetEarnings retrieves current earnings
func (c *RevenueClient) GetEarnings() (float64, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/api/v1/earnings", nil)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("get earnings: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Total float64 `json:"total"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("decode response: %w", err)
	}

	return result.Total, nil
}

// VerifyPayment verifies a payment was received
func (c *RevenueClient) VerifyPayment(taskID string) (bool, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/api/v1/payments/"+taskID, nil)
	if err != nil {
		return false, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("verify payment: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
