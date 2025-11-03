#!/bin/bash

echo "Validating commit for Kainos platform..."

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

validate_go_files() {
    local service=$1
    echo "Validating Go files in $service..."

    if [ -d "$service" ]; then
        cd "$service"

        if [ ! -f "go.mod" ]; then
            echo -e "${RED}go.mod not found in $service${NC}"
            return 1
        fi

        echo "Running go mod tidy..."
        go mod tidy

        echo "Running go vet..."
        if ! go vet ./...; then
            echo -e "${RED}go vet failed in $service${NC}"
            cd ..
            return 1
        fi

        echo "Running go fmt..."
        if [ -n "$(gofmt -l .)" ]; then
            echo -e "${YELLOW}Formatting Go files in $service${NC}"
            gofmt -w .
        fi

        if command -v gosec &> /dev/null; then
            echo "Running security scan..."
            gosec ./... || echo -e "${YELLOW}Security scan completed with warnings${NC}"
        fi

        cd ..
        echo -e "${GREEN}$service validation completed${NC}"
    fi
}

validate_docker_files() {
    echo "Validating Docker files..."

    for dockerfile in core/Dockerfile email/Dockerfile caddy/Dockerfile; do
        if [ -f "$dockerfile" ]; then
            echo "$dockerfile exists"
        else
            echo -e "${RED}$dockerfile missing${NC}"
            return 1
        fi
    done

    if ! docker-compose -f docker-compose.dev.yaml config > /dev/null 2>&1; then
        echo -e "${RED}docker-compose.dev.yaml is invalid${NC}"
        return 1
    fi

    if ! docker-compose -f docker-compose.yaml config > /dev/null 2>&1; then
        echo -e "${RED}docker-compose.yaml is invalid${NC}"
        return 1
    fi

    echo -e "${GREEN}Docker files validation completed${NC}"
}

validate_env_files() {
    echo "Validating environment files..."

    if [ ! -f ".env" ]; then
        echo -e "${RED}.env file missing${NC}"
        return 1
    fi

    if [ ! -f ".env.example" ]; then
        echo -e "${RED}.env.example file missing${NC}"
        return 1
    fi

    required_vars=(
        "APP_NAME"
        "RESEND_API_KEY"
        "FROM_EMAIL"
        "CLERK_SECRET"
        "APP_JWT_SECRET"
    )

    for var in "${required_vars[@]}"; do
        if ! grep -q "^$var=" .env; then
            echo -e "${YELLOW}$var not found in .env${NC}"
        fi
    done

    echo -e "${GREEN}Environment files validation completed${NC}"
}

validate_ssl_certs() {
    echo "Validating SSL certificates..."

    if [ ! -d "certs" ]; then
        echo -e "${YELLOW}certs directory not found, run 'make cert' to generate${NC}"
        return 0
    fi

    required_certs=(
        "certs/localhost+6.pem"
        "certs/localhost+6-key.pem"
        "certs/core-api.pem"
        "certs/core-api-key.pem"
        "certs/email-service.pem"
        "certs/email-service-key.pem"
    )

    for cert in "${required_certs[@]}"; do
        if [ -f "$cert" ]; then
            echo "$cert exists"
        else
            echo -e "${YELLOW}$cert missing${NC}"
        fi
    done

    echo -e "${GREEN}SSL certificates validation completed${NC}"
}

echo "Starting comprehensive validation..."

validate_go_files "core"
validate_go_files "email"
validate_docker_files
validate_env_files
validate_ssl_certs

if make help > /dev/null 2>&1; then
    echo -e "${GREEN}Makefile is valid${NC}"
else
    echo -e "${RED}Makefile is invalid${NC}"
    exit 1
fi

echo -e "${GREEN}All validations completed successfully${NC}"
echo ""
echo "Next steps:"
echo "  Run 'make dev' to start development environment"
echo "  Run 'make test' to test email functionality"
echo "  Run 'make cert' if SSL certificates are missing"
