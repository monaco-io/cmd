# Fanyi Makefile

.PHONY: all build install clean test run help deps

# Default target
all: build

# Build the binary
build:
	@echo "Building fanyi..."
	@go build -o fanyi .
	@echo "✓ Build complete: ./fanyi"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies installed"

# Install to system
install: build
	@echo "Installing fanyi to /usr/local/bin..."
	@sudo mv fanyi /usr/local/bin/
	@echo "✓ Installed successfully"

# Install with go install
install-go:
	@echo "Installing with go install..."
	@go install
	@echo "✓ Installed successfully"

# Run the application
run:
	@go run . $(ARGS)

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f fanyi
	@echo "✓ Clean complete"

# Setup configuration
setup:
	@echo "Setting up configuration..."
	@mkdir -p ~/.config/fanyi
	@if [ ! -f ~/.config/fanyi/config.yaml ]; then \
		cp config.example.yaml ~/.config/fanyi/config.yaml; \
		echo "✓ Config file created at ~/.config/fanyi/config.yaml"; \
		echo "  Please edit it and add your API key"; \
	else \
		echo "✓ Config file already exists"; \
	fi

# Show help
help:
	@echo "Fanyi - AI-Powered CLI Translation Tool"
	@echo ""
	@echo "Available targets:"
	@echo "  make build       - Build the binary"
	@echo "  make deps        - Install dependencies"
	@echo "  make install     - Build and install to /usr/local/bin"
	@echo "  make install-go  - Install using go install"
	@echo "  make run ARGS=.. - Run with arguments"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make setup       - Setup configuration file"
	@echo "  make help        - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make setup"
	@echo "  make run ARGS='hello world'"
	@echo "  make install"
