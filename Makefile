.PHONY: gen gen-be gen-fe run start stop help

# Default target
.DEFAULT_GOAL := help

## gen: Generate code for both backend and frontend from OpenAPI spec
gen: gen-be gen-fe
	@echo ""
	@echo "✅ All code generation completed!"

## gen-be: Generate Go code for backend
gen-be:
	@echo "📦 Generating backend code..."
	@cd crypto-trading-connector-be && make gen

## gen-fe: Generate TypeScript types for frontend
gen-fe:
	@echo "📦 Generating frontend types..."
	@cd crypto-trading-connector-fe && make gen

## run: Start both backend and frontend servers
run:
	@echo "🚀 Starting crypto trading connector..."
	@echo "Backend will be available at: http://10.221.1.162:8090"
	@echo "Frontend will be available at: http://10.221.1.162:3000"
	@echo ""
	@./scripts/start-services.sh

## start: Start services as systemd daemons (requires sudo)
start:
	@echo "🔧 Starting crypto trading connector as systemd services..."
	@sudo systemctl start crypto-trading-be
	@sudo systemctl start crypto-trading-fe
	@echo "✅ Services started!"
	@echo "Check status with: make status"

## stop: Stop systemd services (requires sudo)
stop:
	@echo "🛑 Stopping crypto trading connector services..."
	@sudo systemctl stop crypto-trading-fe
	@sudo systemctl stop crypto-trading-be
	@echo "✅ Services stopped!"

## status: Check status of systemd services
status:
	@echo "📊 Service Status:"
	@echo ""
	@echo "Backend Service:"
	@sudo systemctl status crypto-trading-be --no-pager -l
	@echo ""
	@echo "Frontend Service:"
	@sudo systemctl status crypto-trading-fe --no-pager -l

## install-services: Install systemd service files (requires sudo)
install-services:
	@echo "📦 Installing systemd service files..."
	@sudo cp systemd/crypto-trading-be.service /etc/systemd/system/
	@sudo cp systemd/crypto-trading-fe.service /etc/systemd/system/
	@sudo systemctl daemon-reload
	@sudo systemctl enable crypto-trading-be
	@sudo systemctl enable crypto-trading-fe
	@echo "✅ Services installed and enabled!"
	@echo "Start with: make start"

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
	@echo ""
	@echo "Examples:"
	@echo "  make gen              # Generate code for both backend and frontend"
	@echo "  make run              # Start both services (development mode)"
	@echo "  make install-services # Install systemd services (requires sudo)"
	@echo "  make start            # Start systemd services (requires sudo)"
	@echo "  make stop             # Stop systemd services (requires sudo)"
	@echo "  make status           # Check service status"
