.PHONY: help test test-coverage test-short build fmt vet lint security clean install deps

# Default target
help:
	@echo "Available targets:"
	@echo "  make test          - Run all tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make test-short    - Run short tests"
	@echo "  make build         - Build the project"
	@echo "  make fmt           - Format code"
	@echo "  make vet           - Run go vet"
	@echo "  make lint          - Run golangci-lint (if installed)"
	@echo "  make security      - Run security scan with gosec"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make deps          - Download dependencies"
	@echo "  make ci            - Run all CI checks (fmt, vet, test)"

# Run all tests
test:
	go test -v ./...

# Run short tests (useful for quick checks)
test-short:
	go test -v -short ./...

# Build the project
build:
	go build -v ./...

# Format code
fmt:
	go fmt ./...
	@echo "Code formatted"

# Run go vet
vet:
	go vet ./...
	@echo "go vet passed"

# Run golangci-lint (install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	@command -v golangci-li

# Run security scan with gosec (install with: go install github.com/securego/gosec/v2/cmd/gosec@latest)
security:
	@command -v gosec >/dev/null 2>&1 || { echo "gosec not installed. Run: go install github.com/securego/gosec/v2/cmd/gosec@latest"; exit 1; }
	gosec ./...nt >/dev/null 2>&1 || { echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1; }
	golangci-lint run ./...

# Download dependencies
deps:
	go mod download
	go mod verify

# Tidy dependencies
tidy:
	go mod tidy

# Clean build artifacts
clean:
	go clean

# Run all CI checks locally
ci: fmt vet test
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "\n✅ All CI checks passed!"

# Install tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed"
