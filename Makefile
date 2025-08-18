# Go parameters
BINARY_NAME=gitx
CMD_PATH=./cmd/gitx
BUILD_DIR=./build

# Default target executed when you run `make`
all: build

# Syncs dependencies
sync:
	@echo "Syncing dependencies..."
	@go mod tidy
	@echo "Dependencies synced."

# Builds the binary
build: sync
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Binary available at $(BUILD_DIR)/$(BINARY_NAME)"

# Runs the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Runs all tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Runs golangci-lint
ci:
	@echo "Running golangci-lint..."
	@golangci-lint run

# Installs the binary to /usr/local/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@sudo install $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin
	@echo "$(BINARY_NAME) installed successfully to /usr/local/bin"

# Cleans the build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Cleanup complete."

#PHONY targets are not files
.PHONY: all sync build run test ci install clean
