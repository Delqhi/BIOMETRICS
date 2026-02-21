# Test Generator Workflow Template

## Overview

The Test Generator workflow template provides automated test creation for software projects. This template leverages AI agents to analyze codebases and generate comprehensive unit tests, integration tests, and end-to-end tests that achieve high code coverage while focusing on edge cases and critical paths.

The workflow is designed to understand code structure, identify testable units, and produce meaningful tests that validate expected behavior while catching potential regressions. It supports multiple testing frameworks and can adapt to various programming languages and project structures.

This template is essential for organizations seeking to:
- Rapidly increase test coverage
- Ensure code quality through testing
- Prevent regressions
- Accelerate development velocity
- Meet coverage requirements

## Purpose

The primary purpose of the Test Generator template is to:

1. **Automate Test Creation** - Generate tests without manual effort
2. **Maximize Coverage** - Achieve high code coverage efficiently
3. **Identify Edge Cases** - Find and test boundary conditions
4. **Follow Best Practices** - Produce maintainable, readable tests
5. **Integrate with CI/CD** - Seamlessly fit into existing pipelines

### Key Use Cases

- **New Feature Testing** - Generate tests for new features automatically
- **Legacy Code Coverage** - Add tests to untested code
- **Refactoring Safety** - Ensure refactored code maintains behavior
- **Bug Reproduction** - Create tests that reproduce reported bugs
- **API Testing** - Generate integration tests for APIs

## Input Parameters

The Test Generator template accepts the following input parameters:

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `repository` | string | Yes | - | Git repository URL |
| `files` | array | No | all | Specific files to generate tests for |
| `test_type` | string | No | unit | Type of tests (unit, integration, e2e) |
| `framework` | string | No | auto | Testing framework (auto-detect from project) |
| `coverage_target` | number | No | 80 | Target code coverage percentage |
| `focus_areas` | array | No | [] | Priority areas for testing |

### Input Examples

```yaml
# Example 1: Generate unit tests for entire project
inputs:
  repository: https://github.com/org/app
  test_type: unit
  coverage_target: 80

# Example 2: Generate tests for specific files
inputs:
  repository: https://github.com/org/app
  files:
    - src/utils/parser.ts
    - src/services/user.service.ts
  test_type: unit

# Example 3: Integration tests for API
inputs:
  repository: https://github.com/org/app
  test_type: integration
  focus_areas:
    - API endpoints
    - Database operations

# Example 4: E2E tests for critical flows
inputs:
  repository: https://github.com/org/app
  test_type: e2e
  focus_areas:
    - Checkout flow
    - Authentication
```

## Output Results

The template produces comprehensive test outputs:

| Output | Type | Description |
|--------|------|-------------|
| `test_files` | array | Generated test file paths |
| `test_count` | number | Total number of tests generated |
| `coverage` | object | Code coverage statistics |
| `issues` | array | Potential issues found in tests |

### Output Report Structure

```json
{
  "generated": {
    "timestamp": "2026-02-19T10:30:00Z",
    "test_type": "unit",
    "framework": "jest",
    "repository": "https://github.com/org/app"
  },
  "test_files": [
    {
      "path": "src/utils/parser.test.ts",
      "tests": 24,
      "coverage_percent": 85
    },
    {
      "path": "src/services/user.service.test.ts",
      "tests": 42,
      "coverage_percent": 78
    }
  ],
  "statistics": {
    "total_test_files": 8,
    "total_tests": 156,
    "total_coverage_percent": 82,
    "lines_covered": 2450,
    "lines_total": 2990
  },
  "edge_cases": [
    {
      "file": "src/utils/parser.ts",
      "test": "handles empty string input",
      "type": "boundary"
    },
    {
      "file": "src/utils/parser.ts", 
      "test": "handles unicode characters",
      "type": "special_chars"
    }
  ],
  "issues": []
}
```

## Workflow Steps

### Step 1: Analyze Code Structure

**ID:** `analyze-structure`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Analyzes code to understand:
- File structure
- Exportable functions
- Dependencies
- Testable units

### Step 2: Identify Test Cases

**ID:** `identify-cases`  
**Type:** agent  
**Timeout:** 15 minutes  
**Provider:** opencode-zen

Identifies test cases:
- Happy path scenarios
- Edge cases
- Error conditions
- Boundary values

### Step 3: Generate Tests

**ID:** `generate-tests`  
**Type:** agent  
**Timeout:** 30 minutes  
**Provider:** opencode-zen

Generates test code:
- Unit tests for functions
- Integration tests for modules
- E2E tests for flows

### Step 4: Validate Tests

**ID:** `validate-tests`  
**Type:** agent  
**Timeout:** 10 minutes  
**Provider:** opencode-zen

Validates generated tests:
- Syntax correctness
- Test execution
- Coverage measurement

## Usage Examples

### CLI Usage

```bash
# Generate unit tests
biometrics workflow run test-generator \
  --repository https://github.com/org/app \
  --test_type unit \
  --coverage_target 80

# Generate for specific files
biometrics workflow run test-generator \
  --repository https://github.com/org/app \
  --files '["src/utils/helper.ts", "src/api/handler.ts"]'

# Generate integration tests
biometrics workflow run test-generator \
  --repository https://github.com/org/app \
  --test_type integration
```

### Programmatic Usage

```go
engine := workflows.NewWorkflowEngine("./templates")
template, _ := engine.LoadTemplate("test-generator")

instance, _ := engine.CreateInstance(template, map[string]interface{}{
    "repository":     "https://github.com/org/app",
    "test_type":      "unit",
    "coverage_target": 80,
    "files":          []string{"src/utils/parser.ts"},
})

result, _ := engine.Execute(context.Background(), instance)
```

## Configuration

### Framework Configuration

```yaml
inputs:
  framework: jest  # or pytest, go test, etc.
```

### Test Structure

```yaml
options:
  test_structure:
    file_pattern: "{name}.test.{ext}"
    function_pattern: "Test{Name}"
    setup_file: "setup.test.{ext}"
```

## Troubleshooting

### Issue: Framework Not Detected

**Solution:**
```yaml
inputs:
  framework: jest  # Explicitly specify
```

### Issue: Low Coverage

**Solution:**
Increase coverage target or run multiple iterations:
```yaml
inputs:
  coverage_target: 90
```

## Best Practices

### 1. Start with High-Coverage Targets

Aim for 80%+ coverage on critical modules.

### 2. Focus on Critical Paths

Prioritize business-critical functionality.

### 3. Review Generated Tests

AI-generated tests should be reviewed by developers.

### 4. Maintain Test Suite

Keep tests updated as code evolves.

## Related Templates

- **Code Review** (`code-review/`) - Quality validation
- **Bug Fix** (`bug-fix/`) - Fix issues found by tests
- **Refactor** (`refactor/`) - Improve code maintainability

---

**Template Version:** 1.0.0  
**Author:** BIOMETRICS Team  
**Category:** Testing  
**Tags:** testing, unit-tests, integration-tests, coverage, automation

*Last Updated: February 2026*
