# Contributing to BIOMETRICS CLI

Thank you for your interest in contributing to BIOMETRICS CLI! This document provides comprehensive guidelines for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Prerequisites](#prerequisites)
- [Running Tests](#running-tests)
- [Code Style Guidelines](#code-style-guidelines)
- [Commit Message Format](#commit-message-format)
- [Branch Naming Convention](#branch-naming-convention)
- [Submitting Pull Requests](#submitting-pull-requests)
- [Development Workflow](#development-workflow)
- [Testing Requirements](#testing-requirements)
- [Documentation Standards](#documentation-standards)
- [Questions?](#questions)

---

## Code of Conduct

This project adheres to a strict code of conduct. By participating, you are expected to uphold this code:

- **Be respectful**: Treat all contributors with respect and professionalism
- **Be inclusive**: Welcome diverse perspectives and backgrounds
- **Be constructive**: Provide helpful, actionable feedback
- **Be collaborative**: Work together to achieve common goals
- **Be accountable**: Take responsibility for your contributions

---

## Getting Started

### 1. Fork the Repository

```bash
# Click "Fork" on GitHub to create your fork
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/BIOMETRICS.git
cd BIOMETRICS/biometrics-cli
```

### 2. Add Upstream Remote

```bash
# Add the original repository as upstream
git remote add upstream https://github.com/Delqhi/BIOMETRICS.git

# Verify remotes
git remote -v
```

### 3. Create a Branch

```bash
# Always create a new branch for your work
git checkout -b feature/your-feature-name
```

---

## Development Setup

### Prerequisites

Ensure you have the following installed:

#### Required Tools

- **Go**: Version 1.21 or higher
  ```bash
  go version
  # Should output: go version go1.21.x or higher
  ```

- **Docker**: Latest stable version
  ```bash
  docker --version
  docker-compose --version
  ```

- **Git**: Latest version
  ```bash
  git --version
  ```

- **Make**: For running build tasks
  ```bash
  make --version
  ```

#### Development Dependencies

- **golangci-lint**: For code linting
  ```bash
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```

- **gotestsum**: For enhanced test output (optional)
  ```bash
  go install gotest.tools/gotestsum@latest
  ```

### Local Environment Setup

#### 1. Install Go Dependencies

```bash
cd biometrics-cli
go mod download
```

#### 2. Start Development Services

```bash
# Start Redis and PostgreSQL
docker-compose up -d redis postgres

# Verify services are running
docker-compose ps

# Check service health
docker-compose exec redis redis-cli ping
docker-compose exec postgres pg_isready -U biometrics -d biometrics
```

#### 3. Verify Setup

```bash
# Run tests to verify setup
go test ./...

# Run linters
golangci-lint run

# Build the CLI
go build -o bin/biometrics ./cmd/biometrics
```

---

## Running Tests

### Run All Tests

```bash
# Run all tests with verbose output
go test -v ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out -o coverage.html
```

### Run Specific Tests

```bash
# Run tests in a specific package
go test ./pkg/auth/...

# Run a specific test function
go test -run TestOAuth2Flow ./pkg/auth/oauth2/

# Run tests matching a pattern
go test -run "Test.*Validation" ./pkg/validation/
```

### Run Benchmarks

```bash
# Run all benchmarks
go test -bench=. -benchmem ./...

# Run specific benchmark
go test -bench=BenchmarkAuditLog ./pkg/audit/
```

### Race Detection

```bash
# Run tests with race detector
go test -race ./...
```

---

## Code Style Guidelines

### General Principles

1. **Follow Go Best Practices**: Adhere to [Effective Go](https://golang.org/doc/effective_go)
2. **Keep it Simple**: Write clear, readable code
3. **Be Consistent**: Match existing code style
4. **Document Public APIs**: Use godoc comments for exported symbols
5. **Handle Errors Properly**: Never ignore errors

### Naming Conventions

#### Variables and Functions

```go
// Use camelCase for variables and functions
var userName string
func calculateTotal() int

// Use underscores for test functions
func TestOAuth2_Flow_ValidToken(t *testing.T)

// Use acronyms in all caps
var apiKey string
func GetHTTPClient() *http.Client
```

#### Types and Interfaces

```go
// Use PascalCase for exported types
type UserProfile struct{}
type DataProcessor interface{}

// Use lowercase for unexported types
type internalConfig struct{}

// Interface names should be descriptive
type Reader interface{}
type Writer interface{}
type ReadWriter interface{}
```

#### Constants

```go
// Use PascalCase for exported constants
const MaxRetries = 3
const DefaultTimeout = 30 * time.Second

// Use camelCase for unexported constants
const maxBufferSize = 1024
```

### Code Organization

#### File Structure

```
package-name/
├── package.go          # Main package code
├── package_test.go     # Tests
├── package_bench_test.go # Benchmarks
├── README.md           # Package documentation
└── examples/           # Usage examples
```

#### Import Order

```go
import (
    // Standard library
    "context"
    "fmt"
    "time"
    
    // Third-party packages
    "github.com/go-redis/redis/v8"
    "go.uber.org/zap"
    
    // Internal packages
    "biometrics-cli/pkg/auth"
    "biometrics-cli/pkg/cache"
)
```

### Error Handling

```go
// Always check errors
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Use error wrapping for context
if err := validateInput(input); err != nil {
    return fmt.Errorf("invalid input: %w", err)
}

// Don't ignore errors
_, err := someFunction()
if err != nil {
    log.Printf("Warning: someFunction failed: %v", err)
}
```

### Comments and Documentation

```go
// Godoc comments for exported symbols
// UserService handles user-related operations.
type UserService struct {
    // cache stores user data temporarily
    cache *redis.Client
    
    // logger for debugging
    logger *zap.Logger
}

// GetUser retrieves a user by ID.
// Returns an error if the user is not found.
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    // Implementation
}

// Inline comments explain WHY, not WHAT
// Use atomic operation to prevent race conditions
atomic.AddInt64(&counter, 1)
```

---

## Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Code style changes (formatting, semicolons, etc.)
- **refactor**: Code refactoring (no functional changes)
- **perf**: Performance improvements
- **test**: Adding or updating tests
- **chore**: Maintenance tasks
- **ci**: CI/CD configuration
- **build**: Build system changes
- **revert**: Reverting previous commits

### Examples

```bash
# Feature
feat(auth): implement OAuth2 token refresh flow

# Bug fix
fix(cache): resolve Redis connection timeout issue

# Documentation
docs(readme): add installation instructions

# Refactor
refactor(validation): simplify input sanitization logic

# Performance
perf(audit): reduce database queries by 50%

# Tests
test(oauth2): add integration tests for token validation

# Chore
chore(deps): update go-redis to v8.11.5
```

### Scope Examples

- `auth`: Authentication and authorization
- `cache`: Caching layer
- `validation`: Input validation
- `audit`: Audit logging
- `cli`: Command-line interface
- `config`: Configuration management
- `deps`: Dependencies

### Body and Footer

```bash
feat(auth): implement mTLS certificate validation

Add mutual TLS authentication for secure service-to-service
communication. Certificates are validated against the CA
bundle stored in Vault.

Closes #123
BREAKING CHANGE: Requires clients to provide valid certificates

Signed-off-by: Your Name <your.email@example.com>
```

---

## Branch Naming Convention

Use descriptive branch names following this pattern:

```
<type>/<description>
```

### Types

- **feature**: New features
- **fix**: Bug fixes
- **docs**: Documentation
- **refactor**: Code refactoring
- **test**: Tests
- **chore**: Maintenance
- **hotfix**: Urgent production fixes

### Examples

```bash
# Good branch names
feature/oauth2-token-refresh
fix/redis-connection-timeout
docs/installation-guide
refactor/validation-logic
test/middleware-integration
chore/update-dependencies
hotfix/critical-security-patch

# Bad branch names
patch-1
new-feature
fix-stuff
test
```

---

## Submitting Pull Requests

### Before Submitting

1. **Update Your Branch**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run All Tests**
   ```bash
   go test -v -race ./...
   ```

3. **Run Linters**
   ```bash
   golangci-lint run --fix
   ```

4. **Check Code Coverage**
   ```bash
   go test -coverprofile=coverage.out ./...
   go tool cover -func=coverage.out
   # Ensure coverage is >= 95%
   ```

5. **Build Successfully**
   ```bash
   go build -o bin/biometrics ./cmd/biometrics
   ```

### Pull Request Template

When creating a PR, include:

```markdown
## Description
Brief description of the changes

## Related Issue
Fixes #<issue-number>

## Type of Change
- [ ] 🐛 Bug fix
- [ ] ✨ New feature
- [ ] 💥 Breaking change
- [ ] 📝 Documentation
- [ ] 🔧 Configuration
- [ ] ♻️ Refactoring
- [ ] ✅ Tests

## Checklist
- [ ] Code follows project guidelines
- [ ] Self-review completed
- [ ] Tests added/updated (95%+ coverage)
- [ ] Documentation updated
- [ ] No new warnings
- [ ] Tested locally

## Testing Instructions
How to test your changes:
1. Step 1
2. Step 2
3. Step 3

## Screenshots (if applicable)
Add screenshots to help explain changes
```

### Review Process

1. **Automated Checks**: CI/CD pipeline must pass
2. **Code Review**: At least one maintainer approval required
3. **Testing**: Changes must be tested by reviewer
4. **Merge**: Squash and merge when approved

---

## Development Workflow

### Daily Workflow

```bash
# Start your day
git fetch upstream
git rebase upstream/main

# Create feature branch
git checkout -b feature/your-feature

# Make changes, commit frequently
git add .
git commit -m "feat: implement feature part 1"
git commit -m "feat: implement feature part 2"

# Before pushing
git fetch upstream
git rebase upstream/main
go test ./...
golangci-lint run

# Push to your fork
git push origin feature/your-feature
```

### Handling Merge Conflicts

```bash
# Fetch upstream changes
git fetch upstream

# Rebase your branch
git rebase upstream/main

# Resolve conflicts manually
# Edit conflicted files

# Continue rebase
git add <resolved-files>
git rebase --continue

# Force push (carefully!)
git push --force-with-lease origin feature/your-feature
```

---

## Testing Requirements

### Unit Tests

- **Coverage**: Minimum 95% code coverage required
- **Isolation**: Tests must be independent and isolated
- **Repeatability**: Tests must produce same results every time
- **Speed**: Unit tests should be fast (< 100ms each)

### Integration Tests

- **Real Dependencies**: Use actual Redis, PostgreSQL instances
- **Cleanup**: Clean up test data after each test
- **Timeouts**: Set appropriate timeouts for external calls

### Example Test Structure

```go
func TestUserService_GetUser(t *testing.T) {
    // Arrange
    ctx := context.Background()
    userID := "test-user-123"
    
    // Act
    user, err := service.GetUser(ctx, userID)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, userID, user.ID)
}
```

---

## Documentation Standards

### Code Comments

- **Exported Symbols**: Require godoc comments
- **Complex Logic**: Explain the "why", not the "what"
- **Public APIs**: Include usage examples

### README Files

Every package must have a README.md with:

- Purpose and overview
- Usage examples
- Configuration options
- Common issues and solutions

### API Documentation

- Use godoc for API documentation
- Include examples for all public functions
- Document error cases and edge cases

---

## Questions?

### Getting Help

- **GitHub Issues**: Open an issue for bugs or feature requests
- **Discussions**: Use GitHub Discussions for questions
- **Email**: Contact maintainers directly for sensitive issues

### Common Issues

**Q: Tests are failing locally but pass in CI**

A: Ensure your local environment matches CI:
```bash
docker-compose up -d redis postgres
go clean -testcache
go test -v ./...
```

**Q: Linters report errors I don't understand**

A: Run with verbose output:
```bash
golangci-lint run -v
```

**Q: How do I add a new dependency?**

A: Use go get and update go.mod:
```bash
go get github.com/new/dependency@latest
go mod tidy
git add go.mod go.sum
```

---

## Recognition

Contributors will be recognized in:

- README.md contributors section
- Release notes for significant contributions
- Annual contributor highlights

Thank you for contributing to BIOMETRICS CLI! 🎉

---

**Last Updated**: February 2026  
**Version**: 1.0.0  
**Maintained By**: BIOMETRICS Team
