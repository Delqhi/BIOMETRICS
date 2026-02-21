package revenue

import (
	"context"
	"fmt"
	"os"
	"time"

	"biometrics-cli/internal/metrics"
)

// TestingPlatform represents a user testing platform
type TestingPlatform struct {
	Name         string
	BaseURL      string
	Username     string
	Password     string
	Payout       float64 // Fixed payout per test
	TestDuration int     // Average test duration in minutes
	Active       bool
}

// TestingWorker handles automated website testing
type TestingWorker struct {
	Platforms []TestingPlatform
	Client    *RevenueClient
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewTestingWorker creates a new website testing worker
func NewTestingWorker() *TestingWorker {
	ctx, cancel := context.WithCancel(context.Background())

	worker := &TestingWorker{
		Platforms: []TestingPlatform{
			{
				Name:         "UserTesting",
				BaseURL:      "https://www.usertesting.com",
				Username:     os.Getenv("USERTESTING_USERNAME"),
				Password:     os.Getenv("USERTESTING_PASSWORD"),
				Payout:       10.00,
				TestDuration: 15,
				Active:       true,
			},
			{
				Name:         "TryMyUI",
				BaseURL:      "https://www.trymyui.com",
				Username:     os.Getenv("TRYMYUI_USERNAME"),
				Password:     os.Getenv("TRYMYUI_PASSWORD"),
				Payout:       10.00,
				TestDuration: 20,
				Active:       true,
			},
		},
		Client: NewRevenueClient(),
		ctx:    ctx,
		cancel: cancel,
	}

	return worker
}

// Start begins the website testing automation loop
func (w *TestingWorker) Start() error {
	metrics.RevenueActiveWorkers.Inc()
	defer metrics.RevenueActiveWorkers.Dec()

	ticker := time.NewTicker(10 * time.Minute) // Check less frequently
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
				err := w.checkAndCompleteTests(platform)
				duration := time.Since(startTime).Seconds()

				if err != nil {
					metrics.RevenueAPIErrors.WithLabelValues(platform.Name, "test_check").Inc()
					metrics.RevenueTaskDuration.WithLabelValues("testing", "error").Observe(duration)
				} else {
					metrics.RevenueTaskDuration.WithLabelValues("testing", "success").Observe(duration)
				}
			}
		}
	}
}

// checkAndCompleteTests checks for available tests and completes them
func (w *TestingWorker) checkAndCompleteTests(platform TestingPlatform) error {
	// TODO: Implement Steel Browser automation
	// 1. Navigate to platform
	// 2. Login with session persistence
	// 3. Check available tests
	// 4. Qualify (demographics screener)
	// 5. Record screen + voice (FFmpeg)
	// 6. Complete test tasks
	// 7. Generate AI feedback
	// 8. Submit video + feedback
	// 9. Track payment (7-day delay)

	// Placeholder for now
	fmt.Printf("[%s] Checking for tests...\n", platform.Name)

	// Simulate test completion (rare - maybe 1-5 per day)
	metrics.RevenueTasksCompleted.WithLabelValues("testing").Inc()
	metrics.RevenueEarningsTotal.Add(platform.Payout)
	metrics.RevenueAverageEarningsPerTask.Set(platform.Payout)

	return nil
}

// Stop gracefully stops the worker
func (w *TestingWorker) Stop() {
	w.cancel()
}

// GetActivePlatforms returns list of active testing platforms
func (w *TestingWorker) GetActivePlatforms() []string {
	var active []string
	for _, p := range w.Platforms {
		if p.Active {
			active = append(active, p.Name)
		}
	}
	return active
}

// ConfigurePlatform updates platform configuration
func (w *TestingWorker) ConfigurePlatform(name string, active bool) error {
	for i, p := range w.Platforms {
		if p.Name == name {
			w.Platforms[i].Active = active
			return nil
		}
	}
	return fmt.Errorf("platform not found: %s", name)
}

// EstimateDailyEarnings calculates potential daily earnings
func (w *TestingWorker) EstimateDailyEarnings() map[string]float64 {
	estimates := make(map[string]float64)

	for _, p := range w.Platforms {
		if !p.Active {
			continue
		}

		// Conservative estimate: 1 test per day
		// Optimistic: 5 tests per day
		dailyTests := 1.0
		estimates[p.Name] = dailyTests * p.Payout
	}

	return estimates
}
