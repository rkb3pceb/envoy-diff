# Makefile for envoy-diff

.PHONY: all build test clean install lint fmt vet coverage help

# Binary name
BINARY_NAME=envoy-diff
BINARY_PATH=./bin/$(BINARY_NAME)

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-s -w"

all: test build ## Run tests and build binary

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) ./cmd/envoy-diff
	@echo "Binary built: $(BINARY_PATH)"

test: ## Run all tests
	@echo "Running tests..."
	$(GOTEST) -v -race ./...

coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

fmt: ## Format Go code
	@echo "Formatting code..."
	$(GOFMT) ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

lint: fmt vet ## Run formatting and vetting

install: ## Install the binary
	@echo "Installing $(BINARY_NAME)..."
	$(GOINSTALL) ./cmd/envoy-diff

clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/
	rm -f coverage.out coverage.html

run: build ## Build and run example
	@echo "Running example..."
	$(BINARY_PATH) examples/old.env examples/new.env

run-audit: build ## Build and run with audit mode
	@echo "Running with audit mode..."
	$(BINARY_PATH) --audit examples/old.env examples/new.env

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
