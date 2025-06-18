.PHONY: build run test clean tidy migrate

# Build the application
build:
	@echo "Building application..."
	go build -o bin/blog-management ./cmd/server

# Run the application
run:
	@echo "Starting application..."
	go run ./cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build files
clean:
	@echo "Cleaning up..."
	rm -rf bin/

# Download and tidy dependencies
tidy:
	@echo "Downloading dependencies..."
	go mod tidy

# Run database migrations
migrate:
	@echo "Running database migrations..."
	# Add your migration command here
	# Example: migrate -path ./migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet examines Go source code and reports suspicious constructs
vet:
	@echo "Running go vet..."
	go vet ./...

# Lint the code
lint:
	@echo "Running golangci-lint..."
	golangci-lint run

# Build for production
build-prod: tidy vet test
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/blog-management ./cmd/server

# Help command to display available targets
help:
	@echo "Available targets:"
	@echo "  build     - Build the application"
	@echo "  run       - Run the application"
	@echo "  test      - Run tests"
	@echo "  clean     - Remove build files"
	@echo "  tidy      - Download and tidy dependencies"
	@echo "  deps      - Install dependencies"
	@echo "  fmt       - Format code"
	@echo "  vet       - Run go vet"
	@echo "  lint      - Run linter"
	@echo "  build-prod- Build optimized binary for production"
	@echo "  help      - Show this help message"