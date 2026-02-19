# Tutorial Scripts

**Purpose:** Automated tutorial and onboarding scripts for new users

## Overview

This directory contains executable tutorial scripts that guide users through various aspects of the BIOMETRICS system, from initial setup to advanced features.

## Tutorial Categories

### Getting Started
- `setup-wizard.sh` - Initial setup walkthrough
- `first-command.sh` - Your first command
- `configuration-basics.sh` - Basic configuration

### Intermediate
- `workflow-automation.sh` - Automating workflows
- `agent-configuration.sh` - Setting up agents
- `integration-setup.sh` - Third-party integrations

### Advanced
- `custom-templates.sh` - Creating custom templates
- `performance-tuning.sh` - Optimizing performance
- `troubleshooting.sh` - Common issues and solutions

## Usage

### Run Tutorial
```bash
cd /Users/jeremy/dev/BIOMETRICS
./docs/tutorials/scripts/setup-wizard.sh
```

### Interactive Mode
```bash
./docs/tutorials/scripts/interactive.sh --tutorial getting-started
```

### Specific Topic
```bash
./docs/tutorials/scripts/first-command.sh
```

## Tutorial Structure

Each tutorial includes:
1. **Introduction** - What you'll learn
2. **Prerequisites** - What's needed
3. **Step-by-step** - Instructions
4. **Practice** - Hands-on exercises
5. **Quiz** - Knowledge check

## Creating Custom Tutorials

### Template
```bash
#!/usr/bin/env bash
# Tutorial: Custom Tutorial Name
# Estimated time: 15 minutes

set -e

echo "=== Tutorial: Custom Tutorial ==="
echo "Prerequisites: None"
echo ""

# Step 1
echo "Step 1: Doing something..."
command_to_run

# Step 2
echo "Step 2: Doing more..."
another_command

echo "=== Tutorial Complete ==="
```

## Best Practices

- Test all tutorials locally
- Include error handling
- Provide clear output
- Add progress indicators

## Related Documentation

- [Setup Guide](../setup/)
- [User Guide](../user-guide.md)
- [Examples](../examples/)
