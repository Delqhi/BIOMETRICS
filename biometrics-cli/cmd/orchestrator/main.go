package main

import (
	"biometrics-cli/internal/cache"
	"biometrics-cli/internal/circuit"
	"biometrics-cli/internal/metrics"
	"biometrics-cli/internal/models"
	"biometrics-cli/internal/orchestrator"
	"biometrics-cli/internal/scheduler"
	"biometrics-cli/internal/selfhealing"
	"biometrics-cli/internal/skills"
	"biometrics-cli/internal/state"
	"biometrics-cli/internal/tracker"
	"biometrics-cli/internal/webhook"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	orchestratorRunning bool
	lastAgentHeartbeat  time.Time
	heartbeatTimeout    = 5 * time.Minute
	stuckDetectionChan  chan string
	mainAgentPID        int
	cycleCount          int
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	fmt.Println("=== BIOMETRICS 24/7 ORCHESTRATOR ===")
	fmt.Println("Starting autonomous agent orchestration...")

	state.GlobalState.InitDB()
	scheduler.Start()

	orchestrator.Init()
	skills.SelfTrain()

	webhook.RegisterDefaultHandlers()

	circuit.GetOrCreate("opencode-api", &circuit.CircuitBreakerConfig{
		Name: "opencode-api", MaxFailures: 5, Timeout: 30 * time.Second,
		HalfOpenMax: 3, ResetTimeout: 60 * time.Second,
	})

	cache.New(&cache.CacheConfig{
		DiskPath: "./cache", TTL: 5 * time.Minute, CleanupInterval: 1 * time.Minute,
	})

	modelPool := tracker.NewModelPool()

	go startMetricsServer()
	go heartbeatMonitor()
	go stuckAgentDetector()
	go webhookHealthPinger()

	fmt.Println("\n=== 24/7 ORCHESTRATOR RUNNING ===")
	fmt.Println("Monitoring: heartbeat, stuck detection, self-healing")
	fmt.Println("Press Ctrl+C to stop\n")

	for {
		cycleStart := time.Now()
		cycleCount++
		metrics.CyclesTotal.Inc()

		orchestratorRunning = true

		fmt.Printf("[%s] Cycle #%d starting...\n", time.Now().Format("15:04:05"), cycleCount)

		runSelfHealing()

		verifyOpenCodeHealth()

		processBoulderWithCircuitBreaker(modelPool)

		duration := time.Since(cycleStart).Seconds()
		metrics.CycleDuration.Observe(duration)
		fmt.Printf("[%s] Cycle completed in %.2fs\n", time.Now().Format("15:04:05"), duration)

		orchestratorRunning = false

		select {
		case agent := <-stuckDetectionChan:
			fmt.Printf("[ALERT] Stuck agent: %s - recovering...\n", agent)
			handleStuckAgent(agent)
		default:
		}

		time.Sleep(60 * time.Second)
	}
}

func runSelfHealing() {
	healer := selfhealing.NewSelfHealer()
	healer.RunDiagnostics()
}

func verifyOpenCodeHealth() {
	cb := circuit.GetOrCreate("opencode-api", &circuit.CircuitBreakerConfig{
		Name: "opencode-api", MaxFailures: 3,
	})

	err := cb.Execute(func() error {
		cmd := exec.Command("opencode", "--version")
		return cmd.Run()
	})

	if err != nil {
		metrics.HealingFailures.WithLabelValues("opencode").Inc()
		fmt.Println("[WARN] OpenCode unavailable")
	} else {
		metrics.HealingSuccesses.WithLabelValues("opencode").Inc()
	}
}

func processBoulderWithCircuitBreaker(mt *tracker.ModelPool) {
	boulder, err := readBoulderSafe()
	if err != nil {
		fmt.Printf("[WARN] Cannot read boulder: %v\n", err)
		return
	}

	if boulder.ActivePlan == "" {
		fmt.Println("[IDLE] No active plan - creating tasks...")
		orchestrator.DefaultOrchestrator.AutoCreateTodos()
		return
	}

	agent := orchestrator.DefaultOrchestrator.GetAgentForTask("default")
	if agent == nil {
		fmt.Println("[ERROR] No agent available")
		return
	}

	state.GlobalState.CurrentAgent = agent.Name
	state.GlobalState.PlanName = boulder.PlanName
	lastAgentHeartbeat = time.Now()

	models := []string{"qwen3.5", "minimax", "kimi", "gemini"}
	var acquiredModel string
	for _, m := range models {
		acquired, err := mt.AcquireModel(m)
		if acquired != nil && err == nil {
			acquiredModel = m
			break
		}
	}

	if acquiredModel == "" {
		fmt.Println("[WARN] No models available")
		return
	}

	metrics.ModelAcquisitions.WithLabelValues(acquiredModel).Inc()
	fmt.Printf("[AGENT] Running: %s with %s\n", agent.Name, acquiredModel)

	runAgentWithMonitoring(agent.Name, boulder)

	mt.ReleaseModel(acquiredModel)
	state.GlobalState.ActiveModel = "NONE"
}

func readBoulderSafe() (*models.Boulder, error) {
	if state.GlobalState.GetChaos() && rand.Intn(10) < 3 {
		return nil, fmt.Errorf("chaos")
	}

	data, err := os.ReadFile("/Users/jeremy/.sisyphus/boulder.json")
	if err != nil {
		return nil, err
	}
	var b models.Boulder
	err = json.Unmarshal(data, &b)
	return &b, err
}

func runAgentWithMonitoring(agentName string, boulder *models.Boulder) {
	cmd := exec.Command("opencode", boulder.PlanName, "--agent", agentName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	startTime := time.Now()
	err := cmd.Start()
	if err != nil {
		fmt.Printf("[ERROR] Failed to start agent: %v\n", err)
		return
	}

	mainAgentPID = cmd.Process.Pid

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			fmt.Printf("[AGENT] Finished with error: %v\n", err)
			metrics.TasksFailedTotal.Inc()
		} else {
			fmt.Println("[AGENT] Completed successfully")
			metrics.TasksCompletedTotal.Inc()
		}
	case <-time.After(10 * time.Minute):
		fmt.Println("[TIMEOUT] Agent running > 10 min - may be stuck")
		stuckDetectionChan <- agentName
	case <-time.After(heartbeatTimeout):
		fmt.Printf("[HEARTBEAT] No heartbeat for %v\n", heartbeatTimeout)
		stuckDetectionChan <- agentName
	}

	metrics.AgentDuration.WithLabelValues(agentName, "completed").Observe(time.Since(startTime).Seconds())
}

func handleStuckAgent(agentName string) {
	fmt.Printf("[RECOVERY] Handling stuck agent: %s\n", agentName)

	if mainAgentPID > 0 {
		syscall.Kill(mainAgentPID, syscall.SIGKILL)
		fmt.Printf("[RECOVERY] Killed PID: %d\n", mainAgentPID)
	}

	metrics.HealingFailures.WithLabelValues("stuck-agent").Inc()

	webhook.SendTaskFailed(agentName, "orchestrator", "Agent stuck - auto-recovered")

	fmt.Printf("[RECOVERY] Agent %s recovered\n", agentName)
}

func heartbeatMonitor() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		if orchestratorRunning && time.Since(lastAgentHeartbeat) > heartbeatTimeout {
			fmt.Printf("[HEARTBEAT] WARNING - No heartbeat for %v\n", time.Since(lastAgentHeartbeat))
			stuckDetectionChan <- "main-agent"
		}
	}
}

func stuckAgentDetector() {
	stuckDetectionChan = make(chan string, 1)
	for {
		select {
		case agent := <-stuckDetectionChan:
			handleStuckAgent(agent)
		case <-time.After(30 * time.Second):
			if orchestratorRunning {
				runSicherCheck()
			}
		}
	}
}

func webhookHealthPinger() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		req, _ := http.NewRequestWithContext(ctx, "GET", "http://localhost:59002/health", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			metrics.HealingFailures.WithLabelValues("webhook-health").Inc()
		} else {
			resp.Body.Close()
			metrics.HealingSuccesses.WithLabelValues("webhook-health").Inc()
		}
		cancel()
	}
}

func runSicherCheck() error {
	cmd := exec.Command("opencode", "--yes", "Sicher? Check if everything works.", "--agent", "sisyphus")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func startMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","orchestrator_running":true}`))
	})
	http.HandleFunc("/stuck-agents", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"last_heartbeat": lastAgentHeartbeat,
			"timeout":        heartbeatTimeout,
			"is_stuck":       time.Since(lastAgentHeartbeat) > heartbeatTimeout,
		})
	})
	go http.ListenAndServe(":59002", nil)
}
