#!/bin/bash
# BIOMETRICS CLI - Development Setup Verification Script
# Enterprise Practices Feb 2026

set -e

echo "🔍 BIOMETRICS CLI - Development Setup Verification"
echo "=================================================="
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check Go installation
echo "📦 Checking Go installation..."
if command -v go &> /dev/null; then
    GO_VERSION=$(go version)
    echo -e "${GREEN}✓${NC} $GO_VERSION"
else
    echo -e "${RED}✗${NC} Go is not installed"
    echo "  Install from: https://golang.org/dl/"
    exit 1
fi

# Check Docker installation
echo ""
echo "🐳 Checking Docker installation..."
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version)
    echo -e "${GREEN}✓${NC} $DOCKER_VERSION"
else
    echo -e "${RED}✗${NC} Docker is not installed"
    echo "  Install from: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check Docker Compose
echo ""
echo "🚀 Checking Docker Compose..."
if command -v docker-compose &> /dev/null; then
    COMPOSE_VERSION=$(docker-compose --version)
    echo -e "${GREEN}✓${NC} $COMPOSE_VERSION"
else
    echo -e "${YELLOW}⚠${NC} docker-compose not found (optional)"
fi

# Check golangci-lint
echo ""
echo "🔍 Checking golangci-lint..."
if command -v golangci-lint &> /dev/null; then
    LINT_VERSION=$(golangci-lint --version)
    echo -e "${GREEN}✓${NC} $LINT_VERSION"
else
    echo -e "${YELLOW}⚠${NC} golangci-lint not installed"
    echo "  Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
fi

# Validate docker-compose.yml
echo ""
echo "📋 Validating docker-compose.yml..."
if docker-compose config --quiet 2>&1; then
    echo -e "${GREEN}✓${NC} docker-compose.yml is valid"
else
    echo -e "${RED}✗${NC} docker-compose.yml has errors"
    exit 1
fi

# Check Go modules
echo ""
echo "📦 Checking Go modules..."
if [ -f "go.mod" ]; then
    echo -e "${GREEN}✓${NC} go.mod found"
    go mod verify 2>&1 | head -1
else
    echo -e "${RED}✗${NC} go.mod not found"
    exit 1
fi

# Check required files
echo ""
echo "📄 Checking required files..."
FILES=(".golangci.yml" "docker-compose.yml" "CONTRIBUTING.md" "go.mod" "README.md")
for file in "${FILES[@]}"; do
    if [ -f "$file" ]; then
        echo -e "${GREEN}✓${NC} $file exists"
    else
        echo -e "${RED}✗${NC} $file missing"
    fi
done

# Test Docker services (optional)
echo ""
echo "🧪 Testing Docker services..."
if docker-compose ps 2>&1 | grep -q "redis\|postgres"; then
    echo -e "${GREEN}✓${NC} Services are running"
else
    echo -e "${YELLOW}⚠${NC} Services not running (start with: docker-compose up -d)"
fi

# Summary
echo ""
echo "=================================================="
echo "✅ Verification Complete!"
echo ""
echo "Next Steps:"
echo "  1. Start services: docker-compose up -d"
echo "  2. Install linter: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
echo "  3. Run tests: go test -v ./..."
echo "  4. Build CLI: go build -o bin/biometrics ./cmd/biometrics"
echo ""
echo "Documentation:"
echo "  - CONTRIBUTING.md: Development guidelines"
echo "  - README.md: Project overview"
echo "  - .golangci.yml: Linter configuration"
echo ""
