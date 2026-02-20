#!/bin/bash

# 🚀 BIOMETRICS Bootstrap Script
# Complete setup in one command
# Usage: ./bootstrap.sh

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_header() {
    echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║${NC} $1"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
    echo
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

check_command() {
    if ! command -v $1 &> /dev/null; then
        print_error "$1 is not installed"
        return 1
    fi
    print_success "$1 is installed ($( $1 --version | head -n1 ))"
    return 0
}

# Main script
print_header "🚀 BIOMETRICS Bootstrap Script"

echo "This script will:"
echo "  1. Check system requirements"
echo "  2. Install dependencies"
echo "  3. Setup OpenCode and providers"
echo "  4. Configure environment"
echo "  5. Verify installation"
echo

read -p "Continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_info "Aborted"
    exit 1
fi

echo

# Step 1: Check system requirements
print_header "Step 1: Checking System Requirements"

# Check OS
if [[ "$OSTYPE" == "darwin"* ]]; then
    print_success "macOS detected: $(sw_vers -productVersion)"
else
    print_error "Unsupported OS: $OSTYPE (only macOS supported)"
    exit 1
fi

# Check Node.js
if check_command node; then
    NODE_VERSION=$(node --version | cut -d'v' -f2 | cut -d'.' -f1)
    if [ $NODE_VERSION -lt 20 ]; then
        print_error "Node.js v20+ required (found v$NODE_VERSION)"
        print_info "Install with: brew install node@20"
        exit 1
    fi
fi

# Check Go
if check_command go; then
    GO_VERSION=$(go version | cut -d' ' -f3 | cut -d'o' -f2)
    print_success "Go version: $GO_VERSION"
fi

# Check Git
check_command git

# Check Docker (optional)
if check_command docker; then
    print_info "Docker is installed (optional)"
else
    print_info "Docker not installed (optional, can be skipped)"
fi

echo

# Step 2: Install dependencies
print_header "Step 2: Installing Dependencies"

# Install Node.js dependencies
print_info "Installing Node.js dependencies..."
if [ -f "package.json" ]; then
    npm install
    print_success "Node.js dependencies installed"
else
    print_error "package.json not found"
    exit 1
fi

# Install Go dependencies
print_info "Installing Go dependencies..."
if [ -f "biometrics-cli/go.mod" ]; then
    cd biometrics-cli
    go mod download
    cd ..
    print_success "Go dependencies installed"
else
    print_error "biometrics-cli/go.mod not found"
    exit 1
fi

# Install Python dependencies (optional)
if [ -f "requirements.txt" ]; then
    print_info "Installing Python dependencies..."
    pip install -r requirements.txt
    print_success "Python dependencies installed"
else
    print_info "requirements.txt not found (skipping Python)"
fi

echo

# Step 3: Setup OpenCode
print_header "Step 3: Setting Up OpenCode"

# Install OpenCode
print_info "Installing OpenCode..."
npm install -g opencode
print_success "OpenCode installed ($(opencode --version))"

# Authenticate providers
print_info "Setting up provider authentication..."

echo
print_info "You will need API keys for:"
echo "  - NVIDIA NIM (Qwen 3.5 397B) - FREE"
echo "  - Moonshot AI (Kimi K2.5, MiniMax M2.5) - FREE"
echo

read -p "Do you have API keys ready? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    print_info "Authenticating NVIDIA NIM..."
    opencode auth add nvidia-nim
    
    print_info "Authenticating Moonshot AI..."
    opencode auth add moonshot-ai
    
    print_success "Providers authenticated"
else
    print_info "You can authenticate later with:"
    echo "  opencode auth add nvidia-nim"
    echo "  opencode auth add moonshot-ai"
fi

echo

# Step 4: Configure environment
print_header "Step 4: Configuring Environment"

# Create .env if not exists
if [ ! -f ".env" ]; then
    print_info "Creating .env file..."
    cp .env.example .env 2>/dev/null || echo "# Add your API keys here" > .env
    print_success ".env file created"
else
    print_info ".env already exists"
fi

# Add API keys to .env
print_info "Edit .env and add your API keys:"
echo "  NVIDIA_API_KEY=your-key-here"
echo "  MOONSHOT_API_KEY=your-key-here"
echo

# Add to shell profile (optional)
read -p "Add API keys to ~/.zshrc? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    print_info "Adding to ~/.zshrc..."
    echo "" >> ~/.zshrc
    echo "# BIOMETRICS API Keys" >> ~/.zshrc
    echo "export NVIDIA_API_KEY=\"your-key-here\"" >> ~/.zshrc
    echo "export MOONSHOT_API_KEY=\"your-key-here\"" >> ~/.zshrc
    print_success "Added to ~/.zshrc (edit with your actual keys)"
    print_info "Reload shell with: source ~/.zshrc"
fi

echo

# Step 5: Verify installation
print_header "Step 5: Verifying Installation"

print_info "Checking available models..."
echo

if opencode models | grep -q "nvidia"; then
    print_success "NVIDIA NIM models available"
    opencode models | grep nvidia | head -n3
else
    print_error "NVIDIA NIM models not found"
    print_info "Authenticate with: opencode auth add nvidia-nim"
fi

echo

if opencode models | grep -q "kimi\|minimax"; then
    print_success "OpenCode ZEN models available"
    opencode models | grep -E "kimi|minimax" | head -n3
else
    print_error "OpenCode ZEN models not found"
    print_info "Authenticate with: opencode auth add moonshot-ai"
fi

echo

# Test run
print_info "Running test agent..."
TEST_OUTPUT=$(opencode "Say 'BIOMETRICS is ready!' in one sentence" 2>&1)
if [ $? -eq 0 ]; then
    print_success "Test agent completed successfully"
else
    print_error "Test agent failed"
    print_info "Check error above and fix before continuing"
fi

echo

# Final summary
print_header "🎉 Setup Complete!"

echo "BIOMETRICS is now ready to use!"
echo
echo "Next steps:"
echo "  1. Edit .env and add your actual API keys"
echo "  2. Read AGENTS-PLAN.md to see current tasks"
echo "  3. Read ONBOARDING.md for complete guide"
echo "  4. Run your first agent:"
echo "     opencode \"Your task description\""
echo
echo "Essential commands:"
echo "  opencode models              # List available models"
echo "  opencode agents              # List available agents"
echo "  opencode \"task\"              # Run agent with task"
echo "  ./biometrics-cli/orchestrator start  # Start 24/7 loop"
echo
echo "Documentation:"
echo "  README.md                    # Main documentation"
echo "  ONBOARDING.md                # Complete onboarding guide"
echo "  AGENTS-PLAN.md               # Current tasks and infinity loop"
echo "  docs/                        # Full documentation"
echo
print_success "Happy coding! 🚀"
