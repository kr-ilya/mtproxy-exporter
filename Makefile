.PHONY: all build test test-coverage clean run docker-build docker-run lint

# Binary name
BINARY_NAME=mtproxy-exporter
DOCKER_IMAGE=mtproxy-exporter:latest

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet
GOFMT=gofmt

# Build directory
BUILD_DIR=./bin

all: test build

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/mtproxy-exporter

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

run: build
	@echo "Running..."
	$(BUILD_DIR)/$(BINARY_NAME)

lint:
	@echo "Running linters..."
	$(GOVET) ./...
	$(GOFMT) -s -w .

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	@echo "Running in Docker..."
	docker run -p 9330:9330 -e MTPROXY_URL=http://127.0.0.1:8888 $(DOCKER_IMAGE)

# Development helpers
dev: clean test build run

.DEFAULT_GOAL := build
