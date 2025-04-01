#!/bin/bash
set -e

# Change to project root directory
cd "$(dirname "$0")/.."

# Ensure dependencies are updated
go mod tidy

# Build the tool
echo "Building Pomme..."
go build -o bin/pomme ./cmd/pomme

# Run the tool with passed arguments
echo "Running Pomme..."
./bin/pomme "$@"