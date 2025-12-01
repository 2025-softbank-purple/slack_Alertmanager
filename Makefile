.PHONY: build test install clean lint fmt

# Build the node-watcher controller
build:
	go build -o bin/node-watcher ./cmd/node-watcher

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...
	gofmt -s -w .

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Install the controller to Kubernetes cluster
install: build
	kubectl apply -f charts/node-watcher/

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run integration tests (requires kind cluster)
test-integration:
	go test -v ./tests/integration/...

# Run unit tests only
test-unit:
	go test -v ./tests/unit/...

