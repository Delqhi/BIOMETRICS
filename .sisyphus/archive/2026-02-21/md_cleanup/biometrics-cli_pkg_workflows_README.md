# Workflows Package

**Purpose:** Workflow management and orchestration utilities

## Overview

This package provides workflow management capabilities for the biometrics CLI, enabling complex multi-step operations with proper state management and error handling.

## Features

- Workflow definition and execution
- State management
- Retry policies
- Progress tracking
- Cancellation support

## Usage

### Creating a Workflow

```go
import "github.com/delqhi/biometrics/pkg/workflows"

workflow := workflows.NewWorkflow(ctx, workflows.Config{
    Name: "process-biometrics",
    Steps: []workflows.Step{
        {
            Name: "validate-input",
            Action: validateInput,
            Retry: 3,
        },
        {
            Name: "process-data",
            Action: processData,
            Timeout: 5 * time.Minute,
        },
        {
            Name: "save-results",
            Action: saveResults,
        },
    },
})
```

### Executing

```go
// Execute workflow
result, err := workflow.Execute(ctx, input)

// Stream progress
stream, err := workflow.ExecuteStream(ctx, input)
for update := range stream {
    fmt.Printf("Progress: %d%%\n", update.Progress)
}
```

## Workflow Configuration

```yaml
workflows:
  process:
    name: "Biometrics Processing"
    max_retries: 3
    timeout: 30m
    steps:
      - name: validate
        action: validate
        retry: 2
      - name: process
        action: process
        parallel: 4
      - name: save
        action: save
```

## Error Handling

```go
result, err := workflow.Execute(ctx, input)
if err != nil {
    switch {
    case errors.Is(err, workflows.ErrStepFailed):
        // Specific step failed
        fmt.Printf("Step %s failed: %v\n", err.(workflows.StepError).Step, err)
    case errors.Is(err, workflows.ErrTimeout):
        // Workflow timed out
    case errors.Is(err, workflows.ErrCancelled):
        // Workflow was cancelled
    }
}
```

## Related Documentation

- [Workflow Examples](../examples/workflows/)
- [CLI Commands](../commands/)
