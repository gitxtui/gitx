# Go parameters
BINARY_NAME=gitx
CMD_PATH=./cmd/gitx
BUILD_DIR=./build

# Default target executed when you run `make`
all: build

# Builds the binary
build:
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Binary available at $(BUILD_DIR)/$(BINARY_NAME)"

# Runs all tests
test:
	@echo "Running tests..."
	@go test -v ./...

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
.PHONY: all build test install clean
