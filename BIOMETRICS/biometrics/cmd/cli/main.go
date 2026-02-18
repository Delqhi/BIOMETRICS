package main

import (
	"fmt"
	"os"

	"biometrics/internal/config"
	"biometrics/pkg/utils"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger := utils.NewLogger(cfg.LogLevel, cfg.Environment)
	_ = logger

	fmt.Println("Biometrics CLI - Manage your biometrics platform")
	fmt.Println("")
	fmt.Println("Available commands:")
	fmt.Println("  serve    - Start the API server")
	fmt.Println("  migrate  - Run database migrations")
	fmt.Println("  worker   - Start background workers")
	fmt.Println("")
	fmt.Println("Use biometrics --help for more information")

	os.Exit(0)
}
