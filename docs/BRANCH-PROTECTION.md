# BRANCH PROTECTION & CODE REVIEW STANDARDS

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-21  
**Status:** ACTIVE - MANDATORY COMPLIANCE  
**Owner:** @biometrics-devops @biometrics-core  

---

## TABLE OF CONTENTS

<!-- TOC start (generated with https://github.com/derlin/bitdowntoc) -->

- [1. EXECUTIVE SUMMARY](#1-executive-summary)
  - [1.1 Purpose](#11-purpose)
  - [1.2 Scope](#12-scope)
  - [1.3 Compliance Requirements](#13-compliance-requirements)
- [2. BRANCH PROTECTION RULES](#2-branch-protection-rules)
  - [2.1 Main Branch Protection](#21-main-branch-protection)
  - [2.2 Develop Branch Protection](#22-develop-branch-protection)
  - [2.3 Feature Branch Guidelines](#23-feature-branch-guidelines)
  - [2.4 Release Branch Strategy](#24-release-branch-strategy)
  - [2.5 Hotfix Branch Strategy](#25-hotfix-branch-strategy)
- [3. CODEOWNERS CONFIGURATION](#3-codeowners-configuration)
  - [3.1 File Path Coverage](#31-file-path-coverage)
  - [3.2 Team Definitions](#32-team-definitions)
  - [3.3 Review Assignment Rules](#33-review-assignment-rules)
- [4. PULL REQUEST REQUIREMENTS](#4-pull-request-requirements)
  - [4.1 PR Template](#41-pr-template)
  - [4.2 Required Reviews](#42-required-reviews)
  - [4.3 Status Checks](#43-status-checks)
  - [4.4 Commit Standards](#44-commit-standards)
- [5. MERGE STRATEGIES](#5-merge-strategies)
  - [5.1 Merge Commit Strategy](#51-merge-commit-strategy)
  - [5.2 Squash and Merge Strategy](#52-squash-and-merge-strategy)
  - [5.3 Rebase and Merge Strategy](#53-rebase-and-merge-strategy)
  - [5.4 Strategy Selection Matrix](#54-strategy-selection-matrix)
- [6. CONFLICT RESOLUTION](#6-conflict-resolution)
  - [6.1 Prevention Strategies](#61-prevention-strategies)
  - [6.2 Resolution Process](#62-resolution-process)
  - [6.3 Escalation Path](#63-escalation-path)
- [7. GITHUB ACTIONS INTEGRATION](#7-github-actions-integration)
  - [7.1 Required Status Checks](#71-required-status-checks)
  - [7.2 Optional Status Checks](#72-optional-status-checks)
  - [7.3 Custom Workflows](#73-custom-workflows)
- [8. SECURITY & COMPLIANCE](#8-security--compliance)
  - [8.1 Signed Commits](#81-signed-commits)
  - [8.2 Secret Scanning](#82-secret-scanning)
  - [8.3 Dependency Review](#83-dependency-review)
  - [8.4 Audit Trail](#84-audit-trail)
- [9. EXCEPTIONS & OVERRIDES](#9-exceptions--overrides)
  - [9.1 Emergency Bypass Process](#91-emergency-bypass-process)
  - [9.2 Admin Override Policy](#92-admin-override-policy)
  - [9.3 Documentation Requirements](#93-documentation-requirements)
- [10. MONITORING & METRICS](#10-monitoring--metrics)
  - [10.1 PR Metrics Dashboard](#101-pr-metrics-dashboard)
  - [10.2 Review Time SLAs](#102-review-time-slas)
  - [10.3 Quality Metrics](#103-quality-metrics)
- [11. TROUBLESHOOTING](#11-troubleshooting)
  - [11.1 Common Issues](#111-common-issues)
  - [11.2 FAQ](#112-faq)
  - [11.3 Support Contacts](#113-support-contacts)
- [APPENDIX A: GITHUB API AUTOMATION](#appendix-a-github-api-automation)
- [APPENDIX B: TEMPLATE FILES](#appendix-b-template-files)
- [APPENDIX C: CHECKLISTS](#appendix-c-checklists)

<!-- TOC end -->

---

## 1. EXECUTIVE SUMMARY

### 1.1 Purpose

This document establishes enterprise-grade branch protection rules and code review standards for the BIOMETRICS repository. These standards ensure:

- **Code Quality:** All changes undergo rigorous review before merging
- **Security:** Sensitive changes require appropriate approvals
- **Stability:** Main branch remains production-ready at all times
- **Compliance:** Full audit trail of all changes
- **Collaboration:** Clear ownership and review responsibilities

### 1.2 Scope

These rules apply to:

- **All Contributors:** Internal team members and external contributors
- **All Code Changes:** Features, bug fixes, documentation, configuration
- **All Branches:** Main, develop, feature, release, and hotfix branches
- **All Repositories:** BIOMETRICS primary repository and related sub-modules

### 1.3 Compliance Requirements

**MANDATORY COMPLIANCE:** All contributors MUST adhere to these standards. Violations may result in:

- Pull request rejection
- Required rework
- Team notification
- Process review meetings

**Automatic Enforcement:** GitHub branch protection rules enforce compliance automatically. Bypassing these rules requires admin privileges and documented justification.

---

## 2. BRANCH PROTECTION RULES

### 2.1 Main Branch Protection

The `main` branch is the production-ready branch. All changes MUST go through pull requests with strict review requirements.

#### 2.1.1 Protection Settings (GitHub UI)

**Location:** Settings → Branches → Add branch protection rule  
**Branch name pattern:** `main`

#### 2.1.2 Required Settings

```
✅ Require a pull request before merging
   ├─ Require approvals: 2
   ├─ Dismiss stale pull request approvals when new commits are pushed: ✅
   ├─ Require review from Code Owners: ✅
   └─ Require review from at least one team: @biometrics-core

✅ Require status checks to pass before merging
   ├─ Required status checks:
   │   ├─ ci/lint (ESLint, Prettier)
   │   ├─ ci/typecheck (TypeScript)
   │   ├─ ci/test (Unit Tests)
   │   ├─ ci/build (Build Success)
   │   └─ security/scan (Security Scan)
   ├─ Require branches to be up to date before merging: ✅
   └─ Allow merges even if some status checks are missing: ❌

✅ Require conversation resolution before merging
   └─ All comments must be resolved: ✅

✅ Include administrators
   └─ Apply rules to repository admins: ✅

❌ Allow force pushes
   └─ Permit force pushes: Disabled

❌ Allow deletions
   └─ Permit branch deletions: Disabled

✅ Require signed commits
   └─ Verify commit signatures: Recommended (Optional)
```

#### 2.1.3 Enforcement Level

| Setting | Enforcement | Override Allowed |
|---------|-------------|------------------|
| PR Required | **STRICT** | Admin Only |
| Approvals (2) | **STRICT** | No |
| Code Owners | **STRICT** | No |
| Status Checks | **STRICT** | No |
| Force Push | **BLOCKED** | Admin Only |
| Branch Delete | **BLOCKED** | Admin Only |

### 2.2 Develop Branch Protection

The `develop` branch is the integration branch for features. Moderate protection with faster iteration.

#### 2.2.1 Protection Settings

**Location:** Settings → Branches → Add branch protection rule  
**Branch name pattern:** `develop`

#### 2.2.2 Required Settings

```
✅ Require a pull request before merging
   ├─ Require approvals: 1
   ├─ Dismiss stale pull request approvals when new commits are pushed: ✅
   ├─ Require review from Code Owners: ✅
   └─ Require review from at least one team: @biometrics-dev

✅ Require status checks to pass before merging
   ├─ Required status checks:
   │   ├─ ci/lint
   │   ├─ ci/typecheck
   │   └─ ci/test
   ├─ Require branches to be up to date before merging: ✅
   └─ Allow merges even if some status checks are missing: ❌

✅ Require conversation resolution before merging
   └─ All comments must be resolved: ✅

❌ Include administrators
   └─ Apply rules to repository admins: Optional

❌ Allow force pushes
   └─ Permit force pushes: Allowed for maintainers only

❌ Allow deletions
   └─ Permit branch deletions: Allowed for maintainers only
```

#### 2.2.3 GitFlow Integration

When using GitFlow workflow:

```
main (protected) ← develop (protected) ← feature/* (unprotected)
                          ↑
                    release/* (moderate protection)
                          ↑
                    hotfix/* (moderate protection)
```

### 2.3 Feature Branch Guidelines

Feature branches are short-lived development branches for individual features.

#### 2.3.1 Naming Convention

```bash
# Format: feature/{JIRA-ID}-{short-description}
feature/BIOM-123-user-authentication
feature/BIOM-456-payment-integration

# Alternative (no JIRA):
feature/user-profile-page
feature/search-optimization
```

#### 2.3.2 Branch Creation

```bash
# Always branch from develop
git checkout develop
git pull origin develop
git checkout -b feature/BIOM-123-feature-name

# Push to remote
git push -u origin feature/BIOM-123-feature-name
```

#### 2.3.3 Protection Level

**Default:** No branch protection (flexible development)

**Optional Protection for Long-Running Features:**

```
✅ Require a pull request before merging
   └─ Require approvals: 1

✅ Require status checks to pass before merging
   ├─ Required: ci/lint, ci/typecheck
   └─ Optional: ci/test
```

#### 2.3.4 Lifetime Management

- **Maximum Lifetime:** 14 days
- **After 7 Days:** Team lead notification
- **After 14 Days:** Escalation to project manager
- **Stale Branches:** Auto-close after 30 days (with warning)

### 2.4 Release Branch Strategy

Release branches are created when preparing for a production release.

#### 2.4.1 Naming Convention

```bash
# Format: release/v{MAJOR}.{MINOR}.{PATCH}
release/v1.2.0
release/v1.2.1
release/v2.0.0
```

#### 2.4.2 Protection Settings

```
✅ Require a pull request before merging
   ├─ Require approvals: 2
   ├─ Require review from Code Owners: ✅
   └─ Required reviewers: @biometrics-core, @biometrics-qa

✅ Require status checks to pass before merging
   ├─ Required status checks:
   │   ├─ ci/lint
   │   ├─ ci/typecheck
   │   ├─ ci/test
   │   ├─ ci/build
   │   ├─ e2e/test
   │   └─ security/scan
   └─ Require branches to be up to date before merging: ✅

❌ Allow force pushes
   └─ Permit force pushes: Disabled

❌ Allow deletions
   └─ Permit branch deletions: Disabled (until merged)
```

#### 2.4.3 Release Process

```
1. Create release branch from develop
   git checkout develop
   git checkout -b release/v1.2.0

2. Freeze feature additions
   └─ Only bug fixes allowed

3. Run full test suite
   └─ All tests must pass

4. Update version numbers
   └─ package.json, CHANGELOG.md, etc.

5. Create PR to main
   └─ Requires 2 approvals + all status checks

6. Merge to main and tag
   git tag -a v1.2.0 -m "Release v1.2.0"

7. Merge back to develop
   └─ Sync release changes
```

### 2.5 Hotfix Branch Strategy

Hotfix branches are for urgent production fixes.

#### 2.5.1 Naming Convention

```bash
# Format: hotfix/{JIRA-ID}-{short-description}
hotfix/BIOM-999-critical-security-fix
hotfix/production-crash-fix
```

#### 2.5.2 Protection Settings

```
✅ Require a pull request before merging
   ├─ Require approvals: 1 (expedited)
   ├─ Require review from Code Owners: ✅
   └─ Required reviewers: @biometrics-core

✅ Require status checks to pass before merging
   ├─ Required status checks:
   │   ├─ ci/lint
   │   ├─ ci/typecheck
   │   └─ ci/test (critical tests only)
   └─ Require branches to be up to date before merging: ❌ (optional for speed)

❌ Allow force pushes
   └─ Permit force pushes: Disabled

❌ Allow deletions
   └─ Permit branch deletions: Allowed after merge
```

#### 2.5.3 Emergency Hotfix Process

For CRITICAL production issues only:

```
1. Create hotfix branch from main (not develop!)
   git checkout main
   git checkout -b hotfix/critical-fix

2. Implement minimal fix
   └─ NO feature additions
   └─ NO refactoring
   └─ ONLY the critical fix

3. Create PR to main
   └─ Label: "hotfix", "critical"
   └─ Notify: @biometrics-core (Slack/Teams)

4. Expedited review (1 hour SLA)
   └─ Minimum 1 approval from core team

5. Merge and tag immediately
   git tag -a v1.2.1-hotfix -m "Hotfix v1.2.1"

6. Cherry-pick to develop
   └─ Ensure develop stays in sync
```

---

## 3. CODEOWNERS CONFIGURATION

### 3.1 File Path Coverage

The `.github/CODEOWNERS` file defines automatic reviewer assignment based on file paths.

#### 3.1.1 Coverage Matrix

| Path Pattern | Owner(s) | Review Required |
|--------------|----------|-----------------|
| `*` | @biometrics-core | Default fallback |
| `/biometrics-cli/` | @biometrics-go @go-reviewers | Always |
| `**/*.go` | @biometrics-go @go-reviewers | Always |
| `**/*.ts` | @biometrics-frontend | Always |
| `**/*.tsx` | @biometrics-frontend | Always |
| `/docs/` | @biometrics-docs | Always |
| `**/*.md` | @biometrics-docs | Optional |
| `/.github/` | @biometrics-devops | Always |
| `/assets/` | @biometrics-design | Always |
| `**/*.yaml` | @biometrics-devops | Always |
| `**/*.yml` | @biometrics-devops | Always |
| `**/*.toml` | @biometrics-devops | Always |
| `**/*.json` | @biometrics-devops | Optional |
| `**/Dockerfile` | @biometrics-devops | Always |
| `**/docker-compose.yml` | @biometrics-devops | Always |
| `/backups/` | @biometrics-devops | Always |

#### 3.1.2 Pattern Matching Rules

GitHub CODEOWNERS uses gitignore-style patterns:

```bash
# Exact path match
/docs/                    # Matches /docs directory

# Glob patterns
**/*.go                   # Matches all .go files anywhere
**/test/**                # Matches anything in test directories

# Anchor patterns (root only)
/README.md                # Only root README.md
/package.json             # Only root package.json

# Negation patterns
!*.md                     # Exclude all .md files
!/docs/important.md       # Except this specific file
```

### 3.2 Team Definitions

Teams must be created in GitHub Organization Settings.

#### 3.2.1 Required Teams

| Team Name | Purpose | Members | Maintain