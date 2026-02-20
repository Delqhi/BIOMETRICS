# Refactor Template

## Overview

The `refactor` template provides a structured approach for code refactoring activities. This template ensures refactoring is documented, tested, and maintains backward compatibility.

## Features

- **Impact Analysis**: Assess code changes and dependencies
- **Migration Guides**: Auto-generate migration documentation
- **Deprecation Warnings**: Track deprecated APIs
- **Breaking Changes**: Document and communicate breaking changes
- **Rollback Plans**: Prepare rollback strategies

## Usage

### Start Refactoring
```bash
biometrics-cli refactor start \
  --target biometrics-core \
  --scope "authentication" \
  --type "restructure"
```

### Track Progress
```bash
biometrics-cli refactor status
biometrics-cli refactor progress
```

### Complete Refactoring
```bash
biometrics-cli refactor complete
```

## Refactoring Types

### Code Organization
- Module restructuring
- File organization
- Naming conventions

### Performance
- Algorithm optimization
- Caching implementation
- Query optimization

### Maintainability
- Reduce complexity
- Improve readability
- Remove technical debt

### Security
- Fix vulnerabilities
- Improve encryption
- Enhance authentication

## Template Structure

```markdown
# Refactoring: REF-XXX

## Overview
Brief description of refactoring

## Motivation
Why is this refactoring needed?

## Scope
What files/components are affected?

## Changes

### Before
```python
# Old code
def authenticate(user):
    return user.authenticated
```

### After
```python
# New code
async def authenticate(user, context):
    return await user.verify(context)
```

## Impact Analysis

### Affected Components
| Component | Impact Level | Changes |
|-----------|--------------|---------|
| API | High | Breaking |
| CLI | Medium | Update required |
| SDK | High | Version bump |

### Dependencies
- Dependency A
- Dependency B

## Migration Guide

### Step 1: Update Dependencies
```bash
pip install biometrics-core>=2.0.0
```

### Step 2: Update Code
```python
# Old
result = authenticate(user)

# New
result = await authenticate(user, context)
```

### Step 3: Test
```bash
pytest tests/
```

## Breaking Changes

| Old API | New API | Migration |
|---------|---------|-----------|
| authenticate() | authenticate_async() | Add await |
| get_user(id) | get_user(id, opts) | Add opts parameter |

## Deprecation Timeline

| Feature | Deprecated | Removal | Alternative |
|---------|------------|---------|-------------|
| v1 API | 2026-01 | 2026-07 | v2 API |

## Testing Strategy

### Unit Tests
- [ ] Test all affected functions
- [ ] Test edge cases
- [ ] Test error handling

### Integration Tests
- [ ] Test component interactions
- [ ] Test API endpoints

### Migration Tests
- [ ] Test backward compatibility
- [ ] Test upgrade path

## Rollback Plan

### If Issues Occur
1. Revert code changes
2. Deploy previous version
3. Investigate issues

### Rollback Command
```bash
git revert <commit>
```

## Performance Impact

### Before
- Function calls: 1000/sec
- Latency: 50ms

### After
- Function calls: 5000/sec
- Latency: 10ms

## Documentation Updates

### Files Changed
- README.md
- API Documentation
- Migration Guide

## Approval

- [ ] Code Review
- [ ] Security Review
- [ ] Product Approval

## Timeline

- **Start Date**: 2026-02-01
- **Target Completion**: 2026-02-15
- **Actual Completion**: TBD
```

## Checklist

### Pre-Refactoring
- [ ] Analyze dependencies
- [ ] Identify impact areas
- [ ] Create backup
- [ ] Plan rollback

### During Refactoring
- [ ] Make incremental changes
- [ ] Run tests frequently
- [ ] Update documentation

### Post-Refactoring
- [ ] Run full test suite
- [ ] Update migration guide
- [ ] Communicate changes
- [ ] Monitor production

## CI/CD Integration

### Pre-deployment Check
```yaml
- name: Refactoring Check
  run: |
    biometrics-cli refactor validate \
      --changes ./refactor/changes.json
```

### Migration Test
```yaml
- name: Migration Test
  run: |
    biometrics-cli refactor test-migration \
      --from v1.0.0 \
      --to v2.0.0
```

## Best Practices

1. **Small Iterations**: Make small, testable changes
2. **Document Everything**: Write migration guides
3. **Test Thoroughly**: Run comprehensive tests
4. **Communicate Early**: Notify stakeholders
5. **Plan Rollback**: Always have rollback plan

## Tools

### Code Analysis
```bash
biometrics-cli refactor analyze \
  --target biometrics-core \
  --metrics complexity,coverage
```

### Dependency Graph
```bash
biometrics-cli refactor graph \
  --target biometrics-core \
  --output deps.svg
```

## Common Patterns

### Extract Method
```python
# Before
def process(data):
    # 50 lines of code
    return result

# After
def process(data):
    validate(data)
    transform(data)
    save(data)
    return result
```

### Replace Conditional
```python
# Before
if type == 'A':
    do_a()
elif type == 'B':
    do_b()

# After
handlers[type]()
```

## See Also

- [CLI Commands](../cmd/README.md)
- [Templates Overview](./README.md)
- [Bug Fix Template](./bug-fix/README.md)
