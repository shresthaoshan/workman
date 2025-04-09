# Variables
APP_NAME := workman
GO := go
GO_BUILD_FLAGS := -v
GO_TEST_FLAGS := -v
BINARY_LINUX := bin/$(APP_NAME)-linux-arm64
BINARY_MACOS := bin/$(APP_NAME)-darwin-arm64
BINARY_WINDOWS := bin/$(APP_NAME)-windows-arm64.exe
BINARY_LINUX_AMD := bin/$(APP_NAME)-linux-amd64
BINARY_MACOS_AMD := bin/$(APP_NAME)-darwin-amd64
BINARY_WINDOWS_AMD := bin/$(APP_NAME)-windows-amd64.exe
GO_FILES := $(shell find . -type f -name '*.go')

# Default target
all: build

# Build the binary for the current platform
build: $(APP_NAME)

$(APP_NAME):
	$(GO) build $(GO_BUILD_FLAGS) -o $(APP_NAME) cmd/workman/main.go

# Cross-platform builds
build-all: build-linux build-macos build-windows build-linux-amd build-macos-amd build-windows-amd

build-linux:
	GOOS=linux GOARCH=arm64 $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_LINUX) cmd/workman/main.go

build-macos:
	GOOS=darwin GOARCH=arm64 $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_MACOS) cmd/workman/main.go

build-windows:
	GOOS=windows GOARCH=arm64 $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_WINDOWS) cmd/workman/main.go

build-linux-amd:
	GOOS=linux GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_LINUX_AMD) cmd/workman/main.go

build-macos-amd:
	GOOS=darwin GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_MACOS_AMD) cmd/workman/main.go

build-windows-amd:
	GOOS=windows GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o $(BINARY_WINDOWS_AMD) cmd/workman/main.go

# Install dependencies
deps:
	$(GO) mod download

# Run tests
test:
	$(GO) test $(GO_TEST_FLAGS) ./...

# Clean up build artifacts
clean:
	rm -f $(APP_NAME) $(BINARY_LINUX) $(BINARY_MACOS) $(BINARY_WINDOWS) $(BINARY_LINUX_AMD) $(BINARY_MACOS_AMD) $(BINARY_WINDOWS_AMD)

# Format the code
fmt:
	$(GO) fmt ./...

# Run the application
run: build
	./$(APP_NAME) --help

# Phony targets (targets that don't produce files)
.PHONY: all build build-all build-linux build-macos build-windows deps test clean lint fmt run docs