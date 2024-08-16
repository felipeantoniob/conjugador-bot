
# Variables
BINARY_NAME=conjugationbot
BUILD_DIR=bin
GO=go
GO_BUILD_FLAGS=
GO_TEST_FLAGS=-v

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	$(GO) test $(GO_TEST_FLAGS) ./...

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)

# Format code
format:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Show help
help:
	@echo "Makefile commands:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Build and run the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make format      - Format the code"
	@echo "  make lint     - Lint the code"
	@echo "  make help     - Show this help message"

.PHONY: all build run test clean format lint help
