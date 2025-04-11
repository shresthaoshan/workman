# Variables
APP_NAME := workman
VERSION := $(shell git describe --tags --always)

GO := go
GO_TEST_FLAGS := -v
GO_IMAGE := golang:1.24

BUILD_DIR := bin
MAIN_GO_FILE := cmd/$(APP_NAME)/main.go

PLATFORMS := linux/amd64 linux/arm64 windows/amd64 windows/arm64 darwin/amd64 darwin/arm64

# Default target
all: build

build:
	@make clean
	@echo ">> Building in Docker using $(GO_IMAGE)..."
	@mkdir -p $(BUILD_DIR)
	@$(foreach platform,$(PLATFORMS), \
		GOOS=$(word 1,$(subst /, ,$(platform))); \
		GOARCH=$(word 2,$(subst /, ,$(platform))); \
		EXT=$$(if [ "$$GOOS" = "windows" ]; then echo ".exe"; else echo ""; fi); \
		OUT_FILE=$(APP_NAME)-$(VERSION)-$$GOOS-$$GOARCH$$EXT; \
		echo "Building $$OUT_FILE..."; \
		docker run --rm \
			-e GOOS=$$GOOS \
			-e GOARCH=$$GOARCH \
			-e CGO_ENABLED=0 \
			-v $$PWD:/app -w /app \
			$(GO_IMAGE) \
			go build -o $(BUILD_DIR)/$$OUT_FILE /app/$(MAIN_GO_FILE); \
	)

# Install dependencies
deps:
	$(GO) mod download

# Run tests
test:
	$(GO) test $(GO_TEST_FLAGS) ./...

# Clean up build artifacts
clean:
	rm -f $(APP_NAME)
	rm -fr $(BUILD_DIR)

# Format the code
fmt:
	$(GO) fmt ./...

# Phony targets (targets that don't produce files)
.PHONY: all build deps test clean fmt