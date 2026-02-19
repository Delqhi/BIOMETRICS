# CHANGELOG Configuration

This document defines the guidelines and configuration for maintaining the BIOMETRICS CHANGELOG.md.

---

## Overview

The CHANGELOG.md is auto-generated from Git history and follows industry best practices.

**Standards:**
- [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
- [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
- [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)

---

## Git Commands for CHANGELOG Generation

### Extract Full Commit History

```bash
git log --pretty=format:"%H|%h|%an|%ae|%ad|%s|%b" --date=iso-strict
```

### Extract Commit Statistics

```bash
# Contributor stats
git shortlog -sn --all

# Commit activity by date
git log --all --format="%ad" --date=short | sort | uniq -c | sort -rn

# Commit type distribution
git log --all --grep="^feat:" --oneline | wc -l
git log --all --grep="^fix:" --oneline | wc -l
git log --all --grep="^docs:" --oneline | wc -l
git log --all --grep="^chore:" --oneline | wc -l
```

### Generate Release Notes

```bash
# Between two tags
git log v1.0.0..v1.1.0 --pretty=format:"* %s by @%an in !%h"

# With file stats
git log v1.0.0..v1.1.0 --stat
```

---

## CHANGELOG Structure

### Required Sections

1. **Header**
   - Project Information
   - Repository URL
   - License
   - Best Practices compliance

2. **Contributors**
   - All contributors with commit counts
   - Percentage breakdown
   - Active development period

3. **[Unreleased]**
   - Changes not yet in a release
   - Categorized (Added, Changed, Fixed, etc.)

4. **[VERSION] - DATE**
   - Semantic version number
   - ISO 8601 date format (YYYY-MM-DD)
   - All changes with author attribution

5. **Git Statistics**
   - Commit activity
   - Commit types
   - Top changed files
   - Documentation coverage

6. **Breaking Changes**
   - Clearly marked breaking changes
   - Migration guide

7. **Known Issues**
   - Any known issues in the release

8. **Upgrade Guide**
   - Step-by-step upgrade instructions

9. **References**
   - Git commit SHA256 references
   - Related issues
   - External links

---

## Commit Categorization

### Conventional Commit Types

| Type | Section | Description |
|------|---------|-------------|
| `feat:` | Added | New features |
| `fix:` | Fixed | Bug fixes |
| `docs:` | Changed | Documentation changes |
| `style:` | Changed | Code style (formatting) |
| `refactor:` | Changed | Code refactoring |
| `test:` | Changed | Adding/updating tests |
| `chore:` | Changed | Maintenance tasks |
| `perf:` | Changed | Performance improvements |
| `security:` | Security | Security-related changes |
| `deps:` | Changed | Dependency updates |

### Manual Categorization

Some commits require manual categorization:

```bash
# Port Sovereignty fixes
git log --grep="Port Sovereignty" --oneline
# → Categorize as "Fixed"

# Qwen 3.5 integration
git log --grep="Qwen 3.5" --oneline
# → Categorize as "Added" or "Changed"

# Documentation expansion
git log --grep="5000+ lines" --oneline
# → Categorize as "Changed"
```

---

## Version Numbering

### Semantic Versioning (SemVer)

Format: `MAJOR.MINOR.PATCH`

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### BIOMETRICS Versioning

```
v0.9.0  - Initial beta release
v1.0.0  - First stable release
v1.1.0  - Minor features
v1.1.1  - Patch fixes
v2.0.0  - Major breaking changes
```

---

## Release Process

### 1. Prepare Release

```bash
# Update CHANGELOG.md
# Move [Unreleased] to [VERSION] - DATE
# Add new [Unreleased] section

# Commit CHANGELOG update
git add CHANGELOG.md
git commit -m "docs: update CHANGELOG.md for v1.0.0"
```

### 2. Create Git Tag

```bash
# Annotated tag
git tag -a v1.0.0 -m "Release v1.0.0 - Initial stable release"

# Push tag
git push origin v1.0.0
```

### 3. Create GitHub Release

1. Go to GitHub Releases
2. Click "Create a new release"
3. Select tag: `v1.0.0`
4. Copy CHANGELOG.md content
5. Click "Publish release"

### 4. Post-Release

```bash
# Update version in package.json / VERSION file
# Push to main branch
git push origin main
```

---

## Automation Scripts

### Generate CHANGELOG from Git

```bash
#!/bin/bash
# generate-changelog.sh

# Get all commits since last tag
LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")

if [ -z "$LAST_TAG" ]; then
    echo "No tags found. Generating from initial commit."
    COMMITS=$(git log --pretty=format:"%H|%h|%an|%ae|%ad|%s|%b" --date=iso-strict)
else
    echo "Generating changelog since $LAST_TAG"
    COMMITS=$(git log $LAST_TAG..HEAD --pretty=format:"%H|%h|%an|%ae|%ad|%s|%b" --date=iso-strict)
fi

# Process commits
echo "$COMMITS" | while IFS='|' read -r hash short_hash author email date subject body; do
    # Categorize by commit type
    if [[ $subject =~ ^feat: ]]; then
        echo "### Added"
        echo "- $subject by @$author in !$short_hash"
    elif [[ $subject =~ ^fix: ]]; then
        echo "### Fixed"
        echo "- $subject by @$author in !$short_hash"
    elif [[ $subject =~ ^docs: ]]; then
        echo "### Changed"
        echo "- $subject by @$author in !$short_hash"
    fi
done
```

### Commit Statistics

```bash
#!/bin/bash
# commit-stats.sh

echo "=== Contributor Statistics ==="
git shortlog -sn --all

echo -e "\n=== Commit Activity by Date ==="
git log --all --format="%ad" --date=short | sort | uniq -c | sort -rn | head -10

echo -e "\n=== Commit Type Distribution ==="
echo "feat: $(git log --all --grep='^feat:' --oneline | wc -l)"
echo "fix: $(git log --all --grep='^fix:' --oneline | wc -l)"
echo "docs: $(git log --all --grep='^docs:' --oneline | wc -l)"
echo "chore: $(git log --all --grep='^chore:' --oneline | wc -l)"
```

---

## Best Practices

### DO ✅

- Write clear, concise commit messages
- Use Conventional Commits format
- Categorize changes correctly
- Include author attribution
- Reference Git commit SHA256
- Link to related issues/PRs
- Document breaking changes prominently
- Include upgrade guide
- Update CHANGELOG before release
- Keep [Unreleased] section current

### DON'T ❌

- Don't merge multiple unrelated changes in one commit
- Don't use vague commit messages ("fix stuff", "update code")
- Don't omit author attribution
- Don't forget to categorize security fixes
- Don't hide breaking changes
- Don't release without CHANGELOG update
- Don't use future dates for releases
- Don't exceed 100 characters per commit message
- Don't include commit message bodies in CHANGELOG (summarize instead)

---

## Examples

### Good Commit Message

```
feat: Add Qwen 3.5 integration with NVIDIA NIM

- Added Qwen 3.5 model configuration
- Added 5 skills (vision, code, OCR, video, conversation)
- Added fallback chain (Qwen → Kimi → Claude)
- Updated ARCHITECTURE.md with Qwen 3.5 Brain

Refs: QWEN-3.5, NVIDIA-NIM
```

### Bad Commit Message

```
update stuff
```

### Good CHANGELOG Entry

```markdown
### Added
- Qwen 3.5 integration with NVIDIA NIM by @Jeremy in !QWEN-3.5
  - 5 skills: vision, code, OCR, video, conversation
  - Fallback chain: Qwen → Kimi → Claude
  - 262K context window, 32K output
```

### Bad CHANGELOG Entry

```markdown
### Added
- stuff
```

---

## Tools

### Recommended Tools

1. **git-changelog** - Auto-generate CHANGELOG from Git
   ```bash
   npm install -g git-changelog
   git-changelog -o CHANGELOG.md
   ```

2. **conventional-changelog** - Parse Conventional Commits
   ```bash
   npm install -g conventional-changelog-cli
   conventional-changelog -p angular -i CHANGELOG.md -s
   ```

3. **github-changelog-generator** - Generate from GitHub releases
   ```bash
   gem install github_changelog_generator
   github_changelog_generator -u BIOMETRICS -p BIOMETRICS
   ```

### BIOMETRICS Custom Script

```bash
#!/bin/bash
# biometrics-changelog.sh

OUTPUT="CHANGELOG.md"
REPO="BIOMETRICS/BIOMETRICS"

echo "# Changelog" > $OUTPUT
echo "" >> $OUTPUT
echo "All notable changes to this project will be documented in this file." >> $OUTPUT
echo "" >> $OUTPUT

# Get latest tag
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null)

if [ -z "$LATEST_TAG" ]; then
    RANGE="HEAD"
else
    RANGE="$LATEST_TAG..HEAD"
fi

# Process commits
git log $RANGE --pretty=format:"%h|%an|%ad|%s" --date=short | while IFS='|' read hash author date message; do
    # Extract type
    type=$(echo $message | grep -oE "^(feat|fix|docs|style|refactor|test|chore)" || echo "other")
    
    # Add to appropriate section
    case $type in
        feat)
            echo "### Added" >> $OUTPUT
            echo "- $message by @$author in !$hash" >> $OUTPUT
            ;;
        fix)
            echo "### Fixed" >> $OUTPUT
            echo "- $message by @$author in !$hash" >> $OUTPUT
            ;;
        docs)
            echo "### Changed" >> $OUTPUT
            echo "- $message by @$author in !$hash" >> $OUTPUT
            ;;
        *)
            echo "### Changed" >> $OUTPUT
            echo "- $message by @$author in !$hash" >> $OUTPUT
            ;;
    esac
done
```

---

## Maintenance

### Weekly Tasks

- [ ] Review [Unreleased] section
- [ ] Categorize new commits
- [ ] Update contributor stats
- [ ] Check for missing attributions

### Monthly Tasks

- [ ] Generate release notes
- [ ] Create Git tag if ready
- [ ] Publish GitHub release
- [ ] Archive old releases
- [ ] Review and cleanup

### Quarterly Tasks

- [ ] Audit CHANGELOG quality
- [ ] Update this configuration
- [ ] Review automation scripts
- [ ] Contributor recognition update

---

## References

- [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
- [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
- [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github)
- [Git Log Documentation](https://git-scm.com/docs/git-log)

---

**Last Updated**: 2026-02-19  
**Maintainer**: Jeremy (@jeremy)  
**Version**: 1.0.0
