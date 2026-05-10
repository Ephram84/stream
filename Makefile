.PHONY: test test-short build fmt vet lint clean install deps vulncheck fix tidy

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

# Download dependencies
deps:
	go get ./...

# Tidy dependencies
tidy:
	go mod tidy

# Check for vulnerabilities
vulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

fix:
	go fix ./...

# Clean build artifacts
clean:
	go clean
