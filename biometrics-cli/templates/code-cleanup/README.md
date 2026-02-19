# Code Cleanup Workflow Template

## Overview

The Code Cleanup workflow template provides automated code maintenance and organization. This template identifies and removes dead code, fixes formatting issues, organizes imports, and improves overall code structure without changing functionality.

The workflow is designed to improve codebase maintainability by addressing technical debt, removing unused code, and enforcing consistent coding standards. It can run on entire repositories or specific directories.

This template is essential for organizations seeking to:
- Reduce codebase complexity
- Remove unused code
- Improve maintainability
- Enforce code standards
- Reduce technical debt

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

## Input Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `repository` | string | Yes | - | Git repository URL |
| `target_path` | string | No | all | Specific path to clean |
| `cleanup_types` | array | No | all | Types: dead-code, formatting, imports, comments |
| `dry_run` | boolean | No | true | Preview changes without applying |

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
```

## Output Results

```json
{
  "cleanup": {
    "files_modified": 45,
    "dead_code_removed": 12,
    "imports_organized": 23,
    "formatting_fixed": 156
  },
  "changes": [
    {
      "file": "src/utils/helper.ts",
      "changes": ["removed unused function", "sorted imports"]
    }
  ]
}
```

## Workflow Steps

### Step 1: Analyze Codebase

Identifies cleanup opportunities.

### Step 2: Remove Dead Code

Eliminates unused code.

### Step 3: Fix Formatting

Applies consistent formatting.

### Step 4: Organize Imports

Sorts and consolidates imports.

### Step 5: Validate

Ensures changes don't break functionality.

## Usage

```bash
# Preview cleanup
biometrics workflow run code-cleanup \
  --repository https://github.com/org/app

# Apply cleanup
biometrics workflow run code-cleanup \
  --repository https://github.com/org/app \
  --dry_run false
```

## Troubleshooting

### Issue: Breaking Changes

Always run with dry_run first to preview changes.

### Issue: Large Repository

Process in smaller batches using target_path.

## Best Practices

### 1. Always Preview

Use dry_run before applying changes.

### 2. Review Changes

Have developers review cleanup changes.

### 3. Run Regularly

Schedule periodic cleanup to prevent debt accumulation.

## Related Templates

- **Code Review** (`code-review/`) - Quality validation
- **Refactor** (`refactor/`) - Structural improvements

---

**Template Version:** 1.0.0  
**Author:** BIOMETRICS Team  
**Category:** Code Quality  

*Last Updated: February 2026*
