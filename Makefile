
# Variables
BINARY_NAME=conjugationbot
BUILD_DIR=bin
MAIN_PATH=cmd/bot/main.go
GO=go
GO_BUILD_FLAGS=
GO_TEST_FLAGS=-v

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Watch for changes and automatically rebuild and run the application
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                $(GO) install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# Run tests
test:
	@echo "Running tests..."
	$(GO) test $(GO_TEST_FLAGS) ./...

# Target to run tests and generate coverage profile
coverage: 
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

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
	@echo "  make watch    - Watch for changes and automatically rebuild and run the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make format   - Format the code"
	@echo "  make lint     - Lint the code"
	@echo "  make help     - Show this help message"

.PHONY: all build run watch test coverage clean format lint help
