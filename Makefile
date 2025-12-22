.PHONY: help test test-verbose test-coverage lint fmt vet build clean install deps example

# Default target
help:
	@echo "Available targets:"
	@echo "  make test           - Run all tests"
	@echo "  make test-verbose   - Run tests with verbose output"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make lint           - Run linters"
	@echo "  make fmt            - Format code"
	@echo "  make vet            - Run go vet"
	@echo "  make build          - Build the package"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make install        - Install dependencies"
	@echo "  make deps           - Download dependencies"
	@echo "  make example        - Run the basic example"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linters
lint:
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2" && exit 1)
	golangci-lint run --timeout 5m

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Build the package
build:
	@echo "Building..."
	go build ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f coverage.txt coverage.html
	go clean -cache -testcache

# Install dependencies
install: deps

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	go mod verify

# Run the basic example
example:
	@echo "Running basic example..."
	@if [ -z "$$CLOUDFLARE_ACCOUNT_ID" ] || [ -z "$$CLOUDFLARE_API_TOKEN" ] || [ -z "$$D1_DATABASE_ID" ]; then \
		echo "Error: Please set CLOUDFLARE_ACCOUNT_ID, CLOUDFLARE_API_TOKEN, and D1_DATABASE_ID environment variables"; \
		exit 1; \
	fi
	cd examples/basic && go run main.go

# Check for common issues
check: fmt vet lint test
	@echo "All checks passed!"

# CI target
ci: deps check test-coverage
	@echo "CI checks completed!"

