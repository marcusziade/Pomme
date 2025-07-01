#!/bin/bash
# Simple local server for Pomme landing page

echo "🍎 Starting Pomme Landing Page Preview..."
echo "────────────────────────────────────────"

# Check if Python 3 is installed
if ! command -v python3 &> /dev/null; then
    echo "❌ Python 3 is required but not installed."
    echo "   Please install Python 3 to use this preview server."
    exit 1
fi

# Run the Python server
cd "$(dirname "$0")"
python3 serve.py