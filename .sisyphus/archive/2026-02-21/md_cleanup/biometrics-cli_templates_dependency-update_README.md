# Dependency Update Template

## Overview

The `dependency-update` template provides automated dependency management and updating for the biometrics CLI project. This template handles version checking, compatibility analysis, and safe dependency updates.

## Features

- **Automatic Version Checking**: Scan for outdated dependencies
- **Compatibility Analysis**: Verify version compatibility
- **Security Scanning**: Check for known vulnerabilities
- **Update Automation**: Apply safe updates automatically
- **Rollback Support**: Revert failed updates

## Usage

### Check for Updates
```bash
biometrics-cli deps check
```

Output:
```
Checking dependencies...
┌──────────────────┬─────────┬─────────┬─────────┐
│ Package          │ Current │ Latest  │ Update  │
├──────────────────┼─────────┼─────────┼─────────┤
│ biometrics-core  │ 1.2.0   │ 1.3.0   │ Available│
│ auth-lib         │ 2.0.1   │ 2.0.1   │ Up-to-date│
│ utils            │ 0.5.0   │ 0.6.0   │ Available│
└──────────────────┴─────────┴─────────┴─────────┘
```

### Update Dependencies
```bash
# Update all dependencies
biometrics-cli deps update

# Update specific package
biometrics-cli deps update biometrics-core

# Update with version constraints
biometrics-cli deps update biometrics-core --to 1.3.0
```

### Security Updates Only
```bash
# Apply only security patches
biometrics-cli deps update --security-only
```

## Configuration

```yaml
dependencies:
  update:
    # Auto-update minor versions
    minor: true
    
    # Auto-update patch versions
    patch: true
    
    # Require confirmation for major updates
    major_confirm: true
    
    # Maximum concurrent updates
    parallel: 4
  
  security:
    # Enable security scanning
    scan: true
    
    # Severity threshold (low, medium, high, critical)
    severity_threshold: medium
    
    # Auto-apply security patches
    auto_patch: true
  
  backup:
    # Create backup before updates
    enabled: true
    
    # Backup directory
    directory: "./backups"
    
    # Retention days
    retention: 30
```

## Version Constraints

### Semantic Versioning
The template follows semantic versioning:
- `MAJOR`: Breaking changes
- `MINOR`: New features (backward compatible)
- `PATCH`: Bug fixes (backward compatible)

### Constraint Syntax
```bash
# Exact version
biometrics-cli deps add biometrics-core==1.2.0

# Range
biometrics-cli deps add biometrics-core">=1.0.0 <2.0.0"

# Latest compatible
biometrics-cli deps add biometrics-core@latest
```

## Update Strategy

### Safe Update Order
1. Update dependencies without dependents first
2. Update packages with fewer dependents
3. Run tests after each update
4. Rollback if tests fail

### Update Workflow
```
┌─────────────┐    ┌──────────────┐    ┌───────────┐
│ Check       │───►│ Analyze      │───►│ Update    │
│ Versions    │    │ Compatibility│    │ Package   │
└─────────────┘    └──────────────┘    └───────────┘
                                             │
                                             ▼
                                      ┌───────────┐
                                      │ Verify    │
                                      │ Tests     │
                                      └───────────┘
                                             │
                    ┌────────────────────────┤
                    ▼                        ▼
              ┌───────────┐          ┌───────────┐
              │ Success   │          │ Rollback  │
              │ Commit    │          │ Restore   │
              └───────────┘          └───────────┘
```

## Lock Files

### Generation
```bash
# Generate lock file
biometrics-cli deps lock
```

### Lock File Format
```json
{
  "biometrics-core": {
    "version": "1.2.0",
    "resolved": "https://registry.example.com/biometrics-core-1.2.0.tgz",
    "integrity": "sha256-abc123..."
  }
}
```

## Vulnerability Scanning

### Check for Vulnerabilities
```bash
# Full security scan
biometrics-cli deps audit

# Specific severity
biometrics-cli deps audit --severity critical
```

### Vulnerability Report
```
Vulnerability Report
====================
Package: biometrics-core@1.2.0
Severity: HIGH
CVE: CVE-2026-0123
Description: Remote code execution in biometrics-core
Recommendation: Update to 1.3.0
```

## CI/CD Integration

### GitHub Actions
```yaml
name: Dependency Update

on:
  schedule:
    - cron: '0 0 * * 0'  # Weekly
  workflow_dispatch:

jobs:
  update:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Check Updates
        run: biometrics-cli deps check --json > updates.json
      - name: Create PR
        if: contains(updates.json, 'available')
        run: |
          biometrics-cli deps update
          git checkout -b dependabot/updates
          git commit -a -m "Update dependencies"
          git push
```

### Renovate Integration
```json
{
  "packageRules": [
    {
      "matchPackagePatterns": ["*"],
      "separateMajorMinor": true,
      "schedule": ["before 5am"]
    }
  ]
}
```

## Troubleshooting

### Update Failed
```bash
# View update logs
biometrics-cli deps update --verbose

# Rollback to previous state
biometrics-cli deps rollback
```

### Dependency Conflicts
```bash
# Resolve conflicts
biometrics-cli deps resolve

# Show dependency tree
biometrics-cli deps tree biometrics-core
```

## Maintenance

### Clean Old Backups
```bash
biometrics-cli deps clean --older-than 30d
```

### Update Registry Cache
```bash
biometrics-cli deps cache refresh
```

## Best Practices

1. **Regular Updates**: Check for updates weekly
2. **Test First**: Run tests before committing
3. **Small Batches**: Update in small groups
4. **Review Changes**: Review changelogs before updating

## See Also

- [CLI Commands](../cmd/README.md)
- [Configuration](../docs/configuration.md)
- [Security Guide](./security/README.md)
