package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cyclesTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_orchestrator_cycles_total",
		Help: "The total number of processed cycles",
	})
	cycleDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "biometrics_orchestrator_cycle_duration_seconds",
		Help:    "Duration of cycles in seconds",
		Buckets: prometheus.DefBuckets,
	})
	modelAcquisitions = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "biometrics_orchestrator_model_acquisitions_total",
		Help: "The total number of model acquisitions by model name",
	}, []string{"model"})
	serenaSessionsCleaned = promauto.NewCounter(prometheus.CounterOpts{
		Name: "biometrics_orchestrator_serena_sessions_cleaned_total",
		Help: "Total number of Serena sessions cleaned up",
	})
)

type Boulder struct {
	ActivePlan string   `json:"active_plan"`
	StartedAt  string   `json:"started_at"`
	SessionIDs []string `json:"session_ids"`
	PlanName   string   `json:"plan_name"`
	Agent      string   `json:"agent"`
}

type AppState struct {
	mu           sync.Mutex
	ActivePlan   string
	PlanName     string
	CurrentAgent string
	ActiveModel  string
	ModelStatus  map[string]string
	Logs         []string
	db           *sql.DB
}

var state = &AppState{
	ModelStatus: make(map[string]string),
	Logs:        make([]string, 0),
}

func (s *AppState) InitDB() {
	dbPath := "/Users/jeremy/.sisyphus/logs.db"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return
	}
	query := "CREATE TABLE IF NOT EXISTS logs (id INTEGER PRIMARY KEY AUTOINCREMENT, timestamp TEXT, level TEXT, agent TEXT, plan TEXT, message TEXT)"
	_, _ = db.Exec(query)
	s.db = db
}

func (s *AppState) Log(level, msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ts := time.Now().Format("15:04:05")
	s.Logs = append(s.Logs, fmt.Sprintf("[%s] %s: %s", ts, level, msg))
	if len(s.Logs) > 10 {
		s.Logs = s.Logs[1:]
	}
	if s.db != nil {
		_, _ = s.db.Exec("INSERT INTO logs (timestamp, level, agent, plan, message) VALUES (?, ?, ?, ?, ?)",
			time.Now().Format(time.RFC3339), level, s.CurrentAgent, s.PlanName, msg)
	}
}

func displayDashboard() {
	for {
		state.mu.Lock()
		fmt.Print("\033[H\033[2J")
		fmt.Println("==============================================================")
		fmt.Println("         BIOMETRICS ENTERPRISE ORCHESTRATOR DASHBOARD         ")
		fmt.Println("==============================================================")
		fmt.Printf("STATUS:     RUNNING (24/7 MODE)\n")
		fmt.Printf("METRICS:    :59002/metrics\n")
		fmt.Printf("PLAN:       %s\n", state.PlanName)
		fmt.Printf("AGENT:      %s\n", state.CurrentAgent)
		fmt.Printf("MODEL:      %s\n", state.ActiveModel)
		fmt.Println("--------------------------------------------------------------")
		fmt.Println("MODEL STATUS:")
		for m, s := range state.ModelStatus {
			fmt.Printf("  %-15s : %s\n", m, s)
		}
		fmt.Println("--------------------------------------------------------------")
		fmt.Println("RECENT LOGS:")
		for _, l := range state.Logs {
			fmt.Println("  " + l)
		}
		fmt.Println("==============================================================")
		state.mu.Unlock()
		time.Sleep(2 * time.Second)
	}
}

type ModelTracker struct {
	mu     sync.Mutex
	models map[string]bool
}

func NewModelTracker() *ModelTracker {
	return &ModelTracker{models: make(map[string]bool)}
}

func (mt *ModelTracker) Acquire(model string) error {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	if mt.models[model] {
		return fmt.Errorf("in use")
	}
	mt.models[model] = true
	return nil
}

func (mt *ModelTracker) Release(model string) {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	delete(mt.models, model)
}

func runDoctor() {
	fmt.Println("=== BIOMETRICS DOCTOR ===")
	paths := []string{"/Users/jeremy/.sisyphus", "/Users/jeremy/.config/opencode/opencode.json"}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			fmt.Printf("OK: %s\n", p)
		} else {
			fmt.Printf("ERROR: %s\n", p)
		}
	}
}

func runSicherCheck(agent string) {
	prompt := "Sicher? Führe eine vollständige Selbstreflexion durch."
	_ = exec.Command("opencode", "prompt", prompt, "--agent", agent).Run()
}

func getModelForAgent(agent string) string {
	switch agent {
	case "sisyphus", "build", "atlas", "deep", "oracle", "ultrabrain", "visual-engineering":
		return "qwen3.5"
	case "librarian", "explore", "quick", "metis", "momus":
		return "minimax"
	default:
		return "qwen3.5"
	}
}

func readBoulder(path string) (*Boulder, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var b Boulder
	err = json.Unmarshal(data, &b)
	return &b, err
}

type SerenaSession struct {
	ID       string `json:"id"`
	Project  string `json:"project"`
	LastUsed string `json:"last_used"`
	Status   string `json:"status"`
}

func cleanupInactiveSerenaSessions() {
	cmd := exec.Command("uvx", "--from", "git+https://github.com/oraios/serena", "serena", "session", "list", "--json")
	output, err := cmd.Output()
	if err != nil {
		return
	}
	var sessions []SerenaSession
	_ = json.Unmarshal(output, &sessions)
	for _, s := range sessions {
		lastUsed, _ := time.Parse(time.RFC3339, s.LastUsed)
		if s.Status == "inactive" || time.Since(lastUsed) > 7*24*time.Hour || s.Project == "" || s.Project == "default" {
			_ = exec.Command("uvx", "--from", "git+https://github.com/oraios/serena", "serena", "session", "archive", s.ID).Run()
			serenaSessionsCleaned.Inc()
		}
	}
}

func verifySerenaProcess() error {
	return exec.Command("pgrep", "-f", "serena.*start-mcp-server").Run()
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "doctor" {
		runDoctor()
		return
	}
	state.InitDB()
	go displayDashboard()
	
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":59002", nil)
	}()

	tracker := NewModelTracker()
	state.Log("INFO", "Started")
	for {
		start := time.Now()
		cyclesTotal.Inc()

		if err := verifySerenaProcess(); err != nil {
			state.Log("ERROR", "No Serena")
			time.Sleep(10 * time.Second)
			continue
		}
		cleanupInactiveSerenaSessions()
		b, err := readBoulder("/Users/jeremy/.sisyphus/boulder.json")
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}
		state.mu.Lock()
		state.PlanName = b.PlanName
		state.CurrentAgent = b.Agent
		state.mu.Unlock()
		if b.ActivePlan == "" {
			time.Sleep(10 * time.Second)
			continue
		}
		model := getModelForAgent(b.Agent)
		if err := tracker.Acquire(model); err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		state.mu.Lock()
		state.ActiveModel = model
		state.mu.Unlock()
		
		modelAcquisitions.WithLabelValues(model).Inc()
		state.Log("SUCCESS", "Acquired "+model)
		
		runSicherCheck(b.Agent)
		tracker.Release(model)
		
		state.mu.Lock()
		state.ActiveModel = "NONE"
		state.mu.Unlock()

		duration := time.Since(start).Seconds()
		cycleDuration.Observe(duration)

		time.Sleep(60 * time.Second)
	}
}
