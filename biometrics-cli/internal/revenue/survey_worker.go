package revenue

import (
	"context"
	"fmt"
	"os"
	"time"

	"biometrics-cli/internal/metrics"
)

// SurveyPlatform represents a survey platform integration
type SurveyPlatform struct {
	Name      string
	BaseURL   string
	Username  string
	Password  string
	MinPayout float64
	Active    bool
}

// SurveyWorker handles automated survey completion
type SurveyWorker struct {
	Platforms []SurveyPlatform
	Client    *RevenueClient
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewSurveyWorker creates a new survey automation worker
func NewSurveyWorker() *SurveyWorker {
	ctx, cancel := context.WithCancel(context.Background())

	worker := &SurveyWorker{
		Platforms: []SurveyPlatform{
			{
				Name:      "Prolific",
				BaseURL:   "https://app.prolific.com",
				Username:  os.Getenv("PROLIFIC_USERNAME"),
				Password:  os.Getenv("PROLIFIC_PASSWORD"),
				MinPayout: 0.50,
				Active:    true,
			},
			{
				Name:      "Swagbucks",
				BaseURL:   "https://www.swagbucks.com",
				Username:  os.Getenv("SWAGBUCKS_USERNAME"),
				Password:  os.Getenv("SWAGBUCKS_PASSWORD"),
				MinPayout: 0.25,
				Active:    true,
			},
		},
		Client: NewRevenueClient(),
		ctx:    ctx,
		cancel: cancel,
	}

	return worker
}

// Start begins the survey automation loop
func (w *SurveyWorker) Start() error {
	metrics.RevenueActiveWorkers.Inc()
	defer metrics.RevenueActiveWorkers.Dec()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return fmt.Errorf("worker stopped")
		case <-ticker.C:
			for _, platform := range w.Platforms {
				if !platform.Active {
					continue
				}

				startTime := time.Now()
				err := w.checkAndCompleteSurveys(platform)
				duration := time.Since(startTime).Seconds()

				if err != nil {
					metrics.RevenueAPIErrors.WithLabelValues(platform.Name, "survey_check").Inc()
					metrics.RevenueTaskDuration.WithLabelValues("survey", "error").Observe(duration)
				} else {
					metrics.RevenueTaskDuration.WithLabelValues("survey", "success").Observe(duration)
				}
			}
		}
	}
}

// checkAndCompleteSurveys checks for available surveys and completes them
func (w *SurveyWorker) checkAndCompleteSurveys(platform SurveyPlatform) error {
	// TODO: Implement Steel Browser automation
	// 1. Navigate to platform
	// 2. Login with session persistence
	// 3. Check available surveys
	// 4. Qualify for surveys
	// 5. Complete with AI assistance
	// 6. Verify payment

	// Placeholder for now
	fmt.Printf("[%s] Checking for surveys...\n", platform.Name)

	// Simulate survey completion
	metrics.RevenueTasksCompleted.WithLabelValues("survey").Inc()
	metrics.RevenueEarningsTotal.Add(0.50) // Average survey payout

	return nil
}

// Stop gracefully stops the worker
func (w *SurveyWorker) Stop() {
	w.cancel()
}

// GetActivePlatforms returns list of active platforms
func (w *SurveyWorker) GetActivePlatforms() []string {
	var active []string
	for _, p := range w.Platforms {
		if p.Active {
			active = append(active, p.Name)
		}
	}
	return active
}

// ConfigurePlatform updates platform configuration
func (w *SurveyWorker) ConfigurePlatform(name string, active bool) error {
	for i, p := range w.Platforms {
		if p.Name == name {
			w.Platforms[i].Active = active
			return nil
		}
	}
	return fmt.Errorf("platform not found: %s", name)
}
