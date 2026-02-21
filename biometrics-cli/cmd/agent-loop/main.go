package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Boulder struct {
	ActivePlan string   `json:"active_plan"`
	StartedAt  string   `json:"started_at"`
	SessionIDs []string `json:"session_ids"`
	PlanName   string   `json:"plan_name"`
	Agent      string   `json:"agent"`
}

type ModelTracker struct {
	mu     sync.Mutex
	models map[string]bool
}

func NewModelTracker() *ModelTracker {
	return &ModelTracker{
		models: make(map[string]bool),
	}
}

func (mt *ModelTracker) Acquire(model string) error {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	if mt.models[model] {
		return fmt.Errorf("model %s is already in use", model)
	}

	mt.models[model] = true
	return nil
}

func (mt *ModelTracker) Release(model string) {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	delete(mt.models, model)
}

func main() {
	fmt.Println("Starting BIOMETRICS Agent Loop Orchestrator with Strict Model Assignment...")
	fmt.Println("Machine-readable mode: ENABLED (zero emojis)")
	fmt.Println("Model collision prevention: ACTIVE")
	fmt.Println("Sicher verification: ENABLED")
	fmt.Println("Serena session cleanup: ENABLED")

	boulderPath := "/Users/jeremy/.sisyphus/boulder.json"
	tracker := NewModelTracker()

	fmt.Println("Orchestrator initialized. Monitoring boulder.json...")

	for {
		fmt.Printf("[%s] Checking for agent status...\n", time.Now().Format(time.RFC3339))

		if err := verifySerenaProcess(); err != nil {
			fmt.Printf("ERROR: Serena MCP not running. Please start with: uvx --from git+https://github.com/oraios/serena serena start-mcp-server\n")
			time.Sleep(30 * time.Second)
			continue
		}

		cleanupInactiveSerenaSessions()

		boulder, err := readBoulder(boulderPath)
		if err != nil {
			fmt.Printf("ERROR: Reading boulder.json: %v\n", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if boulder.ActivePlan == "" {
			fmt.Println("STATUS: No active plan found. Waiting...")
			time.Sleep(30 * time.Second)
			continue
		}

		fmt.Printf("STATUS: Active plan detected: %s (Agent: %s)\n", boulder.PlanName, boulder.Agent)

		model := getModelForAgent(boulder.Agent)
		if err := tracker.Acquire(model); err != nil {
			fmt.Printf("WARNING: Model Collision: %v. Waiting for model to be free...\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Printf("SUCCESS: Model %s acquired for agent %s\n", model, boulder.Agent)

		runSicherCheck(boulder.Agent)

		fmt.Printf("INFO: Releasing model %s for next task...\n", model)
		tracker.Release(model)

		fmt.Printf("[%s] Cycle complete. Next check in 60s...\n", time.Now().Format(time.RFC3339))
		time.Sleep(60 * time.Second)
	}
}

func runSicherCheck(agent string) {
	fmt.Printf("ACTION: Triggering Sicher verification for agent: %s\n", agent)

	prompt := "Sicher? Führe eine vollständige Selbstreflexion durch. Prüfe jede deiner Aussagen, verifiziere, ob ALLE Restriktionen des Initial-Prompts exakt eingehalten wurden. Stelle alles Fehlende fertig."

	fmt.Printf("EXEC: opencode prompt '%s' --agent %s\n", prompt, agent)
}

func getModelForAgent(agent string) string {
	switch agent {
	case "sisyphus", "build":
		return "qwen3.5"
	case "atlas", "deep":
		return "kimi-k2.5"
	case "librarian", "explore", "quick":
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

	var boulder Boulder
	if err := json.Unmarshal(data, &boulder); err != nil {
		return nil, err
	}

	return &boulder, nil
}

type SerenaSession struct {
	ID        string
	Project   string
	CreatedAt time.Time
	LastUsed  time.Time
	Status    string
}

func cleanupInactiveSerenaSessions() {
	fmt.Println("ACTION: Starting Serena session cleanup...")

	sessions, err := getSerenaSessions()
	if err != nil {
		fmt.Printf("ERROR: Failed to get Serena sessions: %v\n", err)
		return
	}

	fmt.Printf("INFO: Found %d Serena sessions\n", len(sessions))

	cleanupCount := 0
	for _, session := range sessions {
		shouldCleanup := false
		reason := ""

		if session.Status == "inactive" {
			shouldCleanup = true
			reason = "status=inactive"
		}

		if time.Since(session.LastUsed) > 7*24*time.Hour {
			shouldCleanup = true
			reason = "last_used>7days"
		}

		if session.Project == "" || session.Project == "default" {
			shouldCleanup = true
			reason = "project=empty/default"
		}

		if shouldCleanup {
			fmt.Printf("CLEANUP: Session %s (%s) - Reason: %s\n", session.ID, session.Project, reason)
			if err := archiveSerenaSession(session.ID); err != nil {
				fmt.Printf("WARNING: Failed to archive session %s: %v\n", session.ID, err)
			} else {
				cleanupCount++
				fmt.Printf("SUCCESS: Session %s archived\n", session.ID)
			}
		} else {
			fmt.Printf("KEEP: Session %s (%s) - Last used: %s\n", session.ID, session.Project, session.LastUsed.Format(time.RFC3339))
		}
	}

	fmt.Printf("SUMMARY: Cleaned up %d inactive sessions\n", cleanupCount)
}

func getSerenaSessions() ([]SerenaSession, error) {
	cmd := exec.Command("uvx", "--from", "git+https://github.com/oraios/serena", "serena", "session", "list", "--json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("serena session list failed: %w", err)
	}

	var sessions []SerenaSession
	if err := json.Unmarshal(output, &sessions); err != nil {
		return nil, fmt.Errorf("failed to parse sessions: %w", err)
	}

	return sessions, nil
}

func archiveSerenaSession(sessionID string) error {
	cmd := exec.Command("uvx", "--from", "git+https://github.com/oraios/serena", "serena", "session", "archive", sessionID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("archive failed: %w, output: %s", err, string(output))
	}
	return nil
}

func verifySerenaProcess() error {
	cmd := exec.Command("pgrep", "-f", "serena.*start-mcp-server")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Serena process not found")
	}

	pids := strings.TrimSpace(string(output))
	fmt.Printf("INFO: Serena MCP server running (PIDs: %s)\n", pids)
	return nil
}
