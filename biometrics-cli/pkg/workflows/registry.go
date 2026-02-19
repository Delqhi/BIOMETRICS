package workflows

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// TemplateRegistry manages workflow templates
type TemplateRegistry struct {
	templates map[string]*Workflow
	basePath  string
}

// NewTemplateRegistry creates a new template registry
func NewTemplateRegistry(basePath string) *TemplateRegistry {
	return &TemplateRegistry{
		templates: make(map[string]*Workflow),
		basePath:  basePath,
	}
}

// LoadTemplates loads all templates from the base path
func (r *TemplateRegistry) LoadTemplates() error {
	entries, err := os.ReadDir(r.basePath)
	if err != nil {
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		templatePath := filepath.Join(r.basePath, entry.Name(), "template.yaml")
		data, err := os.ReadFile(templatePath)
		if err != nil {
			continue
		}

		workflow, err := LoadWorkflow(data)
		if err != nil {
			continue
		}

		r.templates[workflow.Name] = workflow
	}

	return nil
}

// GetTemplate retrieves a template by name
func (r *TemplateRegistry) GetTemplate(name string) (*Workflow, bool) {
	wf, ok := r.templates[name]
	return wf, ok
}

// ListTemplates returns all template names
func (r *TemplateRegistry) ListTemplates() []string {
	names := make([]string, 0, len(r.templates))
	for name := range r.templates {
		names = append(names, name)
	}
	return names
}

// AddTemplate adds a new template
func (r *TemplateRegistry) AddTemplate(workflow *Workflow) {
	r.templates[workflow.Name] = workflow
}

// LoadWorkflow parses YAML into Workflow
func LoadWorkflow(data []byte) (*Workflow, error) {
	var workflow Workflow
	err := yaml.Unmarshal(data, &workflow)
	if err != nil {
		return nil, fmt.Errorf("failed to parse workflow: %w", err)
	}
	return &workflow, nil
}
