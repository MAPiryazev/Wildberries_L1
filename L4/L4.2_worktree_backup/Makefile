.PHONY: build run-local run-docker test benchmark clean help

BINARY_NAME=mycut
DOCKER_IMAGE=mycut:latest
GO=go

build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build -o $(BINARY_NAME) ./cmd

build-docker:
	@echo "Building Docker image..."
	docker build -f Dockerfile.worker -t $(DOCKER_IMAGE) .

run-local: build
	@echo "Running in local mode..."
	echo "field1,field2,field3\nvalue1,value2,value3\ntest1,test2,test3" | \
		./$(BINARY_NAME) -f 1,3 -d "," --mode=local

run-docker: build-docker
	@echo "Starting Docker Compose..."
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	sleep 3
	@echo "Starting coordinator..."
	echo "field1,field2,field3\nvalue1,value2,value3\ntest1,test2,test3" | \
		./$(BINARY_NAME) -f 1,3 -d "," --mode=coordinator

run-docker-logs:
	docker-compose logs -f

run-docker-stop:
	docker-compose down

test:
	@echo "Running tests..."
	$(GO) test -v -cover ./...

test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

benchmark: build
	@echo "Running benchmark..."
	@bash scripts/benchmark.sh

lint:
	@echo "Linting code..."
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME) coverage.out coverage.html
	docker-compose down -v 2>/dev/null || true

help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  build-docker   - Build Docker image"
	@echo "  run-local      - Run in local mode"
	@echo "  run-docker     - Run with Docker Compose"
	@echo "  run-docker-logs - Show Docker logs"
	@echo "  run-docker-stop - Stop Docker containers"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  benchmark      - Run performance benchmark"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  clean          - Clean up"
	@echo "  help           - Show this help"
