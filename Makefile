.PHONY: build build-all install clean test help

# Build variables
BINARY_NAME=gobig
VERSION=1.0.0
BUILD_DIR=bin
MAIN_PATH=./cmd/gobig

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -s -w"

# Default target
all: build

# Build for current platform
build:
	@echo "Building $(BINARY_NAME) for current platform..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Cross-compile for all platforms
build-all:
	@echo "Building $(BINARY_NAME) for all platforms..."
	@mkdir -p $(BUILD_DIR)

	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

	@echo "Building for Linux (arm64)..."
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)

	@echo "Building for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)

	@echo "Building for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)

	@echo "Building for Windows (amd64)..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

	@echo "Build complete! Binaries in $(BUILD_DIR)/"
	@ls -lh $(BUILD_DIR)

# Install to $GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) $(MAIN_PATH)
	@echo "Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build for current platform"
	@echo "  build-all  - Cross-compile for all platforms (Linux, macOS, Windows)"
	@echo "  install    - Install to \$$GOPATH/bin"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  help       - Show this help message"
