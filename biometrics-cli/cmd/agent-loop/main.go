package main

import (
	"encoding/json"
	"fmt"
	"os"
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
	fmt.Println("🚀 Starting BIOMETRICS Agent Loop Orchestrator with Strict Model Assignment...")
	
	boulderPath := "/Users/jeremy/.sisyphus/boulder.json"
	tracker := NewModelTracker()
	
	for {
		fmt.Printf("[%s] Checking for agent status...\n", time.Now().Format(time.RFC3339))
		
		boulder, err := readBoulder(boulderPath)
		if err != nil {
			fmt.Printf("Error reading boulder.json: %v\n", err)
			time.Sleep(10 * time.Second)
			continue
		}
		
		if boulder.ActivePlan == "" {
			fmt.Println("No active plan found. Waiting...")
			time.Sleep(30 * time.Second)
			continue
		}
		
		// Simulate model assignment based on agent
		model := getModelForAgent(boulder.Agent)
		if err := tracker.Acquire(model); err != nil {
			fmt.Printf("⚠️ Model Collision: %v. Waiting for model to be free...\n", err)
		} else {
			fmt.Printf("✅ Model %s acquired for agent %s\n", model, boulder.Agent)
			// In a real scenario, we would run the agent here
			// tracker.Release(model) // Release when done
		}
		
		time.Sleep(60 * time.Second)
	}
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
