.PHONY: build test test-coverage run clean help lint install-deps

BUILD_DIR=./bin
BINARY_NAME=mycut
MAIN_PATH=./cmd
INTERNAL_PATH=./internal

help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  build          - Build the binary"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  run            - Run with example data"
	@echo "  clean          - Remove build artifacts"
	@echo "  install-deps   - Install Go dependencies"
	@echo "  lint           - Run linter"

install-deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

build: install-deps
	@echo "Building $(BINARY_NAME)..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✓ Built: $(BUILD_DIR)/$(BINARY_NAME)"

test: install-deps
	@echo "Running tests..."
	go test -v -race ./...
	@echo "✓ All tests passed"

run: build
	@echo "Running mycut..."
	echo "name,age,city\nAlice,30,NYC\nBob,25,London" | \
		$(BUILD_DIR)/$(BINARY_NAME) -d "," -f "1,3"

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR) coverage.out coverage.html
	go clean
	@echo "✓ Cleaned"

lint:
	go fmt ./...
	go vet ./...
