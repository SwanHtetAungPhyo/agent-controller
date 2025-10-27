#!/usr/bin/env bash
set -e

echo "🚀 Building the core API Docker image..."
cd core
docker build -t core-service:latest -f Dockerfile .
cd ..

echo "📧 Building the email service Docker image..."
cd email
docker build -t email-service:latest -f Dockerfile .
cd ..

echo "✅ Available service images:"
docker images | grep service
