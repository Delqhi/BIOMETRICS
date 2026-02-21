# Thinking Logs Directory

## Overview

This directory contains thinking logs, reasoning traces, and decision documentation generated during development and problem-solving sessions. These logs capture the AI's thought process and decision-making rationale.

## Contents

### Session Logs

| Directory | Description |
|-----------|-------------|
| sessions/ | Individual session logs |
| decisions/ | Key decision records |
| research/ | Research notes |

### Thinking Traces

| File | Description |
|------|-------------|
| architecture-thought.md | Architecture decisions |
| implementation-thought.md | Implementation reasoning |
| debugging-thought.md | Debug process logs |

### Analysis Logs

| File | Description |
|------|-------------|
| performance-analysis.md | Performance investigations |
| security-analysis.md | Security reviews |
| code-review.md | Code review notes |

## Log Format

### Session Header
```markdown
# Session: [Date] - [Topic]

**Date**: 2026-02-20
**Duration**: 45 minutes
**Goal**: [Objective]
**Participants**: [Team/AI]
```

### Thought Entry
```markdown
## [Timestamp] - Thinking

**Observation**: What was noticed
**Hypothesis**: What was considered
**Action**: What was decided
**Result**: What happened
```

## Usage

### Review
```bash
# View recent logs
ls -la logs/thinking/

# Search logs
grep -r "authentication" logs/thinking/
```

### Archival
```bash
# Archive old logs
mkdir -p logs/thinking/archive/2026-01
mv logs/thinking/*.md logs/thinking/archive/2026-01/
```

## Best Practices

### Documentation
- **Timestamp**: Always include date/time
- **Context**: Explain background
- **Rationale**: Document reasoning
- **Outcome**: Note results

### Maintenance
- Review weekly
- Archive monthly
- Remove after 6 months

## Privacy

### Sensitive Data
- Remove API keys
- Anonymize user data
- Don't log credentials

### Access Control
- Limit access to team
- Review before sharing

## Integration

### AI Sessions
```bash
# Save thinking log
biometrics-cli session save --thinking
```

### Code Reviews
```bash
# Generate review log
biometrics-cli review generate --log
```

## Tools

### Log Analysis
```bash
# Analyze thinking patterns
biometrics-cli analyze logs --pattern "hypothesis"
```

## Maintenance

- Archive quarterly
- Clean temporary files
- Verify completeness

## See Also

- [Logs Overview](../logs/)
- [Scripts Directory](../scripts/)
- [Documentation](../docs/)
