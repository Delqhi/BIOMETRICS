# Code Cleanup Workflow Template

## Overview

The Code Cleanup workflow template provides automated code maintenance and organization. This template identifies and removes dead code, fixes formatting issues, organizes imports, and improves overall code structure without changing functionality.

The workflow is designed to improve codebase maintainability by addressing technical debt, removing unused code, and enforcing consistent coding standards. It can run on entire repositories or specific directories, making it flexible for both comprehensive cleanups and targeted fixes.

This template is essential for organizations seeking to:
- Reduce codebase complexity
- Remove unused code
- Improve maintainability
- Enforce code standards
- Reduce technical debt
- Speed up build times

## Purpose

The primary purpose of the Code Cleanup template is to:

1. **Remove Dead Code** - Identify and remove unused functions, variables, and imports
2. **Fix Formatting** - Apply consistent code formatting
3. **Organize Imports** - Sort and consolidate import statements
4. **Improve Structure** - Enhance code organization
5. **Preserve Behavior** - Ensure no functional changes

### Key Use Cases

- **Pre-release Cleanup** - Clean up code before releases
- **Regular Maintenance** - Scheduled code hygiene
- **Onboarding** - Clean up legacy codebases
- **Refactoring Prep** - Prepare code for refactoring
- **Performance Improvement** - Remove code bloat

## Input Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `repository` | string | Yes | - | Git repository URL |
| `target_path` | string | No | all | Specific path to clean |
| `cleanup_types` | array | No | all | Types: dead-code, formatting, imports, comments |
| `dry_run` | boolean | No | true | Preview changes without applying |
| `language` | string | No | auto | Programming language |
| `exclude_patterns` | array | No | [] | Patterns to exclude |

### Input Examples

```yaml
# Example 1: Full cleanup
inputs:
  repository: https://github.com/org/app
  dry_run: false

# Example 2: Preview cleanup
inputs:
  repository: https://github.com/org/app
  target_path: src/utils
  dry_run: true

# Example 3: Specific cleanup types
inputs:
  repository: https://github.com/org/app
  cleanup_types:
    - dead-code
    - imports
  dry_run: false

# Example 4: Language-specific
inputs:
  repository: https://github.com/org/python-app
  language: python
  cleanup_types:
    - formatting
    - imports
  exclude_patterns:
    - "**/migrations/**"
    - "**/generated/**"
```

## Output Results

The template produces detailed cleanup reports:

| Output | Type | Description |
|--------|------|-------------|
| `files_modified` | number | Number of files changed |
| `changes_summary` | object | Summary of all changes |
| `issues_fixed` | number | Total issues resolved |

### Output Report Structure

```json
{
  "cleanup": {
    "timestamp": "2026-02-19T10:30:00Z",
    "repository": "https://github.com/org/app",
    "target_path": "src",
    "dry_run": false,
    "duration_seconds": 145
  },
  "statistics": {
    "files_analyzed": 156,
    "files_modified": 45,
    "dead_code_removed": 12,
    "imports_organized": 23,
    "formatting_fixed": 156,
    "comments_cleaned": 34
  },
  "changes": [
    {
      "file": "src/utils/helper.ts",
      "changes": [
        "Removed unused function 'formatDate'",
        "Sorted imports alphabetically",
        "Fixed indentation"
      ]
    },
    {
      "file": "src/services/user.service.ts",
      "changes": [
        "Removed unused variable 'tempData'",
        "Organized imports"
      ]
    }
  ],
  "warnings": [
    {
      "file": "src/legacy/old-module.ts",
      "severity": "medium",
      "message": "File appears unused but may have side effects"
    }
  ]
}
```

## Workflow Steps

### Step 1: Analyze Codebase

**ID:** `analyze-codebase`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Identifies cleanup opportunities:
- Parses all source files
- Builds dependency graph
- Identifies dead code
- Detects formatting issues

### Step 2: Remove Dead Code

**ID:** `remove-dead-code`  
**Type:** agent  
**Timeout:** 15 minutes  
**Provider:** opencode-zen

Eliminates unused code:
- Removes unused functions
- Removes unused variables
- Removes unreachable code
- Removes duplicate code

### Step 3: Fix Formatting

**ID:** `fix-formatting`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Applies consistent formatting:
- Fixes indentation
- Fixes line length issues
- Fixes whitespace problems
- Applies code style

### Step 4: Organize Imports

**ID:** `organize-imports`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Sorts and consolidates imports:
- Removes unused imports
- Groups imports by type
- Sorts alphabetically
- Adds missing imports

### Step 5: Validate Changes

**ID:** `validate-changes`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Ensures changes don't break functionality:
- Runs tests
- Verifies syntax
- Checks types
- Validates build

## Usage Examples

### CLI Usage

```bash
# Preview cleanup (default dry_run: true)
biometrics workflow run code-cleanup \
  --repository https://github.com/org/app

# Apply cleanup
biometrics workflow run code-cleanup \
  --repository https://github.com/org/app \
  --dry_run false

# Target specific path
biometrics workflow run code-cleanup \
  --repository https://github.com/org/app \
  --target_path src/utils

# Specific cleanup types
biometrics workflow run code-cleanup \
  --repository https://github.com/org/app \
  --cleanup_types '["dead-code", "imports"]' \
  --dry_run false
```

### Programmatic Usage

```go
import "github.com/biometrics/biometrics-cli/pkg/workflows"

engine := workflows.NewWorkflowEngine("./templates")
template, _ := engine.LoadTemplate("code-cleanup")

instance, _ := engine.CreateInstance(template, map[string]interface{}{
    "repository":   "https://github.com/org/app",
    "target_path":   "src/utils",
    "cleanup_types": []string{"dead-code", "imports", "formatting"},
    "dry_run":       false,
})

result, err := engine.Execute(context.Background(), instance)
```

## Configuration

### Cleanup Rules

```yaml
options:
  rules:
    dead_code:
      remove_unused_functions: true
      remove_unused_variables: true
      remove_unused_imports: true
      
    formatting:
      indent_size: 2
      max_line_length: 100
      trailing_comma: true
      
    imports:
      group_external_first: true
      sort_alphabetically: true
      remove_duplicates: true
```

### Exclusions

```yaml
options:
  exclude:
    paths:
      - "**/node_modules/**"
      - "**/dist/**"
      - "**/generated/**"
      - "**/migrations/**"
    files:
      - "*.test.ts"
      - "*.spec.ts"
```

## Troubleshooting

### Common Issues

#### Issue: Breaking Changes

**Symptom:** Tests fail after cleanup

**Solution:**
- Always run with dry_run first
- Review changes before applying
- Run tests after cleanup
- Use version control to revert

#### Issue: Large Repository

**Symptom:** Timeout or memory issues

**Solution:**
- Process in smaller batches using target_path
- Exclude large directories
- Run multiple cleanup passes

#### Issue: False Positives

**Symptom:** Removes code that's actually used

**Solution:**
- Review generated changes carefully
- Add exclusions for dynamic usage
- Use function annotations to mark usage

### Debug Mode

```yaml
options:
  debug: true
  verbose: true
```

## Best Practices

### 1. Always Preview First

Use dry_run=true to preview changes before applying:
```bash
biometrics workflow run code-cleanup --repository ...  # dry_run default: true
```

### 2. Review Changes

Have developers review cleanup changes before merging.

### 3. Run Tests

Always run tests after cleanup to ensure nothing broke.

### 4. Use Version Control

Commit changes in small batches for easy rollback.

### 5. Run Regularly

Schedule periodic cleanup to prevent debt accumulation:
```yaml
trigger:
  type: schedule
  cron: "0 2 * * 0"  # Weekly
```

### 6. Exclude Generated Code

Never clean generated or migration code:
```yaml
options:
  exclude:
    paths:
      - "**/migrations/**"
      - "**/generated/**"
```

## Related Templates

- **Code Review** (`code-review/`) - Quality validation
- **Refactor** (`refactor/`) - Structural improvements
- **Test Generator** (`test-generator/`) - Add tests after cleanup

---

**Template Version:** 1.0.0  
**Author:** BIOMETRICS Team  
**Category:** Code Quality  
**Tags:** cleanup, code-quality, maintenance, formatting, dead-code

*Last Updated: February 2026*
