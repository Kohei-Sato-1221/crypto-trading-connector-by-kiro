.PHONY: run test fmt help e2e-test unit-test get-balance buy-order

# Default target
.DEFAULT_GOAL := help

# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

## run: Start the backend server locally
run:
	@echo "Starting backend server..."
	go run cmd/server/main.go

## test: Run all tests (unit tests only, excluding e2e)
test:
	@echo "Running unit tests..."
	go test -v ./...

## e2e-test: Run E2E tests (requires running server and actual DB/API)
e2e-test:
	@echo "Running E2E tests..."
	go test -v -tags=e2e ./test/e2e/...

## fmt: Format all Go files
fmt:
	@echo "Formatting Go files..."
	go fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

## lint: Run golangci-lint (if installed)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping..."; \
	fi

## tidy: Tidy go modules
tidy:
	@echo "Tidying go modules..."
	go mod tidy

## build: Build the binary
build:
	@echo "Building binary..."
	go build -o bin/server cmd/server/main.go

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

## curl: Test API endpoints with curl (usage: make curl a=market|bitcoin|ethereum|bitcoin-chart)
curl:
	@if [ -z "$(a)" ]; then \
		echo "Error: Please specify an API endpoint with a=<endpoint>"; \
		echo ""; \
		echo "Available endpoints:"; \
		echo "  market         - GET /api/v1/crypto/market"; \
		echo "  bitcoin        - GET /api/v1/crypto/bitcoin"; \
		echo "  ethereum       - GET /api/v1/crypto/ethereum"; \
		echo "  bitcoin-chart  - GET /api/v1/crypto/bitcoin/chart?period=7d"; \
		echo "  ethereum-chart - GET /api/v1/crypto/ethereum/chart?period=7d"; \
		echo ""; \
		echo "Example: make curl a=market"; \
		exit 1; \
	fi
	@if [ "$(a)" = "market" ]; then \
		echo "Testing GET /api/v1/crypto/market"; \
		curl -X GET "http://localhost:8080/api/v1/crypto/market" \
			-H "Content-Type: application/json" \
			-w "\n\nHTTP Status: %{http_code}\n" \
			| jq . || curl -X GET "http://localhost:8080/api/v1/crypto/market" -w "\n\nHTTP Status: %{http_code}\n"; \
	elif [ "$(a)" = "bitcoin" ]; then \
		echo "Testing GET /api/v1/crypto/bitcoin"; \
		curl -X GET "http://localhost:8080/api/v1/crypto/bitcoin" \
			-H "Content-Type: application/json" \
			-w "\n\nHTTP Status: %{http_code}\n" \
			| jq . || curl -X GET "http://localhost:8080/api/v1/crypto/bitcoin" -w "\n\nHTTP Status: %{http_code}\n"; \
	elif [ "$(a)" = "ethereum" ]; then \
		echo "Testing GET /api/v1/crypto/ethereum"; \
		curl -X GET "http://localhost:8080/api/v1/crypto/ethereum" \
			-H "Content-Type: application/json" \
			-w "\n\nHTTP Status: %{http_code}\n" \
			| jq . || curl -X GET "http://localhost:8080/api/v1/crypto/ethereum" -w "\n\nHTTP Status: %{http_code}\n"; \
	elif [ "$(a)" = "bitcoin-chart" ]; then \
		echo "Testing GET /api/v1/crypto/bitcoin/chart?period=7d"; \
		curl -X GET "http://localhost:8080/api/v1/crypto/bitcoin/chart?period=7d" \
			-H "Content-Type: application/json" \
			-w "\n\nHTTP Status: %{http_code}\n" \
			| jq . || curl -X GET "http://localhost:8080/api/v1/crypto/bitcoin/chart?period=7d" -w "\n\nHTTP Status: %{http_code}\n"; \
	elif [ "$(a)" = "ethereum-chart" ]; then \
		echo "Testing GET /api/v1/crypto/ethereum/chart?period=7d"; \
		curl -X GET "http://localhost:8080/api/v1/crypto/ethereum/chart?period=7d" \
			-H "Content-Type: application/json" \
			-w "\n\nHTTP Status: %{http_code}\n" \
			| jq . || curl -X GET "http://localhost:8080/api/v1/crypto/ethereum/chart?period=7d" -w "\n\nHTTP Status: %{http_code}\n"; \
	else \
		echo "Error: Unknown endpoint '$(a)'"; \
		echo "Available endpoints: market, bitcoin, ethereum, bitcoin-chart, ethereum-chart"; \
		exit 1; \
	fi

## gen: Generate Go code from OpenAPI spec
gen:
	@echo "Generating Go code from OpenAPI spec..."
	@mkdir -p internal/generated
	@oapi-codegen -package generated -generate types ../openapi.yaml > internal/generated/models.go
	@echo "Generated: internal/generated/models.go"
	@go fmt ./internal/generated/...

## get-balance: Fetch JPY balance from bitFlyer API
get-balance:
	@echo "Fetching balance from bitFlyer API..."
	@go run cmd/get-balance/main.go

## buy-order: Place buy orders for BTC and ETH at 97% of current price
buy-order:
	@echo "Placing buy orders for BTC and ETH..."
	@go run cmd/buy-order/main.go

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
	@echo ""
	@echo "Curl API test endpoints (usage: make curl a=<endpoint>):"
	@echo "  market          - GET /api/v1/crypto/market"
	@echo "  bitcoin         - GET /api/v1/crypto/bitcoin"
	@echo "  ethereum        - GET /api/v1/crypto/ethereum"
	@echo "  bitcoin-chart   - GET /api/v1/crypto/bitcoin/chart?period=7d"
	@echo "  ethereum-chart  - GET /api/v1/crypto/ethereum/chart?period=7d"
	@echo ""
	@echo "bitFlyer API commands:"
	@echo "  make get-balance - Fetch JPY balance from bitFlyer API"
	@echo "  make buy-order   - Place buy orders for BTC and ETH at 97% of current price"
	@echo ""
	@echo "Example: make curl a=market"
