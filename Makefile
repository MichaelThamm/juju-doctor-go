APP_NAME := juju-doctor
SRC_DIR := ./cmd/jujudoctor/main.go
BUILD_DIR := ./bin
GO := go

.PHONY: build

# Build the Go application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)