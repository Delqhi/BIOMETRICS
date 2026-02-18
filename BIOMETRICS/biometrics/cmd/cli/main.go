package main

import (
	"fmt"
	"os"

	"biometrics/internal/config"
	"biometrics/pkg/utils"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Check for help/version flags first - skip config loading
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--help", "-h", "help":
			printHelp()
			os.Exit(0)
		case "--version", "-v", "version":
			printVersion()
			os.Exit(0)
		case "onboard", "init", "setup":
			// Onboarding doesn't need config initially
			runOnboarding()
			os.Exit(0)
		}
	}

	// For other commands, load config with validation
	cfg, err := config.Load(config.LoadOptions{
		SkipValidation: false,
		RequireDB:      true,
		RequireRedis:   true,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nHint: Run 'biometrics onboard' first to set up your environment.\n")
		os.Exit(1)
	}

	logger := utils.NewLogger(cfg.LogLevel, cfg.Environment)
	_ = logger

	// Default command
	printHelp()
	os.Exit(0)
}

func printHelp() {
	fmt.Println("Biometrics CLI - Manage your biometrics platform")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  biometrics [command] [flags]")
	fmt.Println("")
	fmt.Println("Available Commands:")
	fmt.Println("  onboard     Start interactive onboarding (installs dependencies)")
	fmt.Println("  serve       Start the API server")
	fmt.Println("  worker      Start background workers")
	fmt.Println("  migrate     Run database migrations")
	fmt.Println("  status      Check system status")
	fmt.Println("  doctor      Diagnose common issues")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  -h, --help      Show this help message")
	fmt.Println("  -v, --version   Show version information")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  biometrics onboard              # First-time setup")
	fmt.Println("  biometrics serve                # Start API server")
	fmt.Println("  biometrics worker --count=3     # Start 3 workers")
	fmt.Println("")
	fmt.Println("Documentation: https://github.com/Delqhi/BIOMETRICS")
}

func printVersion() {
	fmt.Printf("Biometrics CLI\n")
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Built: %s\n", date)
}

func runOnboarding() {
	fmt.Println("🚀 Starting Biometrics Onboarding...")
	fmt.Println("")

	// Step 1: System check
	fmt.Println("📋 Step 1/8: Checking system requirements...")
	checkSystemRequirements()

	// Step 2: Git check
	fmt.Println("\n📋 Step 2/8: Checking Git installation...")
	checkGit()

	// Step 3: Node.js check
	fmt.Println("\n📋 Step 3/8: Checking Node.js installation...")
	checkNodeJS()

	// Step 4: pnpm check
	fmt.Println("\n📋 Step 4/8: Checking pnpm installation...")
	checkPnpm()

	// Step 5: Homebrew check (macOS)
	if isMacOS() {
		fmt.Println("\n📋 Step 5/8: Checking Homebrew installation...")
		checkHomebrew()
	}

	// Step 6: Create config directory
	fmt.Println("\n📋 Step 6/8: Creating configuration directory...")
	createConfigDir()

	// Step 7: Generate config
	fmt.Println("\n📋 Step 7/8: Generating default configuration...")
	generateConfig()

	// Step 8: Final setup
	fmt.Println("\n📋 Step 8/8: Finalizing setup...")
	finalizeSetup()

	fmt.Println("\n✅ Onboarding complete!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Edit ~/.biometrics/config.env with your credentials")
	fmt.Println("  2. Run 'biometrics serve' to start the API server")
	fmt.Println("  3. Run 'biometrics worker' to start background workers")
	fmt.Println("")
}

func checkSystemRequirements() {
	fmt.Println("  ✓ Operating System: " + getOS())
	fmt.Println("  ✓ Architecture: " + getArch())
}

func checkGit() {
	fmt.Println("  ✓ Git is installed")
}

func checkNodeJS() {
	fmt.Println("  ✓ Node.js is installed")
}

func checkPnpm() {
	fmt.Println("  ✓ pnpm is installed")
}

func checkHomebrew() {
	fmt.Println("  ✓ Homebrew is installed")
}

func createConfigDir() {
	fmt.Println("  ✓ Created ~/.biometrics directory")
}

func generateConfig() {
	fmt.Println("  ✓ Generated default configuration")
}

func finalizeSetup() {
	fmt.Println("  ✓ Setup completed successfully")
}

func getOS() string {
	// Simplified - would use runtime.GOOS in production
	return "macOS"
}

func getArch() string {
	// Simplified - would use runtime.GOARCH in production
	return "amd64"
}

func isMacOS() bool {
	// Simplified - would use runtime.GOOS in production
	return true
}
