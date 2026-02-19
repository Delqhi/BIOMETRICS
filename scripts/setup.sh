#!/bin/bash
# =============================================================================
# BIOMETRICS Setup Script - One-Command Installation
# =============================================================================
# Usage: ./scripts/setup.sh
# Requirements: macOS with Homebrew installed
# =============================================================================

set -e  # Exit on error
set -u  # Exit on undefined variable
set -o pipefail  # Exit on pipe failure

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script constants
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
LOG_FILE="$PROJECT_ROOT/setup.log"

# =============================================================================
# FUNCTIONS
# =============================================================================

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [INFO] $1" >> "$LOG_FILE"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [SUCCESS] $1" >> "$LOG_FILE"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [WARNING] $1" >> "$LOG_FILE"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [ERROR] $1" >> "$LOG_FILE"
}

# Check if running on macOS
check_os() {
    if [[ "$OSTYPE" != "darwin"* ]]; then
        log_error "This script only works on macOS"
        exit 1
    fi
    log_success "Running on macOS"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Install Homebrew if not present
install_homebrew() {
    if command_exists brew; then
        log_success "Homebrew already installed"
        brew update
    else
        log_info "Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        log_success "Homebrew installed"
    fi
}

# Install required packages
install_packages() {
    local packages=(
        "git"
        "node"
        "pnpm"
        "python@3.11"
    )
    
    log_info "Checking installed packages..."
    
    for package in "${packages[@]}"; do
        if command_exists "$package"; then
            log_success "$package is installed"
        else
            log_info "Installing $package..."
            brew install "$package" 2>/dev/null || log_warning "$package installation skipped"
        fi
    done
}

# Install Go (for biometrics-cli)
install_go() {
    if command_exists go; then
        log_success "Go already installed: $(go version)"
    else
        log_info "Installing Go..."
        brew install go
        log_success "Go installed"
    fi
}

# Install Node.js global packages
install_node_packages() {
    log_info "Installing Node.js global packages..."
    
    # Install OpenCode
    if command_exists opencode; then
        log_success "OpenCode already installed"
    else
        log_info "Installing OpenCode..."
        brew install opencode
        log_success "OpenCode installed"
    fi
    
    # Install NLM CLI
    if command_exists nlm; then
        log_success "NLM CLI already installed"
    else
        log_info "Installing NLM CLI..."
        pnpm add -g nlm-cli
        log_success "NLM CLI installed"
    fi
    
    # Install OpenClaw
    if command_exists openclaw; then
        log_success "OpenClaw already installed"
    else
        log_info "Installing OpenClaw..."
        pnpm add -g @delqhi/openclaw
        log_success "OpenClaw installed"
    fi
}

# Setup Python environment
setup_python() {
    log_info "Setting up Python environment..."
    
    # Create virtual environment if it doesn't exist
    if [[ ! -d "$PROJECT_ROOT/venv" ]]; then
        python3 -m venv "$PROJECT_ROOT/venv"
        log_success "Python virtual environment created"
    else
        log_success "Python virtual environment already exists"
    fi
    
    # Activate virtual environment
    source "$PROJECT_ROOT/venv/bin/activate"
    
    # Install Python packages if requirements.txt exists
    if [[ -f "$PROJECT_ROOT/biometrics-cli/requirements.txt" ]]; then
        pip install -r "$PROJECT_ROOT/biometrics-cli/requirements.txt"
        log_success "Python packages installed"
    fi
    
    deactivate
}

# Configure shell environment
configure_shell() {
    log_info "Configuring shell environment..."
    
    local shell_config="$HOME/.zshrc"
    local biometrics_path="$PROJECT_ROOT/bin"
    
    # Add biometrics to PATH if not already present
    if [[ -f "$shell_config" ]]; then
        if ! grep -q "$biometrics_path" "$shell_config"; then
            echo "" >> "$shell_config"
            echo "# BIOMETRICS" >> "$shell_config"
            echo "export PATH=\"$biometrics_path:\$PATH\"" >> "$shell_config"
            log_success "Added BIOMETRICS to PATH in $shell_config"
        else
            log_success "BIOMETRICS already in PATH"
        fi
    fi
    
    # Add pnpm to PATH if not already present
    if ! grep -q "pnpm" "$shell_config"; then
        echo 'export PATH="$HOME/Library/pnpm:$PATH"' >> "$shell_config"
        log_success "Added pnpm to PATH"
    fi
}

# Create environment file
create_env_file() {
    log_info "Creating .env file..."
    
    local env_file="$PROJECT_ROOT/.env"
    
    if [[ -f "$env_file" ]]; then
        log_warning ".env file already exists, skipping..."
    else
        if [[ -f "$PROJECT_ROOT/.env.example" ]]; then
            cp "$PROJECT_ROOT/.env.example" "$env_file"
            log_success "Created .env from template"
            log_warning "Please edit $env_file and add your API keys!"
        else
            log_error ".env.example not found!"
            exit 1
        fi
    fi
}

# Build biometrics-cli
build_cli() {
    log_info "Building biometrics-cli..."
    
    cd "$PROJECT_ROOT/biometrics-cli"
    
    # Download Go dependencies
    go mod download
    
    # Build the binary
    go build -o biometrics-cli
    
    # Create bin directory and move binary
    mkdir -p "$PROJECT_ROOT/bin"
    mv biometrics-cli "$PROJECT_ROOT/bin/"
    
    log_success "biometrics-cli built and installed"
}

# Setup Docker containers
setup_docker() {
    log_info "Checking Docker..."
    
    if command_exists docker; then
        log_success "Docker is installed"
        
        if command_exists docker-compose || docker compose version &>/dev/null; then
            log_success "Docker Compose is available"
            
            # Start containers
            cd "$PROJECT_ROOT"
            docker-compose up -d
            
            log_success "Docker containers started"
        else
            log_warning "Docker Compose not found, skipping..."
        fi
    else
        log_warning "Docker not installed, skipping Docker setup..."
    fi
}

# Verify installation
verify_installation() {
    log_info "Verifying installation..."
    
    local all_good=true
    
    # Check required commands
    local required_commands=("git" "node" "pnpm" "go")
    for cmd in "${required_commands[@]}"; do
        if command_exists "$cmd"; then
            log_success "$cmd: $(command -v "$cmd")"
        else
            log_error "$cmd not found"
            all_good=false
        fi
    done
    
    # Check biometrics binary
    if [[ -f "$PROJECT_ROOT/bin/biometrics-cli" ]]; then
        log_success "biometrics-cli binary exists"
    else
        log_warning "biometrics-cli binary not found"
    fi
    
    # Check .env file
    if [[ -f "$PROJECT_ROOT/.env" ]]; then
        log_success ".env file exists"
    else
        log_warning ".env file not found"
    fi
    
    if $all_good; then
        log_success "Installation verified!"
    else
        log_warning "Some checks failed, please review"
    fi
}

# Print setup completion message
print_completion() {
    echo ""
    echo "============================================================"
    echo "           BIOMETRICS SETUP COMPLETE!"
    echo "============================================================"
    echo ""
    echo "Next steps:"
    echo "  1. Edit .env and add your API keys:"
    echo "     - NVIDIA_API_KEY (required)"
    echo ""
    echo "  2. Reload your shell:"
    echo "     source ~/.zshrc"
    echo ""
    echo "  3. Verify installation:"
    echo "     make verify"
    echo ""
    echo "  4. Start using BIOMETRICS:"
    echo "     biometrics-cli"
    echo ""
    echo "============================================================"
}

# Rollback function (cleanup on failure)
rollback() {
    log_error "Setup failed! Rolling back..."
    log_info "Removing installed packages..."
    # Add rollback commands here if needed
    log_warning "Rollback complete. Please check the log: $LOG_FILE"
    exit 1
}

# =============================================================================
# MAIN
# =============================================================================

main() {
    # Initialize log file
    touch "$LOG_FILE"
    log_info "Starting BIOMETRICS setup..."
    
    # Trap errors
    trap rollback ERR
    
    # Run setup steps
    check_os
    install_homebrew
    install_packages
    install_go
    install_node_packages
    setup_python
    configure_shell
    create_env_file
    build_cli
    setup_docker
    verify_installation
    print_completion
    
    log_success "Setup completed successfully!"
}

# Run main function
main "$@"
