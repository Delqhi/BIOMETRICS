# Bug Fix Template

## Overview

The `bug-fix` template provides a standardized format for reporting, tracking, and fixing bugs in the biometrics system. This template ensures consistent bug documentation and resolution.

## Features

- **Bug Reporting**: Structured bug reports
- **Reproduction Steps**: Clear reproduction instructions
- **Impact Assessment**: Severity and priority classification
- **Fix Verification**: Test case templates
- **Regression Prevention**: Checklist for preventing recurrence

## Usage

### Create Bug Report
```bash
biometrics-cli bug new \
  --title "Login fails with special characters" \
  --severity high \
  --component authentication
```

### List Bugs
```bash
biometrics-cli bug list
biometrics-cli bug list --severity critical
biometrics-cli bug list --status open
```

### Update Bug Status
```bash
biometrics-cli bug update BUG-123 \
  --status in-progress \
  --assignee developer@example.com
```

## Bug Report Format

### Template
```markdown
# Bug Report: BUG-XXX

## Summary
One-line description of the bug

## Environment
- **Version**: 1.5.0
- **OS**: Ubuntu 22.04
- **Browser**: Chrome 120
- **Database**: PostgreSQL 15

## Severity
- [ ] Critical - Data loss, security breach
- [ ] High - Major feature broken
- [ ] Medium - Feature partially working
- [ ] Low - Minor issue, workaround available

## Priority
- [ ] P0 - Immediate fix required
- [ ] P1 - Fix in current sprint
- [ ] P2 - Fix in next sprint
- [ ] P3 - Backlog

## Component
- Authentication
- User Management
- Biometrics Processing
- API
- CLI
- Documentation

## Description
Detailed description of the bug

## Steps to Reproduce
1. Navigate to login page
2. Enter credentials with special chars (e.g., user@domain.com!)
3. Click submit
4. Observe error

## Expected Behavior
User should be able to login with any valid characters

## Actual Behavior
Error message: "Invalid characters in password"

## Workaround
Encode special characters before submission

## Root Cause
Input validation rejects all special characters instead of dangerous ones only

## Fix Description
1. Update validation regex to allow safe special characters
2. Add unit test for special character handling
3. Update error message

## Code Changes

### File: auth/validator.py
```python
# Before
def validate_password(password):
    return re.match(r'^[a-zA-Z0-9]+$', password)

# After  
def validate_password(password):
    return re.match(r'^[a-zA-Z0-9!@#$%^&*()]+$', password)
```

## Test Cases

### Unit Tests
```python
def test_password_special_characters():
    assert validate_password("password!") == True
    assert validate_password("pass@word") == True
    assert validate_password("p@ssw0rd!") == True
```

### Integration Tests
- [ ] Login with special characters
- [ ] Login with all allowed special chars
- [ ] Login with disallowed characters

## Verification

### Pre-Fix Behavior
```
$ curl -X POST /api/login -d '{"email":"test@example.com","password":"test!"}'
{"error": "Invalid characters"}
```

### Post-Fix Behavior
```
$ curl -X POST /api/login -d '{"email":"test@example.com","password":"test!"}'
{"token": "eyJhbGc..."}
```

## Impact Analysis

### Users Affected
All users attempting to login with special characters in password

### Data Integrity
No data corruption risk

### Security Implications
- Positive: More permissive (no security issue)
- Negative: None

## Dependencies
- None

## Related Issues
- Related to: ENH-456 (Password validation improvement)

## Time Tracking
- **Estimated**: 2 hours
- **Actual**: 1.5 hours

## Status History
- 2026-02-01: Created
- 2026-02-01: Assigned
- 2026-02-01: In Progress
- 2026-02-01: Fixed
- 2026-02-01: Verified
- 2026-02-01: Closed
```

## Severity Guidelines

### Critical (S1)
- Complete system failure
- Data loss or corruption
- Security vulnerability
- **Response Time**: 1 hour

### High (S2)
- Major feature not working
- Workaround difficult
- **Response Time**: 4 hours

### Medium (S3)
- Feature partially working
- Workaround available
- **Response Time**: 1 day

### Low (S4)
- Minor issue
- Cosmetic issue
- **Response Time**: 1 week

## Bug Lifecycle

```
┌────────┐    ┌─────────┐    ┌────────────┐    ┌──────────┐    ┌────────┐
│ Created │───►│Assigned │───►│ In-Progress│───►│  Fixed   │───►│Closed │
└────────┘    └─────────┘    └────────────┘    └──────────┘    └────────┘
     │                                                    │
     │                       ┌──────────┐                 │
     └─────────────────────►│ Reopened │◄─────────────────┘
                             └──────────┘
```

## Integration

### GitHub Issues
```bash
biometrics-cli bug sync --github
```

### Slack Notifications
```yaml
# Config
notifications:
  slack:
    critical: "#critical-bugs"
    high: "#bugs"
    medium: "#bugs"
    low: "#backlog"
```

## Testing Checklist

### Before Fix
- [ ] Reproduce the bug
- [ ] Document exact steps

### During Fix
- [ ] Make minimal changes
- [ ] Add test cases

### After Fix
- [ ] Verify fix works
- [ ] Run full test suite
- [ ] Check for regressions
- [ ] Update documentation

## Best Practices

1. **Reproduce First**: Always reproduce before fixing
2. **Minimal Changes**: Fix only what's broken
3. **Test Thoroughly**: Cover edge cases
4. **Document Well**: Clear reproduction steps
5. **Verify Fix**: Confirm the fix works

## Tools

### Debug Commands
```bash
biometrics-cli debug logs --bug BUG-123
biometrics-cli debug trace --request-id req-456
```

### Test Commands
```bash
biometrics-cli test run --bug BUG-123
biometrics-cli test coverage --component auth
```

## See Also

- [CLI Commands](../cmd/README.md)
- [Templates Overview](./README.md)
- [Feature Request Template](./feature-request/README.md)
