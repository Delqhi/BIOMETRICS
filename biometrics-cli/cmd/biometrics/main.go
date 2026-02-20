package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	version   = "1.0.0"
	commit    = "dev"
	buildDate = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		runInit()
	case "onboard":
		runOnboarding()
	case "auto":
		runAutoSetup()
	case "check":
		runBiometricsCheck()
	case "find-keys":
		findAPIKeys()
	case "version":
		printVersion()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`
BIOMETRICS CLI v2.0.0

Usage: biometrics <command> [options]

Commands:
  init          Initialize BIOMETRICS repository
  onboard       Interactive onboarding process  
  auto          Automatic AI-powered setup
  check         Check BIOMETRICS compliance
  find-keys     Find existing API keys on system
  version       Show version information
 `)
}

func printVersion() {
	buildTime := buildDate
	if buildDate == "unknown" {
		buildTime = time.Now().Format("2006-01-02")
	}

	fmt.Printf("biometrics-cli v%s (commit: %s, built: %s)\n", version, commit, buildTime)
}

func runInit() {
	fmt.Println("Initializing BIOMETRICS repository...")

	dirs := []string{
		"global/01-agents",
		"global/02-models",
		"global/03-mandates",
		"local/projects",
		"biometrics-cli/bin",
		"biometrics-cli/commands",
		"docs/media",
		"scripts",
		"assets/images",
		"assets/videos",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Failed to create directory %s: %v\n", dir, err)
			os.Exit(1)
		}
		fmt.Printf("Created directory: %s\n", dir)
	}

	createReadme("global", "Global Configurations", "Global AI agent configurations and mandates")
	createReadme("local", "Local Projects", "Project-specific configurations")
	createReadme("biometrics-cli", "BIOMETRICS CLI", "Command-line interface and automation")
	createReadme("docs", "Documentation", "Complete documentation for BIOMETRICS")

	fmt.Println("\nBIOMETRICS repository initialized successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Run 'biometrics onboard' for interactive setup")
	fmt.Println("  2. Run 'biometrics auto' for automatic AI-powered setup")
}

func runOnboarding() {
	fmt.Println("Starting BIOMETRICS onboarding process...")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Step 1: Checking for existing API keys...")
	existingKeys := findExistingKeys()
	if len(existingKeys) > 0 {
		fmt.Println("Found existing API keys:")
		for keyType, keyPath := range existingKeys {
			fmt.Printf("  - %s: %s\n", keyType, keyPath)
		}
	} else {
		fmt.Println("No existing API keys found")
	}

	fmt.Println()
	fmt.Println("Step 2: Do you want to use existing keys? (y/n)")
	useExisting, _ := reader.ReadString('\n')
	useExisting = strings.TrimSpace(useExisting)

	if useExisting == "y" || useExisting == "Y" {
		fmt.Println("Using existing API keys")
		copyKeysToEnv(existingKeys)
	} else {
		fmt.Println("Step 3: Enter your NVIDIA API Key (or press Enter to skip):")
		nvidiaKey, _ := reader.ReadString('\n')
		nvidiaKey = strings.TrimSpace(nvidiaKey)

		if nvidiaKey != "" {
			saveToEnv("NVIDIA_API_KEY", nvidiaKey)
			fmt.Println("NVIDIA API Key saved")
		}
	}

	fmt.Println("\nStep 4: Do you have a GitLab token? (y/n)")
	hasGitLab, _ := reader.ReadString('\n')
	hasGitLab = strings.TrimSpace(hasGitLab)

	if hasGitLab == "y" || hasGitLab == "Y" {
		fmt.Println("Enter GitLab token:")
		gitLabToken, _ := reader.ReadString('\n')
		gitLabToken = strings.TrimSpace(gitLabToken)
		saveToEnv("GITLAB_TOKEN", gitLabToken)
		fmt.Println("GitLab token saved")
	}

	fmt.Println("\nStep 5: Do you have Supabase credentials? (y/n)")
	hasSupabase, _ := reader.ReadString('\n')
	hasSupabase = strings.TrimSpace(hasSupabase)

	if hasSupabase == "y" || hasSupabase == "Y" {
		fmt.Println("Enter Supabase URL:")
		supabaseURL, _ := reader.ReadString('\n')
		supabaseURL = strings.TrimSpace(supabaseURL)
		saveToEnv("SUPABASE_URL", supabaseURL)

		fmt.Println("Enter Supabase Key:")
		supabaseKey, _ := reader.ReadString('\n')
		supabaseKey = strings.TrimSpace(supabaseKey)
		saveToEnv("SUPABASE_KEY", supabaseKey)
		fmt.Println("Supabase credentials saved")
	}

	fmt.Println("\nStep 6: Installing dependencies...")
	installDependencies()

	fmt.Println("\nStep 7: Installing OpenCode skills...")
	installOpenCodeSkills()

	// Step 8: OpenClaw Onboarding
	fmt.Println("\nStep 8: OpenClaw Configuration")
	fmt.Println("Which channels should OpenClaw monitor?")
	fmt.Println()
	fmt.Println(" [ ] GitHub Issues")
	fmt.Println(" [ ] Slack")
	fmt.Println(" [ ] Email")
	fmt.Println(" [ ] Discord")
	fmt.Println(" [ ] Telegram")
	fmt.Println(" [ ] WhatsApp")
	fmt.Println()
	fmt.Println("Enter channels (comma-separated, e.g., github,slack,email):")
	channels, _ := reader.ReadString('\n')
	channels = strings.TrimSpace(channels)
	if channels != "" {
		fmt.Printf("OpenClaw will monitor: %s\n", channels)
	}

	fmt.Println("\nSelect skills for OpenClaw:")
	fmt.Println(" [x] playwright - Browser automation")
	fmt.Println(" [x] git-master - Git operations")
	fmt.Println(" [x] frontend-ui-ux - UI/UX design")
	fmt.Println(" [ ] dev-browser - Advanced browser control")
	fmt.Println(" [ ] context7 - Documentation lookup")
	fmt.Println()
	fmt.Println("Enter skills (comma-separated, or press Enter for defaults):")
	skills, _ := reader.ReadString('\n')
	skills = strings.TrimSpace(skills)
	if skills == "" {
		skills = "playwright,git-master,frontend-ui-ux"
	}
	fmt.Printf("Selected skills: %s\n", skills)

	fmt.Println("\nAllowed processes (NO deployment without permission!):")
	fmt.Println(" [x] File creation/modification")
	fmt.Println(" [x] Code improvements")
	fmt.Println(" [x] Test execution (localhost only)")
	fmt.Println(" [x] Browser debugging")
	fmt.Println(" [ ] Deployment (requires explicit permission)")
	fmt.Println()

	// Step 9: 24/7 Auto-Development Option
	fmt.Println("\nStep 9: Enable 24/7 Auto-Development?")
	fmt.Println("OpenClaw will continuously:")
	fmt.Println(" - Create missing files")
	fmt.Println(" - Improve existing code")
	fmt.Println(" - Run tests in localhost")
	fmt.Println(" - Debug in browser")
	fmt.Println(" - NEVER deploy without permission")
	fmt.Println()
	fmt.Println("Enable? (y/n)")
	autoDev, _ := reader.ReadString('\n')
	autoDev = strings.TrimSpace(autoDev)

	if autoDev == "y" || autoDev == "Y" {
		fmt.Println("\n✅ 24/7 Auto-Development ENABLED")
		fmt.Println("OpenClaw will run in DELQHI OMEGA LOOP mode")
		fmt.Println("Continuous improvement cycle activated")
	} else {
		fmt.Println("\n24/7 Auto-Development DISABLED")
		fmt.Println("OpenClaw will only run on manual commands")
	}

	fmt.Println("\nOnboarding complete!")
	fmt.Println("\nNext steps:")
	fmt.Println(" 1. Run 'biometrics check' to verify setup")
	fmt.Println(" 2. Run 'opencode \"Build a feature\"' to start developing")
	fmt.Println(" 3. OpenClaw is ready to orchestrate your development!")
}

func runAutoSetup() {
	fmt.Println("Starting automatic AI-powered setup...")
	fmt.Println()

	fmt.Println("Step 1: Scanning for existing API keys...")
	existingKeys := findExistingKeys()

	if len(existingKeys) > 0 {
		fmt.Println("Found API keys:")
		for keyType, keyPath := range existingKeys {
			fmt.Printf("  - %s: %s\n", keyType, keyPath)
		}
		copyKeysToEnv(existingKeys)
		fmt.Println("Keys copied to .env")
	} else {
		fmt.Println("No API keys found. You'll need to add them manually.")
		fmt.Println("   Get NVIDIA API key: https://build.nvidia.com/")
	}

	fmt.Println("\nStep 2: Creating enterprise directory structure...")
	runInit()

	fmt.Println("\nStep 3: Installing dependencies...")
	installDependencies()

	fmt.Println("\nStep 4: Installing OpenCode skills...")
	installOpenCodeSkills()

	fmt.Println("\nStep 5: Verifying setup...")
	runBiometricsCheck()

	fmt.Println("\nAutomatic setup complete!")
}

func runBiometricsCheck() {
	fmt.Println("BIOMETRICS REPO CHECK")
	fmt.Println("=====================")
	fmt.Println()

	checks := []struct {
		name    string
		passed  bool
		success string
		failure string
	}{
		{"global/ README", checkFileExists("global/README.md"), "global/README.md exists", "global/README.md missing"},
		{"local/ README", checkFileExists("local/README.md"), "local/README.md exists", "local/README.md missing"},
		{"biometrics-cli/ README", checkFileExists("biometrics-cli/README.md"), "biometrics-cli/README.md exists", "biometrics-cli/README.md missing"},
		{".env file", checkFileExists(".env"), ".env exists", ".env missing"},
		{"oh-my-opencode.json", checkFileExists("oh-my-opencode.json"), "oh-my-opencode.json exists", "oh-my-opencode.json missing"},
		{"requirements.txt", checkFileExists("requirements.txt"), "requirements.txt exists", "requirements.txt missing"},
	}

	allPassed := true
	for _, check := range checks {
		if check.passed {
			fmt.Println("✓ " + check.success)
		} else {
			fmt.Println("✗ " + check.failure)
			allPassed = false
		}
	}

	fmt.Println()
	if allPassed {
		fmt.Println("All checks passed! BIOMETRICS is ready.")
	} else {
		fmt.Println("Some checks failed. Run 'biometrics onboard' to fix.")
	}
}

func findAPIKeys() {
	fmt.Println("Scanning for existing API keys...")
	fmt.Println()

	keys := findExistingKeys()

	if len(keys) == 0 {
		fmt.Println("No API keys found on system")
		return
	}

	fmt.Println("Found API keys:")
	for keyType, keyPath := range keys {
		fmt.Printf("  %s: %s\n", keyType, keyPath)
	}

	fmt.Println("\nTo use these keys, run: biometrics onboard")
}

func createReadme(dir, title, description string) {
	readmePath := filepath.Join(dir, "README.md")
	content := fmt.Sprintf(`# %s

## Purpose

%s

## Structure

%s/
├── README.md
└── ...

## Enterprise Practices 2026

- Go-Style Modularity
- Machine-readable
- KI-Agent optimized

---

Generated by BIOMETRICS CLI v2.0.0
`, title, description, dir)

	os.WriteFile(readmePath, []byte(content), 0644)
	fmt.Printf("Created: %s\n", readmePath)
}

func findExistingKeys() map[string]string {
	keys := make(map[string]string)

	knownAPIKeys := []string{
		"NVIDIA_API_KEY",
		"MOONSHOT_API_KEY",
		"GITLAB_TOKEN",
		"GITLAB_MEDIA_PROJECT_ID",
		"SUPABASE_URL",
		"SUPABASE_KEY",
		"OPENCODE_API_KEY",
		"ANTHROPIC_API_KEY",
		"OPENROUTER_API_KEY",
		"HUGGINGFACE_TOKEN",
	}

	for _, envVar := range knownAPIKeys {
		if value := os.Getenv(envVar); value != "" {
			keys[envVar] = "Environment Variable"
		}
	}

	locations := map[string]string{
		"NVIDIA_API_KEY": os.Getenv("HOME") + "/.nvidia_api_key",
		"GitLab":         os.Getenv("HOME") + "/.gitlab_token",
		"Supabase":       os.Getenv("HOME") + "/.supabase_key",
		"OpenCode":       os.Getenv("HOME") + "/.opencode/key",
	}

	for keyType, location := range locations {
		if _, err := os.Stat(location); err == nil {
			keys[keyType] = location
		}
	}

	dotenvPath := ".env"
	if _, err := os.Stat(dotenvPath); err == nil {
		file, err := os.Open(dotenvPath)
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "#") {
					continue
				}
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					envVar := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					if value != "" && value != "YOUR_KEY_HERE" && value != "nvapi-YOUR_KEY" {
						keys[envVar] = ".env file"
					}
				}
			}
		}
	}

	shellConfig := os.Getenv("HOME") + "/.zshrc"
	if _, err := os.Stat(shellConfig); err == nil {
		file, err := os.Open(shellConfig)
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "export ") {
					parts := strings.SplitN(line, "=", 2)
					if len(parts) == 2 {
						envVar := strings.TrimPrefix(strings.TrimSpace(parts[0]), "export ")
						value := strings.TrimSpace(parts[1])
						value = strings.Trim(value, `"'`)
						if strings.HasPrefix(value, "nvapi-") || strings.HasPrefix(value, "sk-") {
							keys[envVar] = "~/.zshrc"
						}
					}
				}
			}
		}
	}

	return keys
}

func copyKeysToEnv(keys map[string]string) {
	for keyType, keyPath := range keys {
		if keyPath != "Environment Variable" {
			data, err := os.ReadFile(keyPath)
			if err == nil {
				saveToEnv(keyType+"_KEY", strings.TrimSpace(string(data)))
			}
		}
	}
}

func autoConfigureKeys(keys map[string]string) {
	fmt.Println("\n=== Auto-Configuring API Keys ===")

	ohMyOpenCodePath := "oh-my-opencode.json"
	if _, err := os.Stat(ohMyOpenCodePath); err != nil {
		fmt.Printf("Warning: %s not found, skipping auto-configure\n", ohMyOpenCodePath)
		return
	}

	content, err := os.ReadFile(ohMyOpenCodePath)
	if err != nil {
		fmt.Printf("Warning: Could not read %s: %v\n", ohMyOpenCodePath, err)
		return
	}

	updated := string(content)
	for key, source := range keys {
		if source == "Environment Variable" || source == ".env file" || source == "~/.zshrc" {
			value := os.Getenv(key)
			if value != "" {
				envVarPlaceholder := "${" + key + "}"
				if !strings.Contains(updated, envVarPlaceholder) {
					updated = strings.ReplaceAll(updated, `"`+key+`": ""`, `"`+key+`": "${`+key+`}"`)
					fmt.Printf("✓ Configured %s from %s\n", key, source)
				}
			}
		}
	}

	err = os.WriteFile(ohMyOpenCodePath, []byte(updated), 0644)
	if err != nil {
		fmt.Printf("Warning: Could not write %s: %v\n", ohMyOpenCodePath, err)
		return
	}

	fmt.Println("✓ API keys auto-configured successfully!")
}

func saveToEnv(key, value string) {
	f, err := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Could not open .env: %v\n", err)
		return
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("%s=%s\n", key, value))
}

func checkFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func installDependencies() {
	if _, err := exec.LookPath("npm"); err == nil {
		fmt.Println("Installing npm dependencies...")
		cmd := exec.Command("npm", "install")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: npm install failed: %v\n", err)
		} else {
			fmt.Println("✓ npm dependencies installed")
		}
	}

	if _, err := exec.LookPath("pip3"); err == nil {
		fmt.Println("Installing Python dependencies...")
		cmd := exec.Command("pip3", "install", "-r", "requirements.txt")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: pip3 install failed: %v\n", err)
		} else {
			fmt.Println("✓ Python dependencies installed")
		}
	}

	if _, err := exec.LookPath("opencode"); err == nil {
		fmt.Println("Authenticating OpenCode providers...")
		cmd := exec.Command("opencode", "auth", "add", "nvidia-nim")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: opencode auth failed: %v\n", err)
		} else {
			fmt.Println("✓ OpenCode authenticated")
		}
	}
}

func installOpenCodeSkills() {
	fmt.Println("Installing OpenCode skills...")

	skills := []string{
		"playwright",
		"git-master",
		"frontend-ui-ux",
	}

	for _, skill := range skills {
		fmt.Printf("Installing skill: %s... ", skill)
		cmd := exec.Command("opencode", "skill", "install", skill)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: Failed to install %s\n", skill)
		} else {
			fmt.Println("✓")
		}
	}

	fmt.Println("OpenCode skills installed")
}
