## include .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Tools
OAPI_CODEGEN := $(shell go env GOPATH)/bin/oapi-codegen
GOLANGCI_LINT := $(shell go env GOPATH)/bin/golangci-lint

.PHONY: install-tools
install-tools:
	@echo "Installing tools from tools.go..."
	cd tools && go mod download
	# Check if oapi-codegen is already installed
	if [ ! -f "$(OAPI_CODEGEN)" ]; then \
		echo "Installing oapi-codegen..."; \
		go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest; \
	else \
		echo "oapi-codegen is already installed."; \
	fi
	# Check if golangci-lint is already installed
	if [ ! -f "$(GOLANGCI_LINT)" ]; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	else \
		echo "golangci-lint is already installed."; \
	fi

$(OAPI_CODEGEN) $(GOLANGCI_LINT): install-tools

# Generate API code
.PHONY: generate-api
generate-api: $(OAPI_CODEGEN)
	@echo "Generating API from OpenAPI spec..."
	$(OAPI_CODEGEN) -package api -generate types spec/weather-service.yaml > internal/api/types.go
	$(OAPI_CODEGEN) -package api -generate echo-server -o internal/api/server.go spec/weather-service.yaml	
	
# Lint
.PHONY: lint
lint: $(GOLANGCI_LINT)
	gofmt -w .
	$(GOLANGCI_LINT) run --fix	

# Build the application
.PHONY: build
build:
	@echo "Building the application..."
	go build -o api ./cmd/api

# Run api
.PHONY: run-api
run-api:
	@echo "Running the application..."
	go run ./cmd/api/main.go

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...	

# Default target
all: install-tools generate-api lint build test

# Set default make
.DEFAULT_GOAL := all