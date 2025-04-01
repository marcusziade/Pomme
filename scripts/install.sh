#!/bin/bash
set -e

# Change to project root directory
cd "$(dirname "$0")/.."

# Ensure dependencies are updated
go mod tidy

# Build the tool
echo "Building Pomme..."
go build -o bin/pomme ./cmd/pomme

# Install the binary
echo "Installing Pomme to /usr/local/bin/pomme..."
sudo cp bin/pomme /usr/local/bin/pomme

echo "Pomme has been installed successfully!"
echo "Run 'pomme --help' to get started."