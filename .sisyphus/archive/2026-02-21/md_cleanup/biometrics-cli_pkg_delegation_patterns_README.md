# Delegation Patterns Directory

## Overview

This directory contains reusable delegation patterns for the biometrics CLI. These patterns provide templates for implementing various delegation and authorization workflows.

## Contents

### Authorization Patterns

| File | Description |
|------|-------------|
| role-based.md | Role-based access control |
| policy-based.md | Policy-based delegation |
| hierarchical.md | Hierarchical delegation |

### Implementation Patterns

| File | Description |
|------|-------------|
| api-delegation.md | API request delegation |
| batch-delegation.md | Batch operation delegation |
| async-delegation.md | Async task delegation |

## Role-Based Access

### Pattern Overview
```yaml
roles:
  admin:
    permissions:
      - read
      - write
      - delete
      - delegate
  manager:
    permissions:
      - read
      - write
  user:
    permissions:
      - read
```

### Implementation
```python
from biometrics.delegation import RoleBasedAccess

access = RoleBasedAccess()

# Check permission
if access.has_permission('user123', 'admin', 'write'):
    perform_action()
```

## Policy-Based Delegation

### Pattern Overview
Policies define delegation rules:
```yaml
policies:
  - name: "team-lead"
    effect: "allow"
    actions:
      - "users:read"
      - "users:write"
    conditions:
      - resource.team == user.team

  - name: "cross-team"
    effect: "allow"
    actions:
      - "metrics:read"
    conditions:
      - user.role == "manager"
```

### Implementation
```python
from biometrics.delegation import PolicyEngine

engine = PolicyEngine(policies)

# Evaluate policy
result = engine.evaluate(
    user=user,
    action='users:write',
    resource=team_resource
)
```

## Hierarchical Delegation

### Pattern Overview
```
Organization
├── CEO
│   ├── CTO
│   │   ├── Engineering Manager
│   │   │   └── Senior Engineer
│   │   └── Product Manager
│   └── CFO
```

### Implementation
```python
from biometrics.delegation import HierarchicalDelegation

delegation = HierarchicalDelegation()

# Delegate upward
delegation.delegate(
    from_user='senior_eng',
    to_user='eng_manager',
    scope=['code:read', 'code:write']
)

# Delegate downward
delegation.delegate(
    from_user='cto',
    to_user='eng_manager',
    scope=['deploy:read']
)
```

## API Delegation

### Pattern
```python
from biometrics.delegation import APIDelegation

api = APIDelegation()

# Create delegated token
token = api.create_delegated_token(
    delegator='user123',
    delegatee='service456',
    scopes=['users:read'],
    expires_in=3600
)

# Use delegated token
response = api.request(
    '/api/users',
    token=token
)
```

## Batch Delegation

### Pattern
```python
from biometrics.delegation import BatchDelegation

batch = BatchDelegation()

# Create batch operation
operation = batch.create(
    operations=[
        {'action': 'update', 'resource': 'user1'},
        {'action': 'delete', 'resource': 'user2'},
        {'action': 'create', 'resource': 'user3'}
    ],
    delegate='service456',
    auth='user123'
)

# Execute batch
results = batch.execute(operation)
```

## Async Delegation

### Pattern
```python
from biometrics.delegation import AsyncDelegation

async_delegation = AsyncDelegation()

# Create async task
task = async_delegation.create_task(
    operation='process_biometrics',
    params={'batch_id': 'batch123'},
    delegate='worker456',
    callback='https://example.com/webhook'
)

# Poll for results
while not task.complete:
    task = async_delegation.get_status(task.id)
    await asyncio.sleep(1)
```

## Security Considerations

### Token Management
- Use short-lived tokens
- Implement token rotation
- Monitor token usage
- Revoke compromised tokens

### Audit Logging
```python
# Log all delegation actions
logger.info('delegation_created', {
    'delegator': delegator,
    'delegatee': delegatee,
    'scopes': scopes,
    'timestamp': datetime.utcnow()
})
```

## Best Practices

1. **Least Privilege**: Grant minimum necessary permissions
2. **Time Limits**: Set expiration on delegations
3. **Audit**: Log all delegation activities
4. **Review**: Regularly review delegation policies

## Testing

### Unit Tests
```python
def test_delegation_permissions():
    access = RoleBasedAccess()
    
    # Test role assignment
    assert access.has_permission('user', 'admin', 'delete')
    
    # Test delegation
    assert access.has_permission('delegate', 'user', 'read')
```

## Maintenance

- Review policies quarterly
- Remove expired delegations
- Update role assignments

## See Also

- [CLI Commands](../cmd/README.md)
- [Documentation](../docs/)
- [Configuration](../docs/configuration.md)
