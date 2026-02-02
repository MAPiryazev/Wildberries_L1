.PHONY: help lint lint-install lint-fix build run docker-up docker-down migrate clean test

help:
	@echo "Available commands:"
	@echo "  make lint           - Run linter"
	@echo "  make lint-install   - Install golangci-lint"
	@echo "  make lint-fix       - Fix linting issues automatically"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make docker-up      - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make test           - Run tests"

lint:
	golangci-lint run ./...

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint-fix:
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found. Run 'make lint-install' first"; \
		exit 1; \
	fi
	golangci-lint run ./... --fix

build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

docker-up:
	docker-compose up -d


docker-down:
	docker-compose down


clean:
	rm -rf bin/
	go clean

test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

deps:
	go mod tidy
	go mod download

.DEFAULT_GOAL := help
