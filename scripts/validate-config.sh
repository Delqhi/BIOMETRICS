#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

ERRORS=0
WARNINGS=0

log_error() { echo -e "${RED}[ERROR]${NC} $1"; ((ERRORS++)); }
log_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; ((WARNINGS++)); }
log_success() { echo -e "${GREEN}[OK]${NC} $1"; }

validate_yaml() {
    local file=$1
    if command -v python3 &>/dev/null; then
        if python3 -c "import yaml; yaml.safe_load(open('$file'))" 2>/dev/null; then
            log_success "$file: Valid YAML"
            return 0
        else
            log_error "$file: Invalid YAML syntax"
            return 1
        fi
    elif command -v yq &>/dev/null; then
        if yq '.' "$file" &>/dev/null; then
            log_success "$file: Valid YAML"
            return 0
        else
            log_error "$file: Invalid YAML syntax"
            return 1
        fi
    else
        log_warning "Skipping YAML validation (no python3 or yq)"
        return 0
    fi
}

validate_json() {
    local file=$1
    if python3 -c "import json; json.load(open('$file'))" 2>/dev/null; then
        log_success "$file: Valid JSON"
        return 0
    else
        log_error "$file: Invalid JSON syntax"
        return 1
    fi
}

check_required_fields() {
    local file=$1
    local required_field=$2
    
    if grep -q "$required_field" "$file" 2>/dev/null; then
        log_success "$file: Contains $required_field"
        return 0
    else
        log_error "$file: Missing required field: $required_field"
        return 1
    fi
}

detect_secrets() {
    local file=$1
    local secrets=(
        "nvapi-[a-zA-Z0-9_-]{20,}"
        "sk-[a-zA-Z0-9_-]{20,}"
        "ghp_[a-zA-Z0-9_-]{36}"
        "AKIA[0-9A-Z]{16}"
    )
    
    for pattern in "${secrets[@]}"; do
        if grep -E "$pattern" "$file" 2>/dev/null | grep -v "^#" | grep -v "^[[:space:]]*#" >/dev/null; then
            log_warning "$file: Possible secret detected (pattern: $pattern)"
        fi
    done
}

echo "=========================================="
echo "BIOMETRICS Configuration Validator"
echo "=========================================="
echo ""

echo "=== Checking docker-compose.yml ==="
if [[ -f "$PROJECT_ROOT/docker-compose.yml" ]]; then
    validate_yaml "$PROJECT_ROOT/docker-compose.yml"
else
    log_error "docker-compose.yml not found"
fi

echo ""
echo "=== Checking Makefile ==="
if [[ -f "$PROJECT_ROOT/Makefile" ]]; then
    log_success "Makefile exists"
else
    log_error "Makefile not found"
fi

echo ""
echo "=== Checking .env.example ==="
if [[ -f "$PROJECT_ROOT/.env.example" ]]; then
    log_success ".env.example exists"
    detect_secrets "$PROJECT_ROOT/.env.example"
else
    log_error ".env.example not found"
fi

echo ""
echo "=== Checking .env (if exists) ==="
if [[ -f "$PROJECT_ROOT/.env" ]]; then
    detect_secrets "$PROJECT_ROOT/.env"
    
    if grep -q "NVIDIA_API_KEY=nvapi-YOUR_KEY" "$PROJECT_ROOT/.env" 2>/dev/null; then
        log_warning ".env still contains placeholder NVIDIA_API_KEY"
    fi
else
    log_warning ".env not found (run: make env)"
fi

echo ""
echo "=== Checking Go configuration ==="
if [[ -f "$PROJECT_ROOT/biometrics-cli/go.mod" ]]; then
    log_success "go.mod exists"
else
    log_error "go.mod not found"
fi

echo ""
echo "=== Checking GitHub workflows ==="
if [[ -d "$PROJECT_ROOT/.github/workflows" ]]; then
    count=$(find "$PROJECT_ROOT/.github/workflows" -name "*.yml" -o -name "*.yaml" | wc -l)
    if [[ $count -gt 0 ]]; then
        log_success "Found $count workflow file(s)"
    else
        log_warning "No workflow files found"
    fi
else
    log_warning ".github/workflows directory not found"
fi

echo ""
echo "=== Checking scripts ==="
for script in setup.sh validate-config.sh; do
    if [[ -f "$PROJECT_ROOT/scripts/$script" ]]; then
        if [[ -x "$PROJECT_ROOT/scripts/$script" ]]; then
            log_success "$script is executable"
        else
            log_warning "$script is not executable"
        fi
    else
        log_error "$script not found"
    fi
done

echo ""
echo "=========================================="
echo "Summary"
echo "=========================================="
echo -e "Errors:   ${RED}$ERRORS${NC}"
echo -e "Warnings: ${YELLOW}$WARNINGS${NC}"
echo ""

if [[ $ERRORS -gt 0 ]]; then
    echo "Validation FAILED"
    exit 1
elif [[ $WARNINGS -gt 0 ]]; then
    echo "Validation passed with warnings"
    exit 0
else
    echo "Validation PASSED"
    exit 0
fi
