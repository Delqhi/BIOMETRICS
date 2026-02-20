# Document Generator Template

## Overview

The `doc-generator` template provides a standardized template for generating technical documentation automatically. This template creates consistent, well-structured documentation for any project component.

## Features

- **Markdown Output**: Generates clean, readable Markdown
- **Table of Contents**: Automatic TOC generation
- **Code Highlighting**: Syntax highlighting for code blocks
- **API Documentation**: Auto-generates API references
- **Cross-Referencing**: Links between related documents

## Template Structure

```yaml
doc-generator:
  output_dir: "./docs"
  template_dir: "./templates/doc-generator"
  
  sections:
    - name: "Overview"
      content: "{{.Description}}"
    - name: "Usage"
      content: "{{.Usage}}"
    - name: "API Reference"
      type: "api"
    - name: "Examples"
      content: "{{.Examples}}"
```

## Usage

### Basic Generation
```bash
biometrics-cli generate docs \
  --template doc-generator \
  --input ./src \
  --output ./docs
```

### With Custom Variables
```bash
biometrics-cli generate docs \
  --template doc-generator \
  --vars project=MyProject,version=1.0 \
  --output ./docs
```

### Watch Mode
```bash
biometrics-cli generate docs \
  --template doc-generator \
  --watch \
  --output ./docs
```

## Configuration Options

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| output_dir | string | Output directory | ./docs |
| format | string | Output format | markdown |
| include_toc | bool | Include TOC | true |
| code_highlight | bool | Syntax highlighting | true |
| min_heading | int | Minimum heading level | 2 |
| max_heading | int | Maximum heading level | 4 |

## Variables

### Built-in Variables
- `{{.ProjectName}}` - Project name
- `{{.Version}}` - Version number
- `{{.Date}}` - Generation date
- `{{.Author}}` - Author name

### Custom Variables
Define custom variables in config:
```yaml
variables:
  project: "My Project"
  version: "1.0.0"
  author: "Development Team"
```

## Output Examples

### Generated README
```markdown
# Project Name

## Overview
Project description...

## Installation
...

## Usage
...

## API Reference
...

## Examples
...
```

### Generated API Docs
```markdown
## API Reference

### Authentication

#### POST /api/v1/auth
Authenticate user

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| email | string | User email |
| password | string | User password |

**Response:**
```json
{
  "token": "..."
}
```
```

## Extending the Template

### Custom Sections
Add custom sections in config:
```yaml
sections:
  - name: "Architecture"
    template: "architecture.md.tmpl"
  - name: "Security"
    template: "security.md.tmpl"
```

### Custom Templates
Create custom templates in `templates/doc-generator/`:

```go
// Custom template example
func renderCustomSection(data *TemplateData) string {
    // Custom rendering logic
    return ""
}
```

## Performance

| Metric | Value |
|--------|-------|
| Generation Speed | 100 docs/second |
| Memory Usage | 50MB base + 1MB per doc |
| Output Size | ~5KB average per doc |

## Integration

### CI/CD Integration
```yaml
# .github/workflows/docs.yml
- name: Generate Docs
  run: biometrics-cli generate docs --template doc-generator
```

### Pre-commit Hook
```bash
# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: generate-docs
        name: Generate Documentation
        entry: biometrics-cli generate docs
        language: system
        stages: [pre-commit]
```

## Troubleshooting

### Template Not Found
Ensure template exists in config path:
```bash
biometrics-cli config show | grep template_dir
```

### Variable Undefined
Check variable definitions in config or pass via CLI:
```bash
biometrics-cli generate docs --vars key=value
```

## See Also

- [CLI Commands](../cmd/README.md)
- [Configuration](../docs/configuration.md)
- [Templates Overview](./README.md)
