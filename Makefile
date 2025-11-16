.PHONY: build install clean test run help

# Variables
BINARY_NAME=nmd
INSTALL_PATH=/usr/local/bin
GO=go

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build -o $(BINARY_NAME) ./cmd/notesmd
	@echo "Build complete! Binary: ./$(BINARY_NAME)"

# Install to system
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@sudo mv $(BINARY_NAME) $(INSTALL_PATH)/
	@echo "Installation complete! Run '$(BINARY_NAME)' to start."

# Install to user bin (no sudo required)
install-user: build
	@echo "Installing $(BINARY_NAME) to ~/bin..."
	@mkdir -p ~/bin
	@mv $(BINARY_NAME) ~/bin/
	@echo "Installation complete!"
	@echo "Make sure ~/bin is in your PATH. Add this to your ~/.bashrc or ~/.zshrc:"
	@echo 'export PATH="$$HOME/bin:$$PATH"'

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v -race -coverprofile=coverage.out -covermode=atomic ./cmd/notesmd/...

# Run tests without coverage (for quick checks)
test-quick:
	@echo "Running quick tests..."
	$(GO) test -v ./cmd/notesmd/...

# Run the application
run: build
	@./$(BINARY_NAME)

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run || echo "Install golangci-lint: https://golangci-lint.run/usage/install/"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Help
help:
	@echo "NotesMD Makefile commands:"
	@echo ""
	@echo "  make build        - Build the binary"
	@echo "  make install      - Install to /usr/local/bin (requires sudo)"
	@echo "  make install-user - Install to ~/bin (no sudo)"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make run          - Build and run"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Lint code"
	@echo "  make deps         - Download dependencies"
	@echo "  make help         - Show this help"
