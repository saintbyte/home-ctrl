# Makefile for home-ctrl project

# Configuration
BIN_DIR := bin
BIN_NAME := home-ctrl
VERSION := 0.1.0

# Build targets
.PHONY: all build clean test run migrate

all: build

# Build for current platform
build:
	@echo "Building home-ctrl..."
	@mkdir -p $(BIN_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BIN_NAME)-linux-amd64 ./cmd/home-ctrl
	@GOOS=linux GOARCH=386 go build -o $(BIN_DIR)/$(BIN_NAME)-linux-386 ./cmd/home-ctrl
	@echo "Build complete! Binaries available in $(BIN_DIR)/"

# Build for specific architecture
build-linux-386:
	@echo "Building for Linux i386..."
	@mkdir -p $(BIN_DIR)
	@GOOS=linux GOARCH=386 go build -o $(BIN_DIR)/$(BIN_NAME)-linux-386 ./cmd/home-ctrl
	@echo "Build complete: $(BIN_DIR)/$(BIN_NAME)-linux-386"

build-linux-amd64:
	@echo "Building for Linux amd64..."
	@mkdir -p $(BIN_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BIN_NAME)-linux-amd64 ./cmd/home-ctrl
	@echo "Build complete: $(BIN_DIR)/$(BIN_NAME)-linux-amd64"

# Run database migrations
migrate:
	@echo "Running database migrations..."
	@mkdir -p $(BIN_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/home-ctrl-migrate ./cmd/home-ctrl
	@$(BIN_DIR)/home-ctrl-migrate migrate
	@echo "Migrations complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run the application
run:
	@echo "Running home-ctrl..."
	@go run ./cmd/home-ctrl

# Show help
help:
	@echo "Makefile commands:"
	@echo "  make build          - Build for all supported platforms"
	@echo "  make build-linux-386 - Build for Linux i386"
	@echo "  make build-linux-amd64 - Build for Linux amd64"
	@echo "  make migrate        - Run database migrations"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make test           - Run tests"
	@echo "  make run            - Run the application"
	@echo "  make help           - Show this help message"