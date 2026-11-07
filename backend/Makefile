# PayVue Backend - Makefile

# Variables
BINARY_DIR=bin
READER_BINARY=$(BINARY_DIR)/reader
WRITER_BINARY=$(BINARY_DIR)/writer

# Go build flags
BUILD_FLAGS=
ifeq ($(OS),Windows_NT)
	READER_BINARY := $(BINARY_DIR)/reader.exe
	WRITER_BINARY := $(BINARY_DIR)/writer.exe
endif

# Default target
.PHONY: all
all: reader writer

# Build reader
.PHONY: reader
reader:
	@echo "Building reader..."
	@CGO_ENABLED=1 go build -o $(READER_BINARY) ./cmd/reader

# Build writer
.PHONY: writer
writer:
	@echo "Building writer..."
	@CGO_ENABLED=1 go build -o $(WRITER_BINARY) ./cmd/writer

# Run reader
.PHONY: run-reader
run-reader:
	@echo "Starting reader on port 8080..."
	@CGO_ENABLED=1 go run ./cmd/reader/main.go

# Run writer
.PHONY: run-writer
run-writer:
	@echo "Starting writer on port 8081..."
	@CGO_ENABLED=1 PORT=8081 go run ./cmd/writer/main.go

# Clean binaries
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BINARY_DIR)

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make reader      - Build reader binary"
	@echo "  make writer      - Build writer binary"
	@echo "  make all         - Build both binaries"
	@echo "  make run-reader  - Run reader service"
	@echo "  make run-writer  - Run writer service"
	@echo "  make clean       - Remove binaries"
	@echo "  make deps        - Install dependencies"
