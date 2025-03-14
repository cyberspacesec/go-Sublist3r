VERSION := 1.0.0
BINARY_NAME := go-sublist3r
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

.PHONY: build install clean test api-server

build:
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Build complete"

install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin/..."
	@cp $(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete"

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "Clean complete"

test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests complete"

api-server: build
	@echo "Starting API server on port 8080..."
	@./$(BINARY_NAME) api --port 8080
