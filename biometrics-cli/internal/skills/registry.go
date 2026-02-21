package skills

import (
	"strings"
)

type Skill struct {
	Name        string
	Keywords    []string
	Description string
}

var Registry = []Skill{
	{
		Name:        "playwright",
		Keywords:    []string{"browser", "ui", "scrape", "screenshot", "test"},
		Description: "Browser automation via Playwright MCP",
	},
	{
		Name:        "git-master",
		Keywords:    []string{"git", "commit", "push", "pull", "merge", "branch"},
		Description: "Advanced git operations and workflow management",
	},
	{
		Name:        "frontend-ui-ux",
		Keywords:    []string{"css", "react", "nextjs", "tailwind", "design"},
		Description: "UI/UX engineering and styling skills",
	},
	{
		Name:        "dev-browser",
		Keywords:    []string{"navigate", "click", "fill", "scroll", "browser", "automation"},
		Description: "Browser automation with persistent page state",
	},
	{
		Name:        "database",
		Keywords:    []string{"sql", "postgres", "mysql", "query", "database", "migration"},
		Description: "Database operations and schema management",
	},
	{
		Name:        "api-integration",
		Keywords:    []string{"api", "rest", "http", "endpoint", "webhook", "json"},
		Description: "REST API integration and webhook handling",
	},
	{
		Name:        "docker",
		Keywords:    []string{"docker", "container", "image", "compose", "kubernetes"},
		Description: "Docker container management and deployment",
	},
	{
		Name:        "security",
		Keywords:    []string{"auth", "token", "password", "encrypt", "oauth", "jwt"},
		Description: "Security and authentication patterns",
	},
	{
		Name:        "testing",
		Keywords:    []string{"test", "mock", "assert", "coverage", "unit", "integration"},
		Description: "Testing strategies and test automation",
	},
	{
		Name:        "performance",
		Keywords:    []string{"optimize", "performance", "cache", "memory", "latency"},
		Description: "Performance optimization and profiling",
	},
	{
		Name:        "ai-ml",
		Keywords:    []string{"model", "train", "predict", "embedding", "vector", "llm"},
		Description: "AI/ML model integration and training",
	},
	{
		Name:        "monitoring",
		Keywords:    []string{"metrics", "logs", "prometheus", "grafana", "alert"},
		Description: "Monitoring, logging, and observability",
	},
}

func MatchSkills(prompt string) []string {
	matched := make([]string, 0)
	lowerPrompt := strings.ToLower(prompt)

	for _, skill := range Registry {
		for _, kw := range skill.Keywords {
			if strings.Contains(lowerPrompt, kw) {
				matched = append(matched, skill.Name)
				break
			}
		}
	}
	return matched
}
