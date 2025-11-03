#!/bin/bash

echo "Setting up Kainos Git Hooks..."

if ! command -v pre-commit &> /dev/null; then
    echo "Installing pre-commit..."

    if command -v pip &> /dev/null; then
        pip install pre-commit
    elif command -v pip3 &> /dev/null; then
        pip3 install pre-commit
    elif command -v brew &> /dev/null; then
        brew install pre-commit
    else
        echo "Could not install pre-commit automatically"
        echo "Please install pre-commit manually: pip install pre-commit"
        exit 1
    fi
fi

echo "Installing pre-commit hooks..."
pre-commit install

echo "Setting up git hooks directory..."
git config core.hooksPath .githooks
chmod +x .githooks/pre-commit

echo "Checking gosec installation..."
if command -v gosec &> /dev/null; then
    echo "gosec is already installed"
else
    echo "gosec not found, security scanning will be skipped"
fi

echo "Testing pre-commit setup..."
pre-commit run --all-files || echo "Some hooks failed, but setup is complete"

echo "Git hooks setup completed"
echo ""
echo "Available commands:"
echo "  make help           - Show all available commands"
echo "  pre-commit run      - Run all hooks manually"
echo "  pre-commit autoupdate - Update hook versions"
echo ""
echo "Hooks will now run automatically on git commit"
