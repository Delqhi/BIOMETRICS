#!/bin/bash

# BIOMETRICS Security Audit Script
# Performs comprehensive security audit of Go dependencies and code
# Usage: ./scripts/security-audit.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Counters
VULNERABILITIES_FOUND=0
WARNINGS_FOUND=0
INFO_FOUND=0

echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     BIOMETRICS Security Audit Script                   ║${NC}"
echo -e "${BLUE}║     Comprehensive Security scanning                    ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

# Function to print section headers
print_section() {
    echo -e "\n${BLUE}════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}════════════════════════════════════════════════════════${NC}\n"
}

# Function to print status
print_status() {
    local status=$1
    local message=$2
    
    case $status in
        "PASS")
            echo -e "${GREEN}✓${NC} $message"
            ;;
        "WARN")
            echo -e "${YELLOW}⚠${NC} $message"
            ((WARNINGS_FOUND++))
            ;;
        "FAIL")
            echo -e "${RED}✗${NC} $message"
            ((VULNERABILITIES_FOUND++))
            ;;
        "INFO")
            echo -e "${BLUE}ℹ${NC} $message"
            ((INFO_FOUND++))
            ;;
    esac
}

# Check if Go is installed
check_go() {
    print_section "Checking Go Installation"
    
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version)
        print_status "PASS" "$GO_VERSION"
    else
        print_status "FAIL" "Go is not installed"
        exit 1
    fi
}

# Check for govulncheck
check_govulncheck() {
    print_section "Checking govulncheck Installation"
    
    if command -v govulncheck &> /dev/null; then
        VULNCHECK_VERSION=$(govulncheck -version 2>&1 | head -1)
        print_status "PASS" "$VULNCHECK_VERSION"
    else
        print_status "WARN" "govulncheck not found, installing..."
        go install golang.org/x/vuln/cmd/govulncheck@latest
        print_status "PASS" "govulncheck installed successfully"
    fi
}

# Audit Go modules
audit_go_modules() {
    print_section "Auditing Go Modules"
    
    local modules=("biometrics-cli" "BIOMETRICS/biometrics")
    
    for module in "${modules[@]}"; do
        if [ -d "$PROJECT_ROOT/$module" ]; then
            echo -e "${BLUE}Scanning: $module${NC}"
            cd "$PROJECT_ROOT/$module"
            
            # Check for go.mod
            if [ -f "go.mod" ]; then
                print_status "PASS" "go.mod found"
                
                # Run govulncheck
                echo "Running govulncheck..."
                if govulncheck ./... 2>&1 | tee /tmp/govulncheck-output.txt; then
                    print_status "PASS" "No vulnerabilities found in $module"
                else
                    print_status "FAIL" "Vulnerabilities found in $module - see /tmp/govulncheck-output.txt"
                fi
                
                # Check for outdated dependencies
                echo "Checking for outdated dependencies..."
                if go list -m -u all 2>&1 | grep -E '\[.*\]' > /tmp/outdated-deps.txt; then
                    OUTDATED_COUNT=$(wc -l < /tmp/outdated-deps.txt)
                    if [ "$OUTDATED_COUNT" -gt 0 ]; then
                        print_status "WARN" "$OUTDATED_COUNT outdated dependencies found"
                        echo "Outdated dependencies:"
                        cat /tmp/outdated-deps.txt | head -10
                        if [ "$OUTDATED_COUNT" -gt 10 ]; then
                            echo "... and $((OUTDATED_COUNT - 10)) more"
                        fi
                    else
                        print_status "PASS" "All dependencies are up to date"
                    fi
                else
                    print_status "PASS" "All dependencies are up to date"
                fi
                
                # Check for unused dependencies
                echo "Checking for unused dependencies..."
                if command -v go-mod-tidy &> /dev/null; then
                    go mod tidy
                    print_status "PASS" "go mod tidy completed"
                else
                    go mod tidy
                    print_status "PASS" "go mod tidy completed"
                fi
                
            else
                print_status "FAIL" "go.mod not found in $module"
            fi
            
            cd "$PROJECT_ROOT"
        else
            print_status "WARN" "Module directory not found: $module"
        fi
    done
}

# Check for secrets in code
check_secrets() {
    print_section "Checking for Secrets in Code"
    
    # Check for hardcoded API keys
    echo "Scanning for hardcoded secrets..."
    
    # Patterns to search for
    local patterns=(
        "nvapi-[a-zA-Z0-9]"
        "sk_live_"
        "sk_test_"
        "ghp_"
        "AIza"
        "-----BEGIN.*PRIVATE KEY-----"
        "postgres://.*:.*@"
        "redis://.*:.*@"
    )
    
    local found_secrets=0
    
    for pattern in "${patterns[@]}"; do
        if grep -r --include="*.go" --include="*.ts" --include="*.js" --include="*.py" \
             --exclude-dir=node_modules --exclude-dir=vendor \
             "$pattern" "$PROJECT_ROOT" 2>/dev/null | grep -v ".git" | grep -v "test" | grep -v "example"; then
            found_secrets=1
        fi
    done
    
    if [ $found_secrets -eq 0 ]; then
        print_status "PASS" "No hardcoded secrets found"
    else
        print_status "FAIL" "Potential hardcoded secrets detected - review output above"
    fi
    
    # Check .env files
    echo "Checking .env files..."
    if find "$PROJECT_ROOT" -name ".env" -not -path "*/node_modules/*" -not -path "*/vendor/*" | grep -q .; then
        print_status "WARN" ".env files found - ensure they are in .gitignore"
        find "$PROJECT_ROOT" -name ".env" -not -path "*/node_modules/*" -not -path "*/vendor/*"
    else
        print_status "PASS" "No .env files found (good if using .env.example)"
    fi
}

# Check file permissions
check_permissions() {
    print_section "Checking File Permissions"
    
    # Check for world-writable files
    WORLD_WRITABLE=$(find "$PROJECT_ROOT" -type f -perm -002 -not -path "*/.git/*" 2>/dev/null | wc -l)
    if [ "$WORLD_WRITABLE" -gt 0 ]; then
        print_status "FAIL" "$WORLD_WRITABLE world-writable files found"
        find "$PROJECT_ROOT" -type f -perm -002 -not -path "*/.git/*" 2>/dev/null | head -10
    else
        print_status "PASS" "No world-writable files found"
    fi
    
    # Check for sensitive files with wrong permissions
    for file in "*.key" "*.pem" "*.crt" "id_rsa" "id_ed25519"; do
        if find "$PROJECT_ROOT" -name "$file" -not -path "*/.git/*" 2>/dev/null | grep -q .; then
            print_status "WARN" "Sensitive file pattern found: $file"
            find "$PROJECT_ROOT" -name "$file" -not -path "*/.git/*" 2>/dev/null
        fi
    done
}

# Check git configuration
check_git() {
    print_section "Checking Git Configuration"
    
    # Check if .gitignore exists
    if [ -f "$PROJECT_ROOT/.gitignore" ]; then
        print_status "PASS" ".gitignore found"
        
        # Check if sensitive patterns are in .gitignore
        local required_patterns=(".env" "*.key" "*.pem" "id_rsa" ".DS_Store")
        for pattern in "${required_patterns[@]}"; do
            if grep -q "$pattern" "$PROJECT_ROOT/.gitignore"; then
                print_status "PASS" ".gitignore contains: $pattern"
            else
                print_status "WARN" ".gitignore missing pattern: $pattern"
            fi
        done
    else
        print_status "FAIL" ".gitignore not found"
    fi
    
    # Check for secrets in git history (basic check)
    echo "Checking git history for secrets (basic scan)..."
    if command -v gitleaks &> /dev/null; then
        if gitleaks detect --source "$PROJECT_ROOT" --no-git 2>/dev/null; then
            print_status "PASS" "Gitleaks scan passed"
        else
            print_status "FAIL" "Gitleaks found potential secrets"
        fi
    else
        print_status "INFO" "Gitleaks not installed - skipping deep git history scan"
        echo "Install with: go install github.com/gitleaks/gitleaks@latest"
    fi
}

# Check for security headers in code
check_security_headers() {
    print_section "Checking Security Headers (if applicable)"
    
    # Check for HTTP security headers in Go code
    if grep -r "X-Frame-Options" --include="*.go" "$PROJECT_ROOT" > /dev/null 2>&1; then
        print_status "PASS" "X-Frame-Options header found"
    else
        print_status "INFO" "X-Frame-Options header not found (may not be applicable)"
    fi
    
    if grep -r "X-Content-Type-Options" --include="*.go" "$PROJECT_ROOT" > /dev/null 2>&1; then
        print_status "PASS" "X-Content-Type-Options header found"
    else
        print_status "INFO" "X-Content-Type-Options header not found (may not be applicable)"
    fi
    
    if grep -r "Strict-Transport-Security" --include="*.go" "$PROJECT_ROOT" > /dev/null 2>&1; then
        print_status "PASS" "HSTS header found"
    else
        print_status "INFO" "HSTS header not found (may not be applicable)"
    fi
}

# Generate report
generate_report() {
    print_section "Security Audit Summary"
    
    echo "Vulnerabilities Found: ${RED}$VULNERABILITIES_FOUND${NC}"
    echo "Warnings Found:        ${YELLOW}$WARNINGS_FOUND${NC}"
    echo "Info Messages:         ${BLUE}$INFO_FOUND${NC}"
    echo ""
    
    if [ $VULNERABILITIES_FOUND -gt 0 ]; then
        echo -e "${RED}╔════════════════════════════════════════════════════════╗${NC}"
        echo -e "${RED}║  ⚠ CRITICAL: Security vulnerabilities detected!        ║${NC}"
        echo -e "${RED}║  Please review and fix before deployment.              ║${NC}"
        echo -e "${RED}╚════════════════════════════════════════════════════════╝${NC}"
        exit 1
    elif [ $WARNINGS_FOUND -gt 0 ]; then
        echo -e "${YELLOW}╔════════════════════════════════════════════════════════╗${NC}"
        echo -e "${YELLOW}║  ⚠ WARNING: Some issues need attention                 ║${NC}"
        echo -e "${YELLOW}║  Review warnings above                                 ║${NC}"
        echo -e "${YELLOW}╚════════════════════════════════════════════════════════╝${NC}"
        exit 0
    else
        echo -e "${GREEN}╔════════════════════════════════════════════════════════╗${NC}"
        echo -e "${GREEN}║  ✓ SUCCESS: No critical issues found                   ║${NC}"
        echo -e "${GREEN}║  Security audit passed                                 ║${NC}"
        echo -e "${GREEN}╚════════════════════════════════════════════════════════╝${NC}"
        exit 0
    fi
}

# Main execution
main() {
    cd "$PROJECT_ROOT"
    
    check_go
    check_govulncheck
    audit_go_modules
    check_secrets
    check_permissions
    check_git
    check_security_headers
    generate_report
}

# Run main function
main "$@"
