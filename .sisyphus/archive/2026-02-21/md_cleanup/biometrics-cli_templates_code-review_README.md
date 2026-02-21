# Code Review Workflow Template

## Overview

The Code Review workflow template provides automated code quality analysis with comprehensive security, performance, and best practices validation. This template is designed to run on pull requests, pushes, or manually to ensure code meets quality standards before merging or deployment.

The workflow performs multi-dimensional analysis including security vulnerability scanning, performance profiling, code style validation, and architectural assessment. It generates detailed reports with actionable recommendations and overall quality scores.

This template is critical for organizations seeking to:
- Maintain code quality standards
- Catch security vulnerabilities early
- Ensure consistent coding practices
- Prevent performance regressions
- Improve overall codebase health

## Purpose

The primary purpose of the Code Review template is to:

1. **Automate Reviews** - Replace manual code reviews with AI-powered analysis
2. **Catch Issues Early** - Identify problems before they reach production
3. **Ensure Standards** - Validate adherence to coding standards
4. **Provide Feedback** - Give developers actionable improvement suggestions
5. **Track Quality** - Monitor code quality trends over time

### Key Use Cases

- **Pull Request Reviews** - Automated review on PR creation
- **Pre-commit Checks** - Validate code before commits
- **Scheduled Reviews** - Periodic comprehensive reviews
- **Security Audits** - Focused security analysis
- **Performance Reviews** - Performance-focused assessments

## Input Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `repository` | string | Yes | - | Git repository URL |
| `branch` | string | Yes | - | Branch to review |
| `files` | array | No | all | Specific files to review |
| `focus_areas` | array | No | all | Areas: security, performance, style, best-practices |

### Input Examples

```yaml
# Example 1: Full code review
inputs:
  repository: https://github.com/org/app
  branch: feature/new-feature
  focus_areas:
    - security
    - performance
    - style
    - best-practices

# Example 2: Security-focused review
inputs:
  repository: https://github.com/org/app
  branch: main
  focus_areas:
    - security

# Example 3: Specific files
inputs:
  repository: https://github.com/org/app
  branch: fix/auth-bug
  files:
    - src/auth/login.ts
    - src/auth/token.ts
```

## Output Results

| Output | Type | Description |
|--------|------|-------------|
| `report` | object | Detailed review report |
| `score` | number | Overall quality score (0-100) |
| `issues` | array | List of identified issues |

### Output Report Structure

```json
{
  "report": {
    "timestamp": "2026-02-19T10:30:00Z",
    "repository": "https://github.com/org/app",
    "branch": "feature/new-feature",
    "files_reviewed": 12,
    "lines_analyzed": 3450
  },
  "score": 85,
  "issues": [
    {
      "id": "SEC-001",
      "severity": "critical",
      "category": "security",
      "title": "SQL Injection Vulnerability",
      "location": "src/db/query.ts:45",
      "description": "User input directly concatenated into SQL query",
      "recommendation": "Use parameterized queries"
    }
  ],
  "summary": {
    "critical": 1,
    "high": 5,
    "medium": 12,
    "low": 8
  }
}
```

## Workflow Steps

### Step 1: Fetch Code Changes

**ID:** `fetch-code`  
**Type:** agent  
**Timeout:** 5 minutes

Retrieves code from repository.

### Step 2: Analyze Structure

**ID:** `analyze-structure`  
**Type:** agent  
**Timeout:** 10 minutes

Evaluates code architecture and organization.

### Step 3: Security Scan

**ID:** `security-scan`  
**Type:** agent  
**Timeout:** 10 minutes

Checks for security vulnerabilities:
- SQL injection
- XSS vulnerabilities
- Authentication issues
- Sensitive data exposure

### Step 4: Performance Analysis

**ID:** `performance-analysis`  
**Type:** agent  
**Timeout:** 10 minutes

Identifies performance issues:
- N+1 queries
- Inefficient algorithms
- Missing caching

### Step 5: Style Check

**ID:** `style-check`  
**Type:** agent  
**Timeout:** 5 minutes

Validates code style and conventions.

### Step 6: Generate Report

**ID:** `generate-report`  
**Type:** agent  
**Timeout:** 5 minutes

Compiles comprehensive review report.

### Step 7: Notify

**ID:** `notify`  
**Type:** agent  
**Timeout:** 2 minutes

Sends notifications with results.

## Usage Examples

### CLI Usage

```bash
# Full review on PR
biometrics workflow run code-review \
  --repository https://github.com/org/app \
  --branch feature/new-feature

# Quick security check
biometrics workflow run code-review \
  --repository https://github.com/org/app \
  --branch main \
  --focus_areas '["security"]'
```

### Programmatic Usage

```go
engine := workflows.NewWorkflowEngine("./templates")
template, _ := engine.LoadTemplate("code-review")

instance, _ := engine.CreateInstance(template, map[string]interface{}{
    "repository": "https://github.com/org/app",
    "branch":    "feature/new-feature",
    "focus_areas": []string{"security", "performance"},
})

result, _ := engine.Execute(ctx, instance)
```

## Configuration

### Custom Rules

```yaml
options:
  rules:
    - id: CUSTOM-001
      category: security
      check: "no_console_log"
      severity: warning
```

### Threshold Configuration

```yaml
options:
  thresholds:
    min_score: 70
    critical_allowed: 0
```

## Troubleshooting

### Issue: Timeout on Large PRs

**Solution:**
```yaml
options:
  timeout: 60m
```

### Issue: False Positives

Adjust severity in configuration or exclude false positives with comments.

## Best Practices

### 1. Run on Every PR

Automate reviews on all pull requests.

### 2. Set Quality Gates

Require minimum scores for merging.

### 3. Track Trends

Monitor quality scores over time.

### 4. Prioritize Security

Always include security scanning.

## Related Templates

- **Security Audit** (`security-audit/`) - More detailed security review
- **Test Generator** (`test-generator/`) - Generate tests for issues found
- **Bug Fix** (`bug-fix/`) - Fix issues identified

---

**Template Version:** 1.0.0  
**Author:** BIOMETRICS Team  
**Category:** Code Quality  
**Tags:** code-review, security, quality, linting, best-practices

*Last Updated: February 2026*
