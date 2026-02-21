package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	statStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Bold(true)
)

type metrics struct {
	cycles      int
	modelAcq    int
	goroutines  int
	chaosEvents int
	uptime      string
	lastUpdate  time.Time
	loading     bool
	error       string
}

type tickMsg time.Time

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

type model struct {
	metrics metrics
	spinner spinner.Model
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))

	return model{
		metrics: metrics{loading: true},
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		spinner.Tick,
		tickCmd(),
		fetchMetricsCmd(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tickMsg:
		return m, tea.Batch(tickCmd(), fetchMetricsCmd())

	case metricsMsg:
		m.metrics = msg.metrics
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("🚀 BIOMETRICS ORCHESTRATOR"))
	b.WriteString("\n\n")

	if m.metrics.error != "" {
		b.WriteString(errorStyle.Render("ERROR: " + m.metrics.error))
		b.WriteString("\n")
	} else if m.metrics.loading {
		b.WriteString(m.spinner.View())
		b.WriteString(" Loading metrics...")
	} else {
		// Status
		status := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575")).
			Render("● RUNNING")
		b.WriteString(fmt.Sprintf("Status: %s (PID: %d)\n\n", status, os.Getpid()))

		// Metrics grid
		b.WriteString("📊 LIVE METRICS:\n")
		b.WriteString("┌──────────────────────────────────────┐\n")
		b.WriteString(fmt.Sprintf("│ %-30s %6d │\n", "Cycles Completed:", m.metrics.cycles))
		b.WriteString(fmt.Sprintf("│ %-30s %6d │\n", "Model Acquisitions:", m.metrics.modelAcq))
		b.WriteString(fmt.Sprintf("│ %-30s %6d │\n", "Active Goroutines:", m.metrics.goroutines))
		b.WriteString(fmt.Sprintf("│ %-30s %6d │\n", "Chaos Events:", m.metrics.chaosEvents))
		b.WriteString(fmt.Sprintf("│ %-30s %11s │\n", "Last Update:", m.metrics.lastUpdate.Format("15:04:05")))
		b.WriteString("└──────────────────────────────────────┘\n\n")

		// Endpoints
		b.WriteString("🔗 ENDPOINTS:\n")
		b.WriteString("  • Metrics: http://localhost:59002/metrics\n")
		b.WriteString("  • Logs: /tmp/biometrics-orchestrator.log\n\n")

		// Instructions
		b.WriteString("Press 'q' to quit\n")
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(strings.Repeat("─", 50)))

	return b.String()
}

type metricsMsg struct{ metrics metrics }

func tickCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchMetricsCmd() tea.Cmd {
	return func() tea.Msg {
		m := metrics{loading: true, lastUpdate: time.Now()}

		resp, err := http.Get("http://localhost:59002/metrics")
		if err != nil {
			m.error = fmt.Sprintf("Cannot connect: %v", err)
			m.loading = false
			return metricsMsg{m}
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			m.error = fmt.Sprintf("Read error: %v", err)
			m.loading = false
			return metricsMsg{m}
		}

		lines := strings.Split(string(body), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "biometrics_orchestrator_cycles_total ") {
				parts := strings.Fields(line)
				if len(parts) > 1 {
					m.cycles, _ = strconv.Atoi(parts[1])
				}
			} else if strings.HasPrefix(line, "biometrics_orchestrator_model_acquisitions_total") {
				if strings.Contains(line, "qwen3.5") {
					parts := strings.Fields(line)
					if len(parts) > 1 {
						m.modelAcq, _ = strconv.Atoi(parts[1])
					}
				}
			} else if strings.HasPrefix(line, "go_goroutines ") {
				parts := strings.Fields(line)
				if len(parts) > 1 {
					m.goroutines, _ = strconv.Atoi(parts[1])
				}
			} else if strings.HasPrefix(line, "biometrics_orchestrator_chaos_events_total") {
				parts := strings.Split(line, "}")
				if len(parts) > 1 {
					val := strings.TrimSpace(parts[1])
					if v, err := strconv.Atoi(val); err == nil {
						m.chaosEvents += v
					}
				}
			}
		}

		m.loading = false
		return metricsMsg{m}
	}
}
