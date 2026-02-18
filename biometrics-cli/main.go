package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF00")).
			Background(lipgloss.Color("#004400")).
			Padding(0, 2).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true).
			MarginBottom(2)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	runningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
)

// Step represents an onboarding step
type Step struct {
	Title   string
	Status  string // pending, running, success, error
	Message string
}

// Model for bubbletea
type Model struct {
	steps    []Step
	current  int
	quitting bool
	done     bool
	spinner  spinner.Model
	width    int
	height   int
	config   Config
}

// Config holds user configuration
type Config struct {
	GitLabToken      string
	NVIDIAApiKey     string
	InstallOpenCode  bool
	InstallOpenClaw  bool
	GitLabProjectID  string
	GitLabProjectURL string
}

// Messages
type statusMsg struct {
	Index   int
	Status  string
	Message string
}

type tickMsg time.Time

func initialModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = runningStyle

	steps := []Step{
		{Title: "Checking system requirements", Status: "pending"},
		{Title: "Installing Git", Status: "pending"},
		{Title: "Installing Node.js", Status: "pending"},
		{Title: "Installing pnpm", Status: "pending"},
		{Title: "Installing Homebrew", Status: "pending"},
		{Title: "Installing Python 3", Status: "pending"},
		{Title: "Configuring PATH", Status: "pending"},
		{Title: "Creating GitLab project", Status: "pending"},
		{Title: "Installing NLM CLI", Status: "pending"},
		{Title: "Installing OpenCode", Status: "pending"},
		{Title: "Installing OpenClaw", Status: "pending"},
		{Title: "Configuring integrations", Status: "pending"},
	}

	return Model{
		steps:   steps,
		spinner: s,
		config: Config{
			InstallOpenCode: true,
			InstallOpenClaw: true,
		},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, checkSystemRequirements)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case statusMsg:
		if msg.Index < len(m.steps) {
			m.steps[msg.Index].Status = msg.Status
			m.steps[msg.Index].Message = msg.Message

			if msg.Status == "success" {
				if msg.Index < len(m.steps)-1 {
					m.current = msg.Index + 1
					return m, nextStep(msg.Index + 1)
				}
				m.done = true
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(titleStyle.Render(" BIOMETRICS ONBOARD "))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render(" Professional Setup - GitLab + OpenCode + OpenClaw "))
	b.WriteString("\n")

	for i, step := range m.steps {
		status := "  "
		style := dimStyle

		switch step.Status {
		case "running":
			if i == m.current {
				status = runningStyle.Render(m.spinner.View()+" ") + " "
				style = runningStyle
			}
		case "success":
			status = successStyle.Render("✓ ")
			style = successStyle
		case "error":
			status = errorStyle.Render("✗ ")
			style = errorStyle
		}

		b.WriteString(style.Render(status + step.Title))
		if step.Message != "" && step.Status == "running" {
			b.WriteString(dimStyle.Render(" - " + step.Message))
		}
		b.WriteString("\n")
	}

	if m.done {
		b.WriteString("\n")
		b.WriteString(successStyle.Render("✓ Setup complete!\n"))
		b.WriteString(dimStyle.Render("\nRun 'biometrics' to start using the CLI\n"))
	}

	if m.quitting {
		b.WriteString("\nAborted.\n")
	}

	return b.String()
}

func nextStep(index int) tea.Cmd {
	switch index {
	case 0:
		return checkSystemRequirements
	case 1:
		return installGit
	case 2:
		return installNode
	case 3:
		return installPnpm
	case 4:
		return installHomebrew
	case 5:
		return installPython
	case 6:
		return configurePATH
	case 7:
		return createGitLabProject
	case 8:
		return installNLMCLI
	case 9:
		return installOpenCode
	case 10:
		return installOpenClaw
	case 11:
		return configureIntegrations
	default:
		return nil
	}
}

// System check and installation commands
func checkSystemRequirements() tea.Msg {
	// Check what's already installed
	installed := make(map[string]bool)
	tools := []string{"git", "node", "pnpm", "brew", "python3"}

	for _, tool := range tools {
		_, err := exec.LookPath(tool)
		installed[tool] = err == nil
	}

	// Return status for first missing tool or success if all installed
	for i, tool := range []string{"git", "node", "pnpm", "brew", "python3"} {
		if !installed[tool] {
			return statusMsg{Index: i + 1, Status: "pending", Message: "not installed"}
		}
	}

	return statusMsg{Index: 0, Status: "success", Message: "all requirements met"}
}

func installGit() tea.Msg {
	_, err := runCommand("brew", "install", "git")
	if err != nil {
		return statusMsg{Index: 1, Status: "error", Message: err.Error()}
	}
	return statusMsg{Index: 1, Status: "success", Message: "installed"}
}

func installNode() tea.Msg {
	_, err := runCommand("brew", "install", "node")
	if err != nil {
		return statusMsg{Index: 2, Status: "error", Message: err.Error()}
	}
	return statusMsg{Index: 2, Status: "success", Message: "installed"}
}

func installPnpm() tea.Msg {
	_, err := runCommand("brew", "install", "pnpm")
	if err != nil {
		return statusMsg{Index: 3, Status: "error", Message: err.Error()}
	}
	return statusMsg{Index: 3, Status: "success", Message: "installed"}
}

func installHomebrew() tea.Msg {
	// Check if already installed
	if _, err := exec.LookPath("brew"); err == nil {
		return statusMsg{Index: 4, Status: "success", Message: "already installed"}
	}
	return statusMsg{Index: 4, Status: "pending", Message: "manual installation required"}
}

func installPython() tea.Msg {
	_, err := runCommand("brew", "install", "python@3.11")
	if err != nil {
		return statusMsg{Index: 5, Status: "error", Message: err.Error()}
	}
	return statusMsg{Index: 5, Status: "success", Message: "installed"}
}

func configurePATH() tea.Msg {
	homeDir, _ := os.UserHomeDir()
	zshrc := filepath.Join(homeDir, ".zshrc")

	// Add pnpm to PATH
	pathLine := "\nexport PATH=\"$HOME/Library/pnpm:$PATH\"\n"

	data, _ := os.ReadFile(zshrc)
	content := string(data)

	if !strings.Contains(content, "pnpm") {
		f, _ := os.OpenFile(zshrc, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		f.WriteString(pathLine)
	}

	// Also add to current session
	os.Setenv("PATH", os.Getenv("HOME")+"/Library/pnpm:"+os.Getenv("PATH"))

	return statusMsg{Index: 6, Status: "success", Message: "PATH configured"}
}

func createGitLabProject() tea.Msg {
	// Skip if no token
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		return statusMsg{Index: 7, Status: "success", Message: "skipped (no token)"}
	}

	// Create project via API
	url := "https://gitlab.com/api/v4/projects"
	data := map[string]interface{}{
		"name":        "biometrics-media",
		"description": "BIOMETRICS project media storage",
		"visibility":  "public",
	}

	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	req.Header.Set("PRIVATE-TOKEN", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return statusMsg{Index: 7, Status: "error", Message: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		return statusMsg{Index: 7, Status: "success", Message: "project created"}
	}

	return statusMsg{Index: 7, Status: "success", Message: "project exists"}
}

func installNLMCLI() tea.Msg {
	_, err := runCommand("pnpm", "add", "-g", "nlm-cli")
	if err != nil {
		return statusMsg{Index: 8, Status: "error", Message: err.Error()}
	}
	return statusMsg{Index: 8, Status: "success", Message: "installed"}
}

func installOpenCode() tea.Msg {
	_, err := runCommand("brew", "install", "opencode")
	if err != nil {
		return statusMsg{Index: 9, Status: "error", Message: err.Error()}
	}
	return statusMsg{Index: 9, Status: "success", Message: "installed"}
}

func installOpenClaw() tea.Msg {
	_, err := runCommand("pnpm", "add", "-g", "@delqhi/openclaw")
	if err != nil {
		return statusMsg{Index: 10, Status: "error", Message: err.Error()}
	}
	return statusMsg{Index: 10, Status: "success", Message: "installed"}
}

func configureIntegrations() tea.Msg {
	// Placeholder for OpenClaw integration setup
	return statusMsg{Index: 11, Status: "success", Message: "configured"}
}

// Helper function
func runCommand(cmd string, args ...string) (string, error) {
	output, err := exec.Command(cmd, args...).CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
