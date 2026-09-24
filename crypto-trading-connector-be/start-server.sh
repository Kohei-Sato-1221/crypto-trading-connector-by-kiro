#!/bin/bash

# Crypto Trading Backend Startup Script for systemd
# This script loads environment variables and starts the server

set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Always rebuild so the running binary matches the checked-out source.
# (Skipping the build when bin/server exists silently keeps stale code running.)
echo "Building server binary..."
go build -o bin/server ./cmd/server

# Load environment variables from .env file
if [ -f ".env" ]; then
    echo "Loading environment variables from .env..."
    set -a
    source .env
    set +a
    
    # Debug: Print loaded variables
    echo "Loaded environment variables:"
    echo "  DB_HOST = $DB_HOST"
    echo "  DB_PORT = $DB_PORT"
    echo "  DB_USER = $DB_USER"
    echo "  DB_NAME = $DB_NAME"
    echo "  SERVER_PORT = $SERVER_PORT"
else
    echo "Warning: .env file not found"
fi

# Start the server
echo "Starting crypto trading backend server..."
exec ./bin/server