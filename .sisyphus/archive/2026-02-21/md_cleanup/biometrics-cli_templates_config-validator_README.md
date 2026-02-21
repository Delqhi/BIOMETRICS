# Config Validator Workflow Template

## Overview

The Config Validator workflow template provides automated configuration validation. This template verifies configuration files for correctness, best practices, security issues, and compatibility across your entire application stack.

The workflow supports various configuration formats including YAML, JSON, TOML, and environment files. It can validate application configs, infrastructure configs, deployment configs, and supports custom validation schemas. The template is essential for preventing config-related bugs before they reach production.

This template is essential for organizations seeking to:
- Prevent config-related bugs
- Enforce configuration standards
- Catch misconfigurations early
- Ensure security compliance
- Validate deployment readiness

## Purpose

The primary purpose of the Config Validator template is to:

1. **Validate Syntax** - Ensure config files are valid
2. **Check Best Practices** - Verify recommended patterns
3. **Security Scan** - Identify security issues
4. **Schema Validation** - Validate against schemas
5. **Compatibility Check** - Ensure config compatibility

### Key Use Cases

- **Pre-deployment Validation** - Validate configs before deployment
- **CI/CD Integration** - Automated validation in pipelines
- **Config Review** - Manual config reviews
- **Security Audits** - Configuration security checks
- **Environment Parity** - Ensure dev/staging/prod parity

## Input Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `config_files` | array | Yes | - | Config files to validate |
| `schema` | string | No | - | JSON schema for validation |
| `validation_type` | array | No | all | Types: syntax, schema, security, best-practices |
| `environment` | string | No | - | Target environment |

### Input Examples

```yaml
# Example 1: Validate all configs
inputs:
  config_files:
    - config/app.yaml
    - config/database.yaml

# Example 2: Security-focused validation
inputs:
  config_files:
    - config/production.yaml
  validation_type:
    - security
    - syntax

# Example 3: With schema
inputs:
  config_files:
    - config/api.yaml
  schema: schemas/api-schema.json

# Example 4: Environment-specific
inputs:
  config_files:
    - config/app.yaml
    - config/db.yaml
  environment: production
  validation_type:
    - syntax
    - schema
    - security
    - best-practices
```

## Output Results

The template produces comprehensive validation reports:

| Output | Type | Description |
|--------|------|-------------|
| `overall_status` | string | Validation pass/fail |
| `results` | array | Results per config file |
| `issues` | array | All identified issues |

### Output Report Structure

```json
{
  "validation": {
    "timestamp": "2026-02-19T10:30:00Z",
    "files_validated": 3,
    "environment": "production",
    "duration_seconds": 12
  },
  "overall_status": "failed",
  "summary": {
    "passed": 2,
    "failed": 1,
    "warnings": 5,
    "errors": 3
  },
  "results": [
    {
      "file": "config/app.yaml",
      "status": "passed",
      "issues": []
    },
    {
      "file": "config/database.yaml",
      "status": "failed",
      "issues": [
        {
          "severity": "error",
          "type": "missing_field",
          "message": "Missing required field: host",
          "line": 5,
          "path": "database.host"
        },
        {
          "severity": "warning",
          "type": "insecure_value",
          "message": "Debug mode enabled in production",
          "line": 12,
          "path": "app.debug"
        }
      ]
    },
    {
      "file": "config/cache.yaml",
      "status": "passed",
      "issues": [
        {
          "severity": "info",
          "type": "optimization",
          "message": "Consider enabling compression",
          "path": "cache.compression"
        }
      ]
    }
  ],
  "security_scan": {
    "secrets_detected": false,
    "insecure_protocols": ["http"],
    "weak_encryption": false
  },
  "best_practices": {
    "score": 85,
    "issues": [
      "Missing health check endpoint"
    ]
  }
}
```

## Workflow Steps

### Step 1: Parse Configs

**ID:** `parse-configs`  
**Type:** agent  
**Timeout:** 5 minutes  
**Provider:** opencode-zen

Parses configuration files:
- Detects file format
- Handles parsing errors
- Normalizes structure

### Step 2: Validate Syntax

**ID:** `validate-syntax`  
**Type:** agent  
**Timeout:** 5 minutes  
**Provider:** opencode-zen

Checks for syntax errors:
- Valid JSON/YAML/TOML
- Proper formatting
- Type correctness

### Step 3: Validate Schema

**ID:** `validate-schema`  
**Type:** agent  
**Timeout:** 5 minutes  
**Provider:** opencode-zen

Validates against schemas if provided:
- Required fields
- Field types
- Value constraints
- Custom rules

### Step 4: Security Scan

**ID:** `security-scan`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Checks for security issues:
- Hardcoded secrets
- Insecure protocols
- Weak authentication
- Missing encryption

### Step 5: Best Practices Check

**ID:** `best-practices`  
**Type:** agent  
**Timeout:** 5 minutes  
**Provider:** opencode-zen

Validates best practices:
- Configuration patterns
- Recommended values
- Performance settings
- Monitoring setup

### Step 6: Generate Report

**ID:** `generate-report`  
**Type:** agent  
**Timeout:** 3 minutes  
**Provider:** opencode-zen

Creates validation report.

## Usage Examples

### CLI Usage

```bash
# Basic validation
biometrics workflow run config-validator \
  --config_files '["config/app.yaml", "config/db.yaml"]'

# Security-focused
biometrics workflow run config-validator \
  --config_files '["config/production.yaml"]' \
  --validation_type '["security"]'

# With schema
biometrics workflow run config-validator \
  --config_files '["config/api.yaml"]' \
  --schema schemas/api-schema.json

# Environment-specific
biometrics workflow run config-validator \
  --config_files '["config/app.yaml"]' \
  --environment production
```

### Programmatic Usage

```go
engine := workflows.NewWorkflowEngine("./templates")
template, _ := engine.LoadTemplate("config-validator")

instance, _ := engine.CreateInstance(template, map[string]interface{}{
    "config_files": []string{
        "config/app.yaml",
        "config/database.yaml",
    },
    "validation_type": []string{
        "syntax",
        "schema",
        "security",
        "best-practices",
    },
    "environment": "production",
})

result, err := engine.Execute(ctx, instance)
```

## Configuration

### Custom Rules

```yaml
options:
  rules:
    - id: NO-DEBUG-PROD
      severity: error
      check: "debug_mode must be false in production"
      condition:
        environment: production
        path: app.debug
        value: true
        
    - id: REQUIRED-HEALTH
      severity: warning
      check: "Health check endpoint must be configured"
      condition:
        path: app.healthCheck
        required: true
        
    - id: MAX_CONNECTIONS
      severity: warning
      check: "Database connections should be limited"
      condition:
        path: database.maxConnections
        max: 100
```

### Schema Example

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["host", "port"],
  "properties": {
    "host": {
      "type": "string",
      "pattern": "^[a-z0-9.-]+$"
    },
    "port": {
      "type": "integer",
      "minimum": 1,
      "maximum": 65535
    },
    "ssl": {
      "type": "boolean",
      "default": true
    }
  }
}
```

### Security Rules

```yaml
options:
  security:
    check_secrets: true
    check_protocols: true
    check_auth: true
    allowed_protocols:
      - https
      - ssh
      - sftp
    secret_patterns:
      - "password\\s*=\\s*[\"']"
      - "api[_-]?key\\s*=\\s*[\"']"
      - "secret\\s*=\\s*[\"']"
```

## Troubleshooting

### Common Issues

#### Issue: Schema Validation Fails

**Symptom:** Schema errors

**Solution:** Verify schema is correct and matches config structure.

#### Issue: False Positives

**Symptom:** Too many warnings

**Solution:** Adjust rule severity or add exclusions:
```yaml
options:
  exclude:
    - path: app.debug  # Allow in development
      when: environment == development
```

#### Issue: Large Config Files

**Symptom:** Timeout

**Solution:** Increase timeout:
```yaml
options:
  timeout: 60
```

### Debug Mode

```yaml
options:
  debug: true
  verbose: true
  include_context: true
```

## Best Practices

### 1. Validate in CI/CD

Run validation in your pipeline:
```bash
# In CI pipeline
biometrics workflow run config-validator \
  --config_files '["config/production.yaml"]' \
  --environment production
```

### 2. Use Schemas

Define schemas for type safety and auto-completion.

### 3. Check Security

Always include security validation, especially for production.

### 4. Environment Parity

Validate that configurations work across environments.

### 5. Version Control

Keep configs in version control for audit trail.

## Related Templates

- **Code Review** (`code-review/`) - Code quality validation
- **Deployment** (`deployment/`) - Deployment with config validation
- **Security Audit** (`security-audit/`) - Comprehensive security review

---

**Template Version:** 1.0.0  
**Author:** BIOMETRICS Team  
**Category:** Configuration  
**Tags:** configuration, validation, security, schema, best-practices

*Last Updated: February 2026*
