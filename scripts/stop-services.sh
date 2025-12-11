#!/bin/bash

# Crypto Trading Connector - Service Stop Script
# This script stops both backend and frontend services

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

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

# Function to stop service by PID file
stop_service() {
    local service_name=$1
    local pid_file="$PROJECT_ROOT/logs/${service_name}.pid"
    
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if kill -0 "$pid" 2>/dev/null; then
            log "Stopping $service_name (PID: $pid)..."
            kill -TERM "$pid" 2>/dev/null || true
            
            # Wait for graceful shutdown
            local count=0
            while kill -0 "$pid" 2>/dev/null && [ $count -lt 10 ]; do
                sleep 1
                count=$((count + 1))
            done
            
            # Force kill if still running
            if kill -0 "$pid" 2>/dev/null; then
                warn "Force killing $service_name..."
                kill -9 "$pid" 2>/dev/null || true
            fi
            
            success "$service_name stopped"
        else
            warn "$service_name was not running"
        fi
        rm -f "$pid_file"
    else
        warn "No PID file found for $service_name"
    fi
}

# Function to kill processes on specific ports
kill_port_processes() {
    log "Checking for processes on ports 8090 and 3000..."
    
    # Kill backend (port 8090)
    local backend_pid=$(lsof -Pi :8090 -sTCP:LISTEN -t 2>/dev/null || true)
    if [ ! -z "$backend_pid" ]; then
        log "Killing backend process on port 8090 (PID: $backend_pid)"
        kill -9 "$backend_pid" 2>/dev/null || true
    fi
    
    # Kill frontend (port 3000)
    local frontend_pid=$(lsof -Pi :3000 -sTCP:LISTEN -t 2>/dev/null || true)
    if [ ! -z "$frontend_pid" ]; then
        log "Killing frontend process on port 3000 (PID: $frontend_pid)"
        kill -9 "$frontend_pid" 2>/dev/null || true
    fi
}

# Main execution
main() {
    log "Stopping Crypto Trading Connector services..."
    
    # Stop services by PID files
    stop_service "backend"
    stop_service "frontend"
    
    # Kill any remaining processes on the ports
    kill_port_processes
    
    success "All services stopped successfully!"
}

# Run main function
main "$@"