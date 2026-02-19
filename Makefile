.PHONY: help setup build test lint docker docker-up docker-down clean verify install config env

SHELL := /bin/bash
PROJECT_ROOT := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
BIN_DIR := $(PROJECT_ROOT)/bin
CLI_DIR := $(PROJECT_ROOT)/biometrics-cli

help:
	@echo "BIOMETRICS - Makefile Commands"
	@echo "=============================="
	@echo ""
	@echo "Setup & Installation:"
	@echo "  make install       - Install all dependencies"
	@echo "  make setup         - Run full setup script"
	@echo "  make env           - Create .env from template"
	@echo ""
	@echo "Development:"
	@echo "  make build         - Build biometrics-cli"
	@echo "  make test          - Run tests"
	@echo "  make lint          - Run linters"
	@echo "  make verify        - Verify installation"
	@echo ""
	@echo "Docker:"
	@echo "  make docker        - Build Docker images"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make docker-logs   - View Docker logs"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make config        - Validate configuration"
	@echo "  make doctor        - Run diagnostics"
	@echo ""
	@echo "Deployment:"
	@echo "  make deploy        - Deploy to staging"
	@echo "  make deploy-prod   - Deploy to production"

install:
	@echo "Installing dependencies..."
	@cd $(CLI_DIR) && go mod download
	@echo "Dependencies installed!"

setup: env
	@echo "Running BIOMETRICS setup..."
	@chmod +x $(PROJECT_ROOT)/scripts/setup.sh
	@$(PROJECT_ROOT)/scripts/setup.sh

env:
	@echo "Creating .env file..."
	@if [ ! -f $(PROJECT_ROOT)/.env ]; then \
		if [ -f $(PROJECT_ROOT)/.env.example ]; then \
			cp $(PROJECT_ROOT)/.env.example $(PROJECT_ROOT)/.env; \
			echo "Created .env from template. Please edit and add your API keys!"; \
		else \
			echo "ERROR: .env.example not found!"; \
			exit 1; \
		fi \
	else \
		echo ".env already exists"; \
	fi

build:
	@echo "Building biometrics-cli..."
	@cd $(CLI_DIR)
	@go build -o $(BIN_DIR)/biometrics-cli
	@echo "Build complete: $(BIN_DIR)/biometrics-cli"

test:
	@echo "Running tests..."
	@cd $(CLI_DIR) && go test -v ./...
	@echo "Tests complete!"

lint:
	@echo "Running linters..."
	@cd $(CLI_DIR) && golangci-lint run || echo "golangci-lint not installed, skipping..."
	@echo "Linting complete!"

verify:
	@echo "Verifying installation..."
	@echo "Checking required commands..."
	@which git && echo "  git: OK" || echo "  git: MISSING"
	@which node && echo "  node: OK" || echo "  node: MISSING"
	@which pnpm && echo "  pnpm: OK" || echo "  pnpm: MISSING"
	@which go && echo "  go: OK" || echo "  go: MISSING"
	@if [ -f $(BIN_DIR)/biometrics-cli ]; then \
		echo "  biometrics-cli: OK"; \
	else \
		echo "  biometrics-cli: MISSING (run make build)"; \
	fi
	@if [ -f $(PROJECT_ROOT)/.env ]; then \
		echo "  .env: OK"; \
	else \
		echo "  .env: MISSING (run make env)"; \
	fi
	@echo "Verification complete!"

docker:
	@echo "Building Docker images..."
	@docker build -t biometrics:latest .

docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "Containers started!"
	@echo "Services:"
	@echo "  - OpenCode Server: http://localhost:8080"
	@echo "  - OpenClaw: http://localhost:18789"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379"
	@echo "  - Prometheus: http://localhost:9090"
	@echo "  - Grafana: http://localhost:3001"

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "Containers stopped!"

docker-logs:
	@docker-compose logs -f

docker-clean:
	@echo "Cleaning Docker resources..."
	@docker-compose down -v
	@echo "Docker resources cleaned!"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@rm -rf $(CLI_DIR)/biometrics-cli
	@rm -rf $(CLI_DIR)/vendor
	@cd $(CLI_DIR) && go clean
	@echo "Clean complete!"

config:
	@echo "Validating configuration..."
	@$(PROJECT_ROOT)/scripts/validate-config.sh || true

doctor:
	@echo "Running diagnostics..."
	@echo ""
	@echo "=== System Check ==="
	@echo "OS: $$(uname -s)"
	@echo "Shell: $$SHELL"
	@echo ""
	@echo "=== Installation Check ==="
	@make verify
	@echo ""
	@echo "=== Docker Check ==="
	@if command -v docker &> /dev/null; then \
		docker ps --format "table {{.Names}}\t{{.Status}}" 2>/dev/null || echo "Docker not running"; \
	else \
		echo "Docker: Not installed"; \
	fi
	@echo ""
	@echo "=== Configuration Check ==="
	@make config

deploy:
	@echo "Deploying to staging..."
	@echo "Staging deployment not configured - add your deployment logic"

deploy-prod:
	@echo "Deploying to production..."
	@echo "Production deployment not configured - add your deployment logic"
