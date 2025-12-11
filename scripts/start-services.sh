#!/bin/bash

# Crypto Trading Connector - Service Startup Script
# This script starts both backend and frontend services

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BE_DIR="$PROJECT_ROOT/crypto-trading-connector-be"
FE_DIR="$PROJECT_ROOT/crypto-trading-connector-fe"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Function to check if port is in use
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        return 0  # Port is in use
    else
        return 1  # Port is free
    fi
}

# Function to kill process on port
kill_port() {
    local port=$1
    local pid=$(lsof -Pi :$port -sTCP:LISTEN -t 2>/dev/null)
    if [ ! -z "$pid" ]; then
        log "Killing process $pid on port $port"
        kill -9 $pid 2>/dev/null || true
        sleep 2
    fi
}

# Function to start backend
start_backend() {
    log "Starting backend service..."
    
    cd "$BE_DIR"
    
    # Check if backend is already running
    if check_port 8090; then
        warn "Backend port 8090 is already in use"
        kill_port 8090
    fi
    
    # Build backend
    log "Building backend..."
    if ! make build; then
        error "Failed to build backend"
        exit 1
    fi
    
    # Start backend in background
    log "Starting backend server on port 8090..."
    nohup ./bin/server > ../logs/backend.log 2>&1 &
    echo $! > ../logs/backend.pid
    
    # Wait for backend to start
    sleep 3
    if check_port 8090; then
        success "Backend started successfully on port 8090"
    else
        error "Failed to start backend"
        exit 1
    fi
}

# Function to start frontend
start_frontend() {
    log "Starting frontend service..."
    
    cd "$FE_DIR"
    
    # Check if frontend is already running
    if check_port 3000; then
        warn "Frontend port 3000 is already in use"
        kill_port 3000
    fi
    
    # Install dependencies if needed
    if [ ! -d "node_modules" ]; then
        log "Installing frontend dependencies..."
        npm install
    fi
    
    # Build frontend
    log "Building frontend..."
    if ! npm run build; then
        error "Failed to build frontend"
        exit 1
    fi
    
    # Start frontend in background
    log "Starting frontend server on port 3000..."
    nohup npm run preview -- --host 0.0.0.0 --port 3000 > ../logs/frontend.log 2>&1 &
    echo $! > ../logs/frontend.pid
    
    # Wait for frontend to start
    sleep 5
    if check_port 3000; then
        success "Frontend started successfully on port 3000"
    else
        error "Failed to start frontend"
        exit 1
    fi
}

# Main execution
main() {
    log "Starting Crypto Trading Connector services..."
    
    # Create logs directory
    mkdir -p "$PROJECT_ROOT/logs"
    
    # Update environment files for Raspberry Pi
    log "Updating environment configuration for Raspberry Pi (10.221.1.162)..."
    
    # Update backend CORS settings
    if [ -f "$BE_DIR/.env" ]; then
        sed -i 's/CORS_ALLOWED_ORIGINS=.*/CORS_ALLOWED_ORIGINS=http:\/\/10.221.1.162:3000,http:\/\/localhost:3000/' "$BE_DIR/.env"
    fi
    
    # Update frontend API URL
    if [ -f "$FE_DIR/.env" ]; then
        sed -i 's/API_BASE_URL=.*/API_BASE_URL=http:\/\/10.221.1.162:8090/' "$FE_DIR/.env"
    fi
    
    # Start services
    start_backend
    start_frontend
    
    success "All services started successfully!"
    echo ""
    log "Service URLs:"
    echo "  Backend API:  http://10.221.1.162:8090"
    echo "  Frontend Web: http://10.221.1.162:3000"
    echo ""
    log "Log files:"
    echo "  Backend:  $PROJECT_ROOT/logs/backend.log"
    echo "  Frontend: $PROJECT_ROOT/logs/frontend.log"
    echo ""
    log "To stop services, run: $PROJECT_ROOT/scripts/stop-services.sh"
}

# Run main function
main "$@"