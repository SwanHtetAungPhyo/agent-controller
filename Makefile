.PHONY: help install-tools setup build test clean dev prod deploy cert logs status stop restart setup-hooks lint validate

DOCKER_COMPOSE_DEV = docker-compose -f docker-compose.dev.yaml
DOCKER_COMPOSE_PROD = docker-compose -f docker-compose.yaml

help:
	@echo "Kainos Microservices Management"
	@echo "================================"
	@echo "Setup:"
	@echo "  install-tools    - Install all required development tools"
	@echo "  setup            - Complete project setup for new developers"
	@echo ""
	@echo "Development:"
	@echo "  dev              - Start development environment"
	@echo "  build            - Build all services"
	@echo "  test             - Run all tests"
	@echo "  logs             - Show service logs"
	@echo "  status           - Show service status"
	@echo "  stop             - Stop all services"
	@echo "  restart          - Restart all services"
	@echo ""
	@echo "Code Quality:"
	@echo "  setup-hooks      - Setup git hooks and pre-commit"
	@echo "  lint             - Run linting and formatting"
	@echo "  validate         - Run comprehensive validation"
	@echo ""
	@echo "Production:"
	@echo "  prod             - Start production environment"
	@echo "  deploy           - Deploy to production"
	@echo ""
	@echo "Utilities:"
	@echo "  cert             - Generate SSL certificates"
	@echo "  clean            - Clean up containers and volumes"

install-tools:
	@echo "Installing required development tools..."
	@echo "Checking Docker..."
	@command -v docker >/dev/null 2>&1 || { echo "Docker is required but not installed. Please install Docker first."; exit 1; }
	@echo "Checking Docker Compose..."
	@command -v docker-compose >/dev/null 2>&1 || { echo "Docker Compose is required but not installed. Please install Docker Compose first."; exit 1; }
	@echo "Installing mkcert for SSL certificates..."
	@if command -v brew >/dev/null 2>&1; then \
		brew list mkcert >/dev/null 2>&1 || brew install mkcert; \
	elif command -v apt-get >/dev/null 2>&1; then \
		sudo apt-get update && sudo apt-get install -y libnss3-tools && \
		curl -JLO "https://dl.filippo.io/mkcert/latest?for=linux/amd64" && \
		chmod +x mkcert-v*-linux-amd64 && sudo mv mkcert-v*-linux-amd64 /usr/local/bin/mkcert; \
	elif command -v yum >/dev/null 2>&1; then \
		sudo yum install -y nss-tools && \
		curl -JLO "https://dl.filippo.io/mkcert/latest?for=linux/amd64" && \
		chmod +x mkcert-v*-linux-amd64 && sudo mv mkcert-v*-linux-amd64 /usr/local/bin/mkcert; \
	else \
		echo "Please install mkcert manually: https://github.com/FiloSottile/mkcert#installation"; \
	fi
	@echo "Installing pre-commit..."
	@if command -v pip3 >/dev/null 2>&1; then \
		pip3 install pre-commit; \
	elif command -v pip >/dev/null 2>&1; then \
		pip install pre-commit; \
	elif command -v brew >/dev/null 2>&1; then \
		brew install pre-commit; \
	else \
		echo "Please install pre-commit manually: https://pre-commit.com/#installation"; \
	fi
	@echo "Installing gosec for security scanning..."
	@if command -v go >/dev/null 2>&1; then \
		go install github.com/securecodewarrior/gosec/cmd/gosec@latest || \
		echo "gosec installation failed, security scanning will be skipped"; \
	else \
		echo "Go not found, skipping gosec installation"; \
	fi
	@echo "Installing golangci-lint..."
	@if command -v go >/dev/null 2>&1; then \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest || \
		echo "golangci-lint installation failed"; \
	fi
	@echo "Tool installation completed"

setup: install-tools
	@echo "Setting up Kainos development environment..."
	@echo "Configuring environment..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file from template"; \
		echo "Please update .env with your API keys and configuration"; \
	fi
	@echo "Generating SSL certificates..."
	@$(MAKE) cert
	@echo "Setting up git hooks..."
	@$(MAKE) setup-hooks
	@echo "Adding local domain to hosts file..."
	@if ! grep -q "api.kainos.local" /etc/hosts; then \
		echo "127.0.0.1 api.kainos.local kainos.local" | sudo tee -a /etc/hosts; \
	fi
	@echo "Building services..."
	@$(MAKE) build
	@echo "Setup completed successfully"
	@echo ""
	@echo "Next steps:"
	@echo "1. Update .env file with your API keys"
	@echo "2. Run 'make dev' to start development environment"
	@echo "3. Run 'make test' to verify everything works"

build:
	@echo "Building all services..."
	@$(DOCKER_COMPOSE_DEV) build --no-cache

dev:
	@echo "Starting development environment..."
	@$(DOCKER_COMPOSE_DEV) up -d
	@echo "Services starting..."
	@sleep 10
	@echo "Development environment ready:"
	@echo "  Gateway: https://localhost:9443"
	@echo "  Domain:  https://api.kainos.local"
	@echo "  Core:    https://localhost:8443"
	@echo "  Email:   https://localhost:8444"

prod:
	@echo "Starting production environment..."
	@$(DOCKER_COMPOSE_PROD) up -d

test:
	@echo "Running email functionality tests..."
	@./test-email-functionality.sh

cert:
	@echo "Generating SSL certificates..."
	@mkdir -p certs
	@mkcert -install
	@mkcert -cert-file certs/localhost+6.pem -key-file certs/localhost+6-key.pem localhost 127.0.0.1 ::1 core-api email-service kainos.local api.kainos.local
	@mkcert -cert-file certs/core-api.pem -key-file certs/core-api-key.pem localhost 127.0.0.1 core-api kainos.local api.kainos.local
	@mkcert -cert-file certs/email-service.pem -key-file certs/email-service-key.pem localhost 127.0.0.1 email-service kainos.local email.kainos.local
	@echo "SSL certificates generated in ./certs/"

logs:
	@echo "Service logs:"
	@echo "============="
	@echo "Caddy Gateway:"
	@$(DOCKER_COMPOSE_DEV) logs --tail=20 caddy
	@echo ""
	@echo "Core API:"
	@$(DOCKER_COMPOSE_DEV) logs --tail=20 core-api
	@echo ""
	@echo "Email Service:"
	@$(DOCKER_COMPOSE_DEV) logs --tail=20 email-service

status:
	@echo "Service Status:"
	@echo "==============="
	@$(DOCKER_COMPOSE_DEV) ps

stop:
	@echo "Stopping all services..."
	@$(DOCKER_COMPOSE_DEV) down

restart:
	@echo "Restarting all services..."
	@$(DOCKER_COMPOSE_DEV) restart
	@sleep 5
	@$(MAKE) status

clean:
	@echo "Cleaning up containers and volumes..."
	@$(DOCKER_COMPOSE_DEV) down -v --remove-orphans
	@docker system prune -f
	@echo "Cleanup complete"

deploy:
	@echo "Deploying to production..."
	@$(DOCKER_COMPOSE_PROD) pull
	@$(DOCKER_COMPOSE_PROD) up -d --force-recreate
	@echo "Production deployment complete"

setup-hooks:
	@echo "Setting up git hooks..."
	@chmod +x scripts/setup-hooks.sh
	@chmod +x scripts/validate-commit.sh
	@chmod +x .githooks/pre-commit
	@git config core.hooksPath .githooks
	@pre-commit install
	@echo "Git hooks configured"

lint:
	@echo "Running linting and formatting..."
	@cd core && go fmt ./... && go vet ./...
	@cd email && go fmt ./... && go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd core && golangci-lint run; \
		cd email && golangci-lint run; \
	fi
	@echo "Linting completed"

validate:
	@echo "Running comprehensive validation..."
	@./scripts/validate-commit.sh

check-env:
	@echo "Checking environment configuration..."
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found"; \
		echo "Run 'make setup' to create it"; \
		exit 1; \
	fi
	@echo "Environment configuration OK"

check-tools:
	@echo "Checking required tools..."
	@command -v docker >/dev/null 2>&1 || { echo "Docker not found"; exit 1; }
	@command -v docker-compose >/dev/null 2>&1 || { echo "Docker Compose not found"; exit 1; }
	@command -v mkcert >/dev/null 2>&1 || { echo "mkcert not found"; exit 1; }
	@command -v pre-commit >/dev/null 2>&1 || { echo "pre-commit not found"; exit 1; }
	@echo "All required tools are installed"

health:
	@echo "Checking service health..."
	@curl -k -s https://localhost:9443/health || echo "Gateway not responding"
	@curl -k -s https://localhost:8443/healthz || echo "Core API not responding"
	@curl -k -s https://localhost:8444/healthz || echo "Email service not responding"

quick-start: check-tools check-env cert dev
	@echo "Quick start completed"
	@echo "Services are running at:"
	@echo "  Gateway: https://localhost:9443"
	@echo "  Domain:  https://api.kainos.local"
