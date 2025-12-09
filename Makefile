.PHONY: gen gen-be gen-fe help

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

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
	@echo ""
	@echo "Examples:"
	@echo "  make gen-all    # Generate code for both backend and frontend"
	@echo "  make gen-be     # Generate backend code only"
	@echo "  make gen-fe     # Generate frontend types only"
