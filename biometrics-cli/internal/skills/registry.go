package skills

import (
	"fmt"
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
