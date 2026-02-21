# Feature Request Template

## Overview

The `feature-request` template provides a standardized format for creating and managing feature requests within the biometrics project. This template ensures consistent feature documentation and tracking.

## Template Structure

```yaml
feature-request:
  template_file: "feature-request.md"
  output_dir: "./feature-requests"
  auto_number: true
```

## Usage

### Create New Feature Request
```bash
biometrics-cli feature new \
  --title "Add Face Recognition" \
  --description "Implement face recognition" \
  --priority high
```

### List Feature Requests
```bash
biometrics-cli feature list

# Filter by status
biometrics-cli feature list --status proposed
biometrics-cli feature list --status accepted
biometrics-cli feature list --status in-progress
```

### Update Feature Request
```bash
biometrics-cli feature update FR-001 \
  --status in-progress \
  --assignee developer@example.com
```

## Feature Request Format

### Template
```markdown
# Feature Request: FR-001

## Summary
Brief description of the feature

## Problem Statement
What problem does this feature solve?

## Proposed Solution
How will this feature work?

## User Stories
- As a [user], I want [feature] so that [benefit]

## Requirements
### Must Have
- [ ] Requirement 1
- [ ] Requirement 2

### Should Have
- [ ] Requirement 3

### Nice to Have
- [ ] Requirement 4

## Design
### UI Mockups
Links to design mockups

### API Changes
```json
{
  "endpoint": "/api/v1/face/recognize",
  "method": "POST"
}
```

## Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2

## Dependencies
- Dependency 1
- Dependency 2

## Risks
| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Risk 1 | Low | High | Mitigation |

## Timeline
- **Estimation**: 2 weeks
- **Start Date**: TBD
- **End Date**: TBD

## Discussion
Link to discussion thread
```

## Status Workflow

```
┌──────────┐    ┌───────────┐    ┌────────────┐    ┌─────────────┐
│ Proposed │───►│ Accepted  │───►│ In-Progress│───►│ Completed  │
└──────────┘    └───────────┘    └────────────┘    └─────────────┘
     │                                      │
     └──────────────┐                       │
                    ▼                       ▼
              ┌──────────┐           ┌──────────┐
              │ Rejected │           │  Blocked │
              └──────────┘           └──────────┘
```

## Priority Levels

| Level | Description | SLA |
|-------|-------------|-----|
| Critical | Security, data loss | 1 week |
| High | Core functionality | 2 weeks |
| Medium | Important features | 1 month |
| Low | Nice to have | Quarter |

## Categories

| Category | Description |
|----------|-------------|
| Security | Security-related features |
| Performance | Performance improvements |
| UX | User experience |
| API | API changes |
| Integration | Third-party integrations |
| Infrastructure | DevOps, deployment |

## Example Feature Request

### FR-015: Add Biometric Liveness Detection

**Problem:**
Currently, the system cannot distinguish between real biometric samples and spoofing attempts using photos or recordings.

**Solution:**
Implement liveness detection that analyzes micro-movements, reflection patterns, and other physiological signals to verify the biometric sample is from a live person.

**Requirements:**
- [ ] Real-time liveness detection
- [ ] Support for face and voice biometrics
- [ ] 99% accuracy on live vs spoof detection
- [ ] < 500ms response time
- [ ] Configurable sensitivity levels

**API Changes:**
```json
{
  "POST /api/v1/biometrics/verify": {
    "liveness_check": true,
    "liveness_threshold": 0.95
  }
}
```

## Voting

### Upvote Features
```bash
biometrics-cli feature vote FR-015
```

### View Top Features
```bash
biometrics-cli feature top

# Top Features
1. FR-015: Liveness Detection (45 votes)
2. FR-023: Mobile SDK (38 votes)
3. FR-031: WebAuthn Support (32 votes)
```

## Roadmap Integration

Features are prioritized into roadmap quarters:

```bash
# Add to roadmap
biometrics-cli feature roadmap FR-015 --quarter Q2-2026
```

## Search and Filter

### Search by Keyword
```bash
biometrics-cli feature search "liveness"
```

### Filter by Multiple Criteria
```bash
biometrics-cli feature list \
  --category security \
  --priority high \
  --status accepted
```

## Export

### Export to CSV
```bash
biometrics-cli feature export --format csv --output features.csv
```

### Export to Markdown
```bash
biometrics-cli feature export --format markdown --output features.md
```

## Integration

### GitHub Issues
```bash
biometrics-cli feature sync --github
```

### Linear Integration
```bash
biometrics-cli feature sync --linear
```

## Best Practices

1. **Clear Problem Statement**: Explain the problem clearly
2. **Detailed Requirements**: Be specific about requirements
3. **Include Examples**: Show concrete use cases
4. **Consider Edge Cases**: Document edge cases
5. **Estimate Effort**: Provide time estimates

## Review Process

### Review Checklist
- [ ] Clear problem statement
- [ ] Feasibility assessed
- [ ] Security implications reviewed
- [ ] Performance impact analyzed
- [ ] Dependencies identified

### Review Meeting
- Weekly feature review meetings
- Stakeholder input required
- Final approval by tech lead

## See Also

- [CLI Commands](../cmd/README.md)
- [Templates Overview](./README.md)
- [Bug Fix Template](./bug-fix/README.md)
