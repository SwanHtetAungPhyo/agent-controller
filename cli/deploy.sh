#!/bin/bash

echo "Building Financial Agent CLI for Production"
echo "==========================================="
echo ""

# Clean previous build
echo "Cleaning previous build..."
rm -rf dist
echo "✓ Cleaned"
echo ""

# Install dependencies
echo "Installing dependencies..."
pnpm install --prod=false
echo "✓ Dependencies installed"
echo ""

# Build
echo "Building TypeScript..."
pnpm build
echo "✓ Build complete"
echo ""

# Make executable
echo "Making CLI executable..."
chmod +x dist/index.js
echo "✓ CLI is executable"
echo ""

echo "Production build complete!"
echo ""
echo "To use the CLI:"
echo "  node dist/index.js <command>"
echo ""
echo "To install globally:"
echo "  pnpm link --global"
echo "  financial-agent <command>"
echo ""
