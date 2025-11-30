.PHONY: help build run test test-coverage clean docker-build docker-run docker-stop lint fmt vet install-tools

# Variables
BINARY_NAME=api
MAIN_PATH=./cmd/api/main.go
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Default target
help:
	@echo "Available targets:"
	@echo "  build          - Build the application binary"
	@echo "  run            - Run the application"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Remove build artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-stop    - Stop Docker container"
	@echo "  lint           - Run golangci-lint"
	@echo "  fmt            - Format code with gofmt"
	@echo "  vet            - Run go vet"
	@echo "  install-tools  - Install development tools"

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_NAME)"

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	@go run $(MAIN_PATH)

# Run all tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"
	@go tool cover -func=$(COVERAGE_FILE)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@rm -f *.log *.db
	@echo "Clean complete"

# Docker build
docker-build:
	@echo "Building Docker image..."
	@docker build -t maze-solution-api:latest .
	@echo "Docker image built successfully"

# Docker run with compose
docker-run:
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "Containers started successfully"

# Docker stop
docker-stop:
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "Containers stopped successfully"

# Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted"

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "Vet complete"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed successfully"

# Run all quality checks
check: fmt vet lint test
	@echo "All checks passed!"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "Multi-platform build complete"
