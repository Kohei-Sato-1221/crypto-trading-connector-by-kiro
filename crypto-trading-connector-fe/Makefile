.PHONY: dev build test lint fmt clean install help

# Default target
.DEFAULT_GOAL := help

# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

## run: Start the development server with API connection
run:
	@LOCAL_IP=$$(ipconfig getifaddr en0 2>/dev/null || hostname -I 2>/dev/null | awk '{print $$1}' || echo "10.221.1.199"); \
	echo "Starting development server with API connection..."; \
	echo "Backend API: http://$$LOCAL_IP:8080"; \
	echo "Frontend Local: http://localhost:3000"; \
	echo "Frontend Network: http://$$LOCAL_IP:3000"; \
	echo ""; \
	echo "Note: Make sure backend is running (cd ../crypto-trading-connector-be && make run)"; \
	echo "Note: Access from other devices: http://$$LOCAL_IP:3000/market"; \
	echo ""; \
	USE_MOCK_DATA=false API_BASE_URL=http://$$LOCAL_IP:8080 npm run dev -- --host 0.0.0.0

## run-mock: Start dev server with mock data
run-mock:
	@LOCAL_IP=$$(ipconfig getifaddr en0 2>/dev/null || hostname -I 2>/dev/null | awk '{print $$1}' || echo "10.221.1.199"); \
	echo "Starting development server with mock data..."; \
	echo "Frontend Local: http://localhost:3000"; \
	echo "Frontend Network: http://$$LOCAL_IP:3000"; \
	echo ""; \
	USE_MOCK_DATA=true npm run dev -- --host 0.0.0.0


## build: Build the application for production
build:
	@echo "Building application..."
	npm run build

## preview: Preview the production build
preview:
	@echo "Starting preview server..."
	npm run preview

## test: Run all tests
test:
	@echo "Running tests..."
	npm run test

## test-watch: Run tests in watch mode
test-watch:
	@echo "Running tests in watch mode..."
	npm run test:watch

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	npm run test:coverage

## lint: Run linter
lint:
	@echo "Running linter..."
	npm run lint

## fmt: Format code
fmt:
	@echo "Formatting code..."
	npm run format || echo "Format script not found, skipping..."

## install: Install dependencies
install:
	@echo "Installing dependencies..."
	npm install

## clean: Clean build artifacts and node_modules
clean:
	@echo "Cleaning..."
	rm -rf .nuxt .output dist node_modules

## clean-cache: Clean Nuxt cache only
clean-cache:
	@echo "Cleaning Nuxt cache..."
	rm -rf .nuxt .output

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
	@echo ""
	@echo "Environment variables:"
	@echo "  USE_MOCK_DATA   - Use mock data (true) or API (false)"
	@echo "  API_BASE_URL    - Backend API base URL (default: http://localhost:8080)"
	@echo ""
	@echo "Examples:"
	@echo "  make run              # Start dev server with API connection"
	@echo "  make run-mock         # Start dev server with mock data"
	@echo "  make test             # Run all tests"
	@echo "  make build            # Build for production"
