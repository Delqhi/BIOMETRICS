# Compliance Check Workflow Template

## Overview

The Compliance Check workflow template provides automated regulatory compliance validation for codebases and infrastructure. This template enables organizations to ensure their software systems adhere to industry standards and legal requirements such as GDPR, SOC2, HIPAA, and PCI-DSS.

This workflow is designed to run as a scheduled job (weekly by default) or on-demand, providing continuous compliance monitoring without manual intervention. It leverages AI agents to analyze code patterns, configuration files, and infrastructure setup against predefined compliance frameworks.

## Purpose

The primary purpose of the Compliance Check template is to:

1. **Automate Compliance Validation** - Eliminate manual compliance reviews by automating the entire process
2. **Continuous Monitoring** - Run scheduled checks to catch compliance drift early
3. **Multi-Framework Support** - Support multiple compliance standards simultaneously (GDPR, SOC2, HIPAA, etc.)
4. **Actionable Reporting** - Generate detailed reports with specific remediation steps
5. **Audit Trail** - Maintain comprehensive logs of all compliance checks for audit purposes

### Key Use Cases

- **GDPR Compliance** - Verify data handling, privacy controls, and user consent mechanisms
- **SOC2 Compliance** - Check security controls, access management, and logging
- **HIPAA Compliance** - Validate protected health information (PHI) safeguards
- **PCI-DSS Compliance** - Ensure payment card data protection
- **Custom Standards** - Apply organization-specific compliance rules

## Input Parameters

The Compliance Check template accepts the following input parameters:

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `standards` | array | Yes | ["gdpr", "SOC2", "HIPAA"] | List of compliance standards to check |
| `repository` | string | No | - | Git repository URL for code analysis |
| `scan_scope` | string | No | "full" | Scope of scan: "full", "incremental", or "critical" |
| `auto_remediate` | boolean | No | false | Attempt automatic fixes where possible |
| `report_format` | string | No | "json" | Output format: "json", "html", or "pdf" |
| `notification_channels` | array | No | [] | List of channels to notify on completion |

### Input Examples

```yaml
# Example 1: Basic GDPR check
inputs:
  standards:
    - gdpr
  scan_scope: full
  auto_remediate: false
  report_format: json

# Example 2: Multi-framework compliance
inputs:
  standards:
    - gdpr
    - SOC2
    - HIPAA
  repository: https://github.com/org/enterprise-app
  scan_scope: incremental
  auto_remediate: true
  report_format: html
  notification_channels:
    - slack:#compliance
    - email:security@company.com

# Example 3: Critical path only
inputs:
  standards:
    - PCI-DSS
  scan_scope: critical
  report_format: pdf
```

## Output Results

The template produces comprehensive compliance reports with the following outputs:

| Output | Type | Description |
|--------|------|-------------|
| `report` | object | Detailed compliance report with findings |
| `passed` | boolean | Overall compliance status (true if all checks pass) |
| `score` | number | Compliance score (0-100) |
| `issues` | array | List of identified compliance issues |
| `remediations` | array | Suggested remediation steps |
| `evidence` | object | Evidence collected during the scan |

### Output Report Structure

```json
{
  "report": {
    "standards_checked": ["gdpr", "SOC2", "HIPAA"],
    "timestamp": "2026-02-19T10:30:00Z",
    "duration_seconds": 245,
    "findings": {
      "critical": 2,
      "high": 5,
      "medium": 12,
      "low": 8
    }
  },
  "passed": false,
  "score": 72,
  "issues": [
    {
      "id": "GDPR-001",
      "severity": "critical",
      "standard": "gdpr",
      "title": "Missing data encryption at rest",
      "description": "Customer PII stored without encryption",
      "location": "src/database/user_model.go:45",
      "remediation": "Enable encryption for database tables containing PII"
    }
  ],
  "remediations": [
    {
      "issue_id": "GDPR-001",
      "steps": [
        "1. Enable TDE on PostgreSQL database",
        "2. Update data migration scripts",
        "3. Verify encryption is active"
      ],
      "estimated_effort": "2 hours"
    }
  ]
}
```

## Workflow Steps

The Compliance Check workflow executes in sequential phases:

### Step 1: Initialize Compliance Scan

**ID:** `init-scan`  
**Type:** agent  
**Timeout:** 5 minutes

Initializes the compliance scan environment, loads framework definitions, and prepares the analysis context.

### Step 2: Scan Code Compliance

**ID:** `scan-code`  
**Type:** agent  
**Timeout:** 30 minutes  
**Provider:** opencode-zen

Analyzes source code against compliance requirements:
- Data handling patterns
- Authentication/authorization logic
- Logging and auditing
- Cryptographic implementations
- Third-party dependencies

### Step 3: Scan Infrastructure

**ID:** `scan-infra`  
**Type:** agent  
**Timeout:** 20 minutes  
**Provider:** opencode-zen

Evaluates infrastructure configuration:
- Cloud resource settings
- Network security groups
- Access control policies
- Encryption configurations
- Backup and recovery procedures

### Step 4: Generate Compliance Report

**ID:** `generate-report`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Compiles findings into comprehensive report with:
- Executive summary
- Detailed findings by standard
- Risk assessments
- Remediation recommendations
- Evidence artifacts

### Step 5: Notify Stakeholders

**ID:** `notify`  
**Type:** agent  
**Timeout:** 5 minutes  
**Provider:** opencode-zen

Distributes reports to configured channels:
- Email notifications
- Slack/Teams messages
- Webhook payloads
- Dashboard updates

## Usage Examples

### CLI Usage

```bash
# Run compliance check for GDPR only
biometrics workflow run compliance-check \
  --standards '["gdpr"]'

# Full multi-framework compliance scan
biometrics workflow run compliance-check \
  --standards '["gdpr", "SOC2", "HIPAA"]' \
  --repository https://github.com/org/app \
  --scan_scope full \
  --report_format html

# Incremental scan with auto-remediation
biometrics workflow run compliance-check \
  --standards '["PCI-DSS"]' \
  --scan_scope incremental \
  --auto_remediate true
```

### Programmatic Usage

```go
import "github.com/biometrics/biometrics-cli/pkg/workflows"

engine := workflows.NewWorkflowEngine("./templates")

template, _ := engine.LoadTemplate("compliance-check")

instance, _ := engine.CreateInstance(template, map[string]interface{}{
    "standards":       []string{"gdpr", "SOC2"},
    "repository":      "https://github.com/org/app",
    "scan_scope":      "full",
    "auto_remediate":  false,
    "report_format":   "json",
})

ctx := context.Background()
result, err := engine.Execute(ctx, instance)
```

### API Usage

```bash
curl -X POST http://localhost:8080/api/v1/workflows/run \
  -H "Content-Type: application/json" \
  -d '{
    "template": "compliance-check",
    "inputs": {
      "standards": ["gdpr", "SOC2", "HIPAA"],
      "repository": "https://github.com/org/app",
      "scan_scope": "full"
    }
  }'
```

## Configuration

### Framework Configuration

Customize compliance standards by creating framework definition files:

```yaml
# config/frameworks/my-custom-standard.yaml
name: my-custom-standard
version: 1.0.0
description: Custom internal security standard

requirements:
  - id: CUSTOM-001
    title: No hardcoded secrets
    severity: critical
    check_type: pattern_match
    patterns:
      - "password\\s*=\\s*[\"']"
      - "api[_-]?key\\s*=\\s*[\"']"
      - "secret\\s*=\\s*[\"']"

  - id: CUSTOM-002
    title: Secure random generation
    severity: high
    check_type: imports
    required_imports:
      - "crypto/rand"
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `COMPLIANCE_FRAMEWORKS_PATH` | Path to custom framework definitions | `./config/frameworks` |
| `COMPLIANCE_REPORT_PATH` | Output directory for reports | `./outputs/compliance` |
| `COMPLIANCE_PROVIDER` | AI provider for analysis | `opencode-zen` |
| `COMPLIANCE_MODEL` | AI model to use | `opencode/minimax-m2.5-free` |

### Scheduling Configuration

The template supports scheduled execution via cron:

```yaml
trigger:
  type: schedule
  cron: "0 0 * * 0"  # Weekly on Sunday at midnight
```

Common schedule patterns:
- `"0 0 * * 0"` - Weekly (Sunday midnight)
- `"0 0 * * 1-5"` - Weekdays at midnight
- `"0 */6 * * *"` - Every 6 hours
- `"0 0 1 * *"` - Monthly on 1st

## Troubleshooting

### Common Issues

#### Issue: Scan Timeout

**Symptom:** Workflow times out during large repository scans

**Solution:**
```bash
# Increase timeout in template
options:
  timeout: 60m  # Increase from default 30m

# Or limit scan scope
inputs:
  scan_scope: incremental  # Only scan changed files
```

#### Issue: Missing Framework

**Symptom:** "Framework not found" error

**Solution:**
```bash
# Verify framework exists
ls config/frameworks/

# Create missing framework
cat > config/frameworks/NEW-STANDARD.yaml << 'EOF'
name: new-standard
version: 1.0.0
requirements: []
EOF
```

#### Issue: Permission Denied

**Symptom:** Cannot access repository or create outputs

**Solution:**
```bash
# Check repository access
git ls-remote <repository-url>

# Verify output directory permissions
chmod -R 755 outputs/
```

#### Issue: High False Positives

**Symptom:** Too many non-critical findings

**Solution:**
Adjust severity thresholds in configuration:
```yaml
inputs:
  severity_threshold: high  # Only report high/critical
```

### Debug Mode

Enable detailed logging:

```yaml
options:
  debug: true
  log_level: trace
```

### Getting Help

For additional support:
- Documentation: `/docs/compliance/`
- Issues: Report via GitHub Issues
- Slack: `#compliance-help` channel

## Best Practices

### 1. Run Regularly

Schedule weekly compliance scans to catch issues early:
```yaml
trigger:
  type: schedule
  cron: "0 0 * * 0"
```

### 2. Prioritize Critical Issues

Focus remediation efforts on critical and high severity findings first.

### 3. Integrate into CI/CD

Add compliance checks to your pipeline:
```bash
# In CI pipeline
biometrics workflow run compliance-check --standards '["SOC2"]'
```

### 4. Maintain Evidence

Retain compliance reports for audit purposes (minimum 2 years recommended).

### 5. Customize Frameworks

Adapt standards to your organization's specific requirements.

## Related Templates

- **Security Audit** (`security-audit/`) - More detailed security scanning
- **Config Validator** (`config-validator/`) - Configuration validation
- **Code Review** (`code-review/`) - General code quality checks

---

**Template Version:** 1.0.0  
**Author:** BIOMETRICS Team  
**Category:** Compliance & Security  
**Tags:** compliance, security, GDPR, SOC2, HIPAA, audit

*Last Updated: February 2026*
