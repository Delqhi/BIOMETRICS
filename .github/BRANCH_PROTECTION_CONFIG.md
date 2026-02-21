# BRANCH PROTECTION CONFIGURATION

**Quick Reference Guide for GitHub Settings**

---

## MAIN BRANCH CONFIGURATION

### Settings Path
`Settings → Branches → Add branch protection rule`

### Rule: `main`

```yaml
Branch name pattern: main

# Pull Request Requirements
pull_request_reviews:
  required_approving_review_count: 2
  dismiss_stale_reviews: true
  require_code_owner_reviews: true
  required_review_teams:
    - biometrics-core

# Status Checks
required_status_checks:
  strict: true  # Require branch up-to-date
  contexts:
    - "ci/lint"
    - "ci/typecheck"
    - "ci/test"
    - "ci/build"
    - "security/scan"

# Restrictions
restrictions: null  # No user/team restrictions
allow_force_pushes: false
allow_deletions: false
enforce_admins: true
required_linear_history: false
required_conversation_resolution: true
lock_branch: false
allow_fork_syncing: false

# Signed Commits (Optional but Recommended)
required_signatures: true  # Recommended for security
```

---

## DEVELOP BRANCH CONFIGURATION

### Rule: `develop`

```yaml
Branch name pattern: develop

# Pull Request Requirements
pull_request_reviews:
  required_approving_review_count: 1
  dismiss_stale_reviews: true
  require_code_owner_reviews: true
  required_review_teams:
    - biometrics-dev

# Status Checks
required_status_checks:
  strict: true
  contexts:
    - "ci/lint"
    - "ci/typecheck"
    - "ci/test"

# Restrictions
allow_force_pushes: true  # Allowed for maintainers
allow_deletions: true  # Allowed for maintainers
enforce_admins: false
required_conversation_resolution: true
```

---

## RELEASE BRANCH CONFIGURATION

### Rule: `release/*`

```yaml
Branch name pattern: release/.*

# Pull Request Requirements
pull_request_reviews:
  required_approving_review_count: 2
  dismiss_stale_reviews: true
  require_code_owner_reviews: true
  required_review_teams:
    - biometrics-core
    - biometrics-qa

# Status Checks
required_status_checks:
  strict: true
  contexts:
    - "ci/lint"
    - "ci/typecheck"
    - "ci/test"
    - "ci/build"
    - "e2e/test"
    - "security/scan"

# Restrictions
allow_force_pushes: false
allow_deletions: false
enforce_admins: true
required_conversation_resolution: true
```

---

## HOTFIX BRANCH CONFIGURATION

### Rule: `hotfix/*`

```yaml
Branch name pattern: hotfix/.*

# Pull Request Requirements
pull_request_reviews:
  required_approving_review_count: 1
  dismiss_stale_reviews: true
  require_code_owner_reviews: true
  required_review_teams:
    - biometrics-core

# Status Checks
required_status_checks:
  strict: false  # Optional for speed
  contexts:
    - "ci/lint"
    - "ci/typecheck"
    - "ci/test-critical"

# Restrictions
allow_force_pushes: false
allow_deletions: true  # Allowed after merge
enforce_admins: true
required_conversation_resolution: true
```

---

## FEATURE BRANCH CONFIGURATION (Optional)

### Rule: `feature/*` (Long-running features only)

```yaml
Branch name pattern: feature/.*

# Pull Request Requirements (Optional)
pull_request_reviews:
  required_approving_review_count: 1
  dismiss_stale_reviews: false
  require_code_owner_reviews: false

# Status Checks (Minimal)
required_status_checks:
  strict: false
  contexts:
    - "ci/lint"
    - "ci/typecheck"

# Permissive
allow_force_pushes: true
allow_deletions: true
enforce_admins: false
```

---

## GITHUB API AUTOMATION

### Using GitHub CLI

```bash
# Main Branch Protection
gh api \
  --method PUT \
  /repos/{owner}/{repo}/branches/main/protection \
  -f required_status_checks='{"strict":true,"contexts":["ci/lint","ci/typecheck","ci/test","ci/build","security/scan"]}' \
  -f required_pull_request_reviews='{"required_approving_review_count":2,"dismiss_stale_reviews":true,"require_code_owner_reviews":true}' \
  -f enforce_admins=true \
  -f allow_force_pushes=false \
  -f allow_deletions=false

# Develop Branch Protection
gh api \
  --method PUT \
  /repos/{owner}/{repo}/branches/develop/protection \
  -f required_status_checks='{"strict":true,"contexts":["ci/lint","ci/typecheck","ci/test"]}' \
  -f required_pull_request_reviews='{"required_approving_review_count":1,"dismiss_stale_reviews":true}' \
  -f enforce_admins=false \
  -f allow_force_pushes=true
```

### Using Terraform

```hcl
resource "github_branch_protection" "main" {
  repository_id  = github_repository.biometrics.id
  pattern        = "main"

  required_pull_request_reviews {
    required_approving_review_count = 2
    dismiss_stale_reviews           = true
    require_code_owner_reviews      = true
  }

  required_status_checks {
    strict = true
    contexts = [
      "ci/lint",
      "ci/typecheck",
      "ci/test",
      "ci/build",
      "security/scan"
    ]
  }

  enforce_admins               = true
  is_admin_enforced            = true
  allows_deletions             = false
  allows_force_pushes          = false
  requires_commit_signatures   = true
  requires_conversation_resolution = true
}

resource "github_branch_protection" "develop" {
  repository_id  = github_repository.biometrics.id
  pattern        = "develop"

  required_pull_request_reviews {
    required_approving_review_count = 1
    dismiss_stale_reviews           = true
  }

  required_status_checks {
    strict = true
    contexts = [
      "ci/lint",
      "ci/typecheck",
      "ci/test"
    ]
  }

  enforce_admins               = false
  allows_deletions             = true
  allows_force_pushes          = true
  requires_commit_signatures   = false
}
```

---

## REQUIRED STATUS CHECKS DEFINITIONS

### CI/CD Workflows

All status checks are defined in `.github/workflows/`:

| Check Name | Workflow File | Purpose | Required For |
|------------|---------------|---------|--------------|
| `ci/lint` | `ci.yml` | ESLint, Prettier, Go fmt | main, develop |
| `ci/typecheck` | `ci.yml` | TypeScript compilation | main, develop |
| `ci/test` | `ci.yml` | Unit tests (Jest, Go test) | main, develop |
| `ci/build` | `ci.yml` | Production build | main, release |
| `security/scan` | `security.yml` | Dependency audit, CodeQL | main, release |
| `e2e/test` | `e2e.yml` | End-to-end tests | release |

### Workflow Triggers

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Linters
        run: pnpm run lint

  typecheck:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: TypeScript Check
        run: pnpm run typecheck

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Tests
        run: pnpm test

  build:
    runs-on: ubuntu-latest
    needs: [lint, typecheck, test]
    steps:
      - uses: actions/checkout@v4
      - name: Build
        run: pnpm run build
```

---

## MERGE STRATEGY CONFIGURATION

### Repository Settings

`Settings → General → Pull Requests`

```
✅ Allow merge commits
   └─ Default: Yes (for feature branches)

✅ Allow squashing commits
   └─ Default: Yes (for simple features)

❌ Allow rebasing commits
   └─ Disabled (to preserve merge history)
```

### Strategy Selection

| Branch Type | Strategy | Reason |
|-------------|----------|--------|
| `main` ← `release/*` | **Merge Commit** | Preserves release history |
| `develop` ← `feature/*` | **Squash** | Clean history |
| `release/*` ← `hotfix/*` | **Merge Commit** | Track hotfixes |
| `main` ← `hotfix/*` | **Squash** | Single commit for hotfix |

---

## EXCEPTIONS & OVERRIDES

### Emergency Bypass

For critical production fixes only:

1. **Approval Required:**
   - @biometrics-core team member
   - Documented justification in PR

2. **Process:**
   ```
   - Create hotfix branch
   - Implement minimal fix
   - Create PR with "emergency" label
   - Notify team via Slack/Teams
   - Minimum 1 approval (expedited)
   - Merge immediately
   - Post-merge review within 24h
   ```

3. **Documentation:**
   - Add entry to `CHANGELOG.md`
   - Create incident report
   - Schedule post-mortem

### Admin Override

Repository administrators can override protections in emergencies:

```bash
# Force push (if absolutely necessary)
git push origin +main:main

# Merge without status checks (GitHub UI)
# Click "Merge" button with override confirmation
```

**WARNING:** Admin overrides are logged and audited. Use sparingly.

---

## MONITORING & ALERTS

### Branch Protection Violations

GitHub automatically logs all protection violations:

`Insights → Pulse → Merge stats`

### Key Metrics to Monitor

| Metric | Target | Alert Threshold |
|--------|--------|-----------------|
| PR Review Time | < 24h | > 48h |
| PR Merge Time | < 48h | > 72h |
| Stale PRs | < 5 | > 10 |
| Override Count | 0 | > 1 per month |

### Slack/Teams Integration

```yaml
# .github/workflows/notifications.yml
name: Branch Protection Alerts

on:
  pull_request:
    types: [closed]

jobs:
  notify:
    runs-on: ubuntu-latest
    steps:
      - name: Send Slack notification
        uses: 8398a7/action-slack@v3
        with:
          status: ${{ job.status }}
          text: |
            PR Merged: ${{ github.event.pull_request.title }}
            Branch: ${{ github.event.pull_request.head.ref }} → ${{ github.event.pull_request.base.ref }}
            Author: ${{ github.event.pull_request.user.login }}
            Merged by: ${{ github.event.pull_request.merged_by.login }}
          webhook_url: ${{ secrets.SLACK_WEBHOOK }}
```

---

## TROUBLESHOOTING

### Common Issues

#### Issue: "Required status checks are missing"

**Solution:**
```bash
# Check which checks are required
gh api /repos/{owner}/{repo}/branches/main/protection | jq .required_status_checks.contexts

# Re-run failed checks
gh workflow run ci.yml --ref <branch-name>
```

#### Issue: "Code owner review required"

**Solution:**
```bash
# Check CODEOWNERS file
cat .github/CODEOWNERS

# Verify team membership
gh api /orgs/{org}/teams/{team}/members
```

#### Issue: "Branch is not up to date"

**Solution:**
```bash
# Update branch
git fetch origin
git rebase origin/main

# Or merge (if rebase not allowed)
git merge origin/main
```

#### Issue: "Force push blocked"

**Solution:**
```bash
# NEVER force push to protected branches
# Use regular push instead
git push origin <branch-name>

# If you MUST (admin only):
git push origin +<branch-name>:<branch-name>
```

---

## AUDIT & COMPLIANCE

### Audit Log Access

`Settings → Archives → Audit log`

Download audit logs for compliance:

```bash
# Export audit log (last 90 days)
gh api \
  /orgs/{org}/audit-log \
  -f phrase="action:branch.protection_rule" \
  --paginate > audit-log.json
```

### Required Audit Events

| Event | Description | Retention |
|-------|-------------|-----------|
| `branch.protection_rule.created` | Rule created | 7 years |
| `branch.protection_rule.edited` | Rule modified | 7 years |
| `branch.protection_rule.deleted` | Rule deleted | 7 years |
| `pull_request.merged` | PR merged | 7 years |
| `branch.force_push` | Force push (admin) | 7 years |

### SOC2 Compliance

For SOC2 audits, maintain:

1. **Change Management:**
   - All changes via PR
   - Minimum 2 approvals
   - Status checks pass

2. **Access Control:**
   - Code owners enforced
   - Admin access logged
   - Team membership reviewed quarterly

3. **Audit Trail:**
   - Complete git history
   - PR comments preserved
   - Status check results archived

---

## VERSION HISTORY

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| 1.0.0 | 2026-02-21 | Initial release | @biometrics-devops |

---

## RELATED DOCUMENTATION

- [CODEOWNERS File](/.github/CODEOWNERS)
- [Branch Protection Documentation](/docs/BRANCH-PROTECTION.md)
- [Contributing Guidelines](/.github/CONTRIBUTING.md)
- [Pull Request Template](/.github/PULL_REQUEST_TEMPLATE.md)
- [Security Policy](/.github/SECURITY.md)

---

**Last Updated:** 2026-02-21  
**Owner:** @biometrics-devops  
**Status:** ACTIVE
