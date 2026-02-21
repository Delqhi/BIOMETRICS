# Git Workflow Guide - BIOMETRICS Project

**Version:** 1.0  
**Status:** Active  
**Last Updated:** 2026-02-20  
**Project:** BIOMETRICS

---

## Table of Contents

1. [Commit Standards](#1-commit-standards)
2. [Branch Strategy](#2-branch-strategy)
3. [Pull Request Standards](#3-pull-request-standards)
4. [Git Commands](#4-git-commands)
5. [CI/CD Integration](#5-cicd-integration)
6. [Git Configuration](#6-git-configuration)
7. [Collaboration](#7-collaboration)
8. [Disaster Recovery](#8-disaster-recovery)

---

## 1. Commit Standards

### 1.1 Conventional Commits Format

The BIOMETRICS project follows the Conventional Commits specification. This standardized format enables automated changelog generation, semantic versioning, and clear commit history.

#### Format Structure

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

#### Type Categories

| Type | Description | Example |
|------|-------------|---------|
| `feat` | A new feature | `feat(auth): add biometric authentication` |
| `fix` | A bug fix | `fix(api): resolve null pointer exception` |
| `docs` | Documentation only changes | `docs(readme): update installation guide` |
| `style` | Code style changes (formatting, semicolons) | `style(utils): format code with prettier` |
| `refactor` | Code change that neither fixes nor adds | `refactor(core): simplify data processing` |
| `perf` | Performance improvement | `perf(database): add indexing for queries` |
| `test` | Adding or updating tests | `test(auth): add unit tests for login` |
| `build` | Build system or dependency changes | `build(package): upgrade dependencies` |
| `ci` | CI/CD configuration changes | `ci(github): add workflow for testing` |
| `chore` | Maintenance tasks | `chore(deps): update npm packages` |
| `revert` | Reverting a previous commit | `revert: feat(old): deprecated feature` |

#### Scope Categories

Scopes should be specific to the project structure:

- `auth` - Authentication related
- `api` - API endpoints and handlers
- `ui` - User interface components
- `core` - Core business logic
- `database` - Database models and migrations
- `config` - Configuration files
- `docs` - Documentation
- `ci` - CI/CD pipelines
- `deps` - Dependencies

#### Description Rules

- Use imperative mood: "add" not "added" or "adds"
- Start with lowercase
- No period at the end
- Maximum 72 characters
- Be specific and descriptive

**Examples:**
```
feat(auth): add fingerprint recognition support
fix(database): resolve connection pool exhaustion
refactor(api): simplify request validation
perf(utils): optimize image processing pipeline
```

### 1.2 Commit Message Rules

#### Subject Line Rules

1. **Type must be lowercase**: `feat:` not `FEAT:`
2. **Scope is optional but recommended**: Use parentheses
3. **Description**: Maximum 72 characters
4. **No trailing period**: Use `add feature` not `add feature.`
5. **Use imperative mood**: "fix" not "fixes" or "fixed"

#### Body Rules

- Use present tense: "fix" not "fixed"
- Separate subject from body with a blank line
- Wrap body at 72 characters
- Explain "what" and "why", not "how"

#### Footer Rules

- Use `BREAKING CHANGE:` for breaking changes
- Reference issues: `Closes #123`, `Fixes #456`
- Use `Co-authored-by:` for multiple authors

#### Complete Example

```
feat(auth): implement biometric login flow

Add biometric authentication using fingerprint and face recognition.
Integrate with iOS LocalAuthentication and Android Biometric APIs.

BREAKING CHANGE: Authentication middleware now requires biometric
configuration in config.yaml

Closes #123
Fixes #456
Co-authored-by: John Doe <john@example.com>
```

### 1.3 Atomic Commits

#### What are Atomic Commits?

Atomic commits are small, focused commits that represent a single logical change. Each commit should be:

- **Complete**: The code compiles and tests pass
- **Focused**: Addresses one specific concern
- **Independent**: Can be reverted without affecting other changes

#### Benefits

1. **Easier Code Review**: Reviewers can understand small changes
2. **Better Debugging**: Bisect to find problematic commits
3. **Flexible History**: Cherry-pick or revert specific changes
4. **Clear Documentation**: Tell a coherent story of changes

#### Rules for Atomic Commits

1. **One task per commit**: Don't mix features and fixes
2. **Test before commit**: Ensure code compiles and tests pass
3. **Review diff size**: If diff > 400 lines, consider splitting
4. **Logical grouping**: Group related changes together

#### Examples

**Good Atomic Commits:**
```
commit a1b2c3d4
Author: Developer <dev@biometrics.local>
Date:   Mon Feb 20 10:00:00 2026

    feat(ui): add fingerprint scanner component

    - Create FingerprintScanner component
    - Add visual feedback animation
    - Integrate with authentication service
    
commit e5f6g7h8
Author: Developer <dev@biometrics.local>
Date:   Mon Feb 20 10:15:00 2026

    fix(auth): resolve token expiration issue
    
    - Extend token validity from 1h to 24h
    - Add refresh token mechanism
    - Update validation logic
```

**Non-Atomic (Avoid):**
```
commit i9j0k1l2
Author: Developer <dev@biometrics.local>
Date:   Mon Feb 20 10:30:00 2026

    multiple changes
    
    - Added fingerprint auth
    - Fixed some bugs
    - Updated README
    - Changed database schema
    - Removed unused code
```

### 1.4 Commit Frequency

#### When to Commit

1. **After every significant change**: Don't wait until end of day
2. **Before switching branches**: Always commit or stash changes
3. **After completing a logical unit**: Even if small
4. **Before pull/rebase**: Ensure clean working state

#### Commit Frequency Guidelines

| Scenario | Minimum Frequency |
|----------|------------------|
| Active development | Every 1-2 hours |
| Bug fix | After each fix |
| Documentation | After each section |
| After testing | Before and after |

#### Best Practices

1. **Commit early, commit often**: Small commits are better
2. **Don't break the build**: Always test before committing
3. **Write meaningful messages**: Future you will thank present you
4. **Review before commit**: Use `git diff` to verify changes

---

## 2. Branch Strategy

### 2.1 Main Branch Protection

#### Protected Branches

The BIOMETRICS project protects the following branches:

| Branch | Purpose | Protection Rules |
|--------|---------|-----------------|
| `main` | Production code | Require PR, require reviews, require status checks |
| `develop` | Integration branch | Require PR, require reviews |

#### Protection Rules Configuration

```yaml
# .github/branch-protection.yml
main:
  required_reviewers: 2
  require_status_checks: true
  require_up_to_date: true
  allow_force_push: false
  allow_deletion: false
  
develop:
  required_reviewers: 1
  require_status_checks: true
  require_up_to_date: true
  allow_force_push: false
  allow_deletion: false
```

#### Main Branch Guidelines

1. **Never commit directly**: All changes via PR
2. **Require tests passing**: CI must pass before merge
3. **Require approvals**: At least 2 approvals for main
4. **Keep it clean**: Only production-ready code

### 2.2 Feature Branches

#### Feature Branch Lifecycle

```
create branch → develop → test → pull request → review → merge → delete
```

#### Naming Convention

Format: `feature/<ticket-id>-<short-description>`

Examples:
```
feature/BIOM-123-fingerprint-auth
feature/BIOM-456-face-recognition
feature/BIOM-789-voice-authentication
```

#### Feature Branch Workflow

```bash
# 1. Create from develop
git checkout develop
git pull origin develop
git checkout -b feature/BIOM-123-new-feature

# 2. Develop and commit
git add .
git commit -m "feat: add initial feature implementation"

# 3. Keep branch updated
git fetch origin
git rebase origin/develop

# 4. Push and create PR
git push -u origin feature/BIOM-123-new-feature
```

### 2.3 Release Branches

#### When to Create Release Branches

- When `develop` has all required features
- When preparing for production deployment
- When code freeze for testing begins

#### Naming Convention

Format: `release/<version>`

Examples:
```
release/2.0.0
release/2.1.0
release/2.1.1
```

#### Release Branch Workflow

```bash
# Create release branch
git checkout develop
git pull origin develop
git checkout -b release/2.0.0

# Make release-specific changes
# Update version numbers
# Fix last-minute bugs

# Merge to main
git checkout main
git merge release/2.0.0
git tag -a v2.0.0 -m "Release version 2.0.0"
git push origin main --tags

# Merge back to develop
git checkout develop
git merge release/2.0.0
git push origin develop

# Delete release branch
git branch -d release/2.0.0
git push origin --delete release/2.0.0
```

### 2.4 Hotfix Branches

#### When to Use Hotfix Branches

- Critical production bugs requiring immediate fix
- Security vulnerabilities
- Service outages

#### Naming Convention

Format: `hotfix/<ticket-id>-<short-description>`

Examples:
```
hotfix/BIOM-999-security-patch
hotfix/BIOM-1000-login-fix
```

#### Hotfix Workflow

```bash
# Create hotfix from main
git checkout main
git pull origin main
git checkout -b hotfix/BIOM-999-fix

# Fix the issue
git add .
git commit -m "fix: resolve critical security vulnerability"

# Merge to main
git checkout main
git merge hotfix/BIOM-999-fix
git tag -a v1.0.1 -m "Hotfix version 1.0.1"
git push origin main --tags

# Merge to develop
git checkout develop
git merge hotfix/BIOM-999-fix
git push origin develop

# Delete hotfix branch
git branch -d hotfix/BIOM-999-fix
```

### 2.5 Branch Naming Conventions

#### Complete Naming Guide

| Type | Format | Example |
|------|--------|---------|
| Feature | `feature/<ticket>-<description>` | `feature/BIOM-123-user-auth` |
| Bugfix | `bugfix/<ticket>-<description>` | `bugfix/BIOM-456-login-error` |
| Hotfix | `hotfix/<ticket>-<description>` | `hotfix/BIOM-789-security` |
| Release | `release/<version>` | `release/2.0.0` |
| Experiment | `experiment/<description>` | `experiment/new-algorithm` |
| Documentation | `docs/<description>` | `docs/api-reference` |

#### Branch Naming Rules

1. **Use lowercase**: `feature/auth` not `Feature/Auth`
2. **Use hyphens**: Separate words with `-`
3. **Include ticket ID**: Reference project management
4. **Be descriptive**: Name reflects the branch purpose
5. **Keep it short**: Maximum 50 characters

#### Git Alias for Branch Creation

```bash
# Add to ~/.gitconfig
[alias]
  fb = "!f() { git checkout -b \"feature/BIOM-$1-${2:-feature}\"; }; f"
  bb = "!f() { git checkout -b \"bugfix/BIOM-$1-${2:-fix}\"; }; f"
  hb = "!f() { git checkout -b \"hotfix/BIOM-$1-${2:-hotfix}\"; }; f"

# Usage
git fb 123 fingerprint-auth
git bb 456 login-error
git hb 789 security-patch
```

---

## 3. Pull Request Standards

### 3.1 PR Template

#### Template File: `.github/PULL_REQUEST_TEMPLATE.md`

```markdown
## Description
<!-- Describe your changes in detail -->

## Related Issue
<!-- Link to related issue: Fixes #123 -->

## Type of Change
- [ ] 🐛 Bug fix (non-breaking change)
- [ ] ✨ New feature (non-breaking change)
- [ ] 💥 Breaking change (migration required)
- [ ] 📝 Documentation update
- [ ] 🔧 Configuration change
- [ ] ♻️ Refactoring (no functional changes)

## Testing Performed
<!-- Describe testing done -->
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed
- [ ] Performance testing

## Screenshots
<!-- Add screenshots if applicable -->

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Code commented where necessary
- [ ] Documentation updated
- [ ] No new warnings
- [ ] Tests added/updated
- [ ] All tests pass

## Reviewers Required
- [ ] @reviewer1
- [ ] @reviewer2
```

### 3.2 PR Checklist

#### Before Creating PR

- [ ] All tests pass locally
- [ ] Code follows style guidelines
- [ ] No console.log or debug code
- [ ] No sensitive data in commits
- [ ] Commit messages are descriptive
- [ ] Branch is up-to-date with target
- [ ] Documentation updated if needed

#### Before Requesting Review

- [ ] PR description complete
- [ ] Related issue linked
- [ ] Screenshots added if UI changes
- [ ] Breaking changes documented
- [ ] Migration guide prepared (if breaking)

#### Before Merging

- [ ] All reviews addressed
- [ ] All CI checks pass
- [ ] No merge conflicts
- [ ] At least 2 approvals (for main)
- [ ] Branch up-to-date with target

### 3.3 Code Review Requirements

#### Review Checklist

**Code Quality:**
- [ ] Code is readable and well-structured
- [ ] No duplicate code
- [ ] Functions are small and focused
- [ ] Error handling is proper
- [ ] No hardcoded values

**Security:**
- [ ] No security vulnerabilities
- [ ] Sensitive data protected
- [ ] Input validation present
- [ ] Authentication/authorization correct

**Performance:**
- [ ] No memory leaks
- [ ] Database queries optimized
- [ ] No unnecessary computations

**Testing:**
- [ ] Tests are meaningful
- [ ] Edge cases covered
- [ ] Test coverage maintained

#### Review Comments Guidelines

**Good Comments:**
```
// Good: Specific and actionable
// Consider using a hash map here for O(1) lookup instead of O(n)
```

**Bad Comments:**
```
// Bad: Vague and not helpful
// This could be better
```

#### Review Response Time

| Priority | Response Time |
|----------|---------------|
| Hotfix/Urgent | Within 1 hour |
| Regular | Within 24 hours |
| Documentation | Within 48 hours |

### 3.4 Merge Strategies

#### Available Merge Strategies

| Strategy | When to Use | Command |
|----------|-------------|---------|
| Squash Merge | Feature branches | `git merge --squash` |
| Rebase | Linear history | `git rebase` |
| Merge Commit | Collaborative branches | `git merge` |

#### Squash Merge (Recommended for Features)

```bash
# Before merging
git checkout main
git pull origin main
git merge --squash feature/BIOM-123

# Commit with descriptive message
git commit -m "feat(auth): add biometric authentication

- Implement fingerprint scanner component
- Add face recognition support
- Integrate with authentication service

Closes #123"
```

#### Rebase (For Clean History)

```bash
# Update feature branch with latest develop
git checkout feature/BIOM-123
git fetch origin
git rebase origin/develop

# Force push (only for feature branches!)
git push --force-with-lease
```

#### Merge Commit (For Large Features)

```bash
# Standard merge
git checkout main
git pull origin main
git merge feature/BIOM-123
```

#### Merge Strategy Rules

1. **Feature branches**: Use squash merge
2. **Release branches**: Use merge commit
3. **Hotfix branches**: Use squash merge
4. **Never rebase**: Main and develop branches

---

## 4. Git Commands

### 4.1 Daily Workflow Commands

#### Starting New Work

```bash
# Update develop
git checkout develop
git pull origin develop

# Create new feature branch
git checkout -b feature/BIOM-123-description
```

#### Saving Work in Progress

```bash
# Stage specific files
git add file1.txt file2.txt

# Stage all changes
git add -A

# Stage by pattern
git add "*.js"
git add src/

# Check status
git status
git status -s  # Short format
```

#### Committing Changes

```bash
# Commit with message
git commit -m "feat: add new feature"

# Commit with extended message
git commit -m "feat: add new feature" -m "Detailed description..."

# Amend last commit
git commit --amend -m "Updated message"

# Amend without changing message
git commit --amend --no-edit
```

#### Syncing with Remote

```bash
# Push to remote
git push origin branch-name

# Push with upstream tracking
git push -u origin branch-name

# Pull with rebase (recommended)
git pull --rebase origin develop

# Fetch updates
git fetch origin
```

### 4.2 Advanced Commands

#### Rebase Operations

```bash
# Interactive rebase (last 3 commits)
git rebase -i HEAD~3

# Rebase onto develop
git rebase develop

# Continue after resolving conflicts
git rebase --continue

# Abort rebase
git rebase --abort

# Skip current commit
git rebase --skip
```

**Interactive Rebase Commands:**
```
pick   - use commit
reword - change commit message
edit   - stop to amend
squash - combine with previous
drop   - remove commit
```

#### Cherry-Pick Operations

```bash
# Cherry-pick a single commit
git cherry-pick abc1234

# Cherry-pick without committing
git cherry-pick -n abc1234

# Cherry-pick a range
git cherry-pick abc1234..def5678

# Continue after resolving conflicts
git cherry-pick --continue
```

#### Stash Operations

```bash
# Save work in progress
git stash
git stash save "Work in progress"

# List stashes
git stash list

# Apply latest stash
git stash pop

# Apply specific stash
git stash pop stash@{2}

# Drop stash
git stash drop stash@{0}

# Clear all stashes
git stash clear
```

### 4.3 Undo Operations

#### Undo Staged Changes

```bash
# Unstage a file
git reset HEAD file.txt

# Unstage all
git reset HEAD
```

#### Undo Uncommitted Changes

```bash
# Discard changes in working directory
git checkout -- file.txt
git checkout -- .  # All files

# Or use restore (Git 2.23+)
git restore file.txt
git restore .
```

#### Undo Committed Changes

```bash
# Undo last commit (keep changes staged)
git reset --soft HEAD~1

# Undo last commit (keep changes unstaged)
git reset HEAD~1

# Undo last commit (discard changes)
git reset --hard HEAD~1
```

#### Revert Committed Changes

```bash
# Create new commit that undoes changes
git revert abc1234

# Revert without auto-commit
git revert -n abc1234
```

### 4.4 Cleanup Commands

#### Remove Untracked Files

```bash
# Preview what will be removed
git clean -n

# Remove untracked files
git clean -f

# Remove untracked files and directories
git clean -fd

# Remove ignored files too
git clean -fdx
```

#### Remove Merged Branches

```bash
# Delete local branch
git branch -d branch-name

# Force delete local branch
git branch -D branch-name

# Delete remote branch
git push origin --delete branch-name

# Prune deleted remote branches
git fetch --prune
```

#### Remove Old Tags

```bash
# Delete local tag
git tag -d v1.0.0

# Delete remote tag
git push origin --delete v1.0.0

# Delete tags matching pattern
git tag -l "v1.*" | xargs git tag -d
```

---

## 5. CI/CD Integration

### 5.1 GitHub Actions Workflows

#### CI Workflow: `.github/workflows/ci.yml`

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
      
      - name: Install dependencies
        run: npm ci
      
      - name: Run linter
        run: pnpm run lint
      
      - name: Run type check
        run: pnpm run typecheck

  test:
    runs-on: ubuntu-latest
    needs: lint
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
      
      - name: Install dependencies
        run: npm ci
      
      - name: Run tests
        run: npm test -- --coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          token: ${{ secrets.CODECOV_TOKEN }}

  build:
    runs-on: ubuntu-latest
    needs: test
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
      
      - name: Install dependencies
        run: npm ci
      
      - name: Build
        run: pnpm run build
      
      - name: Upload artifacts
        uses: actions/upload-artifact@v4
        with:
          name: build
          path: dist/
```

#### Release Workflow: `.github/workflows/release.yml`

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          registry-url: 'https://registry.npmjs.org'
      
      - name: Install dependencies
        run: npm ci
      
      - name: Build
        run: pnpm run build
      
      - name: Publish to npm
        run: npm publish
        env:
          NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}
```

### 5.2 Pre-commit Hooks

#### Setup

```bash
# Install Husky
pnpm install husky --save-dev

# Initialize Husky
npx husky install

# Add to package.json
npm pkg set scripts.prepare="husky install"
```

#### Pre-commit Hook Example

Create `.husky/pre-commit`:

```bash
#!/usr/bin/env sh
. "$(dirname -- "$0")/_/husky.sh"

# Run linter on staged files
npx lint-staged

# Run type check
pnpm run typecheck

# Run tests on changed files
npm test -- --changedSince=HEAD
```

#### lint-staged Configuration

```json
{
  "lint-staged": {
    "*.{js,ts,jsx,tsx}": [
      "eslint --fix",
      "prettier --write"
    ],
    "*.{json,md,yml,yaml}": [
      "prettier --write"
    ]
  }
}
```

### 5.3 Post-commit Hooks

#### Useful Post-commit Hooks

Create `.husky/post-commit`:

```bash
#!/usr/bin/env sh
. "$(dirname -- "$0")/_/husky.sh"

# Notify team of new commit
# (Optional: integrate with Slack/Discord)

# Run post-commit tasks
echo "Commit completed: $(git rev-parse HEAD)"
```

### 5.4 Automated Testing

#### Test Commands in CI

```bash
# Unit tests
npm test

# Integration tests
pnpm run test:integration

# E2E tests
pnpm run test:e2e

# Test with coverage
npm test -- --coverage

# Watch mode
npm test -- --watch
```

#### Test Reporting

```yaml
# Add to CI workflow
- name: Generate test report
  if: always()
  run: npm test -- --json > test-results.json
  
- name: Upload test results
  if: always()
  uses: actions/upload-artifact@v4
  with:
    name: test-results
    path: test-results.json
```

---

## 6. Git Configuration

### 6.1 .gitignore Best Practices

#### Essential .gitignore Structure

```
# Dependencies
node_modules/
vendor/

# Build outputs
dist/
build/
*.egg-info/

# Environment files
.env
.env.local
.env.*.local

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Logs
*.log
npm-debug.log*

# Testing
coverage/
.nyc_output/

# Secrets
*.pem
*.key
credentials.json
```

#### Project-Specific .gitignore

```
# BIOMETRICS specific
data/raw/*
data/processed/*
!data/.gitkeep
models/*.pth
models/*.h5
*.onnx
```

### 6.2 Global Git Config

#### Setup

```bash
# Set global config
git config --global init.defaultBranch main
git config --global pull.rebase true
git config --global fetch.prune true
git config --global core.autocrlf input
```

#### Complete Global Config

```ini
[user]
    name = Your Name
    email = your.email@biometrics.local

[core]
    editor = code --wait
    autocrlf = input
    safecrlf = warn

[pull]
    rebase = true
    ff = only

[fetch]
    prune = true
    pruneTags = true

[push]
    default = simple
    followTags = true

[init]
    defaultBranch = main

[merge]
    tool = vscode
    conflictstyle = diff3

[diff]
    tool = vscode
    colorMoved = default

[commit]
    template = ~/.gitmessage
```

### 6.3 Git Aliases

#### Essential Aliases

```ini
[alias]
    # Status
    s = status -sb
    st = status
    
    # Branch
    b = branch
    ba = branch -a
    bd = branch -d
    bD = branch -D
    
    # Log
    l = log --oneline -20
    lg = log --oneline --graph --all
    ll = log --oneline -10 --graph
    
    # Commit
    c = commit
    ca = commit --amend
    can = commit --amend --no-edit
    
    # Diff
    d = diff
    dc = diff --cached
    ds = diff --stat
    
    # Reset
    u = reset HEAD
    unstage = reset HEAD --
    
    # Stash
    ss = stash save
    sl = stash list
    sp = stash pop
    
    # Checkout
    co = checkout
    cob = checkout -b
    
    # Rebase
    ri = rebase -i
    rc = rebase --continue
    ra = rebase --abort
    
    # Other
    fp = fetch --prune
    tags = tag -l
    branches = branch -a
```

### 6.4 Git Hooks

#### Available Hooks

| Hook | When | Use Case |
|------|------|----------|
| `pre-commit` | Before commit | Linting, formatting |
| `prepare-commit-msg` | Before message editor | Auto-fill template |
| `commit-msg` | After message | Validate format |
| `post-commit` | After commit | Notifications |
| `pre-push` | Before push | Run tests |
| `post-checkout` | After checkout | Cleanup |

#### Example: Commit Message Validation

Create `.husky/commit-msg`:

```bash
#!/usr/bin/env sh
COMMIT_MSG=$(cat "$1")
PATTERN="^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\(.+\))?: .{1,72}"

if ! echo "$COMMIT_MSG" | grep -qE "$PATTERN"; then
    echo "Invalid commit message format."
    echo "Expected: <type>(<scope>): <description>"
    echo "Types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert"
    exit 1
fi
```

---

## 7. Collaboration

### 7.1 Code Review Process

#### Review Workflow

```
Author creates PR → Reviewers notified → Reviewers examine code 
→ Author addresses feedback → Reviewers approve → PR merged
```

#### Review Guidelines for Reviewers

1. **Be timely**: Review within 24 hours
2. **Be constructive**: Provide actionable feedback
3. **Be specific**: Point to exact location
4. **Be thorough**: Check all aspects
5. **Be respectful**: Professional tone always

#### Review Guidelines for Authors

1. **Keep PRs small**: Under 400 lines preferred
2. **Provide context**: Explain what and why
3. **Respond promptly**: Address feedback quickly
4. **Be receptive**: Consider all suggestions
5. **Ask questions**: Clarify if needed

### 7.2 Conflict Resolution

#### Merge Conflict Resolution

```bash
# When conflict occurs during merge/rebase
git status

# View conflicting files
git diff --name-only --diff-filter=U

# Open in editor and resolve
# Remove conflict markers
# Keep desired changes

# Mark as resolved
git add resolved-file.txt

# Continue operation
git rebase --continue
# or
git merge --continue
```

#### Conflict Resolution Strategies

1. **Understand both sides**: Read the conflicting code
2. **Talk to author**: Discuss the best solution
3. **Take the best**: Combine the best parts
4. **Test thoroughly**: Verify resolution works

### 7.3 Pair Programming with Git

#### Workflow

```bash
# One person hosts the session
git branch session/pairing-$USER

# Share branch name with pair
# Both work on same branch

# Commit with both authors
git commit --author="Partner <partner@email>" -m "feat: pair programming session"

# Merge when done
```

#### Co-authoring Commits

```bash
# Add co-author to commit
git commit -m "feat: feature implementation

Co-authored-by: Partner Name <partner@email>"
```

### 7.4 Git Etiquette

#### Do's

- ✅ Commit often with clear messages
- ✅ Keep branches up-to-date
- ✅ Review PRs promptly
- ✅ Test before pushing
- ✅ Communicate about blockers

#### Don'ts

- ❌ Push broken code
- ❌ Merge without review
- ❌ Force push to main/develop
- ❌ Leave commented-out code
- ❌ Commit secrets or sensitive data

---

## 8. Disaster Recovery

### 8.1 Undoing Commits

#### Soft Reset (Keep Changes)

```bash
# Undo last commit, keep changes staged
git reset --soft HEAD~1

# Undo last 3 commits
git reset --soft HEAD~3
```

#### Mixed Reset (Keep Working Directory)

```bash
# Undo last commit, keep changes unstaged
git reset HEAD~1

# Undo specific commit
git reset abc1234
```

#### Hard Reset (Discard Everything)

```bash
# Danger: This deletes all changes!
git reset --hard HEAD~1

# Hard reset to specific commit
git reset --hard abc1234
```

### 8.2 Recovering Lost Work

#### Using Reflog

```bash
# View reflog
git reflog
git reflog show HEAD

# Find lost commit
git reflog | grep "commit message"

# Recover branch
git checkout -b recovery-branch abc1234

# Or recover directly
git checkout abc1234
```

#### Recovering Deleted Branch

```bash
# Find the commit
git reflog

# Recreate branch
git checkout -b recovered-branch abc1234
```

#### Recovering Lost Stash

```bash
# Find stash in reflog
git fsck --unreachable | grep commit

# Recover stash
git stash apply abc1234
```

### 8.3 Reset vs Revert

#### When to Use Reset

- Undo unmerged commits
- Local commits not pushed
- Want to rewrite history

```bash
# Local commits only
git reset --hard HEAD~1
```

#### When to Use Revert

- Already pushed commits
- Collaborative branches
- Preserve history

```bash
# Create new commit that undoes
git revert abc1234
```

### 8.4 Reflog Usage

#### Common Reflog Operations

```bash
# View HEAD reflog
git reflog HEAD

# View branch reflog
git reflog show develop

# Time-based recovery
git checkout HEAD@{2.hours.ago}

# Find specific state
git reflog | grep "checkout: moving to"

# Cleanup reflog (older than 90 days)
git reflog expire --expire=90.days.all
```

---

## Appendix A: Quick Reference

### Common Commands

```bash
# Setup
git config --global user.name "Name"
git config --global user.email "email"
git clone url
git init

# Daily Work
git status
git add .
git commit -m "message"
git push
git pull --rebase

# Branching
git branch
git checkout -b branch-name
git checkout branch-name

# History
git log --oneline
git log --graph
git diff

# Undo
git reset --soft HEAD~
git revert commit
git stash
```

### Alias Reference

```bash
git s          # Status
git l          # Log
git c          # Commit
git d          # Diff
git co         # Checkout
git cob        # Create branch
git u          # Unstage
git fp         # Fetch prune
```

---

## Appendix B: Resources

- [Conventional Commits](https://www.conventionalcommits.org/)
- [GitHub Flow](https://guides.github.com/introduction/flow/)
- [Pro Git Book](https://git-scm.com/book)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)

---

**Document Version:** 1.0  
**Last Updated:** 2026-02-20  
**Maintained By:** BIOMETRICS Development Team  
**Review Cycle:** Quarterly
